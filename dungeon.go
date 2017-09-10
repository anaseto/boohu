// many ideas here from articles found at http://www.roguebasin.com/

package main

import (
	"sort"
)

type dungeon struct {
	Cells  []cell
	Width  int
	Heigth int
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

func (m *dungeon) Cell(pos position) cell {
	return m.Cells[pos.Y*m.Width+pos.X]
}

func (m *dungeon) Valid(pos position) bool {
	return pos.X < m.Width && pos.Y < m.Heigth && pos.X >= 0 && pos.Y >= 0
}

func (m *dungeon) SetCell(pos position, t terrain) {
	m.Cells[pos.Y*m.Width+pos.X].T = t
}

func (m *dungeon) SetExplored(pos position) {
	m.Cells[pos.Y*m.Width+pos.X].Explored = true
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

func (m *dungeon) connectRooms(r1, r2 room) {
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
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X {
			x--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y < r2.pos.Y {
			y++
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y > r2.pos.Y {
			y--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		break
	}
}

func (m *dungeon) connectRoomsDiagonally(r1, r2 room) {
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
			m.SetCell(position{x, y}, FreeCell)
			y++
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X && y < r2.pos.Y {
			x--
			m.SetCell(position{x, y}, FreeCell)
			y++
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X && y > r2.pos.Y {
			x--
			m.SetCell(position{x, y}, FreeCell)
			y--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x < r2.pos.X && y > r2.pos.Y {
			x++
			m.SetCell(position{x, y}, FreeCell)
			y--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x < r2.pos.X {
			x++
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if x > r2.pos.X {
			x--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y < r2.pos.Y {
			y++
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		if y > r2.pos.Y {
			y--
			m.SetCell(position{x, y}, FreeCell)
			continue
		}
		break
	}
}

func (d *dungeon) Neighbors(pos position) []position {
	neighbors := [8]position{pos.E(), pos.W(), pos.N(), pos.S(), pos.NE(), pos.NW(), pos.SE(), pos.SW()}
	validNeighbors := []position{}
	for _, pos := range neighbors {
		if d.Valid(pos) {
			validNeighbors = append(validNeighbors, pos)
		}
	}
	return validNeighbors
}

func (d *dungeon) CardinalNeighbors(pos position) []position {
	neighbors := [4]position{pos.E(), pos.W(), pos.N(), pos.S()}
	validNeighbors := []position{}
	for _, pos := range neighbors {
		if d.Valid(pos) {
			validNeighbors = append(validNeighbors, pos)
		}
	}
	return validNeighbors
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
	} else {
		return 1
	}
}

func (dp *dungeonPath) Estimation(from, to position) int {
	return from.Distance(to)
}

func (m *dungeon) ConnectRoomsShortestPath(r1, r2 room) error {
	mp := &dungeonPath{dungeon: m}
	path, _, _ := AstarPath(mp, r1.pos, r2.pos)
	for _, pos := range path {
		m.SetCell(pos, FreeCell)
	}
	return nil
}

func (m *dungeon) PutRoom(r room) {
	for i := r.pos.X; i < r.pos.X+r.w; i++ {
		for j := r.pos.Y; j < r.pos.Y+r.h; j++ {
			if m.Valid(position{i, j}) {
				m.SetCell(position{i, j}, FreeCell)
			}
		}
	}
}

func (m *dungeon) CellPosition(i int) position {
	return position{i - (i/m.Width)*m.Width, i / m.Width}
}

func (g *game) GenRuinsMap(h, w int) {
	m := &dungeon{}
	m.Cells = make([]cell, h*w)
	m.Width = w
	m.Heigth = h
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

		m.PutRoom(ro)
		if len(rooms) > 0 {
			r := RandInt(100)
			if r > 75 {
				m.connectRooms(nearRoom(rooms, ro), ro)
			} else if r > 25 {
				m.ConnectRoomsShortestPath(nearRoom(rooms, ro), ro)
			} else {
				m.connectRoomsDiagonally(nearRoom(rooms, ro), ro)
			}
		}
		rooms = append(rooms, ro)
	}
	g.Dungeon = m
}

type roomSlice []room

func (rs roomSlice) Len() int      { return len(rs) }
func (rs roomSlice) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs roomSlice) Less(i, j int) bool {
	return rs[i].pos.Y < rs[j].pos.Y || rs[i].pos.Y == rs[j].pos.Y && rs[i].pos.X < rs[j].pos.X
}

func (g *game) GenRoomMap(h, w int) {
	m := &dungeon{}
	m.Cells = make([]cell, h*w)
	m.Width = w
	m.Heigth = h
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

		m.PutRoom(ro)
		rooms = append(rooms, ro)
	}
	sort.Sort(roomSlice(rooms))
	for i, ro := range rooms {
		if i == 0 {
			continue
		}
		r := RandInt(100)
		if r > 50 {
			m.connectRooms(nearestRoom(rooms[:i], ro), ro)
		} else if r > 25 {
			m.ConnectRoomsShortestPath(nearRoom(rooms[:i], ro), ro)
		} else {
			m.connectRoomsDiagonally(nearestRoom(rooms[:i], ro), ro)
		}
	}
	g.Dungeon = m
}

func (m *dungeon) FreeCell() position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			return pos
		}
	}
}

func (m *dungeon) WallCell() position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("WallCell")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == WallCell {
			return pos
		}
	}
}

func (g *game) GenCaveMap(h, w int) {
	m := &dungeon{}
	m.Cells = make([]cell, h*w)
	m.Width = w
	m.Heigth = h
	pos := position{40, 10}
	max := 21 * 40
	m.SetCell(pos, FreeCell)
	cells := 1
	notValid := 0
	lastValid := pos
	diag := RandInt(4) == 0
	for cells < max {
		npos := pos.RandomNeighbor(diag)
		if !m.Valid(pos) && m.Valid(npos) && m.Cell(npos).T == WallCell {
			pos = lastValid
			continue
		}
		pos = npos
		if m.Valid(pos) {
			if m.Cell(pos).T != FreeCell {
				m.SetCell(pos, FreeCell)
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
	g.Dungeon = m
}

func (g *game) GenCaveMapTree(h, w int) {
	m := &dungeon{}
	m.Cells = make([]cell, h*w)
	m.Width = w
	m.Heigth = h
	center := position{40, 10}
	m.SetCell(center, FreeCell)
	m.SetCell(center.E(), FreeCell)
	m.SetCell(center.NE(), FreeCell)
	m.SetCell(center.S(), FreeCell)
	m.SetCell(center.SE(), FreeCell)
	m.SetCell(center.N(), FreeCell)
	m.SetCell(center.NW(), FreeCell)
	m.SetCell(center.W(), FreeCell)
	m.SetCell(center.SW(), FreeCell)
	max := 21 * 50
	cells := 1
	diag := RandInt(2) == 0
loop:
	for cells < max {
		pos := m.WallCell()
		block := []position{}
		for {
			block = append(block, pos)
			pos = pos.RandomNeighbor(diag)
			if m.Valid(pos) {
				if m.Cell(pos).T == FreeCell {
					break
				}
			} else {
				continue loop
			}
		}
		for _, pos := range block {
			m.SetCell(pos, FreeCell)
			cells++
		}
	}
	g.Dungeon = m
}

func (d *dungeon) WallNeighborsCount(pos position) int {
	neighbors := d.Neighbors(pos)
	count := 0
	for _, pos := range neighbors {
		if d.Cell(pos).T == WallCell {
			count++
		}
	}
	return count
}

func (d *dungeon) WallAreaCount(pos position, radius int) int {
	neighbors := d.Area(pos, radius)
	count := 0
	for _, pos := range neighbors {
		if d.Cell(pos).T == WallCell {
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
		for _, n := range neighbors {
			if !conn[n] {
				conn[n] = true
				stack = append(stack, n)
			}
		}
	}
	return conn, count
}

func (g *game) RunCellularAutomataCave(h, w int) bool {
	m := &dungeon{}
	m.Cells = make([]cell, h*w)
	m.Width = w
	m.Heigth = h
	for i, _ := range m.Cells {
		r := RandInt(100)
		pos := m.CellPosition(i)
		if r >= 45 {
			m.SetCell(pos, FreeCell)
		} else {
			m.SetCell(pos, WallCell)
		}
	}
	for i := 0; i < 5; i++ {
		bufm := &dungeon{}
		bufm.Cells = make([]cell, h*w)
		bufm.Width = w
		bufm.Heigth = h
		copy(bufm.Cells, m.Cells)
		for j, _ := range bufm.Cells {
			pos := m.CellPosition(j)
			c1 := m.WallAreaCount(pos, 1)
			if c1 >= 5 {
				bufm.SetCell(pos, WallCell)
			} else {
				bufm.SetCell(pos, FreeCell)
			}
			if i == 3 {
				c2 := m.WallAreaCount(pos, 2)
				if c2 <= 2 {
					bufm.SetCell(pos, WallCell)
				}
			}
		}
		m.Cells = bufm.Cells
	}
	var conn map[position]bool
	var count int
	var winner position
	for i := 0; i < 15; i++ {
		pos := m.FreeCell()
		if conn[pos] {
			continue
		}
		var ncount int
		conn, ncount = m.Connected(pos)
		if ncount > count {
			count = ncount
			winner = pos
		}
		if count >= 37*m.Heigth*m.Width/100 {
			break
		}
	}
	conn, count = m.Connected(winner)
	if count <= 37*m.Heigth*m.Width/100 {
		return false
	}
	for i, c := range m.Cells {
		pos := m.CellPosition(i)
		if c.T == FreeCell && !conn[pos] {
			m.SetCell(pos, WallCell)
		}
	}
	g.Dungeon = m
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
