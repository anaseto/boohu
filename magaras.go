package main

import (
	"errors"
	"fmt"
)

type magara int

const (
	NoMagara magara = iota
	BlinkMagara
	DigMagara
	TeleportMagara
	TeleportOtherMagara
	HealWoundsMagara
	SwiftnessMagara
	SwappingMagara
	ShadowsMagara
	FogMagara
	WallsMagara
	SlowingMagara
	SleepingMagara
	NoiseMagara
	ObstructionMagara
	FireMagara
	ConfusionMagara
	LignificationMagara
)

const NumMagaras = int(LignificationMagara)

func (g *game) RandomMagara() (mag magara) {
loop:
	for {
		mag = magara(1 + RandInt(NumMagaras))
		for _, m := range g.GeneratedMagaras {
			if m == mag {
				continue loop
			}
		}
		break
	}
	return mag
}

func (g *game) EquipMagara(i int, ev event) (err error) {
	omagara := g.Player.Magaras[i]
	g.Player.Magaras[i] = g.Objects.Magaras[g.Player.Pos]
	g.Objects.Magaras[g.Player.Pos] = omagara
	g.Printf("You equip %s, leaving %s on the ground.", g.Player.Magaras[i], omagara)
	g.StoryPrintf("You equip %s, leaving %s.", g.Player.Magaras[i], omagara)
	ev.Renew(g, 5)
	return nil
}

func (g *game) UseMagara(n int, ev event) (err error) {
	if g.Player.HasStatus(StatusNausea) {
		return errors.New("You cannot use magaras while sick.")
	}
	mag := g.Player.Magaras[n]
	if mag.MPCost(g) > g.Player.MP {
		return errors.New("Not enough magic points for using this rod.")
	}
	switch mag {
	case NoMagara:
		err = errors.New("You cannot evoke an empty slot!")
	case BlinkMagara:
		err = g.EvokeBlink(ev)
	case TeleportMagara:
		err = g.EvokeTeleport(ev)
	case DigMagara:
		err = g.EvokeDig(ev)
	case TeleportOtherMagara:
		err = g.EvokeTeleportOther(ev)
	case HealWoundsMagara:
		err = g.EvokeHealWounds(ev)
	//case MagicMagara:
	//err = g.EvokeRefillMagic(ev)
	case SwiftnessMagara:
		err = g.EvokeSwiftness(ev)
	case SwappingMagara:
		err = g.EvokeSwapping(ev)
	case ShadowsMagara:
		err = g.EvokeShadows(ev)
	case FogMagara:
		err = g.EvokeFog(ev)
	case WallsMagara:
		err = g.EvokeWalls(ev)
	case SlowingMagara:
		err = g.EvokeSlowing(ev)
	case SleepingMagara:
		err = g.EvokeSleeping(ev)
	case NoiseMagara:
		err = g.EvokeNoise(ev)
	case ConfusionMagara:
		err = g.EvokeConfusion(ev)
	case ObstructionMagara:
		err = g.EvokeObstruction(ev)
	case FireMagara:
		err = g.EvokeFire(ev)
	case LignificationMagara:
		err = g.EvokeLignification(ev)
		//case MagicMappingMagara:
		//err = g.EvokeMagicMapping(ev)
		//case SensingMagara:
		//err = g.EvokeSensing(ev)
		//case DescentMagara:
		//err = g.EvokeDescent(ev)
	}
	if err != nil {
		return err
	}
	g.Stats.MagarasUsed++ // TODO
	// TODO: animation
	g.Player.MP -= mag.MPCost(g)
	g.StoryPrintf("You evoked your %s.", mag)
	ev.Renew(g, 5)
	return nil
}

func (mag magara) String() (desc string) {
	switch mag {
	case NoMagara:
		desc = "empty slot"
	case BlinkMagara:
		desc = "magara of blinking"
	case TeleportMagara:
		desc = "magara of teleportation"
	case DigMagara:
		desc = "magara of digging"
	case TeleportOtherMagara:
		desc = "magara of teleport other"
	case HealWoundsMagara:
		desc = "magara of heal wounds"
	//case MagicMagara:
	//desc = "magara of refill magic"
	case SwiftnessMagara:
		desc = "magara of swiftness"
	case SwappingMagara:
		desc = "magara of swapping"
	case ShadowsMagara:
		desc = "magara of shadows"
	case FogMagara:
		desc = "magara of fog"
	case WallsMagara:
		desc = "magara of walls"
	case SlowingMagara:
		desc = "magara of slowing"
	case SleepingMagara:
		desc = "magara of sleeping"
	case NoiseMagara:
		desc = "magara of noise"
	case ObstructionMagara:
		desc = "magara of obstruction"
	case ConfusionMagara:
		desc = "magara of confusion"
	case FireMagara:
		desc = "magara of fire"
	case LignificationMagara:
		desc = "magara of lignification"
		//case DescentMagara:
		//desc = "magara of descent"
		//case MagicMappingMagara:
		//desc = "magara of magic mapping"
		//case SensingMagara:
		//desc = "magara of sensing"
	}
	return desc
}

func (mag magara) Desc(g *game) (desc string) {
	// TODO
	switch mag {
	case NoMagara:
		desc = "can be used for a new magara."
	case BlinkMagara:
		desc = "makes you blink away within your line of sight. The rod is more susceptible to send you to the cells thar are most far from you."
	case TeleportMagara:
		desc = "makes you teleport far away."
	case DigMagara:
		desc = "makes you dig walls by walking into them like an earth dragon."
	case TeleportOtherMagara:
		desc = "teleports up to two random monsters in sight."
	case HealWoundsMagara:
		desc = "heals you a good deal."
	//case MagicMagara:
	//desc = "replenishes your magical reserves."
	case SwiftnessMagara:
		desc = "makes you move faster and better at avoiding blows for a short time." // XXX
	case SwappingMagara:
		desc = "makes you swap positions with the farthest monster in sight. If there is more than one at the same distance, it will be chosen randomly."
	case ShadowsMagara:
		desc = "reduces your line of sight range to 1. Because monsters only can see you if you see them, this makes it easier to get out of sight of monsters so that they eventually stop chasing you."
	case FogMagara:
		desc = ""
	case WallsMagara:
		desc = "replaces free cells around you with temporary walls."
	case SlowingMagara:
		desc = "induces slow movement and attack for monsters in sight."
	case SleepingMagara:
		desc = "induces deep sleeping and exhaustion for up to two random monsters in sight."
	case NoiseMagara:
		desc = "tricks monsters in a 10-range area with sounds, making them go away from you for a few turns. It only works on monsters that are not already seeing you."
	case ObstructionMagara:
		desc = "creates temporal walls between you and up to 3 monsters."
	case ConfusionMagara:
		desc = "confuses monsters in sight, leaving them unable to attack you."
	case FireMagara:
		desc = "produces a small magical fire that will extend to neighbour flammable terrain. The smoke it generates will induce sleep in monsters. As a gawalt monkey, you resist sleepiness, but you will still feel slowed."
	case LignificationMagara:
		desc = "lignifies up to 2 monsters in view, so that it cannot move. The monster can still fight."
		//case DescentMagara:
		//desc = "makes you go deeper in the Underground."
		//case MagicMappingMagara:
		//desc = "shows you the map layout and item locations."
		//case SensingMagara:
		//desc = "shows you the current position of monsters in the map."
	}
	return fmt.Sprintf("The %s %s", mag, desc)
}

func (mag magara) MPCost(g *game) int {
	if mag == NoMagara {
		return 0
	}
	cost := 1
	switch mag {
	case HealWoundsMagara:
		cost = 2
	}
	if g.Player.HasStatus(StatusConfusion) {
		cost++
	}
	return cost
}

func (g *game) EvokeBlink(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	g.Blink(ev)
	return nil
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

func (g *game) BlinkPos() position {
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		// TODO: skip if not seen?
		if !b {
			continue
		}
		if !g.Dungeon.Cell(pos).IsPassable() {
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

func (g *game) EvokeTeleport(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	g.Teleportation(ev)
	g.Print("You teleported away.")
	return nil
}

func (g *game) EvokeDig(ev event) error {
	g.Player.Statuses[StatusDig] = 1
	end := ev.Rank() + DurationDigging
	g.PushEvent(&simpleEvent{ERank: end, EAction: DigEnd})
	g.Player.Expire[StatusDig] = end
	g.Print("You feel like an earth dragon.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) MonstersInLOS() []*monster {
	ms := []*monster{}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			ms = append(ms, mons)
		}
	}
	// shuffle before, because the order could be unnaturally predicted
	for i := 0; i < len(ms); i++ {
		j := i + RandInt(len(ms)-i)
		ms[i], ms[j] = ms[j], ms[i]
	}
	return ms
}

func (g *game) EvokeTeleportOther(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	for i := 0; i < max; i++ {
		ms[i].TeleportAway(g)
	}

	return nil
}

func (g *game) EvokeHealWounds(ev event) error {
	g.Player.HP = g.Player.HPMax()
	g.Print("Your feel healthy again.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

func (g *game) EvokeRefillMagic(ev event) error {
	g.Player.MP = g.Player.MPMax()
	g.Print("Your magic forces return.")
	g.ui.PlayerGoodEffectAnimation()
	return nil
}

//func (g *game) EvokeDescent(ev event) error {
//if g.Depth >= MaxDepth {
//return errors.New("You cannot descend any deeper!")
//}
//g.Printf("You fall through the ground.")
//g.LevelStats()
//g.StoryPrint("You descended deeper into the dungeon.")
//g.Depth++
//g.DepthPlayerTurn = 0
//g.InitLevel()
//g.Save()
//return nil
//}

func (g *game) EvokeSwiftness(ev event) error {
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + DurationSwiftness
	g.PushEvent(&simpleEvent{ERank: end, EAction: HasteEnd})
	g.Player.Expire[StatusSwift] = end
	g.Printf("You feel speedy and agile.")
	g.ui.PlayerGoodEffectAnimation()
	// XXX do something with agile?
	return nil
}

func (g *game) EvokeSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	ms := g.MonstersInLOS()
	var mons *monster
	best := 0
	for _, m := range ms {
		if m.Status(MonsLignified) {
			continue
		}
		if m.Pos.Distance(g.Player.Pos) > best {
			best = m.Pos.Distance(g.Player.Pos)
			mons = m
		}
	}
	if !mons.Exists() {
		return errors.New("No monsters suitable for swapping in view.")
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

func (g *game) EvokeShadows(ev event) error {
	if g.Player.HasStatus(StatusShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Player.Statuses[StatusShadows] = 1
	end := ev.Rank() + DurationShadows
	g.PushEvent(&simpleEvent{ERank: end, EAction: ShadowsEnd})
	g.Player.Expire[StatusShadows] = end
	g.Printf("You feel surrounded by shadows.")
	g.ui.PlayerGoodEffectAnimation()
	g.ComputeLOS()
	return nil
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
	CloudNight
)

func (g *game) EvokeFog(ev event) error {
	g.Fog(g.Player.Pos, 3, ev)
	g.Print("You are surrounded by a dense fog.")
	return nil
}

func (g *game) Fog(at position, radius int, ev event) {
	dij := &noisePath{game: g}
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

func (g *game) EvokeWalls(ev event) error {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		g.CreateTemporalWallAt(pos, ev)
	}
	g.Print("You feel surrounded by temporary walls.")
	g.ui.PlayerGoodEffectAnimation()
	g.ComputeLOS()
	return nil
}

func (g *game) EvokeSlowing(ev event) error {
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.Statuses[MonsSlow]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + DurationSlow, NMons: mons.Index, EAction: MonsSlowEnd})
	}
	g.Print("Whoosh! A slowing luminous wave emerges.")
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveLOS)

	return nil
}

func (g *game) EvokeSleeping(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	targets := []position{}
	// XXX: maybe use noise distance instead of LOS?
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.State != Resting {
			g.Printf("%s falls asleep.", mons.Kind.Definite(true))
		}
		mons.State = Resting
		mons.Dir = NoDir
		mons.ExhaustTime(g, 40+RandInt(10))
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if max == 1 {
		g.Print("A beam of sleeping emerges.")
	} else {
		g.Print("Two beams of sleeping emerge.")
	}
	g.ui.BeamsAnimation(targets)

	return nil
}

func (g *game) EvokeLignification(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 2
	if max > len(ms) {
		max = len(ms)
	}
	targets := []position{}
	for i := 0; i < max; i++ {
		mons := ms[i]
		if mons.Status(MonsLignified) {
			continue
		}
		mons.EnterLignification(g, ev)
		targets = append(targets, g.Ray(mons.Pos)...)
	}
	if max == 1 {
		g.Print("A beam of lignification emerges.")
	} else {
		g.Print("Two beams of lignification emerge.")
	}
	g.ui.BeamsAnimation(targets)
	return nil
}

func (g *game) EvokeNoise(ev event) error {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, 23)
	nmlosr := Dijkstra(dij, []position{g.Player.Pos}, DefaultLOSRange)
	noises := []position{}
	g.NoiseIllusion = map[position]bool{}
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			continue
		}
		_, ok := nmlosr[mons.Pos]
		if !ok {
			continue
		}
		if mons.SeesPlayer(g) || mons.Kind == MonsSatowalgaPlant {
			continue
		}
		mp := &monPath{game: g, monster: mons}
		target := mons.Pos
		best := nm[target].Cost
		for {
			ncost := best
			for _, pos := range mp.Neighbors(target) {
				node, ok := nm[pos]
				if !ok {
					continue
				}
				ncost := node.Cost
				if ncost > best {
					target = pos
					best = ncost
				}
			}
			if ncost == best {
				break
			}
		}
		if mons.State != Hunting {
			mons.State = Wandering
		}
		mons.Target = target
		noises = append(noises, target)
		g.NoiseIllusion[target] = true
	}
	g.ui.NoiseAnimation(noises)
	g.Print("Monsters are tricked by magical sounds.")
	return nil
}

func (g *game) EvokeConfusion(ev event) error {
	g.Print("Whoosh! A confusing luminous wave emerges.")
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		mons.Statuses[MonsConfused]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + DurationConfusion, NMons: mons.Index, EAction: MonsConfusionEnd})
	}
	g.ui.LOSWavesAnimation(DefaultLOSRange, WaveLOS)
	return nil
}

func (g *game) EvokeFire(ev event) error {
	if !g.Dungeon.Cell(g.Player.Pos).Flammable() {
		return errors.New("You are not standing on flammable terrain.")
	}
	g.Burn(g.Player.Pos, ev)
	g.Print("A small fire appears.")
	return nil
}

func (g *game) EvokeObstruction(ev event) error {
	ms := g.MonstersInLOS()
	if len(ms) == 0 {
		return errors.New("There are no monsters in view.")
	}
	max := 3
	if max > len(ms) {
		max = len(ms)
	}
	targets := []position{}
	for i := 0; i < max; i++ {
		ray := g.Ray(ms[i].Pos)
		for i, pos := range ray[1:] {
			if pos == g.Player.Pos {
				break
			}
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			g.TemporalWallAt(pos, ev)
			if len(ray) == 0 {
				break
			}
			ray = ray[i+1:]
			targets = append(targets, ray...)
			break
		}
	}
	g.Print("Magical walls emerged.")
	g.ui.BeamsAnimation(targets)
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
	t := g.Dungeon.Cell(pos).T
	g.Dungeon.SetCell(pos, WallCell)
	delete(g.Clouds, pos)
	g.TemporalWalls[pos] = t
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationTemporalWall + RandInt(DurationTemporalWall/2), Pos: pos, EAction: ObstructionEnd})
}
