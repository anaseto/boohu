// +build js,dom

package main

import (
	"fmt"

	"github.com/gopherjs/gopherjs/js"
)

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    *js.Object
	cache      []*js.Object
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
	doc := js.Global.Get("document")
	pre := doc.Call("getElementById", "game")
	for y := 0; y < UIHeight; y++ {
		for x := 0; x < UIWidth; x++ {
			sp := doc.Call("createElement", "span")
			//sp.Call("setAttribute", "id", fmt.Sprintf("p%d-%d", x, y))
			ui.cache = append(ui.cache, sp)
			pre.Call("appendChild", sp)
		}
		text := doc.Call("createTextNode", "\n")
		pre.Call("appendChild", text)
	}
	return nil
}

func uiidx(x, y int) int {
	return y*UIWidth + x
}

func (ui *termui) Draw(cell UICell, x, y int) {
	c := ui.cache[uiidx(x, y)]
	c.Set("textContent", string(cell.r))
	c.Call("setAttribute", "class", fmt.Sprintf("fg%d bg%d", cell.fg, cell.bg))
}
