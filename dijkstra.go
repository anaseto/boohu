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

func (g *game) drawDijkstra(nm nodeMap) string {
	b := &bytes.Buffer{}
	for y := 0; y < g.Dungeon.Heigth; y++ {
		for x := 0; x < g.Dungeon.Width; x++ {
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
