package main

import (
	"errors"
	"fmt"
)

type player struct {
	HP          int
	HPbonus     int
	MP          int
	Simellas    int
	Dir         direction
	Armour      armour
	Weapon      weapon
	Consumables map[consumable]int
	Rods        map[rod]rodProps
	Aptitudes   map[aptitude]bool
	Statuses    map[status]int
	Expire      map[status]int
	Pos         position
	Target      position
	LOS         map[position]bool
	Rays        rayMap
	Bored       int
}

const DefaultHealth = 4

func (p *player) HPMax() int {
	hpmax := DefaultHealth
	if p.Aptitudes[AptHealthy] {
		hpmax += 1
	}
	hpmax -= p.Bored
	if p.Weapon == FinalBlade {
		hpmax -= 1
	}
	if hpmax < 2 {
		hpmax = 2
	}
	return hpmax
}

const DefaultMPmax = 3

func (p *player) MPMax() int {
	mpmax := DefaultMPmax
	if p.Aptitudes[AptMagic] {
		mpmax += 2
	}
	if p.Armour == CelmistRobe {
		mpmax += 2
	}
	return mpmax
}

func (p *player) Attack() int {
	attack := p.Weapon.Attack()
	if p.HasStatus(StatusCorrosion) {
		penalty := p.Statuses[StatusCorrosion]
		if penalty > 2 {
			penalty = 2
		}
		attack -= penalty
	}
	if attack <= 0 {
		attack = 0
	}
	return attack
}

func (p *player) HasStatus(st status) bool {
	return p.Statuses[st] > 0
}

func (p *player) AptitudeCount() int {
	count := 0
	for _, b := range p.Aptitudes {
		if b {
			count++
		}
	}
	return count
}

func (g *game) AutoToDir(ev event) bool {
	if g.MonsterInLOS() == nil {
		err := g.MovePlayer(g.Player.Pos.To(g.AutoDir), ev)
		if err != nil {
			g.Print(err.Error())
			g.AutoDir = NoDir
			return false
		}
		return true
	}
	g.AutoDir = NoDir
	return false
}

func (g *game) GoToDir(dir direction, ev event) error {
	if g.MonsterInLOS() != nil {
		g.AutoDir = NoDir
		return errors.New("You cannot travel while there are monsters in view.")
	}
	err := g.MovePlayer(g.Player.Pos.To(dir), ev)
	if err != nil {
		return err
	}
	g.AutoDir = dir
	return nil
}

func (g *game) MoveToTarget(ev event) bool {
	if !g.AutoTarget.valid() {
		return false
	}
	path := g.PlayerPath(g.Player.Pos, g.AutoTarget)
	if g.MonsterInLOS() != nil {
		g.AutoTarget = InvalidPos
	}
	if len(path) < 1 {
		g.AutoTarget = InvalidPos
		return false
	}
	var err error
	if len(path) > 1 {
		err = g.MovePlayer(path[len(path)-2], ev)
		if g.ExclusionsMap[path[len(path)-2]] {
			g.AutoTarget = InvalidPos
		}
	} else {
		g.WaitTurn(ev)
	}
	if err != nil {
		g.Print(err.Error())
		g.AutoTarget = InvalidPos
		return false
	}
	if g.AutoTarget.valid() && g.Player.Pos == g.AutoTarget {
		g.AutoTarget = InvalidPos
	}
	return true
}

func (g *game) WaitTurn(ev event) {
	grade := 1
	if len(g.Noise) > 0 || g.StatusRest() {
		grade = 1
	}
	g.BoredomAction(ev, grade)
	ev.Renew(g, 10)
}

func (g *game) MonsterCount() (count int) {
	for _, mons := range g.Monsters {
		if mons.Exists() {
			count++
		}
	}
	return count
}

func (g *game) BoredomAction(ev event, grade int) {
	obor := g.Boredom
	if g.MonsterInLOS() == nil {
		g.Boredom += grade
	} else {
		g.Boredom--
		if g.Boredom < 0 {
			g.Boredom = 0
			g.Player.Bored = 0
		}
		return
	}
	if g.MonsterCount() <= 4 {
		return
	}
	if g.Boredom >= 500 && obor < 500 {
		g.PrintStyled("You feel a little bored, your health may decline.", logCritic)
		g.StopAuto()
	}
	if g.Boredom >= 510 && (obor/10 != g.Boredom/10) {
		g.Player.Bored++
		g.PrintStyled("You feel unhealthy.", logCritic)
		g.StopAuto()
		if g.Player.HP > g.Player.HPMax() {
			g.Player.HP--
		}
	}
}

func (g *game) FunAction() {
	g.Boredom -= 15
	if g.Boredom < 0 {
		g.Boredom = 0
		g.Player.Bored = 0
	}
}

func (g *game) Rest(ev event) error {
	if g.Dungeon.Cell(g.Player.Pos).T != BarrelCell {
		return fmt.Errorf("This place is not safe for sleeping.")
	}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		return errors.New("You cannot rest on flames.")
	}
	if !g.NeedsRegenRest() && !g.StatusRest() {
		return errors.New("You do not need to rest.")
	}
	g.WaitTurn(ev)
	g.Resting = true
	g.RestingTurns = 0
	if g.StatusRest() {
		g.RestingTurns = -1 // not true resting, just waiting for status end
	}
	g.FunAction()
	return nil
}

func (g *game) StatusRest() bool {
	for _, q := range g.Player.Statuses {
		if q > 0 {
			return true
		}
	}
	return false
}

func (g *game) NeedsRegenRest() bool {
	return g.Player.HP < g.Player.HPMax() || g.Player.MP < g.Player.MPMax()
}

func (g *game) Equip(ev event) error {
	obj, ok := g.Objects[g.Player.Pos]
	if !ok {
		return errors.New("Found nothing to equip here.")
	}
	if eq, ok := obj.(equipable); ok {
		eq.Equip(g)
		ev.Renew(g, 10)
		g.BoredomAction(ev, 1)
		return nil
	}
	return errors.New("Found nothing to equip here.")
}

func (g *game) Teleportation(ev event) {
	var pos position
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("Teleportation")
		}
		pos = g.FreeCell()
		if pos.Distance(g.Player.Pos) < 15 && i < 1000 {
			i++
			continue
		}
		break

	}
	if pos.valid() {
		// should always happen
		opos := g.Player.Pos
		g.Print("You teleport away.")
		g.ui.TeleportAnimation(opos, pos, true)
		g.PlacePlayerAt(pos)
	} else {
		// should not happen
		g.Print("Something went wrong with the teleportation.")
	}
}

func (g *game) CollectGround() {
	pos := g.Player.Pos
	if obj, ok := g.Objects[pos]; ok {
		g.DijkstraMapRebuild = true
		switch o := obj.(type) {
		case rod:
			delete(g.Objects, pos)
			g.Player.Rods[o] = rodProps{Charge: o.MaxCharge() - 1}
			g.Printf("You take %s.", obj.ShortDesc(g))
			g.StoryPrintf("You found and took %s.", obj.ShortDesc(g))
		case collectable:
			delete(g.Objects, pos)
			g.Player.Consumables[o.Consumable] += o.Quantity
			if o.Quantity > 1 {
				g.Printf("You take %d %s.", o.Quantity, o.Consumable.Plural())
			} else {
				g.Printf("You take %s.", obj.ShortDesc(g))
			}
		case simella:
			delete(g.Objects, pos)
			g.Player.Simellas += int(o)
			if o == simella(1) {
				g.Print("You pick up a simella.")
			} else {
				g.Printf("You pick up %d simellas.", o)
			}
			g.DijkstraMapRebuild = true
		default:
			g.Printf("You are standing over %s.", obj.ShortDesc(g))
		}
	} else if g.Dungeon.Cell(pos).T == DoorCell {
		g.Print("You stand at the door.")
	}
}

func (g *game) MovePlayer(pos position, ev event) error {
	//if g.Player.Dir != pos.Dir(g.Player.Pos) {
	//g.Player.Dir = pos.Dir(g.Player.Pos)
	//ev.Renew(g, 5)
	//g.ComputeLOS() // TODO: not really needed
	//return nil
	//}
	if !pos.valid() {
		return errors.New("You cannot move there.")
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot move into a wall.")
	} else if c.T == BarrelCell && g.MonsterLOS[g.Player.Pos] {
		return errors.New("You cannot enter a barrel while seen.")
	}
	if g.Player.HasStatus(StatusConfusion) {
		switch pos.Dir(g.Player.Pos) {
		case E, N, W, S:
		default:
			return errors.New("You cannot use diagonal movements while confused.")
		}
	}
	delay := 10
	mons := g.MonsterAt(pos)
	if !mons.Exists() {
		if g.Player.HasStatus(StatusLignification) {
			return errors.New("You cannot move while lignified")
		}
		if c.T == BarrelCell {
			g.Print("You hide yourself inside the barrel.")
		}
		if c.T == WallCell {
			g.Dungeon.SetCell(pos, GroundCell)
			g.MakeNoise(WallNoise, pos)
			g.Print(g.CrackSound())
			g.Fog(pos, 1, ev)
			g.Stats.Digs++
		}
		switch g.Player.Armour {
		case SmokingScales:
			_, ok := g.Clouds[g.Player.Pos]
			if !ok {
				g.Clouds[g.Player.Pos] = CloudFog
				g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationSmokingScalesFog, EAction: CloudEnd, Pos: g.Player.Pos})
			}
		}
		g.Stats.Moves++
		g.PlacePlayerAt(pos)
		if !g.Autoexploring {
			g.BoredomAction(ev, 1)
		}
	} else if err := g.Jump(mons, ev); err != nil {
		return err
	}
	if g.Player.HasStatus(StatusSwift) {
		// only fast for movement
		delay /= 2
	}
	if g.Player.HasStatus(StatusSlow) {
		delay *= 2
	}
	if delay < 5 {
		delay = 5
	} else if delay > 20 {
		delay = 20
	}
	ev.Renew(g, delay)
	return nil
}

func (g *game) HealPlayer(ev event) {
	if g.Player.HP < g.Player.HPMax() {
		g.Player.HP++
	}
	delay := 50
	ev.Renew(g, delay)
}

func (g *game) MPRegen(ev event) {
	if g.Player.MP < g.Player.MPMax() {
		g.Player.MP++
	}
	delay := 100
	ev.Renew(g, delay)
}

func (g *game) Smoke(ev event) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, 2)
	for pos := range nm {
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudFog
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationFog + RandInt(DurationFog/2), EAction: CloudEnd, Pos: pos})
		}
	}
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + DurationShortSwiftness
	g.PushEvent(&simpleEvent{ERank: end, EAction: HasteEnd})
	g.Player.Expire[StatusSwift] = end
	g.ComputeLOS()
	g.Print("You feel an energy burst and smoke comes out from you.")
}

func (g *game) Corrosion(ev event) {
	g.Player.Statuses[StatusCorrosion]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 80 + RandInt(40), EAction: CorrosionEnd})
	g.Print("Your equipment gets corroded.")
}

func (g *game) Confusion(ev event) {
	if !g.Player.HasStatus(StatusConfusion) {
		g.Player.Statuses[StatusConfusion]++
		g.PushEvent(&simpleEvent{ERank: ev.Rank() + DurationConfusion + RandInt(DurationConfusion/2), EAction: ConfusionEnd})
		g.Print("You feel confused.")
	}
}

func (g *game) PlacePlayerAt(pos position) {
	g.Player.Dir = pos.Dir(g.Player.Pos)
	switch g.Player.Dir {
	case ENE, ESE:
		g.Player.Dir = E
	case NNE, NNW:
		g.Player.Dir = N
	case WNW, WSW:
		g.Player.Dir = W
	case SSW, SSE:
		g.Player.Dir = S
	}
	g.Player.Pos = pos
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
}

func (g *game) EnterLignification(ev event) {
	g.Player.Statuses[StatusLignification]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + DurationLignification + RandInt(DurationLignification/2), EAction: LignificationEnd})
	g.Player.HPbonus += 4
}
