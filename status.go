// confusion idea from: https://crawl.develz.org/tavern/viewtopic.php?f=17&t=24108&sid=cb465fe78aba3b9074a32efc2a835d80#p318813

package main

type status int

const (
	StatusBerserk status = iota
	StatusSlow
	StatusExhausted
	StatusHaste
	StatusEvasion
	StatusLignification
	StatusConfusion
	StatusTele
	StatusResistance
)

func (st status) Good() bool {
	switch st {
	case StatusBerserk, StatusHaste, StatusEvasion:
		return true
	default:
		return false
	}
}

func (st status) Bad() bool {
	switch st {
	case StatusSlow, StatusConfusion:
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
	case StatusHaste:
		return "Swift"
	case StatusLignification:
		return "Lignified"
	case StatusEvasion:
		return "Agile"
	case StatusConfusion:
		return "Confused"
	case StatusTele:
		return "Tele"
	default:
		// should not happen
		return "unknown"
	}
}
