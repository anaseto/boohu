// +build !js,!tk

package main

import "unicode/utf8"

func (ui *termui) ApplyToggleTiles() {
}

func (ui *termui) ColorLine(y int, fg uicolor) {
}

func (ui *termui) SetMapCell(x, y int, r rune, fg, bg uicolor) {
	ui.SetCell(x, y, r, fg, bg)
}

func (ui *termui) KeyToRuneKeyAction(in uiInput) rune {
	if utf8.RuneCountInString(in.key) != 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}