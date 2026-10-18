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

// Game-over screen layout in logical pixels.
const (
	gameOverY     = 50 // "GAME OVER" heading
	gameOverScale = 3
	resultsY      = 110 // score, then wave on the next line
	newHighY      = resultsY + 2*ui.LineH + 4
	promptY       = 184 // "play again" / "title" prompts, one per line
)

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

	ui.BigText(screen, "GAME OVER", gameOverY, gameOverScale, ui.Red)
	ui.TextCentered(screen, fmt.Sprintf("SCORE  %06d", s.score), resultsY)
	ui.TextCentered(screen, fmt.Sprintf("WAVE   %d", s.wave), resultsY+ui.LineH)
	if s.newHigh && (s.t/20)%2 == 0 {
		ui.TextCentered(screen, "NEW HIGH SCORE!", newHighY)
	}
	if s.t >= acceptDelay {
		ui.TextCentered(screen, "CLICK TO PLAY AGAIN", promptY)
		ui.TextCentered(screen, "ESC FOR TITLE", promptY+ui.LineH)
	}
}
