// +build ansi

package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

var ch chan uiInput
var interrupt chan bool

func init() {
	ch = make(chan uiInput, 100)
	interrupt = make(chan bool)
}

type gameui struct {
	g       *game
	bStdin  *bufio.Reader
	bStdout *bufio.Writer
	cursor  position
	stty    string
	// below unused for this backend
	menuHover menu
	itemHover int
}

func (ui *gameui) Init() error {
	ui.bStdin = bufio.NewReader(os.Stdin)
	ui.bStdout = bufio.NewWriter(os.Stdout)
	fmt.Fprint(ui.bStdout, "\x1b[2J")
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
	ui.menuHover = -1
	go func() {
		for {
			r, _, err := ui.bStdin.ReadRune()
			if err == nil {
				ch <- uiInput{key: string(r)}
			}
		}
	}()

	return nil
}

func (ui *gameui) Close() {
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

func (ui *gameui) MoveTo(x, y int) {
	fmt.Fprintf(ui.bStdout, "\x1b[%d;%dH", y+1, x+1)
}

func (ui *gameui) Flush() {
	ui.DrawLogFrame()
	var prevfg, prevbg uicolor
	first := true
	var prevx, prevy int
	for _, cdraw := range ui.g.DrawLog[len(ui.g.DrawLog)-1].Draws {
		cell := cdraw.Cell
		x, y := cdraw.X, cdraw.Y
		pfg := true
		pbg := true
		pxy := true
		if first {
			prevfg = cell.Fg
			prevbg = cell.Bg
			prevx = x
			prevy = y
			first = false
		} else {
			if prevfg == cell.Fg {
				pfg = false
			} else {
				prevfg = cell.Fg
			}
			if prevbg == cell.Bg {
				pbg = false
			} else {
				prevbg = cell.Bg
			}
			if x == prevx+1 && y == prevy {
				pxy = false
			}
		}
		if pxy {
			ui.MoveTo(x, y)
		}
		if pfg {
			fmt.Fprintf(ui.bStdout, "\x1b[38;5;%dm", cell.Fg)
		}
		if pbg {
			fmt.Fprintf(ui.bStdout, "\x1b[48;5;%dm", cell.Bg)
		}
		ui.bStdout.WriteRune(cell.R)
	}
	ui.MoveTo(ui.cursor.X, ui.cursor.Y)
	fmt.Fprintf(ui.bStdout, "\x1b[0m")
	ui.bStdout.Flush()
}

func (ui *gameui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
	if gameConfig.Small {
		ui.Clear()
		ui.Flush()
		UIHeight = 24
		UIWidth = 80
	} else {
		UIHeight = 26
		if CenteredCamera {
			UIWidth = 80
		} else {
			UIWidth = 100
		}
	}
	ui.g.DrawBuffer = make([]UICell, UIWidth*UIHeight)
	ui.Clear()
}

func (ui *gameui) Small() bool {
	return gameConfig.Small
}

func (ui *gameui) Interrupt() {
	interrupt <- true
}

func (ui *gameui) PollEvent() (in uiInput) {
	select {
	case in = <-ch:
	case in.interrupt = <-interrupt:
	}
	return in
}
