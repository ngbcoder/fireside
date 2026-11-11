package settings

import (
	"gioui.org/layout"
	"gioui.org/op/paint"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/profiles"
	"github.com/ngbcoder/fireside/internal/ui/theme"
)

// Screen is the Settings screen, also used for manual onboarding setup.
type Screen struct {
	Theme       *material.Theme
	DeviceName  widget.Editor
	SaveBtn     widget.Clickable
	ScanBtn     widget.Clickable
	Networks    []string
	Status      string
	Profiles    *profiles.Manager
}

// New creates a settings screen.
func New(th *material.Theme) *Screen {
	s := &Screen{Theme: th}
	s.DeviceName.SingleLine = true
	return s
}

// Layout renders settings.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	paint.Fill(gtx.Ops, theme.ColorBackground)

	// Handle save.
	for s.SaveBtn.Clicked(gtx) {
		s.Status = "Saved: " + s.DeviceName.Text()
	}
	// Handle scan.
	for s.ScanBtn.Clicked(gtx) {
		s.Status = "Scanning Wi-Fi..."
		// TODO: call system.Network.AvailableNetworks()
		s.Networks = []string{"(no networks found)"}
	}

	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceStart}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			l := material.H2(s.Theme, "Settings / Setup")
			l.Color = theme.ColorPrimary
			return l.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			l := material.Body1(s.Theme, "Device name")
			l.Color = theme.ColorSecondary
			return l.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			ed := material.Editor(s.Theme, &s.DeviceName, "Fireside")
			return ed.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			btn := material.Button(s.Theme, &s.SaveBtn, "Save")
			return btn.Layout(gtx)
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if s.Profiles == nil {
				return layout.Dimensions{}
			}
			l := material.H3(s.Theme, "Profiles")
			l.Color = theme.ColorPrimary
			dims := l.Layout(gtx)
			active, _ := s.Profiles.Active()
			l2 := material.Body2(s.Theme, "Active: "+active.Name)
			l2.Color = theme.ColorSecondary
			l2.Layout(gtx)
			return dims
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			for _, n := range s.Networks {
				l := material.Body2(s.Theme, n)
				l.Color = theme.ColorSecondary
				l.Layout(gtx)
			}
			if s.Status != "" {
				l := material.Body2(s.Theme, s.Status)
				l.Color = theme.ColorAccent
				return l.Layout(gtx)
			}
			return layout.Dimensions{}
		}),
	)
}
