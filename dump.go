package main

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"
)

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

func (g *game) SortedKilledMonsters() monsSlice {
	var ms monsSlice
	for mk, p := range g.Stats.KilledMons {
		if p == 0 {
			continue
		}
		ms = append(ms, mk)
	}
	sort.Sort(ms)
	return ms
}

func (g *game) Dump() string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, " -- Boohu (stealth) version %s character file --\n\n", Version)
	if g.Wizard {
		fmt.Fprintf(buf, "**WIZARD MODE**\n")
	}
	if g.Player.HP > 0 && g.Depth == -1 {
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
	fmt.Fprintf(buf, "Miscellaneous:\n")
	fmt.Fprintf(buf, "You killed %d monsters.\n", g.Stats.Killed)
	fmt.Fprintf(buf, "You spent %.1f turns in the Underground.\n", float64(g.Turn)/10)
	maxDepth := Max(g.Depth, g.ExploredLevels)
	s := "s"
	if maxDepth == 1 {
		s = ""
	}
	fmt.Fprintf(buf, "You explored %d level%s out of %d.\n", maxDepth, s, MaxDepth)
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Last messages:\n")
	for i := len(g.Log) - 10; i < len(g.Log); i++ {
		if i >= 0 {
			fmt.Fprintf(buf, "%s\n", g.Log[i])
		}
	}
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Dungeon:\n")
	fmt.Fprintf(buf, "┌%s┐\n", strings.Repeat("─", DungeonWidth))
	buf.WriteString(g.DumpDungeon())
	fmt.Fprintf(buf, "└%s┘\n", strings.Repeat("─", DungeonWidth))
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, g.DumpedKilledMonsters())
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "Timeline:\n")
	fmt.Fprintf(buf, g.DumpStory())
	fmt.Fprintf(buf, "\n")
	g.DetailedStatistics(buf)
	return buf.String()
}

func (g *game) DetailedStatistics(w io.Writer) {
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, "Statistics:\n")
	fmt.Fprintf(w, "You had %d hits (%.1f per 100 turns), %d misses (%.1f), and %d moves (%.1f).\n",
		g.Stats.Hits, float64(g.Stats.Hits)*100/float64(g.Stats.Turns+1),
		g.Stats.Misses, float64(g.Stats.Misses)*100/float64(g.Stats.Turns+1),
		g.Stats.Moves, float64(g.Stats.Moves)*100/float64(g.Stats.Turns+1))
	fmt.Fprintf(w, "You got hit %d times, blocked %d times, and dodged %d times.\n", g.Stats.ReceivedHits, g.Stats.Blocks, g.Stats.Dodges)
	fmt.Fprintf(w, "You endured %d damage.\n", g.Stats.Damage)
	fmt.Fprintf(w, "You were lucky %d times.\n", g.Stats.TimesLucky)
	fmt.Fprintf(w, "You activated %d stones.\n", g.Stats.UsedStones)
	fmt.Fprintf(w, "There were %d fires.\n", g.Stats.Burns)
	fmt.Fprintf(w, "There were %d destroyed walls.\n", g.Stats.Digs)
	fmt.Fprintf(w, "You rested %d times (%d interruptions).\n", g.Stats.Rest, g.Stats.RestInterrupt)
	fmt.Fprintf(w, "You spent %.1f%% turns wounded.\n", float64(g.Stats.TWounded)*100/float64(g.Stats.Turns+1))
	fmt.Fprintf(w, "You spent %.1f%% turns with monsters in sight.\n", float64(g.Stats.TMonsLOS)*100/float64(g.Stats.Turns+1))
	fmt.Fprintf(w, "You spent %.1f%% turns wounded with monsters in sight.\n", float64(g.Stats.TMWounded)*100/float64(g.Stats.Turns+1))
	maxDepth := Max(g.Depth-1, g.ExploredLevels)
	if g.Player.HP <= 0 {
		maxDepth++
	}
	if maxDepth >= MaxDepth+1 {
		// should not happen
		maxDepth = -1
	}
	fmt.Fprintf(w, "\n")
	hfmt := "%-23s"
	fmt.Fprintf(w, hfmt, "Quantity/Depth")
	for i := 1; i <= maxDepth; i++ {
		fmt.Fprintf(w, " %3d", i)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, hfmt, "Explored (%)")
	for i, n := range g.Stats.DExplPerc {
		if i == 0 {
			continue
		}
		if i > maxDepth {
			break
		}
		fmt.Fprintf(w, " %3d", n)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, hfmt, "Sleeping monsters (%)")
	for i, n := range g.Stats.DSleepingPerc {
		if i == 0 {
			continue
		}
		if i > maxDepth {
			break
		}
		fmt.Fprintf(w, " %3d", n)
	}
	fmt.Fprintf(w, "\n")
	fmt.Fprintf(w, hfmt, "Dead monsters (%)")
	for i, n := range g.Stats.DKilledPerc {
		if i == 0 {
			continue
		}
		if i > maxDepth {
			break
		}
		fmt.Fprintf(w, " %3d", n)
	}
	fmt.Fprintf(w, "\n")
	//fmt.Fprintf(w, "Legend:")
	//for i, c := range []dungen{GenCaveMap, GenRoomMap, GenCellularAutomataCaveMap, GenCaveMapTree, GenRuinsMap, GenBSPMap} {
	//if i == 4 {
	//fmt.Fprintf(w, "\n       ")
	//}
	//fmt.Fprintf(w, " %s (%s)", c.Description(), c.String())
	//}
}

func (g *game) DumpStory() string {
	return strings.Join(g.Stats.Story, "\n")
}

func (g *game) DumpDungeon() string {
	buf := bytes.Buffer{}
	for i, c := range g.Dungeon.Cells {
		if i%DungeonWidth == 0 {
			if i == 0 {
				buf.WriteRune('│')
			} else {
				buf.WriteString("│\n│")
			}
		}
		pos := idxtopos(i)
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
		default:
			switch {
			case pos == g.Player.Pos:
				r = '@'
			default:
				r, _ = c.Style(g, pos)
				if _, ok := g.Clouds[pos]; ok && g.Player.LOS[pos] {
					r = '§'
				}
				m := g.MonsterAt(pos)
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
		fmt.Fprintf(buf, "- %s: %d\n", mk, g.Stats.KilledMons[mk])
	}
	return buf.String()
}

func (g *game) SimplifedDump(err error) string {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, " ♣ Boohu (stealth) version %s play summary ♣\n\n", Version)
	if g.Wizard {
		fmt.Fprintf(buf, "**WIZARD MODE**\n")
	}
	if g.Player.HP > 0 && g.Depth == -1 {
		fmt.Fprintf(buf, "You escaped from Hareka's Underground alive!\n")
	} else if g.Player.HP <= 0 {
		fmt.Fprintf(buf, "You died while exploring depth %d of Hareka's Underground.\n", g.Depth)
	} else {
		fmt.Fprintf(buf, "You are exploring depth %d of Hareka's Underground.\n", g.Depth)
	}
	fmt.Fprintf(buf, "You collected %d simellas.\n", g.Player.Simellas)
	fmt.Fprintf(buf, "You killed %d monsters.\n", g.Stats.Killed)
	fmt.Fprintf(buf, "You spent %.1f turns in the Underground.\n", float64(g.Turn)/10)
	maxDepth := Max(g.Depth, g.ExploredLevels)
	s := "s"
	if maxDepth == 1 {
		s = ""
	}
	fmt.Fprintf(buf, "You explored %d level%s out of %d.\n", maxDepth, s, MaxDepth+1)
	fmt.Fprintf(buf, "\n")
	if err != nil {
		fmt.Fprintf(buf, "Error writing dump: %v.\n", err)
	} else {
		dataDir, err := g.DataDir()
		if err == nil {
			if dataDir == "" {
				fmt.Fprintf(buf, "Full game statistics written below.\n")
			} else {
				fmt.Fprintf(buf, "Full game statistics dump written to %s.\n", filepath.Join(dataDir, "dump"))
			}
		}
	}
	fmt.Fprintf(buf, "\n\n")
	fmt.Fprintf(buf, "───Press esc or space to quit───")
	return buf.String()
}
