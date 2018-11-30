// +build tcell

package main

import (
	"runtime"

	"github.com/gdamore/tcell"
)

type gameui struct {
	g *game
	tcell.Screen
	cursor position
	small  bool
	// below unused for this backend
	menuHover menu
	itemHover int
}

func (ui *gameui) Init() error {
	screen, err := tcell.NewScreen()
	ui.Screen = screen
	if err != nil {
		return err
	}
	err = ui.Screen.Init()
	if err != nil {
		return err
	}
	ui.Screen.SetStyle(tcell.StyleDefault)
	if runtime.GOOS != "openbsd" {
		ui.Screen.EnableMouse()
	}
	ui.Screen.HideCursor()
	ui.HideCursor()
	ui.menuHover = -1
	return nil
}

func (ui *gameui) Close() {
	ui.Screen.Fini()
}

var SmallScreen = false

func (ui *gameui) Flush() {
	ui.DrawLogFrame()
	for _, cdraw := range ui.g.DrawLog[len(ui.g.DrawLog)-1].Draws {
		cell := cdraw.Cell
		st := tcell.StyleDefault
		fg := cell.Fg
		bg := cell.Bg
		if Only8Colors {
			fg = Map16ColorTo8Color(fg)
			bg = Map16ColorTo8Color(bg)
		}
		st = st.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
		ui.Screen.SetContent(cdraw.X, cdraw.Y, cell.R, nil, st)
	}
	//ui.g.Printf("%d %d %d", ui.g.DrawFrame, ui.g.DrawFrameStart, len(ui.g.DrawLog))
	ui.Screen.Show()
	w, h := ui.Screen.Size()
	if w <= UIWidth-8 || h <= UIHeight-2 {
		SmallScreen = true
	} else {
		SmallScreen = false
	}
}

func (ui *gameui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
	if gameConfig.Small {
		ui.Clear()
		ui.Flush()
		UIHeight = 24
		UIWidth = 80
	} else {
		UIHeight = 26
		UIWidth = 100
	}
	ui.g.DrawBuffer = make([]UICell, UIWidth*UIHeight)
	ui.Clear()
}

func (ui *gameui) Small() bool {
	return gameConfig.Small || SmallScreen
}

func (ui *gameui) Interrupt() {
	ui.Screen.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ui *gameui) PollEvent() (in uiInput) {
	switch tev := ui.Screen.PollEvent().(type) {
	case *tcell.EventKey:
		switch tev.Key() {
		case tcell.KeyEsc:
			in.key = " "
		case tcell.KeyLeft:
			// TODO: will not work if user changes keybindings
			in.key = "4"
		case tcell.KeyDown:
			in.key = "2"
		case tcell.KeyUp:
			in.key = "8"
		case tcell.KeyRight:
			in.key = "6"
		case tcell.KeyHome:
			in.key = "7"
		case tcell.KeyEnd:
			in.key = "1"
		case tcell.KeyPgUp:
			in.key = "9"
		case tcell.KeyPgDn:
			in.key = "3"
		case tcell.KeyDelete:
			in.key = "5"
		case tcell.KeyCtrlW:
			in.key = "W"
		case tcell.KeyCtrlQ:
			in.key = "Q"
		case tcell.KeyCtrlP:
			in.key = "m"
		}
		if tev.Rune() != 0 && in.key == "" {
			in.key = string(tev.Rune())
		}
	case *tcell.EventMouse:
		in.mouseX, in.mouseY = tev.Position()
		switch tev.Buttons() {
		case tcell.Button1:
			in.mouse = true
			in.button = 0
		case tcell.Button2:
			in.mouse = true
			in.button = 1
		case tcell.Button3:
			in.mouse = true
			in.button = 2
		}
	case *tcell.EventInterrupt:
		in.interrupt = true
	}
	return in
}
