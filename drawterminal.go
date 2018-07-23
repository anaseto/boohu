// +build !js

package main

func (ui *termui) DrawMenus(g *game) {
	line := DungeonHeight
	for i, cols := range MenuCols {
		if cols[0] >= 0 {
			ui.DrawColoredText(menu(i).String(), cols[0], line, ColorViolet)
		}
	}
}
