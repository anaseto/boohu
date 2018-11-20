// +build js tk

package main

import (
	"errors"
	"fmt"
	"image/color"
	"strings"
	"unicode"
	"unicode/utf8"
)

type UICell struct {
	fg    uicolor
	bg    uicolor
	r     rune
	inMap bool
}

func (ui *termui) GetIndex(x, y int) int {
	return y*UIWidth + x
}

func (ui *termui) GetPos(i int) (int, int) {
	return i - (i/UIWidth)*UIWidth, i / UIWidth
}

func (ui *termui) ResetCells() {
	for i := 0; i < len(ui.cells); i++ {
		ui.cells[i].r = ' '
		ui.cells[i].bg = ColorBg
	}
}

func (ui *termui) ApplyToggleTiles() {
	gameConfig.Tiles = !gameConfig.Tiles
	for c, _ := range ui.cache {
		if c.inMap {
			delete(ui.cache, c)
		}
	}
	for i := 0; i < len(ui.backBuffer); i++ {
		ui.backBuffer[i] = UICell{}
	}
}

func (c uicolor) String() string {
	color := "#002b36"
	switch c {
	case 0:
		color = "#073642"
	case 1:
		color = "#dc322f"
	case 2:
		color = "#859900"
	case 3:
		color = "#b58900"
	case 4:
		color = "#268bd2"
	case 5:
		color = "#d33682"
	case 6:
		color = "#2aa198"
	case 7:
		color = "#eee8d5"
	case 8:
		color = "#002b36"
	case 9:
		color = "#cb4b16"
	case 10:
		color = "#586e75"
	case 11:
		color = "#657b83"
	case 12:
		color = "#839496"
	case 13:
		color = "#6c71c4"
	case 14:
		color = "#93a1a1"
	case 15:
		color = "#fdf6e3"
	}
	return color
}

func (c uicolor) Color() color.Color {
	cl := color.RGBA{}
	opaque := uint8(255)
	switch c {
	case 0:
		cl = color.RGBA{7, 54, 66, opaque}
	case 1:
		cl = color.RGBA{220, 50, 47, opaque}
	case 2:
		cl = color.RGBA{133, 153, 0, opaque}
	case 3:
		cl = color.RGBA{181, 137, 0, opaque}
	case 4:
		cl = color.RGBA{38, 139, 210, opaque}
	case 5:
		cl = color.RGBA{211, 54, 130, opaque}
	case 6:
		cl = color.RGBA{42, 161, 152, opaque}
	case 7:
		cl = color.RGBA{238, 232, 213, opaque}
	case 8:
		cl = color.RGBA{0, 43, 54, opaque}
	case 9:
		cl = color.RGBA{203, 75, 22, opaque}
	case 10:
		cl = color.RGBA{88, 110, 117, opaque}
	case 11:
		cl = color.RGBA{101, 123, 131, opaque}
	case 12:
		cl = color.RGBA{131, 148, 150, opaque}
	case 13:
		cl = color.RGBA{108, 113, 196, opaque}
	case 14:
		cl = color.RGBA{147, 161, 161, opaque}
	case 15:
		cl = color.RGBA{253, 246, 227, opaque}
	}
	return cl
}

func (ui *termui) DrawMenus(g *game) {
	line := DungeonHeight
	for i, cols := range MenuCols[0 : len(MenuCols)-1] {
		if cols[0] >= 0 {
			if menu(i) == ui.menuHover {
				ui.DrawColoredText(menu(i).String(), cols[0], line, ColorBlue)
			} else {
				ui.DrawColoredText(menu(i).String(), cols[0], line, ColorViolet)
			}
		}
	}
	interactMenu := ui.UpdateInteractButton(g)
	if interactMenu == "" {
		return
	}
	i := len(MenuCols) - 1
	cols := MenuCols[i]
	if menu(i) == ui.menuHover {
		ui.DrawColoredText(interactMenu, cols[0], line, ColorBlue)
	} else {
		ui.DrawColoredText(interactMenu, cols[0], line, ColorViolet)
	}
}

var TileImgs map[string][]byte

var MapNames = map[rune]string{
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'☻':  "dreaming",
	'♫':  "footsteps",
	'#':  "wall",
	'@':  "player",
	'§':  "fog",
	'♣':  "simella",
	'+':  "door",
	'.':  "ground",
	'"':  "foliage",
	'•':  "tick",
	'●':  "rock",
	'×':  "times",
	',':  "comma",
	'}':  "rbrace",
	'%':  "percent",
	':':  "colon",
	'\\': "backslash",
	'~':  "tilde",
	'☼':  "sun",
	'*':  "asterisc",
	'—':  "hbar",
	'/':  "slash",
	'|':  "vbar",
	'∞':  "kill",
	' ':  "space",
	'[':  "lbracket",
	']':  "rbracket",
	')':  "rparen",
	'(':  "lparen",
	'>':  "stairs",
	'!':  "potion",
	';':  "semicolon",
	'_':  "stone",
}

var LetterNames = map[rune]string{
	'(':  "lparen",
	')':  "rparen",
	'@':  "player",
	'{':  "lbrace",
	'}':  "rbrace",
	'[':  "lbracket",
	']':  "rbracket",
	'♪':  "music1",
	'♫':  "music2",
	'•':  "tick",
	'♣':  "simella",
	' ':  "space",
	'!':  "exclamation",
	'?':  "interrogation",
	',':  "comma",
	':':  "colon",
	';':  "semicolon",
	'\'': "quote",
	'—':  "longhyphen",
	'-':  "hyphen",
	'|':  "pipe",
	'/':  "slash",
	'\\': "backslash",
	'%':  "percent",
	'┐':  "boxne",
	'┤':  "boxe",
	'│':  "vbar",
	'┘':  "boxse",
	'─':  "hbar",
	'►':  "arrow",
	'×':  "times",
	'.':  "dot",
	'#':  "hash",
	'"':  "quotes",
	'+':  "plus",
	'“':  "lquotes",
	'”':  "rquotes",
	'=':  "equal",
	'>':  "gt",
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'§':  "fog",
	'●':  "rock",
	'~':  "tilde",
	'☼':  "sun",
	'*':  "asterisc",
	'∞':  "kill",
	'☻':  "dreaming",
	'…':  "dots",
	'_':  "stone",
}

func (ui *termui) Interrupt() {
	interrupt <- true
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Small() bool {
	return gameConfig.Small
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	i := ui.GetIndex(x, y)
	if i >= len(ui.cells) {
		return
	}
	ui.cells[i] = UICell{fg: fg, bg: bg, r: r}
}

func (ui *termui) SetMapCell(x, y int, r rune, fg, bg uicolor) {
	i := ui.GetIndex(x, y)
	if i >= len(ui.cells) {
		return
	}
	ui.cells[i] = UICell{fg: fg, bg: bg, r: r, inMap: true}
}

func (ui *termui) ReadKey(s string) (r rune) {
	bs := strings.NewReader(s)
	r, _, _ = bs.ReadRune()
	return r
}

func (ui *termui) PollEvent() (in uiInput) {
	select {
	case in = <-ch:
	case in.interrupt = <-interrupt:
	}
	return in
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
		switch in.key {
		case "Enter":
			in.key = "."
		case "ArrowLeft":
			in.key = "4"
		case "ArrowRight":
			in.key = "6"
		case "ArrowUp":
			in.key = "8"
		case "ArrowDown":
			in.key = "2"
		case "Home":
			in.key = "7"
		case "End":
			in.key = "1"
		case "PageUp":
			in.key = "9"
		case "PageDown":
			in.key = "3"
		case "Numpad5", "Delete":
			in.key = "5"
		}
		if utf8.RuneCountInString(in.key) > 1 {
			err = fmt.Errorf("Invalid key: “%s”.", in.key)
		} else {
			err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: ui.ReadKey(in.key)})
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
		return
	case "Enter":
		in.key = "."
	case "":
		if in.mouse {
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
		}
		return err, again, quit, notarg
	case "ArrowLeft":
		in.key = "4"
	case "ArrowRight":
		in.key = "6"
	case "ArrowUp":
		in.key = "8"
	case "ArrowDown":
		in.key = "2"
	case "Home":
		in.key = "7"
	case "End":
		in.key = "1"
	case "PageUp":
		in.key = "9"
	case "PageDown":
		in.key = "3"
	case "Numpad5", "Delete":
		in.key = "5"
	}
	if utf8.RuneCountInString(in.key) > 1 {
		g.Printf("Invalid key: “%s”.", in.key)
		notarg = true
		return
	}
	return ui.CursorKeyAction(g, targ, runeKeyAction{r: ui.ReadKey(in.key)}, data)
}

func (ui *termui) ColorLine(y int, fg uicolor) {
	for x := 0; x < DungeonWidth; x++ {
		i := ui.GetIndex(x, y)
		c := ui.cells[i]
		ui.SetCell(x, y, c.r, fg, c.bg)
	}
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
