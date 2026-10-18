package world

import (
	"testing"

	"retro-shooter/internal/input"
)

// busyWorld has a full screen of enemies and the player firing constantly.
func busyWorld() *World {
	w := New()
	for i := range maxEnemies {
		w.Spawn(EnemyKind(i%int(NumEnemyKinds)), Vec{float64(i * 5 % ScreenW), -10})
	}
	w.Player.HP = 1 << 30 // never dies during the benchmark
	return w
}

func BenchmarkWorldUpdate(b *testing.B) {
	w := busyWorld()
	in := input.State{Fire: true, Right: true, MouseX: ScreenW, MouseY: 0}
	b.ReportAllocs()
	for b.Loop() {
		w.Update(in)
		w.DrainEvents()
		if len(w.Enemies) < maxEnemies/2 {
			for i := len(w.Enemies); i < maxEnemies; i++ {
				w.Spawn(Grunt, Vec{float64(i * 5 % ScreenW), -10})
			}
		}
	}
}
