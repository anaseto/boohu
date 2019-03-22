package main

import "fmt"

type monsterState int

const (
	Resting monsterState = iota
	Hunting
	Wandering
	Watching
	Waiting
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
	case Watching:
		st = "watching"
	case Waiting:
		st = "waiting"
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
	MonsGuard monsterKind = iota
	MonsYack
	MonsSatowalgaPlant
	MonsMadNixe
	MonsBlinkingFrog
	MonsWorm
	MonsMirrorSpecter
	//MonsTinyHarpy
	//MonsOgre
	MonsCyclop
	//MonsBrizzia
	MonsHound
	//MonsGiantBee
	//MonsGoblinWarrior
	//MonsHydra
	//MonsSkeletonWarrior
	//MonsSpider
	MonsWingedMilfid
	MonsLich
	MonsEarthDragon
	//MonsAcidMound
	//MonsExplosiveNadre
	//MonsMindCelmist
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
	//case MonsLich, MonsCyclop, MonsGoblinWarrior, MonsSatowalgaPlant, MonsMadNixe, MonsVampire, MonsTreeMushroom:
	case MonsLich, MonsCyclop, MonsSatowalgaPlant, MonsMadNixe, MonsVampire, MonsTreeMushroom:
		return true
	default:
		return false
	}
}

func (mk monsterKind) Smiting() bool {
	switch mk {
	//case MonsMirrorSpecter, MonsMindCelmist:
	case MonsMirrorSpecter:
		return true
	default:
		return false
	}
}

func (mk monsterKind) CanOpenDoors() bool {
	switch mk {
	case MonsGuard, MonsMadNixe, MonsCyclop, MonsLich, MonsVampire, MonsWingedMilfid, MonsMarevorHelith:
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
	// TODO: useless
	switch mk {
	//case MonsLich, MonsSkeletonWarrior, MonsMarevorHelith:
	//return false
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
	MonsGuard: {10, 1, 10, 2, 'g', "goblin", 3},
	//MonsTinyHarpy:       {10, 1, 10, 2, 't', "tiny harpy", 4},
	//MonsOgre:            {10, 2, 20, 3, 'O', "ogre", 7},
	MonsCyclop: {10, 2, 20, 3, 'C', "cyclops", 9},
	MonsWorm:   {15, 1, 10, 3, 'w', "farmer worm", 4},
	//MonsBrizzia:         {15, 1, 10, 3, 'z', "brizzia", 6},
	//MonsAcidMound:       {10, 1, 10, 2, 'a', "acid mound", 6},
	MonsHound: {10, 1, 10, 2, 'h', "hound", 5},
	MonsYack:  {10, 1, 10, 2, 'y', "yack", 5},
	//MonsGiantBee:        {5, 1, 10, 1, 'B', "giant bee", 6},
	//MonsGoblinWarrior:   {10, 1, 10, 2, 'G', "goblin warrior", 5},
	//MonsHydra:           {10, 1, 10, 4, 'H', "hydra", 10},
	//MonsSkeletonWarrior: {10, 1, 10, 3, 'S', "skeleton warrior", 6},
	//MonsSpider:          {10, 1, 10, 2, 's', "spider", 6},
	MonsWingedMilfid:  {10, 1, 10, 2, 'W', "winged milfid", 6},
	MonsBlinkingFrog:  {10, 1, 10, 2, 'F', "blinking frog", 6},
	MonsLich:          {10, 1, 10, 2, 'L', "lich", 15},
	MonsEarthDragon:   {10, 2, 10, 4, 'D', "earth dragon", 20},
	MonsMirrorSpecter: {10, 1, 10, 2, 'm', "mirror specter", 11},
	//MonsExplosiveNadre:  {10, 1, 10, 1, 'n', "explosive nadre", 6},
	MonsSatowalgaPlant: {10, 1, 10, 3, 'P', "satowalga plant", 7},
	MonsMadNixe:        {10, 1, 10, 2, 'N', "mad nixe", 14},
	//MonsMindCelmist:     {10, 1, 20, 2, 'c', "mind celmist", 12},
	MonsVampire:       {10, 1, 10, 2, 'V', "vampire", 13},
	MonsTreeMushroom:  {20, 1, 20, 4, 'T', "tree mushroom", 16},
	MonsMarevorHelith: {10, 0, 10, 10, 'M', "Marevor Helith", 18},
}

var monsDesc = []string{
	MonsGuard: "Goblins are little humanoid creatures. They often appear in a group.",
	//MonsTinyHarpy:       "Tiny harpies are little humanoid flying creatures. They blink away when hurt. They often appear in a group.",
	//MonsOgre:            "Ogres are big clunky humanoids that can hit really hard.",
	MonsCyclop: "Cyclopes are very similar to ogres, but they also like to throw rocks at their foes (2 damage). The rocks can block your way for a while.",
	MonsWorm:   "Farmer worms are ugly slow moving creatures, but surprisingly hardy at times, and they furrow as they move, helping new foliage to grow.",
	//MonsBrizzia:         "Brizzias are big slow moving biped creatures. They are quite hardy, and when hurt they can cause nausea, impeding the use of potions.",
	//MonsAcidMound:       "Acid mounds are acidic creatures. They can temporarily corrode your equipment.",
	MonsHound: "Hounds are fast moving carnivore quadrupeds. They can bark, and smell you.",
	MonsYack:  "Yacks are quite large herbivorous quadrupeds. They tend to form large groups, and can push you one cell away.",
	//MonsGiantBee:        "Giant bees are fragile but extremely fast moving creatures. Their bite can sometimes enrage you.",
	//MonsGoblinWarrior:   "Goblin warriors are goblins that learned to fight, and got equipped with leather armour. They can throw javelins.",
	//MonsHydra:           "Hydras are enormous creatures with four heads that can hit you each at once.",
	//MonsSkeletonWarrior: "Skeleton warriors are good fighters, clad in chain mail.",
	//MonsSpider:          "Spiders are fast moving fragile creatures, whose bite can confuse you.",
	MonsWingedMilfid:  "Winged milfids are fast moving humanoids that can fly over you and make you swap positions. They tend to be very agressive creatures.",
	MonsBlinkingFrog:  "Blinking frogs are big frog-like creatures, whose bite can make you blink away.",
	MonsLich:          "Liches are non-living mages wearing a leather armour. They can throw a bolt of torment at you, halving your HP.",
	MonsEarthDragon:   "Earth dragons are big and hardy creatures that wander in the Underground. It is said they can be credited for many of the tunnels.",
	MonsMirrorSpecter: "Mirror specters are very insubstantial creatures, which can absorb your mana.",
	//MonsExplosiveNadre:  "Explosive nadres are very frail creatures that explode upon dying, halving HP of any adjacent creatures and occasionally destroying walls.",
	MonsSatowalgaPlant: "Satowalga Plants are immobile bushes that throw acidic projectiles at you, sometimes corroding and confusing you.",
	MonsMadNixe:        "Mad nixes are magical humanoids that can attract you to them.",
	//MonsMindCelmist:     "Mind celmists are mages that use magical smitting mind attacks that bypass armour. They can occasionally confuse or slow you. They try to avoid melee.",
	MonsVampire:       "Vampires are humanoids that drink blood to survive. Their spitting can cause nausea, impeding the use of potions.",
	MonsTreeMushroom:  "Tree mushrooms are big clunky slow-moving creatures. They can throw lignifying spores at you.",
	MonsMarevorHelith: "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past.",
}

type bandInfo struct {
	Path []position
	I    int
	Kind monsterBand
}

type monsterBand int

const (
	LoneGuard monsterBand = iota
	LoneYack
	LoneCyclop
	LoneSatowalgaPlant
	LoneBlinkingFrog
	LoneWorm
	LoneMirrorSpecter
	LoneHound
	LoneWingedMilfid
	LoneMadNixe
	LoneTreeMushroom
	LoneEarthDragon
	LoneMarevorHelith
)

type monsterBandData struct {
	Distribution map[monsterKind]int
	Band         bool
	Monster      monsterKind
	Unique       bool
}

var MonsBands = []monsterBandData{
	LoneGuard:          {Monster: MonsGuard},
	LoneYack:           {Monster: MonsYack},
	LoneCyclop:         {Monster: MonsCyclop},
	LoneSatowalgaPlant: {Monster: MonsSatowalgaPlant},
	LoneBlinkingFrog:   {Monster: MonsBlinkingFrog},
	LoneWorm:           {Monster: MonsWorm},
	LoneMirrorSpecter:  {Monster: MonsMirrorSpecter},
	LoneHound:          {Monster: MonsHound},
	LoneWingedMilfid:   {Monster: MonsWingedMilfid},
	LoneMadNixe:        {Monster: MonsMadNixe},
	LoneTreeMushroom:   {Monster: MonsTreeMushroom},
	LoneEarthDragon:    {Monster: MonsEarthDragon},
	LoneMarevorHelith:  {Monster: MonsMarevorHelith},
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
	FireReady     bool
	Seen          bool
	LOS           map[position]bool
	LastSeenState monsterState
	Swapped       bool
}

func (m *monster) Init() {
	m.HPmax = MonsData[m.Kind].maxHP
	m.Attack = MonsData[m.Kind].baseAttack
	m.HP = m.HPmax
	m.Pos = InvalidPos
	m.LastKnownPos = InvalidPos
	switch m.Kind {
	case MonsMarevorHelith:
		m.State = Wandering
	case MonsSatowalgaPlant:
		m.State = Watching
	}
}

func (m *monster) Status(st monsterStatus) bool {
	return m.Statuses[st] > 0
}

func (m *monster) Exists() bool {
	return m != nil && m.HP > 0
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
		// TODO: change the target or state?
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
	recomputeLOS := g.Player.Sees(m.Pos) && g.Dungeon.Cell(m.Pos).T == DoorCell ||
		g.Player.Sees(pos) && g.Dungeon.Cell(pos).T == DoorCell
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
		npos := m.RandomFreeNeighbor(g)
		if npos != m.Pos {
			m.Dir = npos.Dir(m.Pos)
		} else {
			m.Dir = E
		}
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
	m.Dir = g.Player.Pos.Dir(m.Pos)
	if m.Kind == MonsMarevorHelith {
		m.TeleportPlayer(g, ev)
	} else {
		m.HitPlayer(g, ev)
	}
	adelay := m.Kind.AttackDelay()
	if m.Status(MonsSlow) {
		adelay += 10
	}
	ev.Renew(g, adelay)
}

func (m *monster) NaturalAwake(g *game) {
	m.Target = m.NextTarget(g)
	m.State = Wandering
	m.GatherBand(g)
}

func (m *monster) RandomFreeNeighbor(g *game) position {
	pos := m.Pos
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	fnb := []position{}
	for _, nbpos := range neighbors {
		if !nbpos.valid() {
			continue
		}
		c := g.Dungeon.Cell(nbpos)
		if c.IsFree() {
			fnb = append(fnb, nbpos)
		}
	}
	if len(fnb) == 0 {
		return m.Pos
	}
	return fnb[RandInt(len(fnb))]
}

func (m *monster) NextTarget(g *game) position {
	band := g.Bands[m.Band]
	if len(band.Path) == 0 {
		return m.RandomFreeNeighbor(g)
	} else if len(band.Path) == 1 {
		if m.Pos.Distance(band.Path[0]) < 7+RandInt(7) {
			return m.RandomFreeNeighbor(g)
		} else {
			return band.Path[0]
		}
	}
	if band.Path[0] == m.Target {
		return band.Path[1]
	}
	return band.Path[0]
}

func (m *monster) HandleTurn(g *game, ev event) {
	if m.Swapped {
		m.Swapped = false
		ev.Renew(g, m.Kind.MovementDelay())
		return
	}
	ppos := g.Player.Pos
	mpos := m.Pos
	m.MakeAware(g)
	movedelay := m.Kind.MovementDelay()
	if m.Status(MonsSlow) {
		movedelay += 3
	}
	if m.State == Resting {
		if RandInt(3000) == 0 {
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
		switch m.State {
		case Hunting:
			if !m.SeesPlayer(g) {
				m.Dir = m.Dir.Alternate()
				if RandInt(5) == 0 {
					m.State = Watching
				}
			}
		default:
			if RandInt(4) > 0 {
				m.Dir = m.Dir.Alternate()
			}
		}
		// oklob plants are static ranged-only
		ev.Renew(g, movedelay)
		return
	}
	if mpos.Distance(ppos) == 1 && g.Dungeon.Cell(ppos).T != BarrelCell {
		if m.Status(MonsConfused) {
			ev.Renew(g, 10) // wait
			return
		}
		m.AttackAction(g, ev)
		return
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
	switch m.State {
	case Watching:
		if RandInt(5) > 0 {
			m.Dir = m.Dir.Alternate()
		} else {
			// pick a random cell: more escape strategies for the player
			if m.Kind == MonsHound && m.Pos.Distance(g.Player.Pos) <= 6 {
				m.Target = g.Player.Pos
			} else {
				m.Target = m.NextTarget(g)
			}
			m.State = Wandering
			m.GatherBand(g)
		}
		ev.Renew(g, movedelay)
		return
	case Waiting:
		if len(m.Path) < 2 {
			m.Target = m.NextTarget(g)
			m.State = Wandering
			ev.Renew(g, movedelay)
			return
		}
		if RandInt(2) == 0 {
			m.Dir = m.Dir.Alternate()
		} else if RandInt(4) == 0 {
			m.Target = m.NextTarget(g)
			m.State = Wandering
		}
		ev.Renew(g, movedelay)
		return
	}
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
			m.State = Waiting
			m.Dir = m.Dir.Alternate()
		case Hunting:
			m.State = Watching
			m.Dir = m.Dir.Alternate()
		}
		ev.Renew(g, movedelay)
		return
	}
	target := m.Path[len(m.Path)-2]
	mons := g.MonsterAt(target)
	monstarget := InvalidPos
	if mons.Exists() && len(mons.Path) >= 2 {
		monstarget = mons.Path[len(mons.Path)-2]
	}
	c := g.Dungeon.Cell(target)
	switch {
	case !mons.Exists():
		if m.Kind == MonsEarthDragon && c.T == WallCell {
			g.Dungeon.SetCell(target, GroundCell)
			g.Stats.Digs++
			if !g.Player.Sees(target) {
				g.TerrainKnowledge[m.Pos] = WallCell
			}
			g.MakeNoise(WallNoise, m.Pos)
			g.Fog(m.Pos, 1, ev)
			if g.Player.Pos.Distance(target) < 12 {
				// XXX use dijkstra distance ?
				if c.T == WallCell {
					g.Printf("%s You hear an earth-splitting noise.", g.CrackSound())
				} else if c.T == BarrelCell {
					g.Printf("%s You hear an wood-splitting noise.", g.CrackSound())
				}
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
	case mons.Pos == target && m.Pos == monstarget:
		m.MoveTo(g, target)
		m.Path = m.Path[:len(m.Path)-1]
		mons.MoveTo(g, monstarget)
		mons.Path = mons.Path[:len(mons.Path)-1]
		g.MonstersPosCache[m.Pos.idx()] = m.Index + 1
		mons.Swapped = true
		// XXX this is perhaps not the optimal to handle that case.
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
	c := g.Dungeon.Cell(m.Pos)
	if c.T == GroundCell {
		g.Dungeon.SetCell(m.Pos, FungusCell)
		invert = true
	} else if c.T == FungusCell {
		g.Dungeon.SetCell(m.Pos, GroundCell)
		invert = true
	}
	if !g.Player.Sees(m.Pos) && invert {
		_, ok := g.TerrainKnowledge[m.Pos]
		if !ok {
			g.TerrainKnowledge[m.Pos] = c.T
		}
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
	//case MonsSpider:
	//if RandInt(2) == 0 {
	//g.Confusion(ev)
	//}
	//case MonsGiantBee:
	//if RandInt(5) == 0 && !g.Player.HasStatus(StatusBerserk) && !g.Player.HasStatus(StatusExhausted) {
	//g.Player.Statuses[StatusBerserk] = 1
	//g.Player.HP += 2
	//end := ev.Rank() + DurationShortBerserk
	//g.PushEvent(&simpleEvent{ERank: end, EAction: BerserkEnd})
	//g.Player.Expire[StatusBerserk] = end
	//g.Print("You feel a sudden urge to kill things.")
	//}
	case MonsBlinkingFrog:
		if RandInt(2) == 0 {
			g.Blink(ev)
		}
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
		pos.valid() && g.Dungeon.Cell(pos).IsFree() {
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
	if m.Status(MonsConfused) {
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
	//case MonsGoblinWarrior:
	//return m.ThrowJavelin(g, ev)
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
	if len(ray) < 2 {
		// XXX see why this can happen
		return false
	}
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
			g.Confusion(ev)
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
	if m.Status(MonsConfused) {
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
		//case MonsMindCelmist:
		//return m.MindAttack(g, ev)
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
		if c.IsFree() {
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
			m.InflictDamage(g, dmg, 2)
		} else if c.T == WallCell && RandInt(2) == 0 {
			g.Dungeon.SetCell(pos, GroundCell)
			g.Stats.Digs++
			if !g.Player.Sees(pos) {
				g.TerrainKnowledge[m.Pos] = WallCell
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

func (m *monster) MakeAware(g *game) {
	if !m.SeesPlayer(g) {
		return
	}
	if m.State == Resting {
		// XXX maybe in some rare cases you could be able to move near them unnoticed
		if RandInt(3) == 0 {
			return
		}
	}
	if m.State == Resting {
		g.Printf("%s awakens.", m.Kind.Definite(true))
	} else if m.State == Wandering || m.State == Watching {
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
	if !MonsBands[g.Bands[m.Band].Kind].Band {
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

func (g *game) GenBand(band monsterBand) []monsterKind {
	mbd := MonsBands[band]
	if g.GeneratedUniques[band] > 0 && mbd.Unique {
		return nil
	}
	if !mbd.Band {
		return []monsterKind{mbd.Monster}
	}
	bandMonsters := []monsterKind{}
	for m, n := range mbd.Distribution {
		for i := 0; i < n; i++ {
			bandMonsters = append(bandMonsters, m)
		}
	}
	return bandMonsters
}

func (dg *dgen) BandInfoPatrol(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 5000 {
			pos = dg.InsideCell(g)
			break
		}
		pos = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(PlacePatrol)
	}
	target := InvalidPos
	count = 0
	for target == InvalidPos {
		// TODO: only find place in other room?
		count++
		if count > 5000 {
			pos = dg.InsideCell(g)
			break
		}
		target = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(PlacePatrol)
	}
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Path = append(bandinfo.Path, target)
	return bandinfo
}

func (dg *dgen) BandInfoOutsideGround(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideGroundCell(g))
	return bandinfo
}

func (dg *dgen) BandInfoOutside(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	return bandinfo
}

func (dg *dgen) BandInfoFoliage(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.FoliageCell(g))
	return bandinfo
}

func (dg *dgen) PutMonsterBand(g *game, band monsterBand) bool {
	monsters := g.GenBand(band)
	if monsters == nil {
		return false
	}
	var bdinf bandInfo
	switch band {
	case LoneYack, LoneWorm:
		bdinf = dg.BandInfoFoliage(g, band)
	case LoneHound, LoneEarthDragon:
		bdinf = dg.BandInfoOutsideGround(g, band)
	case LoneBlinkingFrog, LoneMirrorSpecter:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneTreeMushroom:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneSatowalgaPlant:
		bdinf = dg.BandInfoOutsideGround(g, band)
	default:
		bdinf = dg.BandInfoPatrol(g, band)
	}
	g.Bands = append(g.Bands, bdinf)
	awake := RandInt(4) > 0
	var pos position
	if len(bdinf.Path) == 0 {
		// should not happen now
		pos = g.FreeCellForMonster()
	} else {
		pos = bdinf.Path[0]
	}
	for _, mk := range monsters {
		mons := &monster{Kind: mk}
		if awake {
			mons.State = Wandering
		}
		mons.Init()
		mons.Index = len(g.Monsters)
		mons.Band = len(g.Bands) - 1
		mons.PlaceAt(g, pos)
		g.Monsters = append(g.Monsters, mons)
		pos = g.FreeCellForBandMonster(pos)
	}
	return true
}

func (dg *dgen) PutRandomBand(g *game, bands []monsterBand) bool {
	return dg.PutMonsterBand(g, bands[RandInt(len(bands))])
}

func (dg *dgen) GenMonsters(g *game) {
	g.Monsters = []*monster{}
	g.Bands = []bandInfo{}
	// TODO, just for testing now
	bandsL1 := []monsterBand{LoneGuard}
	bandsL2 := []monsterBand{LoneYack, LoneWorm, LoneHound}
	bandsL3 := []monsterBand{LoneCyclop, LoneSatowalgaPlant, LoneBlinkingFrog, LoneMirrorSpecter, LoneWingedMilfid}
	bandsL4 := []monsterBand{LoneTreeMushroom, LoneEarthDragon, LoneMadNixe}
	mlevel := 1 + RandInt(MaxDepth)
	for i := 0; i < 5; i++ {
		if !dg.PutRandomBand(g, bandsL1) {
			i--
		}
	}
	dg.PutRandomBand(g, bandsL2)
	if g.Depth > 1 {
		dg.PutRandomBand(g, bandsL2)
	}
	if g.Depth > 2 {
		dg.PutRandomBand(g, bandsL3)
	}
	if g.Depth > 3 {
		dg.PutRandomBand(g, bandsL2)
	}
	if g.Depth > 4 {
		dg.PutRandomBand(g, bandsL3)
	}
	if g.Depth > 5 {
		dg.PutRandomBand(g, bandsL2)
	}
	if g.Depth > 6 {
		dg.PutRandomBand(g, bandsL4)
	}
	if g.Depth > 7 {
		dg.PutRandomBand(g, bandsL2)
	}
	if g.Depth > 8 {
		dg.PutRandomBand(g, bandsL3)
	}
	if g.Depth > 9 {
		dg.PutRandomBand(g, bandsL2)
	}
	if g.Depth > 10 {
		dg.PutRandomBand(g, bandsL4)
		dg.PutRandomBand(g, bandsL3)
	}
	if mlevel == g.Depth {
		// XXX should really Marevor appear in more than one level?
		dg.PutMonsterBand(g, LoneMarevorHelith)
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
