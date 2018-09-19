package main

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	UIWidth                = 100
	UIHeight               = 26
	DisableAnimations bool = false
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
	ColorBgLOSalt,
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
	ColorFgStatusExpire,
	ColorFgStatusOther,
	ColorFgTargetMode,
	ColorFgWanderingMonster uicolor
)

func LinkColors() {
	ColorBg = ColorBase03
	ColorBgBorder = ColorBase02
	ColorBgLOSalt = ColorBase2
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
	ColorFgStatusExpire = ColorViolet
	ColorFgStatusOther = ColorYellow
	ColorFgTargetMode = ColorCyan
	ColorFgWanderingMonster = ColorOrange
}

func ApplyDarkLOS() {
	if ColorBg == Black && ColorBgLOS == Silver {
		ColorFgLOS = Green
		ColorBgLOS = Black
		ColorBgLOSalt = Black
	} else {
		ColorBgLOSalt = ColorBase02
		ColorBgLOS = ColorBase02
		ColorFgLOS = ColorBase1
	}
}

func ApplyLightLOS() {
	if ColorBg == Black && ColorBgLOS == Black {
		ColorFgLOS = Black
		ColorBgLOS = Silver
		ColorBgLOSalt = Silver
	} else {
		ColorBgLOSalt = ColorBase2
		ColorBgLOS = ColorBase3
		ColorFgLOS = ColorBase00
	}
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

func Simple8ColorPalette() {
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
}

func (ui *termui) DrawWelcome() {
	ui.Clear()
	col := 10
	line := 5
	rcol := col + 20
	ColorText := ColorFgHPok
	ui.DrawDark(fmt.Sprintf("       Boohu %s", Version), col, line-2, ColorText, false)
	ui.DrawDark("────│\\/\\/\\/\\/\\/\\/\\/│────", col, line, ColorText, false)
	line++
	ui.DrawDark("##", col, line, ColorFgDark, true)
	ui.DrawLight("#", col+2, line, ColorFgLOS, true)
	ui.DrawLightAlt("#", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("####", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark("#.", col, line, ColorFgDark, true)
	ui.DrawLightAlt(".", col+2, line, ColorFgLOS, true)
	ui.DrawLight(".", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark(".", rcol, line, ColorFgDark, true)
	ui.DrawDark("♣", rcol+1, line, ColorFgSimellas, true)
	ui.DrawDark(".#", rcol+2, line, ColorFgDark, true)
	line++
	ui.DrawDark("##", col, line, ColorFgDark, true)
	ui.DrawLight("!", col+2, line, ColorFgCollectable, true)
	ui.DrawLightAlt(".", col+3, line, ColorFgLOS, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("│  BREAK       │", col+4, line, ColorText, false)
	ui.DrawDark(".###", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark(" #", col, line, ColorFgDark, true)
	ui.DrawLightAlt("g", col+2, line, ColorFgMonster, true)
	ui.DrawLight("G", col+3, line, ColorFgMonster, true)
	ui.DrawDark("│  OUT OF      │", col+4, line, ColorText, false)
	ui.DrawDark("##  ", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLight("#", col, line, ColorFgLOS, true)
	ui.DrawLightAlt("#", col+1, line, ColorFgLOS, true)
	ui.DrawLight("D", col+2, line, ColorFgMonster, true)
	ui.DrawLightAlt("g", col+3, line, ColorFgMonster, true)
	ui.DrawDark("│  HAREKA'S    │", col+4, line, ColorText, false)
	ui.DrawDark(".## ", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLightAlt("#", col, line, ColorFgLOS, true)
	ui.DrawLight("@", col+1, line, ColorFgPlayer, true)
	ui.DrawLightAlt("#", col+2, line, ColorFgLOS, true)
	ui.DrawDark("#", col+3, line, ColorFgDark, true)
	ui.DrawDark("│  UNDERGROUND │", col+4, line, ColorText, false)
	ui.DrawDark("\".##", rcol, line, ColorFgDark, true)
	line++
	ui.DrawLight("#", col, line, ColorFgLOS, true)
	ui.DrawLightAlt(".", col+1, line, ColorFgLOS, true)
	ui.DrawLight("#", col+2, line, ColorFgLOS, true)
	ui.DrawDark("#", col+3, line, ColorFgDark, true)
	ui.DrawDark("│              │", col+4, line, ColorText, false)
	ui.DrawDark("#.", rcol, line, ColorFgDark, true)
	ui.DrawDark(">", rcol+2, line, ColorFgPlace, true)
	ui.DrawDark("#", rcol+3, line, ColorFgDark, true)
	line++
	ui.DrawLightAlt("#", col, line, ColorFgLOS, true)
	ui.DrawLight("[", col+1, line, ColorFgCollectable, true)
	ui.DrawLightAlt(".", col+2, line, ColorFgLOS, true)
	ui.DrawDark("##", col+3, line, ColorFgDark, true)
	ui.DrawDark("│              │", col+4, line, ColorFgHPok, false)
	ui.DrawDark("\"\"##", rcol, line, ColorFgDark, true)
	line++
	ui.DrawDark("────│/\\/\\/\\/\\/\\/\\/\\│────", col, line, ColorText, false)
	line++
	line++
	if runtime.GOOS == "js" {
		ui.DrawDark("───Click on the image to play───", col-3, line, ColorFg, false)
	} else {
		ui.DrawDark("───Press any key to continue───", col-3, line, ColorFg, false)
	}
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

func (ui *termui) DrawDark(text string, x, y int, fg uicolor, inmap bool) {
	col := 0
	for _, r := range text {
		if inmap {
			ui.SetMapCell(x+col, y, r, fg, ColorBgDark)
		} else {
			ui.SetCell(x+col, y, r, fg, ColorBgDark)
		}
		col++
	}
}

func (ui *termui) DrawLight(text string, x, y int, fg uicolor, inmap bool) {
	col := 0
	for _, r := range text {
		if inmap {
			ui.SetMapCell(x+col, y, r, fg, ColorBgLOS)
		} else {
			ui.SetCell(x+col, y, r, fg, ColorBgLOS)
		}
		col++
	}
}

func (ui *termui) DrawLightAlt(text string, x, y int, fg uicolor, inmap bool) {
	col := 0
	for _, r := range text {
		if inmap {
			ui.SetMapCell(x+col, y, r, fg, ColorBgLOSalt)
		} else {
			ui.SetCell(x+col, y, r, fg, ColorBgLOSalt)
		}
		col++
	}
}

const DoNothing = "Do nothing, then."

type uiMode int

const (
	NormalMode uiMode = iota
	TargetingMode
	NoFlushMode
)

func (ui *termui) EnterWizard(g *game) {
	if ui.Wizard(g) {
		g.WizardMode()
		ui.DrawDungeonView(g, NoFlushMode)
	} else {
		g.Print(DoNothing)
	}
}

func (ui *termui) CleanError(err error) error {
	if err != nil && err.Error() == DoNothing {
		err = errors.New("")
	}
	return err
}

type keyAction int

const (
	KeyNothing keyAction = iota
	KeyW
	KeyS
	KeyN
	KeyE
	KeyNW
	KeyNE
	KeySW
	KeySE
	KeyRunW
	KeyRunS
	KeyRunN
	KeyRunE
	KeyRunNW
	KeyRunNE
	KeyRunSW
	KeyRunSE
	KeyRest
	KeyWaitTurn
	KeyDescend
	KeyGoToStairs
	KeyExplore
	KeyExamine
	KeyEquip
	KeyDrink
	KeyThrow
	KeyEvoke
	KeyCharacterInfo
	KeyLogs
	KeyDump
	KeyHelp
	KeySave
	KeyQuit
	KeyWizard
	KeyWizardInfo

	KeyPreviousMonster
	KeyNextMonster
	KeyNextObject
	KeyDescription
	KeyTarget
	KeyExclude
	KeyEscape

	KeyConfigure
	KeyMenu
	KeyNextStairs
	KeyMenuCommandHelp
	KeyMenuTargetingHelp
)

var configurableKeyActions = [...]keyAction{
	KeyW,
	KeyS,
	KeyN,
	KeyE,
	KeyNW,
	KeyNE,
	KeySW,
	KeySE,
	KeyRunW,
	KeyRunS,
	KeyRunN,
	KeyRunE,
	KeyRunNW,
	KeyRunNE,
	KeyRunSW,
	KeyRunSE,
	KeyRest,
	KeyWaitTurn,
	KeyDescend,
	KeyGoToStairs,
	KeyExplore,
	KeyExamine,
	KeyEquip,
	KeyDrink,
	KeyThrow,
	KeyEvoke,
	KeyCharacterInfo,
	KeyLogs,
	KeyDump,
	KeySave,
	KeyQuit,
	KeyPreviousMonster,
	KeyNextMonster,
	KeyNextObject,
	KeyNextStairs,
	KeyDescription,
	KeyTarget,
	KeyExclude}

var CustomKeys bool

func FixedRuneKey(r rune) bool {
	switch r {
	case ' ', '?', '=':
		return true
	default:
		return false
	}
}

func (k keyAction) NormalModeKey() bool {
	switch k {
	case KeyW, KeyS, KeyN, KeyE,
		KeyNW, KeyNE, KeySW, KeySE,
		KeyRunW, KeyRunS, KeyRunN, KeyRunE,
		KeyRunNW, KeyRunNE, KeyRunSW, KeyRunSE,
		KeyRest,
		KeyWaitTurn,
		KeyDescend,
		KeyGoToStairs,
		KeyExplore,
		KeyExamine,
		KeyEquip,
		KeyDrink,
		KeyThrow,
		KeyEvoke,
		KeyCharacterInfo,
		KeyLogs,
		KeyDump,
		KeyHelp,
		KeyMenuCommandHelp,
		KeyMenuTargetingHelp,
		KeySave,
		KeyQuit,
		KeyConfigure,
		KeyWizard,
		KeyWizardInfo:
		return true
	default:
		return false
	}
}

func (k keyAction) NormalModeDescription() (text string) {
	switch k {
	case KeyW:
		text = "Move west"
	case KeyS:
		text = "Move south"
	case KeyN:
		text = "Move north"
	case KeyE:
		text = "Move east"
	case KeyNW:
		text = "Move north west"
	case KeyNE:
		text = "Move north east"
	case KeySW:
		text = "Move south west"
	case KeySE:
		text = "Move south east"
	case KeyRunW:
		text = "Travel west"
	case KeyRunS:
		text = "Travel south"
	case KeyRunN:
		text = "Travel north"
	case KeyRunE:
		text = "Travel east"
	case KeyRunNW:
		text = "Travel north west"
	case KeyRunNE:
		text = "Travel north east"
	case KeyRunSW:
		text = "Travel south west"
	case KeyRunSE:
		text = "Travel south east"
	case KeyRest:
		text = "Rest (until status free or regen)"
	case KeyWaitTurn:
		text = "Wait a turn"
	case KeyDescend:
		text = "Descend stairs"
	case KeyGoToStairs:
		text = "Go to nearest stairs"
	case KeyExplore:
		text = "Autoexplore"
	case KeyExamine:
		text = "Examine"
	case KeyEquip:
		text = "Equip weapon/armour/..."
	case KeyDrink:
		text = "Quaff potion"
	case KeyThrow:
		text = "Throw item"
	case KeyEvoke:
		text = "Evoke rod"
	case KeyCharacterInfo:
		text = "View Character and Quest Information"
	case KeyLogs:
		text = "View previous messages"
	case KeyDump:
		text = "Write game statistics to file"
	case KeySave:
		text = "Save and Quit"
	case KeyQuit:
		text = "Quit without saving"
	case KeyHelp:
		text = "Help (keys and mouse)"
	case KeyMenuCommandHelp:
		text = "Help (general commands)"
	case KeyMenuTargetingHelp:
		text = "Help (targeting commands)"
	case KeyConfigure:
		text = "Settings and key bindings"
	case KeyWizard:
		text = "Wizard (debug) mode"
	case KeyWizardInfo:
		text = "Wizard (debug) mode information"
	case KeyMenu:
		text = "Action Menu"
	}
	return text
}

func (k keyAction) TargetingModeDescription() (text string) {
	switch k {
	case KeyW:
		text = "Move cursor west"
	case KeyS:
		text = "Move cursor south"
	case KeyN:
		text = "Move cursor north"
	case KeyE:
		text = "Move cursor east"
	case KeyNW:
		text = "Move cursor north west"
	case KeyNE:
		text = "Move cursor north east"
	case KeySW:
		text = "Move cursor south west"
	case KeySE:
		text = "Move cursor south east"
	case KeyRunW:
		text = "Big move cursor west"
	case KeyRunS:
		text = "Big move cursor south"
	case KeyRunN:
		text = "Big move north"
	case KeyRunE:
		text = "Big move east"
	case KeyRunNW:
		text = "Big move north west"
	case KeyRunNE:
		text = "Big move north east"
	case KeyRunSW:
		text = "Big move south west"
	case KeyRunSE:
		text = "Big move south east"
	case KeyDescend:
		text = "Target next stair"
	case KeyPreviousMonster:
		text = "Target previous monster"
	case KeyNextMonster:
		text = "Target next monster"
	case KeyNextObject:
		text = "Target next object"
	case KeyNextStairs:
		text = "Target next stairs"
	case KeyDescription:
		text = "View target description"
	case KeyTarget:
		text = "Go to/select target"
	case KeyExclude:
		text = "Toggle exclude area from auto-travel"
	case KeyEscape:
		text = "Quit targeting mode"
	case KeyMenu:
		text = "Action Menu"
	}
	return text
}

func (k keyAction) TargetingModeKey() bool {
	switch k {
	case KeyW, KeyS, KeyN, KeyE,
		KeyNW, KeyNE, KeySW, KeySE,
		KeyRunW, KeyRunS, KeyRunN, KeyRunE,
		KeyRunNW, KeyRunNE, KeyRunSW, KeyRunSE,
		KeyDescend,
		KeyPreviousMonster,
		KeyNextMonster,
		KeyNextObject,
		KeyNextStairs,
		KeyDescription,
		KeyTarget,
		KeyExclude,
		KeyEscape:
		return true
	default:
		return false
	}
}

var gameConfig config

func ApplyDefaultKeyBindings() {
	gameConfig.RuneNormalModeKeys = map[rune]keyAction{
		'h': KeyW,
		'j': KeyS,
		'k': KeyN,
		'l': KeyE,
		'y': KeyNW,
		'u': KeyNE,
		'b': KeySW,
		'n': KeySE,
		'4': KeyW,
		'2': KeyS,
		'8': KeyN,
		'6': KeyE,
		'7': KeyNW,
		'9': KeyNE,
		'1': KeySW,
		'3': KeySE,
		'H': KeyRunW,
		'J': KeyRunS,
		'K': KeyRunN,
		'L': KeyRunE,
		'Y': KeyRunNW,
		'U': KeyRunNE,
		'B': KeyRunSW,
		'N': KeyRunSE,
		'.': KeyWaitTurn,
		'5': KeyWaitTurn,
		'r': KeyRest,
		'>': KeyDescend,
		'D': KeyDescend,
		'G': KeyGoToStairs,
		'o': KeyExplore,
		'x': KeyExamine,
		'e': KeyEquip,
		'g': KeyEquip,
		',': KeyEquip,
		'q': KeyDrink,
		'd': KeyDrink,
		't': KeyThrow,
		'f': KeyThrow,
		'v': KeyEvoke,
		'z': KeyEvoke,
		'%': KeyCharacterInfo,
		'C': KeyCharacterInfo,
		'm': KeyLogs,
		'#': KeyDump,
		'?': KeyHelp,
		'S': KeySave,
		'Q': KeyQuit,
		'W': KeyWizard,
		'@': KeyWizardInfo,
		'=': KeyConfigure,
	}
	gameConfig.RuneTargetModeKeys = map[rune]keyAction{
		'h':    KeyW,
		'j':    KeyS,
		'k':    KeyN,
		'l':    KeyE,
		'y':    KeyNW,
		'u':    KeyNE,
		'b':    KeySW,
		'n':    KeySE,
		'4':    KeyW,
		'2':    KeyS,
		'8':    KeyN,
		'6':    KeyE,
		'7':    KeyNW,
		'9':    KeyNE,
		'1':    KeySW,
		'3':    KeySE,
		'H':    KeyRunW,
		'J':    KeyRunS,
		'K':    KeyRunN,
		'L':    KeyRunE,
		'Y':    KeyRunNW,
		'U':    KeyRunNE,
		'B':    KeyRunSW,
		'N':    KeyRunSE,
		'>':    KeyNextStairs,
		'-':    KeyPreviousMonster,
		'+':    KeyNextMonster,
		'o':    KeyNextObject,
		']':    KeyNextObject,
		')':    KeyNextObject,
		'(':    KeyNextObject,
		'[':    KeyNextObject,
		'_':    KeyNextObject,
		'v':    KeyDescription,
		'd':    KeyDescription,
		'.':    KeyTarget,
		'e':    KeyExclude,
		' ':    KeyEscape,
		'\x1b': KeyEscape,
		'?':    KeyHelp,
	}
	CustomKeys = false
}

type runeKeyAction struct {
	r rune
	k keyAction
}

func (ui *termui) HandleKeyAction(g *game, rka runeKeyAction) (err error, again bool, quit bool) {
	if rka.r != 0 {
		var ok bool
		rka.k, ok = gameConfig.RuneNormalModeKeys[rka.r]
		if !ok {
			switch rka.r {
			case 's':
				err = errors.New("Unknown key. Did you mean capital S for save and quit?")
			default:
				err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
			}
			return err, again, quit
		}
	}
	if rka.k == KeyMenu {
		rka.k, err = ui.SelectAction(g, menuActions, g.Ev)
		if err != nil {
			err = ui.CleanError(err)
			return err, again, quit
		}
	}
	return ui.HandleKey(g, rka)
}

func (ui *termui) OptionalDescendConfirmation(g *game, st stair) (err error) {
	if g.Depth == WinDepth && st == NormalStair {
		g.Print("Do you really want to dive into optional depths? [y/N]")
		ui.DrawDungeonView(g, NormalMode)
		dive := ui.PromptConfirmation(g)
		if !dive {
			err = errors.New("Keep going in the current level, then.")
		}
	}
	return err

}

func (ui *termui) HandleKey(g *game, rka runeKeyAction) (err error, again bool, quit bool) {
	switch rka.k {
	case KeyW, KeyS, KeyN, KeyE, KeyNW, KeyNE, KeySW, KeySE:
		err = g.MovePlayer(g.Player.Pos.To(KeyToDir(rka.k)), g.Ev)
	case KeyRunW, KeyRunS, KeyRunN, KeyRunE, KeyRunNW, KeyRunNE, KeyRunSW, KeyRunSE:
		err = g.GoToDir(KeyToDir(rka.k), g.Ev)
	case KeyWaitTurn:
		g.WaitTurn(g.Ev)
	case KeyRest:
		err = g.Rest(g.Ev)
		ui.MenuSelectedAnimation(g, MenuRest, err == nil)
	case KeyDescend:
		if st, ok := g.Stairs[g.Player.Pos]; ok {
			ui.MenuSelectedAnimation(g, MenuInteract, true)
			err = ui.OptionalDescendConfirmation(g, st)
			if err != nil {
				break
			}
			if g.Descend() {
				ui.Win(g)
				quit = true
				return err, again, quit
			}
			ui.DrawDungeonView(g, NormalMode)
		} else {
			err = errors.New("No stairs here.")
		}
	case KeyGoToStairs:
		stairs := g.StairsSlice()
		sortedStairs := g.SortedNearestTo(stairs, g.Player.Pos)
		if len(sortedStairs) > 0 {
			stair := sortedStairs[0]
			if g.Player.Pos == stair {
				err = errors.New("You are already on the stairs.")
				break
			}
			ex := &examiner{stairs: true}
			err = ex.Action(g, stair)
			if err == nil && !g.MoveToTarget(g.Ev) {
				err = errors.New("You could not move toward stairs.")
			}
			if ex.Done() {
				g.Targeting = InvalidPos
			}
		} else {
			err = errors.New("You cannot go to any stairs.")
		}
	case KeyEquip:
		err = ui.Equip(g, g.Ev)
		ui.MenuSelectedAnimation(g, MenuInteract, err == nil)
	case KeyDrink:
		err = ui.SelectPotion(g, g.Ev)
		err = ui.CleanError(err)
	case KeyThrow:
		err = ui.SelectProjectile(g, g.Ev)
		err = ui.CleanError(err)
	case KeyEvoke:
		err = ui.SelectRod(g, g.Ev)
		err = ui.CleanError(err)
	case KeyExplore:
		err = g.Autoexplore(g.Ev)
		ui.MenuSelectedAnimation(g, MenuExplore, err == nil)
	case KeyExamine:
		err, again, quit = ui.Examine(g, nil)
	case KeyHelp, KeyMenuCommandHelp:
		ui.KeysHelp(g)
		again = true
	case KeyMenuTargetingHelp:
		ui.ExamineHelp(g)
		again = true
	case KeyCharacterInfo:
		ui.CharacterInfo(g)
		again = true
	case KeyLogs:
		ui.DrawPreviousLogs(g)
		again = true
	case KeySave:
		g.Ev.Renew(g, 0)
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save game.", logError)
		} else {
			quit = true
		}
	case KeyDump:
		errdump := g.WriteDump()
		if errdump != nil {
			g.PrintfStyled("Error: %v", logError, errdump)
			g.PrintStyled("Could not write character dump.", logError)
		} else {
			dataDir, _ := g.DataDir()
			if dataDir != "" {
				g.Printf("Dump written to %s.", filepath.Join(dataDir, "dump"))
			} else {
				g.Print("Dump written.")
			}
		}
		again = true
	case KeyWizardInfo:
		if g.Wizard {
			err = ui.HandleWizardAction(g)
			again = true
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case KeyWizard:
		ui.EnterWizard(g)
		return nil, true, false
	case KeyQuit:
		if ui.Quit(g) {
			return nil, false, true
		}
		return nil, true, false
	case KeyConfigure:
		err = ui.HandleSettingAction(g)
		again = true
	case KeyDescription:
		//ui.MenuSelectedAnimation(g, MenuView, false)
		err = fmt.Errorf("You must choose a target to describe.")
	case KeyExclude:
		err = fmt.Errorf("You must choose a target for exclusion.")
	default:
		err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
	}
	if err != nil {
		again = true
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

func (ui *termui) ExaminePos(g *game, ev event, pos position) (err error, again, quit bool) {
	var start *position
	if pos.valid() {
		start = &pos
	}
	err, again, quit = ui.Examine(g, start)
	return err, again, quit
}

func (ui *termui) DrawKeysDescription(g *game, title string, actions []string) {
	ui.DrawDungeonView(g, NoFlushMode)

	if CustomKeys {
		ui.DrawStyledTextLine(fmt.Sprintf(" Default %s ", title), 0, HeaderLine)
	} else {
		ui.DrawStyledTextLine(fmt.Sprintf(" %s ", title), 0, HeaderLine)
	}
	for i := 0; i < len(actions)-1; i += 2 {
		bg := ui.ListItemBG(i / 2)
		ui.ClearLineWithColor(i/2+1, bg)
		ui.DrawColoredTextOnBG(fmt.Sprintf(" %-36s %s", actions[i], actions[i+1]), 0, i/2+1, ColorFg, bg)
	}
	lines := 1 + len(actions)/2
	ui.DrawTextLine(" press esc or space to continue ", lines)
	ui.Flush()

	ui.WaitForContinue(g, lines)
}

func (ui *termui) KeysHelp(g *game) {
	ui.DrawKeysDescription(g, "Commands", []string{
		"Movement", "h/j/k/l/y/u/b/n or numpad or mouse left",
		"Wait a turn", "“.” or 5 or mouse left on @",
		"Rest (until status free or regen)", "r",
		"Descend stairs", "> or D",
		"Go to nearest stairs", "G",
		"Autoexplore", "o",
		"Examine", "x or mouse left",
		"Equip/Get weapon/armour/...", "e or g",
		"Quaff/Drink potion", "q or d",
		"Throw/Fire item", "t or f",
		"Evoke/Zap rod", "v or z",
		"View Character and Quest Information", `% or C`,
		"View previous messages", "m",
		"Write game statistics to file", "#",
		"Save and Quit", "S",
		"Quit without saving", "Q",
		"Change settings and key bindings", "=",
	})
}

func (ui *termui) ExamineHelp(g *game) {
	ui.DrawKeysDescription(g, "Examine/Travel/Targeting Commands", []string{
		"Move cursor", "h/j/k/l/y/u/b/n or numpad or mouse left",
		"Cycle through monsters", "+",
		"Cycle through stairs", ">",
		"Cycle through objects", "o",
		"Go to/select target", "“.” or enter or mouse left",
		"View target description", "v or d or mouse right",
		"Toggle exclude area from auto-travel", "e or mouse middle",
	})
}

func (ui *termui) Equip(g *game, ev event) error {
	return g.Equip(ev)
}

const TextWidth = DungeonWidth - 2

func (ui *termui) CharacterInfo(g *game) {
	ui.DrawDungeonView(g, NoFlushMode)

	b := bytes.Buffer{}
	b.WriteString(formatText("Every year, the elders send someone to collect medicinal simella plants in the Underground.  This year, the honor fell upon you, and so here you are.  According to the elders, deep in the Underground, magical stairs will lead you back to your village.  Along the way, you will collect simellas, as well as various items that will help you deal with monsters, which you may fight or flee...", TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(
		fmt.Sprintf("You are wielding %s. %s", Indefinite(g.Player.Weapon.String(), false), g.Player.Weapon.Desc()), TextWidth))
	b.WriteString("\n\n")
	b.WriteString(formatText(fmt.Sprintf("You are wearing %s. %s", g.Player.Armour.StringIndefinite(), g.Player.Armour.Desc()), TextWidth))
	b.WriteString("\n\n")
	if g.Player.Shield != NoShield {
		b.WriteString(formatText(fmt.Sprintf("You are wearing a %s. %s", g.Player.Shield, g.Player.Shield.Desc()), TextWidth))
		b.WriteString("\n\n")
	}
	b.WriteString(ui.AptitudesText(g))

	desc := b.String()
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		if i >= DungeonWidth {
			ui.SetCell(DungeonWidth, i, '│', ColorFg, ColorBg)
		}
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	escspace := " press esc or space to continue "
	if lines+2 >= DungeonHeight {
		ui.DrawTextLine(escspace, lines+2)
		ui.SetCell(DungeonWidth, lines+2, '┘', ColorFg, ColorBg)
	} else {
		ui.DrawTextLine(escspace, lines+2)
	}

	ui.Flush()
	ui.WaitForContinue(g, lines+2)
}

func (ui *termui) WizardInfo(g *game) {
	ui.Clear()
	b := &bytes.Buffer{}
	fmt.Fprintf(b, "Monsters: %d (%d)\n", len(g.Monsters), g.MaxMonsters())
	fmt.Fprintf(b, "Danger: %d (%d)\n", g.Danger(), g.MaxDanger())
	ui.DrawText(b.String(), 0, 0)
	ui.Flush()
	ui.WaitForContinue(g, -1)
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

func (ui *termui) AddComma(see, s string) string {
	if len(s) > 0 {
		return s + ", "
	}
	return fmt.Sprintf("You %s %s", see, s)
}
func (ui *termui) DescribePosition(g *game, pos position, targ Targeter) {
	var desc string
	switch {
	case !g.Dungeon.Cell(pos).Explored:
		desc = "You do not know what is in there."
		g.InfoEntry = desc
		return
	case !targ.Reachable(g, pos):
		desc = "This is out of reach."
		g.InfoEntry = desc
		return
	}
	mons := g.MonsterAt(pos)
	c, okCollectable := g.Collectables[pos]
	eq, okEq := g.Equipables[pos]
	rod, okRod := g.Rods[pos]
	if pos == g.Player.Pos {
		desc = "This is you"
	}
	see := "see"
	if !g.Player.LOS[pos] {
		see = "saw"
	}
	if g.Dungeon.Cell(pos).T == WallCell && !g.WrongWall[pos] || g.Dungeon.Cell(pos).T == FreeCell && g.WrongWall[pos] {
		desc = ui.AddComma(see, "")
		desc += fmt.Sprintf("a wall")
		g.InfoEntry = desc + "."
		return
	}
	if mons.Exists() && g.Player.LOS[pos] {
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("%s (%s)", mons.Kind.Indefinite(false), ui.MonsterInfo(mons))
	}
	strt, okStair := g.Stairs[pos]
	stn, okStone := g.MagicalStones[pos]
	switch {
	case g.Simellas[pos] > 0:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("some simellas (%d)", g.Simellas[pos])
	case okCollectable:
		if c.Quantity > 1 {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("%d %s", c.Quantity, c.Consumable)
		} else {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("%s", Indefinite(c.Consumable.String(), false))
		}
	case okEq:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("%s", Indefinite(eq.String(), false))
	case okRod:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("a %v", rod)
	case okStair:
		if strt == WinStair {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("glowing stairs")
		} else {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("stairs downwards")
		}
	case okStone:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprint(Indefinite(stn.String(), false))
	case g.Doors[pos] || g.WrongDoor[pos]:
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("a door")
	}
	if cld, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
		if cld == CloudFire {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("burning flames")
		} else if cld == CloudNight {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("night clouds")
		} else {
			desc = ui.AddComma(see, desc)
			desc += fmt.Sprintf("a dense fog")
		}
	} else if _, ok := g.Fungus[pos]; ok && !g.WrongFoliage[pos] || !ok && g.WrongFoliage[pos] {
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("foliage")
	} else if desc == "" {
		desc = ui.AddComma(see, desc)
		desc += fmt.Sprintf("the ground")
	}
	g.InfoEntry = desc + "."
}

func (ui *termui) Examine(g *game, start *position) (err error, again, quit bool) {
	ex := &examiner{}
	err, again, quit = ui.CursorAction(g, ex, start)
	return err, again, quit
}

func (ui *termui) ChooseTarget(g *game, targ Targeter) error {
	err, _, _ := ui.CursorAction(g, targ, nil)
	if err != nil {
		return err
	}
	if !targ.Done() {
		return errors.New(DoNothing)
	}
	return nil
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
		for p := range g.MagicalStones {
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

func (ui *termui) CursorMouseLeft(g *game, targ Targeter, pos position, data *examineData) (again, notarg bool) {
	again = true
	if data.npos == pos {
		err := targ.Action(g, pos)
		if err != nil {
			g.Print(err.Error())
		} else {
			if g.MoveToTarget(g.Ev) {
				again = false
			}
			if targ.Done() {
				notarg = true
			}
		}
	} else {
		data.npos = pos
	}
	return again, notarg
}

func (ui *termui) CursorKeyAction(g *game, targ Targeter, rka runeKeyAction, data *examineData) (err error, again, quit, notarg bool) {
	pos := data.npos
	again = true
	if rka.r != 0 {
		var ok bool
		rka.k, ok = gameConfig.RuneTargetModeKeys[rka.r]
		if !ok {
			err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
			return err, again, quit, notarg
		}
	}
	if rka.k == KeyMenu {
		rka.k, err = ui.SelectAction(g, menuActions, g.Ev)
		if err != nil {
			err = ui.CleanError(err)
			return err, again, quit, notarg
		}
	}
	switch rka.k {
	case KeyW, KeyS, KeyN, KeyE, KeyNW, KeyNE, KeySW, KeySE:
		data.npos = pos.To(KeyToDir(rka.k))
	case KeyRunW, KeyRunS, KeyRunN, KeyRunE, KeyRunNW, KeyRunNE, KeyRunSW, KeyRunSE:
		for i := 0; i < 5; i++ {
			p := data.npos.To(KeyToDir(rka.k))
			if !p.valid() {
				break
			}
			data.npos = p
		}
	case KeyNextStairs:
		ui.NextStair(g, data)
	case KeyDescend:
		if strt, ok := g.Stairs[g.Player.Pos]; ok {
			ui.MenuSelectedAnimation(g, MenuInteract, true)
			err = ui.OptionalDescendConfirmation(g, strt)
			if err != nil {
				break
			}
			again = false
			g.Targeting = InvalidPos
			notarg = true
			if g.Descend() {
				ui.Win(g)
				quit = true
				return err, again, quit, notarg
			}
		} else {
			err = errors.New("No stairs here.")
		}
	case KeyPreviousMonster, KeyNextMonster:
		ui.NextMonster(g, rka.r, pos, data)
	case KeyNextObject:
		ui.NextObject(g, pos, data)
	case KeyHelp, KeyMenuTargetingHelp:
		ui.HideCursor()
		ui.ExamineHelp(g)
		ui.SetCursor(pos)
	case KeyMenuCommandHelp:
		ui.HideCursor()
		ui.KeysHelp(g)
		ui.SetCursor(pos)
	case KeyTarget:
		err = targ.Action(g, pos)
		if err != nil {
			break
		}
		g.Targeting = InvalidPos
		if g.MoveToTarget(g.Ev) {
			again = false
		}
		if targ.Done() {
			notarg = true
		}
	case KeyDescription:
		ui.HideCursor()
		ui.ViewPositionDescription(g, pos)
		ui.SetCursor(pos)
	case KeyExclude:
		ui.ExcludeZone(g, pos)
	case KeyEscape:
		g.Targeting = InvalidPos
		notarg = true
		err = errors.New(DoNothing)
	case KeyExplore, KeyRest, KeyThrow, KeyDrink, KeyEvoke, KeyLogs, KeyEquip, KeyCharacterInfo:
		if _, ok := targ.(*examiner); !ok {
			break
		}
		err, again, quit = ui.HandleKey(g, rka)
		if err != nil {
			notarg = true
		}
		g.Targeting = InvalidPos
	case KeyConfigure:
		err = ui.HandleSettingAction(g)
	case KeySave:
		g.Ev.Renew(g, 0)
		g.Highlight = nil
		g.Targeting = InvalidPos
		errsave := g.Save()
		if errsave != nil {
			g.PrintfStyled("Error: %v", logError, errsave)
			g.PrintStyled("Could not save game.", logError)
		} else {
			notarg = true
			again = false
			quit = true
		}
	case KeyQuit:
		if ui.Quit(g) {
			quit = true
			again = false
		}
	default:
		err = fmt.Errorf("Invalid targeting mode key '%c'. Type ? for help.", rka.r)
	}
	return err, again, quit, notarg
}

type examineData struct {
	npos         position
	nmonster     int
	objects      []position
	nobject      int
	sortedStairs []position
	stairIndex   int
}

var InvalidPos = position{-1, -1}

func (ui *termui) CursorAction(g *game, targ Targeter, start *position) (err error, again, quit bool) {
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
	data := &examineData{
		npos:    pos,
		objects: []position{},
	}
	if _, ok := targ.(*examiner); ok && pos == g.Player.Pos && start == nil {
		ui.NextObject(g, InvalidPos, data)
		if !data.npos.valid() {
			ui.NextStair(g, data)
		}
		if data.npos.valid() {
			pos = data.npos
		}
	}
	opos := InvalidPos
loop:
	for {
		err = nil
		if pos != opos {
			ui.DescribePosition(g, pos, targ)
		}
		opos = pos
		targ.ComputeHighlight(g, pos)
		ui.SetCursor(pos)
		ui.DrawDungeonView(g, TargetingMode)
		ui.DrawInfoLine(g.InfoEntry)
		if !ui.Small() {
			st := " Examine/Travel (? for help) "
			if _, ok := targ.(*examiner); !ok {
				st = " Targeting (? for help) "
			}
			ui.DrawStyledTextLine(st, DungeonHeight+2, FooterLine)
		}
		ui.SetCell(DungeonWidth, DungeonHeight, '┤', ColorFg, ColorBg)
		ui.Flush()
		data.npos = pos
		var notarg bool
		err, again, quit, notarg = ui.TargetModeEvent(g, targ, data)
		if err != nil {
			err = ui.CleanError(err)
		}
		if !again || notarg {
			break loop
		}
		if err != nil {
			g.Print(err.Error())
		}
		if data.npos.valid() {
			pos = data.npos
		}
	}
	g.Highlight = nil
	ui.HideCursor()
	return err, again, quit
}

func (ui *termui) ViewPositionDescription(g *game, pos position) {
	if !g.Dungeon.Cell(pos).Explored {
		ui.DrawDescription(g, "This place is unknown to you.")
		return
	}
	mons := g.MonsterAt(pos)
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
	} else if strt, ok := g.Stairs[pos]; ok {
		if strt == WinStair {
			desc := "These shiny-looking stairs are in fact a magical monolith. It is said they were made some centuries ago by Marevor Helith. They will lead you back to your village."
			if g.Depth < MaxDepth {
				desc += " Note that this is not the last floor, so you may want to find a normal stair and continue collecting simellas, if you're courageous enough."
			}
			ui.DrawDescription(g, desc)
		} else {
			desc := "Stairs lead to the next level of the Underground. There's no way back. Monsters do not follow you."
			if g.Depth == WinDepth {
				desc += " If you're afraid, you could instead just win by taking the magical stairs somewhere in the same map."
			}
			ui.DrawDescription(g, desc)
		}
	} else if stn, ok := g.MagicalStones[pos]; ok {
		ui.DrawDescription(g, stn.Description())
	} else if g.Doors[pos] {
		ui.DrawDescription(g, "A closed door blocks your line of sight. Doors open automatically when you or a monster stand on them. Doors are flammable.")
	} else if g.Simellas[pos] > 0 {
		ui.DrawDescription(g, "A simella is a plant with big white flowers which are used in the Underground for their medicinal properties. They can also make tasty infusions. You were actually sent here by your village to collect as many as possible of those plants.")
	} else if _, ok := g.Fungus[pos]; ok && g.Dungeon.Cell(pos).T == FreeCell {
		ui.DrawDescription(g, "Blue dense foliage grows in the Underground. It is difficult to see through, and is flammable.")
	} else if g.Dungeon.Cell(pos).T == WallCell {
		ui.DrawDescription(g, "A wall is an impassable pile of rocks. It can be destructed by using some items.")
	} else {
		ui.DrawDescription(g, "This is just plain ground.")
	}
}

func (ui *termui) MonsterInfo(m *monster) string {
	infos := []string{}
	state := m.State.String()
	if m.Kind == MonsSatowalgaPlant && m.State == Wandering {
		state = "awaken"
	}
	infos = append(infos, state)
	for st, i := range m.Statuses {
		if i > 0 {
			infos = append(infos, monsterStatus(st).String())
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
		ui.SetMapCell(x, y, r, fg, bg)
		if ui.InViewBorder(g, pos, targeting) && g.Dungeon.Border(pos) {
			for _, opos := range pos.OutsideNeighbors() {
				xo, yo := ui.CameraOffset(g, opos, targeting)
				ui.SetMapCell(xo, yo, '#', ColorFg, ColorBgBorder)
			}
		}
		return
	}
	ui.SetMapCell(pos.X, pos.Y, r, fg, bg)
}

const BarCol = DungeonWidth + 2

func (ui *termui) DrawDungeonView(g *game, m uiMode) {
	ui.Clear()
	d := g.Dungeon
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, DungeonHeight, '─', ColorFg, ColorBg)
	}
	for i := 0; i < DungeonHeight; i++ {
		ui.SetCell(DungeonWidth, i, '│', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, DungeonHeight, '┘', ColorFg, ColorBg)
	for i := range d.Cells {
		pos := idxtopos(i)
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, m == TargetingMode, r, fgColor, bgColor)
	}
	line := 0
	if !ui.Small() {
		ui.SetMapCell(BarCol, line, '[', ColorFg, ColorBg)
		ui.DrawText(fmt.Sprintf(" %v", g.Player.Armour), BarCol+1, line)
		line++
		ui.SetMapCell(BarCol, line, ')', ColorFg, ColorBg)
		ui.DrawText(fmt.Sprintf(" %v", g.Player.Weapon), BarCol+1, line)
		line++
		if g.Player.Shield != NoShield {
			if g.Player.Weapon.TwoHanded() {
				ui.SetMapCell(BarCol, line, ']', ColorFg, ColorBg)
				ui.DrawText(" (unusable)", BarCol+1, line)
			} else {
				ui.SetMapCell(BarCol, line, ']', ColorFg, ColorBg)
				ui.DrawText(fmt.Sprintf(" %v", g.Player.Shield), BarCol+1, line)
			}
		}
		line++
		line++
	}
	if ui.Small() {
		ui.DrawStatusLine(g)
	} else {
		ui.DrawStatusBar(g, line)
		ui.DrawMenus(g)
	}
	if ui.Small() {
		ui.DrawLog(g, 2)
	} else {
		ui.DrawLog(g, 4)
	}
	if m != TargetingMode && m != NoFlushMode {
		ui.Flush()
	}
}

func (ui *termui) SwappingAnimation(g *game, mpos, ppos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	_, fgm, bgColorm := ui.PositionDrawing(g, mpos)
	_, _, bgColorp := ui.PositionDrawing(g, ppos)
	ui.DrawAtPosition(g, mpos, true, 'Φ', fgm, bgColorp)
	ui.DrawAtPosition(g, ppos, true, 'Φ', ColorFgPlayer, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, mpos, true, 'Φ', ColorFgPlayer, bgColorp)
	ui.DrawAtPosition(g, ppos, true, 'Φ', fgm, bgColorm)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
}

func (ui *termui) TeleportAnimation(g *game, from, to position, showto bool) {
	if DisableAnimations {
		return
	}
	_, _, bgColorf := ui.PositionDrawing(g, from)
	_, _, bgColort := ui.PositionDrawing(g, to)
	ui.DrawAtPosition(g, from, true, 'Φ', ColorCyan, bgColorf)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	if showto {
		ui.DrawAtPosition(g, from, true, 'Φ', ColorBlue, bgColorf)
		ui.DrawAtPosition(g, to, true, 'Φ', ColorCyan, bgColort)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
	}
}

type explosionStyle int

const (
	FireExplosion explosionStyle = iota
	WallExplosion
	AroundWallExplosion
)

func (ui *termui) ProjectileTrajectoryAnimation(g *game, ray []position, fg uicolor) {
	if DisableAnimations {
		return
	}
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, '•', fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
}

func (ui *termui) MonsterProjectileAnimation(g *game, ray []position, r rune, fg uicolor) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		or, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, or, fgColor, bgColor)
	}
}

func (ui *termui) ExplosionAnimationAt(g *game, pos position, fg uicolor) {
	_, _, bgColor := ui.PositionDrawing(g, pos)
	mons := g.MonsterAt(pos)
	r := ';'
	switch RandInt(9) {
	case 0, 6:
		r = ','
	case 1:
		r = '}'
	case 2:
		r = '%'
	case 3, 7:
		r = ':'
	case 4:
		r = '\\'
	case 5:
		r = '~'
	}
	if mons.Exists() || g.Player.Pos == pos {
		r = '√'
	}
	//ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
	ui.DrawAtPosition(g, pos, true, r, bgColor, fg)
}

func (ui *termui) ExplosionAnimation(g *game, es explosionStyle, pos position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	if es == WallExplosion || es == AroundWallExplosion {
		colors[0] = ColorFgExplosionWallStart
		colors[1] = ColorFgExplosionWallEnd
	}
	for i := 0; i < 3; i++ {
		nb := g.Dungeon.FreeNeighbors(pos)
		if es != AroundWallExplosion {
			nb = append(nb, pos)
		}
		for _, npos := range nb {
			fg := colors[RandInt(2)]
			if !g.Player.LOS[npos] {
				continue
			}
			ui.ExplosionAnimationAt(g, npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

func (ui *termui) TormentExplosionAnimation(g *game) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(20 * time.Millisecond)
	colors := [3]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd, ColorFgMagicPlace}
	for i := 0; i < 3; i++ {
		for npos, b := range g.Player.LOS {
			if !b {
				continue
			}
			fg := colors[RandInt(3)]
			ui.ExplosionAnimationAt(g, npos, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(20 * time.Millisecond)
}

func (ui *termui) WallExplosionAnimation(g *game, pos position) {
	if DisableAnimations {
		return
	}
	colors := [2]uicolor{ColorFgExplosionWallStart, ColorFgExplosionWallEnd}
	for _, fg := range colors {
		_, _, bgColor := ui.PositionDrawing(g, pos)
		//ui.DrawAtPosition(g, pos, true, '☼', fg, bgColor)
		ui.DrawAtPosition(g, pos, true, '☼', bgColor, fg)
		ui.Flush()
		time.Sleep(25 * time.Millisecond)
	}
}

func (ui *termui) LightningBoltAnimation(g *game, ray []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 3; j++ {
		for i := len(ray) - 1; i >= 0; i-- {
			fg := colors[RandInt(2)]
			pos := ray[i]
			_, _, bgColor := ui.PositionDrawing(g, pos)
			mons := g.MonsterAt(pos)
			r := '*'
			if RandInt(2) == 0 {
				r = '×'
			}
			if mons.Exists() {
				r = '√'
			}
			//ui.DrawAtPosition(g, pos, true, r, fg, bgColor)
			ui.DrawAtPosition(g, pos, true, r, bgColor, fg)
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
	time.Sleep(25 * time.Millisecond)
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
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := len(ray) - 1; i >= 0; i-- {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgProjectile, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
	if hit {
		pos := ray[0]
		ui.HitAnimation(g, pos, true)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *termui) MonsterJavelinAnimation(g *game, ray []position, hit bool) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	for i := 0; i < len(ray); i++ {
		pos := ray[i]
		r, fgColor, bgColor := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, true, ui.ProjectileSymbol(pos.Dir(g.Player.Pos)), ColorFgMonster, bgColor)
		ui.Flush()
		time.Sleep(30 * time.Millisecond)
		ui.DrawAtPosition(g, pos, true, r, fgColor, bgColor)
	}
	time.Sleep(30 * time.Millisecond)
}

func (ui *termui) HitAnimation(g *game, pos position, targeting bool) {
	if DisableAnimations {
		return
	}
	if !g.Player.LOS[pos] {
		return
	}
	ui.DrawDungeonView(g, NoFlushMode)
	_, _, bgColor := ui.PositionDrawing(g, pos)
	mons := g.MonsterAt(pos)
	if mons.Exists() || pos == g.Player.Pos {
		ui.DrawAtPosition(g, pos, targeting, '√', ColorFgAnimationHit, bgColor)
	} else {
		ui.DrawAtPosition(g, pos, targeting, '∞', ColorFgAnimationHit, bgColor)
	}
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
}

func (ui *termui) LightningHitAnimation(g *game, targets []position) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(25 * time.Millisecond)
	colors := [2]uicolor{ColorFgExplosionStart, ColorFgExplosionEnd}
	for j := 0; j < 2; j++ {
		for _, pos := range targets {
			_, _, bgColor := ui.PositionDrawing(g, pos)
			mons := g.MonsterAt(pos)
			if mons.Exists() || pos == g.Player.Pos {
				ui.DrawAtPosition(g, pos, false, '√', bgColor, colors[RandInt(2)])
			} else {
				ui.DrawAtPosition(g, pos, false, '∞', bgColor, colors[RandInt(2)])
			}
		}
		ui.Flush()
		time.Sleep(100 * time.Millisecond)
	}
}

var WizardMap = false

func (ui *termui) PositionDrawing(g *game, pos position) (r rune, fgColor, bgColor uicolor) {
	m := g.Dungeon
	c := m.Cell(pos)
	fgColor = ColorFg
	bgColor = ColorBg
	if !c.Explored && !g.Wizard {
		r = ' '
		bgColor = ColorBgDark
		if g.HasFreeExploredNeighbor(pos) {
			r = '¤'
			fgColor = ColorFgDark
		}
		if g.DreamingMonster[pos] {
			r = '☻'
			fgColor = ColorFgSleepingMonster
		}
		if g.Noise[pos] {
			r = '♫'
			fgColor = ColorFgWanderingMonster
		}
		return
	}
	if g.Wizard {
		if !c.Explored && g.HasFreeExploredNeighbor(pos) && !WizardMap {
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
	if g.Player.LOS[pos] && !WizardMap {
		fgColor = ColorFgLOS
		bgColor = ColorBgLOS
		if pos.X%2 == 1 && pos.Y%2 == 0 {
			bgColor = ColorBgLOSalt
		} else if pos.X%2 == 0 && pos.Y%2 == 1 {
			bgColor = ColorBgLOSalt
		}
	} else {
		fgColor = ColorFgDark
		bgColor = ColorBgDark
	}
	if g.ExclusionsMap[pos] && c.T != WallCell {
		fgColor = ColorFgExcluded
	}
	switch {
	case c.T == WallCell && (!g.WrongWall[pos] || g.Wizard) || c.T == FreeCell && g.WrongWall[pos] && !g.Wizard:
		r = '#'
		if g.TemporalWalls[pos] {
			fgColor = ColorFgMagicPlace
		}
	case pos == g.Player.Pos && !WizardMap:
		r = '@'
		fgColor = ColorFgPlayer
	default:
		r = '.'
		if _, ok := g.Fungus[pos]; ok && !g.WrongFoliage[pos] || !ok && g.WrongFoliage[pos] {
			r = '"'
		}
		if cld, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
			r = '§'
			if cld == CloudFire {
				fgColor = ColorFgWanderingMonster
			} else if cld == CloudNight {
				fgColor = ColorFgSleepingMonster
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
		} else if strt, ok := g.Stairs[pos]; ok {
			r = '>'
			if strt == WinStair {
				fgColor = ColorFgMagicPlace
			} else {
				fgColor = ColorFgPlace
			}
		} else if stn, ok := g.MagicalStones[pos]; ok {
			r = '_'
			if stn == InertStone {
				fgColor = ColorFgPlace
			} else {
				fgColor = ColorFgMagicPlace
			}
		} else if _, ok := g.Simellas[pos]; ok {
			r = '♣'
			fgColor = ColorFgSimellas
		} else if _, ok := g.Doors[pos]; ok {
			r = '+'
			fgColor = ColorFgPlace
		}
		if (g.Player.LOS[pos] || g.Wizard) && !WizardMap {
			m := g.MonsterAt(pos)
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
		} else if !g.Wizard && g.DreamingMonster[pos] {
			r = '☻'
			fgColor = ColorFgSleepingMonster
		}
	}
	return
}

func (ui *termui) DrawStatusBar(g *game, line int) {
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
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
	ui.DrawColoredText(fmt.Sprintf("HP: %d", g.Player.HP), BarCol, line, hpColor)
	line++
	ui.DrawColoredText(fmt.Sprintf("MP: %d", g.Player.MP), BarCol, line, mpColor)
	line++
	line++
	ui.DrawText(fmt.Sprintf("Simellas: %d", g.Player.Simellas), BarCol, line)
	line++
	if g.Depth == -1 {
		ui.DrawText("Depth: Out!", BarCol, line)
	} else {
		ui.DrawText(fmt.Sprintf("Depth: %d", g.Depth), BarCol, line)
	}
	line++
	ui.DrawText(fmt.Sprintf("Turns: %.1f", float64(g.Turn)/10), BarCol, line)
	line++
	for _, st := range sts {
		fg := ColorFgStatusOther
		if st.Good() {
			fg = ColorFgStatusGood
			t := 13
			if g.Player.Statuses[StatusBerserk] > 0 {
				t -= 3
			}
			if g.Player.Statuses[StatusSlow] > 0 {
				t += 3
			}
			if g.Player.Expire[st] >= g.Ev.Rank() && g.Player.Expire[st]-g.Ev.Rank() <= t {
				fg = ColorFgStatusExpire
			}
		} else if st.Bad() {
			fg = ColorFgStatusBad
		}
		if g.Player.Statuses[st] > 1 {
			ui.DrawColoredText(fmt.Sprintf("%s(%d)", st, g.Player.Statuses[st]), BarCol, line, fg)
		} else {
			ui.DrawColoredText(st.String(), BarCol, line, fg)
		}
		line++
	}
}

func (ui *termui) DrawStatusLine(g *game) {
	sts := statusSlice{}
	if cld, ok := g.Clouds[g.Player.Pos]; ok && cld == CloudFire {
		g.Player.Statuses[StatusFlames] = 1
		defer func() {
			g.Player.Statuses[StatusFlames] = 0
		}()
	}
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
	line := DungeonHeight
	col := 2
	ui.DrawText(" ", col, line)
	col++
	ui.SetMapCell(col, line, ')', ColorFg, ColorBg)
	col++
	weapon := fmt.Sprintf("%s ", g.Player.Weapon.Short())
	ui.DrawText(weapon, col, line)
	col += utf8.RuneCountInString(weapon)
	ui.SetMapCell(col, line, '[', ColorFg, ColorBg)
	col++
	armour := fmt.Sprintf("%s ", g.Player.Armour.Short())
	ui.DrawText(armour, col, line)
	col += utf8.RuneCountInString(armour)
	if g.Player.Shield != NoShield {
		ui.SetMapCell(col, line, ']', ColorFg, ColorBg)
		col++
		shield := fmt.Sprintf("%s ", g.Player.Shield.Short())
		ui.DrawText(shield, col, line)
		col += utf8.RuneCountInString(shield)
	}
	ui.SetMapCell(col, line, '♣', ColorFg, ColorBg)
	col++
	simellas := fmt.Sprintf(":%d ", g.Player.Simellas)
	ui.DrawText(simellas, col, line)
	col += utf8.RuneCountInString(simellas)
	var depth string
	if g.Depth == -1 {
		depth = "D: Out! "
	} else {
		depth = fmt.Sprintf("D:%d ", g.Depth)
	}
	ui.DrawText(depth, col, line)
	col += utf8.RuneCountInString(depth)
	turns := fmt.Sprintf("T:%.1f ", float64(g.Turn)/10)
	ui.DrawText(turns, col, line)
	col += utf8.RuneCountInString(turns)
	hp := fmt.Sprintf("HP:%2d ", g.Player.HP)
	ui.DrawColoredText(hp, col, line, hpColor)
	col += utf8.RuneCountInString(hp)
	mp := fmt.Sprintf("MP:%d ", g.Player.MP)
	ui.DrawColoredText(mp, col, line, mpColor)
	col += utf8.RuneCountInString(mp)
	if len(sts) > 0 {
		ui.DrawText("| ", col, line)
		col += 2
	}
	for _, st := range sts {
		fg := ColorFgStatusOther
		if st.Good() {
			fg = ColorFgStatusGood
		} else if st.Bad() {
			fg = ColorFgStatusBad
		}
		var sttext string
		if g.Player.Statuses[st] > 1 {
			sttext = fmt.Sprintf("%s(%d) ", st.Short(), g.Player.Statuses[st])
		} else {
			sttext = fmt.Sprintf("%s ", st.Short())
		}
		ui.DrawColoredText(sttext, col, line, fg)
		col += utf8.RuneCountInString(sttext)
	}
}

type menu int

const (
	MenuRest menu = iota
	MenuExplore
	MenuThrow
	MenuDrink
	MenuEvoke
	MenuOther
	MenuInteract
)

func (m menu) String() (text string) {
	switch m {
	case MenuRest:
		text = "rest"
	case MenuExplore:
		text = "explore"
	case MenuThrow:
		text = "throw"
	case MenuDrink:
		text = "drink"
	case MenuEvoke:
		text = "evoke"
	case MenuOther:
		text = "menu"
	case MenuInteract:
		text = "interact"
	}
	return "[" + text + "]"
}

func (m menu) Key(g *game) (key keyAction) {
	switch m {
	case MenuRest:
		key = KeyRest
	case MenuExplore:
		key = KeyExplore
	case MenuThrow:
		key = KeyThrow
	case MenuDrink:
		key = KeyDrink
	case MenuEvoke:
		key = KeyEvoke
	case MenuOther:
		key = KeyMenu
	case MenuInteract:
		if _, ok := g.Equipables[g.Player.Pos]; ok {
			key = KeyEquip
		} else if _, ok := g.Stairs[g.Player.Pos]; ok {
			key = KeyDescend
		}
	}
	return key
}

var MenuCols = [][2]int{
	MenuRest:     {0, 0},
	MenuExplore:  {0, 0},
	MenuThrow:    {0, 0},
	MenuDrink:    {0, 0},
	MenuEvoke:    {0, 0},
	MenuOther:    {0, 0},
	MenuInteract: {0, 0}}

func init() {
	for i := range MenuCols {
		runes := utf8.RuneCountInString(menu(i).String())
		if i == 0 {
			MenuCols[0] = [2]int{7, 7 + runes}
			continue
		}
		MenuCols[i] = [2]int{MenuCols[i-1][1] + 2, MenuCols[i-1][1] + 2 + runes}
	}
}

func (ui *termui) WhichButton(g *game, col int) (menu, bool) {
	if ui.Small() {
		return MenuOther, false
	}
	end := len(MenuCols) - 1
	if _, ok := g.Equipables[g.Player.Pos]; ok {
		end++
	} else if _, ok := g.Stairs[g.Player.Pos]; ok {
		end++
	}
	for i, cols := range MenuCols[0:end] {
		if cols[0] >= 0 && col >= cols[0] && col < cols[1] {
			return menu(i), true
		}
	}
	return MenuOther, false
}

func (ui *termui) UpdateInteractButton(g *game) string {
	var interactMenu string
	var show bool
	if _, ok := g.Equipables[g.Player.Pos]; ok {
		interactMenu = "[equip]"
		show = true
	} else if _, ok := g.Stairs[g.Player.Pos]; ok {
		interactMenu = "[descend]"
		show = true
	}
	if !show {
		return ""
	}
	i := len(MenuCols) - 1
	runes := utf8.RuneCountInString(interactMenu)
	MenuCols[i][1] = MenuCols[i][0] + runes
	return interactMenu
}

func (ui *termui) LogColor(e logEntry) uicolor {
	fg := ColorFg
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

func (ui *termui) DrawLog(g *game, lines int) {
	min := len(g.Log) - lines
	if min < 0 {
		min = 0
	}
	l := len(g.Log) - 1
	if l < lines {
		lines = l + 1
	}
	for i := lines; i > 0 && l >= 0; i-- {
		cols := 0
		first := true
		to := l
		for l >= 0 {
			e := g.Log[l]
			el := utf8.RuneCountInString(e.String())
			if e.Tick {
				el += 2
			}
			cols += el + 1
			if !first && cols > DungeonWidth {
				l++
				break
			}
			if e.Tick || l <= i {
				break
			}
			first = false
			l--
		}
		if l < 0 {
			l = 0
		}
		col := 0
		for ln := l; ln <= to; ln++ {
			e := g.Log[ln]
			fguicolor := ui.LogColor(e)
			if e.Tick {
				ui.DrawColoredText("•", 0, DungeonHeight+i, ColorYellow)
				col += 2
			}
			ui.DrawColoredText(e.String(), col, DungeonHeight+i, fguicolor)
			col += utf8.RuneCountInString(e.String()) + 1
		}
		l--
	}
}

func InRuneSlice(r rune, s []rune) bool {
	for _, rr := range s {
		if r == rr {
			return true
		}
	}
	return false
}

func (ui *termui) RunesForKeyAction(k keyAction) string {
	runes := []rune{}
	for r, ka := range gameConfig.RuneNormalModeKeys {
		if k == ka && !InRuneSlice(r, runes) {
			runes = append(runes, r)
		}
	}
	for r, ka := range gameConfig.RuneTargetModeKeys {
		if k == ka && !InRuneSlice(r, runes) {
			runes = append(runes, r)
		}
	}
	chars := strings.Split(string(runes), "")
	sort.Strings(chars)
	text := strings.Join(chars, " or ")
	return text
}

type keyConfigAction int

const (
	NavigateKeys keyConfigAction = iota
	ChangeKeys
	ResetKeys
	QuitKeyConfig
)

func (ui *termui) ChangeKeys(g *game) {
	lines := DungeonHeight
	nmax := len(configurableKeyActions) - lines
	n := 0
	s := 0
loop:
	for {
		ui.DrawDungeonView(g, NoFlushMode)
		if n >= nmax {
			n = nmax
		}
		if n < 0 {
			n = 0
		}
		to := n + lines
		if to >= len(configurableKeyActions) {
			to = len(configurableKeyActions)
		}
		for i := n; i < to; i++ {
			ka := configurableKeyActions[i]
			desc := ka.NormalModeDescription()
			if !ka.NormalModeKey() {
				desc = ka.TargetingModeDescription()
			}
			bg := ui.ListItemBG(i)
			ui.ClearLineWithColor(i-n, bg)
			desc = fmt.Sprintf(" %-36s %s", desc, ui.RunesForKeyAction(ka))
			if i == s {
				ui.DrawColoredTextOnBG(desc, 0, i-n, ColorYellow, bg)
			} else {
				ui.DrawColoredTextOnBG(desc, 0, i-n, ColorFg, bg)
			}
		}
		ui.ClearLine(lines)
		ui.DrawStyledTextLine(" add key (a) up/down (arrows/u/d) reset (R) quit (esc or space) ", lines, FooterLine)
		ui.Flush()

		var action keyConfigAction
		s, action = ui.KeyMenuAction(s)
		if s >= len(configurableKeyActions) {
			s = len(configurableKeyActions) - 1
		}
		if s < 0 {
			s = 0
		}
		if s < n+1 {
			n -= 12
		}
		if s > n+lines-2 {
			n += 12
		}
		switch action {
		case ChangeKeys:
			ui.DrawStyledTextLine(" insert new key ", lines, FooterLine)
			ui.Flush()
			r := ui.ReadRuneKey()
			if r == 0 {
				continue loop
			}
			if FixedRuneKey(r) {
				g.Printf("You cannot rebind “%c”.", r)
				continue loop
			}
			CustomKeys = true
			ka := configurableKeyActions[s]
			if ka.NormalModeKey() {
				gameConfig.RuneNormalModeKeys[r] = ka
			} else {
				delete(gameConfig.RuneNormalModeKeys, r)
			}
			if ka.TargetingModeKey() {
				gameConfig.RuneTargetModeKeys[r] = ka
			} else {
				delete(gameConfig.RuneTargetModeKeys, r)
			}
			err := g.SaveConfig()
			if err != nil {
				g.Print(err.Error())
			}
		case QuitKeyConfig:
			break loop
		case ResetKeys:
			ApplyDefaultKeyBindings()
			err := g.SaveConfig()
			//err := g.RemoveDataFile("config.gob")
			if err != nil {
				g.Print(err.Error())
			}
		}
	}
}

func (ui *termui) DrawPreviousLogs(g *game) {
	bottom := 4
	if ui.Small() {
		bottom = 2
	}
	lines := DungeonHeight + bottom
	nmax := len(g.Log) - lines
	n := nmax
loop:
	for {
		ui.DrawDungeonView(g, NoFlushMode)
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
		for i := 0; i < bottom; i++ {
			ui.SetCell(DungeonWidth, DungeonHeight+i, '│', ColorFg, ColorBg)
		}
		for i := n; i < to; i++ {
			e := g.Log[i]
			fguicolor := ui.LogColor(e)
			ui.ClearLine(i - n)
			rc := utf8.RuneCountInString(e.String())
			if e.Tick {
				rc += 2
			}
			if rc >= DungeonWidth {
				for j := DungeonWidth; j < 103; j++ {
					ui.SetCell(j, i-n, ' ', ColorFg, ColorBg)
				}
			}
			if e.Tick {
				ui.DrawColoredText("•", 0, i-n, ColorYellow)
				ui.DrawColoredText(e.String(), 2, i-n, fguicolor)
			} else {
				ui.DrawColoredText(e.String(), 0, i-n, fguicolor)
			}
		}
		for i := len(g.Log); i < DungeonHeight+bottom; i++ {
			ui.ClearLine(i - n)
		}
		ui.ClearLine(lines)
		s := fmt.Sprintf(" half-page up/down (u/d) quit (esc or space) — (%d/%d) \n", len(g.Log)-to, len(g.Log))
		ui.DrawStyledTextLine(s, lines, FooterLine)
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
	ui.DrawDungeonView(g, NoFlushMode)
	desc = formatText(desc, TextWidth)
	lines := strings.Count(desc, "\n")
	for i := 0; i <= lines+2; i++ {
		ui.ClearLine(i)
	}
	ui.DrawText(desc, 0, 0)
	ui.DrawTextLine(" press esc or space to continue ", lines+2)
	ui.Flush()
	ui.WaitForContinue(g, lines+2)
	ui.DrawDungeonView(g, NoFlushMode)
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
		if x+col >= UIWidth {
			break
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
	FooterLine
)

func (ui *termui) DrawInfoLine(text string) {
	ui.ClearLineWithColor(DungeonHeight+1, ColorBase02)
	ui.DrawColoredTextOnBG(text, 0, DungeonHeight+1, ColorBlue, ColorBase02)
}

func (ui *termui) DrawStyledTextLine(text string, lnum int, st linestyle) {
	nchars := utf8.RuneCountInString(text)
	dist := (DungeonWidth - nchars) / 2
	for i := 0; i < dist; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.DrawColoredText(text, dist, lnum, ColorYellow)
	case FooterLine:
		ui.DrawColoredText(text, dist, lnum, ColorCyan)
	default:
		ui.DrawColoredText(text, dist, lnum, ColorFg)
	}
	for i := dist + nchars; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, '─', ColorFg, ColorBg)
	}
	switch st {
	case HeaderLine:
		ui.SetCell(DungeonWidth, lnum, '┐', ColorFg, ColorBg)
	case FooterLine:
		ui.SetCell(DungeonWidth, lnum, '┘', ColorFg, ColorBg)
	default:
		ui.SetCell(DungeonWidth, lnum, '┤', ColorFg, ColorBg)
	}
}

func (ui *termui) ClearLine(lnum int) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, ColorBg)
	}
	ui.SetCell(DungeonWidth, lnum, '│', ColorFg, ColorBg)
}

func (ui *termui) ClearLineWithColor(lnum int, bg uicolor) {
	for i := 0; i < DungeonWidth; i++ {
		ui.SetCell(i, lnum, ' ', ColorFg, bg)
	}
	ui.SetCell(DungeonWidth, lnum, '│', ColorFg, ColorBg)
}

func (ui *termui) ListItemBG(i int) uicolor {
	bg := ColorBase03
	if i%2 == 1 {
		bg = ColorBase02
	}
	return bg
}

func (ui *termui) ConsumableItem(g *game, i, lnum int, c consumable, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d available)", rune(i+97), c, g.Player.Consumables[c]), 0, lnum, fg, bg)
}

func (ui *termui) SelectProjectile(g *game, ev event) error {
	desc := false
	for {
		cs := g.SortedProjectiles()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuThrow.String(), MenuCols[MenuThrow][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which projectile? (press ? or click here for throwing menu)", col, 0)
		} else {
			ui.DrawColoredText("Throw", 0, 0, ColorOrange)
			col := utf8.RuneCountInString("Throw")
			ui.DrawText(" which projectile? (press ? or click here for describe menu)", col, 0)
		}
		for i, c := range cs {
			ui.ConsumableItem(g, i, i+1, c, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(cs)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(cs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.ConsumableItem(g, index, index+1, cs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(g, cs[index].Desc())
				continue
			}
			err = cs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *termui) SelectPotion(g *game, ev event) error {
	desc := false
	for {
		cs := g.SortedPotions()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuDrink.String(), MenuCols[MenuDrink][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which potion? (press ? or click here for quaff menu)", col, 0)
		} else {
			ui.DrawColoredText("Drink", 0, 0, ColorGreen)
			col := utf8.RuneCountInString("Drink")
			ui.DrawText(" which potion? (press ? or click here for description menu)", col, 0)
		}
		for i, c := range cs {
			ui.ConsumableItem(g, i, i+1, c, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(cs)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(cs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.ConsumableItem(g, index, index+1, cs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(g, cs[index].Desc())
				continue
			}
			err = cs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *termui) RodItem(g *game, i, lnum int, r rod, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s (%d/%d charges, %d mana cost)",
		rune(i+97), r, g.Player.Rods[r].Charge, r.MaxCharge(), r.MPCost()), 0, lnum, fg, bg)
}

func (ui *termui) SelectRod(g *game, ev event) error {
	desc := false
	for {
		rs := g.SortedRods()
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuEvoke.String(), MenuCols[MenuEvoke][0], DungeonHeight, ColorCyan)
		}
		if desc {
			ui.DrawColoredText("Describe", 0, 0, ColorBlue)
			col := utf8.RuneCountInString("Describe")
			ui.DrawText(" which rod? (press ? or click here for evocation menu)", col, 0)
		} else {
			ui.DrawColoredText("Evoke", 0, 0, ColorCyan)
			col := utf8.RuneCountInString("Evoke")
			ui.DrawText(" which rod? (press ? or click here for description menu)", col, 0)
		}
		for i, r := range rs {
			ui.RodItem(g, i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(rs)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(rs))
		if alt {
			desc = !desc
			continue
		}
		if err == nil {
			ui.RodItem(g, index, index+1, rs[index], ColorYellow)
			ui.Flush()
			time.Sleep(75 * time.Millisecond)
			if desc {
				ui.DrawDescription(g, rs[index].Desc())
				continue
			}
			err = rs[index].Use(g, ev)
		}
		return err
	}
}

func (ui *termui) ActionItem(g *game, i, lnum int, ka keyAction, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	desc := ka.NormalModeDescription()
	if !ka.NormalModeKey() {
		desc = ka.TargetingModeDescription()
	}
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), desc), 0, lnum, fg, bg)
}

var menuActions = []keyAction{
	KeyCharacterInfo,
	KeyLogs,
	KeyMenuCommandHelp,
	KeyMenuTargetingHelp,
	KeyConfigure,
	KeySave,
	KeyQuit,
}

func (ui *termui) SelectAction(g *game, actions []keyAction, ev event) (keyAction, error) {
	for {
		ui.ClearLine(0)
		if !ui.Small() {
			ui.DrawColoredText(MenuOther.String(), MenuCols[MenuOther][0], DungeonHeight, ColorCyan)
		}
		ui.DrawColoredText("Choose", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Choose")
		ui.DrawText(" which action?", col, 0)
		for i, r := range actions {
			ui.ActionItem(g, i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(g, NoFlushMode)
			return KeyExamine, err
		}
		ui.ActionItem(g, index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(g, NoFlushMode)
		return actions[index], nil
	}
}

type setting int

const (
	setKeys setting = iota
	invertLOS
	toggleLayout
	toggleTiles
)

func (s setting) String() (text string) {
	switch s {
	case setKeys:
		text = "Change key bindings"
	case invertLOS:
		text = "Toggle dark/light LOS"
	case toggleLayout:
		text = "Toggle normal/compact layout"
	case toggleTiles:
		text = "Toggle Tiles/Ascii display"
	}
	return text
}

var settingsActions = []setting{
	setKeys,
	invertLOS,
	toggleLayout,
}

func (ui *termui) ConfItem(g *game, i, lnum int, s setting, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), s), 0, lnum, fg, bg)
}

func (ui *termui) SelectConfigure(g *game, actions []setting) (setting, error) {
	for {
		ui.ClearLine(0)
		ui.DrawColoredText("Perform", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Perform")
		ui.DrawText(" which change?", col, 0)
		for i, r := range actions {
			ui.ConfItem(g, i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(g, NoFlushMode)
			return setKeys, err
		}
		ui.ConfItem(g, index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(g, NoFlushMode)
		return actions[index], nil
	}
}

func (ui *termui) HandleSettingAction(g *game) error {
	s, err := ui.SelectConfigure(g, settingsActions)
	if err != nil {
		return err
	}
	switch s {
	case setKeys:
		ui.ChangeKeys(g)
	case invertLOS:
		gameConfig.DarkLOS = !gameConfig.DarkLOS
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
		if gameConfig.DarkLOS {
			ApplyDarkLOS()
		} else {
			ApplyLightLOS()
		}
	case toggleLayout:
		ui.ApplyToggleLayout()
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
	case toggleTiles:
		ui.ApplyToggleTiles()
		err := g.SaveConfig()
		if err != nil {
			g.Print(err.Error())
		}
	}
	return nil
}

type wizardAction int

const (
	WizardInfoAction wizardAction = iota
	WizardToggleMap
)

func (a wizardAction) String() (text string) {
	switch a {
	case WizardInfoAction:
		text = "Info"
	case WizardToggleMap:
		text = "toggle see/hide monsters"
	}
	return text
}

var wizardActions = []wizardAction{
	WizardInfoAction,
	WizardToggleMap,
}

func (ui *termui) WizardItem(g *game, i, lnum int, s wizardAction, fg uicolor) {
	bg := ui.ListItemBG(i)
	ui.ClearLineWithColor(lnum, bg)
	ui.DrawColoredTextOnBG(fmt.Sprintf("%c - %s", rune(i+97), s), 0, lnum, fg, bg)
}

func (ui *termui) SelectWizardMagic(g *game, actions []wizardAction) (wizardAction, error) {
	for {
		ui.ClearLine(0)
		ui.DrawColoredText("Evoke", 0, 0, ColorCyan)
		col := utf8.RuneCountInString("Evoke")
		ui.DrawText(" which magic?", col, 0)
		for i, r := range actions {
			ui.WizardItem(g, i, i+1, r, ColorFg)
		}
		ui.DrawTextLine(" press esc or space to cancel ", len(actions)+1)
		ui.Flush()
		index, alt, err := ui.Select(g, len(actions))
		if alt {
			continue
		}
		if err != nil {
			ui.DrawDungeonView(g, NoFlushMode)
			return WizardInfoAction, err
		}
		ui.WizardItem(g, index, index+1, actions[index], ColorYellow)
		ui.Flush()
		time.Sleep(75 * time.Millisecond)
		ui.DrawDungeonView(g, NoFlushMode)
		return actions[index], nil
	}
}

func (ui *termui) HandleWizardAction(g *game) error {
	s, err := ui.SelectWizardMagic(g, wizardActions)
	if err != nil {
		return err
	}
	switch s {
	case WizardInfoAction:
		ui.WizardInfo(g)
	case WizardToggleMap:
		WizardMap = !WizardMap
		ui.DrawDungeonView(g, NoFlushMode)
	}
	return nil
}

func (ui *termui) Death(g *game) {
	g.Print("You die... --press esc or space to continue--")
	ui.DrawDungeonView(g, NormalMode)
	ui.WaitForContinue(g, -1)
	err := g.WriteDump()
	ui.Dump(g, err)
	ui.WaitForContinue(g, -1)
}

func (ui *termui) Win(g *game) {
	err := g.RemoveSaveFile()
	if err != nil {
		g.PrintfStyled("Error removing save file: %v", logError, err)
	}
	if g.Wizard {
		g.Print("You escape by the magic stairs! **WIZARD** --press esc or space to continue--")
	} else {
		g.Print("You escape by the magic stairs! You win. --press esc or space to continue--")
	}
	ui.DrawDungeonView(g, NormalMode)
	ui.WaitForContinue(g, -1)
	err = g.WriteDump()
	ui.Dump(g, err)
	ui.WaitForContinue(g, -1)
}

func (ui *termui) Dump(g *game, err error) {
	ui.Clear()
	ui.DrawText(g.SimplifedDump(err), 0, 0)
	ui.Flush()
}

func (ui *termui) CriticalHPWarning(g *game) {
	g.PrintStyled("*** CRITICAL HP WARNING *** --press esc or space to continue--", logCritic)
	ui.DrawDungeonView(g, NormalMode)
	ui.WaitForContinue(g, DungeonHeight)
	g.Print("Ok. Be careful, then.")
}

func (ui *termui) WoundedAnimation(g *game) {
	if DisableAnimations {
		return
	}
	r, _, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorFgHPwounded, bg)
	ui.Flush()
	time.Sleep(50 * time.Millisecond)
	if g.Player.HP <= 15 {
		ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorFgHPcritical, bg)
		ui.Flush()
		time.Sleep(50 * time.Millisecond)
	}
}

func (ui *termui) DrinkingPotionAnimation(g *game) {
	if DisableAnimations {
		return
	}
	ui.DrawDungeonView(g, NormalMode)
	time.Sleep(50 * time.Millisecond)
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorGreen, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorYellow, bg)
	ui.Flush()
	time.Sleep(75 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *termui) StatusEndAnimation(g *game) {
	if DisableAnimations {
		return
	}
	r, fg, bg := ui.PositionDrawing(g, g.Player.Pos)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, ColorViolet, bg)
	ui.Flush()
	time.Sleep(100 * time.Millisecond)
	ui.DrawAtPosition(g, g.Player.Pos, false, r, fg, bg)
	ui.Flush()
}

func (ui *termui) MenuSelectedAnimation(g *game, m menu, ok bool) {
	if DisableAnimations {
		return
	}
	if !ui.Small() {
		var message string
		if m == MenuInteract {
			message = ui.UpdateInteractButton(g)
		} else {
			message = m.String()
		}
		if message == "" {
			return
		}
		if ok {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorCyan)
		} else {
			ui.DrawColoredText(message, MenuCols[m][0], DungeonHeight, ColorMagenta)
		}
		ui.Flush()
		time.Sleep(25 * time.Millisecond)
		ui.DrawColoredText(m.String(), MenuCols[m][0], DungeonHeight, ColorViolet)
	}
}

func (ui *termui) MagicMappingAnimation(g *game, border []int) {
	if DisableAnimations {
		return
	}
	for _, i := range border {
		pos := idxtopos(i)
		r, fg, bg := ui.PositionDrawing(g, pos)
		ui.DrawAtPosition(g, pos, false, r, fg, bg)
	}
	ui.Flush()
	time.Sleep(12 * time.Millisecond)
}

func (ui *termui) Quit(g *game) bool {
	g.Print("Do you really want to quit without saving? [y/N]")
	ui.DrawDungeonView(g, NormalMode)
	quit := ui.PromptConfirmation(g)
	if quit {
		err := g.RemoveSaveFile()
		if err != nil {
			g.PrintfStyled("Error removing save file: %v ——press any key to quit——", logError, err)
			ui.DrawDungeonView(g, NormalMode)
			ui.PressAnyKey()
		}
	} else {
		g.Print(DoNothing)
	}
	return quit
}

func (ui *termui) Wizard(g *game) bool {
	g.Print("Do you really want to enter wizard mode (no return)? [y/N]")
	ui.DrawDungeonView(g, NormalMode)
	return ui.PromptConfirmation(g)
}

func (ui *termui) HandlePlayerTurn(g *game, ev event) bool {
getKey:
	for {
		var err error
		var again, quit bool
		if g.Targeting.valid() {
			err, again, quit = ui.ExaminePos(g, ev, g.Targeting)
		} else {
			ui.DrawDungeonView(g, NormalMode)
			err, again, quit = ui.PlayerTurnEvent(g, ev)
		}
		if err != nil && err.Error() != "" {
			g.Print(err.Error())
		}
		if again {
			continue getKey
		}
		return quit
	}
}
