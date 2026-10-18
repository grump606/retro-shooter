package level

import "retro-shooter/internal/world"

// Group spawns Count enemies of one kind, one every Interval ticks.
type Group struct {
	Kind     world.EnemyKind
	Count    int
	Interval int
	Delay    int // ticks after wave start before the first spawn
}

// Wave is a set of groups that spawn in parallel.
type Wave []Group

// hand-made opening waves; after these, waves are generated (see waveFor).
var waves = []Wave{
	{{Kind: world.Grunt, Count: 6, Interval: 45}},
	{{Kind: world.Grunt, Count: 8, Interval: 35}, {Kind: world.Fast, Count: 3, Interval: 60, Delay: 120}},
	{{Kind: world.Grunt, Count: 10, Interval: 30}, {Kind: world.Fast, Count: 6, Interval: 40, Delay: 60}},
	{{Kind: world.Grunt, Count: 8, Interval: 30}, {Kind: world.Tank, Count: 1, Interval: 1, Delay: 90}},
	{{Kind: world.Grunt, Count: 12, Interval: 25}, {Kind: world.Fast, Count: 8, Interval: 30, Delay: 60}, {Kind: world.Tank, Count: 2, Interval: 240, Delay: 120}},
}

// waveFor returns the definition for wave n (1-based). Past the hand-made
// waves it keeps scaling up, so the game is endless.
func waveFor(n int) Wave {
	if n <= len(waves) {
		return waves[n-1]
	}
	k := n - len(waves)
	interval := max(12, 25-k)
	return Wave{
		{Kind: world.Grunt, Count: 12 + 2*k, Interval: interval},
		{Kind: world.Fast, Count: 8 + k, Interval: interval + 5, Delay: 60},
		{Kind: world.Tank, Count: 2 + k/2, Interval: 200, Delay: 120},
	}
}
