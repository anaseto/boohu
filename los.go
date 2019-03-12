package main

type raynode struct {
	Cost int
}

type rayMap map[position]raynode

func (g *game) bestParent(rm rayMap, from, pos position, losrange int) (position, int) {
	p := pos.Parents(from)
	b := p[0]
	if len(p) > 1 && rm[p[1]].Cost+g.losCost(from, p[1], pos, losrange) < rm[b].Cost+g.losCost(from, b, pos, losrange) {
		b = p[1]
	}
	return b, rm[b].Cost + g.losCost(from, b, pos, losrange)
}

func (g *game) DiagonalWall(from, to position) bool {
	p := make([]position, 0, 2)
	switch to.Dir(from) {
	case NE:
		p = append(p, to.S(), to.W())
	case NW:
		p = append(p, to.S(), to.E())
	case SW:
		p = append(p, to.N(), to.E())
	case SE:
		p = append(p, to.N(), to.W())
	}
	count := 0
	for _, pos := range p {
		if pos.valid() && g.Dungeon.Cell(pos).T == WallCell {
			count++
		}
	}
	return count > 1
}

func (g *game) losCost(from, pos, to position, losrange int) int {
	if g.DiagonalWall(pos, to) {
		return losrange
	}
	if from == pos {
		return to.Distance(pos) - 1
	}
	c := g.Dungeon.Cell(pos)
	if c.T == WallCell {
		return losrange
	}
	if _, ok := g.Clouds[pos]; ok {
		return losrange
	}
	if c.T == DoorCell {
		if pos != from {
			mons := g.MonsterAt(pos)
			if !mons.Exists() {
				return losrange
			}
		}
	}
	if c.T == FungusCell {
		return losrange - 2
	}
	return to.Distance(pos)
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
				_, c := g.bestParent(rm, from, pos, distance)
				rm[pos] = raynode{Cost: c}
			}
		}
		for y := -d + 1 + from.Y; y <= d-1+from.Y; y++ {
			for _, pos := range []position{{from.X + d, y}, {from.X - d, y}} {
				if !pos.valid() {
					continue
				}
				_, c := g.bestParent(rm, from, pos, distance)
				rm[pos] = raynode{Cost: c}
			}
		}
	}
	return rm
}

const DefaultLOSRange = 10
const DefaultMonsterLOSRange = 7

func (g *game) LosRange() int {
	losRange := DefaultLOSRange
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
	losRange := DefaultMonsterLOSRange
	if losRange < 1 {
		losRange = 1
	}
	rays := g.buildRayMap(m.Pos, losRange)
	for pos, n := range rays {
		if n.Cost < losRange && g.Dungeon.Cell(pos).T != BarrelCell {
			mlos[pos] = true
		}
	}
	m.LOS = mlos
}

func (g *game) SeePosition(pos position) {
	c := g.Dungeon.Cell(pos)
	t, okT := g.TerrainKnowledge[pos]
	if !c.Explored {
		see := "see"
		if c.IsNotable() {
			g.Printf("You %s %s.", see, c.ShortDesc(g, pos))
			g.StopAuto()
		}
		g.FunAction()
		g.Dungeon.SetExplored(pos)
		g.DijkstraMapRebuild = true
	} else {
		if okT && t == WallCell && c.T != WallCell {
			g.Printf("There is no longer a wall there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
		if cld, ok := g.Clouds[pos]; ok && cld == CloudFire && okT && (t == FungusCell || t == DoorCell) {
			g.Printf("There are flames there.")
			g.StopAuto()
			g.DijkstraMapRebuild = true
		}
	}
	if okT {
		delete(g.TerrainKnowledge, pos)
		if c.IsFree() {
			delete(g.TemporalWalls, pos)
		}
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
		pos, _ = g.bestParent(g.Player.Rays, g.Player.Pos, pos, g.LosRange())
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
	rg := DefaultLOSRange
	if g.Player.Aptitudes[AptHear] {
		rg++
	}
	nm := Dijkstra(dij, []position{g.Player.Pos}, rg)
	count := 0
	noise := map[position]bool{}
	rmax := 2
	if g.Player.Aptitudes[AptHear] {
		rmax += 2
	}
	for pos := range nm {
		if g.Player.Sees(pos) {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() && mons.State != Resting && RandInt(rmax) > 0 {
			switch mons.Kind {
			case MonsMirrorSpecter, MonsSatowalgaPlant:
				// no footsteps
			//case MonsTinyHarpy, MonsWingedMilfid, MonsGiantBee:
			//noise[pos] = true
			//g.Print("You hear the flapping of wings.")
			//count++
			//case MonsOgre, MonsCyclop, MonsBrizzia, MonsHydra, MonsEarthDragon, MonsTreeMushroom:
			//noise[pos] = true
			//g.Print("You hear heavy footsteps.")
			//count++
			//case MonsWorm, MonsAcidMound:
			case MonsWorm:
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
	return m.Sees(g, g.Player.Pos)
}

func (m *monster) Sees(g *game, pos position) bool {
	return m.LOS[pos] && m.Dir.InViewCone(m.Pos, pos)
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
