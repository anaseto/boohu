package main

import "errors"

var AutoexploreMap [DungeonNCells]int

func (g *game) Autoexplore(ev event) error {
	if mons := g.MonsterInLOS(); mons.Exists() {
		return errors.New("You cannot auto-explore while there are monsters in view.")
	}
	if g.ExclusionsMap[g.Player.Pos] {
		return errors.New("You cannot auto-explore while in an excluded area.")
	}
	sources := g.AutoexploreSources()
	if g.AllExplored() {
		return errors.New("Nothing left to explore.")
	}
	if len(sources) == 0 {
		return errors.New("Some excluded places remain unexplored.")
	}
	g.BuildAutoexploreMap(sources)
	n, finished := g.NextAuto()
	if finished || n == nil {
		return errors.New("You cannot reach safely some places.")
	}
	g.Autoexploring = true
	g.AutoHalt = false
	return g.MovePlayer(*n, ev)
}

func (g *game) AllExplored() bool {
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if !c.Explored || g.Gold[pos] > 0 || g.Collectables[pos] != nil {
			return false
		} else if _, ok := g.Rods[pos]; ok {
			return false
		}
	}
	return true
}

func (g *game) AutoexploreSources() []int {
	sources := []int{}
	np := &normalPath{game: g}
	for i, c := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if c.T == WallCell {
			if len(np.Neighbors(pos)) == 0 {
				continue
			}
		}
		if g.ExclusionsMap[pos] {
			continue
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
	g.AutoExploreDijkstra(ap, sources)
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
