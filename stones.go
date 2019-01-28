package main

import "fmt"

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
	g.StoryPrintf("You activated %s.", g.Objects[pos].ShortDesc(g))
	g.Objects[pos] = InertStone
	g.Stats.UsedStones++
	g.Print("The stone becomes inert.")
}
