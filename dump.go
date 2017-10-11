package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

type monsSlice []monsterKind

func (ms monsSlice) Len() int      { return len(ms) }
func (ms monsSlice) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }
func (ms monsSlice) Less(i, j int) bool {
	return ms[i].Dangerousness() > ms[j].Dangerousness()
}

func (g *game) DumpAptitudes() string {
	apts := []string{}
	for apt, b := range g.Player.Aptitudes {
		if b {
			apts = append(apts, apt.String())
		}
	}
	sort.Strings(apts)
	if len(apts) == 0 {
		return "You do not have any special aptitudes."
	}
	return "Aptitudes:\n" + strings.Join(apts, "\n")
}

func (g *game) DumpStatuses() string {
	sts := sort.StringSlice{}
	for st, c := range g.Player.Statuses {
		if c > 0 {
			sts = append(sts, st.String())
		}
	}
	sort.Sort(sts)
	if len(sts) == 0 {
		return "You are free of any status effects."
	}
	return "Statuses:\n" + strings.Join(sts, "\n")
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

func (g *game) SortedKilledMonsters() monsSlice {
	var ms monsSlice
	for mk, p := range g.KilledMons {
		if p == 0 {
			continue
		}
		ms = append(ms, mk)
	}
	sort.Sort(ms)
	return ms
}

func (g *game) SortedPotions() consumableSlice {
	var cs consumableSlice
	for k := range g.Player.Consumables {
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
	for k := range g.Player.Consumables {
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
	fmt.Fprintf(buf, " -- Boohu version %s character file --\n\n", Version)
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
	fmt.Fprintf(buf, "You have %d/%d HP, and %d/%d MP.\n", g.Player.HP, g.Player.HPMax(), g.Player.MP, g.Player.MPMax())
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, g.DumpAptitudes())
	fmt.Fprintf(buf, "\n\n")
	fmt.Fprintf(buf, g.DumpStatuses())
	fmt.Fprintf(buf, "\n\n")
	fmt.Fprintf(buf, "Equipment:\n")
	fmt.Fprintf(buf, "You are wearing %s.\n", Indefinite(g.Player.Armour.String(), false))
	fmt.Fprintf(buf, "You are wielding %s.\n", Indefinite(g.Player.Weapon.String(), false))
	if g.Player.Shield != NoShield {
		if g.Player.Weapon.TwoHanded() {
			fmt.Fprintf(buf, "You have %s (unused).\n", Indefinite(g.Player.Shield.String(), false))
		} else {
			fmt.Fprintf(buf, "You are wearing %s.\n", Indefinite(g.Player.Shield.String(), false))
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
	for i := len(g.Log) - 10; i < len(g.Log); i++ {
		if i >= 0 {
			fmt.Fprintf(buf, "%s\n", g.Log[i])
		}
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Dungeon:\n")
	fmt.Fprintf(buf, "┌%s┐\n", strings.Repeat("─", g.Dungeon.Width))
	buf.WriteString(g.DumpDungeon())
	fmt.Fprintf(buf, "└%s┘\n", strings.Repeat("─", g.Dungeon.Width))
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, g.DumpedKilledMonsters())
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Timeline:\n")
	fmt.Fprintf(buf, g.DumpStory())
	return buf.String()
}

func (g *game) DumpStory() string {
	return strings.Join(g.Story, "\n")
}

func (g *game) DumpDungeon() string {
	buf := bytes.Buffer{}
	for i, c := range g.Dungeon.Cells {
		if i%g.Dungeon.Width == 0 {
			if i == 0 {
				buf.WriteRune('│')
			} else {
				buf.WriteString("│\n│")
			}
		}
		pos := g.Dungeon.CellPosition(i)
		if !c.Explored {
			buf.WriteRune(' ')
			if i == len(g.Dungeon.Cells)-1 {
				buf.WriteString("│\n")
			}
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
		if i == len(g.Dungeon.Cells)-1 {
			buf.WriteString("│\n")
		}
	}
	return buf.String()
}

func (g *game) DumpedKilledMonsters() string {
	buf := &bytes.Buffer{}
	fmt.Fprint(buf, "Killed Monsters:\n")
	ms := g.SortedKilledMonsters()
	for _, mk := range ms {
		fmt.Fprintf(buf, "- %s: %d\n", mk, g.KilledMons[mk])
	}
	return buf.String()
}

func (g *game) SimplifedDump() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, " ♣ Boohu version %s play summary ♣\n\n", Version)
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
