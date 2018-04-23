// many ideas here from articles found at http://www.roguebasin.com/

package main

import (
	"sort"
)

type dungeon struct {
	Cells []cell
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
	return d.Cells[pos.idx()]
}

func (d *dungeon) Border(pos position) bool {
	return pos.X == DungeonWidth-1 || pos.Y == DungeonHeight-1 || pos.X == 0 || pos.Y == 0
}

func (d *dungeon) SetCell(pos position, t terrain) {
	d.Cells[pos.idx()].T = t
}

func (d *dungeon) SetExplored(pos position) {
	d.Cells[pos.idx()].Explored = true
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

func (d *dungeon) Area(area []position, pos position, radius int) []position {
	area = area[:0]
	for x := pos.X - radius; x <= pos.X+radius; x++ {
		for y := pos.Y - radius; y <= pos.Y+radius; y++ {
			pos := position{x, y}
			if pos.valid() {
				area = append(area, pos)
			}
		}
	}
	return area
}

func (d *dungeon) ConnectRoomsShortestPath(r1, r2 room) {
	mp := &dungeonPath{dungeon: d}
	path, _, _ := AstarPath(mp, r1.pos, r2.pos)
	for _, pos := range path {
		d.SetCell(pos, FreeCell)
	}
}

func (d *dungeon) DigRoom(r room) {
	for i := r.pos.X; i < r.pos.X+r.w; i++ {
		for j := r.pos.Y; j < r.pos.Y+r.h; j++ {
			rpos := position{i, j}
			if rpos.valid() {
				d.SetCell(rpos, FreeCell)
			}
		}
	}
}

func (d *dungeon) IsAreaFree(pos position, h, w int) bool {
	for i := pos.X; i < pos.X+w; i++ {
		for j := pos.Y; j < pos.Y+h; j++ {
			rpos := position{i, j}
			if !rpos.valid() || d.Cell(rpos).T != FreeCell {
				return false
			}
		}
	}
	return true
}

func (d *dungeon) BuildRoom(pos position, w, h int) bool {
	spos := position{pos.X - 1, pos.Y - 1}
	if !d.IsAreaFree(spos, h+2, w+2) {
		return false
	}
	for i := pos.X; i < pos.X+w; i++ {
		d.SetCell(position{i, pos.Y}, WallCell)
		d.SetCell(position{i, pos.Y + h - 1}, WallCell)
	}
	for i := pos.Y; i < pos.Y+h; i++ {
		d.SetCell(position{pos.X, i}, WallCell)
		d.SetCell(position{pos.X + w - 1, i}, WallCell)
	}
	doors := [4]position{
		position{pos.X + w/2, pos.Y},
		position{pos.X + w/2, pos.Y + h - 1},
		position{pos.X, pos.Y + h/2},
		position{pos.X + w - 1, pos.Y + h/2},
	}
	for i := 0; i < 4; i++ {
		d.SetCell(doors[RandInt(4)], FreeCell)
	}
	return true
}

func (d *dungeon) BuildSomeRoom(w, h int) bool {
	for i := 0; i < 100; i++ {
		pos := d.FreeCell()
		if d.BuildRoom(pos, w, h) {
			return true
		}
	}
	return false
}

func (g *game) GenRuinsMap(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	rooms := []room{}
	for i := 0; i < 45; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				pos: position{RandInt(w - 1), RandInt(h - 1)},
				w:   3 + RandInt(5),
				h:   2 + RandInt(3)}
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.DigRoom(ro)
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
	rooms := []room{}
	for i := 0; i < 35; i++ {
		var ro room
		count := 100
		for count > 0 {
			count--
			ro = room{
				pos: position{RandInt(w - 1), RandInt(h - 1)},
				w:   5 + RandInt(4),
				h:   3 + RandInt(3)}
			if !intersectsRoom(rooms, ro) {
				break
			}
		}

		d.DigRoom(ro)
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
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
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
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
	pos := position{40, 10}
	max := 21 * 42
	d.SetCell(pos, FreeCell)
	cells := 1
	notValid := 0
	lastValid := pos
	diag := RandInt(4) == 0
	for cells < max {
		npos := pos.RandomNeighbor(diag)
		if !pos.valid() && npos.valid() && d.Cell(npos).T == WallCell {
			pos = lastValid
			continue
		}
		pos = npos
		if pos.valid() {
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
	max = DungeonHeight * 1
	digs := 0
	i := 0
	block := make([]position, 0, 64)
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
		block = d.DigBlock(block, diag)
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
	//d.BuildSomeRoom(9, 7)
	g.Dungeon = d
	g.PutDoors(5)
}

func (d *dungeon) HasFreeNeighbor(pos position) bool {
	neighbors := pos.ValidNeighbors()
	for _, pos := range neighbors {
		if d.Cell(pos).T == FreeCell {
			return true
		}
	}
	return false
}

func (g *game) HasFreeExploredNeighbor(pos position) bool {
	d := g.Dungeon
	neighbors := pos.ValidNeighbors()
	for _, pos := range neighbors {
		c := d.Cell(pos)
		if c.T == FreeCell && c.Explored && !g.UnknownDig[pos] {
			return true
		}
	}
	return false
}

func (d *dungeon) DigBlock(block []position, diag bool) []position {
	pos := d.WallCell()
	block = block[:0]
	for {
		block = append(block, pos)
		if d.HasFreeNeighbor(pos) {
			break
		}
		pos = pos.RandomNeighbor(diag)
		if !pos.valid() {
			block = block[:0]
			pos = d.WallCell()
			continue
		}
		if !pos.valid() {
			return nil
		}
	}
	return block
}

func (g *game) GenCaveMapTree(h, w int) {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
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
	block := make([]position, 0, 64)
loop:
	for cells < max {
		block = d.DigBlock(block, diag)
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

func (d *dungeon) WallAreaCount(area []position, pos position, radius int) int {
	area = d.Area(area, pos, radius)
	count := 0
	for _, npos := range area {
		if d.Cell(npos).T == WallCell {
			count++
		}
	}
	switch radius {
	case 1:
		count += 9 - len(area)
	case 2:
		count += 25 - len(area)
	}
	return count
}

func (d *dungeon) Connected(pos position) (map[position]bool, int) {
	conn := map[position]bool{}
	stack := []position{pos}
	count := 0
	conn[pos] = true
	nb := make([]position, 0, 8)
	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		nb = pos.Neighbors(nb, func(npos position) bool {
			return npos.valid() && d.Cell(npos).T != WallCell
		})
		for _, npos := range nb {
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
		if c.T == FreeCell && !conn[idxtopos(i)] {
			return false
		}
	}
	return true
}

func (g *game) RunCellularAutomataCave(h, w int) bool {
	d := &dungeon{}
	d.Cells = make([]cell, h*w)
	for i := range d.Cells {
		r := RandInt(100)
		pos := idxtopos(i)
		if r >= 45 {
			d.SetCell(pos, FreeCell)
		} else {
			d.SetCell(pos, WallCell)
		}
	}
	bufm := &dungeon{}
	bufm.Cells = make([]cell, h*w)
	area := make([]position, 0, 25)
	for i := 0; i < 5; i++ {
		for j := range bufm.Cells {
			pos := idxtopos(j)
			c1 := d.WallAreaCount(area, pos, 1)
			if c1 >= 5 {
				bufm.SetCell(pos, WallCell)
			} else {
				bufm.SetCell(pos, FreeCell)
			}
			if i == 3 {
				c2 := d.WallAreaCount(area, pos, 2)
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
		if count >= 37*DungeonHeight*DungeonWidth/100 {
			break
		}
	}
	conn, count = d.Connected(winner)
	if count <= 37*DungeonHeight*DungeonWidth/100 {
		return false
	}
	for i, c := range d.Cells {
		pos := idxtopos(i)
		if c.T == FreeCell && !conn[pos] {
			d.SetCell(pos, WallCell)
		}
	}
	max := DungeonHeight * 1
	cells := 1
	digs := 0
	i := 0
	block := make([]position, 0, 64)
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
		block = d.DigBlock(block, diag)
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
	for i := range d.Cells {
		r := RandInt(100)
		pos := idxtopos(i)
		if r >= 43 {
			d.SetCell(pos, WallCell)
		} else {
			d.SetCell(pos, FreeCell)
		}
	}
	area := make([]position, 0, 25)
	for i := 0; i < 6; i++ {
		bufm := &dungeon{}
		bufm.Cells = make([]cell, h*w)
		copy(bufm.Cells, d.Cells)
		for j := range bufm.Cells {
			pos := idxtopos(j)
			c1 := d.WallAreaCount(area, pos, 1)
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
				c2 := d.WallAreaCount(area, pos, 2)
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
			fungus[idxtopos(i)] = foliage
		}
	}
	return fungus
}

func (g *game) DoorCandidate(pos position) bool {
	d := g.Dungeon
	if !pos.valid() || d.Cell(pos).T != FreeCell {
		return false
	}
	return pos.W().valid() && pos.E().valid() &&
		d.Cell(pos.W()).T == FreeCell && d.Cell(pos.E()).T == FreeCell &&
		!g.Doors[pos.W()] && !g.Doors[pos.E()] &&
		(!pos.N().valid() || d.Cell(pos.N()).T == WallCell) &&
		(!pos.S().valid() || d.Cell(pos.S()).T == WallCell) &&
		((pos.NW().valid() && d.Cell(pos.NW()).T == FreeCell) ||
			(pos.SW().valid() && d.Cell(pos.SW()).T == FreeCell) ||
			(pos.NE().valid() && d.Cell(pos.NE()).T == FreeCell) ||
			(pos.SE().valid() && d.Cell(pos.SE()).T == FreeCell)) ||
		pos.N().valid() && pos.S().valid() &&
			d.Cell(pos.N()).T == FreeCell && d.Cell(pos.S()).T == FreeCell &&
			!g.Doors[pos.N()] && !g.Doors[pos.S()] &&
			(!pos.E().valid() || d.Cell(pos.E()).T == WallCell) &&
			(!pos.W().valid() || d.Cell(pos.W()).T == WallCell) &&
			((pos.NW().valid() && d.Cell(pos.NW()).T == FreeCell) ||
				(pos.SW().valid() && d.Cell(pos.SW()).T == FreeCell) ||
				(pos.NE().valid() && d.Cell(pos.NE()).T == FreeCell) ||
				(pos.SE().valid() && d.Cell(pos.SE()).T == FreeCell))
}

func (g *game) PutDoors(percentage int) {
	g.Doors = map[position]bool{}
	for i := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if g.DoorCandidate(pos) && RandInt(100) < percentage {
			g.Doors[pos] = true
		}
	}
}
