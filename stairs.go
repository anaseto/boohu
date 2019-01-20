package main

import "fmt"

type stair int

const (
	NormalStair stair = iota
	WinStair
)

func (strt stair) Style(g *game) (r rune, fg uicolor) {
	r = '>'
	if strt == WinStair {
		fg = ColorFgMagicPlace
	} else {
		fg = ColorFgPlace
	}
	return r, fg
}

func (strt stair) ShortDesc(g *game) (desc string) {
	if strt == WinStair {
		desc = fmt.Sprintf("glowing stairs")
	} else {
		desc = fmt.Sprintf("stairs downwards")
	}
	return desc
}

func (strt stair) Desc(g *game) (desc string) {
	if strt == WinStair {
		desc = "These shiny-looking stairs are in fact a magical monolith. It is said they were made some centuries ago by Marevor Helith. They will lead you back to your village."
		if g.Depth < MaxDepth {
			desc += " Note that this is not the last floor, so you may want to find a normal stair and continue collecting simellas, if you're courageous enough."
		}
	} else {
		desc = "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
		if g.Depth == WinDepth {
			desc += " If you're afraid, you could instead just win by taking the magical stairs somewhere in the same map."
		}
	}
	return desc
}
