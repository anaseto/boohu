package main

type raynode struct {
	Cost int
}

type rayMap map[position]raynode

func (g *game) bestParent(rm rayMap, from, pos position) (position, int) {
	p := pos.Parents(from)
	b := p[0]
	if len(p) > 1 && rm[p[1]].Cost+g.losCost(p[1]) < rm[b].Cost+g.losCost(b) {
		b = p[1]
	}
	return b, rm[b].Cost + g.losCost(b)
}

func (g *game) losCost(pos position) int {
	if g.Player.Pos == pos {
		return 0
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell {
		return g.LosRange()
	}
	if _, ok := g.Clouds[pos]; ok {
		return g.LosRange()
	}
	if _, ok := g.Doors[pos]; ok {
		if pos != g.Player.Pos {
			mons := g.MonsterAt(pos)
			if !mons.Exists() {
				return g.LosRange()
			}
		}
	}
	if _, ok := g.Fungus[pos]; ok {
		return g.LosRange() - 1
	}
	return 1
}

func (g *game) buildRayMap(from position, distance int) rayMap {
	rm := rayMap{}
	rm[from] = raynode{Cost: 0}
	for d := 1; d <= distance; d++ {
		for x := -d + from.X; x <= d+from.X; x++ {
			for _, pos := range []position{{x, from.Y + d}, {x, from.Y - d}} {
				if !pos.valid() {
					continue
				}
				_, c := g.bestParent(rm, from, pos)
				rm[pos] = raynode{Cost: c}
			}
		}
		for y := -d + 1 + from.Y; y <= d-1+from.Y; y++ {
			for _, pos := range []position{{from.X + d, y}, {from.X - d, y}} {
				if !pos.valid() {
					continue
				}
				_, c := g.bestParent(rm, from, pos)
				rm[pos] = raynode{Cost: c}
			}
		}
	}
	return rm
}

func (g *game) LosRange() int {
	losRange := 6
	if g.Player.Aptitudes[AptStealthyLOS] {
		losRange -= 2
	}
	return losRange
}

func (g *game) StopAuto() {
	g.AutoHalt = true
	g.AutoDir = NoDir
	g.AutoTarget = InvalidPos
	g.Resting = false
}

func (g *game) ComputeLOS() {
	m := map[position]bool{}
	losRange := g.LosRange()
	g.Player.Rays = g.buildRayMap(g.Player.Pos, losRange)
	for pos, n := range g.Player.Rays {
		if n.Cost < g.LosRange() {
			m[pos] = true
			g.SeePosition(pos)
		}
	}
	g.Player.LOS = m
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.Pos] {
			g.StopAuto()
			if mons.Seen {
				continue
			}
			mons.Seen = true
			g.Printf("You see %s (%v).", mons.Kind.Indefinite(false), mons.State)
			if mons.Kind.Dangerousness() > 10 {
				g.StoryPrint(mons.Kind.SeenStoryText())
			}
		}
	}
}

func (g *game) SeePosition(pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		see := "see"
		if c, ok := g.Collectables[pos]; ok {
			g.StopAuto()
			if c.Quantity > 1 {
				g.Printf("You %s %d %s.", see, c.Quantity, c.Consumable.Plural())
			} else {
				g.Printf("You %s %s.", see, Indefinite(c.Consumable.String(), false))
			}
		} else if _, ok := g.Stairs[pos]; ok {
			g.StopAuto()
			g.Printf("You %s stairs.", see)
		} else if eq, ok := g.Equipables[pos]; ok {
			g.StopAuto()
			g.Printf("You %s %s.", see, Indefinite(eq.String(), false))
		} else if rod, ok := g.Rods[pos]; ok {
			g.StopAuto()
			g.Printf("You %s %s.", see, Indefinite(rod.String(), false))
		}
		g.FairAction()
		g.Dungeon.SetExplored(pos)
		g.DijkstraMapRebuild = true
	} else {
		if g.UnknownDig[pos] {
			g.StopAuto()
			g.Printf("There is no more a wall there.")
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[pos]; ok && cld == CloudFire && g.UnknownBurn[pos] {
			g.StopAuto()
			g.Printf("There are flames there.")
			g.DijkstraMapRebuild = true
		}
	}
	if g.UnknownDig[pos] {
		delete(g.UnknownDig, pos)
		delete(g.TemporalWalls, pos)
	}
	if g.UnknownBurn[pos] {
		delete(g.UnknownBurn, pos)
	}
}

func (g *game) ComputeExclusion(pos position, toggle bool) {
	exclusionRange := g.LosRange()
	g.ExclusionsMap[pos] = toggle
	for d := 1; d <= exclusionRange; d++ {
		for x := -d + pos.X; x <= d+pos.X; x++ {
			for _, pos := range []position{{x, pos.Y + d}, {x, pos.Y - d}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
		for y := -d + 1 + pos.Y; y <= d-1+pos.Y; y++ {
			for _, pos := range []position{{pos.X + d, y}, {pos.X - d, y}} {
				if !pos.valid() {
					continue
				}
				g.ExclusionsMap[pos] = toggle
			}
		}
	}
}

func (g *game) Ray(pos position) []position {
	if !g.Player.LOS[pos] {
		return nil
	}
	ray := []position{}
	for pos != g.Player.Pos {
		ray = append(ray, pos)
		pos, _ = g.bestParent(g.Player.Rays, g.Player.Pos, pos)
	}
	return ray
}

func (g *game) ComputeRayHighlight(pos position) {
	g.Highlight = map[position]bool{}
	ray := g.Ray(pos)
	for _, p := range ray {
		g.Highlight[p] = true
	}
}

func (g *game) ComputeNoise() {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, g.LosRange()+2)
	count := 0
	noise := map[position]bool{}
	for pos := range nm {
		if g.Player.LOS[pos] {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && mons.State != Resting && RandInt(3) == 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsGiantBee, MonsSatowalgaPlant:
				// no footsteps
			default:
				noise[pos] = true
				count++
			}
		}
	}
	if count > 0 {
		if count > 1 {
			g.Printf("You hear some footsteps (%dx).", count)
		} else {
			g.Print("You hear some footsteps.")
		}
		g.StopAuto()
	}
	g.Noise = noise
}
