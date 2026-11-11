package onboarding

import (
	"image"

	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget/material"

	"github.com/skip2/go-qrcode"
)

// Overlay shows a QR code for onboarding after network connect.
type Overlay struct {
	Theme *material.Theme
	URL   string
	img   image.Image
}

// Ensure generates the QR image for URL.
func (o *Overlay) Ensure() {
	if o.img != nil || o.URL == "" {
		return
	}
	qr, err := qrcode.New(o.URL, qrcode.Medium)
	if err != nil {
		return
	}
	o.img = qr.Image(256)
}

// Layout renders the onboarding overlay.
func (o *Overlay) Layout(gtx layout.Context) layout.Dimensions {
	o.Ensure()
	if o.img == nil {
		l := material.Body1(o.Theme, "Onboarding unavailable")
		return l.Layout(gtx)
	}
	op := paint.NewImageOp(o.img)
	op.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	l := material.Body1(o.Theme, "Scan to set up: "+o.URL)
	return l.Layout(gtx)
}
