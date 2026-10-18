package world

// Bullet is a projectile fired by the player or an enemy.
type Bullet struct {
	Pos        Vec
	Vel        Vec
	Radius     float64
	FromPlayer bool

	dead bool
}

// bulletMargin lets bullets fly slightly off-screen before despawning.
const bulletMargin = 16

func (b *Bullet) update() {
	b.Pos = b.Pos.Add(b.Vel)
	if !onScreen(b.Pos, bulletMargin) {
		b.dead = true
	}
}

// onScreen reports whether p is inside the screen expanded by margin.
func onScreen(p Vec, margin float64) bool {
	return p.X >= -margin && p.X <= ScreenW+margin &&
		p.Y >= -margin && p.Y <= ScreenH+margin
}
