package session

import "context"

// FakeProtocol is a test double for Protocol.
type FakeProtocol struct {
	NameVal  string
	EventsCh chan SessionEvent
}

func (f *FakeProtocol) ensure() {
	if f.EventsCh == nil {
		f.EventsCh = make(chan SessionEvent, 16)
	}
}

func (f *FakeProtocol) Name() string { return f.NameVal }

func (f *FakeProtocol) Start(ctx context.Context, deviceName string) error { return nil }

func (f *FakeProtocol) Stop() error { return nil }

func (f *FakeProtocol) Events() <-chan SessionEvent {
	f.ensure()
	return f.EventsCh
}

// Send emits an event on the fake protocol.
func (f *FakeProtocol) Send(e SessionEvent) {
	f.ensure()
	f.EventsCh <- e
}
