package state

// Mode represents who owns the screen.
type Mode int

const (
	ModeShell Mode = iota
	ModeApp
	ModeCast
)

// State is the shell state machine.
type State struct {
	mode Mode
}

// New creates a state in Shell mode.
func New() *State { return &State{mode: ModeShell} }

// Mode returns the current mode.
func (s *State) Mode() Mode { return s.mode }

// ToApp transitions to App mode.
func (s *State) ToApp() { s.mode = ModeApp }

// ToCast transitions to Cast mode.
func (s *State) ToCast() { s.mode = ModeCast }

// ToShell returns to Shell mode.
func (s *State) ToShell() { s.mode = ModeShell }

// IsShell reports whether the shell owns the screen.
func (s *State) IsShell() bool { return s.mode == ModeShell }
