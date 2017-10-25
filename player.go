package main

import (
	"errors"
	"fmt"
)

type player struct {
	LOS         map[position]bool
	Rays        rayMap
	Pos         position
	HP          int
	MP          int
	Consumables map[consumable]int
	Gold        int
	Target      position
	Statuses    map[status]int
	Armour      armour
	Weapon      weapon
	Shield      shield
	Aptitudes   map[aptitude]bool
	Rods        map[rod]*rodProps
}

func (p *player) HPMax() int {
	hpmax := 40
	if p.Aptitudes[AptHealthy] {
		hpmax += 10
	}
	return hpmax
}

func (p *player) MPMax() int {
	mpmax := 10
	if p.Aptitudes[AptMagic] {
		mpmax += 5
	}
	return mpmax
}

func (p *player) Accuracy() int {
	acc := 15
	if p.Aptitudes[AptAccurate] {
		acc += 3
	}
	return acc
}

func (p *player) Armor() int {
	ar := 0
	switch p.Armour {
	case LeatherArmour:
		ar += 3
	case ChainMail:
		ar += 4
	case PlateArmour:
		ar += 6
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
		err := g.MovePlayer(g.Player.Pos.To(*g.AutoDir), ev)
		if err != nil {
			g.Print(err.Error())
			g.AutoDir = nil
			return false
		}
		return true
	}
	g.AutoDir = nil
	return false
}

func (g *game) GoToDir(dir direction, ev event) error {
	if g.MonsterInLOS() != nil {
		g.AutoDir = nil
		return errors.New("You cannot travel while there are monsters in view.")
	}
	err := g.MovePlayer(g.Player.Pos.To(dir), ev)
	if err != nil {
		return err
	}
	g.AutoDir = &dir
	return nil
}

func (g *game) MoveToTarget(ev event) bool {
	if g.MonsterInLOS() == nil {
		path := g.PlayerPath(g.Player.Pos, *g.AutoTarget)
		if len(path) > 1 {
			err := g.MovePlayer(path[len(path)-2], ev)
			if err != nil {
				g.Print(err.Error())
				g.AutoTarget = nil
				return false
			}
			return true
		}
	}
	g.AutoTarget = nil
	return false
}

func (g *game) WaitTurn(ev event) {
	// XXX Really wait for 10 ?
	g.ScummingAction(ev)
	ev.Renew(g, 10)
}

func (g *game) ExistsMonster() bool {
	for _, mons := range g.Monsters {
		if mons.Exists() {
			return true
		}
	}
	return false
}

func (g *game) ScummingAction(ev event) {
	if g.Player.HP == g.Player.HPMax() && g.Player.MP == g.Player.MPMax() {
		g.Scumming++
	}
	if g.Scumming == 100 {
		if g.ExistsMonster() {
			g.Print("You feel a little bored.")
		}
	}
	if g.Scumming > 120 {
		if !g.ExistsMonster() {
			g.Scumming = 0
			return
		}
		g.Player.HP = g.Player.HP / 2
		if RandInt(2) == 0 {
			g.MakeNoise(100, g.Player.Pos)
			neighbors := g.Dungeon.Neighbors(g.Player.Pos)
			for _, pos := range neighbors {
				if RandInt(3) != 0 {
					g.Dungeon.SetCell(pos, FreeCell)
				}
			}
			g.Print("You hear a terrible explosion coming from the ground. You are lignified.")
			g.Player.Statuses[StatusLignification]++
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 240 + RandInt(10), EAction: LignificationEnd})
		} else {
			delay := 20 + RandInt(5)
			g.Player.Statuses[StatusTele]++
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + delay, EAction: Teleportation})
			g.Print("Something hurt you! You feel unstable.")
		}
		g.Scumming = 0
	}
}

func (g *game) FairAction() {
	g.Scumming -= 10
	if g.Scumming < 0 {
		g.Scumming = 0
	}
}

func (g *game) Rest(ev event) error {
	if g.MonsterInLOS() != nil {
		return fmt.Errorf("You cannot sleep while monsters are in view.")
	}
	if g.Player.HP == g.Player.HPMax() && g.Player.MP == g.Player.MPMax() && !g.Player.HasStatus(StatusExhausted) &&
		!g.Player.HasStatus(StatusConfusion) && !g.Player.HasStatus(StatusLignification) {
		return errors.New("You do not need to rest.")
	}
	g.WaitTurn(ev)
	g.Resting = true
	return nil
}

func (g *game) Equip(ev event) error {
	if eq, ok := g.Equipables[g.Player.Pos]; ok {
		eq.Equip(g)
		ev.Renew(g, 10)
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
	g.Player.Statuses[StatusTele]--
	if g.Dungeon.Valid(pos) {
		// should always happen
		g.Player.Pos = pos
		g.Print("You feel yourself teleported away.")
		g.CollectGround()
		g.ComputeLOS()
		g.MakeMonstersAware()
	} else {
		g.Print("Something went wrong with the teleportation.")
	}
}

func (g *game) CollectGround() {
	pos := g.Player.Pos
	if g.Gold[pos] > 0 {
		g.Player.Gold += g.Gold[pos]
		g.Printf("You pick up %d gold.", g.Gold[pos])
		delete(g.Gold, pos)
	}
	if c, ok := g.Collectables[pos]; ok && c != nil {
		g.Player.Consumables[c.Consumable] += c.Quantity
		delete(g.Collectables, pos)
		if c.Quantity > 1 {
			g.Printf("You take %d %s.", c.Quantity, c.Consumable.Plural())
		} else {
			g.Printf("You take %s.", Indefinite(c.Consumable.String(), false))
		}
	}
	if r, ok := g.Rods[pos]; ok {
		g.Player.Rods[r] = &rodProps{Charge: r.MaxCharge() - 1}
		delete(g.Rods, pos)
		g.Printf("You take a %s.", r)
		g.StoryPrintf("You found and took a %s.", r)
	}
}

func (g *game) MovePlayer(pos position, ev event) error {
	if !g.Dungeon.Valid(pos) || g.Dungeon.Cell(pos).T == WallCell {
		return errors.New("You cannot move there.")
	}
	if g.Player.HasStatus(StatusConfusion) {
		switch pos.Dir(g.Player.Pos) {
		case E, N, W, S:
		default:
			return errors.New("You cannot use diagonal movements while confused.")
		}
	}
	delay := 10
	switch g.Dungeon.Cell(pos).T {
	case FreeCell:
		mons, _ := g.MonsterAt(pos)
		if !mons.Exists() {
			if g.Player.HasStatus(StatusLignification) {
				return errors.New("You cannot move while lignified")
			}
			g.Player.Pos = pos
			g.CollectGround()
			g.ComputeLOS()
			if g.Autoexploring {
				mons := g.MonsterInLOS()
				if mons.Exists() {
					g.Printf("You see %s (%v).", Indefinite(mons.Kind.String(), false), mons.State)
				}
				g.FairAction()
			} else {
				g.ScummingAction(ev)
			}
			g.MakeMonstersAware()
			if g.Player.Aptitudes[AptFast] {
				// only fast for movement
				delay -= 2
			}
			if g.Player.HasStatus(StatusSwift) {
				// only fast for movement
				delay -= 3
			}
		} else {
			g.FairAction()
			g.AttackMonster(mons, ev)
		}
	}
	if g.Player.HasStatus(StatusBerserk) {
		delay -= 3
	}
	if g.Player.HasStatus(StatusSlow) {
		delay += 3
	}
	ev.Renew(g, delay)
	return nil
}

func (g *game) HealPlayer(ev event) {
	if g.Player.HP < g.Player.HPMax() {
		g.Player.HP++
	}
	delay := 50
	if g.Player.Aptitudes[AptRegen] {
		delay = 25
	}
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
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 20 + RandInt(10), EAction: HasteEnd})
	g.ComputeLOS()
	g.Print("You feel an energy burst and smoking coming out from you.")
}
