package search

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Screen is the Search screen.
type Screen struct {
	Theme *material.Theme
	Query string
}

// New creates a search screen.
func New(th *material.Theme) *Screen { return &Screen{Theme: th} }

// Layout renders search.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ColorBackground)
	l := material.H2(s.Theme, "Search")
	l.Color = theme.ColorPrimary
	return l.Layout(gtx)
}
