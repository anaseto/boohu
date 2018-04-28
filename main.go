// +build !js

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

var Version string = "v0.6.1"

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to `file`")
	opt := flag.Bool("s", false, "Use true 16-uicolor solarized palette")
	optVersion := flag.Bool("v", false, "print version number")
	optCenteredCamera := flag.Bool("c", false, "centered camera")
	flag.Parse()
	if *cpuprofile != "" {
		// profiling
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	if *opt {
		SolarizedPalette()
	}
	if *optVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if *optCenteredCamera {
		CenteredCamera = true
	}

	tui := &termui{}
	err := tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()

	if runtime.GOOS == "windows" {
		WindowsPalette()
	}
	tui.PostInit()

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
