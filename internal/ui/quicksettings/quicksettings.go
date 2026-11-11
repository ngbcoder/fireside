package quicksettings

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Overlay is the quick-settings panel.
type Overlay struct {
	Theme      *material.Theme
	Visible    bool
	Volume     int
	Brightness int
}

// Layout renders the quick settings overlay when visible.
func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	if !o.Visible {
		return layout.Dimensions{}
	}
	panel := image.Rect(0, 0, gtx.Constraints.Max.X, 200)
	defer clip.UniformRRect(panel, theme.RadiusLG).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, theme.ColorGlass)

	l := material.H3(o.Theme, "Quick Settings")
	l.Color = theme.ColorPrimary
	return l.Layout(gtx)
}
