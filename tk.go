// +build tk

package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"unicode/utf8"

	"github.com/nsf/gothic"
)

func main() {
	opt := flag.Bool("s", false, "Use true 16-color solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	optCenteredCamera := flag.Bool("c", false, "centered camera")
	optMinimalUI := flag.Bool("m", false, "80x24 minimal UI")
	optNoAnim := flag.Bool("n", false, "no animations")
	flag.Parse()
	if *opt {
		SolarizedPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optCenteredCamera {
		CenteredCamera = true
	}
	if *optNoAnim {
		DisableAnimations = true
	}

	tui := &termui{}
	if *optMinimalUI {
		gameConfig.Small = true
		UIHeight = 24
		UIWidth = 80
	}
	err := tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	ApplyDefaultKeyBindings()
	tui.PostInit()
	gameConfig.Tiles = true
	LinkColors()
	gameConfig.DarkLOS = true

	g := &game{}
	load, err := g.LoadConfig()
	if load && err != nil {
		g.Print("Error loading config file.")
	} else if load {
		CustomKeys = true
	}

	tui.DrawWelcome()
	load, err = g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Print("Error loading saved game… starting new game.")
	}
	g.ui = tui
	g.EventLoop()
}

type termui struct {
	ir         *gothic.Interpreter
	cells      []UICell
	backBuffer []UICell
	cursor     position
	stty       string
	cache      map[UICell]string
	width      int
	height     int
	mousepos   position
	menuHover  menu
	itemHover  int
}

func (ui *termui) Init() error {
	ui.cells = make([]UICell, UIWidth*UIHeight)
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
	ui.ir = gothic.NewInterpreter(`
set width [expr {16 * 100}]
set height [expr {24 * 26}]
set can [canvas .c -width $width -height $height -background black]
grid $can -row 0 -column 0
focus $can
image create photo gamescreen -width $width -height $height -palette 256/256/256
image create photo bufscreen -width $width -height $height -palette 256/256/256
$can create image 0 0 -anchor nw -image gamescreen
bind $can <Key> {
	GetKey %A %K
}
bind $can <Motion> {
	MouseMotion %x %y
}
bind $can <ButtonPress> {
	MouseDown %x %y %b
}
`)
	ui.ir.RegisterCommand("GetKey", func(c, keysym string) {
		var s string
		if c != "" {
			s = c
		} else {
			s = keysym
		}
		//fmt.Printf("“%s” “%s”\n", c, keysym)
		ch <- uiInput{key: s}
	})
	ui.ir.RegisterCommand("MouseDown", func(x, y, b int) {
		ch <- uiInput{mouse: true, mouseX: (x - 1) / ui.width, mouseY: (y - 1) / ui.height, button: b - 1}
	})
	ui.ir.RegisterCommand("MouseMotion", func(x, y int) {
		nx := (x - 1) / ui.width
		ny := (y - 1) / ui.height
		if nx != ui.mousepos.X || ny != ui.mousepos.Y {
			ui.mousepos.X = nx
			ui.mousepos.Y = ny
			ch <- uiInput{mouse: true, mouseX: nx, mouseY: ny, button: -1}
		}
	})
	ui.menuHover = -1
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
	ui.InitElements()
	return nil
}

func (ui *termui) InitElements() error {
	ui.width = 16
	ui.height = 24
	ui.cache = make(map[UICell]string)
	return nil
}

var ch chan uiInput
var interrupt chan bool

func init() {
	ch = make(chan uiInput, 100)
	interrupt = make(chan bool)
}

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *termui) Close() {
}

func (ui *termui) PostInit() {
	SolarizedPalette()
	ui.HideCursor()
	settingsActions = append(settingsActions, toggleTiles)
}

func (ui *termui) Flush() {
	xmin := UIWidth - 1
	xmax := 0
	ymin := UIHeight - 1
	ymax := 0
	for i := 0; i < len(ui.cells); i++ {
		if ui.cells[i] == ui.backBuffer[i] {
			continue
		}
		cell := ui.cells[i]
		x, y := ui.GetPos(i)
		ui.Draw(cell, x, y)
		ui.backBuffer[i] = cell
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
	}
	ui.ir.Eval("gamescreen copy bufscreen -from %{0%d} %{1%d} %{2%d} %{3%d} -to %{0%d} %{1%d} %{2%d} %{3%d}",
		xmin*16, ymin*24, (xmax+1)*16, (ymax+1)*24) // TODO: optimize this more
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
	ui.cache = make(map[UICell]string)
	ui.cells = make([]UICell, UIWidth*UIHeight)
	ui.ResetCells()
	ui.backBuffer = make([]UICell, UIWidth*UIHeight)
}

func getImage(cell UICell) string {
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
		return ""
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
	pngbuf := &bytes.Buffer{}
	png.Encode(pngbuf, rgbaimg)
	return base64.StdEncoding.EncodeToString(pngbuf.Bytes())
}

func (ui *termui) Draw(cell UICell, x, y int) {
	var img string
	if im, ok := ui.cache[cell]; ok {
		img = im
	} else {
		img = getImage(cell)
		ui.cache[cell] = img
	}
	ui.ir.Eval("::bufscreen put %{%s} -format png -to %{%d} %{%d}", img, x*ui.width, ui.height*y)
}

func (ui *termui) KeyToRuneKeyAction(in uiInput) rune {
	switch in.key {
	case "Enter":
		in.key = "."
	case "Left", "KP_Left":
		in.key = "4"
	case "Right", "KP_Right":
		in.key = "6"
	case "Up", "KP_Up":
		in.key = "8"
	case "Down", "KP_Down":
		in.key = "2"
	case "KP_Home":
		in.key = "7"
	case "KP_End":
		in.key = "1"
	case "KP_Prior":
		in.key = "9"
	case "KP_Next":
		in.key = "3"
	case "KP_Begin", "KP_Delete":
		in.key = "5"
	}
	if utf8.RuneCountInString(in.key) > 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}
