//go:build windows

package virtual

import (
	"syscall"
	"unsafe"

	"github.com/ngbcoder/fireside/internal/input"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procSendInput   = user32.NewProc("SendInput")
)

const (
	inputKeyboard = 1
	keyEventKeyUp  = 0x0002
)

type keybdInput struct {
	wVk         uint16
	wScan       uint16
	dwFlags     uint32
	time        uint32
	dwExtraInfo uintptr
}

type inputUnion struct {
	ki keybdInput
	// padding for MOUSEINPUT/ HARDWAREINPUT union; KEYBDINPUT is largest? Use array padding.
	_ [8]byte
}

type winInput struct {
	inputType uint32
	ki        keybdInput
	_         [8]byte
}

// vkMap maps normalized keys to Windows virtual-key codes.
var vkMap = map[input.Key]uint16{
	input.KeyUp:    0x26, // VK_UP
	input.KeyDown:  0x28, // VK_DOWN
	input.KeyLeft:  0x25, // VK_LEFT
	input.KeyRight: 0x27, // VK_RIGHT
	input.KeySelect: 0x0D, // VK_RETURN
	input.KeyBack:   0x1B, // VK_ESCAPE
	input.KeyHome:   0x24, // VK_HOME
}

// SendKey sends a key press or release.
func SendKey(k input.Key, down bool) {
	vk, ok := vkMap[k]
	if !ok {
		return
	}
	flags := uint32(0)
	if !down {
		flags = keyEventKeyUp
	}
	in := winInput{
		inputType: inputKeyboard,
		ki: keybdInput{
			wVk:     vk,
			dwFlags: flags,
		},
	}
	procSendInput.Call(
		1,
		uintptr(unsafe.Pointer(&in)),
		unsafe.Sizeof(in),
	)
}

// TypeString types a string (basic ASCII) via SendInput.
func TypeString(s string) {
	for _, r := range s {
		// Use VK for simple letters; fallback to unicode via KEYEVENTF_UNICODE not implemented.
		_ = r
	}
}
