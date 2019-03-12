package main

import (
	"errors"
	"sort"
)

func (g *game) QuaffTeleportation(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	if g.Player.HasStatus(StatusTele) {
		return errors.New("You already quaffed a potion of teleportation.")
	}
	delay := DurationTeleportationDelay
	g.Player.Statuses[StatusTele] = 1
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + delay, EAction: Teleportation})
	//g.Printf("You quaff the %s. You feel unstable.", TeleportationPotion)
	return nil
}

//func (g *game) QuaffBerserk(ev event) error {
//if g.Player.HasStatus(StatusExhausted) {
//return errors.New("You are too exhausted to berserk.")
//}
//if g.Player.HasStatus(StatusBerserk) {
//return errors.New("You are already berserk.")
//}
//g.Player.Statuses[StatusBerserk] = 1
//end := ev.Rank() + DurationBerserk
//g.PushEvent(&simpleEvent{ERank: end, EAction: BerserkEnd})
//g.Player.Expire[StatusBerserk] = end
//g.Printf("You quaff the %s. You feel a sudden urge to kill things.", BerserkPotion)
//g.Player.HPbonus += 2
//return nil
//}

func (g *game) QuaffHealWounds(ev event) error {
	//hp := g.Player.HP
	g.Player.HP = g.Player.HPMax()
	//g.Printf("You quaff the %s (%d -> %d).", HealWoundsPotion, hp, g.Player.HP)
	return nil
}

func (g *game) QuaffMagic(ev event) error {
	//mp := g.Player.MP
	g.Player.MP += 2 * g.Player.MPMax() / 3
	if g.Player.MP > g.Player.MPMax() {
		g.Player.MP = g.Player.MPMax()
	}
	//g.Printf("You quaff the %s (%d -> %d).", MagicPotion, mp, g.Player.MP)
	return nil
}

func (g *game) QuaffDescent(ev event) error {
	// why not?
	//if g.Player.HasStatus(StatusLignification) {
	//return errors.New("You cannot descend while lignified.")
	//}
	if g.Depth >= MaxDepth {
		return errors.New("You cannot descend any deeper!")
	}
	//g.Printf("You quaff the %s. You fall through the ground.", DescentPotion)
	g.LevelStats()
	g.StoryPrint("You descended deeper into the dungeon.")
	g.Depth++
	g.DepthPlayerTurn = 0
	g.InitLevel()
	g.Save()
	return nil
}

func (g *game) QuaffSwiftness(ev event) error {
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + DurationSwiftness
	g.PushEvent(&simpleEvent{ERank: end, EAction: HasteEnd})
	g.Player.Expire[StatusSwift] = end
	g.Player.Statuses[StatusAgile]++
	g.PushEvent(&simpleEvent{ERank: end, EAction: EvasionEnd})
	g.Player.Expire[StatusAgile] = end
	//g.Printf("You quaff the %s. You feel speedy and agile.", SwiftnessPotion)
	return nil
}

func (g *game) QuaffDigPotion(ev event) error {
	g.Player.Statuses[StatusDig] = 1
	end := ev.Rank() + DurationDigging
	g.PushEvent(&simpleEvent{ERank: end, EAction: DigEnd})
	g.Player.Expire[StatusDig] = end
	//g.Printf("You quaff the %s. You feel like an earth dragon.", DigPotion)
	return nil
}

func (g *game) QuaffSwapPotion(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot drink this potion while lignified.")
	}
	g.Player.Statuses[StatusSwap] = 1
	end := ev.Rank() + DurationSwap
	g.PushEvent(&simpleEvent{ERank: end, EAction: SwapEnd})
	g.Player.Expire[StatusSwap] = end
	//g.Printf("You quaff the %s. You feel light-footed.", SwapPotion)
	return nil
}

func (g *game) QuaffShadowsPotion(ev event) error {
	if g.Player.HasStatus(StatusShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Player.Statuses[StatusShadows] = 1
	end := ev.Rank() + DurationShadows
	g.PushEvent(&simpleEvent{ERank: end, EAction: ShadowsEnd})
	g.Player.Expire[StatusShadows] = end
	//g.Printf("You quaff the %s. You feel surrounded by shadows.", ShadowsPotion)
	g.ComputeLOS()
	return nil
}

//func (g *game) QuaffLignification(ev event) error {
//if g.Player.HasStatus(StatusLignification) {
//return errors.New("You are already lignified.")
//}
//g.EnterLignification(ev)
//g.Printf("You quaff the %s. You feel rooted to the ground.", LignificationPotion)
//return nil
//}

func (g *game) QuaffMagicMapping(ev event) error {
	dp := &dungeonPath{dungeon: g.Dungeon}
	g.AutoExploreDijkstra(dp, []int{g.Player.Pos.idx()})
	cdists := make(map[int][]int)
	for i, dist := range DijkstraMapCache {
		cdists[dist] = append(cdists[dist], i)
	}
	var dists []int
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	g.ui.DrawDungeonView(NormalMode)
	for _, d := range dists {
		draw := false
		for _, i := range cdists[d] {
			pos := idxtopos(i)
			c := g.Dungeon.Cell(pos)
			if (c.IsFree() || g.Dungeon.HasFreeNeighbor(pos)) && !c.Explored {
				g.Dungeon.SetExplored(pos)
				draw = true
			}
		}
		if draw {
			g.ui.MagicMappingAnimation(cdists[d])
		}
	}
	//g.Printf("You quaff the %s. You feel aware of your surroundings..", MagicMappingPotion)
	return nil
}

func (g *game) QuaffTormentPotion(ev event) error {
	//g.Printf("You quaff the %s. %s It hurts!", TormentPotion, g.ExplosionSound())
	g.DamagePlayer(g.Player.HP / 2)
	g.ui.WoundedAnimation()
	g.MakeNoise(ExplosionNoise+10, g.Player.Pos)
	g.ui.TormentExplosionAnimation()
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		g.ExplosionAt(ev, pos)
	}
	return nil
}

func (g *game) QuaffDreamPotion(ev event) error {
	for _, mons := range g.Monsters {
		if mons.Exists() && mons.State == Resting && !g.Player.Sees(mons.Pos) {
			mons.UpdateKnowledge(g, mons.Pos)
		}
	}
	//g.Printf("You quaff the %s. You perceive monsters' dreams.", DreamPotion)
	return nil
}

func (g *game) QuaffWallPotion(ev event) error {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		g.CreateTemporalWallAt(pos, ev)
	}
	//g.Printf("You quaff the %s. You feel surrounded by temporary walls.", WallPotion)
	g.ComputeLOS()
	return nil
}

func (g *game) QuaffCBlinkPotion(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	if err := g.ui.ChooseTarget(&chooser{free: true}); err != nil {
		return err
	}
	//g.Printf("You quaff the %s. You blink.", CBlinkPotion)
	g.PlacePlayerAt(g.Player.Target)
	return nil
}

func (g *game) ThrowConfusingDart(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	attack := 1
	mons.HP -= attack
	if mons.HP > 0 {
		mons.EnterConfusion(g, ev)
		//g.PrintfStyled("Your %s hits the %s (%d dmg), who appears confused.", logPlayerHit, ConfusingDart, mons.Kind, attack)
		g.ui.ThrowAnimation(g.Ray(mons.Pos), true)
		mons.MakeHuntIfHurt(g)
	} else {
		//g.PrintfStyled("Your %s kills the %s.", logPlayerHit, ConfusingDart, mons.Kind)
		g.ui.ThrowAnimation(g.Ray(mons.Pos), true)
		g.HandleKill(mons, ev)
	}
	g.HandleStone(mons)
	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) ExplosionAt(ev event, pos position) {
	g.Burn(pos, ev)
	mons := g.MonsterAt(pos)
	if mons.Exists() {
		mons.HP /= 2
		if mons.HP <= 0 {
			g.HandleKill(mons, ev)
			if g.Player.Sees(mons.Pos) {
				g.Printf("%s dies.", mons.Kind.Definite(true))
			}
		}
		g.MakeNoise(ExplosionHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	} else if c := g.Dungeon.Cell(pos); !c.IsFree() && RandInt(2) == 0 {
		g.Dungeon.SetCell(pos, GroundCell)
		g.Stats.Digs++
		if !g.Player.Sees(pos) {
			g.TerrainKnowledge[pos] = c.T
		} else {
			g.ui.WallExplosionAnimation(pos)
		}
		g.MakeNoise(WallNoise, pos)
		g.Fog(pos, 1, ev)
	}
}

func (g *game) ThrowExplosiveMagara(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true, flammable: true, wall: true}); err != nil {
		return err
	}
	neighbors := g.Player.Target.ValidNeighbors()
	g.Printf("You throw the explosive magara... %s", g.ExplosionSound())
	g.MakeNoise(ExplosionNoise, g.Player.Target)
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgPlayer)
	g.ui.ExplosionAnimation(FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.ExplosionAt(ev, pos)
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) ThrowTeleportMagara(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{area: true, minDist: true}); err != nil {
		return err
	}
	neighbors := g.Player.Target.ValidNeighbors()
	g.Print("You throw the teleport magara.")
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgPlayer)
	for _, pos := range append(neighbors, g.Player.Target) {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.TeleportAway(g)
		}
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) ThrowSlowingMagara(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{}); err != nil {
		return err
	}
	ray := g.Ray(g.Player.Target)
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.Print("Whoosh! A bolt of slowing emerges out of the magara.")
	g.ui.SlowingMagaraAnimation(ray)
	for _, pos := range ray {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		mons.Statuses[MonsSlow]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + DurationSlow, NMons: mons.Index, EAction: MonsSlowEnd})
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) ThrowConfuseMagara(ev event) error {
	//g.Printf("You activate the %s. A harmonic light confuses monsters.", ConfuseMagara)
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.EnterConfusion(g, ev)
		}
	}

	ev.Renew(g, DurationThrowItem)
	return nil
}

func (g *game) NightFog(at position, radius int, ev event) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	for pos := range nm {
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudNight
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: NightProgression, Pos: pos})
			g.MakeCreatureSleep(pos, ev)
		}
	}
	g.ComputeLOS()
}

func (g *game) ThrowNightMagara(ev event) error {
	if err := g.ui.ChooseTarget(&chooser{needsFreeWay: true}); err != nil {
		return err
	}
	g.Print("You throw the night magara… Clouds come out of it.")
	g.ui.ProjectileTrajectoryAnimation(g.Ray(g.Player.Target), ColorFgSleepingMonster)
	g.NightFog(g.Player.Target, 2, ev)

	ev.Renew(g, DurationThrowItem)
	return nil
}
