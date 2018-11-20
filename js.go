// +build js

package main

import (
	"encoding/base64"
	"errors"
	"log"
	"runtime"

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

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}
