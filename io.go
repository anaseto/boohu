// +build !js

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func Replay(file string) error {
	tui := &gameui{}
	g := &game{}
	tui.g = g
	g.ui = tui
	err := g.LoadReplay(file)
	if err != nil {
		return fmt.Errorf("loading replay: %v", err)
	}
	err = tui.Init()
	if err != nil {
		fmt.Fprintf(os.Stderr, "boohu: %v\n", err)
		os.Exit(1)
	}
	defer tui.Close()
	tui.DrawBufferInit()
	tui.Replay()
	tui.Close()
	return nil
}

func (g *game) DataDir() (string, error) {
	var xdg string
	if os.Getenv("GOOS") == "windows" {
		xdg = os.Getenv("LOCALAPPDATA")
	} else {
		xdg = os.Getenv("XDG_DATA_HOME")
	}
	if xdg == "" {
		xdg = filepath.Join(os.Getenv("HOME"), ".local", "share")
	}
	dataDir := filepath.Join(xdg, "boohu")
	_, err := os.Stat(dataDir)
	if err != nil {
		err = os.MkdirAll(dataDir, 0755)
		if err != nil {
			return "", fmt.Errorf("%v\n", err)
		}
	}
	return dataDir, nil
}

func (g *game) Save() error {
	dataDir, err := g.DataDir()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	saveFile := filepath.Join(dataDir, "save")
	data, err := g.GameSave()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		g.Print(err.Error())
		return err
	}
	return nil
}

func (g *game) RemoveSaveFile() error {
	return g.RemoveDataFile("save")
}

func (g *game) Load() (bool, error) {
	dataDir, err := g.DataDir()
	if err != nil {
		return false, err
	}
	saveFile := filepath.Join(dataDir, "save")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new game
		return false, err
	}
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return true, err
	}
	lg, err := g.DecodeGameSave(data)
	if err != nil {
		return true, err
	}
	if lg.Version != Version {
		return true, fmt.Errorf("saved game for previous version %s.", lg.Version)
	}
	*g = *lg
	return true, nil
}

func (g *game) SaveConfig() error {
	dataDir, err := g.DataDir()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	saveFile := filepath.Join(dataDir, "config.gob")
	data, err := gameConfig.ConfigSave()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		g.Print(err.Error())
		return err
	}
	return nil
}

func (g *game) LoadConfig() (bool, error) {
	dataDir, err := g.DataDir()
	if err != nil {
		return false, err
	}
	saveFile := filepath.Join(dataDir, "config.gob")
	_, err = os.Stat(saveFile)
	if err != nil {
		// no save file, new game
		return false, err
	}
	data, err := ioutil.ReadFile(saveFile)
	if err != nil {
		return true, err
	}
	c, err := g.DecodeConfigSave(data)
	if err != nil {
		return true, err
	}
	gameConfig = *c
	if gameConfig.RuneNormalModeKeys == nil || gameConfig.RuneTargetModeKeys == nil {
		ApplyDefaultKeyBindings()
	}
	if gameConfig.DarkLOS {
		ApplyDarkLOS()
	}
	if gameConfig.Small {
		UIHeight = 24
		UIWidth = 80
	}
	return true, nil
}

func (g *game) RemoveDataFile(file string) error {
	dataDir, err := g.DataDir()
	if err != nil {
		return err
	}
	dataFile := filepath.Join(dataDir, file)
	_, err = os.Stat(dataFile)
	if err == nil {
		err := os.Remove(dataFile)
		if err != nil {
			return err
		}
	}
	return nil
}

func (g *game) SaveReplay() error {
	dataDir, err := g.DataDir()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	saveFile := filepath.Join(dataDir, "replay")
	data, err := g.EncodeDrawLog()
	if err != nil {
		g.Print(err.Error())
		return err
	}
	err = ioutil.WriteFile(saveFile, data, 0644)
	if err != nil {
		g.Print(err.Error())
		return err
	}
	return nil
}

func (g *game) LoadReplay(file string) error {
	dataDir, err := g.DataDir()
	if err != nil {
		return err
	}
	replayFile := filepath.Join(dataDir, "replay")
	if file != "" {
		replayFile = file
	}
	_, err = os.Stat(replayFile)
	if err != nil {
		// no save file, new game
		return err
	}
	data, err := ioutil.ReadFile(replayFile)
	if err != nil {
		return err
	}
	dl, err := g.DecodeDrawLog(data)
	if err != nil {
		return err
	}
	g.DrawLog = dl
	return nil
}

func (g *game) WriteDump() error {
	dataDir, err := g.DataDir()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dataDir, "dump"), []byte(g.Dump()), 0644)
	if err != nil {
		return fmt.Errorf("writing game statistics: %v", err)
	}
	err = g.SaveReplay()
	if err != nil {
		return fmt.Errorf("writing replay: %v", err)
	}
	return nil
}
