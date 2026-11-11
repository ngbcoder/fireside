//go:build windows

package system

import (
	"syscall"
	"time"
	"unsafe"
)

var (
	kernel32 = syscall.NewLazyDLL("kernel32.dll")
	user32k  = syscall.NewLazyDLL("user32.dll")
	procFindWindowW      = user32k.NewProc("FindWindowW")
	procSetWindowLongW   = user32k.NewProc("SetWindowLongW")
	procSetWindowPos     = user32k.NewProc("SetWindowPos")
	procShowWindow       = user32k.NewProc("ShowWindow")
	procGetSystemMetrics = user32k.NewProc("GetSystemMetrics")
)

const (
	gwlStyle       = ^uintptr(15) // -16
	wsOverlapped   = 0x00000000
	wsCaption      = 0x00C00000
	wsThickFrame   = 0x00040000
	wsMaximizeBox  = 0x00010000
	wsMinimizeBox  = 0x00020000
	wsSysMenu      = 0x00080000
	swShowMaximized = 3
	smCxScreen     = 0
	smCyScreen     = 1
	hwndTopMost    = ^uintptr(0) // -1
	swpFrameChanged = 0x0020
	swpShowWindow   = 0x0040
)

// EnterKioskMode makes the Fireside window fullscreen and borderless on Windows.
func EnterKioskMode(title string) {
	go func() {
		time.Sleep(1 * time.Second)
		titlePtr, _ := syscall.UTF16PtrFromString(title)
		hwnd, _, _ := procFindWindowW.Call(0, uintptr(unsafe.Pointer(titlePtr)))
		if hwnd == 0 {
			return
		}
		style, _, _ := procSetWindowLongW.Call(hwnd, gwlStyle, wsOverlapped)
		_ = style
		// remove caption/thickframe by clearing those bits
		newStyle := uint32(wsOverlapped)
		procSetWindowLongW.Call(hwnd, gwlStyle, uintptr(newStyle))
		w, _, _ := procGetSystemMetrics.Call(smCxScreen)
		h, _, _ := procGetSystemMetrics.Call(smCyScreen)
		procSetWindowPos.Call(hwnd, hwndTopMost, 0, 0, w, h, swpFrameChanged|swpShowWindow)
		procShowWindow.Call(hwnd, swShowMaximized)
	}()
}
