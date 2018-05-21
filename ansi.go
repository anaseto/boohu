// +build ansi

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"unicode"
)

type ansiInput struct {
	r         rune
	interrupt bool
}

var ch chan ansiInput
var interrupt chan bool

func init() {
	ch = make(chan ansiInput, 100)
	interrupt = make(chan bool)
}

func main() {
	opt := flag.Bool("s", false, "Use true 16-uicolor solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	optCenteredCamera := flag.Bool("c", false, "centered camera")
	optMinimalUI := flag.Bool("m", false, "80x24 minimal UI")
	flag.Parse()
	if *opt {
		SolarizedPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optCenteredCamera {
		CenteredCamera = true
	}

	tui := &termui{}
	err := tui.Init()
	if *optMinimalUI {
		tui.minimal = true
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	ApplyDefaultKeyBindings()
	tui.PostInit()
	LinkColors()

	go func() {
		for {
			r, _, err := tui.bStdin.ReadRune()
			if err == nil {
				ch <- ansiInput{r: r}
			}
		}
	}()

	tui.DrawWelcome()
	g := &game{}
	load, err := g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Print("Error loading saved game… starting new game.")
	}
	load, err = g.LoadConfig()
	if load && err != nil {
		g.Print("Error loading config file.")
	} else if load {
		CustomKeys = true
	}

	g.ui = tui
	g.EventLoop()
}

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
	minimal    bool
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
	var prevfg, prevbg uicolor
	first := true
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		x, y := ui.GetPos(i)
		ui.MoveTo(x, y)
		pfg := true
		pbg := true
		if first {
			prevfg = cell.fg
			prevbg = cell.bg
			first = false
		} else {
			if prevfg == cell.fg {
				pfg = false
			} else {
				prevfg = cell.fg
			}
			if prevbg == cell.bg {
				pbg = false
			} else {
				prevbg = cell.bg
			}
		}
		if pfg {
			fmt.Fprintf(ui.bStdout, "\x1b[38;5;%dm", cell.fg)
		}
		if pbg {
			fmt.Fprintf(ui.bStdout, "\x1b[48;5;%dm", cell.bg)
		}
		ui.bStdout.WriteRune(cell.r)
		ui.backBuffer[i] = cell
	}
	ui.MoveTo(ui.cursor.X, ui.cursor.Y)
	fmt.Fprintf(ui.bStdout, "\x1b[0m")
	ui.bStdout.Flush()
}

func (ui *termui) Small() bool {
	return ui.minimal
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
	ui.cells[ui.GetIndex(x, y)] = AnsiCell{fg: fg, bg: bg, r: r}

}

func (ui *termui) Interrupt() {
	interrupt <- true
}

func (ui *termui) ReadChar() (in ansiInput) {
	select {
	case in = <-ch:
	case in.interrupt = <-interrupt:
	}
	return in
}

func (ui *termui) WaitForContinue(g *game, line int) {
loop:
	for {
		in := ui.ReadChar()
		switch in.r {
		case '\x1b', ' ':
			break loop
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		in := ui.ReadChar()
		switch in.r {
		case 'Y', 'y':
			return true
		default:
			return false
		}
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		in := ui.ReadChar()
		if in.interrupt {
			return errors.New("interrupted")
		}
		return nil
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	in := ui.ReadChar()
	switch in.r {
	case 'W':
		ui.EnterWizard(g)
		return nil, true, false
	case 'Q':
		if ui.Quit(g) {
			return nil, false, true
		}
		return nil, true, false
	}
	err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: in.r})
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	in := ui.ReadChar()
	switch in.r {
	case '\x1b', ' ':
		quit = true
	case 'u':
		n -= 12
	case 'd':
		n += 12
	case 'j', '2':
		n++
	case 'k', '8':
		n--
	}
	return n, quit
}

func (ui *termui) ReadRuneKey() rune {
	for {
		in := ui.ReadChar()
		if in.r == ' ' || in.r == '\x1b' {
			return 0
		}
		if unicode.IsPrint(in.r) {
			return in.r
		}
	}
}

func (ui *termui) MenuAction(n int) (m int, action configAction) {
	in := ui.ReadChar()
	switch in.r {
	case 'a':
		action = ChangeConfig
	case '\x1b', ' ':
		action = QuitConfig
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
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	in := ui.ReadChar()
	return ui.CursorKeyAction(g, targ, runeKeyAction{r: in.r}, data)
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		in := ui.ReadChar()
		switch {
		case in.r == '\x1b' || in.r == ' ':
			return -1, false, errors.New(DoNothing)
		case 97 <= in.r && int(in.r) < 97+l:
			return int(in.r - 97), false, nil
		case in.r == '?':
			return -1, true, nil
		case in.r == ' ':
			return -1, false, errors.New(DoNothing)
		}
	}
}
