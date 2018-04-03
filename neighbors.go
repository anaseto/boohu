package main

func (pos position) Neighbors(nb []position, keep func(position) bool) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (pos position) CardinalNeighbors(nb []position, keep func(position) bool) []position {
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	nb = nb[:0]
	for _, npos := range neighbors {
		if keep(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (d *dungeon) OutsideNeighbors(pos position) []position {
	nb := make([]position, 0, 8)
	nb = pos.Neighbors(nb, func(npos position) bool {
		return !d.Valid(npos)
	})
	return nb
}

func (d *dungeon) Neighbors(pos position) []position {
	nb := make([]position, 0, 8)
	nb = pos.Neighbors(nb, d.Valid)
	return nb
}

func (d *dungeon) CardinalNeighbors(pos position) []position {
	nb := make([]position, 0, 4)
	nb = pos.CardinalNeighbors(nb, d.Valid)
	return nb
}

func (d *dungeon) FreeNeighbors(pos position) []position {
	nb := make([]position, 0, 8)
	nb = pos.Neighbors(nb, func(npos position) bool {
		return d.Valid(npos) && d.Cell(npos).T != WallCell
	})
	return nb
}

func (d *dungeon) CardinalFreeNeighbors(pos position) []position {
	nb := make([]position, 0, 4)
	nb = pos.CardinalNeighbors(nb, func(npos position) bool {
		return d.Valid(npos) && d.Cell(npos).T != WallCell
	})
	return nb
}
