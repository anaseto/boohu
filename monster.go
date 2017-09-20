package main

import "container/heap"

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

func (st monsterStatus) String() (text string) {
	switch st {
	case MonsConfused:
		text = "confused"
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

func (mk monsterKind) Ranged() bool {
	switch mk {
	case MonsLich, MonsCyclop:
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
	MonsOgre:            {10, 15, 12, 28, 13, 0, 8, 'O', "ogre", 5},
	MonsCyclop:          {10, 12, 12, 28, 13, 0, 8, 'C', "cyclop", 9},
	MonsWorm:            {12, 9, 10, 25, 13, 0, 10, 'w', "worm", 3},
	MonsHound:           {8, 9, 10, 15, 14, 0, 12, 'h', "hound", 4},
	MonsYack:            {10, 11, 10, 21, 14, 0, 10, 'y', "yack", 5},
	MonsGiantBee:        {6, 10, 10, 11, 15, 0, 15, 'B', "giant bee", 6},
	MonsGoblinWarrior:   {10, 11, 10, 25, 15, 3, 12, 'G', "goblin warrior", 8},
	MonsHydra:           {10, 9, 10, 45, 13, 0, 6, 'H', "hydra", 15},
	MonsSkeletonWarrior: {10, 12, 10, 25, 15, 4, 12, 'S', "skeleton warrior", 10},
	MonsSpider:          {8, 7, 10, 13, 17, 0, 15, 's', "spider", 6},
	MonsBlinkingFrog:    {10, 10, 10, 20, 15, 0, 12, 'F', "blinking frog", 7},
	MonsLich:            {10, 10, 10, 23, 15, 3, 12, 'L', "lich", 17},
	MonsEarthDragon:     {10, 14, 10, 40, 14, 6, 8, 'D', "earth dragon", 20},
	MonsMirrorSpecter:   {10, 9, 10, 18, 15, 0, 17, 'm', "mirror specter", 11},
}

var monsDesc = []string{
	MonsGoblin:          "Goblins are little humanoid creatures. They often appear in group.",
	MonsOgre:            "Ogres are big clunky humanoids that can hit really hard.",
	MonsCyclop:          "Cyclops are very similar to ogres, but they also like to throw rocks at their foes, sometimes confusing them.",
	MonsWorm:            "Worms are ugly slow moving creatures, but surprisingly hardy at times.",
	MonsHound:           "Hounds are fast moving carnivore quadrupeds. They sometimes attack in group.",
	MonsYack:            "Yacks are quite large herbivorous quadrupeds. They tend to form large groups.",
	MonsGiantBee:        "Giant bees are fragile, but extremely fast moving creatures. Their bite can sometimes enrage you.",
	MonsGoblinWarrior:   "Goblin warriors are goblins that learned to fight, and got equipped with a leather armour.",
	MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	MonsSkeletonWarrior: "Skeleton warriors are good fighters, and are equipped with a chain mail.",
	MonsSpider:          "Spiders are fast moving fragile creatures, whose bite can confuse you.",
	MonsBlinkingFrog:    "Blinking frogs are big frog-like unstable creatures, whose bite can make you blink away.",
	MonsLich:            "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you.",
	MonsEarthDragon:     "Earth dragons are big and hardy creatures that wander in the Underground. It is said they are to credit for many tunnels.",
	MonsMirrorSpecter:   "Mirror specters are very insubstancial creatures. They can absorb your mana.",
}

type monsterBand int

const (
	LoneGoblin monsterBand = iota
	LoneOgre
	LoneWorm
	LoneHound
	LoneHydra
	LoneSpider
	LoneBlinkingFrog
	LoneCyclop
	LoneLich
	LoneEarthDragon
	LoneSpecter
	BandGoblins
	BandGoblinsWithWarriors
	BandGoblinWarriors
	BandHounds
	BandYacks
	BandSpiders
	BandBlinkingFrogs
	BandGiantBees
	BandSkeletonWarrior
	UBandWorms
	UBandGoblinsEasy
	UBandFrogs
	UBandOgres
	UBandGoblins
	UBandBeeYacks
	UHydras
	ULich
	UDragon
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
	LoneGoblin:       {rarity: 10, minDepth: 0, maxDepth: 5, monster: MonsGoblin},
	LoneOgre:         {rarity: 15, minDepth: 2, maxDepth: 11, monster: MonsOgre},
	LoneWorm:         {rarity: 10, minDepth: 0, maxDepth: 6, monster: MonsWorm},
	LoneHound:        {rarity: 20, minDepth: 1, maxDepth: 8, monster: MonsHound},
	LoneHydra:        {rarity: 45, minDepth: 8, maxDepth: 13, monster: MonsHydra},
	LoneSpider:       {rarity: 20, minDepth: 3, maxDepth: 13, monster: MonsSpider},
	LoneBlinkingFrog: {rarity: 50, minDepth: 5, maxDepth: 13, monster: MonsBlinkingFrog},
	LoneCyclop:       {rarity: 45, minDepth: 5, maxDepth: 13, monster: MonsCyclop},
	LoneLich:         {rarity: 70, minDepth: 9, maxDepth: 13, monster: MonsLich},
	LoneEarthDragon:  {rarity: 80, minDepth: 10, maxDepth: 13, monster: MonsEarthDragon},
	LoneSpecter:      {rarity: 70, minDepth: 6, maxDepth: 13, monster: MonsMirrorSpecter},
	BandGoblins: {
		distribution: map[monsterKind]monsInterval{MonsGoblin: monsInterval{2, 4}},
		rarity:       10, minDepth: 1, maxDepth: 7, band: true,
	},
	BandGoblinsWithWarriors: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        monsInterval{3, 5},
			MonsGoblinWarrior: monsInterval{0, 2}},
		rarity: 10, minDepth: 5, maxDepth: 9, band: true,
	},
	BandGoblinWarriors: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        monsInterval{0, 1},
			MonsGoblinWarrior: monsInterval{2, 4}},
		rarity: 45, minDepth: 10, maxDepth: 13, band: true,
	},
	BandHounds: {
		distribution: map[monsterKind]monsInterval{MonsHound: monsInterval{2, 3}},
		rarity:       20, minDepth: 2, maxDepth: 10, band: true,
	},
	BandSpiders: {
		distribution: map[monsterKind]monsInterval{MonsSpider: monsInterval{2, 4}},
		rarity:       25, minDepth: 5, maxDepth: 13, band: true,
	},
	BandBlinkingFrogs: {
		distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: monsInterval{2, 4}},
		rarity:       70, minDepth: 9, maxDepth: 13, band: true,
	},
	BandYacks: {
		distribution: map[monsterKind]monsInterval{MonsYack: monsInterval{2, 5}},
		rarity:       15, minDepth: 5, maxDepth: 11, band: true,
	},
	BandGiantBees: {
		distribution: map[monsterKind]monsInterval{MonsGiantBee: monsInterval{2, 5}},
		rarity:       30, minDepth: 6, maxDepth: 13, band: true,
	},
	BandSkeletonWarrior: {
		distribution: map[monsterKind]monsInterval{MonsSkeletonWarrior: monsInterval{2, 3}},
		rarity:       45, minDepth: 8, maxDepth: 13, band: true,
	},
	UBandWorms: {
		distribution: map[monsterKind]monsInterval{MonsWorm: monsInterval{3, 4}, MonsSpider: monsInterval{1, 1}},
		rarity:       50, minDepth: 4, maxDepth: 4, band: true, unique: true,
	},
	UBandGoblinsEasy: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin: monsInterval{3, 5},
			MonsHound:  monsInterval{1, 2},
		},
		rarity: 30, minDepth: 5, maxDepth: 5, band: true, unique: true,
	},
	UBandFrogs: {
		distribution: map[monsterKind]monsInterval{MonsBlinkingFrog: monsInterval{2, 3}},
		rarity:       60, minDepth: 6, maxDepth: 6, band: true, unique: true,
	},
	UBandOgres: {
		distribution: map[monsterKind]monsInterval{MonsOgre: monsInterval{2, 3}, MonsCyclop: monsInterval{1, 1}},
		rarity:       35, minDepth: 7, maxDepth: 7, band: true, unique: true,
	},
	UBandGoblins: {
		distribution: map[monsterKind]monsInterval{
			MonsGoblin:        monsInterval{3, 5},
			MonsGoblinWarrior: monsInterval{1, 2},
			MonsHound:         monsInterval{1, 2},
		},
		rarity: 30, minDepth: 8, maxDepth: 8, band: true, unique: true,
	},
	UBandBeeYacks: {
		distribution: map[monsterKind]monsInterval{
			MonsYack:     monsInterval{2, 5},
			MonsGiantBee: monsInterval{1, 3},
		},
		rarity: 30, minDepth: 9, maxDepth: 9, band: true, unique: true,
	},
	UHydras: {
		distribution: map[monsterKind]monsInterval{
			MonsHydra:  monsInterval{2, 3},
			MonsSpider: monsInterval{1, 2},
		},
		rarity: 35, minDepth: 10, maxDepth: 10, band: true, unique: true,
	},
	ULich: {
		distribution: map[monsterKind]monsInterval{
			MonsSkeletonWarrior: monsInterval{1, 3},
			MonsLich:            monsInterval{1, 1},
			MonsMirrorSpecter:   monsInterval{0, 1},
		},
		rarity: 50, minDepth: 11, maxDepth: 11, band: true, unique: true,
	},
	UDragon: {
		distribution: map[monsterKind]monsInterval{
			MonsEarthDragon: monsInterval{2, 2},
		},
		rarity: 60, minDepth: 12, maxDepth: 12, band: true, unique: true,
	},
}

type monster struct {
	Kind        monsterKind
	Band        int
	Attack      int
	Accuracy    int
	Armor       int
	Evasion     int
	HPmax       int
	HP          int
	Pos         position
	State       monsterState
	Statuses    map[monsterStatus]int
	Target      position
	Path        []position // cache
	Obstructing bool
}

func (m *monster) Init() {
	m.HPmax = MonsData[m.Kind].maxHP - 1 + RandInt(3)
	m.Attack = MonsData[m.Kind].baseAttack - RandInt(2)
	m.HP = m.HPmax
	m.Accuracy = MonsData[m.Kind].accuracy
	m.Armor = MonsData[m.Kind].armor
	m.Evasion = MonsData[m.Kind].evasion
	m.Statuses = map[monsterStatus]int{}
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
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return &pos
	}
	return nil
}

func (m *monster) AttackAction(g *game, ev event) {
	switch {
	case m.Obstructing:
		m.Obstructing = false
		pos := m.AlternatePlacement(g)
		if pos != nil {
			m.Pos = *pos
			ev.Renew(g, m.Kind.MovementDelay())
			return
		}
		fallthrough
	default:
		if m.Kind == MonsHydra {
			for i := 0; i <= 3; i++ {
				m.HitPlayer(g, ev)
			}
		} else {
			m.HitPlayer(g, ev)
		}
		ev.Renew(g, m.Kind.AttackDelay())
	}
}

func (m *monster) HandleTurn(g *game, ev event) {
	ppos := g.Player.Pos
	mpos := m.Pos
	m.MakeAware(g)
	if m.State == Resting {
		wander := RandInt(1500)
		if wander == 0 {
			m.Target = g.FreeCell()
			m.State = Wandering
			m.GatherBand(g)
		}
		ev.Renew(g, m.Kind.MovementDelay())
		return
	}
	if m.RangedAttack(g, ev) {
		return
	}
	if m.SmitingAttack(g, ev) {
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
	m.Obstructing = false
	if !(len(m.Path) > 0 && m.Path[0] == m.Target && m.Path[len(m.Path)-1] == mpos) {
		m.Path = m.APath(g, mpos, m.Target)
	}
	if m.Path == nil || len(m.Path) < 2 {
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
			if RandInt(5) == 0 && m.Pos.Distance(g.Player.Pos) < 10 {
				// make hunting monsters sometimes smart
				m.Target = g.Player.Pos
			} else {
				m.Target = g.FreeCell()
			}
			m.State = Wandering
			m.GatherBand(g)
		}
		ev.Renew(g, m.Kind.MovementDelay())
		return
	}
	target := m.Path[len(m.Path)-2]
	mons, _ := g.MonsterAt(target)
	switch {
	case !mons.Exists():
		m.Pos = m.Path[len(m.Path)-2]
		if m.Kind == MonsEarthDragon && g.Dungeon.Cell(m.Pos).T == WallCell {
			g.Dungeon.SetCell(m.Pos, FreeCell)
			if !g.Player.LOS[m.Pos] {
				g.UnknownDig[m.Pos] = true
			}
			g.MakeNoise(18, m.Pos)
			if g.Player.Pos.Distance(m.Pos) < 10 {
				// XXX use dijkstra distance ?
				g.Print("You hear an earth-breaking noise.")
				g.AutoHalt = true
			}
		}
		m.Path = m.Path[:len(m.Path)-1]
	case !g.Player.LOS[mons.Pos] && g.Player.Pos.Distance(mons.Target) > 2 && mons.State != Hunting:
		r := RandInt(10)
		if r == 0 {
			m.Target = g.FreeCell()
			m.GatherBand(g)
		} else if (r == 1 || r == 2) && mons.State == Resting {
			mons.Target = g.FreeCell()
			mons.State = Wandering
			mons.GatherBand(g)
		}
	case mons.Pos.Distance(g.Player.Pos) == 1:
		m.Path = m.APath(g, mpos, m.Target)
		if len(m.Path) < 2 || m.Path[len(m.Path)-2] == mons.Pos {
			mons.Obstructing = true
		}
	default:
		m.Path = m.APath(g, mpos, m.Target)
	}
	ev.Renew(g, m.Kind.MovementDelay())
}

func (m *monster) HitPlayer(g *game, ev event) {
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	if acc > evasion {
		if m.Blocked(g) {
			g.Printf("You block the %s's attack with your %s.", m.Kind, g.Player.Shield)
			return
		}
		noise := 12
		noise += g.Player.Armor() / 2
		g.MakeNoise(noise, g.Player.Pos)
		attack := g.HitDamage(m.Attack, g.Player.Armor())
		g.Player.HP -= attack
		g.Printf("The %s hits you (%d damage).", m.Kind, attack)
		m.HitSideEffects(g, ev)
	} else {
		g.Printf("The %s misses you.", m.Kind)
	}
}

func (m *monster) HitSideEffects(g *game, ev event) {
	switch m.Kind {
	case MonsSpider:
		if RandInt(2) == 0 && !g.Player.HasStatus(StatusConfusion) {
			g.Player.Statuses[StatusConfusion]++
			heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: ConfusionEnd})
			g.Print("You feel confused.")
		}
	case MonsGiantBee:
		if RandInt(5) == 0 && !g.Player.HasStatus(StatusBerserk) {
			g.Player.Statuses[StatusBerserk]++
			heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 25 + RandInt(40), EAction: BerserkEnd})
			g.Print("You feel a sudden urge to kill things.")
		}
	case MonsBlinkingFrog:
		if RandInt(2) == 0 {
			g.Blink(ev)
		}
		//if RandInt(3) == 0 && !g.Player.HasStatus(StatusNausea) {
		//g.Player.Statuses[StatusNausea]++
		//heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 20 + RandInt(30), EAction: NauseaEnd})
		//g.Print("You feel sick.")
		//}
	}

}

func (m *monster) RangedAttack(g *game, ev event) bool {
	if !m.Kind.Ranged() {
		return false
	}
	rdist := 5
	if g.Player.Aptitudes[AptStealthyLOS] {
		rdist = 4
	}
	if m.Pos.Distance(g.Player.Pos) <= 1 || m.Pos.Distance(g.Player.Pos) > rdist || !g.Player.LOS[m.Pos] {
		return false
	}
	if m.Status(MonsExhausted) {
		return false
	}
	switch m.Kind {
	case MonsLich:
		return m.TormentBolt(g, ev)
	case MonsCyclop:
		return m.ThrowRock(g, ev)
	case MonsMirrorSpecter:
	}
	return false
}

func (m *monster) RangeBlocked(g *game) bool {
	ray := g.Ray(m.Pos)
	blocked := false
	for _, pos := range ray[1:] {
		mons, _ := g.MonsterAt(pos)
		if mons == nil {
			continue
		}
		blocked = true
		break
	}
	return blocked
}

func (m *monster) Index(g *game) int {
	for i, mons := range g.Monsters {
		if mons.Pos == m.Pos {
			return i
		}
	}
	// not reached
	return 0
}

func (m *monster) TormentBolt(g *game, ev event) bool {
	blocked := m.RangeBlocked(g)
	if blocked {
		return false
	}
	//g.Player.Statuses[StatusSlow]++
	//heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 50 + RandInt(50), EAction: SlowEnd})
	hit := !m.Blocked(g)
	g.MakeNoise(9, m.Pos)
	if hit {
		g.MakeNoise(12, g.Player.Pos)
		g.Player.HP = g.Player.HP / 2
		g.Printf("The %s throws a bolt of torment at you.", m.Kind)
	} else {
		g.Printf("You block the %s's bolt of torment.", m.Kind)
	}
	m.Statuses[MonsExhausted]++
	heap.Push(g.Events, &monsterEvent{ERank: ev.Rank() + 100 + RandInt(50), NMons: m.Index(g), EAction: MonsExhaustionEnd})
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (g *game) HitDamage(base int, armor int) int {
	min := base / 2
	attack := min + RandInt(base-min+1)
	attack -= RandInt(armor + 1)
	if attack < 0 {
		attack = 0
	}
	return attack
}

func (m *monster) Blocked(g *game) bool {
	blocked := false
	if g.Player.Shield != NoShield && !g.Player.Weapon.TwoHanded() {
		block := RandInt(g.Player.Shield.Block())
		acc := RandInt(m.Accuracy)
		if block >= acc {
			g.MakeNoise(12+g.Player.Shield.Block()/2, g.Player.Pos)
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
	hit := true
	evasion := RandInt(g.Player.Evasion())
	acc := RandInt(m.Accuracy)
	if 3*acc/2 <= evasion {
		// rocks are big and do not miss so often
		hit = false
	} else {
		hit = !m.Blocked(g)
	}
	if hit {
		noise := 12
		noise += g.Player.Armor() / 2
		g.MakeNoise(noise, g.Player.Pos)
		attack := g.HitDamage(15, g.Player.Armor())
		g.Player.HP -= attack
		g.Printf("The %s throws a rock at you (%d damage).", m.Kind, attack)
		if RandInt(4) == 0 {
			g.Player.Statuses[StatusConfusion]++
			heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 100 + RandInt(100), EAction: ConfusionEnd})
			g.Print("You feel confused.")
		}
	} else {
		g.Printf("You block the %s's rock.", m.Kind)
	}
	ev.Renew(g, 2*m.Kind.AttackDelay())
	return true
}

func (m *monster) SmitingAttack(g *game, ev event) bool {
	if !m.Kind.Smiting() {
		return false
	}
	rdist := 5
	if g.Player.Aptitudes[AptStealthyLOS] {
		rdist = 4
	}
	if m.Pos.Distance(g.Player.Pos) > rdist || !g.Player.LOS[m.Pos] {
		return false
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
	g.Player.MP = 2 * g.Player.MP / 3
	g.Printf("The %s absorbs your mana.", m.Kind)
	m.Statuses[MonsExhausted]++
	heap.Push(g.Events, &monsterEvent{ERank: ev.Rank() + 10 + RandInt(20), NMons: m.Index(g), EAction: MonsExhaustionEnd})
	ev.Renew(g, m.Kind.AttackDelay())
	return true
}

func (m *monster) MakeHuntIfHurt(g *game) {
	if m.State != Hunting {
		m.Target = g.Player.Pos
		m.State = Hunting
		if m.State == Resting {
			g.Printf("The %s awakes.", m.Kind)
		}
		if m.Kind == MonsHound {
			g.Printf("The %s barks.", m.Kind)
			g.MakeNoise(12, m.Pos)
		}
	}
}

func (m *monster) MakeAware(g *game) {
	if !g.Player.LOS[m.Pos] {
		return
	}
	if m.State == Resting {
		adjust := (m.Pos.Distance(g.Player.Pos) - g.LosRange()/2 + 1)
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
		adjust *= adjust
		r := RandInt(30 + adjust)
		if r >= 25 {
			return
		}
	}
	if m.State == Resting {
		g.Printf("The %s awakes.", m.Kind)
	}
	if m.State == Wandering {
		g.Printf("The %s notices you.", m.Kind)
	}
	if m.State != Hunting && m.Kind == MonsHound {
		g.Printf("The %s barks.", m.Kind)
		g.MakeNoise(12, m.Pos)
	}
	m.Target = g.Player.Pos
	m.State = Hunting
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
			n, ok := nm[mons.Pos]
			if !ok || n.Cost > 4 {
				continue
			}
			r := RandInt(100)
			if r > 60 || mons.State == Wandering && r > 10 {
				mons.Target = m.Target
				mons.State = m.State
			}

		}
	}
}

func (g *game) MakeMonstersAware() {
	for _, m := range g.Monsters {
		if m.HP <= 0 {
			continue
		}
		if g.Player.LOS[m.Pos] {
			m.MakeAware(g)
			if m.State != Resting {
				m.GatherBand(g)
			}
		}
	}
}

func (g *game) MakeNoise(noise int, at position) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, noise)
	for _, m := range g.Monsters {
		if !m.Exists() {
			continue
		}
		if m.State == Hunting {
			continue
		}
		n, ok := nm[m.Pos]
		if !ok {
			continue
		}
		d := n.Cost
		v := noise - d
		if v <= 0 {
			continue
		}
		v *= 3
		if v > 90 {
			v = 90
		}
		r := RandInt(100)
		if m.State == Resting {
			r += 10
		}
		if v > r {
			m.Target = at
			if g.Player.LOS[m.Pos] {
				m.State = Hunting
			} else {
				m.State = Wandering
			}
			m.GatherBand(g)
		}
	}
}

func (g *game) MonsterAt(pos position) (*monster, int) {
	var mons *monster
	var index int
	for i, m := range g.Monsters {
		if m.Pos == pos && m.HP > 0 {
			mons = m
			index = i
			break
		}
	}
	return mons, index
}

func (g *game) GenMonsters() {
	g.Monsters = []*monster{}
	g.Bands = []monsterBand{}
	danger := 20 + 10*g.Depth + g.Depth*g.Depth/3
	nmons := 15 + 3*g.Depth
	nmons += RandInt(3)
	if nmons > 40 {
		nmons = 40 + RandInt(5)
	}
	nband := 0
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
				danger -= MonsData[mk].dangerousness
				nmons--
				if danger <= 0 || nmons <= 0 {
					return
				}
				mons := &monster{Kind: mk}
				mons.Init()
				mons.Pos = pos
				mons.Band = nband
				g.Monsters = append(g.Monsters, mons)
				pos = g.FreeCellForBandMonster(pos)
			}
			nband++
		}
	}
}
