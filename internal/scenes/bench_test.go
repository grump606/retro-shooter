package scenes

import (
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/assets"
	"retro-shooter/internal/world"
)

// BenchmarkDrawFrame measures CPU-side cost (mostly allocations) of drawing a
// busy frame: starfield, a full set of particles, enemies and bullets.
func BenchmarkDrawFrame(b *testing.B) {
	a := assets.Load()
	stars := NewStarfield(1)
	fx := NewParticles(1)
	for range 20 {
		fx.Burst(world.Vec{X: 160, Y: 120}, 25, 2, assets.Color('o'))
	}
	w := world.New()
	for i := range 64 {
		w.Spawn(world.EnemyKind(i%int(world.NumEnemyKinds)), world.Vec{X: float64(i * 5), Y: float64(i * 3)})
	}
	dst := ebiten.NewImage(world.ScreenW, world.ScreenH)
	b.ReportAllocs()
	for b.Loop() {
		stars.Draw(dst, a.Pixel)
		fx.Draw(dst, a.Pixel)
		drawWorld(dst, a, w)
	}
}
