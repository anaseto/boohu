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
	Plural() string
	Desc() string
	Letter() rune
	Int() int
}

func (g *game) UseConsumable(c consumable) {
	g.Player.Consumables[c]--
	g.StoryPrintf("You used %s.", Indefinite(c.String(), false))
	if g.Player.Consumables[c] <= 0 {
		delete(g.Player.Consumables, c)
	}
	g.FairAction()
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
	WallPotion
	CBlinkPotion
	// below unimplemented
	ResistancePotion
)

func (p potion) String() (text string) {
	text = "potion"
	switch p {
	case HealWoundsPotion:
		text += " of heal wounds"
	case TeleportationPotion:
		text += " of teleportation"
	case DescentPotion:
		text += " of descent"
	case EvasionPotion:
		text += " of evasion"
	case MagicMappingPotion:
		text += " of magic mapping"
	case MagicPotion:
		text += " of refill magic"
	case BerserkPotion:
		text += " of berserk"
	case RunningPotion:
		text += " of running"
	case LignificationPotion:
		text += " of lignification"
	case WallPotion:
		text += " of walls"
	case CBlinkPotion:
		text += " of controlled blink"
	case ResistancePotion:
		text += " of resistance"
	}
	return text
}

func (p potion) Plural() (text string) {
	// never used for potions
	return p.String()
}

func (p potion) Desc() (text string) {
	switch p {
	case HealWoundsPotion:
		text = "heals you a good deal."
	case TeleportationPotion:
		text = "teleports you away after a short delay."
	case DescentPotion:
		text = "makes you go to deeper in the Underground."
	case EvasionPotion:
		text = "makes you better at avoiding blows."
	case MagicMappingPotion:
		text = "shows you the map."
	case MagicPotion:
		text = "replenishes your magical reserves."
	case BerserkPotion:
		text = "makes you enter a crazy rage, temporarily making you faster, stronger and healthier. You cannot drink potions while berserk, and afterwards it leaves you slow and exhausted."
	case RunningPotion:
		text = "makes you move faster."
	case LignificationPotion:
		text = "makes you more resistant to physical blows, but you are attached to the ground while the effect lasts."
	case WallPotion:
		text = "replaces free cells around you with temporal walls."
	case CBlinkPotion:
		text = "makes you blink to a targetted cell in your line of sight."
	case ResistancePotion:
		text = "makes you resistent to the elements."
	}
	return fmt.Sprintf("The %s %s", p, text)
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
	if g.Player.HasStatus(StatusNausea) {
		return errors.New("You cannot drink potions while sick.")
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
	case WallPotion:
		err = g.QuaffWallPotion(ev)
	case CBlinkPotion:
		err = g.QuaffCBlinkPotion(ev)
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
	g.Printf("You quaff a %s. You feel unstable.", TeleportationPotion)
	return nil
}

func (g *game) QuaffBerserk(ev event) error {
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You are too exhausted to berserk.")
	}
	g.Player.Statuses[StatusBerserk]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 65 + RandInt(20), EAction: BerserkEnd})
	g.Printf("You quaff a %s. You feel a sudden urge to kill things.", BerserkPotion)
	g.Player.HP += 10
	return nil
}

func (g *game) QuaffHealWounds(ev event) error {
	hp := g.Player.HP
	g.Player.HP += 2 * g.Player.HPMax() / 3
	if g.Player.HP > g.Player.HPMax() {
		g.Player.HP = g.Player.HPMax()
	}
	g.Printf("You quaff a %s (%d -> %d).", HealWoundsPotion, hp, g.Player.HP)
	return nil
}

func (g *game) QuaffMagic(ev event) error {
	mp := g.Player.MP
	g.Player.MP += 2 * g.Player.MPMax() / 3
	if g.Player.MP > g.Player.MPMax() {
		g.Player.MP = g.Player.MPMax()
	}
	g.Printf("You quaff the %s (%d -> %d).", MagicPotion, mp, g.Player.MP)
	return nil
}

func (g *game) QuaffDescent(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot descend while lignified.")
	}
	if g.Depth >= g.MaxDepth() {
		return errors.New("You cannot descend more!")
	}
	g.Printf("You quaff the %s. You feel yourself falling through the ground.", DescentPotion)
	g.Depth++
	g.InitLevel()
	g.Save()
	return nil
}

func (g *game) QuaffHaste(ev event) error {
	g.Player.Statuses[StatusSwift]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 80 + RandInt(20), EAction: HasteEnd})
	g.Printf("You quaff the %s. You feel speedy.", RunningPotion)
	return nil
}

func (g *game) QuaffEvasion(ev event) error {
	g.Player.Statuses[StatusAgile]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 90 + RandInt(20), EAction: EvasionEnd})
	g.Printf("You quaff the %s. You feel agile.", EvasionPotion)
	return nil
}

func (g *game) QuaffLignification(ev event) error {
	g.Player.Statuses[StatusLignification]++
	heap.Push(g.Events, &simpleEvent{ERank: ev.Rank() + 150 + RandInt(100), EAction: LignificationEnd})
	g.Printf("You quaff the %s. You feel attuned with the ground.", LignificationPotion)
	return nil
}

func (g *game) QuaffMagicMapping(ev event) error {
	for i, c := range g.Dungeon.Cells {
		pos := g.Dungeon.CellPosition(i)
		if c.T == FreeCell || g.Dungeon.WallNeighborsCount(pos) < 8 {
			g.Dungeon.SetExplored(pos)
		}
	}
	g.Printf("You quaff the %s. You feel wiser.", MagicMappingPotion)
	return nil
}

func (g *game) QuaffWallPotion(ev event) error {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		posNeighbors := g.Dungeon.FreeNeighbors(pos)
		for _, pos := range posNeighbors {
			if pos == g.Player.Pos {
				continue
			}
			g.MakeNoise(18, pos)
			break
		}
		g.Dungeon.SetCell(pos, WallCell)
		if g.TemporalWalls != nil {
			g.TemporalWalls[pos] = true
		}
		heap.Push(g.Events, &cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), Pos: pos, EAction: ObstructionEnd})
	}
	g.Printf("You quaff the %s. You feel surrounded by temporal walls.", WallPotion)
	g.ComputeLOS()
	return nil
}

func (g *game) QuaffCBlinkPotion(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	if !g.ui.ChooseTarget(g, &chooser{free: true}) {
		return errors.New("Ok, then.")
	}
	g.Player.Pos = g.Player.Target
	g.Printf("You quaff the %s. You blink.", CBlinkPotion)
	g.CollectGround()
	g.ComputeLOS()
	g.MakeMonstersAware()
	return nil
}

type projectile int

const (
	Javelin projectile = iota
	ConfusingDart
	// unimplemented
	Net
)

func (p projectile) String() (text string) {
	switch p {
	case Javelin:
		text = "javelin"
	case ConfusingDart:
		text = "dart of confusion"
	case Net:
		text = "throwing net"
	}
	return text
}

func (p projectile) Plural() (text string) {
	switch p {
	case Javelin:
		text = "javelins"
	case ConfusingDart:
		text = "darts of confusion"
	case Net:
		text = "throwing nets"
	}
	return text
}

func (p projectile) Desc() (text string) {
	switch p {
	case Javelin:
		// XXX
		text = "can be thrown to ennemies, dealing up to 11 damage."
	case ConfusingDart:
		text = "can be thrown to confuse foes. Confused monsters cannot move diagonally."
	case Net:
		text = "can be thrown to emprison your ennemies."
	}
	return fmt.Sprintf("The %s %s", p, text)
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
	switch p {
	case Javelin:
		g.ThrowJavelin(mons, ev)
	case ConfusingDart:
		g.ThrowConfusingDart(mons, ev)
	}
	g.UseConsumable(p)
	return nil
}

func (g *game) ThrowJavelin(mons *monster, ev event) {
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		g.MakeNoise(12, mons.Pos)
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += RandInt(5)
		}
		if g.Player.Aptitudes[AptStrong] {
			bonus += 2
		}
		attack := g.HitDamage(11+bonus, mons.Armor)
		mons.HP -= attack
		if mons.HP > 0 {
			g.Printf("Your %s hits the %s (%d).", Javelin, mons.Kind, attack)
			mons.MakeHuntIfHurt(g)
		} else {
			g.Printf("Your %s kills the %s.", Javelin, mons.Kind)
			g.HandleKill(mons)
		}
	} else {
		g.Printf("Your %s missed the %s.", Javelin, mons.Kind)
		mons.MakeHuntIfHurt(g)
	}
	ev.Renew(g, 10)
}

func (g *game) ThrowConfusingDart(mons *monster, ev event) {
	acc := RandInt(g.Player.Accuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		mons.EnterConfusion(g, ev)
		g.Printf("Your %s hits the %s. The %s appears confused.", ConfusingDart, mons.Kind, mons.Kind)
	} else {
		g.Printf("Your %s missed the %s.", ConfusingDart, mons.Kind)
	}
	mons.MakeHuntIfHurt(g)
	ev.Renew(g, 10)
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
	EvasionPotion:       {rarity: 10, quantity: 1},
	LignificationPotion: {rarity: 8, quantity: 1},
	MagicMappingPotion:  {rarity: 15, quantity: 1},
	MagicPotion:         {rarity: 10, quantity: 1},
	WallPotion:          {rarity: 12, quantity: 1},
	CBlinkPotion:        {rarity: 12, quantity: 1},
	Javelin:             {rarity: 3, quantity: 3},
	ConfusingDart:       {rarity: 5, quantity: 2},
}

type equipable interface {
	Equip(g *game)
	String() string
	Letter() rune
	Desc() string
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
	if !g.FoundEquipables[ar] {
		if g.FoundEquipables == nil {
			g.FoundEquipables = map[equipable]bool{}
		}
		g.StoryPrintf("You found and put on %s.", Indefinite(ar.String(), false))
		g.FoundEquipables[ar] = true
	}
	g.Printf("You put the %s on and leave your %s on the ground.", ar, oar)
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

func (ar armour) Desc() string {
	var text string
	switch ar {
	case Robe:
		text = "A robe provides no special protection, and will not help you much in your journey."
	case LeatherArmour:
		text = "A leather armour provides some protection against blows."
	case ChainMail:
		text = "A chain mail provides more protection than a leather armour, but the blows you receive are louder."
	case PlateArmour:
		text = "A plate armour provides great protection against blows, but blows you receive are quite noisy."
	}
	return text
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
	if !g.FoundEquipables[wp] {
		if g.FoundEquipables == nil {
			g.FoundEquipables = map[equipable]bool{}
		}
		g.StoryPrintf("You found and took %s.", Indefinite(wp.String(), false))
		g.FoundEquipables[wp] = true
	}
	g.Printf("You take the %s and leave your %s on the ground.", wp, owp)
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

func (wp weapon) Desc() string {
	var text string
	switch wp {
	case Dagger:
		text = "A dagger is the most basic weapon. Great against sleeping monsters, but that's all."
	case Axe:
		text = "An axe is a one-handed weapon that can hit at once any foes adjacent to you."
	case BattleAxe:
		text = "A battle axe is a big two-handed weapon that can hit at once any foes adjacent to you."
	case Spear:
		text = "A spear is a one-handed weapon that can hit two opponents in a row at once. Useful in corridors."
	case Halberd:
		text = "An halberd is a big two-handed weapon that can hit two opponents in a row at once. Useful in corridors."
	case Sword:
		text = "A sword is a one-handed weapon that occasionally gets additional free hits."
	case DoubleSword:
		text = "A double sword is a big two-handed weapon that occasionally gets additional free hits."
	}
	return fmt.Sprintf("%s It can hit for up to %d damage.", text, wp.Attack())
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
	if !g.FoundEquipables[sh] {
		if g.FoundEquipables == nil {
			g.FoundEquipables = map[equipable]bool{}
		}
		g.StoryPrintf("You found and put on %s.", Indefinite(sh.String(), false))
		g.FoundEquipables[sh] = true
	}
	if osh != NoShield {
		g.Equipables[g.Player.Pos] = osh
		g.Printf("You put the %s on and leave your %s on the ground.", sh, osh)
	} else {
		delete(g.Equipables, g.Player.Pos)
		g.Printf("You put the %s on.", sh)
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

func (sh shield) Desc() (text string) {
	switch sh {
	case Buckler:
		text = "A buckler is a small shield that can sometimes block attacks, including some magical attacks. You cannot use it if you are wielding a two-handed weapon."
	case Shield:
		text = "A shield can block attacks, including some magical attacks. You cannot use it if you are wielding a two-handed weapon."
	}
	return text
}

func (sh shield) Letter() rune {
	return ']'
}

func (sh shield) Block() (block int) {
	switch sh {
	case Buckler:
		block += 6
	case Shield:
		block += 9
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
