package main

import (
	"fmt"
	"os"
	"time"
)

func Replay(file string) error {
	tui := &termui{}
	g := &game{}
	tui.g = g
	g.ui = tui
	err := g.LoadReplay()
	if err != nil {
		return fmt.Errorf("loading replay: %v", err)
	}
	err = tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()
	tui.PostInit()
	tui.DrawBufferInit()
	tui.Replay()
	tui.Close()
	return nil
}

func (ui *termui) Replay() {
	g := ui.g
	dl := g.DrawLog
	if len(dl) == 0 {
		return
	}
	g.DrawLog = nil
	rep := &replay{ui: ui, frames: dl, frame: 0}
	rep.Run()
}

func (rep *replay) Run() {
	rep.auto = true
	rep.speed = 1
	rep.evch = make(chan repEvent, 100)
	go func(r *replay) {
		r.PollKeyboardEvents()
	}(rep)
	for {
		e := rep.PollEvent()
		switch e {
		case ReplayNext:
			if rep.frame >= len(rep.frames) {
				break
			} else if rep.frame < 0 {
				rep.frame = 0
			}
			rep.DrawFrame(rep.frame)
			rep.frame++
		case ReplayQuit:
			return
		case ReplayTogglePause:
			rep.auto = !rep.auto
		case ReplaySpeedMore:
			rep.speed++
			if rep.speed > 15 {
				rep.speed = 15
			}
		case ReplaySpeedLess:
			rep.speed--
			if rep.speed < 1 {
				rep.speed = 1
			}
		}
	}
}

func (rep *replay) DrawFrame(i int) {
	ui := rep.ui
	df := rep.frames[i]
	for _, dr := range df.Draws {
		x, y := ui.GetPos(dr.I)
		ui.SetGenCell(x, y, dr.Cell.R, dr.Cell.Fg, dr.Cell.Bg, dr.Cell.InMap)
	}
	ui.Flush()
}

type replay struct {
	ui     *termui
	frames []drawFrame
	undo   [][]cellDraw
	frame  int
	auto   bool
	speed  time.Duration
	evch   chan repEvent
}

type repEvent int

const (
	ReplayNext repEvent = iota
	ReplayPrevious
	ReplayTogglePause
	ReplayQuit
	ReplaySpeedMore
	ReplaySpeedLess
)

func (rep *replay) PollEvent() (in repEvent) {
	if rep.auto && rep.frame < len(rep.frames)-1 && rep.frame >= 0 {
		d := rep.frames[rep.frame+1].Time.Sub(rep.frames[rep.frame].Time)
		if d >= 2*time.Second {
			d = 2 * time.Second
		}
		d = d / rep.speed
		if d <= 10*time.Millisecond {
			d = 10 * time.Millisecond
		}
		t := time.NewTimer(d / rep.speed)
		select {
		case in = <-rep.evch:
		case <-t.C:
			in = ReplayNext
		}
		t.Stop()
	} else {
		in = <-rep.evch
	}
	return in
}

func (rep *replay) PollKeyboardEvents() {
	for {
		e := rep.ui.PollEvent()
		if e.interrupt {
			rep.evch <- ReplayNext
			continue
		}
		switch e.key {
		case "Q":
			rep.evch <- ReplayQuit
		case "p", " ":
			rep.evch <- ReplayTogglePause
		case "+", ">":
			rep.evch <- ReplaySpeedMore
		case "-", "<":
			rep.evch <- ReplaySpeedLess
		case ".", "6", "j", "n":
			rep.evch <- ReplayNext
		case "4", "k", "N":
			rep.evch <- ReplayPrevious
		default:
			if !e.mouse {
				break
			}
			switch e.button {
			case 0:
				rep.evch <- ReplayNext
			case 2:
				rep.evch <- ReplayPrevious
			}
		}
	}
}
