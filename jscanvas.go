// +build js,!dom

package main

import "github.com/hajimehoshi/gopherwasm/js"

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    js.Value
	cache      map[UICell]js.Value
	ctx        js.Value
	width      int
}

func (ui *termui) InitElements() error {
	canvas := js.Global.Get("document").Call("getElementById", "gamecanvas")
	canvas.Call("addEventListener", "contextmenu", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
	}), false)
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("font", "18px monospace")
	mesure := ui.ctx.Call("measureText", "W")
	ui.width = mesure.Get("width").Int() + 1
	ui.cache = make(map[UICell]js.Value)
	return nil
}

func (ui *termui) Draw(cell UICell, x, y int) {
	var canvas js.Value
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

func (ui *termui) GetMousePos(evt js.Value) (x, y int) {
	canvas := js.Global.Get("document").Call("getElementById", "gamecanvas")
	rect := canvas.Call("getBoundingClientRect")
	x = evt.Get("clientX").Int() - rect.Get("left").Int()
	y = evt.Get("clientY").Int() - rect.Get("top").Int()
	return (x - ui.width/2) / ui.width, (y - 8) / 22
}
