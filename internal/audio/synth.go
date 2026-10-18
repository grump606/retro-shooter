package audio

import (
	"encoding/binary"
	"math"
	"math/rand/v2"
)

// Sounds are generated in code (chiptune style) instead of loaded from files,
// so the game has no audio asset dependencies. Swap in .wav files from jsfxr
// later by decoding them with ebiten's audio/wav package.

func synthesizeAll() [numSounds][]byte {
	var d [numSounds][]byte
	d[Shoot] = synth(0.09, square(sweep(880, 260)), 0.16)
	d[EnemyShoot] = synth(0.10, triangle(sweep(520, 380)), 0.20)
	d[EnemyHit] = synth(0.05, square(sweep(240, 180)), 0.14)
	d[Explode] = synth(0.35, noise(0.25), 0.40)
	d[PlayerHit] = synth(0.25, square(sweep(320, 70)), 0.22)
	d[PlayerDie] = synth(0.9, noise(0.08), 0.45)
	d[Blip] = synth(0.10, square(sweep(660, 660)), 0.12)
	return d
}

// wave returns a sample in [-1,1] given the progress p in [0,1].
type wave func(p float64) float64

// synth renders dur seconds of a wave with a decaying envelope into 16-bit
// little-endian stereo PCM, the format Ebitengine's audio players expect.
func synth(dur float64, w wave, amp float64) []byte {
	n := int(dur * SampleRate)
	buf := make([]byte, n*4)
	for i := range n {
		p := float64(i) / float64(n)
		env := (1 - p) * (1 - p)
		v := int16(math.Max(-1, math.Min(1, w(p)*env*amp)) * math.MaxInt16)
		binary.LittleEndian.PutUint16(buf[i*4:], uint16(v))
		binary.LittleEndian.PutUint16(buf[i*4+2:], uint16(v))
	}
	return buf
}

// sweep returns a frequency that slides linearly from f0 to f1.
func sweep(f0, f1 float64) func(p float64) float64 {
	return func(p float64) float64 { return f0 + (f1-f0)*p }
}

// square and triangle accumulate phase sample-by-sample so frequency sweeps
// stay smooth. They are stateful and must be used for one render only.
func square(freq func(float64) float64) wave {
	phase := 0.0
	return func(p float64) float64 {
		phase += freq(p) / SampleRate
		if math.Mod(phase, 1) < 0.5 {
			return 1
		}
		return -1
	}
}

func triangle(freq func(float64) float64) wave {
	phase := 0.0
	return func(p float64) float64 {
		phase += freq(p) / SampleRate
		return 4*math.Abs(math.Mod(phase, 1)-0.5) - 1
	}
}

// noise is white noise through a one-pole low-pass filter; smaller cutoff
// (0..1) gives a deeper rumble.
func noise(cutoff float64) wave {
	r := rand.New(rand.NewPCG(1, 2))
	prev := 0.0
	return func(float64) float64 {
		prev += cutoff * (r.Float64()*2 - 1 - prev)
		return prev * 2.5
	}
}
