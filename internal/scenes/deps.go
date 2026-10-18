// Package scenes contains the title, play and game-over screens.
package scenes

import (
	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/assets"
	"retro-shooter/internal/audio"
	"retro-shooter/internal/world"
)

// Deps is shared by every scene. It is built once in main.go and passed to
// each scene's constructor, so there are no globals.
type Deps struct {
	Assets *assets.Assets
	Audio  *audio.Player
	Stars  *Starfield // shared so the background scrolls smoothly across scenes

	HighScore int  // best score this session
	Debug     bool // F3 overlay

	canvas *ebiten.Image // see Canvas
}

// Canvas returns the off-screen image the play scene draws the world onto
// (so it can be shaken). It's created once and reused by every run, instead
// of allocating a new GPU texture on each restart.
func (d *Deps) Canvas() *ebiten.Image {
	if d.canvas == nil {
		d.canvas = ebiten.NewImage(world.ScreenW, world.ScreenH)
	}
	return d.canvas
}
