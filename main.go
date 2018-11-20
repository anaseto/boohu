// +build !js,!ansi,!tk

package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func main() {
	optSolarized := flag.Bool("s", false, "Use true 16-color solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	optCenteredCamera := flag.Bool("c", false, "centered camera")
	color8 := false
	if runtime.GOOS == "windows" {
		color8 = true
	}
	opt8colors := flag.Bool("o", color8, "use only 8-color palette")
	opt256colors := flag.Bool("x", !color8, "use xterm 256-color palette (solarized approximation)")
	optNoAnim := flag.Bool("n", false, "no animations")
	flag.Parse()
	if *optSolarized {
		SolarizedPalette()
	} else if color8 && !*opt256colors || !color8 && *opt8colors {
		Simple8ColorPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optCenteredCamera {
		CenteredCamera = true
	}
	if *optNoAnim {
		DisableAnimations = true
	}

	tui := &termui{}
	err := tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	ApplyDefaultKeyBindings()
	tui.PostInit()
	LinkColors()

	tui.DrawWelcome()
	g := &game{}
	load, err := g.Load()
	if !load {
		g.InitLevel()
	} else if err != nil {
		g.InitLevel()
		g.PrintfStyled("Error: %v", logError, err)
		g.PrintStyled("Could not load saved game… starting new game.", logError)
	}
	load, err = g.LoadConfig()
	if load && err != nil {
		g.Print("Error loading config file.")
	} else if load {
		CustomKeys = true
	}

	g.ui = tui
	g.EventLoop()
}
