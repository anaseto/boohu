package main

import (
	"errors"
	"fmt"
)

func (g *game) Autoexplore(ev event) error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return fmt.Errorf("You cannot auto-explore while there are monsters in view.")
	}
	g.BuildAutoexploreMap()
	n, _ := g.NextAuto()
	if n == nil {
		return errors.New("Nothing left to explore.")
	}
	g.Autoexploring = true
	g.AutoHalt = false
	return g.MovePlayer(n.Pos, ev)
}

func (g *game) AutoexploreSources() []position {
	sources := []position{}
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := g.Dungeon.CellPosition(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if !c.Explored || g.Gold[pos] > 0 || g.Collectables[pos] != nil {
			sources = append(sources, pos)
		} else if _, ok := g.Rods[pos]; ok {
			sources = append(sources, pos)
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap() {
	sources := g.AutoexploreSources()
	np := &normalPath{game: g}
	g.AutoexploreMap = Dijkstra(np, sources, 9999)
}

func (g *game) NextAuto() (*node, bool) {
	rebuild := false
	np := &normalPath{game: g}
	neighbors := np.Neighbors(g.Player.Pos)
	next, ok := g.AutoexploreMap[neighbors[0]]
	if !ok {
		return nil, rebuild
	}
	for _, pos := range neighbors[1:] {
		n := g.AutoexploreMap[pos]
		if n.Cost < next.Cost {
			next = n
		}
	}
	if next.Cost >= g.AutoexploreMap[g.Player.Pos].Cost {
		rebuild = true
	}
	return next, rebuild
}
