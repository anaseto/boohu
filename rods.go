package main

import (
	"errors"
	"fmt"
)

type rod int

const (
	RodDigging rod = iota
	RodBlink
	RodTeleportOther
	RodLightningBolt
	RodFireball
	RodFog
	RodObstruction
	RodShatter
	RodSwapping
	// below unimplemented
	RodConfusingClouds
	RodFear
	RodFreezingClouds
)

const NumRods = int(RodSwapping) + 1

func (r rod) Letter() rune {
	return '/'
}

func (r rod) Rare() bool {
	switch r {
	case RodDigging, RodTeleportOther, RodShatter, RodSwapping:
		return true
	default:
		return false
	}
}

func (r rod) String() string {
	var text string
	switch r {
	case RodDigging:
		text = "rod of digging"
	case RodBlink:
		text = "rod of blinking"
	case RodTeleportOther:
		text = "rod of teleport other"
	case RodFog:
		text = "rod of fog"
	case RodFear:
		text = "rod of fear"
	case RodFireball:
		text = "rod of fireball"
	case RodLightningBolt:
		text = "rod of lightning bolt"
	case RodObstruction:
		text = "rod of obstruction"
	case RodShatter:
		text = "rod of shatter"
	case RodSwapping:
		text = "rod of swapping"
	case RodFreezingClouds:
		text = "rod of freezing clouds"
	case RodConfusingClouds:
		text = "rod of confusing clouds"
	}
	return text
}

func (r rod) Desc() string {
	var text string
	switch r {
	case RodDigging:
		text = "digs through walls."
	case RodBlink:
		text = "makes you blink away within your line of sight."
	case RodTeleportOther:
		text = "teleports away one of your foes."
	case RodFog:
		text = "creates a dense fog that reduces your (and monster's) line of sight."
	case RodFireball:
		text = "throws a 1-radius fireball at your foes."
	case RodLightningBolt:
		text = "throws a lightning bolt through one or more enemies."
	case RodObstruction:
		text = "creates a temporary wall at targeted location."
	case RodShatter:
		text = "induces an explosion around a wall. The wall can disintegrate."
	case RodSwapping:
		text = "makes you swap positions with a targeted monster."
	case RodFear:
		text = "TODO"
	case RodFreezingClouds:
		text = "TODO"
	case RodConfusingClouds:
		text = "TODO"
	}
	return fmt.Sprintf("The %s %s Rods sometimes regain charges as you go deeper. This rod can have up to %d charges.", r, text, r.MaxCharge())
}

type rodProps struct {
	Charge int
}

func (r rod) MaxCharge() (charges int) {
	switch r {
	case RodBlink:
		charges = 5
	case RodTeleportOther, RodDigging:
		charges = 3
	default:
		charges = 4
	}
	return charges
}

func (r rod) Rate() int {
	rate := r.MaxCharge() - 2
	if rate < 1 {
		rate = 1
	}
	return rate
}

func (r rod) MPCost() (mp int) {
	switch r {
	case RodBlink:
		mp = 3
	case RodTeleportOther, RodDigging, RodShatter:
		mp = 5
	default:
		mp = 4
	}
	return mp
}

func (r rod) Use(g *game, ev event) error {
	rods := g.Player.Rods
	if rods[r].Charge <= 0 {
		return errors.New("No charges remaining on this rod.")
	}
	if r.MPCost() > g.Player.MP {
		return errors.New("Not enough magic points for using this rod.")
	}
	var err error
	switch r {
	case RodBlink:
		err = g.EvokeRodBlink(ev)
	case RodTeleportOther:
		err = g.EvokeRodTeleportOther(ev)
	case RodLightningBolt:
		err = g.EvokeRodLightningBolt(ev)
	case RodFireball:
		err = g.EvokeRodFireball(ev)
	case RodFog:
		err = g.EvokeRodFog(ev)
	case RodDigging:
		err = g.EvokeRodDigging(ev)
	case RodObstruction:
		err = g.EvokeRodObstruction(ev)
	case RodShatter:
		err = g.EvokeRodShatter(ev)
	case RodSwapping:
		err = g.EvokeRodSwapping(ev)
	}

	if err != nil {
		return err
	}
	rods[r].Charge--
	g.Player.MP -= r.MPCost()
	g.StoryPrintf("You evoked your %s.", r)
	g.FairAction()
	ev.Renew(g, 7)
	return nil
}

func (g *game) EvokeRodBlink(ev event) error {
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
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		if g.Dungeon.Cell(pos).T != FreeCell {
			continue
		}
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		losPos = append(losPos, pos)
	}
	if len(losPos) == 0 {
		// should not happen
		g.Print("You could not blink.")
		return
	}
	npos := losPos[RandInt(len(losPos))]
	if npos.Distance(g.Player.Pos) <= 3 {
		// Give close cells less chance to make blinking more useful
		npos = losPos[RandInt(len(losPos))]
	}
	g.Player.Pos = npos
	g.Print("You blink away.")
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
}

func (g *game) EvokeRodTeleportOther(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{}) {
		return errors.New(DoNothing)
	}
	mons, _ := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	mons.TeleportAway(g)
	return nil
}

func (g *game) EvokeRodLightningBolt(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{flammable: true}) {
		return errors.New(DoNothing)
	}
	ray := g.Ray(g.Player.Target)
	g.Print("A lightning bolt emerges straight from the rod.")
	g.ui.LightningBoltAnimation(g, ray)
	for _, pos := range ray {
		g.Burn(pos, ev)
		mons, _ := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		mons.HP -= RandInt(21)
		if mons.HP <= 0 {
			g.Printf("%s is killed by the bolt.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(12, mons.Pos)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodFireball(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{area: true, minDist: true, flammable: true}) {
		return errors.New(DoNothing)
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Print("A fireball emerges straight from the rod.")
	g.ui.ExplosionAnimation(g, FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.Burn(pos, ev)
		mons, _ := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		mons.HP -= RandInt(21)
		if mons.HP <= 0 {
			g.Printf("%s is killed by the fireball.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(12, mons.Pos)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

type cloud int

const (
	CloudFog cloud = iota
	CloudFire
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
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: CloudEnd, Pos: pos})
		}
	}
	g.ComputeLOS()
}

func (g *game) EvokeRodDigging(ev event) error {
	if !g.ui.ChooseTarget(g, &wallChooser{}) {
		return errors.New(DoNothing)
	}
	pos := g.Player.Target
	for i := 0; i < 3; i++ {
		g.Dungeon.SetCell(pos, FreeCell)
		g.MakeNoise(17, pos)
		g.Fog(pos, 1, ev)
		pos = pos.To(pos.Dir(g.Player.Pos))
		if !g.Player.LOS[pos] {
			g.UnknownDig[pos] = true
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
	if !g.ui.ChooseTarget(g, &wallChooser{minDist: true}) {
		return errors.New(DoNothing)
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	if RandInt(2) == 0 {
		g.Dungeon.SetCell(g.Player.Target, FreeCell)
		g.ComputeLOS()
		g.MakeMonstersAware()
		g.MakeNoise(19, g.Player.Target)
		g.Print("You see the wall disappear in a noisy explosion.")
		g.ui.ExplosionAnimation(g, WallExplosion, g.Player.Target)
		g.Fog(g.Player.Target, 2, ev)
	} else {
		g.MakeNoise(15, g.Player.Target)
		g.Print("You see an explosion around the wall.")
		g.ui.ExplosionAnimation(g, AroundWallExplosion, g.Player.Target)
	}
	for _, pos := range neighbors {
		mons, _ := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		mons.HP -= RandInt(30)
		if mons.HP <= 0 {
			g.Printf("%s is killed by the explosion.", mons.Kind.Indefinite(true))
			g.HandleKill(mons, ev)
		}
		g.MakeNoise(12, mons.Pos)
		mons.MakeHuntIfHurt(g)
	}
	return nil
}

func (g *game) EvokeRodObstruction(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{needsFreeWay: true, free: true}) {
		return errors.New(DoNothing)
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	for _, pos := range neighbors {
		g.MakeNoise(18, pos)
		break
	}
	g.Dungeon.SetCell(g.Player.Target, WallCell)
	delete(g.Clouds, g.Player.Target)
	if g.TemporalWalls != nil {
		g.TemporalWalls[g.Player.Target] = true
	}
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), Pos: g.Player.Target, EAction: ObstructionEnd})
	g.Printf("You see a wall appear from nothing.")
	g.ComputeLOS()
	return nil
}

func (g *game) EvokeRodSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	if !g.ui.ChooseTarget(g, &chooser{}) {
		return errors.New(DoNothing)
	}
	mons, _ := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	mons.Pos, g.Player.Pos = g.Player.Pos, mons.Pos
	mons.MakeAware(g)
	g.Printf("You swap positions with the %s.", mons.Kind)
	g.ui.SwappingAnimation(g, mons.Pos, g.Player.Pos)
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
	return nil
}

func (g *game) GeneratedRodsCount() int {
	count := 0
	for _, b := range g.GeneratedRods {
		if b {
			count++
		}
	}
	return count
}

func (g *game) GenerateRod() {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("GenerateRod")
		}
		pos := g.FreeCellForStatic()
		r := rod(RandInt(NumRods))
		if r.Rare() && RandInt(3) == 0 {
			r = rod(RandInt(NumRods))
		}
		if g.Player.Rods[r] == nil && !g.GeneratedRods[r] {
			g.GeneratedRods[r] = true
			g.Rods[pos] = r
			return
		}
	}
}
