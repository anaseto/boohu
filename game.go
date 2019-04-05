package main

import "container/heap"

var Version string = "v0.1-dev"

type game struct {
	Dungeon            *dungeon
	Player             *player
	Monsters           []*monster
	MonstersPosCache   []int // monster (dungeon index + 1) / no monster (0)
	Bands              []bandInfo
	Events             *eventQueue
	Ev                 event
	EventIndex         int
	Depth              int
	ExploredLevels     int
	DepthPlayerTurn    int
	Turn               int
	Highlight          map[position]bool // highlighted positions (e.g. targeted ray)
	Objects            objects
	Clouds             map[position]cloud
	TemporalWalls      map[position]terrain
	GeneratedUniques   map[monsterBand]int
	GeneratedLore      map[int]bool
	GeneratedMagaras   []magara
	GenPlan            [MaxDepth + 1]genFlavour
	TerrainKnowledge   map[position]terrain
	ExclusionsMap      map[position]bool
	Noise              map[position]bool
	NoiseIllusion      map[position]bool
	LastMonsterKnownAt map[position]*monster
	MonsterLOS         map[position]bool
	MonsterTargLOS     map[position]bool
	Illuminated        map[position]bool
	Resting            bool
	RestingTurns       int
	Autoexploring      bool
	DijkstraMapRebuild bool
	Targeting          position
	AutoTarget         position
	AutoDir            direction
	AutoHalt           bool
	AutoNext           bool
	DrawBuffer         []UICell
	drawBackBuffer     []UICell
	DrawLog            []drawFrame
	Log                []logEntry
	LogIndex           int
	LogNextTick        int
	InfoEntry          string
	Stats              stats
	Boredom            int
	Quit               bool
	Wizard             bool
	WizardMap          bool
	Version            string
	Places             places
	Params             startParams
	//Opts                startOpts
	ui *gameui
}

type startParams struct {
	Lore map[int]bool
}

type places struct {
	Shaedra  position
	Monolith position
	Marevor  position
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
		if !c.IsPassable() {
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
		if !c.IsPassable() {
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
	ml := AutomataCave
	//ml := RandomWalkTreeCave
	switch g.Depth {
	case 2, 6, 7:
		ml = RandomWalkCave
	case 4, 10, 11:
		ml = RandomWalkTreeCave
	}
	g.GenRoomTunnels(ml)
}

func (g *game) InitPlayer() {
	g.Player = &player{
		HP:        DefaultHealth,
		MP:        DefaultMPmax,
		Bananas:   2,
		Aptitudes: map[aptitude]bool{},
	}
	g.Player.Statuses = map[status]int{}
	g.Player.Expire = map[status]int{}
	g.Player.Magaras = []magara{
		NoMagara,
		NoMagara,
		NoMagara,
	}
	g.GeneratedMagaras = []magara{}
	for i := 0; i < 2; i++ {
		g.Player.Magaras[i] = g.RandomMagara()
		g.GeneratedMagaras = append(g.GeneratedMagaras, g.Player.Magaras[i])
	}
	// Testing
	//g.Player.Magaras[2] = NoiseMagara
	//g.Player.Magaras[2] = SlowingMagara
	//g.Player.Magaras[2] = ConfusionMagara
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
	g.Version = Version
	g.Depth++ // start at 1
	g.InitPlayer()
	g.AutoTarget = InvalidPos
	g.Targeting = InvalidPos
	g.GeneratedUniques = map[monsterBand]int{}
	g.GeneratedLore = map[int]bool{}
	g.Stats.KilledMons = map[monsterKind]int{}
	g.GenPlan = [MaxDepth + 1]genFlavour{ // XXX this is obsolete
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
	g.Params.Lore = map[int]bool{}
	for i := 0; i < 4; i++ {
		g.Params.Lore[RandInt(MaxDepth)] = true
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

func (g *game) InitLevelStructures() {
	g.MonstersPosCache = make([]int, DungeonNCells)
	g.TerrainKnowledge = map[position]terrain{}
	g.ExclusionsMap = map[position]bool{}
	g.TemporalWalls = map[position]terrain{}
	g.LastMonsterKnownAt = map[position]*monster{}
	g.Objects.Magaras = map[position]magara{}
	g.Objects.Lore = map[position]int{}
	g.Clouds = map[position]cloud{}
}

func (g *game) InitLevel() {
	// Starting data
	if g.Depth == 0 {
		g.InitFirstLevel()
	}

	g.InitLevelStructures()

	// Dungeon terrain
	g.GenDungeon()

	// Aptitudes/Mutations
	if g.Depth == 2 || g.Depth == 5 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Magara slots
	if g.Depth == 3 || g.Depth == 6 {
		g.Player.Magaras = append(g.Player.Magaras, NoMagara)
		g.PrintStyled("You have a new empty slot for a magara.", logSpecial)
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
	// TODO: use cache?
	stairs := []position{}
	for i, c := range g.Dungeon.Cells {
		if c.T != StairCell || !c.Explored {
			continue
		}
		pos := idxtopos(i)
		stairs = append(stairs, pos)
	}
	return stairs
}

func (g *game) Descend() bool {
	g.LevelStats()
	c := g.Dungeon.Cell(g.Player.Pos)
	if c.T == StairCell && g.Objects.Stairs[g.Player.Pos] == WinStair {
		g.StoryPrint("You escaped!")
		g.ExploredLevels = g.Depth
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
	//g.Player.Consumables[DescentPotion] = 15
	g.PrintStyled("You are now in wizard mode and cannot obtain winner status.", logSpecial)
}

func (g *game) ApplyRest() {
	g.Player.HP = g.Player.HPMax()
	g.Player.HPbonus = 0
	g.Player.MP = g.Player.MPMax()
	g.Stats.Rest++
	g.PrintStyled("You feel fresh again after eating banana and sleeping.", logStatusEnd)
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		const enoughRestTurns = 25
		if g.RestingTurns < enoughRestTurns {
			g.WaitTurn(ev)
			g.RestingTurns++
			return true
		}
		if g.RestingTurns >= enoughRestTurns {
			g.ApplyRest()
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
