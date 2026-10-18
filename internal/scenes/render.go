package scenes

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/colorm"

	"retro-shooter/assets"
	"retro-shooter/internal/world"
)

// drawWorld renders the World's entities. Like the HUD, it only reads state.
func drawWorld(dst *ebiten.Image, a *assets.Assets, w *world.World) {
	for _, b := range w.Bullets {
		img := a.EnemyBullet
		if b.FromPlayer {
			img = a.PlayerBullet
		}
		drawCentered(dst, img, b.Pos, 0)
	}

	for _, e := range w.Enemies {
		img := a.Enemies[e.Kind]
		if e.Flash > 0 {
			// Flash pure white when hit.
			var cm colorm.ColorM
			cm.Scale(0, 0, 0, 1)
			cm.Translate(1, 1, 1, 0)
			op := &colorm.DrawImageOptions{}
			op.GeoM = centeredGeoM(img, e.Pos, 0)
			colorm.DrawImage(dst, img, cm, op)
			continue
		}
		drawCentered(dst, img, e.Pos, 0)
	}

	p := w.Player
	blinkOff := p.Invuln > 0 && (p.Invuln/4)%2 == 1
	if p.HP > 0 && !blinkOff {
		drawCentered(dst, a.Ship, p.Pos, p.Angle)
	}
}

// centeredGeoM rotates img around its center and places that center at pos.
func centeredGeoM(img *ebiten.Image, pos world.Vec, angle float64) ebiten.GeoM {
	var m ebiten.GeoM
	b := img.Bounds()
	m.Translate(-float64(b.Dx())/2, -float64(b.Dy())/2)
	if angle != 0 {
		m.Rotate(angle)
	}
	// Round so unrotated sprites land on whole pixels and stay sharp.
	m.Translate(math.Round(pos.X), math.Round(pos.Y))
	return m
}

func drawCentered(dst, img *ebiten.Image, pos world.Vec, angle float64) {
	op := &ebiten.DrawImageOptions{GeoM: centeredGeoM(img, pos, angle)}
	dst.DrawImage(img, op)
}
