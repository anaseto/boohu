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
	SlowEnd
	ExhaustionEnd
	HasteEnd
	EvasionEnd
	LignificationEnd
	ConfusionEnd
	NauseaEnd
	DigEnd
	ShadowsEnd
	LevitationEnd
	ShaedraAnimation
	ArtifactAnimation
)

func (g *game) PushEvent(ev event) {
	iev := iEvent{Event: ev, Index: g.EventIndex}
	g.EventIndex++
	heap.Push(g.Events, iev)
}

func (g *game) PushAgainEvent(ev event) {
	iev := iEvent{Event: ev, Index: 0}
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
	if delay == 0 {
		g.PushAgainEvent(sev)
	} else {
		g.PushEvent(sev)
	}
}

const DurationStatusStep = 10

var endmsgs = [...]string{
	SlowEnd:          "You no longer feel slow.",
	ExhaustionEnd:    "You no longer feel exhausted.",
	HasteEnd:         "You no longer feel speedy.",
	LignificationEnd: "You no longer feel attached to the ground.",
	ConfusionEnd:     "You no longer feel confused.",
	NauseaEnd:        "You no longer feel sick.",
	DigEnd:           "You no longer feel like an earth dragon.",
	LevitationEnd:    "You no longer levitate.",
}

var endstatuses = [...]status{
	SlowEnd:          StatusSlow,
	ExhaustionEnd:    StatusExhausted,
	HasteEnd:         StatusSwift,
	LignificationEnd: StatusLignification,
	ConfusionEnd:     StatusConfusion,
	NauseaEnd:        StatusNausea,
	DigEnd:           StatusDig,
	LevitationEnd:    StatusLevitation,
}

var statusEndActions = [...]simpleAction{
	StatusSlow:          SlowEnd,
	StatusExhausted:     ExhaustionEnd,
	StatusSwift:         HasteEnd,
	StatusLignification: LignificationEnd,
	StatusConfusion:     ConfusionEnd,
	StatusNausea:        NauseaEnd,
	StatusDig:           DigEnd,
	StatusLevitation:    LevitationEnd,
}

func (sev *simpleEvent) Action(g *game) {
	switch sev.EAction {
	case PlayerTurn:
		g.ComputeNoise()
		g.ComputeLOS() // TODO: optimize? most of the time almost redundant (unless on a tree)
		g.ComputeMonsterLOS()
		g.LogNextTick = g.LogIndex
		g.AutoNext = g.AutoPlayer(sev)
		if g.AutoNext {
			g.TurnStats()
			return
		}
		g.Quit = g.ui.HandlePlayerTurn(sev)
		if g.Quit {
			return
		}
		g.TurnStats()
	case ShaedraAnimation:
		g.ComputeLOS()
		g.ui.FreeingShaedraAnimation()
	case ArtifactAnimation:
		g.ComputeLOS()
		g.ui.TakingArtifactAnimation()
	case SlowEnd, ExhaustionEnd, HasteEnd, LignificationEnd, ConfusionEnd, NauseaEnd, DigEnd, LevitationEnd:
		g.Player.Statuses[endstatuses[sev.EAction]] -= DurationStatusStep
		if g.Player.Statuses[endstatuses[sev.EAction]] <= 0 {
			g.Player.Statuses[endstatuses[sev.EAction]] = 0
			g.PrintStyled(endmsgs[sev.EAction], logStatusEnd)
			g.ui.StatusEndAnimation()
			switch sev.EAction {
			case LevitationEnd:
				if g.Dungeon.Cell(g.Player.Pos).T == ChasmCell {
					g.FallAbyss()
				}
			}
		} else {
			g.PushEvent(&simpleEvent{ERank: sev.Rank() + DurationStatusStep, EAction: sev.EAction})
		}
	}
}

type monsterAction int

const (
	MonsterTurn monsterAction = iota
	MonsConfusionEnd
	MonsExhaustionEnd
	MonsSlowEnd
	MonsLignificationEnd
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
	case MonsConfusionEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsConfused] = 0
			if g.Player.Sees(mons.Pos) {
				g.Printf("The %s is no longer confused.", mons.Kind)
			}
			mons.Path = mons.APath(g, mons.Pos, mons.Target)
		}
	case MonsLignificationEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsLignified] = 0
			if g.Player.Sees(mons.Pos) {
				g.Printf("%s is no longer lignified.", mons.Kind.Definite(true))
			}
			mons.Path = mons.APath(g, mons.Pos, mons.Target)
		}
	case MonsSlowEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsSlow]--
			if g.Player.Sees(mons.Pos) {
				g.Printf("%s is no longer slowed.", mons.Kind.Definite(true))
			}
		}
	case MonsExhaustionEnd:
		mons := g.Monsters[mev.NMons]
		if mons.Exists() {
			mons.Statuses[MonsExhausted]--
			//if mons.State != Resting && g.Player.LOS[mons.Pos] &&
			//(mons.Kind.Ranged() || mons.Kind.Smiting()) && mons.Pos.Distance(g.Player.Pos) > 1 {
			//g.Printf("%s is ready to fire again.", mons.Kind.Definite(true))
			//}
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
	ObstructionProgression
	FireProgression
	NightProgression
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
		t := g.MagicalBarriers[cev.Pos]
		if !g.Player.Sees(cev.Pos) && g.Dungeon.Cell(cev.Pos).T == BarrierCell {
			// XXX does not handle all cases
			g.UpdateKnowledge(cev.Pos, BarrierCell)
		} else {
			delete(g.MagicalBarriers, cev.Pos)
			delete(g.TerrainKnowledge, cev.Pos)
		}
		// TODO: rework temporal walls so that they preserve doors and foliage
		if g.Dungeon.Cell(cev.Pos).T != BarrierCell {
			break
		}
		g.Dungeon.SetCell(cev.Pos, t)
	case ObstructionProgression:
		pos := g.FreePassableCell()
		g.MagicalBarrierAt(pos, cev)
		if g.Player.Sees(pos) {
			g.Printf("You see an oric barrier appear out of thin air.")
			g.StopAuto()
		}
		g.PushEvent(&cloudEvent{ERank: cev.Rank() + DurationObstructionProgression + RandInt(DurationObstructionProgression/4),
			EAction: ObstructionProgression})
	case FireProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		//g.BurnCreature(cev.Pos, cev)
		for _, pos := range g.Dungeon.FreeNeighbors(cev.Pos) {
			if RandInt(5) > 0 {
				continue
			}
			g.Burn(pos, cev)
		}
		delete(g.Clouds, cev.Pos)
		g.NightFog(cev.Pos, 1, &simpleEvent{ERank: cev.Rank()})
		g.ComputeLOS()
	case NightProgression:
		if _, ok := g.Clouds[cev.Pos]; !ok {
			break
		}
		g.MakeCreatureSleep(cev.Pos, cev)
		if RandInt(20) == 0 {
			delete(g.Clouds, cev.Pos)
			g.ComputeLOS()
			break
		}
		cev.Renew(g, 10)
	}
}

func (g *game) NightFog(at position, radius int, ev event) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	nm.iter(at, func(n *node) {
		pos := n.Pos
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudNight
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: NightProgression, Pos: pos})
			g.MakeCreatureSleep(pos, ev)
		}
	})
	g.ComputeLOS()
}

func (g *game) MakeCreatureSleep(pos position, ev event) {
	if pos == g.Player.Pos {
		if g.PutStatus(StatusSlow, DurationSleepSlow) {
			g.Print("The clouds of night make you sleepy.")
		}
		return
	}
	mons := g.MonsterAt(pos)
	if !mons.Exists() || (RandInt(2) == 0 && mons.Status(MonsExhausted)) {
		// do not always make already exhausted monsters sleep (they were probably awaken)
		return
	}
	if mons.State != Resting && g.Player.Sees(mons.Pos) {
		g.Printf("%s falls asleep.", mons.Kind.Definite(true))
	}
	mons.State = Resting
	mons.Dir = NoDir
	mons.ExhaustTime(g, 40+RandInt(10))
}

//func (g *game) BurnCreature(pos position, ev event) {
//mons := g.MonsterAt(pos)
//if mons.Exists() {
//mons.HP -= 1
//if mons.HP <= 0 {
//if g.Player.Sees(mons.Pos) {
//g.PrintfStyled("%s is killed by the fire.", logPlayerHit, mons.Kind.Definite(true))
//}
//g.HandleKill(mons, ev)
//} else {
//mons.MakeAwareIfHurt(g)
//}
//}
//if pos == g.Player.Pos {
//damage := 1
//g.Player.HP -= damage
//g.PrintfStyled("The fire burns you (%d dmg).", logMonsterHit, damage)
//if g.Player.HP+damage < 10 {
//g.Stats.TimesLucky++
//}
//g.StopAuto()
//}
//}

func (g *game) Burn(pos position, ev event) {
	if _, ok := g.Clouds[pos]; ok {
		return
	}
	c := g.Dungeon.Cell(pos)
	if !c.Flammable() {
		return
	}
	g.Stats.Burns++
	switch c.T {
	case DoorCell:
		g.Print("The door vanishes in magical flames.")
	case TableCell:
		g.Print("The table vanishes in magical flames.")
	case BarrelCell:
		g.Print("The barrel vanishes in magical flames.")
		delete(g.Objects.Barrels, pos)
	case TreeCell:
		g.Print("The tree vanishes in magical flames.")
	}
	g.Dungeon.SetCell(pos, GroundCell)
	g.Clouds[pos] = CloudFire
	if !g.Player.Sees(pos) {
		// TODO: knowledge
	} else {
		g.ComputeLOS()
	}
	g.PushEvent(&cloudEvent{ERank: ev.Rank() + DurationCloudProgression, EAction: FireProgression, Pos: pos})
	//g.BurnCreature(pos, ev)
}

func (cev *cloudEvent) Renew(g *game, delay int) {
	cev.ERank += delay
	g.PushEvent(cev)
}

const (
	DurationSwiftness              = 50
	DurationLevitation             = 180
	DurationShortSwiftness         = 30
	DurationDigging                = 80
	DurationSlow                   = 120
	DurationSleepSlow              = 40
	DurationCloudProgression       = 10
	DurationFog                    = 150
	DurationExhaustion             = 70
	DurationConfusion              = 110
	DurationConfusionPlayer        = 50
	DurationLignification          = 110
	DurationLignificationPlayer    = 30
	DurationMagicalBarrier         = 150
	DurationObstructionProgression = 150
	DurationSmokingCloakFog        = 20
	DurationMonsterExhaustion      = 100
)
