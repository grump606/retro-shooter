// Package assets builds the game's images once at startup.
//
// The sprites are pixel art defined as text grids below, so the game runs with
// no downloaded files. To use real art (e.g. Kenney's CC0 "Pixel Shmup" pack),
// drop PNGs in assets/images, embed them with //go:embed, and decode them in
// Load in place of the sprite() calls. Keep the ship pointing right.
package assets

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Assets holds every image the game draws. Treat it as read-only.
type Assets struct {
	Ship         *ebiten.Image
	Grunt        *ebiten.Image
	Fast         *ebiten.Image
	Tank         *ebiten.Image
	PlayerBullet *ebiten.Image
	EnemyBullet  *ebiten.Image
	Pixel        *ebiten.Image // 1×1 white; scale and tint it for rects and particles
}

// Palette: a small retro set (PICO-8 inspired).
var palette = map[byte]color.RGBA{
	'k': {0x1d, 0x2b, 0x53, 0xff}, // dark blue outline
	'w': {0xff, 0xf1, 0xe8, 0xff}, // white
	'c': {0x29, 0xad, 0xff, 0xff}, // cyan
	'b': {0x83, 0x76, 0x9c, 0xff}, // lavender grey
	'r': {0xff, 0x00, 0x4d, 0xff}, // red
	'o': {0xff, 0xa3, 0x00, 0xff}, // orange
	'y': {0xff, 0xec, 0x27, 0xff}, // yellow
	'g': {0x00, 0xe4, 0x36, 0xff}, // green
	'p': {0x7e, 0x25, 0x53, 0xff}, // plum
	'm': {0xff, 0x77, 0xa8, 0xff}, // pink
}

// Load builds all images. Call it once at startup; it is safe to call before
// ebiten.RunGame.
func Load() *Assets {
	px := ebiten.NewImage(1, 1)
	px.Fill(color.White)
	return &Assets{
		// Points right (angle 0); the game rotates it toward the mouse.
		Ship: sprite(
			"oo.kk......",
			".okbbk.....",
			"..kbbbkk...",
			".okbbbbcck.",
			"ookbbbbccwk",
			".okbbbbcck.",
			"..kbbbkk...",
			".okbbk.....",
			"oo.kk......",
		),
		Grunt: sprite(
			"..kkkkk..",
			".kgggggk.",
			"kgwkgwkgk",
			"kgggggggk",
			"kgggggggk",
			".kgkgkgk.",
			"..k.k.k..",
		),
		Fast: sprite(
			"...k...",
			"..kyk..",
			".kyoyk.",
			"kyoroyk",
			".kyoyk.",
			"..kyk..",
			"...k...",
		),
		Tank: sprite(
			"....kkkkkkk....",
			"..kkpppppppkk..",
			".kppmmmmmmmppk.",
			".kpmmpppppmmpk.",
			"kpmmpkkkkkpmmpk",
			"kpmpkrrrrrkpmpk",
			"kpmpkrwwwrkpmpk",
			"kpmpkrwwwrkpmpk",
			"kpmpkrrrrrkpmpk",
			"kpmmpkkkkkpmmpk",
			".kpmmpppppmmpk.",
			".kppmmmmmmmppk.",
			"..kkpppppppkk..",
			"....kkkkkkk....",
		),
		PlayerBullet: sprite(
			".c.",
			"cwc",
			".c.",
		),
		EnemyBullet: sprite(
			".o.",
			"oyo",
			".o.",
		),
		Pixel: px,
	}
}

// sprite turns rows of palette letters into an image; '.' is transparent.
func sprite(rows ...string) *ebiten.Image {
	img := image.NewRGBA(image.Rect(0, 0, len(rows[0]), len(rows)))
	for y, row := range rows {
		for x := range len(row) {
			if c, ok := palette[row[x]]; ok {
				img.SetRGBA(x, y, c)
			}
		}
	}
	return ebiten.NewImageFromImage(img)
}
