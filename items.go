package main

import (
	"container/heap"
	"errors"
	"fmt"
)

// + consumables (potion-like or throwing dart, strategic + tactical)
// + equipables
// + recharging with depth (rod-like, strategic & a little tactical + mana)
//   - digging, fog, slowing clouds or something, fear,
//     fireball, lightning bolt, shatter, blink, teleport other

type consumable interface {
	Use(*game, event) error
	String() string
	Letter() rune
	Int() int
}

func (g *game) UseConsumable(c consumable) {
	g.Player.Consumables[c]--
	if g.Player.Consumables[c] <= 0 {
		delete(g.Player.Consumables, c)
	}
}

type potion int

const (
	HealWoundsPotion potion = iota
	TeleportationPotion
	BerserkPotion
	DescentPotion
	RunningPotion
	EvasionPotion
	LignificationPotion
	MagicMappingPotion
	MagicPotion
	// below unimplemented
	ResistancePotion
)

func (p potion) String() (text string) {
	switch p {
	case HealWoundsPotion:
		text = "potion of heal wounds"
	case TeleportationPotion:
		text = "potion of teleportation"
	case DescentPotion:
		text = "potion of descent"
	case EvasionPotion:
		text = "potion of evasion"
	case MagicMappingPotion:
		text = "potion of magic mapping"
	case MagicPotion:
		text = "potion of refill magic"
	case ResistancePotion:
		text = "potion of resistance"
	case BerserkPotion:
		text = "potion of berserk"
	case RunningPotion:
		text = "potion of running"
	case LignificationPotion:
		text = "potion of lignification"
	}
	return text
}

func (p potion) Letter() rune {
	return '!'
}

func (p potion) Int() int {
	return int(p)
}

func (p potion) Use(g *game, ev event) error {
	quant, ok := g.Player.Consumables[p]
	if !ok || quant <= 0 {
		// should not happen
		return errors.New("no such consumable: " + p.String())
	}
	if g.Player.HasStatus(StatusBerserk) {
		return errors.New("You cannot drink potions while berserk.")
	}
	var err error
	switch p {
	case HealWoundsPotion:
		err = g.QuaffHealWounds(ev)
	case TeleportationPotion:
		err = g.QuaffTeleportation(ev)
	case BerserkPotion:
		err = g.QuaffBerserk(ev)
	case DescentPotion:
		err = g.QuaffDescent(ev)
	case RunningPotion:
		err = g.QuaffHaste(ev)
	case EvasionPotion:
		err = g.QuaffEvasion(ev)
	case LignificationPotion:
		err = g.QuaffLignification(ev)
	case MagicMappingPotion:
		err = g.QuaffMagicMapping(ev)
	case MagicPotion:
		err = g.QuaffMagic(ev)
	}
	if err != nil {
		return err
	}
	ev.Renew(g, 5)
	g.UseConsumable(p)
	return nil
}

func (g *game) QuaffTeleportation(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot teleport while lignified.")
	}
	if g.Player.HasStatus(StatusTele) {
		return errors.New("You already quaffed a potion of teleportation.")
	}
	delay := 20 + RandInt(30)
	g.Player.Statuses[StatusTele]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + delay, EAction: Teleportation})
	g.Print(fmt.Sprintf("You quaff a %s. You feel unstable.", TeleportationPotion))
	return nil
}

func (g *game) QuaffBerserk(ev event) error {
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You are too exhausted to berserk.")
	}
	g.Player.Statuses[StatusBerserk]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 75, EAction: BerserkEnd})
	g.Print(fmt.Sprintf("You quaff a %s. You feel a sudden urge to kill things.", BerserkPotion))
	return nil
}

func (g *game) QuaffHealWounds(ev event) error {
	hp := g.Player.HP
	g.Player.HP += 2 * g.Player.HPMax() / 3
	if g.Player.HP > g.Player.HPMax() {
		g.Player.HP = g.Player.HPMax()
	}
	g.Print(fmt.Sprintf("You quaff a %s (%d -> %d).", HealWoundsPotion, hp, g.Player.HP))
	return nil
}

func (g *game) QuaffMagic(ev event) error {
	mp := g.Player.MP
	g.Player.MP += 2 * g.Player.MPMax() / 3
	if g.Player.MP > g.Player.MPMax() {
		g.Player.MP = g.Player.MPMax()
	}
	g.Print(fmt.Sprintf("You quaff the %s (%d -> %d).", MagicPotion, mp, g.Player.MP))
	return nil
}

func (g *game) QuaffDescent(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot descend while lignified.")
	}
	if g.Depth >= g.MaxDepth() {
		return errors.New("You cannot descend more!")
	}
	g.Print(fmt.Sprintf("You quaff the %s. You feel yourself falling through the ground.", DescentPotion))
	g.Depth++
	g.InitLevel()
	g.Save()
	return nil
}

func (g *game) QuaffHaste(ev event) error {
	g.Player.Statuses[StatusHaste]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 90, EAction: HasteEnd})
	g.Print(fmt.Sprintf("You quaff the %s. You feel speedy.", RunningPotion))
	return nil
}

func (g *game) QuaffEvasion(ev event) error {
	g.Player.Statuses[StatusEvasion]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 100, EAction: EvasionEnd})
	g.Print(fmt.Sprintf("You quaff the %s. You feel agile.", EvasionPotion))
	return nil
}

func (g *game) QuaffLignification(ev event) error {
	g.Player.Statuses[StatusLignification]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 200, EAction: LignificationEnd})
	g.Print(fmt.Sprintf("You quaff the %s. You feel attuned with the ground.", LignificationPotion))
	return nil
}

func (g *game) QuaffMagicMapping(ev event) error {
	for i, c := range g.Dungeon.Cells {
		pos := g.Dungeon.CellPosition(i)
		if c.T == FreeCell || g.Dungeon.WallNeighborsCount(pos) < 8 {
			g.Dungeon.SetExplored(pos)
		}
	}
	g.Print(fmt.Sprintf("You quaff the %s. You feel wiser.", MagicMappingPotion))
	return nil
}

type projectile int

const (
	Javeline projectile = iota
	ConfusingDart
	// unimplemented
	Net
)

func (p projectile) String() (text string) {
	switch p {
	case Javeline:
		text = "javeline"
	case ConfusingDart:
		text = "dart of confusion"
	case Net:
		text = "throwing net"
	}
	return text
}

func (p projectile) Letter() rune {
	return '('
}

func (p projectile) Int() int {
	return int(p)
}

func (p projectile) Use(g *game, ev event) error {
	quant, ok := g.Player.Consumables[p]
	if !ok || quant <= 0 {
		// should not happen
		return errors.New("no such consumable: " + p.String())
	}
	mons, _ := g.MonsterAt(g.Player.Target)
	if mons == nil {
		// should not happen
		return errors.New("internal error: no monster")
	}
	var err error
	var hit bool
	switch p {
	case Javeline:
		hit, err = g.ThrowJaveline(mons, ev)
	case ConfusingDart:
		hit, err = g.ThrowConfusingDart(mons, ev)
	}
	if !hit {
		g.Print(fmt.Sprintf("Your %s missed.", p))
	} else {
		if mons.HP > 0 {
			g.Print(fmt.Sprintf("Your %s hits the %s.", p, mons.Kind))
			mons.MakeHuntIfHurt(g)
		} else {
			g.Print(fmt.Sprintf("Your %s kills the %s.", p, mons.Kind))
			g.Killed++
		}
	}
	if err == nil {
		g.UseConsumable(p)
	}
	return err
}

func (g *game) ThrowJaveline(mons *monster, ev event) (bool, error) {
	// XXX This should go in the chooser stuff, perhaps
	ray := g.Ray(g.Player.Target)
	for _, pos := range ray {
		m, _ := g.MonsterAt(pos)
		if m != nil {
			mons = m
		}
	}
	hit := false
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if acc > evasion {
		g.MakeNoise(12, mons.Pos)
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += RandInt(5)
		}
		if g.Player.Aptitudes[AptStrong] {
			bonus += 2
		}
		base := 11
		min := base / 2
		attack := min + RandInt(base-min+1) + bonus
		attack -= RandInt(mons.Armor)
		if attack <= 0 {
			attack = 0
		}
		mons.HP -= attack
		hit = true
		mons.MakeHuntIfHurt(g)
	}
	ev.Renew(g, 10)
	return hit, nil
}

func (g *game) ThrowConfusingDart(mons *monster, ev event) (bool, error) {
	// XXX This should go in the chooser stuff, perhaps
	ray := g.Ray(g.Player.Target)
	index := 0
	for _, pos := range ray {
		m, i := g.MonsterAt(pos)
		if m != nil {
			mons = m
			index = i
		}
	}
	hit := false
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if acc > evasion {
		mons.Statuses[MonsConfused]++
		mons.Path = nil
		heap.Push(g.Events, &monsterEvent{
			ERank: ev.Rank() + 50 + RandInt(100), NMons: index, EAction: MonsConfusionEnd})
		mons.MakeHuntIfHurt(g)
		hit = true
	}
	ev.Renew(g, 10)
	return hit, nil
}

type collectable struct {
	Consumable consumable
	Quantity   int
}

type collectData struct {
	rarity   int
	quantity int
}

var ConsumablesCollectData = map[consumable]collectData{
	HealWoundsPotion:    {rarity: 6, quantity: 1},
	TeleportationPotion: {rarity: 4, quantity: 1},
	BerserkPotion:       {rarity: 5, quantity: 1},
	RunningPotion:       {rarity: 10, quantity: 1},
	DescentPotion:       {rarity: 15, quantity: 1},
	EvasionPotion:       {rarity: 8, quantity: 1},
	LignificationPotion: {rarity: 8, quantity: 1},
	MagicMappingPotion:  {rarity: 15, quantity: 1},
	MagicPotion:         {rarity: 10, quantity: 1},
	Javeline:            {rarity: 3, quantity: 3},
	ConfusingDart:       {rarity: 5, quantity: 2},
}

type equipable interface {
	Equip(g *game)
	String() string
	Letter() rune
}

type armour int

const (
	Robe armour = iota
	LeatherArmour
	ChainMail
	PlateArmour
)

func (ar armour) Equip(g *game) {
	oar := g.Player.Armour
	g.Player.Armour = ar
	g.Equipables[g.Player.Pos] = oar
}

func (ar armour) String() string {
	switch ar {
	case Robe:
		return "robe"
	case LeatherArmour:
		return "leather armour"
	case ChainMail:
		return "chain mail"
	case PlateArmour:
		return "plate armour"
	default:
		// should not happen
		return "some piece of armour"
	}
}

func (ar armour) Letter() rune {
	return '['
}

type weapon int

const (
	Dagger weapon = iota
	Axe
	BattleAxe
	Spear
	Halberd
	Sword
	DoubleSword
)

func (wp weapon) Equip(g *game) {
	owp := g.Player.Weapon
	g.Player.Weapon = wp
	g.Equipables[g.Player.Pos] = owp
}

func (wp weapon) String() string {
	switch wp {
	case Dagger:
		return "dagger"
	case Axe:
		return "axe"
	case BattleAxe:
		return "battle axe"
	case Spear:
		return "spear"
	case Halberd:
		return "halberd"
	case Sword:
		return "sword"
	case DoubleSword:
		return "double sword"
	default:
		// should not happen
		return "some weapon"
	}
}

func (wp weapon) Attack() int {
	switch wp {
	case Axe, Spear, Sword:
		return 11
	case BattleAxe, Halberd, DoubleSword:
		return 15
	case Dagger:
		return 8
	default:
		return 0
	}
}

func (wp weapon) TwoHanded() bool {
	switch wp {
	case BattleAxe, Halberd, DoubleSword:
		return true
	default:
		return false
	}
}

func (wp weapon) Letter() rune {
	return ')'
}

func (wp weapon) Cleave() bool {
	switch wp {
	case Axe, BattleAxe:
		return true
	default:
		return false
	}
}

func (wp weapon) Pierce() bool {
	switch wp {
	case Spear, Halberd:
		return true
	default:
		return false
	}
}

type shield int

const (
	NoShield shield = iota
	Buckler
	Shield
)

func (sh shield) Equip(g *game) {
	osh := g.Player.Shield
	g.Player.Shield = sh
	if osh != NoShield {
		g.Equipables[g.Player.Pos] = osh
	} else {
		delete(g.Equipables, g.Player.Pos)
	}
}

func (sh shield) String() (text string) {
	switch sh {
	case Buckler:
		text = "buckler"
	case Shield:
		text = "shield"
	}
	return text
}

func (sh shield) Letter() rune {
	return ']'
}

func (sh shield) Block() (block int) {
	switch sh {
	case Buckler:
		block += 5
	case Shield:
		block += 10
	}
	return block
}

type equipableData struct {
	rarity   int
	minDepth int
}

func (data equipableData) FavorableRoll(lateness int) int {
	ratio := data.rarity / (2 * lateness)
	if ratio < 2 {
		ratio = 2
	}
	r := RandInt(ratio)
	if r != 0 && ratio == 2 && lateness >= 3 {
		r = RandInt(ratio)
	}
	return r
}

var EquipablesRepartitionData = map[equipable]equipableData{
	Robe:          equipableData{5, 0},
	LeatherArmour: equipableData{5, 0},
	ChainMail:     equipableData{10, 3},
	PlateArmour:   equipableData{15, 6},
	Dagger:        equipableData{20, 0},
	Axe:           equipableData{25, 1},
	BattleAxe:     equipableData{30, 3},
	Spear:         equipableData{25, 1},
	Halberd:       equipableData{30, 3},
	Sword:         equipableData{25, 1},
	DoubleSword:   equipableData{30, 3},
	Buckler:       equipableData{10, 2},
	Shield:        equipableData{15, 5},
}
