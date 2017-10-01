package main

import "fmt"

func (g *game) Print(s string) {
	g.Log = append(g.Log, s)
	if len(g.Log) > 1000 {
		g.Log = g.Log[500:]
	}
}

func (g *game) Printf(format string, a ...interface{}) {
	g.Log = append(g.Log, fmt.Sprintf(format, a...))
	if len(g.Log) > 1000 {
		g.Log = g.Log[500:]
	}
}
