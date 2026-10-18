// Package scaling holds the math for fitting the fixed logical screen into
// the window and for mapping the mouse back into logical pixels.
//
// It has no Ebitengine dependency so it can be unit-tested anywhere.
package scaling

import "math"

// Transform maps logical screen pixels to window pixels:
// window = logical*Scale + Offset. No rotation or skew.
type Transform struct {
	Scale      float64
	OffX, OffY float64
}

// Identity is the transform before the first frame is drawn.
var Identity = Transform{Scale: 1}

// Apply maps a logical point to window pixels.
func (t Transform) Apply(x, y float64) (float64, float64) {
	return x*t.Scale + t.OffX, y*t.Scale + t.OffY
}

// Invert maps a window pixel back to logical pixels.
func (t Transform) Invert(x, y float64) (float64, float64) {
	if t.Scale == 0 {
		return x, y
	}
	return (x - t.OffX) / t.Scale, (y - t.OffY) / t.Scale
}

// PixelPerfect returns the largest whole-number scale at which a logical
// screen of logicalW×logicalH fits in a window of windowW×windowH, centered
// on whole pixels. If the window is smaller than the logical screen it
// shrinks to fit instead (then pixels can't stay square, but nothing is cut off).
func PixelPerfect(windowW, windowH, logicalW, logicalH float64) Transform {
	fit := math.Min(windowW/logicalW, windowH/logicalH)
	scale := math.Floor(fit)
	if scale < 1 {
		scale = fit
	}
	return Transform{
		Scale: scale,
		OffX:  math.Floor((windowW - logicalW*scale) / 2),
		OffY:  math.Floor((windowH - logicalH*scale) / 2),
	}
}

// RemapCursor converts a cursor position that Ebitengine reported using its
// default scaling (def) into logical pixels under the scaling we actually
// draw with (drawn). The result is clamped to the logical screen, so the
// cursor sitting on a black bar or outside the window aims at the nearest edge.
func RemapCursor(x, y float64, def, drawn Transform, logicalW, logicalH float64) (float64, float64) {
	wx, wy := def.Apply(x, y)
	lx, ly := drawn.Invert(wx, wy)
	return clamp(lx, 0, logicalW-1), clamp(ly, 0, logicalH-1)
}

func clamp(v, lo, hi float64) float64 { return math.Max(lo, math.Min(hi, v)) }
