// Package game implements ebiten.Game: it polls input once per tick, runs the
// current Scene, and scales the 320×240 screen to the window.
package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"retro-shooter/internal/input"
	"retro-shooter/internal/scaling"
	"retro-shooter/internal/world"
)

// Game is the top-level ebiten.Game.
type Game struct {
	scene Scene

	// Transforms from the last frame, used to map the cursor. See cursor().
	defaultT scaling.Transform // how Ebitengine would draw the screen
	drawnT   scaling.Transform // how we actually draw it (whole-number scale)
}

// New returns a Game starting on the given scene.
func New(first Scene) *Game {
	return &Game{scene: first, defaultT: scaling.Identity, drawnT: scaling.Identity}
}

// Update runs one fixed 60 TPS tick. All game logic counts ticks, never
// wall-clock time, so speed is the same at any frame rate.
func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) ||
		(inpututil.IsKeyJustPressed(ebiten.KeyEnter) && ebiten.IsKeyPressed(ebiten.KeyAlt)) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	next := g.scene.Update(g.pollInput())
	if next == Quit {
		// Clean shutdown: RunGame returns nil and the program exits.
		return ebiten.Termination
	}
	if next != nil {
		g.scene = next
	}
	return nil
}

// Draw renders the current scene onto the 320×240 logical screen.
func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
	g.drawCrosshair(screen)
}

// drawCrosshair replaces the hidden system cursor, so the aim point always
// matches the game's pixel grid.
func (g *Game) drawCrosshair(screen *ebiten.Image) {
	x, y := g.cursor()
	x, y = math.Floor(x), math.Floor(y)
	c := color.RGBA{0xff, 0xf1, 0xe8, 0xff}
	for _, r := range [][4]float64{{-4, 0, 3, 1}, {2, 0, 3, 1}, {0, -4, 1, 3}, {0, 2, 1, 3}} {
		vector.FillRect(screen, float32(x+r[0]), float32(y+r[1]), float32(r[2]), float32(r[3]), c, false)
	}
}

// Layout always returns the fixed logical size. Ebitengine keeps it no matter
// how the window is resized; DrawFinalScreen handles scaling.
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return world.ScreenW, world.ScreenH
}

// DrawFinalScreen scales the logical screen by the largest whole number that
// fits, so pixels stay square and crisp at any window size or display scale.
// The rest is filled with black bars.
func (g *Game) DrawFinalScreen(screen ebiten.FinalScreen, offscreen *ebiten.Image, geoM ebiten.GeoM) {
	sw, sh := float64(screen.Bounds().Dx()), float64(screen.Bounds().Dy())
	t := scaling.PixelPerfect(sw, sh, world.ScreenW, world.ScreenH)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(t.Scale, t.Scale)
	op.GeoM.Translate(t.OffX, t.OffY)
	op.Filter = ebiten.FilterNearest
	if t.Scale < 1 {
		op.Filter = ebiten.FilterLinear // shrinking: smooth instead of dropping pixels
	}
	screen.Fill(color.Black)
	screen.DrawImage(offscreen, op)

	// Ebitengine's default transform is a plain scale + translate.
	g.defaultT = scaling.Transform{Scale: geoM.Element(0, 0), OffX: geoM.Element(0, 2), OffY: geoM.Element(1, 2)}
	g.drawnT = t
}

// pollInput is the only place that reads the keyboard and mouse.
func (g *Game) pollInput() input.State {
	key := ebiten.IsKeyPressed
	pressed := inpututil.IsKeyJustPressed
	mx, my := g.cursor()
	return input.State{
		Up:           key(ebiten.KeyW) || key(ebiten.KeyArrowUp),
		Down:         key(ebiten.KeyS) || key(ebiten.KeyArrowDown),
		Left:         key(ebiten.KeyA) || key(ebiten.KeyArrowLeft),
		Right:        key(ebiten.KeyD) || key(ebiten.KeyArrowRight),
		MouseX:       mx,
		MouseY:       my,
		Fire:         ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) || key(ebiten.KeySpace),
		FirePressed:  inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || inpututil.IsKeyJustPressed(ebiten.KeySpace),
		MenuUp:       pressed(ebiten.KeyW) || pressed(ebiten.KeyArrowUp),
		MenuDown:     pressed(ebiten.KeyS) || pressed(ebiten.KeyArrowDown),
		Confirm:      pressed(ebiten.KeyEnter) || pressed(ebiten.KeyNumpadEnter) || pressed(ebiten.KeySpace),
		Click:        inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft),
		PausePressed: inpututil.IsKeyJustPressed(ebiten.KeyP) || inpututil.IsKeyJustPressed(ebiten.KeyEscape),
		DebugPressed: inpututil.IsKeyJustPressed(ebiten.KeyF3),
		Focused:      ebiten.IsFocused(),
	}
}

// cursor returns the mouse position in logical pixels, clamped to the screen.
// Ebitengine reports it using its own default scaling, but we draw with
// whole-number scaling, so it's remapped (see scaling.RemapCursor).
func (g *Game) cursor() (float64, float64) {
	x, y := ebiten.CursorPositionF()
	return scaling.RemapCursor(x, y, g.defaultT, g.drawnT, world.ScreenW, world.ScreenH)
}
