package session

import "testing"

func TestFakeProtocol(t *testing.T) {
	f := &FakeProtocol{NameVal: "fake"}
	if f.Name() != "fake" {
		t.Fatal("bad name")
	}
	f.Send(SessionEvent{Kind: SessionStarted, Protocol: "fake"})
	e := <-f.Events()
	if e.Kind != SessionStarted {
		t.Fatal("bad event")
	}
}
