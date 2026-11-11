package home

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/apps"
	"github.com/ngbcoder/fireside/internal/ui/focus"
	"github.com/ngbcoder/fireside/internal/ui/theme"
)

var homeTag = new(int)

// Screen is the Home screen.
type Screen struct {
	Theme *material.Theme
	Focus *focus.Manager

	Tabs []string
	Apps []AppItem

	Registry apps.Registry
	Folders  []apps.Folder

	tabFocusIDs []string
	appFocusIDs []string
}

// AppItem represents an app on the home row.
type AppItem struct {
	ID   string
	Name string
}

// New creates a Home screen with default content.
func New(th *material.Theme) *Screen {
	s := &Screen{Theme: th, Focus: focus.NewManager(), Tabs: []string{"Home", "Library", "App Store", "Search", "Settings"}}
	s.Apps = []AppItem{
		{ID: "youtube", Name: "YouTube"},
		{ID: "disney", Name: "Disney+"},
		{ID: "appletv", Name: "Apple TV"},
		{ID: "media", Name: "Media Player"},
		{ID: "music", Name: "Music"},
		{ID: "settings", Name: "Settings"},
	}
	s.rebuildFocusNodes()
	return s
}

func (s *Screen) rebuildFocusNodes() {
	var nodes []focus.Node
	s.tabFocusIDs = make([]string, len(s.Tabs))
	for i := range s.Tabs {
		id := "tab:" + s.Tabs[i]
		s.tabFocusIDs[i] = id
		nodes = append(nodes, focus.Node{ID: id, Rect: focus.Rect{X: theme.OverscanMargin(1920) + i*180, Y: theme.OverscanMargin(1080), W: 160, H: 60}})
	}
	s.appFocusIDs = make([]string, len(s.Apps))
	baseY := theme.OverscanMargin(1080) + 260
	for i := range s.Apps {
		id := "app:" + s.Apps[i].ID
		s.appFocusIDs[i] = id
		nodes = append(nodes, focus.Node{ID: id, Rect: focus.Rect{X: theme.OverscanMargin(1920) + i*150, Y: baseY, W: 120, H: 120}})
	}
	s.Focus.SetNodes(nodes)
}

// Layout renders the home screen.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	if !gtx.Focused(homeTag) {
		gtx.Execute(key.FocusCmd{Tag: homeTag})
	}
	s.handleKeys(gtx)

	paint.Fill(gtx.Ops, theme.ColorBackground)

	margins := layout.Inset{
		Top:    unit.Dp(theme.OverscanMargin(1080)),
		Bottom: unit.Dp(theme.OverscanMargin(1080)),
		Left:   unit.Dp(theme.OverscanMargin(1920)),
		Right:  unit.Dp(theme.OverscanMargin(1920)),
	}
	return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceBetween}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutTabs(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutHero(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutAppRow(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutFolders(gtx)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.layoutCastShelf(gtx)
			}),
		)
	})
}

func (s *Screen) handleKeys(gtx layout.Context) {
	for {
		ev, ok := gtx.Event(
			key.Filter{Focus: homeTag, Name: key.NameLeftArrow},
			key.Filter{Focus: homeTag, Name: key.NameRightArrow},
			key.Filter{Focus: homeTag, Name: key.NameUpArrow},
			key.Filter{Focus: homeTag, Name: key.NameDownArrow},
		)
		if !ok {
			break
		}
		ke, ok := ev.(key.Event)
		if !ok || ke.State != key.Press {
			continue
		}
		switch ke.Name {
		case key.NameLeftArrow:
			s.Focus.Move(focus.Left)
		case key.NameRightArrow:
			s.Focus.Move(focus.Right)
		case key.NameUpArrow:
			s.Focus.Move(focus.Up)
		case key.NameDownArrow:
			s.Focus.Move(focus.Down)
		}
	}
}

func (s *Screen) layoutTabs(gtx layout.Context) layout.Dimensions {
	cur, _ := s.Focus.Current()
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(s.Tabs))
			for i, t := range s.Tabs {
				i, t := i, t
				focused := cur.ID == s.tabFocusIDs[i]
				children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.drawTab(gtx, t, focused)
				})
			}
			return children
		}()...,
	)
}

func (s *Screen) drawTab(gtx layout.Context, label string, focused bool) layout.Dimensions {
	// glass panel
	defer clip.UniformRRect(imageRect(gtx.Constraints.Min.X, gtx.Constraints.Min.Y, 150, 50), theme.RadiusSM).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, theme.ColorGlass)
	if focused {
		paint.Fill(gtx.Ops, theme.ColorAccent)
	}
	l := material.Body1(s.Theme, label)
	if focused {
		l.Color = theme.ColorPrimary
	} else {
		l.Color = theme.ColorSecondary
	}
	return l.Layout(gtx)
}

func (s *Screen) layoutHero(gtx layout.Context) layout.Dimensions {
	// Simple hero card
	defer clip.UniformRRect(imageRect(0, 0, gtx.Constraints.Max.X, 200), theme.RadiusLG).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, theme.ColorSurface)
	l := material.H1(s.Theme, "Continue Watching")
	l.Color = theme.ColorPrimary
	return l.Layout(gtx)
}

func (s *Screen) layoutAppRow(gtx layout.Context) layout.Dimensions {
	cur, _ := s.Focus.Current()
	return layout.Flex{Axis: layout.Horizontal, Spacing: layout.SpaceStart}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(s.Apps))
			for i, a := range s.Apps {
				i, a := i, a
				focused := cur.ID == s.appFocusIDs[i]
				children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return s.drawApp(gtx, a, focused)
				})
			}
			return children
		}()...,
	)
}

func (s *Screen) drawApp(gtx layout.Context, a AppItem, focused bool) layout.Dimensions {
	size := 100
	if focused {
		size = 110
	}
	defer clip.UniformRRect(imageRect(0, 0, size, size), theme.RadiusMD).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, theme.ColorGlass)
	if focused {
		paint.Fill(gtx.Ops, theme.ColorAccent)
	}
	l := material.Body2(s.Theme, a.Name)
	l.Color = theme.ColorPrimary
	l.Alignment = text.Middle
	return l.Layout(gtx)
}

func (s *Screen) layoutFolders(gtx layout.Context) layout.Dimensions {
	if len(s.Folders) == 0 {
		return layout.Dimensions{}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, func() []layout.FlexChild {
		children := make([]layout.FlexChild, len(s.Folders))
		for i, f := range s.Folders {
			i, f := i, f
			children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				l := material.H3(s.Theme, f.Name)
				l.Color = theme.ColorPrimary
				return l.Layout(gtx)
			})
		}
		return children
	}()...)
}

func (s *Screen) layoutCastShelf(gtx layout.Context) layout.Dimensions {
	l := material.Body1(s.Theme, "Cast to [DEVICE NAME]")
	l.Color = theme.ColorSecondary
	return l.Layout(gtx)
}

func imageRect(x, y, w, h int) image.Rectangle {
	return image.Rect(x, y, x+w, y+h)
}

