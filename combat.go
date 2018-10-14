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
	g.Stats.ReceivedHits++
	g.Stats.Damage += damage
	oldHP := g.Player.HP
	g.Player.HP -= damage
	g.ui.WoundedAnimation(g)
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical HP: %d (hit by %s)", g.Player.HP, m.Kind.Indefinite(false))
		g.ui.CriticalHPWarning(g)
	}
	if g.Player.HP <= 0 {
		return
	}
	stn, ok := g.MagicalStones[g.Player.Pos]
	if !ok {
		return
	}
	switch stn {
	case TeleStone:
		g.UseStone(g.Player.Pos)
		g.Teleportation(g.Ev)
	case FogStone:
		g.Fog(g.Player.Pos, 3, g.Ev)
		g.UseStone(g.Player.Pos)
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, g.Player.Pos)
		dij := &normalPath{game: g}
		nm := Dijkstra(dij, []position{g.Player.Pos}, 2)
		for _, m := range g.Monsters {
			if !m.Exists() {
				continue
			}
			if m.State == Resting {
				continue
			}
			_, ok := nm[m.Pos]
			if !ok {
				continue
			}
			m.EnterConfusion(g, g.Ev)
		}
		//g.Confusion(g.Ev)
		g.UseStone(g.Player.Pos)
	case ObstructionStone:
		neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
		for _, pos := range neighbors {
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			g.CreateTemporalWallAt(pos, g.Ev)
		}
		g.Printf("You see walls appear out of thin air around the stone.")
		g.UseStone(g.Player.Pos)
		g.ComputeLOS()
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
	g.SwapWithMonster(mons)
}

func (g *game) AttractMonster(pos position) *monster {
	dir := pos.Dir(g.Player.Pos)
	for cpos := pos.To(dir); g.Player.LOS[cpos]; cpos = cpos.To(dir) {
		mons := g.MonsterAt(cpos)
		if mons.Exists() {
			mons.MoveTo(g, pos)
			g.ui.TeleportAnimation(g, cpos, pos, false)
			return mons
		}
	}
	return nil
}

func (g *game) HitNoise(clang bool) int {
	noise := BaseHitNoise
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

func (g *game) HandleStone(mons *monster) {
	stn, ok := g.MagicalStones[mons.Pos]
	if !ok {
		return
	}
	switch stn {
	case TeleStone:
		if mons.Exists() {
			g.UseStone(mons.Pos)
			mons.TeleportAway(g)
		}
	case FogStone:
		g.Fog(mons.Pos, 3, g.Ev)
		g.UseStone(mons.Pos)
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, mons.Pos)
		dij := &normalPath{game: g}
		nm := Dijkstra(dij, []position{mons.Pos}, 2)
		for _, m := range g.Monsters {
			if !m.Exists() {
				continue
			}
			if m.State == Resting {
				continue
			}
			_, ok := nm[m.Pos]
			if !ok {
				continue
			}
			m.EnterConfusion(g, g.Ev)
		}
		// _, ok := nm[g.Player.Pos]
		// if ok {
		// 	g.Confusion(g.Ev)
		// }
		g.UseStone(mons.Pos)
	case ObstructionStone:
		if !mons.Exists() {
			g.CreateTemporalWallAt(mons.Pos, g.Ev)
		}
		neighbors := g.Dungeon.FreeNeighbors(mons.Pos)
		for _, pos := range neighbors {
			if pos == g.Player.Pos {
				continue
			}
			m := g.MonsterAt(pos)
			if m.Exists() {
				continue
			}
			g.CreateTemporalWallAt(pos, g.Ev)
		}
		g.Printf("You see walls appear out of thin air around the stone.")
		g.UseStone(mons.Pos)
		g.ComputeLOS()
	}
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
	TemporalWallNoise   = 13
	ExplosionHitNoise   = 13
	ExplosionNoise      = 18
	MagicHitNoise       = 15
	BarkNoise           = 13
	MagicExplosionNoise = 16
	MagicCastNoise      = 16
	BaseHitNoise        = 11
	ShieldBlockNoise    = 17
	QueenStoneNoise     = 19
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
	g.Stats.Blocks++
	switch g.Player.Shield {
	case EarthShield:
		dir := m.Pos.Dir(g.Player.Pos)
		lat := g.Player.Pos.Laterals(dir)
		for _, pos := range lat {
			if !pos.valid() {
				continue
			}
			if RandInt(3) == 0 && g.Dungeon.Cell(pos).T == WallCell {
				g.Dungeon.SetCell(pos, FreeCell)
				g.Stats.Digs++
				g.MakeNoise(WallNoise+3, pos)
				g.Fog(pos, 1, g.Ev)
				g.Printf("%s The sound of blocking breaks a wall.", g.CrackSound())
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
			g.Printf("%s is repelled.", m.Kind.Definite(true))
		}
	case ConfusingShield:
		if m.Pos.Distance(g.Player.Pos) > 1 {
			break
		}
		if RandInt(4) == 0 {
			m.EnterConfusion(g, g.Ev)
			g.Printf("%s appears confused.", m.Kind.Definite(true))
		}
	case FireShield:
		dir := m.Pos.Dir(g.Player.Pos)
		burnpos := g.Player.Pos.To(dir)
		if RandInt(4) == 0 {
			g.Print("Sparks emerge out of the shield.")
			g.Burn(burnpos, g.Ev)
		}
	}
}
