package search

import (
	"gioui.org/io/key"
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
	Selected int
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
	q := s.Editor.Text()
	if s.Searcher != nil {
		s.Results = s.Searcher.Search(q)
	}
	if s.Selected >= len(s.Results) {
		s.Selected = 0
	}

	// Key handling for navigation/launch.
	for {
		ev, ok := gtx.Event(
			key.Filter{Name: key.NameDownArrow},
			key.Filter{Name: key.NameUpArrow},
			key.Filter{Name: key.NameReturn},
		)
		if !ok {
			break
		}
		if ke, ok := ev.(key.Event); ok && ke.State == key.Press {
			switch ke.Name {
			case key.NameDownArrow:
				if s.Selected < len(s.Results)-1 {
					s.Selected++
				}
			case key.NameUpArrow:
				if s.Selected > 0 {
					s.Selected--
				}
			case key.NameReturn:
				if len(s.Results) > 0 && s.OnLaunch != nil {
					s.OnLaunch(s.Results[s.Selected].App)
				}
			}
		}
	}

	paint.Fill(gtx.Ops, theme.ColorBackground)
	// Editor
	ed := material.Editor(s.Theme, &s.Editor, "Search apps...")
	edDims := ed.Layout(gtx)

	// Results or empty state.
	if len(s.Results) == 0 {
		l := material.Body1(s.Theme, "No results")
		l.Color = theme.ColorSecondary
		l.Layout(gtx)
		return edDims
	}
	for i, r := range s.Results {
		l := material.Body1(s.Theme, r.App.Name)
		if i == s.Selected {
			l.Color = theme.ColorAccent
		} else {
			l.Color = theme.ColorPrimary
		}
		l.Layout(gtx)
	}
	return edDims
}
