package main

type aptitude int

const (
	AptAccurate aptitude = iota
	AptAgile
	AptFast
	AptHealthy
	AptStealthyMovement
	AptScales
	AptRegen
	AptStealthyLOS
	AptMagic
	AptStrong
	// below unimplemented
	AptVampiric
	AptReflect
)

func (ap aptitude) String() string {
	var text string
	switch ap {
	case AptAccurate:
		text = "You are unusually accurate."
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
	case AptStrong:
		text = "You are strong."
	case AptMagic:
		text = "You have big magic reserves."
	case AptVampiric:
		text = "You sometimes steal the life of those you strike."
	case AptRegen:
		text = "You regenerate quickly."
	case AptStealthyLOS:
		text = "The shadows follow you."
	case AptReflect:
		text = "You sometimes reflect damage in combat."
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
		g.Print("Hm… You already have that aptitude. " + ap.String())
		return
	}
	g.Player.Aptitudes[ap] = true
	g.Print("You feel a little different. " + ap.String())
}
