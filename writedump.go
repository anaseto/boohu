// +build !ebiten

package main

import (
	"io/ioutil"
	"path/filepath"
)

func (g *game) WriteDump() error {
	dataDir, err := g.DataDir()
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filepath.Join(dataDir, "dump"), []byte(g.Dump()), 0644)
	if err != nil {
		return err
	}
	return nil
}
