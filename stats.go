package main

type stats struct {
	Story          []string
	Killed         int
	KilledMons     map[monsterKind]int
	Moves          int
	Jumps          int
	ReceivedHits   int
	Dodges         int
	MagarasUsed    int
	UsedStones     int
	UsedMagaras    map[magara]int
	Damage         int
	DExplPerc      [MaxDepth + 1]int
	DSleepingPerc  [MaxDepth + 1]int
	DKilledPerc    [MaxDepth + 1]int
	Burns          int
	Digs           int
	Rest           int
	Turns          int
	TWounded       int
	TMWounded      int
	TMonsLOS       int
	NSpotted       int
	NUSpotted      int
	DSpotted       [MaxDepth + 1]int
	DUSpotted      [MaxDepth + 1]int
	Achievements   map[achievement]bool
	HarmonicMagUse int
	OricMagUse     int
	FireUse        int
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
}

type achievement string

// Achievements.
const (
	NoAchievement        achievement = "Stupid Death"
	AchBananaCollector   achievement = "Banana Collector"
	AchHarmonist         achievement = "Harmonist"
	AchOricCelmist       achievement = "Oric Celmist"
	AchUnstealthy        achievement = "Unstealthy Gawalt"
	AchStealthNovice     achievement = "Stealth Novice"
	AchStealthInitiate   achievement = "Stealth Initiate"
	AchStealthMaster     achievement = "Stealth Master"
	AchPyromancer        achievement = "Pyromancer"
	AchDestructor        achievement = "Destructor"
	AchTeleport          achievement = "Teleport Maniac"
	AchCloak             achievement = "Dressed Gawalt"
	AchAmulet            achievement = "Protective Charm"
	AchRescuedShaedra    achievement = "Rescuer"
	AchRetrievedArtifact achievement = "Artifact Finding"
	AchAcrobat           achievement = "Acrobat"
	AchTree              achievement = "Tree Climber"
	AchTable             achievement = "Table Hiding"
	AchHole              achievement = "Hole Crawler"
	AchExtinguisher      achievement = "Light Extinguisher"
	AchLoremaster        achievement = "Loremaster"
	AchExplorer          achievement = "Explorer"
	AchKiller            achievement = "Killer"
	AchInsomnia          achievement = "Insomnia"
	AchAntimagic         achievement = "Antimagic"
	AchWinInsomnia       achievement = "Insomnia Win"
	AchWinNoDamage       achievement = "Unhurt Win"
)

func (ach achievement) Get(g *game) {
	if !g.Stats.Achievements[ach] {
		g.Stats.Achievements[ach] = true
		g.PrintfStyled("Achievement: %s.", logSpecial, ach)
		g.StoryPrintf("Achievement: %s.", ach)
	}
}
