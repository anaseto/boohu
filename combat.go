// combat utility functions

package main

import "container/heap"

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
			if g.Player.LOS[m.Pos] {
				m.MakeHunt(g)
			} else {
				m.Target = at
				m.State = Wandering
			}
			m.GatherBand(g)
		}
	}
}

func (g *game) AttackMonster(mons *monster, ev event) {
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
				g.HitMonster(mons, ev)
			}
		}
	case g.Player.Weapon.Pierce():
		g.HitMonster(mons, ev)
		deltaX := mons.Pos.X - g.Player.Pos.X
		deltaY := mons.Pos.Y - g.Player.Pos.Y
		behind := position{g.Player.Pos.X + 2*deltaX, g.Player.Pos.Y + 2*deltaY}
		if g.Dungeon.Valid(behind) {
			mons, _ := g.MonsterAt(behind)
			if mons.Exists() {
				g.HitMonster(mons, ev)
			}
		}
	default:
		g.HitMonster(mons, ev)
		if (g.Player.Weapon == Sword || g.Player.Weapon == DoubleSword) && RandInt(4) == 0 {
			g.HitMonster(mons, ev)
		}
	}
}

func (g *game) HitMonster(mons *monster, ev event) {
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		g.MakeNoise(12, mons.Pos)
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += 2 + RandInt(4)
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
			g.HandleKill(mons)
		}
		if mons.Kind == MonsBrizzia && RandInt(4) == 0 && !g.Player.HasStatus(StatusNausea) {
			g.Player.Statuses[StatusNausea]++
			heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 30 + RandInt(20), EAction: NauseaEnd})
			g.Print("The brizzia's corpse releases a nauseous gas. You feel sick.")
		}
	} else {
		g.Printf("You miss the %v.", mons.Kind)
	}
	mons.MakeHuntIfHurt(g)
}

func (g *game) HandleKill(mons *monster) {
	g.Killed++
	if g.KilledMons == nil {
		g.KilledMons = map[monsterKind]int{}
	}
	g.KilledMons[mons.Kind]++
	if mons.Kind == MonsExplosiveNadre {
		mons.Explode(g)
	}
	if mons.Kind.Dangerousness() > 10 {
		g.StoryPrintf("You killed %s.", Indefinite(mons.Kind.String(), false))
	}
}
