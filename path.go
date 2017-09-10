package main

import (
	"log"
)

func (d *dungeon) FreeNeighbors(pos position) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	freeNeighbors := []position{}
	for _, c := range neighbors {
		if d.Valid(c) && d.Cell(c).T != WallCell {
			freeNeighbors = append(freeNeighbors, c)
		}
	}
	return freeNeighbors
}

func (d *dungeon) CardinalFreeNeighbors(pos position) []position {
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	freeNeighbors := []position{}
	for _, c := range neighbors {
		if d.Valid(c) && d.Cell(c).T != WallCell {
			freeNeighbors = append(freeNeighbors, c)
		}
	}
	return freeNeighbors
}

type playerPath struct {
	game *game
}

func (pp *playerPath) Neighbors(pos position) []position {
	m := pp.game.Dungeon
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	freeNeighbors := []position{}
	for _, c := range neighbors {
		if m.Valid(c) && m.Cell(c).T != WallCell && m.Cell(c).Explored {
			freeNeighbors = append(freeNeighbors, c)
		}
	}
	return freeNeighbors
}

func (pp *playerPath) Cost(from, to position) int {
	return 1
}

func (pp *playerPath) Estimation(from, to position) int {
	return from.Distance(to)
}

type monPath struct {
	game    *game
	monster *monster
	wall    bool
}

func (mp *monPath) Neighbors(pos position) []position {
	if mp.monster.Status(MonsConfused) {
		if mp.wall {
			return mp.game.Dungeon.CardinalNeighbors(pos)
		} else {
			return mp.game.Dungeon.CardinalFreeNeighbors(pos)
		}
	} else {
		if mp.wall {
			return mp.game.Dungeon.Neighbors(pos)
		} else {
			return mp.game.Dungeon.FreeNeighbors(pos)
		}
	}
}

func (mp *monPath) Cost(from, to position) int {
	g := mp.game
	mons, _ := g.MonsterAt(to)
	if !mons.Exists() {
		if mp.wall && g.Dungeon.Cell(to).T == WallCell && mp.monster.State != Hunting {
			return 6
		}
		return 1
	}
	return 4
}

func (mp *monPath) Estimation(from, to position) int {
	return from.Distance(to)
}

func (m *monster) APath(g *game, from, to position) []position {
	mp := &monPath{game: g, monster: m}
	if m.Kind == MonsEarthDragon {
		mp.wall = true
	}
	path, _, found := AstarPath(mp, from, to)
	if !found {
		return nil
	}
	return path
}

func (g *game) PlayerPath(from, to position) []position {
	pp := &playerPath{game: g}
	path, _, found := AstarPath(pp, from, to)
	if !found {
		log.Printf("no path from %+v to %+v\n", from, to)
		g.Print("No path found to there.")
		return nil
	}
	return path
}
