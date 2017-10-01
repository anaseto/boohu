// combat utility functions

package main

func (g *game) HitDamage(base int, armor int) int {
	min := base / 2
	attack := min + RandInt(base-min+1)
	attack -= RandInt(armor + 1)
	if attack < 0 {
		attack = 0
	}
	return attack
}

func (m *monster) InflictDamage(g *game, damage, max int) {
	oldHP := g.Player.HP
	g.Player.HP -= damage
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical HP: %d (hit by %s)", g.Player.HP, Indefinite(m.Kind.String(), false))
		g.ui.CriticalHPWarning(g)
	}
}

func (g *game) MakeMonstersAware() {
	for _, m := range g.Monsters {
		if m.HP <= 0 {
			continue
		}
		if g.Player.LOS[m.Pos] {
			m.MakeAware(g)
			if m.State != Resting {
				m.GatherBand(g)
			}
		}
	}
}

func (g *game) MakeNoise(noise int, at position) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, noise)
	for _, m := range g.Monsters {
		if !m.Exists() {
			continue
		}
		if m.State == Hunting {
			continue
		}
		n, ok := nm[m.Pos]
		if !ok {
			continue
		}
		d := n.Cost
		v := noise - d
		if v <= 0 {
			continue
		}
		v *= 3
		if v > 90 {
			v = 90
		}
		r := RandInt(100)
		if m.State == Resting {
			r += 10
		}
		if v > r {
			m.Target = at
			if g.Player.LOS[m.Pos] {
				m.State = Hunting
			} else {
				m.State = Wandering
			}
			m.GatherBand(g)
		}
	}
}

func (g *game) AttackMonster(mons *monster) {
	switch {
	case g.Player.Weapon.Cleave():
		var neighbors []position
		if g.Player.HasStatus(StatusConfusion) {
			neighbors = g.Dungeon.CardinalFreeNeighbors(g.Player.Pos)
		} else {
			neighbors = g.Dungeon.FreeNeighbors(g.Player.Pos)
		}
		for _, pos := range neighbors {
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				g.HitMonster(mons)
			}
		}
	case g.Player.Weapon.Pierce():
		g.HitMonster(mons)
		deltaX := mons.Pos.X - g.Player.Pos.X
		deltaY := mons.Pos.Y - g.Player.Pos.Y
		behind := position{g.Player.Pos.X + 2*deltaX, g.Player.Pos.Y + 2*deltaY}
		if g.Dungeon.Valid(behind) {
			mons, _ := g.MonsterAt(behind)
			if mons.Exists() {
				g.HitMonster(mons)
			}
		}
	default:
		g.HitMonster(mons)
		if (g.Player.Weapon == Sword || g.Player.Weapon == DoubleSword) && RandInt(4) == 0 {
			g.HitMonster(mons)
		}
	}
}

func (g *game) HitMonster(mons *monster) {
	acc := RandInt(g.Player.Accuracy())
	ev := RandInt(mons.Evasion)
	if mons.State == Resting {
		ev /= 2 + 1
	}
	if acc > ev {
		g.MakeNoise(12, mons.Pos)
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += RandInt(5)
		}
		attack := g.HitDamage(g.Player.Attack()+bonus, mons.Armor)
		if mons.State == Resting {
			if g.Player.Weapon == Dagger {
				attack *= 4
			} else {
				attack *= 2
			}
		}
		oldHP := mons.HP
		mons.HP -= attack
		if mons.HP > 0 {
			g.Printf("You hit the %v (%d damage).", mons.Kind, attack)
		} else if oldHP > 0 {
			// test oldHP > 0 because of sword special attack
			g.Printf("You kill the %v (%d damage).", mons.Kind, attack)
			g.KillStats(mons)
		}
	} else {
		g.Printf("You miss the %v.", mons.Kind)
	}
	mons.MakeHuntIfHurt(g)
}
