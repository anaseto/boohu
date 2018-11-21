package main

import (
	"errors"
	"strings"
	"unicode"
)

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *termui) WaitForContinue(g *game, line int) {
loop:
	for {
		in := ui.PollEvent()
		switch in.key {
		case "Escape", " ":
			break loop
		}
		if in.mouse && in.button == -1 {
			continue
		}
		if in.mouse && line >= 0 {
			if in.mouseY > line || in.mouseX > DungeonWidth {
				break loop
			}
		} else if in.mouse {
			break loop
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "Y", "y":
			return true
		default:
			return false
		}
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		e := ui.PollEvent()
		if e.interrupt {
			return errors.New("interrupted")
		}
		if e.key != "" || (e.mouse && e.button != -1) {
			return nil
		}
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "":
		if in.mouse {
			pos := position{X: in.mouseX, Y: in.mouseY}
			switch in.button {
			case -1:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(g, in.mouseX)
					omh := ui.menuHover
					if ok {
						ui.menuHover = m
					} else {
						ui.menuHover = -1
					}
					if ui.menuHover != omh {
						ui.DrawMenus(g)
						ui.Flush()
					}
					break
				}
				ui.menuHover = -1
				if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
					break
				}
				fallthrough
			case 0:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(g, in.mouseX)
					if !ok {
						again = true
						break
					}
					err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: m.Key(g)})
					if err != nil {
						again = true
					}
					return err, again, quit
				} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
				} else {
					err, again, quit = ui.ExaminePos(g, ev, pos)
				}
			case 2:
				err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: KeyMenu})
				if err != nil {
					again = true
				}
				return err, again, quit
			}
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r == 0 {
			again = true
		} else {
			err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: r})
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	in := ui.PollEvent()
	switch in.key {
	case "Escape", "\x1b", " ":
		quit = true
	case "u":
		n -= 12
	case "d":
		n += 12
	case "j", "2":
		n++
	case "k", "8":
		n--
	case "":
		if in.mouse {
			switch in.button {
			case 0:
				y := in.mouseY
				x := in.mouseX
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

func (ui *termui) GetIndex(x, y int) int {
	return y*UIWidth + x
}

func (ui *termui) Select(g *game, l int) (index int, alternate bool, err error) {
	for {
		in := ui.PollEvent()
		r := ui.ReadKey(in.key)
		switch {
		case in.key == "\x1b" || in.key == "Escape" || in.key == " ":
			return -1, false, errors.New(DoNothing)
		case in.key == "?":
			return -1, true, nil
		case 97 <= r && int(r) < 97+l:
			return int(r - 97), false, nil
		case in.key == "" && in.mouse:
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case -1:
				oih := ui.itemHover
				if y <= 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					if oih != -1 {
						ui.ColorLine(oih, ColorFg)
						ui.Flush()
					}
					break
				}
				if y == oih {
					break
				}
				ui.itemHover = y
				ui.ColorLine(y, ColorYellow)
				if oih != -1 {
					ui.ColorLine(oih, ColorFg)
				}
				ui.Flush()
			case 0:
				if y < 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					return -1, false, errors.New(DoNothing)
				}
				if y == 0 {
					ui.itemHover = -1
					return -1, true, nil
				}
				ui.itemHover = -1
				return y - 1, false, nil
			case 2:
				ui.itemHover = -1
				return -1, true, nil
			case 1:
				ui.itemHover = -1
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}

func (ui *termui) KeyMenuAction(n int) (m int, action keyConfigAction) {
	in := ui.PollEvent()
	switch in.key {
	case "a":
		action = ChangeKeys
	case "\x1b", "Escape", " ":
		action = QuitKeyConfig
	case "u":
		n -= DungeonHeight / 2
	case "d":
		n += DungeonHeight / 2
	case "j", "2", "ArrowDown":
		n++
	case "k", "8", "ArrowUp":
		n--
	case "R":
		action = ResetKeys
	case "":
		if in.mouse {
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case 0:
				if x > DungeonWidth || y > DungeonHeight {
					action = QuitKeyConfig
				}
			case 1:
				action = QuitKeyConfig
			}
		}
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "\x1b", "Escape", " ":
		g.Targeting = InvalidPos
		notarg = true
	case "":
		if !in.mouse {
			return
		}
		switch in.button {
		case -1:
			if in.mouseY == DungeonHeight {
				m, ok := ui.WhichButton(g, in.mouseX)
				omh := ui.menuHover
				if ok {
					ui.menuHover = m
				} else {
					ui.menuHover = -1
				}
				if ui.menuHover != omh {
					ui.DrawMenus(g)
					ui.Flush()
				}
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			ui.menuHover = -1
			if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			mpos := position{in.mouseX, in.mouseY}
			if g.Targeting == mpos {
				break
			}
			g.Targeting = InvalidPos
			fallthrough
		case 0:
			if in.mouseY == DungeonHeight {
				m, ok := ui.WhichButton(g, in.mouseX)
				if !ok {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
					break
				}
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: m.Key(g)}, data)
			} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
			} else {
				again, notarg = ui.CursorMouseLeft(g, targ, position{X: in.mouseX, Y: in.mouseY}, data)
			}
		case 2:
			if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyMenu}, data)
			} else {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyDescription}, data)
			}
		case 1:
			err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyExclude}, data)
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r != 0 {
			return ui.CursorKeyAction(g, targ, runeKeyAction{r: r}, data)
		}
		again = true
		notarg = true
	}
	return
}

func (ui *termui) ReadRuneKey() rune {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "\x1b", "Escape", " ":
			return 0
		case "Enter":
			return '.'
		}
		r := ui.ReadKey(in.key)
		if unicode.IsPrint(r) {
			return r
		}
	}
}

func (ui *termui) ReadKey(s string) (r rune) {
	bs := strings.NewReader(s)
	r, _, _ = bs.ReadRune()
	return r
}
