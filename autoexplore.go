package main

import "errors"

func (g *game) Autoexplore(ev event) error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.Pos] {
		return errors.New("You cannot auto-explore while in an excluded area.")
	}
	sources := g.AutoexploreSources()
	if len(sources) == 0 {
		return errors.New("Nothing left to explore.")
	}
	g.BuildAutoexploreMap(sources)
	n, _ := g.NextAuto()
	if n == nil {
		return errors.New("Some unexplored parts not safely reachable remain.")
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

func (g *game) BuildAutoexploreMap(sources []position) {
	ap := &autoexplorePath{game: g}
	g.AutoexploreMap = Dijkstra(ap, sources, 9999)
}

func (g *game) NextAuto() (*node, bool) {
	rebuild := false
	ap := &autoexplorePath{game: g}
	neighbors := ap.Neighbors(g.Player.Pos)
	if len(neighbors) == 0 || g.ExclusionsMap[g.Player.Pos] {
		return nil, false
	}
	next, ok := g.AutoexploreMap[neighbors[0]]
	if !ok {
		return nil, rebuild
	}
	for _, pos := range neighbors[1:] {
		n := g.AutoexploreMap[pos]
		if n != nil && n.Cost < next.Cost {
			next = n
		}
	}
	if next.Cost >= g.AutoexploreMap[g.Player.Pos].Cost {
		rebuild = true
	}
	return next, rebuild
}
