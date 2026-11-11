package store

import "context"

// AppEntry is a store listing.
type AppEntry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	URL         string `json:"url"` // web app URL
	Category    string `json:"category"`
}

// State represents install state.
type State int

const (
	StateAvailable State = iota
	StateInstalling
	StateInstalled
	StateUpdateAvailable
)

// Backend is a store backend.
type Backend interface {
	List(ctx context.Context) ([]AppEntry, error)
	Install(ctx context.Context, id string) error
	Remove(ctx context.Context, id string) error
	State(id string) State
}

// RegistryBackend is a curated web app registry.
type RegistryBackend struct {
	Entries []AppEntry
	States  map[string]State
}

func NewRegistryBackend(entries []AppEntry) *RegistryBackend {
	return &RegistryBackend{Entries: entries, States: map[string]State{}}
}

func (r *RegistryBackend) List(ctx context.Context) ([]AppEntry, error) {
	return r.Entries, nil
}

func (r *RegistryBackend) Install(ctx context.Context, id string) error {
	r.States[id] = StateInstalled
	return nil
}

func (r *RegistryBackend) Remove(ctx context.Context, id string) error {
	r.States[id] = StateAvailable
	return nil
}

func (r *RegistryBackend) State(id string) State {
	if s, ok := r.States[id]; ok {
		return s
	}
	return StateAvailable
}
