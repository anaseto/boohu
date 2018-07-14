// combat utility functions

package main

func (g *game) Absorb(armor int) int {
	absorb := 0
	for i := 0; i <= 2; i++ {
		absorb += RandInt(armor + 1)
	}
	q := absorb / 3
	r := absorb % 3
	if r == 2 {
		q++
	}
	return q
}

func (g *game) HitDamage(dt dmgType, base int, armor int) (attack int, clang bool) {
	min := base / 2
	attack = min + RandInt(base-min+1)
	absorb := g.Absorb(armor)
	if dt == DmgMagical {
		absorb /= 2
	}
	attack -= absorb
	if absorb > 0 && absorb >= 2*armor/3 && RandInt(2) == 0 {
		clang = true
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
		if m.Status(MonsExhausted) {
			v = 2 * v / 3
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

func (g *game) InOpenMons(mons *monster) bool {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		if pos.Distance(mons.Pos) > 1 {
			continue
		}
		if g.Dungeon.Cell(pos).T == WallCell {
			return false
		}
	}
	return true
}

func (g *game) AttackMonster(mons *monster, ev event) {
	switch {
	case g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification):
		g.SwapWithMonster(mons)
	case g.Player.Weapon == Frundis:
		if !g.HitMonster(DmgPhysical, mons, ev) {
			break
		}
		if RandInt(2) == 0 {
			mons.EnterConfusion(g, ev)
			g.PrintfStyled("Frundis glows… %s appears confused.", logPlayerHit, mons.Kind.Definite(false))
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
		dir := mons.Pos.Dir(g.Player.Pos)
		behind := g.Player.Pos.To(dir).To(dir)
		if behind.valid() {
			mons := g.MonsterAt(behind)
			if mons.Exists() {
				g.HitMonster(DmgPhysical, mons, ev)
			}
		}
	case g.Player.Weapon == ElecWhip:
		g.HitConnected(mons.Pos, DmgMagical, ev)
	case g.Player.Weapon == DancingRapier:
		g.HitMonster(DmgPhysical, mons, ev)
		if mons.Exists() {
			dir := mons.Pos.Dir(g.Player.Pos)
			behind := g.Player.Pos.To(dir).To(dir)
			if behind.valid() {
				mons := g.MonsterAt(behind)
				if mons.Exists() {
					g.HitMonster(DmgPhysical, mons, ev)
				}
			}
			if !g.Player.HasStatus(StatusLignification) {
				ompos := mons.Pos
				mons.MoveTo(g, g.Player.Pos)
				g.PlacePlayerAt(ompos)
			}
		} else if !g.Player.HasStatus(StatusLignification) {
			g.PlacePlayerAt(mons.Pos)
		}
	case g.Player.Weapon == HarKarGauntlets:
		g.HarKarAttack(mons, ev)
	case g.Player.Weapon == BerserkSword:
		g.HitMonster(DmgPhysical, mons, ev)
		if RandInt(20) == 0 && !g.Player.HasStatus(StatusExhausted) && !g.Player.HasStatus(StatusBerserk) {
			g.Player.Statuses[StatusBerserk] = 1
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 65 + RandInt(20), EAction: BerserkEnd})
			g.Printf("Your sword insurges you to kill things.", BerserkPotion)
		}
	default:
		g.HitMonster(DmgPhysical, mons, ev)
	}
}

func (g *game) HarKarAttack(mons *monster, ev event) {
	dir := mons.Pos.Dir(g.Player.Pos)
	pos := g.Player.Pos
	for {
		pos = pos.To(dir)
		if !pos.valid() || g.Dungeon.Cell(pos).T != FreeCell {
			break
		}
		m := g.MonsterAt(pos)
		if !m.Exists() {
			break
		}
	}
	if pos.valid() && g.Dungeon.Cell(pos).T == FreeCell {
		pos = g.Player.Pos
		for {
			pos = pos.To(dir)
			if !pos.valid() || g.Dungeon.Cell(pos).T != FreeCell {
				break
			}
			m := g.MonsterAt(pos)
			if !m.Exists() {
				break
			}
			g.HitMonster(DmgPhysical, m, ev)
		}
		if !g.Player.HasStatus(StatusLignification) {
			g.PlacePlayerAt(pos)
		}
	} else {
		g.HitMonster(DmgPhysical, mons, ev)
	}
}

func (g *game) HitConnected(pos position, dt dmgType, ev event) {
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
		noise -= 4
	}
	if g.Player.Armour == HarmonistRobe {
		noise -= 3
	}
	if g.Player.Armour == Robe {
		noise -= 1
	}
	if clang {
		arnoise := g.Player.Armor()
		if arnoise > 7 {
			arnoise = 7
		}
		noise += arnoise
	}
	return noise
}

type dmgType int

const (
	DmgPhysical dmgType = iota
	DmgMagical
)

func (g *game) HitMonster(dt dmgType, mons *monster, ev event) (hit bool) {
	maxacc := g.Player.Accuracy()
	if g.Player.Weapon == Sabre && mons.HP > 0 {
		maxacc += int(6 * (-1 + float64(mons.HPmax)/float64(mons.HP)))
	}
	acc := RandInt(maxacc)
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		hit = true
		noise := BaseHitNoise
		if g.Player.Weapon == Dagger {
			noise -= 2
		}
		if g.Player.Armour == HarmonistRobe {
			noise -= 3
		}
		if g.Player.Weapon == Frundis {
			noise -= 4
		}
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += 2 + RandInt(4)
		}
		pa := g.Player.Attack() + bonus
		if g.Player.Weapon.Cleave() && g.InOpenMons(mons) {
			if g.Player.Attack() >= 15 {
				pa += 1 + RandInt(3)
			} else {
				pa += 1 + RandInt(2)
			}
		}
		attack, clang := g.HitDamage(dt, pa, mons.Armor)
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

func (g *game) BlockEffects(m *monster) {
	switch g.Player.Shield {
	case EarthShield:
		dir := m.Pos.Dir(g.Player.Pos)
		lat := g.Player.Pos.Laterals(dir)
		for _, pos := range lat {
			if !pos.valid() {
				continue
			}
			if RandInt(4) == 0 && g.Dungeon.Cell(pos).T == WallCell {
				g.Dungeon.SetCell(pos, FreeCell)
				g.Stats.Digs++
				g.MakeNoise(WallNoise, pos)
				g.Fog(pos, 1, g.Ev)
			}
		}
	case BashingShield:
		if m.Kind == MonsSatowalgaPlant || m.Pos.Distance(g.Player.Pos) > 1 {
			break
		}
		if RandInt(3) == 0 {
			break
		}
		dir := m.Pos.Dir(g.Player.Pos)
		pos := m.Pos
		for i := 0; i < 3; i++ {
			npos := pos.To(dir)
			if !npos.valid() || g.Dungeon.Cell(npos).T == WallCell {
				break
			}
			mons := g.MonsterAt(npos)
			if mons.Exists() {
				break
			}
			pos = npos
		}
		if !m.Status(MonsExhausted) {
			m.Statuses[MonsExhausted] = 1
			g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + 100 + RandInt(50), NMons: m.Index, EAction: MonsExhaustionEnd})
		}
		if pos != m.Pos {
			m.MoveTo(g, pos)
		}
	case ConfusingShield:
		if m.Pos.Distance(g.Player.Pos) > 1 {
			break
		}
		if RandInt(4) == 0 {
			m.EnterConfusion(g, g.Ev)
			g.Printf("%s appears confused.", m.Kind.Definite(true))
		}
	}
}
