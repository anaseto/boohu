// +build ebiten

package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

var Version string = "v0.4"

func main() {
	tui := &termui{}
	err := tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	tui.PostInit()

	go func() {
		tui.DrawWelcome()
		g := &game{}
		g.InitLevel()
		g.ui = tui
		g.EventLoop()
		// TODO: exit
	}()
	if err := ebiten.Run(update, screenWidth, screenHeight, 2.0, "Boohu"); err != nil {
		log.Fatal(err)
	}
}

// io compatibility functions

func (g *game) DataDir() (string, error) {
	return "", nil
}

func (g *game) Save() error {
	return nil
}

func (g *game) RemoveSaveFile() error {
	return nil
}

func (g *game) Load() (bool, error) {
	return false, nil
}

func (g *game) WriteDump() error {
	return nil
}

// End of io compatibility functions

const (
	screenWidth  = 320
	screenHeight = 240
)

var canvasImage *ebiten.Image
var flush chan bool
var flushDone chan bool
var interrupt chan bool

type termui struct {
}

func WindowsPalette() {
}

func (ui *termui) Init() error {
	var err error
	canvasImage, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)
	if err != nil {
		return err
	}
	canvasImage.Fill(color.Black)
	return nil
}

func update(screen *ebiten.Image) error {
	if ebiten.IsRunningSlowly() {
		return nil
	}
	select {
	case <-flush:
		screen.DrawImage(canvasImage, nil)
		flushDone <- true
	default:
	}
	return nil
}

func (ui *termui) Close() {
}

func (ui *termui) PostInit() {
}

func (ui *termui) Clear() {
	canvasImage.Fill(color.Black)
}

func (ui *termui) Flush() {
	flush <- true
	<-flushDone
}

func (ui *termui) Interrupt() {
	interrupt <- true
}

func (ui *termui) HideCursor() {
}

func (ui *termui) SetCursor(pos position) {
}

var uiFont font.Face

func init() {
	tt, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}
	uiFont = truetype.NewFace(tt, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	// TODO: background?
	text.Draw(canvasImage, string(r), uiFont, x, y, color.White)
}

func (ui *termui) Reverse(c uicolor) uicolor {
	return uicolor(1)
}

func (ui *termui) WaitForContinue(g *game) {
}

func (ui *termui) PromptConfirmation(g *game) bool {
	return false
}

func (ui *termui) PressAnyKey() error {
	return nil
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	return nil, false, false
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	return 0, false
}

func (ui *termui) TargetModeEvent(g *game, targ Targetter, pos position, data *examineData) bool {
	return false
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	return 0, false, nil
}
