// Package scenes contains the title, play and game-over screens.
package scenes

import (
	"retro-shooter/assets"
	"retro-shooter/internal/audio"
)

// Deps is shared by every scene. It is built once in main.go and passed to
// each scene's constructor, so there are no globals.
type Deps struct {
	Assets *assets.Assets
	Audio  *audio.Player
	Stars  *Starfield // shared so the background scrolls smoothly across scenes

	HighScore int  // best score this session
	Debug     bool // F3 overlay
}
