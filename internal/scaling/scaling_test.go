package scaling

import (
	"math"
	"testing"
)

const lw, lh = 320, 240

func TestPixelPerfect(t *testing.T) {
	cases := []struct {
		name string
		w, h float64
		want Transform
	}{
		{"exact 3x", 960, 720, Transform{Scale: 3}},
		{"wide window: bars left/right", 1000, 720, Transform{Scale: 3, OffX: 20}},
		{"fractional fit rounds down", 1100, 900, Transform{Scale: 3, OffX: 70, OffY: 90}},
		{"1080p fullscreen", 1920, 1080, Transform{Scale: 4, OffX: 320, OffY: 60}},
		{"odd leftover floors offset", 961, 721, Transform{Scale: 3, OffX: 0, OffY: 0}},
	}
	for _, c := range cases {
		if got := PixelPerfect(c.w, c.h, lw, lh); got != c.want {
			t.Errorf("%s: got %+v, want %+v", c.name, got, c.want)
		}
	}
}

func TestPixelPerfectShrinksTinyWindow(t *testing.T) {
	got := PixelPerfect(160, 120, lw, lh)
	if got.Scale != 0.5 {
		t.Fatalf("scale = %v, want 0.5", got.Scale)
	}
}

func TestInvertUndoesApply(t *testing.T) {
	tr := Transform{Scale: 3.25, OffX: 17, OffY: -4}
	x, y := tr.Invert(tr.Apply(123.5, 45.25))
	if math.Abs(x-123.5) > 1e-9 || math.Abs(y-45.25) > 1e-9 {
		t.Fatalf("round trip gave (%v, %v)", x, y)
	}
}

// The real-world case: Ebitengine maps the cursor with a fractional "fit"
// scale, but we draw with whole-number scaling. The remapped cursor must land
// on the logical pixel that is actually drawn under the mouse.
func TestRemapCursorMatchesWhatIsDrawn(t *testing.T) {
	const ww, wh = 1100, 900
	fit := math.Min(ww/lw, wh/lh) // 3.4375
	def := Transform{Scale: fit, OffX: (ww - lw*fit) / 2, OffY: (wh - lh*fit) / 2}
	drawn := PixelPerfect(ww, wh, lw, lh)

	for _, win := range [][2]float64{{100, 120}, {550, 450}, {1000, 800}} {
		// What Ebitengine would report for a mouse at this window pixel.
		cx, cy := def.Invert(win[0], win[1])
		gx, gy := RemapCursor(cx, cy, def, drawn, lw, lh)
		// What's really drawn there.
		wantX, wantY := drawn.Invert(win[0], win[1])
		wantX, wantY = clamp(wantX, 0, lw-1), clamp(wantY, 0, lh-1)
		if math.Abs(gx-wantX) > 1e-9 || math.Abs(gy-wantY) > 1e-9 {
			t.Errorf("window %v: got (%v, %v), want (%v, %v)", win, gx, gy, wantX, wantY)
		}
	}
}

func TestRemapCursorClampsToScreen(t *testing.T) {
	drawn := PixelPerfect(1000, 720, lw, lh) // 20 px bars left and right
	x, y := RemapCursor(-50, 9999, Identity, drawn, lw, lh)
	if x != 0 || y != lh-1 {
		t.Fatalf("got (%v, %v), want (0, %v)", x, y, lh-1)
	}
}

func TestRemapCursorBeforeFirstFrame(t *testing.T) {
	x, y := RemapCursor(10, 20, Identity, Identity, lw, lh)
	if x != 10 || y != 20 {
		t.Fatalf("got (%v, %v), want (10, 20)", x, y)
	}
}
