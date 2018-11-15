package main

import "time"

func (ui *termui) ExploreStep(g *game) bool {
	next := make(chan bool)
	var stop bool
	go func() {
		time.Sleep(10 * time.Millisecond)
		ui.Interrupt()
	}()
	go func() {
		err := ui.PressAnyKey()
		interrupted := err != nil
		next <- !interrupted
	}()
	stop = <-next
	ui.DrawDungeonView(g, NormalMode)
	return stop
}
