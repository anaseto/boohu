// NOTE: The js version is way too slow, don't know why.
// NOTE: This is just an early draft, nothing finished.

// +build ebiten

package main

import (
	"errors"
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
		log.Fatal(os.Stderr, "boohu: %v\n", err)
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
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Boohu"); err != nil {
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
	screenWidth  = 820
	screenHeight = 240
)

type ebitenKey struct {
	key ebiten.Key
	ch  rune
}

var (
	canvasImage *ebiten.Image

	flush     chan bool
	flushDone chan bool
	Interrupt chan bool
	wantKey   chan bool
	EbKey     chan ebitenKey
	InputEnd  chan bool

	Keys = []ebiten.Key{
		ebiten.KeyEscape,
		ebiten.KeySpace,
		ebiten.KeyEnter,
		ebiten.KeyControl,
	}
)

func init() {
	flush = make(chan bool)
	flushDone = make(chan bool)
	Interrupt = make(chan bool)
	wantKey = make(chan bool)
	EbKey = make(chan ebitenKey)
	InputEnd = make(chan bool)
}

type termui struct {
}

func (ui *termui) Init() error {
	var err error
	canvasImage, err = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterNearest)
	if err != nil {
		return err
	}
	canvasImage.Fill(color.RGBA{0, 43, 54, 0xff})
	return nil
}

func UpdateKeys() bool {
	//keyStates := map[ebiten.Key]int{}
	var ekey ebiten.Key
	for _, key := range Keys {
		if !ebiten.IsKeyPressed(key) {
			continue
		}
		ekey = key
		break
		//keyStates[key]++
	}
	chars := ebiten.InputChars()
	found := false
	var ch rune
	if len(ebiten.InputChars()) > 0 {
		ch = chars[0]
		found = true
	}
	//for key, value := range keyStates {
	//if value > 0 {
	//ekey = key
	//found = true
	//break
	//}
	//}
	if found {
		go func() {
			EbKey <- ebitenKey{key: ekey, ch: ch}
		}()
		select {
		case <-InputEnd:
		}
	}
	return found
}

var WantsEbKey bool

func update(screen *ebiten.Image) error {
	select {
	case <-wantKey:
		WantsEbKey = true
	default:
	}
	if WantsEbKey {
		found := UpdateKeys()
		if found {
			WantsEbKey = false
		}
	}
	if ebiten.IsRunningSlowly() {
		return nil
	}
	screen.DrawImage(canvasImage, nil)
	return nil
}

func (ui *termui) Close() {
}

func (ui *termui) PostInit() {
}

func (ui *termui) Clear() {
	canvasImage.Fill(color.RGBA{0, 43, 54, 0xff})
}

func (ui *termui) Flush() {
	//flush <- true
	//<-flushDone
}

func (ui *termui) Interrupt() {
	Interrupt <- true
}

func (ui *termui) HideCursor() {
}

func (ui *termui) SetCursor(pos position) {
}

var uiFont font.Face

const FontSize = 8

func init() {
	tt, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Fatal(err)
	}
	uiFont = truetype.NewFace(tt, &truetype.Options{
		Size: FontSize,
		DPI:  72,
		//DPI:     100,
		Hinting: font.HintingFull,
	})
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	// TODO: background?
	text.Draw(canvasImage, string(r), uiFont, FontSize*x, FontSize*y, color.White)
}

func (ui *termui) Reverse(c uicolor) uicolor {
	return uicolor(1)
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		wantKey <- true
		select {
		case ekey := <-EbKey:
			InputEnd <- true
			if ekey.key == ebiten.KeyEscape || ekey.key == ebiten.KeySpace {
				break loop
			}
			if ekey.ch == ' ' {
				break loop
			}
		}
		// TODO: mouse
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		wantKey <- true
		select {
		case ekey := <-EbKey:
			InputEnd <- true
			if ekey.ch == 'Y' || ekey.ch == 'y' {
				return true
			}
			return false
		}
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		wantKey <- true
		select {
		case <-EbKey:
			InputEnd <- true
			return nil
		case <-Interrupt:
			InputEnd <- true
			return errors.New("interrupted")
		}
		// TODO: mouse
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	wantKey <- true
	select {
	case ekey := <-EbKey:
		InputEnd <- true
		// TODO: other cases
		switch ekey.ch {
		case 'S', 'Q', '#':
			err = errors.New("Command not available for the js version.")
		default:
			err, again, quit = ui.HandleCharacter(g, ev, ekey.ch)
		}
		// TODO: mouse
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	wantKey <- true
	select {
	case ekey := <-EbKey:
		InputEnd <- true
		if ekey.ch == 0 && (ekey.key == ebiten.KeyEscape || ekey.key == ebiten.KeySpace) {
			quit = true
			return n, quit
		}
		switch ekey.ch {
		case 'u':
			n -= 12
		case 'd':
			n += 12
		case 'j':
			n++
		case 'k':
			n--
		case ' ':
			quit = true
		}
		// TODO: mouse
	}
	return n, quit
}

func (ui *termui) TargetModeEvent(g *game, targ Targetter, pos position, data *examineData) bool {
	wantKey <- true
	select {
	case ekey := <-EbKey:
		InputEnd <- true
		if ekey.ch == 0 {
			if ekey.key == ebiten.KeyEscape {
				return true
			}
		}
		if ui.CursorCharAction(g, targ, ekey.ch, pos, data) {
			return true
		}
		// TODO: mouse
	}
	return false
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		wantKey <- true
		select {
		case ekey := <-EbKey:
			InputEnd <- true
			if ekey.ch == 0 && (ekey.key == ebiten.KeyEscape || ekey.key == ebiten.KeySpace) {
				return -1, false, errors.New("Ok, then.")
			}
			if 97 <= ekey.ch && int(ekey.ch) < 97+l {
				return int(ekey.ch - 97), false, nil
			}
			if ekey.ch == '?' {
				return -1, true, nil
			}
			if ekey.ch == ' ' {
				return -1, false, errors.New("Ok, then.")
			}
			// TODO: mouse
		}
	}
}
