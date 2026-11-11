//go:build linux

package virtual

import "github.com/ngbcoder/fireside/internal/input"

// SendKey is a no-op on Linux (use uinput forwarder instead).
func SendKey(k input.Key, down bool) {}

// TypeString is a no-op on Linux.
func TypeString(s string) {}