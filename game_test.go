package main

import "testing"

func TestInitLevel(t *testing.T) {
	for i := 0; i < 10; i++ {
		g := &game{}
		for depth := 0; depth < 13; depth++ {
			g.Depth = depth
			g.InitLevel()
		}
	}
}
