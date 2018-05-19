package main

import "fmt"

type logStyle int

const (
	logNormal logStyle = iota
	logCritic
	logPlayerHit
	logMonsterHit
	logSpecial
	logStatusEnd
	logError
)

type logEntry struct {
	Text  string
	Index int
	Tick  bool
	Style logStyle
}

func (e logEntry) String() string {
	return e.Text
}

func (g *game) Print(s string) {
	e := logEntry{Text: s, Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *game) PrintStyled(s string, style logStyle) {
	e := logEntry{Text: s, Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *game) Printf(format string, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex}
	g.PrintEntry(e)
}

func (g *game) PrintfStyled(format string, style logStyle, a ...interface{}) {
	e := logEntry{Text: fmt.Sprintf(format, a...), Index: g.LogIndex, Style: style}
	g.PrintEntry(e)
}

func (g *game) PrintEntry(e logEntry) {
	if e.Index == g.LogNextTick {
		e.Tick = true
	}
	g.Log = append(g.Log, e)
	g.LogIndex++
	if len(g.Log) > 10000 {
		g.Log = g.Log[5000:]
	}
}

func (g *game) StoryPrint(s string) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %7.1f| %s", g.Depth, float64(g.Turn)/10, s))
}

func (g *game) StoryPrintf(format string, a ...interface{}) {
	g.Stats.Story = append(g.Stats.Story, fmt.Sprintf("Depth %2d|Turn %7.1f| %s", g.Depth, float64(g.Turn)/10, fmt.Sprintf(format, a...)))
}

func (g *game) CrackSound() (text string) {
	switch RandInt(4) {
	case 0:
		text = "Crack!"
	case 1:
		text = "Crash!"
	case 2:
		text = "Crunch!"
	case 3:
		text = "Creak!"
	}
	return text
}

func (g *game) ExplosionSound() (text string) {
	switch RandInt(3) {
	case 0:
		text = "Bang!"
	case 1:
		text = "Pop!"
	case 2:
		text = "Boom!"
	}
	return text
}
