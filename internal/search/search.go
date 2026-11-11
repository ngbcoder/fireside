package search

import (
	"strings"

	"github.com/ngbcoder/fireside/internal/apps"
)

// Result is a search hit.
type Result struct {
	App apps.App
}

// Searcher searches apps by name/id.
type Searcher struct {
	Apps []apps.App
}

// New creates a searcher.
func New(apps []apps.App) *Searcher { return &Searcher{Apps: apps} }

// Search returns apps whose name or id contains query.
func (s *Searcher) Search(query string) []Result {
	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return nil
	}
	var out []Result
	for _, a := range s.Apps {
		if strings.Contains(strings.ToLower(a.Name), q) || strings.Contains(strings.ToLower(a.ID), q) {
			out = append(out, Result{App: a})
		}
	}
	return out
}
