// +build js

package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"image/color"
	"log"
	"runtime"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gopherjs/gopherwasm/js"
)

func main() {
	tui := &termui{}
	err := tui.Init()
	if err != nil {
		log.Fatalf("boohu: %v\n", err)
	}
	defer tui.Close()

	ApplyDefaultKeyBindings()
	gameConfig.Tiles = true
	tui.PostInit()
	LinkColors()
	gameConfig.DarkLOS = true
	ApplyDarkLOS()

	tui.DrawWelcome()
	g := &game{}
	load, err := g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Printf("Error loading saved game… starting new game. (%v)", err)
	}
	load, err = g.LoadConfig()
	if load && err != nil {
		g.Print("Error loading config file.")
	} else if load {
		CustomKeys = true
	}
	g.ui = tui
	g.EventLoop()
	tui.Clear()
	tui.DrawText("Refresh the page to start again.\nYou can find last game statistics below.", 0, 0)
	tui.DrawText(SaveError, 0, 1)
	tui.Flush()
}

var SaveError string

// io compatibility functions

func (g *game) DataDir() (string, error) {
	return "", nil
}

func (g *game) Save() error {
	if runtime.GOARCH != "wasm" {
		return errors.New("Saving games is not available in the web html version.") // TODO remove when it works
	}
	save, err := g.GameSave()
	if err != nil {
		SaveError = err.Error()
		return err
	}
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		SaveError = "localStorage not found"
		return errors.New("localStorage not found")
	}
	s := base64.StdEncoding.EncodeToString(save)
	storage.Call("setItem", "boohusave", s)
	SaveError = "no errors"
	return nil
}

func (g *game) SaveConfig() error {
	if runtime.GOARCH != "wasm" {
		return nil
	}
	conf, err := gameConfig.ConfigSave()
	if err != nil {
		SaveError = err.Error()
		return err
	}
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		SaveError = "localStorage not found"
		return errors.New("localStorage not found")
	}
	s := base64.StdEncoding.EncodeToString(conf)
	storage.Call("setItem", "boohuconfig", s)
	SaveError = "no errors"
	return nil
}

func (g *game) RemoveSaveFile() error {
	storage := js.Global().Get("localStorage")
	storage.Call("removeItem", "boohusave")
	return nil
}

func (g *game) RemoveDataFile(file string) error {
	storage := js.Global().Get("localStorage")
	storage.Call("removeItem", file)
	return nil
}

func (g *game) Load() (bool, error) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return true, errors.New("localStorage not found")
	}
	save := storage.Call("getItem", "boohusave")
	if save.Type() != js.TypeString || runtime.GOARCH != "wasm" {
		return false, nil
	}
	s, err := base64.StdEncoding.DecodeString(save.String())
	if err != nil {
		return true, err
	}
	lg, err := g.DecodeGameSave(s)
	if err != nil {
		return true, err
	}
	*g = *lg

	// // XXX: gob encoding works badly with gopherjs, it seems, some maps get broken

	return true, nil
}

func (g *game) LoadConfig() (bool, error) {
	storage := js.Global().Get("localStorage")
	if storage.Type() != js.TypeObject {
		return true, errors.New("localStorage not found")
	}
	conf := storage.Call("getItem", "boohuconfig")
	if conf.Type() != js.TypeString || runtime.GOARCH != "wasm" {
		return false, nil
	}
	s, err := base64.StdEncoding.DecodeString(conf.String())
	if err != nil {
		return true, err
	}
	c, err := g.DecodeConfigSave(s)
	if err != nil {
		return true, err
	}
	gameConfig = *c
	if gameConfig.RuneNormalModeKeys == nil || gameConfig.RuneTargetModeKeys == nil {
		ApplyDefaultKeyBindings()
	}
	if !gameConfig.DarkLOS {
		ApplyLightLOS()
	}
	return true, nil
}

func (g *game) WriteDump() error {
	//storage := js.Global.Get("localStorage")
	//storage.Call("setItem", "boohudump", g.Dump())
	pre := js.Global().Get("document").Call("getElementById", "dump")
	pre.Set("innerHTML", g.Dump())
	return nil
}

// End of io compatibility functions

func (ui *termui) Init() error {
	ui.cells = make([]UICell, UIWidth*UIHeight)
	js.Global().Get("document").Call("getElementById", "gamecanvas").Call(
		"addEventListener", "keypress", js.NewEventCallback(0, func(e js.Value) {
			s := e.Get("key").String()
			ch <- jsInput{key: s}
		}))
	js.Global().Get("document").Call(
		"addEventListener", "keypress", js.NewEventCallback(0, func(e js.Value) {
			s := e.Get("key").String()
			if s == " " {
				e.Call("preventDefault")
			}
		}))
	js.Global().Get("document").Call("getElementById", "gamecanvas").Call(
		"addEventListener", "mousedown", js.NewEventCallback(0, func(e js.Value) {
			x, y := ui.GetMousePos(e)
			ch <- jsInput{mouse: true, mouseX: x, mouseY: y, button: e.Get("button").Int()}
		}))
	js.Global().Get("document").Call("getElementById", "gamecanvas").Call(
		"addEventListener", "mousemove", js.NewEventCallback(0, func(e js.Value) {
			x, y := ui.GetMousePos(e)
			if x != ui.mousepos.X || y != ui.mousepos.Y {
				ui.mousepos.X = x
				ui.mousepos.Y = y
				ch <- jsInput{mouse: true, mouseX: x, mouseY: y, button: -1}
			}
		}))
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
	ui.InitElements()
	return nil
}

type UICell struct {
	fg    uicolor
	bg    uicolor
	r     rune
	inMap bool
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
		cl = color.RGBA{211, 1, 2, opaque}
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

var ch chan jsInput
var interrupt chan bool

func init() {
	ch = make(chan jsInput, 100)
	interrupt = make(chan bool)
}

func (ui *termui) Interrupt() {
	interrupt <- true
}

func (ui *termui) Close() {
	// TODO
}

func (ui *termui) PostInit() {
	SolarizedPalette()
	ui.HideCursor()
	//MenuCols[MenuOther] = MenuCols[MenuView]
	//MenuCols[MenuView] = [2]int{-1, -1}
	settingsActions = append(settingsActions, toggleTiles)
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Flush() {
	js.Global().Get("window").Call("requestAnimationFrame", js.NewEventCallback(0, ui.FlushCallback))
}

func (ui *termui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
	if gameConfig.Small {
		ui.ResetCells()
		ui.Flush()
		UIHeight = 24
		UIWidth = 80
	} else {
		UIHeight = 26
		UIWidth = 100
	}
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Set("height", 24*UIHeight)
	canvas.Set("width", 16*UIWidth)
	ui.cache = make(map[UICell]js.Value)
	ui.cells = make([]UICell, UIWidth*UIHeight)
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
}

func (ui *termui) Small() bool {
	return gameConfig.Small
}

func (ui *termui) FlushCallback(obj js.Value) {
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		//if cell.r == ' ' {
		//cell.r = ' '
		//}
		x, y := ui.GetPos(i)
		ui.Draw(cell, x, y)
		ui.backBuffer[i] = cell
	}
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

//func (ui *termui) IsMapCell(x, y int) bool {
//	i := ui.GetIndex(x, y)
//	return i < len(ui.cells) && ui.cells[i].inMap
//}

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

type jsInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *termui) ReadKey(s string) (r rune) {
	bs := strings.NewReader(s)
	r, _, _ = bs.ReadRune()
	return r
}

func (ui *termui) PollEvent() (in jsInput) {
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
				if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
					break
				}
				fallthrough
			case 0:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(in.mouseX)
					if !ok {
						again = true
						break
					}
					err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: m.Key()})
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
		action = ChangeConfig
	case "\x1b", "Escape", " ":
		action = QuitConfig
	case "u":
		n -= DungeonHeight / 2
	case "d":
		n += DungeonHeight / 2
	case "j", "2":
		n++
	case "k", "8":
		n--
	case "R":
		action = ResetConfig
	case "":
		if in.mouse {
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case 0:
				if x > DungeonWidth || y > DungeonHeight {
					action = QuitConfig
				}
			case 1:
				action = QuitConfig
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
				if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
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
					m, ok := ui.WhichButton(in.mouseX)
					if !ok {
						g.Targeting = InvalidPos
						notarg = true
						err = errors.New(DoNothing)
						break
					}
					err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: m.Key()}, data)
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
	}
	if utf8.RuneCountInString(in.key) > 1 {
		g.Printf("Invalid key: “%s”.", in.key)
		notarg = true
		return
	}
	return ui.CursorKeyAction(g, targ, runeKeyAction{r: ui.ReadKey(in.key)}, data)
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
			case 0:
				if y < 0 || y > l || x >= DungeonWidth {
					return -1, false, errors.New(DoNothing)
				}
				if y == 0 {
					return -1, true, nil
				}
				return y - 1, false, nil
			case 2:
				return -1, true, nil
			case 1:
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}
