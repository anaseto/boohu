package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type uicolor int

// uicolors: http://ethanschoonover.com/solarized
var (
	ColorBgLOS              uicolor = 230
	ColorBgDark             uicolor = 234
	ColorBgBorder           uicolor = 235
	ColorBg                 uicolor = 234
	ColorBgCloud            uicolor = 235
	ColorFgLOS              uicolor = 241
	ColorFgDark             uicolor = 240
	ColorFg                 uicolor = 245
	ColorFgPlayer           uicolor = 33
	ColorFgMonster          uicolor = 160
	ColorFgSleepingMonster  uicolor = 61
	ColorFgWanderingMonster uicolor = 166
	ColorFgConfusedMonster  uicolor = 64
	ColorFgCollectable      uicolor = 136
	ColorFgStairs           uicolor = 125
	ColorFgGold             uicolor = 136
	ColorFgHPok             uicolor = 64
	ColorFgHPwounded        uicolor = 136
	ColorFgHPcritical       uicolor = 160
	ColorFgMPok             uicolor = 33
	ColorFgMPpartial        uicolor = 61
	ColorFgMPcritical       uicolor = 125
	ColorFgStatusGood       uicolor = 33
	ColorFgStatusBad        uicolor = 160
	ColorFgStatusOther      uicolor = 136
	ColorFgExcluded         uicolor = 160
	ColorFgTargetMode       uicolor = 37
	ColorFgTemporalWall     uicolor = 37
)

func SolarizedPalette() {
	ColorBgLOS = 15
	ColorBgDark = 8
	ColorBgBorder = 0
	ColorBg = 8
	ColorBgCloud = 7
	ColorFgLOS = 11
	ColorFgDark = 10
	ColorFg = 12
	ColorFgPlayer = 4
	ColorFgMonster = 1
	ColorFgSleepingMonster = 13
	ColorFgWanderingMonster = 9
	ColorFgConfusedMonster = 2
	ColorFgCollectable = 3
	ColorFgStairs = 5
	ColorFgGold = 3
	ColorFgHPok = 2
	ColorFgHPwounded = 3
	ColorFgHPcritical = 1
	ColorFgMPok = 4
	ColorFgMPpartial = 13
	ColorFgMPcritical = 5
	ColorFgStatusGood = 4
	ColorFgStatusBad = 1
	ColorFgStatusOther = 3
	ColorFgExcluded = 1
	ColorFgTargetMode = 6
	ColorFgTemporalWall = 6
}

func FixColor() {
	ColorBgLOS = ColorBgLOS + 1
	ColorBgDark = ColorBgDark + 1
	ColorBgBorder = ColorBgBorder + 1
	ColorBg = ColorBg + 1
	ColorBgCloud = ColorBgCloud + 1
	ColorFgLOS = ColorFgLOS + 1
	ColorFgDark = ColorFgDark + 1
	ColorFg = ColorFg + 1
	ColorFgPlayer = ColorFgPlayer + 1
	ColorFgMonster = ColorFgMonster + 1
	ColorFgSleepingMonster = ColorFgSleepingMonster + 1
	ColorFgWanderingMonster = ColorFgWanderingMonster + 1
	ColorFgConfusedMonster = ColorFgConfusedMonster + 1
	ColorFgCollectable = ColorFgCollectable + 1
	ColorFgStairs = ColorFgStairs + 1
	ColorFgGold = ColorFgGold + 1
	ColorFgHPok = ColorFgHPok + 1
	ColorFgHPwounded = ColorFgHPwounded + 1
	ColorFgHPcritical = ColorFgHPcritical + 1
	ColorFgMPok = ColorFgMPok + 1
	ColorFgMPpartial = ColorFgMPpartial + 1
	ColorFgMPcritical = ColorFgMPcritical + 1
	ColorFgStatusGood = ColorFgStatusGood + 1
	ColorFgStatusBad = ColorFgStatusBad + 1
	ColorFgStatusOther = ColorFgStatusOther + 1
	ColorFgExcluded = ColorFgExcluded + 1
	ColorFgTargetMode = ColorFgTargetMode + 1
	ColorFgTemporalWall = ColorFgTemporalWall + 1
}

const (
	Black uicolor = iota
	Maroon
	Green
	Olive
	Navy
	Purple
	Teal
	Silver
)

func WindowsPalette() {
	ColorBgLOS = Silver
	ColorBgDark = Black
	ColorBgBorder = Black
	ColorBg = Black
	ColorBgCloud = Silver
	ColorFgLOS = Black
	ColorFgDark = Silver
	ColorFg = Silver
	ColorFgPlayer = Navy
	ColorFgMonster = Maroon
	ColorFgSleepingMonster = Teal
	ColorFgWanderingMonster = Purple
	ColorFgConfusedMonster = Green
	ColorFgCollectable = Olive
	ColorFgStairs = Purple
	ColorFgGold = Olive
	ColorFgHPok = Green
	ColorFgHPwounded = Olive
	ColorFgHPcritical = Maroon
	ColorFgMPok = Navy
	ColorFgMPpartial = Purple
	ColorFgMPcritical = Maroon
	ColorFgStatusGood = Navy
	ColorFgStatusBad = Maroon
	ColorFgStatusOther = Olive
	ColorFgTargetMode = Teal
	ColorFgTemporalWall = Teal
}

func (ui *termui) DrawWelcome() {
	ui.Clear()
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
	ui.Flush()
	ui.PressAnyKey()
}

func (ui *termui) DrawColored(text string, x, y int, fg, bg uicolor) {
	col := 0
	for _, r := range text {
		ui.SetCell(x+col, y, r, fg, bg)
		col++
	}
}

func (ui *termui) DrawDark(text string, x, y int, fg uicolor) {
	col := 0
	for _, r := range text {
		ui.SetCell(x+col, y, r, fg, ColorBgDark)
		col++
	}
}

func (ui *termui) DrawLight(text string, x, y int, fg uicolor) {
	col := 0
	for _, r := range text {
		ui.SetCell(x+col, y, r, fg, ColorBgLOS)
		col++
	}
}

func (ui *termui) EnterWizard(g *game) {
	if ui.Wizard(g) {
		g.WizardMode()
		ui.DrawDungeonView(g, false)
	} else {
		g.Print("Ok, then.")
	}
}

func (ui *termui) HandleCharacter(g *game, ev event, c rune) (err error, again bool, quit bool) {
	switch c {
	case 'h', '4', 'l', '6', 'j', '2', 'k', '8',
		'y', '7', 'b', '1', 'u', '9', 'n', '3':
		err = g.MovePlayer(g.Player.Pos.To(KeyToDir(c)), ev)
	case 'H', 'L', 'J', 'K', 'Y', 'B', 'U', 'N':
		err = g.GoToDir(KeyToDir(c), ev)
	case '.', '5':
		g.WaitTurn(ev)
	case 'r':
		err = g.Rest(ev)
	case '>', 'D':
		if g.Stairs[g.Player.Pos] {
			if g.Descend(ev) {
				ui.Win(g)
				quit = true
				return err, again, quit
			}
			ui.DrawDungeonView(g, false)
		} else {
			err = errors.New("No stairs here.")
		}
	case 'G':
		stairs := g.StairsSlice()
		sortedStairs := g.SortedNearestTo(stairs, g.Player.Pos)
		if len(sortedStairs) > 0 {
			g.AutoTarget = &sortedStairs[0]
			if !g.MoveToTarget(ev) {
				err = errors.New("Cannot travel to stairs now.")
			}
		} else {
			err = errors.New("You cannot go to any stairs.")
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
		b := ui.Examine(g, nil)
		ui.DrawDungeonView(g, false)
		if !b {
			again = true
		} else if !g.MoveToTarget(ev) {
			again = true
		}
	case '?':
		ui.KeysHelp(g)
		again = true
	case '%', 'C':
		ui.CharacterInfo(g)
		again = true
	case 'm':
		ui.DrawPreviousLogs(g)
		again = true
	case 'S':
		ev.Renew(g, 0)
		err := g.Save()
		if err != nil {
			g.PrintStyled("Could not save game. --press any key to continue--", logError)
			ui.DrawDungeonView(g, false)
			ui.PressAnyKey()
		}
		quit = true
	case 's':
		err = errors.New("Unknown key. Did you mean capital S for save and quit?")
	case '#':
		err := g.WriteDump()
		if err != nil {
			g.PrintStyled("Error writting dump to file.", logError)
		} else {
			dataDir, _ := g.DataDir()
			g.Printf("Dump written to %s.", filepath.Join(dataDir, "dump"))
		}
		again = true
	case '@':
		if g.Wizard {
			ui.WizardInfo(g)
			again = true
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	default:
		err = fmt.Errorf("Unknown key '%c'. Type ? for help.", c)
	}
	return err, again, quit
}

func (ui *termui) GoToPos(g *game, ev event, pos position) (err error, again bool) {
	if !g.Dungeon.Valid(pos) {
		return errors.New("Invalid location."), true
	}
	switch pos.Distance(g.Player.Pos) {
	case 0:
		g.WaitTurn(ev)
	case 1:
		dir := pos.Dir(g.Player.Pos)
		err = g.MovePlayer(g.Player.Pos.To(dir), ev)
		if err != nil {
			again = true
		}
	default:
		ex := &examiner{}
		err = ex.Action(g, pos)
		if !ex.done || !g.MoveToTarget(ev) {
			again = true
		}
	}
	return err, again
}

func (ui *termui) ExaminePos(g *game, ev event, pos position) (again bool) {
	var start *position
	if g.Dungeon.Valid(pos) {
		start = &pos
	}
	b := ui.Examine(g, start)
	ui.DrawDungeonView(g, false)
	if !b || !g.MoveToTarget(ev) {
		again = true
	}
	return again
}

func (ui *termui) DrawKeysDescription(g *game, actions []string) {
	ui.Clear()
	help := &bytes.Buffer{}
	help.WriteString("┌────────────── Keys ────────────────────────────────────────────────────────\n")
	help.WriteString("│\n")
	for i := 0; i < len(actions)-1; i += 2 {
		fmt.Fprintf(help, "│ %s: %s\n", actions[i], actions[i+1])
	}
	help.WriteString("│\n")
	help.WriteString("└──── press esc or space to return to the game ──────────────────────────────\n")
	ui.DrawText(help.String(), 0, 0)
	ui.Flush()
	ui.WaitForContinue(g)
}

func (ui *termui) KeysHelp(g *game) {
	ui.DrawKeysDescription(g, []string{
		"Movement", "h/j/k/l/y/u/b/n or numpad or mouse left-click",
		"Rest", "r",
		"Wait", "“.” or 5",
		"Use stairs", "> or D",
		"Go to nearest stairs", "G",
		"Quaff potion", "q or a",
		"Equip weapon/armour/...", "e or g",
		"Autoexplore", "o",
		"Examine", "x or mouse right-click",
		"Throw item", "t or f",
		"Evoke rod", "v or z",
		"View Character Information", `% or C`,
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
	ui.Clear()
	b := bytes.Buffer{}
	b.WriteString(formatText(
		fmt.Sprintf("You are wielding %s. %s", Indefinite(g.Player.Weapon.String(), false), g.Player.Weapon.Desc()), 79))
	b.WriteString("\n\n")
	b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Armour, g.Player.Armour.Desc()), 79))
	b.WriteString("\n\n")
	if g.Player.Shield != NoShield {
		b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Shield, g.Player.Shield.Desc()), 79))
		b.WriteString("\n\n")
	}
	b.WriteString(ui.AptitudesText(g))
	ui.DrawText(b.String(), 0, 0)
	ui.Flush()
	ui.WaitForContinue(g)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) WizardInfo(g *game) {
	ui.Clear()
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "Monsters: %d (%d)\n", len(g.Monsters), g.MaxMonsters())
	fmt.Fprintf(b, "Danger: %d (%d)\n", g.Danger(), g.MaxDanger())
	ui.DrawText(b.String(), 0, 0)
	ui.Flush()
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
	var text string
	if len(apts) > 0 {
		text = "Aptitudes:\n" + strings.Join(apts, "\n")
	} else {
		text = "You do not have any special aptitudes."
	}
	text += "\n\n--press esc or space to continue--"
	return text
}

func (ui *termui) DescribePosition(g *game, pos position, targ Targeter) {
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
		desc += fmt.Sprintf("You see %s (%s).", Indefinite(mons.Kind.String(), false), ui.MonsterInfo(mons))
	case g.Gold[pos] > 0:
		desc += fmt.Sprintf("You see some gold (%d).", g.Gold[pos])
	case okCollectable && c != nil:
		if c.Quantity > 1 {
			desc += fmt.Sprintf("You see %d %s there.", c.Quantity, c.Consumable)
		} else {
			desc += fmt.Sprintf("You see %s there.", Indefinite(c.Consumable.String(), false))
		}
	case okEq:
		desc += fmt.Sprintf("You see %s.", Indefinite(eq.String(), false))
	case okRod:
		desc += fmt.Sprintf("You see a %v.", rod)
	case g.Stairs[pos]:
		desc += "You see stairs downwards."
	case g.Doors[pos]:
		desc += "You see a door."
	case g.Dungeon.Cell(pos).T == WallCell:
		desc += "You see a wall."
	default:
		if _, ok := g.Fungus[pos]; ok {
			desc += "You see dense foliage there."
		} else {
			desc += "You see the ground."
		}
	}
	g.Print(desc)
}

func (ui *termui) Examine(g *game, start *position) bool {
	ex := &examiner{}
	err := ui.CursorAction(g, ex, start)
	if err != nil {
		g.Print(err.Error())
		return false
	}
	return ex.done
}

func (ui *termui) ChooseTarget(g *game, targ Targeter) bool {
	err := ui.CursorAction(g, targ, nil)
	if err != nil {
		g.Print(err.Error())
		return false
	}
	return targ.Done()
}

func (ui *termui) NextMonster(g *game, r rune, pos position, nmonster int) (position, int) {
	for i := 0; i < len(g.Monsters); i++ {
		if r == '+' {
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
			pos = mons.Pos
			break
		}
	}
	return pos, nmonster
}

func (ui *termui) NextObject(g *game, pos position, nobject int, objects *[]position) (position, int) {
	if len(*objects) == 0 {
		for p := range g.Collectables {
			*objects = append(*objects, p)
		}
		for p := range g.Rods {
			*objects = append(*objects, p)
		}
		for p := range g.Equipables {
			*objects = append(*objects, p)
		}
		for p := range g.Gold {
			*objects = append(*objects, p)
		}
	}
	for i := 0; i < len(*objects); i++ {
		nobject++
		if nobject > len(*objects)-1 {
			nobject = 0
		}
		p := (*objects)[nobject]
		if g.Dungeon.Cell(p).Explored {
			pos = p
			break
		}
	}
	return pos, nobject
}

func (ui *termui) ExcludeZone(g *game, pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}

func (ui *termui) CursorMouseLeft(g *game, targ Targeter, pos position) bool {
	err := targ.Action(g, pos)
	if err != nil {
		g.Print(err.Error())
	} else {
		return true
	}
	return false
}

func (ui *termui) CursorCharAction(g *game, targ Targeter, r rune, pos position, data *examineData) bool {
	switch r {
	case 'h', '4', 'l', '6', 'j', '2', 'k', '8',
		'y', '7', 'b', '1', 'u', '9', 'n', '3':
		data.npos = pos.To(KeyToDir(r))
	case 'H', 'L', 'J', 'K', 'Y', 'B', 'U', 'N':
		for i := 0; i < 5; i++ {
			p := data.npos.To(KeyToDir(r))
			if !g.Dungeon.Valid(p) {
				break
			}
			data.npos = p
		}
	case '>', 'D':
		if data.sortedStairs == nil {
			stairs := g.StairsSlice()
			data.sortedStairs = g.SortedNearestTo(stairs, g.Player.Pos)
		}
		if data.stairIndex >= len(data.sortedStairs) {
			data.stairIndex = 0
		}
		if len(data.sortedStairs) > 0 {
			data.npos = data.sortedStairs[data.stairIndex]
			data.stairIndex++
		}
	case '+', '-':
		data.npos, data.nmonster = ui.NextMonster(g, r, pos, data.nmonster)
	case 'o':
		data.npos, data.nobject = ui.NextObject(g, pos, data.nobject, &data.objects)
	case 'v', 'd':
		ui.HideCursor()
		ui.ViewPositionDescription(g, pos)
		ui.SetCursor(pos)
	case '?':
		ui.HideCursor()
		ui.ExamineHelp(g)
		ui.SetCursor(pos)
	case '.':
		err := targ.Action(g, pos)
		if err != nil {
			g.Print(err.Error())
		} else {
			return true
		}
	case 'e':
		ui.ExcludeZone(g, pos)
	default:
		g.Print("Invalid key. Type ? for help.")
	}
	return false
}

type examineData struct {
	npos         position
	nmonster     int
	objects      []position
	nobject      int
	sortedStairs []position
	stairIndex   int
}

func (ui *termui) CursorAction(g *game, targ Targeter, start *position) error {
	pos := g.Player.Pos
	if start != nil {
		pos = *start
	} else {
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
	}
	var err error
	data := &examineData{
		npos:    pos,
		objects: []position{},
	}
	opos := position{-1, -1}
loop:
	for {
		err = nil
		if pos != opos {
			ui.DescribePosition(g, pos, targ)
		}
		opos = pos
		targ.ComputeHighlight(g, pos)
		ui.SetCursor(pos)
		ui.DrawDungeonView(g, true)
		data.npos = pos
		b := ui.TargetModeEvent(g, targ, pos, data)
		if b {
			break loop
		}
		if g.Dungeon.Valid(data.npos) {
			pos = data.npos
		}
	}
	g.Highlight = nil
	ui.HideCursor()
	return err
}

func (ui *termui) ViewPositionDescription(g *game, pos position) {
	mons, _ := g.MonsterAt(pos)
	if mons.Exists() {
		ui.HideCursor()
		ui.DrawMonsterDescription(g, mons)
		ui.SetCursor(pos)
	} else if c, ok := g.Collectables[pos]; ok {
		ui.DrawDescription(g, c.Consumable.Desc())
	} else if r, ok := g.Rods[pos]; ok {
		ui.DrawDescription(g, r.Desc())
	} else if eq, ok := g.Equipables[pos]; ok {
		ui.DrawDescription(g, eq.Desc())
	} else if g.Stairs[pos] {
		ui.DrawDescription(g, "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you.")
	} else if g.Doors[pos] {
		ui.DrawDescription(g, "A closed door blocks your line of sight. Doors open automatically when you or a monster stand on them.")
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

var CenteredCamera bool

func (ui *termui) InView(g *game, pos position, targeting bool) bool {
	if targeting {
		return pos.DistanceY(ui.cursor) <= 10 && pos.DistanceX(ui.cursor) <= 39
	}
	return pos.DistanceY(g.Player.Pos) <= 10 && pos.DistanceX(g.Player.Pos) <= 39
}

func (ui *termui) CameraOffset(g *game, pos position, targeting bool) (int, int) {
	if targeting {
		return pos.X + 39 - ui.cursor.X, pos.Y + 10 - ui.cursor.Y
	}
	return pos.X + 39 - g.Player.Pos.X, pos.Y + 10 - g.Player.Pos.Y
}

func (ui *termui) InViewBorder(g *game, pos position, targeting bool) bool {
	if targeting {
		return pos.DistanceY(ui.cursor) != 10 && pos.DistanceX(ui.cursor) != 39
	}
	return pos.DistanceY(g.Player.Pos) != 10 && pos.DistanceX(g.Player.Pos) != 39
}

func (ui *termui) DrawDungeonView(g *game, targeting bool) {
	ui.Clear()
	m := g.Dungeon
	for i := 0; i < g.Dungeon.Width; i++ {
		ui.SetCell(i, g.Dungeon.Heigth, '─', ColorFg, ColorBg)
	}
	for i := 0; i < g.Dungeon.Heigth; i++ {
		ui.SetCell(g.Dungeon.Width, i, '│', ColorFg, ColorBg)
	}
	ui.SetCell(g.Dungeon.Width, g.Dungeon.Heigth, '┘', ColorFg, ColorBg)
	for i := range m.Cells {
		pos := m.CellPosition(i)
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		if g.Highlight[pos] || pos == ui.cursor {
			bgColor, fgColor = fgColor, bgColor
		}
		if CenteredCamera {
			if !ui.InView(g, pos, targeting) {
				continue
			}
			x, y := ui.CameraOffset(g, pos, targeting)
			ui.SetCell(x, y, r, fgColor, bgColor)
			if ui.InViewBorder(g, pos, targeting) && g.Dungeon.Border(pos) {
				for _, opos := range g.Dungeon.OutsideNeighbors(pos) {
					xo, yo := ui.CameraOffset(g, opos, targeting)
					ui.SetCell(xo, yo, '#', ColorFg, ColorBgBorder)
				}
			}
		} else {
			ui.SetCell(pos.X, pos.Y, r, fgColor, bgColor)
		}
	}
	ui.DrawText(fmt.Sprintf("[ %v (%d)", g.Player.Armour, g.Player.Armor()), 81, 0)
	ui.DrawText(fmt.Sprintf(") %v (%d)", g.Player.Weapon, g.Player.Attack()), 81, 1)
	if g.Player.Shield != NoShield {
		if g.Player.Weapon.TwoHanded() {
			ui.DrawText(fmt.Sprintf("] %v (unusable)", g.Player.Shield), 81, 2)
		} else {
			ui.DrawText(fmt.Sprintf("] %v (%d)", g.Player.Shield, g.Player.Block()), 81, 2)
		}
	}
	if targeting {
		ui.DrawColoredText("Targeting", 81, 20, ColorFgTargetMode)
		ui.DrawColoredText("(? for help)", 81, 21, ColorFgTargetMode)
	}
	ui.DrawStatusLine(g)
	ui.DrawLog(g)
	ui.Flush()
}

func (ui *termui) PositionDrawing(g *game, pos position) (r rune, fgColor, bgColor uicolor) {
	m := g.Dungeon
	c := m.Cell(pos)
	fgColor = ColorFg
	bgColor = ColorBg
	if !c.Explored && !g.Wizard {
		r = ' '
		if g.HasFreeExploredNeighbor(pos) {
			r = '¤'
			fgColor = ColorFgDark
			bgColor = ColorBgDark
		}
		if g.Noise[pos] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
			bgColor = ColorBgDark
		}
		return
	}
	if g.Wizard {
		if !c.Explored && g.HasFreeExploredNeighbor(pos) {
			r = '¤'
			fgColor = ColorFgDark
			bgColor = ColorBgDark
			return
		}
		if c.T == WallCell {
			if len(g.Dungeon.FreeNeighbors(pos)) == 0 {
				r = ' '
				return
			}
		}
	}
	if g.Player.LOS[pos] {
		fgColor = ColorFgLOS
		bgColor = ColorBgLOS
		if _, ok := g.Clouds[pos]; ok {
			bgColor = ColorBgCloud
		}
	} else {
		fgColor = ColorFgDark
		bgColor = ColorBgDark
	}
	if g.ExclusionsMap[pos] {
		fgColor = ColorFgExcluded
	}
	switch c.T {
	case WallCell:
		r = '#'
		if g.TemporalWalls[pos] {
			fgColor = ColorFgTemporalWall
		}
	case FreeCell:
		if g.UnknownDig[pos] {
			r = '#'
			if g.TemporalWalls[pos] {
				fgColor = ColorFgTemporalWall
			}
			break
		}
		switch {
		case pos == g.Player.Pos:
			r = '@'
			fgColor = ColorFgPlayer
		default:
			r = '.'
			if _, ok := g.Fungus[pos]; ok {
				r = '"'
			}
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
			} else if _, ok := g.Doors[pos]; ok {
				r = '+'
				fgColor = ColorFgStairs
			}
			if g.Player.LOS[pos] || g.Wizard {
				m, _ := g.MonsterAt(pos)
				if m.Exists() {
					r = m.Kind.Letter()
					if m.Status(MonsConfused) {
						fgColor = ColorFgConfusedMonster
					} else if m.State == Resting {
						fgColor = ColorFgSleepingMonster
					} else if m.State == Wandering {
						fgColor = ColorFgWanderingMonster
					} else {
						fgColor = ColorFgMonster
					}
				}
			} else if !g.Wizard && g.Noise[pos] {
				r = '♫'
				fgColor = ColorFgWanderingMonster
			}
		}
	}
	return
}

func (ui *termui) DrawStatusLine(g *game) {
	sts := statusSlice{}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st)
		}
	}
	sort.Sort(sts)
	hpColor := ColorFgHPok
	switch {
	case g.Player.HP*100/g.Player.HPMax() < 30:
		hpColor = ColorFgHPcritical
	case g.Player.HP*100/g.Player.HPMax() < 70:
		hpColor = ColorFgHPwounded
	}
	mpColor := ColorFgMPok
	switch {
	case g.Player.MP*100/g.Player.MPMax() < 30:
		mpColor = ColorFgMPcritical
	case g.Player.MP*100/g.Player.MPMax() < 70:
		mpColor = ColorFgMPpartial
	}
	ui.DrawColoredText(fmt.Sprintf("HP: %d", g.Player.HP), 81, 4, hpColor)
	ui.DrawColoredText(fmt.Sprintf("MP: %d", g.Player.MP), 81, 5, mpColor)
	ui.DrawText(fmt.Sprintf("Gold: %d", g.Player.Gold), 81, 7)
	if g.Depth > g.MaxDepth() {
		ui.DrawText("Depth: Out!", 81, 8)
	} else {
		ui.DrawText(fmt.Sprintf("Depth: %d", g.Depth), 81, 8)
	}
	ui.DrawText(fmt.Sprintf("Turns: %.1f", float64(g.Turn)/10), 81, 9)

	for i, st := range sts {
		fg := ColorFgStatusOther
		if st.Good() {
			fg = ColorFgStatusGood
		} else if st.Bad() {
			fg = ColorFgStatusBad
		}
		if g.Player.Statuses[st] > 1 {
			ui.DrawColoredText(fmt.Sprintf("%s (%d)", st, g.Player.Statuses[st]), 81, 10+i, fg)
		} else {
			ui.DrawColoredText(st.String(), 81, 10+i, fg)
		}
	}
}

func (ui *termui) LogColor(e logEntry) uicolor {
	fg := ColorFg
	// TODO: define uicolors?
	switch e.Style {
	case logCritic:
		fg = ColorFgHPcritical
	case logPlayerHit:
		fg = ColorFgHPok
	case logMonsterHit:
		fg = ColorFgHPwounded
	case logSpecial:
		fg = ColorFgStairs
	case logStatusEnd:
		fg = ColorFgSleepingMonster
	case logError:
		fg = ColorFgHPcritical
	}
	return fg
}

func (ui *termui) DrawLog(g *game) {
	min := len(g.Log) - 4
	if min < 0 {
		min = 0
	}
	for i, e := range g.Log[min:] {
		fguicolor := ui.LogColor(e)
		if e.Tick {
			ui.DrawColoredText("•", 0, g.Dungeon.Heigth+1+i, ColorFgCollectable)
			ui.DrawColoredText(e.String(), 2, g.Dungeon.Heigth+1+i, fguicolor)
		} else {
			ui.DrawColoredText(e.String(), 0, g.Dungeon.Heigth+1+i, fguicolor)
		}
	}
}

func (ui *termui) DrawPreviousLogs(g *game) {
	lines := 23
	nmax := len(g.Log) - lines
	n := nmax
loop:
	for {
		ui.Clear()
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
			e := g.Log[i]
			fguicolor := ui.LogColor(e)
			if e.Tick {
				ui.DrawColoredText("•", 0, i-n, ColorFgCollectable)
				ui.DrawColoredText(e.String(), 2, i-n, fguicolor)
			} else {
				ui.DrawColoredText(e.String(), 0, i-n, fguicolor)
			}
		}
		s := fmt.Sprintf("─────────(%d/%d)───────────────────────────────────────────────────────────────\n", len(g.Log)-to, len(g.Log))
		ui.DrawText(s, 0, to-n)
		ui.DrawText("Keys: half-page up (u), half-page down (d), quit (esc or space)", 0, to+1-n)
		ui.Flush()
		var quit bool
		n, quit = ui.Scroll(n)
		if quit {
			break loop
		}
	}
}

func (ui *termui) DrawMonsterDescription(g *game, mons *monster) {
	s := mons.Kind.Desc()
	s += " " + fmt.Sprintf("They can hit for up to %d damage.", mons.Kind.BaseAttack())
	s += " " + fmt.Sprintf("They have around %d HP.", mons.Kind.MaxHP())
	ui.DrawDescription(g, s)
}

func (ui *termui) DrawConsumableDescription(g *game, c consumable) {
	ui.DrawDescription(g, c.Desc())
}

func (ui *termui) DrawDescription(g *game, desc string) {
	ui.Clear()
	desc = formatText(desc, 79)
	lines := strings.Count(desc, "\n")
	ui.DrawText(desc, 0, 0)
	ui.DrawText("--press esc or space to continue--", 0, lines+2)
	ui.Flush()
	ui.WaitForContinue(g)
}

func (ui *termui) DrawText(text string, x, y int) {
	ui.DrawColoredText(text, x, y, ColorFg)
}

func (ui *termui) DrawColoredText(text string, x, y int, fg uicolor) {
	col := 0
	for _, r := range text {
		if r == '\n' {
			y++
			col = 0
			continue
		}
		ui.SetCell(x+col, y, r, fg, ColorBg)
		col++
	}
}

func (ui *termui) SelectProjectile(g *game, ev event) error {
	desc := false
	for {
		ui.Clear()
		cs := g.SortedProjectiles()
		if desc {
			ui.DrawText("Describe which projectile? (press ? for throwing menu, esc to return to game)", 0, 0)
		} else {
			ui.DrawText("Throw which projectile? (press ? for describe menu, esc to return to game)", 0, 0)
		}
		for i, c := range cs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1)
		}
		ui.Flush()
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
			b := ui.ChooseTarget(g, &chooser{needsFreeWay: true})
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
		ui.Clear()
		cs := g.SortedPotions()
		if desc {
			ui.DrawText("Describe which potion? (press ? for quaff menu, esc to return to game)", 0, 0)
		} else {
			ui.DrawText("Drink which potion? (press ? for description menu, esc to return to game)", 0, 0)
		}
		for i, c := range cs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1)
		}
		ui.Flush()
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
		ui.Clear()
		rs := g.SortedRods()
		if desc {
			ui.DrawText("Describe which rod? (press ? for evocation menu, esc to return to game)", 0, 0)
		} else {
			ui.DrawText("Evoke which rod? (press ? for description menu, esc to return to game)", 0, 0)
		}
		for i, c := range rs {
			ui.DrawText(fmt.Sprintf("%c - %s (%d/%d charges, %d mana cost)",
				rune(i+97), c, g.Player.Rods[c].Charge, c.MaxCharge(), c.MPCost()), 0, i+1)
		}
		ui.Flush()
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

func (ui *termui) Death(g *game) {
	g.Print("You die... --press esc or space to continue--")
	ui.DrawDungeonView(g, false)
	ui.WaitForContinue(g)
	err := g.WriteDump()
	ui.Dump(g, err)
	ui.WaitForContinue(g)
}

func (ui *termui) Win(g *game) {
	err := g.RemoveSaveFile()
	if err != nil {
		g.PrintfStyled("Error removing saved file: %v", logError, err)
	}
	if g.Wizard {
		g.Print("You escape by the magic stairs! **WIZARD** --press esc or space to continue--")
	} else {
		g.Print("You escape by the magic stairs! You win. --press esc or space to continue--")
	}
	ui.DrawDungeonView(g, false)
	ui.WaitForContinue(g)
	err = g.WriteDump()
	ui.Dump(g, err)
	ui.WaitForContinue(g)
}

func (ui *termui) Dump(g *game, err error) {
	ui.Clear()
	ui.DrawText(g.SimplifedDump(err), 0, 0)
	ui.Flush()
}

func (ui *termui) CriticalHPWarning(g *game) {
	g.PrintStyled("*** CRITICAL HP WARNING *** --press esc or space to continue--", logCritic)
	ui.DrawDungeonView(g, false)
	ui.WaitForContinue(g)
	g.Print("Ok. Be careful, then.")
}

func (ui *termui) Quit(g *game) bool {
	g.Print("Do you really want to quit without saving? [y/N]")
	ui.DrawDungeonView(g, false)
	quit := ui.PromptConfirmation(g)
	if quit {
		err := g.RemoveSaveFile()
		if err != nil {
			g.PrintfStyled("Error removing save file: %v ——press any key to quit——", logError, err)
			ui.DrawDungeonView(g, false)
			ui.PressAnyKey()
		}
	} else {
		g.Print("Ok, then.")
	}
	return quit
}

func (ui *termui) Wizard(g *game) bool {
	g.Print("Do you really want to enter wizard mode (no return)? [y/N]")
	ui.DrawDungeonView(g, false)
	return ui.PromptConfirmation(g)
}

func (ui *termui) HandlePlayerTurn(g *game, ev event) bool {
getKey:
	for {
		ui.DrawDungeonView(g, false)
		err, again, quit := ui.PlayerTurnEvent(g, ev)
		if err != nil {
			g.Print(err.Error())
		}
		if again {
			continue getKey
		}
		return quit
	}
}
