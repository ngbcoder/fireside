package spotify

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/ngbcoder/fireside/internal/cast/session"
)

// Receiver wraps librespot as a Spotify Connect receiver.
type Receiver struct {
	Bin        string
	DeviceName string
	events     chan session.SessionEvent
	cancel     context.CancelFunc
	cmd        *exec.Cmd
}

// New creates a Spotify Connect receiver.
func New(deviceName string) *Receiver {
	if deviceName == "" {
		deviceName = "Fireside"
	}
	return &Receiver{Bin: "librespot", DeviceName: deviceName, events: make(chan session.SessionEvent, 16)}
}

func (r *Receiver) Name() string { return "spotify" }

func (r *Receiver) Start(ctx context.Context, deviceName string) error {
	if deviceName != "" {
		r.DeviceName = deviceName
	}
	runCtx, cancel := context.WithCancel(ctx)
	r.cancel = cancel
	r.cmd = exec.CommandContext(runCtx, r.Bin, "--name", r.DeviceName, "--backend", "pipe")
	if err := r.cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("start librespot: %w", err)
	}
	r.events <- session.SessionEvent{Kind: session.SessionStarted, Protocol: "spotify", Peer: r.DeviceName}
	go func() {
		_ = r.cmd.Wait()
		r.events <- session.SessionEvent{Kind: session.SessionEnded, Protocol: "spotify"}
	}()
	return nil
}

func (r *Receiver) Stop() error {
	if r.cancel != nil {
		r.cancel()
	}
	if r.cmd != nil && r.cmd.Process != nil {
		_ = r.cmd.Process.Kill()
	}
	return nil
}

func (r *Receiver) Events() <-chan session.SessionEvent { return r.events }
