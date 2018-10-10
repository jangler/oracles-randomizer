package logic

// these nodes do not define items, only which items can kill which enemies
// under what circumstances, assuming that you've arrived in the room
// containing the enemy.
//
// anything that can be destroyed in more than one way is also included in
// here. bushes, flowers, mushrooms, etc.
//
// don't worry about thrown objects, sword beams, mystery seeds, or punch.
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
// - what pushes them into pits (if applicable)
//   - sword
//   - shield
//   - boomerangs (they work on hardhats!)
//   - seeds (satchel first, slingshot if satchel doesn't work)
//   - rod
//   - bombs (hard only)
//   - NOT magnet ball; it kills anything pittable
//   - fool's ore

var killNodes = map[string]*Node{
	"satchel kill normal": And("satchel",
		Or("ember seeds", HardOr("scent seeds", "gale seeds"))),
	"slingshot kill normal": And("slingshot",
		Or("ember seeds", "scent seeds", "gale seeds")),
	"scent kill normal": And("scent seeds", Or("slingshot", Hard("satchel"))),
	"seed kill normal":  Or("slingshot kill normal", "satchel kill normal"),
	"jump kill normal":  And("jump 2", "kill normal"),
	"jump pit normal":   And("jump 2", "pit kill normal"),

	// the "safe" version is for areas where you can't possibly get stuck from
	// being on the wrong side of a bush.
	"remove bush safe": Or("sword", "boomerang L-2", "bracelet",
		"ember seeds", "gale slingshot", "bombs"),
	"remove bush": Or("sword", "boomerang L-2", "bracelet"),

	"kill normal": Or("sword", "seed kill normal", "fool's ore",
		Hard("bombs")),
	"pit kill normal": Or("sword", "shield", "rod", "fool's ore",
		Hard("bombs"), "scent kill normal"),
	"kill stalfos": Or("kill normal", "rod"),
	"hit lever": Or("sword", "boomerang", "rod", "ember seeds",
		"scent seeds", "any slingshot", "fool's ore", "shovel"),
	"kill goriya bros":  Or("sword", Hard("bombs"), "fool's ore"),
	"kill goriya":       Or("kill normal"),
	"kill goriya (pit)": Or("kill goriya", "pit kill normal"),
	"kill aquamentus": Or("sword", "fool's ore", Hard("bombs"),
		"scent kill normal"),
	"hit far switch": Or("boomerang", "bombs", "any slingshot"),
	"kill rope":      Or("kill normal"),
	"kill hardhat (pit)": Or("sword", "boomerang", "shield", "rod",
		"fool's ore", Hard("bombs"), And(
			Or("slingshot", Hard("satchel")), Or("scent seeds", "gale seeds"))),
	"kill moblin (gap)": Or("sword", "scent seeds", "slingshot kill normal",
		Hard("bombs"), "fool's ore", "jump kill normal", "jump pit normal"),
	"kill zol":           Or("kill normal"),
	"remove pot":         Or("sword L-2", "bracelet"),
	"kill facade":        Or("bombs"),
	"flip spiked beetle": Or("shield", "shovel"),
	"hit spiked beetle": Or("sword", Hard("bombs"), "fool's ore",
		"seed kill normal"),
	"flip kill spiked beetle": And("flip spiked beetle", "hit spiked beetle"),
	"kill spiked beetle": Or("flip kill spiked beetle", "gale slingshot",
		Hard("gale seeds")),
	"kill mimic": Or("kill normal"),
	"damage omuai": Or("sword", Hard("bombs"), "scent kill normal",
		"fool's ore"),
	"kill omuai": And("damage omuai", "bracelet"),
	"damage mothula": Or("sword", Hard("bombs"), "scent kill normal",
		"fool's ore"),
	"kill mothula":  And("damage mothula", Or("jump 2", Hard("start"))),
	"remove flower": Or("sword", "boomerang L-2"),
	"damage agunima": Or("sword", "scent kill normal", Hard("bombs"),
		"fool's ore"),
	"kill agunima":       And("ember seeds", "damage agunima"),
	"hit very far lever": Or("boomerang L-2", "any slingshot"),
	"hit far lever": Or("boomerang", "any slingshot",
		HardAnd("jump 2", Or("sword", "rod", "fool's ore"))),
	"kill gohma":      Or("scent seeds", "ember seeds"),
	"remove mushroom": Or("boomerang L-2", "bracelet"),
	"kill moldorm": Or("sword", Hard("bombs"), "scent kill normal",
		"fool's ore"),
	"kill iron mask": Or("kill normal"),
	"kill armos": Or("sword", Hard("bombs"), "boomerang L-2",
		"scent kill normal", "fool's ore"),
	"kill gibdo": Or("kill normal", "boomerang L-2", "rod"),
	"kill darknut": Or("sword", Hard("bombs"), "scent kill normal",
		"fool's ore"),
	"kill darknut (pit)": Or("kill darknut", "shield"),
	"kill syger": Or("sword", Hard("bombs"), "scent kill normal",
		"fool's ore"),
	"break crystal": Or("sword", "bombs", "bracelet"),
	"kill hardhat (magnet)": Or("magnet gloves", "gale slingshot",
		Hard("gale satchel")),
	"kill vire": Or("sword", Hard("bombs"), "fool's ore"),
	"finish manhandla": Or("sword", Hard("bombs"), "any slingshot",
		"fool's ore"),
	"kill manhandla": And("boomerang L-2", "finish manhandla"),
	"kill wizzrobe":  Or("kill normal"),
	"kill magunesu":  Or("sword", "fool's ore"), // even bombs don't work!
	"kill poe sister": Or("sword", "ember seeds", "scent kill normal",
		Hard("bombs"), "fool's ore"),
	"kill darknut (across pit)": Or("scent slingshot", "magnet gloves",
		And("jump 4", "kill darknut (pit)")),
	"kill gleeok": Or("sword", Hard("bombs"), "fool's ore"),
	"kill frypolar": Or(And("bracelet",
		Or("mystery slingshot", Hard("mystery satchel"))),
		Or("ember slingshot", Hard("ember satchel"))),
	"kill medusa head": Or("sword", "fool's ore"),
	"kill floormaster": Or("kill normal"),
	"kill onox":        And("sword", "jump 2"),
}
