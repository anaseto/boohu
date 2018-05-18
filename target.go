package main

import "errors"

type Targeter interface {
	ComputeHighlight(*game, position)
	Action(*game, position) error
	Reachable(*game, position) bool
	Done() bool
}

type examiner struct {
	done   bool
	stairs bool
}

func (ex *examiner) ComputeHighlight(g *game, pos position) {
	g.ComputePathHighlight(pos)
}

func (g *game) ComputePathHighlight(pos position) {
	path := g.PlayerPath(g.Player.Pos, pos)
	g.Highlight = map[position]bool{}
	for _, p := range path {
		g.Highlight[p] = true
	}
}

func (ex *examiner) Action(g *game, pos position) error {
	if g.ExclusionsMap[g.Player.Pos] {
		return errors.New("You cannot travel while in an excluded area.")
	}
	if !g.Dungeon.Cell(pos).Explored {
		return errors.New("You do not know this place.")
	}
	if g.ExclusionsMap[pos] {
		return errors.New("You cannot travel to an excluded area.")
	}
	if g.Dungeon.Cell(pos).T == WallCell {
		return errors.New("You cannot travel into a wall.")
	}
	path := g.PlayerPath(g.Player.Pos, pos)
	if path == nil {
		if ex.stairs {
			return errors.New("There is no safe path to the nearest stairs.")
		}
		return errors.New("There is no safe path to this place.")
	}
	if c := g.Dungeon.Cell(pos); c.Explored && c.T == FreeCell {
		g.AutoTarget = &pos
		g.Targeting = &pos
		ex.done = true
		return nil
	}
	return errors.New("Invalid destination.")
}

func (ex *examiner) Reachable(g *game, pos position) bool {
	return true
}

func (ex *examiner) Done() bool {
	return ex.done
}

type chooser struct {
	done         bool
	area         bool
	minDist      bool
	needsFreeWay bool
	free         bool
	flammable    bool
	wall         bool
}

func (ch *chooser) ComputeHighlight(g *game, pos position) {
	g.ComputeRayHighlight(pos)
	if !ch.area {
		return
	}
	neighbors := g.Dungeon.FreeNeighbors(pos)
	for _, pos := range neighbors {
		g.Highlight[pos] = true
	}
}

func (ch *chooser) Reachable(g *game, pos position) bool {
	return g.Player.LOS[pos]
}

func (ch *chooser) Action(g *game, pos position) error {
	if !ch.Reachable(g, pos) {
		return errors.New("You cannot target that place.")
	}
	if ch.minDist && pos.Distance(g.Player.Pos) <= 1 {
		return errors.New("Invalid target: too close.")
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell {
		return errors.New("You cannot target a wall.")
	}
	if (ch.area || ch.needsFreeWay) && !ch.freeWay(g, pos) {
		return errors.New("Invalid target: there are monsters in the way.")
	}
	mons := g.MonsterAt(pos)
	if ch.free {
		if mons.Exists() {
			return errors.New("Invalid target: there is a monster there.")
		}
		if g.Player.Pos == pos {
			return errors.New("Invalid target: you are here.")
		}
		g.Player.Target = pos
		ch.done = true
		return nil
	}
	if mons.Exists() || ch.flammable && ch.flammableInWay(g, pos) {
		g.Player.Target = pos
		ch.done = true
		return nil
	}
	if ch.flammable && ch.flammableInWay(g, pos) {
		g.Player.Target = pos
		ch.done = true
		return nil
	}
	if !ch.area {
		return errors.New("You must target a monster.")
	}
	neighbors := pos.ValidNeighbors()
	for _, npos := range neighbors {
		nc := g.Dungeon.Cell(npos)
		if !ch.wall && nc.T == WallCell {
			continue
		}
		mons := g.MonsterAt(npos)
		_, okFungus := g.Fungus[pos]
		_, okDoors := g.Doors[pos]
		if ch.flammable && (okFungus || okDoors) || mons.Exists() || nc.T == WallCell {
			g.Player.Target = pos
			ch.done = true
			return nil
		}
	}
	if ch.flammable && ch.wall {
		return errors.New("Invalid target: no monsters, walls nor flammable terrain in the area.")
	}
	if ch.flammable {
		return errors.New("Invalid target: no monsters nor flammable terrain in the area.")
	}
	if ch.wall {
		return errors.New("Invalid target: no monsters nor walls in the area.")
	}
	return errors.New("Invalid target: no monsters in the area.")
}

func (ch *chooser) Done() bool {
	return ch.done
}

func (ch *chooser) freeWay(g *game, pos position) bool {
	ray := g.Ray(pos)
	tpos := pos
	for _, rpos := range ray {
		mons := g.MonsterAt(rpos)
		if !mons.Exists() {
			continue
		}
		tpos = mons.Pos
	}
	return tpos == pos
}

func (ch *chooser) flammableInWay(g *game, pos position) bool {
	ray := g.Ray(pos)
	for _, rpos := range ray {
		if rpos == g.Player.Pos {
			continue
		}
		if _, ok := g.Fungus[rpos]; ok {
			return true
		}
		if _, ok := g.Doors[rpos]; ok {
			return true
		}
	}
	return false
}

type wallChooser struct {
	done    bool
	minDist bool
}

func (ch *wallChooser) ComputeHighlight(g *game, pos position) {
	g.ComputeRayHighlight(pos)
}

func (ch *wallChooser) Reachable(g *game, pos position) bool {
	return g.Player.LOS[pos]
}

func (ch *wallChooser) Action(g *game, pos position) error {
	if !ch.Reachable(g, pos) {
		return errors.New("You cannot target that place.")
	}
	ray := g.Ray(pos)
	if len(ray) == 0 {
		return errors.New("You are not a wall.")
	}
	if g.Dungeon.Cell(ray[0]).T != WallCell {
		return errors.New("You must target a wall.")
	}
	if ch.minDist && g.Player.Pos.Distance(pos) <= 1 {
		return errors.New("You cannot target an adjacent wall.")
	}
	for _, pos := range ray[1:] {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			return errors.New("There are monsters in the way.")
		}
	}
	g.Player.Target = pos
	ch.done = true
	return nil
}

func (ch *wallChooser) Done() bool {
	return ch.done
}
