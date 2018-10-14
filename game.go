package main

import "container/heap"

var Version string = "v0.11-dev"

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	MonstersPosCache    []int // monster (dungeon index + 1) / no monster (0)
	Bands               []monsterBand
	BandData            []monsterBandData
	Events              *eventQueue
	Ev                  event
	EventIndex          int
	Depth               int
	ExploredLevels      int
	DepthPlayerTurn     int
	Turn                int
	Highlight           map[position]bool // highlighted positions (e.g. targeted ray)
	Collectables        map[position]collectable
	CollectableScore    int
	UnstableLevel       int
	StoneLevel          int
	Equipables          map[position]equipable
	Rods                map[position]rod
	Stairs              map[position]stair
	Clouds              map[position]cloud
	Fungus              map[position]vegetation
	Doors               map[position]bool
	TemporalWalls       map[position]bool
	MagicalStones       map[position]stone
	GeneratedUniques    map[monsterBand]int
	SpecialBands        map[int][]monsterBandData
	GeneratedEquipables map[equipable]bool
	GeneratedRods       map[rod]bool
	FoundEquipables     map[equipable]bool
	Simellas            map[position]int
	WrongWall           map[position]bool
	WrongFoliage        map[position]bool
	WrongDoor           map[position]bool
	ExclusionsMap       map[position]bool
	Noise               map[position]bool
	DreamingMonster     map[position]bool
	Resting             bool
	RestingTurns        int
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
	Boredom             int
	Quit                bool
	Wizard              bool
	WizardMap           bool
	Version             string
	ui                  Renderer
}

type Renderer interface {
	ExploreStep(*game) bool
	HandlePlayerTurn(*game, event) bool
	Death(*game)
	ChooseTarget(*game, Targeter) error
	CriticalHPWarning(*game)
	ExplosionAnimation(*game, explosionStyle, position)
	TormentExplosionAnimation(*game)
	LightningBoltAnimation(*game, []position)
	ThrowAnimation(*game, []position, bool)
	MonsterJavelinAnimation(*game, []position, bool)
	MonsterProjectileAnimation(*game, []position, rune, uicolor)
	DrinkingPotionAnimation(*game)
	SwappingAnimation(*game, position, position)
	TeleportAnimation(*game, position, position, bool)
	MagicMappingAnimation(*game, []int)
	HitAnimation(*game, position, bool)
	LightningHitAnimation(*game, []position)
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

func (g *game) FreeCellForStair(dist int) position {
	for {
		pos := g.FreeCellForStatic()
		adjust := 0
		for i := 0; i < 4; i++ {
			adjust += RandInt(dist)
		}
		adjust /= 4
		if pos.Distance(g.Player.Pos) > 6+adjust {
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
		if _, ok := g.Stairs[pos]; ok {
			continue
		}
		if _, ok := g.Rods[pos]; ok {
			continue
		}
		if _, ok := g.Equipables[pos]; ok {
			continue
		}
		if _, ok := g.MagicalStones[pos]; ok {
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

const MaxDepth = 15
const WinDepth = 12

const (
	DungeonHeight = 21
	DungeonWidth  = 79
	DungeonNCells = DungeonWidth * DungeonHeight
)

func (g *game) GenDungeon() {
	g.Fungus = make(map[position]vegetation)
	dg := GenRuinsMap
	switch RandInt(7) {
	//switch 4 {
	case 0:
		dg = GenCaveMap
	case 1:
		dg = GenRoomMap
	case 2:
		dg = GenCellularAutomataCaveMap
	case 3:
		dg = GenCaveMapTree
	case 4:
		dg = GenBSPMap
	}
	dg.Use(g)
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        42,
		MP:        3,
		Simellas:  0,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Consumables = map[consumable]int{
		HealWoundsPotion: 1,
	}
	switch RandInt(4) {
	case 0:
		g.Player.Consumables[ExplosiveMagara] = 1
	case 1:
		g.Player.Consumables[NightMagara] = 1
	default:
		g.Player.Consumables[ConfusingDart] = 3
	}
	switch RandInt(12) {
	case 0, 1:
		g.Player.Consumables[TeleportationPotion] = 1
	case 2, 3, 4:
		g.Player.Consumables[SwiftnessPotion] = 1
	case 6:
		g.Player.Consumables[WallPotion] = 1
	case 7:
		g.Player.Consumables[CBlinkPotion] = 1
	case 5, 8:
		g.Player.Consumables[DigPotion] = 1
	case 9:
		g.Player.Consumables[SwapPotion] = 1
	case 10:
		g.Player.Consumables[ShadowsPotion] = 1
	case 11:
		g.Player.Consumables[ConfusePotion] = 1
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
	//g.Player.Consumables[DigPotion] = 5
	//g.Player.Consumables[SwapPotion] = 5
	//g.Player.Consumables[DreamPotion] = 5
	//g.Player.Consumables[ShadowsPotion] = 5
	//g.Player.Consumables[ConfusePotion] = 5
	//g.Player.Consumables[TormentPotion] = 5
	//g.Player.Shield = EarthShield
	//g.Player.Shield = FireShield
	//g.Player.Armour = TurtlePlates
	//g.Player.Armour = HarmonistRobe
	//g.Player.Armour = CelmistRobe
	//g.Player.Armour = ShinyPlates
	//g.Player.Armour = SmokingScales
}

func (g *game) InitSpecialBands() {
	g.SpecialBands = map[int][]monsterBandData{}
	sb := MonsSpecialBands[RandInt(len(MonsSpecialBands))]
	depth := sb.minDepth + RandInt(sb.maxDepth-sb.minDepth+1)
	g.SpecialBands[depth] = sb.bands
	seb := MonsSpecialEndBands[RandInt(len(MonsSpecialEndBands))]
	if RandInt(4) == 0 {
		if RandInt(5) > 1 {
			g.SpecialBands[13] = seb.bands
		} else {
			g.SpecialBands[12] = seb.bands
		}
	} else if RandInt(5) > 0 {
		if RandInt(3) > 0 {
			g.SpecialBands[15] = seb.bands
		} else {
			g.SpecialBands[14] = seb.bands
		}
	}
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
		g.FoundEquipables = map[equipable]bool{Robe: true}
		g.GeneratedUniques = map[monsterBand]int{}
		g.Stats.KilledMons = map[monsterKind]int{}
		g.InitSpecialBands()
		if RandInt(3) > 0 {
			g.UnstableLevel = 1 + RandInt(15)
		}
		if RandInt(2) == 0 || RandInt(2) == 0 && g.UnstableLevel == 0 {
			g.StoneLevel = 1 + RandInt(15)
		}
		g.Version = Version
	}

	g.MonstersPosCache = make([]int, DungeonNCells)
	g.Player.Pos = g.FreeCell()

	g.WrongWall = map[position]bool{}
	g.WrongFoliage = map[position]bool{}
	g.WrongDoor = map[position]bool{}
	g.ExclusionsMap = map[position]bool{}
	g.TemporalWalls = map[position]bool{}
	g.DreamingMonster = map[position]bool{}

	// Monsters
	g.BandData = MonsBands
	if bd, ok := g.SpecialBands[g.Depth]; ok {
		g.BandData = bd
	}
	g.GenMonsters()

	// Collectables
	g.Collectables = make(map[position]collectable)
	g.GenCollectables()

	// Equipment
	g.Equipables = make(map[position]equipable)
	g.GenArmour()
	g.GenShield()

	// Rods
	g.Rods = map[position]rod{}
	r := g.Depth - 3*g.GeneratedRodsCount()
	if r > 0 && RandInt((6-r)*3) == 0 && g.GeneratedRodsCount() < 4 ||
		g.GeneratedRodsCount() == 0 && g.Depth > 0 ||
		g.GeneratedRodsCount() == 1 && g.Depth > 4 ||
		g.GeneratedRodsCount() == 2 && g.Depth > 8 ||
		g.GeneratedRodsCount() == 3 && g.Depth > 11 {
		g.GenerateRod()
	}

	// Aptitudes/Mutations
	r = 15 + 3*g.Player.AptitudeCount() - g.Depth
	if RandInt(r) == 0 && g.Depth > 0 && g.Player.AptitudeCount() < 2 ||
		g.Player.AptitudeCount() == 0 && g.Depth > 1 ||
		g.Player.AptitudeCount() == 1 && g.Depth > 5 {
		//g.Player.AptitudeCount() == 2 && g.Depth > 8 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Stairs
	g.Stairs = make(map[position]stair)
	nstairs := 2
	if RandInt(3) == 0 {
		if RandInt(2) == 0 {
			nstairs++
		} else {
			nstairs--
		}
	}
	if g.Depth >= WinDepth {
		nstairs = 1
	} else if g.Depth == WinDepth-1 && nstairs > 2 {
		nstairs = 2
	}
	for i := 0; i < nstairs; i++ {
		var pos position
		if g.Depth >= WinDepth && g.Depth != 14 {
			pos = g.FreeCellForStair(60)
			g.Stairs[pos] = WinStair
		}
		if g.Depth < MaxDepth {
			if g.Depth > 9 {
				pos = g.FreeCellForStair(40)
			} else {
				pos = g.FreeCellForStair(0)
			}
			g.Stairs[pos] = NormalStair
		}
	}

	// Magical Stones
	g.MagicalStones = map[position]stone{}
	nstones := 1
	switch RandInt(7) {
	case 0:
		nstones = 0
	case 1, 2:
		nstones = 2
	}
	ustone := stone(0)
	if g.Depth > 0 && g.Depth == g.StoneLevel {
		ustone = stone(1 + RandInt(NumStones-1))
		nstones = 10 + g.Depth/2
		if RandInt(4) == 0 {
			g.StoneLevel = g.StoneLevel + RandInt(MaxDepth-g.StoneLevel) + 1
		}
	}
	for i := 0; i < nstones; i++ {
		pos := g.FreeCellForStatic()
		var st stone
		if ustone != stone(0) {
			st = ustone
		} else {
			st = stone(1 + RandInt(NumStones-1))
		}
		g.MagicalStones[pos] = st
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
		g.PrintStyled("► Press ? for help on keys or use the mouse and [buttons].", logSpecial)
	}
	if g.Depth == WinDepth {
		g.PrintStyled("You feel magic in the air. A first way out is close!", logSpecial)
	} else if g.Depth == MaxDepth {
		g.PrintStyled("If rumors are true, you have reached the bottom!", logSpecial)
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// recharge rods
	if g.Depth > 0 {
		g.RechargeRods()
	}

	// clouds
	g.Clouds = map[position]cloud{}

	// Events
	if g.Depth == 0 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		g.PushEvent(&simpleEvent{ERank: 0, EAction: PlayerTurn})
	} else {
		g.CleanEvents()
	}
	for i := range g.Monsters {
		g.PushEvent(&monsterEvent{ERank: g.Turn + RandInt(10), EAction: MonsterTurn, NMons: i})
	}
	if g.Depth > 0 && g.Depth == g.UnstableLevel {
		g.PrintStyled("You sense magic instability on this level.", logSpecial)
		for i := 0; i < 15; i++ {
			g.PushEvent(&cloudEvent{ERank: g.Turn + 100 + RandInt(900), EAction: ObstructionProgression})
		}
		if RandInt(4) == 0 {
			g.UnstableLevel = g.UnstableLevel + RandInt(MaxDepth-g.UnstableLevel) + 1
		}
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
	for stairPos, _ := range g.Stairs {
		if g.Dungeon.Cell(stairPos).Explored {
			stairs = append(stairs, stairPos)
		}
	}
	return stairs
}

func (g *game) GenCollectables() {
	rounds := 100
	for i := 0; i < rounds; i++ {
		for c, data := range ConsumablesCollectData {
			var r int
			dfactor := g.Depth + 1
			if g.Depth >= WinDepth {
				// more items in last levels
				dfactor += g.Depth - WinDepth + 1
			}
			if g.CollectableScore >= 5*dfactor/3 {
				r = RandInt(data.rarity * rounds * 4)
			} else if g.CollectableScore < 4*dfactor/3 {
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

func (g *game) SeenGoodArmour() (count int) {
	for eq, b := range g.GeneratedEquipables {
		ar, ok := eq.(armour)
		if ok && b && ar != Robe && ar != LeatherArmour {
			count++
		}
	}
	return count
}

func (g *game) SeenGoodShield() (count int) {
	for eq, b := range g.GeneratedEquipables {
		sh, ok := eq.(shield)
		if ok && b && sh != Buckler {
			count++
		}
	}
	return count
}

func (g *game) GenShield() {
	ars := [5]shield{Buckler, ConfusingShield, BashingShield, EarthShield, FireShield}
	n := 12 + 5*g.SeenGoodShield()
	if g.SeenGoodShield() == 2 {
		return
	}
	if g.SeenGoodShield() == 0 {
		n -= 2 * g.Depth
		if n < 2 {
			if g.Depth < 8 {
				n = 2
			} else {
				n = 1
			}
		}
	} else if g.SeenGoodShield() == 1 {
		n -= 4 * (g.Depth - 7)
		if n < 2 {
			if g.Depth < 12 {
				n = 2
			} else {
				n = 1
			}
		}
	} else if g.Player.Shield != NoShield && g.Player.Shield != Buckler {
		n += 10
	} else if g.Depth > WinDepth {
		n = 2
	}
	r := RandInt(n)
	if r != 0 {
		if !g.GeneratedEquipables[Buckler] && (g.Depth > 0 && RandInt(2) == 0 || g.Depth > 3) {
			pos := g.FreeCellForStatic()
			g.Equipables[pos] = Buckler
			g.GeneratedEquipables[Buckler] = true
		}
		return
	}
loop:
	for {
		for i := 0; i < len(ars); i++ {
			if g.GeneratedEquipables[ars[i]] {
				// do not generate duplicates
				continue
			}
			n := 50
			r := RandInt(n)
			if r == 0 {
				pos := g.FreeCellForStatic()
				g.Equipables[pos] = ars[i]
				g.GeneratedEquipables[ars[i]] = true
				break loop
			}
		}
	}
}

func (g *game) GenArmour() {
	ars := [8]armour{Robe, LeatherArmour, SmokingScales, ShinyPlates, TurtlePlates, SpeedRobe, CelmistRobe, HarmonistRobe}
	n := 11 + 5*g.SeenGoodArmour()
	if g.SeenGoodArmour() > 2 {
		return
	}
	if g.SeenGoodArmour() == 0 {
		n -= 2 * g.Depth
		if n < 2 {
			if g.Depth < 7 {
				n = 2
			} else {
				n = 1
			}
		}
	} else if g.SeenGoodArmour() == 1 {
		n -= 4 * (g.Depth - 7)
		if n < 2 {
			if g.Depth < 12 {
				n = 2
			} else {
				n = 1
			}
		}
	} else if g.Player.Armour != Robe && g.Player.Armour != LeatherArmour {
		n += 10
	} else if g.Depth > WinDepth {
		n = 2
	}
	r := RandInt(n)
	if r != 0 {
		if !g.GeneratedEquipables[LeatherArmour] && (RandInt(2) == 0 || g.Depth >= 3) {
			pos := g.FreeCellForStatic()
			g.Equipables[pos] = LeatherArmour
			g.GeneratedEquipables[LeatherArmour] = true
		}
		return
	}
loop:
	for {
		for i := 0; i < len(ars); i++ {
			if g.GeneratedEquipables[ars[i]] {
				// do not generate duplicates
				continue
			}
			n := 50
			if ars[i] == Robe {
				n *= 2
			}
			r := RandInt(n)
			if r == 0 {
				pos := g.FreeCellForStatic()
				g.Equipables[pos] = ars[i]
				g.GeneratedEquipables[ars[i]] = true
				break loop
			}
		}
	}
}

func (g *game) Descend() bool {
	g.LevelStats()
	if strt, ok := g.Stairs[g.Player.Pos]; ok && strt == WinStair {
		g.StoryPrint("You escaped!")
		g.ExploredLevels = g.Depth + 1
		g.Depth = -1
		return true
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
		if !g.ui.ExploreStep(g) && g.MoveToTarget(ev) {
			return true
		} else {
			g.AutoTarget = InvalidPos
		}
	} else if g.AutoDir != NoDir {
		if !g.ui.ExploreStep(g) && g.AutoToDir(ev) {
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
