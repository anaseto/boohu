// +build !ebiten

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var Version string = "v0.4"

func main() {
	opt := flag.Bool("s", false, "Use true 16-uicolor solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	flag.Parse()
	if *opt {
		SolarizedPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}

	tui := &termui{}
	err := tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	tui.PostInit()
	if runtime.GOOS == "windows" {
		WindowsPalette()
	}

	tui.DrawWelcome()
	g := &game{}
	load, err := g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.Print("Error loading saved game… starting new game.")
	}
	g.ui = tui
	g.EventLoop()
}
