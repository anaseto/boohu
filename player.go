package main

import (
	"errors"
	"fmt"
)

type player struct {
	HP          int
	MP          int
	Simellas    int
	Armour      armour
	Weapon      weapon
	Shield      shield
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

const DefaultHealth = 42

func (p *player) HPMax() int {
	hpmax := DefaultHealth
	if p.Aptitudes[AptHealthy] {
		hpmax += 10
	}
	hpmax -= 3 * p.Bored
	if p.Weapon == FinalBlade {
		hpmax = 2 * hpmax / 3
	}
	if hpmax < 21 {
		hpmax = 21
	}
	return hpmax
}

func (p *player) MPMax() int {
	mpmax := 3
	if p.Aptitudes[AptMagic] {
		mpmax += 2
	}
	if p.Armour == CelmistRobe {
		mpmax += 2
	}
	return mpmax
}

func (p *player) Accuracy() int {
	acc := 15
	return acc
}

func (p *player) RangedAccuracy() int {
	acc := 15
	return acc
}

func (p *player) Armor() int {
	ar := 0
	switch p.Armour {
	case LeatherArmour:
		ar += 3
	case SmokingScales:
		ar += 4
	case ShinyPlates:
		ar += 6
	case TurtlePlates:
		ar += 9
	}
	if p.Aptitudes[AptScales] {
		ar += 2
	}
	if p.HasStatus(StatusLignification) {
		ar = 9 + ar/2
	}
	if p.HasStatus(StatusCorrosion) {
		ar -= 2 * p.Statuses[StatusCorrosion]
		if ar < 0 {
			ar = 0
		}
	}
	return ar
}

func (p *player) Attack() int {
	attack := p.Weapon.Attack()
	if p.Aptitudes[AptStrong] {
		attack += attack / 5
	}
	if p.HasStatus(StatusCorrosion) {
		penalty := p.Statuses[StatusCorrosion]
		if penalty > 5 {
			penalty = 5
		}
		attack -= penalty
	}
	return attack
}

func (p *player) Block() int {
	block := p.Shield.Block()
	if p.HasStatus(StatusDisabledShield) {
		block /= 3
	}
	return block
}

func (p *player) Evasion() int {
	ev := 15
	if p.Aptitudes[AptAgile] {
		ev += 3
	}
	switch p.Armour {
	case ShinyPlates:
		ev -= 1
	case TurtlePlates:
		ev -= 2
	case HarmonistRobe:
		ev += 1
	}
	if p.HasStatus(StatusAgile) {
		ev += 7
	}
	return ev
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
	// XXX Really wait for 10 ?
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
	if g.Boredom >= 120 && obor < 120 {
		if g.MonsterCount() > 4 {
			g.PrintStyled("You feel a little bored, your health may decline.", logCritic)
			g.StopAuto()
		}
	}
	if g.Boredom >= 130 && (obor/10 != g.Boredom/10) {
		if g.MonsterCount() > 4 {
			g.Player.Bored++
			g.PrintStyled("You feel unhealthy.", logCritic)
			g.StopAuto()
			if g.Player.HP > g.Player.HPMax() {
				g.Player.HP -= 3
			}
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
	if g.MonsterInLOS() != nil {
		return fmt.Errorf("You cannot sleep while monsters are in view.")
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
	if eq, ok := g.Equipables[g.Player.Pos]; ok {
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
		g.ui.TeleportAnimation(g, opos, pos, true)
		g.PlacePlayerAt(pos)
	} else {
		// should not happen
		g.Print("Something went wrong with the teleportation.")
	}
}

func (g *game) CollectGround() {
	pos := g.Player.Pos
	if g.Simellas[pos] > 0 {
		g.Player.Simellas += g.Simellas[pos]
		if g.Simellas[pos] == 1 {
			g.Print("You pick up a simella.")
		} else {
			g.Printf("You pick up %d simellas.", g.Simellas[pos])
		}
		g.DijkstraMapRebuild = true
		delete(g.Simellas, pos)
	}
	if c, ok := g.Collectables[pos]; ok {
		g.Player.Consumables[c.Consumable] += c.Quantity
		g.DijkstraMapRebuild = true
		delete(g.Collectables, pos)
		if c.Quantity > 1 {
			g.Printf("You take %d %s.", c.Quantity, c.Consumable.Plural())
		} else {
			g.Printf("You take %s.", Indefinite(c.Consumable.String(), false))
		}
	}
	if r, ok := g.Rods[pos]; ok {
		g.Player.Rods[r] = rodProps{Charge: r.MaxCharge() - 1}
		g.DijkstraMapRebuild = true
		delete(g.Rods, pos)
		g.Printf("You take a %s.", r)
		g.StoryPrintf("You found and took a %s.", r)
	}
	if eq, ok := g.Equipables[pos]; ok {
		g.Printf("You are standing over %s.", Indefinite(eq.String(), false))
	} else if _, ok := g.Stairs[pos]; ok {
		g.Print("You are standing on a staircase.")
	} else if stn, ok := g.MagicalStones[pos]; ok {
		g.Printf("You are standing on %s.", Indefinite(stn.String(), false))
	} else if g.Doors[pos] {
		g.Print("You stand at the door.")
	}
}

func (g *game) MovePlayer(pos position, ev event) error {
	if !pos.valid() {
		return errors.New("You cannot move there.")
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell && !g.Player.HasStatus(StatusDig) {
		return errors.New("You cannot move into a wall.")
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
	if g.Player.Weapon == DefenderFlail && !mons.Exists() {
		mons = g.AttractMonster(pos)
	}
	if !mons.Exists() {
		if g.Player.HasStatus(StatusLignification) {
			return errors.New("You cannot move while lignified")
		}
		if c.T == WallCell {
			g.Dungeon.SetCell(pos, FreeCell)
			g.MakeNoise(WallNoise, pos)
			g.Print(g.CrackSound())
			g.Fog(pos, 1, ev)
			g.Stats.Digs++
		}
		if g.Player.Aptitudes[AptFast] {
			// only fast for movement
			delay -= 2
		}
		switch g.Player.Armour {
		case TurtlePlates:
			delay += 3
		case SpeedRobe:
			delay -= 3
		case SmokingScales:
			_, ok := g.Clouds[g.Player.Pos]
			if !ok {
				g.Clouds[g.Player.Pos] = CloudFog
				g.PushEvent(&cloudEvent{ERank: ev.Rank() + 15 + RandInt(10), EAction: CloudEnd, Pos: g.Player.Pos})
			}
		}
		if g.Player.HasStatus(StatusSwift) {
			// only fast for movement
			delay -= 3
		}
		g.Stats.Moves++
		g.PlacePlayerAt(pos)
		if !g.Autoexploring {
			g.BoredomAction(ev, 1)
		}
		if g.Player.Statuses[StatusSlay] > 0 {
			g.Player.Statuses[StatusSlay] /= 2
		}
	} else {
		g.FunAction()
		g.AttackMonster(mons, ev)
	}
	if g.Player.HasStatus(StatusBerserk) {
		delay -= 3
	}
	if g.Player.HasStatus(StatusSlow) {
		delay += 3 * g.Player.Statuses[StatusSlow]
	}
	if delay < 3 {
		delay = 3
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
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: CloudEnd, Pos: pos})
		}
	}
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + 20 + RandInt(10)
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
		g.PushEvent(&simpleEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: ConfusionEnd})
		g.Print("You feel confused.")
	}
}

func (g *game) PlacePlayerAt(pos position) {
	g.Player.Pos = pos
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
}

func (g *game) EnterLignification(ev event) {
	g.Player.Statuses[StatusLignification]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 150 + RandInt(100), EAction: LignificationEnd})
	g.Player.HP += 10
}
