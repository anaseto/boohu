// +build js

// TODO: adapt this from ansi.go for js

package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var Version string = "v0.4"

func main() {
	tui := &termui{}
	err := tui.Init()
	if err != nil {
		log.Fatal("boohu: %v\n", err)
	}
	defer tui.Close()

	tui.PostInit()

	tui.DrawWelcome()
	g := &game{}
	//load, err := g.Load()
	//if !load {
	//g.InitLevel()
	//} else if err != nil {
	//g.InitLevel()
	//g.Print("Error loading saved game… starting new game.")
	//}
	g.InitLevel()
	g.ui = tui
	g.EventLoop()
}

// io compatibility functions

func (g *game) DataDir() (string, error) {
	return "", nil
}

func (g *game) Save() error {
	return nil
}

func (g *game) RemoveSaveFile() error {
	return nil
}

func (g *game) Load() (bool, error) {
	return false, nil
}

func (g *game) WriteDump() error {
	return nil
}

// End of io compatibility functions

type UICell struct {
	fg uicolor
	bg uicolor
	r  rune
}

func (c uicolor) String() string {
	// TODO
	return "color"
}

type termui struct {
	cells  []UICell
	cursor position
}

const (
	UIWidth = 103
	//UIWidth = 10
	UIHeigth = 27
	//UIHeigth = 5
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

var ch chan string
var wants chan bool

func init() {
	ch = make(chan string)
	wants = make(chan bool)
}

func (ui *termui) Init() error {
	ui.cells = make([]UICell, UIWidth*UIHeigth)
	js.Global.Get("document").Call("addEventListener", "keypress", func(e *js.Object) {
		select {
		case <-wants:
			s := e.Get("key").String()
			ch <- s
		default:
		}
	})
	ui.ResetCells()
	js.Global.Get("document").Call("getElementById", "game").Set("innerHTML", "game screen")
	// TODO: init pre?
	return nil
}

func (ui *termui) Close() {
	// TODO
}

func (ui *termui) PostInit() {
	ui.HideCursor()
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Flush() {
	buf := &bytes.Buffer{}
	for i := 0; i < len(ui.cells); i++ {
		cell := ui.cells[i]
		// TODO: print newlines and html markup
		if i%UIWidth == 0 {
			fmt.Fprintf(buf, "\n")
		}
		if cell.r == ' ' {
			cell.r = ' '
		}
		//fmt.Fprintf(buf, "<span class=\"%s %s\">%c</span>", cell.fg, cell.bg, cell.r)
		fmt.Fprintf(buf, "%c", cell.r)
	}
	//ui.MoveTo(ui.cursor.X, ui.cursor.Y)
	//if ui.cursor.X >= 0 && ui.cursor.Y >= 0 {
	//fmt.Fprintf(buf, "\x1b[?25h")
	//} else {
	//fmt.Fprintf(buf, "\x1b[?25l")
	//}
	js.Global.Get("document").Call("getElementById", "game").Set("innerHTML", buf.String())
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
	ui.cells[ui.GetIndex(x, y)] = UICell{fg: fg, bg: bg, r: r}
}

func (ui *termui) ReadChar() rune {
	wants <- true
	s := <-ch
	bs := strings.NewReader(s)
	r, _, _ := bs.ReadRune()
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

func (ui *termui) TargetModeEvent(g *game, targ Targetter, pos position, data *examineData) bool {
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
