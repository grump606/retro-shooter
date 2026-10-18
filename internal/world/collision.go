package world

func overlaps(a Vec, ra float64, b Vec, rb float64) bool {
	d := a.Sub(b)
	r := ra + rb
	return d.X*d.X+d.Y*d.Y <= r*r
}

func (w *World) resolveCollisions() {
	// Player bullets vs enemies. Each bullet hits at most one enemy.
	for _, b := range w.Bullets {
		if b.dead || !b.FromPlayer {
			continue
		}
		for _, e := range w.Enemies {
			if e.dead || !overlaps(b.Pos, b.Radius, e.Pos, e.Radius) {
				continue
			}
			b.dead = true
			e.HP--
			e.Flash = 4
			if e.HP <= 0 {
				e.dead = true
				w.Score += ScoreValue(e.Kind)
				w.emit(Event{Kind: EvEnemyKilled, Pos: e.Pos, Enemy: e.Kind})
			} else {
				w.emit(Event{Kind: EvEnemyHit, Pos: e.Pos, Enemy: e.Kind})
			}
			break
		}
	}

	// Enemies and enemy bullets vs player. However many things touch the
	// player this tick, they lose at most one HP, then get invulnerability.
	p := w.Player
	if p.Invuln > 0 || p.HP <= 0 {
		return
	}
	hit := false
	for _, e := range w.Enemies {
		if !e.dead && overlaps(e.Pos, e.Radius, p.Pos, p.Radius) {
			hit = true
			// Small enemies are destroyed by ramming (no score); tanks survive.
			if e.Kind != Tank {
				e.dead = true
				w.emit(Event{Kind: EvEnemyKilled, Pos: e.Pos, Enemy: e.Kind})
			}
		}
	}
	for _, b := range w.Bullets {
		if !b.dead && !b.FromPlayer && overlaps(b.Pos, b.Radius, p.Pos, p.Radius) {
			hit = true
			b.dead = true
		}
	}
	if !hit {
		return
	}
	p.HP--
	p.Invuln = playerInvulnTicks
	w.emit(Event{Kind: EvPlayerHit, Pos: p.Pos})
	if p.HP <= 0 {
		w.emit(Event{Kind: EvPlayerDied, Pos: p.Pos})
	}
}
