// Package audio owns the single Ebitengine audio context and plays the
// game's sound effects.
package audio

import (
	eaudio "github.com/hajimehoshi/ebiten/v2/audio"
)

// SampleRate for the audio context and all generated sounds.
const SampleRate = 44100

// Sound identifies a sound effect.
type Sound int

const (
	Shoot Sound = iota
	EnemyShoot
	EnemyHit
	Explode
	PlayerHit
	PlayerDie
	Blip // menus and wave start
	numSounds
)

// Player plays sound effects. A muted Player never touches the audio device.
type Player struct {
	ctx    *eaudio.Context
	data   [numSounds][]byte
	last   [numSounds]int
	tick   int
	volume float64
}

// New creates the audio context (exactly once; Ebitengine panics on a second
// context) unless muted is true. Sounds are synthesized up front.
func New(muted bool) *Player {
	p := &Player{volume: 0.6}
	for i := range p.last {
		p.last[i] = -1
	}
	if muted {
		return p
	}
	p.ctx = eaudio.NewContext(SampleRate)
	p.data = synthesizeAll()
	return p
}

// Play starts a sound. Each sound plays at most once per tick, so ten enemies
// dying on the same frame don't produce a deafening spike.
func (p *Player) Play(s Sound) {
	if p.ctx == nil || p.last[s] == p.tick {
		return
	}
	p.last[s] = p.tick
	pl := p.ctx.NewPlayerFromBytes(p.data[s])
	pl.SetVolume(p.volume)
	pl.Play()
}

// NextTick must be called once per game tick (rate-limiting window).
func (p *Player) NextTick() { p.tick++ }
