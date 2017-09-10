package main

import "errors"

type Targetter interface {
	ComputeHighlight(*game, position)
	Action(*game, position) error
	Reachable(*game, position) bool
	Done() bool
}

type examiner struct {
	done bool
}

func (ex *examiner) ComputeHighlight(g *game, pos position) {
	g.ComputeRayHighlight(pos)
}

func (ex *examiner) Action(g *game, pos position) error {
	if c := g.Dungeon.Cell(pos); c.Explored && c.T == FreeCell {
		g.AutoTarget = &pos
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
	done    bool
	area    bool
	minDist bool
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
	if c := g.Dungeon.Cell(pos); c.Explored && c.T == FreeCell {
		mons, _ := g.MonsterAt(pos)
		if mons != nil {
			g.Player.Target = pos
			ch.done = true
			return nil
		}
		if !ch.area {
			return errors.New("You must target a monster.")
		}
		ray := g.Ray(pos)
		tpos := pos
		for _, rpos := range ray {
			mons, _ := g.MonsterAt(rpos)
			if mons == nil {
				continue
			}
			tpos = mons.Pos
		}
		if tpos != pos {
			return errors.New("Invalid target: there are monsters in the way.")
		}
		neighbors := g.Dungeon.FreeNeighbors(pos)
		for _, npos := range neighbors {
			mons, _ := g.MonsterAt(npos)
			if mons != nil {
				g.Player.Target = pos
				ch.done = true
				return nil
			}
		}
		return errors.New("Invalid target: no monsters in the area.")
	}
	if !g.Dungeon.Cell(pos).Explored {
		return errors.New("You do not know what is in there.")
	}
	return errors.New("You cannot target a wall.")
}

func (ch *chooser) Done() bool {
	return ch.done
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
		mons, _ := g.MonsterAt(pos)
		if mons != nil {
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
