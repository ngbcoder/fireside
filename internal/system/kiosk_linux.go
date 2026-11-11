//go:build linux

package system

// EnterKioskMode is a no-op on Linux (use cage).
func EnterKioskMode(title string) {}