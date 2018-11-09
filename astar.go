// code of this file is a modified version of code from
// github.com/beefsack/go-astar, which has the following license:
//
// Copyright (c) 2014 Michael Charles Alexander
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"container/heap"
)

type node struct {
	Pos    position
	Cost   int
	Rank   int
	Parent *position
	Open   bool
	Closed bool
	Index  int
	Num    int
}

type nodeMap map[position]*node

var nodeCache []node

func init() {
	nodeCache = make([]node, 0, DungeonNCells)
}

func (nm nodeMap) get(p position) *node {
	n, ok := nm[p]
	if !ok {
		nodeCache = append(nodeCache, node{Pos: p})
		n = &nodeCache[len(nodeCache)-1]
		nm[p] = n
	}
	return n
}

type Astar interface {
	Neighbors(position) []position
	Cost(position, position) int
	Estimation(position, position) int
}

func AstarPath(ast Astar, from, to position) (path []position, length int, found bool) {
	nodeCache = nodeCache[:0]
	nm := nodeMap{}
	nq := &priorityQueue{}
	heap.Init(nq)
	fromNode := nm.get(from)
	fromNode.Open = true
	num := 0
	fromNode.Num = num
	heap.Push(nq, fromNode)
	for {
		if nq.Len() == 0 {
			// There's no path, return found false.
			return
		}
		current := heap.Pop(nq).(*node)
		current.Open = false
		current.Closed = true

		if current.Pos == to {
			// Found a path to the goal.
			p := []position{}
			curr := current
			for {
				p = append(p, curr.Pos)
				if curr.Parent == nil {
					break
				}
				curr = nm[*curr.Parent]
			}
			return p, current.Cost, true
		}

		for _, neighbor := range ast.Neighbors(current.Pos) {
			cost := current.Cost + ast.Cost(current.Pos, neighbor)
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
				neighborNode.Open = true
				neighborNode.Rank = cost + ast.Estimation(neighbor, to)
				neighborNode.Parent = &current.Pos
				num++
				neighborNode.Num = num
				heap.Push(nq, neighborNode)
			}
		}
	}
}

// A priorityQueue implements heap.Interface and holds Nodes.  The
// priorityQueue is used to track open nodes by rank.
type priorityQueue []*node

func (pq priorityQueue) Len() int {
	return len(pq)
}

func (pq priorityQueue) Less(i, j int) bool {
	//return pq[i].Rank < pq[j].Rank
	return pq[i].Rank < pq[j].Rank || pq[i].Rank == pq[j].Rank && pq[i].Num < pq[j].Num
}

func (pq priorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *priorityQueue) Push(x interface{}) {
	n := len(*pq)
	no := x.(*node)
	no.Index = n
	*pq = append(*pq, no)
}

func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	no := old[n-1]
	no.Index = -1
	*pq = old[0 : n-1]
	return no
}
