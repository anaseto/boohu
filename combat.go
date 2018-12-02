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

func (g *game) HitDamage(base int, armor int) (attack int, clang bool) {
	attack = base
	//absorb := g.Absorb(armor)
	//attack -= absorb
	//if absorb > 0 && absorb >= 2*armor/3 && RandInt(2) == 0 {
	//clang = true
	//}
	//if attack < 0 {
	//attack = 0
	//}
	return attack, clang
}

func (m *monster) InflictDamage(g *game, damage, max int) {
	g.Stats.ReceivedHits++
	g.Stats.Damage += damage
	oldHP := g.Player.HP
	g.Player.HP -= damage
	g.ui.WoundedAnimation()
	if oldHP > max && g.Player.HP <= max {
		g.StoryPrintf("Critical HP: %d (hit by %s)", g.Player.HP, m.Kind.Indefinite(false))
		g.ui.CriticalHPWarning()
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
	case TreeStone:
		if !g.Player.HasStatus(StatusLignification) {
			g.UseStone(g.Player.Pos)
			g.EnterLignification(g.Ev)
			g.Print("You feel rooted to the ground.")
		}
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
			if m.SeesPlayer(g) {
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
	case g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification) && !mons.Status(MonsLignified):
		g.SwapWithMonster(mons)
	case g.Player.Weapon == Frundis:
		if !g.HitMonster(mons) {
			break
		}
		if RandInt(2) == 0 {
			mons.EnterConfusion(g, ev)
			g.PrintfStyled("Frundis glows… %s appears confused.", logPlayerHit, mons.Kind.Definite(false))
		}
	case g.Player.Weapon == DancingRapier:
		ompos := mons.Pos
		g.HitMonster(mons)
		if g.Player.HasStatus(StatusLignification) || mons.Status(MonsLignified) {
			break
		}
		dir := ompos.Dir(g.Player.Pos)
		behind := g.Player.Pos.To(dir).To(dir)
		if behind.valid() {
			m := g.MonsterAt(behind)
			if m.Exists() {
				g.HitMonster(m)
			}
		}
		if mons.Exists() {
			mons.MoveTo(g, g.Player.Pos)
		}
		g.PlacePlayerAt(ompos)
	case g.Player.Weapon == HarKarGauntlets:
		g.HarKarAttack(mons, ev)
	case g.Player.Weapon == DefenderFlail:
		g.HitMonster(mons)
	default:
		g.HitMonster(mons)
	}
}

func (g *game) AttractMonster(pos position) *monster {
	dir := pos.Dir(g.Player.Pos)
	for cpos := pos.To(dir); g.Player.LOS[cpos]; cpos = cpos.To(dir) {
		mons := g.MonsterAt(cpos)
		if mons.Exists() {
			mons.MoveTo(g, pos)
			g.ui.TeleportAnimation(cpos, pos, false)
			return mons
		}
	}
	return nil
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
	if pos.valid() && g.Dungeon.Cell(pos).T == FreeCell && !g.Player.HasStatus(StatusLignification) {
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
			g.HitMonster(m)
		}
		g.PlacePlayerAt(pos)
		behind := pos.To(dir)
		m := g.MonsterAt(behind)
		if m.Exists() {
			g.HitMonster(m)
		}
	} else {
		g.HitMonster(mons)
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
		g.HitMonster(mons)
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
		noise -= 5
	}
	if g.Player.Armour == HarmonistRobe {
		noise -= 3
	}
	if g.Player.Armour == Robe {
		noise -= 1
	}
	if clang {
		noise += 5
	}
	return noise
}

type dmgType int

const (
	DmgPhysical dmgType = iota
	DmgMagical
)

func (g *game) HitMonster(mons *monster) (hit bool) {
	ev := g.Ev
	dmg := 2
	if mons.State == Resting {
		dmg += 1
	}
	hit = true
	noise := BaseHitNoise
	if g.Player.Weapon == Dagger || g.Player.Weapon == VampDagger {
		noise -= 2
	}
	if g.Player.Armour == HarmonistRobe {
		noise -= 3
	}
	if g.Player.Weapon == Frundis {
		noise -= 5
	}
	bonus := 0
	if g.Player.HasStatus(StatusBerserk) {
		bonus += 2 + RandInt(4)
	}
	pa := dmg + bonus
	attack, clang := g.HitDamage(pa, 0)
	if clang {
		// TODO
		noise += 3
	}
	g.MakeNoise(noise, mons.Pos)
	if mons.State == Resting {
		if g.Player.Weapon == Dagger || g.Player.Weapon == VampDagger {
			attack *= 4
		} else {
			attack *= 2
		}
	}
	var sclang string
	if clang {
		if RandInt(2) == 0 {
			sclang = " ♫ Clang!"
		} else {
			sclang = " ♪ Clang!"
		}
	}
	oldHP := mons.HP
	mons.HP -= attack
	g.ui.HitAnimation(mons.Pos, false)
	if mons.HP > 0 {
		g.PrintfStyled("You hit %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
	} else if oldHP > 0 {
		// test oldHP > 0 because of sword special attack
		g.PrintfStyled("You kill %s (%d dmg).%s", logPlayerHit, mons.Kind.Definite(false), attack, sclang)
		g.HandleKill(mons, ev) // TODO
	}
	if mons.Kind == MonsBrizzia && RandInt(4) == 0 && !g.Player.HasStatus(StatusNausea) &&
		mons.Pos.Distance(g.Player.Pos) == 1 {
		g.Player.Statuses[StatusNausea]++
		g.PushEvent(&simpleEvent{ERank: ev.Rank() + 30 + RandInt(20), EAction: NauseaEnd})
		g.Print("The brizzia's corpse releases some nauseating gas. You feel sick.")
	}
	if mons.Kind == MonsTinyHarpy && mons.HP > 0 {
		mons.Blink(g)
	}
	g.HandleStone(mons)
	g.Stats.Hits++
	mons.MakeHuntIfHurt(g)
	return hit
}

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
	case TreeStone:
		if mons.Exists() {
			g.UseStone(mons.Pos)
			mons.EnterLignification(g, g.Ev)
		}
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
	if RandInt(2) == 0 {
		sclang = " Clang!"
	} else {
		sclang = " Smash!"
	}
	return sclang
}
