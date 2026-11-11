package apps

import (
	"encoding/json"
	"os"
)

// Registry holds the set of known apps.
type Registry struct {
	Apps []App `json:"apps"`
}

// LoadRegistry reads an app registry from a JSON file.
func LoadRegistry(path string) (Registry, error) {
	var r Registry
	f, err := os.Open(path)
	if err != nil {
		return r, err
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&r); err != nil {
		return r, err
	}
	return r, nil
}

// Find returns an app by id.
func (r Registry) Find(id string) (App, bool) {
	for _, a := range r.Apps {
		if a.ID == id {
			return a, true
		}
	}
	return App{}, false
}
