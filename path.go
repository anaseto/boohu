package main

import "sort"

func (d *dungeon) FreeNeighbors(pos position) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	nb := make([]position, 0, 8)
	for _, npos := range neighbors {
		if d.Valid(npos) && d.Cell(npos).T != WallCell {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(pos position) []position {
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	nb := make([]position, 0, 4)
	for _, npos := range neighbors {
		if d.Valid(npos) && d.Cell(npos).T != WallCell {
			nb = append(nb, npos)
		}
	}
	return nb
}

type playerPath struct {
	game *game
}

func (pp *playerPath) Neighbors(pos position) []position {
	m := pp.game.Dungeon
	var neighbors []position
	if pp.game.Player.HasStatus(StatusConfusion) {
		neighbors = m.CardinalFreeNeighbors(pos)
	} else {
		neighbors = m.FreeNeighbors(pos)
	}
	freeNeighbors := make([]position, 0, len(neighbors))
	for _, npos := range neighbors {
		if m.Cell(npos).Explored && !pp.game.UnknownDig[npos] && !pp.game.ExclusionsMap[npos] {
			freeNeighbors = append(freeNeighbors, npos)
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

type noisePath struct {
	game *game
}

func (fp *noisePath) Neighbors(pos position) []position {
	return fp.game.Dungeon.FreeNeighbors(pos)
}

func (fp *noisePath) Cost(from, to position) int {
	return 1
}

type normalPath struct {
	game *game
}

func (np *normalPath) Neighbors(pos position) []position {
	if np.game.Player.HasStatus(StatusConfusion) {
		return np.game.Dungeon.CardinalFreeNeighbors(pos)
	}
	return np.game.Dungeon.FreeNeighbors(pos)
}

func (np *normalPath) Cost(from, to position) int {
	return 1
}

type autoexplorePath struct {
	game *game
}

func (ap *autoexplorePath) Neighbors(pos position) []position {
	if ap.game.ExclusionsMap[pos] {
		return nil
	}
	var neighbors []position
	if ap.game.Player.HasStatus(StatusConfusion) {
		neighbors = ap.game.Dungeon.CardinalFreeNeighbors(pos)
	} else {
		neighbors = ap.game.Dungeon.FreeNeighbors(pos)
	}
	var suitableNeighbors []position
	for _, pos := range neighbors {
		if !ap.game.ExclusionsMap[pos] {
			suitableNeighbors = append(suitableNeighbors, pos)
		}
	}
	return suitableNeighbors
}

func (ap *autoexplorePath) Cost(from, to position) int {
	return 1
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
		}
		return mp.game.Dungeon.CardinalFreeNeighbors(pos)
	}
	if mp.wall {
		return mp.game.Dungeon.Neighbors(pos)
	}
	return mp.game.Dungeon.FreeNeighbors(pos)
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
		return nil
	}
	return path
}

func (g *game) SortedNearestTo(cells []position, to position) []position {
	ps := posSlice{}
	for _, pos := range cells {
		pp := &playerPath{game: g}
		_, cost, found := AstarPath(pp, pos, to)
		if found {
			ps = append(ps, posCost{pos, cost})
		}
	}
	sort.Sort(ps)
	sorted := []position{}
	for _, pc := range ps {
		sorted = append(sorted, pc.pos)
	}
	return sorted
}

type posCost struct {
	pos  position
	cost int
}

type posSlice []posCost

func (ps posSlice) Len() int           { return len(ps) }
func (ps posSlice) Swap(i, j int)      { ps[i], ps[j] = ps[j], ps[i] }
func (ps posSlice) Less(i, j int) bool { return ps[i].cost < ps[j].cost }
