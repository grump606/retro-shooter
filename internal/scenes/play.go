package scenes

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/audio"
	"retro-shooter/internal/game"
	"retro-shooter/internal/input"
	"retro-shooter/internal/level"
	"retro-shooter/internal/ui"
	"retro-shooter/internal/world"
)

const deathDelay = 90 // ticks to watch the explosion before game over

// Play owns everything for one run: the World, the wave LevelManager and
// visual effects. A new Play (and so a fresh World) is made for every run.
type Play struct {
	d      *Deps
	w      *world.World
	levels *level.Manager
	fx     *Particles
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
		canvas: ebiten.NewImage(world.ScreenW, world.ScreenH),
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
			n, c := 14, colorGreen
			switch ev.Enemy {
			case world.Fast:
				c = colorYellow
			case world.Tank:
				n, c = 40, colorPink
				s.shake = max(s.shake, 10)
			}
			s.fx.Burst(ev.Pos, n, 2, c)
			s.fx.Burst(ev.Pos, n/2, 1, colorOrange)
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

	// Screen shake: draw the canvas at a small random offset.
	op := &ebiten.DrawImageOptions{}
	if s.shake > 0 && !s.paused {
		mag := float64(min(s.shake, 6)) / 2
		op.GeoM.Translate(float64(rand.IntN(3)-1)*mag, float64(rand.IntN(3)-1)*mag)
	}
	screen.Fill(ui.Night)
	screen.DrawImage(c, op)

	ui.DrawHUD(screen, s.w, s.levels.Wave())

	if s.levels.InIntermission() && !s.w.GameOver() {
		ui.BigText(screen, fmt.Sprintf("WAVE %d", s.levels.Wave()), 90, 3, ui.Yellow)
		ui.TextCentered(screen, "GET READY", 140)
	}
	if s.paused {
		ui.Rect(screen, 0, 0, world.ScreenW, world.ScreenH, color.RGBA{0, 0, 0, 0xa0})
		ui.BigText(screen, "PAUSED", 90, 3, ui.White)
		ui.TextCentered(screen, "P OR CLICK TO RESUME", 140)
	}
	if s.d.Debug {
		ui.Text(screen, fmt.Sprintf("TPS %.0f FPS %.0f\nENEMIES %d\nBULLETS %d\nFX %d",
			ebiten.ActualTPS(), ebiten.ActualFPS(), len(s.w.Enemies), len(s.w.Bullets), s.fx.Len()), 4, 160)
	}
}

var (
	colorWhite  = color.RGBA{0xff, 0xf1, 0xe8, 0xff}
	colorGreen  = color.RGBA{0x00, 0xe4, 0x36, 0xff}
	colorYellow = color.RGBA{0xff, 0xec, 0x27, 0xff}
	colorOrange = color.RGBA{0xff, 0xa3, 0x00, 0xff}
	colorPink   = color.RGBA{0xff, 0x77, 0xa8, 0xff}
	colorRed    = color.RGBA{0xff, 0x00, 0x4d, 0xff}
	colorCyan   = color.RGBA{0x29, 0xad, 0xff, 0xff}
)
