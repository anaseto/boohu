package main

import "container/heap"

var Version string = "v0.1-dev"

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	MonstersPosCache    []int // monster (dungeon index + 1) / no monster (0)
	Bands               []bandInfo
	Events              *eventQueue
	Ev                  event
	EventIndex          int
	Depth               int
	ExploredLevels      int
	DepthPlayerTurn     int
	Turn                int
	Highlight           map[position]bool // highlighted positions (e.g. targeted ray)
	CollectableScore    int
	LastConsumables     []consumable
	Object              map[position]object
	Clouds              map[position]cloud
	Fungus              map[position]vegetation
	Doors               map[position]bool
	TemporalWalls       map[position]bool
	GeneratedUniques    map[monsterBand]int
	GeneratedEquipables map[equipable]bool
	GeneratedRods       map[rod]bool
	GenPlan             [MaxDepth + 1]genFlavour
	FoundEquipables     map[equipable]bool
	Simellas            map[position]int
	WrongWall           map[position]bool
	WrongFoliage        map[position]bool
	WrongDoor           map[position]bool
	ExclusionsMap       map[position]bool
	Noise               map[position]bool
	LastMonsterKnownAt  map[position]*monster
	MonsterLOS          map[position]bool
	Resting             bool
	RestingTurns        int
	Autoexploring       bool
	DijkstraMapRebuild  bool
	Targeting           position
	AutoTarget          position
	AutoDir             direction
	AutoHalt            bool
	AutoNext            bool
	DrawBuffer          []UICell
	drawBackBuffer      []UICell
	DrawLog             []drawFrame
	Log                 []logEntry
	LogIndex            int
	LogNextTick         int
	InfoEntry           string
	Stats               stats
	Boredom             int
	Quit                bool
	Wizard              bool
	WizardMap           bool
	Version             string
	//Opts                startOpts
	ui *gameui
}

//type startOpts struct {
//Alternate     monsterKind
//StoneLevel    int
//UnstableLevel int
//}

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

func (g *game) FreeCellForPlayer() position {
	center := position{DungeonWidth / 2, DungeonHeight / 2}
	bestpos := g.FreeCell()
	for i := 0; i < 2; i++ {
		pos := g.FreeCell()
		if pos.Distance(center) > bestpos.Distance(center) {
			bestpos = pos
		}
	}
	return bestpos
}

func (g *game) FreeCellForStair(dist int) position {
	iters := 0
	bestpos := g.Player.Pos
	for {
		pos := g.FreeCellForStatic()
		adjust := 0
		for i := 0; i < 4; i++ {
			adjust += RandInt(dist)
		}
		adjust /= 4
		if pos.Distance(g.Player.Pos) <= 6+adjust {
			continue
		}
		iters++
		if pos.Distance(g.Player.Pos) > bestpos.Distance(g.Player.Pos) {
			bestpos = pos
		}
		if iters == 2 {
			return bestpos
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
		if _, ok := g.Object[pos]; ok {
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

const MaxDepth = 11
const WinDepth = 8

const (
	DungeonHeight = 21
	DungeonWidth  = 79
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	g.GenRoomTunnels()
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        DefaultHealth,
		MP:        DefaultMPmax,
		Simellas:  0,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Consumables = map[consumable]int{
		HealWoundsPotion: 1,
	}
	switch RandInt(7) {
	case 0:
		g.Player.Consumables[ExplosiveMagara] = 1
	case 1:
		g.Player.Consumables[NightMagara] = 1
	case 2:
		g.Player.Consumables[TeleportMagara] = 1
	case 3:
		g.Player.Consumables[SlowingMagara] = 1
	case 4:
		g.Player.Consumables[ConfuseMagara] = 1
	default:
		g.Player.Consumables[ConfusingDart] = 2
	}
	switch RandInt(11) {
	case 0, 1, 5:
		g.Player.Consumables[TeleportationPotion] = 1
	case 2, 3, 4:
		g.Player.Consumables[SwiftnessPotion] = 1
	case 6:
		g.Player.Consumables[WallPotion] = 1
	case 7:
		g.Player.Consumables[CBlinkPotion] = 1
	case 8:
		g.Player.Consumables[DigPotion] = 1
	case 9:
		g.Player.Consumables[SwapPotion] = 1
	case 10:
		g.Player.Consumables[ShadowsPotion] = 1
	}
	r := g.RandomRod()
	g.Player.Rods = map[rod]rodProps{r: rodProps{r.MaxCharge() - 1}}
	g.Player.Statuses = map[status]int{}
	g.Player.Expire = map[status]int{}

	// Testing
	//g.Player.Aptitudes[AptStealthyLOS] = true
	//g.Player.Aptitudes[AptStealthyMovement] = true
	//g.Player.Rods[RodSwapping] = rodProps{Charge: 3}
	//g.Player.Rods[RodFireball] = rodProps{Charge: 3}
	//g.Player.Rods[RodLightning] = rodProps{Charge: 3}
	//g.Player.Rods[RodLightningBolt] = rodProps{Charge: 3}
	//g.Player.Rods[RodShatter] = rodProps{Charge: 3}
	//g.Player.Rods[RodFog] = rodProps{Charge: 3}
	//g.Player.Rods[RodSleeping] = rodProps{Charge: 3}
	//g.Player.Consumables[BerserkPotion] = 5
	//g.Player.Consumables[MagicMappingPotion] = 1
	//g.Player.Consumables[ExplosiveMagara] = 5
	//g.Player.Consumables[NightMagara] = 5
	//g.Player.Consumables[SlowingMagara] = 5
	//g.Player.Consumables[ConfuseMagara] = 5
	//g.Player.Consumables[DigPotion] = 5
	//g.Player.Consumables[SwapPotion] = 5
	//g.Player.Consumables[DreamPotion] = 5
	//g.Player.Consumables[ShadowsPotion] = 5
	//g.Player.Consumables[TormentPotion] = 5
	//g.Player.Consumables[AccuracyPotion] = 5
	//g.Player.Weapon = ElecWhip
	//g.Player.Weapon = DancingRapier
	//g.Player.Weapon = Sabre
	//g.Player.Weapon = HarKarGauntlets
	//g.Player.Weapon = DefenderFlail
	//g.Player.Weapon = HopeSword
	//g.Player.Weapon = DragonSabre
	//g.Player.Weapon = FinalBlade
	//g.Player.Weapon = VampDagger
	//g.Player.Shield = EarthShield
	//g.Player.Shield = FireShield
	//g.Player.Shield = BashingShield
	//g.Player.Armour = TurtlePlates
	//g.Player.Armour = HarmonistRobe
	//g.Player.Armour = CelmistRobe
	//g.Player.Armour = ShinyPlates
	//g.Player.Armour = SmokingScales
}

type genFlavour int

const (
	GenRod genFlavour = iota
	//GenWeapon
	GenArmour
	GenWpArm
	GenExtraCollectables
)

func (g *game) InitFirstLevel() {
	g.Depth++ // start at 1
	g.InitPlayer()
	g.AutoTarget = InvalidPos
	g.Targeting = InvalidPos
	g.GeneratedRods = map[rod]bool{}
	g.GeneratedEquipables = map[equipable]bool{}
	g.FoundEquipables = map[equipable]bool{Robe: true, Dagger: true}
	g.GeneratedUniques = map[monsterBand]int{}
	g.Stats.KilledMons = map[monsterKind]int{}
	g.Version = Version
	g.GenPlan = [MaxDepth + 1]genFlavour{
		1:  GenRod,
		2:  GenArmour,
		3:  GenExtraCollectables,
		4:  GenRod,
		5:  GenExtraCollectables,
		6:  GenRod,
		7:  GenExtraCollectables,
		8:  GenExtraCollectables,
		9:  GenRod,
		10: GenExtraCollectables,
		11: GenExtraCollectables,
	}
	permi := RandInt(7)
	switch permi {
	case 0, 1, 2, 3:
		g.GenPlan[permi+1], g.GenPlan[permi+2] = g.GenPlan[permi+2], g.GenPlan[permi+1]
	}
	if RandInt(4) == 0 {
		g.GenPlan[6], g.GenPlan[7] = g.GenPlan[7], g.GenPlan[6]
	}
}

func (g *game) InitLevel() {
	// Starting data
	if g.Depth == 0 {
		g.InitFirstLevel()
	}

	g.MonstersPosCache = make([]int, DungeonNCells)
	g.WrongWall = map[position]bool{}
	g.WrongFoliage = map[position]bool{}
	g.WrongDoor = map[position]bool{}
	g.ExclusionsMap = map[position]bool{}
	g.TemporalWalls = map[position]bool{}
	g.LastMonsterKnownAt = map[position]*monster{}
	g.Object = make(map[position]object)

	// Dungeon terrain
	g.GenDungeon()

	// Aptitudes/Mutations
	if g.Depth == 2 || g.Depth == 5 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Magical Stones
	nstones := 1
	switch RandInt(8) {
	case 0:
		nstones = 0
	case 1, 2, 3:
		nstones = 2
	case 4, 5, 6:
		nstones = 3
	}
	ustone := stone(0)
	for i := 0; i < nstones; i++ {
		pos := g.FreeCellForStatic()
		var st stone
		if ustone != stone(0) {
			st = ustone
		} else {
			st = stone(1 + RandInt(NumStones-1))
		}
		g.Object[pos] = st
	}

	// Simellas
	g.Simellas = make(map[position]int)
	for i := 0; i < 5; i++ {
		pos := g.FreeCellForStatic()
		const rounds = 5
		for j := 0; j < rounds; j++ {
			g.Simellas[pos] += 1 + RandInt(g.Depth+g.Depth*g.Depth/6)
		}
		g.Simellas[pos] /= rounds
		if g.Simellas[pos] == 0 {
			g.Simellas[pos] = 1
		}
	}

	// initialize LOS
	if g.Depth == 1 {
		g.Print("You're in Hareka's Underground searching for medicinal simellas. Good luck!")
		g.PrintStyled("► Press ? for help on keys or use the mouse and [buttons].", logSpecial)
	}
	if g.Depth == WinDepth {
		g.PrintStyled("You feel magic in the air. A first way out is close!", logSpecial)
	} else if g.Depth == MaxDepth {
		g.PrintStyled("If rumors are true, you have reached the bottom!", logSpecial)
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// Frundis is somewhere in the level
	if g.FrundisInLevel() {
		g.PrintStyled("You hear some faint music… ♫ larilon, larila ♫ ♪", logSpecial)
	}

	// recharge rods
	if g.Depth > 1 {
		g.RechargeRods()
	}

	// clouds
	g.Clouds = map[position]cloud{}

	// Events
	if g.Depth == 1 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		g.PushEvent(&monsterEvent{ERank: g.Turn + RandInt(10), EAction: MonsterTurn, NMons: i})
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
	for pos, obj := range g.Object {
		_, ok := obj.(stair)
		if ok && g.Dungeon.Cell(pos).Explored {
			stairs = append(stairs, pos)
		}
	}
	return stairs
}

func (dg *dgen) GenCollectable(g *game) {
	rounds := 100
	if len(g.LastConsumables) > 3 {
		g.LastConsumables = g.LastConsumables[1:]
	}
	for {
	loopcons:
		for c, data := range ConsumablesCollectData {
			r := RandInt(data.rarity * rounds)
			if r != 0 {
				continue
			}

			// avoid too many of the same
			for _, co := range g.LastConsumables {
				if co == c && RandInt(4) > 0 {
					continue loopcons
				}
			}
			g.LastConsumables = append(g.LastConsumables, c)
			g.CollectableScore++
			pos := InvalidPos
			for pos == InvalidPos {
				pos = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(PlaceItem)
			}
			g.Object[pos] = collectable{Consumable: c, Quantity: data.quantity}
			return
		}
	}

}

func (dg *dgen) GenCollectables(g *game) {
	score := g.CollectableScore - 2*(g.Depth-1)
	n := 2
	if score >= 0 && RandInt(4) == 0 {
		n--
	}
	if score <= 0 && RandInt(4) == 0 {
		n++
	}
	if score > 0 && n >= 2 {
		n--
	}
	if score < 0 && n <= -2 {
		n++
	}
	for i := 0; i < n; i++ {
		dg.GenCollectable(g)
	}
}

func (g *game) GenArmour() {
	ars := [3]armour{SmokingScales, CelmistRobe, HarmonistRobe}
	for {
		i := RandInt(len(ars))
		if g.GeneratedEquipables[ars[i]] {
			// do not generate duplicates
			continue
		}
		pos := g.FreeCellForStatic()
		g.Object[pos] = ars[i]
		g.GeneratedEquipables[ars[i]] = true
		break
	}
}

func (g *game) GenWeapon() {
	wps := [4]weapon{DancingRapier, Frundis, HarKarGauntlets, DefenderFlail}
	for {
		i := RandInt(len(wps))
		if g.GeneratedEquipables[wps[i]] {
			// do not generate duplicates
			continue
		}
		pos := g.FreeCellForStatic()
		g.Object[pos] = wps[i]
		g.GeneratedEquipables[wps[i]] = true
		break
	}
}

func (g *game) FrundisInLevel() bool {
	for _, obj := range g.Object {
		eq, ok := obj.(equipable)
		if !ok {
			continue
		}
		if wp, ok := eq.(weapon); ok && wp == Frundis {
			return true
		}
	}
	return false
}

func (g *game) Descend() bool {
	g.LevelStats()
	if obj, ok := g.Object[g.Player.Pos]; ok {
		if strt, ok := obj.(stair); ok && strt == WinStair {
			g.StoryPrint("You escaped!")
			g.ExploredLevels = g.Depth
			g.Depth = -1
			return true
		}
	}
	g.Print("You descend deeper in the dungeon.")
	g.StoryPrint("You descended deeper in the dungeon.")
	g.Depth++
	g.DepthPlayerTurn = 0
	g.Boredom = 0
	g.PushEvent(&simpleEvent{ERank: g.Ev.Rank(), EAction: PlayerTurn})
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) WizardMode() {
	g.Wizard = true
	g.Player.Consumables[DescentPotion] = 15
	g.PrintStyled("You are now in wizard mode and cannot obtain winner status.", logSpecial)
}

func (g *game) ApplyRest() {
	g.Player.HP = g.Player.HPMax()
	g.Player.HPbonus = 0
	g.Player.MP = g.Player.MPMax()
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			continue
		}
		mons.HP = mons.HPmax
	}
	adjust := 0
	if g.Player.Armour == HarmonistRobe {
		// the harmonist robe mitigates the sound of your snorts
		adjust = 100
	}
	if g.DepthPlayerTurn < 100+adjust && RandInt(5) > 2 || g.DepthPlayerTurn >= 100+adjust && g.DepthPlayerTurn < 250+adjust && RandInt(2) == 0 ||
		g.DepthPlayerTurn >= 250+adjust && RandInt(3) > 0 {
		rmons := []int{}
		for i, mons := range g.Monsters {
			if mons.Exists() && mons.State == Resting {
				rmons = append(rmons, i)
			}
		}
		if len(rmons) > 0 {
			g.Monsters[rmons[RandInt(len(rmons))]].NaturalAwake(g)
		}
	}
	g.Stats.Rest++
	g.PrintStyled("You feel fresh again. Some monsters might have awoken.", logStatusEnd)
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		const enoughRestTurns = 15
		mons := g.MonsterInLOS()
		sr := g.StatusRest()
		if mons == nil && (sr || g.NeedsRegenRest() && g.RestingTurns >= 0) && g.RestingTurns < enoughRestTurns {
			g.WaitTurn(ev)
			if !sr && g.RestingTurns >= 0 {
				g.RestingTurns++
			}
			return true
		}
		if g.RestingTurns >= enoughRestTurns {
			g.ApplyRest()
		} else if mons != nil {
			g.Stats.RestInterrupt++
			g.Print("You could not sleep.")
		}
		g.Resting = false
	} else if g.Autoexploring {
		if g.ui.ExploreStep() {
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
				g.Print("You could not safely reach some places.")
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
		if !g.ui.ExploreStep() && g.MoveToTarget(ev) {
			return true
		} else {
			g.AutoTarget = InvalidPos
		}
	} else if g.AutoDir != NoDir {
		if !g.ui.ExploreStep() && g.AutoToDir(ev) {
			return true
		} else {
			g.AutoDir = NoDir
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
				g.ui.Death()
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
