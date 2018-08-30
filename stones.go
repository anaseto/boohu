package main

type stone int

const (
	InertStone stone = iota
	TeleStone
	FogStone
	QueenStone
	ObstructionStone
)

const NumStones = int(ObstructionStone) + 1

func (s stone) String() (text string) {
	switch s {
	case InertStone:
		text = "inert stone"
	case TeleStone:
		text = "teleport stone"
	case FogStone:
		text = "fog stone"
	case QueenStone:
		text = "queenstone"
	case ObstructionStone:
		text = "obstruction stone"
	}
	return text
}

func (s stone) Description() (text string) {
	switch s {
	case InertStone:
		text = "This stone has been depleted of magical energies."
	case TeleStone:
		text = "Any creature standing on the teleport stone will teleport away when hit in combat."
	case FogStone:
		text = "Fog will appear if a creature is hurt while standing on the fog stone."
	case QueenStone:
		text = "If a creature is hurt while standing on queenstone, a loud boom will resonate, leaving nearby monsters in a 2-range distance confused. You know how to avoid the effect yourself."
	case ObstructionStone:
		text = "When a creature is hurt while standing on the obstruction stone, temporal walls appear around it."
	}
	return text
}

func (g *game) UseStone(pos position) {
	g.MagicalStones[pos] = InertStone
	g.Print("The stone becomes inert.")
}
