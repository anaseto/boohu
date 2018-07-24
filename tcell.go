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
	small  bool
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

var SmallScreen = false

func (ui *termui) Flush() {
	ui.Screen.Show()
	w, h := ui.Screen.Size()
	if w <= UIWidth-8 || h <= UIHeight-2 {
		SmallScreen = true
	} else {
		SmallScreen = false
	}
}

func (ui *termui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
}

func (ui *termui) Small() bool {
	return gameConfig.Small || SmallScreen
}

func (ui *termui) Interrupt() {
	ui.Screen.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
	ui.Screen.SetContent(x, y, r, nil, st)
}

func (ui *termui) SetMapCell(x, y int, r rune, fg, bg uicolor) {
	ui.SetCell(x, y, r, fg, bg)
}

func (ui *termui) WaitForContinue(g *game, line int) {
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
			switch tev.Buttons() {
			case tcell.Button1:
				x, y := tev.Position()
				if line >= 0 {
					if y > line || x > DungeonWidth {
						break loop
					}
				} else {
					break loop
				}
			case tcell.Button2:
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
		r := tev.Rune()
		switch tev.Key() {
		case tcell.KeyLeft:
			// TODO: will not work if user changes keybindings
			r = '4'
		case tcell.KeyDown:
			r = '2'
		case tcell.KeyUp:
			r = '8'
		case tcell.KeyRight:
			r = '6'
		case tcell.KeyCtrlW:
			ui.EnterWizard(g)
			return nil, true, false
		case tcell.KeyCtrlQ:
			if ui.Quit(g) {
				return nil, false, true
			}
			return nil, true, false
		case tcell.KeyCtrlP:
			r = 'm'
		}
		err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: r})
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.ButtonNone:
		case tcell.Button1:
			x, y := tev.Position()
			pos := position{X: x, Y: y}
			if y == DungeonHeight {
				m, ok := ui.WhichButton(g, x)
				if !ok {
					again = true
					break
				}
				err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: m.Key(g)})
			} else if x >= DungeonWidth || y >= DungeonHeight {
				again = true
			} else {
				err, again, quit = ui.ExaminePos(g, ev, pos)
			}
		case tcell.Button3:
			err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: KeyMenu})
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
		case tcell.Button1:
			x, y := tev.Position()
			if x >= DungeonWidth {
				quit = true
				break
			}
			if y > UIHeight {
				break
			}
			n += y - (DungeonHeight+3)/2
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

func (ui *termui) KeyMenuAction(n int) (m int, action keyConfigAction) {
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		r := tev.Rune()
		switch tev.Key() {
		case tcell.KeyEsc:
			action = QuitKeyConfig
			return n, action
		case tcell.KeyDown:
			r = '2'
		case tcell.KeyUp:
			r = '8'
		}
		switch r {
		case 'a':
			action = ChangeKeys
		case 'u':
			n -= DungeonHeight / 2
		case 'd':
			n += DungeonHeight / 2
		case 'j', '2':
			n++
		case 'k', '8':
			n--
		case 'R':
			action = ResetKeys
		case ' ':
			action = QuitKeyConfig
		}
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.Button1:
			x, y := tev.Position()
			if x > DungeonWidth || y > DungeonHeight {
				action = QuitKeyConfig
			}
		case tcell.WheelUp:
			n -= 2
		case tcell.WheelDown:
			n += 2
		case tcell.Button2:
			action = QuitKeyConfig
		}
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	again = true
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		r := tev.Rune()
		switch tev.Key() {
		case tcell.KeyLeft:
			r = '4'
		case tcell.KeyDown:
			r = '2'
		case tcell.KeyUp:
			r = '8'
		case tcell.KeyRight:
			r = '6'
		case tcell.KeyEsc:
			g.Targeting = InvalidPos
			notarg = true
			return
		case tcell.KeyEnter:
			r = '.'
		}
		err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{r: r}, data)
	case *tcell.EventMouse:
		switch tev.Buttons() {
		case tcell.Button1:
			x, y := tev.Position()
			if y == DungeonHeight {
				m, ok := ui.WhichButton(g, x)
				if !ok {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
					break
				}
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: m.Key(g)}, data)
			} else if x >= DungeonWidth || y >= DungeonHeight {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
			} else {
				again, notarg = ui.CursorMouseLeft(g, targ, position{X: x, Y: y}, data)
			}
		case tcell.Button3:
			x, y := tev.Position()
			if y >= DungeonHeight || x >= DungeonWidth {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyMenu}, data)
			} else {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyDescription}, data)
			}
		case tcell.Button2:
			err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyExclude}, data)
		}
	}
	return err, again, quit, notarg
}

func (ui *termui) Select(g *game, l int) (index int, alternate bool, err error) {
	for {
		switch tev := ui.Screen.PollEvent().(type) {
		case *tcell.EventKey:
			if tev.Key() == tcell.KeyEsc {
				return -1, false, errors.New(DoNothing)
			}
			r := tev.Rune()
			if 97 <= r && int(r) < 97+l {
				return int(r - 97), false, nil
			}
			if r == '?' {
				return -1, true, nil
			}
			if r == ' ' {
				return -1, false, errors.New(DoNothing)
			}
		case *tcell.EventMouse:
			switch tev.Buttons() {
			case tcell.Button1:
				x, y := tev.Position()
				if y < 0 || y > l || x >= DungeonWidth {
					return -1, false, errors.New(DoNothing)
				}
				if y == 0 {
					return -1, true, nil
				}
				return y - 1, false, nil
			case tcell.Button3:
				return -1, true, nil
			case tcell.Button2:
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}
