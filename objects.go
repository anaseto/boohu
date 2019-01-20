package main

type object interface {
	Desc(g *game) string
	ShortDesc(g *game) string
	Style(g *game) (rune, uicolor)
}
