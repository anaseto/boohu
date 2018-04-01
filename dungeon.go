// many ideas here from articles found at http://www.roguebasin.com/

package main

import (
	"sort"
)

type dungeon struct {
	Cells  []cell
	Width  int
	Height int
}

type cell struct {
	T        terrain
	Explored bool
}

type terrain int

const (
	WallCell terrain = iota
	FreeCell
)

type room struct {
	pos position
	w   int
	h   int
}

func (d *dungeon) Cell(pos position) cell {
	return d.Cells[pos.Y*d.Width+pos.X]
}

func (d *dungeon) Valid(pos position) bool {
	return pos.X < d.Width && pos.Y < d.Height && pos.X >= 0 && pos.Y >= 0
}

func (d *dungeon) Border(pos position) bool {
	return pos.X == d.Width-1 || pos.Y == d.Height-1 || pos.X == 0 || pos.Y == 0
}

func (d *dungeon) OutsideNeighbors(pos position) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	nb := make([]position, 0, 8)
	for _, npos := range neighbors {
		if !d.Valid(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (d *dungeon) SetCell(pos position, t terrain) {
	d.Cells[pos.Y*d.Width+pos.X].T = t
}

func (d *dungeon) SetExplored(pos position) {
	d.Cells[pos.Y*d.Width+pos.X].Explored = true
}

func roomDistance(r1, r2 room) int {
	return Abs(r1.pos.X-r2.pos.X) + Abs(r1.pos.Y-r2.pos.Y)
}

func nearRoom(rooms []room, r room) room {
	closest := rooms[0]
	d := roomDistance(r, closest)
	for _, nextRoom := range rooms {
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := RandInt(10)
			if n > 3 {
				d = nd
				closest = nextRoom
			}
		}
	}
	return closest
}

func nearestRoom(rooms []room, r room) room {
	closest := rooms[0]
	d := roomDistance(r, closest)
	for _, nextRoom := range rooms {
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := RandInt(10)
			if n > 0 {
				d = nd
				closest = nextRoom
			}
		}
	}
	return closest
}

func intersectsRoom(rooms []room, r room) bool {
	for _, rr := range rooms {
		if (r.pos.X+r.w-1 >= rr.pos.X && rr.pos.X+rr.w-1 >= r.pos.X) &&
			(r.pos.Y+r.h-1 >= rr.pos.Y && rr.pos.Y+rr.h-1 >= r.pos.Y) {
			return true
		}
	}
	return false
}

func (d *dungeon) connectRooms(r1, r2 room) {
	x := r1.pos.X
	y := r1.pos.Y
	count := 0
	for {
		count++
		if count > 1000 {
			panic("ConnectRooms")
		}
		if x < r2.pos.X {
			x++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X {
			x--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y < r2.pos.Y {
			y++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y > r2.pos.Y {
			y--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		break
	}
}

func (d *dungeon) connectRoomsDiagonally(r1, r2 room) {
	x := r1.pos.X
	y := r1.pos.Y
	count := 0
	for {
		count++
		if count > 1000 {
			panic("ConnectRooms")
		}
		if x < r2.pos.X && y < r2.pos.Y {
			x++
			d.SetCell(position{x, y}, FreeCell)
			y++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X && y < r2.pos.Y {
			x--
			d.SetCell(position{x, y}, FreeCell)
			y++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X && y > r2.pos.Y {
			x--
			d.SetCell(position{x, y}, FreeCell)
			y--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x < r2.pos.X && y > r2.pos.Y {
			x++
			d.SetCell(position{x, y}, FreeCell)
			y--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x < r2.pos.X {
			x++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X {
			x--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y < r2.pos.Y {
			y++
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y > r2.pos.Y {
			y--
			d.SetCell(position{x, y}, FreeCell)
			continue
		}
		break
	}
}

func (d *dungeon) Neighbors(pos position) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	nb := make([]position, 0, 8)
	for _, npos := range neighbors {
		if d.Valid(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (d *dungeon) CardinalNeighbors(pos position) []position {
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	nb := make([]position, 0, 4)
	for _, npos := range neighbors {
		if d.Valid(npos) {
			nb = append(nb, npos)
		}
	}
	return nb
}

func (d *dungeon) Area(pos position, radius int) []position {
	area := []position{}
	for x := pos.X - radius; x <= pos.X+radius; x++ {
		for y := pos.Y - radius; y <= pos.Y+radius; y++ {
			pos := position{x, y}
			if d.Valid(pos) {
				area = append(area, pos)
			}
		}
	}
	return area
}

type dungeonPath struct {
	dungeon *dungeon
}

func (dp *dungeonPath) Neighbors(pos position) []position {
	return dp.dungeon.Neighbors(pos)
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

func (d *dungeon) ConnectRoomsShortestPath(r1, r2 room) {
	mp := &dungeonPath{dungeon: d}
	path, _, _ := AstarPath(mp, r1.pos, r2.pos)
	for _, pos := range path {
		d.SetCell(pos, FreeCell)
	}
}

func (d *dungeon) PutRoom(r room) {
	for i := r.pos.X; i < r.pos.X+r.w; i++ {
		for j := r.pos.Y; j < r.pos.Y+r.h; j++ {
			if d.Valid(position{i, j}) {
				d.SetCell(position{i, j}, FreeCell)
			}
		}
	}
}

func (d *dungeon) CellPosition(i int) position {
	return position{i % d.Width, i / d.Width}
}

func (g *game) GenRuinsMap(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	rooms := []room{}
	noIntersect := true
	//if randInt(100) > 50 {
	//noIntersect = false
	//}
	for i := 0; i < 45; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				pos: position{RandInt(w - 1), RandInt(h - 1)},
				w:   3 + RandInt(5),
				h:   2 + RandInt(3)}
			if !noIntersect {
				break
			}
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.PutRoom(ro)
		if len(rooms) > 0 {
			r := RandInt(100)
			if r > 75 {
				d.connectRooms(nearRoom(rooms, ro), ro)
			} else if r > 25 {
				d.ConnectRoomsShortestPath(nearRoom(rooms, ro), ro)
			} else {
				d.connectRoomsDiagonally(nearRoom(rooms, ro), ro)
			}
		}
		rooms = append(rooms, ro)
	}
	g.Dungeon = d
	g.PutDoors(20)
}

type roomSlice []room

func (rs roomSlice) Len() int      { return len(rs) }
func (rs roomSlice) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs roomSlice) Less(i, j int) bool {
	return rs[i].pos.Y < rs[j].pos.Y || rs[i].pos.Y == rs[j].pos.Y && rs[i].pos.X < rs[j].pos.X
}

func (g *game) GenRoomMap(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	rooms := []room{}
	noIntersect := true
	for i := 0; i < 35; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				pos: position{RandInt(w - 1), RandInt(h - 1)},
				w:   5 + RandInt(4),
				h:   3 + RandInt(3)}
			if !noIntersect {
				break
			}
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.PutRoom(ro)
		rooms = append(rooms, ro)
	}
	sort.Sort(roomSlice(rooms))
	for i, ro := range rooms {
		if i == 0 {
			continue
		}
		r := RandInt(100)
		if r > 50 {
			d.connectRooms(nearestRoom(rooms[:i], ro), ro)
		} else if r > 25 {
			d.ConnectRoomsShortestPath(nearRoom(rooms[:i], ro), ro)
		} else {
			d.connectRoomsDiagonally(nearestRoom(rooms[:i], ro), ro)
		}
	}
	g.Dungeon = d
	g.PutDoors(90)
}

func (d *dungeon) FreeCell() position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(d.Width)
		y := RandInt(d.Height)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T == FreeCell {
			return pos
		}
	}
}

func (d *dungeon) WallCell() position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("WallCell")
		}
		x := RandInt(d.Width)
		y := RandInt(d.Height)
		pos := position{x, y}
		c := d.Cell(pos)
		if c.T == WallCell {
			return pos
		}
	}
}

func (g *game) GenCaveMap(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	pos := position{40, 10}
	max := 21 * 42
	d.SetCell(pos, FreeCell)
	cells := 1
	notValid := 0
	lastValid := pos
	diag := RandInt(4) == 0
	for cells < max {
		npos := pos.RandomNeighbor(diag)
		if !d.Valid(pos) && d.Valid(npos) && d.Cell(npos).T == WallCell {
			pos = lastValid
			continue
		}
		pos = npos
		if d.Valid(pos) {
			if d.Cell(pos).T != FreeCell {
				d.SetCell(pos, FreeCell)
				cells++
			}
			lastValid = pos
		} else {
			notValid++
		}
		if notValid > 200 {
			notValid = 0
			pos = lastValid
		}
	}
	cells = 1
	max = d.Height * 1
	digs := 0
	i := 0
loop:
	for cells < max {
		i++
		if digs > 3 {
			break
		}
		if i > 1000 {
			break
		}
		diag = RandInt(2) == 0
		block := d.DigBlock(diag)
		if len(block) == 0 {
			continue loop
		}
		if len(block) < 4 || len(block) > 10 {
			continue loop
		}
		for _, pos := range block {
			d.SetCell(pos, FreeCell)
			cells++
		}
		digs++
	}
	g.Dungeon = d
	g.PutDoors(5)
}

func (d *dungeon) HasFreeNeighbor(pos position) bool {
	neighbors := d.Neighbors(pos)
	for _, pos := range neighbors {
		if d.Cell(pos).T == FreeCell {
			return true
		}
	}
	return false
}

func (g *game) HasFreeExploredNeighbor(pos position) bool {
	d := g.Dungeon
	neighbors := d.Neighbors(pos)
	for _, pos := range neighbors {
		c := d.Cell(pos)
		if c.T == FreeCell && c.Explored && !g.UnknownDig[pos] {
			return true
		}
	}
	return false
}

func (d *dungeon) DigBlock(diag bool) []position {
	pos := d.WallCell()
	block := []position{}
	for {
		block = append(block, pos)
		if d.HasFreeNeighbor(pos) {
			break
		}
		pos = pos.RandomNeighbor(diag)
		if !d.Valid(pos) {
			block = block[:0]
			pos = d.WallCell()
			continue
		}
		if !d.Valid(pos) {
			return nil
		}
	}
	return block
}

func (g *game) GenCaveMapTree(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	center := position{40, 10}
	d.SetCell(center, FreeCell)
	d.SetCell(center.E(), FreeCell)
	d.SetCell(center.NE(), FreeCell)
	d.SetCell(center.S(), FreeCell)
	d.SetCell(center.SE(), FreeCell)
	d.SetCell(center.N(), FreeCell)
	d.SetCell(center.NW(), FreeCell)
	d.SetCell(center.W(), FreeCell)
	d.SetCell(center.SW(), FreeCell)
	max := 21 * 23
	cells := 1
	diag := RandInt(2) == 0
loop:
	for cells < max {
		block := d.DigBlock(diag)
		if len(block) == 0 {
			continue loop
		}
		for _, pos := range block {
			if d.Cell(pos).T != FreeCell {
				d.SetCell(pos, FreeCell)
				cells++
			}
		}
	}
	g.Dungeon = d
	g.PutDoors(5)
}

func (d *dungeon) WallAreaCount(pos position, radius int) int {
	neighbors := d.Area(pos, radius)
	count := 0
	for _, npos := range neighbors {
		if d.Cell(npos).T == WallCell {
			count++
		}
	}
	switch radius {
	case 1:
		count += 9 - len(neighbors)
	case 2:
		count += 25 - len(neighbors)
	}
	return count
}

func (d *dungeon) Connected(pos position) (map[position]bool, int) {
	conn := map[position]bool{}
	stack := []position{pos}
	count := 0
	conn[pos] = true
	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		neighbors := d.FreeNeighbors(pos)
		for _, npos := range neighbors {
			if !conn[npos] {
				conn[npos] = true
				stack = append(stack, npos)
			}
		}
	}
	return conn, count
}

func (d *dungeon) connex() bool {
	pos := d.FreeCell()
	conn, _ := d.Connected(pos)
	for i, c := range d.Cells {
		if c.T == FreeCell && !conn[d.CellPosition(i)] {
			return false
		}
	}
	return true
}

func (g *game) RunCellularAutomataCave(h, w int) bool {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	for i := range d.Cells {
		r := RandInt(100)
		pos := d.CellPosition(i)
		if r >= 45 {
			d.SetCell(pos, FreeCell)
		} else {
			d.SetCell(pos, WallCell)
		}
	}
	bufm := &dungeon{}
	bufm.Cells = make([]cell, h*w)
	bufm.Width = w
	bufm.Height = h
	for i := 0; i < 5; i++ {
		for j := range bufm.Cells {
			pos := d.CellPosition(j)
			c1 := d.WallAreaCount(pos, 1)
			if c1 >= 5 {
				bufm.SetCell(pos, WallCell)
			} else {
				bufm.SetCell(pos, FreeCell)
			}
			if i == 3 {
				c2 := d.WallAreaCount(pos, 2)
				if c2 <= 2 {
					bufm.SetCell(pos, WallCell)
				}
			}
		}
		copy(d.Cells, bufm.Cells)
	}
	var conn map[position]bool
	var count int
	var winner position
	for i := 0; i < 15; i++ {
		pos := d.FreeCell()
		if conn[pos] {
			continue
		}
		var ncount int
		conn, ncount = d.Connected(pos)
		if ncount > count {
			count = ncount
			winner = pos
		}
		if count >= 37*d.Height*d.Width/100 {
			break
		}
	}
	conn, count = d.Connected(winner)
	if count <= 37*d.Height*d.Width/100 {
		return false
	}
	for i, c := range d.Cells {
		pos := d.CellPosition(i)
		if c.T == FreeCell && !conn[pos] {
			d.SetCell(pos, WallCell)
		}
	}
	max := d.Height * 1
	cells := 1
	digs := 0
	i := 0
loop:
	for cells < max {
		i++
		if digs > 4 {
			break
		}
		if i > 1000 {
			break
		}
		diag := RandInt(2) == 0
		block := d.DigBlock(diag)
		if len(block) == 0 {
			continue loop
		}
		if len(block) < 4 || len(block) > 10 {
			continue loop
		}
		for _, pos := range block {
			d.SetCell(pos, FreeCell)
			cells++
		}
		digs++
	}
	g.Dungeon = d
	g.PutDoors(10)
	return true
}

func (g *game) GenCellularAutomataCaveMap(h, w int) {
	count := 0
	for {
		count++
		if count > 100 {
			panic("genCellularAutomataCaveMap")
		}
		if g.RunCellularAutomataCave(h, w) {
			break
		}
	}
}

type vegetation int

const (
	foliage vegetation = iota
)

func (g *game) Foliage(h, w int) map[position]vegetation {
	// use same structure as for the dungeon
	// walls will become foliage
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	d.Width = w
	d.Height = h
	for i := range d.Cells {
		r := RandInt(100)
		pos := d.CellPosition(i)
		if r >= 43 {
			d.SetCell(pos, WallCell)
		} else {
			d.SetCell(pos, FreeCell)
		}
	}
	for i := 0; i < 6; i++ {
		bufm := &dungeon{}
		bufm.Cells = make([]cell, h*w)
		bufm.Width = w
		bufm.Height = h
		copy(bufm.Cells, d.Cells)
		for j := range bufm.Cells {
			pos := d.CellPosition(j)
			c1 := d.WallAreaCount(pos, 1)
			if i < 4 {
				if c1 <= 4 {
					bufm.SetCell(pos, FreeCell)
				} else {
					bufm.SetCell(pos, WallCell)
				}
			}
			if i == 4 {
				if c1 > 6 {
					bufm.SetCell(pos, WallCell)
				}
			}
			if i == 5 {
				c2 := d.WallAreaCount(pos, 2)
				if c2 < 5 && c1 <= 2 {
					bufm.SetCell(pos, FreeCell)
				}
			}
		}
		d.Cells = bufm.Cells
	}
	fungus := make(map[position]vegetation)
	for i, c := range d.Cells {
		if c.T == FreeCell {
			fungus[d.CellPosition(i)] = foliage
		}
	}
	return fungus
}

func (g *game) DoorCandidate(pos position) bool {
	d := g.Dungeon
	if !d.Valid(pos) || d.Cell(pos).T != FreeCell {
		return false
	}
	return d.Valid(pos.W()) && d.Valid(pos.E()) &&
		d.Cell(pos.W()).T == FreeCell && d.Cell(pos.E()).T == FreeCell &&
		!g.Doors[pos.W()] && !g.Doors[pos.E()] &&
		(!d.Valid(pos.N()) || d.Cell(pos.N()).T == WallCell) &&
		(!d.Valid(pos.S()) || d.Cell(pos.S()).T == WallCell) &&
		((d.Valid(pos.NW()) && d.Cell(pos.NW()).T == FreeCell) ||
			(d.Valid(pos.SW()) && d.Cell(pos.SW()).T == FreeCell) ||
			(d.Valid(pos.NE()) && d.Cell(pos.NE()).T == FreeCell) ||
			(d.Valid(pos.SE()) && d.Cell(pos.SE()).T == FreeCell)) ||
		d.Valid(pos.N()) && d.Valid(pos.S()) &&
			d.Cell(pos.N()).T == FreeCell && d.Cell(pos.S()).T == FreeCell &&
			!g.Doors[pos.N()] && !g.Doors[pos.S()] &&
			(!d.Valid(pos.E()) || d.Cell(pos.E()).T == WallCell) &&
			(!d.Valid(pos.W()) || d.Cell(pos.W()).T == WallCell) &&
			((d.Valid(pos.NW()) && d.Cell(pos.NW()).T == FreeCell) ||
				(d.Valid(pos.SW()) && d.Cell(pos.SW()).T == FreeCell) ||
				(d.Valid(pos.NE()) && d.Cell(pos.NE()).T == FreeCell) ||
				(d.Valid(pos.SE()) && d.Cell(pos.SE()).T == FreeCell))
}

func (g *game) PutDoors(percentage int) {
	g.Doors = map[position]bool{}
	for i := range g.Dungeon.Cells {
		pos := g.Dungeon.CellPosition(i)
		if g.DoorCandidate(pos) && RandInt(100) < percentage {
			g.Doors[pos] = true
		}
	}
}
