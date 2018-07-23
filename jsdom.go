// +build js,dom

// THIS FILE IS NO LONGER UP TO DATE

package main

import (
	"fmt"

	"github.com/gopherjs/gopherwasm/js"
)

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    *js.Object
	cache      []*js.Object
}

func (ui *termui) InitElements() error {
	doc := js.Global().Get("document")
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

func (ui *termui) GetMousePos(evt *js.Object) (x, y int) {
	// TODO
	return
}
