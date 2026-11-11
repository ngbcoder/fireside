package miracast

import (
	"context"

	"github.com/ngbcoder/fireside/internal/cast/session"
)

// Receiver is an experimental Miracast receiver (stretch goal).
type Receiver struct {
	events chan session.SessionEvent
}

// New creates a Miracast receiver stub.
func New() *Receiver {
	return &Receiver{events: make(chan session.SessionEvent, 16)}
}

func (r *Receiver) Name() string { return "miracast" }

func (r *Receiver) Start(ctx context.Context, deviceName string) error {
	// TODO: Wi-Fi Direct P2P, RTSP negotiation, RTP/H.264 receive.
	return nil
}

func (r *Receiver) Stop() error { return nil }

func (r *Receiver) Events() <-chan session.SessionEvent { return r.events }
