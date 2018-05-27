package main

import (
	"bytes"
	"fmt"
	"testing"
)

func BenchmarkCellularAutomataCaveMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		g := &game{}
		g.GenCellularAutomataCaveMap(DungeonHeight, DungeonWidth)
	}
}

func TestCellularAutomataCaveMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCellularAutomataCaveMap(DungeonHeight, DungeonWidth)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestCaveMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMap(DungeonHeight, DungeonWidth)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestCaveMapTree(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMapTree(DungeonHeight, DungeonWidth)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRuinsMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenRuinsMap(DungeonHeight, DungeonWidth)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestBSPMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenBSPMap(DungeonHeight, DungeonWidth)
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
		g.GenRoomMap(DungeonHeight, DungeonWidth)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}
