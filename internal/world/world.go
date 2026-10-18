// Package world holds the game simulation: player, enemies, bullets, score.
//
// It has no rendering or Ebitengine code, so it can be unit-tested with fake
// input snapshots. Scenes own a *World, feed it input each tick, and read its
// Events to trigger sounds, effects and level progression.
package world

import (
	"math"

	"retro-shooter/internal/input"
)

// Logical screen size. Everything is simulated and drawn at this resolution
// and scaled up by whole-number factors for a crisp pixel look.
const (
	ScreenW = 320
	ScreenH = 240
)

// Caps keep entity counts bounded no matter what happens.
const (
	maxBullets = 256
	maxEnemies = 64
)

// Vec is a 2D vector in logical pixels.
type Vec struct{ X, Y float64 }

func (v Vec) Add(o Vec) Vec        { return Vec{v.X + o.X, v.Y + o.Y} }
func (v Vec) Sub(o Vec) Vec        { return Vec{v.X - o.X, v.Y - o.Y} }
func (v Vec) Scale(s float64) Vec  { return Vec{v.X * s, v.Y * s} }
func (v Vec) Len() float64         { return math.Hypot(v.X, v.Y) }
func FromAngle(a float64) Vec      { return Vec{math.Cos(a), math.Sin(a)} }
func (v Vec) DistTo(o Vec) float64 { return v.Sub(o).Len() }

// Normalized returns a unit vector, or the zero vector if v is (near) zero.
func (v Vec) Normalized() Vec {
	l := v.Len()
	if l < 1e-9 {
		return Vec{}
	}
	return v.Scale(1 / l)
}

// EventKind identifies something that happened during a tick.
type EventKind int

const (
	EvShot EventKind = iota
	EvEnemyShot
	EvEnemyHit
	EvEnemyKilled
	EvPlayerHit
	EvPlayerDied
)

// Event is reported upward to the owning scene.
type Event struct {
	Kind  EventKind
	Pos   Vec
	Enemy EnemyKind // for enemy events
}

// World is all live game state for one run.
type World struct {
	Player  *Player
	Enemies []*Enemy
	Bullets []*Bullet
	Score   int
	Tick    int

	events []Event
}

// New creates a fresh world with the player in the center.
func New() *World {
	return &World{Player: newPlayer(Vec{ScreenW / 2, ScreenH / 2})}
}

// GameOver reports whether the player has died.
func (w *World) GameOver() bool { return w.Player.HP <= 0 }

// Update advances the simulation by one tick.
func (w *World) Update(in input.State) {
	if w.GameOver() {
		return
	}
	w.Player.update(in, w)
	for _, e := range w.Enemies {
		e.update(w)
	}
	for _, b := range w.Bullets {
		b.update()
	}
	w.resolveCollisions()
	w.removeDead()
	w.Tick++
}

// Spawn adds an enemy unless the enemy cap is reached.
func (w *World) Spawn(kind EnemyKind, pos Vec) bool {
	if len(w.Enemies) >= maxEnemies {
		return false
	}
	w.Enemies = append(w.Enemies, newEnemy(kind, pos))
	return true
}

// DrainEvents returns events from the last update(s) and clears the queue.
func (w *World) DrainEvents() []Event {
	ev := w.events
	w.events = nil
	return ev
}

func (w *World) emit(e Event) { w.events = append(w.events, e) }

func (w *World) addBullet(b *Bullet) bool {
	if len(w.Bullets) >= maxBullets {
		return false
	}
	w.Bullets = append(w.Bullets, b)
	return true
}

// removeDead filters dead entities after the update pass, so nothing is
// removed from a slice while it is being iterated.
func (w *World) removeDead() {
	enemies := w.Enemies[:0]
	for _, e := range w.Enemies {
		if !e.dead {
			enemies = append(enemies, e)
		}
	}
	clear(w.Enemies[len(enemies):])
	w.Enemies = enemies

	bullets := w.Bullets[:0]
	for _, b := range w.Bullets {
		if !b.dead {
			bullets = append(bullets, b)
		}
	}
	clear(w.Bullets[len(bullets):])
	w.Bullets = bullets
}
