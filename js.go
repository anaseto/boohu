// +build js

package main

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
)

var Version string = "v0.7"

func main() {
	tui := &termui{}
	err := tui.Init()
	if err != nil {
		log.Fatalf("boohu: %v\n", err)
	}
	defer tui.Close()

	tui.PostInit()
	LinkColors()

	tui.DrawWelcome()
	g := &game{}
	//load, err := g.Load()
	//if !load {
	g.InitLevel()
	//} else if err != nil {
	//g.InitLevel()
	//g.Printf("Error loading saved game… starting new game. (%v)", err)
	//}
	g.ui = tui
	g.EventLoop()
	tui.Clear()
	tui.DrawText("Refresh the page to start again", 0, 0)
	tui.Flush()
}

// io compatibility functions

func (g *game) DataDir() (string, error) {
	return "", nil
}

func (g *game) Save() error {
	return nil
	//save, err := g.GameSave()
	//if err != nil {
	//return err
	//}
	//storage := js.Global.Get("localStorage")
	//if !storage.Bool() {
	//return errors.New("localStorage not found")
	//}
	//s := base64.StdEncoding.EncodeToString(save)
	//storage.Call("setItem", "boohusave", s)
	//return nil
}

func (g *game) RemoveSaveFile() error {
	//storage := js.Global.Get("localStorage")
	//storage.Call("removeItem", "boohusave")
	return nil
}

func (g *game) Load() (bool, error) {
	//storage := js.Global.Get("localStorage")
	//if !storage.Bool() {
	//return true, errors.New("localStorage not found")
	//}
	//save := storage.Call("getItem", "boohusave")
	//if !save.Bool() {
	//return false, nil
	//}
	//s, err := base64.StdEncoding.DecodeString(save.String())
	//if err != nil {
	//return true, err
	//}
	//lg, err := g.DecodeGameSave(s)
	//if err != nil {
	//return true, err
	//}
	//*g = *lg

	// // XXX: gob encoding works badly with gopherjs, it seems, some maps get broken
	// g.GeneratedRods = map[rod]bool{}
	// g.GeneratedEquipables = map[equipable]bool{}
	// g.FoundEquipables = map[equipable]bool{Robe: true, Dagger: true}
	// g.GeneratedBands = map[monsterBand]int{}
	// g.KilledMons = map[monsterKind]int{}
	// g.Simellas = make(map[position]int)
	// g.Stairs = make(map[position]bool)
	// g.Collectables = make(map[position]*collectable)
	// g.UnknownDig = map[position]bool{}
	// g.ExclusionsMap = map[position]bool{}
	// g.TemporalWalls = map[position]bool{}
	// g.Clouds = map[position]cloud{}
	// g.Highlight = map[position]bool{}

	// g.Equipables = map[position]equipable{}
	// g.Player.Consumables = map[consumable]int{
	// 	HealWoundsPotion: 1,
	// 	Javelin:          3,
	// }
	// g.Player.Statuses = map[status]int{}
	// g.Player.Aptitudes = map[aptitude]bool{}
	// g.ComputeLOS()

	// g.Rods = map[position]rod{}
	// g.Fungus = map[position]vegetation{}
	// g.Doors = map[position]bool{}

	return true, nil
}

func (g *game) WriteDump() error {
	//storage := js.Global.Get("localStorage")
	//storage.Call("setItem", "boohudump", g.Dump())
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
	cache      map[UICell]*js.Object
	ctx        *js.Object
	width      int
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
	canvas := js.Global.Get("document").Call("getElementById", "gamecanvas")
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("font", "18px monospace")
	mesure := ui.ctx.Call("measureText", "W")
	ui.width = mesure.Get("width").Int() + 1
	ui.cache = make(map[UICell]*js.Object)
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

func (ui *termui) Draw(cell UICell, x, y int) {
	var canvas *js.Object
	if cv, ok := ui.cache[cell]; ok {
		canvas = cv
	} else {
		canvas = js.Global.Get("document").Call("createElement", "canvas")
		ctx := canvas.Call("getContext", "2d")
		canvas.Set("width", ui.width)
		canvas.Set("height", 22)
		ctx.Set("font", ui.ctx.Get("font"))
		ctx.Set("fillStyle", cell.bg.String())
		ctx.Call("fillRect", 0, 0, ui.width, 22)
		ctx.Set("fillStyle", cell.fg.String())
		ctx.Call("fillText", string(cell.r), 0, 18)
		ui.cache[cell] = canvas
	}
	ui.ctx.Call("drawImage", canvas, x*ui.width, 22*y)
}

func (ui *termui) Flush() {
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		if cell.r == ' ' {
			cell.r = ' '
		}
		x, y := ui.GetPos(i)
		ui.Draw(cell, x, y)
		ui.backBuffer[i] = cell
	}
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
	time.Sleep(5 * time.Millisecond)
	ui.DrawDungeonView(g, false)
	return false
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		r := ui.ReadChar()
		switch r {
		case '\x1b', 'E', ' ':
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
	case 'S':
		err = errors.New("Command not available for the web html5 version.")
		return err, true, false
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
	case '\x1b', 'E', ' ':
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
	if r == '\x1b' || r == 'E' || r == ' ' {
		return true
	}
	return ui.CursorCharAction(g, targ, r, pos, data)
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		r := ui.ReadChar()
		switch {
		case r == '\x1b' || r == 'E' || r == ' ':
			return -1, false, errors.New("Do nothing, then.")
		case 97 <= r && int(r) < 97+l:
			return int(r - 97), false, nil
		case r == '?':
			return -1, true, nil
		case r == ' ':
			return -1, false, errors.New("Do nothing, then.")
		}
	}
}
