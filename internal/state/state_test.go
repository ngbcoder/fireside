package state

import "testing"

func TestStateMachine(t *testing.T) {
	s := New()
	if !s.IsShell() {
		t.Fatal("expected shell")
	}
	s.ToApp()
	if s.Mode() != ModeApp {
		t.Fatal("expected app")
	}
	s.ToShell()
	if !s.IsShell() {
		t.Fatal("expected shell")
	}
	s.ToCast()
	if s.Mode() != ModeCast {
		t.Fatal("expected cast")
	}
}
