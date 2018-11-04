package main

type aptitude int

const (
	AptSwap aptitude = iota
	AptAgile
	AptFast
	AptHealthy
	AptStealthyMovement
	AptScales
	AptStealthyLOS
	AptMagic
	AptConfusingGas
	AptSmoke
	AptHear
	AptStrong
)

func (ap aptitude) String() string {
	var text string
	switch ap {
	case AptSwap:
		text = "You occasionally get light-footed when hurt."
	case AptAgile:
		text = "You are agile."
	case AptFast:
		text = "You move fast."
	case AptHealthy:
		text = "You are healthy."
	case AptStealthyMovement:
		text = "You move stealthily."
	case AptScales:
		text = "You are covered by scales."
	case AptHear:
		text = "You have good ears."
	case AptStrong:
		text = "You are strong."
	case AptMagic:
		text = "You have big magic reserves."
	case AptStealthyLOS:
		text = "The shadows follow you. (reduced LOS)"
	case AptConfusingGas:
		text = "You occasionally release some confusing gas when hurt."
	case AptSmoke:
		text = "You occasionally get energetic and emit smoke clouds when hurt."
	}
	return text
}

func (g *game) RandomApt() (aptitude, bool) {
	// XXX use less uniform probability ?
	max := int(AptStrong)
	count := 0
	var apt aptitude
	for {
		count++
		if count > 1000 {
			break
		}
		r := RandInt(max + 1)
		apt = aptitude(r)
		if g.Player.Aptitudes[apt] {
			continue
		}
		return apt, true
	}
	return apt, false
}

func (g *game) ApplyAptitude(ap aptitude) {
	if g.Player.Aptitudes[ap] {
		// should not happen
		g.PrintStyled("Hm… You already have that aptitude. "+ap.String(), logError)
		return
	}
	g.Player.Aptitudes[ap] = true
	g.PrintStyled("You feel different. "+ap.String(), logSpecial)
}
