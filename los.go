package main

type raynode struct {
	Cost int
}

type rayMap map[position]raynode

func (g *game) bestParent(rm rayMap, from, pos position) (position, int) {
	p := pos.Parents(from)
	b := p[0]
	if len(p) > 1 && rm[p[1]].Cost+g.losCost(from, p[1]) < rm[b].Cost+g.losCost(from, b) {
		b = p[1]
	}
	return b, rm[b].Cost + g.losCost(from, b)
}

func (g *game) losCost(from, pos position) int {
	if from == pos {
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
		if pos != from {
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
	losRange := 7
	if g.Player.Aptitudes[AptStealthyLOS] {
		losRange -= 2
	}
	if g.Player.Armour == HarmonistRobe {
		losRange -= 1
	}
	if g.Player.Weapon == Frundis {
		losRange -= 1
	}
	if g.Player.HasStatus(StatusShadows) {
		losRange = 1
	}
	if losRange < 1 {
		losRange = 1
	}
	return losRange
}

func (g *game) StopAuto() {
	if g.Autoexploring && !g.AutoHalt {
		g.Print("You stop exploring.")
	} else if g.AutoDir != NoDir {
		g.Print("You stop.")
	} else if g.AutoTarget != InvalidPos {
		g.Print("You stop.")
	}
	g.AutoHalt = true
	g.AutoDir = NoDir
	g.AutoTarget = InvalidPos
	if g.Resting {
		g.Stats.RestInterrupt++
		g.Resting = false
		g.Print("You could not sleep.")
	}
}

func (g *game) ComputeLOS() {
	m := map[position]bool{}
	losRange := g.LosRange()
	g.Player.Rays = g.buildRayMap(g.Player.Pos, losRange)
	for pos, n := range g.Player.Rays {
		if n.Cost < g.LosRange() {
			m[pos] = true
		}
	}
	g.Player.LOS = m
	for pos := range g.Player.LOS {
		if g.Player.Sees(pos) {
			g.SeePosition(pos)
		}
	}
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.Sees(mons.Pos) {
			mons.UpdateKnowledge(g, mons.Pos)
			if mons.Seen {
				g.StopAuto()
				continue
			}
			mons.Seen = true
			g.Printf("You see %s (%v).", mons.Kind.Indefinite(false), mons.State)
			if mons.Kind.Dangerousness() > 10 {
				g.StoryPrint(mons.Kind.SeenStoryText())
			}
			g.StopAuto()
		}
	}
}

func (m *monster) ComputeLOS(g *game) {
	mlos := map[position]bool{}
	losRange := g.LosRange()
	rays := g.buildRayMap(m.Pos, losRange)
	for pos, n := range rays {
		if n.Cost < g.LosRange() {
			mlos[pos] = true
		}
	}
	m.LOS = mlos
}

func (g *game) SeePosition(pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		see := "see"
		if c, ok := g.Collectables[pos]; ok {
			if c.Quantity > 1 {
				g.Printf("You %s %d %s.", see, c.Quantity, c.Consumable.Plural())
			} else {
				g.Printf("You %s %s.", see, Indefinite(c.Consumable.String(), false))
			}
			g.StopAuto()
		} else if _, ok := g.Stairs[pos]; ok {
			g.Printf("You %s stairs.", see)
			g.StopAuto()
		} else if eq, ok := g.Equipables[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(eq.String(), false))
			g.StopAuto()
		} else if rd, ok := g.Rods[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(rd.String(), false))
			g.StopAuto()
		} else if stn, ok := g.MagicalStones[pos]; ok {
			g.Printf("You %s %s.", see, Indefinite(stn.String(), false))
			g.StopAuto()
		}
		g.FunAction()
		g.Dungeon.SetExplored(pos)
		g.DijkstraMapRebuild = true
	} else {
		if g.WrongWall[pos] {
			g.Printf("There is no longer a wall there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[pos]; ok && cld == CloudFire && (g.WrongDoor[pos] || g.WrongFoliage[pos]) {
			g.Printf("There are flames there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
	}
	if g.WrongWall[pos] {
		delete(g.WrongWall, pos)
		if g.Dungeon.Cell(pos).T == FreeCell {
			delete(g.TemporalWalls, pos)
		}
	}
	if _, ok := g.WrongDoor[pos]; ok {
		delete(g.WrongDoor, pos)
	}
	if _, ok := g.WrongFoliage[pos]; ok {
		delete(g.WrongFoliage, pos)
	}
	if mons, ok := g.LastMonsterKnownAt[pos]; ok && (mons.Pos != pos || !mons.Exists()) {
		delete(g.LastMonsterKnownAt, pos)
		mons.LastKnownPos = InvalidPos
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
	rg := g.LosRange() + 2
	if rg <= 5 {
		rg++
	}
	if g.Player.Aptitudes[AptHear] {
		rg++
	}
	nm := Dijkstra(dij, []position{g.Player.Pos}, rg)
	count := 0
	noise := map[position]bool{}
	rmax := 3
	if g.Player.Aptitudes[AptHear] {
		rmax--
	}
	for pos := range nm {
		if g.Player.Sees(pos) {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && mons.State != Resting && RandInt(rmax) == 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsSatowalgaPlant:
				// no footsteps
			case MonsTinyHarpy, MonsWingedMilfid, MonsGiantBee:
				noise[pos] = true
				g.Print("You hear the flapping of wings.")
				count++
			case MonsOgre, MonsCyclop, MonsBrizzia, MonsHydra, MonsEarthDragon, MonsTreeMushroom:
				noise[pos] = true
				g.Print("You hear heavy footsteps.")
				count++
			case MonsWorm, MonsAcidMound:
				noise[pos] = true
				g.Print("You hear a creep noise.")
				count++
			default:
				noise[pos] = true
				g.Print("You hear footsteps.")
				count++
			}
		}
	}
	if count > 0 {
		g.StopAuto()
	}
	g.Noise = noise
}

func (p *player) Sees(pos position) bool {
	//return pos == p.Pos || p.LOS[pos] && p.Dir.InViewCone(p.Pos, pos)
	return p.LOS[pos]
}

func (m *monster) SeesPlayer(g *game) bool {
	return g.Player.LOS[m.Pos] && m.Dir.InViewCone(m.Pos, g.Player.Pos)
}

func (m *monster) Sees(g *game, pos position) bool {
	return pos == m.Pos || m.LOS[pos] && m.Dir.InViewCone(m.Pos, pos)
}

func (g *game) ComputeMonsterLOS() {
	g.MonsterLOS = make(map[position]bool)
	for _, mons := range g.Monsters {
		if !mons.Exists() || !g.Player.Sees(mons.Pos) {
			continue
		}
		for pos, _ := range g.Player.LOS {
			if !g.Player.Sees(pos) {
				continue
			}
			if mons.Sees(g, pos) {
				g.MonsterLOS[pos] = true
			}
		}
		// unoptimized version for testing:
		//if !mons.Exists() {
		//continue
		//}
		//for pos := range mons.LOS {
		//if mons.Sees(g, pos) {
		//g.MonsterLOS[pos] = true
		//}
		//}
	}
}

func (m *monster) UpdateKnowledge(g *game, pos position) {
	if mons, ok := g.LastMonsterKnownAt[pos]; ok {
		mons.LastKnownPos = InvalidPos
	}
	if m.LastKnownPos != InvalidPos {
		delete(g.LastMonsterKnownAt, m.LastKnownPos)
	}
	g.LastMonsterKnownAt[pos] = m
	m.LastSeenState = m.State
	m.LastKnownPos = pos
}
