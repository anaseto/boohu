package main

import (
	"bytes"
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
	return data.Bytes(), nil
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
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	lg := &game{}
	err := dec.Decode(lg)
	if err != nil {
		return nil, err
	}
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
