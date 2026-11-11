package system

import (
	"context"
	"net"
	"time"
)

// Connected reports whether a non-loopback interface has an IP address.
func Connected() bool {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, _ := iface.Addrs()
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				return true
			}
		}
	}
	return false
}

// WatchNetwork polls for connectivity and calls cb when it transitions to connected.
func WatchNetwork(ctx context.Context, interval time.Duration, cb func()) {
	if interval <= 0 {
		interval = 2 * time.Second
	}
	wasConnected := Connected()
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := Connected()
			if now && !wasConnected {
				if cb != nil {
					cb()
				}
			}
			wasConnected = now
		}
	}
}
