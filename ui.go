package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

type UICell struct {
	Fg    uicolor
	Bg    uicolor
	R     rune
	InMap bool
}

type uiInput struct {
	key       string
	mouse     bool
	mouseX    int
	mouseY    int
	button    int
	interrupt bool
}

func (ui *gameui) HideCursor() {
	ui.cursor = InvalidPos
}

func (ui *gameui) SetCursor(pos position) {
	ui.cursor = pos
}

func (ui *gameui) KeyToRuneKeyAction(in uiInput) rune {
	if utf8.RuneCountInString(in.key) != 1 {
		return 0
	}
	return ui.ReadKey(in.key)
}

func (ui *gameui) WaitForContinue(line int) {
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

func (ui *gameui) PromptConfirmation() bool {
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

func (ui *gameui) PressAnyKey() error {
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

type startAction int

const (
	StartPlay startAction = iota
	StartWatchReplay
)

func (ui *gameui) StartMenu(l int) startAction {
	for {
		in := ui.PollEvent()
		switch in.key {
		case "P", "p":
			ui.ColorLine(l, ColorYellow)
			ui.Flush()
			time.Sleep(10 * time.Millisecond)
			return StartPlay
		case "W", "w":
			ui.ColorLine(l+1, ColorYellow)
			ui.Flush()
			time.Sleep(10 * time.Millisecond)
			return StartWatchReplay
		}
		if in.key != "" && !in.mouse {
			continue
		}
		y := in.mouseY
		switch in.button {
		case -1:
			oih := ui.itemHover
			if y < l || y >= l+2 {
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
			if y < l || y >= l+2 {
				ui.itemHover = -1
				break
			}
			ui.itemHover = -1
			switch y - l {
			case 0:
				return StartPlay
			case 1:
				return StartWatchReplay
			}
		}
	}
}

func (ui *gameui) PlayerTurnEvent(ev event) (err error, again, quit bool) {
	g := ui.g
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "":
		if in.mouse {
			pos := position{X: in.mouseX, Y: in.mouseY}
			switch in.button {
			case -1:
				if in.mouseY == DungeonHeight {
					m, ok := ui.WhichButton(in.mouseX)
					omh := ui.menuHover
					if ok {
						ui.menuHover = m
					} else {
						ui.menuHover = -1
					}
					if ui.menuHover != omh {
						ui.DrawMenus()
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
					m, ok := ui.WhichButton(in.mouseX)
					if !ok {
						again = true
						break
					}
					err, again, quit = ui.HandleKeyAction(runeKeyAction{k: m.Key(g)})
					if err != nil {
						again = true
					}
					return err, again, quit
				} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
					again = true
				} else {
					err, again, quit = ui.ExaminePos(ev, pos)
				}
			case 2:
				err, again, quit = ui.HandleKeyAction(runeKeyAction{k: KeyMenu})
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
			err, again, quit = ui.HandleKeyAction(runeKeyAction{r: r})
		}
	}
	if err != nil {
		again = true
	}
	return err, again, quit
}

func (ui *gameui) Scroll(n int) (m int, quit bool) {
	in := ui.PollEvent()
	switch in.key {
	case "Escape", "\x1b", " ", "x":
		quit = true
	case "u", "9", "b":
		n -= 12
	case "d", "3", "f":
		n += 12
	case "j", "2", ".":
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

func (ui *gameui) GetIndex(x, y int) int {
	return y*UIWidth + x
}

func (ui *gameui) GetPos(i int) (int, int) {
	return i - (i/UIWidth)*UIWidth, i / UIWidth
}

func (ui *gameui) Select(l int) (index int, alternate bool, err error) {
	for {
		in := ui.PollEvent()
		r := ui.ReadKey(in.key)
		switch {
		case in.key == "\x1b" || in.key == "Escape" || in.key == " " || in.key == "x":
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

func (ui *gameui) KeyMenuAction(n int) (m int, action keyConfigAction) {
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

func (ui *gameui) TargetModeEvent(targ Targeter, data *examineData) (err error, again, quit, notarg bool) {
	g := ui.g
	again = true
	in := ui.PollEvent()
	switch in.key {
	case "\x1b", "Escape", " ", "x":
		g.Targeting = InvalidPos
		notarg = true
	case "":
		if !in.mouse {
			return
		}
		switch in.button {
		case -1:
			if in.mouseY == DungeonHeight {
				m, ok := ui.WhichButton(in.mouseX)
				omh := ui.menuHover
				if ok {
					ui.menuHover = m
				} else {
					ui.menuHover = -1
				}
				if ui.menuHover != omh {
					ui.DrawMenus()
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
				m, ok := ui.WhichButton(in.mouseX)
				if !ok {
					g.Targeting = InvalidPos
					notarg = true
					err = errors.New(DoNothing)
					break
				}
				err, again, quit, notarg = ui.CursorKeyAction(targ, runeKeyAction{k: m.Key(g)}, data)
			} else if in.mouseX >= DungeonWidth || in.mouseY >= DungeonHeight {
				g.Targeting = InvalidPos
				notarg = true
				err = errors.New(DoNothing)
			} else {
				again, notarg = ui.CursorMouseLeft(targ, position{X: in.mouseX, Y: in.mouseY}, data)
			}
		case 2:
			if in.mouseY >= DungeonHeight || in.mouseX >= DungeonWidth {
				err, again, quit, notarg = ui.CursorKeyAction(targ, runeKeyAction{k: KeyMenu}, data)
			} else {
				err, again, quit, notarg = ui.CursorKeyAction(targ, runeKeyAction{k: KeyDescription}, data)
			}
		case 1:
			err, again, quit, notarg = ui.CursorKeyAction(targ, runeKeyAction{k: KeyExclude}, data)
		}
	default:
		r := ui.KeyToRuneKeyAction(in)
		if r != 0 {
			return ui.CursorKeyAction(targ, runeKeyAction{r: r}, data)
		}
		again = true
	}
	return
}

func (ui *gameui) ReadRuneKey() rune {
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

func (ui *gameui) ReadKey(s string) (r rune) {
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

func (ui *gameui) EnterWizard() {
	g := ui.g
	if ui.Wizard() {
		g.WizardMode()
		ui.DrawDungeonView(NoFlushMode)
	} else {
		g.Print(DoNothing)
	}
}

func (ui *gameui) CleanError(err error) error {
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
	KeyRunW
	KeyRunS
	KeyRunN
	KeyRunE
	KeyRest
	KeyWaitTurn
	KeyDescend
	KeyGoToStairs
	KeyExplore
	KeyExamine
	KeyEvoke
	KeyInteract
	KeyInventory
	KeyLogs
	KeyDump
	KeyHelp
	KeySave
	KeyQuit
	KeyWizard
	KeyWizardInfo
	KeyWizardDescend

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
	KeyRunW,
	KeyRunS,
	KeyRunN,
	KeyRunE,
	KeyWaitTurn,
	KeyRest,
	KeyEvoke,
	KeyInteract,
	KeyInventory,
	KeyExamine,
	KeyGoToStairs,
	KeyExplore,
	KeyLogs,
	KeyDump,
	KeySave,
	KeyQuit,
	KeyMenu,
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
		KeyRunW, KeyRunS, KeyRunN, KeyRunE,
		KeyRest,
		KeyWaitTurn,
		KeyDescend,
		KeyGoToStairs,
		KeyExplore,
		KeyExamine,
		KeyEvoke,
		KeyInteract,
		KeyInventory,
		KeyLogs,
		KeyDump,
		KeyHelp,
		KeyMenu,
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
	case KeyRunW:
		text = "Travel west"
	case KeyRunS:
		text = "Travel south"
	case KeyRunN:
		text = "Travel north"
	case KeyRunE:
		text = "Travel east"
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
	case KeyEvoke:
		text = "Evoke card"
	case KeyInteract:
		text = "Interact"
	case KeyInventory:
		text = "Inventory"
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
	case KeyRunW:
		text = "Big move cursor west"
	case KeyRunS:
		text = "Big move cursor south"
	case KeyRunN:
		text = "Big move north"
	case KeyRunE:
		text = "Big move east"
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
		text = "Go to"
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
		KeyRunW, KeyRunS, KeyRunN, KeyRunE,
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
		'a': KeyW,
		's': KeyS,
		'w': KeyN,
		'd': KeyE,
		'4': KeyW,
		'2': KeyS,
		'8': KeyN,
		'6': KeyE,
		'H': KeyRunW,
		'J': KeyRunS,
		'K': KeyRunN,
		'L': KeyRunE,
		'.': KeyWaitTurn,
		'5': KeyWaitTurn,
		'G': KeyGoToStairs,
		'o': KeyExplore,
		'x': KeyExamine,
		'v': KeyEvoke,
		'z': KeyEvoke,
		'e': KeyInteract, // intEract
		'i': KeyInventory,
		'm': KeyLogs,
		'M': KeyMenu,
		'#': KeyDump,
		'?': KeyHelp,
		'S': KeySave,
		'Q': KeyQuit,
		'W': KeyWizard,
		'@': KeyWizardInfo,
		'>': KeyWizardDescend,
		'=': KeyConfigure,
	}
	gameConfig.RuneTargetModeKeys = map[rune]keyAction{
		'h':    KeyW,
		'j':    KeyS,
		'k':    KeyN,
		'l':    KeyE,
		'a':    KeyW,
		's':    KeyS,
		'w':    KeyN,
		'd':    KeyE,
		'4':    KeyW,
		'2':    KeyS,
		'8':    KeyN,
		'6':    KeyE,
		'H':    KeyRunW,
		'J':    KeyRunS,
		'K':    KeyRunN,
		'L':    KeyRunE,
		'>':    KeyNextStairs,
		'-':    KeyPreviousMonster,
		'+':    KeyNextMonster,
		'o':    KeyNextObject,
		']':    KeyNextObject,
		')':    KeyNextObject,
		'(':    KeyNextObject,
		'[':    KeyNextObject,
		'_':    KeyNextObject,
		'=':    KeyNextObject,
		'v':    KeyDescription,
		'.':    KeyTarget,
		't':    KeyTarget,
		'g':    KeyTarget,
		'e':    KeyExclude,
		' ':    KeyEscape,
		'\x1b': KeyEscape,
		'x':    KeyEscape,
		'?':    KeyHelp,
	}
	CustomKeys = false
}

type runeKeyAction struct {
	r rune
	k keyAction
}

func (ui *gameui) HandleKeyAction(rka runeKeyAction) (err error, again bool, quit bool) {
	g := ui.g
	if rka.r != 0 {
		var ok bool
		rka.k, ok = gameConfig.RuneNormalModeKeys[rka.r]
		if !ok {
			switch rka.r {
			case 's':
				err = errors.New("Unknown key. Did you mean capital S for save and quit?")
			case 'q':
				err = errors.New("Unknown key. Did you mean capital Q for quit without saving?")
			default:
				err = fmt.Errorf("Unknown key '%c'. Type ? for help.", rka.r)
			}
			return err, again, quit
		}
	}
	if rka.k == KeyMenu {
		rka.k, err = ui.SelectAction(menuActions, g.Ev)
		if err != nil {
			err = ui.CleanError(err)
			return err, again, quit
		}
	}
	return ui.HandleKey(rka)
}

func (ui *gameui) OptionalDescendConfirmation(st stair) (err error) {
	g := ui.g
	if g.Depth == WinDepth && st == NormalStair && g.Dungeon.Cell(g.Places.Shaedra).T == StoryCell {
		err = errors.New("You have to rescue Shaedra first!")
	}
	return err

}

func (ui *gameui) HandleKey(rka runeKeyAction) (err error, again bool, quit bool) {
	g := ui.g
	switch rka.k {
	case KeyW, KeyS, KeyN, KeyE:
		err = g.MovePlayer(g.Player.Pos.To(KeyToDir(rka.k)), g.Ev)
	case KeyRunW, KeyRunS, KeyRunN, KeyRunE:
		err = g.GoToDir(KeyToDir(rka.k), g.Ev)
	case KeyWaitTurn:
		g.WaitTurn(g.Ev)
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
	case KeyInteract:
		c := g.Dungeon.Cell(g.Player.Pos)
		switch c.T {
		case StairCell:
			if g.Dungeon.Cell(g.Player.Pos).T == StairCell {
				ui.MenuSelectedAnimation(MenuInteract, true)
				strt := g.Objects.Stairs[g.Player.Pos]
				err = ui.OptionalDescendConfirmation(strt)
				if err != nil {
					break
				}
				if g.Descend() {
					ui.Win()
					quit = true
					return err, again, quit
				}
				ui.DrawDungeonView(NormalMode)
			} else {
				err = errors.New("No stairs here.")
			}
		case BarrelCell:
			ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.Rest(g.Ev)
			if err != nil {
				ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case MagaraCell:
			err = ui.EquipMagara(g.Ev)
			err = ui.CleanError(err)
		case StoneCell:
			ui.MenuSelectedAnimation(MenuInteract, true)
			err = g.ActivateStone()
			if err != nil {
				ui.MenuSelectedAnimation(MenuInteract, false)
			}
		case ScrollCell:
			err = ui.ReadScroll()
			err = ui.CleanError(err)
		case ItemCell:
			err = ui.g.EquipItem()
		default:
			err = errors.New("You cannot interact with anything here.")
		}
	case KeyEvoke:
		err = ui.SelectMagara(g.Ev)
		err = ui.CleanError(err)
	case KeyInventory:
		err = ui.SelectItem(g.Ev)
		err = ui.CleanError(err)
	case KeyExplore:
		//ui.MenuSelectedAnimation(MenuExplore, true)
		err = g.Autoexplore(g.Ev)
		if err != nil {
			//ui.MenuSelectedAnimation(MenuExplore, false)
		}
	case KeyExamine:
		err, again, quit = ui.Examine(nil)
	case KeyHelp, KeyMenuCommandHelp:
		ui.KeysHelp()
		again = true
	case KeyMenuTargetingHelp:
		ui.ExamineHelp()
		again = true
	case KeyLogs:
		ui.DrawPreviousLogs()
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
		} else {
			dataDir, _ := g.DataDir()
			if dataDir != "" {
				g.Printf("Game statistics written to %s.", filepath.Join(dataDir, "dump"))
			} else {
				g.Print("Game statistics written.")
			}
		}
		again = true
	case KeyWizardInfo:
		if g.Wizard {
			err = ui.HandleWizardAction()
			again = true
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case KeyWizardDescend:
		if g.Wizard && g.Depth < MaxDepth {
			if g.Descend() {
				ui.Win()
				quit = true
				return err, again, quit
			}
		} else {
			err = errors.New("Unknown key. Type ? for help.")
		}
	case KeyWizard:
		ui.EnterWizard()
		return nil, true, false
	case KeyQuit:
		if ui.Quit() {
			return nil, false, true
		}
		return nil, true, false
	case KeyConfigure:
		err = ui.HandleSettingAction()
		again = true
	case KeyDescription:
		//ui.MenuSelectedAnimation(MenuView, false)
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

func (ui *gameui) ExaminePos(ev event, pos position) (err error, again, quit bool) {
	var start *position
	if pos.valid() {
		start = &pos
	}
	err, again, quit = ui.Examine(start)
	return err, again, quit
}

func (ui *gameui) Examine(start *position) (err error, again, quit bool) {
	ex := &examiner{}
	err, again, quit = ui.CursorAction(ex, start)
	return err, again, quit
}

func (ui *gameui) ChooseTarget(targ Targeter) error {
	err, _, _ := ui.CursorAction(targ, nil)
	if err != nil {
		return err
	}
	if !targ.Done() {
		return errors.New(DoNothing)
	}
	return nil
}

func (ui *gameui) NextMonster(r rune, pos position, data *examineData) {
	g := ui.g
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

func (ui *gameui) NextStair(data *examineData) {
	g := ui.g
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

func (ui *gameui) NextObject(pos position, data *examineData) {
	g := ui.g
	nobject := data.nobject
	if len(data.objects) == 0 {
		for p := range g.Objects.Stairs {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Stones {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Barrels {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Magaras {
			data.objects = append(data.objects, p)
		}
		for p := range g.Objects.Bananas {
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

func (ui *gameui) ExcludeZone(pos position) {
	g := ui.g
	if !g.Dungeon.Cell(pos).Explored {
		g.Print("You cannot choose an unexplored cell for exclusion.")
	} else {
		toggle := !g.ExclusionsMap[pos]
		g.ComputeExclusion(pos, toggle)
	}
}

func (ui *gameui) CursorMouseLeft(targ Targeter, pos position, data *examineData) (again, notarg bool) {
	g := ui.g
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

func (ui *gameui) CursorKeyAction(targ Targeter, rka runeKeyAction, data *examineData) (err error, again, quit, notarg bool) {
	g := ui.g
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
		rka.k, err = ui.SelectAction(menuActions, g.Ev)
		if err != nil {
			err = ui.CleanError(err)
			return err, again, quit, notarg
		}
	}
	switch rka.k {
	case KeyW, KeyS, KeyN, KeyE:
		data.npos = pos.To(KeyToDir(rka.k))
	case KeyRunW, KeyRunS, KeyRunN, KeyRunE:
		for i := 0; i < 5; i++ {
			p := data.npos.To(KeyToDir(rka.k))
			if !p.valid() {
				break
			}
			data.npos = p
		}
	case KeyNextStairs:
		ui.NextStair(data)
	case KeyDescend:
		if g.Dungeon.Cell(g.Player.Pos).T == StairCell {
			ui.MenuSelectedAnimation(MenuInteract, true)
			strt := g.Objects.Stairs[g.Player.Pos]
			err = ui.OptionalDescendConfirmation(strt)
			if err != nil {
				break
			}
			again = false
			g.Targeting = InvalidPos
			notarg = true
			if g.Descend() {
				ui.Win()
				quit = true
				return err, again, quit, notarg
			}
		} else {
			err = errors.New("No stairs here.")
		}
	case KeyPreviousMonster, KeyNextMonster:
		ui.NextMonster(rka.r, pos, data)
	case KeyNextObject:
		ui.NextObject(pos, data)
	case KeyHelp, KeyMenuTargetingHelp:
		ui.HideCursor()
		ui.ExamineHelp()
		ui.SetCursor(pos)
	case KeyMenuCommandHelp:
		ui.HideCursor()
		ui.KeysHelp()
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
		ui.ViewPositionDescription(pos)
		ui.SetCursor(pos)
	case KeyExclude:
		ui.ExcludeZone(pos)
	case KeyEscape:
		g.Targeting = InvalidPos
		notarg = true
		err = errors.New(DoNothing)
	case KeyExplore, KeyRest, KeyLogs, KeyEvoke, KeyInventory:
		// XXX: hm, this is only useful with mouse in terminal, rarely tested.
		if _, ok := targ.(*examiner); !ok {
			break
		}
		err, again, quit = ui.HandleKey(rka)
		if err != nil {
			notarg = true
		}
		g.Targeting = InvalidPos
	case KeyConfigure:
		err = ui.HandleSettingAction()
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
		if ui.Quit() {
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

func (ui *gameui) CursorAction(targ Targeter, start *position) (err error, again, quit bool) {
	g := ui.g
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
		ui.NextObject(InvalidPos, data)
		if !data.npos.valid() {
			ui.NextStair(data)
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
			ui.DescribePosition(pos, targ)
		}
		opos = pos
		targ.ComputeHighlight(g, pos)
		ui.SetCursor(pos)
		m := g.MonsterAt(pos)
		if m.Exists() && g.Player.Sees(pos) {
			g.ComputeMonsterCone(m)
		} else {
			g.MonsterTargLOS = nil
		}
		ui.DrawDungeonView(TargetingMode)
		ui.DrawInfoLine(g.InfoEntry)
		if !ui.Small() {
			st := " Examine/Travel mode "
			if _, ok := targ.(*examiner); !ok {
				st = " Targeting mode "
			}
			ui.DrawStyledTextLine(st, DungeonHeight+2, FooterLine)
		}
		ui.SetCell(DungeonWidth, DungeonHeight, '┤', ColorFg, ColorBg)
		ui.Flush()
		data.npos = pos
		var notarg bool
		err, again, quit, notarg = ui.TargetModeEvent(targ, data)
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
	g.MonsterTargLOS = nil
	ui.HideCursor()
	return err, again, quit
}

type menu int

const (
	//MenuExplore menu = iota
	MenuOther menu = iota
	MenuInventory
	MenuEvoke
	MenuInteract
)

func (m menu) String() (text string) {
	switch m {
	//case MenuExplore:
	//text = "explore"
	case MenuEvoke:
		text = "evoke"
	case MenuInventory:
		text = "inventory"
	case MenuOther:
		text = "menu"
	case MenuInteract:
		text = "interact"
	}
	return "[" + text + "]"
}

func (m menu) Key(g *game) (key keyAction) {
	switch m {
	//case MenuExplore:
	//key = KeyExplore
	case MenuOther:
		key = KeyMenu
	case MenuInventory:
		key = KeyInventory
	case MenuEvoke:
		key = KeyEvoke
	case MenuInteract:
		key = KeyInteract
	}
	return key
}

var MenuCols = [][2]int{
	//MenuExplore:  {0, 0},
	MenuOther:     {0, 0},
	MenuInventory: {0, 0},
	MenuEvoke:     {0, 0},
	MenuInteract:  {0, 0}}

func init() {
	for i := range MenuCols {
		runes := utf8.RuneCountInString(menu(i).String())
		if i == 0 {
			MenuCols[0] = [2]int{25, 25 + runes}
			continue
		}
		MenuCols[i] = [2]int{MenuCols[i-1][1] + 2, MenuCols[i-1][1] + 2 + runes}
	}
}

func (ui *gameui) WhichButton(col int) (menu, bool) {
	g := ui.g
	if ui.Small() {
		return MenuOther, false
	}
	end := len(MenuCols) - 1
	switch g.Dungeon.Cell(g.Player.Pos).T {
	case StairCell, BarrelCell, ScrollCell, MagaraCell, StoneCell:
		end++
	}
	for i, cols := range MenuCols[0:end] {
		if cols[0] >= 0 && col >= cols[0] && col < cols[1] {
			return menu(i), true
		}
	}
	return MenuOther, false
}

func (ui *gameui) UpdateInteractButton() string {
	g := ui.g
	var interactMenu string
	var show bool
	switch g.Dungeon.Cell(g.Player.Pos).T {
	case StairCell:
		interactMenu = "[descend]"
		show = true
	case BarrelCell:
		interactMenu = "[rest]"
		show = true
	case MagaraCell:
		interactMenu = "[equip]"
		show = true
	case StoneCell:
		interactMenu = "[activate]"
		show = true
	case ScrollCell:
		interactMenu = "[read]"
		show = true
	case ItemCell:
		interactMenu = "[equip]"
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

func (ui *gameui) HandleWizardAction() error {
	g := ui.g
	s, err := ui.SelectWizardMagic(wizardActions)
	if err != nil {
		return err
	}
	switch s {
	case WizardInfoAction:
		ui.WizardInfo()
	case WizardToggleMap:
		g.WizardMap = !g.WizardMap
		ui.DrawDungeonView(NoFlushMode)
	}
	return nil
}

func (ui *gameui) Death() {
	g := ui.g
	g.Print("You die... [(x) to continue]")
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(-1)
	err := g.WriteDump()
	ui.Dump(err)
	ui.WaitForContinue(-1)
}

func (ui *gameui) Win() {
	g := ui.g
	err := g.RemoveSaveFile()
	if err != nil {
		g.PrintfStyled("Error removing save file: %v", logError, err)
	}
	if g.Wizard {
		g.Print("You escape by the magic stairs! **WIZARD** [(x) to continue]")
	} else {
		g.Print("You escape by the magic stairs! You win. [(x) to continue]")
	}
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(-1)
	err = g.WriteDump()
	ui.Dump(err)
	ui.WaitForContinue(-1)
}

func (ui *gameui) Dump(err error) {
	g := ui.g
	ui.Clear()
	ui.DrawText(g.SimplifedDump(err), 0, 0)
	ui.Flush()
}

func (ui *gameui) CriticalHPWarning() {
	g := ui.g
	g.PrintStyled("*** CRITICAL HP WARNING *** [(x) to continue]", logCritic)
	ui.DrawDungeonView(NormalMode)
	ui.WaitForContinue(DungeonHeight)
	g.Print("Ok. Be careful, then.")
}

func (ui *gameui) Quit() bool {
	g := ui.g
	g.Print("Do you really want to quit without saving? [y/N]")
	ui.DrawDungeonView(NormalMode)
	quit := ui.PromptConfirmation()
	if quit {
		err := g.RemoveSaveFile()
		if err != nil {
			g.PrintfStyled("Error removing save file: %v ——press any key to quit——", logError, err)
			ui.DrawDungeonView(NormalMode)
			ui.PressAnyKey()
		}
	} else {
		g.Print(DoNothing)
	}
	return quit
}

func (ui *gameui) Wizard() bool {
	g := ui.g
	g.Print("Do you really want to enter wizard mode (no return)? [y/N]")
	ui.DrawDungeonView(NormalMode)
	return ui.PromptConfirmation()
}

func (ui *gameui) HandlePlayerTurn(ev event) bool {
	g := ui.g
getKey:
	for {
		var err error
		var again, quit bool
		if g.Targeting.valid() {
			err, again, quit = ui.ExaminePos(ev, g.Targeting)
		} else {
			ui.DrawDungeonView(NormalMode)
			err, again, quit = ui.PlayerTurnEvent(ev)
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

func (ui *gameui) ExploreStep() bool {
	next := make(chan bool)
	var stop bool
	go func() {
		time.Sleep(10 * time.Millisecond)
		ui.Interrupt()
	}()
	go func() {
		err := ui.PressAnyKey()
		interrupted := err != nil
		next <- !interrupted
	}()
	stop = <-next
	ui.DrawDungeonView(NormalMode)
	return stop
}

func (ui *gameui) Clear() {
	for i := 0; i < UIHeight*UIWidth; i++ {
		x, y := ui.GetPos(i)
		ui.SetCell(x, y, ' ', ColorFg, ColorBg)
	}
}

func (ui *gameui) DrawBufferInit() {
	if len(ui.g.DrawBuffer) == 0 {
		ui.g.DrawBuffer = make([]UICell, UIHeight*UIWidth)
	} else if len(ui.g.DrawBuffer) != UIHeight*UIWidth {
		ui.g.DrawBuffer = make([]UICell, UIHeight*UIWidth)
	}
}

func ApplyConfig() {
	if gameConfig.RuneNormalModeKeys == nil || gameConfig.RuneTargetModeKeys == nil {
		ApplyDefaultKeyBindings()
	}
	if gameConfig.DarkLOS {
		ApplyDarkLOS()
	} else {
		ApplyLightLOS()
	}
}
