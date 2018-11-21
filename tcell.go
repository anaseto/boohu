// +build tcell

package main

import (
	"runtime"

	"github.com/gdamore/tcell"
)

type termui struct {
	tcell.Screen
	cursor position
	small  bool
	// below unused for this backend
	menuHover menu
	itemHover int
}

func (ui *termui) Init() error {
	screen, err := tcell.NewScreen()
	ui.Screen = screen
	if err != nil {
		return err
	}
	return ui.Screen.Init()
}

func (ui *termui) Close() {
	ui.Screen.Fini()
}

func (ui *termui) PostInit() {
	ui.Screen.SetStyle(tcell.StyleDefault)
	if runtime.GOOS != "openbsd" {
		ui.Screen.EnableMouse()
	}
	ui.Screen.HideCursor()
	ui.HideCursor()
	ui.menuHover = -1
}

func (ui *termui) Clear() {
	w, h := ui.Screen.Size()
	if w > UIWidth {
		w = UIWidth
	}
	if h > UIHeight {
		h = UIHeight
	}
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(ColorFg)).Background(tcell.Color(ColorBg))
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			ui.Screen.SetContent(col, row, ' ', nil, st)
		}
	}
}

var SmallScreen = false

func (ui *termui) Flush() {
	ui.Screen.Show()
	w, h := ui.Screen.Size()
	if w <= UIWidth-8 || h <= UIHeight-2 {
		SmallScreen = true
	} else {
		SmallScreen = false
	}
}

func (ui *termui) ApplyToggleLayout() {
	gameConfig.Small = !gameConfig.Small
}

func (ui *termui) Small() bool {
	return gameConfig.Small || SmallScreen
}

func (ui *termui) Interrupt() {
	ui.Screen.PostEvent(tcell.NewEventInterrupt(nil))
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	st := tcell.StyleDefault
	st = st.Foreground(tcell.Color(fg)).Background(tcell.Color(bg))
	ui.Screen.SetContent(x, y, r, nil, st)
}

func (ui *termui) PollEvent() (in uiInput) {
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
