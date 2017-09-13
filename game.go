package main

import (
	"bytes"
	"container/heap"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// idea: punish (mark) the player if it remains too many turns not in combat
// and not exploring new cells and in good health
// monster induced berserker

type game struct {
	Dungeon             *dungeon
	Player              *player
	Monsters            []*monster
	Bands               []monsterBand
	Events              *eventQueue
	Highlight           map[position]bool // highlighted positions (e.g. targeted ray)
	Collectables        map[position]*collectable
	Equipables          map[position]equipable
	Rods                map[position]rod
	Stairs              map[position]bool
	Clouds              map[position]cloud
	GeneratedBands      map[monsterBand]int
	GeneratedEquipables map[equipable]bool
	GeneratedRods       map[rod]bool
	Gold                map[position]int
	Resting             bool
	Autoexploring       bool
	AutoexploreMap      nodeMap
	AutoTarget          *position
	AutoHalt            bool
	AutoNext            bool
	Quit                bool
	ui                  Renderer
	Depth               int
	Wizard              bool
	Log                 []string
	Turn                int
	Killed              int
	Scumming            int
}

func init() {
	gob.Register(potion(0))
	gob.Register(projectile(0))
	gob.Register(&simpleEvent{})
	gob.Register(&monsterEvent{})
	gob.Register(&cloudEvent{})
	gob.Register(armour(0))
	gob.Register(weapon(0))
	gob.Register(shield(0))
}

func (g *game) DataDir() (string, error) {
	var xdg string
	if os.Getenv("GOOS") == "windows" {
		xdg = os.Getenv("LOCALAPPDATA")
	} else {
		xdg = os.Getenv("XDG_DATA_HOME")
	}
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	dataDir := filepath.Join(xdg, "boohu")
	_, err := os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return "", fmt.Errorf("%v\n", err)
		}
	}
	return dataDir, nil
}

func (g *game) Save() {
	dataDir, err := g.DataDir()
	if err != nil {
		g.Print(err.Error())
		return
	}
	saveFile := filepath.Join(dataDir, "save.gob")
	var data bytes.Buffer
	enc := gob.NewEncoder(&data)
	err = enc.Encode(g)
	if err != nil {
		g.Print(err.Error())
		return
	}
	err = ioutil.WriteFile(saveFile, data.Bytes(), 0644)
	if err != nil {
		g.Print(err.Error())
	}
}

func (g *game) RemoveSaveFile() {
	dataDir, err := g.DataDir()
	if err != nil {
		g.Print(err.Error())
		return
	}
	saveFile := filepath.Join(dataDir, "save.gob")
	_, err = os.Stat(saveFile)
	if err == nil {
		err := os.Remove(saveFile)
		if err != nil {
			fmt.Fprint(os.Stderr, "Error removing old save file")
		}
	}
}

func (g *game) Load() (bool, error) {
	dataDir, err := g.DataDir()
	if err != nil {
		return false, err
	}
	saveFile := filepath.Join(dataDir, "save.gob")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new game
		return false, err
	}
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return true, err
	}
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	var lg game
	err = dec.Decode(&lg)
	if err != nil {
		return true, err
	}
	*g = lg
	return true, nil
}

func (g *game) DumpAptitudes() string {
	apts := []string{}
	for apt, b := range g.Player.Aptitudes {
		if b {
			apts = append(apts, apt.String())
		}
	}
	sort.Strings(apts)
	if len(apts) > 0 {
		return "Aptitudes:\n" + strings.Join(apts, "\n")
	} else {
		return "You do not have any special aptitudes."
	}
}

func (g *game) SortedRods() rodSlice {
	var rs rodSlice
	for k, p := range g.Player.Rods {
		if p == nil {
			continue
		}
		rs = append(rs, k)
	}
	sort.Sort(rs)
	return rs
}

func (g *game) SortedPotions() consumableSlice {
	var cs consumableSlice
	for k, _ := range g.Player.Consumables {
		switch k := k.(type) {
		case potion:
			cs = append(cs, k)
		}
	}
	sort.Sort(cs)
	return cs
}

func (g *game) SortedProjectiles() consumableSlice {
	var cs consumableSlice
	for k, _ := range g.Player.Consumables {
		switch k := k.(type) {
		case projectile:
			cs = append(cs, k)
		}
	}
	sort.Sort(cs)
	return cs
}

func (g *game) Dump() string {
	buf := &bytes.Buffer{}
	if g.Wizard {
		fmt.Fprintf(buf, "**WIZARD MODE**\n")
	}
	if g.Player.HP > 0 && g.Depth > 12 {
		fmt.Fprintf(buf, "You escaped from Hareka's Underground alive!\n")
	} else if g.Player.HP <= 0 {
		fmt.Fprintf(buf, "You died while exploring depth %d of Hareka's Underground.\n", g.Depth)
	} else {
		fmt.Fprintf(buf, "You are exploring depth %d of Hareka's Underground.\n", g.Depth)
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, g.DumpAptitudes())
	fmt.Fprintf(buf, "\n\n")
	fmt.Fprintf(buf, "Equipment:\n")
	fmt.Fprintf(buf, "You are wearing a %v.\n", g.Player.Armour)
	fmt.Fprintf(buf, "You are wielding a %v.\n", g.Player.Weapon)
	if g.Player.Shield != NoShield {
		if g.Player.Weapon.TwoHanded() {
			fmt.Fprintf(buf, "You had a %v.\n", g.Player.Shield)
		} else {
			fmt.Fprintf(buf, "You are wearing a %v.\n", g.Player.Shield)
		}
	}
	fmt.Fprintf(buf, "\n")
	rs := g.SortedRods()
	if len(rs) > 0 {
		fmt.Fprintf(buf, "Rods:\n")
		for _, c := range rs {
			fmt.Fprintf(buf, "- %s (%d/%d charges)\n",
				c, g.Player.Rods[c].Charge, c.MaxCharge())
		}
	} else {
		fmt.Fprintf(buf, "You do not have any rods.\n")
	}
	fmt.Fprintf(buf, "\n")
	ps := g.SortedPotions()
	if len(ps) > 0 {
		fmt.Fprintf(buf, "Potions:\n")
		for _, c := range ps {
			fmt.Fprintf(buf, "- %s (%d available)\n", c, g.Player.Consumables[c])
		}
	} else {
		fmt.Fprintf(buf, "You do not have any potions.\n")
	}
	fmt.Fprintf(buf, "\n")
	ps = g.SortedProjectiles()
	if len(ps) > 0 {
		fmt.Fprintf(buf, "Projectiles:\n")
		for _, c := range ps {
			fmt.Fprintf(buf, "- %s (%d available)\n", c, g.Player.Consumables[c])
		}
	} else {
		fmt.Fprintf(buf, "You do not have any projectiles.\n")
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Miscelaneous:\n")
	fmt.Fprintf(buf, "You collected %d gold coins.\n", g.Player.Gold)
	fmt.Fprintf(buf, "You killed %d monsters.\n", g.Killed)
	fmt.Fprintf(buf, "You spent %.1f turns in the Underground.\n", float64(g.Turn)/10)
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Last messages:\n")
	for i := len(g.Log) - 10; i < len(g.Log) && i >= 0; i++ {
		fmt.Fprintf(buf, "%s\n", g.Log[i])
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Dungeon:\n")
	fmt.Fprintf(buf, "\n")
	buf.WriteString(g.DumpDungeon())
	return buf.String()
}

func (g *game) DumpDungeon() string {
	buf := bytes.Buffer{}
	for i, c := range g.Dungeon.Cells {
		if i%g.Dungeon.Width == 0 {
			buf.WriteRune('\n')
		}
		pos := g.Dungeon.CellPosition(i)
		if !c.Explored {
			buf.WriteRune(' ')
			continue
		}
		var r rune
		switch c.T {
		case WallCell:
			r = '#'
		case FreeCell:
			switch {
			case pos == g.Player.Pos:
				r = '@'
			default:
				r = '.'
				if _, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
					r = '§'
				}
				if c, ok := g.Collectables[pos]; ok {
					r = c.Consumable.Letter()
				} else if eq, ok := g.Equipables[pos]; ok {
					r = eq.Letter()
				} else if rod, ok := g.Rods[pos]; ok {
					r = rod.Letter()
				} else if _, ok := g.Stairs[pos]; ok {
					r = '>'
				} else if _, ok := g.Gold[pos]; ok {
					r = '$'
				}
				m, _ := g.MonsterAt(pos)
				if m.Exists() && (g.Player.LOS[m.Pos] || g.Wizard) {
					r = m.Kind.Letter()
				}
			}
		}
		buf.WriteRune(r)
	}
	return buf.String()
}

func (g *game) SimplifedDump() string {
	buf := &bytes.Buffer{}
	if g.Wizard {
		fmt.Fprintf(buf, "**WIZARD MODE**\n")
	}
	if g.Player.HP > 0 && g.Depth > 12 {
		fmt.Fprintf(buf, "You escaped from Hareka's Underground alive!\n")
	} else if g.Player.HP <= 0 {
		fmt.Fprintf(buf, "You died while exploring depth %d of Hareka's Underground.\n", g.Depth)
	} else {
		fmt.Fprintf(buf, "You are exploring depth %d of Hareka's Underground.\n", g.Depth)
	}
	fmt.Fprintf(buf, "You collected %d gold coins.\n", g.Player.Gold)
	fmt.Fprintf(buf, "You killed %d monsters.\n", g.Killed)
	fmt.Fprintf(buf, "You spent %.1f turns in the Underground.\n", float64(g.Turn)/10)
	fmt.Fprintf(buf, "\n")
	dataDir, err := g.DataDir()
	if err == nil {
		fmt.Fprintf(buf, "Full dump written to %s.\n", filepath.Join(dataDir, "dump"))
	}
	fmt.Fprintf(buf, "\n\n")
	fmt.Fprintf(buf, "───Press esc or space to quit───")
	return buf.String()
}

func (g *game) WriteDump() error {
	dataDir, err := g.DataDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing dump: %s", err)
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dataDir, "dump"), []byte(g.Dump()), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing dump: %s", err)
		return err
	}
	return nil
}

func (g *game) SeenGoodWeapon() bool {
	return g.GeneratedEquipables[Sword] || g.GeneratedEquipables[DoubleSword] || g.GeneratedEquipables[Spear] || g.GeneratedEquipables[Halberd] ||
		g.GeneratedEquipables[Axe] || g.GeneratedEquipables[BattleAxe]
}

func (g *game) Print(s string) {
	g.Log = append(g.Log, s)
	if len(g.Log) > 1000 {
		g.Log = g.Log[500:]
	}
}

func (g *game) Printf(format string, a ...interface{}) {
	g.Log = append(g.Log, fmt.Sprintf(format, a...))
	if len(g.Log) > 1000 {
		g.Log = g.Log[500:]
	}
}

type Renderer interface {
	AutoExploreStep(*game)
	HandlePlayerTurn(*game, event) bool
	Death(*game)
	ChooseTarget(*game, Targetter) bool
}

func (g *game) FreeCell() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCell")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos == pos {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForImportantStair() position {
	for {
		pos := g.FreeCellForStatic()
		if pos.Distance(g.Player.Pos) > 12 {
			return pos
		}
	}
}

func (g *game) FreeCellForStatic() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForStatic")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos == pos {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			if g.Gold[pos] > 0 {
				continue
			}
			if g.Collectables[pos] != nil {
				continue
			}
			if g.Stairs[pos] {
				continue
			}
			if _, ok := g.Rods[pos]; ok {
				continue
			}
			if _, ok := g.Equipables[pos]; ok {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForMonster() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForMonster")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
				continue
			}
			mons, _ := g.MonsterAt(pos)
			if mons.Exists() {
				continue
			}
			return pos
		}
	}
}

func (g *game) FreeCellForBandMonster(pos position) position {
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeCellForBandMonster")
		}
		neighbors := g.Dungeon.FreeNeighbors(pos)
		r := RandInt(len(neighbors))
		pos = neighbors[r]
		if g.Player != nil && g.Player.Pos.Distance(pos) < 8 {
			continue
		}
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		return pos
	}
}

func (g *game) FreeForStairs() position {
	m := g.Dungeon
	count := 0
	for {
		count++
		if count > 1000 {
			panic("FreeForStairs")
		}
		x := RandInt(m.Width)
		y := RandInt(m.Heigth)
		pos := position{x, y}
		c := m.Cell(pos)
		if c.T == FreeCell {
			_, ok := g.Collectables[pos]
			if ok {
				continue
			}
			return pos
		}
	}
}

func (g *game) MaxDepth() int {
	return 12
}

func (g *game) GenDungeon() {
	switch RandInt(6) {
	case 0:
		g.GenCaveMap(21, 79)
	case 1:
		g.GenRoomMap(21, 79)
	case 2:
		g.GenCellularAutomataCaveMap(21, 79)
	case 3:
		g.GenCaveMapTree(21, 79)
	default:
		g.GenRuinsMap(21, 79)
	}
}

func (g *game) InitLevel() {
	// Dungeon terrain
	g.GenDungeon()

	// Player
	if g.Depth == 0 {
		g.Player = &player{
			HP:        40,
			MP:        10,
			Gold:      0,
			Aptitudes: map[aptitude]bool{},
		}
		g.Player.Consumables = map[consumable]int{
			HealWoundsPotion:    1,
			TeleportationPotion: 1,
			Javeline:            3,
		}
		g.GeneratedRods = map[rod]bool{}
		g.Player.Rods = map[rod]*rodProps{}
		g.Player.Statuses = map[status]int{}
		g.GeneratedEquipables = map[equipable]bool{}
		g.GeneratedBands = map[monsterBand]int{}
	}
	g.Player.Pos = g.FreeCell()

	// Monsters
	g.GenMonsters()

	// Collectables
	g.Collectables = make(map[position]*collectable)
	for c, data := range ConsumablesCollectData {
		r := RandInt(data.rarity)
		if r == 0 {
			pos := g.FreeCellForStatic()
			g.Collectables[pos] = &collectable{Consumable: c, Quantity: data.quantity}
		}
	}

	// Equipment
	g.Equipables = make(map[position]equipable)
	for eq, data := range EquipablesRepartitionData {
		g.GenEquip(eq, data)
	}

	// Rods
	g.Rods = map[position]rod{}
	r := 5*g.GeneratedRodsCount() - g.Depth + 4
	if r < 2 {
		r = 1
	}
	if RandInt(r) == 0 && g.GeneratedRodsCount() < 4 {
		g.GenerateRod()
	}

	// Aptitudes/Mutations
	r = 5*g.Player.AptitudeCount() - g.Depth + 2
	if r < 2 {
		r = 1
	}
	if RandInt(r) == 0 && g.Depth > 0 && g.Player.AptitudeCount() < 3 {
		apt, ok := g.RandomApt()
		if ok {
			g.ApplyAptitude(apt)
		}
	}

	// Stairs
	g.Stairs = make(map[position]bool)
	nstairs := 1 + RandInt(3)
	if g.Depth == g.MaxDepth() {
		nstairs = 1
	} else if g.Depth == g.MaxDepth()-1 && nstairs > 2 {
		nstairs = 2
	}
	for i := 0; i < nstairs; i++ {
		var pos position
		if g.Depth > 9 {
			pos = g.FreeCellForImportantStair()
		} else {
			pos = g.FreeCellForStatic()
		}
		g.Stairs[pos] = true
	}

	// Gold
	g.Gold = make(map[position]int)
	for i := 0; i < 5; i++ {
		pos := g.FreeCellForStatic()
		g.Gold[pos] = 1 + RandInt(g.Depth+g.Depth*g.Depth/10)
	}

	// initialize LOS
	if g.Depth == 0 {
		g.Print("You're in Hareka's Underground. Good luck! Press ? for help.")
	}
	if g.Depth == g.MaxDepth() {
		g.Print("You feel magic in the air. The way out is close.")
	}
	g.ComputeLOS()
	g.MakeMonstersAware()

	// recharge rods
	for r, props := range g.Player.Rods {
		if props.Charge < r.MaxCharge() {
			props.Charge += RandInt(1 + r.Rate())
		}
		if props.Charge > r.MaxCharge() {
			props.Charge = r.MaxCharge()
		}
	}

	// clouds
	g.Clouds = map[position]cloud{}

	// Events
	if g.Depth == 0 {
		g.Events = &eventQueue{}
		heap.Init(g.Events)
		heap.Push(g.Events, &simpleEvent{ERank: 0, EAction: PlayerTurn})
		heap.Push(g.Events, &simpleEvent{ERank: 50, EAction: HealPlayer})
		heap.Push(g.Events, &simpleEvent{ERank: 100, EAction: MPRegen})
	} else {
		g.CleanEvents()
	}
	for i, _ := range g.Monsters {
		heap.Push(g.Events, &monsterEvent{ERank: g.Turn + 1, EAction: MonsterTurn, NMons: i})
		heap.Push(g.Events, &monsterEvent{ERank: g.Turn + 50, EAction: HealMonster, NMons: i})
	}
}

func (g *game) CleanEvents() {
	evq := &eventQueue{}
	for g.Events.Len() > 0 {
		ev := heap.Pop(g.Events).(event)
		switch ev.(type) {
		case *monsterEvent:
		case *cloudEvent:
		default:
			heap.Push(evq, ev)
		}
	}
	g.Events = evq
}

func (g *game) GenEquip(eq equipable, data equipableData) {
	depthAdjust := data.minDepth - g.Depth
	var r int
	if depthAdjust >= 0 {
		r = RandInt(data.rarity * (depthAdjust + 1) * (depthAdjust + 1))
	} else {
		switch eq.(type) {
		case shield:
			if !g.GeneratedEquipables[eq] {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				r = RandInt(data.rarity * 2)
			}
		case armour:
			if !g.GeneratedEquipables[eq] && eq != Robe {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				r = RandInt(data.rarity * 2)
			}
		case weapon:
			if !g.SeenGoodWeapon() && eq != Dagger {
				r = data.FavorableRoll(-depthAdjust)
			} else {
				if g.Player.Weapon != Dagger {
					r = RandInt(data.rarity * 4)
				} else {
					r = RandInt(data.rarity * 2)
				}
			}
		default:
			// not reached
			r = RandInt(data.rarity)
		}
	}
	if r == 0 {
		pos := g.FreeCellForStatic()
		g.Equipables[pos] = eq
		g.GeneratedEquipables[eq] = true
	}

}

func (g *game) Descend(ev event) bool {
	if g.Depth >= g.MaxDepth() {
		g.Depth++
		// win
		g.RemoveSaveFile()
		return true
	}
	g.Print("You descend deeper in the dungeon.")
	g.Depth++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank(), EAction: PlayerTurn})
	g.InitLevel()
	g.Save()
	return false
}

func (g *game) AutoPlayer(ev event) bool {
	if g.Resting {
		if g.MonsterInLOS() == nil && (g.Player.HP < g.Player.HPMax() || g.Player.HasStatus(StatusExhausted)) {
			g.WaitTurn(ev)
			return true
		}
		g.Resting = false
	} else if g.Autoexploring {
		g.ui.AutoExploreStep(g)
		mons := g.MonsterInLOS()
		switch {
		case mons.Exists():
			g.Print("You stop exploring.")
		case g.AutoHalt:
			// stop exploring for other reasons
			g.Print("You stop exploring.")
		default:
			var n *node
			var b bool
			count := 0
			for {
				count++
				if count > 100 {
					// should not happen
					g.Print("Hm… something went wrong with auto-explore.")
					n = nil
					break
				}
				n, b = g.NextAuto()
				if b {
					g.BuildAutoexploreMap()
				} else {
					break
				}
			}
			if n != nil {
				g.MovePlayer(n.Pos, ev)
				return true
			}
			g.Print("You finished exploring.")
		}
		g.Autoexploring = false
	} else if g.AutoTarget != nil {
		if g.MoveToTarget(ev) {
			return true
		}
	}
	return false
}

func (g *game) EventLoop() {
loop:
	for {
		if g.Player.HP <= 0 {
			if g.Wizard {
				g.Player.HP = g.Player.HPMax()
			} else {
				g.ui.Death(g)
				g.RemoveSaveFile()
				break loop
			}
		}
		if g.Events.Len() == 0 {
			break loop
		}
		ev := heap.Pop(g.Events).(event)
		g.Turn = ev.Rank()
		ev.Action(g)
		if g.AutoNext {
			continue loop
		}
		if g.Quit {
			break loop
		}
	}
}
