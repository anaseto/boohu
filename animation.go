package main

// TODO: many animations are obsolete so remove them

import (
	"sort"
	"time"
)

func (ui *gameui) SwappingAnimation(mpos, ppos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	_, fgm, bgColorm := ui.PositionDrawing(mpos)
	_, _, bgColorp := ui.PositionDrawing(ppos)
	ui.DrawAtPosition(mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
}

func (ui *gameui) TeleportAnimation(from, to position, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(from)
	_, _, bgColort := ui.PositionDrawing(to)
	ui.DrawAtPosition(from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	if showto {
		ui.DrawAtPosition(from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *gameui) ProjectileTrajectoryAnimation(ray []position, fg uicolor) {
	if DisableAnimations {
		return
	}
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, '•', fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
}

func (ui *gameui) MonsterProjectileAnimation(ray []position, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, r, fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(pos, true, or, fgColor, bgColor)
	}
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
	//ui.DrawAtPosition(pos, true, r, fg, bgColor)
	ui.DrawAtPosition(pos, true, r, bgColor, fg)
}

func (ui *gameui) NoiseAnimation(noises []position) {
	if DisableAnimations {
		return
	}
	ui.LOSWavesAnimation(DefaultLOSRange, WaveNoise)
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
		time.Sleep(50 * time.Millisecond)
	}

}

func (ui *gameui) ExplosionAnimation(es explosionStyle, pos position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(20 * time.Millisecond)
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
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

func (g *game) Waves(maxCost int) (dists []int, cdists map[int][]int) {
	dij := &noisePath{game: g}
	nm := Dijkstra(dij, []position{g.Player.Pos}, maxCost)
	cdists = make(map[int][]int)
	for pos, n := range nm {
		cdists[n.Cost] = append(cdists[n.Cost], pos.idx())
	}
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	return dists, cdists
}

func (ui *gameui) LOSWavesAnimation(r int, ws wavestyle) {
	dists, cdists := ui.g.Waves(r)
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
	WaveLOS wavestyle = iota
	WaveNoise
)

func (ui *gameui) WaveAnimation(wave []int, ws wavestyle) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	colors := [2]uicolor{ColorFgMagicPlace, ColorFgSleepingMonster}
	for _, i := range wave {
		pos := idxtopos(i)
		switch ws {
		case WaveLOS:
			fg := colors[RandInt(2)]
			if ui.g.Player.Sees(pos) {
				ui.ExplosionDrawAt(pos, fg)
			}
		case WaveNoise:
			fg := colors[RandInt(2)]
			ui.ExplosionDrawAt(pos, fg)
		}
	}
	ui.Flush()
	time.Sleep(25 * time.Millisecond)
}

func (ui *gameui) TormentExplosionAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [3]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd, ColorFgMagicPlace}
	for i := 0; i < 3; i++ {
		for npos, b := range g.Player.LOS {
			if !b {
				continue
			}
			fg := colors[RandInt(3)]
			ui.ExplosionDrawAt(npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

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
		time.Sleep(25 * time.Millisecond)
	}
}

func (ui *gameui) BeamsAnimation(ray []position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	// change colors depending on effect
	colors := [2]uicolor{ColorFgSleepingMonster, ColorFgSlowedMonster}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
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
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(25 * time.Millisecond)
}

func (ui *gameui) FireBoltAnimation(ray []position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
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
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(25 * time.Millisecond)
}

func (ui *gameui) SlowingMagaraAnimation(ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
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
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(25 * time.Millisecond)
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

func (ui *gameui) ThrowAnimation(ray []position, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgProjectile, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	if hit {
		pos := ray[0]
		ui.HitAnimation(pos, true)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *gameui) MonsterJavelinAnimation(ray []position, hit bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(pos)
		ui.DrawAtPosition(pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgMonster, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(pos, true, r, fgColor, bgColor)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *gameui) HitAnimation(pos position, targeting bool) {
	g := ui.g
	if DisableAnimations {
		return
	}
	if !g.Player.LOS[pos] {
		return
	}
	ui.DrawDungeonView(NoFlushMode)
	_, _, bgColor := ui.PositionDrawing(pos)
	mons := g.MonsterAt(pos)
	if mons.Exists() || pos == g.Player.Pos {
		ui.DrawAtPosition(pos, targeting, '√', ColorFgAnimationHit, bgColor)
	} else {
		ui.DrawAtPosition(pos, targeting, '∞', ColorFgAnimationHit, bgColor)
	}
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
}

func (ui *gameui) LightningHitAnimation(targets []position) {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 2; j++ {
		for _, pos := range targets {
			_, _, bgColor := ui.PositionDrawing(pos)
			mons := g.MonsterAt(pos)
			if mons.Exists() || pos == g.Player.Pos {
				ui.DrawAtPosition(pos, false, '√', bgColor, colors[RandInt(2)])
			} else {
				ui.DrawAtPosition(pos, false, '∞', bgColor, colors[RandInt(2)])
			}
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) WoundedAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, _, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		time.Sleep(50 * time.Millisecond)
	}
}

func (ui *gameui) PlayerGoodEffectAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(NormalMode)
	time.Sleep(50 * time.Millisecond)
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) StatusEndAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g.Player.Pos)
	ui.DrawAtPosition(g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	time.Sleep(100 * time.Millisecond)
	ui.DrawAtPosition(g.Player.Pos, false, r, fg, bg)
	ui.Flush()
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
		ui.DrawColoredText(m.String(), MenuCols[m][0], DungeonHeight, ColorViolet)
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
	time.Sleep(12 * time.Millisecond)
}

func (ui *gameui) FreeingShaedraAnimation() {
	g := ui.g
	if DisableAnimations {
		return
	}
	g.Print("You see Shaedra. She is wounded!")
	g.PrintStyled("Shaedra: “By Ruyale, thank you Syu! Let's flee with Marevor's magara!”", logSpecial)
	g.Print("[esc|space]...")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	_, _, bg := ui.PositionDrawing(g.Places.Monolith)
	ui.DrawAtPosition(g.Places.Monolith, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(400 * time.Millisecond)
	g.Objects.Stairs[g.Places.Monolith] = WinStair
	g.Dungeon.SetCell(g.Places.Monolith, StairCell)
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	time.Sleep(400 * time.Millisecond)
	_, _, bg = ui.PositionDrawing(g.Places.Marevor)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.Flush()
	time.Sleep(400 * time.Millisecond)
	g.Dungeon.SetCell(g.Places.Marevor, StoryCell)
	g.Objects.Story[g.Places.Marevor] = StoryMarevor
	g.PrintStyled("Marevor: “And what about the mission?”", logSpecial)
	g.PrintStyled("Shaedra: “Pff, don't be reckless!”", logSpecial)
	g.Print("[esc/space]...")
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	ui.WaitForContinue(-1)
	ui.DrawAtPosition(g.Places.Marevor, false, 'Φ', ColorFgMagicPlace, bg)
	ui.DrawAtPosition(g.Places.Shaedra, false, 'Φ', ColorFgMagicPlace, bg)
	time.Sleep(400 * time.Millisecond)
	ui.Flush()
	g.Dungeon.SetCell(g.Places.Shaedra, GroundCell)
	g.Dungeon.SetCell(g.Places.Marevor, ScrollCell)
	g.Objects.Scrolls[g.Places.Marevor] = ScrollExtended
	ui.DrawDungeonView(NoFlushMode)
	ui.Flush()
	time.Sleep(12 * time.Millisecond)
}
