package dial

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/ngbcoder/fireside/internal/cast/session"
)

// Server is a DIAL receiver.
type Server struct {
	DeviceName string
	Port       int // HTTP port
	events     chan session.SessionEvent
	httpSrv    *http.Server
	listener   net.Listener
	ssdpConn   *net.UDPConn
	mu         sync.Mutex
}

// New creates a DIAL server.
func New(deviceName string, port int) *Server {
	if port == 0 {
		port = 8080
	}
	if deviceName == "" {
		deviceName = "Fireside"
	}
	return &Server{DeviceName: deviceName, Port: port, events: make(chan session.SessionEvent, 16)}
}

func (s *Server) Name() string { return "dial" }

func (s *Server) Start(ctx context.Context, deviceName string) error {
	if deviceName != "" {
		s.DeviceName = deviceName
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/ssdp/device-desc.xml", s.handleDeviceDesc)
	mux.HandleFunc("/apps/", s.handleApp)
	s.httpSrv = &http.Server{Addr: fmt.Sprintf(":%d", s.Port), Handler: mux}
	ln, err := net.Listen("tcp", s.httpSrv.Addr)
	if err != nil {
		return err
	}
	s.listener = ln
	s.Port = ln.Addr().(*net.TCPAddr).Port
	go func() { _ = s.httpSrv.Serve(ln) }()

	// SSDP multicast responder.
	go s.ssdpLoop(ctx)

	s.events <- session.SessionEvent{Kind: session.SessionStarted, Protocol: "dial", Peer: s.DeviceName}
	return nil
}

func (s *Server) Stop() error {
	if s.httpSrv != nil {
		_ = s.httpSrv.Close()
	}
	if s.ssdpConn != nil {
		_ = s.ssdpConn.Close()
	}
	return nil
}

func (s *Server) Events() <-chan session.SessionEvent { return s.events }

func (s *Server) handleDeviceDesc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	fmt.Fprintf(w, `<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
  <specVersion><major>1</major><minor>0</minor></specVersion>
  <device>
    <deviceType>urn:dial-multiscreen-org:device:dial:1</deviceType>
    <friendlyName>%s</friendlyName>
    <manufacturer>Fireside</manufacturer>
    <modelName>Fireside</modelName>
    <UDN>uuid:fireside-dial</UDN>
    <serviceList></serviceList>
  </device>
</root>`, s.DeviceName)
}

func (s *Server) handleApp(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/apps/")
	name = strings.TrimSuffix(name, "/")
	if name == "" {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		// Launch app: emit session event.
		s.events <- session.SessionEvent{Kind: session.SessionStarted, Protocol: "dial", Title: name}
		w.WriteHeader(http.StatusCreated)
		return
	}
	// GET: return running state.
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	fmt.Fprintf(w, `<service xmlns="urn:dial-multiscreen-org:schemas:dial"><state>stopped</state></service>`)
}

func (s *Server) ssdpLoop(ctx context.Context) {
	addr, err := net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if err != nil {
		return
	}
	conn, err := net.ListenMulticastUDP("udp4", nil, addr)
	if err != nil {
		return
	}
	s.ssdpConn = conn
	defer conn.Close()

	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}
		n, src, err := conn.ReadFromUDP(buf)
		if err != nil {
			return
		}
		msg := string(buf[:n])
		if strings.Contains(msg, "M-SEARCH") && strings.Contains(msg, "dial") || strings.Contains(msg, "ssdp:all") || strings.Contains(msg, "upnp:rootdevice") {
			resp := fmt.Sprintf("HTTP/1.1 200 OK\r\n"+
				"CACHE-CONTROL: max-age=1800\r\n"+
				"EXT:\r\n"+
				"LOCATION: http://%s:%d/ssdp/device-desc.xml\r\n"+
				"SERVER: Fireside/1.0 UPnP/1.0\r\n"+
				"ST: urn:dial-multiscreen-org:service:dial:1\r\n"+
				"USN: uuid:fireside-dial::urn:dial-multiscreen-org:service:dial:1\r\n\r\n", localIP(), s.Port)
			_ = sendUDP(src, resp)
		}
		_ = bufio.NewScanner
	}
}

func localIP() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}

func sendUDP(dst *net.UDPAddr, msg string) error {
	conn, err := net.DialUDP("udp4", nil, dst)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Write([]byte(msg))
	return err
}
