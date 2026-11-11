package home

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/io/key"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"

	"github.com/ngbcoder/fireside/internal/apps"
	"github.com/ngbcoder/fireside/internal/ui/focus"
	"github.com/ngbcoder/fireside/internal/ui/theme"
)

var homeTag = new(int)

// baseline is the reference resolution the theme sizes are designed for.
const baseline = 1080

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

	// Measured rectangles for the current frame, used by the focus engine.
	tabRects map[string]focus.Rect
	appRects map[string]focus.Rect
}

// AppItem represents an app on the home row.
type AppItem struct {
	ID   string
	Name string
}

// New creates a Home screen with default content.
func New(th *material.Theme) *Screen {
	s := &Screen{
		Theme:    th,
		Focus:    focus.NewManager(),
		Tabs:     []string{"Home", "Library", "App Store", "Search", "Settings"},
		tabRects: map[string]focus.Rect{},
		appRects: map[string]focus.Rect{},
	}
	s.Apps = []AppItem{
		{ID: "youtube", Name: "YouTube"},
		{ID: "disney", Name: "Disney+"},
		{ID: "appletv", Name: "Apple TV"},
		{ID: "media", Name: "Media Player"},
		{ID: "music", Name: "Music"},
		{ID: "settings", Name: "Settings"},
	}
	s.rebuildFocusIDs()
	return s
}

func (s *Screen) rebuildFocusIDs() {
	s.tabFocusIDs = make([]string, len(s.Tabs))
	for i := range s.Tabs {
		s.tabFocusIDs[i] = "tab:" + s.Tabs[i]
	}
	s.appFocusIDs = make([]string, len(s.Apps))
	for i := range s.Apps {
		s.appFocusIDs[i] = "app:" + s.Apps[i].ID
	}
}

// scale returns the UI scale factor for the current window height.
func (s *Screen) scale(gtx layout.Context) float32 {
	h := gtx.Constraints.Max.Y
	if h <= 0 {
		return 1
	}
	sc := float32(h) / baseline
	if sc < 0.5 {
		sc = 0.5
	}
	return sc
}

// sp converts a theme font size to scaled sp.
func (s *Screen) sp(gtx layout.Context, size int) unit.Sp {
	return unit.Sp(float32(size) * s.scale(gtx))
}

// dp converts a theme size to scaled dp.
func (s *Screen) dp(gtx layout.Context, size int) unit.Dp {
	return unit.Dp(float32(size) * s.scale(gtx))
}

// text lays out a single-line label with the given scaled size.
func (s *Screen) label(gtx layout.Context, str string, size int, c color.NRGBA, align text.Alignment) layout.Dimensions {
	return widget.Label{Alignment: align, MaxLines: 1}.Layout(
		gtx, s.Theme.Shaper, font.Font{Typeface: s.Theme.Face}, s.sp(gtx, size), str, op.CallOp{})
}

// Layout renders the home screen.
func (s *Screen) Layout(gtx layout.Context) layout.Dimensions {
	if !gtx.Focused(homeTag) {
		gtx.Execute(key.FocusCmd{Tag: homeTag})
	}
	s.handleKeys(gtx)

	paint.Fill(gtx.Ops, theme.ColorBackground)

	s.tabRects = map[string]focus.Rect{}
	s.appRects = map[string]focus.Rect{}

	margin := int(theme.OverscanMargin(gtx.Constraints.Max.X))
	margins := layout.Inset{
		Top:    unit.Dp(margin),
		Bottom: unit.Dp(margin),
		Left:   unit.Dp(margin),
		Right:  unit.Dp(margin),
	}
	dims := margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
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
	s.syncFocusNodes()
	return dims
}

// syncFocusNodes updates the focus graph with the measured rectangles.
func (s *Screen) syncFocusNodes() {
	var nodes []focus.Node
	for _, id := range s.tabFocusIDs {
		nodes = append(nodes, focus.Node{ID: id, Rect: s.tabRects[id]})
	}
	for _, id := range s.appFocusIDs {
		nodes = append(nodes, focus.Node{ID: id, Rect: s.appRects[id]})
	}
	s.Focus.SetNodes(nodes)
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
	gap := int(s.dp(gtx, theme.SpaceSM))
	x := 0
	return layout.Flex{Axis: layout.Horizontal, Gap: gap}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(s.Tabs))
			for i, t := range s.Tabs {
				i, t := i, t
				focused := cur.ID == s.tabFocusIDs[i]
				children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := s.drawTab(gtx, t, focused)
					s.tabRects[s.tabFocusIDs[i]] = focus.Rect{X: x, Y: 0, W: d.Size.X, H: d.Size.Y}
					x += d.Size.X + gap
					return d
				})
			}
			return children
		}()...,
	)
}

func (s *Screen) drawTab(gtx layout.Context, label string, focused bool) layout.Dimensions {
	size := theme.FontBody
	col := theme.ColorSecondary
	bg := theme.ColorGlass
	if focused {
		col = theme.ColorPrimary
		bg = theme.ColorAccent
	}
	// Measure the text so the panel fits it.
	macro := op.Record(gtx.Ops)
	dims := s.label(gtx, label, size, col, text.Middle)
	macro.Stop()

	padX := int(s.dp(gtx, theme.SpaceSM))
	padY := int(s.dp(gtx, theme.SpaceXS))
	w := dims.Size.X + 2*padX
	h := dims.Size.Y + 2*padY

	// Background panel.
	defer clip.UniformRRect(image.Rect(0, 0, w, h), int(s.dp(gtx, theme.RadiusSM))).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, bg)

	// Text on top, centered in the panel.
	gtx2 := gtx
	gtx2.Constraints.Min = image.Point{}
	gtx2.Constraints.Max = image.Pt(w, h)
	layout.Center.Layout(gtx2, func(gtx layout.Context) layout.Dimensions {
		return s.label(gtx, label, size, col, text.Middle)
	})

	return layout.Dimensions{Size: image.Pt(w, h)}
}

func (s *Screen) layoutHero(gtx layout.Context) layout.Dimensions {
	title := "Continue Watching"
	sub := "Resume where you left off"

	titleSize := theme.FontHero
	titleCol := theme.ColorPrimary
	subCol := theme.ColorSecondary

	// Measure.
	macro := op.Record(gtx.Ops)
	titleDims := s.label(gtx, title, titleSize, titleCol, text.Start)
	subDims := s.label(gtx, sub, theme.FontBody, subCol, text.Start)
	macro.Stop()

	pad := int(s.dp(gtx, theme.SpaceLG))
	w := titleDims.Size.X
	if subDims.Size.X > w {
		w = subDims.Size.X
	}
	w += 2 * pad
	h := titleDims.Size.Y + subDims.Size.Y + 3*pad

	// Panel.
	defer clip.UniformRRect(image.Rect(0, 0, w, h), int(s.dp(gtx, theme.RadiusLG))).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, theme.ColorSurface)

	// Content.
	gtx2 := gtx
	gtx2.Constraints.Min = image.Point{}
	gtx2.Constraints.Max = image.Pt(w, h)
	layout.UniformInset(unit.Dp(pad)).Layout(gtx2, func(gtx layout.Context) layout.Dimensions {
		return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.label(gtx, title, titleSize, titleCol, text.Start)
			}),
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return s.label(gtx, sub, theme.FontBody, subCol, text.Start)
			}),
		)
	})

	return layout.Dimensions{Size: image.Pt(w, h)}
}

func (s *Screen) layoutAppRow(gtx layout.Context) layout.Dimensions {
	cur, _ := s.Focus.Current()
	gap := int(s.dp(gtx, theme.SpaceSM))
	x := 0
	return layout.Flex{Axis: layout.Horizontal, Gap: gap}.Layout(gtx,
		func() []layout.FlexChild {
			children := make([]layout.FlexChild, len(s.Apps))
			for i, a := range s.Apps {
				i, a := i, a
				focused := cur.ID == s.appFocusIDs[i]
				children[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					d := s.drawApp(gtx, a, focused)
					s.appRects[s.appFocusIDs[i]] = focus.Rect{X: x, Y: 0, W: d.Size.X, H: d.Size.Y}
					x += d.Size.X + gap
					return d
				})
			}
			return children
		}()...,
	)
}

func (s *Screen) drawApp(gtx layout.Context, a AppItem, focused bool) layout.Dimensions {
	tile := int(s.dp(gtx, 100))
	if focused {
		tile = int(s.dp(gtx, 112))
	}
	labelSize := theme.FontSmall
	labelCol := theme.ColorSecondary
	if focused {
		labelCol = theme.ColorPrimary
	}
	bg := theme.ColorGlass
	if focused {
		bg = theme.ColorAccent
	}

	macro := op.Record(gtx.Ops)
	labelDims := s.label(gtx, a.Name, labelSize, labelCol, text.Middle)
	macro.Stop()

	w := labelDims.Size.X
	if tile > w {
		w = tile
	}
	labelGap := int(s.dp(gtx, theme.SpaceXS))
	h := tile + labelGap + labelDims.Size.Y

	// Tile.
	defer clip.UniformRRect(image.Rect(0, 0, tile, tile), int(s.dp(gtx, theme.RadiusMD))).Push(gtx.Ops).Pop()
	paint.Fill(gtx.Ops, bg)

	// Tile initial.
	gtx2 := gtx
	gtx2.Constraints.Min = image.Point{}
	gtx2.Constraints.Max = image.Pt(tile, tile)
	layout.Center.Layout(gtx2, func(gtx layout.Context) layout.Dimensions {
		initial := string([]rune(a.Name)[0])
		return s.label(gtx, initial, theme.FontH2, theme.ColorPrimary, text.Middle)
	})

	// Label under the tile, centered.
	gtx3 := gtx
	gtx3.Constraints.Min = image.Point{}
	gtx3.Constraints.Max = image.Pt(w, h)
	layout.Stack{Alignment: layout.NW}.Layout(gtx3,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(w, h)}
		}),
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			gtx4 := gtx
			gtx4.Constraints.Min = image.Point{X: 0, Y: tile + labelGap}
			gtx4.Constraints.Max = image.Pt(w, tile+labelGap+labelDims.Size.Y)
			return layout.Center.Layout(gtx4, func(gtx layout.Context) layout.Dimensions {
				return s.label(gtx, a.Name, labelSize, labelCol, text.Middle)
			})
		}),
	)

	return layout.Dimensions{Size: image.Pt(w, h)}
}

func (s *Screen) layoutFolders(gtx layout.Context) layout.Dimensions {
	if len(s.Folders) == 0 {
		return layout.Dimensions{}
	}
	return layout.Flex{Axis: layout.Vertical}.Layout(gtx, func() []layout.FlexChild {
		children := make([]layout.FlexChild, 0, len(s.Folders))
		for _, f := range s.Folders {
			f := f
			children = append(children, layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				dims := s.label(gtx, f.Name, theme.FontH2, theme.ColorPrimary, text.Start)
				items := f.AppsInFolder(s.Registry)
				row := layout.Flex{Axis: layout.Horizontal, Gap: int(s.dp(gtx, theme.SpaceSM))}
				row.Layout(gtx, func() []layout.FlexChild {
					ch := make([]layout.FlexChild, len(items))
					for i, a := range items {
						i, a := i, a
						ch[i] = layout.Rigid(func(gtx layout.Context) layout.Dimensions {
							return s.label(gtx, a.Name, theme.FontSmall, theme.ColorSecondary, text.Start)
						})
					}
					return ch
				}()...)
				return dims
			}))
		}
		return children
	}()...)
}

func (s *Screen) layoutCastShelf(gtx layout.Context) layout.Dimensions {
	return s.label(gtx, "Cast to [DEVICE NAME]", theme.FontBody, theme.ColorSecondary, text.Start)
}
