package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

type uicolor int

// uicolors: http://ethanschoonover.com/solarized
var (
	ColorBase03  uicolor = 234
	ColorBase02  uicolor = 235
	ColorBase01  uicolor = 240
	ColorBase00  uicolor = 241 // for dark on light background
	ColorBase0   uicolor = 244
	ColorBase1   uicolor = 245
	ColorBase2   uicolor = 254
	ColorBase3   uicolor = 230
	ColorYellow  uicolor = 136
	ColorOrange  uicolor = 166
	ColorRed     uicolor = 160
	ColorMagenta uicolor = 125
	ColorViolet  uicolor = 61
	ColorBlue    uicolor = 33
	ColorCyan    uicolor = 37
	ColorGreen   uicolor = 64
)

var (
	ColorBg,
	ColorBgBorder,
	ColorBgCloud,
	ColorBgDark,
	ColorBgLOS,
	ColorFg,
	ColorFgAnimationHit,
	ColorFgCollectable,
	ColorFgConfusedMonster,
	ColorFgDark,
	ColorFgExcluded,
	ColorFgExplosionEnd,
	ColorFgExplosionStart,
	ColorFgExplosionWallEnd,
	ColorFgExplosionWallStart,
	ColorFgHPcritical,
	ColorFgHPok,
	ColorFgHPwounded,
	ColorFgLOS,
	ColorFgMPcritical,
	ColorFgMPok,
	ColorFgMPpartial,
	ColorFgMagicPlace,
	ColorFgMonster,
	ColorFgPlace,
	ColorFgPlayer,
	ColorFgProjectile,
	ColorFgSimellas,
	ColorFgSleepingMonster,
	ColorFgStatusBad,
	ColorFgStatusGood,
	ColorFgStatusOther,
	ColorFgTargetMode,
	ColorFgWanderingMonster uicolor
)

func LinkColors() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgCloud = ColorBase2
	ColorBgDark = ColorBase03
	ColorBgLOS = ColorBase3
	ColorFg = ColorBase0
	ColorFgAnimationHit = ColorMagenta
	ColorFgCollectable = ColorYellow
	ColorFgConfusedMonster = ColorGreen
	ColorFgDark = ColorBase01
	ColorFgExcluded = ColorRed
	ColorFgExplosionEnd = ColorOrange
	ColorFgExplosionStart = ColorYellow
	ColorFgExplosionWallEnd = ColorMagenta
	ColorFgExplosionWallStart = ColorViolet
	ColorFgHPcritical = ColorRed
	ColorFgHPok = ColorGreen
	ColorFgHPwounded = ColorYellow
	ColorFgLOS = ColorBase00
	ColorFgMPcritical = ColorMagenta
	ColorFgMPok = ColorBlue
	ColorFgMPpartial = ColorViolet
	ColorFgMagicPlace = ColorCyan
	ColorFgMonster = ColorRed
	ColorFgPlace = ColorMagenta
	ColorFgPlayer = ColorBlue
	ColorFgProjectile = ColorBlue
	ColorFgSimellas = ColorYellow
	ColorFgSleepingMonster = ColorViolet
	ColorFgStatusBad = ColorRed
	ColorFgStatusGood = ColorBlue
	ColorFgStatusOther = ColorYellow
	ColorFgTargetMode = ColorCyan
	ColorFgWanderingMonster = ColorOrange
}

func SolarizedPalette() {
	ColorBase03 = 8
	ColorBase02 = 0
	ColorBase01 = 10
	ColorBase00 = 11
	ColorBase0 = 12
	ColorBase1 = 14
	ColorBase2 = 7
	ColorBase3 = 15
	ColorYellow = 3
	ColorOrange = 9
	ColorRed = 1
	ColorMagenta = 5
	ColorViolet = 13
	ColorBlue = 4
	ColorCyan = 6
	ColorGreen = 2
}

func FixColor() {
	ColorBase03++
	ColorBase02++
	ColorBase01++
	ColorBase00++
	ColorBase0++
	ColorBase1++
	ColorBase2++
	ColorBase3++
	ColorYellow++
	ColorOrange++
	ColorRed++
	ColorMagenta++
	ColorViolet++
	ColorBlue++
	ColorCyan++
	ColorGreen++
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
	ColorBase03 = Black
	ColorBase02 = Black
	ColorBase01 = Silver
	ColorBase00 = Black
	ColorBase0 = Silver
	ColorBase1 = Silver
	ColorBase2 = Silver
	ColorBase3 = Silver
	ColorYellow = Olive
	ColorOrange = Purple
	ColorRed = Maroon
	ColorMagenta = Purple
	ColorViolet = Teal
	ColorBlue = Navy
	ColorCyan = Teal
	ColorGreen = Green

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
	ColorFgPlace = Purple
	ColorFgSimellas = Olive
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
	ColorFgMagicPlace = Teal
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
	ui.DrawDark(">", rcol+2, line, ColorFgPlace)
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
			g.PrintfStyled("Error: %v", logError, err)
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
			g.PrintStyled("Error writing dump to file.", logError)
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
	case 'W':
		ui.EnterWizard(g)
		return nil, true, false
	case 'Q':
		if ui.Quit(g) {
			return nil, false, true
		}
		return nil, true, false
	default:
		err = fmt.Errorf("Unknown key '%c'. Type ? for help.", c)
	}
	return err, again, quit
}

func (ui *termui) GoToPos(g *game, ev event, pos position) (err error, again bool) {
	if !pos.valid() {
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
	if pos.valid() {
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
	ui.DrawDungeonView(g, false)

	ui.DrawStyledTextLine(" Keys ", 0, HeaderLine)
	for i := 0; i < len(actions)-1; i += 2 {
		bg := ColorBase03
		if i%4 == 2 {
			bg = ColorBase02
		}
		ui.ClearLineWithColor(i/2+1, bg)
		ui.DrawColoredTextOnBG(fmt.Sprintf(" %-36s %s", actions[i], actions[i+1]), 0, i/2+1, ColorFg, bg)
	}
	lines := 1 + len(actions)/2
	ui.DrawTextLine("press esc or space to continue", lines)
	ui.Flush()

	ui.WaitForContinue(g)
}

func (ui *termui) KeysHelp(g *game) {
	ui.DrawKeysDescription(g, []string{
		"Movement", "h/j/k/l/y/u/b/n or numpad or mouse left",
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
		"View Character and Quest Information", `% or C`,
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
		"Toggle exclude area from auto-travel", "e",
	})
}

func (ui *termui) Equip(g *game, ev event) error {
	return g.Equip(ev)
}

const TextWidth = 78

func (ui *termui) CharacterInfo(g *game) {
	ui.DrawDungeonView(g, false)

	b := bytes.Buffer{}
	b.WriteString(formatText("Every year, your village sends someone to collect medicinal simella plants in the Underground. This year, the duty fell upon you, and so here you are. Your heart is teared between your will to be as helpful as possible to your village and your will to make it out alive.", TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(
		fmt.Sprintf("You are wielding %s. %s", Indefinite(g.Player.Weapon.String(), false), g.Player.Weapon.Desc()), TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Armour, g.Player.Armour.Desc()), TextWidth))
	b.WriteString("\n\n")
	if g.Player.Shield != NoShield {
		b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Shield, g.Player.Shield.Desc()), TextWidth))
		b.WriteString("\n\n")
	}
	b.WriteString(ui.AptitudesText(g))

	desc := b.String()
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	ui.DrawTextLine("press esc or space to continue", lines+2)

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
		desc += fmt.Sprintf("You see %s (%s).", mons.Kind.Indefinite(false), ui.MonsterInfo(mons))
	case g.Simellas[pos] > 0:
		desc += fmt.Sprintf("You see some simellas (%d).", g.Simellas[pos])
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
		if g.Depth == g.MaxDepth() {
			desc += "You see some glowing stairs."
		} else {
			desc += "You see stairs downwards."
		}
	case g.Doors[pos]:
		desc += "You see a door."
	case g.Dungeon.Cell(pos).T == WallCell:
		desc += "You see a wall."
	default:
		if cld, ok := g.Clouds[pos]; ok {
			if cld == CloudFire {
				desc += "You see burning flames."
			} else {
				desc += "You see a dense fog."
			}
		} else if _, ok := g.Fungus[pos]; ok {
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

func (ui *termui) NextMonster(g *game, r rune, pos position, data *examineData) {
	nmonster := data.nmonster
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
	data.npos = pos
	data.nmonster = nmonster
}

func (ui *termui) NextStair(g *game, data *examineData) {
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
}

func (ui *termui) NextObject(g *game, pos position, data *examineData) {
	nobject := data.nobject
	if len(data.objects) == 0 {
		for p := range g.Collectables {
			data.objects = append(data.objects, p)
		}
		for p := range g.Rods {
			data.objects = append(data.objects, p)
		}
		for p := range g.Equipables {
			data.objects = append(data.objects, p)
		}
		for p := range g.Simellas {
			data.objects = append(data.objects, p)
		}
		data.objects = g.SortedNearestTo(data.objects, g.Player.Pos)
	}
	for i := 0; i < len(data.objects); i++ {
		p := data.objects[nobject]
		nobject++
		if nobject > len(data.objects)-1 {
			nobject = 0
		}
		if g.Dungeon.Cell(p).Explored {
			pos = p
			break
		}
	}
	data.npos = pos
	data.nobject = nobject
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
			if !p.valid() {
				break
			}
			data.npos = p
		}
	case '>', 'D':
		ui.NextStair(g, data)
	case '+', '-':
		ui.NextMonster(g, r, pos, data)
	case 'o':
		ui.NextObject(g, pos, data)
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
	if _, ok := targ.(*examiner); ok && pos == g.Player.Pos {
		ui.NextObject(g, position{-1, -1}, data)
		if !data.npos.valid() {
			ui.NextStair(g, data)
		}
		if data.npos.valid() {
			pos = data.npos
		}
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
		if data.npos.valid() {
			pos = data.npos
		}
	}
	g.Highlight = nil
	ui.HideCursor()
	return err
}

func (ui *termui) ViewPositionDescription(g *game, pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("No description: unknown place.")
		return
	}
	mons, _ := g.MonsterAt(pos)
	if mons.Exists() && g.Player.LOS[mons.Pos] {
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
		if g.Depth == g.MaxDepth() {
			ui.DrawDescription(g, "These shiny-looking stairs are in fact a magical monolith. It is said they were made some centuries ago by Marevor Helith. They will lead you back to your village.")
		} else {
			ui.DrawDescription(g, "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you.")
		}
	} else if g.Doors[pos] {
		ui.DrawDescription(g, "A closed door blocks your line of sight. Doors open automatically when you or a monster stand on them. Doors are flammable.")
	} else if g.Simellas[pos] > 0 {
		ui.DrawDescription(g, "A simella is a plant with big white flowers which are used in the Underground for their medicinal properties. They can also make tasty infusions. You were actually sent here by your village to collect as many as possible of those plants.")
	} else if _, ok := g.Fungus[pos]; ok {
		ui.DrawDescription(g, "Dense foliage is difficult to see through. It is flammable.")
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

func (ui *termui) DrawAtPosition(g *game, pos position, targeting bool, r rune, fg, bg uicolor) {
	if g.Highlight[pos] || pos == ui.cursor {
		bg, fg = fg, bg
	}
	if CenteredCamera {
		if !ui.InView(g, pos, targeting) {
			return
		}
		x, y := ui.CameraOffset(g, pos, targeting)
		ui.SetCell(x, y, r, fg, bg)
		if ui.InViewBorder(g, pos, targeting) && g.Dungeon.Border(pos) {
			for _, opos := range pos.OutsideNeighbors() {
				xo, yo := ui.CameraOffset(g, opos, targeting)
				ui.SetCell(xo, yo, '#', ColorFg, ColorBgBorder)
			}
		}
		return
	}
	ui.SetCell(pos.X, pos.Y, r, fg, bg)
}

func (ui *termui) DrawDungeonView(g *game, targeting bool) {
	ui.Clear()
	m := g.Dungeon
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, DungeonHeight, '─', ColorFg, ColorBg)
	}
	for i := 0; i < DungeonHeight; i++ {
		ui.SetCell(DungeonWidth, i, '│', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, DungeonHeight, '┘', ColorFg, ColorBg)
	for i := range m.Cells {
		pos := idxtopos(i)
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, targeting, r, fgColor, bgColor)
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

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *termui) ExplosionAnimation(g *game, es explosionStyle, pos position) {
	ui.DrawDungeonView(g, false)
	// TODO: use new specific variables for colors
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for _, fg := range colors {
		if es != AroundWallExplosion {
			_, _, bgColor := ui.PositionDrawing(g, pos)
			ui.DrawAtPosition(g, pos, true, '☼', fg, bgColor)
			ui.Flush()
			time.Sleep(15 * time.Millisecond)
		}
		for _, npos := range g.Dungeon.FreeNeighbors(pos) {
			if !g.Player.LOS[npos] {
				continue
			}
			_, _, bgColor := ui.PositionDrawing(g, npos)
			ui.DrawAtPosition(g, npos, true, '¤', fg, bgColor)
			ui.Flush()
			time.Sleep(7 * time.Millisecond)
		}
	}
	time.Sleep(25 * time.Millisecond)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) LightningBoltAnimation(g *game, ray []position) {
	ui.DrawDungeonView(g, false)
	time.Sleep(10 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for _, fg := range colors {
		for i := len(ray) - 1; i >= 0; i-- {
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(g, pos)
			ui.DrawAtPosition(g, pos, true, '☼', fg, bgColor)
			ui.Flush()
			time.Sleep(7 * time.Millisecond)
		}
	}
	time.Sleep(25 * time.Millisecond)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) ProjectileSymbol(dir direction) (r rune) {
	switch dir {
	case E, ENE, ESE, WNW, W, WSW:
		r = '—'
	case NE, SW:
		r = '/'
	case NNE, N, NNW, SSW, S, SSE:
		r = '|'
	case NW, SE:
		r = '\\'
	}
	return r
}

func (ui *termui) ThrowAnimation(g *game, ray []position, hit bool) {
	ui.DrawDungeonView(g, false)
	time.Sleep(10 * time.Millisecond)
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgProjectile, bgColor)
		ui.Flush()
		time.Sleep(15 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
	if hit {
		pos := ray[0]
		_, _, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, '¤', ColorFgAnimationHit, bgColor)
		ui.Flush()
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
	ui.DrawDungeonView(g, false)
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
			fgColor = ColorFgMagicPlace
		}
	case FreeCell:
		if g.UnknownDig[pos] {
			r = '#'
			if g.TemporalWalls[pos] {
				fgColor = ColorFgMagicPlace
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
			if _, ok := g.UnknownBurn[pos]; ok {
				r = '"'
			}
			if cld, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
				r = '§'
				if cld == CloudFire {
					fgColor = ColorFgWanderingMonster
				}
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
				if g.Depth == g.MaxDepth() {
					fgColor = ColorFgMagicPlace
				} else {
					fgColor = ColorFgPlace
				}
			} else if _, ok := g.Simellas[pos]; ok {
				r = '♣'
				fgColor = ColorFgSimellas
			} else if _, ok := g.Doors[pos]; ok {
				r = '+'
				fgColor = ColorFgPlace
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
	ui.DrawText(fmt.Sprintf("Simellas: %d", g.Player.Simellas), 81, 7)
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
		fg = ColorRed
	case logPlayerHit:
		fg = ColorGreen
	case logMonsterHit:
		fg = ColorOrange
	case logSpecial:
		fg = ColorMagenta
	case logStatusEnd:
		fg = ColorViolet
	case logError:
		fg = ColorRed
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
			ui.DrawColoredText("•", 0, DungeonHeight+1+i, ColorYellow)
			ui.DrawColoredText(e.String(), 2, DungeonHeight+1+i, fguicolor)
		} else {
			ui.DrawColoredText(e.String(), 0, DungeonHeight+1+i, fguicolor)
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
				ui.DrawColoredText("•", 0, i-n, ColorYellow)
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
	ui.DrawDungeonView(g, false)
	desc = formatText(desc, TextWidth)
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	ui.DrawTextLine("press esc or space to continue", lines+2)
	ui.Flush()
	ui.WaitForContinue(g)
	ui.DrawDungeonView(g, false)
}

func (ui *termui) DrawText(text string, x, y int) {
	ui.DrawColoredText(text, x, y, ColorFg)
}

func (ui *termui) DrawColoredText(text string, x, y int, fg uicolor) {
	ui.DrawColoredTextOnBG(text, x, y, fg, ColorBg)
}

func (ui *termui) DrawColoredTextOnBG(text string, x, y int, fg, bg uicolor) {
	col := 0
	for _, r := range text {
		if r == '\n' {
			y++
			col = 0
			continue
		}
		ui.SetCell(x+col, y, r, fg, bg)
		col++
	}
}

func (ui *termui) DrawLine(lnum int) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, lnum, '┤', ColorFg, ColorBg)
}

func (ui *termui) DrawTextLine(text string, lnum int) {
	ui.DrawStyledTextLine(text, lnum, NormalLine)
}

type linestyle int

const (
	NormalLine linestyle = iota
	HeaderLine
)

func (ui *termui) DrawStyledTextLine(text string, lnum int, st linestyle) {
	nchars := utf8.RuneCountInString(text)
	dist := (DungeonWidth - nchars) / 2
	for i := 0; i < dist; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.DrawColoredText(text, dist, lnum, ColorYellow)
	default:
		ui.DrawColoredText(text, dist, lnum, ColorFg)
	}
	for i := dist + nchars; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.SetCell(DungeonWidth, lnum, '┐', ColorFg, ColorBg)
	default:
		ui.SetCell(DungeonWidth, lnum, '┤', ColorFg, ColorBg)
	}
}

func (ui *termui) ClearLine(lnum int) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, ColorBg)
	}
}

func (ui *termui) ClearLineWithColor(lnum int, bg uicolor) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, bg)
	}
}

func (ui *termui) SelectProjectile(g *game, ev event) error {
	desc := false
	for {
		cs := g.SortedProjectiles()
		ui.ClearLine(0)
		if desc {
			ui.DrawText("Describe which projectile? (press ? for throwing menu, esc or space to cancel)", 0, 0)
		} else {
			ui.DrawText("Throw which projectile? (press ? for describe menu, esc or space to cancel)", 0, 0)
		}
		for i, c := range cs {
			bg := ColorBase03
			if i%2 == 1 {
				bg = ColorBase02
			}
			ui.ClearLineWithColor(i+1, bg)
			ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1, ColorFg, bg)
		}
		ui.DrawLine(len(cs) + 1)
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
		cs := g.SortedPotions()
		ui.ClearLine(0)
		if desc {
			ui.DrawText("Describe which potion? (press ? for quaff menu, esc or space to cancel)", 0, 0)
		} else {
			ui.DrawText("Drink which potion? (press ? for description menu, esc or space to cancel)", 0, 0)
		}
		for i, c := range cs {
			bg := ColorBase03
			if i%2 == 1 {
				bg = ColorBase02
			}
			ui.ClearLineWithColor(i+1, bg)
			ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, i+1, ColorFg, bg)
		}
		ui.DrawLine(len(cs) + 1)
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
		rs := g.SortedRods()
		ui.ClearLine(0)
		if desc {
			ui.DrawText("Describe which rod? (press ? for evocation menu, esc or space to cancel)", 0, 0)
		} else {
			ui.DrawText("Evoke which rod? (press ? for description menu, esc or space to cancel)", 0, 0)
		}
		for i, c := range rs {
			bg := ColorBase03
			if i%2 == 1 {
				bg = ColorBase02
			}
			ui.ClearLineWithColor(i+1, bg)
			ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d/%d charges, %d mana cost)",
				rune(i+97), c, g.Player.Rods[c].Charge, c.MaxCharge(), c.MPCost()), 0, i+1, ColorFg, bg)
		}
		ui.DrawLine(len(rs) + 1)
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
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, ColorFgHPcritical, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, fg, bg)
	ui.Flush()
	ui.WaitForContinue(g)
	g.Print("Ok. Be careful, then.")
}

func (ui *termui) DrinkingPotionAnimation(g *game) {
	ui.DrawDungeonView(g, false)
	time.Sleep(50 * time.Millisecond)
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, true, r, fg, bg)
	ui.Flush()
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
