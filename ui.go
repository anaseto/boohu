package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"
)

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *termui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *termui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *termui) WaitForContinue(g *game, line int) {
loop:
	for {
		in := ui.PollEvent()
		r := ui.KeyToRuneKeyAction(in)
		switch r {
		case '\x1b', ' ':
			break loop
		}
		if in.mouse && in.button == -1 {
			continue
		}
		if in.mouse && line >= 0 {
			if in.mouseY > line || in.mouseX > DungeonWidth {
				break loop
			}
		} else if in.mouse {
			break loop
		}
	}
}

func (ui *termui) PromptConfirmation(g *game) bool {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "Y", "y":
			return true
		default:
			return false
		}
	}
}

func (ui *termui) PressAnyKey() error {
	for {
		e := ui.PollEvent()
		if e.interrupt {
			return errors.New("interrupted")
		}
		if e.key != "" || (e.mouse && e.button != -1) {
			return nil
		}
	}
}

func (ui *termui) PlayerTurnEvent(g *game, ev event) (err error, again, quit bool) {
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "":
		if in.mouse {
			pos := position{X: in.mouseX, Y: in.mouseY}
			switch in.button {
			case -1:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(g, in.mouseX)
					omh := ui.menuHover
					if ok {
						ui.menuHover = m
					} else {
						ui.menuHover = -1
					}
					if ui.menuHover != omh {
						ui.DrawMenus(g)
						ui.Flush()
					}
					break
				}
				ui.menuHover = -1
				if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
					break
				}
				fallthrough
			case 0:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(g, in.mouseX)
					if !ok {
						again = true
						break
					}
					err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: m.Key(g)})
					if err != nil {
						again = true
					}
					return err, again, quit
				} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
				} else {
					err, again, quit = ui.ExaminePos(g, ev, pos)
				}
			case 2:
				err, again, quit = ui.HandleKeyAction(g, runeKeyAction{k: KeyMenu})
				if err != nil {
					again = true
				}
				return err, again, quit
			}
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r == 0 {
			again = true
		} else {
			err, again, quit = ui.HandleKeyAction(g, runeKeyAction{r: r})
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *termui) Scroll(n int) (m int, quit bool) {
	in := ui.PollEvent()
	switch in.key {
	case "Escape", "\x1b", " ":
		quit = true
	case "u":
		n -= 12
	case "d":
		n += 12
	case "j", "2":
		n++
	case "k", "8":
		n--
	case "":
		if in.mouse {
			switch in.button {
			case 0:
				y := in.mouseY
				x := in.mouseX
				if x >= DungeonWidth {
					quit = true
					break
				}
				if y > UIHeight {
					break
				}
				n += y - (DungeonHeight+3)/2
			}
		}
	}
	return n, quit
}

func (ui *termui) GetIndex(x, y int) int {
	return y*UIWidth + x
}

func (ui *termui) Select(g *game, l int) (index int, alternate bool, err error) {
	for {
		in := ui.PollEvent()
		r := ui.ReadKey(in.key)
		switch {
		case in.key == "\x1b" || in.key == "Escape" || in.key == " ":
			return -1, false, errors.New(DoNothing)
		case in.key == "?":
			return -1, true, nil
		case 97 <= r && int(r) < 97+l:
			return int(r - 97), false, nil
		case in.key == "" && in.mouse:
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case -1:
				oih := ui.itemHover
				if y <= 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					if oih != -1 {
						ui.ColorLine(oih, ColorFg)
						ui.Flush()
					}
					break
				}
				if y == oih {
					break
				}
				ui.itemHover = y
				ui.ColorLine(y, ColorYellow)
				if oih != -1 {
					ui.ColorLine(oih, ColorFg)
				}
				ui.Flush()
			case 0:
				if y < 0 || y > l || x >= DungeonWidth {
					ui.itemHover = -1
					return -1, false, errors.New(DoNothing)
				}
				if y == 0 {
					ui.itemHover = -1
					return -1, true, nil
				}
				ui.itemHover = -1
				return y - 1, false, nil
			case 2:
				ui.itemHover = -1
				return -1, true, nil
			case 1:
				ui.itemHover = -1
				return -1, false, errors.New(DoNothing)
			}
		}
	}
}

func (ui *termui) KeyMenuAction(n int) (m int, action keyConfigAction) {
	in := ui.PollEvent()
	r := ui.KeyToRuneKeyAction(in)
	switch string(r) {
	case "a":
		action = ChangeKeys
	case "\x1b", " ":
		action = QuitKeyConfig
	case "u":
		n -= DungeonHeight / 2
	case "d":
		n += DungeonHeight / 2
	case "j", "2", "ArrowDown":
		n++
	case "k", "8", "ArrowUp":
		n--
	case "R":
		action = ResetKeys
	default:
		if r == 0 && in.mouse {
			y := in.mouseY
			x := in.mouseX
			switch in.button {
			case 0:
				if x > DungeonWidth || y > DungeonHeight {
					action = QuitKeyConfig
				}
			case 1:
				action = QuitKeyConfig
			}
		}
	}
	return n, action
}

func (ui *termui) TargetModeEvent(g *game, targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "\x1b", "Escape", " ":
		g.Targeting = InvalidPos
		notarg = true
	case "":
		if !in.mouse {
			return
		}
		switch in.button {
		case -1:
			if in.mouseY == DungeonHeight {
				m, ok := ui.WhichButton(g, in.mouseX)
				omh := ui.menuHover
				if ok {
					ui.menuHover = m
				} else {
					ui.menuHover = -1
				}
				if ui.menuHover != omh {
					ui.DrawMenus(g)
					ui.Flush()
				}
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			ui.menuHover = -1
			if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
				break
			}
			mpos := position{in.mouseX, in.mouseY}
			if g.Targeting == mpos {
				break
			}
			g.Targeting = InvalidPos
			fallthrough
		case 0:
			if in.mouseY == DungeonHeight {
				m, ok := ui.WhichButton(g, in.mouseX)
				if !ok {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
					break
				}
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: m.Key(g)}, data)
			} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
			} else {
				again, notarg = ui.CursorMouseLeft(g, targ, position{X: in.mouseX, Y: in.mouseY}, data)
			}
		case 2:
			if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyMenu}, data)
			} else {
				err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyDescription}, data)
			}
		case 1:
			err, again, quit, notarg = ui.CursorKeyAction(g, targ, runeKeyAction{k: KeyExclude}, data)
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r != 0 {
			return ui.CursorKeyAction(g, targ, runeKeyAction{r: r}, data)
		}
		again = true
		notarg = true
	}
	return
}

func (ui *termui) ReadRuneKey() rune {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "\x1b", "Escape", " ":
			return 0
		case "Enter":
			return '.'
		}
		r := ui.ReadKey(in.key)
		if unicode.IsPrint(r) {
			return r
		}
	}
}

func (ui *termui) ReadKey(s string) (r rune) {
	bs := strings.NewReader(s)
	r, _, _ = bs.ReadRune()
	return r
}

type uiMode int

const (
	NormalMode uiMode = iota
	TargetingMode
	NoFlushMode
)

const DoNothing = "Do nothing, then."

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
		text = "Throw/Fire item"
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
		err = g.Equip(g.Ev)
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

func (ui *termui) HandleWizardAction(g *game) error {
	s, err := ui.SelectWizardMagic(g, wizardActions)
	if err != nil {
		return err
	}
	switch s {
	case WizardInfoAction:
		ui.WizardInfo(g)
	case WizardToggleMap:
		g.WizardMap = !g.WizardMap
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
