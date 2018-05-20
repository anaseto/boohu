package main

import "sort"

type dungeonPath struct {
	dungeon   *dungeon
	neighbors [8]position
}

func (dp *dungeonPath) Neighbors(pos position) []position {
	nb := dp.neighbors[:0]
	return pos.Neighbors(nb, position.valid)
}

func (dp *dungeonPath) Cost(from, to position) int {
	if dp.dungeon.Cell(to).T == WallCell {
		return 4
	}
	return 1
}

func (dp *dungeonPath) Estimation(from, to position) int {
	return from.Distance(to)
}

type playerPath struct {
	game      *game
	neighbors [8]position
}

func (pp *playerPath) Neighbors(pos position) []position {
	d := pp.game.Dungeon
	nb := pp.neighbors[:0]
	keep := func(npos position) bool {
		if cld, ok := pp.game.Clouds[npos]; ok && cld == CloudFire && pp.game.UnknownBurn[npos] == NoUnknownBurn {
			return false
		}
		return npos.valid() && d.Cell(npos).T != WallCell &&
			d.Cell(npos).Explored && !pp.game.UnknownDig[npos] && !pp.game.ExclusionsMap[npos]
	}
	if pp.game.Player.HasStatus(StatusConfusion) {
		nb = pos.CardinalNeighbors(nb, keep)
	} else {
		nb = pos.Neighbors(nb, keep)
	}
	return nb
}

func (pp *playerPath) Cost(from, to position) int {
	return 1
}

func (pp *playerPath) Estimation(from, to position) int {
	return from.Distance(to)
}

type noisePath struct {
	game      *game
	neighbors [8]position
}

func (fp *noisePath) Neighbors(pos position) []position {
	nb := fp.neighbors[:0]
	d := fp.game.Dungeon
	keep := func(npos position) bool {
		return npos.valid() && d.Cell(npos).T != WallCell
	}
	return pos.Neighbors(nb, keep)
}

func (fp *noisePath) Cost(from, to position) int {
	return 1
}

type normalPath struct {
	game      *game
	neighbors [8]position
}

func (np *normalPath) Neighbors(pos position) []position {
	nb := np.neighbors[:0]
	d := np.game.Dungeon
	keep := func(npos position) bool {
		return npos.valid() && d.Cell(npos).T != WallCell
	}
	if np.game.Player.HasStatus(StatusConfusion) {
		return pos.CardinalNeighbors(nb, keep)
	}
	return pos.Neighbors(nb, keep)
}

func (np *normalPath) Cost(from, to position) int {
	return 1
}

type autoexplorePath struct {
	game      *game
	neighbors [8]position
}

func (ap *autoexplorePath) Neighbors(pos position) []position {
	if ap.game.ExclusionsMap[pos] {
		return nil
	}
	d := ap.game.Dungeon
	nb := ap.neighbors[:0]
	keep := func(npos position) bool {
		if cld, ok := ap.game.Clouds[npos]; ok && cld == CloudFire && ap.game.UnknownBurn[npos] == NoUnknownBurn {
			// XXX little info leak
			return false
		}
		return npos.valid() && d.Cell(npos).T != WallCell && !ap.game.ExclusionsMap[npos]
	}
	if ap.game.Player.HasStatus(StatusConfusion) {
		nb = pos.CardinalNeighbors(nb, keep)
	} else {
		nb = pos.Neighbors(nb, keep)
	}
	return nb
}

func (ap *autoexplorePath) Cost(from, to position) int {
	return 1
}

type monPath struct {
	game      *game
	monster   *monster
	wall      bool
	neighbors [8]position
}

func (mp *monPath) Neighbors(pos position) []position {
	nb := mp.neighbors[:0]
	d := mp.game.Dungeon
	keep := func(npos position) bool {
		return npos.valid() && (d.Cell(npos).T != WallCell || mp.wall)
	}
	if mp.monster.Status(MonsConfused) {
		return pos.CardinalNeighbors(nb, keep)
	}
	return pos.Neighbors(nb, keep)
}

func (mp *monPath) Cost(from, to position) int {
	g := mp.game
	mons := g.MonsterAt(to)
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
		pp := &dungeonPath{dungeon: g.Dungeon}
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
