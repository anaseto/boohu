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
}

type stair int

const (
	NormalStair stair = iota
	WinStair
)

func (st stair) ShortDesc(g *game) (desc string) {
	if st == WinStair {
		desc = fmt.Sprintf("glowing stairs")
	} else {
		desc = fmt.Sprintf("stairs downwards")
	}
	return desc
}

func (st stair) Desc(g *game) (desc string) {
	if st == WinStair {
		desc = "These shiny-looking stairs are in fact a magical monolith. It is said they were made some centuries ago by Marevor Helith. They will lead you back to your village."
		if g.Depth < MaxDepth {
			desc += " Note that this is not the last floor, so you may want to find a normal stair and continue collecting simellas, if you're courageous enough."
		}
	} else {
		desc = "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
		if g.Depth == WinDepth {
			desc += " If you're afraid, you could instead just win by taking the magical stairs somewhere in the same map."
		}
	}
	return desc
}

func (st stair) Style(g *game) (r rune, fg uicolor) {
	r = '>'
	if st == WinStair {
		fg = ColorFgMagicPlace
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
	MappingDistance    = 25
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
	dp := &dungeonPath{dungeon: g.Dungeon}
	g.AutoExploreDijkstra(dp, []int{g.Player.Pos.idx()})
	cdists := make(map[int][]int)
	for i, dist := range DijkstraMapCache {
		cdists[dist] = append(cdists[dist], i)
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
			if (c.IsFree() || g.Dungeon.HasFreeNeighbor(pos)) && !c.Explored {
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
	switch stn {
	case InertStone:
		err = errors.New("Stone is inert.")
	case BarrelStone:
		oppos := g.Player.Pos
		g.Print("You teleport away.")
		g.TeleportToBarrel()
		g.UseStone(oppos)
	case FogStone:
		g.Fog(g.Player.Pos, FogStoneDistance, g.Ev)
		g.Print("You are surrounded by fog.")
		g.UseStone(g.Player.Pos)
	case QueenStone:
		g.MakeNoise(QueenStoneNoise, g.Player.Pos)
		dij := &normalPath{game: g}
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
		g.UseStone(g.Player.Pos)
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
		} else {
			g.UseStone(g.Player.Pos)
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
			g.UseStone(g.Player.Pos)
		}
	case MappingStone:
		err = g.MagicMapping(g.Ev, MappingDistance)
	case SensingStone:
		err = g.Sensing(g.Ev)
	}
	if err != nil {
		return err
	}
	g.Ev.Renew(g, 5)
	return nil
}
