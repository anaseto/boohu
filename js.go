// +build js

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image"
	"image/draw"
	"image/png"
	"log"
	"runtime"
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
	for {
		newGame(tui)
	}
}

func newGame(tui *termui) {
	g := &game{}
	load, err := g.LoadConfig()
	if load && err != nil {
		g.Print("Error loading config file.")
	} else if load {
		CustomKeys = true
		if gameConfig.Small {
			gameConfig.Small = false
			tui.ApplyToggleLayout()
		}
	}
	tui.DrawWelcome()
	load, err = g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Printf("Error loading saved game… starting new game. (%v)", err)
	}
	g.ui = tui
	g.EventLoop()
	tui.Clear()
	tui.DrawColoredText("Do you want to collect some more simellas today?\n\n───Click or press any key to play again───", 7, 5, ColorFg)
	tui.DrawText(SaveError, 0, 10)
	tui.Flush()
	tui.PressAnyKey()
}

var SaveError string

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    js.Value
	cache      map[UICell]js.Value
	ctx        js.Value
	width      int
	height     int
	mousepos   position
	menuHover  menu
	itemHover  int
}

func (ui *termui) InitElements() error {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Call("addEventListener", "contextmenu", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
	}), false)
	canvas.Call("setAttribute", "tabindex", "1")
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("imageSmoothingEnabled", false)
	ui.width = 16
	ui.height = 24
	canvas.Set("height", 24*UIHeight)
	canvas.Set("width", 16*UIWidth)
	ui.cache = make(map[UICell]js.Value)
	return nil
}

func getImage(cell UICell) []byte {
	var pngImg []byte
	if cell.inMap && gameConfig.Tiles {
		pngImg = TileImgs["map-notile"]
		if im, ok := TileImgs["map-"+string(cell.r)]; ok {
			pngImg = im
		} else if im, ok := TileImgs["map-"+MapNames[cell.r]]; ok {
			pngImg = im
		}
	} else {
		pngImg = TileImgs["map-notile"]
		if im, ok := TileImgs["letter-"+string(cell.r)]; ok {
			pngImg = im
		} else if im, ok := TileImgs["letter-"+LetterNames[cell.r]]; ok {
			pngImg = im
		}
	}
	buf := make([]byte, len(pngImg))
	base64.StdEncoding.Decode(buf, pngImg) // TODO: check error
	br := bytes.NewReader(buf)
	img, err := png.Decode(br)
	if err != nil {
		js.Global().Get("console").Call("log", "could not decode png")
	}
	rect := img.Bounds()
	rgbaimg := image.NewRGBA(rect)
	draw.Draw(rgbaimg, rect, img, rect.Min, draw.Src)
	bgc := cell.bg.Color()
	fgc := cell.fg.Color()
	for y := 0; y < rect.Max.Y; y++ {
		for x := 0; x < rect.Max.X; x++ {
			c := rgbaimg.At(x, y)
			r, _, _, _ := c.RGBA()
			if r == 0 {
				rgbaimg.Set(x, y, bgc)
			} else {
				rgbaimg.Set(x, y, fgc)
			}
		}
	}
	buf = rgbaimg.Pix
	return buf
}

func (ui *termui) Draw(cell UICell, x, y int) {
	var canvas js.Value
	if cv, ok := ui.cache[cell]; ok {
		canvas = cv
	} else {
		canvas = js.Global().Get("document").Call("createElement", "canvas")
		canvas.Set("width", 16)
		canvas.Set("height", 24)
		ctx := canvas.Call("getContext", "2d")
		ctx.Set("imageSmoothingEnabled", false)
		buf := getImage(cell)
		ta := js.TypedArrayOf(buf)
		ca := js.Global().Get("Uint8ClampedArray").New(ta)
		imgdata := js.Global().Get("ImageData").New(ca, 16, 24)
		ctx.Call("putImageData", imgdata, 0, 0)
		ta.Release()
		ui.cache[cell] = canvas
	}
	ui.ctx.Call("drawImage", canvas, x*ui.width, ui.height*y)
}

func (ui *termui) GetMousePos(evt js.Value) (x, y int) {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	rect := canvas.Call("getBoundingClientRect")
	x = evt.Get("clientX").Int() - rect.Get("left").Int()
	y = evt.Get("clientY").Int() - rect.Get("top").Int()
	return (x - 1) / ui.width, (y - 1) / ui.height
}

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
	SaveError = ""
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
	SaveError = ""
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
	pre := js.Global().Get("document").Call("getElementById", "dump")
	pre.Set("innerHTML", g.Dump())
	return nil
}

// End of io compatibility functions

func (ui *termui) Init() error {
	ui.cells = make([]UICell, UIWidth*UIHeight)
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Call(
		"addEventListener", "keypress", js.NewEventCallback(0, func(e js.Value) {
			s := e.Get("key").String()
			if s == "Unidentified" {
				s = e.Get("code").String()
			}
			ch <- uiInput{key: s}
		}))
	js.Global().Get("document").Call(
		"addEventListener", "keypress", js.NewEventCallback(0, func(e js.Value) {
			if !e.Get("ctrlKey").Bool() && !e.Get("metaKey").Bool() && js.Global().Get("Object").Call("is", js.Global().Get("document").Get("activeElement"), canvas).Bool() {
				e.Call("preventDefault")
			}
		}))
	canvas.Call(
		"addEventListener", "mousedown", js.NewEventCallback(0, func(e js.Value) {
			x, y := ui.GetMousePos(e)
			ch <- uiInput{mouse: true, mouseX: x, mouseY: y, button: e.Get("button").Int()}
		}))
	canvas.Call(
		"addEventListener", "mousemove", js.NewEventCallback(0, func(e js.Value) {
			x, y := ui.GetMousePos(e)
			if x != ui.mousepos.X || y != ui.mousepos.Y {
				ui.mousepos.X = x
				ui.mousepos.Y = y
				ch <- uiInput{mouse: true, mouseX: x, mouseY: y, button: -1}
			}
		}))
	ui.menuHover = -1
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
	ui.InitElements()
	return nil
}

var ch chan uiInput
var interrupt chan bool

func init() {
	ch = make(chan uiInput, 100)
	interrupt = make(chan bool)
}

func (ui *termui) Close() {
	// TODO
}

func (ui *termui) PostInit() {
	SolarizedPalette()
	ui.HideCursor()
	settingsActions = append(settingsActions, toggleTiles)
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

func (ui *termui) FlushCallback(obj js.Value) {
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		x, y := ui.GetPos(i)
		ui.Draw(cell, x, y)
		ui.backBuffer[i] = cell
	}
}

func (ui *termui) KeyToRuneKeyAction(in uiInput) rune {
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
	if utf8.RuneCountInString(in.key) != 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}
