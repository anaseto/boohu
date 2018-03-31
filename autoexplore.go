package main

import "errors"

var AutoexploreMap []int

func init() {
	AutoexploreMap = make([]int, DungeonNCells)
}

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
	return g.MovePlayer(*n, ev)
}

func (g *game) AutoexploreSources() []int {
	sources := []int{}
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := g.Dungeon.CellPosition(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if !c.Explored || g.Gold[pos] > 0 || g.Collectables[pos] != nil {
			sources = append(sources, i)
		} else if _, ok := g.Rods[pos]; ok {
			sources = append(sources, i)
		}

	}
	return sources
}

func (g *game) BuildAutoexploreMap(sources []int) {
	ap := &autoexplorePath{game: g}
	g.DijkstraFast(ap, sources)
	g.DijkstraMapRebuild = false
}

func (g *game) NextAuto() (next *position, finished bool) {
	ap := &autoexplorePath{game: g}
	neighbors := ap.Neighbors(g.Player.Pos)
	if len(neighbors) == 0 {
		return nil, false
	}
	n := neighbors[0]
	ncost := AutoexploreMap[n.idx()]
	for _, pos := range neighbors[1:] {
		cost := AutoexploreMap[pos.idx()]
		if cost < ncost {
			n = pos
			ncost = cost
		}
	}
	if ncost >= AutoexploreMap[g.Player.Pos.idx()] {
		finished = true
	}
	next = &n
	return next, finished
}
