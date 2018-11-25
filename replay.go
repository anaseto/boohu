package main

import (
	"time"
)

func (ui *gameui) Replay() {
	g := ui.g
	dl := g.DrawLog
	if len(dl) == 0 {
		return
	}
	g.DrawLog = nil
	rep := &replay{ui: ui, frames: dl, frame: 0}
	rep.Run()
}

type replay struct {
	ui     *gameui
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

func (rep *replay) Run() {
	rep.auto = true
	rep.speed = 1
	rep.evch = make(chan repEvent, 100)
	rep.undo = [][]cellDraw{}
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
			rep.DrawFrame()
			rep.frame++
		case ReplayPrevious:
			if rep.frame <= 1 {
				break
			} else if rep.frame >= len(rep.frames) {
				rep.frame = len(rep.frames)
			}
			rep.frame--
			rep.UndoFrame()
		case ReplayQuit:
			return
		case ReplayTogglePause:
			rep.auto = !rep.auto
		case ReplaySpeedMore:
			rep.speed *= 2
			if rep.speed > 16 {
				rep.speed = 16
			}
		case ReplaySpeedLess:
			rep.speed /= 2
			if rep.speed < 1 {
				rep.speed = 1
			}
		}
	}
}

func (rep *replay) DrawFrame() {
	ui := rep.ui
	df := rep.frames[rep.frame]
	rep.undo = append(rep.undo, []cellDraw{})
	j := len(rep.undo) - 1
	for _, dr := range df.Draws {
		i := ui.GetIndex(dr.X, dr.Y)
		c := ui.g.DrawBuffer[i]
		rep.undo[j] = append(rep.undo[j], cellDraw{Cell: c, X: dr.X, Y: dr.Y})
		ui.SetGenCell(dr.X, dr.Y, dr.Cell.R, dr.Cell.Fg, dr.Cell.Bg, dr.Cell.InMap)
	}
	ui.Flush()
	ui.g.DrawLog = nil
}

func (rep *replay) UndoFrame() {
	ui := rep.ui
	df := rep.undo[len(rep.undo)-1]
	for _, dr := range df {
		ui.SetGenCell(dr.X, dr.Y, dr.Cell.R, dr.Cell.Fg, dr.Cell.Bg, dr.Cell.InMap)
	}
	rep.undo = rep.undo[:len(rep.undo)-1]
	ui.Flush()
	ui.g.DrawLog = nil
}

func (rep *replay) PollEvent() (in repEvent) {
	if rep.auto && rep.frame <= len(rep.frames)-1 && rep.frame >= 0 {
		var d time.Duration
		if rep.frame > 0 {
			d = rep.frames[rep.frame].Time.Sub(rep.frames[rep.frame-1].Time)
		} else {
			d = 0
		}
		if d >= 2*time.Second {
			d = 2 * time.Second
		}
		d = d / rep.speed
		if d <= 10*time.Millisecond {
			d = 10 * time.Millisecond
		}
		t := time.NewTimer(d)
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
			return
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
			case 1:
				rep.evch <- ReplayTogglePause
			case 2:
				rep.evch <- ReplayPrevious
			}
		}
	}
}
