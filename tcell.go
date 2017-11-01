// +build tcell

package main

import (
	"errors"

	"github.com/gdamore/tcell"
)

type termui struct {
	tcell.Screen
}

func FixColor() {
	ColorBgLOS = ColorBgLOS - 1
	ColorBgDark = ColorBgDark - 1
	ColorBg = ColorBg - 1
	ColorBgCloud = ColorBgCloud - 1
	ColorFgLOS = ColorFgLOS - 1
	ColorFgDark = ColorFgDark - 1
	ColorFg = ColorFg - 1
	ColorFgPlayer = ColorFgPlayer - 1
	ColorFgMonster = ColorFgMonster - 1
	ColorFgSleepingMonster = ColorFgSleepingMonster - 1
	ColorFgWanderingMonster = ColorFgWanderingMonster - 1
	ColorFgConfusedMonster = ColorFgConfusedMonster - 1
	ColorFgCollectable = ColorFgCollectable - 1
	ColorFgStairs = ColorFgStairs - 1
	ColorFgGold = ColorFgGold - 1
	ColorFgHPok = ColorFgHPok - 1
	ColorFgHPwounded = ColorFgHPwounded - 1
	ColorFgHPcritical = ColorFgHPcritical - 1
	ColorFgMPok = ColorFgMPok - 1
	ColorFgMPpartial = ColorFgMPpartial - 1
	ColorFgMPcritical = ColorFgMPcritical - 1
	ColorFgStatusGood = ColorFgStatusGood - 1
	ColorFgStatusBad = ColorFgStatusBad - 1
	ColorFgStatusOther = ColorFgStatusOther - 1
	ColorFgTargetMode = ColorFgTargetMode - 1
	ColorFgTemporalWall = ColorFgTemporalWall - 1
}

func WindowsPalette() {
	ColorBgLOS = color(tcell.ColorSilver)
	ColorBgDark = color(tcell.ColorBlack)
	ColorBg = color(tcell.ColorBlack)
	ColorBgCloud = color(tcell.ColorSilver)
	ColorFgLOS = color(tcell.ColorBlack)
	ColorFgDark = color(tcell.ColorSilver)
	ColorFg = color(tcell.ColorSilver)
	ColorFgPlayer = color(tcell.ColorNavy)
	ColorFgMonster = color(tcell.ColorMaroon)
	ColorFgSleepingMonster = color(tcell.ColorTeal)
	ColorFgWanderingMonster = color(tcell.ColorPurple)
	ColorFgConfusedMonster = color(tcell.ColorGreen)
	ColorFgCollectable = color(tcell.ColorOlive)
	ColorFgStairs = color(tcell.ColorPurple)
	ColorFgGold = color(tcell.ColorOlive)
	ColorFgHPok = color(tcell.ColorGreen)
	ColorFgHPwounded = color(tcell.ColorOlive)
	ColorFgHPcritical = color(tcell.ColorMaroon)
	ColorFgMPok = color(tcell.ColorNavy)
	ColorFgMPpartial = color(tcell.ColorPurple)
	ColorFgMPcritical = color(tcell.ColorMaroon)
	ColorFgStatusGood = color(tcell.ColorNavy)
	ColorFgStatusBad = color(tcell.ColorMaroon)
	ColorFgStatusOther = color(tcell.ColorOlive)
	ColorFgTargetMode = color(tcell.ColorTeal)
	ColorFgTemporalWall = color(tcell.ColorTeal)
}

func (ui *termui) Init() error {
	screen, err := tcell.NewScreen()
	ui.Screen = screen
	if err != nil {
		return err
	}
	return ui.Screen.Init()
}

func (ui *termui) Close() {
	ui.Screen.Fini()
}

func (ui *termui) PostInit() {
	//termbox.SetOutputMode(termbox.Output256)
	//termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	FixColor()
	ui.Screen.SetStyle(tcell.StyleDefault)
	ui.Screen.EnableMouse()
}

func (ui *termui) Clear() {
	ui.Screen.Clear()
	w, h := ui.Screen.Size()
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(ColorFg)).Background(tcell.Color(ColorBg))
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			ui.Screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

func (ui *termui) Flush() {
	ui.Screen.Show()
}

func (ui *termui) Interrupt() {
	ui.Screen.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ui *termui) HideCursor() {
	ui.Screen.ShowCursor(-1, -1)
}

func (ui *termui) SetCursor(pos position) {
	ui.Screen.ShowCursor(pos.X, pos.Y)
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg color) {
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
	ui.Screen.SetContent(x, y, r, nil, st)
}

func (ui *termui) Reverse(c color) color {
	st := tcell.StyleDefault
	// TODO: true reverse of actual foreground
	st = st.Background(tcell.Color(c)).Foreground(tcell.Color(ColorFg))
	bg, _, _ := st.Reverse(true).Decompose()
	return color(bg)
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			if tev.Key() == tcell.KeyEsc {
				break loop
			}
			if tev.Rune() == ' ' {
				break loop
			}
		case *tcell.EventMouse:
			if tev.Buttons() == tcell.Button2 {
				break loop
			}
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			if tev.Rune() == 'Y' || tev.Rune() == 'y' {
				return true
			}
		}
		return false
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			return nil
		case *tcell.EventInterrupt:
			return errors.New("interrupted")
		case *tcell.EventMouse:
			switch tev.Buttons() {
			case tcell.Button1, tcell.Button2, tcell.Button3:
				return nil
			}
		}
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		again = false
		key := tev.Rune()
		switch tev.Key() {
		case tcell.KeyUp:
			key = 'k'
		case tcell.KeyRight:
			key = 'l'
		case tcell.KeyDown:
			key = 'j'
		case tcell.KeyLeft:
			key = 'h'
		case tcell.KeyCtrlW:
			ui.EnterWizard(g)
			return nil, true, false
		case tcell.KeyCtrlQ:
			if ui.Quit(g) {
				return nil, false, true
			}
			return nil, true, false
		case tcell.KeyCtrlP:
			key = 'm'
		}
		err, again, quit = ui.HandleCharacter(g, ev, key)
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.ButtonNone:
		case tcell.Button1:
			x, y := tev.Position()
			pos := position{X: x, Y: y}
			err, again = ui.GoToPos(g, ev, pos)
		case tcell.Button3:
			x, y := tev.Position()
			pos := position{X: x, Y: y}
			again = ui.ExaminePos(g, ev, pos)
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		if tev.Key() == tcell.KeyEsc {
			quit = true
			return n, quit
		}
		switch tev.Rune() {
		case 'u':
			n -= 12
		case 'd':
			n += 12
		case 'j':
			n++
		case 'k':
			n--
		case ' ':
			quit = true
		}
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.WheelUp:
			n -= 2
		case tcell.WheelDown:
			n += 2
		case tcell.Button2:
			quit = true
		}
	}
	return n, quit
}

func (ui *termui) TargetModeEvent(g *game, targ Targetter, pos position, data *examineData) bool {
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		key := tev.Rune()
		switch tev.Key() {
		case tcell.KeyUp:
			key = 'k'
		case tcell.KeyRight:
			key = 'l'
		case tcell.KeyDown:
			key = 'j'
		case tcell.KeyLeft:
			key = 'h'
		case tcell.KeyEsc:
			return true
		case tcell.KeyEnter:
			key = '.'
		}
		if ui.CursorCharAction(g, targ, key, pos, data) {
			return true
		}
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.Button1:
			if ui.CursorMouseLeft(g, targ, pos) {
				return true
			}
		case tcell.Button3:
			x, y := tev.Position()
			data.npos = position{X: x, Y: y}
		case tcell.Button2:
			return true
		}
	}
	return false
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			if tev.Key() == tcell.KeyEsc {
				return -1, false, errors.New("Ok, then.")
			}
			key := tev.Rune()
			if 97 <= key && int(key) < 97+l {
				return int(key - 97), false, nil
			}
			if key == '?' {
				return -1, true, nil
			}
			if key == ' ' {
				return -1, false, errors.New("Ok, then.")
			}
		case *tcell.EventMouse:
			switch tev.Buttons() {
			case tcell.Button1:
				_, y := tev.Position()
				if y > 0 && y <= l {
					return y - 1, false, nil
				}
			case tcell.Button3:
				return -1, true, nil
			case tcell.Button2:
				return -1, false, errors.New("Ok, then.")
			}
		}
	}
}
