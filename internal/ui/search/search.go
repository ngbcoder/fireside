package search

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/apps"
	"github.com/ngbcoder/fireside/internal/search"
	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Screen is the Search screen.
type Screen struct {
	Theme    *material.Theme
	Editor   widget.Editor
	Searcher *search.Searcher
	Results  []search.Result
	OnLaunch func(apps.App)
}

// New creates a search screen.
func New(th *material.Theme, searcher *search.Searcher) *Screen {
	s := &Screen{Theme: th, Searcher: searcher}
	s.Editor.SingleLine = true
	s.Editor.Submit = true
	return s
}

// Layout renders search.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	// Update results on query change.
	q := s.Editor.Text()
	s.Results = nil
	if s.Searcher != nil {
		s.Results = s.Searcher.Search(q)
	}

	paint.Fill(gtx.Ops, theme.ColorBackground)
	// Editor
	editor := material.Editor(s.Theme, &s.Editor, "Search apps...")
	edDims := editor.Layout(gtx)

	// Results list (simple vertical)
	for _, r := range s.Results {
		l := material.Body1(s.Theme, r.App.Name)
		l.Color = theme.ColorPrimary
		l.Layout(gtx)
	}
	return edDims
}
