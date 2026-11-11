package system

import (
	"context"
	"time"
)

// Watchdog restarts a component if it stops heartbeating.
type Watchdog struct {
	timeout  time.Duration
	lastBeat time.Time
	onTimeout func()
}

// NewWatchdog creates a watchdog with the given timeout and callback.
func NewWatchdog(timeout time.Duration, onTimeout func()) *Watchdog {
	return &Watchdog{timeout: timeout, lastBeat: time.Now(), onTimeout: onTimeout}
}

// Heartbeat records a liveness signal.
func (w *Watchdog) Heartbeat() { w.lastBeat = time.Now() }

// Run monitors liveness and invokes onTimeout when stale.
func (w *Watchdog) Run(ctx context.Context) {
	ticker := time.NewTicker(w.timeout / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if time.Since(w.lastBeat) > w.timeout {
				if w.onTimeout != nil {
					w.onTimeout()
				}
				w.lastBeat = time.Now()
			}
		}
	}
}
