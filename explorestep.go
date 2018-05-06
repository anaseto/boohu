// +build !ansi,!js

package main

import "time"

func (ui *termui) ExploreStep(g *game) bool {
	next := make(chan bool)
	var stop bool
	//if runtime.GOOS != "windows" {
	// strange bugs it seems, cannot test myself, so disable on windows
	// -> Enable this again, works at least on one particular windows cmd
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
	//} else {
	//time.Sleep(10 * time.Millisecond)
	//}
	ui.DrawDungeonView(g, NormalMode)
	return stop
}
