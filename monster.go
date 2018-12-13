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
	MonsSlow
	MonsLignified
)

const NMonsStatus = int(MonsLignified) + 1

func (st monsterStatus) String() (text string) {
	switch st {
	case MonsConfused:
		text = "confused"
	case MonsExhausted:
		text = "exhausted"
	case MonsSlow:
		text = "slowed"
	case MonsLignified:
		text = "lignified"
	}
	return text
}

type monsterKind int

const (
	MonsGoblin monsterKind = iota
	MonsTinyHarpy
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
	MonsWingedMilfid
	MonsBlinkingFrog
	MonsLich
	MonsEarthDragon
	MonsMirrorSpecter
	MonsAcidMound
	MonsExplosiveNadre
	MonsSatowalgaPlant
	MonsMadNixe
	MonsMindCelmist
	MonsVampire
	MonsTreeMushroom
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
	case MonsLich, MonsCyclop, MonsGoblinWarrior, MonsSatowalgaPlant, MonsMadNixe, MonsVampire, MonsTreeMushroom:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Smiting() bool {
	switch mk {
	case MonsMirrorSpecter, MonsMindCelmist:
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

func (mk monsterKind) Living() bool {
	switch mk {
	case MonsLich, MonsSkeletonWarrior, MonsMarevorHelith:
		return false
	default:
		return true
	}
}

type monsterData struct {
	movementDelay int
	baseAttack    int
	attackDelay   int
	maxHP         int
	letter        rune
	name          string
	dangerousness int
}

var MonsData = []monsterData{
	MonsGoblin:          {10, 1, 10, 2, 'g', "goblin", 3},
	MonsTinyHarpy:       {10, 1, 10, 2, 't', "tiny harpy", 4},
	MonsOgre:            {10, 2, 20, 3, 'O', "ogre", 7},
	MonsCyclop:          {10, 2, 20, 3, 'C', "cyclops", 9},
	MonsWorm:            {15, 1, 10, 3, 'w', "farmer worm", 4},
	MonsBrizzia:         {15, 1, 10, 3, 'z', "brizzia", 6},
	MonsAcidMound:       {10, 1, 10, 2, 'a', "acid mound", 6},
	MonsHound:           {10, 1, 10, 2, 'h', "hound", 5},
	MonsYack:            {10, 1, 10, 2, 'y', "yack", 5},
	MonsGiantBee:        {5, 1, 10, 1, 'B', "giant bee", 6},
	MonsGoblinWarrior:   {10, 1, 10, 2, 'G', "goblin warrior", 5},
	MonsHydra:           {10, 1, 10, 4, 'H', "hydra", 10},
	MonsSkeletonWarrior: {10, 1, 10, 3, 'S', "skeleton warrior", 6},
	MonsSpider:          {10, 1, 10, 2, 's', "spider", 6},
	MonsWingedMilfid:    {10, 1, 10, 2, 'W', "winged milfid", 6},
	MonsBlinkingFrog:    {10, 1, 10, 2, 'F', "blinking frog", 6},
	MonsLich:            {10, 1, 10, 2, 'L', "lich", 15},
	MonsEarthDragon:     {10, 2, 10, 4, 'D', "earth dragon", 20},
	MonsMirrorSpecter:   {10, 1, 10, 2, 'm', "mirror specter", 11},
	MonsExplosiveNadre:  {10, 1, 10, 1, 'n', "explosive nadre", 6},
	MonsSatowalgaPlant:  {10, 1, 10, 3, 'P', "satowalga plant", 7},
	MonsMadNixe:         {10, 1, 10, 2, 'N', "mad nixe", 14},
	MonsMindCelmist:     {10, 1, 20, 2, 'c', "mind celmist", 12},
	MonsVampire:         {10, 1, 10, 2, 'V', "vampire", 13},
	MonsTreeMushroom:    {15, 2, 20, 4, 'T', "tree mushroom", 16},
	MonsMarevorHelith:   {10, 0, 10, 10, 'M', "Marevor Helith", 18},
}

var monsDesc = []string{
	MonsGoblin:          "Goblins are little humanoid creatures. They often appear in a group.",
	MonsTinyHarpy:       "Tiny harpies are little humanoid flying creatures. They blink away when hurt. They often appear in a group.",
	MonsOgre:            "Ogres are big clunky humanoids that can hit really hard.",
	MonsCyclop:          "Cyclopes are very similar to ogres, but they also like to throw rocks at their foes (2 damage). The rocks can block your way for a while.",
	MonsWorm:            "Farmer worms are ugly slow moving creatures, but surprisingly hardy at times, and they furrow as they move, helping new foliage to grow.",
	MonsBrizzia:         "Brizzias are big slow moving biped creatures. They are quite hardy, and when hurt they can cause nausea, impeding the use of potions.",
	MonsAcidMound:       "Acid mounds are acidic creatures. They can temporarily corrode your equipment.",
	MonsHound:           "Hounds are fast moving carnivore quadrupeds. They can bark, and smell you.",
	MonsYack:            "Yacks are quite large herbivorous quadrupeds. They tend to form large groups, and can push you one cell away.",
	MonsGiantBee:        "Giant bees are fragile but extremely fast moving creatures. Their bite can sometimes enrage you.",
	MonsGoblinWarrior:   "Goblin warriors are goblins that learned to fight, and got equipped with leather armour. They can throw javelins.",
	MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	MonsSkeletonWarrior: "Skeleton warriors are good fighters, clad in chain mail.",
	MonsSpider:          "Spiders are fast moving fragile creatures, whose bite can confuse you.",
	MonsWingedMilfid:    "Winged milfids are fast moving humanoids that can fly over you and make you swap positions. They tend to be very agressive creatures.",
	MonsBlinkingFrog:    "Blinking frogs are big frog-like creatures, whose bite can make you blink away.",
	MonsLich:            "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you, halving your HP.",
	MonsEarthDragon:     "Earth dragons are big and hardy creatures that wander in the Underground. It is said they can be credited for many of the tunnels.",
	MonsMirrorSpecter:   "Mirror specters are very insubstantial creatures, which can absorb your mana.",
	MonsExplosiveNadre:  "Explosive nadres are very frail creatures that explode upon dying, halving HP of any adjacent creatures and occasionally destroying walls.",
	MonsSatowalgaPlant:  "Satowalga Plants are immobile bushes that throw acidic projectiles at you, sometimes corroding and confusing you.",
	MonsMadNixe:         "Mad nixes are magical humanoids that can attract you to them.",
	MonsMindCelmist:     "Mind celmists are mages that use magical smitting mind attacks that bypass armour. They can occasionally confuse or slow you. They try to avoid melee.",
	MonsVampire:         "Vampires are humanoids that drink blood to survive. Their spitting can cause nausea, impeding the use of potions.",
	MonsTreeMushroom:    "Tree mushrooms are big clunky slow-moving creatures. They can throw lignifying spores at you.",
	MonsMarevorHelith:   "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past.",
}

type bandInfo struct {
	Path []position
	I    int
	Kind monsterBand
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
	LoneMilfid
	LoneBlinkingFrog
	LoneCyclop
	LoneLich
	LoneEarthDragon
	LoneSpecter
	LoneAcidMound
	LoneExplosiveNadre
	LoneSatowalgaPlant
	LoneMindCelmist
	LoneVampire
	LoneTreeMushroom
	LoneEarlyNixe
	LoneEarlyAcidMound
	LoneEarlyBrizzia
	LoneEarlySpecter
	LoneEarlySatowalgaPlant
	LoneEarlyEarthDragon
	LoneEarlyHydra
	LoneEarlyLich
	LoneEarlyMindCelmist
	LoneEarlyVampire
	LoneEarlyTreeMushroom
	BandGoblins
	BandGoblinsMany
	BandGoblinsHound
	BandGoblinsOgre
	BandGoblinsWithWarriors
	BandGoblinsWithWarriorsMilfid
	BandGoblinsWithWarriorsHound
	BandGoblinsWithWarriorsOgre
	BandGoblinWarriors
	BandGoblinWarriorsMilfid
	BandHounds
	BandHoundsMany
	BandYacksGoblin
	BandYacksMilfid
	BandYacksMany
	BandSpiders
	BandSpidersMilfid
	BandWingedMilfids
	BandSatowalga
	BandBlinkingFrogs
	BandExplosiveFrog
	BandExplosiveBrizzia
	BandGiantBees
	BandGiantBeesMany
	BandSkeletonWarrior
	BandTreeMushroomWorms
	BandTreeMushrooms
	BandMindCelmists
	BandMindCelmistsLich
	BandMindCelmistsMadNixe
	BandMadNixes
	BandMadNixesDragon
	BandMadNixesHydra
	BandMadNixesFrogs
	BandVampires
	BandVampireNixe
	BandVampireCelmist
	UBandTinyHarpy
	UBandWorms
	UBandGoblinsEasy
	UBandFrogs
	UBandOgres
	UBandGoblins
	UBandBeeYacks
	UBandMadNixes
	UBandMindCelmist
	UHydras
	UExplosiveNadres
	ULich
	UVampires
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
	UXSatowalgaNixe
	UXSpecters
	UXDisabling
	UXMadNixeSpecter
	UXMadNixeCyclop
	UXMadNixeHydra
	UXMadNixes
	UXVampires
	UXTreeMushrooms
	UXMindCelmists
	UXMilfidYack
	UXYacks
	UXVariedWarriors
)

type monsInterval struct {
	Min int
	Max int
}

type monsterBandData struct {
	Distribution map[monsterKind]monsInterval
	Rarity       int
	MinDepth     int
	MaxDepth     int
	Band         bool
	Monster      monsterKind
	Unique       bool
}

func (g *game) GenBand(mbd monsterBandData, band monsterBand) []monsterKind {
	if g.GeneratedUniques[band] > 0 && mbd.Unique {
		return nil
	}
	if g.Depth > mbd.MaxDepth {
		return nil
	}
	if g.Depth < mbd.MinDepth {
		return nil
	}
	if !mbd.Band {
		return []monsterKind{mbd.Monster}
	}
	bandMonsters := []monsterKind{}
	for m, interval := range mbd.Distribution {
		for i := 0; i < interval.Min+RandInt(interval.Max-interval.Min+1); i++ {
			bandMonsters = append(bandMonsters, m)
		}
	}
	return bandMonsters
}

var MonsBands = []monsterBandData{
	LoneGoblin:              {Rarity: 2, MinDepth: 1, MaxDepth: 2, Monster: MonsGoblin},
	LoneOgre:                {Rarity: 4, MinDepth: 2, MaxDepth: 7, Monster: MonsOgre},
	LoneWorm:                {Rarity: 2, MinDepth: 1, MaxDepth: 3, Monster: MonsWorm},
	LoneRareWorm:            {Rarity: 13, MinDepth: 4, MaxDepth: WinDepth + 1, Monster: MonsWorm},
	LoneBrizzia:             {Rarity: 13, MinDepth: 4, MaxDepth: WinDepth + 1, Monster: MonsBrizzia},
	LoneHound:               {Rarity: 5, MinDepth: 1, MaxDepth: 5, Monster: MonsHound},
	LoneHydra:               {Rarity: 10, MinDepth: 5, MaxDepth: WinDepth + 1, Monster: MonsHydra},
	LoneSpider:              {Rarity: 3, MinDepth: 3, MaxDepth: WinDepth + 1, Monster: MonsSpider},
	LoneMilfid:              {Rarity: 13, MinDepth: 3, MaxDepth: WinDepth + 1, Monster: MonsWingedMilfid},
	LoneBlinkingFrog:        {Rarity: 7, MinDepth: 3, MaxDepth: WinDepth + 1, Monster: MonsBlinkingFrog},
	LoneCyclop:              {Rarity: 4, MinDepth: 3, MaxDepth: WinDepth + 1, Monster: MonsCyclop},
	LoneLich:                {Rarity: 8, MinDepth: 5, MaxDepth: WinDepth + 1, Monster: MonsLich},
	LoneEarthDragon:         {Rarity: 9, MinDepth: 6, MaxDepth: WinDepth + 1, Monster: MonsEarthDragon},
	LoneSpecter:             {Rarity: 7, MinDepth: 4, MaxDepth: WinDepth + 1, Monster: MonsMirrorSpecter},
	LoneAcidMound:           {Rarity: 7, MinDepth: 4, MaxDepth: WinDepth + 1, Monster: MonsAcidMound},
	LoneExplosiveNadre:      {Rarity: 5, MinDepth: 2, MaxDepth: 4, Monster: MonsExplosiveNadre},
	LoneSatowalgaPlant:      {Rarity: 9, MinDepth: 4, MaxDepth: WinDepth + 1, Monster: MonsSatowalgaPlant},
	LoneMindCelmist:         {Rarity: 12, MinDepth: 5, MaxDepth: WinDepth + 1, Monster: MonsMindCelmist},
	LoneVampire:             {Rarity: 12, MinDepth: 5, MaxDepth: WinDepth + 1, Monster: MonsVampire},
	LoneTreeMushroom:        {Rarity: 15, MinDepth: 6, MaxDepth: WinDepth + 1, Monster: MonsTreeMushroom},
	LoneEarlyNixe:           {Rarity: 20, MinDepth: 1, MaxDepth: 4, Monster: MonsMadNixe, Unique: true},
	LoneEarlyVampire:        {Rarity: 30, MinDepth: 2, MaxDepth: 4, Monster: MonsVampire, Unique: true},
	LoneEarlyAcidMound:      {Rarity: 20, MinDepth: 1, MaxDepth: 3, Monster: MonsAcidMound, Unique: true},
	LoneEarlyBrizzia:        {Rarity: 20, MinDepth: 1, MaxDepth: 3, Monster: MonsBrizzia, Unique: true},
	LoneEarlySpecter:        {Rarity: 20, MinDepth: 1, MaxDepth: 3, Monster: MonsMirrorSpecter, Unique: true},
	LoneEarlySatowalgaPlant: {Rarity: 20, MinDepth: 1, MaxDepth: 3, Monster: MonsSatowalgaPlant, Unique: true},
	LoneEarlyEarthDragon:    {Rarity: 30, MinDepth: 4, MaxDepth: 5, Monster: MonsEarthDragon, Unique: true},
	LoneEarlyHydra:          {Rarity: 30, MinDepth: 3, MaxDepth: 4, Monster: MonsHydra, Unique: true},
	LoneEarlyLich:           {Rarity: 30, MinDepth: 3, MaxDepth: 4, Monster: MonsLich, Unique: true},
	LoneEarlyMindCelmist:    {Rarity: 30, MinDepth: 3, MaxDepth: 4, Monster: MonsMindCelmist, Unique: true},
	LoneEarlyTreeMushroom:   {Rarity: 30, MinDepth: 4, MaxDepth: 5, Monster: MonsTreeMushroom, Unique: true},
	BandGoblins: {
		Distribution: map[monsterKind]monsInterval{MonsGoblin: {2, 3}},
		Rarity:       2, MinDepth: 1, MaxDepth: 3, Band: true,
	},
	BandGoblinsMany: {
		Distribution: map[monsterKind]monsInterval{MonsGoblin: {4, 4}},
		Rarity:       7, MinDepth: 2, MaxDepth: 3, Band: true,
	},
	BandGoblinsHound: {
		Distribution: map[monsterKind]monsInterval{MonsGoblin: {2, 2}, MonsHound: {1, 1}},
		Rarity:       4, MinDepth: 1, MaxDepth: 3, Band: true,
	},
	BandGoblinsOgre: {
		Distribution: map[monsterKind]monsInterval{MonsGoblin: {1, 1}, MonsOgre: {1, 1}},
		Rarity:       7, MinDepth: 2, MaxDepth: 3, Band: true,
	},
	BandGoblinsWithWarriors: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 3},
			MonsGoblinWarrior: {2, 2}},
		Rarity: 7, MinDepth: 4, MaxDepth: 5, Band: true,
	},
	BandGoblinsWithWarriorsMilfid: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 3},
			MonsGoblinWarrior: {1, 1},
			MonsWingedMilfid:  {1, 1}},
		Rarity: 8, MinDepth: 4, MaxDepth: 5, Band: true,
	},
	BandGoblinsWithWarriorsHound: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 3},
			MonsGoblinWarrior: {1, 1},
			MonsHound:         {1, 1}},
		Rarity: 7, MinDepth: 4, MaxDepth: 5, Band: true,
	},
	BandGoblinsWithWarriorsOgre: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 3},
			MonsGoblinWarrior: {1, 1},
			MonsOgre:          {1, 1}},
		Rarity: 7, MinDepth: 4, MaxDepth: 5, Band: true,
	},
	BandGoblinWarriors: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {1, 1},
			MonsGoblinWarrior: {3, 3}},
		Rarity: 10, MinDepth: 6, MaxDepth: WinDepth + 1, Band: true,
	},
	BandGoblinWarriorsMilfid: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {1, 1},
			MonsGoblinWarrior: {2, 2},
			MonsWingedMilfid:  {1, 1}},
		Rarity: 10, MinDepth: 6, MaxDepth: WinDepth + 1, Band: true,
	},
	BandHounds: {
		Distribution: map[monsterKind]monsInterval{MonsHound: {2, 2}, MonsGoblin: {1, 1}},
		Rarity:       6, MinDepth: 2, MaxDepth: 6, Band: true,
	},
	BandHoundsMany: {
		Distribution: map[monsterKind]monsInterval{MonsHound: {3, 3}},
		Rarity:       10, MinDepth: 2, MaxDepth: 6, Band: true,
	},
	BandSpiders: {
		Distribution: map[monsterKind]monsInterval{MonsSpider: {2, 3}},
		Rarity:       4, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandSpidersMilfid: {
		Distribution: map[monsterKind]monsInterval{MonsSpider: {2, 2}, MonsWingedMilfid: {1, 1}},
		Rarity:       7, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandWingedMilfids: {
		Distribution: map[monsterKind]monsInterval{MonsWingedMilfid: {2, 3}},
		Rarity:       9, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandBlinkingFrogs: {
		Distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: {2, 4}},
		Rarity:       7, MinDepth: 5, MaxDepth: WinDepth + 1, Band: true,
	},
	BandSatowalga: {
		Distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {2, 2},
		},
		Rarity: 10, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandExplosiveFrog: {
		Distribution: map[monsterKind]monsInterval{
			MonsBlinkingFrog:   {1, 1},
			MonsExplosiveNadre: {2, 2},
			MonsGiantBee:       {1, 1},
		},
		Rarity: 10, MinDepth: 5, MaxDepth: WinDepth + 1, Band: true,
	},
	BandExplosiveBrizzia: {
		Distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {2, 2},
			MonsGiantBee:       {1, 1},
			MonsBrizzia:        {1, 1},
		},
		Rarity: 10, MinDepth: 5, MaxDepth: WinDepth + 1, Band: true,
	},
	BandYacksGoblin: {
		Distribution: map[monsterKind]monsInterval{MonsYack: {2, 2}, MonsGoblin: {1, 1}},
		Rarity:       5, MinDepth: 3, MaxDepth: WinDepth - 1, Band: true,
	},
	BandYacksMilfid: {
		Distribution: map[monsterKind]monsInterval{MonsYack: {2, 2}, MonsWingedMilfid: {1, 1}},
		Rarity:       8, MinDepth: 3, MaxDepth: WinDepth - 1, Band: true,
	},
	BandYacksMany: {
		Distribution: map[monsterKind]monsInterval{MonsYack: {4, 5}},
		Rarity:       5, MinDepth: 4, MaxDepth: WinDepth - 1, Band: true,
	},
	BandGiantBees: {
		Distribution: map[monsterKind]monsInterval{MonsGiantBee: {2, 3}},
		Rarity:       6, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandGiantBeesMany: {
		Distribution: map[monsterKind]monsInterval{MonsGiantBee: {4, 5}},
		Rarity:       9, MinDepth: 4, MaxDepth: WinDepth + 1, Band: true,
	},
	BandSkeletonWarrior: {
		Distribution: map[monsterKind]monsInterval{MonsSkeletonWarrior: {2, 3}},
		Rarity:       7, MinDepth: 5, MaxDepth: WinDepth + 1, Band: true,
	},
	BandTreeMushroomWorms: {
		Distribution: map[monsterKind]monsInterval{
			MonsTreeMushroom: {1, 1},
			MonsWorm:         {2, 2},
		},
		Rarity: 10, MinDepth: 6, MaxDepth: WinDepth, Band: true,
	},
	BandVampires: {
		Distribution: map[monsterKind]monsInterval{
			MonsVampire: {2, 2},
		},
		Rarity: 10, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandVampireNixe: {
		Distribution: map[monsterKind]monsInterval{
			MonsVampire: {1, 1},
			MonsMadNixe: {1, 1},
		},
		Rarity: 10, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandVampireCelmist: {
		Distribution: map[monsterKind]monsInterval{
			MonsVampire:     {1, 1},
			MonsMindCelmist: {1, 1},
		},
		Rarity: 10, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandTreeMushrooms: {
		Distribution: map[monsterKind]monsInterval{
			MonsTreeMushroom: {2, 2},
			MonsWorm:         {1, 1},
		},
		Rarity: 10, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMindCelmists: {
		Distribution: map[monsterKind]monsInterval{
			MonsMindCelmist:   {1, 1},
			MonsGoblinWarrior: {1, 1},
		},
		Rarity: 8, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMindCelmistsLich: {
		Distribution: map[monsterKind]monsInterval{
			MonsMindCelmist: {2, 2},
		},
		Rarity: 8, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMindCelmistsMadNixe: {
		Distribution: map[monsterKind]monsInterval{
			MonsMindCelmist: {1, 1},
			MonsMadNixe:     {1, 1},
		},
		Rarity: 8, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMadNixes: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe: {1, 1},
			MonsSpider:  {1, 1},
			MonsHound:   {1, 1},
		},
		Rarity: 4, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMadNixesDragon: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe:     {1, 1},
			MonsEarthDragon: {1, 1},
		},
		Rarity: 4, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMadNixesHydra: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe: {1, 1},
			MonsHydra:   {1, 1},
		},
		Rarity: 4, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	BandMadNixesFrogs: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe:      {1, 1},
			MonsBlinkingFrog: {2, 2},
		},
		Rarity: 4, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true,
	},
	UBandTinyHarpy: {
		Distribution: map[monsterKind]monsInterval{
			MonsTinyHarpy:    {3, 3},
			MonsWingedMilfid: {1, 1},
		},
		Rarity: 6, MinDepth: 2, MaxDepth: 2, Band: true, Unique: true,
	},
	UBandWorms: {
		Distribution: map[monsterKind]monsInterval{MonsWorm: {3, 4}, MonsSpider: {1, 1}},
		Rarity:       8, MinDepth: 2, MaxDepth: 3, Band: true, Unique: true,
	},
	UBandGoblinsEasy: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin: {3, 3},
			MonsHound:  {2, 2},
		},
		Rarity: 4, MinDepth: 3, MaxDepth: 3, Band: true, Unique: true,
	},
	UBandFrogs: {
		Distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: {2, 3}},
		Rarity:       7, MinDepth: 4, MaxDepth: 4, Band: true, Unique: true,
	},
	UBandOgres: {
		Distribution: map[monsterKind]monsInterval{MonsOgre: {2, 3}, MonsCyclop: {1, 1}},
		Rarity:       4, MinDepth: 4, MaxDepth: 4, Band: true, Unique: true,
	},
	UBandGoblins: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblin:        {3, 3},
			MonsGoblinWarrior: {2, 2},
			MonsHound:         {1, 1},
		},
		Rarity: 4, MinDepth: 5, MaxDepth: 5, Band: true, Unique: true,
	},
	UBandBeeYacks: {
		Distribution: map[monsterKind]monsInterval{
			MonsYack:     {3, 4},
			MonsGiantBee: {2, 2},
		},
		Rarity: 5, MinDepth: 5, MaxDepth: 5, Band: true, Unique: true,
	},
	UBandMadNixes: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe: {2, 2},
			MonsSpider:  {1, 1},
		},
		Rarity: 5, MinDepth: 5, MaxDepth: 5, Band: true, Unique: true,
	},
	UVampires: {
		Distribution: map[monsterKind]monsInterval{
			MonsVampire:      {2, 2},
			MonsWingedMilfid: {1, 1},
		},
		Rarity: 10, MinDepth: 5, MaxDepth: 5, Band: true, Unique: true,
	},
	UHydras: {
		Distribution: map[monsterKind]monsInterval{
			MonsHydra:  {2, 2},
			MonsSpider: {2, 2},
		},
		Rarity: 5, MinDepth: 6, MaxDepth: 6, Band: true, Unique: true,
	},
	UExplosiveNadres: {
		Distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {2, 3},
			MonsBrizzia:        {1, 2},
		},
		Rarity: 6, MinDepth: 6, MaxDepth: 6, Band: true, Unique: true,
	},
	ULich: {
		Distribution: map[monsterKind]monsInterval{
			MonsSkeletonWarrior: {2, 2},
			MonsLich:            {1, 1},
			MonsMirrorSpecter:   {0, 1},
		},
		Rarity: 6, MinDepth: WinDepth - 1, MaxDepth: WinDepth - 1, Band: true, Unique: true,
	},
	UBrizzias: {
		Distribution: map[monsterKind]monsInterval{
			MonsBrizzia: {3, 4},
		},
		Rarity: 8, MinDepth: WinDepth - 1, MaxDepth: WinDepth - 1, Band: true, Unique: true,
	},
	UBandMindCelmist: {
		Distribution: map[monsterKind]monsInterval{
			MonsMindCelmist: {2, 2},
			MonsHound:       {1, 1},
		},
		Rarity: 10, MinDepth: WinDepth - 1, MaxDepth: WinDepth - 1, Band: true, Unique: true,
	},
	UAcidMounds: {
		Distribution: map[monsterKind]monsInterval{
			MonsAcidMound: {3, 4},
		},
		Rarity: 8, MinDepth: WinDepth, MaxDepth: WinDepth, Band: true, Unique: true,
	},
	USatowalga: {
		Distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {3, 3},
		},
		Rarity: 8, MinDepth: WinDepth, MaxDepth: WinDepth, Band: true, Unique: true,
	},
	UDragon: {
		Distribution: map[monsterKind]monsInterval{
			MonsEarthDragon: {2, 2},
		},
		Rarity: 6, MinDepth: WinDepth, MaxDepth: WinDepth, Band: true, Unique: true,
	},
	UMarevorHelith: {
		Distribution: map[monsterKind]monsInterval{
			MonsMarevorHelith: {1, 1},
			MonsLich:          {0, 1},
			MonsVampire:       {0, 1},
		},
		Rarity: 13, MinDepth: 2, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXCyclops: {
		Distribution: map[monsterKind]monsInterval{
			MonsCyclop: {3, 3},
		},
		Rarity: 6, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXLiches: {
		Distribution: map[monsterKind]monsInterval{
			MonsLich: {2, 2},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXFrogRanged: {
		Distribution: map[monsterKind]monsInterval{
			MonsBlinkingFrog: {2, 2},
			MonsCyclop:       {1, 1},
			MonsLich:         {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXExplosive: {
		Distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {5, 5},
		},
		Rarity: 6, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXWarriors: {
		Distribution: map[monsterKind]monsInterval{
			MonsHound:         {2, 2},
			MonsGoblinWarrior: {3, 3},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXSatowalgaNixe: {
		Distribution: map[monsterKind]monsInterval{
			MonsSatowalgaPlant: {2, 2},
			MonsMadNixe:        {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXSpecters: {
		Distribution: map[monsterKind]monsInterval{
			MonsMirrorSpecter: {3, 3},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXDisabling: {
		Distribution: map[monsterKind]monsInterval{
			MonsExplosiveNadre: {1, 1},
			MonsSpider:         {1, 1},
			MonsBrizzia:        {1, 1},
			MonsGiantBee:       {1, 1},
			MonsMirrorSpecter:  {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMadNixeSpecter: {
		Distribution: map[monsterKind]monsInterval{
			MonsMirrorSpecter: {1, 1},
			MonsMadNixe:       {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMadNixeCyclop: {
		Distribution: map[monsterKind]monsInterval{
			MonsCyclop:  {1, 1},
			MonsMadNixe: {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMadNixeHydra: {
		Distribution: map[monsterKind]monsInterval{
			MonsHydra:   {1, 1},
			MonsMadNixe: {1, 1},
		},
		Rarity: 6, MinDepth: MaxDepth, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXVampires: {
		Distribution: map[monsterKind]monsInterval{
			MonsVampire: {3, 3},
		},
		Rarity: 10, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMadNixes: {
		Distribution: map[monsterKind]monsInterval{
			MonsMadNixe: {3, 3},
		},
		Rarity: 10, MinDepth: MaxDepth - 2, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMindCelmists: {
		Distribution: map[monsterKind]monsInterval{
			MonsMindCelmist: {2, 2},
			MonsCyclop:      {1, 1},
		},
		Rarity: 8, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXTreeMushrooms: {
		Distribution: map[monsterKind]monsInterval{
			MonsTreeMushroom: {3, 3},
		},
		Rarity: 10, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXMilfidYack: {
		Distribution: map[monsterKind]monsInterval{
			MonsWingedMilfid: {2, 2},
			MonsYack:         {3, 3},
		},
		Rarity: 6, MinDepth: MaxDepth - 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXYacks: {
		Distribution: map[monsterKind]monsInterval{
			MonsYack: {7, 7},
		},
		Rarity: 8, MinDepth: MaxDepth - 2, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
	UXVariedWarriors: {
		Distribution: map[monsterKind]monsInterval{
			MonsGoblinWarrior:   {1, 1},
			MonsWingedMilfid:    {1, 1},
			MonsSkeletonWarrior: {1, 1},
		},
		Rarity: 6, MinDepth: WinDepth + 1, MaxDepth: MaxDepth, Band: true, Unique: true,
	},
}

type specialBands struct {
	bands    []monsterBandData
	minDepth int
	maxDepth int
}

var MonsSpecialBands []specialBands
var MonsSpecialEndBands []specialBands

func init() {
	MonsSpecialBands = []specialBands{
		{bands: []monsterBandData{ // ogres easy
			{Monster: MonsOgre, Rarity: 20},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsOgre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {1, 1}, MonsOgre: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 4, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // spiders
			{Monster: MonsSpider, Rarity: 40},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSpider: {4, 4},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSpider: {2, 2}, MonsBrizzia: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSpider: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 8, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // milfids
			{Monster: MonsWingedMilfid, Rarity: 50},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsWingedMilfid: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsWingedMilfid: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsWingedMilfid: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsWingedMilfid: {1, 1}, MonsYack: {3, 3},
			}, Rarity: 4, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // Bees
			{Monster: MonsGiantBee, Rarity: 50},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsGiantBee: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsCyclop: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsGiantBee: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsYack: {3, 3},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsBrizzia: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsHydra: {1, 1},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 8, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // goblins
			{Monster: MonsGoblin, Rarity: 4},
			{Monster: MonsGoblinWarrior, Rarity: 5},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsGoblin: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {3, 3}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {2, 2}, MonsGoblinWarrior: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsGoblin: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblin: {2, 2}, MonsYack: {3, 3},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 10, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // explosive nadres
			{Monster: MonsExplosiveNadre, Rarity: 4},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsExplosiveNadre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsGiantBee: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsGoblinWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsExplosiveNadre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsYack: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 7, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsEarthDragon: {1, 1},
			}, Rarity: 10, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // plants
			{Monster: MonsSatowalgaPlant, Rarity: 4},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsWorm: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {3, 3}, MonsSatowalgaPlant: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {1, 1}, MonsGiantBee: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsSatowalgaPlant: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {2, 2}, MonsSpider: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsSatowalgaPlant: {2, 2},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {2, 2}, MonsWingedMilfid: {1, 1},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {1, 1}, MonsMadNixe: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {2, 2}, MonsBlinkingFrog: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {1, 1}, MonsTreeMushroom: {1, 1},
			}, Rarity: 10, Band: true},
		},
			minDepth: 4,
			maxDepth: 7,
		},
		{bands: []monsterBandData{ // acid mounds
			{Monster: MonsAcidMound, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsAcidMound: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {3, 3}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsGoblinWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsAcidMound: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsYack: {2, 2},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 8, Band: true},
		},
			minDepth: 4,
			maxDepth: WinDepth,
		},
		{bands: []monsterBandData{ // blinking frogs
			{Monster: MonsBlinkingFrog, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {3, 3}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2}, MonsGoblinWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2}, MonsYack: {2, 2},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBlinkingFrog: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 10, Band: true},
		},
			minDepth: 4,
			maxDepth: WinDepth,
		},
		{bands: []monsterBandData{ // hydras
			{Monster: MonsHydra, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsWorm: {3, 3}, MonsSpider: {2, 2},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsGoblin: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsSkeletonWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsMadNixe: {1, 1},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {2, 2}, MonsMirrorSpecter: {1, 1},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsTreeMushroom: {1, 1},
			}, Rarity: 8, Band: true},
		},
			minDepth: 5,
			maxDepth: WinDepth,
		},
		{bands: []monsterBandData{ // liches
			{Monster: MonsLich, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsSkeletonWarrior: {1, 2}, MonsHound: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSkeletonWarrior: {1, 2}, MonsAcidMound: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsGoblin: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsSkeletonWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsVampire: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsBlinkingFrog: {1, 1}, MonsWingedMilfid: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsMadNixe: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsMirrorSpecter: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {2, 2}, MonsSkeletonWarrior: {2, 2},
			}, Rarity: 8, Band: true},
		},
			minDepth: 6,
			maxDepth: WinDepth,
		},
		{bands: []monsterBandData{ // dragons
			{Monster: MonsEarthDragon, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {3, 3}, MonsHound: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {3, 3}, MonsAcidMound: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsSpider: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsWingedMilfid: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsExplosiveNadre: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsGoblin: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsSkeletonWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsBlinkingFrog: {1, 1}, MonsWingedMilfid: {1, 1},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsMadNixe: {1, 1},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {2, 2}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 10, Band: true},
		},
			minDepth: 6,
			maxDepth: WinDepth,
		},
	}
	for _, sb := range MonsSpecialBands {
		for i, _ := range sb.bands {
			sb.bands[i].MaxDepth = MaxDepth
		}
	}
	MonsSpecialEndBands = []specialBands{
		{bands: []monsterBandData{ // ogres terrible
			{Monster: MonsOgre, Rarity: 5},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsOgre: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {1, 1}, MonsOgre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsCyclop: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsEarthDragon: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsHydra: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsAcidMound: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsMirrorSpecter: {1, 1}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 3, Band: true},
		}},
		{bands: []monsterBandData{ // ranged terrible
			{Monster: MonsCyclop, Rarity: 5},
			{Monster: MonsLich, Rarity: 5},
			{Distribution: map[monsterKind]monsInterval{
				MonsCyclop: {2, 2}, MonsOgre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsCyclop: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {2, 2}, MonsWingedMilfid: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {2, 2}, MonsGoblinWarrior: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsCyclop: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsEarthDragon: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {3, 3}, MonsWingedMilfid: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {2, 2}, MonsTreeMushroom: {1, 1},
			}, Rarity: 5, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsCyclop: {2, 2}, MonsLich: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {2, 2}, MonsLich: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsCyclop: {1, 1},
			}, Rarity: 3, Band: true},
		}},
		{bands: []monsterBandData{ // mind celmists
			{Monster: MonsMindCelmist, Rarity: 5},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {2, 2}, MonsHound: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsMadNixe: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsLich: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsOgre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsCyclop: {1, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsYack: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsVampire: {1, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {3, 3},
			}, Rarity: 10, Band: true},
		}},
		{bands: []monsterBandData{ // nixe trap
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {2, 2}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsSatowalgaPlant: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsAcidMound: {3, 3},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsOgre: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsEarthDragon: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsHydra: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsHydra: {1, 1}, MonsEarthDragon: {1, 1},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsGiantBee: {3, 3},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {4, 4},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {2, 2}, MonsMindCelmist: {1, 1},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {2, 2}, MonsVampire: {1, 1},
			}, Rarity: 6, Band: true},
		}},
		{bands: []monsterBandData{ // blinking frogs terrible
			{Distribution: map[monsterKind]monsInterval{
				MonsMadNixe: {1, 1}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSatowalgaPlant: {1, 1}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSpider: {2, 2}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {1, 1}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsCyclop: {1, 1}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsWingedMilfid: {2, 2}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {2, 2}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGiantBee: {2, 2}, MonsBlinkingFrog: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMindCelmist: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsTreeMushroom: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 8, Band: true},
		}},
		{bands: []monsterBandData{ // yacks and brizzias terrible
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {4, 4}, MonsExplosiveNadre: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {4, 4}, MonsSpider: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {3, 3}, MonsSpider: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {3, 3}, MonsAcidMound: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {2, 2}, MonsExplosiveNadre: {1, 1}, MonsMirrorSpecter: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {1, 1}, MonsHydra: {1, 1}, MonsYack: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {3, 3}, MonsWorm: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {3, 3}, MonsBrizzia: {3, 3},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsYack: {1, 1}, MonsBrizzia: {1, 1}, MonsBlinkingFrog: {1, 1}, MonsHound: {1, 1},
			}, Rarity: 2, Band: true},
		}},
		{bands: []monsterBandData{ // terrible undead
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsSkeletonWarrior: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsMadNixe: {1, 1}, MonsSkeletonWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {2, 2}, MonsSkeletonWarrior: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSkeletonWarrior: {3, 3}, MonsMadNixe: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsSkeletonWarrior: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {1, 1}, MonsSkeletonWarrior: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsVampire: {2, 2},
			}, Rarity: 6, Band: true},
		}},
		{bands: []monsterBandData{ // terrible vampires
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsVampire: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsLich: {1, 1}, MonsMadNixe: {1, 1}, MonsVampire: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsVampire: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsVampire: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsVampire: {2, 2}, MonsMindCelmist: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsVampire: {4, 4},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsMirrorSpecter: {1, 1}, MonsVampire: {1, 1},
			}, Rarity: 2, Band: true},
		}},
		{bands: []monsterBandData{ // terrible dragon and hydras
			{Distribution: map[monsterKind]monsInterval{
				MonsBrizzia: {2, 2}, MonsExplosiveNadre: {2, 2},
			}, Rarity: 10, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsHydra: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {2, 2}, MonsSpider: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsBlinkingFrog: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsExplosiveNadre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsBrizzia: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsMirrorSpecter: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsEarthDragon: {1, 1}, MonsMindCelmist: {1, 1},
			}, Rarity: 8, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsHydra: {1, 1}, MonsTreeMushroom: {1, 1},
			}, Rarity: 8, Band: true},
		}},
		{bands: []monsterBandData{ // terrible goblin warriors
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsHound: {4, 4},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsHydra: {1, 1},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsBrizzia: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsSpider: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsMadNixe: {1, 1},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsWingedMilfid: {2, 2},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsGoblinWarrior: {2, 2}, MonsYack: {3, 3},
			}, Rarity: 3, Band: true},
		}},
		{bands: []monsterBandData{ // terrible acid mounds
			{Monster: MonsAcidMound, Rarity: 2},
			{Distribution: map[monsterKind]monsInterval{
				MonsHound: {1, 1}, MonsAcidMound: {3, 3},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {3, 3}, MonsExplosiveNadre: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsHydra: {1, 1},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsOgre: {2, 2}, MonsAcidMound: {2, 2},
			}, Rarity: 2, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsSpider: {3, 3},
			}, Rarity: 3, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {3, 3}, MonsWingedMilfid: {2, 2},
			}, Rarity: 6, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsBrizzia: {2, 2},
			}, Rarity: 4, Band: true},
			{Distribution: map[monsterKind]monsInterval{
				MonsAcidMound: {2, 2}, MonsMadNixe: {1, 1}, MonsSatowalgaPlant: {1, 1},
			}, Rarity: 8, Band: true},
		}},
	}
	for _, sb := range MonsSpecialEndBands {
		for i, _ := range sb.bands {
			sb.bands[i].MaxDepth = MaxDepth
		}
	}
}

type monster struct {
	Kind          monsterKind
	Band          int
	Index         int
	Dir           direction
	Attack        int
	HPmax         int
	HP            int
	State         monsterState
	Statuses      [NMonsStatus]int
	Pos           position
	LastKnownPos  position
	Target        position
	Path          []position // cache
	Obstructing   bool
	FireReady     bool
	Seen          bool
	LOS           map[position]bool
	LastSeenState monsterState
}

func (m *monster) Init() {
	m.HPmax = MonsData[m.Kind].maxHP
	m.Attack = MonsData[m.Kind].baseAttack
	m.HP = m.HPmax
	m.Pos = InvalidPos
	m.LastKnownPos = InvalidPos
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
	if m.Status(MonsLignified) {
		return nil
	}
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

func (m *monster) AlternateConfusedPlacement(g *game) *position {
	var neighbors []position
	neighbors = g.Dungeon.CardinalFreeNeighbors(m.Pos)
	npos := InvalidPos
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if mons.Exists() || g.Player.Pos == pos {
			continue
		}
		npos = pos
		if npos.Distance(g.Player.Pos) == 1 {
			return &npos
		}
	}
	if npos.valid() {
		return &npos
	}
	return nil
}

func (m *monster) SafePlacement(g *game) *position {
	var neighbors []position
	if m.Status(MonsConfused) {
		neighbors = g.Dungeon.CardinalFreeNeighbors(m.Pos)
	} else {
		neighbors = g.Dungeon.FreeNeighbors(m.Pos)
	}
	spos := InvalidPos
	sbest := 9
	area := make([]position, 9)
	for _, pos := range neighbors {
		if pos.Distance(g.Player.Pos) <= 1 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		// simple heuristic
		nsbest := g.Dungeon.WallAreaCount(area, pos, 1)
		if nsbest < sbest {
			sbest = nsbest
			spos = pos
		} else if nsbest == sbest {
			switch pos.Dir(g.Player.Pos) {
			case N, W, E, S:
			default:
				sbest = nsbest
				spos = pos
			}
		}
	}
	if spos.valid() {
		return &spos
	}
	return nil
}

func (m *monster) TeleportPlayer(g *game, ev event) {
	if RandInt(2) == 0 {
		g.Print("Marevor pushes you through a monolith.")
		g.StoryPrint("Marevor pushed you through a monolith.")
		g.Teleportation(ev)
	} else {
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
	if g.Player.Sees(m.Pos) {
		g.Printf("%s teleports away.", m.Kind.Definite(true))
	}
	opos := m.Pos
	m.MoveTo(g, pos)
	if g.Player.Sees(opos) {
		g.ui.TeleportAnimation(opos, pos, false)
	}
}

func (m *monster) MoveTo(g *game, pos position) {
	if g.Player.Sees(pos) {
		m.UpdateKnowledge(g, pos)
	} else if g.Player.Sees(m.Pos) {
		delete(g.LastMonsterKnownAt, m.Pos)
		m.LastKnownPos = InvalidPos
	}
	if !g.Player.Sees(m.Pos) && g.Player.Sees(pos) {
		if !m.Seen {
			m.Seen = true
			g.Printf("%s (%v) comes into view.", m.Kind.Indefinite(true), m.State)
		}
		g.StopAuto()
	}
	recomputeLOS := g.Player.Sees(m.Pos) && g.Doors[m.Pos] || g.Player.Sees(pos) && g.Doors[pos]
	m.PlaceAt(g, pos)
	if recomputeLOS {
		g.ComputeLOS()
	}
}

func (m *monster) PlaceAt(g *game, pos position) {
	if !m.Pos.valid() {
		m.Pos = pos
		g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
		m.ComputeLOS(g)
		return
	}
	if pos == m.Pos {
		// should not happen
		return
	}
	m.Dir = pos.Dir(m.Pos)
	switch m.Dir {
	case ENE, ESE:
		m.Dir = E
	case NNE, NNW:
		m.Dir = N
	case WNW, WSW:
		m.Dir = W
	case SSW, SSE:
		m.Dir = S
	}
	g.MonstersPosCache[m.Pos.idx()] = 0
	m.Pos = pos
	g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
	m.ComputeLOS(g)
}

func (m *monster) TeleportMonsterAway(g *game) bool {
	neighbors := g.Dungeon.FreeNeighbors(m.Pos)
	for _, pos := range neighbors {
		if pos == m.Pos || RandInt(3) != 0 {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			if g.Player.Sees(m.Pos) {
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
		m.Dir = g.Player.Pos.Dir(m.Pos)
		if m.Kind == MonsHydra {
			for i := 0; i <= 3; i++ {
				if RandInt(3) == 0 {
					// XXX: hydras now use a random number of heads (solve: probably remove hydras).
					m.HitPlayer(g, ev)
				}
			}
		} else if m.Kind == MonsMarevorHelith {
			m.TeleportPlayer(g, ev)
		} else {
			m.HitPlayer(g, ev)
		}
		adelay := m.Kind.AttackDelay()
		if m.Status(MonsSlow) {
			adelay += 3
		}
		ev.Renew(g, adelay)
	}
}

func (m *monster) NaturalAwake(g *game) {
	m.Target = m.NextTarget(g)
	m.State = Wandering
	m.GatherBand(g)
}

func (m *monster) NextTarget(g *game) position {
	// TODO: improve this to handle more varied cases
	band := g.Bands[m.Band]
	if len(band.Path) == 0 {
		return g.FreeCell()
	} else if len(band.Path) == 1 {
		return band.Path[0]
	}
	if band.Path[0] == m.Target {
		return band.Path[1]
	}
	return band.Path[0]
}

func (m *monster) HandleTurn(g *game, ev event) {
	ppos := g.Player.Pos
	mpos := m.Pos
	m.MakeAware(g)
	if !m.SeesPlayer(g) && m.State == Hunting {
		if g.Player.Armour == HarmonistRobe && RandInt(5) == 0 ||
			g.Player.Aptitudes[AptStealthyMovement] && RandInt(5) == 0 ||
			RandInt(10) == 0 {
			m.State = Wandering
		}
	}
	movedelay := m.Kind.MovementDelay()
	if m.Status(MonsSlow) {
		movedelay += 3
	}
	if m.State == Resting {
		if g.DepthPlayerTurn > 600 || g.DepthPlayerTurn > m.Band*50+25 {
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
	switch m.Kind {
	case MonsSatowalgaPlant:
		ev.Renew(g, movedelay)
		// oklob plants are static ranged-only
		return
	case MonsMindCelmist:
		if m.State == Hunting && !m.SeesPlayer(g) && m.Pos.Distance(g.Player.Pos) <= 2 {
			// “smart” wait at short distance
			ev.Renew(g, movedelay)
			return
		}
	}
	if mpos.Distance(ppos) == 1 {
		attack := true
		if m.Status(MonsConfused) {
			switch m.Pos.Dir(g.Player.Pos) {
			case E, N, W, S:
			default:
				attack = false
				m.Path = nil
				safepos := m.AlternateConfusedPlacement(g)
				if safepos != nil {
					m.Target = *safepos
				}
			}
		} else if m.Kind == MonsMindCelmist {
			// we try to avoid melee
			safepos := m.SafePlacement(g)
			if safepos != nil {
				attack = false
				m.Path = nil
				m.Target = *safepos
			}
		}
		if attack {
			m.AttackAction(g, ev)
			return
		}
	}
	if m.Status(MonsLignified) {
		ev.Renew(g, 10) // wait
		return
	}
	if m.Kind == MonsMarevorHelith {
		if m.TeleportMonsterAway(g) {
			ev.Renew(g, movedelay)
			return
		}
	}
	m.Obstructing = false
	if !(len(m.Path) > 0 && m.Path[0] == m.Target && m.Path[len(m.Path)-1] == mpos) {
		m.Path = m.APath(g, mpos, m.Target)
		if len(m.Path) == 0 && !m.Status(MonsConfused) {
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
			m.Target = m.NextTarget(g)
			m.GatherBand(g)
		case Hunting:
			if RandInt(4) > 0 {
				m.Dir = m.Dir.Alternate()
			} else {
				// pick a random cell: more escape strategies for the player
				if m.Kind == MonsHound && m.Pos.Distance(g.Player.Pos) <= 6 &&
					!(g.Player.Aptitudes[AptStealthyMovement] && RandInt(2) == 0) {
					m.Target = g.Player.Pos
				} else {
					m.Target = m.NextTarget(g)
				}
				m.State = Wandering
				m.GatherBand(g)
			}
		}
		ev.Renew(g, movedelay)
		return
	}
	target := m.Path[len(m.Path)-2]
	mons := g.MonsterAt(target)
	switch {
	case !mons.Exists():
		if m.Kind == MonsEarthDragon && g.Dungeon.Cell(target).T == WallCell {
			g.Dungeon.SetCell(target, FreeCell)
			g.Stats.Digs++
			if !g.Player.Sees(target) {
				g.WrongWall[m.Pos] = true
			}
			g.MakeNoise(WallNoise, m.Pos)
			g.Fog(m.Pos, 1, ev)
			if g.Player.Pos.Distance(target) < 12 {
				// XXX use dijkstra distance ?
				g.Printf("%s You hear an earth-splitting noise.", g.CrackSound())
				g.StopAuto()
			}
			m.MoveTo(g, target)
			m.Path = m.Path[:len(m.Path)-1]
		} else if g.Dungeon.Cell(target).T == WallCell {
			m.Path = m.APath(g, mpos, m.Target)
		} else {
			m.InvertFoliage(g)
			m.MoveTo(g, target)
			if (m.Kind.Ranged() || m.Kind.Smiting()) && !m.FireReady && m.SeesPlayer(g) {
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
			mons.Target = mons.NextTarget(g)
			mons.State = Wandering
			mons.GatherBand(g)
		} else {
			m.Path = m.APath(g, mpos, m.Target)
		}
	case !mons.SeesPlayer(g) && g.Player.Pos.Distance(mons.Target) > 2 && mons.State != Hunting:
		r := RandInt(5)
		if r == 0 {
			m.Target = m.NextTarget(g)
			m.GatherBand(g)
		} else if (r == 1 || r == 2) && mons.State == Resting {
			mons.Target = m.NextTarget(g)
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
	case mons.State == Hunting && m.State == Hunting || !g.Player.LOS[m.Target]:
		if RandInt(4) == 0 {
			m.Target = mons.Target
			m.Path = m.APath(g, mpos, m.Target)
		} else {
			m.Path = m.APath(g, mpos, m.Target)
		}
	default:
		m.Path = m.APath(g, mpos, m.Target)
	}
	ev.Renew(g, movedelay)
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
	if !g.Player.Sees(m.Pos) && invert {
		g.WrongFoliage[m.Pos] = !g.WrongFoliage[m.Pos]
	} else if invert {
		g.ComputeLOS()
	}
}

func (m *monster) Exhaust(g *game) {
	m.ExhaustTime(g, DurationMonsterExhaustion+RandInt(DurationMonsterExhaustion/2))
}

func (m *monster) ExhaustTime(g *game, t int) {
	m.Statuses[MonsExhausted]++
	g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + t, NMons: m.Index, EAction: MonsExhaustionEnd})
}

func (m *monster) HitPlayer(g *game, ev event) {
	if g.Player.HP <= 0 || g.Player.Pos.Distance(m.Pos) > 1 {
		return
	}
	dmg := m.Attack
	clang := RandInt(4) == 0
	if g.Player.HasStatus(StatusSwap) && !g.Player.HasStatus(StatusLignification) && !m.Status(MonsLignified) {
		g.SwapWithMonster(m)
		return
	}
	noise := g.HitNoise(clang)
	g.MakeNoise(noise, g.Player.Pos)
	var sclang string
	if clang {
		sclang = g.ArmourClang()
	}
	g.PrintfStyled("%s hits you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), dmg, sclang)
	m.InflictDamage(g, dmg, m.Attack)
	if m.Kind == MonsVampire {
		m.HP += 1
		if m.HP > m.HPmax {
			m.HP = m.HPmax
		}
	}
	if g.Player.HP <= 0 {
		return
	}
	m.HitSideEffects(g, ev)
	const HeavyWoundHP = 2
	if g.Player.Aptitudes[AptConfusingGas] && g.Player.HP < HeavyWoundHP {
		m.EnterConfusion(g, ev)
		g.Printf("You release some confusing gas against the %s.", m.Kind)
	}
	if g.Player.Aptitudes[AptSmoke] && g.Player.HP < HeavyWoundHP {
		g.Smoke(ev)
	}
	if g.Player.Aptitudes[AptObstruction] && g.Player.HP <= HeavyWoundHP {
		opos := m.Pos
		m.Blink(g)
		if opos != m.Pos {
			g.TemporalWallAt(opos, ev)
			g.Print("A temporal wall emerges.")
			m.Exhaust(g)
		}
	}
	if g.Player.Aptitudes[AptTeleport] && g.Player.HP < HeavyWoundHP {
		m.TeleportAway(g)
	}
	if g.Player.Aptitudes[AptLignification] && g.Player.HP < HeavyWoundHP {
		m.EnterLignification(g, ev)
	}
}

func (m *monster) EnterConfusion(g *game, ev event) {
	if !m.Status(MonsConfused) {
		m.Statuses[MonsConfused] = 1
		m.Path = m.Path[:0]
		g.PushEvent(&monsterEvent{
			ERank: ev.Rank() + DurationConfusion + RandInt(DurationConfusion/4), NMons: m.Index, EAction: MonsConfusionEnd})
	}
}

func (m *monster) EnterLignification(g *game, ev event) {
	if !m.Status(MonsLignified) {
		m.Statuses[MonsLignified] = 1
		m.Path = m.Path[:0]
		g.PushEvent(&monsterEvent{
			ERank: ev.Rank() + DurationLignification + RandInt(DurationLignification/2), NMons: m.Index, EAction: MonsLignificationEnd})
		if g.Player.Sees(m.Pos) {
			g.Printf("%s is rooted to the ground.", m.Kind.Definite(true))
		}
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
			g.Player.HP += 2
			end := ev.Rank() + DurationShortBerserk
			g.PushEvent(&simpleEvent{ERank: end, EAction: BerserkEnd})
			g.Player.Expire[StatusBerserk] = end
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
	case MonsWingedMilfid:
		if m.Status(MonsExhausted) || g.Player.HasStatus(StatusLignification) {
			break
		}
		ompos := m.Pos
		m.MoveTo(g, g.Player.Pos)
		g.PlacePlayerAt(ompos)
		g.Print("The flying milfid makes you swap positions.")
		m.ExhaustTime(g, 50+RandInt(50))
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
	if !m.SeesPlayer(g) {
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
	case MonsMadNixe:
		return m.NixeAttraction(g, ev)
	case MonsVampire:
		return m.VampireSpit(g, ev)
	case MonsTreeMushroom:
		return m.ThrowSpores(g, ev)
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
	g.MakeNoise(9, m.Pos)
	if RandInt(3) > 0 {
		g.MakeNoise(MagicHitNoise, g.Player.Pos)
		damage := g.Player.HP / 2
		g.PrintfStyled("%s throws a bolt of torment at you.", logMonsterHit, m.Kind.Definite(true))
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '*', ColorCyan)
		m.InflictDamage(g, damage, 1)
	} else {
		g.Printf("You dodge the %s's bolt of torment.", m.Kind)
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '*', ColorCyan)
		// TODO: hit monster behind?
	}
	m.Exhaust(g)
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowRock(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgExtra
	clang := RandInt(3) == 0
	if RandInt(2) == 0 {
		noise := g.HitNoise(clang)
		g.MakeNoise(noise, g.Player.Pos)
		var sclang string
		if clang {
			sclang = g.ArmourClang()
		}
		g.PrintfStyled("%s throws a rock at you (%d dmg).%s", logMonsterHit, m.Kind.Definite(true), dmg, sclang)
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '●', ColorMagenta)
		oppos := g.Player.Pos
		if m.PushPlayer(g) {
			g.TemporalWallAt(oppos, ev)
		} else {
			ray := g.Ray(m.Pos)
			if len(ray) > 0 {
				g.TemporalWallAt(ray[len(ray)-1], ev)
			}
		}
		m.InflictDamage(g, dmg, dmg)
	} else {
		g.Stats.Dodges++
		g.Printf("You dodge %s's rock.", m.Kind.Indefinite(false))
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '●', ColorMagenta)
		dir := g.Player.Pos.Dir(m.Pos)
		pos := g.Player.Pos.To(dir)
		if pos.valid() {
			mons := g.MonsterAt(pos)
			if mons.Exists() {
				mons.HP -= RandInt(15)
				if mons.HP <= 0 {
					g.HandleKill(mons, ev)
				} else {
					mons.Blink(g)
					if mons.Pos != pos {
						g.TemporalWallAt(pos, ev)
					}
				}
			} else {
				g.TemporalWallAt(pos, ev)
			}
		}
	}
	ev.Renew(g, 2*m.Kind.AttackDelay())
	return true
}

func (m *monster) VampireSpit(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusNausea) {
		return false
	}
	g.Player.Statuses[StatusNausea]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + DurationSick, EAction: NauseaEnd})
	g.Print("The vampire spits at you. You feel sick.")
	m.Exhaust(g)
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowSpores(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked || g.Player.HasStatus(StatusLignification) {
		return false
	}
	g.EnterLignification(ev)
	g.Print("The tree mushroom releases spores. You feel rooted to the ground.")
	m.Exhaust(g)
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowJavelin(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgNormal
	clang := RandInt(4) == 0
	if RandInt(2) == 0 {
		noise := g.HitNoise(clang)
		g.MakeNoise(noise, g.Player.Pos)
		var sclang string
		if clang {
			sclang = g.ArmourClang()
		}
		g.Printf("%s throws %s at you (%d dmg).%s", m.Kind.Definite(true), Indefinite("javelin", false), dmg, sclang)
		g.ui.MonsterJavelinAnimation(g.Ray(m.Pos), true)
		m.InflictDamage(g, dmg, dmg)
	} else {
		g.Stats.Dodges++
		g.Printf("You dodge %s's %s.", m.Kind.Indefinite(false), "javelin")
		g.ui.MonsterJavelinAnimation(g.Ray(m.Pos), false)
	}
	m.ExhaustTime(g, 50+RandInt(50))
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) ThrowAcid(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	dmg := DmgNormal
	if RandInt(2) == 0 {
		noise := g.HitNoise(false) // no clang with acid projectiles
		g.MakeNoise(noise, g.Player.Pos)
		g.Printf("%s throws acid at you (%d dmg).", m.Kind.Definite(true), dmg)
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '*', ColorGreen)
		m.InflictDamage(g, dmg, dmg)
		if RandInt(2) == 0 {
			g.Corrosion(ev)
			if RandInt(2) == 0 {
				g.Confusion(ev)
			}
		}
	} else {
		g.Stats.Dodges++
		g.Printf("You dodge %s's acid projectile.", m.Kind.Indefinite(false))
		g.ui.MonsterProjectileAnimation(g.Ray(m.Pos), '*', ColorGreen)
	}
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) NixeAttraction(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	g.MakeNoise(9, m.Pos)
	g.PrintfStyled("%s lures you to her.", logMonsterHit, m.Kind.Definite(true))
	ray := g.Ray(m.Pos)
	g.ui.MonsterProjectileAnimation(ray, 'θ', ColorCyan) // TODO: improve
	if len(ray) > 1 {
		// should always be the case
		g.ui.TeleportAnimation(g.Player.Pos, ray[1], true)
		g.PlacePlayerAt(ray[1])
	}
	m.Exhaust(g)
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) SmitingAttack(g *game, ev event) bool {
	if !m.Kind.Smiting() {
		return false
	}
	if !m.SeesPlayer(g) {
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
	case MonsMindCelmist:
		return m.MindAttack(g, ev)
	}
	return false
}

func (m *monster) AbsorbMana(g *game, ev event) bool {
	if g.Player.MP == 0 {
		return false
	}
	g.Player.MP -= 1
	g.Printf("%s absorbs your mana.", m.Kind.Definite(true))
	m.ExhaustTime(g, 10+RandInt(10))
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) MindAttack(g *game, ev event) bool {
	if g.Player.Pos.Distance(m.Pos) == 1 {
		return false
	}
	g.Print("The celmist mage attacks your mind.")
	if RandInt(2) == 0 {
		g.Player.Statuses[StatusSlow]++
		g.PushEvent(&simpleEvent{ERank: ev.Rank() + DurationSleepSlow, EAction: SlowEnd})
		g.Print("You feel slow.")
	} else {
		g.Confusion(ev)
	}
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) Explode(g *game, ev event) {
	neighbors := m.Pos.ValidNeighbors()
	g.MakeNoise(WallNoise, m.Pos)
	g.Printf("%s %s explodes with a loud boom.", g.ExplosionSound(), m.Kind.Definite(true))
	g.ui.ExplosionAnimation(FireExplosion, m.Pos)
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
			g.HandleStone(mons)
			mons.MakeHuntIfHurt(g)
		} else if g.Player.Pos == pos {
			dmg := g.Player.HP / 2
			m.InflictDamage(g, dmg, 2)
		} else if c.T == WallCell && RandInt(2) == 0 {
			g.Dungeon.SetCell(pos, FreeCell)
			g.Stats.Digs++
			if !g.Player.Sees(pos) {
				g.WrongWall[pos] = true
			} else {
				g.ui.WallExplosionAnimation(pos)
			}
			g.MakeNoise(WallNoise, pos)
			g.Fog(pos, 1, ev)
		}
	}
}

func (m *monster) Blink(g *game) {
	npos := g.BlinkPos()
	if !npos.valid() || npos == g.Player.Pos || npos == m.Pos {
		return
	}
	opos := m.Pos
	g.Printf("The %s blinks away.", m.Kind)
	g.ui.TeleportAnimation(opos, npos, true)
	m.MoveTo(g, npos)
}

func (m *monster) MakeHunt(g *game) {
	m.State = Hunting
	m.Target = g.Player.Pos
}

func (m *monster) MakeHuntIfHurt(g *game) {
	if m.Exists() && m.State != Hunting {
		m.MakeHunt(g)
		if m.State == Resting {
			g.Printf("%s awakens.", m.Kind.Definite(true))
		}
		if m.Kind == MonsHound {
			g.Printf("%s barks.", m.Kind.Definite(true))
			g.MakeNoise(BarkNoise, m.Pos)
		}
	}
}

func (m *monster) MakeAwareIfHurt(g *game) {
	if m.SeesPlayer(g) && m.State != Hunting {
		m.MakeHuntIfHurt(g)
		return
	}
	if m.State != Resting {
		return
	}
	m.State = Wandering
	m.Target = m.NextTarget(g)
}

func (m *monster) MakeAware(g *game) {
	if !m.SeesPlayer(g) {
		return
	}
	if m.State == Resting {
		if m.Status(MonsExhausted) && (m.Pos.Distance(g.Player.Pos) > 1 || RandInt(3) > 0) {
			return
		}
		adjust := g.LosRange() - m.Pos.Distance(g.Player.Pos)
		max := 28
		if g.Player.Aptitudes[AptStealthyMovement] {
			max += 3
		}
		if g.Player.Armour == HarmonistRobe {
			max += 10
		}
		stealth := max - 4*adjust
		fact := 2
		if m.Pos.Distance(g.Player.Pos) > 1 {
			fact = 3
		} else if stealth > 15 {
			stealth = 15
		}
		r := RandInt(stealth)
		if g.Player.Aptitudes[AptStealthyMovement] {
			r *= fact
		}
		if g.Player.Armour == HarmonistRobe {
			r *= fact
		}
		if r >= 5 {
			return
		}
	}
	if m.State == Wandering {
		adjust := g.LosRange() - m.Pos.Distance(g.Player.Pos)
		max := 37
		if g.Player.Aptitudes[AptStealthyMovement] {
			max += 5
		}
		if g.Player.Armour == HarmonistRobe {
			max += 10
		}
		stealth := max - 4*adjust
		r := RandInt(stealth)
		if g.Player.Aptitudes[AptStealthyMovement] {
			r *= 2
		}
		if g.Player.Armour == HarmonistRobe {
			r *= 2
			r += 5
		}
		if r >= 25 && m.Pos.Distance(g.Player.Pos) > 1 {
			return
		}
	}
	if m.State == Resting {
		g.Printf("%s awakens.", m.Kind.Definite(true))
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
	if !g.BandData[g.Bands[m.Band].Kind].Band {
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
			if !ok || n.Cost > 4 || mons.State == Resting && mons.Status(MonsExhausted) {
				continue
			}
			mons.Target = m.Target
			if mons.State == Resting {
				mons.State = Wandering
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
	danger := [MaxDepth + 1]int{
		1:  20,
		2:  42,
		3:  65,
		4:  90,
		5:  115,
		6:  140,
		7:  165,
		8:  190,
		9:  215,
		10: 245,
		11: 285,
	}
	max := danger[g.Depth]
	adjust := -2 * g.Depth
	for c, q := range g.Player.Consumables {
		switch c {
		case HealWoundsPotion, CBlinkPotion:
			adjust += Min(5, g.Depth) * Min(q, Min(5, g.Depth))
		case TeleportationPotion, DigPotion, WallPotion:
			adjust += Min(3, g.Depth) * Min(q, 3)
		case SwiftnessPotion, LignificationPotion, MagicPotion, BerserkPotion, ExplosiveMagara, ShadowsPotion, TormentPotion, TeleportMagara, NightMagara:
			adjust += Min(2, g.Depth) * Min(q, 3)
		case ConfusingDart:
			adjust += Min(1, g.Depth) * Min(q, 7)
		}
	}
	for _, props := range g.Player.Rods {
		adjust += Min(props.Charge, 2) * Min(2, g.Depth-1)
	}
	if g.Depth < MaxDepth && g.Player.Consumables[DescentPotion] > 0 {
		adjust += g.Depth
	}
	if max+adjust < max-max/3 {
		max = max - max/3
	} else if max+adjust > max+max/3 {
		max = max + max/3
	} else {
		max = max + adjust
	}
	if g.Depth > 3 && g.Player.Weapon == Dagger {
		max -= 3 * g.Depth
	}
	if g.Depth > 4 && g.Player.Armour == Robe {
		max -= 2 * g.Depth
	}
	if g.Player.Consumables[MagicMappingPotion] > 0 && WinDepth-g.Depth < g.Player.Consumables[MagicMappingPotion] {
		max = max * 110 / 100
	}
	if g.Player.Consumables[DreamPotion] > 0 && WinDepth-g.Depth < g.Player.Consumables[DreamPotion] {
		max = max * 105 / 100
	}
	//switch g.Dungeon.Gen {
	//case GenCaveMapTree:
	//max = max * 90 / 100
	//case GenCaveMap:
	//max = max * 95 / 100
	//case GenRoomMap:
	//max = max * 105 / 100
	//case GenRuinsMap:
	//max = max * 108 / 100
	//case GenBSPMap:
	//max = max * 115 / 100
	//}
	return max
}

func (g *game) MaxMonsters() int {
	nmons := [MaxDepth + 1]int{
		1:  11,
		2:  15,
		3:  20,
		4:  26,
		5:  31,
		6:  31,
		7:  31,
		8:  33,
		9:  33,
		10: 36,
		11: 39,
	}
	max := nmons[g.Depth]
	//switch g.Dungeon.Gen {
	//case GenCaveMapTree, GenCaveMap:
	//max = max * 90 / 100
	//case GenBSPMap:
	//max = max * 110 / 100
	//}
	return max
}

func (g *game) GenMonsters() {
	g.Monsters = []*monster{}
	g.Bands = []bandInfo{}
	danger := g.MaxDanger()
	nmons := g.MaxMonsters()
	nband := 0
	i := 0
	repeat := 0
loop:
	for danger > 0 && nmons > 0 {
		for band, data := range g.BandData {
			if RandInt(data.Rarity*50) != 0 {
				continue
			}
			monsters := g.GenBand(data, monsterBand(band))
			if monsters == nil {
				continue
			}
			if data.Unique {
				g.GeneratedUniques[monsterBand(band)]++
			}
			bandinfo := bandInfo{Kind: monsterBand(band)}
			pos := g.FreeCellForMonster()
			bandinfo.Path = append(bandinfo.Path, pos)
			bandinfo.Path = append(bandinfo.Path, g.FreeCellForMonster())
			g.Bands = append(g.Bands, bandinfo)
			for _, mk := range monsters {
				if mk == MonsGoblin {
					mk = g.Opts.Alternate
				}
				if nmons-1 <= 0 {
					return
				}
				if danger-mk.Dangerousness() <= 0 {
					if repeat > 15 {
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
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			return mons
		}
	}
	return nil
}
