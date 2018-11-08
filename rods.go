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
	RodFireBolt
	RodFireBall
	RodLightning
	RodFog
	RodObstruction
	RodShatter
	RodSleeping
	RodLignification
	RodHope
	RodSwapping
)

const NumRods = int(RodSwapping) + 1

func (r rod) Letter() rune {
	return '/'
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
	case RodFireBall:
		text = "rod of fireball"
	case RodFireBolt:
		text = "rod of fire bolt"
	case RodLightning:
		text = "rod of lightning"
	case RodObstruction:
		text = "rod of obstruction"
	case RodShatter:
		text = "rod of shatter"
	case RodSleeping:
		text = "rod of sleeping"
	case RodLignification:
		text = "rod of lignification"
	case RodHope:
		text = "rod of last hope"
	case RodSwapping:
		text = "rod of swapping"
	}
	return text
}

func (r rod) Desc() string {
	var text string
	switch r {
	case RodDigging:
		text = "digs through up to 3 walls in a given direction."
	case RodBlink:
		text = "makes you blink away within your line of sight."
	case RodTeleportOther:
		text = "teleports away one of your foes. Note that the monster remembers where it saw you last time."
	case RodFog:
		text = "creates a dense fog that reduces your (and monster's) line of sight."
	case RodFireBall:
		text = "throws a 1-radius fireball at your foes. You cannot use it at melee range. It can burn foliage and doors."
	case RodFireBolt:
		text = "throws a fire bolt through one or more enemies. It can burn foliage and doors."
	case RodLightning:
		text = "deals electrical damage to foes connected to you. It can burn foliage and doors."
	case RodObstruction:
		text = "creates a temporary wall at targeted location."
	case RodShatter:
		text = "induces an explosion around a wall, hurting adjacent monsters. The wall can disintegrate. You cannot use it at melee range."
	case RodSleeping:
		text = "induces sleeping and exhaustion for monsters in the targeted area."
	case RodLignification:
		text = "lignifies a monster, so that it cannot move, but can still fight with improved resistance."
	case RodHope:
		text = "creates an energy channel against a targeted monster. The damage done is inversely proportional to your health. It can burn foliage and doors."
	case RodSwapping:
		text = "makes you swap positions with a targeted monster."
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
	case RodDigging, RodShatter:
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
	return 1
	//switch r {
	//case RodBlink:
	//mp = 3
	//case RodTeleportOther, RodDigging, RodShatter:
	//mp = 5
	//default:
	//mp = 4
	//}
	//return mp
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
	case RodFireBolt:
		err = g.EvokeRodFireBolt(ev)
	case RodFireBall:
		err = g.EvokeRodFireball(ev)
	case RodLightning:
		err = g.EvokeRodLightning(ev)
	case RodFog:
		err = g.EvokeRodFog(ev)
	case RodDigging:
		err = g.EvokeRodDigging(ev)
	case RodObstruction:
		err = g.EvokeRodObstruction(ev)
	case RodShatter:
		err = g.EvokeRodShatter(ev)
	case RodSleeping:
		err = g.EvokeRodSleeping(ev)
	case RodLignification:
		err = g.EvokeRodLignification(ev)
	case RodHope:
		err = g.EvokeRodHope(ev)
	case RodSwapping:
		err = g.EvokeRodSwapping(ev)
	}

	if err != nil {
		return err
	}
	rp := rods[r]
	rp.Charge--
	rods[r] = rp
	g.Player.MP -= r.MPCost()
	g.StoryPrintf("You evoked your %s.", r)
	g.Stats.UsedRod[r]++
	g.Stats.Evocations++
	g.FunAction()
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

func (g *game) BlinkPos() position {
	losPos := []position{}
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		if g.Dungeon.Cell(pos).T != FreeCell {
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
	if npos == InvalidPos {
		// should not happen
		g.Print("You could not blink.")
		return
	}
	opos := g.Player.Pos
	g.Print("You blink away.")
	g.ui.TeleportAnimation(g, opos, npos, true)
	g.PlacePlayerAt(npos)
}

func (g *game) EvokeRodTeleportOther(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	mons.TeleportAway(g)
	return nil
}

func (g *game) EvokeRodSleeping(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{area: true, minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Print("A sleeping ball emerges straight out of the rod.")
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgSleepingMonster)
	for _, pos := range append(neighbors, g.Player.Target) {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		mons.State = Resting
		if !mons.Status(MonsExhausted) {
			mons.Statuses[MonsExhausted] = 1
			g.PushEvent(&monsterEvent{ERank: ev.Rank() + 30 + RandInt(20), NMons: mons.Index, EAction: MonsExhaustionEnd})
		}
	}
	return nil
}

func (g *game) EvokeRodFireBolt(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{flammable: true}); err != nil {
		return err
	}
	ray := g.Ray(g.Player.Target)
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.Print("Whoosh! A fire bolt emerges straight out of the rod.")
	g.ui.LightningBoltAnimation(g, ray)
	for _, pos := range ray {
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(21)
		}
		dmg /= 2
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
	if err := g.ui.ChooseTarget(g, &chooser{area: true, minDist: true, flammable: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.MakeNoise(MagicExplosionNoise, g.Player.Target)
	g.Printf("A fireball emerges straight out of the rod... %s", g.ExplosionSound())
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgExplosionStart)
	g.ui.ExplosionAnimation(g, FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(21)
		}
		dmg /= 2
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
		dmg := 0
		for i := 0; i < 2; i++ {
			dmg += RandInt(17)
		}
		dmg /= 2
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
	g.ui.LightningHitAnimation(g, targets)

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
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: CloudEnd, Pos: pos})
		}
	}
	g.ComputeLOS()
}

func (g *game) EvokeRodDigging(ev event) error {
	if err := g.ui.ChooseTarget(g, &wallChooser{}); err != nil {
		return err
	}
	pos := g.Player.Target
	for i := 0; i < 3; i++ {
		g.Dungeon.SetCell(pos, FreeCell)
		g.Stats.Digs++
		g.MakeNoise(WallNoise, pos)
		g.Fog(pos, 1, ev)
		pos = pos.To(pos.Dir(g.Player.Pos))
		if !g.Player.LOS[pos] {
			g.WrongWall[pos] = true
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
	if err := g.ui.ChooseTarget(g, &wallChooser{minDist: true}); err != nil {
		return err
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Dungeon.SetCell(g.Player.Target, FreeCell)
	g.Stats.Digs++
	g.ComputeLOS()
	g.MakeMonstersAware()
	g.MakeNoise(WallNoise, g.Player.Target)
	g.Printf("%s The wall disappeared.", g.CrackSound())
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgExplosionWallStart)
	g.ui.ExplosionAnimation(g, WallExplosion, g.Player.Target)
	g.Fog(g.Player.Target, 2, ev)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		dmg := 0
		for i := 0; i < 3; i++ {
			dmg += RandInt(30)
		}
		dmg /= 3
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
	if err := g.ui.ChooseTarget(g, &chooser{free: true}); err != nil {
		return err
	}
	g.TemporalWallAt(g.Player.Target, ev)
	g.Printf("You see a wall appear out of thin air.")
	return nil
}

func (g *game) EvokeRodLignification(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{}); err != nil {
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
	if !g.Player.LOS[pos] {
		g.WrongWall[pos] = true
	}
	g.CreateTemporalWallAt(pos, ev)
	g.ComputeLOS()
}

func (g *game) CreateTemporalWallAt(pos position, ev event) {
	g.Dungeon.SetCell(pos, WallCell)
	delete(g.Clouds, pos)
	g.TemporalWalls[pos] = true
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), Pos: pos, EAction: ObstructionEnd})
}

func (g *game) EvokeRodHope(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{needsFreeWay: true}); err != nil {
		return err
	}
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgExplosionStart)
	mons := g.MonsterAt(g.Player.Target)
	// mons not nil (check done in the targeter)
	attack := -20 + 30*DefaultHealth/g.Player.HP
	if attack > 130 {
		attack = 130
	}
	dmg := 0
	for i := 0; i < 5; i++ {
		dmg += RandInt(attack)
	}
	dmg /= 5
	if dmg < 0 {
		// should not happen
		dmg = 0
	}
	mons.HP -= dmg
	g.ui.HitAnimation(g, g.Player.Target, true)
	g.Printf("An energy channel hits %s (%d dmg).", mons.Kind.Definite(false), dmg)
	g.Burn(g.Player.Target, ev)
	return nil
}

func (g *game) EvokeRodSwapping(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot use this rod while lignified.")
	}
	if err := g.ui.ChooseTarget(g, &chooser{}); err != nil {
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
	g.ui.SwappingAnimation(g, mons.Pos, g.Player.Pos)
	mons.MoveTo(g, g.Player.Pos)
	g.PlacePlayerAt(ompos)
	mons.MakeAware(g)
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

func (g *game) RandomRod() rod {
	r := rod(RandInt(NumRods))
	return r
}

func (g *game) GenerateRod() {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("GenerateRod")
		}
		pos := g.FreeCellForStatic()
		r := g.RandomRod()
		if _, ok := g.Player.Rods[r]; !ok && !g.GeneratedRods[r] {
			g.GeneratedRods[r] = true
			g.Rods[pos] = r
			return
		}
	}
}

func (g *game) RechargeRods() {
	for r, props := range g.Player.Rods {
		max := r.MaxCharge()
		if g.Player.Armour == CelmistRobe {
			max++
		}
		if props.Charge < max {
			rchg := RandInt(1 + r.Rate())
			if rchg == 0 && RandInt(2) == 0 {
				rchg++
			}
			if g.Player.Armour == CelmistRobe {
				if RandInt(5) > 0 {
					rchg++
				}
				if RandInt(3) == 0 {
					rchg++
				}
			}
			props.Charge += rchg
			g.Player.Rods[r] = props
		}
		if props.Charge > max {
			props.Charge = max
			g.Player.Rods[r] = props
		}
	}
}
