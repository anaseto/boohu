// +build !js,!tk

package main

func (ui *termui) DrawMenus(g *game) {
	line := DungeonHeight
	for i, cols := range MenuCols[0 : len(MenuCols)-1] {
		if cols[0] >= 0 {
			ui.DrawColoredText(menu(i).String(), cols[0], line, ColorViolet)
		}
	}
	interactMenu := ui.UpdateInteractButton(g)
	if interactMenu == "" {
		return
	}
	i := len(MenuCols) - 1
	cols := MenuCols[i]
	ui.DrawColoredText(interactMenu, cols[0], line, ColorViolet)
}
