// combat utility functions

package main

func (g *game) HitDamage(dt dmgType, base int, armor int) (attack int, clang bool) {
	min := base / 2
	attack = min + RandInt(base-min+1)
	if dt == DmgPhysical {
		absorb := RandInt(armor + 1)
		if absorb > 0 && absorb >= 2*armor/3 && RandInt(2) == 0 {
			clang = true
		}
		attack -= absorb
	}
	if attack < 0 {
		attack = 0
	}
	return attack, clang
}

func (m *monster) InflictDamage(g *game, damage, max int) {
	oldHP := g.Player.HP
	g.Player.HP -= damage
	g.ui.WoundedAnimation(g)
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical HP: %d (hit by %s)", g.Player.HP, m.Kind.Indefinite(false))
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
		if v > 25 {
			v = 25
		}
		r := RandInt(30)
		if m.State == Resting {
			v /= 2
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
	case g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification):
		g.SwapWithMonster(mons)
	case g.Player.Weapon == Frundis:
		if !g.HitMonster(DmgPhysical, mons, ev) {
			break
		}
		if RandInt(4) == 0 {
			mons.EnterConfusion(g, ev)
			g.PrintfStyled("Frundis glows… the %s appears confused.", logPlayerHit, mons.Kind)
		}
	case g.Player.Weapon.Cleave():
		var neighbors []position
		if g.Player.HasStatus(StatusConfusion) {
			neighbors = g.Dungeon.CardinalFreeNeighbors(g.Player.Pos)
		} else {
			neighbors = g.Dungeon.FreeNeighbors(g.Player.Pos)
		}
		for _, pos := range neighbors {
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				g.HitMonster(DmgPhysical, mons, ev)
			}
		}
	case g.Player.Weapon.Pierce():
		g.HitMonster(DmgPhysical, mons, ev)
		deltaX := mons.Pos.X - g.Player.Pos.X
		deltaY := mons.Pos.Y - g.Player.Pos.Y
		behind := position{g.Player.Pos.X + 2*deltaX, g.Player.Pos.Y + 2*deltaY}
		if behind.valid() {
			mons := g.MonsterAt(behind)
			if mons.Exists() {
				g.HitMonster(DmgPhysical, mons, ev)
			}
		}
	case g.Player.Weapon == ElecWhip:
		g.HitConnected(mons.Pos, DmgMagical, ev)
	default:
		g.HitMonster(DmgPhysical, mons, ev)
		if (g.Player.Weapon == Sword || g.Player.Weapon == DoubleSword) && RandInt(4) == 0 {
			g.HitMonster(DmgPhysical, mons, ev)
		}
	}
}

func (g *game) HitConnected(pos position, dt dmgType, ev event) {
	// inspired from TOME4's Harbinger addon class
	d := g.Dungeon
	conn := map[position]bool{}
	stack := []position{pos}
	conn[pos] = true
	nb := make([]position, 0, 8)
	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		g.HitMonster(dt, mons, ev)
		nb = pos.Neighbors(nb, func(npos position) bool {
			return npos.valid() && d.Cell(npos).T != WallCell
		})
		for _, npos := range nb {
			if !conn[npos] {
				conn[npos] = true
				stack = append(stack, npos)
			}
		}
	}
}

func (g *game) HitNoise(clang bool) int {
	noise := BaseHitNoise
	if g.Player.Weapon == Frundis {
		noise -= 2
	}
	if g.Player.Armour == Robe {
		noise -= 1
	}
	if clang {
		noise += g.Player.Armor()
	}
	return noise
}

type dmgType int

const (
	DmgPhysical dmgType = iota
	DmgMagical
)

func (g *game) HitMonster(dt dmgType, mons *monster, ev event) (hit bool) {
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		hit = true
		noise := BaseHitNoise
		if g.Player.Weapon == Frundis || g.Player.Weapon == Dagger {
			noise -= 2
		}
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += 2 + RandInt(4)
		}
		attack, clang := g.HitDamage(dt, g.Player.Attack()+bonus, mons.Armor)
		if clang {
			noise += mons.Armor
		}
		g.MakeNoise(noise, mons.Pos)
		if mons.State == Resting {
			if g.Player.Weapon == Dagger {
				attack *= 4
			} else {
				attack *= 2
			}
		}
		var sclang string
		if clang {
			if mons.Armor > 3 {
				sclang = " ♫ Clang!"
			} else {
				sclang = " ♪ Clang!"
			}
		}
		oldHP := mons.HP
		mons.HP -= attack
		g.ui.HitAnimation(g, mons.Pos, false)
		if mons.HP > 0 {
			g.PrintfStyled("You hit %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
		} else if oldHP > 0 {
			// test oldHP > 0 because of sword special attack
			g.PrintfStyled("You kill %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
			g.HandleKill(mons, ev)
		}
		if mons.Kind == MonsBrizzia && RandInt(4) == 0 && !g.Player.HasStatus(StatusNausea) &&
			mons.Pos.Distance(g.Player.Pos) == 1 {
			g.Player.Statuses[StatusNausea]++
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 30 + RandInt(20), EAction: NauseaEnd})
			g.Print("The brizzia's corpse releases a nauseous gas. You feel sick.")
		}
		g.Stats.Hits++
	} else {
		g.Printf("You miss %s.", mons.Kind.Definite(false))
		g.Stats.Misses++
	}
	mons.MakeHuntIfHurt(g)
	return hit
}

func (g *game) HandleKill(mons *monster, ev event) {
	g.Stats.Killed++
	g.Stats.KilledMons[mons.Kind]++
	if mons.Kind == MonsExplosiveNadre {
		mons.Explode(g, ev)
	}
	if g.Doors[mons.Pos] {
		g.ComputeLOS()
	}
	if mons.Kind.Dangerousness() > 10 {
		g.StoryPrintf("You killed %s.", mons.Kind.Indefinite(false))
	}
}

const (
	WallNoise           = 18
	TemporalWallNoise   = 16
	ExplosionHitNoise   = 13
	ExplosionNoise      = 18
	MagicHitNoise       = 15
	BarkNoise           = 13
	MagicExplosionNoise = 16
	MagicCastNoise      = 16
	BaseHitNoise        = 11
	ShieldBlockNoise    = 15
)

func (g *game) ArmourClang() (sclang string) {
	if g.Player.Armor() > 3 {
		sclang = " Clang!"
	} else {
		sclang = " Smash!"
	}
	return sclang
}
