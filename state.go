package main

type state struct {
	// equip items, in the order you normally receive them
	sword          bool
	swordLevel     int // 1-3 if sword
	shield         bool
	shieldLevel    int // 1-3 if shield
	bombs          bool
	satchel        bool
	boomerang      bool
	boomerangLevel int // 1-2 if boomerang
	rod            bool
	winter         bool
	spring         bool
	summer         bool
	autumn         bool // less ambiguous than "fall"
	shovel         bool
	bracelet       bool
	flute          bool
	fluteLevel     int // if flute, 0 = strange flute, 1 = companion
	feather        bool
	featherLevel   int // 1-2 if feather
	slingshot      bool
	slingshotLevel int // 1-2 if slingshot
	magnetGloves   bool

	// trees, in the order you normally access them. just being able to access
	// a tree doesn't mean you get the seeds; you need a satchel and a harvest
	// item
	emberTree   bool
	mysteryTree bool
	scentTree   bool
	pegasusTree bool
	galeTree    bool

	// collection items, in the order you normally receive them
	gnarledKey   bool
	floodgateKey bool
	flippers     bool
	dragonKey    bool
	// TODO: others?

	// potential route-altering rings
	//
	// ring box as a collection item can be safely ignored since you get it for
	// free from the very start
	//
	// you *are* guaranteed to have rupees by the time you get these rings,
	// since the hero's cave item must be able to destroy bushes and/or burn
	// down trees
	fistRing      bool
	fistRingLevel int // 1 = fist ring, 2 = expert's ring
	energyRing    bool
	tossRing      bool
}

// magnet glove is not included in any of this (yet) because it's so rare

func (s *state) sustainRupees() bool {
	// fist ring is not included since you need rupees to appraise it in the
	// first place. bombs are not included since you might need rupees to buy
	// them
	return s.sword ||
		(s.satchel &&
			(s.sustainEmber() ||
				s.sustainMystery() ||
				s.sustainScent())) ||
		s.boomerang ||
		s.rod ||
		s.shovel ||
		s.bracelet ||
		(s.flute && s.fluteLevel > 0) ||
		(s.slingshot &&
			(s.sustainEmber() ||
				s.sustainMystery() ||
				s.sustainScent()))
}

func (s *state) popMakuBubble() bool {
	return s.sword || s.rod ||
		(s.satchel && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent())) ||
		(s.slingshot && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent() || s.sustainPegasus() || s.sustainGale()))
}

func (s *state) breakBushes(animal bool) bool {
	return s.sword || s.sustainBombs() ||
		(s.boomerang && s.boomerangLevel == 2) || s.bracelet ||
		(animal && s.flute && s.fluteLevel > 0) ||
		(s.satchel && (s.sustainEmber() || s.sustainMystery())) ||
		(s.slingshot && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainGale()))
}

func (s *state) killStalfos(thrownObj bool) bool {
	return s.sword || s.sustainBombs() || s.rod || (thrownObj && s.bracelet) ||
		(s.fistRing && s.sustainRupees()) ||
		(s.satchel && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent())) ||
		(s.slingshot && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent() || s.sustainGale()))
}

func (s *state) hitLever() bool {
	return s.sword || s.boomerang || s.rod || s.fistRing ||
		(s.satchel && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent())) ||
		(s.slingshot && (s.sustainEmber() || s.sustainMystery() ||
			s.sustainScent() || s.sustainPegasus() || s.sustainGale()))
}

func (s *state) sustainBombs() bool {
	return s.sustainRupees()
	// TODO: check for sustainable areas
	// - bracelet room in d2 (requires killRope && lightTorch && killHardhat &&
	//   && harvestPlant)
	// - bomb wall room in d2 (requires ((lightTorch && killRope && killMoblin
	//   && killGel) || (removeRock && pushRoller)) && breakBush)
}

func (s *state) destroyCracked(animal bool) bool {
	return s.bombs || (animal && s.flute && s.fluteLevel > 0)
}

func (s *state) killGoriyaBros() bool {
	return s.sword || s.sustainBombs() || (s.fistRing && s.sustainRupees())
}

func (s *state) lightTorch() bool {
	return (s.satchel || s.slingshot) && s.sustainEmber()
}

func (s *state) killNormalEnemy(animal, pit, thrownObj bool) bool {
	return s.sword || s.sustainBombs() || (pit && s.pushItem(thrownObj)) ||
		(s.fistRing && s.sustainRupees()) || (thrownObj && s.bracelet) ||
		(animal && s.flute && s.fluteLevel > 0) ||
		((s.satchel || s.slingshot) && (s.sustainEmber() ||
			s.sustainMystery() || s.sustainScent() || s.sustainGale()))
}

func (s *state) killGoriya(animal, pit, thrownObj bool) bool {
	return s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killAquamentus() bool {
	return s.sword || s.sustainBombs() || (s.fistRing && s.sustainRupees()) ||
		((s.satchel || s.slingshot) &&
			(s.sustainMystery() || s.sustainScent()))
}

func (s *state) destroyTree() bool {
	return s.lightTorch()
}

func (s *state) killRope(animal, pit, thrownObj bool) bool {
	return s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killHardhat(pit, thrownObj bool) bool {
	// still going to ignore magnet…
	return (pit && s.pushItem(thrownObj)) || ((s.satchel || s.slingshot) &&
		s.sustainGale())
}

// this refers to the renewable plants that drop an item (or not) and
// regenerate after some seconds
//
// do *not* call this if you want to find out if you can harvest bombs, since
// it will lead to a dependency loop
func (s *state) harvestPlant() bool {
	return s.sword || s.sustainBombs()
}

// for the bracelet room in d2
func (s *state) killGapMoblin() bool {
	// the bracelet works because you can throw a pot from the previous room
	return s.sword || s.sustainBombs() || (s.fistRing && s.sustainRupees()) ||
		s.bracelet || (s.satchel && s.sustainScent()) || (s.slingshot &&
		(s.sustainEmber() || s.sustainMystery() || s.sustainScent() ||
			s.sustainGale())) ||
		(s.feather && (s.pushItem(false) || (s.satchel &&
			// scent and mystery seeds are push items already
			(s.sustainEmber() || (s.sustainGale())))))
}

func (s *state) removeRock() bool {
	return s.bracelet
}

func (s *state) pushRoller() bool {
	return s.bracelet
}

func (s *state) removePot() bool {
	return s.bracelet || (s.sword && s.swordLevel > 1)
}

func (s *state) killFacade() bool {
	return s.sustainBombs() || s.killBeetle()
}

// spawned by facade
func (s *state) killBeetle() bool {
	return s.killNormalEnemy(false, false, false)
}

func (s *state) killDodongo() bool {
	return s.sustainBombs() && s.bracelet
}

func (s *state) killMoblin(animal, pit, thrownObj bool) bool {
	return s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killGel(pit, thrownObj bool) bool {
	// gels are immune to satchel gale seeds (but not slingshot ones) and
	// pushing for some reason, but can be lured into pits without any item at
	// all
	return s.sword || s.sustainBombs() || pit ||
		(s.fistRing && s.sustainRupees()) || (thrownObj && s.bracelet) ||
		((s.satchel || s.slingshot) && (s.sustainEmber() ||
			s.sustainMystery() || s.sustainScent())) ||
		(s.slingshot && s.sustainGale())
}

// true iff the player has an item that can harvest seeds from trees
func (s *state) harvestItem() bool {
	return s.sword || s.rod || s.fistRing
}

func (s *state) sustainEmber() bool {
	return s.satchel && s.emberTree && s.harvestItem()
	// TODO: check for sustainable areas
	// - in d1, as long as you can kill stalfos
	// - bushes outside d2
}

func (s *state) sustainMystery() bool {
	return s.satchel && s.emberTree && s.harvestItem()
	// TODO: check for sustainable areas
	// - bomb wall room in d2 (requires ((lightTorch && killRope && killMoblin
	//   && killGel) || (removeRock && pushRoller)) && breakBush)
}

func (s *state) sustainScent() bool {
	return s.satchel && s.scentTree && s.harvestItem()
	// TODO: check for sustainable areas
}

func (s *state) sustainPegasus() bool {
	return s.satchel && s.pegasusTree && s.harvestItem()
	// TODO: check for sustainable areas
	// - pots before jumps in d4
	// - pot before "indiana jones" in d6
	// - various pots in d7
	// - grass on the way to d8
	// - pots before the boss in d8
}

func (s *state) sustainGale() bool {
	return s.satchel && s.galeTree && s.harvestItem()
	// TODO: check for sustainable areas
}

func (s *state) pushItem(thrownObj bool) bool {
	// fist ring doesn't push hardhats for some reason
	return s.sword || s.shield || s.sustainBombs() || s.rod || s.shovel ||
		(thrownObj && s.bracelet) || ((s.satchel || s.slingshot) &&
		(s.sustainMystery() || s.sustainScent()))
}
