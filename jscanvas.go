// +build js,!dom

package main

import "github.com/gopherjs/gopherjs/js"

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    *js.Object
	cache      map[UICell]*js.Object
	ctx        *js.Object
	width      int
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
