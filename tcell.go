// +build tcell

package main

import (
	"errors"
	"runtime"
	"unicode"

	"github.com/gdamore/tcell"
)

type termui struct {
	tcell.Screen
	cursor position
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
	ui.Screen.SetStyle(tcell.StyleDefault)
	if runtime.GOOS != "openbsd" {
		ui.Screen.EnableMouse()
	}
	ui.Screen.HideCursor()
	ui.HideCursor()
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
	ui.cursor = position{-1, -1}
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
	ui.Screen.SetContent(x, y, r, nil, st)
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
		case tcell.KeyLeft:
			// TODO: will not work if user changes keybindings
			key = '4'
		case tcell.KeyDown:
			key = '2'
		case tcell.KeyUp:
			key = '8'
		case tcell.KeyRight:
			key = '6'
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
		r := tev.Rune()
		switch tev.Key() {
		case tcell.KeyEsc:
			quit = true
			return n, quit
		case tcell.KeyDown:
			r = '2'
		case tcell.KeyUp:
			r = '8'
		}
		switch r {
		case 'u':
			n -= 12
		case 'd':
			n += 12
		case 'j', '2':
			n++
		case 'k', '8':
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

func (ui *termui) ReadRuneKey() rune {
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			r := tev.Rune()
			if r == ' ' {
				return 0
			}
			if unicode.IsPrint(r) {
				return r
			}
			switch tev.Key() {
			case tcell.KeyEsc:
				return 0
			}
		}
	}
}

func (ui *termui) MenuAction(n int) (m int, action configAction) {
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		r := tev.Rune()
		switch tev.Key() {
		case tcell.KeyEsc:
			action = QuitConfig
			return n, action
		case tcell.KeyDown:
			r = '2'
		case tcell.KeyUp:
			r = '8'
		}
		switch r {
		case 'a':
			action = ChangeConfig
		case 'u':
			n -= DungeonHeight / 2
		case 'd':
			n += DungeonHeight / 2
		case 'j', '2':
			n++
		case 'k', '8':
			n--
		case 'R':
			action = ResetConfig
		case ' ':
			action = QuitConfig
		}
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.WheelUp:
			n -= 2
		case tcell.WheelDown:
			n += 2
		case tcell.Button2:
			action = QuitConfig
		}
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) bool {
	pos := data.npos
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		key := tev.Rune()
		switch tev.Key() {
		case tcell.KeyLeft:
			key = '4'
		case tcell.KeyDown:
			key = '2'
		case tcell.KeyUp:
			key = '8'
		case tcell.KeyRight:
			key = '6'
		case tcell.KeyEsc:
			return true
		case tcell.KeyEnter:
			key = '.'
		}
		if ui.CursorCharAction(g, targ, key, data) {
			return true
		}
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.Button1:
			if ui.CursorMouseLeft(g, targ, pos, data) {
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
				return -1, false, errors.New(DoNothing)
			}
			key := tev.Rune()
			if 97 <= key && int(key) < 97+l {
				return int(key - 97), false, nil
			}
			if key == '?' {
				return -1, true, nil
			}
			if key == ' ' {
				return -1, false, errors.New(DoNothing)
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
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}
