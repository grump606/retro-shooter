// Package assets builds the game's images once at startup.
//
// Sprites are pixel art defined as text grids (enemy sprites live with the
// rest of each enemy's settings in world.EnemyDefs), so the game runs with no
// downloaded files. To use real art (e.g. Kenney's CC0 "Pixel Shmup" pack),
// embed PNGs with //go:embed and decode them in Load in place of sprite().
// Keep the ship pointing right.
package assets

import (
	"fmt"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/world"
)

// Assets holds every image the game draws. Treat it as read-only.
type Assets struct {
	Ship         *ebiten.Image
	Enemies      [world.NumEnemyKinds]*ebiten.Image // indexed by world.EnemyKind
	PlayerBullet *ebiten.Image
	EnemyBullet  *ebiten.Image
	Pixel        *ebiten.Image // 1×1 white; scale and tint it for rects and particles
}

// Palette maps the letters used in sprite grids (and EnemyDef.Debris) to
// colors: a small retro set inspired by PICO-8.
var Palette = map[byte]color.RGBA{
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

// ShipSprite is the player ship. It points right (angle 0); the game rotates
// it toward the mouse. Also used to generate the app icon (tools/mkicon).
var ShipSprite = []string{
	"oo.kk......",
	".okbbk.....",
	"..kbbbkk...",
	".okbbbbcck.",
	"ookbbbbccwk",
	".okbbbbcck.",
	"..kbbbkk...",
	".okbbk.....",
	"oo.kk......",
}

// SpriteRGBA renders a sprite grid to a plain image (no GPU needed).
func SpriteRGBA(rows []string) *image.RGBA { return spriteRGBA(rows) }

// Color returns the palette color for a letter. It panics on an unknown
// letter, which the tests catch before the game ships.
func Color(letter byte) color.RGBA {
	c, ok := Palette[letter]
	if !ok {
		panic(fmt.Sprintf("assets: unknown palette letter %q", letter))
	}
	return c
}

// Load builds all images. Call it once at startup; it is safe to call before
// ebiten.RunGame.
func Load() *Assets {
	px := ebiten.NewImage(1, 1)
	px.Fill(color.White)
	a := &Assets{
		Ship: sprite(ShipSprite...),
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
	for k := range world.NumEnemyKinds {
		a.Enemies[k] = sprite(k.Def().Sprite...)
	}
	return a
}

// sprite turns rows of palette letters into an image; '.' is transparent.
func sprite(rows ...string) *ebiten.Image {
	return ebiten.NewImageFromImage(spriteRGBA(rows))
}

func spriteRGBA(rows []string) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, len(rows[0]), len(rows)))
	for y, row := range rows {
		for x := range len(row) {
			if row[x] != '.' {
				img.SetRGBA(x, y, Color(row[x]))
			}
		}
	}
	return img
}
