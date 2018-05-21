package main

import "container/heap"

var Version string = "v0.8-dev"

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	MonstersPosCache    []int // monster (dungeon index + 1) / no monster (0)
	Bands               []monsterBand
	Events              *eventQueue
	Ev                  event
	EventIndex          int
	Depth               int
	Turn                int
	Highlight           map[position]bool // highlighted positions (e.g. targeted ray)
	Collectables        map[position]collectable
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
	Simellas            map[position]int
	UnknownDig          map[position]bool
	UnknownBurn         map[position]burn
	ExclusionsMap       map[position]bool
	Noise               map[position]bool
	Resting             bool
	Autoexploring       bool
	DijkstraMapRebuild  bool
	Targeting           position
	AutoTarget          position
	AutoDir             direction
	AutoHalt            bool
	AutoNext            bool
	Log                 []logEntry
	LogIndex            int
	LogNextTick         int
	InfoEntry           string
	Stats               stats
	Scumming            int
	Quit                bool
	Wizard              bool
	ui                  Renderer
}

type Renderer interface {
	ExploreStep(*game) bool
	HandlePlayerTurn(*game, event) bool
	Death(*game)
	ChooseTarget(*game, Targeter) error
	CriticalHPWarning(*game)
	ExplosionAnimation(*game, explosionStyle, position)
	LightningBoltAnimation(*game, []position)
	ThrowAnimation(*game, []position, bool)
	DrinkingPotionAnimation(*game)
	SwappingAnimation(*game, position, position)
	TeleportAnimation(*game, position, position, bool)
	MagicMappingAnimation(*game, []int)
	HitAnimation(*game, position, bool)
	WoundedAnimation(*game)
	WallExplosionAnimation(*game, position)
	ProjectileTrajectoryAnimation(*game, []position, uicolor)
	StatusEndAnimation(*game)
	DrawDungeonView(*game, uiMode)
}

func (g *game) FreeCell() position {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && g.Player.Pos == pos {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
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
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForStatic")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && g.Player.Pos == pos {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		if g.Doors[pos] {
			continue
		}
		if g.Simellas[pos] > 0 {
			continue
		}
		if _, ok := g.Collectables[pos]; ok {
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

func (g *game) FreeCellForMonster() position {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForMonster")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T != FreeCell {
			continue
		}
		if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
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
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
	}
}

func (g *game) FreeForStairs() position {
	d := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeForStairs")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T != FreeCell {
			continue
		}
		_, ok := g.Collectables[pos]
		if ok {
			continue
		}
		return pos
	}
}

const MaxDepth = 12

const (
	DungeonHeight = 21
	DungeonWidth  = 79
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	g.Fungus = make(map[position]vegetation)
	dg := GenRuinsMap
	switch RandInt(6) {
	//switch 0 {
	case 0:
		dg = GenCaveMap
	case 1:
		dg = GenRoomMap
	case 2:
		dg = GenCellularAutomataCaveMap
	case 3:
		dg = GenCaveMapTree
	}
	dg.Use(g)
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        40,
		MP:        10,
		Simellas:  0,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Consumables = map[consumable]int{
		HealWoundsPotion: 1,
	}
	switch RandInt(3) {
	case 1:
		g.Player.Consumables[ExplosiveMagara] = 1
	default:
		g.Player.Consumables[ConfusingDart] = 3
	}
	switch RandInt(10) {
	case 0, 1:
		g.Player.Consumables[TeleportationPotion] = 1
	case 2, 3:
		g.Player.Consumables[BerserkPotion] = 1
	case 4:
		g.Player.Consumables[SwiftnessPotion] = 1
	case 5:
		g.Player.Consumables[LignificationPotion] = 1
	case 6:
		g.Player.Consumables[WallPotion] = 1
	case 7:
		g.Player.Consumables[CBlinkPotion] = 1
	case 8:
		g.Player.Consumables[DigPotion] = 1
	case 9:
		g.Player.Consumables[SwapPotion] = 1
	}
	g.Player.Rods = map[rod]rodProps{}
	g.Player.Statuses = map[status]int{}

	// Testing
	// g.Player.Aptitudes[AptSmoke] = true
	//g.Player.Rods[RodSwapping] = rodProps{Charge: 3}
	g.Player.Rods[RodFireball] = rodProps{Charge: 3}
	//g.Player.Rods[RodFog] = rodProps{Charge: 3}
	//g.Player.Consumables[MagicMappingPotion] = 1
	g.Player.Consumables[ExplosiveMagara] = 5
	//g.Player.Consumables[DigPotion] = 5
	//g.Player.Consumables[SwapPotion] = 5
	//g.Player.Weapon = ElecWhip
}

func (g *game) InitLevel() {
	// Dungeon terrain
	g.GenDungeon()

	// Starting data
	if g.Depth == 0 {
		g.InitPlayer()
		g.AutoTarget = InvalidPos
		g.Targeting = InvalidPos
		g.GeneratedRods = map[rod]bool{}
		g.GeneratedEquipables = map[equipable]bool{}
		g.FoundEquipables = map[equipable]bool{Robe: true, Dagger: true}
		g.GeneratedBands = map[monsterBand]int{}
		g.Stats.KilledMons = map[monsterKind]int{}
	}

	g.MonstersPosCache = make([]int, DungeonNCells)
	g.Player.Pos = g.FreeCell()

	g.UnknownDig = map[position]bool{}
	g.UnknownBurn = map[position]burn{}
	g.ExclusionsMap = map[position]bool{}
	g.TemporalWalls = map[position]bool{}

	// Monsters
	g.GenMonsters()

	// Collectables
	g.Collectables = make(map[position]collectable)
	g.GenCollectables()

	// Equipment
	g.Equipables = make(map[position]equipable)
	for eq, data := range EquipablesRepartitionData {
		if _, ok := eq.(weapon); ok {
			continue
		}
		g.GenEquip(eq, data)
	}
	g.GenWeapon()

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
	if g.Depth == MaxDepth {
		nstairs = 1
	} else if g.Depth == MaxDepth-1 && nstairs > 2 {
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

	// Simellas
	g.Simellas = make(map[position]int)
	for i := 0; i < 5; i++ {
		pos := g.FreeCellForStatic()
		const rounds = 5
		for j := 0; j < rounds; j++ {
			g.Simellas[pos] += 1 + RandInt(g.Depth+g.Depth*g.Depth/10)
		}
		g.Simellas[pos] /= rounds
		if g.Simellas[pos] == 0 {
			g.Simellas[pos] = 1
		}
	}

	// initialize LOS
	if g.Depth == 0 {
		g.Print("You're in Hareka's Underground searching for medicinal simellas. Good luck!")
		g.PrintStyled("► Press ? for help or click on right upper →Menu button.", logSpecial)
	}
	if g.Depth == MaxDepth {
		g.PrintStyled("You feel magic in the air. The way out is close.", logSpecial)
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// Frundis is somewhere in the level
	if g.FrundisInLevel() {
		g.PrintStyled("You hear some faint music… ♫ larilon, larila ♫ ♪", logSpecial)
	}

	// recharge rods
	for r, props := range g.Player.Rods {
		if props.Charge < r.MaxCharge() {
			rchg := RandInt(1 + r.Rate())
			if rchg == 0 && RandInt(2) == 0 {
				rchg++
			}
			props.Charge += rchg
			g.Player.Rods[r] = props
		}
		if props.Charge > r.MaxCharge() {
			props.Charge = r.MaxCharge()
			g.Player.Rods[r] = props
		}
	}

	// clouds
	g.Clouds = map[position]cloud{}

	// Events
	if g.Depth == 0 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
		g.PushEvent(&simpleEvent{ERank: 50, EAction: HealPlayer})
		g.PushEvent(&simpleEvent{ERank: 100, EAction: MPRegen})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		g.PushEvent(&monsterEvent{ERank: g.Turn + 1, EAction: MonsterTurn, NMons: i})
		g.PushEvent(&monsterEvent{ERank: g.Turn + 50, EAction: HealMonster, NMons: i})
	}
}

func (g *game) CleanEvents() {
	evq := &eventQueue{}
	for g.Events.Len() > 0 {
		iev := g.PopIEvent()
		switch iev.Event.(type) {
		case *monsterEvent:
		case *cloudEvent:
		default:
			heap.Push(evq, iev)
		}
	}
	g.Events = evq
}

func (g *game) StairsSlice() []position {
	stairs := []position{}
	for stairPos, b := range g.Stairs {
		if b && g.Dungeon.Cell(stairPos).Explored {
			stairs = append(stairs, stairPos)
		}
	}
	return stairs
}

func (g *game) GenCollectables() {
	rounds := 50
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
				g.Collectables[pos] = collectable{Consumable: c, Quantity: data.quantity}
			}
		}
	}
}

func (g *game) SeenGoodWeapon() bool {
	return g.GeneratedEquipables[Sword] || g.GeneratedEquipables[DoubleSword] || g.GeneratedEquipables[Spear] || g.GeneratedEquipables[Halberd] ||
		g.GeneratedEquipables[Axe] || g.GeneratedEquipables[BattleAxe] || g.GeneratedEquipables[Frundis] || g.GeneratedEquipables[ElecWhip]
}

func (g *game) GenWeapon() {
	wps := [9]weapon{Dagger, Axe, BattleAxe, Spear, Halberd, Sword, DoubleSword, Frundis, ElecWhip}
	n := 11
	if !g.SeenGoodWeapon() {
		n -= 4 * g.Depth
		if n < 2 {
			n = 2
		}
	} else if g.Player.Weapon != Dagger {
		n *= 2
	}
	r := RandInt(n)
	if r != 0 && !g.SeenGoodWeapon() && g.Depth > 3 {
		r = RandInt(n)
	}
	if r != 0 {
		return
	}
loop:
	for {
		for i := 0; i < len(wps); i++ {
			if wps[i] == Frundis && g.GeneratedEquipables[Frundis] {
				// unique
				return
			}
			n := 30
			if wps[i].TwoHanded() && g.Depth < 3 {
				n *= (3 - g.Depth)
			}
			if wps[i] == Dagger {
				n *= 2
			}
			r := RandInt(n)
			if r == 0 {
				pos := g.FreeCellForStatic()
				g.Equipables[pos] = wps[i]
				g.GeneratedEquipables[wps[i]] = true
				break loop
			}
		}

	}
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
		default:
			// not reached
			return
		}
	}
	if r == 0 {
		pos := g.FreeCellForStatic()
		g.Equipables[pos] = eq
		g.GeneratedEquipables[eq] = true
	}

}

func (g *game) FrundisInLevel() bool {
	for _, eq := range g.Equipables {
		if wp, ok := eq.(weapon); ok && wp == Frundis {
			return true
		}
	}
	return false
}

func (g *game) Descend() bool {
	g.LevelStats()
	if g.Depth >= MaxDepth {
		g.Depth++
		return true
	}
	g.Print("You descend deeper in the dungeon.")
	g.Depth++
	g.Scumming = 0
	g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: PlayerTurn})
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) WizardMode() {
	g.Wizard = true
	g.Player.Consumables[DescentPotion] = 12
	g.PrintStyled("You are now in wizard mode and cannot obtain winner status.", logSpecial)
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		mons := g.MonsterInLOS()
		if mons == nil &&
			(g.Player.HP < g.Player.HPMax() || g.Player.MP < g.Player.MPMax() || g.Player.HasStatus(StatusExhausted) ||
				g.Player.HasStatus(StatusConfusion) || g.Player.HasStatus(StatusLignification)) {
			g.WaitTurn(ev)
			return true
		}
		if mons != nil && g.Resting {
			g.Stats.RestInterrupt++
		}
		g.Resting = false
	} else if g.Autoexploring {
		if g.ui.ExploreStep(g) {
			g.AutoHalt = true
			g.Print("Stopping, then.")
		}
		switch {
		case g.AutoHalt:
			// stop exploring
		default:
			var n *position
			var finished bool
			if g.DijkstraMapRebuild {
				if g.AllExplored() {
					g.Print("You finished exploring.")
					break
				}
				sources := g.AutoexploreSources()
				g.BuildAutoexploreMap(sources)
			}
			n, finished = g.NextAuto()
			if finished {
				n = nil
			}
			if finished && g.AllExplored() {
				g.Print("You finished exploring.")
			} else if n == nil {
				g.Print("You could not reach safely some places.")
			}
			if n != nil {
				err := g.MovePlayer(*n, ev)
				if err != nil {
					g.Print(err.Error())
					break
				}
				return true
			}
		}
		g.Autoexploring = false
	} else if g.AutoTarget.valid() {
		if !g.ui.ExploreStep(g) && g.MoveToTarget(ev) {
			return true
		}
	} else if g.AutoDir != NoDir {
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
				g.LevelStats()
				err := g.RemoveSaveFile()
				if err != nil {
					g.PrintfStyled("Error removing save file: %v", logError, err.Error())
				}
				g.ui.Death(g)
				break loop
			}
		}
		if g.Events.Len() == 0 {
			break loop
		}
		ev := g.PopIEvent().Event
		g.Turn = ev.Rank()
		g.Ev = ev
		ev.Action(g)
		if g.AutoNext {
			continue loop
		}
		if g.Quit {
			break loop
		}
	}
}
