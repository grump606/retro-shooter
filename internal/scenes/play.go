package scenes

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/assets"
	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
	"retro-shooter/internal/level"
	"retro-shooter/internal/ui"
	"retro-shooter/internal/world"
)

const deathDelay = 90 // ticks to watch the explosion before game over

const (
	// maxShakePx caps how far (in whole pixels) the screen moves while shaking.
	maxShakePx = 3

	// debugY puts the F3 overlay's 4 lines one line above the bottom edge.
	debugY = world.ScreenH - 5*ui.LineH

	pauseDimAlpha = 0xa0 // how dark the pause overlay is (0–255)
)

// Play owns everything for one run: the World, the wave LevelManager and
// visual effects. A new Play (and so a fresh World) is made for every run.
type Play struct {
	d      *Deps
	w      *world.World
	levels *level.Manager
	fx     *Particles
	rng    *rand.Rand    // visual randomness (screen shake)
	canvas *ebiten.Image // world is drawn here first so it can be shaken

	paused      bool
	waitRelease bool // ignore fire until the button is released once
	shake       int
	deathTimer  int
	lastWave    int
}

func NewPlay(d *Deps) *Play {
	seed := uint64(time.Now().UnixNano())
	return &Play{
		d:      d,
		w:      world.New(),
		levels: level.New(seed),
		fx:     NewParticles(seed),
		rng:    rand.New(rand.NewPCG(seed, seed+13)),
		canvas: d.Canvas(),
		// The click that started the game shouldn't also fire a shot.
		waitRelease: true,
		// The title screen already blipped for wave 1; only blip for later waves.
		lastWave: 1,
	}
}

func (s *Play) Update(in input.State) game.Scene {
	s.d.Audio.NextTick()
	if in.DebugPressed {
		s.d.Debug = !s.d.Debug
	}

	// Pause: P/Esc toggles; losing window focus pauses automatically.
	if !in.Focused && !s.paused {
		s.paused = true
	}
	if s.paused {
		if in.PausePressed || in.FirePressed {
			s.paused = false
			s.waitRelease = true
			s.d.Audio.Play(audio.Blip)
		}
		return nil
	}
	if in.PausePressed && !s.w.GameOver() {
		s.paused = true
		s.d.Audio.Play(audio.Blip)
		return nil
	}

	if s.waitRelease {
		if in.Fire {
			in.Fire = false
		} else {
			s.waitRelease = false
		}
	}

	s.d.Stars.Update()
	s.w.Update(in)
	// Waves stop once the player is dead; otherwise new enemies would pop in
	// frozen during the death animation and the next wave could start.
	if !s.w.GameOver() {
		s.levels.Update(s.w)
		if s.levels.Wave() != s.lastWave {
			s.lastWave = s.levels.Wave()
			s.d.Audio.Play(audio.Blip)
		}
	}
	s.handleEvents()
	s.fx.Update()
	if s.shake > 0 {
		s.shake--
	}

	if s.w.GameOver() {
		s.deathTimer++
		if s.deathTimer >= deathDelay {
			newHigh := s.w.Score > s.d.HighScore
			if newHigh {
				s.d.HighScore = s.w.Score
			}
			return NewGameOver(s.d, s.w.Score, s.levels.Wave(), newHigh)
		}
	}
	return nil
}

// handleEvents turns world events into sound, particles and screen shake.
func (s *Play) handleEvents() {
	for _, ev := range s.w.DrainEvents() {
		switch ev.Kind {
		case world.EvShot:
			s.d.Audio.Play(audio.Shoot)
		case world.EvEnemyShot:
			s.d.Audio.Play(audio.EnemyShoot)
		case world.EvEnemyHit:
			s.d.Audio.Play(audio.EnemyHit)
			s.fx.Burst(ev.Pos, 4, 1.5, colorWhite)
		case world.EvEnemyKilled:
			s.d.Audio.Play(audio.Explode)
			d := ev.Enemy.Def()
			s.fx.Burst(ev.Pos, d.DebrisCount, 2, assets.Color(d.Debris))
			s.fx.Burst(ev.Pos, d.DebrisCount/2, 1, colorOrange)
			s.shake = max(s.shake, d.Shake)
		case world.EvPlayerHit:
			s.d.Audio.Play(audio.PlayerHit)
			s.fx.Burst(ev.Pos, 12, 2, colorRed)
			s.shake = max(s.shake, 8)
		case world.EvPlayerDied:
			s.d.Audio.Play(audio.PlayerDie)
			s.fx.Burst(ev.Pos, 80, 3, colorCyan)
			s.fx.Burst(ev.Pos, 60, 1.5, colorOrange)
			s.shake = 25
		}
	}
}

func (s *Play) Draw(screen *ebiten.Image) {
	a := s.d.Assets
	c := s.canvas
	c.Fill(ui.Night)
	s.d.Stars.Draw(c, a.Pixel)
	drawWorld(c, a, s.w)
	s.fx.Draw(c, a.Pixel)

	// Screen shake: draw the canvas at a small random whole-pixel offset, so
	// the pixel art never lands between pixels and blurs or shimmers.
	op := &ebiten.DrawImageOptions{}
	if !s.paused {
		dx, dy := shakeOffset(s.shake, s.rng)
		op.GeoM.Translate(float64(dx), float64(dy))
	}
	screen.Fill(ui.Night)
	screen.DrawImage(c, op)

	ui.DrawHUD(screen, s.w, s.levels.Wave())

	if s.levels.InIntermission() && !s.w.GameOver() {
		ui.BigText(screen, fmt.Sprintf("WAVE %d", s.levels.Wave()), ui.BannerY, ui.BannerScale, ui.Yellow)
		ui.TextCentered(screen, "GET READY", ui.BannerSubY)
	}
	if s.paused {
		ui.Rect(screen, 0, 0, world.ScreenW, world.ScreenH, color.RGBA{0, 0, 0, pauseDimAlpha})
		ui.BigText(screen, "PAUSED", ui.BannerY, ui.BannerScale, ui.White)
		ui.TextCentered(screen, "P OR CLICK TO RESUME", ui.BannerSubY)
	}
	if s.d.Debug {
		ui.Text(screen, fmt.Sprintf("TPS %.0f FPS %.0f\nENEMIES %d\nBULLETS %d\nFX %d",
			ebiten.ActualTPS(), ebiten.ActualFPS(), len(s.w.Enemies), len(s.w.Bullets), s.fx.Len()), ui.Margin, debugY)
	}
}

// Effect colors come from the shared sprite palette.
var (
	colorWhite  = assets.Color('w')
	colorOrange = assets.Color('o')
	colorRed    = assets.Color('r')
	colorCyan   = assets.Color('c')
)

// shakeOffset returns a random whole-pixel offset for screen shake. The
// strength fades as the remaining shake ticks run out: up to maxShakePx at
// first, down to 1 px at the end, and none once ticks reaches 0.
func shakeOffset(ticks int, r *rand.Rand) (dx, dy int) {
	if ticks <= 0 {
		return 0, 0
	}
	mag := max(1, min(ticks/2, maxShakePx))
	return r.IntN(2*mag+1) - mag, r.IntN(2*mag+1) - mag
}
