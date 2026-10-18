package ui

import (
	"fmt"
	"testing"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/world"
)

func TestBigTextCacheIsBounded(t *testing.T) {
	dst := ebiten.NewImage(world.ScreenW, world.ScreenH)
	// An endless run: one new banner per wave.
	for wave := 1; wave <= 500; wave++ {
		BigText(dst, fmt.Sprintf("WAVE %d", wave), 90, 3, Yellow)
		if len(bigCache) > bigCacheMax {
			t.Fatalf("cache grew to %d entries at wave %d (max %d)", len(bigCache), wave, bigCacheMax)
		}
	}
	// A string drawn again is served from the cache.
	BigText(dst, "TITLE", 0, 1, White)
	img := bigCache["TITLE"]
	BigText(dst, "TITLE", 0, 1, White)
	if bigCache["TITLE"] != img {
		t.Fatal("repeated text was re-rendered instead of cached")
	}
}
