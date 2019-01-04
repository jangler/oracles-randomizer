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
// - boomerang L-1
// - boomerang L-2
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
		Or("ember seeds", HardOr("scent seeds", "gale seeds"))),
	"slingshot kill normal": And("slingshot",
		Or("ember seeds", "scent seeds", "gale seeds")),
	"jump kill normal": And("jump 2", "kill normal"),
	"jump pit normal":  And("jump 2", "pit kill normal"),

	// enemies vulnerable to scent seeds are always vulnerable to sword, bombs,
	// and fool's ore (and punches?).
	"scent kill normal": Or("sword", Hard("bombs"), "fool's ore", "punch enemy",
		And("scent seeds", Or("slingshot", Hard("satchel")))),

	// the "safe" version is for areas where you can't possibly get stuck from
	// being on the wrong side of a bush.
	"remove bush safe": Or("sword", "boomerang L-2", "bracelet",
		"ember seeds", "gale slingshot", "bombs"),
	"remove bush": Or("sword", "boomerang L-2", "bracelet"),

	"kill normal": Or("sword", "satchel kill normal", "slingshot kill normal",
		"fool's ore", "punch enemy", Hard("bombs")),
	"pit kill normal": Or("sword", "shield", "rod", "fool's ore",
		Hard("bombs"), "scent kill normal"),
	"kill stalfos": Or("kill normal", "rod"),
	"hit lever from minecart": Or("sword", "boomerang", "rod", "ember seeds",
		"scent seeds", "mystery seeds", "any slingshot", "fool's ore",
		"punch object"),
	"hit lever":         Or("shovel", "hit lever from minecart"),
	"kill goriya bros":  Or("scent kill normal"),
	"kill goriya":       Or("kill normal"),
	"kill goriya (pit)": Or("kill goriya", "pit kill normal"),
	"kill aquamentus":   Or("scent kill normal"),
	"hit far switch": Or("boomerang", "bombs", "any slingshot",
		And("sword", "energy ring")),
	"kill rope": Or("kill normal"),
	"kill hardhat (pit)": Or("sword", "boomerang", "shield", "rod",
		"fool's ore", Hard("bombs"), And(
			Or("slingshot", Hard("satchel")), Or("scent seeds", "gale seeds"))),
	"kill moblin (gap)": Or("sword", "scent seeds", "slingshot kill normal",
		"fool's ore", "jump kill normal", "jump pit normal",
		HardOr("bombs", "punch enemy")), // need to fall into pit to punch (?)
	"kill zol":           Or("kill normal"),
	"remove pot":         Or("sword L-2", "bracelet"),
	"kill facade":        Or("bombs"),
	"flip spiked beetle": Or("shield", "shovel"),
	// spiked beetles can't be punched for some reason
	"flip kill spiked beetle": And("flip spiked beetle", Or("sword", "fool's ore",
		"satchel kill normal", "slingshot kill normal", Hard("bombs"))),
	"kill spiked beetle": Or("flip kill spiked beetle", "gale slingshot",
		Hard("gale seeds")),
	"kill mimic":         Or("kill normal"),
	"damage omuai":       Or("scent kill normal"),
	"kill omuai":         And("damage omuai", "bracelet"),
	"kill mothula":       Or("scent kill normal"),
	"remove flower":      Or("sword", "boomerang L-2"),
	"kill agunima":       And("ember seeds", "scent kill normal"),
	"hit very far lever": Or("boomerang L-2", "any slingshot"),
	"hit far lever": Or("boomerang", "any slingshot",
		HardAnd("jump 2", Or("sword", "rod", "fool's ore"))),
	"kill gohma":      Or("scent seeds", "ember seeds"),
	"remove mushroom": Or("boomerang L-2", "bracelet"),
	"kill moldorm":    Or("scent kill normal"),
	"kill iron mask":  Or("kill normal"),
	// armos are an exception to the "bombs are hard logic" rule, since you're
	// intended to kill them with bombs in ages.
	"kill armos":         Or("scent kill normal", "boomerang L-2", "bombs"),
	"kill gibdo":         Or("kill normal", "boomerang L-2", "rod"),
	"kill darknut":       Or("scent kill normal"),
	"kill darknut (pit)": Or("kill darknut", "shield"),
	"kill syger":         Or("scent kill normal"),
	"break crystal":      Or("sword", "bombs", "bracelet", "expert's ring"),
	"kill hardhat (magnet)": Or("magnet gloves", "gale slingshot",
		Hard("gale satchel")),
	"kill vire": Or("sword", Hard("bombs"), "fool's ore", "expert's ring"),
	"finish manhandla": Or("sword", Hard("bombs"), "any slingshot",
		"fool's ore", "expert's ring"),
	"kill manhandla":  And("boomerang L-2", "finish manhandla"),
	"kill wizzrobe":   Or("kill normal"),
	"kill magunesu":   Or("sword", "fool's ore", "expert's ring"),
	"kill poe sister": Or("scent kill normal", "ember seeds"),
	"kill darknut (across pit)": Or("scent slingshot", "magnet gloves",
		And("jump 4", "kill darknut (pit)"), And("sword", "energy ring"),
		HardAnd("toss ring", "bombs")),
	"kill gleeok": Or("sword", Hard("bombs"), "fool's ore", "punch enemy"),
	"kill frypolar": Or(And("bracelet",
		Or("mystery slingshot", Hard("mystery satchel"))),
		Or("ember slingshot", Hard("ember satchel"))),
	"kill medusa head": Or("sword", "fool's ore"),
	"kill floormaster": Or("kill normal"),
	"kill onox":        And("sword", "jump 2"),
}
