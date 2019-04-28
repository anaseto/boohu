package main

// TODO: many animations are obsolete so remove them

import (
	"sort"
	"time"
)

const (
	AnimDurShort       = 25
	AnimDurShortMedium = 50
	AnimDurMedium      = 75
	AnimDurMediumLong  = 100
	AnimDurLong        = 200
	AnimDurExtraLong   = 300
)

func (ui *gameui) SwappingAnimation(mpos, ppos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	_, fgm, bgColorm := ui.PositionDrawing(mpos)
	_, _, bgColorp := ui.PositionDrawing(ppos)
	ui.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	time.Sleep(AnimDurMedium)
	ui.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	time.Sleep(AnimDurMedium)
}

func (ui *gameui) TeleportAnimation(from, to position, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(from)
	_, _, bgColort := ui.PositionDrawing(to)
	ui.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	time.Sleep(AnimDurMediumLong)
	if showto {
		ui.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		time.Sleep(AnimDurMedium)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

//func (ui *gameui) ProjectileTrajectoryAnimation(ray []position, fg uicolor) {
//if DisableAnimations {
//return
//}
//for i := len(ray) - 1; i >= 0; i-- {
//pos := ray[i]
//r, fgColor, bgColor := ui.PositionDrawing(pos)
//ui.DrawAtPosition(pos, true, '•', fg, bgColor)
//ui.Flush()
//time.Sleep(30 * time.Millisecond)
//ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
//}
//}

func (ui *gameui) MonsterProjectileAnimation(ray []position, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, r, fg, bgColor)
		ui.Flush()
		time.Sleep(AnimDurMedium)
		ui.DrawAtPosition(pos, true, or, fgColor, bgColor)
	}
}

func (ui *gameui) WaveDrawAt(pos position, fg uicolor) {
	r, _, bgColor := ui.PositionDrawing(pos)
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *gameui) ExplosionDrawAt(pos position, fg uicolor) {
	g := ui.g
	_, _, bgColor := ui.PositionDrawing(pos)
	mons := g.MonsterAt(pos)
	r := ';'
	switch RandInt(9) {
	case 0, 6:
		r = ','
	case 1:
		r = '}'
	case 2:
		r = '%'
	case 3, 7:
		r = ':'
	case 4:
		r = '\\'
	case 5:
		r = '~'
	}
	if mons.Exists() || g.Player.Pos == pos {
		r = '√'
	}
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *gameui) NoiseAnimation(noises []position) {
	if DisableAnimations {
		return
	}
	ui.LOSWavesAnimation(DefaultLOSRange, WaveMagicNoise, ui.g.Player.Pos)
	colors := []uicolor{ColorFgSleepingMonster, ColorFgMagicPlace}
	for i := 0; i < 2; i++ {
		for _, pos := range noises {
			r := '♫'
			_, _, bgColor := ui.PositionDrawing(pos)
			ui.DrawAtPosition(pos, false, r, bgColor, colors[i])
		}
		_, _, bgColor := ui.PositionDrawing(ui.g.Player.Pos)
		ui.DrawAtPosition(ui.g.Player.Pos, false, '@', bgColor, colors[i])
		ui.Flush()
		time.Sleep(AnimDurShortMedium)
	}

}

func (ui *gameui) ExplosionAnimation(es explosionStyle, pos position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for i := 0; i < 3; i++ {
		nb := g.Dungeon.FreeNeighbors(pos)
		if es != AroundWallExplosion {
			nb = append(nb, pos)
		}
		for _, npos := range nb {
			fg := colors[RandInt(2)]
			if !g.Player.LOS[npos] {
				continue
			}
			ui.ExplosionDrawAt(npos, fg)
		}
		ui.Flush()
		time.Sleep(AnimDurMediumLong)
	}
}

func (g *game) Waves(maxCost int, ws wavestyle, center position) (dists []int, cdists map[int][]int) {
	var dij Dijkstrer
	switch ws {
	case WaveMagicNoise:
		dij = &gridPath{dungeon: g.Dungeon}
	default:
		dij = &noisePath{game: g}
	}
	nm := Dijkstra(dij, []position{center}, maxCost)
	cdists = make(map[int][]int)
	nm.iter(g.Player.Pos, func(n *node) {
		pos := n.Pos
		cdists[n.Cost] = append(cdists[n.Cost], pos.idx())
	})
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	return dists, cdists
}

func (ui *gameui) LOSWavesAnimation(r int, ws wavestyle, center position) {
	dists, cdists := ui.g.Waves(r, ws, center)
	for _, d := range dists {
		wave := cdists[d]
		if len(wave) == 0 {
			break
		}
		ui.WaveAnimation(wave, ws)
	}
}

type wavestyle int

const (
	WaveMagicNoise wavestyle = iota
	WaveNoise
	WaveConfusion
	WaveSlowing
)

func (ui *gameui) WaveAnimation(wave []int, ws wavestyle) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	//colors := [2]uicolor{ColorFgMagicPlace, ColorFgSleepingMonster}
	for _, i := range wave {
		pos := idxtopos(i)
		switch ws {
		case WaveConfusion:
			fg := ColorFgConfusedMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveSlowing:
			fg := ColorFgSlowedMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveNoise:
			fg := ColorFgWanderingMonster
			if ui.g.Player.Sees(pos) {
				ui.WaveDrawAt(pos, fg)
			}
		case WaveMagicNoise:
			fg := ColorFgMagicPlace
			ui.WaveDrawAt(pos, fg)
		}
	}
	ui.Flush()
	time.Sleep(AnimDurShort)
}

//func (ui *gameui) TormentExplosionAnimation() {
//g := ui.g
//if DisableAnimations {
//return
//}
//ui.DrawDungeonView(NormalMode)
//time.Sleep(AnimDurShort)
//colors := [3]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd, ColorFgMagicPlace}
//for i := 0; i < 3; i++ {
//for npos, b := range g.Player.LOS {
//if !b {
//continue
//}
//fg := colors[RandInt(3)]
//ui.ExplosionDrawAt(npos, fg)
//}
//ui.Flush()
//time.Sleep(AnimDurMediumLong)
//}
//}

func (ui *gameui) WallExplosionAnimation(pos position) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(pos)
		//ui.DrawAtPosition(pos, true, '☼', fg, bgColor)
		ui.DrawAtPosition(pos, true, '%', bgColor, fg)
		ui.Flush()
		time.Sleep(AnimDurShort)
	}
}

type beamstyle int

const (
	BeamSleeping beamstyle = iota
	BeamLignification
	BeamObstruction
)

func (ui *gameui) BeamsAnimation(ray []position, bs beamstyle) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	// change colors depending on effect
	var fg uicolor
	switch bs {
	case BeamSleeping:
		fg = ColorFgSleepingMonster
	case BeamLignification:
		fg = ColorFgLignifiedMonster
	case BeamObstruction:
		fg = ColorFgMagicPlace
	}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(pos)
			mons := g.MonsterAt(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			if mons.Exists() {
				r = '√'
			}
			//ui.DrawAtPosition(pos, true, r, fg, bgColor)
			ui.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(AnimDurMediumLong)
	}
}

//func (ui *gameui) FireBoltAnimation(ray []position) {
//g := ui.g
//if DisableAnimations {
//return
//}
//ui.DrawDungeonView(NormalMode)
//time.Sleep(AnimDurShort)
//colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
//for j := 0; j < 3; j++ {
//for i := len(ray) - 1; i >= 0; i-- {
//fg := colors[RandInt(2)]
//pos := ray[i]
//_, _, bgColor := ui.PositionDrawing(pos)
//mons := g.MonsterAt(pos)
//r := '*'
//if RandInt(2) == 0 {
//r = '×'
//}
//if mons.Exists() {
//r = '√'
//}
////ui.DrawAtPosition(pos, true, r, fg, bgColor)
//ui.DrawAtPosition(pos, true, r, bgColor, fg)
//}
//ui.Flush()
//time.Sleep(AnimDurMediumLong)
//}
//}

func (ui *gameui) SlowingMagaraAnimation(ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	colors := [2]uicolor{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			ui.DrawAtPosition(pos, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(AnimDurMediumLong)
	}
}

func (ui *gameui) ProjectileSymbol(dir direction) (r rune) {
	switch dir {
	case E, ENE, ESE, WNW, W, WSW:
		r = '—'
	case NE, SW:
		r = '/'
	case NNE, N, NNW, SSW, S, SSE:
		r = '|'
	case NW, SE:
		r = '\\'
	}
	return r
}

//func (ui *gameui) ThrowAnimation(ray []position, hit bool) {
//g := ui.g
//if DisableAnimations {
//return
//}
//ui.DrawDungeonView(NormalMode)
//time.Sleep(AnimDurShort)
//for i := len(ray) - 1; i >= 0; i-- {
//pos := ray[i]
//r, fgColor, bgColor := ui.PositionDrawing(pos)
//ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgProjectile, bgColor)
//ui.Flush()
//time.Sleep(30 * time.Millisecond)
//ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
//}
//if hit {
//pos := ray[0]
//ui.HitAnimation(pos, true)
//}
//time.Sleep(30 * time.Millisecond)
//}

func (ui *gameui) MonsterJavelinAnimation(ray []position, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShort)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgMonster, bgColor)
		ui.Flush()
		time.Sleep(AnimDurShort)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	time.Sleep(AnimDurShort)
}

//func (ui *gameui) HitAnimation(pos position, targeting bool) {
//g := ui.g
//if DisableAnimations {
//return
//}
//if !g.Player.LOS[pos] {
//return
//}
//ui.DrawDungeonView(NoFlushMode)
//_, _, bgColor := ui.PositionDrawing(pos)
//mons := g.MonsterAt(pos)
//if mons.Exists() || pos == g.Player.Pos {
//ui.DrawAtPosition(pos, targeting, '√', ColorFgAnimationHit, bgColor)
//} else {
//ui.DrawAtPosition(pos, targeting, '∞', ColorFgAnimationHit, bgColor)
//}
//ui.Flush()
//time.Sleep(AnimDurShortMedium)
//}

//func (ui *gameui) LightningHitAnimation(targets []position) {
//g := ui.g
//if DisableAnimations {
//return
//}
//ui.DrawDungeonView(NormalMode)
//time.Sleep(AnimDurShort)
//colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
//for j := 0; j < 2; j++ {
//for _, pos := range targets {
//_, _, bgColor := ui.PositionDrawing(pos)
//mons := g.MonsterAt(pos)
//if mons.Exists() || pos == g.Player.Pos {
//ui.DrawAtPosition(pos, false, '√', bgColor, colors[RandInt(2)])
//} else {
//ui.DrawAtPosition(pos, false, '∞', bgColor, colors[RandInt(2)])
//}
//}
//ui.Flush()
//time.Sleep(AnimDurMediumLong)
//}
//}

func (ui *gameui) WoundedAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	r, _, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(AnimDurShortMedium)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		time.Sleep(AnimDurShortMedium)
	}
}

func (ui *gameui) PlayerGoodEffectAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(AnimDurShortMedium)
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(AnimDurMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(AnimDurMedium)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) StatusEndAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	time.Sleep(AnimDurMediumLong)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) PushAnimation(path []position) {
	if DisableAnimations {
		return
	}
	if len(path) == 0 {
		// should not happen
		return
	}
	_, _, bg := ui.PositionDrawing(path[0])
	for _, pos := range path[:len(path)-1] {
		ui.DrawAtPosition(pos, false, '×', ColorFgPlayer, bg)
	}
	ui.DrawAtPosition(path[len(path)-1], false, '@', ColorFgPlayer, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
}

func (ui *gameui) MenuSelectedAnimation(m menu, ok bool) {
	if DisableAnimations {
		return
	}
	if !ui.Small() {
		var message string
		if m == MenuInteract {
			message = ui.UpdateInteractButton()
		} else {
			message = m.String()
		}
		if message == "" {
			return
		}
		if ok {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorCyan)
		} else {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorMagenta)
		}
		ui.Flush()
		var t time.Duration = 25
		if !ok {
			t += 25
		}
		time.Sleep(t * time.Millisecond)
		ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *gameui) MagicMappingAnimation(border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, false, r, fg, bg)
	}
	ui.Flush()
}

func (ui *gameui) FreeingShaedraAnimation() {
	g := ui.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.Print("You see Shaedra. She is wounded!")
	g.PrintStyled("Shaedra: “Oh, it's you, Syu! Let's flee with Marevor's magara!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	time.Sleep(AnimDurExtraLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “And what about the mission?”", logSpecial)
	g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawDungeonView(NoFlushMode)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.DrawAtPosition(g.Places.Shaedra, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
	g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
	g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	g.Player.Magaras = append(g.Player.Magaras, NoMagara)
	g.Player.Inventory.Misc = NoItem
	g.PrintStyled("You have a new empty slot for a magara.", logSpecial)
	AchRescuedShaedra.Get(g)
}

func (ui *gameui) TakingArtifactAnimation() {
	g := ui.g
	//if DisableAnimations {
	// TODO this animation cannot be disabled as-is, because code is mixed with it...
	//return
	//}
	g.Print("You take and use the artifact.")
	g.Dungeon.SetCell(g.Places.Artifact, GroundCell)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	time.Sleep(AnimDurExtraLong)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “Great! Let's escape and find some bones to celebrate!”", logSpecial)
	g.PrintStyled("Syu: “Sorry, but I prefer bananas!”", logSpecial)
	g.Print("[(x) to continue]")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawDungeonView(NoFlushMode)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(AnimDurLong)
	g.Dungeon.SetCell(g.Places.Marevor, GroundCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	AchRetrievedArtifact.Get(g)
}
