package main

import "container/heap"

type event interface {
	Rank() int
	Action(*game)
	Renew(*game, int)
}

type iEvent struct {
	Event event
	Index int
}

type eventQueue []iEvent

func (evq eventQueue) Len() int {
	return len(evq)
}

func (evq eventQueue) Less(i, j int) bool {
	return evq[i].Event.Rank() < evq[j].Event.Rank() ||
		evq[i].Event.Rank() == evq[j].Event.Rank() && evq[i].Index < evq[j].Index
}

func (evq eventQueue) Swap(i, j int) {
	evq[i], evq[j] = evq[j], evq[i]
}

func (evq *eventQueue) Push(x interface{}) {
	no := x.(iEvent)
	*evq = append(*evq, no)
}

func (evq *eventQueue) Pop() interface{} {
	old := *evq
	n := len(old)
	no := old[n-1]
	*evq = old[0 : n-1]
	return no
}

type simpleAction int

const (
	PlayerTurn simpleAction = iota
	HealPlayer
	MPRegen
	Teleportation
	BerserkEnd
	SlowEnd
	ExhaustionEnd
	HasteEnd
	EvasionEnd
	LignificationEnd
	ConfusionEnd
	NauseaEnd
	DisabledShieldEnd
	CorrosionEnd
	DigEnd
)

func (g *game) PushEvent(ev event) {
	iev := iEvent{Event: ev, Index: g.EventIndex}
	g.EventIndex++
	heap.Push(g.Events, iev)
}

func (g *game) PopIEvent() iEvent {
	iev := heap.Pop(g.Events).(iEvent)
	return iev
}

type simpleEvent struct {
	ERank   int
	EAction simpleAction
}

func (sev *simpleEvent) Rank() int {
	return sev.ERank
}

func (sev *simpleEvent) Renew(g *game, delay int) {
	sev.ERank += delay
	g.PushEvent(sev)
}

func (sev *simpleEvent) Action(g *game) {
	switch sev.EAction {
	case PlayerTurn:
		g.ComputeNoise()
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer(sev)
		if g.AutoNext {
			g.Stats.Turns++
			return
		}
		g.Quit = g.ui.HandlePlayerTurn(g, sev)
		if g.Quit {
			return
		}
		g.TurnStats()
	case HealPlayer:
		g.HealPlayer(sev)
	case MPRegen:
		g.MPRegen(sev)
	case Teleportation:
		g.Teleportation(sev)
		g.Player.Statuses[StatusTele] = 0
	case BerserkEnd:
		g.Player.Statuses[StatusBerserk] = 0
		g.Player.Statuses[StatusSlow] = 1
		g.Player.Statuses[StatusExhausted] = 1
		g.Player.HP -= int(10 * g.Player.HP / Max(g.Player.HPMax(), g.Player.HP))
		g.PrintStyled("You are no longer berserk.", logStatusEnd)
		g.PushEvent(&simpleEvent{ERank: sev.Rank() + 90 + RandInt(40), EAction: SlowEnd})
		g.PushEvent(&simpleEvent{ERank: sev.Rank() + 270 + RandInt(60), EAction: ExhaustionEnd})
	case SlowEnd:
		g.PrintStyled("You feel no longer slow.", logStatusEnd)
		g.Player.Statuses[StatusSlow] = 0
	case ExhaustionEnd:
		g.PrintStyled("You feel no longer exhausted.", logStatusEnd)
		g.Player.Statuses[StatusExhausted] = 0
	case HasteEnd:
		g.Player.Statuses[StatusSwift]--
		if g.Player.Statuses[StatusSwift] == 0 {
			g.PrintStyled("You feel no longer speedy.", logStatusEnd)
		}
	case EvasionEnd:
		g.Player.Statuses[StatusAgile]--
		if g.Player.Statuses[StatusAgile] == 0 {
			g.PrintStyled("You feel no longer agile.", logStatusEnd)
		}
	case LignificationEnd:
		g.Player.Statuses[StatusLignification]--
		g.Player.HP -= int(10 * g.Player.HP / Max(g.Player.HPMax(), g.Player.HP))
		if g.Player.Statuses[StatusLignification] == 0 {
			g.PrintStyled("You feel no longer attached to the ground.", logStatusEnd)
		}
	case ConfusionEnd:
		g.PrintStyled("You feel no longer confused.", logStatusEnd)
		g.Player.Statuses[StatusConfusion] = 0
	case NauseaEnd:
		g.PrintStyled("You feel no longer sick.", logStatusEnd)
		g.Player.Statuses[StatusNausea] = 0
	case DisabledShieldEnd:
		g.PrintStyled("You manage to free your shield from the projectile.", logStatusEnd)
		g.Player.Statuses[StatusDisabledShield] = 0
	case CorrosionEnd:
		g.Player.Statuses[StatusCorrosion]--
		if g.Player.Statuses[StatusCorrosion] == 0 {
			g.PrintStyled("Your equipment is now free from acid.", logStatusEnd)
		}
	case DigEnd:
		g.Player.Statuses[StatusDig] = 0
		if g.Player.Statuses[StatusDig] == 0 {
			g.PrintStyled("You feel no longer like an earth dragon.", logStatusEnd)
		}
	}
}

type monsterAction int

const (
	MonsterTurn monsterAction = iota
	HealMonster
	MonsConfusionEnd
	MonsExhaustionEnd
)

type monsterEvent struct {
	ERank   int
	NMons   int
	EAction monsterAction
}

func (mev *monsterEvent) Rank() int {
	return mev.ERank
}

func (mev *monsterEvent) Action(g *game) {
	switch mev.EAction {
	case MonsterTurn:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.HandleTurn(g, mev)
		}
	case HealMonster:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Heal(g, mev)
		}
	case MonsConfusionEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsConfused] = 0
			g.Printf("The %s is no longer confused.", mons.Kind)
		}
	case MonsExhaustionEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsExhausted] = 0
		}
	}
}

func (mev *monsterEvent) Renew(g *game, delay int) {
	mev.ERank += delay
	g.PushEvent(mev)
}

type cloudAction int

const (
	CloudEnd cloudAction = iota
	ObstructionEnd
	FireProgression
)

type cloudEvent struct {
	ERank   int
	Pos     position
	EAction cloudAction
}

func (cev *cloudEvent) Rank() int {
	return cev.ERank
}

func (cev *cloudEvent) Action(g *game) {
	switch cev.EAction {
	case CloudEnd:
		delete(g.Clouds, cev.Pos)
		g.ComputeLOS()
	case ObstructionEnd:
		if !g.Player.LOS[cev.Pos] && g.Dungeon.Cell(cev.Pos).T == WallCell {
			g.UnknownDig[cev.Pos] = true
		} else {
			delete(g.TemporalWalls, cev.Pos)
		}
		if g.Dungeon.Cell(cev.Pos).T == FreeCell {
			break
		}
		g.Dungeon.SetCell(cev.Pos, FreeCell)
		g.MakeNoise(15, cev.Pos)
		g.Fog(cev.Pos, 1, &simpleEvent{ERank: cev.Rank()})
		g.ComputeLOS()
	case FireProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		g.BurnCreature(cev.Pos, cev)
		if RandInt(10) == 0 {
			delete(g.Clouds, cev.Pos)
			g.Fog(cev.Pos, 1, &simpleEvent{ERank: cev.Rank()})
			g.ComputeLOS()
			break
		}
		for _, pos := range g.Dungeon.FreeNeighbors(cev.Pos) {
			if RandInt(3) > 0 {
				continue
			}
			g.Burn(pos, cev)
		}
		cev.Renew(g, 10)
	}
}

func (g *game) BurnCreature(pos position, ev event) {
	mons := g.MonsterAt(pos)
	if mons.Exists() {
		mons.HP -= 1 + RandInt(10)
		if mons.HP <= 0 {
			g.PrintfStyled("%s is killed by the fire.", logPlayerHit, mons.Kind.Definite(true))
			g.HandleKill(mons, ev)
		} else {
			mons.MakeAwareIfHurt(g)
		}
	}
	if pos == g.Player.Pos {
		damage := 1 + RandInt(10)
		if damage > g.Player.HP {
			damage = 1 + RandInt(10)
		}
		g.Player.HP -= damage
		g.PrintfStyled("The fire burns you (%d damage).", logMonsterHit, damage)
		if g.Player.HP+damage < 10 {
			g.Stats.TimesLucky++
		}
	}
}

func (g *game) Burn(pos position, ev event) {
	if _, ok := g.Clouds[pos]; ok {
		return
	}
	_, okFungus := g.Fungus[pos]
	_, okDoor := g.Doors[pos]
	if !okFungus && !okDoor {
		return
	}
	g.Stats.Burns++
	delete(g.Fungus, pos)
	if _, ok := g.Doors[pos]; ok {
		delete(g.Doors, pos)
		g.Print("The door vanishes in flames.")
	}
	g.Clouds[pos] = CloudFire
	if !g.Player.LOS[pos] {
		g.UnknownBurn[pos] = true
	}
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + 10, EAction: FireProgression, Pos: pos})
	g.BurnCreature(pos, ev)
}

func (cev *cloudEvent) Renew(g *game, delay int) {
	cev.ERank += delay
	g.PushEvent(cev)
}
