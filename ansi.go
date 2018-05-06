// +build ansi

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type AnsiCell struct {
	fg uicolor
	bg uicolor
	r  rune
}

type termui struct {
	bStdin     *bufio.Reader
	bStdout    *bufio.Writer
	cells      []AnsiCell
	backBuffer []AnsiCell
	cursor     position
	stty       string
}

const (
	UIWidth  = 103
	UIHeight = 27
)

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

func (ui *termui) Init() error {
	ui.bStdin = bufio.NewReader(os.Stdin)
	ui.bStdout = bufio.NewWriter(os.Stdout)
	ui.cells = make([]AnsiCell, UIWidth*UIHeight)
	ui.ResetCells()
	ui.backBuffer = make([]AnsiCell, UIWidth*UIHeight)
	fmt.Fprint(ui.bStdout, "\x1b[2J")
	return nil
}

func (ui *termui) Close() {
	fmt.Fprint(ui.bStdout, "\x1b[2J")
	fmt.Fprintf(ui.bStdout, "\x1b[?25h")
	ui.bStdout.Flush()
	cmd := exec.Command("stty", ui.stty)
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		cmd = exec.Command("stty", "sane")
		cmd.Stdin = os.Stdin
		cmd.Run()
	}
}

func (ui *termui) PostInit() {
	ui.HideCursor()
	fmt.Fprintf(ui.bStdout, "\x1b[?25l")
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	save, err := cmd.Output()
	if err != nil {
		save = []byte("sane")
	}
	ui.stty = string(save)
	cmd = exec.Command("stty", "raw", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func (ui *termui) MoveTo(x, y int) {
	fmt.Fprintf(ui.bStdout, "\x1b[%d;%dH", y+1, x+1)
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Flush() {
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		x, y := ui.GetPos(i)
		ui.MoveTo(x, y)
		fmt.Fprintf(ui.bStdout, "\x1b[38;5;%dm", cell.fg)
		fmt.Fprintf(ui.bStdout, "\x1b[48;5;%dm", cell.bg)
		ui.bStdout.WriteRune(cell.r)
		fmt.Fprintf(ui.bStdout, "\x1b[0m")
		ui.backBuffer[i] = cell
	}
	ui.MoveTo(ui.cursor.X, ui.cursor.Y)
	ui.bStdout.Flush()
}

func (ui *termui) HideCursor() {
	ui.cursor = position{-1, -1}
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	i := ui.GetIndex(x, y)
	if i >= len(ui.cells) {
		return
	}
	ui.cells[ui.GetIndex(x, y)] = AnsiCell{fg: fg, bg: bg, r: r}

}

func (ui *termui) ReadChar() rune {
	r, _, _ := ui.bStdin.ReadRune()
	return r
}

func (ui *termui) ExploreStep(g *game) bool {
	time.Sleep(10 * time.Millisecond)
	ui.DrawDungeonView(g, NormalMode)
	return false
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		r := ui.ReadChar()
		switch r {
		case '\x1b', ' ':
			break loop
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		r := ui.ReadChar()
		switch r {
		case 'Y', 'y':
			return true
		default:
			return false
		}
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		ui.ReadChar()
		return nil
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	r := ui.ReadChar()
	switch r {
	case 'W':
		ui.EnterWizard(g)
		return nil, true, false
	case 'Q':
		if ui.Quit(g) {
			return nil, false, true
		}
		return nil, true, false
	}
	err, again, quit = ui.HandleCharacter(g, ev, r)
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	r := ui.ReadChar()
	switch r {
	case '\x1b', ' ':
		quit = true
	case 'u':
		n -= 12
	case 'd':
		n += 12
	case 'j':
		n++
	case 'k':
		n--
	}
	return n, quit
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, pos position, data *examineData) bool {
	r := ui.ReadChar()
	if r == '\x1b' || r == ' ' {
		return true
	}
	return ui.CursorCharAction(g, targ, r, pos, data)
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		r := ui.ReadChar()
		switch {
		case r == '\x1b' || r == ' ':
			return -1, false, errors.New(DoNothing)
		case 97 <= r && int(r) < 97+l:
			return int(r - 97), false, nil
		case r == '?':
			return -1, true, nil
		case r == ' ':
			return -1, false, errors.New(DoNothing)
		}
	}
}
