package settings

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Screen is the Settings screen.
type Screen struct {
	Theme *material.Theme
	Items []string
}

// New creates a settings screen.
func New(th *material.Theme) *Screen {
	return &Screen{Theme: th, Items: []string{
		"Network", "Bluetooth", "Audio output", "Display", "Device name", "About",
	}}
}

// Layout renders settings.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ColorBackground)
	l := material.H2(s.Theme, "Settings")
	l.Color = theme.ColorPrimary
	return l.Layout(gtx)
}
