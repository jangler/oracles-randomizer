package logic

// these nodes do not define items, only which items can kill which enemies
// under what circumstances, assuming that you've arrived in the room
// containing the enemy.
//
// anything that can be destroyed in more than one way is also included in
// here. bushes, flowers, mushrooms, etc.
//
// don't worry about thrown objects, sword beams, or mystery seeds
//
// animal companions are not included in this logic, since they're only
// available in certain areas.

// when testing how to kill enemies, remember to try:
// - sword
// - boomerang
// - magic boomerang
// - rod
// - seeds (satchel first, slingshot if satchel doesn't work)
// - bombs (hard only)
// - magnet ball (if applicable)
// - fool's ore
// - punch
// - what pushes them into pits (if applicable)
//   - sword
//   - shield
//   - boomerangs (they work on hardhats!)
//   - seeds (satchel first, slingshot if satchel doesn't work)
//   - rod
//   - bombs (hard only)
//   - NOT magnet ball; it kills anything pittable
//   - fool's ore

var seasonsKillNodes = map[string]*Node{
	"satchel kill normal": And("satchel",
		Or("ember seeds", And("hard", Or("scent seeds", "gale seeds")))),
	"slingshot kill normal": And("slingshot",
		Or("ember seeds", "scent seeds", "gale seeds")),
	"jump kill normal": And("jump 2", "kill normal"),
	"jump pit normal":  And("jump 2", "pit kill normal"),

	// enemies vulnerable to scent seeds are always vulnerable to sword, bombs,
	// and fool's ore (and punches?).
	"kill armored": Or("sword", "fool's ore", "punch enemy",
		And("hard", "bombs"),
		And("scent seeds", Or("slingshot", And("hard", "satchel")))),

	// the "safe" version is for areas where you can't possibly get stuck from
	// being on the wrong side of a bush.
	"remove bush safe": Or("sword", "magic boomerang", "bracelet",
		"ember seeds", "gale slingshot", "bombs"),
	"remove bush": Or("sword", "magic boomerang", "bracelet"),

	"kill normal": Or("sword", "satchel kill normal", "slingshot kill normal",
		"fool's ore", "punch enemy", And("hard", "bombs")),
	"pit kill normal": Or("sword", "shield", "rod", "fool's ore",
		And("hard", "bombs"), "kill armored"),
	"kill stalfos": Or("kill normal", "rod"),
	"hit lever from minecart": Or("sword", "boomerang", "rod", "scent seeds",
		"mystery seeds", "any slingshot", "fool's ore", "punch object"),
	"hit lever":         Or("shovel", "hit lever from minecart"),
	"kill normal (pit)": Or("kill normal", "pit kill normal"),
	"hit far switch": Or("boomerang", "bombs", "any slingshot",
		And("sword", "energy ring")),
	"kill hardhat (pit)": Or("sword", "boomerang", "shield", "rod",
		"fool's ore", And("hard", "bombs"), And(
			Or("slingshot", And("hard", "satchel")),
			Or("scent seeds", "gale seeds"))),
	"kill moblin (gap)": Or("sword", "scent seeds", "slingshot kill normal",
		"fool's ore", "jump kill normal", "jump pit normal",
		And("hard", Or("bombs", "punch enemy", "ember seeds"))),
	"remove pot":         Or("noble sword", "bracelet"),
	"flip spiked beetle": Or("shield", "shovel"),
	// spiked beetles can't be punched for some reason
	"flip kill spiked beetle": And("flip spiked beetle", Or("sword", "fool's ore",
		"satchel kill normal", "slingshot kill normal", And("hard", "bombs"))),
	"kill spiked beetle": Or("flip kill spiked beetle", "gale slingshot",
		And("hard", "gale seeds")),
	"kill omuai":         And("kill armored", "bracelet"),
	"remove flower":      Or("sword", "magic boomerang"),
	"kill agunima":       And("ember seeds", "kill armored"),
	"hit very far lever": Or("magic boomerang", "any slingshot"),
	"hit far lever": Or("boomerang", "any slingshot",
		And("hard", "jump 2", Or("sword", "rod", "fool's ore"))),
	"kill gohma":         Or("scent seeds", "ember seeds"),
	"remove mushroom":    Or("magic boomerang", "bracelet"),
	"kill moldorm":       Or("kill armored"),
	"kill armored (pit)": Or("kill armored", "shield"),
	"break crystal":      Or("sword", "bombs", "bracelet", "expert's ring"),
	"kill hardhat (magnet)": Or("magnet gloves", "gale slingshot",
		And("hard", "gale satchel")),
	"kill vire": Or("sword", "fool's ore", "expert's ring",
		And("hard", "bombs")),
	"finish manhandla": Or("sword", And("hard", "bombs"), "any slingshot",
		"fool's ore", "expert's ring"),
	"kill manhandla":  And("magic boomerang", "finish manhandla"),
	"kill poe sister": Or("kill armored", "rod", "ember seeds"),
	"kill magunesu":   Or("sword", "fool's ore", "expert's ring"),
	"kill armored (across pit)": Or("scent slingshot", "magnet gloves",
		And("jump 4", "kill armored (pit)"), And("sword", "energy ring"),
		And("hard", "toss ring", "bombs")),
	"kill gleeok": Or("sword", "fool's ore", "punch enemy"),
	"kill frypolar": Or(And("bracelet",
		Or("mystery slingshot", And("hard", "mystery satchel"))),
		Or("ember slingshot", And("hard", "ember satchel"))),
	"kill medusa head": Or("sword", "fool's ore"),
	"kill onox":        And("sword", "jump 2"),
}
