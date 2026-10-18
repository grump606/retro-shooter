// Package ui draws text, bars and the in-game HUD. It only reads game state.
package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"retro-shooter/internal/world"
)

// Ebitengine's built-in debug font: 6 px per character, 16 px per line.
const (
	CharW = 6
	LineH = 16
)

// Retro colors shared by the HUD and scenes.
var (
	White  = color.RGBA{0xff, 0xf1, 0xe8, 0xff}
	Grey   = color.RGBA{0x83, 0x76, 0x9c, 0xff}
	Dark   = color.RGBA{0x1d, 0x2b, 0x53, 0xff}
	Red    = color.RGBA{0xff, 0x00, 0x4d, 0xff}
	Yellow = color.RGBA{0xff, 0xec, 0x27, 0xff}
	Cyan   = color.RGBA{0x29, 0xad, 0xff, 0xff}
	Night  = color.RGBA{0x0b, 0x0d, 0x1a, 0xff} // background
)

// Text draws s with its top-left corner at (x, y).
func Text(dst *ebiten.Image, s string, x, y int) {
	ebitenutil.DebugPrintAt(dst, s, x, y)
}

// TextCentered draws s horizontally centered on the screen.
func TextCentered(dst *ebiten.Image, s string, y int) {
	Text(dst, s, (world.ScreenW-len(s)*CharW)/2, y)
}

// bigCache holds pre-rendered text so BigText doesn't redraw glyphs every frame.
var bigCache = map[string]*ebiten.Image{}

// BigText draws s scaled up by a whole number and tinted, centered on the
// screen at height y. Used for titles and banners.
func BigText(dst *ebiten.Image, s string, y int, scale float64, clr color.Color) {
	img, ok := bigCache[s]
	if !ok {
		img = ebiten.NewImage(len(s)*CharW, LineH)
		ebitenutil.DebugPrint(img, s)
		bigCache[s] = img
	}
	w := float64(img.Bounds().Dx()) * scale
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate((world.ScreenW-w)/2, float64(y))
	op.ColorScale.ScaleWithColor(clr)
	dst.DrawImage(img, op)
}

// Rect fills a rectangle.
func Rect(dst *ebiten.Image, x, y, w, h float64, clr color.Color) {
	vector.FillRect(dst, float32(x), float32(y), float32(w), float32(h), clr, false)
}
