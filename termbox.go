// +build !tcell,!ansi,!js

package main

import (
	"errors"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

type termui struct {
	cursor position
}

func (ui *termui) Init() error {
	return termbox.Init()
}

func (ui *termui) Close() {
	termbox.Close()
}

func (ui *termui) PostInit() {
	FixColor()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.HideCursor()
	ui.HideCursor()
}

func (ui *termui) Clear() {
	termbox.Clear(termbox.Attribute(ColorFg), termbox.Attribute(ColorBg))
}

func (ui *termui) Flush() {
	termbox.Flush()
}

func (ui *termui) Interrupt() {
	termbox.Interrupt()
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	termbox.SetCell(x, y, r, termbox.Attribute(fg), termbox.Attribute(bg))
}

func (ui *termui) WaitForContinue(g *game, line int) {
loop:
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					break loop
				}
			}
			if tev.Ch == ' ' {
				break loop
			}
		case termbox.EventMouse:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.MouseMiddle:
					break loop
				case termbox.MouseLeft:
					if line >= 0 {
						if tev.MouseY > line || tev.MouseX > DungeonWidth {
							break loop
						}

					} else {
						break loop
					}
				}
			}
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 'Y' || tev.Ch == 'y' {
				return true
			}
		}
		return false
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			return nil
		case termbox.EventInterrupt:
			return errors.New("interrupted")
		case termbox.EventMouse:
			if tev.Ch == 0 && tev.Key == termbox.MouseLeft ||
				tev.Key == termbox.MouseMiddle || tev.Key == termbox.MouseRight {
				return nil
			}
		}
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		again = false
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyArrowLeft:
				// TODO: will not work if user changes keybindings
				tev.Ch = '4'
			case termbox.KeyArrowDown:
				tev.Ch = '2'
			case termbox.KeyArrowUp:
				tev.Ch = '8'
			case termbox.KeyArrowRight:
				tev.Ch = '6'
			case termbox.KeyCtrlW:
				ui.EnterWizard(g)
				return nil, true, false
			case termbox.KeyCtrlQ:
				if ui.Quit(g) {
					return nil, false, true
				}
				return nil, true, false
			case termbox.KeyCtrlP:
				tev.Ch = 'm'
			}
		}
		err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: tev.Ch})
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseLeft:
				pos := position{X: tev.MouseX, Y: tev.MouseY}
				if pos.X > DungeonWidth && pos.Y == 0 {
					err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: KeyMenu})
				} else if pos.X > DungeonWidth || pos.Y > DungeonHeight {
					again = true
				} else {
					err, again, quit = ui.ExaminePos(g, ev, pos)
				}
			case termbox.MouseRight:
				err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: KeyMenu})
			}
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyEsc, termbox.KeySpace:
				quit = true
				return n, quit
			case termbox.KeyArrowDown:
				tev.Ch = '2'
			case termbox.KeyArrowUp:
				tev.Ch = '8'
			}
		}
		switch tev.Ch {
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
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseMiddle:
				quit = true
			case termbox.MouseWheelUp:
				n -= 2
			case termbox.MouseWheelDown:
				n += 2
			case termbox.MouseLeft:
				y := tev.MouseY
				x := tev.MouseX
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
	}
	return n, quit
}

func (ui *termui) ReadRuneKey() rune {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == ' ' {
				return 0
			}
			if unicode.IsPrint(tev.Ch) {
				return tev.Ch
			}
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					return 0
				}
			}
		}
	}
}

func (ui *termui) MenuAction(n int) (m int, action configAction) {
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyEsc, termbox.KeySpace:
				action = QuitConfig
				return n, action
			case termbox.KeyArrowDown:
				tev.Ch = '2'
			case termbox.KeyArrowUp:
				tev.Ch = '8'
			}
		}
		switch tev.Ch {
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
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseMiddle:
				action = QuitConfig
			case termbox.MouseWheelUp:
				n -= 2
			case termbox.MouseWheelDown:
				n += 2
			}
		}
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	again = true
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyArrowLeft:
				tev.Ch = '4'
			case termbox.KeyArrowDown:
				tev.Ch = '2'
			case termbox.KeyArrowUp:
				tev.Ch = '8'
			case termbox.KeyArrowRight:
				tev.Ch = '6'
			case termbox.KeyEsc, termbox.KeySpace:
				g.Targeting = InvalidPos
				notarg = true
				return
			case termbox.KeyEnter:
				tev.Ch = '.'
			}
		}
		err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{r: tev.Ch}, data)
	case termbox.EventMouse:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.MouseLeft:
				if tev.MouseX > DungeonWidth && tev.MouseY == 0 {
					err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyMenu}, data)
				} else if tev.MouseX > DungeonWidth || tev.MouseY > DungeonHeight {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
				} else {
					again, notarg = ui.CursorMouseLeft(g, targ, position{X: tev.MouseX, Y: tev.MouseY}, data)
				}
			case termbox.MouseRight:
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyMenu}, data)
			case termbox.MouseMiddle:
				g.Targeting = InvalidPos
				notarg = true
			}
		}
	}
	return err, again, quit, notarg
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					return -1, false, errors.New(DoNothing)
				}
			}
			if 97 <= tev.Ch && int(tev.Ch) < 97+l {
				return int(tev.Ch - 97), false, nil
			}
			if tev.Ch == '?' {
				return -1, true, nil
			}
			if tev.Ch == ' ' {
				return -1, false, errors.New(DoNothing)
			}
		case termbox.EventMouse:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.MouseLeft:
					y := tev.MouseY
					x := tev.MouseX
					if y < 0 || y > l || x >= DungeonWidth {
						return -1, false, errors.New(DoNothing)
					}
					if y == 0 {
						return -1, true, nil
					}
					return y - 1, false, nil
				case termbox.MouseRight:
					return -1, true, nil
				case termbox.MouseMiddle:
					return -1, false, errors.New(DoNothing)
				}
			}
		}
	}
}
