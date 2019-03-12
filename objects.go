package main

import "fmt"

type stair int

const (
	NormalStair stair = iota
	WinStair
)

func (st stair) ShortDesc(g *game) (desc string) {
	if st == WinStair {
		desc = fmt.Sprintf("glowing stairs")
	} else {
		desc = fmt.Sprintf("stairs downwards")
	}
	return desc
}

func (st stair) Desc(g *game) (desc string) {
	if st == WinStair {
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

func (st stair) Style(g *game) (r rune, fg uicolor) {
	r = '>'
	if st == WinStair {
		fg = ColorFgMagicPlace
	} else {
		fg = ColorFgPlace
	}
	return r, fg
}

type stone int

const (
	InertStone stone = iota
	TeleStone
	FogStone
	QueenStone
	TreeStone
	ObstructionStone
)

const NumStones = int(ObstructionStone) + 1

func (stn stone) String() (text string) {
	switch stn {
	case InertStone:
		text = "inert stone"
	case TeleStone:
		text = "teleport stone"
	case FogStone:
		text = "fog stone"
	case QueenStone:
		text = "queenstone"
	case TreeStone:
		text = "tree stone"
	case ObstructionStone:
		text = "obstruction stone"
	}
	return text
}

func (stn stone) Desc(g *game) (text string) {
	switch stn {
	case InertStone:
		text = "This stone has been depleted of magical energies."
	case TeleStone:
		text = "Any creature standing on the teleport stone will teleport away when hit in combat."
	case FogStone:
		text = "Fog will appear if a creature is hurt while standing on the fog stone."
	case QueenStone:
		text = "If a creature is hurt while standing on queenstone, a loud boom will resonate, leaving nearby monsters in a 2-range distance confused. You know how to avoid the effect yourself."
	case TreeStone:
		text = "Any creature hurt while standing on a tree stone will be lignified."
	case ObstructionStone:
		text = "When a creature is hurt while standing on the obstruction stone, temporal walls appear around it."
	}
	return text
}

func (stn stone) ShortDesc(g *game) string {
	return fmt.Sprintf("%s", Indefinite(stn.String(), false))
}

func (stn stone) Style(g *game) (r rune, fg uicolor) {
	r = '_'
	if stn == InertStone {
		fg = ColorFgPlace
	} else {
		fg = ColorFgMagicPlace
	}
	return r, fg
}

func (g *game) UseStone(pos position) {
	g.StoryPrintf("You activated %s.", g.Objects.Stones[pos].ShortDesc(g))
	g.Objects.Stones[pos] = InertStone
	g.Stats.UsedStones++
	g.Print("The stone becomes inert.")
}
