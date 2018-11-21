// +build js tk

package main

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
)

type UICell struct {
	fg    uicolor
	bg    uicolor
	r     rune
	inMap bool
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

func (ui *termui) ApplyToggleTiles() {
	gameConfig.Tiles = !gameConfig.Tiles
	for c, _ := range ui.cache {
		if c.inMap {
			delete(ui.cache, c)
		}
	}
	for i := 0; i < len(ui.backBuffer); i++ {
		ui.backBuffer[i] = UICell{}
	}
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
		cl = color.RGBA{220, 50, 47, opaque}
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
	'×':  "times",
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
	';':  "semicolon",
	'_':  "stone",
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
	'…':  "dots",
	'_':  "stone",
}

func (ui *termui) Interrupt() {
	interrupt <- true
}

func (ui *termui) Clear() {
	ui.ResetCells()
}

func (ui *termui) Small() bool {
	return gameConfig.Small
}

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

func (ui *termui) PollEvent() (in uiInput) {
	select {
	case in = <-ch:
	case in.interrupt = <-interrupt:
	}
	return in
}

func (ui *termui) ColorLine(y int, fg uicolor) {
	for x := 0; x < DungeonWidth; x++ {
		i := ui.GetIndex(x, y)
		c := ui.cells[i]
		ui.SetCell(x, y, c.r, fg, c.bg)
	}
}

func getImage(cell UICell) *image.RGBA {
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
		log.Printf("Could not decode png: %v", err)
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
	return rgbaimg
}
