package airplay

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/ngbcoder/fireside/internal/cast/session"
)

// Config for the AirPlay supervisor.
type Config struct {
	Bin        string // uxplay binary, default "uxplay"
	DeviceName string
	Args       []string
}

// AirPlay wraps UxPlay as a managed subprocess and implements session.Protocol.
type AirPlay struct {
	cfg     Config
	events  chan session.SessionEvent
	cancel  context.CancelFunc
	cmd     *exec.Cmd
}

// New creates an AirPlay supervisor.
func New(cfg Config) *AirPlay {
	if cfg.Bin == "" {
		cfg.Bin = "uxplay"
	}
	if cfg.DeviceName == "" {
		cfg.DeviceName = "Fireside"
	}
	return &AirPlay{cfg: cfg, events: make(chan session.SessionEvent, 16)}
}

func (a *AirPlay) Name() string { return "airplay" }

func (a *AirPlay) Start(ctx context.Context, deviceName string) error {
	if deviceName != "" {
		a.cfg.DeviceName = deviceName
	}
	runCtx, cancel := context.WithCancel(ctx)
	a.cancel = cancel

	args := append([]string{}, a.cfg.Args...)
	args = append(args, "-n", a.cfg.DeviceName)
	a.cmd = exec.CommandContext(runCtx, a.cfg.Bin, args...)

	stdout, _ := a.cmd.StdoutPipe()
	stderr, _ := a.cmd.StderrPipe()

	if err := a.cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start uxplay: %w", err)
	}

	a.events <- session.SessionEvent{Kind: session.SessionStarted, Protocol: "airplay", Peer: a.cfg.DeviceName}

	// Parse UxPlay output for session lifecycle hints.
	go func() {
		sc := bufio.NewScanner(stdout)
		for sc.Scan() {
			line := sc.Text()
			lower := strings.ToLower(line)
			if strings.Contains(lower, "client") && (strings.Contains(lower, "connected") || strings.Contains(lower, "playing")) {
				a.events <- session.SessionEvent{Kind: session.SessionMedia, Protocol: "airplay", Title: strings.TrimSpace(line)}
			}
		}
	}()
	go func() { bufio.NewScanner(stderr).Scan() }()

	go func() {
		_ = a.cmd.Wait()
		a.events <- session.SessionEvent{Kind: session.SessionEnded, Protocol: "airplay"}
	}()
	return nil
}

func (a *AirPlay) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	if a.cmd != nil && a.cmd.Process != nil {
		_ = a.cmd.Process.Kill()
	}
	return nil
}

func (a *AirPlay) Events() <-chan session.SessionEvent { return a.events }
