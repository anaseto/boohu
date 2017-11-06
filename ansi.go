// +build ansi

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type termui struct {
	bStdin  *bufio.Reader
	bStdout *bufio.Writer
}

func WindowsPalette() {
	ColorBgLOS = uicolor(termbox.ColorWhite)
	ColorBgDark = uicolor(termbox.ColorBlack)
	ColorBg = uicolor(termbox.ColorBlack)
	ColorBgCloud = uicolor(termbox.ColorWhite)
	ColorFgLOS = uicolor(termbox.ColorBlack)
	ColorFgDark = uicolor(termbox.ColorWhite)
	ColorFg = uicolor(termbox.ColorWhite)
	ColorFgPlayer = uicolor(termbox.ColorBlue)
	ColorFgMonster = uicolor(termbox.ColorRed)
	ColorFgSleepingMonster = uicolor(termbox.ColorCyan)
	ColorFgWanderingMonster = uicolor(termbox.ColorMagenta)
	ColorFgConfusedMonster = uicolor(termbox.ColorGreen)
	ColorFgCollectable = uicolor(termbox.ColorYellow)
	ColorFgStairs = uicolor(termbox.ColorMagenta)
	ColorFgGold = uicolor(termbox.ColorYellow)
	ColorFgHPok = uicolor(termbox.ColorGreen)
	ColorFgHPwounded = uicolor(termbox.ColorYellow)
	ColorFgHPcritical = uicolor(termbox.ColorRed)
	ColorFgMPok = uicolor(termbox.ColorBlue)
	ColorFgMPpartial = uicolor(termbox.ColorMagenta)
	ColorFgMPcritical = uicolor(termbox.ColorRed)
	ColorFgStatusGood = uicolor(termbox.ColorBlue)
	ColorFgStatusBad = uicolor(termbox.ColorRed)
	ColorFgStatusOther = uicolor(termbox.ColorYellow)
	ColorFgTargetMode = uicolor(termbox.ColorCyan)
	ColorFgTemporalWall = uicolor(termbox.ColorCyan)
}

func (ui *termui) Init() error {
	ui.bStdin = bufio.NewReader(os.Stdin)
	ui.bStdout = bufio.NewWriter(os.Stdout)
	// TODO: stty
	return nil
}

func (ui *termui) Close() {
	// TODO: stty
}

func (ui *termui) PostInit() {
	//SolarizedPalette()
	FixColor()
}

func (ui *termui) MoveTo(x, y int) {
	fmt.Fprintf(ui.bStdout, "\x1b[%d;%dH", y, x)
}

func (ui *termui) Clear() {
	// TODO: avoid complete clear
	fmt.Fprintf(ui.bStdout, "\x1b[2J")
	ui.MoveTo(1, 1)
}

func (ui *termui) Flush() {
	ui.bStdout.Flush()
}

func (ui *termui) HideCursor() {
	fmt.Fprintf(ui.bStdout, "\x1b[?25l")
}

func (ui *termui) SetCursor(pos position) {
	fmt.Fprintf(ui.bStdout, "\x1b[?25h")
	ui.MoveTo(pos.X, pos.Y)
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	//var fgAttr string
	//if fg <= 7 {
	//fgAttr = fmt.Sprintf("%d", fg)
	//} else {
	//fgAttr = fmt.Sprintf("1;%d", fg)
	//}
	//var bgAttr string
	//if bg <= 7 {
	//bgAttr = fmt.Sprintf("%d", 40+bg)
	//} else {
	//bgAttr = fmt.Sprintf("%d", 100+bg)
	//}

	ui.MoveTo(x, y)
	//fmt.Fprintf(ui.bStdout, "\x1b[%s;%sm", fgAttr, bgAttr)
	fmt.Fprintf(ui.bStdout, "\x1b[38;5;%dm", fg)
	fmt.Fprintf(ui.bStdout, "\x1b[48;5;%dm", bg)
	ui.bStdout.WriteRune(r)
	fmt.Fprintf(ui.bStdout, "\x1b[0m")
}

func (ui *termui) ReadChar() rune {
	cmd := exec.Command("stty", "raw", "-echo")
	cmd.Stdin = os.Stdin
	cmd.Run()
	r, _, _ := ui.bStdin.ReadRune()
	cmd = exec.Command("stty", "sane")
	cmd.Stdin = os.Stdin
	cmd.Run()
	return r
}

func (ui *termui) ExploreStep(g *game) bool {
	time.Sleep(10 * time.Millisecond)
	ui.DrawDungeonView(g, false)
	return false
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		r := ui.ReadChar()
		switch r {
		case '\xb1', ' ':
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
	err, again, quit = ui.HandleCharacter(g, ev, r)
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	r := ui.ReadChar()
	switch r {
	case '\xb1', ' ':
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

func (ui *termui) TargetModeEvent(g *game, targ Targetter, pos position, data *examineData) bool {
	r := ui.ReadChar()
	if r == '\xb1' {
		return true
	}
	return ui.CursorCharAction(g, targ, r, pos, data)
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		r := ui.ReadChar()
		switch {
		case r == '\xb1' || r == ' ':
			return -1, false, errors.New("Ok, then.")
		case 97 <= r && int(r) < 97+l:
			return int(r - 97), false, nil
		case r == '?':
			return -1, true, nil
		case r == ' ':
			return -1, false, errors.New("Ok, then.")
		}
	}
}
