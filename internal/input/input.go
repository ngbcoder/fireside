package input

// Key is a normalized remote key.
type Key int

const (
	KeyUp Key = iota
	KeyDown
	KeyLeft
	KeyRight
	KeySelect
	KeyBack
	KeyHome
	KeyPlayPause
	KeyMenu
	KeyVolumeUp
	KeyVolumeDown
)

// Action is the key action.
type Action int

const (
	Press Action = iota
	Release
	Repeat
)

// InputEvent is a normalized input event.
type InputEvent struct {
	Key    Key
	Action Action
	Source string // "cec", "evdev:/dev/input/event3", "keyboard"
}

// InputSource is an input device.
type InputSource interface {
	Name() string
	Events() <-chan InputEvent
	Close() error
}
