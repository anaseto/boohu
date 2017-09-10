package main

import (
	"bytes"
	"container/heap"
	"fmt"
)

type Dijkstrer interface {
	Neighbors(position) []position
	Cost(position, position) int
}

type normalPath struct {
	game *game
}

func (np *normalPath) Neighbors(pos position) []position {
	if np.game.Player.HasStatus(StatusConfusion) {
		return np.game.Dungeon.CardinalFreeNeighbors(pos)
	} else {
		return np.game.Dungeon.FreeNeighbors(pos)
	}
}

func (np *normalPath) Cost(from, to position) int {
	return 1
}

func (g *game) drawDijkstra(nm nodeMap) string {
	b := &bytes.Buffer{}
	for y := 0; y < 25; y++ {
		for x := 0; x < 70; x++ {
			pos := position{x, y}
			n, ok := nm[pos]
			if ok {
				if pos == g.Player.Pos {
					fmt.Fprintf(b, "%d@", n.Cost)
				} else {
					c := g.Dungeon.Cell(pos)
					if c.T == WallCell {
						fmt.Fprintf(b, "%d#", n.Cost)
					} else {
						fmt.Fprintf(b, "%d ", n.Cost)
					}
				}
			} else {
				fmt.Fprintf(b, "# ")
			}
		}
		fmt.Fprintf(b, "\n")
	}
	return b.String()
}

func Dijkstra(dij Dijkstrer, sources []position, maxCost int) nodeMap {
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	for _, f := range sources {
		n := nm.get(f)
		n.Open = true
		heap.Push(nq, n)
	}
	for {
		if nq.Len() == 0 {
			return nm
		}
		current := heap.Pop(nq).(*node)
		current.Open = false
		current.Closed = true

		for _, neighbor := range dij.Neighbors(current.Pos) {
			cost := current.Cost + dij.Cost(current.Pos, neighbor)
			neighborNode := nm.get(neighbor)
			if cost < neighborNode.Cost {
				if neighborNode.Open {
					heap.Remove(nq, neighborNode.Index)
				}
				neighborNode.Open = false
				neighborNode.Closed = false
			}
			if !neighborNode.Open && !neighborNode.Closed {
				neighborNode.Cost = cost
				if cost < maxCost {
					neighborNode.Open = true
					heap.Push(nq, neighborNode)
				}
			}
		}
	}
}
