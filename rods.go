package main

import (
	"errors"
)

func (g *game) EvokeRodBlink(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink(ev)
	return nil
}

func (g *game) BlinkPos() position {
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		// TODO: skip if not seen?
		if !b {
			continue
		}
		if !g.Dungeon.Cell(pos).IsFree() {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, pos)
	}
	if len(losPos) == 0 {
		return InvalidPos
	}
	npos := losPos[RandInt(len(losPos))]
	for i := 0; i < 4; i++ {
		pos := losPos[RandInt(len(losPos))]
		if npos.Distance(g.Player.Pos) < pos.Distance(g.Player.Pos) {
			npos = pos
		}
	}
	return npos
}

func (g *game) Blink(ev event) {
	if g.Player.HasStatus(StatusLignification) {
		return
	}
	npos := g.BlinkPos()
	if !npos.valid() {
		// should not happen
		g.Print("You could not blink.")
		return
	}
	opos := g.Player.Pos
	g.Print("You blink away.")
	g.ui.TeleportAnimation(opos, npos, true)
	g.PlacePlayerAt(npos)
}

func (g *game) EvokeRodTeleportOther(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	mons.TeleportAway(g)
	return nil
}

func (g *game) EvokeRodSleeping(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Print("A sleeping ball emerges straight out of the rod.")
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgSleepingMonster)
	for _, pos := range append(neighbors, g.Player.Target) {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		}
		mons.State = Resting
		mons.Dir = NoDir
		mons.ExhaustTime(g, 40+RandInt(10))
	}
	return nil
}

func (g *game) EvokeRodFireBolt(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{flammable: true}); err != nil {
		return err
	}
	ray := g.Ray(g.Player.Target)
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.Print("Whoosh! A fire bolt emerges straight out of the rod.")
	g.ui.FireBoltAnimation(ray)
	for _, pos := range ray {
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		dmg := 1
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the bolt.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodFireball(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true, flammable: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.MakeNoise(MagicExplosionNoise, g.Player.Target)
	g.Printf("A fireball emerges straight out of the rod... %s", g.ExplosionSound())
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
	g.ui.ExplosionAnimation(FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		dmg := 1 + RandInt(2)
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the fireball.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodLightning(ev event) error {
	d := g.Dungeon
	conn := map[position]bool{}
	nb := make([]position, 0, 8)
	nb = g.Player.Pos.Neighbors(nb, func(npos position) bool {
		return npos.valid() && d.Cell(npos).T != WallCell
	})
	stack := []position{}
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.Print("Whoosh! Lightning emerges straight out of the rod.")
	for _, pos := range nb {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		stack = append(stack, pos)
		conn[pos] = true
	}
	if len(stack) == 0 {
		return errors.New("There are no adjacent monsters.")
	}
	var pos position
	targets := []position{}
	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		targets = append(targets, pos)
		dmg := 1
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by lightning.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(MagicHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
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
	g.ui.LightningHitAnimation(targets)

	return nil
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeRodFog(ev event) error {
	g.Fog(g.Player.Pos, 3, ev)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at position, radius int, ev event) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	for pos := range nm {
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	}
	g.ComputeLOS()
}

func (g *game) EvokeRodDigging(ev event) error {
	if err := g.ui.ChooseTarget(&wallChooser{}); err != nil {
		return err
	}
	pos := g.Player.Target
	for i := 0; i < 3; i++ {
		g.Dungeon.SetCell(pos, GroundCell)
		g.Stats.Digs++
		g.MakeNoise(WallNoise, pos)
		g.Fog(pos, 1, ev)
		pos = pos.To(pos.Dir(g.Player.Pos))
		if !g.Player.Sees(pos) {
			g.TerrainKnowledge[pos] = WallCell
		}
		if !pos.valid() || g.Dungeon.Cell(pos).T != WallCell {
			break
		}
	}
	g.Print("You see the wall disintegrate with a crash.")
	g.ComputeLOS()
	g.MakeMonstersAware()
	return nil
}

func (g *game) EvokeRodShatter(ev event) error {
	if err := g.ui.ChooseTarget(&wallChooser{minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Dungeon.SetCell(g.Player.Target, GroundCell)
	g.Stats.Digs++
	g.ComputeLOS()
	g.MakeMonstersAware()
	g.MakeNoise(WallNoise, g.Player.Target)
	g.Printf("%s The wall disappeared.", g.CrackSound())
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionWallStart)
	g.ui.ExplosionAnimation(WallExplosion, g.Player.Target)
	g.Fog(g.Player.Target, 2, ev)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		dmg := 2
		mons.HP -= dmg
		if mons.HP <= 0 {
			g.Printf("%s is killed by the explosion.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(ExplosionHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodObstruction(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{free: true}); err != nil {
		return err
	}
	g.TemporalWallAt(g.Player.Target, ev)
	g.Printf("You see a wall appear out of thin air.")
	return nil
}

func (g *game) EvokeRodLignification(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in targeter)
	if mons.Status(MonsLignified) {
		return errors.New("You cannot target a lignified monster.")
	}
	mons.EnterLignification(g, ev)
	return nil
}

func (g *game) TemporalWallAt(pos position, ev event) {
	if g.Dungeon.Cell(pos).T == WallCell {
		return
	}
	if !g.Player.Sees(pos) {
		g.TerrainKnowledge[pos] = g.Dungeon.Cell(pos).T
	}
	g.CreateTemporalWallAt(pos, ev)
	g.ComputeLOS()
}

func (g *game) CreateTemporalWallAt(pos position, ev event) {
	g.Dungeon.SetCell(pos, WallCell)
	delete(g.Clouds, pos)
	g.TemporalWalls[pos] = true
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationTemporalWall + RandInt(DurationTemporalWall/2), Pos: pos, EAction: ObstructionEnd})
}

func (g *game) EvokeRodHope(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
		return err
	}
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgExplosionStart)
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	dmg := DefaultHealth - g.Player.HP + 1
	if dmg <= 0 {
		dmg = 1
	}
	mons.HP -= dmg
	g.Burn(g.Player.Target, ev)
	g.ui.HitAnimation(g.Player.Target, true)
	g.Printf("An energy channel hits %s (%d dmg).", mons.Kind.Definite(false), dmg)
	if mons.HP <= 0 {
		g.Printf("%s dies.", mons.Kind.Indefinite(true))
		g.HandleKill(mons, ev)
	}
	return nil
}

func (g *game) EvokeRodSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	if mons.Status(MonsLignified) {
		return errors.New("You cannot target a lignified monster.")
	}
	g.SwapWithMonster(mons)
	return nil
}

func (g *game) SwapWithMonster(mons *monster) {
	ompos := mons.Pos
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.ui.SwappingAnimation(mons.Pos, g.Player.Pos)
	mons.MoveTo(g, g.Player.Pos)
	g.PlacePlayerAt(ompos)
	mons.MakeAware(g)
}
