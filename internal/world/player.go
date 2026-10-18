package world

import (
	"math"

	"retro-shooter/internal/input"
)

const (
	playerSpeed        = 1.6 // pixels per tick
	playerRadius       = 5
	playerMaxHP        = 5
	playerFireCooldown = 8  // ticks between shots
	playerInvulnTicks  = 60 // ticks of invulnerability after a hit
	playerBulletSpeed  = 4.5
	playerGunOffset    = 7 // bullet spawns this far in front of the ship
)

// Player is the ship the user controls.
type Player struct {
	Pos    Vec
	Angle  float64 // aim direction in radians; 0 points right
	HP     int
	MaxHP  int
	Radius float64
	Invuln int // ticks of invulnerability left

	cooldown int
}

func newPlayer(pos Vec) *Player {
	return &Player{Pos: pos, HP: playerMaxHP, MaxHP: playerMaxHP, Radius: playerRadius}
}

func (p *Player) update(in input.State, w *World) {
	// Movement, normalized so diagonals aren't faster.
	var dir Vec
	if in.Left {
		dir.X--
	}
	if in.Right {
		dir.X++
	}
	if in.Up {
		dir.Y--
	}
	if in.Down {
		dir.Y++
	}
	p.Pos = p.Pos.Add(dir.Normalized().Scale(playerSpeed))
	p.Pos.X = clamp(p.Pos.X, p.Radius, ScreenW-p.Radius)
	p.Pos.Y = clamp(p.Pos.Y, p.Radius, ScreenH-p.Radius)

	// Aim at the mouse. When the cursor is on top of the ship the direction
	// is undefined, so keep the previous angle instead of snapping to 0.
	aim := Vec{in.MouseX, in.MouseY}.Sub(p.Pos)
	if aim.Len() > 1 {
		p.Angle = math.Atan2(aim.Y, aim.X)
	}

	if p.cooldown > 0 {
		p.cooldown--
	}
	if p.Invuln > 0 {
		p.Invuln--
	}

	if in.Fire && p.cooldown == 0 {
		dir := FromAngle(p.Angle)
		b := &Bullet{
			Pos:        p.Pos.Add(dir.Scale(playerGunOffset)),
			Vel:        dir.Scale(playerBulletSpeed),
			Radius:     2,
			FromPlayer: true,
		}
		if w.addBullet(b) {
			p.cooldown = playerFireCooldown
			w.emit(Event{Kind: EvShot, Pos: b.Pos})
		}
	}
}

func clamp(v, lo, hi float64) float64 {
	return math.Max(lo, math.Min(hi, v))
}
