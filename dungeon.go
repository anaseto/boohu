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

func (d *dungeon) FreePassableCell() position {
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
		if c.IsPassable() {
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
		if d.Cell(pos).IsPassable() {
			return true
		}
	}
	return false
}

func (d *dungeon) HasTooManyWallNeighbors(pos position) bool {
	neighbors := pos.ValidNeighbors()
	count := 0
	for _, pos := range neighbors {
		if !d.Cell(pos).IsPassable() {
			count++
		}
	}
	return count > 1
}

func (g *game) HasFreeExploredNeighbor(pos position) bool {
	d := g.Dungeon
	neighbors := pos.ValidCardinalNeighbors()
	for _, pos := range neighbors {
		c := d.Cell(pos)
		if t, ok := g.TerrainKnowledge[pos]; ok {
			c.T = t
		}
		if c.IsPassable() && c.Explored {
			return true
		}
	}
	return false
}

func (dg *dgen) ComputeConnectedComponents(nf func(position) bool) {
	dg.cc = make([]int, DungeonNCells)
	index := 1
	stack := []position{}
	nb := make([]position, 0, 8)
	for i, _ := range dg.d.Cells {
		pos := idxtopos(i)
		if dg.cc[i] != 0 || !nf(pos) {
			continue
		}
		stack = append(stack[:0], pos)
		count := 0
		dg.cc[i] = index
		for len(stack) > 0 {
			pos = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			count++
			nb = pos.CardinalNeighbors(nb, nf)
			for _, npos := range nb {
				if dg.cc[npos.idx()] != index {
					dg.cc[npos.idx()] = index
					stack = append(stack, npos)
				}
			}
		}
	}
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
	pos := d.FreePassableCell()
	conn, _ := d.Connected(pos, d.NotWallCell)
	for i, c := range d.Cells {
		if c.IsPassable() && !conn[idxtopos(i)] {
			return false
		}
	}
	return true
}

func (d *dungeon) IsAreaFree(pos position, h, w int) bool {
	for i := pos.X; i < pos.X+w; i++ {
		for j := pos.Y; j < pos.Y+h; j++ {
			rpos := position{i, j}
			if !rpos.valid() || d.Cell(rpos).IsPassable() {
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
	PlaceStory
	PlacePatrolSpecial
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
	d       *dungeon
	tunnel  map[position]bool
	room    map[position]bool
	rooms   []*room
	fungus  map[position]vegetation
	spl     places
	special specialRoom
	layout  maplayout
	cc      []int
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
		t := dg.d.Cell(pos).T
		if t == WallCell || t == ChasmCell {
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
	RoomRoundSimple = `
??#+#??
?#!.!#?
#_.P._#
+.P#P.+
#_.P._#
?#!.!#?
??#+#??`
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
?##+#??
#_..!##
+.P#P.+
##!.._#
??#+##?`
	RoomLittleTreeDiamond = `
??#+##?
##!.._#
+.PTP.+
#_..!##
?##+#??`
	RoomRound = `
???##+##???
??#".P."#??
##"._#_."##
+.P.###.P.+
##"._#_."##
??#".P."#??
???##+##???`
	RoomRoundTree = `
???##+##???
??#".P."#??
##"._._."##
+.P..T..P.+
##"._._."##
??#".P."#??
???##+##???`
)

var roomNormalTemplates = []string{RoomSquare, RoomRoundSimple, RoomLittle, RoomLittleDiamond, RoomLittleColumnDiamond, RoomRound, RoomLittleTreeDiamond, RoomRoundTree}

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
#"T""".!P."""T"#
#""""">P_>"""""#
#"T""".!P."""T"#
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
#############?
+...#.......>#
#.P.|....P...#
##|###!.....!#
#...>##|######
#!P..|...P...+
?##########+##
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
	RoomHome4 = `
?############?
#.....#.....!#
#.P...|......#
##|####....P.+
#....>#>.##|##
+..P..####...#
#!....|....P.#
?###+########?
`
	RoomHome5 = `
?######+######?
#_...........!#
#..####|####..#
+P.#!..P..>#.P+
+..#>..P..!#..+
#..####|####..#
#!..........._#
?######+######?
`
	RoomCaban = `
???????-??????
?????""""?????
???""""""""???
??"""###."""??
?"""#>!|PP"""?
-""""###.""""-
??""""""""""??
????"""""?????
???????-??????`
	RoomDolmen = `
???????-?????
????.....????
????.!#!.????
???..#>#..???
??....P....??
?...#_._#...?
-.P...P...P.-
??...#>#...??
????.!#!.????
????.....????
???????-?????`
	RoomRuins = `
????????-????
????....P????
???..###.????
???.##>#..???
-.....P....??
?..##_"!##..?
??.#""""#.P.-
??.""####..??
????"#>#.????
????..P..????
??????-??????`
	RoomPillars = `
???????-?????
???.......???
?...#!#P#...?
?.#!......#.?
-.P.#>#>#.P.-
?.#......!#.?
?...#P#!#...?
???.......???
?????-???????`
	RoomRoundColumns = `
???##+##???
??#_..._#??
##!.#P#.!##
+...P>P...+
##!.#P#.!##
??#_..._#??
???##+##???`
	RoomTriangle = `
?????#?????
????#>#????
???#!.!#???
??#_..._#??
##!.#.#.!##
+..P...P..+
##_..P.._##
??###+###??`
	RoomRevTriangle = `
??###+###??
##_..P.._##
+..P...P..+
##!.#.#.!##
??#_..._#??
???#!.!#???
????#>#????
?????#?????`
	RoomSpiraling = `
?#####+#
#!.P...+
#.>#####
#!.P..!#
#####>.#
+...P.!#
#+#####?`
	RoomCircleDouble = `
???####+####???
??#""..P..""#??
?#""..#|#..""#?
#"""!#!P!#!"""#
#"._#.....#_."#
#..#..>#>..#..#
+.P|P.###.P|P.+
#..#..>#>..#..#
#"._#.....#_."#
#"""!#!P!#!"""#
?#""..#|#..""#?
??#""..P..""#??
???####+####???`
	RoomAltar = `
#+#??#######??#+#
+P_##>..!..>##_P+
#...#..#>#..#...#
?##.#!..P..!#.##?
??#..P.....P..#??
???#####+#####???`
	RoomRoundGarden = `
???##+##???
??#>.P.>#??
##!.""".!##
+.P""T""P.+
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
?#################?
#"""""".>!>.""""""#
+....P...P...P....+
#"""""".>!>.""""""#
?#################?`
)

var roomBigTemplates = []string{RoomBigColumns, RoomBigGarden, RoomColumns, RoomRoundColumns, RoomRoundGarden, RoomLongHall,
	RoomGardenHall, RoomHome1, RoomHome2, RoomHome3, RoomHome4, RoomHome5, RoomTriangle, RoomRevTriangle, RoomSpiraling, RoomAltar, RoomCircleDouble, RoomBigRooms, RoomCaban, RoomDolmen, RoomRuins, RoomPillars}

const (
	CellShaedra1 = `
#########
#HMΔ#_!_#
##|###|##
+.G.P.G.+
##|###|##
#_!_#_!_#
#########
`
	CellShaedra2 = `
#########
#_!_#_!_#
##|###|##
+.G.P.G.+
##|###|##
#_!_#HMΔ#
#########
`
	CellShaedra3 = `
#########
#_!_#_!_#
##|###|##
+.G.P.G.+
##|###|##
#HMΔ#_!_#
#########
`
	CellShaedra4 = `
#########
#_!_#HMΔ#
##|###|##
+.G.P.G.+
##|###|##
#_!_#_!_#
#########
`
)

// TODO: add indestructible walls?

var roomCellTemplates = []string{CellShaedra1, CellShaedra2, CellShaedra3, CellShaedra4}

const (
	RoomArtifact = `
????#????
???#A#???
??#.MΔ#??
?###|###?
#!_#.#_!#
+P.G.G.P+
#>..P..>#
?###+###?
`
)

var roomArtifactTemplates = []string{RoomArtifact}

const (
	RoomSpecialVampires = `
???####+####???
??#!...P...!#??
##_#.......#_##
+.P.........P.+
#..#.G...G.#..#
#.....#!#.....#
#..#...>...#..#
#!...G.#.G...!#
?##>.......>##?
???#########???`
	RoomSpecialNixes = `
?#####+#####?
#>.G.#.#.G.>#
?#.........#?
??#!_.P._!#??
???#.....#???
??#!_.P._!#??
?#.........#?
#>.G.#.#.G.>#
?#####+#####?`
	RoomSpecialMilfids = `
?????????-???
???......P.??
??..!?G._?..?
?.?.?#>#.?.??
-P.G.>#>G...?
?.?._#>#.??.?
??.......?..-
????!.G...P??
????????-????`
	RoomSpecialTreeMushrooms = `
?????"--.???????
???"""..G."""???
???""?....!"""??
??.....T>T..G..?
?..G.T..!..T..P-
-...!..T>T._"".-
-P.???..G..""""?
.-???????.-?????`
	RoomSpecialHarpies = `
?-????##??????
?P???#..#?????
?.???#G.>####?
?.??#. .##.>#?
?.?#.._....G.#
-G........!..#
??.#.G_.._#>#?
??P?#.>###?#??
??-??##???????`
	RoomSpecialCelmists = `
?#############+##?
#>#_.......>#.P._#
#...G!#!G..##....#
#....###....|...P+
#_...###....|...P+
##..G!#!G..##....#
#>.........>#.P._#
?#############+##?`
	RoomSpecialCelmists2 = `
?##+#########+##?
#_#.....G.....#_#
#..#!#W#|#W#!#..#
#..###.>.>.###..#
+P.|.G.....G.|.P+
#..###.>.>.###..#
#..#!#W#|#W#!#..#
#_#.....G.....#_#
?##+#########+##?`
	RoomSpecialCelmists3 = `
########-########
+.P|....P....|P.+
##|##W#|||#W##|##
?#....G...G....#?
??#!._#._.#_.!#??
???#..G...G..#???
????#!.....!#????
?????#>#>#>#?????
??????#?#?#??????`
	RoomSpecialMirrorSpecters = `
########-#########
-P.....W.W......P-
##W##_.#.#._###W##
-P......G.......P-
##W##..##W##.##W##
#>.!W..W>!>W.W!.>#
#.G.#..#.G.#.#.G.#
#................#
##################`
)

type specialRoom int

const (
	noSpecialRoom specialRoom = iota
	roomMilfids
	roomNixes
	roomVampires
	roomCelmists
	roomHarpies
	roomTreeMushrooms
	roomMirrorSpecters
	roomShaedra
	roomArtifact
)

func (sr specialRoom) Templates() (tpl []string) {
	switch sr {
	case roomMilfids:
		tpl = append(tpl, RoomSpecialMilfids)
	case roomVampires:
		tpl = append(tpl, RoomSpecialVampires)
	case roomCelmists:
		tpl = append(tpl, RoomSpecialCelmists, RoomSpecialCelmists2, RoomSpecialCelmists3)
	case roomNixes:
		tpl = append(tpl, RoomSpecialNixes)
	case roomHarpies:
		tpl = append(tpl, RoomSpecialHarpies)
	case roomTreeMushrooms:
		tpl = append(tpl, RoomSpecialTreeMushrooms)
	case roomMirrorSpecters:
		tpl = append(tpl, RoomSpecialMirrorSpecters)
	case roomShaedra:
		tpl = roomCellTemplates
	case roomArtifact:
		tpl = roomArtifactTemplates
	}
	return tpl
}

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
		case '.', '>', '!', 'P', '_', '|', 'M', 'Δ', 'G', '-':
			if pos.valid() {
				dg.d.SetCell(pos, GroundCell)
			}
		case '#', '+':
			if pos.valid() {
				dg.d.SetCell(pos, WallCell)
			}
		case 'T':
			if pos.valid() {
				dg.d.SetCell(pos, TreeCell)
			}
		case 'W':
			if pos.valid() {
				dg.d.SetCell(pos, WindowCell)
			}
		}
		switch c {
		case '>':
			r.places = append(r.places, place{pos: pos, kind: PlaceSpecialStatic})
		case '!':
			r.places = append(r.places, place{pos: pos, kind: PlaceItem})
		case 'P':
			r.places = append(r.places, place{pos: pos, kind: PlacePatrol})
		case 'G':
			r.places = append(r.places, place{pos: pos, kind: PlacePatrolSpecial})
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
				dg.d.SetCell(pos, FoliageCell)
			}
		case 'H':
			r.places = append(r.places, place{pos: pos, kind: PlaceStory})
			dg.spl.Shaedra = pos
			dg.d.SetCell(pos, StoryCell)
		case 'M':
			r.places = append(r.places, place{pos: pos, kind: PlaceStory})
			dg.spl.Marevor = pos
			dg.d.SetCell(pos, StoryCell)
		case 'Δ':
			r.places = append(r.places, place{pos: pos, kind: PlaceStory})
			dg.spl.Monolith = pos
			dg.d.SetCell(pos, StoryCell)
		case 'A':
			r.places = append(r.places, place{pos: pos, kind: PlaceStory})
			dg.spl.Artifact = pos
			dg.d.SetCell(pos, StoryCell)
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
	if !pos.valid() || d.Cell(pos).IsPassable() {
		return false
	}
	return pos.W().valid() && pos.E().valid() &&
		d.Cell(pos.W()).IsGround() && d.Cell(pos.E()).IsGround() &&
		(!pos.N().valid() || d.Cell(pos.N()).T == WallCell) &&
		(!pos.S().valid() || d.Cell(pos.S()).T == WallCell) &&
		((pos.NW().valid() && d.Cell(pos.NW()).IsPassable()) ||
			(pos.SW().valid() && d.Cell(pos.SW()).IsPassable()) ||
			(pos.NE().valid() && d.Cell(pos.NE()).IsPassable()) ||
			(pos.SE().valid() && d.Cell(pos.SE()).IsPassable())) ||
		pos.N().valid() && pos.S().valid() &&
			d.Cell(pos.N()).IsGround() && d.Cell(pos.S()).IsGround() &&
			(!pos.E().valid() || d.Cell(pos.E()).T == WallCell) &&
			(!pos.W().valid() || d.Cell(pos.W()).T == WallCell) &&
			((pos.NW().valid() && d.Cell(pos.NW()).IsPassable()) ||
				(pos.SW().valid() && d.Cell(pos.SW()).IsPassable()) ||
				(pos.NE().valid() && d.Cell(pos.NE()).IsPassable()) ||
				(pos.SE().valid() && d.Cell(pos.SE()).IsPassable()))
}

func (dg *dgen) PutHoledWalls(g *game, n int) {
	candidates := []position{}
	for i, _ := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if dg.room[pos] && g.HoledWallCandidate(pos) {
			candidates = append(candidates, pos)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		pos := candidates[RandInt(len(candidates))]
		g.Dungeon.SetCell(pos, HoledWallCell)
	}
}

func (dg *dgen) PutWindows(g *game, n int) {
	candidates := []position{}
	for i, _ := range g.Dungeon.Cells {
		pos := idxtopos(i)
		if dg.room[pos] && g.HoledWallCandidate(pos) {
			candidates = append(candidates, pos)
		}
	}
	if len(candidates) == 0 {
		return
	}
	for i := 0; i < n; i++ {
		pos := candidates[RandInt(len(candidates))]
		g.Dungeon.SetCell(pos, WindowCell)
	}
}

func (g *game) HoledWallCandidate(pos position) bool {
	d := g.Dungeon
	if !pos.valid() || !d.Cell(pos).IsWall() {
		return false
	}
	return pos.W().valid() && pos.E().valid() &&
		d.Cell(pos.W()).IsWall() && d.Cell(pos.E()).IsWall() &&
		pos.N().valid() && d.Cell(pos.N()).IsPassable() &&
		pos.S().valid() && d.Cell(pos.S()).IsPassable() ||
		(pos.W().valid() && pos.E().valid() &&
			d.Cell(pos.W()).IsPassable() && d.Cell(pos.E()).IsPassable() &&
			pos.N().valid() && d.Cell(pos.N()).IsWall() &&
			pos.S().valid() && d.Cell(pos.S()).IsWall())
}

type placement int

const (
	PlacementRandom placement = iota
	PlacementCenter
	PlacementEdge
)

func (dg *dgen) GenRooms(templates []string, n int, pl placement) (ok bool, ps []position) {
	ok = true
	for i := 0; i < n; i++ {
		var r *room
		count := 200
		var pos position
		var tpl string
		for r == nil && count > 0 {
			count--
			switch pl {
			case PlacementRandom:
				pos = position{RandInt(DungeonWidth - 1), RandInt(DungeonHeight - 1)}
			case PlacementCenter:
				pos = position{DungeonWidth/2 - 4 + RandInt(5), DungeonHeight/2 - 3 + RandInt(4)}
			case PlacementEdge:
				if RandInt(2) == 0 {
					pos = position{RandInt(DungeonWidth / 4), RandInt(DungeonHeight - 1)}
				} else {
					pos = position{3*DungeonWidth/4 + RandInt(DungeonWidth/4) - 1, RandInt(DungeonHeight - 1)}
				}
			}
			tpl = templates[RandInt(len(templates))]
			r = dg.NewRoom(pos, tpl)
		}
		if r != nil {
			dg.rooms = append(dg.rooms, r)
			ps = append(ps, pos)
		} else {
			ok = false
		}
	}
	return ok, ps
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

type maplayout int

const (
	AutomataCave maplayout = iota
	RandomWalkCave
	RandomWalkTreeCave
)

func (g *game) GenRoomTunnels(ml maplayout) {
	dg := dgen{}
	dg.layout = ml
	d := &dungeon{}
	d.Cells = make([]cell, DungeonNCells)
	dg.d = d
	dg.tunnel = make(map[position]bool)
	dg.room = make(map[position]bool)
	dg.rooms = []*room{}
	dg.fungus = make(map[position]vegetation)
	switch ml {
	case AutomataCave:
		dg.GenCellularAutomataCaveMap()
	case RandomWalkCave:
		dg.GenCaveMap()
	case RandomWalkTreeCave:
		dg.GenTreeCaveMap()
	}
	var places []position
	var nspecial = 4
	if sr := g.Params.Special[g.Depth]; sr != noSpecialRoom {
		nspecial--
		pl := PlacementEdge
		if RandInt(2) == 0 {
			pl = PlacementCenter
		}
		dg.special = sr
		dg.GenRooms(sr.Templates(), 1, pl)
	}
	switch ml {
	case RandomWalkCave:
		dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
		dg.GenRooms(roomNormalTemplates, 5, PlacementRandom)
	case RandomWalkTreeCave:
		if g.Depth == MaxDepth {
			g.Objects.Story = map[position]story{}
			g.Places.Artifact = dg.spl.Artifact
			g.Objects.Story[g.Places.Artifact] = StoryArtifactSealed
			g.Places.Monolith = dg.spl.Monolith
			g.Objects.Story[g.Places.Monolith] = NoStory
			g.Places.Marevor = dg.spl.Marevor
			g.Objects.Story[g.Places.Marevor] = NoStory
			dg.GenRooms(roomBigTemplates, nspecial+1, PlacementRandom)
			dg.GenRooms(roomNormalTemplates, 7, PlacementRandom)
		} else {
			dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
			dg.GenRooms(roomNormalTemplates, 7, PlacementRandom)
		}
	default:
		if g.Depth == WinDepth {
			g.Objects.Story = map[position]story{}
			g.Places.Shaedra = dg.spl.Shaedra
			g.Objects.Story[g.Places.Shaedra] = StoryShaedra
			g.Places.Monolith = dg.spl.Monolith
			g.Objects.Story[g.Places.Monolith] = NoStory
			g.Places.Marevor = dg.spl.Marevor
			g.Objects.Story[g.Places.Marevor] = NoStory
			dg.GenRooms(roomBigTemplates, nspecial, PlacementRandom)
			dg.GenRooms(roomNormalTemplates, 7, PlacementRandom)
		} else {
			dg.GenRooms(roomBigTemplates, nspecial-1, PlacementRandom)
			dg.GenRooms(roomNormalTemplates, 7, PlacementRandom)
		}
	}
	dg.ConnectRooms()
	g.Dungeon = d
	dg.PutDoors(g)
	dg.PlayerStartCell(g, places)
	dg.ClearUnconnected(g)
	if RandInt(10) > 0 {
		var t terrain
		if RandInt(5) > 1 {
			t = ChasmCell
		} else {
			t = WaterCell
		}
		if g.Depth == WinDepth || g.Depth == MaxDepth {
			t = WaterCell
		}
		dg.GenLake(t)
		if RandInt(5) == 0 {
			dg.GenLake(t)
		}
	}
	if g.Depth < MaxDepth {
		if g.Params.Blocked[g.Depth] {
			dg.GenStairs(g, BlockedStair)
		} else {
			dg.GenStairs(g, NormalStair)
		}
	}
	for i := 0; i < 4+RandInt(2); i++ {
		dg.GenBarrel(g)
	}
	dg.AddSpecial(g, ml)
	dg.ComputeConnectedComponents(func(pos position) bool {
		return pos.valid() && g.Dungeon.Cell(pos).IsPassable()
	})
	dg.GenMonsters(g)
}

func (dg *dgen) ClearUnconnected(g *game) {
	d := dg.d
	conn, _ := d.Connected(g.Player.Pos, d.IsFreeCell)
	for i, c := range d.Cells {
		pos := idxtopos(i)
		if c.IsPassable() && !conn[pos] {
			d.SetCell(pos, WallCell)
		}
	}
}

func (dg *dgen) AddSpecial(g *game, ml maplayout) {
	// Equipment
	//switch g.GenPlan[g.Depth] {
	//case GenRod:
	////g.GenerateRod()
	//case GenExtraCollectables:
	////for i := 0; i < 2; i++ {
	////dg.GenCollectable(g)
	////g.CollectableScore-- // these are extra
	////}
	//}
	g.Objects.Stones = map[position]stone{}
	if g.Params.Blocked[g.Depth] || g.Depth == MaxDepth {
		dg.GenBarrierStone(g)
	}
	bananas := 2
	if ml == RandomWalkTreeCave && g.Depth < MaxDepth {
		bananas--
	}
	for i := 0; i < bananas; i++ {
		dg.GenBanana(g)
	}
	if !g.Params.NoMagara[g.Depth] {
		dg.GenMagara(g)
	}
	dg.GenItem(g)
	dg.GenStones(g)
	dg.GenLight(g)
	ntables := 2
	switch ml {
	case AutomataCave, RandomWalkCave:
		if RandInt(3) == 0 {
			ntables++
		} else if RandInt(10) == 0 {
			ntables--
		}
	case RandomWalkTreeCave:
		if RandInt(4) > 0 {
			ntables++
		}
		if RandInt(4) > 0 {
			ntables++
		}
	}
	if g.Params.Tables[g.Depth] {
		ntables += 2 + RandInt(2)
	}
	for i := 0; i < 3+RandInt(2); i++ {
		dg.GenTable(g)
	}
	ntrees := 1
	switch ml {
	case AutomataCave:
		if RandInt(4) == 0 {
			ntrees++
		} else if RandInt(8) == 0 {
			ntrees--
		}
	case RandomWalkCave:
		if RandInt(4) > 0 {
			ntrees++
		}
		if RandInt(8) == 0 {
			ntrees++
		}
	case RandomWalkTreeCave:
		if RandInt(2) == 0 {
			ntrees--
		}
	}
	if g.Params.Trees[g.Depth] {
		ntrees += 2 + RandInt(2)
	}
	for i := 0; i < ntrees; i++ {
		dg.GenTree(g)
	}
	nhw := 1 + RandInt(2)
	if g.Params.Holes[g.Depth] {
		nhw += 3 + RandInt(2)
	}
	dg.PutHoledWalls(g, nhw)
	nwin := 1
	if nhw == 1 {
		nwin++
	}
	if g.Params.Windows[g.Depth] {
		nwin += 4 + RandInt(3)
	}
	dg.PutWindows(g, nwin)
	if g.Params.Lore[g.Depth] {
		dg.PutLore(g)
	}
}

func (dg *dgen) PutLore(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 2000 {
			panic("PutLore1")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	count = 0
	for {
		count++
		if count > 1000 {
			panic("PutLore2")
		}
		i := RandInt(len(LoreMessages))
		if g.GeneratedLore[i] {
			continue
		}
		g.GeneratedLore[i] = true
		g.Objects.Lore[pos] = i
		g.Objects.Scrolls[pos] = ScrollLore
		g.Dungeon.SetCell(pos, ScrollCell)
		break
	}
}

func (dg *dgen) GenLight(g *game) {
	lights := []position{}
	for i := 0; i < 2+RandInt(2); i++ {
		pos := dg.OutsideGroundCell(g)
		g.Dungeon.SetCell(pos, LightCell)
		lights = append(lights, pos)
	}
	for i := 0; i < 8+RandInt(4); i++ {
		pos := dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
		if pos != InvalidPos {
			g.Dungeon.SetCell(pos, LightCell)
			lights = append(lights, pos)
		}
	}
	for _, pos := range lights {
		g.Objects.Lights[pos] = true
	}
	g.ComputeLights()
}

func (r *room) RandomPlace(kind placeKind) position {
	var p []int
	for i, pl := range r.places {
		if pl.kind == kind && !pl.used {
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

var PlaceSpecialOrStatic = []placeKind{PlaceSpecialStatic, PlaceStatic}

func (r *room) RandomPlaces(kinds []placeKind) position {
	pos := InvalidPos
	for _, kind := range kinds {
		pos = r.RandomPlace(kind)
		if pos != InvalidPos {
			break
		}
	}
	return pos
}

func (dg *dgen) PlayerStartCell(g *game, places []position) {
	const far = 30
	r := dg.rooms[len(dg.rooms)-1]
loop:
	for i := len(dg.rooms) - 2; i >= 0; i-- {
		for _, pos := range places {
			if r.pos.Distance(pos) < far && dg.rooms[i].pos.Distance(pos) >= far {
				r = dg.rooms[i]
				break loop
			}
		}
	}
	g.Player.Pos = r.RandomPlace(PlacePatrol)
	if g.Depth > 1 {
		return
	}
	itpos := r.RandomPlace(PlaceItem)
	// TODO: make the player only start in rooms with enough places
	if itpos == InvalidPos {
		itpos = r.RandomPlaces(PlaceSpecialOrStatic)
		if itpos == InvalidPos {
			panic("no item")
		}
	}
	g.Dungeon.SetCell(itpos, ScrollCell)
	g.Objects.Scrolls[itpos] = ScrollStory
}

func (dg *dgen) GenBanana(g *game) {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("GenBanana")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := dg.d.Cell(pos)
		if c.T == GroundCell && !dg.room[pos] {
			dg.d.SetCell(pos, BananaCell)
			g.Objects.Bananas[pos] = true
			break
		}
	}
}

func (dg *dgen) OutsideGroundCell(g *game) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("OutsideGroundCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if c.T == GroundCell && !dg.room[pos] {
			return pos
		}
	}
}

func (dg *dgen) OutsideGroundMiddleCell(g *game) position {
	count := 0
	for {
		count++
		if count > 2000 {
			panic("OutsideGroundMiddleCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if c.T == GroundCell && !dg.room[pos] && !dg.d.HasTooManyWallNeighbors(pos) {
			return pos
		}
	}
}

func (dg *dgen) FoliageCell(g *game) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FoliageCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if c.T == FoliageCell {
			return pos
		}
	}
}

func (dg *dgen) OutsideCell(g *game) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("OutsideCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		c := dg.d.Cell(pos)
		if !dg.room[pos] && (c.T == FoliageCell || c.T == GroundCell) {
			return pos
		}
	}
}

func (dg *dgen) InsideCell(g *game) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("InsideCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		if pos.Distance(g.Player.Pos) < DefaultLOSRange {
			continue
		}
		c := dg.d.Cell(pos)
		if dg.room[pos] && (c.T == FoliageCell || c.T == GroundCell) {
			return pos
		}
	}
}

func (dg *dgen) GenItem(g *game) {
	plan := g.GenPlan[g.Depth]
	if plan != GenAmulet && plan != GenCloak {
		return
	}
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenItem")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(pos, ItemCell)
	var it item
	switch plan {
	case GenCloak:
		it = g.RandomCloak()
		g.GeneratedCloaks = append(g.GeneratedCloaks, it)
	case GenAmulet:
		it = g.RandomAmulet()
		g.GeneratedAmulets = append(g.GeneratedAmulets, it)
	}
	g.Objects.Items[pos] = it
}

func (dg *dgen) GenBarrierStone(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenBarrierStone")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(pos, StoneCell)
	g.Objects.Stones[pos] = SealStone
}

func (dg *dgen) GenMagara(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 1000 {
			panic("GenMagara")
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceItem)
	}
	g.Dungeon.SetCell(pos, MagaraCell)
	mag := g.RandomMagara()
	g.Objects.Magaras[pos] = mag
	g.GeneratedMagaras = append(g.GeneratedMagaras, mag)
}

func (dg *dgen) GenStairs(g *game, st stair) {
	var ri, pj int
	best := 0
	for i, r := range dg.rooms {
		for j, pl := range r.places {
			score := pl.pos.Distance(g.Player.Pos) + RandInt(20)
			if !pl.used && pl.kind == PlaceSpecialStatic && score > best {
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

func (dg *dgen) GenBarrel(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 500 {
			return
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceSpecialStatic)
	}
	g.Dungeon.SetCell(pos, BarrelCell)
	g.Objects.Barrels[pos] = true
}

func (dg *dgen) GenTable(g *game) {
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 500 {
			return
		}
		pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlaces(PlaceSpecialOrStatic)
	}
	g.Dungeon.SetCell(pos, TableCell)
}

func (dg *dgen) GenTree(g *game) {
	pos := dg.OutsideGroundMiddleCell(g)
	g.Dungeon.SetCell(pos, TreeCell)
}

func (dg *dgen) CaveGroundCell(g *game) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("CaveGroundCell")
		}
		x := RandInt(DungeonWidth)
		y := RandInt(DungeonHeight)
		pos := position{x, y}
		c := dg.d.Cell(pos)
		if c.T == GroundCell && !dg.room[pos] {
			return pos
		}
	}
}

func (dg *dgen) GenStones(g *game) {
	// Magical Stones
	// TODO: move into dungeon generation
	nstones := 3
	switch RandInt(8) {
	case 1, 2, 3, 4, 5:
		nstones++
	case 6, 7:
		nstones += 2
	}
	inroom := 2
	if g.Params.Stones[g.Depth] {
		nstones += 4 + RandInt(3)
		inroom += 2
	}
	for i := 0; i < nstones; i++ {
		pos := InvalidPos
		if i < inroom {
			count := 0
			for pos == InvalidPos {
				count++
				if count > 1000 {
					panic("GenStones")
				}
				pos = dg.rooms[RandInt(len(dg.rooms))].RandomPlace(PlaceStatic)
			}
		} else {
			pos = dg.CaveGroundCell(g)
		}
		st := stone(1 + RandInt(NumStones-1))
		g.Objects.Stones[pos] = st
		g.Dungeon.SetCell(pos, StoneCell)
	}
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
	dg.Foliage(false)
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

func (dg *dgen) GenLake(t terrain) {
	walls := []position{}
	xshift := 10 + RandInt(5)
	yshift := 5 + RandInt(3)
	for i := 0; i < DungeonNCells; i++ {
		pos := idxtopos(i)
		if pos.X < xshift || pos.Y < yshift || pos.X > DungeonWidth-xshift || pos.Y > DungeonHeight-yshift {
			continue
		}
		c := dg.d.Cell(pos)
		if c.T == WallCell && !dg.room[pos] {
			walls = append(walls, pos)
		}
	}
	count := 0
	var bestpos = walls[RandInt(len(walls))]
	var bestsize int
	d := dg.d
	for {
		pos := walls[RandInt(len(walls))]
		_, size := d.Connected(pos, func(npos position) bool {
			return npos.valid() && dg.d.Cell(npos).T == WallCell && !dg.room[npos] && pos.Distance(npos) < 10+RandInt(10)
		})
		count++
		if Abs(bestsize-90) > Abs(size-90) {
			bestsize = size
			bestpos = pos
		}
		if count > 15 || Abs(size-90) < 25 {
			break
		}
	}
	conn, _ := d.Connected(bestpos, func(npos position) bool {
		return npos.valid() && dg.d.Cell(npos).T == WallCell && !dg.room[npos] && bestpos.Distance(npos) < 10+RandInt(10)
	})
	for pos := range conn {
		d.SetCell(pos, t)
	}
}

//func (dg *dgen) GenChasm() {
//pos := position{20 + RandInt(DungeonWidth-21), 7 + RandInt(DungeonHeight-8)}
//size := 100 + RandInt(100)
//var queue [DungeonNCells]int
//var visited [DungeonNCells]bool
//var qstart, qend int
//visited[pos.idx()] = true
//queue[qend] = pos.idx()
//qend++
//nb := make([]position, 4)
//for qstart < qend && size > 0 {
//cidx := queue[qstart]
//qstart++
//cpos := idxtopos(cidx)
//dg.d.SetCell(cpos, ChasmCell)
//size--
//for _, npos := range cpos.CardinalNeighbors(nb, func(p position) bool { return p.valid() }) {
//nidx := npos.idx()
//if !visited[nidx] {
//if RandInt(3) > 0 || size > 50 || qend-qstart < 4 {
//queue[qend] = nidx
//qend++
//}
//visited[nidx] = true
//}
//}
//}
//}

func (dg *dgen) Foliage(less bool) {
	// use same structure as for the dungeon
	// walls will become foliage
	d := &dungeon{}
	d.Cells = make([]cell, DungeonNCells)
	limit := 43
	if less {
		limit = 40
	}
	for i := range d.Cells {
		r := RandInt(100)
		pos := idxtopos(i)
		if r >= limit {
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
			dg.d.SetCell(idxtopos(i), FoliageCell)
		}
	}
}

func (dg *dgen) GenCaveMap() {
	d := dg.d
	for i := range d.Cells {
		pos := idxtopos(i)
		d.SetCell(pos, WallCell)
	}
	pos := position{40, 10}
	max := 21 * 40
	d.SetCell(pos, GroundCell)
	cells := 1
	notValid := 0
	lastValid := pos
	for cells < max {
		npos := pos.RandomNeighbor(false)
		if !pos.valid() && npos.valid() && d.Cell(npos).T == WallCell {
			pos = lastValid
			continue
		}
		pos = npos
		if pos.valid() {
			if d.Cell(pos).T != GroundCell {
				d.SetCell(pos, GroundCell)
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
	dg.Foliage(false)
}

func (d *dungeon) DigBlock(block []position) []position {
	pos := d.WallCell()
	block = block[:0]
	count := 0
	for {
		count++
		if count > 10000 {
			panic("DigBlock")
		}
		block = append(block, pos)
		if d.HasFreeNeighbor(pos) {
			break
		}
		pos = pos.RandomNeighbor(false)
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

func (dg *dgen) GenTreeCaveMap() {
	d := dg.d
	center := position{40, 10}
	d.SetCell(center, GroundCell)
	d.SetCell(center.E(), GroundCell)
	d.SetCell(center.NE(), GroundCell)
	d.SetCell(center.S(), GroundCell)
	d.SetCell(center.SE(), GroundCell)
	d.SetCell(center.N(), GroundCell)
	d.SetCell(center.NW(), GroundCell)
	d.SetCell(center.W(), GroundCell)
	d.SetCell(center.SW(), GroundCell)
	max := 21 * 19
	cells := 1
	block := make([]position, 0, 64)
loop:
	for cells < max {
		block = d.DigBlock(block)
		if len(block) == 0 {
			continue loop
		}
		for _, pos := range block {
			if d.Cell(pos).T != GroundCell {
				d.SetCell(pos, GroundCell)
				cells++
			}
		}
	}
	dg.Foliage(true)
}

// monster generation

func (g *game) GenBand(band monsterBand) []monsterKind {
	mbd := MonsBands[band]
	if g.GeneratedUniques[band] > 0 && mbd.Unique {
		return nil
	}
	if !mbd.Band {
		return []monsterKind{mbd.Monster}
	}
	bandMonsters := []monsterKind{}
	for m, n := range mbd.Distribution {
		for i := 0; i < n; i++ {
			bandMonsters = append(bandMonsters, m)
		}
	}
	return bandMonsters
}

func (dg *dgen) BandInfoGuard(g *game, band monsterBand, pl placeKind) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
loop:
	for pos == InvalidPos {
		count++
		if count > 3000 {
			pos = dg.InsideCell(g)
			break
		}
		for i := 0; i < 10; i++ {
			r := dg.rooms[RandInt(len(dg.rooms)-1)]
			for _, e := range r.places {
				if e.kind == PlaceSpecialStatic {
					pos = r.RandomPlace(pl)
					if pos != InvalidPos {
						break loop
					}
				}
			}
		}
		r := dg.rooms[RandInt(len(dg.rooms)-1)]
		pos = r.RandomPlace(pl)
	}
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Beh = BehGuard
	return bandinfo
}

func (dg *dgen) BandInfoPatrol(g *game, band monsterBand, pl placeKind) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	pos := InvalidPos
	count := 0
	for pos == InvalidPos {
		count++
		if count > 5000 {
			pos = dg.InsideCell(g)
			break
		}
		pos = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(pl)
	}
	target := InvalidPos
	count = 0
	for target == InvalidPos {
		// TODO: only find place in other room?
		count++
		if count > 5000 {
			pos = dg.InsideCell(g)
			break
		}
		target = dg.rooms[RandInt(len(dg.rooms)-1)].RandomPlace(pl)
	}
	bandinfo.Path = append(bandinfo.Path, pos)
	bandinfo.Path = append(bandinfo.Path, target)
	bandinfo.Beh = BehPatrol
	return bandinfo
}

func (dg *dgen) BandInfoOutsideGround(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideGroundCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoOutsideGroundMiddle(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideGroundMiddleCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoOutside(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) BandInfoOutsideExplore(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	for i := 0; i < 4; i++ {
		for j := 0; j < 100; j++ {
			pos := dg.OutsideCell(g)
			if dg.cc[pos.idx()] == dg.cc[bandinfo.Path[0].idx()] {
				bandinfo.Path = append(bandinfo.Path, pos)
				break
			}
		}
	}
	bandinfo.Beh = BehExplore
	return bandinfo
}

func (dg *dgen) BandInfoOutsideExploreButterfly(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.OutsideCell(g))
	for i := 0; i < 9; i++ {
		for j := 0; j < 100; j++ {
			pos := dg.OutsideCell(g)
			if dg.cc[pos.idx()] == dg.cc[bandinfo.Path[0].idx()] {
				bandinfo.Path = append(bandinfo.Path, pos)
				break
			}
		}
	}
	bandinfo.Beh = BehExplore
	return bandinfo
}

func (dg *dgen) BandInfoFoliage(g *game, band monsterBand) bandInfo {
	bandinfo := bandInfo{Kind: monsterBand(band)}
	bandinfo.Path = append(bandinfo.Path, dg.FoliageCell(g))
	bandinfo.Beh = BehWander
	return bandinfo
}

func (dg *dgen) PutMonsterBand(g *game, band monsterBand) bool {
	monsters := g.GenBand(band)
	if monsters == nil {
		return false
	}
	var bdinf bandInfo
	switch band {
	case LoneYack, LoneWorm, PairYack:
		bdinf = dg.BandInfoFoliage(g, band)
	case LoneDog, LoneHarpy:
		bdinf = dg.BandInfoOutsideGround(g, band)
	case LoneBlinkingFrog, LoneExplosiveNadre, PairExplosiveNadre:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneMirrorSpecter, LoneWingedMilfid, LoneVampire, PairWingedMilfid, LoneEarthDragon:
		bdinf = dg.BandInfoOutsideExplore(g, band)
	case LoneButterfly:
		bdinf = dg.BandInfoOutsideExploreButterfly(g, band)
	case LoneTreeMushroom:
		bdinf = dg.BandInfoOutside(g, band)
	case LoneHighGuard:
		bdinf = dg.BandInfoGuard(g, band, PlacePatrol)
	case LoneSatowalgaPlant:
		bdinf = dg.BandInfoOutsideGroundMiddle(g, band)
	case SpecialLoneVampire, SpecialLoneNixe, SpecialLoneMilfid, SpecialLoneOricCelmist, SpecialLoneHarmonicCelmist, SpecialLoneHighGuard,
		SpecialLoneHarpy, SpecialLoneTreeMushroom:
		if RandInt(5) > 0 {
			bdinf = dg.BandInfoPatrol(g, band, PlacePatrolSpecial)
		} else {
			bdinf = dg.BandInfoGuard(g, band, PlacePatrolSpecial)
		}
	default:
		bdinf = dg.BandInfoPatrol(g, band, PlacePatrol)
	}
	g.Bands = append(g.Bands, bdinf)
	awake := RandInt(5) > 0
	var pos position
	if len(bdinf.Path) == 0 {
		// should not happen now
		pos = g.FreeCellForMonster()
	} else {
		pos = bdinf.Path[0]
	}
	for _, mk := range monsters {
		mons := &monster{Kind: mk}
		if awake {
			mons.State = Wandering
		}
		mons.Init()
		mons.Index = len(g.Monsters)
		mons.Band = len(g.Bands) - 1
		mons.PlaceAt(g, pos)
		mons.Target = mons.NextTarget(g)
		g.Monsters = append(g.Monsters, mons)
		pos = g.FreeCellForBandMonster(pos)
	}
	return true
}

func (dg *dgen) PutRandomBand(g *game, bands []monsterBand) bool {
	return dg.PutMonsterBand(g, bands[RandInt(len(bands))])
}

func (dg *dgen) PutRandomBandN(g *game, bands []monsterBand, n int) {
	for i := 0; i < n; i++ {
		dg.PutMonsterBand(g, bands[RandInt(len(bands))])
	}
}

type monspecial int

const (
	MonsNormal monspecial = iota
	MonsSpecialFrogs
	MonsSpecialPlants
	MonsSpecialButterflies
	MonsSpecialWorms
	MonsSpecialNadres
	MonsSpecialNixes
	MonsSpecialVampires
	MonsSpecialOricCelmists
	MonsSpecialWingedMilfid
)

func (dg *dgen) GenMonsters(g *game) {
	g.Monsters = []*monster{}
	g.Bands = []bandInfo{}
	// common bands
	bandsGuard := []monsterBand{LoneGuard}
	bandsButterfly := []monsterBand{LoneButterfly}
	bandsHighGuard := []monsterBand{LoneHighGuard}
	bandsAnimals := []monsterBand{LoneYack, LoneWorm, LoneDog, LoneBlinkingFrog, LoneExplosiveNadre, LoneHarpy}
	bandsPlants := []monsterBand{LoneSatowalgaPlant}
	bandsBipeds := []monsterBand{LoneOricCelmist, LoneMirrorSpecter, LoneWingedMilfid, LoneMadNixe, LoneVampire, LoneHarmonicCelmist}
	bandsBig := []monsterBand{LoneTreeMushroom, LoneEarthDragon}
	// monster specific bands
	bandNadre := []monsterBand{LoneExplosiveNadre}
	bandFrog := []monsterBand{LoneBlinkingFrog}
	bandDog := []monsterBand{LoneDog}
	bandYack := []monsterBand{LoneYack}
	bandVampire := []monsterBand{LoneVampire}
	bandOricCelmist := []monsterBand{LoneOricCelmist}
	bandHarmonicCelmist := []monsterBand{LoneHarmonicCelmist}
	bandMadNixe := []monsterBand{LoneMadNixe}
	bandMirrorSpecter := []monsterBand{LoneMirrorSpecter}
	bandTreeMushroom := []monsterBand{LoneTreeMushroom}
	bandDragon := []monsterBand{LoneEarthDragon}
	bandGuardPair := []monsterBand{PairGuard}
	bandYackPair := []monsterBand{PairYack}
	bandExplosiveNadrePair := []monsterBand{PairExplosiveNadre}
	bandWingedMilfidPair := []monsterBand{PairWingedMilfid}
	bandNixePair := []monsterBand{PairNixe}
	bandVampirePair := []monsterBand{PairVampire}
	bandOricCelmistPair := []monsterBand{PairOricCelmist}
	bandHarmonicCelmistPair := []monsterBand{PairHarmonicCelmist}
	// special bands
	if g.Params.Special[g.Depth] != noSpecialRoom {
		switch dg.special {
		case roomVampires:
			bandVamps := []monsterBand{SpecialLoneVampire}
			dg.PutRandomBandN(g, bandVamps, 2)
		case roomNixes:
			bandNixes := []monsterBand{SpecialLoneNixe}
			dg.PutRandomBandN(g, bandNixes, 2)
		case roomMilfids:
			bandMilfids := []monsterBand{SpecialLoneMilfid}
			dg.PutRandomBandN(g, bandMilfids, 2)
		case roomCelmists:
			switch RandInt(3) {
			case 0:
				bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
				dg.PutRandomBandN(g, bandOricCelmists, 2)
			case 1:
				bandHarmonicCelmists := []monsterBand{SpecialLoneHarmonicCelmist}
				dg.PutRandomBandN(g, bandHarmonicCelmists, 2)
			case 2:
				bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
				bandHarmonicCelmists := []monsterBand{SpecialLoneHarmonicCelmist}
				dg.PutRandomBandN(g, bandHarmonicCelmists, 1)
				dg.PutRandomBandN(g, bandOricCelmists, 1)
			}
		case roomHarpies:
			bandHarpies := []monsterBand{SpecialLoneHarpy}
			dg.PutRandomBandN(g, bandHarpies, 2)
		case roomTreeMushrooms:
			bandTreeMushrooms := []monsterBand{SpecialLoneTreeMushroom}
			dg.PutRandomBandN(g, bandTreeMushrooms, 2)
		case roomShaedra:
			if RandInt(3) > 0 {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHighGuard})
			} else {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneOricCelmist})
			}
		case roomArtifact:
			if RandInt(3) > 0 {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneOricCelmist})
			} else {
				dg.PutRandomBand(g, []monsterBand{SpecialLoneHighGuard})
			}
		default:
			// XXX not used now
			bandOricCelmists := []monsterBand{SpecialLoneOricCelmist}
			dg.PutRandomBandN(g, bandOricCelmists, 2)
		}
	}
	dg.PutRandomBandN(g, bandsButterfly, 2)
	switch g.Depth {
	case 1:
		// 8
		if RandInt(2) == 0 {
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsAnimals, 3)
		} else {
			dg.PutRandomBandN(g, bandsGuard, 4)
			dg.PutRandomBandN(g, bandsAnimals, 4)
		}
	case 2:
		// 9-11(9)
		dg.PutRandomBandN(g, bandsGuard, 3)
		switch RandInt(5) {
		case 0, 1:
			// 6
			dg.PutRandomBandN(g, bandsBipeds, 1)
			dg.PutRandomBandN(g, bandsAnimals, 3)
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 2, 3:
			// 8
			dg.PutRandomBandN(g, bandsAnimals, 4)
			dg.PutRandomBandN(g, bandsButterfly, 2)
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 4:
			// 7
			dg.PutRandomBandN(g, bandsPlants, 2)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandFrog, 5)
			} else {
				dg.PutRandomBandN(g, bandYack, 5)
			}
		}
	case 3:
		// 10-11
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		dg.PutRandomBandN(g, bandsGuard, 4)
		switch RandInt(5) {
		case 0, 1:
			// 5
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandDog, 3)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 3)
			}
			dg.PutRandomBandN(g, bandsPlants, 2)
		case 2, 3:
			// 4
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
			dg.PutRandomBandN(g, bandsBipeds, 2)
		case 4:
			dg.PutRandomBandN(g, bandsPlants, 1)
			dg.PutRandomBandN(g, bandNadre, 4)
		}
	case 4:
		// 10-12
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		switch RandInt(5) {
		case 0, 1:
			// 8
			dg.PutRandomBandN(g, bandsGuard, 4)
			dg.PutRandomBandN(g, bandsBig, 2)
			dg.PutRandomBandN(g, bandsBipeds, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		case 2, 3:
			// 10
			dg.PutRandomBandN(g, bandsGuard, 8)
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		case 4:
			dg.PutRandomBandN(g, bandsGuard, 4)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandOricCelmist, 4)
			} else {
				dg.PutRandomBandN(g, bandHarmonicCelmist, 4)
			}
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 5:
		// 12
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if RandInt(2) == 0 {
			// 10
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 2)
				dg.PutRandomBandN(g, bandGuardPair, 1)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 4)
			}
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsBig, 2)
			dg.PutRandomBandN(g, bandsBipeds, 2)
			dg.PutRandomBandN(g, bandsPlants, 1)
		} else {
			// 10
			dg.PutRandomBandN(g, bandsGuard, 2)
			dg.PutRandomBandN(g, bandsAnimals, 3)
			dg.PutRandomBandN(g, bandsBipeds, 2)
			dg.PutRandomBandN(g, bandsBig, 1)
			dg.PutRandomBandN(g, bandsPlants, 2)
		}
	case 6:
		// 13-17
		dg.PutRandomBandN(g, bandsHighGuard, 1)
		if RandInt(2) == 0 {
			// 12
			dg.PutRandomBandN(g, bandsGuard, 4)
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsBig, 2)
			dg.PutRandomBandN(g, bandsBipeds, 1)
			dg.PutRandomBandN(g, bandsPlants, 3)
		} else {
			// 16
			dg.PutRandomBandN(g, bandsGuard, 2)
			if RandInt(2) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandYack, 8)
				} else {
					dg.PutRandomBandN(g, bandFrog, 8)
				}
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 8)
			}
			dg.PutRandomBandN(g, bandsButterfly, 1)
			dg.PutRandomBandN(g, bandsPlants, 5)
		}
	case 7:
		// 18-19
		dg.PutRandomBandN(g, bandsHighGuard, 1)
		if RandInt(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 4)
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandDog, 4)
				dg.PutRandomBandN(g, bandsAnimals, 2)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 6)
			}
			dg.PutRandomBandN(g, bandsButterfly, 1)
			dg.PutRandomBandN(g, bandsBig, 1)
			dg.PutRandomBandN(g, bandsBipeds, 4)
			dg.PutRandomBandN(g, bandsPlants, 2)
		} else {
			// 17
			dg.PutRandomBandN(g, bandsGuard, 1)
			dg.PutRandomBandN(g, bandsBig, 2)
			dg.PutRandomBandN(g, bandsButterfly, 1)
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandNadre, 8)
			} else {
				dg.PutRandomBandN(g, bandsAnimals, 8)
			}
			dg.PutRandomBandN(g, bandsPlants, 5)
		}
	case 8:
		// 18
		dg.PutRandomBandN(g, bandsHighGuard, 4)
		if RandInt(2) == 0 {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsBig, 1)
			if RandInt(3) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmist, 6)
				} else {
					dg.PutRandomBandN(g, bandMadNixe, 6)
				}
				dg.PutRandomBandN(g, bandsBipeds, 2)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 8)
			}
		} else {
			// 14
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsBig, 2)
			dg.PutRandomBandN(g, bandsBipeds, 4)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 9:
		// 15-23
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if RandInt(2) == 0 {
			// 13
			dg.PutRandomBandN(g, bandsGuard, 3)
			if RandInt(2) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandTreeMushroom, 6)
				} else {
					dg.PutRandomBandN(g, bandDragon, 6)
					dg.PutRandomBandN(g, bandsAnimals, 2)
				}
			} else {
				dg.PutRandomBandN(g, bandsBig, 6)
			}
			dg.PutRandomBandN(g, bandsBipeds, 3)
			dg.PutRandomBandN(g, bandsPlants, 1)
		} else {
			// 21
			dg.PutRandomBandN(g, bandsButterfly, 2)
			dg.PutRandomBandN(g, bandsGuard, 2)
			dg.PutRandomBandN(g, bandsAnimals, 8)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandExplosiveNadrePair, 2)
				dg.PutRandomBandN(g, bandsBig, 2)
			} else {
				dg.PutRandomBandN(g, bandYackPair, 3)
				dg.PutRandomBandN(g, bandsBig, 1)
			}
			dg.PutRandomBandN(g, bandsPlants, 5)
		}
	case 10:
		// 19-20
		dg.PutRandomBandN(g, bandsHighGuard, 2)
		if RandInt(2) == 0 {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 7)
			dg.PutRandomBandN(g, bandGuardPair, 1)
			dg.PutRandomBandN(g, bandsBig, 1)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandsBipeds, 8)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 4)
				dg.PutRandomBandN(g, bandMirrorSpecter, 4)
			}
		} else {
			// 17
			dg.PutRandomBandN(g, bandGuardPair, 1)
			if RandInt(3) == 0 {
				dg.PutRandomBandN(g, bandsGuard, 4)
				dg.PutRandomBandN(g, bandVampire, 5)
				dg.PutRandomBandN(g, bandsBig, 2)
			} else {
				dg.PutRandomBandN(g, bandsGuard, 6)
				dg.PutRandomBandN(g, bandsBipeds, 4)
				dg.PutRandomBandN(g, bandsAnimals, 2)
			}
			dg.PutRandomBandN(g, bandsAnimals, 2)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	case 11:
		// 23-26
		dg.PutRandomBandN(g, bandsHighGuard, 5)
		if RandInt(2) == 0 {
			// 21
			dg.PutRandomBandN(g, bandsGuard, 5)
			dg.PutRandomBandN(g, bandsBig, 2)
			if RandInt(3) == 0 {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmist, 5)
				} else {
					dg.PutRandomBandN(g, bandHarmonicCelmist, 5)
				}
				dg.PutRandomBandN(g, bandsBipeds, 5)
			} else {
				dg.PutRandomBandN(g, bandsBipeds, 10)
			}
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandVampirePair, 1)
			} else {
				if RandInt(2) == 0 {
					dg.PutRandomBandN(g, bandOricCelmistPair, 1)
				} else {
					dg.PutRandomBandN(g, bandHarmonicCelmistPair, 1)
				}
			}
			dg.PutRandomBandN(g, bandsAnimals, 2)
		} else {
			// 18
			dg.PutRandomBandN(g, bandsGuard, 7)
			dg.PutRandomBandN(g, bandsBig, 1)
			dg.PutRandomBandN(g, bandsBipeds, 6)
			if RandInt(2) == 0 {
				dg.PutRandomBandN(g, bandWingedMilfidPair, 1)
			} else {
				dg.PutRandomBandN(g, bandNixePair, 1)
			}
			dg.PutRandomBandN(g, bandsAnimals, 1)
			dg.PutRandomBandN(g, bandsPlants, 1)
		}
	}
}
