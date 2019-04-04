package main

import (
	"errors"
	"fmt"
	"sort"
)

type objects struct {
	Stairs  map[position]stair
	Stones  map[position]stone
	Magaras map[position]magara
	Barrels map[position]bool
	Bananas map[position]bool
	Lights  map[position]bool
	Scrolls map[position]scroll
	Story   map[position]story
}

type stair int

const (
	NormalStair stair = iota
	WinStair
)

func (st stair) ShortDesc(g *game) (desc string) {
	if st == WinStair {
		desc = fmt.Sprintf("a monolith portal")
	} else {
		desc = fmt.Sprintf("stairs downwards")
	}
	return desc
}

func (st stair) Desc(g *game) (desc string) {
	if st == WinStair {
		desc = "Going through this portal will make you escape from this place, going back to the Surface."
		if g.Depth < MaxDepth {
			desc += " If you're courageous enough, you may skip this portal and continue going deeper in the dungeon, to find Marevor's magara, finishing Shaedra's failed mission."
		}
	} else {
		desc = "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
		if g.Depth == WinDepth {
			desc += " You may want to take those after freeing Shaedra from her cell."
		}
	}
	return desc
}

func (st stair) Style(g *game) (r rune, fg uicolor) {
	r = '>'
	if st == WinStair {
		fg = ColorFgMagicPlace
		r = 'Δ'
	} else {
		fg = ColorFgPlace
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
	}
	return text
}

func (stn stone) Desc(g *game) (text string) {
	switch stn {
	case InertStone:
		text = "This stone has been depleted of magical energies."
	case BarrelStone:
		text = "Activating this stone will teleport you away to a random barrel."
	case FogStone:
		text = "Activating this stone will produce fog in a 4-radius area."
	case QueenStone:
		text = "Activating this stone will produce a sound confusing enemies in a quite large area. This can also attract monsters."
	case TreeStone:
		text = "Activating this stone will lignify monsters in sight."
	case ObstructionStone:
		text = "Activating this stone will create temporal walls around all monsters in sight."
	case MappingStone:
		text = "Activating this stone shows you the map layout and item locations in a wide area."
	case SensingStone:
		text = "Activating this stone shows you the current position of monsters in a wide area."
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
				g.CreateTemporalWallAt(pos, g.Ev)
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
		desc = "Your friend Shaedra got captured by some nasty people while she was trying to retrieve a powerful magara artifact that was stolen from the great magara-specialist Marevor Helith. As a gawalt monkey, you don't understand much why people complicate so much their lives caring about artifacts and the like, but one thing is clear: you have to rescue your friend, somewhere to be found in the eighth floor of this Underground area, if the rumours are true. Marevor did give to you the twin sister of the stolen artifact, saying that he'll be able to create a portal for you to flee once you find Shaedra, though he hopes you'll find the stolen artifact too. Until then, everything is up to you. You are small and have good night vision, so you hope the infiltration will go smoothly..."
	case ScrollExtended:
		desc = "Now that Shaedra's back to safety, you can either follow her advice, and get away from here too using the monolith portal, or you can finish the original mission: going deeper to find Marevor's powerful magara, before those mad people do bad things with it. You honestly didn't understand why it was dangerous, but Shaedra and Marevor had seemed truly concerned. Marevor said that he'll be able to create a new portal for you when you activate the artifact upon finding it."
	default:
		desc = "a message"
	}
	return desc
}

func (sc scroll) Desc(g *game) (desc string) {
	desc = "Messages can be read by using the interact key (by default “e”). Some explain tutorial material, and some others tell story elements."
	return desc
}

func (sc scroll) Style(g *game) (r rune, fg uicolor) {
	r = '?'
	fg = ColorFgMagicPlace
	return r, fg
}

type story int

const (
	StoryShaedra story = iota
	StoryMarevor
)

func (st story) Desc(g *game) (desc string) {
	switch st {
	case StoryShaedra:
		desc = "Shaedra is the friend you came here to rescue, a human-like creature with claws, a ternian. Many other human-like creatures consider them as savages."
	case StoryMarevor:
		desc = "Marevor Helith is an ancient undead nakrus very fond of teleporting people away. He is a well-known expert in the field of magaras - items that many people simply call magical objects. His current research focus is monolith creation. Marevor, a repentant necromancer, is now searching for his old disciple Jaixel in the Underground to help him overcome the past."
	}
	return desc
}

func (st story) ShortDesc(g *game) (desc string) {
	switch st {
	case StoryShaedra:
		desc = "Shaedra"
	case StoryMarevor:
		desc = "Marevor"
	}
	return desc
}

func (st story) Style(g *game) (r rune, fg uicolor) {
	switch st {
	case StoryShaedra:
		r = 'H'
	case StoryMarevor:
		r = 'M'
	}
	fg = ColorFgPlayer
	return r, fg
}
