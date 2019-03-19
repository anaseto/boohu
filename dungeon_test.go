package main

import (
	"bytes"
	"fmt"
	"testing"
)

var Rounds = 100

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

func TestAutomataCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(AutomataCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}

func TestRandomWalkTreeCave(t *testing.T) {
	for i := 0; i < Rounds; i++ {
		g := &game{}
		g.InitFirstLevel()
		g.InitLevelStructures()
		g.GenRoomTunnels(RandomWalkTreeCave)
		if !g.Dungeon.connex() {
			t.Errorf("Not connex:\n%s\n", g.Dungeon.String())
		}
	}
}
