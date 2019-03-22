package main

type status int

const (
	StatusSlow status = iota
	StatusExhausted
	StatusSwift
	StatusLignification
	StatusConfusion
	StatusNausea
	StatusDisabledShield
	StatusFlames // fake status
	StatusHidden
	StatusUnhidden
	StatusDig
	StatusSwap
	StatusShadows
)

func (st status) Good() bool {
	switch st {
	case StatusSwift, StatusDig, StatusSwap, StatusShadows, StatusHidden:
		return true
	default:
		return false
	}
}

func (st status) Bad() bool {
	switch st {
	case StatusSlow, StatusConfusion, StatusNausea, StatusDisabledShield, StatusUnhidden:
		return true
	default:
		return false
	}
}

func (st status) String() string {
	switch st {
	case StatusSlow:
		return "Slow"
	case StatusExhausted:
		return "Exhausted"
	case StatusSwift:
		return "Swift"
	case StatusLignification:
		return "Lignified"
	case StatusConfusion:
		return "Confused"
	case StatusNausea:
		return "Nausea"
	case StatusDisabledShield:
		return "-Shield"
	case StatusFlames:
		return "Flames"
	case StatusHidden:
		return "Hidden"
	case StatusUnhidden:
		return "Unhidden"
	case StatusDig:
		return "Dig"
	case StatusSwap:
		return "Swap"
	case StatusShadows:
		return "Shadows"
	default:
		// should not happen
		return "unknown"
	}
}

func (st status) Short() string {
	switch st {
	case StatusSlow:
		return "Sl"
	case StatusExhausted:
		return "Ex"
	case StatusSwift:
		return "Sw"
	case StatusLignification:
		return "Li"
	case StatusConfusion:
		return "Co"
	case StatusNausea:
		return "Na"
	case StatusDisabledShield:
		return "-S"
	case StatusFlames:
		return "Fl"
	case StatusHidden:
		return "H+"
	case StatusUnhidden:
		return "H-"
	case StatusDig:
		return "Di"
	case StatusSwap:
		return "Sw"
	case StatusShadows:
		return "Sh"
	default:
		// should not happen
		return "?"
	}
}
