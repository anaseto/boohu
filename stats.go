package main

type stats struct {
	Story             []string
	Killed            int
	KilledMons        map[monsterKind]int
	Moves             int
	Jumps             int
	ReceivedHits      int
	Dodges            int
	MagarasUsed       int
	DMagaraUses       [MaxDepth + 1]int
	UsedStones        int
	UsedMagaras       map[magara]int
	Damage            int
	DDamage           [MaxDepth + 1]int
	DExplPerc         [MaxDepth + 1]int
	DSleepingPerc     [MaxDepth + 1]int
	DKilledPerc       [MaxDepth + 1]int
	Burns             int
	Digs              int
	Rest              int
	DRests            [MaxDepth + 1]int
	Turns             int
	TWounded          int
	TMWounded         int
	TMonsLOS          int
	NSpotted          int
	NUSpotted         int
	DSpotted          [MaxDepth + 1]int
	DUSpotted         [MaxDepth + 1]int
	DUSpottedPerc     [MaxDepth + 1]int
	Achievements      map[achievement]bool
	AtNotablePos      map[position]bool
	HarmonicMagUse    int
	OricMagUse        int
	FireUse           int
	DestructionUse    int
	OricTelUse        int
	ClimbedTree       int
	TableHides        int
	HoledWallsCrawled int
	DoorsOpened       int
	BarrelHides       int
	Extinguishments   int
	Lore              map[int]bool
}

func (g *game) TurnStats() {
	g.Stats.Turns++
	g.DepthPlayerTurn++
	if g.Player.HP < g.Player.HPMax() {
		g.Stats.TWounded++
	}
	if g.MonsterInLOS() != nil {
		g.Stats.TMonsLOS++
		if g.Player.HP < g.Player.HPMax() {
			g.Stats.TMWounded++
		}
	}
}

func (g *game) LevelStats() {
	free := 0
	exp := 0
	for _, c := range g.Dungeon.Cells {
		if c.IsWall() || c.T == ChasmCell {
			continue
		}
		free++
		if c.Explored {
			exp++
		}
	}
	g.Stats.DExplPerc[g.Depth] = exp * 100 / free
	if g.Stats.DExplPerc[g.Depth] > 93 {
		AchNoviceExplorer.Get(g)
	}
	if g.Depth >= 5 && g.Stats.DExplPerc[g.Depth] > 93 && g.Stats.DExplPerc[g.Depth-1] > 93 && g.Stats.DExplPerc[g.Depth-2] > 93 {
		AchInitiateExplorer.Get(g)
	}
	if g.Depth >= 8 && g.Stats.DExplPerc[g.Depth] > 93 && g.Stats.DExplPerc[g.Depth-1] > 93 && g.Stats.DExplPerc[g.Depth-2] > 93 &&
		g.Stats.DExplPerc[g.Depth-3] > 93 && g.Stats.DExplPerc[g.Depth-4] > 93 {
		AchMasterExplorer.Get(g)
	}
	//g.Stats.DBurns[g.Depth] = g.Stats.CurBurns // XXX to avoid little dump info leak
	nmons := len(g.Monsters)
	kmons := 0
	smons := 0
	for _, mons := range g.Monsters {
		if !mons.Exists() {
			kmons++
			continue
		}
		if mons.State == Resting {
			smons++
		}
	}
	g.Stats.DSleepingPerc[g.Depth] = smons * 100 / nmons
	g.Stats.DKilledPerc[g.Depth] = kmons * 100 / nmons
	g.Stats.DUSpottedPerc[g.Depth] = g.Stats.DUSpotted[g.Depth] * 100 / nmons
}

type achievement string

// Achievements.
const (
	NoAchievement        achievement = "Pitiful Death"
	AchBananaCollector   achievement = "Banana Collector"
	AchHarmonist         achievement = "Harmonist"
	AchOricCelmist       achievement = "Oric Celmist"
	AchUnstealthy        achievement = "Unstealthy Gawalt"
	AchStealthNovice     achievement = "Stealth Novice"
	AchStealthInitiate   achievement = "Stealth Initiate"
	AchStealthMaster     achievement = "Stealth Master"
	AchPyromancer        achievement = "Pyromancer"
	AchDestructor        achievement = "Destructor"
	AchTeleport          achievement = "Oric Teleport Maniac"
	AchCloak             achievement = "Dressed Gawalt"
	AchAmulet            achievement = "Protective Charm"
	AchRescuedShaedra    achievement = "Rescuer"
	AchRetrievedArtifact achievement = "Artifact Retriever"
	AchAcrobat           achievement = "Acrobat"
	AchTree              achievement = "Tree Climber"
	AchTable             achievement = "Under Table Gawalt"
	AchHole              achievement = "Hole Crawler"
	AchDoors             achievement = "Door Opener"
	AchBarrels           achievement = "Barrel Enthousiast"
	AchExtinguisher      achievement = "Light Extinguisher"
	AchLoreStudent       achievement = "Lore student"
	AchLoremaster        achievement = "Loremaster"
	AchNoviceExplorer    achievement = "Novice Explorer"
	AchInitiateExplorer  achievement = "Initiate Explorer"
	AchMasterExplorer    achievement = "Master Explorer"
	AchAssassin          achievement = "Assassin"
	AchInsomniaNovice    achievement = "Insomnia Novice"
	AchInsomniaInitiate  achievement = "Insomnia Initiate"
	AchInsomniaMaster    achievement = "Insomnia Master"
	AchAntimagicNovice   achievement = "Antimagic Novice"
	AchAntimagicInitiate achievement = "Antimagic Initiate"
	AchAntimagicMaster   achievement = "Antimagic Master"
	AchEscape            achievement = "Escape"
)

func (ach achievement) Get(g *game) {
	if !g.Stats.Achievements[ach] {
		g.Stats.Achievements[ach] = true
		g.PrintfStyled("Achievement: %s.", logSpecial, ach)
		g.StoryPrintf("Achievement: %s.", ach)
	}
}

func (t terrain) ReachNotable() bool {
	switch t {
	case TreeCell, TableCell, HoledWallCell, DoorCell, BarrelCell:
		return true
	default:
		return false
	}
}

func (pos position) Reach(g *game) {
	if g.Stats.AtNotablePos[pos] {
		return
	}
	g.Stats.AtNotablePos[pos] = true
	c := g.Dungeon.Cell(pos)
	switch c.T {
	case TreeCell:
		g.Stats.ClimbedTree++
		if g.Stats.ClimbedTree == 12 {
			AchTree.Get(g)
		}
	case TableCell:
		g.Stats.TableHides++
		if g.Stats.TableHides == 12 {
			AchTable.Get(g)
		}
	case HoledWallCell:
		g.Stats.HoledWallsCrawled++
		if g.Stats.HoledWallsCrawled == 12 {
			AchHole.Get(g)
		}
	case DoorCell:
		g.Stats.DoorsOpened++
		if g.Stats.DoorsOpened == 100 {
			AchDoors.Get(g)
		}
	case BarrelCell:
		g.Stats.BarrelHides++
		if g.Stats.BarrelHides == 20 {
			AchBarrels.Get(g)
		}
	}
}
