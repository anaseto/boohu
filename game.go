package main

import "container/heap"

// idea: punish (mark) the player if it remains too many turns not in combat
// and not exploring new cells and in good health
// monster induced berserker

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	Bands               []monsterBand
	Events              *eventQueue
	Highlight           map[position]bool // highlighted positions (e.g. targeted ray)
	Collectables        map[position]*collectable
	CollectableScore    int
	Equipables          map[position]equipable
	Rods                map[position]rod
	Stairs              map[position]bool
	Clouds              map[position]cloud
	Fungus              map[position]vegetation
	Doors               map[position]bool
	TemporalWalls       map[position]bool
	GeneratedBands      map[monsterBand]int
	GeneratedEquipables map[equipable]bool
	GeneratedRods       map[rod]bool
	FoundEquipables     map[equipable]bool
	Gold                map[position]int
	UnknownDig          map[position]bool
	Resting             bool
	Autoexploring       bool
	AutoexploreMap      nodeMap
	AutoTarget          *position
	AutoDir             *direction
	AutoHalt            bool
	AutoNext            bool
	ExclusionsMap       map[position]bool
	Quit                bool
	ui                  Renderer
	Depth               int
	Wizard              bool
	Log                 []string
	Story               []string
	Turn                int
	Killed              int
	KilledMons          map[monsterKind]int
	Scumming            int
	Noise               map[position]bool
}

type Renderer interface {
	ExploreStep(*game) bool
	HandlePlayerTurn(*game, event) bool
	Death(*game)
	ChooseTarget(*game, Targetter) bool
	CriticalHPWarning(*game)
}

func (g *game) FreeCell() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos == pos {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForImportantStair() position {
	for {
		pos := g.FreeCellForStatic()
		if pos.Distance(g.Player.Pos) > 12 {
			return pos
		}
	}
}

func (g *game) FreeCellForStatic() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForStatic")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos == pos {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			if g.Doors[pos] {
				continue
			}
			if g.Gold[pos] > 0 {
				continue
			}
			if g.Collectables[pos] != nil {
				continue
			}
			if g.Stairs[pos] {
				continue
			}
			if _, ok := g.Rods[pos]; ok {
				continue
			}
			if _, ok := g.Equipables[pos]; ok {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForMonster() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForMonster")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForBandMonster(pos position) position {
	count := 0
	for {
		count++
		if count > 1000 {
			return g.FreeCellForMonster()
		}
		neighbors := g.Dungeon.FreeNeighbors(pos)
		r := RandInt(len(neighbors))
		pos = neighbors[r]
		if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
			continue
		}
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
	}
}

func (g *game) FreeForStairs() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeForStairs")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			_, ok := g.Collectables[pos]
			if ok {
				continue
			}
			return pos
		}
	}
}

func (g *game) MaxDepth() int {
	return 12
}

const (
	DungeonHeigth = 21
	DungeonWidth  = 79
)

func (g *game) GenDungeon() {
	g.Fungus = make(map[position]vegetation)
	switch RandInt(6) {
	//switch 1 {
	case 0:
		g.GenCaveMap(DungeonHeigth, DungeonWidth)
		g.Fungus = g.Foliage(DungeonHeigth, DungeonWidth)
	case 1:
		g.GenRoomMap(DungeonHeigth, DungeonWidth)
	case 2:
		g.GenCellularAutomataCaveMap(DungeonHeigth, DungeonWidth)
		g.Fungus = g.Foliage(DungeonHeigth, DungeonWidth)
	case 3:
		g.GenCaveMapTree(DungeonHeigth, DungeonWidth)
	default:
		g.GenRuinsMap(DungeonHeigth, DungeonWidth)
	}
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        40,
		MP:        10,
		Gold:      0,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Consumables = map[consumable]int{
		HealWoundsPotion: 1,
		Javelin:          3,
	}
	switch RandInt(6) {
	case 0, 1:
		g.Player.Consumables[TeleportationPotion] = 1
	case 2, 3:
		g.Player.Consumables[BerserkPotion] = 1
	case 4:
		g.Player.Consumables[EvasionPotion] = 1
	case 5:
		g.Player.Consumables[LignificationPotion] = 1
	}
	g.Player.Rods = map[rod]*rodProps{}
	g.Player.Statuses = map[status]int{}

	// Testing
	// g.Player.Aptitudes[AptSmoke] = true
}

func (g *game) InitLevel() {
	// Dungeon terrain
	g.GenDungeon()

	// Starting data
	if g.Depth == 0 {
		g.InitPlayer()
		g.GeneratedRods = map[rod]bool{}
		g.GeneratedEquipables = map[equipable]bool{}
		g.FoundEquipables = map[equipable]bool{Robe: true, Dagger: true}
		g.GeneratedBands = map[monsterBand]int{}
		g.KilledMons = map[monsterKind]int{}
	}

	g.Player.Pos = g.FreeCell()

	g.UnknownDig = map[position]bool{}
	g.ExclusionsMap = map[position]bool{}
	g.TemporalWalls = map[position]bool{}

	// Monsters
	g.GenMonsters()

	// Collectables
	g.Collectables = make(map[position]*collectable)
	g.GenCollectables()

	// Equipment
	g.Equipables = make(map[position]equipable)
	for eq, data := range EquipablesRepartitionData {
		g.GenEquip(eq, data)
	}

	// Rods
	g.Rods = map[position]rod{}
	r := 7*(g.GeneratedRodsCount()+1) - 2*(g.Depth+1)
	if r < -3 {
		r = 0
	} else if r < 2 {
		r = 1
	}
	if RandInt(r) == 0 && g.GeneratedRodsCount() < 3 {
		g.GenerateRod()
	}

	// Aptitudes/Mutations
	r = 5*g.Player.AptitudeCount() - g.Depth + 2
	if r < 2 {
		r = 1
	}
	if RandInt(r) == 0 && g.Depth > 0 && g.Player.AptitudeCount() < 3 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Stairs
	g.Stairs = make(map[position]bool)
	nstairs := 1 + RandInt(3)
	if g.Depth == g.MaxDepth() {
		nstairs = 1
	} else if g.Depth == g.MaxDepth()-1 && nstairs > 2 {
		nstairs = 1 + RandInt(2)
	}
	for i := 0; i < nstairs; i++ {
		var pos position
		if g.Depth > 9 {
			pos = g.FreeCellForImportantStair()
		} else {
			pos = g.FreeCellForStatic()
		}
		g.Stairs[pos] = true
	}

	// Gold
	g.Gold = make(map[position]int)
	for i := 0; i < 5; i++ {
		pos := g.FreeCellForStatic()
		g.Gold[pos] = 1 + RandInt(g.Depth+g.Depth*g.Depth/10)
	}

	// initialize LOS
	if g.Depth == 0 {
		g.Print("You're in Hareka's Underground. Good luck! Press ? for help.")
	}
	if g.Depth == g.MaxDepth() {
		g.Print("You feel magic in the air. The way out is close.")
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// recharge rods
	for r, props := range g.Player.Rods {
		if props.Charge < r.MaxCharge() {
			props.Charge += RandInt(1 + r.Rate())
		}
		if props.Charge > r.MaxCharge() {
			props.Charge = r.MaxCharge()
		}
	}

	// clouds
	g.Clouds = map[position]cloud{}

	// Events
	if g.Depth == 0 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		heap.Push(g.Events, &simpleEvent{ERank: 0, EAction: PlayerTurn})
		heap.Push(g.Events, &simpleEvent{ERank: 50, EAction: HealPlayer})
		heap.Push(g.Events, &simpleEvent{ERank: 100, EAction: MPRegen})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		heap.Push(g.Events, &monsterEvent{ERank: g.Turn + 1, EAction: MonsterTurn, NMons: i})
		heap.Push(g.Events, &monsterEvent{ERank: g.Turn + 50, EAction: HealMonster, NMons: i})
	}
}

func (g *game) CleanEvents() {
	evq := &eventQueue{}
	for g.Events.Len() > 0 {
		ev := heap.Pop(g.Events).(event)
		switch ev.(type) {
		case *monsterEvent:
		case *cloudEvent:
		default:
			heap.Push(evq, ev)
		}
	}
	g.Events = evq
}

func (g *game) GenCollectables() {
	rounds := 10
	for i := 0; i < rounds; i++ {
		for c, data := range ConsumablesCollectData {
			var r int
			if g.CollectableScore >= 5*(g.Depth+1)/3 {
				r = RandInt(data.rarity * rounds * 4)
			} else if g.CollectableScore < 4*(g.Depth+1)/3 {
				r = RandInt(data.rarity * rounds / 4)
			} else {
				r = RandInt(data.rarity * rounds)
			}

			if r == 0 {
				g.CollectableScore++
				pos := g.FreeCellForStatic()
				g.Collectables[pos] = &collectable{Consumable: c, Quantity: data.quantity}
			}
		}
	}
}

func (g *game) SeenGoodWeapon() bool {
	return g.GeneratedEquipables[Sword] || g.GeneratedEquipables[DoubleSword] || g.GeneratedEquipables[Spear] || g.GeneratedEquipables[Halberd] ||
		g.GeneratedEquipables[Axe] || g.GeneratedEquipables[BattleAxe]
}

func (g *game) GenEquip(eq equipable, data equipableData) {
	depthAdjust := data.minDepth - g.Depth
	var r int
	if depthAdjust >= 0 {
		r = RandInt(data.rarity * (depthAdjust + 1) * (depthAdjust + 1))
	} else {
		switch eq.(type) {
		case shield:
			if !g.GeneratedEquipables[eq] {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				r = RandInt(data.rarity * 2)
			}
		case armour:
			if !g.GeneratedEquipables[eq] && eq != Robe {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				r = RandInt(5 * data.rarity / 4)
			}
		case weapon:
			if !g.SeenGoodWeapon() && eq != Dagger {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				if g.Player.Weapon != Dagger {
					r = RandInt(data.rarity * 4)
				} else {
					r = RandInt(5 * data.rarity / 4)
				}
			}
		default:
			// not reached
			r = RandInt(data.rarity)
		}
	}
	if r == 0 {
		pos := g.FreeCellForStatic()
		g.Equipables[pos] = eq
		g.GeneratedEquipables[eq] = true
	}

}

func (g *game) Descend(ev event) bool {
	if g.Depth >= g.MaxDepth() {
		g.Depth++
		// win
		g.RemoveSaveFile()
		return true
	}
	g.Print("You descend deeper in the dungeon.")
	g.Depth++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank(), EAction: PlayerTurn})
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) WizardMode() {
	g.Wizard = true
	g.Player.Consumables[DescentPotion] = 12
	g.Print("You are now in wizard mode and cannot obtain winner status.")
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		if g.MonsterInLOS() == nil &&
			(g.Player.HP < g.Player.HPMax() || g.Player.MP < g.Player.MPMax() || g.Player.HasStatus(StatusExhausted) ||
				g.Player.HasStatus(StatusConfusion) || g.Player.HasStatus(StatusLignification)) {
			g.WaitTurn(ev)
			return true
		}
		g.Resting = false
	} else if g.Autoexploring {
		if g.ui.ExploreStep(g) {
			g.AutoHalt = true
		}
		mons := g.MonsterInLOS()
		switch {
		case mons.Exists():
			g.Print("You stop exploring.")
		case g.AutoHalt:
			// stop exploring for other reasons
			g.Print("You stop exploring.")
		default:
			var n *node
			var b bool
			count := 0
			for {
				count++
				if count > 100 {
					// should not happen
					g.Print("Hm… something went wrong with auto-explore. You stop.")
					n = nil
					break
				}
				n, b = g.NextAuto()
				if !b {
					if n == nil {
						g.Print("You could not reach safely some places.")
					}
					break
				}
				sources := g.AutoexploreSources()
				if len(sources) == 0 {
					g.Print("You finished exploring.")
					n = nil
					break
				}
				g.BuildAutoexploreMap(sources)
			}
			if n != nil {
				err := g.MovePlayer(n.Pos, ev)
				if err != nil {
					g.Print(err.Error())
					break
				}
				return true
			}
		}
		g.Autoexploring = false
	} else if g.AutoTarget != nil {
		if !g.ui.ExploreStep(g) && g.MoveToTarget(ev) {
			return true
		}
	} else if g.AutoDir != nil {
		if !g.ui.ExploreStep(g) && g.AutoToDir(ev) {
			return true
		}

	}
	return false
}

func (g *game) EventLoop() {
loop:
	for {
		if g.Player.HP <= 0 {
			if g.Wizard {
				g.Player.HP = g.Player.HPMax()
			} else {
				g.ui.Death(g)
				g.RemoveSaveFile()
				break loop
			}
		}
		if g.Events.Len() == 0 {
			break loop
		}
		ev := heap.Pop(g.Events).(event)
		g.Turn = ev.Rank()
		ev.Action(g)
		if g.AutoNext {
			continue loop
		}
		if g.Quit {
			break loop
		}
	}
}
