package world

import (
	"math"
	"testing"

	"retro-shooter/internal/input"
)

func TestDiagonalMovementIsNotFaster(t *testing.T) {
	straight := New()
	diag := New()
	start := straight.Player.Pos

	straight.Update(input.State{Right: true})
	diag.Update(input.State{Right: true, Down: true})

	ds := straight.Player.Pos.DistTo(start)
	dd := diag.Player.Pos.DistTo(start)
	if math.Abs(ds-dd) > 1e-9 {
		t.Fatalf("straight moved %v, diagonal moved %v", ds, dd)
	}
}

func TestPlayerStaysOnScreen(t *testing.T) {
	w := New()
	for range 1000 {
		w.Update(input.State{Left: true, Up: true})
	}
	p := w.Player
	if p.Pos.X < p.Radius || p.Pos.Y < p.Radius {
		t.Fatalf("player left the screen: %+v", p.Pos)
	}
}

func TestAimKeepsAngleWhenCursorOnShip(t *testing.T) {
	w := New()
	p := w.Player
	w.Update(input.State{MouseX: p.Pos.X, MouseY: p.Pos.Y - 50}) // aim up
	want := p.Angle
	w.Update(input.State{MouseX: p.Pos.X, MouseY: p.Pos.Y}) // cursor on ship
	if p.Angle != want {
		t.Fatalf("angle changed to %v, want %v", p.Angle, want)
	}
}

func TestFireCooldownAndCleanup(t *testing.T) {
	w := New()
	fire := input.State{Fire: true, MouseX: ScreenW, MouseY: ScreenH / 2}
	for range playerFireCooldown*3 + 1 {
		w.Update(fire)
	}
	if got := len(w.Bullets); got != 4 {
		t.Fatalf("got %d bullets, want 4 (one per cooldown)", got)
	}
	// Bullets fly off-screen and are removed.
	for range 200 {
		w.Update(input.State{})
	}
	if got := len(w.Bullets); got != 0 {
		t.Fatalf("got %d bullets after they left the screen, want 0", got)
	}
}

func TestBulletKillsEnemyAndScores(t *testing.T) {
	w := New()
	p := w.Player
	w.Spawn(Grunt, Vec{p.Pos.X + 40, p.Pos.Y})
	fire := input.State{Fire: true, MouseX: p.Pos.X + 100, MouseY: p.Pos.Y}
	for range 30 {
		w.Update(fire)
		if len(w.Enemies) == 0 {
			break
		}
	}
	if len(w.Enemies) != 0 {
		t.Fatal("enemy was not destroyed")
	}
	if w.Score != ScoreValue(Grunt) {
		t.Fatalf("score = %d, want %d", w.Score, ScoreValue(Grunt))
	}
}

func TestOnlyOneHitPerTick(t *testing.T) {
	w := New()
	p := w.Player
	// Three enemies touching the player at once.
	for range 3 {
		w.Spawn(Grunt, p.Pos)
	}
	w.Update(input.State{})
	if p.HP != p.MaxHP-1 {
		t.Fatalf("HP = %d, want %d", p.HP, p.MaxHP-1)
	}
	if p.Invuln == 0 {
		t.Fatal("expected invulnerability after a hit")
	}
}

func TestDeathEndsGame(t *testing.T) {
	w := New()
	w.Player.HP = 1
	w.Spawn(Grunt, w.Player.Pos)
	w.Update(input.State{})
	if !w.GameOver() {
		t.Fatal("expected game over")
	}
	var died bool
	for _, e := range w.DrainEvents() {
		died = died || e.Kind == EvPlayerDied
	}
	if !died {
		t.Fatal("expected EvPlayerDied event")
	}
	// The world is frozen after death.
	tick := w.Tick
	w.Update(input.State{Right: true})
	if w.Tick != tick {
		t.Fatal("world kept updating after game over")
	}
}

func TestEntityCaps(t *testing.T) {
	w := New()
	for range maxEnemies + 10 {
		w.Spawn(Grunt, Vec{-100, -100})
	}
	if len(w.Enemies) != maxEnemies {
		t.Fatalf("got %d enemies, want cap %d", len(w.Enemies), maxEnemies)
	}
}
