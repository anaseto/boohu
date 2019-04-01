package main

import (
	"sort"
)

type dungeonPath struct {
	dungeon   *dungeon
	neighbors [8]position
	wcost     int
}

func (dp *dungeonPath) Neighbors(pos position) []position {
	nb := dp.neighbors[:0]
	return pos.CardinalNeighbors(nb, func(npos position) bool { return npos.valid() })
}

func (dp *dungeonPath) Cost(from, to position) int {
	if dp.dungeon.Cell(to).T == WallCell {
		if dp.wcost > 0 {
			return dp.wcost
		}
		return 4
	}
	return 1
}

func (dp *dungeonPath) Estimation(from, to position) int {
	return from.Distance(to)
}

type mappingPath struct {
	game      *game
	neighbors [8]position
	wcost     int
}

func (dp *mappingPath) Neighbors(pos position) []position {
	d := dp.game.Dungeon
	if d.Cell(pos).T == WallCell {
		return nil
	}
	nb := dp.neighbors[:0]
	keep := func(npos position) bool {
		return npos.valid()
	}
	return pos.CardinalNeighbors(nb, keep)
}

func (dp *mappingPath) Cost(from, to position) int {
	return 1
}

func (dp *mappingPath) Estimation(from, to position) int {
	return from.Distance(to)
}

type tunnelPath struct {
	dg        *dgen
	neighbors [4]position
	area      [9]position
}

func (tp *tunnelPath) Neighbors(pos position) []position {
	nb := tp.neighbors[:0]
	return pos.CardinalNeighbors(nb, func(npos position) bool { return npos.valid() })
}

func (tp *tunnelPath) Cost(from, to position) int {
	if tp.dg.room[from] && !tp.dg.tunnel[from] {
		return 50
	}
	wc := tp.dg.WallAreaCount(tp.area[:0], from, 1)
	if tp.dg.d.Cell(from).IsPassable() {
		return 1
	}
	return 10 - wc
}

func (tp *tunnelPath) Estimation(from, to position) int {
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
		t, okT := pp.game.TerrainKnowledge[npos]
		if cld, ok := pp.game.Clouds[npos]; ok && cld == CloudFire && (!okT || t != FungusCell && t != DoorCell) {
			return false
		}
		return npos.valid() && (d.Cell(npos).T != WallCell && (!okT || t != WallCell) || pp.game.Player.HasStatus(StatusDig)) &&
			d.Cell(npos).Explored
	}
	nb = pos.CardinalNeighbors(nb, keep)
	return nb
}

func (pp *playerPath) Cost(from, to position) int {
	if !pp.game.ExclusionsMap[from] && pp.game.ExclusionsMap[to] {
		return unreachable
	}
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
	return pos.CardinalNeighbors(nb, keep)
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
		return npos.valid() && d.Cell(npos).IsPassable()
	}
	return pos.CardinalNeighbors(nb, keep)
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
		t, okT := ap.game.TerrainKnowledge[npos]
		if cld, ok := ap.game.Clouds[npos]; ok && cld == CloudFire && (!okT || t != FungusCell && t != DoorCell) {
			// XXX little info leak
			return false
		}
		return npos.valid() && (d.Cell(npos).IsPassable() && (!okT || t != WallCell)) &&
			!ap.game.ExclusionsMap[npos]
	}
	//if ap.game.Player.HasStatus(StatusConfusion) {
	nb = pos.CardinalNeighbors(nb, keep)
	//} else {
	//nb = pos.Neighbors(nb, keep)
	//}
	return nb
}

func (ap *autoexplorePath) Cost(from, to position) int {
	return 1
}

type monPath struct {
	game      *game
	monster   *monster
	destruct  bool
	neighbors [8]position
}

func (mp *monPath) Neighbors(pos position) []position {
	nb := mp.neighbors[:0]
	d := mp.game.Dungeon
	keep := func(npos position) bool {
		if !npos.valid() {
			return false
		}
		c := d.Cell(npos)
		return (c.IsPassable() || c.IsDestructible() && mp.destruct) && (c.T != DoorCell || mp.monster.Kind.CanOpenDoors() || mp.destruct)
	}
	ret := pos.CardinalNeighbors(nb, keep)
	// shuffle so that monster movement is not unnaturally predictable
	for i := 0; i < len(ret); i++ {
		j := i + RandInt(len(ret)-i)
		ret[i], ret[j] = ret[j], ret[i]
	}
	return ret
}

func (mp *monPath) Cost(from, to position) int {
	g := mp.game
	mons := g.MonsterAt(to)
	if !mons.Exists() {
		c := g.Dungeon.Cell(to)
		if mp.destruct && (!c.IsPassable() || c.T == DoorCell) && mp.monster.State != Hunting {
			return 6
		}
		if to == g.Player.Pos && mp.monster.Kind.Peaceful() {
			return 4
		}
		return 1
	}
	if mons.Status(MonsLignified) {
		return 8
	}
	return 4
}

func (mp *monPath) Estimation(from, to position) int {
	return from.Distance(to)
}

func (m *monster) APath(g *game, from, to position) []position {
	mp := &monPath{game: g, monster: m}
	if m.Kind == MonsEarthDragon {
		mp.destruct = true
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
		pp := &dungeonPath{dungeon: g.Dungeon, wcost: unreachable}
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
