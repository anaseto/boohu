package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"sort"
	"strings"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type termui struct {
}

// colors: http://ethanschoonover.com/solarized
var (
	ColorBgLOS             termbox.Attribute = 231
	ColorBgDark            termbox.Attribute = 235
	ColorBg                termbox.Attribute = 235
	ColorBgCloud           termbox.Attribute = 236
	ColorFgLOS             termbox.Attribute = 242
	ColorFgDark            termbox.Attribute = 241
	ColorFg                termbox.Attribute = 246
	ColorFgPlayer          termbox.Attribute = 34
	ColorFgMonster         termbox.Attribute = 161
	ColorFgSleepingMonster termbox.Attribute = 62
	ColorFgConfusedMonster termbox.Attribute = 167
	ColorFgCollectable     termbox.Attribute = 137
	ColorFgStairs          termbox.Attribute = 126
	ColorFgGold            termbox.Attribute = 137
	ColorFgHPok            termbox.Attribute = 65
	ColorFgHPwounded       termbox.Attribute = 137
	ColorFgHPcritical      termbox.Attribute = 161
	ColorFgMPok            termbox.Attribute = 34
	ColorFgMPpartial       termbox.Attribute = 62
	ColorFgMPcritical      termbox.Attribute = 126
	ColorFgStatusGood      termbox.Attribute = 34
	ColorFgStatusBad       termbox.Attribute = 161
	ColorFgStatusOther     termbox.Attribute = 137
	ColorFgExcluded        termbox.Attribute = 161
	ColorFgTargetMode      termbox.Attribute = 38
)

func SolarizedPalette() {
	ColorBgLOS = 16
	ColorBgDark = 0
	ColorBg = 0
	ColorBgCloud = 8
	ColorFgLOS = 12
	ColorFgDark = 11
	ColorFg = 13
	ColorFgPlayer = 5
	ColorFgMonster = 2
	ColorFgSleepingMonster = 14
	ColorFgConfusedMonster = 10
	ColorFgCollectable = 4
	ColorFgStairs = 6
	ColorFgGold = 4
	ColorFgHPok = 3
	ColorFgHPwounded = 4
	ColorFgHPcritical = 2
	ColorFgMPok = 5
	ColorFgMPpartial = 14
	ColorFgMPcritical = 6
	ColorFgStatusGood = 5
	ColorFgStatusBad = 2
	ColorFgStatusOther = 4
	ColorFgTargetMode = 7
}

func main() {
	opt := flag.Bool("s", false, "Use true 16-color solarized palette")
	flag.Parse()
	if *opt {
		SolarizedPalette()
	}

	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	termbox.SetOutputMode(termbox.Output256)
	if err != nil {
		log.Println(err)
	}

	tui := &termui{}
	tui.DrawWelcome()
	g := &game{}
	load, err := g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Print("Error loading saved game… starting new game.")
	}
	g.ui = tui
	g.EventLoop()
}

func (ui *termui) DrawWelcome() {
	termbox.Clear(ColorFg, ColorBg)
	col := 10
	line := 5
	rcol := col + 20
	ColorText := ColorFgHPok
	ui.DrawDark("────│\\/\\/\\/\\/\\/\\/\\/│────", col, line, ColorText)
	line++
	ui.DrawDark("##", col, line, ColorFgDark)
	ui.DrawLight("##", col+2, line, ColorFgLOS)
	ui.DrawDark("│              │", col+4, line, ColorText)
	ui.DrawDark("####", rcol, line, ColorFgDark)
	line++
	ui.DrawDark("#.", col, line, ColorFgDark)
	ui.DrawLight("..", col+2, line, ColorFgLOS)
	ui.DrawDark("│              │", col+4, line, ColorText)
	ui.DrawDark("...#", rcol, line, ColorFgDark)
	line++
	ui.DrawDark("##", col, line, ColorFgDark)
	ui.DrawLight("!", col+2, line, ColorFgCollectable)
	ui.DrawLight(".", col+3, line, ColorFgLOS)
	ui.DrawDark("│              │", col+4, line, ColorText)
	ui.DrawDark("│  BREAK       │", col+4, line, ColorText)
	ui.DrawDark(".###", rcol, line, ColorFgDark)
	line++
	ui.DrawDark(" #", col, line, ColorFgDark)
	ui.DrawLight("gG", col+2, line, ColorFgMonster)
	ui.DrawDark("│  OUT OF      │", col+4, line, ColorText)
	ui.DrawDark("##  ", rcol, line, ColorFgDark)
	line++
	ui.DrawLight("##", col, line, ColorFgLOS)
	ui.DrawLight("Dg", col+2, line, ColorFgMonster)
	ui.DrawDark("│  HAREKA'S    │", col+4, line, ColorText)
	ui.DrawDark(".## ", rcol, line, ColorFgDark)
	line++
	ui.DrawLight("#", col, line, ColorFgLOS)
	ui.DrawLight("@", col+1, line, ColorFgPlayer)
	ui.DrawLight("#", col+2, line, ColorFgLOS)
	ui.DrawDark("#", col+3, line, ColorFgDark)
	ui.DrawDark("│  UNDERGROUND │", col+4, line, ColorText)
	ui.DrawDark("..##", rcol, line, ColorFgDark)
	line++
	ui.DrawLight("#.#", col, line, ColorFgLOS)
	ui.DrawDark("#", col+3, line, ColorFgDark)
	ui.DrawDark("│              │", col+4, line, ColorText)
	ui.DrawDark("#.", rcol, line, ColorFgDark)
	ui.DrawDark(">", rcol+2, line, ColorFgStairs)
	ui.DrawDark("#", rcol+3, line, ColorFgDark)
	line++
	ui.DrawLight("#", col, line, ColorFgLOS)
	ui.DrawLight("[", col+1, line, ColorFgCollectable)
	ui.DrawLight(".", col+2, line, ColorFgLOS)
	ui.DrawDark("##", col+3, line, ColorFgDark)
	ui.DrawDark("│              │", col+4, line, ColorFgHPok)
	ui.DrawDark("..##", rcol, line, ColorFgDark)
	line++
	ui.DrawDark("────│/\\/\\/\\/\\/\\/\\/\\│────", col, line, ColorText)
	line++
	line++
	ui.DrawDark("───Press any key to continue───", col-3, line, ColorFg)
	termbox.Flush()
	ui.PressAnyKey()
}

func (ui *termui) DrawColored(text string, x, y int, fg, bg termbox.Attribute) {
	col := 0
	for _, r := range text {
		termbox.SetCell(x+col, y, r, fg, bg)
		col++
	}
}

func (ui *termui) DrawDark(text string, x, y int, fg termbox.Attribute) {
	col := 0
	for _, r := range text {
		termbox.SetCell(x+col, y, r, fg, ColorBgDark)
		col++
	}
}

func (ui *termui) DrawLight(text string, x, y int, fg termbox.Attribute) {
	col := 0
	for _, r := range text {
		termbox.SetCell(x+col, y, r, fg, ColorBgLOS)
		col++
	}
}

func (ui *termui) HandlePlayerTurn(g *game, ev event) bool {
getKey:
	for {
		ui.DrawDungeonView(g, false)
		var err error
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyArrowUp:
					tev.Ch = 'k'
				case termbox.KeyArrowRight:
					tev.Ch = 'l'
				case termbox.KeyArrowDown:
					tev.Ch = 'j'
				case termbox.KeyArrowLeft:
					tev.Ch = 'h'
				case termbox.KeyCtrlW:
					if ui.Wizard(g) {
						g.Wizard = true
						g.Print("You are now in wizard mode and cannot obtain winner status.")
						ui.DrawDungeonView(g, false)
						continue getKey
					}
					g.Print("Ok, then.")
					continue getKey
				case termbox.KeyCtrlQ:
					if ui.Quit(g) {
						g.RemoveSaveFile()
						return true
					}
					g.Print("Ok, then.")
					continue getKey
				case termbox.KeyCtrlP:
					tev.Ch = 'm'
				}
			}
			switch tev.Ch {
			case 'h', '4':
				err = g.MovePlayer(g.Player.Pos.W(), ev)
			case 'l', '6':
				err = g.MovePlayer(g.Player.Pos.E(), ev)
			case 'j', '2':
				err = g.MovePlayer(g.Player.Pos.S(), ev)
			case 'k', '8':
				err = g.MovePlayer(g.Player.Pos.N(), ev)
			case 'y', '7':
				err = g.MovePlayer(g.Player.Pos.NW(), ev)
			case 'b', '1':
				err = g.MovePlayer(g.Player.Pos.SW(), ev)
			case 'u', '9':
				err = g.MovePlayer(g.Player.Pos.NE(), ev)
			case 'n', '3':
				err = g.MovePlayer(g.Player.Pos.SE(), ev)
			case '.', '5':
				g.WaitTurn(ev)
			case 'r':
				err = g.Rest(ev)
			case '>':
				if g.Stairs[g.Player.Pos] {
					if g.Descend(ev) {
						ui.Win(g)
						return true
					}
					ui.DrawDungeonView(g, false)
				} else {
					err = errors.New("No stairs here.")
				}
			case 'e', 'g', ',':
				err = ui.Equip(g, ev)
			case 'q', 'a':
				err = ui.SelectPotion(g, ev)
			case 't', 'f':
				err = ui.SelectProjectile(g, ev)
			case 'v', 'z':
				err = ui.SelectRod(g, ev)
			case 'o':
				err = g.Autoexplore(ev)
			case 'x':
				b := ui.Examine(g)
				ui.DrawDungeonView(g, false)
				if !b {
					continue getKey
				} else if !g.MoveToTarget(ev) {
					continue getKey
				}
			case '?':
				ui.KeysHelp(g)
				continue getKey
			case '%', 'C':
				ui.CharacterInfo(g)
				continue getKey
			case 'm':
				ui.DrawPreviousLogs(g)
				continue getKey
			case 'S':
				ev.Renew(g, 0)
				g.Save()
				return true
			case '#':
				err := g.WriteDump()
				if err != nil {
					g.Print("Error writting dump to file.")
				} else {
					dataDir, _ := g.DataDir()
					g.Printf("Dump written to %s.", filepath.Join(dataDir, "dump"))
				}
				continue getKey
			default:
				err = errors.New("Unknown key.")
			}
			if err != nil {
				g.Print(err.Error())
				continue getKey
			}
			return false
		}
	}
}

func (ui *termui) DrawKeysDescription(g *game, actions []string) {
	termbox.Clear(ColorFg, ColorBg)
	help := &bytes.Buffer{}
	help.WriteString("┌────────────── Keys ────────────────────────────────────────────────────────\n")
	help.WriteString("│\n")
	for i := 0; i < len(actions)-1; i += 2 {
		fmt.Fprintf(help, "│ %s: %s\n", actions[i], actions[i+1])
	}
	help.WriteString("│\n")
	help.WriteString("└──── press esc or space to return to the game ──────────────────────────────\n")
	ui.DrawText(help.String(), 0, 0)
	termbox.Flush()
	ui.WaitForContinue(g)
}

func (ui *termui) KeysHelp(g *game) {
	ui.DrawKeysDescription(g, []string{
		"Movement", "h/j/k/l/y/u/b/n or numpad",
		"Rest", "r",
		"Wait", "“.” or 5",
		"Use stairs", ">",
		"Quaff potion", "q or a",
		"Equip weapon/armour/...", "e or g",
		"Autoexplore", "o",
		"Examine", "x (? for help)",
		"Throw item", "t or f (? for help)",
		"Evoke rod", "v or z (? for help)",
		"View Character Information", "%% or C",
		"View previous messages", "m",
		"Write character dump to file", "#",
		"Save and Quit", "S",
		"Quit without saving", "Ctrl-Q",
	})
}

func (ui *termui) ExamineHelp(g *game) {
	ui.DrawKeysDescription(g, []string{
		"Move cursor", "h/j/k/l/y/u/b/n or numpad",
		"Cycle through monsters", "+",
		"Cycle through stairs", ">",
		"Cycle through objects", "o",
		"Go to/select target", "“.” or enter",
		"View target description", "v or d",
		"Toggle exclude area from automatic travelling", "e",
	})
}

func (ui *termui) Equip(g *game, ev event) error {
	return g.Equip(ev)
}

func (ui *termui) CharacterInfo(g *game) {
	termbox.Clear(ColorFg, ColorBg)
	b := bytes.Buffer{}
	b.WriteString(formatText(fmt.Sprintf("You are wielding a %s. %s", g.Player.Weapon, g.Player.Weapon.Desc()), 79))
	b.WriteString("\n\n")
	b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Armour, g.Player.Armour.Desc()), 79))
	b.WriteString("\n\n")
	if g.Player.Shield != NoShield {
		b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Shield, g.Player.Shield.Desc()), 79))
		b.WriteString("\n\n")
	}
	b.WriteString(ui.AptitudesText(g))
	ui.DrawText(b.String(), 0, 0)
	termbox.Flush()
	ui.WaitForContinue(g)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) AptitudesText(g *game) string {
	apts := []string{}
	for apt, b := range g.Player.Aptitudes {
		if b {
			apts = append(apts, apt.String())
		}
	}
	sort.Strings(apts)
	if len(apts) > 0 {
		return "Aptitudes:\n" + strings.Join(apts, "\n")
	}
	return "You do not have any special aptitudes."
}

func (ui *termui) DescribePosition(g *game, pos position, targ Targetter) {
	mons, _ := g.MonsterAt(pos)
	c, okCollectable := g.Collectables[pos]
	eq, okEq := g.Equipables[pos]
	rod, okRod := g.Rods[pos]
	var desc string
	if pos == g.Player.Pos {
		desc = "This is you. "
	}
	switch {
	case !g.Dungeon.Cell(pos).Explored:
		desc = "You do not know what is in there."
	case !targ.Reachable(g, pos):
		desc = "This is out of reach."
	case mons.Exists() && g.Player.LOS[pos]:
		desc += fmt.Sprintf("You see a %s (%s).", mons.Kind, ui.MonsterInfo(mons))
	case g.Gold[pos] > 0:
		desc += fmt.Sprintf("You see some gold (%d).", g.Gold[pos])
	case okCollectable && c != nil:
		if c.Quantity > 1 {
			desc += fmt.Sprintf("You see %d %ss there.", c.Quantity, c.Consumable)
		} else {
			desc += fmt.Sprintf("You see a %s there.", c.Consumable)
		}
	case okEq:
		desc += fmt.Sprintf("You see a %v.", eq)
	case okRod:
		desc += fmt.Sprintf("You see a %v.", rod)
	case g.Stairs[pos]:
		desc += "You see stairs downwards."
	case g.Dungeon.Cell(pos).T == WallCell:
		desc += "You see a wall."
	default:
		desc += "You see the ground."
	}
	g.Print(desc)
}

func (ui *termui) Examine(g *game) bool {
	ex := &examiner{}
	err := ui.CursorAction(g, ex)
	if err != nil {
		g.Print(err.Error())
		return false
	}
	return ex.done
}

func (ui *termui) ChooseTarget(g *game, targ Targetter) bool {
	err := ui.CursorAction(g, targ)
	if err != nil {
		g.Print(err.Error())
		return false
	}
	return targ.Done()
}

func (ui *termui) CursorAction(g *game, targ Targetter) error {
	pos := g.Player.Pos
	minDist := 999
	for _, mons := range g.Monsters {
		if mons.Exists() && g.Player.LOS[mons.Pos] {
			dist := mons.Pos.Distance(g.Player.Pos)
			if minDist > dist {
				minDist = dist
				pos = mons.Pos
			}
		}
	}
	var err error
	var nstatic position
	nmonster := 0
	objects := []position{}
	nobject := 0
	opos := position{-1, -1}
loop:
	for {
		err = nil
		if pos != opos {
			ui.DescribePosition(g, pos, targ)
		}
		opos = pos
		targ.ComputeHighlight(g, pos)
		termbox.SetCursor(pos.X, pos.Y)
		ui.DrawDungeonView(g, true)
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			npos := pos
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyArrowUp:
					tev.Ch = 'k'
				case termbox.KeyArrowRight:
					tev.Ch = 'l'
				case termbox.KeyArrowDown:
					tev.Ch = 'j'
				case termbox.KeyArrowLeft:
					tev.Ch = 'h'
				case termbox.KeyEsc:
					break loop
				case termbox.KeyEnter:
					tev.Ch = '.'
				}
			}
			switch tev.Ch {
			case 'h', '4':
				npos = pos.W()
			case 'l', '6':
				npos = pos.E()
			case 'j', '2':
				npos = pos.S()
			case 'k', '8':
				npos = pos.N()
			case 'y', '7':
				npos = pos.NW()
			case 'b', '1':
				npos = pos.SW()
			case 'u', '9':
				npos = pos.NE()
			case 'n', '3':
				npos = pos.SE()
			case '>':
			search:
				for i := 0; i < g.Dungeon.Width*g.Dungeon.Heigth; i++ {
					for nstatic.X < g.Dungeon.Width-1 {
						nstatic.X++
						if g.Stairs[nstatic] && g.Dungeon.Cell(nstatic).Explored {
							npos = nstatic
							break search
						}
					}
					nstatic.Y++
					nstatic.X = 0
					if nstatic.Y >= g.Dungeon.Heigth {
						nstatic.Y = 0
					}
				}
			case '+', '-':
				for i := 0; i < len(g.Monsters); i++ {
					if tev.Ch == '+' {
						nmonster++
					} else {
						nmonster--
					}
					if nmonster > len(g.Monsters)-1 {
						nmonster = 0
					} else if nmonster < 0 {
						nmonster = len(g.Monsters) - 1
					}
					mons := g.Monsters[nmonster]
					if mons.Exists() && g.Player.LOS[mons.Pos] && pos != mons.Pos {
						npos = mons.Pos
						break
					}
				}
			case 'o':
				if len(objects) == 0 {
					for p := range g.Collectables {
						objects = append(objects, p)
					}
					for p := range g.Rods {
						objects = append(objects, p)
					}
					for p := range g.Equipables {
						objects = append(objects, p)
					}
					for p := range g.Gold {
						objects = append(objects, p)
					}
				}
				for i := 0; i < len(objects); i++ {
					nobject++
					if nobject > len(objects)-1 {
						nobject = 0
					}
					p := objects[nobject]
					if g.Dungeon.Cell(p).Explored {
						npos = p
						break
					}
				}
			case 'v', 'd':
				termbox.HideCursor()
				ui.ViewPositionDescription(g, pos)
				termbox.SetCursor(pos.X, pos.Y)
			case '?':
				termbox.HideCursor()
				ui.ExamineHelp(g)
				termbox.SetCursor(pos.X, pos.Y)
			case '.':
				err = targ.Action(g, pos)
				if err != nil {
					g.Print(err.Error())
				} else {
					break loop
				}
			case 'e':
				if !g.Dungeon.Cell(pos).Explored {
					g.Print("You cannot choose an unexplored cell for exclusion.")
				} else {
					toggle := !g.ExclusionsMap[pos]
					g.ComputeExclusion(pos, toggle)
				}
			default:
				g.Print("Invalid key.")
			}
			if g.Dungeon.Valid(npos) {
				pos = npos
			}
		}
	}
	g.Highlight = nil
	termbox.HideCursor()
	return err
}

func (ui *termui) ViewPositionDescription(g *game, pos position) {
	mons, _ := g.MonsterAt(pos)
	if mons.Exists() {
		termbox.HideCursor()
		ui.DrawMonsterDescription(g, mons)
		termbox.SetCursor(pos.X, pos.Y)
	} else if c, ok := g.Collectables[pos]; ok {
		ui.DrawDescription(g, c.Consumable.Desc())
	} else if r, ok := g.Rods[pos]; ok {
		ui.DrawDescription(g, r.Desc())
	} else if eq, ok := g.Equipables[pos]; ok {
		ui.DrawDescription(g, eq.Desc())
	} else if g.Stairs[pos] {
		ui.DrawDescription(g, "Stairs lead to the next level of the Underground. There's no way back.")
	} else {
		g.Print("Nothing worth of description here.")
	}

}

func (ui *termui) MonsterInfo(m *monster) string {
	infos := []string{}
	infos = append(infos, m.State.String())
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, st.String())
		}
	}
	p := (m.HP * 100) / m.HPmax
	health := fmt.Sprintf("%d %% HP", p)
	infos = append(infos, health)
	return strings.Join(infos, ", ")
}

func (ui *termui) DrawDungeonView(g *game, targetting bool) {
	err := termbox.Clear(ColorFg, ColorBg)
	if err != nil {
		log.Println(err)
	}
	m := g.Dungeon
	for i := 0; i < g.Dungeon.Width; i++ {
		termbox.SetCell(i, g.Dungeon.Heigth, '─', ColorFg, ColorBg)
	}
	for i := 0; i < g.Dungeon.Heigth; i++ {
		termbox.SetCell(g.Dungeon.Width, i, '│', ColorFg, ColorBg)
	}
	termbox.SetCell(g.Dungeon.Width, g.Dungeon.Heigth, '┘', ColorFg, ColorBg)
	for i := range m.Cells {
		pos := m.CellPosition(i)
		ui.DrawPosition(g, pos)
	}
	ui.DrawText(fmt.Sprintf("[ %v (%d)", g.Player.Armour, g.Player.Armor()), 81, 0)
	ui.DrawText(fmt.Sprintf(") %v (%d)", g.Player.Weapon, g.Player.Attack()), 81, 1)
	if g.Player.Shield != NoShield {
		if g.Player.Weapon.TwoHanded() {
			ui.DrawText(fmt.Sprintf("] %v (unusable)", g.Player.Shield), 81, 2)
		} else {
			ui.DrawText(fmt.Sprintf("] %v (%d)", g.Player.Shield, g.Player.Shield.Block()), 81, 2)
		}
	}
	if targetting {
		ui.DrawColoredText("Targetting", 81, 20, ColorFgTargetMode)
	}
	ui.DrawStatusLine(g)
	ui.DrawLog(g)
	termbox.Flush()
}

func (ui *termui) DrawPosition(g *game, pos position) {
	m := g.Dungeon
	c := m.Cell(pos)
	if !c.Explored && !g.Wizard {
		return
	}
	if g.Wizard {
		if c.T == WallCell {
			if len(g.Dungeon.FreeNeighbors(pos)) == 0 {
				return
			}
		}
	}
	var fgColor termbox.Attribute
	var bgColor termbox.Attribute
	if g.Player.LOS[pos] {
		fgColor = ColorFgLOS
		bgColor = ColorBgLOS
		if _, ok := g.Clouds[pos]; ok {
			bgColor = ColorBgCloud
		}
		if g.Highlight[pos] {
			bgColor = ColorBgLOS | termbox.AttrReverse
			//fgColor = ColorFgRay
			//bgColor = ColorBgRay
		}
	} else {
		fgColor = ColorFgDark
		bgColor = ColorBgDark
	}
	if g.ExclusionsMap[pos] {
		fgColor = ColorFgExcluded
	}
	var r rune
	switch c.T {
	case WallCell:
		r = '#'
	case FreeCell:
		if g.UnknownDig[pos] {
			r = '#'
			break
		}
		switch {
		case pos == g.Player.Pos:
			r = '@'
			fgColor = ColorFgPlayer
		default:
			r = '.'
			if _, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
				r = '§'
			}
			if c, ok := g.Collectables[pos]; ok {
				r = c.Consumable.Letter()
				fgColor = ColorFgCollectable
			} else if eq, ok := g.Equipables[pos]; ok {
				r = eq.Letter()
				fgColor = ColorFgCollectable
			} else if rod, ok := g.Rods[pos]; ok {
				r = rod.Letter()
				fgColor = ColorFgCollectable
			} else if _, ok := g.Stairs[pos]; ok {
				r = '>'
				fgColor = ColorFgStairs
			} else if _, ok := g.Gold[pos]; ok {
				r = '$'
				fgColor = ColorFgGold
			}
			m, _ := g.MonsterAt(pos)
			if m.Exists() && (g.Player.LOS[m.Pos] || g.Wizard) {
				r = m.Kind.Letter()
				if m.Status(MonsConfused) {
					fgColor = ColorFgConfusedMonster
				} else if m.State == Resting {
					fgColor = ColorFgSleepingMonster
				} else {
					fgColor = ColorFgMonster
				}
			}
		}
	}
	termbox.SetCell(pos.X, pos.Y, r, fgColor, bgColor)
}

func (ui *termui) DrawStatusLine(g *game) {
	sts := statusSlice{}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)
	hpColor := termbox.Attribute(ColorFgHPok)
	switch {
	case g.Player.HP*100/g.Player.HPMax() < 30:
		hpColor = ColorFgHPcritical
	case g.Player.HP*100/g.Player.HPMax() < 70:
		hpColor = ColorFgHPwounded
	}
	mpColor := termbox.Attribute(ColorFgMPok)
	switch {
	case g.Player.MP*100/g.Player.MPMax() < 30:
		mpColor = ColorFgMPcritical
	case g.Player.MP*100/g.Player.MPMax() < 70:
		mpColor = ColorFgMPpartial
	}
	ui.DrawColoredText(fmt.Sprintf("HP: %d", g.Player.HP), 81, 4, hpColor)
	ui.DrawColoredText(fmt.Sprintf("MP: %d", g.Player.MP), 81, 5, mpColor)
	ui.DrawText(fmt.Sprintf("Gold: %d", g.Player.Gold), 81, 7)
	ui.DrawText(fmt.Sprintf("Depth: %d", g.Depth), 81, 8)
	ui.DrawText(fmt.Sprintf("Turns: %.1f", float64(g.Turn)/10), 81, 9)

	for i, st := range sts {
		var color termbox.Attribute
		if st.Good() {
			color = ColorFgStatusGood
		} else if st.Bad() {
			color = ColorFgStatusBad
		} else {
			color = ColorFgStatusOther
		}
		ui.DrawColoredText(st.String(), 81, 10+i, color)
	}
}

func (ui *termui) DrawLog(g *game) {
	min := len(g.Log) - 4
	if min < 0 {
		min = 0
	}
	for i, s := range g.Log[min:] {
		ui.DrawText(s, 0, g.Dungeon.Heigth+1+i)
	}
}

func (ui *termui) DrawPreviousLogs(g *game) {
	lines := 23
	nmax := len(g.Log) - lines
	n := nmax
loop:
	for {
		termbox.Clear(ColorFg, ColorBg)
		if n >= nmax {
			n = nmax
		}
		if n < 0 {
			n = 0
		}
		to := n + lines
		if to >= len(g.Log) {
			to = len(g.Log)
		}
		for i := n; i < to; i++ {
			ui.DrawText(g.Log[i], 0, i-n)
		}
		s := fmt.Sprintf("─────────(%d/%d)───────────────────────────────────────────────────────────────\n", len(g.Log)-to, len(g.Log))
		ui.DrawText(s, 0, to-n)
		ui.DrawText("Keys: half-page up (u), half-page down (d), quit (esc or space)", 0, to+1-n)
		termbox.Flush()
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					break loop
				}
			}
			switch tev.Ch {
			case 'u':
				n -= 12
			case 'd':
				n += 12
			case 'j':
				n++
			case 'k':
				n--
			}
		}
	}
}

func (ui *termui) DrawMonsterDescription(g *game, mons *monster) {
	s := mons.Kind.Desc()
	s += " " + fmt.Sprintf("They can hit for up to %d damage.", MonsData[mons.Kind].baseAttack)
	s += " " + fmt.Sprintf("They have around %d HP.", MonsData[mons.Kind].maxHP)
	ui.DrawDescription(g, s)
}

func (ui *termui) DrawConsumableDescription(g *game, c consumable) {
	ui.DrawDescription(g, c.Desc())
}

func (ui *termui) DrawDescription(g *game, desc string) {
	termbox.Clear(ColorFg, ColorBg)
	ui.DrawText(formatText(desc, 79), 0, 0)
	termbox.Flush()
	ui.WaitForContinue(g)
}

func (ui *termui) DrawText(text string, x, y int) {
	ui.DrawColoredText(text, x, y, ColorFg)
}

func (ui *termui) DrawColoredText(text string, x, y int, color termbox.Attribute) {
	col := 0
	for _, r := range text {
		if r == '\n' {
			y++
			col = 0
			continue
		}
		termbox.SetCell(x+col, y, r, color, ColorBg)
		col++
	}
}

type rodSlice []rod

func (rs rodSlice) Len() int           { return len(rs) }
func (rs rodSlice) Swap(i, j int)      { rs[i], rs[j] = rs[j], rs[i] }
func (rs rodSlice) Less(i, j int) bool { return int(rs[i]) < int(rs[j]) }

type consumableSlice []consumable

func (cs consumableSlice) Len() int           { return len(cs) }
func (cs consumableSlice) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs consumableSlice) Less(i, j int) bool { return cs[i].Int() < cs[j].Int() }

type statusSlice []status

func (sts statusSlice) Len() int           { return len(sts) }
func (sts statusSlice) Swap(i, j int)      { sts[i], sts[j] = sts[j], sts[i] }
func (sts statusSlice) Less(i, j int) bool { return sts[i] < sts[j] }

func (ui *termui) SelectProjectile(g *game, ev event) error {
	desc := false
	for {
		termbox.Clear(ColorFg, ColorBg)
		cs := g.SortedProjectiles()
		if desc {
			ui.DrawText("Describe which projectile? (press ? for throwing menu)", 0, 0)
		} else {
			ui.DrawText("Throw which projectile? (press ? for description menu)", 0, 0)
		}
		for i, c := range cs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1)
		}
		termbox.Flush()
		index, alternate, noAction := ui.Select(g, ev, len(cs))
		if alternate {
			desc = !desc
			continue
		}
		if noAction == nil {
			if desc {
				ui.DrawDescription(g, cs[index].Desc())
				continue
			}
			b := ui.ChooseTarget(g, &chooser{single: true})
			if b {
				noAction = cs[index].Use(g, ev)
			} else {
				noAction = errors.New("Ok, then.")
			}
		}
		return noAction
	}
}

func (ui *termui) SelectPotion(g *game, ev event) error {
	desc := false
	for {
		termbox.Clear(ColorFg, ColorBg)
		cs := g.SortedPotions()
		if desc {
			ui.DrawText("Describe which potion? (press ? for quaff menu)", 0, 0)
		} else {
			ui.DrawText("Drink which potion? (press ? for describe menu)", 0, 0)
		}
		for i, c := range cs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1)
		}
		termbox.Flush()
		index, alternate, noAction := ui.Select(g, ev, len(cs))
		if alternate {
			desc = !desc
			continue
		}
		if noAction == nil {
			if desc {
				ui.DrawDescription(g, cs[index].Desc())
				continue
			}
			noAction = cs[index].Use(g, ev)
		}
		return noAction
	}
}

func (ui *termui) SelectRod(g *game, ev event) error {
	desc := false
	for {
		termbox.Clear(ColorFg, ColorBg)
		rs := g.SortedRods()
		if desc {
			ui.DrawText("Describe which rod? (press ? for evocation menu)", 0, 0)
		} else {
			ui.DrawText("Evoke which rod? (press ? for description menu)", 0, 0)
		}
		for i, c := range rs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d/%d charges, %d mana cost)",
				rune(i+97), c, g.Player.Rods[c].Charge, c.MaxCharge(), c.MPCost()), 0, i+1)
		}
		termbox.Flush()
		index, alternate, noAction := ui.Select(g, ev, len(rs))
		if alternate {
			desc = !desc
			continue
		}
		if noAction == nil {
			if desc {
				ui.DrawDescription(g, rs[index].Desc())
				continue
			}
			noAction = rs[index].Use(g, ev)
		}
		ui.DrawDungeonView(g, false)
		return noAction
	}
}

func (ui *termui) Select(g *game, ev event, l int) (index int, alternate bool, err error) {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc:
					return -1, false, errors.New("Ok, then.")
				}
			}
			if 97 <= tev.Ch && int(tev.Ch) < 97+l {
				return int(tev.Ch - 97), false, nil
			}
			if tev.Ch == '?' {
				return -1, true, nil
			}
		}
	}
}

func (ui *termui) ExploreStep(g *game) {
	time.Sleep(10 * time.Millisecond)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) Death(g *game) {
	g.Print("You die... -- press esc or space to continue --")
	ui.DrawDungeonView(g, false)
	ui.WaitForContinue(g)
	ui.Dump(g)
	g.WriteDump()
	ui.WaitForContinue(g)
}

func (ui *termui) Win(g *game) {
	if g.Wizard {
		g.Print("You escape by the magic stairs! **WIZARD** --press esc or space to continue--")
	} else {
		g.Print("You escape by the magic stairs! You win. --press esc or space to continue--")
	}
	ui.DrawDungeonView(g, false)
	ui.WaitForContinue(g)
	ui.Dump(g)
	g.WriteDump()
	ui.WaitForContinue(g)
}

func (ui *termui) Dump(g *game) {
	termbox.Clear(ColorFg, ColorBg)
	ui.DrawText(g.SimplifedDump(), 0, 0)
	termbox.Flush()
}

func (ui *termui) WaitForContinue(g *game) {
loop:
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 0 {
				switch tev.Key {
				case termbox.KeyEsc, termbox.KeySpace:
					break loop
				}
			}
		}
	}
}

func (ui *termui) Quit(g *game) bool {
	g.Print("Do you really want to quit without saving? [Y/n]")
	ui.DrawDungeonView(g, false)
	return ui.PromptConfirmation(g)
}

func (ui *termui) Wizard(g *game) bool {
	g.Print("Do you really want to enter wizard mode (no return)? [Y/n]")
	ui.DrawDungeonView(g, false)
	return ui.PromptConfirmation(g)
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			if tev.Ch == 'Y' {
				return true
			}
		}
		return false
	}
}

func (ui *termui) PressAnyKey() {
	for {
		switch tev := termbox.PollEvent(); tev.Type {
		case termbox.EventKey:
			return
		}
	}
}
