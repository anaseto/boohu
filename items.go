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
	g.FunAction()
}

type potion int

const (
	HealWoundsPotion potion = iota
	TeleportationPotion
	//BerserkPotion
	DescentPotion
	SwiftnessPotion
	//LignificationPotion
	MagicMappingPotion
	MagicPotion
	WallPotion
	CBlinkPotion
	DigPotion
	SwapPotion
	ShadowsPotion
	ConfusePotion
	TormentPotion
	DreamPotion
)

const NumPotions = int(DreamPotion) + 1

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
	//case BerserkPotion:
	//text += " of berserk"
	case SwiftnessPotion:
		text += " of swiftness"
	//case LignificationPotion:
	//text += " of lignification"
	case WallPotion:
		text += " of walls"
	case CBlinkPotion:
		text += " of controlled blink"
	case DigPotion:
		text += " of digging"
	case SwapPotion:
		text += " of swapping"
	case ShadowsPotion:
		text += " of shadows"
	case ConfusePotion:
		text += " of confusion"
	case TormentPotion:
		text += " of torment explosion"
	case DreamPotion:
		text += " of dreams"
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
	//case BerserkPotion:
	//text = "makes you enter a crazy rage, temporarily making you faster, stronger and healthier. You cannot drink potions while berserk, and afterwards it leaves you slow and exhausted."
	case SwiftnessPotion:
		text = "makes you move faster and better at avoiding blows for a short time."
	//case LignificationPotion:
	//text = "makes you more resistant to physical blows, but you are attached to the ground while the effect lasts."
	case WallPotion:
		text = "replaces free cells around you with temporary walls."
	case CBlinkPotion:
		text = "makes you blink to a targeted cell in your line of sight."
	case DigPotion:
		text = "makes you dig walls like an earth dragon."
	case SwapPotion:
		text = "makes you swap positions with monsters instead of attacking."
	case ShadowsPotion:
		text = "reduces your line of sight range to 1."
	case ConfusePotion:
		text = "generates a harmonic light that confuses monsters in your line of sight."
	case TormentPotion:
		text = "halves HP of every creature in sight, including the player, and destroys visible walls. Extremely noisy. It can burn foliage and doors."
	case DreamPotion:
		text = "shows you the position of monsters sleeping at drink time."
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
	//case BerserkPotion:
	//err = g.QuaffBerserk(ev)
	case DescentPotion:
		err = g.QuaffDescent(ev)
	case SwiftnessPotion:
		err = g.QuaffSwiftness(ev)
	//case LignificationPotion:
	//err = g.QuaffLignification(ev)
	case MagicMappingPotion:
		err = g.QuaffMagicMapping(ev)
	case MagicPotion:
		err = g.QuaffMagic(ev)
	case WallPotion:
		err = g.QuaffWallPotion(ev)
	case CBlinkPotion:
		err = g.QuaffCBlinkPotion(ev)
	case DigPotion:
		err = g.QuaffDigPotion(ev)
	case SwapPotion:
		err = g.QuaffSwapPotion(ev)
	case ShadowsPotion:
		err = g.QuaffShadowsPotion(ev)
	case ConfusePotion:
		err = g.QuaffConfusePotion(ev)
	case TormentPotion:
		err = g.QuaffTormentPotion(ev)
	case DreamPotion:
		err = g.QuaffDreamPotion(ev)
	}
	if err != nil {
		return err
	}
	//if p == DescentPotion {
	//g.Stats.UsedPotion[g.Depth-1]++
	//} else {
	//g.Stats.UsedPotion[g.Depth]++
	//}
	ev.Renew(g, 5)
	g.UseConsumable(p)
	g.Stats.Drinks++
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
	g.Printf("You quaff the %s. You feel unstable.", TeleportationPotion)
	return nil
}

func (g *game) QuaffHealWounds(ev event) error {
	hp := g.Player.HP
	g.Player.HP += 2 * g.Player.HPMax() / 3
	if g.Player.HP > g.Player.HPMax() {
		g.Player.HP = g.Player.HPMax()
	}
	g.Printf("You quaff the %s (%d -> %d).", HealWoundsPotion, hp, g.Player.HP)
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
	// why not?
	//if g.Player.HasStatus(StatusLignification) {
	//return errors.New("You cannot descend while lignified.")
	//}
	if g.Depth >= MaxDepth {
		return errors.New("You cannot descend any deeper!")
	}
	g.Printf("You quaff the %s. You fall through the ground.", DescentPotion)
	g.LevelStats()
	g.StoryPrint("You descended deeper into the dungeon.")
	g.Depth++
	g.DepthPlayerTurn = 0
	g.InitLevel()
	g.Save()
	return nil
}

func (g *game) QuaffSwiftness(ev event) error {
	if g.Player.HasStatus(StatusSwift) && g.Player.HasStatus(StatusAgile) {
		return fmt.Errorf("You already quaffed a %s potion.", SwiftnessPotion)
	}
	g.Player.Statuses[StatusSwift]++
	end := ev.Rank() + 85 + RandInt(20)
	g.PushEvent(&simpleEvent{ERank: end, EAction: HasteEnd})
	g.Player.Expire[StatusSwift] = end
	g.Player.Statuses[StatusAgile]++
	g.PushEvent(&simpleEvent{ERank: end, EAction: EvasionEnd})
	g.Player.Expire[StatusAgile] = end
	g.Printf("You quaff the %s. You feel speedy and agile.", SwiftnessPotion)
	return nil
}

func (g *game) QuaffDigPotion(ev event) error {
	if g.Player.HasStatus(StatusDig) {
		return errors.New("You are already digging.")
	}
	g.Player.Statuses[StatusDig] = 1
	end := ev.Rank() + 75 + RandInt(20)
	g.PushEvent(&simpleEvent{ERank: end, EAction: DigEnd})
	g.Player.Expire[StatusDig] = end
	g.Printf("You quaff the %s. You feel like an earth dragon.", DigPotion)
	return nil
}

func (g *game) QuaffSwapPotion(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot drink this potion while lignified.")
	}
	if g.Player.HasStatus(StatusSwap) {
		return errors.New("You are already swapping.")
	}
	g.Player.Statuses[StatusSwap] = 1
	end := ev.Rank() + 130 + RandInt(41)
	g.PushEvent(&simpleEvent{ERank: end, EAction: SwapEnd})
	g.Player.Expire[StatusSwap] = end
	g.Printf("You quaff the %s. You feel light-footed.", SwapPotion)
	return nil
}

func (g *game) QuaffShadowsPotion(ev event) error {
	if g.Player.HasStatus(StatusShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Player.Statuses[StatusShadows] = 1
	end := ev.Rank() + 130 + RandInt(41)
	g.PushEvent(&simpleEvent{ERank: end, EAction: ShadowsEnd})
	g.Player.Expire[StatusShadows] = end
	g.Printf("You quaff the %s. You feel surrounded by shadows.", ShadowsPotion)
	g.ComputeLOS()
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
	g.Printf("You quaff the %s. You feel aware of your surroundings..", MagicMappingPotion)
	return nil
}

func (g *game) QuaffConfusePotion(ev event) error {
	g.Printf("You quaff the %s. A harmonic light confuses monsters.", ConfusePotion)
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.EnterConfusion(g, ev)
		}
	}
	return nil
}

func (g *game) QuaffTormentPotion(ev event) error {
	g.Printf("You quaff the %s. %s It hurts!", TormentPotion, g.ExplosionSound())
	damage := g.Player.HP / 2
	g.Player.HP = g.Player.HP - damage
	g.Stats.Damage += damage
	g.ui.WoundedAnimation(g)
	g.MakeNoise(ExplosionNoise+10, g.Player.Pos)
	g.ui.TormentExplosionAnimation(g)
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		g.ExplosionAt(ev, pos)
	}
	return nil
}

func (g *game) QuaffDreamPotion(ev event) error {
	for _, mons := range g.Monsters {
		if mons.Exists() && mons.State == Resting && !g.Player.LOS[mons.Pos] {
			g.DreamingMonster[mons.Pos] = true
		}
	}
	g.Printf("You quaff the %s. You perceive monsters' dreams.", DreamPotion)
	return nil
}

func (g *game) QuaffWallPotion(ev event) error {
	neighbors := g.Dungeon.FreeNeighbors(g.Player.Pos)
	for _, pos := range neighbors {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			continue
		}
		g.CreateTemporalWallAt(pos, ev)
	}
	g.Printf("You quaff the %s. You feel surrounded by temporary walls.", WallPotion)
	g.ComputeLOS()
	return nil
}

func (g *game) QuaffCBlinkPotion(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You cannot blink while lignified.")
	}
	if err := g.ui.ChooseTarget(g, &chooser{free: true}); err != nil {
		return err
	}
	g.Printf("You quaff the %s. You blink.", CBlinkPotion)
	g.PlacePlayerAt(g.Player.Target)
	return nil
}

type projectile int

const (
	ConfusingDart projectile = iota
	ExplosiveMagara
	TeleportMagara
	NightMagara
)

const NumProjectiles = int(NightMagara) + 1

func (p projectile) String() (text string) {
	switch p {
	case ConfusingDart:
		text = "dart of confusion"
	case ExplosiveMagara:
		text = "explosive magara"
	case TeleportMagara:
		text = "teleport magara"
	case NightMagara:
		text = "night magara"
	}
	return text
}

func (p projectile) Plural() (text string) {
	switch p {
	case ConfusingDart:
		text = "darts of confusion"
	case ExplosiveMagara:
		text = "explosive magaras"
	case TeleportMagara:
		text = "teleport magaras"
	case NightMagara:
		text = "night magaras"
	}
	return text
}

func (p projectile) Desc() (text string) {
	switch p {
	case ConfusingDart:
		text = "can be silently thrown to confuse foes. Confused monsters cannot move diagonally."
	case ExplosiveMagara:
		text = "can be thrown to cause a fire explosion halving HP of monsters in a square area. It can occasionally destruct walls. It can burn doors and foliage."
	case TeleportMagara:
		text = "can be thrown to make monsters in a square area teleport."
	case NightMagara:
		text = "can be thrown at a monster to produce sleep inducing clouds in a 2-radius area. You are affected too by the clouds, but they will slow your actions instead. Can burn doors and foliage."
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
	case ConfusingDart:
		err = g.ThrowConfusingDart(ev)
	case ExplosiveMagara:
		err = g.ThrowExplosiveMagara(ev)
	case TeleportMagara:
		err = g.ThrowTeleportMagara(ev)
	case NightMagara:
		err = g.ThrowNightMagara(ev)
	}
	if err != nil {
		return err
	}
	g.UseConsumable(p)
	g.Stats.Throws++
	return nil
}

func (g *game) ThrowConfusingDart(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{needsFreeWay: true}); err != nil {
		return err
	}
	mons := g.MonsterAt(g.Player.Target)
	mons.EnterConfusion(g, ev)
	g.PrintfStyled("Your %s hits the %s, who appears confused.", logPlayerHit, ConfusingDart, mons.Kind)
	g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
	mons.MakeHuntIfHurt(g)
	g.HandleStone(mons)
	ev.Renew(g, 10)
	return nil
}

func (g *game) ExplosionAt(ev event, pos position) {
	g.Burn(pos, ev)
	mons := g.MonsterAt(pos)
	if mons.Exists() {
		mons.HP /= 2
		if mons.HP == 0 {
			mons.HP = 1
		}
		g.MakeNoise(ExplosionHitNoise, mons.Pos)
		g.HandleStone(mons)
		mons.MakeHuntIfHurt(g)
	} else if g.Dungeon.Cell(pos).T == WallCell {
		g.Dungeon.SetCell(pos, FreeCell)
		g.Stats.Digs++
		if !g.Player.LOS[pos] {
			g.WrongWall[pos] = true
		} else {
			g.ui.WallExplosionAnimation(g, pos)
		}
		g.MakeNoise(WallNoise, pos)
		g.Fog(pos, 1, ev)
	}
}

func (g *game) ThrowExplosiveMagara(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{area: true, minDist: true, flammable: true, wall: true}); err != nil {
		return err
	}
	neighbors := g.Player.Target.ValidNeighbors()
	g.Printf("You throw the explosive magara... %s", g.ExplosionSound())
	g.MakeNoise(ExplosionNoise, g.Player.Target)
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgPlayer)
	g.ui.ExplosionAnimation(g, FireExplosion, g.Player.Target)
	for _, pos := range append(neighbors, g.Player.Target) {
		g.ExplosionAt(ev, pos)
	}

	ev.Renew(g, 10)
	return nil
}

func (g *game) ThrowTeleportMagara(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{area: true, minDist: true}); err != nil {
		return err
	}
	neighbors := g.Player.Target.ValidNeighbors()
	g.Print("You throw the teleport magara.")
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgPlayer)
	for _, pos := range append(neighbors, g.Player.Target) {
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.TeleportAway(g)
		}
	}

	ev.Renew(g, 10)
	return nil
}

func (g *game) NightFog(at position, radius int, ev event) {
	dij := &normalPath{game: g}
	nm := Dijkstra(dij, []position{at}, radius)
	for pos := range nm {
		_, ok := g.Clouds[pos]
		if !ok {
			g.Clouds[pos] = CloudNight
			g.PushEvent(&cloudEvent{ERank: ev.Rank() + 10, EAction: NightProgression, Pos: pos})
		}
	}
	g.ComputeLOS()
}

func (g *game) ThrowNightMagara(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{needsFreeWay: true}); err != nil {
		return err
	}
	g.Print("You throw the night magara… Clouds come out of it.")
	g.ui.ProjectileTrajectoryAnimation(g, g.Ray(g.Player.Target), ColorFgSleepingMonster)
	g.NightFog(g.Player.Target, 2, ev)

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
	ConfusingDart:       {rarity: 4, quantity: 3},
	ExplosiveMagara:     {rarity: 8, quantity: 1},
	NightMagara:         {rarity: 10, quantity: 1},
	TeleportMagara:      {rarity: 8, quantity: 1},
	TeleportationPotion: {rarity: 6, quantity: 1},
	//BerserkPotion:       {rarity: 6, quantity: 1},
	HealWoundsPotion: {rarity: 6, quantity: 1},
	SwiftnessPotion:  {rarity: 6, quantity: 1},
	//LignificationPotion: {rarity: 9, quantity: 1},
	MagicPotion:        {rarity: 9, quantity: 1},
	WallPotion:         {rarity: 12, quantity: 1},
	CBlinkPotion:       {rarity: 12, quantity: 1},
	DigPotion:          {rarity: 12, quantity: 1},
	SwapPotion:         {rarity: 12, quantity: 1},
	ShadowsPotion:      {rarity: 15, quantity: 1},
	ConfusePotion:      {rarity: 15, quantity: 1},
	DescentPotion:      {rarity: 18, quantity: 1},
	MagicMappingPotion: {rarity: 18, quantity: 1},
	DreamPotion:        {rarity: 18, quantity: 1},
	TormentPotion:      {rarity: 30, quantity: 1},
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
	//ChainMail
	SmokingScales
	ShinyPlates
	TurtlePlates
	SpeedRobe
	CelmistRobe
	HarmonistRobe
)

func (ar armour) Equip(g *game) {
	oar := g.Player.Armour
	g.Player.Armour = ar
	if !g.FoundEquipables[ar] {
		g.StoryPrintf("You found and put on %s.", ar.StringIndefinite())
		g.FoundEquipables[ar] = true
	}
	g.Printf("You put the %s on and leave your %s.", ar, oar)
	g.Equipables[g.Player.Pos] = oar
	if oar == CelmistRobe && g.Player.MP > g.Player.MPMax() {
		g.Player.MP = g.Player.MPMax()
	}
}

func (ar armour) String() string {
	switch ar {
	case Robe:
		return "robe"
	case LeatherArmour:
		return "leather armour"
	//case ChainMail:
	//return "chain mail"
	case SmokingScales:
		return "smoking scales"
	case ShinyPlates:
		return "shiny plates"
	case TurtlePlates:
		return "turtle plates"
	case SpeedRobe:
		return "robe of speed"
	case CelmistRobe:
		return "celmist robe"
	case HarmonistRobe:
		return "harmonist robe"
	default:
		// should not happen
		return "?"
	}
}

func (ar armour) StringIndefinite() string {
	switch ar {
	case ShinyPlates, TurtlePlates, SmokingScales:
		return ar.String()
	default:
		return "a " + ar.String()
	}
}

func (ar armour) Short() string {
	switch ar {
	case Robe:
		return "Rb"
	case LeatherArmour:
		return "Lt"
	//case ChainMail:
	//return "Ch"
	case SmokingScales:
		return "Sm"
	case ShinyPlates:
		return "Sh"
	case TurtlePlates:
		return "Tr"
	case SpeedRobe:
		return "Sp"
	case CelmistRobe:
		return "Cl"
	case HarmonistRobe:
		return "Hr"
	default:
		// should not happen
		return "?"
	}
}

func (ar armour) Desc() string {
	var text string
	switch ar {
	case Robe:
		text = "A robe provides no special protection, and will not help you much in your journey."
	case LeatherArmour:
		text = "A leather armour provides some protection against blows."
	//case ChainMail:
	//text = "A chain mail provides good protection against blows, at a minor evasion cost."
	case SmokingScales:
		text = "Smoking scales provide protection against blows. They leave short-lived fog as you move."
	case ShinyPlates:
		text = "Shiny plates provide very good protection against blows, but increase your line of sight range."
	case TurtlePlates:
		text = "Turtle plates provide great protection against blows, but make you move slower and a little less good at evading blows."
	case SpeedRobe:
		text = "The speed robe makes you move faster, but makes you frail."
	case CelmistRobe:
		text = "The celmist robe improves your magic reserves, rod recharge rate, and rods can gain an extra charge."
	case HarmonistRobe:
		text = "The harmonist robe makes you harder to detect (reduced LOS, stealthy, noise mitigation)."
	}
	return text
}

func (ar armour) Letter() rune {
	return '['
}

type shield int

const (
	NoShield shield = iota
	Buckler
	ConfusingShield
	EarthShield
	BashingShield
	FireShield
)

func (sh shield) Equip(g *game) {
	osh := g.Player.Shield
	g.Player.Shield = sh
	if !g.FoundEquipables[sh] {
		g.StoryPrintf("You found and put on %s.", Indefinite(sh.String(), false))
		g.FoundEquipables[sh] = true
	}
	if osh != NoShield {
		g.Equipables[g.Player.Pos] = osh
		g.Printf("You put the %s on and leave your %s.", sh, osh)
	} else {
		delete(g.Equipables, g.Player.Pos)
		g.Printf("You put the %s on.", sh)
	}
}

func (sh shield) String() (text string) {
	switch sh {
	case Buckler:
		text = "buckler"
	case ConfusingShield:
		text = "confusing shield"
	case EarthShield:
		text = "earth shield"
	case BashingShield:
		text = "bashing shield"
	case FireShield:
		text = "fire shield"
	}
	return text
}

func (sh shield) Short() (text string) {
	switch sh {
	case Buckler:
		text = "Bk"
	case ConfusingShield:
		text = "Cn"
	case EarthShield:
		text = "Er"
	case BashingShield:
		text = "Bs"
	case FireShield:
		text = "Fr"
	}
	return text
}

func (sh shield) Desc() (text string) {
	switch sh {
	case Buckler:
		text = "A buckler is a small shield that can block attacks."
	case ConfusingShield:
		text = "A confusing shield blocks attacks, sometimes confusing monsters."
	case EarthShield:
		text = "An earth shield offers great protection, but impact sound can disintegrate nearby walls."
	case BashingShield:
		text = "A bashing shield can block attacks and push ennemies away."
	case FireShield:
		text = "A fire shield blocks attacks, sometimes burning nearby foliage."
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
	case ConfusingShield, BashingShield, FireShield:
		block += 9
	case EarthShield:
		block += 15
	}
	return block
}
