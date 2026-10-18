package level

import (
	"testing"

	"retro-shooter/internal/world"
)

func TestSpawnPointsAreOffScreenAndAwayFromPlayer(t *testing.T) {
	m := New(1)
	corners := []world.Vec{{X: 5, Y: 5}, {X: world.ScreenW - 5, Y: world.ScreenH / 2}, {X: world.ScreenW / 2, Y: world.ScreenH / 2}}
	for _, player := range corners {
		for range 200 {
			p := m.spawnPoint(player)
			inside := p.X > 0 && p.X < world.ScreenW && p.Y > 0 && p.Y < world.ScreenH
			if inside {
				t.Fatalf("spawn %+v is on screen", p)
			}
			if p.DistTo(player) < minSpawnDist {
				t.Fatalf("spawn %+v too close to player %+v", p, player)
			}
		}
	}
}

func TestWaveAdvancesWhenCleared(t *testing.T) {
	m := New(1)
	w := world.New()
	if m.Wave() != 1 || !m.InIntermission() {
		t.Fatal("expected to start in wave 1 intermission")
	}
	// Run until everything in wave 1 has spawned, clearing enemies as they appear.
	for range 10_000 {
		m.Update(w)
		w.Enemies = nil // pretend the player killed everything
		if m.Wave() == 2 {
			return
		}
	}
	t.Fatalf("still on wave %d", m.Wave())
}

func TestWavesKeepGoing(t *testing.T) {
	for n := 1; n <= 30; n++ {
		if len(waveFor(n)) == 0 {
			t.Fatalf("wave %d is empty", n)
		}
	}
}
