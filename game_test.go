package main

import "testing"

func TestInitLevel(t *testing.T) {
	for i := 0; i < 10; i++ {
		g := &game{}
		for depth := 0; depth < 11; depth++ {
			g.Depth = depth
			g.InitLevel()
			for _, m := range g.Monsters {
				if !g.Dungeon.Cell(m.Pos).IsPassable() {
					t.Errorf("Not free: %+v", m.Pos)
				}
			}
		}
	}
}
