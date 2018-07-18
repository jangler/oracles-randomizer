package main

// conditions required to kill enemies
//
// roughly in the order encountered as *necessary* to kill in the any% route,
// with other enemies afterward

// general functions

// for pushing enemies into pits
//
// mind that all of these items don't work on all enemies that can be pushed.
// these are just the common conditions that work on most enemies.
func (s *state) pushItem(thrownObj bool) bool {
	// fist ring doesn't push hardhats for some reason
	return s.sword || s.shield || s.sustainBombs() || s.rod || s.shovel ||
		(thrownObj && s.bracelet) || ((s.satchel || s.slingshot) &&
		(s.sustainMystery() || s.sustainScent()))
}

// a bunch of common enemies are vulnerable to the same things
func (s *state) killNormalEnemy(animal, pit, thrownObj bool) bool {
	return s.sword || s.sustainBombs() || (pit && s.pushItem(thrownObj)) ||
		s.useFists() || (thrownObj && s.bracelet) ||
		(animal && s.callAnimal()) ||
		((s.satchel || s.slingshot) && (s.sustainEmber() ||
			s.sustainMystery() || s.sustainScent() || s.sustainGale()))
}

// any% route enemies

func (s *state) killStalfos(animal, pit, thrownObj bool) bool {
	return s.rod || s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killGoriyaBros() bool {
	return s.sword || s.sustainBombs() || s.useFists()
}

func (s *state) killGoriya(animal, pit, thrownObj bool) bool {
	return s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killAquamentus() bool {
	return s.sword || s.sustainBombs() || s.useFists() ||
		((s.satchel || s.slingshot) &&
			(s.sustainMystery() || s.sustainScent()))
}

func (s *state) killRope(animal, pit, thrownObj bool) bool {
	return s.killNormalEnemy(animal, pit, thrownObj)
}

func (s *state) killHardhat(pit, thrownObj bool) bool {
	// still going to ignore magnet…
	return (pit && s.pushItem(thrownObj)) || ((s.satchel || s.slingshot) &&
		s.sustainGale())
}

// for the bracelet room in d2
func (s *state) killGapMoblin() bool {
	// the bracelet works because you can throw a pot from the previous room
	return s.sword || s.sustainBombs() || s.useFists() ||
		s.bracelet || (s.satchel && s.sustainScent()) || (s.slingshot &&
		(s.sustainEmber() || s.sustainMystery() || s.sustainScent() ||
			s.sustainGale())) ||
		(s.feather && (s.killNormalEnemy(false, true, false)))
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

func (s *state) killMoldorm() bool {
	return s.sword || s.sustainBombs() || s.useFists() ||
		((s.satchel || s.slingshot) && s.sustainScent())
}

// non-any% enemies

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
