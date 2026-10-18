package scenes

import (
	"math/rand/v2"
	"testing"
)

func TestShakeOffsetIsWholePixelsAndBounded(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	for ticks := 0; ticks <= 30; ticks++ {
		limit := max(1, min(ticks/2, maxShakePx))
		if ticks == 0 {
			limit = 0
		}
		for range 200 {
			dx, dy := shakeOffset(ticks, r)
			if dx < -limit || dx > limit || dy < -limit || dy > limit {
				t.Fatalf("ticks=%d: offset (%d,%d) outside ±%d", ticks, dx, dy, limit)
			}
		}
	}
}

func TestShakeOffsetUsesFullRange(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	seen := map[int]bool{}
	for range 1000 {
		dx, _ := shakeOffset(30, r)
		seen[dx] = true
	}
	for v := -maxShakePx; v <= maxShakePx; v++ {
		if !seen[v] {
			t.Errorf("strong shake never produced offset %d", v)
		}
	}
}
