package scenes

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
	"retro-shooter/internal/ui"
)

// acceptDelay stops a panicked click right after dying from skipping this screen.
const acceptDelay = 45

// GameOver shows the final score. Click to play again, Esc for the title.
type GameOver struct {
	d       *Deps
	score   int
	wave    int
	newHigh bool
	t       int
}

func NewGameOver(d *Deps, score, wave int, newHigh bool) *GameOver {
	return &GameOver{d: d, score: score, wave: wave, newHigh: newHigh}
}

func (s *GameOver) Update(in input.State) game.Scene {
	s.d.Audio.NextTick()
	s.d.Stars.Update()
	s.t++
	if s.t < acceptDelay {
		return nil
	}
	if in.FirePressed {
		s.d.Audio.Play(audio.Blip)
		return NewPlay(s.d) // fresh World: nothing carries over
	}
	if in.PausePressed {
		return NewTitle(s.d)
	}
	return nil
}

func (s *GameOver) Draw(screen *ebiten.Image) {
	screen.Fill(ui.Night)
	s.d.Stars.Draw(screen, s.d.Assets.Pixel)

	ui.BigText(screen, "GAME OVER", 50, 3, ui.Red)
	ui.TextCentered(screen, fmt.Sprintf("SCORE  %06d", s.score), 110)
	ui.TextCentered(screen, fmt.Sprintf("WAVE   %d", s.wave), 126)
	if s.newHigh && (s.t/20)%2 == 0 {
		ui.TextCentered(screen, "NEW HIGH SCORE!", 146)
	}
	if s.t >= acceptDelay {
		ui.TextCentered(screen, "CLICK TO PLAY AGAIN", 184)
		ui.TextCentered(screen, "ESC FOR TITLE", 200)
	}
}
