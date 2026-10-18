// Package input defines a per-tick snapshot of player input.
//
// It deliberately does not import Ebitengine: the game loop fills in a State
// once per tick (see game.pollInput), and everything downstream only sees this
// plain struct. That keeps the simulation testable with fake input.
package input

// State is everything the game needs to know about input for one tick.
type State struct {
	Up, Down, Left, Right bool

	// MouseX/MouseY are in logical screen pixels (0..ScreenW, 0..ScreenH).
	MouseX, MouseY float64

	Fire        bool // fire button held (gameplay)
	FirePressed bool // fire button went down this tick (menus)

	PausePressed bool // P or Esc went down this tick
	DebugPressed bool // F3 went down this tick

	Focused bool // window has focus
}
