package session

import "context"

// SessionKind describes a cast session lifecycle event.
type SessionKind int

const (
	SessionStarted SessionKind = iota
	SessionMedia
	SessionPaused
	SessionResumed
	SessionEnded
)

// SessionEvent is emitted by a protocol plugin.
type SessionEvent struct {
	Kind     SessionKind
	Protocol string // "airplay", "dial", "miracast"
	Peer     string // sender device name, if known
	Title    string // optional
}

// Protocol is a casting protocol plugin.
type Protocol interface {
	Name() string
	Start(ctx context.Context, deviceName string) error
	Stop() error
	Events() <-chan SessionEvent
}
