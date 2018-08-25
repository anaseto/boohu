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
	BerserkPotion
	DescentPotion
	SwiftnessPotion
	LignificationPotion
	MagicMappingPotion
	MagicPotion
	WallPotion
	CBlinkPotion
	DigPotion
	SwapPotion
	ShadowsPotion
	ConfusePotion
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
	case DigPotion:
		text += " of digging"
	case SwapPotion:
		text += " of swapping"
	case ShadowsPotion:
		text += " of shadows"
	case ConfusePotion:
		text += " of confusion"
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
	case BerserkPotion:
		text = "makes you enter a crazy rage, temporarily making you faster, stronger and healthier. You cannot drink potions while berserk, and afterwards it leaves you slow and exhausted."
	case SwiftnessPotion:
		text = "makes you move faster and better at avoiding blows for a short time."
	case LignificationPotion:
		text = "makes you more resistant to physical blows, but you are attached to the ground while the effect lasts."
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
	case DigPotion:
		err = g.QuaffDigPotion(ev)
	case SwapPotion:
		err = g.QuaffSwapPotion(ev)
	case ShadowsPotion:
		err = g.QuaffShadowsPotion(ev)
	case ConfusePotion:
		err = g.QuaffConfusePotion(ev)
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

func (g *game) QuaffBerserk(ev event) error {
	if g.Player.HasStatus(StatusExhausted) {
		return errors.New("You are too exhausted to berserk.")
	}
	g.Player.Statuses[StatusBerserk] = 1
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 65 + RandInt(20), EAction: BerserkEnd})
	g.Printf("You quaff the %s. You feel a sudden urge to kill things.", BerserkPotion)
	g.Player.HP += 10
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
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 85 + RandInt(20), EAction: HasteEnd})
	g.Player.Statuses[StatusAgile]++
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 85 + RandInt(20), EAction: EvasionEnd})
	g.Printf("You quaff the %s. You feel speedy and agile.", SwiftnessPotion)
	return nil
}

func (g *game) QuaffDigPotion(ev event) error {
	if g.Player.HasStatus(StatusDig) {
		return errors.New("You are already digging.")
	}
	g.Player.Statuses[StatusDig] = 1
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 75 + RandInt(20), EAction: DigEnd})
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
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 130 + RandInt(41), EAction: SwapEnd})
	g.Printf("You quaff the %s. You feel light-footed.", SwapPotion)
	return nil
}

func (g *game) QuaffShadowsPotion(ev event) error {
	if g.Player.HasStatus(StatusShadows) {
		return errors.New("You are already surrounded by shadows.")
	}
	g.Player.Statuses[StatusShadows] = 1
	g.PushEvent(&simpleEvent{ERank: ev.Rank() + 130 + RandInt(41), EAction: ShadowsEnd})
	g.Printf("You quaff the %s. You feel surrounded by shadows.", ShadowsPotion)
	g.ComputeLOS()
	return nil
}

func (g *game) QuaffLignification(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You are already lignified.")
	}
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
		g.Dungeon.SetCell(pos, WallCell)
		delete(g.Clouds, g.Player.Target)
		if g.TemporalWalls != nil {
			g.TemporalWalls[pos] = true
		}
		g.PushEvent(&cloudEvent{ERank: ev.Rank() + 200 + RandInt(50), Pos: pos, EAction: ObstructionEnd})
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
	NightMagara
)

const NumProjectiles = int(ExplosiveMagara) + 1

func (p projectile) String() (text string) {
	switch p {
	case ConfusingDart:
		text = "dart of confusion"
	case ExplosiveMagara:
		text = "explosive magara"
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
	case NightMagara:
		text = "night magaras"
	}
	return text
}

func (p projectile) Desc() (text string) {
	switch p {
	case ConfusingDart:
		text = "can be silently thrown to confuse foes, dealing up to 7 damage. Confused monsters cannot move diagonally."
	case ExplosiveMagara:
		text = "can be thrown to cause a fire explosion halving HP of monsters in a square area. It can occasionally destruct walls."
	case NightMagara:
		text = "can be thrown at a monster to produce sleep inducing clouds in a 2-radius area. You are affected too by the clouds, but they will slow your actions instead."
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
	acc := RandInt(g.Player.RangedAccuracy())
	evasion := RandInt(mons.Evasion)
	if mons.State == Resting {
		evasion /= 2 + 1
	}
	if acc > evasion {
		bonus := 0
		if g.Player.HasStatus(StatusBerserk) {
			bonus += RandInt(5)
		}
		if g.Player.Aptitudes[AptStrong] {
			bonus += 2
		}
		attack, _ := g.HitDamage(DmgPhysical, 7+bonus, mons.Armor) // no clang with darts
		mons.HP -= attack
		if mons.HP > 0 {
			mons.EnterConfusion(g, ev)
			g.PrintfStyled("Your %s hits the %s (%d dmg), who appears confused.", logPlayerHit, ConfusingDart, mons.Kind, attack)
			g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
			mons.MakeHuntIfHurt(g)
		} else {
			g.PrintfStyled("Your %s kills the %s.", logPlayerHit, ConfusingDart, mons.Kind)
			g.ui.ThrowAnimation(g, g.Ray(mons.Pos), true)
			g.HandleKill(mons, ev)
		}
	} else {
		g.Printf("Your %s missed the %s.", ConfusingDart, mons.Kind)
		g.ui.ThrowAnimation(g, g.Ray(mons.Pos), false)
	}
	ev.Renew(g, 10)
	return nil
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
		g.Burn(pos, ev)
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.HP /= 2
			if mons.HP == 0 {
				mons.HP = 1
			}
			g.MakeNoise(ExplosionHitNoise, mons.Pos)
			mons.MakeHuntIfHurt(g)
		} else if g.Dungeon.Cell(pos).T == WallCell && RandInt(2) == 0 {
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
	TeleportationPotion: {rarity: 6, quantity: 1},
	BerserkPotion:       {rarity: 6, quantity: 1},
	HealWoundsPotion:    {rarity: 6, quantity: 1},
	SwiftnessPotion:     {rarity: 6, quantity: 1},
	LignificationPotion: {rarity: 9, quantity: 1},
	MagicPotion:         {rarity: 9, quantity: 1},
	WallPotion:          {rarity: 12, quantity: 1},
	CBlinkPotion:        {rarity: 12, quantity: 1},
	DigPotion:           {rarity: 12, quantity: 1},
	SwapPotion:          {rarity: 12, quantity: 1},
	ShadowsPotion:       {rarity: 15, quantity: 1},
	ConfusePotion:       {rarity: 15, quantity: 1},
	DescentPotion:       {rarity: 18, quantity: 1},
	MagicMappingPotion:  {rarity: 18, quantity: 1},
	DreamPotion:         {rarity: 18, quantity: 1},
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
	SmokingScales
	ScintillatingPlates
	PonderousnessPlates
	SpeedRobe
	CelmistRobe
	HarmonistRobe
)

func (ar armour) Equip(g *game) {
	oar := g.Player.Armour
	g.Player.Armour = ar
	if !g.FoundEquipables[ar] {
		g.StoryPrintf("You found and put on %s.", Indefinite(ar.String(), false))
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
	case ChainMail:
		return "chain mail"
	case SmokingScales:
		return "smoking scales"
	case ScintillatingPlates:
		return "scintillating plates"
	case PonderousnessPlates:
		return "ponderousness plates"
	case SpeedRobe:
		return "robe of speed"
	case CelmistRobe:
		return "celmist robe"
	case HarmonistRobe:
		return "harmonist robe"
	default:
		// should not happen
		return "some piece of armour"
	}
}

func (ar armour) Short() string {
	switch ar {
	case Robe:
		return "Rb"
	case LeatherArmour:
		return "Lt"
	case ChainMail:
		return "Ch"
	case SmokingScales:
		return "Sm"
	case ScintillatingPlates:
		return "Sc"
	case PonderousnessPlates:
		return "Pl"
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
	case ChainMail:
		text = "A chain mail provides good protection against blows, at a minor evasion cost."
	case SmokingScales:
		text = "Smoking scales provide protection against blows. They leave short-lived fog as you move."
	case ScintillatingPlates:
		text = "Scintillating plates provide very good protection against blows, but increase your line of sight range."
	case PonderousnessPlates:
		text = "Ponderousness plates provide great protection against blows, but make you move slower and a little less good at evading blows."
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

type weapon int

const (
	Dagger weapon = iota
	Axe
	BattleAxe
	Spear
	Halberd
	Sabre
	DancingRapier
	BerserkSword
	Frundis
	ElecWhip
	HarKarGauntlets
	DefenderFlail
)

func (wp weapon) Equip(g *game) {
	owp := g.Player.Weapon
	g.Player.Weapon = wp
	if !g.FoundEquipables[wp] {
		g.StoryPrintf("You found and took %s.", Indefinite(wp.String(), false))
		g.FoundEquipables[wp] = true
	}
	g.Printf("You take the %s and leave your %s.", wp, owp)
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
	case Sabre:
		return "sabre"
	case DancingRapier:
		return "dancing rapier"
	case BerserkSword:
		return "berserk sword"
	case Frundis:
		return "staff Frundis"
	case ElecWhip:
		return "lightning whip"
	case HarKarGauntlets:
		return "har-kar gauntlets"
	case DefenderFlail:
		return "defender flail"
	default:
		// should not happen
		return "some weapon"
	}
}

func (wp weapon) Short() string {
	switch wp {
	case Dagger:
		return "Dg"
	case Axe:
		return "Ax"
	case BattleAxe:
		return "Bt"
	case Spear:
		return "Sp"
	case Halberd:
		return "Hl"
	case Sabre:
		return "Sb"
	case DancingRapier:
		return "Dn"
	case BerserkSword:
		return "Br"
	case Frundis:
		return "Fr"
	case ElecWhip:
		return "Wh"
	case HarKarGauntlets:
		return "Hk"
	case DefenderFlail:
		return "Fl"
	default:
		// should not happen
		return "?"
	}
}

func (wp weapon) Desc() string {
	var text string
	switch wp {
	case Dagger:
		text = "A dagger is the most basic weapon. Great against sleeping monsters, but that's all."
	case Axe:
		text = "An axe is a one-handed weapon that can hit at once any foes adjacent to you, dealing extra damage in the open."
	case BattleAxe:
		text = "A battle axe is a big two-handed weapon that can hit at once any foes adjacent to you, dealing extra damage in the open."
	case Spear:
		text = "A spear is a one-handed weapon that can hit two opponents in a row at once. Useful in corridors."
	case Halberd:
		text = "An halberd is a big two-handed weapon that can hit two opponents in a row at once. Useful in corridors."
	case Sabre:
		text = "A sabre is a one-handed weapon. It is more accurate against injured opponents."
	case DancingRapier:
		text = "A dancing rapier is a one-handed weapon. It makes you swap with your foe and can hit another monster behind with extra damage."
	case BerserkSword:
		text = "A berserk sword is a big two-handed weapon that can make you berserk when attacking while injured."
	case Frundis:
		text = "Frundis is a musician and harmonist, which happens to be a two-handed staff too. It may occasionally confuse monsters on hit. It magically helps reducing noise in combat too."
	case ElecWhip:
		text = "The lightning whip is a one-handed weapon that inflicts electrical damage to a monster and any foes connected to it."
	case HarKarGauntlets:
		text = "Har-kar gauntlets are an unarmed combat weapon. They allow you to make a wind attack, passing over foes in a direction."
	case DefenderFlail:
		text = "The defender flail is a one-handed weapon that moves foes toward you, and hits harder as you keep attacking without moving."
	}
	return fmt.Sprintf("%s It can hit for up to %d damage.", text, wp.Attack())
}

func (wp weapon) Attack() int {
	switch wp {
	case Axe, Spear, Sabre, DancingRapier:
		return 11
	case BerserkSword:
		return 17
	case BattleAxe, Halberd:
		return 15
	case Frundis, HarKarGauntlets:
		return 13
	case DefenderFlail:
		return 10
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
	case BattleAxe, Halberd, BerserkSword, Frundis, HarKarGauntlets:
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
		text = "Fs"
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
		text = "A fire shield blocks attacks, sometimes burning foliage."
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
