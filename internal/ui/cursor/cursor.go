package cursor

import (
	"image"
	"image/color"

	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// Overlay draws a software cursor at (X,Y).
type Overlay struct {
	X, Y int
}

// Layout renders the cursor.
func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	if o == nil {
		return layout.Dimensions{}
	}
	defer clip.UniformRRect(image.Rect(o.X, o.Y, o.X+20, o.Y+20), 10).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xDD})
	return layout.Dimensions{Size: image.Pt(20, 20)}
}
