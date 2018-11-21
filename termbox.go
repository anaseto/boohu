// +build !tcell,!ansi,!js,!tk

package main

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

type termui struct {
	cursor position
	small  bool
	// below unused for this backend
	menuHover menu
	itemHover int
}

func (ui *termui) Init() error {
	return termbox.Init()
}

func (ui *termui) Close() {
	termbox.Close()
}

func (ui *termui) PostInit() {
	FixColor()
	termbox.SetOutputMode(termbox.Output256)
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	termbox.HideCursor()
	ui.HideCursor()
}

func (ui *termui) Clear() {
	termbox.Clear(termbox.Attribute(ColorFg), termbox.Attribute(ColorBg))
}

var SmallScreen = false

func (ui *termui) Flush() {
	termbox.Flush()
	w, h := termbox.Size()
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
	termbox.Interrupt()
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) SetCell(x, y int, r rune, fg, bg uicolor) {
	termbox.SetCell(x, y, r, termbox.Attribute(fg), termbox.Attribute(bg))
}

func (ui *termui) SetMapCell(x, y int, r rune, fg, bg uicolor) {
	ui.SetCell(x, y, r, fg, bg)
}

func (ui *termui) PollEvent() (in uiInput) {
	switch tev := termbox.PollEvent(); tev.Type {
	case termbox.EventKey:
		if tev.Ch == 0 {
			switch tev.Key {
			case termbox.KeyArrowLeft:
				in.key = "4"
			case termbox.KeyArrowDown:
				in.key = "2"
			case termbox.KeyArrowUp:
				in.key = "8"
			case termbox.KeyArrowRight:
				in.key = "6"
			case termbox.KeyHome:
				in.key = "7"
			case termbox.KeyEnd:
				in.key = "1"
			case termbox.KeyPgup:
				in.key = "9"
			case termbox.KeyPgdn:
				in.key = "3"
			case termbox.KeyDelete:
				in.key = "5"
			case termbox.KeyEsc, termbox.KeySpace:
				in.key = " "
			case termbox.KeyEnter:
				in.key = "."
			}
		}
		if tev.Ch != 0 && in.key == "" {
			in.key = string(tev.Ch)
		}
	case termbox.EventMouse:
		if tev.Ch == 0 {
			in.mouseX, in.mouseY = tev.MouseX, tev.MouseY
			switch tev.Key {
			case termbox.MouseLeft:
				in.mouse = true
				in.button = 0
			case termbox.MouseMiddle:
				in.mouse = true
				in.button = 1
			case termbox.MouseRight:
				in.mouse = true
				in.button = 2
			}
		}
	case termbox.EventInterrupt:
		in.interrupt = true
	}
	return in
}

func (ui *termui) KeyToRuneKeyAction(in uiInput) rune {
	if utf8.RuneCountInString(in.key) != 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}
