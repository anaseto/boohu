package main

import "time"

func (ui *gameui) SwappingAnimation(g *game, mpos, ppos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	_, fgm, bgColorm := ui.PositionDrawing(g, mpos)
	_, _, bgColorp := ui.PositionDrawing(g, ppos)
	ui.DrawAtPosition(g, mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(g, ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(g, ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
}

func (ui *gameui) TeleportAnimation(g *game, from, to position, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(g, from)
	_, _, bgColort := ui.PositionDrawing(g, to)
	ui.DrawAtPosition(g, from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	if showto {
		ui.DrawAtPosition(g, from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(g, to, true, 'Φ', ColorCyan, bgColort)
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

func (ui *gameui) ProjectileTrajectoryAnimation(g *game, ray []position, fg uicolor) {
	if DisableAnimations {
		return
	}
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, '•', fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
}

func (ui *gameui) MonsterProjectileAnimation(g *game, ray []position, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, or, fgColor, bgColor)
	}
}

func (ui *gameui) ExplosionAnimationAt(g *game, pos position, fg uicolor) {
	_, _, bgColor := ui.PositionDrawing(g, pos)
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
	//ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
	ui.DrawAtPosition(g, pos, true, r, bgColor, fg)
}

func (ui *gameui) ExplosionAnimation(g *game, es explosionStyle, pos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
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
			ui.ExplosionAnimationAt(g, npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

func (ui *gameui) TormentExplosionAnimation(g *game) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [3]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd, ColorFgMagicPlace}
	for i := 0; i < 3; i++ {
		for npos, b := range g.Player.LOS {
			if !b {
				continue
			}
			fg := colors[RandInt(3)]
			ui.ExplosionAnimationAt(g, npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

func (ui *gameui) WallExplosionAnimation(g *game, pos position) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(g, pos)
		//ui.DrawAtPosition(g, pos, true, '☼', fg, bgColor)
		ui.DrawAtPosition(g, pos, true, '☼', bgColor, fg)
		ui.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (ui *gameui) FireBoltAnimation(g *game, ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(g, pos)
			mons := g.MonsterAt(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			if mons.Exists() {
				r = '√'
			}
			//ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
			ui.DrawAtPosition(g, pos, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(25 * time.Millisecond)
}

func (ui *gameui) SlowingMagaraAnimation(g *game, ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgConfusedMonster, ColorFgMagicPlace}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(g, pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			ui.DrawAtPosition(g, pos, true, r, bgColor, fg)
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

func (ui *gameui) ThrowAnimation(g *game, ray []position, hit bool) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgProjectile, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
	if hit {
		pos := ray[0]
		ui.HitAnimation(g, pos, true)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *gameui) MonsterJavelinAnimation(g *game, ray []position, hit bool) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgMonster, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *gameui) HitAnimation(g *game, pos position, targeting bool) {
	if DisableAnimations {
		return
	}
	if !g.Player.LOS[pos] {
		return
	}
	ui.DrawDungeonView(g, NoFlushMode)
	_, _, bgColor := ui.PositionDrawing(g, pos)
	mons := g.MonsterAt(pos)
	if mons.Exists() || pos == g.Player.Pos {
		ui.DrawAtPosition(g, pos, targeting, '√', ColorFgAnimationHit, bgColor)
	} else {
		ui.DrawAtPosition(g, pos, targeting, '∞', ColorFgAnimationHit, bgColor)
	}
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
}

func (ui *gameui) LightningHitAnimation(g *game, targets []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 2; j++ {
		for _, pos := range targets {
			_, _, bgColor := ui.PositionDrawing(g, pos)
			mons := g.MonsterAt(pos)
			if mons.Exists() || pos == g.Player.Pos {
				ui.DrawAtPosition(g, pos, false, '√', bgColor, colors[RandInt(2)])
			} else {
				ui.DrawAtPosition(g, pos, false, '∞', bgColor, colors[RandInt(2)])
			}
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

func (ui *gameui) WoundedAnimation(g *game) {
	if DisableAnimations {
		return
	}
	r, _, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		time.Sleep(50 * time.Millisecond)
	}
}

func (ui *gameui) DrinkingPotionAnimation(g *game) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(50 * time.Millisecond)
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) StatusEndAnimation(g *game) {
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	time.Sleep(100 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *gameui) MenuSelectedAnimation(g *game, m menu, ok bool) {
	if DisableAnimations {
		return
	}
	if !ui.Small() {
		var message string
		if m == MenuInteract {
			message = ui.UpdateInteractButton(g)
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
		time.Sleep(25 * time.Millisecond)
		ui.DrawColoredText(m.String(), MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *gameui) MagicMappingAnimation(g *game, border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, false, r, fg, bg)
	}
	ui.Flush()
	time.Sleep(12 * time.Millisecond)
}
