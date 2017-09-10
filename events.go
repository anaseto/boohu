package main

import "container/heap"

type event interface {
	Rank() int
	Action(*game)
	Renew(*game, int)
}

type eventQueue []event

func (evq eventQueue) Len() int {
	return len(evq)
}

func (evq eventQueue) Less(i, j int) bool {
	return evq[i].Rank() < evq[j].Rank()
}

func (evq eventQueue) Swap(i, j int) {
	evq[i], evq[j] = evq[j], evq[i]
}

func (evq *eventQueue) Push(x interface{}) {
	no := x.(event)
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
	// not used
	ResistanceEnd
)

type simpleEvent struct {
	ERank   int
	EAction simpleAction
}

func (sev *simpleEvent) Rank() int {
	return sev.ERank
}

func (sev *simpleEvent) Renew(g *game, delay int) {
	sev.ERank += delay
	heap.Push(g.Events, sev)
}

func (sev *simpleEvent) Action(g *game) {
	switch sev.EAction {
	case PlayerTurn:
		g.AutoNext = g.AutoPlayer(sev)
		if g.AutoNext {
			return
		}
		g.Quit = g.ui.HandlePlayerTurn(g, sev)
		if g.Quit {
			return
		}
	case HealPlayer:
		g.HealPlayer(sev)
	case MPRegen:
		g.MPRegen(sev)
	case Teleportation:
		g.Teleportation(sev)
	case BerserkEnd:
		g.Player.Statuses[StatusBerserk]--
		g.Player.Statuses[StatusSlow]++
		g.Player.Statuses[StatusExhausted]++
		g.Print("You are no longer berserk.")
		heap.Push(g.Events, &simpleEvent{ERank: sev.Rank() + 100, EAction: SlowEnd})
		heap.Push(g.Events, &simpleEvent{ERank: sev.Rank() + 300, EAction: ExhaustionEnd})
	case SlowEnd:
		g.Print("You feel no longer slow.")
		g.Player.Statuses[StatusSlow]--
	case ExhaustionEnd:
		g.Print("You feel no longer exhausted.")
		g.Player.Statuses[StatusExhausted]--
	case HasteEnd:
		g.Print("You feel no longer speedy.")
		g.Player.Statuses[StatusHaste]--
	case EvasionEnd:
		g.Print("You feel no longer agile.")
		g.Player.Statuses[StatusEvasion]--
	case ResistanceEnd:
		g.Print("You feel no longer resistant to the elements.")
		g.Player.Statuses[StatusResistance]--
	case LignificationEnd:
		g.Print("Your feel no longer attached to the ground.")
		g.Player.Statuses[StatusLignification]--
	case ConfusionEnd:
		g.Print("Your feel no longer confused.")
		g.Player.Statuses[StatusConfusion]--
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
			mons.Statuses[MonsConfused]--
		}
	case MonsExhaustionEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsExhausted]--
		}
	}
}

func (mev *monsterEvent) Renew(g *game, delay int) {
	mev.ERank += delay
	heap.Push(g.Events, mev)
}

type cloudAction int

const (
	CloudEnd cloudAction = iota
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
	}
}

func (cev *cloudEvent) Renew(g *game, delay int) {
	cev.ERank += delay
	heap.Push(g.Events, cev)
}
