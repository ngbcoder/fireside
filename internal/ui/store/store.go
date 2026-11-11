package store

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Screen is the App Store screen.
type Screen struct {
	Theme *material.Theme
}

// New creates a store screen.
func New(th *material.Theme) *Screen { return &Screen{Theme: th} }

// Layout renders the store.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ColorBackground)
	l := material.H2(s.Theme, "App Store")
	l.Color = theme.ColorPrimary
	return l.Layout(gtx)
}
