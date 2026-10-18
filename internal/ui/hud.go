package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/world"
)

// HUD layout: a row of HP segments at the top left.
const (
	hpY      = 5 // top of the HP segments
	hpSegW   = 8 // segment size, including its dark frame
	hpSegH   = 6
	hpSegGap = 2 // space between segments
	hpInset  = 1 // frame thickness around the red fill
)

// DrawHUD draws HP, wave number and score along the top of the screen.
// It reads the world but never changes it.
func DrawHUD(dst *ebiten.Image, w *world.World, wave int) {
	// HP: one segment per hit point.
	p := w.Player
	for i := range p.MaxHP {
		x := float64(Margin + i*(hpSegW+hpSegGap))
		Rect(dst, x, hpY, hpSegW, hpSegH, Dark)
		if i < p.HP {
			Rect(dst, x+hpInset, hpY+hpInset, hpSegW-2*hpInset, hpSegH-2*hpInset, Red)
		}
	}

	TextCentered(dst, fmt.Sprintf("WAVE %d", wave), 0)

	score := fmt.Sprintf("%06d", w.Score)
	Text(dst, score, world.ScreenW-len(score)*CharW-Margin, 0)
}
