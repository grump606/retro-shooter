package scenes

import (
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"

	"retro-shooter/internal/world"
)

// drawPixel draws one tinted pixel using the shared 1×1 image. Drawing the
// same source image many times lets Ebitengine batch it into one draw call.
//
// It takes color.RGBA, not color.Color: converting a color to the interface
// allocates, and this runs for every star and particle, every frame.
func drawPixel(dst, px *ebiten.Image, x, y float64, c color.RGBA) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(math.Floor(x), math.Floor(y))
	// Same as ScaleWithColor(c): color.RGBA is already alpha-premultiplied.
	op.ColorScale.Scale(float32(c.R)/0xff, float32(c.G)/0xff, float32(c.B)/0xff, float32(c.A)/0xff)
	dst.DrawImage(px, op)
}

// --- Starfield ---------------------------------------------------------------

type star struct {
	x, y, speed float64
	c           color.RGBA
}

// Starfield is a slowly scrolling 3-layer parallax background.
type Starfield struct{ stars []star }

func NewStarfield(seed uint64) *Starfield {
	r := rand.New(rand.NewPCG(seed, seed+1))
	layers := []struct {
		speed float64
		c     color.RGBA
		n     int
	}{
		{0.08, color.RGBA{0x2a, 0x2f, 0x4f, 0xff}, 40},
		{0.20, color.RGBA{0x5f, 0x57, 0x4f, 0xff}, 25},
		{0.45, color.RGBA{0xc2, 0xc3, 0xc7, 0xff}, 12},
	}
	s := &Starfield{}
	for _, l := range layers {
		for range l.n {
			s.stars = append(s.stars, star{
				x: r.Float64() * world.ScreenW, y: r.Float64() * world.ScreenH,
				speed: l.speed, c: l.c,
			})
		}
	}
	return s
}

func (s *Starfield) Update() {
	for i := range s.stars {
		s.stars[i].y += s.stars[i].speed
		if s.stars[i].y >= world.ScreenH {
			s.stars[i].y -= world.ScreenH
		}
	}
}

func (s *Starfield) Draw(dst, px *ebiten.Image) {
	for _, st := range s.stars {
		drawPixel(dst, px, st.x, st.y, st.c)
	}
}

// --- Particles ---------------------------------------------------------------

const maxParticles = 500

type particle struct {
	pos, vel  world.Vec
	life, max int
	c         color.RGBA
}

// Particles are purely visual (explosions, hit sparks). They live in the
// scene, not the World, because they don't affect gameplay.
type Particles struct {
	ps  []particle
	rng *rand.Rand
}

func NewParticles(seed uint64) *Particles {
	return &Particles{rng: rand.New(rand.NewPCG(seed, seed+7))}
}

// Burst spawns n particles flying out from pos.
func (p *Particles) Burst(pos world.Vec, n int, speed float64, c color.RGBA) {
	for range n {
		if len(p.ps) >= maxParticles {
			return
		}
		life := 15 + p.rng.IntN(20)
		v := world.FromAngle(p.rng.Float64() * 2 * math.Pi).Scale(speed * (0.3 + p.rng.Float64()))
		p.ps = append(p.ps, particle{pos: pos, vel: v, life: life, max: life, c: c})
	}
}

func (p *Particles) Update() {
	alive := p.ps[:0]
	for _, q := range p.ps {
		q.pos = q.pos.Add(q.vel)
		q.vel = q.vel.Scale(0.92)
		q.life--
		if q.life > 0 {
			alive = append(alive, q)
		}
	}
	p.ps = alive
}

func (p *Particles) Draw(dst, px *ebiten.Image) {
	for _, q := range p.ps {
		c := q.c
		if q.life < q.max/3 {
			c.R, c.G, c.B = c.R/2, c.G/2, c.B/2 // fade out
		}
		drawPixel(dst, px, q.pos.X, q.pos.Y, c)
	}
}

func (p *Particles) Len() int { return len(p.ps) }
