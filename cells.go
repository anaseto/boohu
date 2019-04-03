package main

type cell struct {
	T        terrain
	Explored bool
}

type terrain int

const (
	WallCell terrain = iota
	GroundCell
	DoorCell
	FungusCell
	BarrelCell
	StairCell
	StoneCell
	MagaraCell
	BananaCell
	LightCell
	TableCell
	TreeCell
	HoledWallCell
	ScrollCell
)

func (c cell) IsPassable() bool {
	switch c.T {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell:
		return false
	default:
		return true
	}
}

func (c cell) IsIlluminable() bool {
	switch c.T {
	case WallCell, BarrelCell, TableCell, TreeCell, HoledWallCell:
		return false
	}
	return true
}

func (c cell) IsDestructible() bool {
	switch c.T {
	case WallCell, BarrelCell, DoorCell, TableCell, TreeCell, HoledWallCell:
		return true
	default:
		return false
	}
}

func (c cell) IsWall() bool {
	switch c.T {
	case WallCell:
		return true
	default:
		return false
	}
}

func (c cell) Flammable() bool {
	switch c.T {
	case FungusCell, DoorCell, BarrelCell, TableCell, TreeCell:
		return true
	default:
		return false
	}
}

func (c cell) IsGround() bool {
	switch c.T {
	case GroundCell:
		return true
	default:
		return false
	}
}

func (c cell) IsNotable() bool {
	switch c.T {
	case StairCell, StoneCell, BarrelCell, MagaraCell, BananaCell, ScrollCell:
		return true
	default:
		return false
	}
}

func (c cell) ShortDesc(g *game, pos position) (desc string) {
	switch c.T {
	case WallCell:
		desc = "a wall"
	case GroundCell:
		desc = "the ground"
	case DoorCell:
		desc = "a door"
	case FungusCell:
		desc = "foliage"
	case BarrelCell:
		desc = "a barrel"
	case StoneCell:
		desc = g.Objects.Stones[pos].ShortDesc(g)
	case StairCell:
		desc = g.Objects.Stairs[pos].ShortDesc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[pos].String()
	case BananaCell:
		desc = "a banana"
	case LightCell:
		desc = "a light"
	case TableCell:
		desc = "a table"
	case TreeCell:
		desc = "a tree"
	case HoledWallCell:
		desc = "a holed wall"
	case ScrollCell:
		desc = g.Objects.Scrolls[pos].ShortDesc(g)
	}
	return desc
}

func (c cell) Desc(g *game, pos position) (desc string) {
	switch c.T {
	case WallCell:
		desc = "A wall is an impassable pile of rocks. It can be destructed by using some items."
	case GroundCell:
		desc = "This is just plain ground."
	case DoorCell:
		desc = "A closed door blocks your line of sight. Doors open automatically when you or a monster stand on them. Doors are flammable."
	case FungusCell:
		desc = "Blue dense foliage grows in the Underground. It is difficult to see through, and is flammable."
	case BarrelCell:
		desc = "A barrel. You can hide yourself inside it when no monsters see you. It is a safe place for resting and recovering."
	case StoneCell:
		desc = g.Objects.Stones[pos].Desc(g)
	case StairCell:
		desc = g.Objects.Stairs[pos].Desc(g)
	case MagaraCell:
		desc = g.Objects.Magaras[pos].Desc(g)
	case BananaCell:
		desc = "A gawalt monkey cannot enter a healthy sleep without eating one of those bananas before."
	case LightCell:
		desc = "A light illuminates surrounding cells. Monsters can spot you in illuminated cells from a greater range."
	case TableCell:
		desc = "You can hide under the table so that only adjacent monsters can see you. Most monsters cannot walk accross the table."
	case TreeCell:
		desc = "You can climb to see farther. Moreover, many monsters will not be able to attack you while you stand on a tree. The top is never illuminated."
	case HoledWallCell:
		desc = "Only very small creatures can pass there. It is difficult to see through."
	case ScrollCell:
		desc = g.Objects.Scrolls[pos].Desc(g)
	}
	return desc
}

func (c cell) Style(g *game, pos position) (r rune, fg uicolor) {
	switch c.T {
	case WallCell:
		r, fg = '#', ColorFgLOS
	case GroundCell:
		r, fg = '.', ColorFgLOS
	case DoorCell:
		r, fg = '+', ColorFgPlace
	case FungusCell:
		r, fg = '"', ColorFgLOS
	case BarrelCell:
		r, fg = '&', ColorFgObject
	case StoneCell:
		r, fg = g.Objects.Stones[pos].Style(g)
	case StairCell:
		st := g.Objects.Stairs[pos]
		r, fg = st.Style(g)
	case MagaraCell:
		r, fg = '/', ColorFgObject
	case BananaCell:
		r, fg = ')', ColorFgObject
	case LightCell:
		r, fg = '☼', ColorFgObject
	case TableCell:
		r, fg = 'π', ColorFgObject
	case TreeCell:
		r, fg = '♣', ColorFgConfusedMonster
	case HoledWallCell:
		r, fg = 'Π', ColorFgPlace
	case ScrollCell:
		r, fg = '?', ColorFgMagicPlace
	}
	return r, fg
}
