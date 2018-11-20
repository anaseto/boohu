package main

import (
	"errors"
	"fmt"
	"sort"
)

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
	TormentPotion
	AccuracyPotion
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
	case TormentPotion:
		text += " of torment explosion"
	case AccuracyPotion:
		text += " of accuracy"
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
		text = "makes you enter a crazy rage, temporarily making you faster, stronger and healthier. You cannot use rods while berserk, and afterwards it leaves you slow and exhausted."
	case SwiftnessPotion:
		text = "makes you move faster and better at avoiding blows for a short time."
	case LignificationPotion:
		text = "makes you more resistant to physical blows, but you are attached to the ground while the effect lasts (you can still descend)."
	case WallPotion:
		text = "replaces free cells around you with temporary walls."
	case CBlinkPotion:
		text = "makes you blink to a targeted cell in your line of sight."
	case DigPotion:
		text = "makes you dig walls like an earth dragon."
	case SwapPotion:
		text = "makes you swap positions with monsters instead of attacking. Ranged monsters can still damage you."
	case ShadowsPotion:
		text = "reduces your line of sight range to 1."
	case TormentPotion:
		text = "halves HP of every creature in sight, including the player, and destroys visible walls. Extremely noisy. It can burn foliage and doors."
	case AccuracyPotion:
		text = "makes you never miss for a few turns."
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
	case TormentPotion:
		err = g.QuaffTormentPotion(ev)
	case AccuracyPotion:
		err = g.QuaffAccuracyPotion(ev)
	case DreamPotion:
		err = g.QuaffDreamPotion(ev)
	}
	if err != nil {
		return err
	}
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
	if g.Player.HasStatus(StatusBerserk) {
		return errors.New("You are already berserk.")
	}
	g.Player.Statuses[StatusBerserk] = 1
	end := ev.Rank() + 65 + RandInt(20)
	g.PushEvent(&simpleEvent{ERank: end, EAction: BerserkEnd})
	g.Player.Expire[StatusBerserk] = end
	g.Printf("You quaff the %s. You feel a sudden urge to kill things.", BerserkPotion)
	g.Player.HP += 10
	return nil
}

func (g *game) QuaffHealWounds(ev event) error {
	hp := g.Player.HP
	g.Player.HP += 2 * DefaultHealth / 3
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

func (g *game) QuaffLignification(ev event) error {
	if g.Player.HasStatus(StatusLignification) {
		return errors.New("You are already lignified.")
	}
	g.EnterLignification(ev)
	g.Printf("You quaff the %s. You feel rooted to the ground.", LignificationPotion)
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

func (g *game) QuaffAccuracyPotion(ev event) error {
	g.Player.Statuses[StatusAccurate]++
	end := ev.Rank() + 85 + RandInt(20)
	g.PushEvent(&simpleEvent{ERank: end, EAction: AccurateEnd})
	g.Player.Expire[StatusAccurate] = end
	g.Printf("You quaff the %s. You feel accurate.", SwiftnessPotion)
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
	SlowingMagara
	ConfuseMagara
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
	case SlowingMagara:
		text = "slowing magara"
	case ConfuseMagara:
		text = "confusion magara"
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
	case SlowingMagara:
		text = "slowing magaras"
	case ConfuseMagara:
		text = "confusion magaras"
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
		text = "can be thrown to cause a fire explosion halving HP of monsters in a square area. It can occasionally destruct walls. It can burn doors and foliage."
	case TeleportMagara:
		text = "can be thrown to make monsters in a square area teleport."
	case SlowingMagara:
		text = "can be activated to release a slowing bolt inducing slow movement and attack in one or more foes."
	case ConfuseMagara:
		text = "generates a harmonic light that confuses monsters in your line of sight."
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
	case SlowingMagara:
		err = g.ThrowSlowingMagara(ev)
	case ConfuseMagara:
		err = g.ThrowConfuseMagara(ev)
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
	g.HandleStone(mons)
	ev.Renew(g, 7)
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

	ev.Renew(g, 7)
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

	ev.Renew(g, 7)
	return nil
}

func (g *game) ThrowSlowingMagara(ev event) error {
	if err := g.ui.ChooseTarget(g, &chooser{}); err != nil {
		return err
	}
	ray := g.Ray(g.Player.Target)
	g.MakeNoise(MagicCastNoise, g.Player.Pos)
	g.Print("Whoosh! A bolt of slowing emerges out of the magara.")
	g.ui.SlowingMagaraAnimation(g, ray)
	for _, pos := range ray {
		mons := g.MonsterAt(pos)
		if !mons.Exists() {
			continue
		}
		mons.Statuses[MonsSlow]++
		g.PushEvent(&monsterEvent{ERank: g.Ev.Rank() + 130 + RandInt(40), NMons: mons.Index, EAction: MonsSlowEnd})
	}

	ev.Renew(g, 7)
	return nil
}

func (g *game) ThrowConfuseMagara(ev event) error {
	g.Printf("You activate the %s. A harmonic light confuses monsters.", ConfuseMagara)
	for pos, b := range g.Player.LOS {
		if !b {
			continue
		}
		mons := g.MonsterAt(pos)
		if mons.Exists() {
			mons.EnterConfusion(g, ev)
		}
	}

	ev.Renew(g, 7)
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
			g.MakeCreatureSleep(pos, ev)
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

	ev.Renew(g, 7)
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
	ConfusingDart:       {rarity: 4, quantity: 2},
	ExplosiveMagara:     {rarity: 6, quantity: 1},
	NightMagara:         {rarity: 9, quantity: 1},
	TeleportMagara:      {rarity: 12, quantity: 1},
	SlowingMagara:       {rarity: 12, quantity: 1},
	ConfuseMagara:       {rarity: 15, quantity: 1},
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
	DescentPotion:       {rarity: 18, quantity: 1},
	MagicMappingPotion:  {rarity: 18, quantity: 1},
	DreamPotion:         {rarity: 18, quantity: 1},
	TormentPotion:       {rarity: 30, quantity: 1},
	AccuracyPotion:      {rarity: 18, quantity: 1},
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
	case SmokingScales:
		text = "Smoking scales provide protection against blows. They leave short-lived fog as you move."
	case ShinyPlates:
		text = "Shiny plates provide good protection against blows, but increase your line of sight range."
	case TurtlePlates:
		text = "Turtle plates provide great protection against blows, but make you move slower and a little less good at evading blows."
	case SpeedRobe:
		text = "The speed robe makes you move faster, with a minor evasion bonus."
	case CelmistRobe:
		text = "The celmist robe improves your magic reserves, rod recharge rate, and rods can gain two extra charges. In Hareka, celmists are what most people would call mages."
	case HarmonistRobe:
		text = "The harmonist robe makes you harder to detect (reduced LOS, stealthy, noise mitigation). Harmonists are mages specialized in manipulation of light and noise."
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
	AssassinSabre
	DancingRapier
	HopeSword
	Frundis
	ElecWhip
	HarKarGauntlets
	VampDagger
	DragonSabre
	FinalBlade
	DefenderFlail
)

const WeaponNum = int(DefenderFlail) + 1

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
	case AssassinSabre:
		return "assassin sabre"
	case DancingRapier:
		return "dancing rapier"
	case HopeSword:
		return "hopeful sword"
	case Frundis:
		return "staff Frundis"
	case ElecWhip:
		return "lightning whip"
	case HarKarGauntlets:
		return "har-kar gauntlets"
	case VampDagger:
		return "vampiric dagger"
	case DragonSabre:
		return "dragon sabre"
	case FinalBlade:
		return "final blade"
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
	case AssassinSabre:
		return "Sb"
	case DancingRapier:
		return "Dn"
	case HopeSword:
		return "Ds"
	case Frundis:
		return "Fr"
	case ElecWhip:
		return "Wh"
	case HarKarGauntlets:
		return "Hk"
	case VampDagger:
		return "Vm"
	case DragonSabre:
		return "Dr"
	case FinalBlade:
		return "Fn"
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
	case AssassinSabre:
		text = "The assassin sabre is a one-handed weapon. It is more accurate against injured opponents."
	case DancingRapier:
		text = "The dancing rapier is a one-handed weapon. It makes you swap with your foe and can hit another monster behind with extra damage."
	case HopeSword:
		text = "The hopeful sword is a big two-handed weapon that hits with extra damage when you are injured."
	case Frundis:
		text = "Frundis is a musician and harmonist, which happens to be a two-handed staff too. It may occasionally confuse monsters on hit. It magically helps reducing noise in combat too."
	case ElecWhip:
		text = "The lightning whip is a one-handed weapon that inflicts electrical damage to a monster and any foes connected to it."
	case HarKarGauntlets:
		text = "Har-kar gauntlets are an unarmed combat weapon. They allow you to make a wind attack, passing over foes in a direction."
	case VampDagger:
		text = "The vampiric dagger is a one-handed weapon that gives you some healing when you hit living monsters."
	case DragonSabre:
		text = "The dragon sabre is a one-handed weapon that inflicts extra damage on healthy big monsters."
	case FinalBlade:
		text = "The final blade is an accurate two-handed weapon that instantly kills monsters at less than half full health. Wielding this weapon hurts your maximum health."
	case DefenderFlail:
		text = "The defender flail is a one-handed weapon that moves foes toward you, and hits harder as you keep attacking without moving."
	}
	return fmt.Sprintf("%s It can hit for up to %d damage.", text, wp.Attack())
}

func (wp weapon) Attack() int {
	switch wp {
	case Axe, Spear, AssassinSabre, DancingRapier, DragonSabre:
		return 11
	case BattleAxe, Halberd, HopeSword, FinalBlade:
		return 15
	case Frundis, HarKarGauntlets:
		return 13
	case DefenderFlail:
		return 10
	case Dagger, VampDagger:
		return 9
	case ElecWhip:
		return 8
	default:
		return 0
	}
}

func (wp weapon) TwoHanded() bool {
	switch wp {
	case BattleAxe, Halberd, HopeSword, Frundis, HarKarGauntlets, FinalBlade:
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
	case ConfusingShield:
		text = "A confusing shield can block an attack, sometimes confusing the monster."
	case EarthShield:
		text = "An earth shield offers great protection, but impact sound can disintegrate nearby walls."
	case BashingShield:
		text = "A bashing shield can block an attack and push the ennemy away."
	case FireShield:
		text = "A fire shield can block an attack, sometimes burning nearby foliage."
	}
	return text
}

func (sh shield) Letter() rune {
	return ']'
}

func (sh shield) Block() (block int) {
	switch sh {
	case ConfusingShield, BashingShield, FireShield:
		block += 10
	case EarthShield:
		block += 15
	}
	return block
}
