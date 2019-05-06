package main

import "testing"

func TestInitLevel(t *testing.T) {
	Testing = true
	for i := 0; i < 50; i++ {
		g := &game{}
		for depth := 0; depth < 11; depth++ {
			g.InitLevel()
			g.Depth++
			for _, m := range g.Monsters {
				if !g.Dungeon.Cell(m.Pos).IsPassable() {
					t.Errorf("Not free: %+v", m.Pos)
				}
			}
		}
	}
}
