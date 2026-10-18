package scenes

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
	"retro-shooter/internal/ui"
)

// Title is the start screen.
//
// Requiring a click here also unlocks audio: browsers (a future web build)
// refuse to play sound until the user interacts with the page, and Ebitengine
// resumes its audio context on that first input. By the time gameplay
// starts, sound is guaranteed to work.
type Title struct {
	d *Deps
	t int
}

func NewTitle(d *Deps) *Title { return &Title{d: d} }

func (s *Title) Update(in input.State) game.Scene {
	s.d.Audio.NextTick()
	s.d.Stars.Update()
	s.t++
	// The short delay stops a click from the previous screen carrying over.
	if s.t > 15 && in.FirePressed {
		s.d.Audio.Play(audio.Blip)
		return NewPlay(s.d)
	}
	return nil
}

func (s *Title) Draw(screen *ebiten.Image) {
	screen.Fill(ui.Night)
	s.d.Stars.Draw(screen, s.d.Assets.Pixel)

	ui.BigText(screen, "RETRO", 40, 4, ui.Cyan)
	ui.BigText(screen, "SHOOTER", 100, 3, ui.Yellow)

	if (s.t/30)%2 == 0 {
		ui.TextCentered(screen, "CLICK TO START", 156)
	}
	ui.TextCentered(screen, "WASD MOVE  MOUSE AIM  CLICK FIRE", 192)
	ui.TextCentered(screen, "P PAUSE   F11 FULLSCREEN", 208)
	if s.d.HighScore > 0 {
		ui.TextCentered(screen, fmt.Sprintf("HI %06d", s.d.HighScore), 4)
	}
}
