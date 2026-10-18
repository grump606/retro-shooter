// Package level runs wave progression: it decides when and where enemies
// spawn into the World, and when the next wave starts.
package level

import (
	"math/rand/v2"
	"sort"

	"retro-shooter/internal/world"
)

const (
	intermissionTicks = 120 // pause + "WAVE N" banner between waves
	spawnMargin       = 12  // enemies appear this far outside the screen
	minSpawnDist      = 90  // never spawn closer than this to the player
)

type spawn struct {
	at   int // tick within the wave
	kind world.EnemyKind
}

// Manager drives waves for one run.
type Manager struct {
	wave         int
	elapsed      int // ticks since the current wave started spawning
	schedule     []spawn
	intermission int
	rng          *rand.Rand
}

// New returns a manager that starts with the wave 1 intermission.
func New(seed uint64) *Manager {
	m := &Manager{rng: rand.New(rand.NewPCG(seed, seed^0x9e3779b97f4a7c15))}
	m.startWave(1)
	return m
}

// Wave is the current wave number (1-based).
func (m *Manager) Wave() int { return m.wave }

// InIntermission reports whether the "get ready" pause before a wave is on.
func (m *Manager) InIntermission() bool { return m.intermission > 0 }

// Update spawns enemies that are due and advances to the next wave once the
// current one is fully spawned and cleared.
func (m *Manager) Update(w *world.World) {
	if m.intermission > 0 {
		m.intermission--
		return
	}
	for len(m.schedule) > 0 && m.schedule[0].at <= m.elapsed {
		// If the world is full, keep it queued and retry next tick.
		if !w.Spawn(m.schedule[0].kind, m.spawnPoint(w.Player.Pos)) {
			break
		}
		m.schedule = m.schedule[1:]
	}
	m.elapsed++

	if len(m.schedule) == 0 && len(w.Enemies) == 0 {
		m.startWave(m.wave + 1)
	}
}

func (m *Manager) startWave(n int) {
	m.wave = n
	m.elapsed = 0
	m.intermission = intermissionTicks
	m.schedule = m.schedule[:0]
	for _, g := range waveFor(n) {
		for i := range g.Count {
			m.schedule = append(m.schedule, spawn{at: g.Delay + i*g.Interval, kind: g.Kind})
		}
	}
	sort.SliceStable(m.schedule, func(i, j int) bool { return m.schedule[i].at < m.schedule[j].at })
}

// spawnPoint picks a point just outside a random screen edge, away from the
// player so nothing appears on top of them.
func (m *Manager) spawnPoint(player world.Vec) world.Vec {
	var p world.Vec
	for range 10 {
		p = m.edgePoint()
		if p.DistTo(player) >= minSpawnDist {
			return p
		}
	}
	// The player is hugging an edge; use the point farthest from them.
	return world.Vec{X: world.ScreenW - player.X, Y: world.ScreenH - player.Y}.
		Sub(world.Vec{X: world.ScreenW / 2, Y: world.ScreenH / 2}).
		Normalized().
		Scale(world.ScreenW).
		Add(world.Vec{X: world.ScreenW / 2, Y: world.ScreenH / 2})
}

func (m *Manager) edgePoint() world.Vec {
	const w, h, mg = world.ScreenW, world.ScreenH, spawnMargin
	switch m.rng.IntN(4) {
	case 0: // top
		return world.Vec{X: m.rng.Float64() * w, Y: -mg}
	case 1: // bottom
		return world.Vec{X: m.rng.Float64() * w, Y: h + mg}
	case 2: // left
		return world.Vec{X: -mg, Y: m.rng.Float64() * h}
	default: // right
		return world.Vec{X: w + mg, Y: m.rng.Float64() * h}
	}
}
