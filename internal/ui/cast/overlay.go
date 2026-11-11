package cast

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Overlay is the casting takeover screen.
type Overlay struct {
	Theme *material.Theme
	Peer  string
	Title string
}

// Layout renders the casting overlay.
func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ColorBackground)
	l := material.H1(o.Theme, "Casting")
	l.Color = theme.ColorPrimary
	dims := l.Layout(gtx)
	sub := material.Body1(o.Theme, "From "+o.Peer)
	sub.Color = theme.ColorSecondary
	sub.Layout(gtx)
	return dims
}
