package main

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
)

func init() {
	gob.Register(potion(0))
	gob.Register(projectile(0))
	gob.Register(&simpleEvent{})
	gob.Register(&monsterEvent{})
	gob.Register(&cloudEvent{})
	gob.Register(armour(0))
	gob.Register(weapon(0))
	gob.Register(shield(0))
}

func (g *game) GameSave() ([]byte, error) {
	data := bytes.Buffer{}
	enc := gob.NewEncoder(&data)
	err := enc.Encode(g)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	w.Write(data.Bytes())
	w.Close()
	return buf.Bytes(), nil
}

type config struct {
	RuneNormalModeKeys map[rune]keyAction
	RuneTargetModeKeys map[rune]keyAction
	DarkLOS            bool
	Small              bool
	Tiles              bool
	Version            string
}

func (c *config) ConfigSave() ([]byte, error) {
	data := bytes.Buffer{}
	enc := gob.NewEncoder(&data)
	err := enc.Encode(c)
	if err != nil {
		return nil, err
	}
	return data.Bytes(), nil
}

func (g *game) DecodeGameSave(data []byte) (*game, error) {
	buf := bytes.NewReader(data)
	r, err := zlib.NewReader(buf)
	if err != nil {
		return nil, err
	}
	dec := gob.NewDecoder(r)
	lg := &game{}
	err = dec.Decode(lg)
	if err != nil {
		return nil, err
	}
	r.Close()
	return lg, nil
}

func (g *game) DecodeConfigSave(data []byte) (*config, error) {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	c := &config{}
	err := dec.Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
