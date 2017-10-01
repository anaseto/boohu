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

func (g *game) StoryPrint(s string) {
	g.Story = append(g.Story, fmt.Sprintf("Depth %2d|Turn %7.1f| %s", g.Depth, float64(g.Turn)/10, s))
}

func (g *game) StoryPrintf(format string, a ...interface{}) {
	g.Story = append(g.Story, fmt.Sprintf("Depth %2d|Turn %7.1f| %s", g.Depth, float64(g.Turn)/10, fmt.Sprintf(format, a...)))
}
