package assets

import (
	"testing"

	"retro-shooter/internal/world"
)

// Every enemy definition must be complete and drawable, so adding an enemy
// with a typo fails here instead of crashing the game.
func TestEnemyDefsAreComplete(t *testing.T) {
	names := map[string]bool{}
	for k := range world.NumEnemyKinds {
		d := k.Def()
		if d.Name == "" {
			t.Errorf("enemy %d: missing Name", k)
			continue
		}
		if names[d.Name] {
			t.Errorf("enemy %q: duplicate Name", d.Name)
		}
		names[d.Name] = true

		if d.HP <= 0 || d.Speed <= 0 || d.Radius <= 0 || d.Score < 0 || d.FireEvery < 0 {
			t.Errorf("enemy %q: invalid stats %+v", d.Name, *d)
		}
		if d.DebrisCount <= 0 {
			t.Errorf("enemy %q: DebrisCount must be positive", d.Name)
		}
		if _, ok := Palette[d.Debris]; !ok {
			t.Errorf("enemy %q: Debris %q is not a palette letter", d.Name, d.Debris)
		}
		checkSprite(t, d.Name, d.Sprite)
	}
}

func checkSprite(t *testing.T, name string, rows []string) {
	t.Helper()
	if len(rows) == 0 || len(rows[0]) == 0 {
		t.Errorf("enemy %q: empty sprite", name)
		return
	}
	for y, row := range rows {
		if len(row) != len(rows[0]) {
			t.Errorf("enemy %q: sprite row %d is %d wide, want %d", name, y, len(row), len(rows[0]))
		}
		for x := range len(row) {
			if c := row[x]; c != '.' {
				if _, ok := Palette[c]; !ok {
					t.Errorf("enemy %q: sprite (%d,%d) has unknown letter %q", name, x, y, c)
				}
			}
		}
	}
}

func TestSpriteRGBA(t *testing.T) {
	img := spriteRGBA([]string{"k.", ".w"})
	if img.Bounds().Dx() != 2 || img.Bounds().Dy() != 2 {
		t.Fatalf("size = %v", img.Bounds())
	}
	if img.RGBAAt(0, 0) != Palette['k'] || img.RGBAAt(1, 1) != Palette['w'] {
		t.Error("colored pixels don't match palette")
	}
	if img.RGBAAt(1, 0).A != 0 {
		t.Error("'.' should be transparent")
	}
}
