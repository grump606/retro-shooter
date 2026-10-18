package world

// EnemyKind selects an enemy's stats and sprite.
type EnemyKind int

const (
	Grunt EnemyKind = iota // slow, weak, comes in numbers
	Fast                   // quick and fragile
	Tank                   // slow, tough, shoots back
)

type enemyStats struct {
	hp        int
	speed     float64
	radius    float64
	score     int
	fireEvery int // ticks between shots; 0 = never shoots
}

var stats = map[EnemyKind]enemyStats{
	Grunt: {hp: 1, speed: 0.6, radius: 5, score: 100},
	Fast:  {hp: 1, speed: 1.25, radius: 4, score: 150},
	Tank:  {hp: 6, speed: 0.35, radius: 8, score: 500, fireEvery: 120},
}

const enemyBulletSpeed = 1.6

// Enemy chases the player.
type Enemy struct {
	Kind   EnemyKind
	Pos    Vec
	HP     int
	Radius float64
	Flash  int // ticks to draw white after being hit

	speed     float64
	fireTimer int
	dead      bool
}

func newEnemy(kind EnemyKind, pos Vec) *Enemy {
	s := stats[kind]
	return &Enemy{
		Kind:      kind,
		Pos:       pos,
		HP:        s.hp,
		Radius:    s.radius,
		speed:     s.speed,
		fireTimer: s.fireEvery,
	}
}

// ScoreValue is the points awarded for destroying this kind of enemy.
func ScoreValue(kind EnemyKind) int { return stats[kind].score }

func (e *Enemy) update(w *World) {
	if e.Flash > 0 {
		e.Flash--
	}
	toPlayer := w.Player.Pos.Sub(e.Pos)
	e.Pos = e.Pos.Add(toPlayer.Normalized().Scale(e.speed))

	if every := stats[e.Kind].fireEvery; every > 0 {
		e.fireTimer--
		// Only shoot once on screen, so tanks can't snipe from off-screen.
		if e.fireTimer <= 0 && onScreen(e.Pos, 0) {
			e.fireTimer = every
			b := &Bullet{
				Pos:    e.Pos,
				Vel:    toPlayer.Normalized().Scale(enemyBulletSpeed),
				Radius: 2,
			}
			if w.addBullet(b) {
				w.emit(Event{Kind: EvEnemyShot, Pos: e.Pos, Enemy: e.Kind})
			}
		}
	}
}
