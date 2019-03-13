// many ideas here from articles found at http://www.roguebasin.com/

package main

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

type dungeon struct {
	Cells []cell
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
		if c.IsFree() {
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

func (d *dungeon) HasFreeNeighbor(pos position) bool {
	neighbors := pos.ValidCardinalNeighbors()
	for _, pos := range neighbors {
		if d.Cell(pos).IsFree() {
			return true
		}
	}
	return false
}

func (g *game) HasFreeExploredNeighbor(pos position) bool {
	d := g.Dungeon
	neighbors := pos.ValidCardinalNeighbors()
	for _, pos := range neighbors {
		c := d.Cell(pos)
		if t, ok := g.TerrainKnowledge[pos]; ok {
			c.T = t
		}
		if c.IsFree() && c.Explored {
			return true
		}
	}
	return false
}

func (d *dungeon) Connected(pos position, nf func(position) bool) (map[position]bool, int) {
	conn := map[position]bool{}
	stack := []position{pos}
	count := 0
	conn[pos] = true
	nb := make([]position, 0, 8)
	for len(stack) > 0 {
		pos = stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		count++
		nb = pos.CardinalNeighbors(nb, nf)
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
	conn, _ := d.Connected(pos, d.IsFreeCell)
	for i, c := range d.Cells {
		if c.IsFree() && !conn[idxtopos(i)] {
			return false
		}
	}
	return true
}

func (d *dungeon) IsAreaFree(pos position, h, w int) bool {
	for i := pos.X; i < pos.X+w; i++ {
		for j := pos.Y; j < pos.Y+h; j++ {
			rpos := position{i, j}
			if !rpos.valid() || d.Cell(rpos).IsFree() {
				return false
			}
		}
	}
	return true
}

func (d *dungeon) IsAreaWall(pos position, h, w int) bool {
	for i := pos.X; i < pos.X+w; i++ {
		for j := pos.Y; j < pos.Y+h; j++ {
			rpos := position{i, j}
			if rpos.valid() && d.Cell(rpos).T != WallCell {
				return false
			}
		}
	}
	return true
}

type rentry struct {
	pos     position
	used    bool
	virtual bool
}

type placeKind int

const (
	PlaceDoor placeKind = iota
	PlacePatrol
	PlaceStatic
	PlaceSpecialStatic
	PlaceItem
)

type place struct {
	pos  position
	kind placeKind
	used bool
}

type room struct {
	pos     position
	w       int
	h       int
	entries []rentry
	places  []place
	kind    string
}

func roomDistance(r1, r2 *room) int {
	// TODO: use the center?
	return Abs(r1.pos.X-r2.pos.X) + Abs(r1.pos.Y-r2.pos.Y)
}

type roomSlice []*room

func (rs roomSlice) Len() int      { return len(rs) }
func (rs roomSlice) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }
func (rs roomSlice) Less(i, j int) bool {
	//return rs[i].pos.Y < rs[j].pos.Y || rs[i].pos.Y == rs[j].pos.Y && rs[i].pos.X < rs[j].pos.X
	center := position{DungeonWidth / 2, DungeonHeight / 2}
	ipos := rs[i].pos
	ipos.X += rs[i].w / 2
	ipos.Y += rs[i].h / 2
	jpos := rs[j].pos
	jpos.X += rs[j].w / 2
	jpos.Y += rs[j].h / 2
	return ipos.Distance(center) <= jpos.Distance(center)
}

type dgen struct {
	d      *dungeon
	tunnel map[position]bool
	room   map[position]bool
	rooms  []*room
	fungus map[position]vegetation
}

func (dg *dgen) WallAreaCount(area []position, pos position, radius int) int {
	d := dg.d
	area = d.Area(area, pos, radius)
	count := 0
	for _, npos := range area {
		if d.Cell(npos).T == WallCell || dg.tunnel[npos] {
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

// UnusedEntry returns an unused entry, if possible, or a random entry
// otherwise.
func (r *room) UnusedEntry() int {
	ens := []int{}
	for i, e := range r.entries {
		if !e.used {
			ens = append(ens, i)
		}
	}
	if len(ens) == 0 {
		return RandInt(len(r.entries))
	}
	return ens[RandInt(len(ens))]
}

func (dg *dgen) ConnectRoomsShortestPath(i, j int) bool {
	if i == j {
		return false
	}
	r1 := dg.rooms[i]
	r2 := dg.rooms[j]
	// TODO: more versatile and hand-made doors locations
	var e1pos, e2pos position
	var e1i, e2i int
	e1i = r1.UnusedEntry()
	e1pos = r1.entries[e1i].pos
	e2i = r2.UnusedEntry()
	e2pos = r2.entries[e2i].pos
	tp := &tunnelPath{dg: dg}
	path, _, found := AstarPath(tp, e1pos, e2pos)
	if !found {
		log.Println(fmt.Sprintf("no path from %v to %v", e1pos, e2pos))
		return false
	}
	for _, pos := range path {
		if !pos.valid() {
			panic(fmt.Sprintf("position %v from %v to %v", pos, e1pos, e2pos))
		}
		if dg.d.Cell(pos).T == WallCell {
			dg.d.SetCell(pos, GroundCell)
			dg.tunnel[pos] = true
		}
	}
	r1.entries[e1i].used = true
	r2.entries[e2i].used = true
	return true
}

const (
	RoomSquare = `
?###+###?
#_..!.._#
+..!P!..+
#_..!.._#
?###+###?`
	RoomLittle = `
?#+#?
#_._#
+.P.+
#_._#
?#+#?`
	RoomLittleDiamond = `
??#+#??
##_._##
+..P..+
##_._##
??#+#??`
	RoomLittleColumnDiamond = `
??#+#??
##_._##
+.P#P.+
##_._##
??#+#??`
	RoomRound = `
???##+##???
??#".P."#??
##"._#_."##
+.P.###.P.+
##"._#_."##
??#".P."#??
???##+##???`
)

var roomNormalTemplates = []string{RoomSquare, RoomLittle, RoomLittleDiamond, RoomLittleColumnDiamond, RoomRound}

const (
	RoomBigColumns = `
?####?#++#?####?
#!..>##..##>..!#
##.P........P.##
+...._####_....+
##.P........P.##
#!..>##..##>..!#
?####?#++#?####?`
	RoomBigGarden = `
?####?#++#?####?
#""""##..##""""#
#""""""!!""""""#
#""""""P>""""""#
#""""""!!""""""#
#""""##..##""""#
?####?#++#?####?`
	RoomBigRooms = `
?####?#++#?####?
#>..!##..##!..>#
#"""..#..#.."""#
#"""P.|..|.P"""#
#"""..#..#.."""#
#>..!##..##!..>#
?####?#++#?####?`
	RoomColumns = `
###+##+###
+..P..P..+
#.#>##!#.#
#.#!##>#.#
+..P..P..+
###+##+###`
	RoomHome1 = `
?##########+#?
#>..P...|..P.#
#......!#!!..#
####|#######|#
#....P....|..#
#>.......!#P.#
?##########+#?
`
	RoomHome2 = `
?############?
+...#.......>#
#.P.|....P...#
##|###!.....!#
#...>##|######
#!P..|...P...+
?##########+#?
`
	RoomHome3 = `
?###############?
#>....|.........#
#..P.!##|##!.P..+
######!...!#....#
+....|.P>._###|##
######!...!#!...#
#!...>##|##..P..+
#..P......|.....#
?###############?
`
	RoomCaban = `
???????-??????
?????""""?????
???""""""""???
??"""###."""??
?"""#>!|P."""?
-""""###.""""-
??""""""""""??
????"""""?????
???????-??????`
	RoomRoundColumns = `
???##+##???
??#_..._#??
##!.#P#.!##
+...P>P...+
##!.#P#.!##
??#_..._#??
???##+##???`
	RoomRoundGarden = `
???##+##???
??#>.P.>#??
##!.""".!##
+.P"""""P.+
##!.""".!##
??#>.P.>#??
???##+##???`
	RoomLongHall = `
####################
+.P..............P.+
#..!##!>.##.>!##!..#
+.P..............P.+
####################`
	RoomGardenHall = `
?##################?
#"""""".>!>."""""""#
+....P........P....+
#"""""".>!>."""""""#
?##################?`
)

var roomSpecialTemplates = []string{RoomBigColumns, RoomBigGarden, RoomColumns, RoomRoundColumns, RoomRoundGarden, RoomLongHall,
	RoomGardenHall, RoomHome1, RoomHome2, RoomHome3, RoomBigRooms, RoomCaban}

func (r *room) ComputeDimensions() {
	x := 0
	y := 0
	for _, c := range r.kind {
		if c == '\n' {
			if x > r.w {
				r.w = x
			}
			x = 0
			y++
		}
		x++
	}
	r.h = y + 1
}

func (r *room) HasSpace(dg *dgen) bool {
	if DungeonWidth-r.pos.X < r.w || DungeonHeight-r.pos.Y < r.h {
		return false
	}
	for i := r.pos.X - 1; i <= r.pos.X+r.w; i++ {
		for j := r.pos.Y - 1; j <= r.pos.Y+r.h; j++ {
			rpos := position{i, j}
			if rpos.valid() && dg.room[rpos] {
				return false
			}
		}
	}
	return true
}

func (r *room) Dig(dg *dgen) {
	x := 0
	y := 0
	for _, c := range r.kind {
		if c == '\n' {
			x = 0
			y++
			continue
		}
		pos := position{X: r.pos.X + x, Y: r.pos.Y + y}
		if pos.valid() && c != '?' {
			dg.room[pos] = true
		}
		switch c {
		case '.', '>', '!', 'P', '_', '|':
			if pos.valid() {
				dg.d.SetCell(pos, GroundCell)
			}
		case '#', '+':
			if pos.valid() {
				dg.d.SetCell(pos, WallCell)
			}
		}
		switch c {
		case '>':
			r.places = append(r.places, place{pos: pos, kind: PlaceSpecialStatic})
		case '!':
			r.places = append(r.places, place{pos: pos, kind: PlaceItem})
		case 'P':
			r.places = append(r.places, place{pos: pos, kind: PlacePatrol})
		case '_':
			r.places = append(r.places, place{pos: pos, kind: PlaceStatic})
		case '|':
			r.places = append(r.places, place{pos: pos, kind: PlaceDoor})
		case '+', '-':
			if pos.X == 0 || pos.X == DungeonWidth-1 || pos.Y == 0 || pos.Y == DungeonHeight-1 {
				break
			}
			e := rentry{}
			e.pos = pos
			if c == '-' {
				e.virtual = true
			}
			r.entries = append(r.entries, e)
		case '"':
			if pos.valid() {
				dg.d.SetCell(pos, FungusCell)
			}
		}
		if c != '"' {
			delete(dg.fungus, pos)
		}
		x++
	}
}

func (dg *dgen) NewRoom(rpos position, kind string) *room {
	r := &room{pos: rpos, kind: kind}
	r.kind = strings.TrimSpace(r.kind)
	r.ComputeDimensions()
	if !r.HasSpace(dg) {
		return nil
	}
	r.Dig(dg)
	return r
}

func (dg *dgen) nearestRoom(i int) (k int) {
	r := dg.rooms[i]
	d := roomDistance(r, dg.rooms[k])
	for j, nextRoom := range dg.rooms[i:] {
		nd := roomDistance(r, nextRoom)
		if nd < d {
			n := RandInt(15)
			if n > 0 {
				d = nd
				k = j
			}
		}
	}
	return k
}

func (dg *dgen) PutDoors(g *game) {
	for _, r := range dg.rooms {
		for _, e := range r.entries {
			//if e.used && g.DoorCandidate(e.pos) {
			if e.used && !e.virtual {
				r.places = append(r.places, place{pos: e.pos, kind: PlaceDoor})
			}
		}
		for _, pl := range r.places {
			if pl.kind != PlaceDoor {
				continue
			}
			dg.d.SetCell(pl.pos, DoorCell)
			r.places = append(r.places, place{pos: pl.pos, kind: PlaceDoor})
			if _, ok := dg.fungus[pl.pos]; ok {
				delete(dg.fungus, pl.pos)
			}
		}
	}
}

func (g *game) DoorCandidate(pos position) bool {
	d := g.Dungeon
	if !pos.valid() || d.Cell(pos).IsFree() {
		return false
	}
	return pos.W().valid() && pos.E().valid() &&
		d.Cell(pos.W()).IsGround() && d.Cell(pos.E()).IsGround() &&
		(!pos.N().valid() || d.Cell(pos.N()).T == WallCell) &&
		(!pos.S().valid() || d.Cell(pos.S()).T == WallCell) &&
		((pos.NW().valid() && d.Cell(pos.NW()).IsFree()) ||
			(pos.SW().valid() && d.Cell(pos.SW()).IsFree()) ||
			(pos.NE().valid() && d.Cell(pos.NE()).IsFree()) ||
			(pos.SE().valid() && d.Cell(pos.SE()).IsFree())) ||
		pos.N().valid() && pos.S().valid() &&
			d.Cell(pos.N()).IsGround() && d.Cell(pos.S()).IsGround() &&
			(!pos.E().valid() || d.Cell(pos.E()).T == WallCell) &&
			(!pos.W().valid() || d.Cell(pos.W()).T == WallCell) &&
			((pos.NW().valid() && d.Cell(pos.NW()).IsFree()) ||
				(pos.SW().valid() && d.Cell(pos.SW()).IsFree()) ||
				(pos.NE().valid() && d.Cell(pos.NE()).IsFree()) ||
				(pos.SE().valid() && d.Cell(pos.SE()).IsFree()))
}

func (dg *dgen) GenRooms(templates []string, n int) {
	for i := 0; i < n; i++ {
		var r *room
		count := 100
		for r == nil && count > 0 {
			count--
			r = dg.NewRoom(position{RandInt(DungeonWidth - 1), RandInt(DungeonHeight - 1)}, templates[RandInt(len(templates))])
		}
		if r != nil {
			dg.rooms = append(dg.rooms, r)
		}
	}
}

func (dg *dgen) ConnectRooms() {
	sort.Sort(roomSlice(dg.rooms))
	for i := range dg.rooms {
		if i == 0 {
			continue
		}
		ok := dg.ConnectRoomsShortestPath(dg.nearestRoom(i), i)
		for !ok {
			ok = dg.ConnectRoomsShortestPath(RandInt(len(dg.rooms)), i)
		}
	}
	for i := 0; i < 4; i++ {
		j := RandInt(len(dg.rooms))
		k := RandInt(len(dg.rooms))
		if j != k {
			dg.ConnectRoomsShortestPath(j, k)
		}
	}
}

func (g *game) GenRoomTunnels() {
	dg := dgen{}
	d := &dungeon{}
	d.Cells = make([]cell, DungeonNCells)
	dg.d = d
	dg.tunnel = make(map[position]bool)
	dg.room = make(map[position]bool)
	dg.rooms = []*room{}
	dg.fungus = make(map[position]vegetation)
	dg.GenCellularAutomataCaveMap()
	dg.GenRooms(roomSpecialTemplates, 3)
	dg.GenRooms(roomNormalTemplates, 7)
	dg.ConnectRooms()
	g.Dungeon = d
	dg.PutDoors(g)
	dg.PlayerStartCell(g)
	dg.ClearUnconnected(g)
	g.Objects.Stairs = map[position]stair{}
	if g.Depth < MaxDepth {
		dg.Stairs(g, NormalStair)
	}
	if g.Depth == WinDepth || g.Depth == MaxDepth {
		dg.Stairs(g, WinStair)
	}
	for i := 0; i < 4+RandInt(2); i++ {
		dg.Barrel(g)
	}
	dg.GenMonsters(g)
	dg.AddSpecial(g)
}

func (dg *dgen) ClearUnconnected(g *game) {
	d := dg.d
	conn, _ := d.Connected(g.Player.Pos, d.IsFreeCell)
	for i, c := range d.Cells {
		pos := idxtopos(i)
		if c.IsFree() && !conn[pos] {
			d.SetCell(pos, WallCell)
		}
	}
}

func (dg *dgen) AddSpecial(g *game) {
	// Equipment
	switch g.GenPlan[g.Depth] {
	//case GenWeapon:
	//g.GenWeapon()
	//case GenArmour:
	//g.GenArmour()
	//case GenWpArm:
	//g.GenWeapon()
	//g.GenArmour()
	case GenRod:
		//g.GenerateRod()
	case GenExtraCollectables:
		//for i := 0; i < 2; i++ {
		//dg.GenCollectable(g)
		//g.CollectableScore-- // these are extra
		//}
	}
	if g.Depth == 1 {
		// extra collectable
		//dg.GenCollectable(g)
		//g.CollectableScore--
	}

}

func (r *room) RandomPlace(kind placeKind) position {
	var p []int
	for i, pl := range r.places {
		if pl.kind == kind && (!pl.used || RandInt(4) == 0) {
			p = append(p, i)
		}
	}
	if len(p) == 0 {
		return InvalidPos
	}
	j := p[RandInt(len(p))]
	r.places[j].used = true
	return r.places[j].pos
}

func (dg *dgen) PlayerStartCell(g *game) {
	g.Player.Pos = dg.rooms[len(dg.rooms)-1].RandomPlace(PlacePatrol)
}

func (dg *dgen) Stairs(g *game, st stair) {
	var ri, pj int
	best := 0
	for i, r := range dg.rooms {
		for j, pl := range r.places {
			if !pl.used && pl.kind == PlaceSpecialStatic && pl.pos.Distance(g.Player.Pos) > best && (RandInt(3) == 0 || best == 0) {
				ri = i
				pj = j
				best = pl.pos.Distance(g.Player.Pos)
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	pos := r.places[pj].pos
	g.Dungeon.SetCell(pos, StairCell)
	g.Objects.Stairs[pos] = st
}

func (dg *dgen) Barrel(g *game) {
	var ri, pj int
	best := 0
	for i, r := range dg.rooms {
		for j, pl := range r.places {
			n := RandInt(100)
			if !pl.used && pl.kind == PlaceSpecialStatic && n >= best {
				ri = i
				pj = j
				best = n
			}
		}
	}
	r := dg.rooms[ri]
	r.places[pj].used = true
	r.places[pj].used = true
	g.Dungeon.SetCell(r.places[pj].pos, BarrelCell)
}

type vegetation int

const (
	foliage vegetation = iota
)

func (dg *dgen) GenCellularAutomataCaveMap() {
	count := 0
	for {
		count++
		if count > 100 {
			panic("genCellularAutomataCaveMap")
		}
		if dg.RunCellularAutomataCave() {
			break
		}
		// refresh cells
		dg.d.Cells = make([]cell, DungeonNCells)
	}
	dg.Foliage()
}

func (dg *dgen) RunCellularAutomataCave() bool {
	d := dg.d // TODO: reset
	for i := range d.Cells {
		r := RandInt(100)
		pos := idxtopos(i)
		if r >= 45 {
			d.SetCell(pos, GroundCell)
		} else {
			d.SetCell(pos, WallCell)
		}
	}
	bufm := &dungeon{}
	bufm.Cells = make([]cell, DungeonNCells)
	area := make([]position, 0, 25)
	for i := 0; i < 5; i++ {
		for j := range bufm.Cells {
			pos := idxtopos(j)
			c1 := d.WallAreaCount(area, pos, 1)
			if c1 >= 5 {
				bufm.SetCell(pos, WallCell)
			} else {
				bufm.SetCell(pos, GroundCell)
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
	return true
}

func (dg *dgen) Foliage() {
	// use same structure as for the dungeon
	// walls will become foliage
	d := &dungeon{}
	d.Cells = make([]cell, DungeonNCells)
	for i := range d.Cells {
		r := RandInt(100)
		pos := idxtopos(i)
		if r >= 43 {
			d.SetCell(pos, WallCell)
		} else {
			d.SetCell(pos, GroundCell)
		}
	}
	area := make([]position, 0, 25)
	for i := 0; i < 6; i++ {
		bufm := &dungeon{}
		bufm.Cells = make([]cell, DungeonNCells)
		copy(bufm.Cells, d.Cells)
		for j := range bufm.Cells {
			pos := idxtopos(j)
			c1 := d.WallAreaCount(area, pos, 1)
			if i < 4 {
				if c1 <= 4 {
					bufm.SetCell(pos, GroundCell)
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
					bufm.SetCell(pos, GroundCell)
				}
			}
		}
		d.Cells = bufm.Cells
	}
	for i, c := range d.Cells {
		if c.T == GroundCell {
			dg.d.SetCell(idxtopos(i), FungusCell)
		}
	}
}
