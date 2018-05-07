package main

import (
	"errors"
	"fmt"
	"sort"
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
	SwiftnessPotion
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
	case MagicMappingPotion:
		text += " of magic mapping"
	case MagicPotion:
		text += " of refill magic"
	case BerserkPotion:
		text += " of berserk"
	case SwiftnessPotion:
		text += " of swiftness"
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
		text = "makes you go deeper in the Underground."
	case MagicMappingPotion:
		text = "shows you the map."
	case MagicPotion:
		text = "replenishes your magical reserves."
	case BerserkPotion:
		text = "makes you enter a crazy rage, temporarily making you faster, stronger and healthier. You cannot drink potions while berserk, and afterwards it leaves you slow and exhausted."
	case SwiftnessPotion:
		text = "makes you move faster and better at avoiding blows for a short time."
	case LignificationPotion:
		text = "makes you more resistant to physical blows, but you are attached to the ground while the effect lasts."
	case WallPotion:
		text = "replaces free cells around you with temporal walls."
	case CBlinkPotion:
		text = "makes you blink to a targeted cell in your line of sight."
	case ResistancePotion:
		text = "makes you resistant to the elements."
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
	case SwiftnessPotion:
		err = g.QuaffSwiftness(ev)
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
	g.ui.DrinkingPotionAnimation(g)
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
	g.Player.Statuses[StatusTele] = 1
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + delay, EAction: Teleportation})
	g.Printf("You quaff a %s. You feel unstable.", TeleportationPotion)
	return nil
}

func (g *game) QuaffBerserk(ev event) error {
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You are too exhausted to berserk.")
	}
	g.Player.Statuses[StatusBerserk]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 65 + RandInt(20), EAction: BerserkEnd})
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

func (g *game) QuaffSwiftness(ev event) error {
	g.Player.Statuses[StatusSwift]++
	g.Player.Statuses[StatusAgile]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 85 + RandInt(20), EAction: HasteEnd})
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 85 + RandInt(20), EAction: EvasionEnd})
	g.Printf("You quaff the %s. You feel speedy and agile.", SwiftnessPotion)
	return nil
}

func (g *game) QuaffLignification(ev event) error {
	g.Player.Statuses[StatusLignification]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 150 + RandInt(100), EAction: LignificationEnd})
	g.Printf("You quaff the %s. You feel rooted to the ground.", LignificationPotion)
	g.Player.HP += 10
	return nil
}

func (g *game) QuaffMagicMapping(ev event) error {
	dp := &dungeonPath{dungeon: g.Dungeon}
	g.AutoExploreDijkstra(dp, []int{g.Player.Pos.idx()})
	cdists := make(map[int][]int)
	for i, dist := range DijkstraMapCache {
		cdists[dist] = append(cdists[dist], i)
	}
	var dists []int
	for dist, _ := range cdists {
		dists = append(dists, dist)
	}
	sort.Ints(dists)
	g.ui.DrawDungeonView(g, NormalMode)
	for _, d := range dists {
		draw := false
		for _, i := range cdists[d] {
			pos := idxtopos(i)
			c := g.Dungeon.Cell(pos)
			if (c.T == FreeCell || g.Dungeon.HasFreeNeighbor(pos)) && !c.Explored {
				g.Dungeon.SetExplored(pos)
				draw = true
			}
		}
		if draw {
			g.ui.MagicMappingAnimation(g, cdists[d])
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
		g.MakeNoise(15, pos)
		g.Dungeon.SetCell(pos, WallCell)
		delete(g.Clouds, g.Player.Target)
		if g.TemporalWalls != nil {
			g.TemporalWalls[pos] = true
		}
		g.PushEvent(&cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), Pos: pos, EAction: ObstructionEnd})
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
		return errors.New(DoNothing)
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
	ExplosiveMagara
	// unimplemented
	Net
)

func (p projectile) String() (text string) {
	switch p {
	case Javelin:
		text = "javelin"
	case ConfusingDart:
		text = "dart of confusion"
	case ExplosiveMagara:
		text = "explosive magara"
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
	case ExplosiveMagara:
		text = "explosive magaras"
	case Net:
		text = "throwing nets"
	}
	return text
}

func (p projectile) Desc() (text string) {
	switch p {
	case Javelin:
		// XXX
		text = "can be thrown to foes, dealing up to 11 damage."
	case ConfusingDart:
		text = "can be thrown to confuse foes. Confused monsters cannot move diagonally."
	case ExplosiveMagara:
		text = "can be thrown to cause a fire explosion halving HP of monsters in a square area."
	case Net:
		text = "can be thrown to emprison your enemies."
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
	var err error
	switch p {
	case Javelin:
		err = g.ThrowJavelin(ev)
	case ConfusingDart:
		err = g.ThrowConfusingDart(ev)
	case ExplosiveMagara:
		err = g.ThrowExplosiveMagara(ev)
	}
	g.UseConsumable(p)
	return err
}

func (g *game) ThrowJavelin(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{needsFreeWay: true}) {
		return errors.New(DoNothing)
	}
	mons, _ := g.MonsterAt(g.Player.Target)
	acc := RandInt(g.Player.RangedAccuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		noise := 12 + mons.Armor/3
		g.MakeNoise(noise, mons.Pos)
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += RandInt(5)
		}
		if g.Player.Aptitudes[AptStrong] {
			bonus += 2
		}
		attack := g.HitDamage(DmgPhysical, 11+bonus, mons.Armor)
		mons.HP -= attack
		if mons.HP > 0 {
			g.PrintfStyled("Your %s hits the %s (%d).", logPlayerHit, Javelin, mons.Kind, attack)
			g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
			mons.MakeHuntIfHurt(g)
		} else {
			g.PrintfStyled("Your %s kills the %s.", logPlayerHit, Javelin, mons.Kind)
			g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
			g.HandleKill(mons, ev)
		}
	} else {
		g.Printf("Your %s missed the %s.", Javelin, mons.Kind)
		g.ui.ThrowAnimation(g, g.Ray(mons.Pos), false)
		mons.MakeHuntIfHurt(g)
	}
	ev.Renew(g, 10)
	return nil
}

func (g *game) ThrowConfusingDart(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{needsFreeWay: true}) {
		return errors.New(DoNothing)
	}
	mons, _ := g.MonsterAt(g.Player.Target)
	acc := RandInt(g.Player.RangedAccuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		mons.EnterConfusion(g, ev)
		g.PrintfStyled("Your %s hits the %s. The %s appears confused.", logPlayerHit, ConfusingDart, mons.Kind, mons.Kind)
		g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
	} else {
		g.Printf("Your %s missed the %s.", ConfusingDart, mons.Kind)
		g.ui.ThrowAnimation(g, g.Ray(mons.Pos), false)
	}
	mons.MakeHuntIfHurt(g)
	ev.Renew(g, 10)
	return nil
}

func (g *game) ThrowExplosiveMagara(ev event) error {
	if !g.ui.ChooseTarget(g, &chooser{area: true, minDist: true, flammable: true}) {
		return errors.New(DoNothing)
	}
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Target)
	g.Print("You throw the explosive magara, which gives a noisy pop.")
	g.MakeNoise(18, g.Player.Target)
	g.ui.ExplosionAnimation(g, FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.Burn(pos, ev)
		mons, _ := g.MonsterAt(pos)
		if mons.Exists() {
			mons.HP /= 2
			if mons.HP == 0 {
				mons.HP = 1
			}
			g.MakeNoise(12, mons.Pos)
			mons.MakeHuntIfHurt(g)
		}
	}

	ev.Renew(g, 10)
	return nil
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
	SwiftnessPotion:     {rarity: 6, quantity: 1},
	DescentPotion:       {rarity: 15, quantity: 1},
	LignificationPotion: {rarity: 8, quantity: 1},
	MagicMappingPotion:  {rarity: 15, quantity: 1},
	MagicPotion:         {rarity: 10, quantity: 1},
	WallPotion:          {rarity: 12, quantity: 1},
	CBlinkPotion:        {rarity: 12, quantity: 1},
	Javelin:             {rarity: 3, quantity: 3},
	ConfusingDart:       {rarity: 5, quantity: 2},
	ExplosiveMagara:     {rarity: 10, quantity: 1},
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
	Frundis
	ElecWhip
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
	if wp == Frundis {
		g.PrintfStyled("♫ ♪ … Oh, you're there, let's fight our way out!", logSpecial)
	}
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
	case Frundis:
		return "staff Frundis"
	case ElecWhip:
		return "lightning whip"
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
	case Frundis:
		text = "Frundis is a musician and harmonist, which happens to be a two-handed staff too. It may occasionally confuse monsters on hit. It magically helps reducing noise in combat too."
	case ElecWhip:
		text = "The lightning whip is a one-handed weapon that inflicts electrical damage to a monster and any foes connected to it."
	}
	return fmt.Sprintf("%s It can hit for up to %d damage.", text, wp.Attack())
}

func (wp weapon) Attack() int {
	switch wp {
	case Axe, Spear, Sword:
		return 11
	case BattleAxe, Halberd, DoubleSword:
		return 15
	case Frundis:
		return 13
	case Dagger:
		return 9
	case ElecWhip:
		return 8
	default:
		return 0
	}
}

func (wp weapon) TwoHanded() bool {
	switch wp {
	case BattleAxe, Halberd, DoubleSword, Frundis:
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
	Robe:          {5, 0},
	LeatherArmour: {5, 0},
	ChainMail:     {10, 3},
	PlateArmour:   {15, 6},
	Buckler:       {10, 2},
	Shield:        {15, 5},
}
