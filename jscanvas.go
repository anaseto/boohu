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
	mousepos   position
}

var Tiles bool = true

func (ui *termui) InitElements() error {
	canvas := js.Global().Get("document").Call("getElementById", "gamecanvas")
	canvas.Call("addEventListener", "contextmenu", js.NewEventCallback(js.PreventDefault, func(e js.Value) {
	}), false)
	canvas.Call("setAttribute", "tabindex", "1")
	ui.ctx = canvas.Call("getContext", "2d")
	ui.ctx.Set("imageSmoothingEnabled", false)
	//if Tiles {
	//ui.ctx.Set("font", "22px monospace")
	//} else {
	//ui.ctx.Set("font", "18px monospace")
	//}
	//if Tiles {
	ui.width = 16
	ui.height = 24
	canvas.Set("height", 24*UIHeight)
	canvas.Set("width", 16*UIWidth)
	//} else {
	//ui.height = 22
	//mesure := ui.ctx.Call("measureText", "W")
	//ui.width = mesure.Get("width").Int() + 1
	//canvas.Set("height", ui.height*UIHeight)
	//canvas.Set("width", ui.width*UIWidth)
	//}
	// seems to be needed again
	//if Tiles {
	//ui.ctx.Set("font", "22px monospace")
	//} else {
	//ui.ctx.Set("font", "18px monospace")
	//}
	ui.cache = make(map[UICell]js.Value)
	return nil
}

var TileImgs map[string][]byte

var MapNames = map[rune]string{
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'☻':  "dreaming",
	'♫':  "footsteps",
	'#':  "wall",
	'@':  "player",
	'§':  "fog",
	'♣':  "simella",
	'+':  "door",
	'.':  "ground",
	'"':  "foliage",
	'•':  "tick",
	'●':  "rock",
	',':  "comma",
	'}':  "rbrace",
	'%':  "percent",
	':':  "colon",
	'\\': "backslash",
	'~':  "tilde",
	'☼':  "sun",
	'*':  "asterisc",
	'—':  "hbar",
	'/':  "slash",
	'|':  "vbar",
	'∞':  "kill",
	' ':  "space",
	'[':  "lbracket",
	']':  "rbracket",
	')':  "rparen",
	'(':  "lparen",
	'>':  "stairs",
	'!':  "potion",
}

var LetterNames = map[rune]string{
	'(':  "lparen",
	')':  "rparen",
	'@':  "player",
	'{':  "lbrace",
	'}':  "rbrace",
	'[':  "lbracket",
	']':  "rbracket",
	'♪':  "music1",
	'♫':  "music2",
	'•':  "tick",
	'♣':  "simella",
	' ':  "space",
	'!':  "exclamation",
	'?':  "interrogation",
	',':  "comma",
	':':  "colon",
	';':  "semicolon",
	'\'': "quote",
	'—':  "longhyphen",
	'-':  "hyphen",
	'|':  "pipe",
	'/':  "slash",
	'\\': "backslash",
	'%':  "percent",
	'┐':  "boxne",
	'┤':  "boxe",
	'│':  "vbar",
	'┘':  "boxse",
	'─':  "hbar",
	'►':  "arrow",
	'×':  "times",
	'.':  "dot",
	'#':  "hash",
	'"':  "quotes",
	'+':  "plus",
	'“':  "lquotes",
	'”':  "rquotes",
	'=':  "equal",
	'>':  "gt",
	'¤':  "frontier",
	'√':  "hit",
	'Φ':  "magic",
	'§':  "fog",
	'●':  "rock",
	'~':  "tilde",
	'☼':  "sun",
	'*':  "asterisc",
	'∞':  "kill",
	'☻':  "dreaming",
}

func getImage(cell UICell) []byte {
	var pngImg []byte
	if cell.inMap && Tiles {
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
		//if Tiles {
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
		//} else {
		//canvas.Set("width", ui.width)
		//canvas.Set("height", ui.height)
		//ctx := canvas.Call("getContext", "2d")
		//ctx.Set("imageSmoothingEnabled", false)
		//ctx.Set("font", ui.ctx.Get("font"))
		//ctx.Set("fillStyle", cell.bg.String())
		//ctx.Call("fillRect", 0, 0, ui.width, ui.height)
		//ctx.Set("fillStyle", cell.fg.String())
		////if Tiles {
		////ctx.Call("fillText", string(cell.r), 0, 18)
		////} else {
		//ctx.Call("fillText", string(cell.r), 0, 18)
		////}
		//}
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
