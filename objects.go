package main

import (
	"errors"
	"fmt"
	"sort"
)

type objects struct {
	Stairs  map[position]stair // TODO simplify? (there's never more than one)
	Stones  map[position]stone
	Magaras map[position]magara // TODO simplify? (there's never more than one)
	Barrels map[position]bool
	Bananas map[position]bool
	Lights  map[position]bool // true: on, false: off
	Scrolls map[position]scroll
	Story   map[position]story
	Lore    map[position]int // TODO simplify? (there's never more than one)
	Items   map[position]item
}

type stair int

const (
	NormalStair stair = iota
	WinStair
	BlockedStair
)

func (st stair) ShortDesc(g *game) (desc string) {
	switch st {
	case NormalStair:
		desc = fmt.Sprintf("stairs downwards")
	case WinStair:
		desc = fmt.Sprintf("a monolith portal")
	case BlockedStair:
		desc = fmt.Sprintf("blocked stairs downwards")
	}
	return desc
}

func (st stair) Desc(g *game) (desc string) {
	switch st {
	case WinStair:
		desc = "Going through this portal will make you escape from this place, going back to the Surface."
		if g.Depth < MaxDepth {
			desc += " If you're courageous enough, you may skip this portal and continue going deeper in the dungeon, to find Marevor's magara, finishing Shaedra's failed mission."
		}
	case NormalStair:
		desc = "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
		if g.Depth == WinDepth {
			desc += " You may want to take those after freeing Shaedra from her cell."
		}
	case BlockedStair:
		desc = "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you. These are blocked by a magical barrier that you have to disable by activating a corresponding stone of barrier."
	}
	return desc
}

func (st stair) Style(g *game) (r rune, fg uicolor) {
	r = '>'
	switch st {
	case WinStair:
		fg = ColorFgMagicPlace
		r = 'Δ'
	case NormalStair:
		fg = ColorFgPlace
	case BlockedStair:
		fg = ColorFgMagicPlace
	}
	return r, fg
}

type stone int

const (
	InertStone stone = iota
	BarrelStone
	FogStone
	QueenStone
	TreeStone
	ObstructionStone
	MappingStone
	SensingStone
	// special
	SealStone
)

const NumStones = int(SensingStone) + 1

func (stn stone) String() (text string) {
	switch stn {
	case InertStone:
		text = "inert stone"
	case BarrelStone:
		text = "barrel stone"
	case FogStone:
		text = "fog stone"
	case QueenStone:
		text = "queenstone"
	case TreeStone:
		text = "tree stone"
	case ObstructionStone:
		text = "obstruction stone"
	case MappingStone:
		text = "mapping stone"
	case SensingStone:
		text = "sensing stone"
	case SealStone:
		text = "seal stone"
	}
	return text
}

func (stn stone) Desc(g *game) (text string) {
	switch stn {
	case InertStone:
		text = "This stone has been depleted of magical energies."
	case BarrelStone:
		text = "Activating this stone will teleport you away to a barrel in the same level."
	case FogStone:
		text = "Activating this stone will produce fog in a 4-radius area using harmonic energies."
	case QueenStone:
		text = "Activating this stone will produce an harmonic sound confusing enemies in a quite large area. This can also attract monsters."
	case TreeStone:
		text = "Activating this stone will lignify monsters in sight."
	case ObstructionStone:
		text = "Activating this stone will create temporal oric-energy based barriers around all monsters in sight."
	case MappingStone:
		text = "Activating this stone shows you the map layout and item locations in a wide area."
	case SensingStone:
		text = "Activating this stone shows you the current position of monsters in a wide area."
	case SealStone:
		text = "Activating this stone will disable a magical barrier somewhere in the same level, usually one blocking stairs."
	}
	return text
}

func (stn stone) ShortDesc(g *game) string {
	return fmt.Sprintf("%s", Indefinite(stn.String(), false))
}

func (stn stone) Style(g *game) (r rune, fg uicolor) {
	r = '_'
	if stn == InertStone {
		fg = ColorFgPlace
	} else if stn == SealStone {
		fg = ColorFgPlayer
	} else {
		fg = ColorFgMagicPlace
	}
	return r, fg
}

func (g *game) UseStone(pos position) {
	g.StoryPrintf("You activated %s.", g.Objects.Stones[pos].ShortDesc(g))
	g.Objects.Stones[pos] = InertStone
	g.Stats.UsedStones++
	g.Print("The stone becomes inert.")
}

const (
	FogStoneDistance   = 4
	QueenStoneDistance = 12
	MappingDistance    = 32
)

func (g *game) TeleportToBarrel() {
	barrels := []position{}
	for pos, _ := range g.Objects.Barrels {
		barrels = append(barrels, pos)
	}
	pos := barrels[RandInt(len(barrels))]
	opos := g.Player.Pos
	g.Print("You teleport away.")
	g.ui.TeleportAnimation(opos, pos, true)
	g.PlacePlayerAt(pos)
}

func (g *game) MagicMapping(ev event, maxdist int) error {
	dp := &mappingPath{game: g}
	nm := Dijkstra(dp, []position{g.Player.Pos}, maxdist)
	cdists := make(map[int][]int)
	for pos, n := range nm {
		cdists[n.Cost] = append(cdists[n.Cost], pos.idx())
	}
	var dists []int
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	g.ui.DrawDungeonView(NormalMode)
	for _, d := range dists {
		if maxdist > 0 && d > maxdist {
			continue
		}
		draw := false
		for _, i := range cdists[d] {
			pos := idxtopos(i)
			c := g.Dungeon.Cell(pos)
			if !c.Explored {
				g.Dungeon.SetExplored(pos)
				draw = true
			}
		}
		if draw {
			g.ui.MagicMappingAnimation(cdists[d])
		}
	}
	g.Printf("You feel aware of your surroundings..")
	return nil
}

func (g *game) Sensing(ev event) error {
	for _, mons := range g.Monsters {
		if mons.Exists() && !g.Player.Sees(mons.Pos) && mons.Pos.Distance(g.Player.Pos) <= MappingDistance {
			mons.UpdateKnowledge(g, mons.Pos)
		}
	}
	g.Printf("You briefly sense monsters around.")
	return nil
}

func (g *game) BarrierStone(ev event) error {
	if g.Depth == MaxDepth {
		g.Objects.Story[g.Places.Artifact] = StoryArtifact
		g.Print("You feel oric energies dissipating.")
		return nil
	}
	for pos, st := range g.Objects.Stairs {
		// actually there is at most only such stair
		if st == BlockedStair {
			g.Objects.Stairs[pos] = NormalStair
		}
	}
	g.Print("You feel oric energies dissipating.")
	return nil
}

func (g *game) ActivateStone() (err error) {
	stn, ok := g.Objects.Stones[g.Player.Pos]
	if !ok {
		return errors.New("No stone to activate here.")
	}
	oppos := g.Player.Pos
	switch stn {
	case InertStone:
		err = errors.New("Stone is inert.")
	case BarrelStone:
		g.Print("You teleport away.")
		g.TeleportToBarrel()
	case FogStone:
		g.Fog(g.Player.Pos, FogStoneDistance, g.Ev)
		g.Print("You are surrounded by fog.")
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, g.Player.Pos)
		dij := &noisePath{game: g}
		nm := Dijkstra(dij, []position{g.Player.Pos}, QueenStoneDistance)
		for _, m := range g.Monsters {
			if !m.Exists() {
				continue
			}
			if m.State == Resting {
				continue
			}
			_, ok := nm[m.Pos]
			if !ok {
				continue
			}
			m.EnterConfusion(g, g.Ev)
		}
		g.Print("The stone releases a confusing sound.")
	case TreeStone:
		count := 0
		for _, mons := range g.Monsters {
			if !mons.Exists() || !g.Player.Sees(mons.Pos) {
				continue
			}
			mons.EnterLignification(g, g.Ev)
			count++
		}
		if count == 0 {
			err = errors.New("There are no monsters to confuse around.")
		}
	case ObstructionStone:
		count := 0
		for _, mons := range g.Monsters {
			if !mons.Exists() || !g.Player.Sees(mons.Pos) {
				continue
			}
			neighbors := g.Dungeon.FreeNeighbors(mons.Pos)
			for _, pos := range neighbors {
				m := g.MonsterAt(pos)
				if m.Exists() || pos == g.Player.Pos {
					continue
				}
				g.CreateMagicalBarrierAt(pos, g.Ev)
				count++
			}
		}
		if count == 0 {
			err = errors.New("There are no monsters to be surrounded by walls.")
		} else {
			g.Print("Walls appear around your foes.")
		}
	case MappingStone:
		err = g.MagicMapping(g.Ev, MappingDistance)
	case SensingStone:
		err = g.Sensing(g.Ev)
	case SealStone:
		err = g.BarrierStone(g.Ev)
	}
	if err != nil {
		return err
	}
	g.UseStone(oppos)
	g.Ev.Renew(g, 5)
	return nil
}

type scroll int

const (
	ScrollBasics scroll = iota
	ScrollStory
	ScrollExtended
	ScrollLore
)

func (sc scroll) ShortDesc(g *game) (desc string) {
	switch sc {
	case ScrollBasics:
		desc = "the basics scroll"
	case ScrollStory, ScrollExtended:
		desc = "a story message"
	default:
		desc = "a message"
	}
	return desc
}

func (sc scroll) Text(g *game) (desc string) {
	switch sc {
	case ScrollBasics:
		desc = "the basics scroll"
	case ScrollStory:
		desc = "Your friend Shaedra got captured by some nasty people while she was trying to retrieve a powerful magara artifact that was stolen from the great magara-specialist Marevor Helith. As a gawalt monkey, you don't understand much why people complicate so much their lives caring about artifacts and the like, but one thing is clear: you have to rescue your friend, somewhere to be found in the eighth floor of this Underground area, if what you heard the guards say is true. You are small and have good night vision, so you hope the infiltration will go smoothly..."
	case ScrollExtended:
		desc = "Now that Shaedra's back to safety, you can either follow her advice, and get away from here too using the monolith portal, or you can finish the original mission: going deeper to find Marevor's powerful magara, before those mad people do bad things with it. You honestly didn't understand why it was dangerous, but Shaedra and Marevor had seemed truly concerned. Marevor said that he'll be able to create a new portal for you when you activate the artifact upon finding it."
	case ScrollLore:
		i, ok := g.Objects.Lore[g.Player.Pos]
		if !ok {
			// should not happen
			desc = "Some unintelligible notes."
			break
		}
		if i < len(LoreMessages) {
			desc = LoreMessages[i]
		}
	default:
		desc = "a message"
	}
	return desc
}

func (sc scroll) Desc(g *game) (desc string) {
	desc = "A message. It can be read by using the interact key (by default “e”)."
	return desc
}

func (sc scroll) Style(g *game) (r rune, fg uicolor) {
	r = '?'
	fg = ColorFgMagicPlace
	if sc == ScrollLore {
		fg = ColorViolet
	}
	return r, fg
}

type story int

const (
	NoStory story = iota
	StoryShaedra
	StoryMarevor
	StoryArtifact
	StoryArtifactSealed
)

func (st story) Desc(g *game) (desc string) {
	switch st {
	case StoryShaedra:
		desc = "Shaedra is the friend you came here to rescue, a human-like creature with claws, a ternian. Many other human-like creatures consider them as savages."
	case StoryMarevor:
		desc = "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past."
	case StoryArtifact:
		desc = "This is the magara that you have to retrieve: the Gem Portal Artifact that was stolen to Marevor Helith."
	case StoryArtifactSealed:
		desc = "This is the magara that you have to retrieve: the Gem Portal Artifact that was stolen to Marevor Helith. Before taking it, you have to release the magical barrier that protects it activating the corresponding protective barrier magical stone."
	}
	return desc
}

func (st story) ShortDesc(g *game) (desc string) {
	switch st {
	case StoryShaedra:
		desc = "Shaedra"
	case StoryMarevor:
		desc = "Marevor"
	case StoryArtifact:
		desc = "Gem Portal Artifact"
	case StoryArtifactSealed:
		desc = "Gem Portal Artifact (sealed)"
	}
	return desc
}

func (st story) Style(g *game) (r rune, fg uicolor) {
	fg = ColorFgPlayer
	switch st {
	case StoryShaedra:
		r = 'H'
	case StoryMarevor:
		r = 'M'
	case StoryArtifact:
		r = '='
	case StoryArtifactSealed:
		r = '='
		fg = ColorFgMagicPlace
	}
	return r, fg
}

type item int

const (
	NoItem item = iota
	CloakMagic
	CloakHear
	CloakVitality
	CloakAcrobat // no exhaustion between jumps?
	CloakShadows // reduce monster los?
	CloakSmoke
	AmuletTeleport
	AmuletConfusion
	AmuletFog
	AmuletLignification
	AmuletObstruction
	MarevorMagara
)

func (it item) IsCloak() bool {
	switch it {
	case CloakMagic,
		CloakHear,
		CloakVitality,
		CloakAcrobat,
		CloakSmoke,
		CloakShadows:
		return true
	}
	return false
}

func (it item) IsAmulet() bool {
	switch it {
	case AmuletTeleport,
		AmuletConfusion,
		AmuletFog,
		AmuletLignification,
		AmuletObstruction:
		return true
	}
	return false
}

func (it item) ShortDesc(g *game) (desc string) {
	switch it {
	case NoItem:
		desc = "empty slot"
	case CloakMagic:
		desc = "cloak of magic"
	case CloakHear:
		desc = "cloak of hearing"
	case CloakVitality:
		desc = "cloak of vitality"
	case CloakAcrobat:
		desc = "cloak of acrobatics"
	case CloakShadows:
		desc = "cloak of shadows"
	case CloakSmoke:
		desc = "cloak of smoking"
	case AmuletTeleport:
		desc = "amulet of teleport"
	case AmuletConfusion:
		desc = "amulet of confusion"
	case AmuletFog:
		desc = "amulet of fog"
	case AmuletLignification:
		desc = "amulet of lignification"
	case AmuletObstruction:
		desc = "amulet of obstruction"
	case MarevorMagara:
		desc = "Moon Portal Artifact"
	}
	return desc
}

func (it item) Desc(g *game) (desc string) {
	switch it {
	case NoItem:
		return "You do not have an item equipped on this slot."
	case CloakMagic:
		desc = "increases your magical reserves."
	case CloakHear:
		desc = "improves your hearing skills."
	case CloakVitality:
		desc = "improves your health."
	case CloakAcrobat:
		desc = "removes exhaustion from jumps."
	case CloakShadows:
		desc = "reduces the range at which foes see you in the dark."
	case CloakSmoke:
		desc = "leaves smoke behind as you move, making you difficult to spot."
	case AmuletTeleport:
		desc = "teleports away foes that critically hit you."
	case AmuletConfusion:
		desc = "confuses foes that critically hit you."
	case AmuletFog:
		desc = "releases fog and makes you swift when critically hurt."
	case AmuletLignification:
		desc = "lignifies foes that critically hit you."
	case AmuletObstruction:
		desc = "uses a magical barrier to blow away monsters that critically hit you."
	case MarevorMagara:
		desc = "magara was given to you by Marevor Helith so that he can create an escape portal when you reach Shaedra. Its sister magara, the Gem Portal Artifact, also crafted by Marevor, is the artifact that was stolen and that Shaedra was trying to retrieve before being captured. This magara needs a lot of time to recharge, so you'll only be able to use it once."
	}
	return "The " + it.ShortDesc(g) + " " + desc
}

func (it item) Style(g *game) (r rune, fg uicolor) {
	fg = ColorFgObject
	if it.IsAmulet() {
		r = '='
	} else if it.IsCloak() {
		r = '['
	}
	return r, fg
}

func (g *game) EquipItem() error {
	it, ok := g.Objects.Items[g.Player.Pos]
	if !ok {
		return errors.New("Nothing to equip here.")
	}
	var oitem item
	switch {
	case it.IsCloak():
		oitem = g.Player.Inventory.Body
		g.Player.Inventory.Body = it
	case it.IsAmulet():
		oitem = g.Player.Inventory.Neck
		g.Player.Inventory.Neck = it
	}
	if oitem != NoItem {
		g.Objects.Items[g.Player.Pos] = oitem
		g.Printf("You equip %s, leaving %s on the ground.", it.ShortDesc(g), oitem.ShortDesc(g))
		g.StoryPrintf("You equip %s, leaving %s.", it.ShortDesc(g), oitem.ShortDesc(g))
	} else {
		delete(g.Objects.Items, g.Player.Pos)
		g.Dungeon.SetCell(g.Player.Pos, GroundCell)
		g.Printf("You equip %s.", it.ShortDesc(g))
		g.StoryPrintf("You equip %s.", it.ShortDesc(g))
	}
	g.Ev.Renew(g, 5)
	return nil
}

func (g *game) RandomCloak() (it item) {
	cloaks := []item{CloakMagic,
		CloakHear,
		CloakVitality,
		CloakAcrobat,
		CloakSmoke,
		CloakShadows}
loop:
	for {
		it = cloaks[RandInt(len(cloaks))]
		for _, cl := range g.GeneratedCloaks {
			if cl == it {
				continue loop
			}
		}
		break
	}
	return it
}

func (g *game) RandomAmulet() (it item) {
	amulets := []item{AmuletTeleport,
		AmuletConfusion,
		AmuletFog,
		AmuletLignification,
		AmuletObstruction}
loop:
	for {
		it = amulets[RandInt(len(amulets))]
		for _, cl := range g.GeneratedAmulets {
			if cl == it {
				continue loop
			}
		}
		break
	}
	return it
}
