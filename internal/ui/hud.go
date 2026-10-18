package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/world"
)

// DrawHUD draws HP, wave number and score along the top of the screen.
// It reads the world but never changes it.
func DrawHUD(dst *ebiten.Image, w *world.World, wave int) {
	// HP: one segment per hit point.
	p := w.Player
	for i := range p.MaxHP {
		x := 4 + float64(i)*10
		Rect(dst, x, 5, 8, 6, Dark)
		if i < p.HP {
			Rect(dst, x+1, 6, 6, 4, Red)
		}
	}

	TextCentered(dst, fmt.Sprintf("WAVE %d", wave), 0)

	score := fmt.Sprintf("%06d", w.Score)
	Text(dst, score, world.ScreenW-len(score)*CharW-4, 0)
}
