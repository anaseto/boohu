package main

type aptitude int

const (
	AptObstruction aptitude = iota
	AptHealthy
	AptMagic
	AptConfusingGas
	AptSmoke
	AptHear
	AptTeleport
	AptLignification
)

const NumApts = int(AptLignification) + 1

func (ap aptitude) String() string {
	var text string
	switch ap {
	case AptObstruction:
		text = "The earth occasionally blows monsters away when hurt."
	case AptHealthy:
		text = "You are healthy."
	case AptHear:
		text = "You have good ears."
	case AptMagic:
		text = "You have big magic reserves."
	case AptConfusingGas:
		text = "You occasionally release some confusing gas when hurt."
	case AptSmoke:
		text = "You occasionally get energetic and emit smoke clouds when hurt."
	case AptLignification:
		text = "Nature occasionally lignifies your foes when hurt."
	case AptTeleport:
		text = "You occasionally teleport your foes when hurt."
	}
	return text
}

func (g *game) RandomApt() (aptitude, bool) {
	count := 0
	var apt aptitude
	for {
		count++
		if count > 1000 {
			break
		}
		r := RandInt(NumApts)
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
