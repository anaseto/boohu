package main

import "fmt"

type monsterState int

const (
	Resting monsterState = iota
	Hunting
	Wandering
)

func (m monsterState) String() string {
	var st string
	switch m {
	case Resting:
		st = "resting"
	case Wandering:
		st = "wandering"
	case Hunting:
		st = "hunting"
	}
	return st
}

type monsterStatus int

const (
	MonsConfused monsterStatus = iota
	MonsExhausted
	// unimplemented
	MonsAfraid
)

const NMonsStatus = 2

func (st monsterStatus) String() (text string) {
	switch st {
	case MonsConfused:
		text = "confused"
	case MonsExhausted:
		text = "exhausted"
	case MonsAfraid:
		text = "afraid"
	}
	return text
}

type monsterKind int

const (
	MonsGoblin monsterKind = iota
	MonsOgre
	MonsCyclop
	MonsWorm
	MonsBrizzia
	MonsHound
	MonsYack
	MonsGiantBee
	MonsGoblinWarrior
	MonsHydra
	MonsSkeletonWarrior
	MonsSpider
	MonsBlinkingFrog
	MonsLich
	MonsEarthDragon
	MonsMirrorSpecter
	MonsAcidMound
	MonsExplosiveNadre
	MonsSatowalgaPlant
	MonsMarevorHelith
)

func (mk monsterKind) String() string {
	return MonsData[mk].name
}

func (mk monsterKind) MovementDelay() int {
	return MonsData[mk].movementDelay
}

func (mk monsterKind) Letter() rune {
	return MonsData[mk].letter
}

func (mk monsterKind) AttackDelay() int {
	return MonsData[mk].attackDelay
}

func (mk monsterKind) BaseAttack() int {
	return MonsData[mk].baseAttack
}

func (mk monsterKind) MaxHP() int {
	return MonsData[mk].maxHP
}

func (mk monsterKind) Dangerousness() int {
	return MonsData[mk].dangerousness
}

func (mk monsterKind) Ranged() bool {
	switch mk {
	case MonsLich, MonsCyclop, MonsGoblinWarrior, MonsSatowalgaPlant:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Smiting() bool {
	switch mk {
	case MonsMirrorSpecter:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Desc() string {
	return monsDesc[mk]
}

func (mk monsterKind) SeenStoryText() (text string) {
	switch mk {
	case MonsMarevorHelith:
		text = "You saw Marevor."
	default:
		text = fmt.Sprintf("You saw %s.", Indefinite(mk.String(), false))
	}
	return text
}

func (mk monsterKind) Indefinite(capital bool) (text string) {
	switch mk {
	case MonsMarevorHelith:
		text = mk.String()
	default:
		text = Indefinite(mk.String(), capital)
	}
	return text
}

func (mk monsterKind) Definite(capital bool) (text string) {
	switch mk {
	case MonsMarevorHelith:
		text = mk.String()
	default:
		if capital {
			text = fmt.Sprintf("The %s", mk.String())
		} else {
			text = fmt.Sprintf("the %s", mk.String())
		}
	}
	return text
}

type monsterData struct {
	movementDelay int
	baseAttack    int
	attackDelay   int
	maxHP         int
	accuracy      int
	armor         int
	evasion       int
	letter        rune
	name          string
	dangerousness int
}

var MonsData = []monsterData{
	MonsGoblin:          {10, 7, 10, 15, 14, 0, 12, 'g', "goblin", 2},
	MonsOgre:            {10, 15, 12, 28, 13, 0, 8, 'O', "ogre", 6},
	MonsCyclop:          {10, 12, 12, 28, 13, 0, 8, 'C', "cyclop", 9},
	MonsWorm:            {12, 9, 10, 25, 13, 0, 10, 'w', "farmer worm", 3},
	MonsBrizzia:         {12, 10, 10, 30, 13, 0, 10, 'z', "brizzia", 7},
	MonsAcidMound:       {10, 9, 10, 19, 15, 0, 8, 'a', "acid mound", 7},
	MonsHound:           {8, 9, 10, 15, 14, 0, 12, 'h', "hound", 4},
	MonsYack:            {10, 11, 10, 21, 14, 0, 10, 'y', "yack", 6},
	MonsGiantBee:        {6, 10, 10, 11, 15, 0, 15, 'B', "giant bee", 6},
	MonsGoblinWarrior:   {10, 11, 10, 22, 15, 3, 12, 'G', "goblin warrior", 8},
	MonsHydra:           {10, 9, 10, 45, 13, 0, 6, 'H', "hydra", 15},
	MonsSkeletonWarrior: {10, 12, 10, 25, 15, 4, 12, 'S', "skeleton warrior", 10},
	MonsSpider:          {8, 7, 10, 13, 17, 0, 15, 's', "spider", 6},
	MonsBlinkingFrog:    {10, 10, 10, 20, 15, 0, 12, 'F', "blinking frog", 7},
	MonsLich:            {10, 10, 10, 23, 15, 3, 12, 'L', "lich", 17},
	MonsEarthDragon:     {10, 14, 10, 40, 14, 6, 8, 'D', "earth dragon", 20},
	MonsMirrorSpecter:   {10, 9, 10, 18, 15, 0, 17, 'm', "mirror specter", 11},
	MonsExplosiveNadre:  {10, 4, 10, 1, 14, 0, 10, 'n', "explosive nadre", 5},
	MonsSatowalgaPlant:  {10, 12, 12, 30, 15, 0, 4, 'P', "satowalga plant", 7},
	MonsMarevorHelith:   {10, 0, 10, 99, 18, 10, 15, 'M', "Marevor Helith", 18},
}

var monsDesc = []string{
	MonsGoblin:          "Goblins are little humanoid creatures. They often appear in group.",
	MonsOgre:            "Ogres are big clunky humanoids that can hit really hard.",
	MonsCyclop:          "Cyclops are very similar to ogres, but they also like to throw rocks at their foes (for up to 15 damage). The rocks can block your way for a while.",
	MonsWorm:            "Farmer worms are ugly slow moving creatures, but surprisingly hardy at times, and they furrow as they move, helping new foliage to grow.",
	MonsBrizzia:         "Brizzias are big slow moving biped creatures. They are quite hardy, and when hurt they can cause nausea, impeding the use of potions.",
	MonsAcidMound:       "Acid mounds are acidic creatures. They can temporally corrode your equipment.",
	MonsHound:           "Hounds are fast moving carnivore quadrupeds. They sometimes attack in group.",
	MonsYack:            "Yacks are quite large herbivorous quadrupeds. They tend to form large groups. They can push you one cell away.",
	MonsGiantBee:        "Giant bees are fragile, but extremely fast moving creatures. Their bite can sometimes enrage you.",
	MonsGoblinWarrior:   "Goblin warriors are goblins that learned to fight, and got equipped with a leather armour. They can throw javelins.",
	MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	MonsSkeletonWarrior: "Skeleton warriors are good fighters, and are equipped with a chain mail.",
	MonsSpider:          "Spiders are fast moving fragile creatures, whose bite can confuse you.",
	MonsBlinkingFrog:    "Blinking frogs are big frog-like unstable creatures, whose bite can make you blink away.",
	MonsLich:            "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you, halving your HP.",
	MonsEarthDragon:     "Earth dragons are big and hardy creatures that wander in the Underground. It is said they are to credit for many tunnels.",
	MonsMirrorSpecter:   "Mirror specters are very insubstantial creatures. They can absorb your mana.",
	MonsExplosiveNadre:  "Explosive nadres are very frail creatures that explode upon dying, halving HP of any adjacent creatures and occasionally destroying walls.",
	MonsSatowalgaPlant:  "Satowalga Plants are static bushes that throw acidic projectiles at you, sometimes corroding and confusing you.",
	MonsMarevorHelith:   "Marevor Helith is an ancient undead nakrus very fond of teleporting people away.",
}

type monsterBand int

const (
	LoneGoblin monsterBand = iota
	LoneOgre
	LoneWorm
	LoneRareWorm
	LoneBrizzia
	LoneHound
	LoneHydra
	LoneSpider
	LoneBlinkingFrog
	LoneCyclop
	LoneLich
	LoneEarthDragon
	LoneSpecter
	LoneAcidMound
	LoneExplosiveNadre
	LoneSatowalgaPlant
	BandGoblins
	BandGoblinsWithWarriors
	BandGoblinWarriors
	BandHounds
	BandYacks
	BandSpiders
	BandSatowalga
	BandBlinkingFrogs
	BandExplosive
	BandGiantBees
	BandSkeletonWarrior
	UBandWorms
	UBandGoblinsEasy
	UBandFrogs
	UBandOgres
	UBandGoblins
	UBandBeeYacks
	UHydras
	UExplosiveNadres
	ULich
	UBrizzias
	UAcidMounds
	USatowalga
	UDragon
	UMarevorHelith
	UXCyclops
	UXLiches
	UXFrogRanged
	UXExplosive
	UXWarriors
	UXSatowalga
	UXSpecters
	UXDisabling
)

type monsInterval struct {
	min int
	max int
}

type monsterBandData struct {
	distribution map[monsterKind]monsInterval
	rarity       int
	minDepth     int
	maxDepth     int
	band         bool
	monster      monsterKind
	unique       bool
}

func (g *game) GenBand(mbd monsterBandData, band monsterBand) []monsterKind {
	if g.GeneratedBands[band] > 0 && mbd.unique {
		return nil
	}
	if g.Depth > mbd.maxDepth+RandInt(3) || RandInt(10) == 0 {
		return nil
	}
	if g.Depth < mbd.minDepth-RandInt(3) {
		return nil
	}
	if !mbd.band {
		return []monsterKind{mbd.monster}
	}
	bandMonsters := []monsterKind{}
	for m, interval := range mbd.distribution {
		for i := 0; i < interval.min+RandInt(interval.max-interval.min+1); i++ {
			bandMonsters = append(bandMonsters, m)
		}
	}
	return bandMonsters
}

var MonsBands = []monsterBandData{
	LoneGoblin:         {rarity: 10, minDepth: 0, maxDepth: 5, monster: MonsGoblin},
	LoneOgre:           {rarity: 15, minDepth: 2, maxDepth: 11, monster: MonsOgre},
	LoneWorm:           {rarity: 10, minDepth: 0, maxDepth: 6, monster: MonsWorm},
	LoneRareWorm:       {rarity: 90, minDepth: 7, maxDepth: 13, monster: MonsWorm},
	LoneBrizzia:        {rarity: 90, minDepth: 7, maxDepth: 13, monster: MonsBrizzia},
	LoneHound:          {rarity: 20, minDepth: 1, maxDepth: 8, monster: MonsHound},
	LoneHydra:          {rarity: 45, minDepth: 8, maxDepth: 13, monster: MonsHydra},
	LoneSpider:         {rarity: 20, minDepth: 3, maxDepth: 13, monster: MonsSpider},
	LoneBlinkingFrog:   {rarity: 50, minDepth: 5, maxDepth: 13, monster: MonsBlinkingFrog},
	LoneCyclop:         {rarity: 35, minDepth: 5, maxDepth: 13, monster: MonsCyclop},
	LoneLich:           {rarity: 70, minDepth: 9, maxDepth: 13, monster: MonsLich},
	LoneEarthDragon:    {rarity: 80, minDepth: 10, maxDepth: 13, monster: MonsEarthDragon},
	LoneSpecter:        {rarity: 70, minDepth: 6, maxDepth: 13, monster: MonsMirrorSpecter},
	LoneAcidMound:      {rarity: 70, minDepth: 6, maxDepth: 13, monster: MonsAcidMound},
	LoneExplosiveNadre: {rarity: 55, minDepth: 4, maxDepth: 7, monster: MonsExplosiveNadre},
	LoneSatowalgaPlant: {rarity: 80, minDepth: 5, maxDepth: 13, monster: MonsSatowalgaPlant},
	BandGoblins: {
		distribution: map[monsterKind]monsInterval{MonsGoblin: {2, 4}},
		rarity:       10, minDepth: 1, maxDepth: 5, band: true,
	},
	BandGoblinsWithWarriors: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 5},
			MonsGoblinWarrior: {0, 2}},
		rarity: 12, minDepth: 5, maxDepth: 9, band: true,
	},
	BandGoblinWarriors: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {0, 1},
			MonsGoblinWarrior: {2, 4}},
		rarity: 45, minDepth: 10, maxDepth: 13, band: true,
	},
	BandHounds: {
		distribution: map[monsterKind]monsInterval{MonsHound: {2, 3}},
		rarity:       20, minDepth: 2, maxDepth: 10, band: true,
	},
	BandSpiders: {
		distribution: map[monsterKind]monsInterval{MonsSpider: {2, 4}},
		rarity:       25, minDepth: 5, maxDepth: 13, band: true,
	},
	BandBlinkingFrogs: {
		distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: {2, 4}},
		rarity:       65, minDepth: 9, maxDepth: 13, band: true,
	},
	BandSatowalga: {
		distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {2, 2},
		},
		rarity: 100, minDepth: 7, maxDepth: 13, band: true,
	},
	BandExplosive: {
		distribution: map[monsterKind]monsInterval{
			MonsBlinkingFrog:   {0, 1},
			MonsExplosiveNadre: {1, 2},
			MonsGiantBee:       {1, 1},
			MonsBrizzia:        {0, 1},
		},
		rarity: 60, minDepth: 8, maxDepth: 13, band: true,
	},
	BandYacks: {
		distribution: map[monsterKind]monsInterval{MonsYack: {2, 5}},
		rarity:       15, minDepth: 5, maxDepth: 11, band: true,
	},
	BandGiantBees: {
		distribution: map[monsterKind]monsInterval{MonsGiantBee: {2, 5}},
		rarity:       30, minDepth: 6, maxDepth: 13, band: true,
	},
	BandSkeletonWarrior: {
		distribution: map[monsterKind]monsInterval{MonsSkeletonWarrior: {2, 3}},
		rarity:       60, minDepth: 8, maxDepth: 13, band: true,
	},
	UBandWorms: {
		distribution: map[monsterKind]monsInterval{MonsWorm: {3, 4}, MonsSpider: {1, 1}},
		rarity:       50, minDepth: 4, maxDepth: 4, band: true, unique: true,
	},
	UBandGoblinsEasy: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin: {3, 5},
			MonsHound:  {1, 2},
		},
		rarity: 30, minDepth: 5, maxDepth: 5, band: true, unique: true,
	},
	UBandFrogs: {
		distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: {2, 3}},
		rarity:       60, minDepth: 6, maxDepth: 6, band: true, unique: true,
	},
	UBandOgres: {
		distribution: map[monsterKind]monsInterval{MonsOgre: {2, 3}, MonsCyclop: {1, 1}},
		rarity:       35, minDepth: 7, maxDepth: 7, band: true, unique: true,
	},
	UBandGoblins: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {2, 3},
			MonsGoblinWarrior: {2, 2},
			MonsHound:         {1, 2},
		},
		rarity: 30, minDepth: 8, maxDepth: 8, band: true, unique: true,
	},
	UBandBeeYacks: {
		distribution: map[monsterKind]monsInterval{
			MonsYack:     {3, 4},
			MonsGiantBee: {2, 2},
		},
		rarity: 30, minDepth: 9, maxDepth: 9, band: true, unique: true,
	},
	UHydras: {
		distribution: map[monsterKind]monsInterval{
			MonsHydra:  {2, 3},
			MonsSpider: {1, 2},
		},
		rarity: 40, minDepth: 10, maxDepth: 10, band: true, unique: true,
	},
	UExplosiveNadres: {
		distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {2, 3},
			MonsBrizzia:        {1, 2},
		},
		rarity: 55, minDepth: 10, maxDepth: 10, band: true, unique: true,
	},
	ULich: {
		distribution: map[monsterKind]monsInterval{
			MonsSkeletonWarrior: {2, 2},
			MonsLich:            {1, 1},
			MonsMirrorSpecter:   {0, 1},
		},
		rarity: 50, minDepth: 11, maxDepth: 11, band: true, unique: true,
	},
	UBrizzias: {
		distribution: map[monsterKind]monsInterval{
			MonsBrizzia: {3, 4},
		},
		rarity: 80, minDepth: 11, maxDepth: 11, band: true, unique: true,
	},
	UAcidMounds: {
		distribution: map[monsterKind]monsInterval{
			MonsAcidMound: {3, 4},
		},
		rarity: 80, minDepth: 12, maxDepth: 12, band: true, unique: true,
	},
	USatowalga: {
		distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {3, 3},
		},
		rarity: 80, minDepth: 12, maxDepth: 12, band: true, unique: true,
	},
	UDragon: {
		distribution: map[monsterKind]monsInterval{
			MonsEarthDragon: {2, 2},
		},
		rarity: 60, minDepth: 12, maxDepth: 12, band: true, unique: true,
	},
	UMarevorHelith: {
		distribution: map[monsterKind]monsInterval{
			MonsMarevorHelith: {1, 1},
			MonsLich:          {0, 1},
		},
		rarity: 100, minDepth: 7, maxDepth: 15, band: true, unique: true,
	},
	UXCyclops: {
		distribution: map[monsterKind]monsInterval{
			MonsCyclop: {3, 3},
		},
		rarity: 100, minDepth: 13, maxDepth: 15, band: true, unique: true,
	},
	UXLiches: {
		distribution: map[monsterKind]monsInterval{
			MonsLich: {2, 2},
		},
		rarity: 100, minDepth: 14, maxDepth: 15, band: true, unique: true,
	},
	UXFrogRanged: {
		distribution: map[monsterKind]monsInterval{
			MonsBlinkingFrog: {2, 2},
			MonsCyclop:       {1, 1},
			MonsLich:         {1, 1},
		},
		rarity: 100, minDepth: 14, maxDepth: 15, band: true, unique: true,
	},
	UXExplosive: {
		distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {5, 5},
		},
		rarity: 100, minDepth: 13, maxDepth: 15, band: true, unique: true,
	},
	UXWarriors: {
		distribution: map[monsterKind]monsInterval{
			MonsHound:         {2, 2},
			MonsGoblinWarrior: {3, 3},
		},
		rarity: 100, minDepth: 14, maxDepth: 15, band: true, unique: true,
	},
	UXSatowalga: {
		distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {3, 3},
		},
		rarity: 100, minDepth: 13, maxDepth: 15, band: true, unique: true,
	},
	UXSpecters: {
		distribution: map[monsterKind]monsInterval{
			MonsMirrorSpecter: {3, 3},
		},
		rarity: 100, minDepth: 14, maxDepth: 15, band: true, unique: true,
	},
	UXDisabling: {
		distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {1, 1},
			MonsSpider:         {1, 1},
			MonsBrizzia:        {1, 1},
			MonsGiantBee:       {1, 1},
			MonsMirrorSpecter:  {1, 1},
		},
		rarity: 100, minDepth: 15, maxDepth: 15, band: true, unique: true,
	},
}

type monster struct {
	Kind        monsterKind
	Band        int
	Index       int
	Attack      int
	Accuracy    int
	Armor       int
	Evasion     int
	HPmax       int
	HP          int
	State       monsterState
	Statuses    [NMonsStatus]int
	Pos         position
	Target      position
	Path        []position // cache
	Obstructing bool
	FireReady   bool
	Seen        bool
}

func (m *monster) Init() {
	m.HPmax = MonsData[m.Kind].maxHP - 1 + RandInt(3)
	m.Attack = MonsData[m.Kind].baseAttack
	m.HP = m.HPmax
	m.Accuracy = MonsData[m.Kind].accuracy
	m.Armor = MonsData[m.Kind].armor
	m.Evasion = MonsData[m.Kind].evasion
	if m.Kind == MonsMarevorHelith {
		m.State = Wandering
	}
}

func (m *monster) Status(st monsterStatus) bool {
	return m.Statuses[st] > 0
}

func (m *monster) Exists() bool {
	return m != nil && m.HP > 0
}

func (m *monster) AlternatePlacement(g *game) *position {
	var neighbors []position
	if m.Status(MonsConfused) {
		neighbors = g.Dungeon.CardinalFreeNeighbors(m.Pos)
	} else {
		neighbors = g.Dungeon.FreeNeighbors(m.Pos)
	}
	for _, pos := range neighbors {
		if pos.Distance(g.Player.Pos) != 1 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return &pos
	}
	return nil
}

func (m *monster) TeleportPlayer(g *game, ev event) {
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	if acc > evasion {
		g.Print("Marevor pushes you through a monolith.")
		g.StoryPrint("Marevor pushed you through a monolith.")
		g.Teleportation(ev)
	} else if RandInt(2) == 0 {
		g.Print("Marevor inadvertently goes into a monolith.")
		m.TeleportAway(g)
	}
}

func (m *monster) TeleportAway(g *game) {
	pos := m.Pos
	i := 0
	count := 0
	for {
		count++
		if count > 1000 {
			panic("TeleportOther")
		}
		pos = g.FreeCell()
		if pos.Distance(m.Pos) < 15 && i < 1000 {
			i++
			continue
		}
		break
	}

	switch m.State {
	case Hunting:
		m.State = Wandering
		// TODO: change the target?
	case Resting, Wandering:
		m.State = Wandering
		m.Target = m.Pos
	}
	if g.Player.LOS[m.Pos] {
		g.Printf("%s teleports away.", m.Kind.Definite(true))
	}
	opos := m.Pos
	m.MoveTo(g, pos)
	if g.Player.LOS[opos] {
		g.ui.TeleportAnimation(g, opos, pos, false)
	}
}

func (m *monster) MoveTo(g *game, pos position) {
	if !g.Player.LOS[m.Pos] && g.Player.LOS[pos] {
		if !m.Seen {
			m.Seen = true
			g.Printf("%s (%v) comes into view.", m.Kind.Indefinite(true), m.State)
		}
		g.StopAuto()
	}
	recomputeLOS := g.Player.LOS[m.Pos] && g.Doors[m.Pos] || g.Player.LOS[pos] && g.Doors[pos]
	m.PlaceAt(g, pos)
	if recomputeLOS {
		g.ComputeLOS()
	}
}

func (m *monster) PlaceAt(g *game, pos position) {
	g.MonstersPosCache[m.Pos.idx()] = 0
	m.Pos = pos
	g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
}

func (m *monster) TeleportMonsterAway(g *game) bool {
	neighbors := g.Dungeon.FreeNeighbors(m.Pos)
	for _, pos := range neighbors {
		if pos == m.Pos || RandInt(3) != 0 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			if g.Player.LOS[m.Pos] {
				g.Print("Marevor makes some strange gestures.")
			}
			mons.TeleportAway(g)
			return true
		}
	}
	return false
}

func (m *monster) AttackAction(g *game, ev event) {
	switch {
	case m.Obstructing:
		m.Obstructing = false
		pos := m.AlternatePlacement(g)
		if pos != nil {
			m.MoveTo(g, *pos)
			ev.Renew(g, m.Kind.MovementDelay())
			return
		}
		fallthrough
	default:
		if m.Kind == MonsHydra {
			for i := 0; i <= 3; i++ {
				m.HitPlayer(g, ev)
			}
		} else if m.Kind == MonsMarevorHelith {
			m.TeleportPlayer(g, ev)
		} else {
			m.HitPlayer(g, ev)
		}
		ev.Renew(g, m.Kind.AttackDelay())
	}
}

func (m *monster) NaturalAwake(g *game) {
	m.Target = g.FreeCell()
	m.State = Wandering
	m.GatherBand(g)
}

func (m *monster) HandleTurn(g *game, ev event) {
	ppos := g.Player.Pos
	mpos := m.Pos
	m.MakeAware(g)
	if m.State == Resting {
		wander := RandInt(100 + 6*Max(800-(g.DepthPlayerTurn+1), 0))
		if wander == 0 {
			m.NaturalAwake(g)
		}
		ev.Renew(g, m.Kind.MovementDelay())
		return
	}
	if m.State == Hunting && m.RangedAttack(g, ev) {
		return
	}
	if m.State == Hunting && m.SmitingAttack(g, ev) {
		return
	}
	if m.Kind == MonsSatowalgaPlant {
		ev.Renew(g, m.Kind.MovementDelay())
		// oklob plants are static ranged-only
		return
	}
	if mpos.Distance(ppos) == 1 {
		attack := true
		if m.Status(MonsConfused) {
			switch m.Pos.Dir(g.Player.Pos) {
			case E, N, W, S:
			default:
				attack = false
			}
		}
		if attack {
			m.AttackAction(g, ev)
			return
		}
	}
	if m.Kind == MonsMarevorHelith {
		if m.TeleportMonsterAway(g) {
			ev.Renew(g, m.Kind.MovementDelay())
			return
		}
	}
	m.Obstructing = false
	if !(len(m.Path) > 0 && m.Path[0] == m.Target && m.Path[len(m.Path)-1] == mpos) {
		m.Path = m.APath(g, mpos, m.Target)
		if len(m.Path) == 0 {
			// if target is not accessible, try free neighbor cells
			for _, npos := range g.Dungeon.FreeNeighbors(m.Target) {
				m.Path = m.APath(g, mpos, npos)
				if len(m.Path) > 0 {
					m.Target = npos
					break
				}
			}
		}
	}
	if len(m.Path) < 2 {
		switch m.State {
		case Wandering:
			keepWandering := RandInt(100)
			if keepWandering > 75 && MonsBands[g.Bands[m.Band]].band {
				for _, mons := range g.Monsters {
					m.Target = mons.Pos
				}
			} else {
				m.Target = g.FreeCell()
			}
			m.GatherBand(g)
		case Hunting:
			// pick a random cell: more escape strategies for the player
			m.Target = g.FreeCell()
			m.State = Wandering
			m.GatherBand(g)
		}
		ev.Renew(g, m.Kind.MovementDelay())
		return
	}
	target := m.Path[len(m.Path)-2]
	mons := g.MonsterAt(target)
	switch {
	case !mons.Exists():
		if m.Kind == MonsEarthDragon && g.Dungeon.Cell(target).T == WallCell {
			g.Dungeon.SetCell(target, FreeCell)
			g.Stats.Digs++
			if !g.Player.LOS[target] {
				g.WrongWall[m.Pos] = true
			}
			g.MakeNoise(WallNoise, m.Pos)
			g.Fog(m.Pos, 1, ev)
			if g.Player.Pos.Distance(target) < 12 {
				// XXX use dijkstra distance ?
				g.Printf("%s You hear an earth-breaking noise.", g.CrackSound())
				g.StopAuto()
			}
			m.MoveTo(g, target)
			m.Path = m.Path[:len(m.Path)-1]
		} else if g.Dungeon.Cell(target).T == WallCell {
			m.Path = m.APath(g, mpos, m.Target)
		} else {
			m.InvertFoliage(g)
			m.MoveTo(g, target)
			if m.Kind.Ranged() && !m.FireReady && g.Player.LOS[m.Pos] {
				m.FireReady = true
			}
			m.Path = m.Path[:len(m.Path)-1]
		}
	case m.State == Hunting && mons.State != Hunting:
		r := RandInt(5)
		if r == 0 {
			mons.Target = m.Target
			mons.State = Wandering
			mons.GatherBand(g)
		} else if (r == 1 || r == 2) && g.Player.Pos.Distance(mons.Target) > 2 {
			mons.Target = g.FreeCell()
			mons.State = Wandering
			mons.GatherBand(g)
		} else {
			m.Path = m.APath(g, mpos, m.Target)
		}
	case !g.Player.LOS[mons.Pos] && g.Player.Pos.Distance(mons.Target) > 2 && mons.State != Hunting:
		r := RandInt(5)
		if r == 0 {
			m.Target = g.FreeCell()
			m.GatherBand(g)
		} else if (r == 1 || r == 2) && mons.State == Resting {
			mons.Target = g.FreeCell()
			mons.State = Wandering
			mons.GatherBand(g)
		} else {
			m.Path = m.APath(g, mpos, m.Target)
		}
	case mons.Pos.Distance(g.Player.Pos) == 1:
		m.Path = m.APath(g, mpos, m.Target)
		if len(m.Path) < 2 || m.Path[len(m.Path)-2] == mons.Pos {
			mons.Obstructing = true
		}
	case mons.State == Hunting && m.State == Hunting:
		if RandInt(5) == 0 {
			m.Target = mons.Target
			m.Path = m.APath(g, mpos, m.Target)
		} else {
			m.Path = m.APath(g, mpos, m.Target)
		}
	default:
		m.Path = m.APath(g, mpos, m.Target)
	}
	ev.Renew(g, m.Kind.MovementDelay())
}

func (m *monster) InvertFoliage(g *game) {
	if m.Kind != MonsWorm {
		return
	}
	invert := false
	if _, ok := g.Fungus[m.Pos]; !ok {
		if _, ok := g.Doors[m.Pos]; !ok {
			g.Fungus[m.Pos] = foliage
			invert = true
		}
	} else {
		delete(g.Fungus, m.Pos)
		invert = true
	}
	if !g.Player.LOS[m.Pos] && invert {
		g.WrongFoliage[m.Pos] = !g.WrongFoliage[m.Pos]
	}
}

func (m *monster) DramaticAdjustment(g *game, baseAttack, attack, evasion, acc int, clang bool) (int, int, bool) {
	if attack >= g.Player.HP {
		// a little dramatic effect
		if RandInt(2) == 0 {
			attack, clang = g.HitDamage(DmgPhysical, baseAttack, g.Player.Armor())
		}
		if attack >= g.Player.HP {
			n := RandInt(g.Player.Evasion())
			if n > evasion {
				evasion = n
			}
		}
	}
	if baseAttack >= g.Player.HP && (acc <= evasion || attack < g.Player.HP) {
		g.Stats.TimesLucky++
	}
	return attack, evasion, clang
}

func (m *monster) HitPlayer(g *game, ev event) {
	if g.Player.HP <= 0 {
		return
	}
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	attack, clang := g.HitDamage(DmgPhysical, m.Attack, g.Player.Armor())
	attack, evasion, clang = m.DramaticAdjustment(g, m.Attack, attack, evasion, acc, clang)
	if acc > evasion {
		if m.Blocked(g) {
			g.Printf("Clang! You block %s's attack.", m.Kind.Definite(false))
			g.MakeNoise(ShieldBlockNoise, g.Player.Pos)
			return
		}
		if g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification) {
			g.SwapWithMonster(m)
			return
		}
		noise := g.HitNoise(clang)
		g.MakeNoise(noise, g.Player.Pos)
		var sclang string
		if clang {
			sclang = g.ArmourClang()
		}
		g.PrintfStyled("%s hits you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), attack, sclang)
		m.InflictDamage(g, attack, m.Attack)
		if g.Player.HP <= 0 {
			return
		}
		m.HitSideEffects(g, ev)
		if g.Player.Aptitudes[AptConfusingGas] && g.Player.HP < g.Player.HPMax()/2 && RandInt(3) == 0 {
			m.EnterConfusion(g, ev)
			g.Printf("You release a confusing gas on the %s.", m.Kind)
		}
		if g.Player.Aptitudes[AptSmoke] && g.Player.HP < g.Player.HPMax()/2 && RandInt(2) == 0 {
			g.Smoke(ev)
		}
	} else {
		g.Printf("%s misses you.", m.Kind.Definite(true))
	}
}

func (m *monster) EnterConfusion(g *game, ev event) {
	if !m.Status(MonsConfused) {
		m.Statuses[MonsConfused] = 1
		m.Path = m.Path[:0]
		g.PushEvent(&monsterEvent{
			ERank: ev.Rank() + 50 + RandInt(100), NMons: m.Index, EAction: MonsConfusionEnd})
	}
}

func (m *monster) HitSideEffects(g *game, ev event) {
	switch m.Kind {
	case MonsSpider:
		if RandInt(2) == 0 {
			g.Confusion(ev)
		}
	case MonsGiantBee:
		if RandInt(5) == 0 && !g.Player.HasStatus(StatusBerserk) && !g.Player.HasStatus(StatusExhausted) {
			g.Player.Statuses[StatusBerserk] = 1
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 25 + RandInt(40), EAction: BerserkEnd})
			g.Print("You feel a sudden urge to kill things.")
		}
	case MonsBlinkingFrog:
		if RandInt(2) == 0 {
			g.Blink(ev)
		}
	case MonsAcidMound:
		g.Corrosion(ev)
	case MonsYack:
		if RandInt(2) == 0 && m.PushPlayer(g) {
			g.Print("The yack pushes you.")
		}
	}
}

func (m *monster) PushPlayer(g *game) (pushed bool) {
	dir := g.Player.Pos.Dir(m.Pos)
	pos := g.Player.Pos.To(dir)
	if !g.Player.HasStatus(StatusLignification) &&
		pos.valid() && g.Dungeon.Cell(pos).T == FreeCell {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			g.PlacePlayerAt(pos)
			pushed = true
		}
	}
	return pushed
}

func (m *monster) RangedAttack(g *game, ev event) bool {
	if !m.Kind.Ranged() {
		return false
	}
	if m.Pos.Distance(g.Player.Pos) <= 1 && m.Kind != MonsSatowalgaPlant {
		return false
	}
	if !g.Player.LOS[m.Pos] {
		m.FireReady = false
		return false
	}
	if !m.FireReady {
		m.FireReady = true
		if m.Pos.Distance(g.Player.Pos) <= 3 {
			ev.Renew(g, m.Kind.AttackDelay())
			return true
		} else {
			return false
		}
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	case MonsLich:
		return m.TormentBolt(g, ev)
	case MonsCyclop:
		return m.ThrowRock(g, ev)
	case MonsGoblinWarrior:
		return m.ThrowJavelin(g, ev)
	case MonsSatowalgaPlant:
		return m.ThrowAcid(g, ev)
	}
	return false
}

func (m *monster) RangeBlocked(g *game) bool {
	ray := g.Ray(m.Pos)
	blocked := false
	for _, pos := range ray[1:] {
		mons := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		blocked = true
		break
	}
	return blocked
}

func (m *monster) TormentBolt(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	hit := !m.Blocked(g)
	g.MakeNoise(9, m.Pos)
	if hit {
		g.MakeNoise(MagicHitNoise, g.Player.Pos)
		damage := g.Player.HP - g.Player.HP/2
		g.PrintfStyled("%s throws a bolt of torment at you.", logMonsterHit, m.Kind.Definite(true))
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '*', ColorCyan)
		m.InflictDamage(g, damage, 15)
	} else {
		g.Printf("You block the %s's bolt of torment.", m.Kind)
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '*', ColorCyan)
	}
	m.Statuses[MonsExhausted] = 1
	g.PushEvent(&monsterEvent{ERank: ev.Rank() + 100 + RandInt(50), NMons: m.Index, EAction: MonsExhaustionEnd})
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) Blocked(g *game) bool {
	blocked := false
	if g.Player.Shield != NoShield && !g.Player.Weapon.TwoHanded() {
		block := RandInt(g.Player.Block())
		acc := RandInt(m.Accuracy)
		if block >= acc {
			g.MakeNoise(12+g.Player.Block()/2, g.Player.Pos)
			blocked = true
		}
	}
	return blocked
}

func (m *monster) ThrowRock(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	block := false
	hit := true
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	const rockdmg = 15
	attack, clang := g.HitDamage(DmgPhysical, rockdmg, g.Player.Armor())
	attack, evasion, clang = m.DramaticAdjustment(g, rockdmg, attack, evasion, acc, clang)
	if 4*acc/3 <= evasion {
		// rocks are big and do not miss so often
		hit = false
	} else {
		block = m.Blocked(g)
		hit = !block
	}
	if hit {
		noise := g.HitNoise(clang)
		g.MakeNoise(noise, g.Player.Pos)
		var sclang string
		if clang {
			sclang = g.ArmourClang()
		}
		g.PrintfStyled("%s throws a rock at you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), attack, sclang)
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '●', ColorMagenta)
		oppos := g.Player.Pos
		if m.PushPlayer(g) {
			g.TemporalWallAt(oppos, ev)
		} else {
			ray := g.Ray(m.Pos)
			if len(ray) > 0 {
				g.TemporalWallAt(ray[len(ray)-1], ev)
			}
		}
		m.InflictDamage(g, attack, rockdmg)
	} else if block {
		g.Printf("You block %s's rock. Clang!", m.Kind.Indefinite(false))
		g.MakeNoise(ShieldBlockNoise, g.Player.Pos)
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '●', ColorMagenta)
		ray := g.Ray(m.Pos)
		if len(ray) > 0 {
			g.TemporalWallAt(ray[len(ray)-1], ev)
		}
	} else {
		g.Printf("You dodge %s's rock.", m.Kind.Indefinite(false))
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '●', ColorMagenta)
		dir := g.Player.Pos.Dir(m.Pos)
		pos := g.Player.Pos.To(dir)
		if pos.valid() {
			g.TemporalWallAt(pos, ev)
		}
	}
	ev.Renew(g, 2*m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowJavelin(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	block := false
	hit := true
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	const jdmg = 11
	attack, clang := g.HitDamage(DmgPhysical, jdmg, g.Player.Armor())
	attack, evasion, clang = m.DramaticAdjustment(g, jdmg, attack, evasion, acc, clang)
	if acc <= evasion {
		hit = false
	} else {
		block = m.Blocked(g)
		hit = !block
	}
	if hit {
		noise := g.HitNoise(clang)
		g.MakeNoise(noise, g.Player.Pos)
		var sclang string
		if clang {
			sclang = g.ArmourClang()
		}
		g.Printf("%s throws %s at you (%d dmg).%s", m.Kind.Definite(true), Indefinite("javelin", false), attack, sclang)
		g.ui.MonsterJavelinAnimation(g, g.Ray(m.Pos), true)
		m.InflictDamage(g, attack, jdmg)
	} else if block {
		if RandInt(3) == 0 {
			g.Printf("You block %s's %s. Clang!", m.Kind.Indefinite(false), "javelin")
			g.MakeNoise(ShieldBlockNoise, g.Player.Pos)
			g.ui.MonsterJavelinAnimation(g, g.Ray(m.Pos), false)
		} else if !g.Player.HasStatus(StatusDisabledShield) {
			g.Player.Statuses[StatusDisabledShield] = 1
			g.PushEvent(&simpleEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: DisabledShieldEnd})
			g.Printf("%s's %s gets fixed on your shield.", m.Kind.Indefinite(true), "javelin")
			g.MakeNoise(ShieldBlockNoise, g.Player.Pos)
			g.ui.MonsterJavelinAnimation(g, g.Ray(m.Pos), false)
		}
	} else {
		g.Printf("You dodge %s's %s.", m.Kind.Indefinite(false), "javelin")
		g.ui.MonsterJavelinAnimation(g, g.Ray(m.Pos), false)
	}
	m.Statuses[MonsExhausted] = 1
	g.PushEvent(&monsterEvent{ERank: ev.Rank() + 50 + RandInt(50), NMons: m.Index, EAction: MonsExhaustionEnd})
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowAcid(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	block := false
	hit := true
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	acdmg := 12
	attack, clang := g.HitDamage(DmgPhysical, acdmg, g.Player.Armor())
	attack, evasion, clang = m.DramaticAdjustment(g, acdmg, attack, evasion, acc, clang)
	if acc <= evasion {
		hit = false
	} else {
		block = m.Blocked(g)
		hit = !block
	}
	if hit {
		noise := g.HitNoise(false) // no clang with acid projectiles
		g.MakeNoise(noise, g.Player.Pos)
		g.Printf("%s throws acid at you (%d dmg).", m.Kind.Definite(true), attack)
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '*', ColorGreen)
		m.InflictDamage(g, attack, acdmg)
		if RandInt(2) == 0 {
			g.Corrosion(ev)
			if RandInt(2) == 0 {
				g.Confusion(ev)
			}
		}
	} else if block {
		g.Printf("You block %s's acid projectile.", m.Kind.Indefinite(false))
		g.MakeNoise(BaseHitNoise, g.Player.Pos) // no real clang
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '*', ColorGreen)
		if RandInt(2) == 0 {
			g.Corrosion(ev)
		}
	} else {
		g.Printf("You dodge %s's acid projectile.", m.Kind.Indefinite(false))
		g.ui.MonsterProjectileAnimation(g, g.Ray(m.Pos), '*', ColorGreen)
	}
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) SmitingAttack(g *game, ev event) bool {
	if !m.Kind.Smiting() {
		return false
	}
	if !g.Player.LOS[m.Pos] {
		m.FireReady = false
		return false
	}
	if !m.FireReady {
		m.FireReady = true
		if m.Pos.Distance(g.Player.Pos) <= 3 {
			ev.Renew(g, m.Kind.AttackDelay())
			return true
		} else {
			return false
		}
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	case MonsMirrorSpecter:
		return m.AbsorbMana(g, ev)
	}
	return false
}

func (m *monster) AbsorbMana(g *game, ev event) bool {
	if g.Player.MP == 0 {
		return false
	}
	g.Player.MP -= 1
	g.Printf("%s absorbs your mana.", m.Kind.Definite(true))
	m.Statuses[MonsExhausted] = 1
	g.PushEvent(&monsterEvent{ERank: ev.Rank() + 10 + RandInt(10), NMons: m.Index, EAction: MonsExhaustionEnd})
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) Explode(g *game, ev event) {
	neighbors := m.Pos.ValidNeighbors()
	g.MakeNoise(WallNoise, m.Pos)
	g.Printf("%s %s blows with a noisy pop.", g.ExplosionSound(), m.Kind.Definite(true))
	g.ui.ExplosionAnimation(g, FireExplosion, m.Pos)
	for _, pos := range append(neighbors, m.Pos) {
		c := g.Dungeon.Cell(pos)
		if c.T == FreeCell {
			g.Burn(pos, ev)
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.HP /= 2
			if mons.HP == 0 {
				mons.HP = 1
			}
			g.MakeNoise(ExplosionHitNoise, mons.Pos)
			mons.MakeHuntIfHurt(g)
		} else if g.Player.Pos == pos {
			dmg := g.Player.HP / 2
			m.InflictDamage(g, dmg, 15)
		} else if c.T == WallCell && RandInt(2) == 0 {
			g.Dungeon.SetCell(pos, FreeCell)
			g.Stats.Digs++
			if !g.Player.LOS[pos] {
				g.WrongWall[pos] = true
			} else {
				g.ui.WallExplosionAnimation(g, pos)
			}
			g.MakeNoise(WallNoise, pos)
			g.Fog(pos, 1, ev)
		}
	}
}

func (m *monster) MakeHunt(g *game) {
	m.State = Hunting
	m.Target = g.Player.Pos
}

func (m *monster) MakeHuntIfHurt(g *game) {
	if m.Exists() && m.State != Hunting {
		m.MakeHunt(g)
		if m.State == Resting {
			g.Printf("%s awakes.", m.Kind.Definite(true))
		}
		if m.Kind == MonsHound {
			g.Printf("%s barks.", m.Kind.Definite(true))
			g.MakeNoise(BarkNoise, m.Pos)
		}
	}
}

func (m *monster) MakeAwareIfHurt(g *game) {
	if g.Player.LOS[m.Pos] && m.State != Hunting {
		m.MakeHuntIfHurt(g)
		return
	}
	if m.State != Resting {
		return
	}
	m.State = Wandering
	m.Target = g.FreeCell()
}

func (m *monster) MakeAware(g *game) {
	if !g.Player.LOS[m.Pos] {
		return
	}
	if m.State == Resting {
		adjust := (m.Pos.Distance(g.Player.Pos) - g.LosRange()/2 + 1)
		if g.Player.Aptitudes[AptStealthyLOS] {
			adjust += 1
		}
		adjust *= adjust
		r := RandInt(25 + 3*adjust)
		if g.Player.Aptitudes[AptStealthyMovement] {
			r *= 2
		}
		if r > 5 {
			return
		}
	}
	if m.State == Wandering {
		adjust := (m.Pos.Distance(g.Player.Pos) - g.LosRange()/2 + 1)
		if g.Player.Aptitudes[AptStealthyLOS] {
			adjust += 1
		}
		adjust *= adjust
		r := RandInt(30 + adjust)
		if g.Player.Aptitudes[AptStealthyMovement] {
			r += 5
		}
		if r >= 25 {
			return
		}
	}
	if m.State == Resting {
		g.Printf("%s awakes.", m.Kind.Definite(true))
	}
	if m.State == Wandering {
		g.Printf("%s notices you.", m.Kind.Definite(true))
	}
	if m.State != Hunting && m.Kind == MonsHound {
		g.Printf("%s barks.", m.Kind.Definite(true))
		g.MakeNoise(BarkNoise, m.Pos)
	}
	m.MakeHunt(g)
}

func (m *monster) Heal(g *game, ev event) {
	if m.HP < m.HPmax {
		m.HP++
	}
	ev.Renew(g, 50)
}

func (m *monster) GatherBand(g *game) {
	if !MonsBands[g.Bands[m.Band]].band {
		return
	}
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{m.Pos}, 4)
	for _, mons := range g.Monsters {
		if mons.Band == m.Band {
			if mons.State == Hunting && m.State != Hunting {
				continue
			}
			n, ok := nm[mons.Pos]
			if !ok || n.Cost > 4 {
				continue
			}
			r := RandInt(100)
			if r > 50 || mons.State == Wandering && r > 10 {
				mons.Target = m.Target
				mons.State = m.State
			}
		}
	}
}

func (g *game) MonsterAt(pos position) *monster {
	if !pos.valid() {
		return nil
	}
	i := g.MonstersPosCache[pos.idx()]
	if i <= 0 {
		return nil
	}
	return g.Monsters[i-1]
}

func (g *game) Danger() int {
	danger := 0
	for _, mons := range g.Monsters {
		danger += mons.Kind.Dangerousness()
	}
	return danger
}

func (g *game) MaxDanger() int {
	max := 18 + 9*g.Depth + g.Depth/2 + g.Depth*g.Depth/3
	adjust := -2 * g.Depth
	for c, q := range g.Player.Consumables {
		switch c {
		case HealWoundsPotion, CBlinkPotion:
			adjust += Min(5, g.Depth) * Min(q, Min(5, g.Depth))
		case TeleportationPotion, DigPotion, WallPotion:
			adjust += Min(3, g.Depth) * Min(q, 3)
		case SwiftnessPotion, LignificationPotion, MagicPotion, BerserkPotion, ExplosiveMagara:
			adjust += Min(2, g.Depth) * Min(q, 3)
		case ConfusingDart:
			adjust += Min(1, g.Depth) * Min(q, 7)
		}
	}
	for _, props := range g.Player.Rods {
		adjust += Min(props.Charge, 2) * Min(2, g.Depth)
	}
	if g.Depth < MaxDepth && g.Player.Consumables[DescentPotion] > 0 {
		adjust += g.Depth
	}
	if g.Player.Weapon == Dagger {
		adjust -= Min(3, g.Depth) * Max(1, g.Depth-2)
	}
	if g.Player.Armour == PlateArmour {
		adjust += WinDepth - Min(g.Depth, WinDepth)
	}
	if g.Depth > 3 && g.Player.Shield == NoShield && !g.Player.Weapon.TwoHanded() {
		adjust -= Min(g.Depth, 6) * 2
	}
	if g.Player.Weapon.TwoHanded() && g.Depth < 4 {
		adjust += (4 - g.Depth) * 2
	}
	if g.Player.Armour == ChainMail || g.Player.Armour == LeatherArmour {
		adjust += WinDepth/2 - g.Depth
	}
	if g.Player.Weapon != Dagger && g.Depth < 3 {
		adjust += 4 + (3-g.Depth)*3
	}
	if g.Player.Armour == Robe {
		adjust -= 3 * g.Depth / 2
	}
	if max+adjust < max-max/3 {
		max = max - max/3
	} else if max+adjust > max+max/3 {
		max = max + max/3
	} else {
		max = max + adjust
	}
	if WinDepth-g.Depth < g.Player.Consumables[MagicMappingPotion] {
		max = max * 110 / 100
	}
	if WinDepth-g.Depth < g.Player.Consumables[DreamPotion] {
		max = max * 105 / 100
	}
	switch g.Dungeon.Gen {
	case GenCaveMapTree:
		max = max * 90 / 100
	case GenCaveMap:
		max = max * 95 / 100
	case GenRoomMap:
		max = max * 105 / 100
	case GenRuinsMap:
		max = max * 108 / 100
	case GenBSPMap:
		max = max * 115 / 100
	}
	return max
}

func (g *game) MaxMonsters() int {
	max := 13 + 3*g.Depth
	if max > 33 && g.Depth <= WinDepth {
		max = 33
	} else if max > 36 {
		max = 36
	}
	switch g.Dungeon.Gen {
	case GenCaveMapTree, GenCaveMap:
		max = max * 90 / 100
	case GenBSPMap:
		max = max * 110 / 100
	}
	return max
}

func (g *game) GenMonsters() {
	g.Monsters = []*monster{}
	g.Bands = []monsterBand{}
	danger := g.MaxDanger()
	nmons := g.MaxMonsters()
	nband := 0
	i := 0
	repeat := 0
loop:
	for danger > 0 && nmons > 0 {
		for band, data := range MonsBands {
			if RandInt(data.rarity*2) != 0 {
				continue
			}
			monsters := g.GenBand(data, monsterBand(band))
			if monsters == nil {
				continue
			}
			g.GeneratedBands[monsterBand(band)]++
			g.Bands = append(g.Bands, monsterBand(band))
			pos := g.FreeCellForMonster()
			for _, mk := range monsters {
				if nmons-1 <= 0 {
					return
				}
				if danger-mk.Dangerousness() <= 0 {
					if repeat > 10 {
						return
					}
					repeat++
					continue loop
				}
				danger -= mk.Dangerousness()
				nmons--
				mons := &monster{Kind: mk}
				mons.Init()
				mons.Index = i
				mons.Band = nband
				mons.PlaceAt(g, pos)
				g.Monsters = append(g.Monsters, mons)
				i++
				pos = g.FreeCellForBandMonster(pos)
			}
			nband++
		}
	}
}

func (g *game) MonsterInLOS() *monster {
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.Pos] {
			return mons
		}
	}
	return nil
}
