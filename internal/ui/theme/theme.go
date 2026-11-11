// Package theme defines the Fireside visual language: dark base, frosted
// glass, large rounded corners, generous spacing. Values are for a 1920x1080
// canvas and are scaled by UIScale for other resolutions.
package theme

import "image/color"

// Palette
var (
	ColorBackground = color.NRGBA{R: 0x0B, G: 0x0D, B: 0x12, A: 0xFF}
	ColorSurface    = color.NRGBA{R: 0x1A, G: 0x1D, B: 0x24, A: 0xCC}
	ColorGlass      = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0x14}
	ColorPrimary    = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
	ColorSecondary  = color.NRGBA{R: 0xA0, G: 0xA4, B: 0xAE, A: 0xFF}
	ColorAccent     = color.NRGBA{R: 0x0A, G: 0x84, B: 0xFF, A: 0xFF}
	ColorFocus      = color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
)

// Spacing in dp at 1080p.
const (
	SpaceXS = 8
	SpaceSM = 16
	SpaceMD = 24
	SpaceLG = 32
	SpaceXL = 48
)

// Corner radii in dp.
const (
	RadiusSM = 12
	RadiusMD = 20
	RadiusLG = 32
)

// Typography sizes in sp at 1080p.
const (
	FontHero  = 72
	FontH1    = 48
	FontH2    = 30
	FontBody  = 24
	FontSmall = 22 // minimum readable on TV
)

// UIScale returns a scale factor for the given screen height (1080 baseline).
func UIScale(screenHeight int) float32 {
	if screenHeight <= 0 {
		return 1
	}
	return float32(screenHeight) / 1080.0
}

// OverscanMargin returns the safe margin in px for a dimension.
func OverscanMargin(size int) int {
	return size * 5 / 100
}
