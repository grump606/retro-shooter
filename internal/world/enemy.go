package world

// EnemyKind selects an enemy's definition in EnemyDefs.
type EnemyKind int

const (
	Grunt EnemyKind = iota // slow, weak, comes in numbers
	Fast                   // quick and fragile
	Tank                   // slow, tough, shoots back

	NumEnemyKinds // keep last: number of enemy kinds
)

// EnemyDef is everything about one kind of enemy: gameplay stats, how it
// looks, and how it explodes. To add an enemy, add a constant above and an
// entry in EnemyDefs; nothing else needs to change (tests check the table).
//
// Visual fields are plain data (palette letters, text grids) so this package
// stays free of rendering code. Palette letters are defined in assets.Palette.
type EnemyDef struct {
	Name string

	// Gameplay
	HP          int
	Speed       float64 // pixels per tick
	Radius      float64 // collision radius in pixels
	Score       int     // points for shooting it down
	FireEvery   int     // ticks between shots; 0 = never shoots
	SurvivesRam bool    // true: touching the player doesn't destroy it

	// Look
	Sprite []string // pixel art rows, one palette letter per pixel, '.' = clear

	// Death effect
	Debris      byte // palette letter for explosion particles
	DebrisCount int  // number of particles
	Shake       int  // screen-shake ticks (0 = none)
}

// EnemyDefs is the single source of truth for every enemy kind.
var EnemyDefs = [NumEnemyKinds]EnemyDef{
	Grunt: {
		Name: "grunt",
		HP:   1, Speed: 0.6, Radius: 5, Score: 100,
		Sprite: []string{
			"..kkkkk..",
			".kgggggk.",
			"kgwkgwkgk",
			"kgggggggk",
			"kgggggggk",
			".kgkgkgk.",
			"..k.k.k..",
		},
		Debris: 'g', DebrisCount: 14,
	},
	Fast: {
		Name: "fast",
		HP:   1, Speed: 1.25, Radius: 4, Score: 150,
		Sprite: []string{
			"...k...",
			"..kyk..",
			".kyoyk.",
			"kyoroyk",
			".kyoyk.",
			"..kyk..",
			"...k...",
		},
		Debris: 'y', DebrisCount: 14,
	},
	Tank: {
		Name: "tank",
		HP:   6, Speed: 0.35, Radius: 8, Score: 500,
		FireEvery: 120, SurvivesRam: true,
		Sprite: []string{
			"....kkkkkkk....",
			"..kkpppppppkk..",
			".kppmmmmmmmppk.",
			".kpmmpppppmmpk.",
			"kpmmpkkkkkpmmpk",
			"kpmpkrrrrrkpmpk",
			"kpmpkrwwwrkpmpk",
			"kpmpkrwwwrkpmpk",
			"kpmpkrrrrrkpmpk",
			"kpmmpkkkkkpmmpk",
			".kpmmpppppmmpk.",
			".kppmmmmmmmppk.",
			"..kkpppppppkk..",
			"....kkkkkkk....",
		},
		Debris: 'm', DebrisCount: 40, Shake: 10,
	},
}

// Def returns the definition for an enemy kind.
func (k EnemyKind) Def() *EnemyDef { return &EnemyDefs[k] }

// ScoreValue is the points awarded for destroying this kind of enemy.
func ScoreValue(kind EnemyKind) int { return kind.Def().Score }

const enemyBulletSpeed = 1.6

// Enemy chases the player.
type Enemy struct {
	Kind   EnemyKind
	Pos    Vec
	HP     int
	Radius float64
	Flash  int // ticks to draw white after being hit

	fireTimer int
	dead      bool
}

func newEnemy(kind EnemyKind, pos Vec) *Enemy {
	d := kind.Def()
	return &Enemy{
		Kind:      kind,
		Pos:       pos,
		HP:        d.HP,
		Radius:    d.Radius,
		fireTimer: d.FireEvery,
	}
}

func (e *Enemy) update(w *World) {
	d := e.Kind.Def()
	if e.Flash > 0 {
		e.Flash--
	}
	toPlayer := w.Player.Pos.Sub(e.Pos)
	e.Pos = e.Pos.Add(toPlayer.Normalized().Scale(d.Speed))

	if d.FireEvery > 0 {
		e.fireTimer--
		// Only shoot once on screen, so enemies can't snipe from off-screen.
		if e.fireTimer <= 0 && onScreen(e.Pos, 0) {
			e.fireTimer = d.FireEvery
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
