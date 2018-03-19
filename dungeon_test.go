package main

import "testing"

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
			t.Errorf("Not connex: %+v\n", g.Dungeon.Cells)
		}
	}
}

func TestCaveMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex: %+v\n", g.Dungeon.Cells)
		}
	}
}

func TestCaveMapTree(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenCaveMapTree(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex: %+v\n", g.Dungeon.Cells)
		}
	}
}

func TestRuinsMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenRuinsMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex: %+v\n", g.Dungeon.Cells)
		}
	}
}

func TestRoomMap(t *testing.T) {
	for i := 0; i < 100; i++ {
		g := &game{}
		g.GenRoomMap(21, 79)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex: %+v\n", g.Dungeon.Cells)
		}
	}
}
