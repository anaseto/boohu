// +build js,!dom

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/draw"
	"image/png"

	"github.com/gopherjs/gopherwasm/js"
)

type termui struct {
	cells      []UICell
	backBuffer []UICell
	cursor     position
	display    js.Value
	cache      map[UICell]js.Value
	ctx        js.Value
	width      int
	height     int
}

var Tiles bool = false

func (ui *termui) InitElements() error {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Call("addEventListener", "contextmenu", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
	}), false)
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("font", "18px monospace")
	if Tiles {
		ui.width = 16
		ui.height = 24
		canvas.Set("height", 24*UIHeight)
		canvas.Set("width", 16*UIWidth)
	} else {
		ui.height = 22
		mesure := ui.ctx.Call("measureText", "W")
		ui.width = mesure.Get("width").Int() + 1
		canvas.Set("height", ui.height*UIHeight)
		canvas.Set("width", ui.width*UIWidth)
	}
	ui.ctx.Set("font", "18px monospace") // seems to be needed again
	ui.cache = make(map[UICell]js.Value)
	return nil
}

var TileImgs map[string][]byte

func getImage(cell UICell) []byte {
	buf := make([]byte, len(TileImgs["img"]))
	base64.StdEncoding.Decode(buf, TileImgs["img"]) // TODO: check error
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
		if Tiles {
			canvas.Set("width", 16)
			canvas.Set("height", 24)
			ctx := canvas.Call("getContext", "2d")
			buf := getImage(cell)
			ta := js.TypedArrayOf(buf)
			ca := js.Global().Get("Uint8ClampedArray").New(ta)
			imgdata := js.Global().Get("ImageData").New(ca, 16, 24)
			ctx.Call("putImageData", imgdata, 0, 0)
			ta.Release()
		} else {
			canvas.Set("width", ui.width)
			canvas.Set("height", ui.height)
			ctx := canvas.Call("getContext", "2d")
			ctx.Set("font", ui.ctx.Get("font"))
			ctx.Set("fillStyle", cell.bg.String())
			ctx.Call("fillRect", 0, 0, ui.width, ui.height)
			ctx.Set("fillStyle", cell.fg.String())
			ctx.Call("fillText", string(cell.r), 0, 18)
		}
		ui.cache[cell] = canvas
	}
	ui.ctx.Call("drawImage", canvas, x*ui.width, ui.height*y)
}

func (ui *termui) GetMousePos(evt js.Value) (x, y int) {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	rect := canvas.Call("getBoundingClientRect")
	x = evt.Get("clientX").Int() - rect.Get("left").Int()
	y = evt.Get("clientY").Int() - rect.Get("top").Int()
	return (x - ui.width/2) / ui.width, (y - (ui.height/2 - 3)) / ui.height
}
