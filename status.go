package main

type status int

const (
	StatusBerserk status = iota
	StatusSlow
	StatusExhausted
	StatusSwift
	StatusAgile
	StatusLignification
	StatusConfusion
	StatusTele
	StatusNausea
	StatusDisabledShield
	StatusCorrosion
	StatusFlames // fake status
	StatusDig
	StatusSwap
	StatusShadows
	StatusSlay
)

func (st status) Good() bool {
	switch st {
	case StatusBerserk, StatusSwift, StatusAgile, StatusDig, StatusSwap, StatusShadows, StatusSlay:
		return true
	default:
		return false
	}
}

func (st status) Bad() bool {
	switch st {
	case StatusSlow, StatusConfusion, StatusNausea, StatusDisabledShield, StatusFlames, StatusCorrosion:
		return true
	default:
		return false
	}
}

func (st status) String() string {
	switch st {
	case StatusBerserk:
		return "Berserk"
	case StatusSlow:
		return "Slow"
	case StatusExhausted:
		return "Exhausted"
	case StatusSwift:
		return "Swift"
	case StatusLignification:
		return "Lignified"
	case StatusAgile:
		return "Agile"
	case StatusConfusion:
		return "Confused"
	case StatusTele:
		return "Tele"
	case StatusNausea:
		return "Nausea"
	case StatusDisabledShield:
		return "-Shield"
	case StatusCorrosion:
		return "Corroded"
	case StatusFlames:
		return "Flames"
	case StatusDig:
		return "Dig"
	case StatusSwap:
		return "Swap"
	case StatusShadows:
		return "Shadows"
	case StatusSlay:
		return "Slay"
	default:
		// should not happen
		return "unknown"
	}
}

func (st status) Short() string {
	switch st {
	case StatusBerserk:
		return "Be"
	case StatusSlow:
		return "Sl"
	case StatusExhausted:
		return "Ex"
	case StatusSwift:
		return "Sw"
	case StatusLignification:
		return "Li"
	case StatusAgile:
		return "Ag"
	case StatusConfusion:
		return "Co"
	case StatusTele:
		return "Te"
	case StatusNausea:
		return "Na"
	case StatusDisabledShield:
		return "-S"
	case StatusCorrosion:
		return "Co"
	case StatusFlames:
		return "Fl"
	case StatusDig:
		return "Di"
	case StatusSwap:
		return "Sw"
	case StatusShadows:
		return "Sh"
	case StatusSlay:
		return "Sl"
	default:
		// should not happen
		return "?"
	}
}
