package game

import (
	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/input"
)

// Scene is one screen of the game (title, play, game over).
//
// Update returns the scene to switch to, or nil to stay. Scenes ask for
// transitions this way instead of reaching into the Game.
type Scene interface {
	Update(in input.State) Scene
	Draw(screen *ebiten.Image)
}
