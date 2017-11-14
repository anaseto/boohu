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

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    *js.Object
}

const (
	UIWidth = 103
	//UIWidth = 10
	UIHeight = 27
	//UIHeight = 5
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
	ui.cells = make([]UICell, UIWidth*UIHeight)
	js.Global.Get("document").Call("addEventListener", "keypress", func(e *js.Object) {
		select {
		case <-wants:
			s := e.Get("key").String()
			ch <- s
		default:
		}
	})
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
	//ui.display = js.Global.Get("ROT").Call("Display", map[string]int{"width": UIWidth, "height": UIHeight})
	//ui.display = js.Global.Get("ROT").Call("Display")
	//container := ui.display.Call("getContainer")
	//js.Global.Get("document").Get("body").Call("appendChild", container)
	//ui.display.Call("draw", 5, 4, "@")
	// TODO: init pre?
	return nil
}

func (ui *termui) Close() {
	// TODO
}

func (ui *termui) PostInit() {
	SolarizedPalette()
	ui.HideCursor()
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Flush() {
	buf := &bytes.Buffer{}
	canvas := js.Global.Get("document").Call("getElementById", "gamecanvas")
	ctx := canvas.Call("getContext", "2d")
	ctx.Set("font", "12px monospace")
	mesure := ctx.Call("measureText", "W")
	width := mesure.Get("width").Int() + 1
	//height := mesure.Get("height")
	//ctx.Set("fillStyle", "#002b36")
	//ctx.Call("fillRect", 0, 0, UIWidth*10, UIHeight*10)
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		// TODO: print newlines and html markup
		if i%UIWidth == 0 {
			fmt.Fprintf(buf, "\n")
		}
		if cell.r == ' ' {
			cell.r = ' '
		}
		x, y := ui.GetPos(i)
		ctx.Set("fillStyle", cell.bg.String())
		ctx.Call("fillRect", width*x, 17*y, width, 17)
		ctx.Set("fillStyle", cell.fg.String())
		ctx.Call("fillText", string(cell.r), width*x, 17*y+12)
		//ui.display.Call("draw", x, y, string(cell.r))
		ui.backBuffer[i] = cell
	}
	//ui.MoveTo(ui.cursor.X, ui.cursor.Y)
	//if ui.cursor.X >= 0 && ui.cursor.Y >= 0 {
	//fmt.Fprintf(buf, "\x1b[?25h")
	//} else {
	//fmt.Fprintf(buf, "\x1b[?25l")
	//}
	//ctx.Call("fillText", buf.String(), 10, 10)
	//js.Global.Get("document").Call("getElementById", "game").Set("innerHTML", buf.String())
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
	case 'S', 'Q', '#':
		err = errors.New("Command not available (still) for the js version.")
		return nil, true, false
	case 'W':
		ui.EnterWizard(g)
		return nil, true, false
		//case 'Q':
		//if ui.Quit(g) {
		//return nil, false, true
		//}
		//return nil, true, false
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
