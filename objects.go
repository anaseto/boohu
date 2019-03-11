package main

import "fmt"

type object interface {
	Desc(g *game) string
	ShortDesc(g *game) string
	Style(g *game) (rune, uicolor)
}

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

type simella int

func (s simella) Style(g *game) (r rune, fg uicolor) {
	r = '♣'
	fg = ColorFgSimellas
	return r, fg
}

func (s simella) ShortDesc(g *game) (desc string) {
	desc = fmt.Sprintf("some simellas (%d)", s)
	return desc
}

func (s simella) Desc(g *game) (desc string) {
	desc = "A simella is a plant with big white flowers which are used in the Underground for their medicinal properties. They can also make tasty infusions. You were actually sent here by your village to collect as many as possible of those plants."
	return desc
}
