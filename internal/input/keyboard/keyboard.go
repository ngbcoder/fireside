package keyboard

import (
	"github.com/ngbcoder/fireside/internal/input"
)

// Source is a keyboard input source. Events are pushed via Send.
type Source struct {
	name   string
	events chan input.InputEvent
}

// New creates a keyboard source.
func New(name string) *Source {
	if name == "" {
		name = "keyboard"
	}
	return &Source{name: name, events: make(chan input.InputEvent, 64)}
}

func (s *Source) Name() string { return s.name }

func (s *Source) Events() <-chan input.InputEvent { return s.events }

func (s *Source) Send(e input.InputEvent) { s.events <- e }

func (s *Source) Close() error {
	close(s.events)
	return nil
}
