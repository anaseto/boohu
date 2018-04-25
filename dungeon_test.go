package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkCellularAutomataCaveMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := &game{}
		g.GenCellularAutomataCaveMap(21, 79)
	}
}

func TestCellularAutomataCaveMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCellularAutomataCaveMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestCaveMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestCaveMapTree(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMapTree(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRuinsMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenRuinsMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func (d *dungeon) String() string {
	b := &bytes.Buffer{}
	for i, c := range d.Cells {
		if i > 0 && i%DungeonWidth == 0 {
			fmt.Fprint(b, "\n")
		}
		if c.T == WallCell {
			fmt.Fprint(b, "#")
		} else {
			fmt.Fprint(b, ".")
		}
	}
	return b.String()
}

func TestRoomMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenRoomMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}
