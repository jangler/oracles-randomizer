package prenode

// these nodes do not define items, only which items can kill which enemies
// under what circumstances, assuming that you've arrived in the room
// containing the enemy.
//
// anything that can be destroyed in more than one way is also included in
// here. bushes, flowers, mushrooms, etc.
//
// technically mystery seeds can be used to kill many enemies that can be
// killed by ember, scent, or gale seeds. mystery seeds are only included as a
// kill option if at all three of these seed types work.
//
// if an enemy is in the same room as a throwable object and is vulnerable to
// thrown objects, than just adding "bracelet" as an OR is sufficient.
//
// animal companions are not included in this logic, since they're only
// available in certain areas.

// when testing how to kill enemies, remember to try:
// - sword
// - beams
// - boomerang L-1
// - boomerang L-2
// - rod
// - seeds (satchel first, slingshot if satchel doesn't work)
// - bombs
// - thrown objects (if applicable)
// - magnet ball (if applicable)
// - fool's ore
// - punch
// - what pushes them into pits (if applicable)
//   - sword
//   - beams
//   - shield
//   - boomerangs (they work on hardhats!)
//   - seeds (satchel first, slingshot if satchel doesn't work)
//   - rod
//   - bombs
//   - shovel
//   - thrown objects (if applicable)
//   - NOT magnet ball; it kills anything pittable
//   - fool's ore
//   - punch

var killPrenodes = map[string]*Prenode{
	"gale seed weapon":        And("gale seeds", Or("slingshot", HardAnd("satchel", "jump"))),
	"gale boomerang":          And("gale satchel", "boomerang"), // stun, then drop from satchel
	"slingshot kill normal":   And("slingshot", "seed kill normal"),
	"jump kill normal":        And("jump", "kill normal"),
	"jump pit normal":         And("jump", "pit kill normal"),
	"slingshot gale seeds":    And("slingshot", "gale seeds"),
	"slingshot mystery seeds": And("slingshot", "mystery seeds"),
	"kill dodongo":            And("bombs", "bracelet"),

	// required enemies in normal route-ish order, but with prereqs first
	"seed kill normal":                Or("ember seeds", "scent seeds", "gale seed weapon", "gale boomerang", "mystery seeds"),
	"pop maku bubble":                 Or("sword", "rod", "seed kill normal", "pegasus slingshot", "bombs", "fool's ore"),
	"remove bush":                     Or("sword", "boomerang L-2", "ember seeds", "gale slingshot", "bracelet", Hard("bombs")),
	"remove bush sustainable":         Or("sword", "boomerang L-2", "bracelet"),
	"kill normal":                     Or("sword", "bombs", "beams", "seed kill normal", "fool's ore", "punch"),
	"pit kill normal":                 Or("sword", "beams", "shield", "scent seeds", "rod", "bombs", Hard("shovel"), "fool's ore", "punch"),
	"kill stalfos":                    Or("kill normal", "rod"),
	"kill stalfos (throw)":            Or("kill stalfos", "bracelet"),
	"hit lever":                       Or("sword", "boomerang", "rod", "ember seeds", "scent seeds", "slingshot", "fool's ore", "punch"),
	"kill goriya bros":                Or("sword", "bombs", "fool's ore", "punch"),
	"kill goriya":                     Or("kill normal"),
	"kill goriya (pit)":               Or("kill goriya", "pit kill normal"),
	"kill aquamentus":                 Or("sword", "beams", "scent seeds", "bombs", "fool's ore", "punch"),
	"hit far switch":                  Or("beams", "boomerang", "bombs", "slingshot"),
	"toss bombs":                      And("bombs", "toss ring"),
	"hit very far switch":             Or("beams", "boomerang", "toss bombs", "slingshot"),
	"kill rope":                       Or("kill normal"),
	"kill hardhat (pit, throw)":       Or("gale seed weapon", "sword", "beams", "boomerang", "shield", "scent seeds", "rod", "bombs", Hard("shovel"), "fool's ore", "bracelet"),
	"kill moblin":                     Or("kill normal"),
	"kill moblin (gap, throw)":        Or("sword", "beams", "scent seeds", "slingshot kill normal", "bombs", "fool's ore", "punch", "jump kill normal", "jump pit normal"),
	"kill zol":                        Or("kill normal"),
	"remove pot":                      Or("sword L-2", "bracelet"),
	"kill facade":                     Or("bombs"),
	"flip spiked beetle":              Or("shield", "shovel"),
	"damage spiked beetle (throw)":    Or("sword", "bombs", "beams", "seed kill normal", "bracelet", "fool's ore"),
	"flip kill spiked beetle (throw)": And("flip spiked beetle", "damage spiked beetle (throw)"),
	"gale kill spiked beetle":         And("gale seed weapon"),
	"kill spiked beetle (throw)":      Or("flip kill spiked beetle (throw)", "gale kill spiked beetle"),
	"kill mimic":                      Or("kill normal"),
	"damage omuai":                    Or("sword", "bombs", "scent seeds", "fool's ore", "punch"),
	"kill omuai":                      And("damage omuai", "bracelet"),
	"damage mothula":                  Or("sword", "bombs", "scent seeds", "fool's ore", "punch"),
	"kill mothula":                    And("damage mothula", "jump"), // you will basically die without feather
	"remove flower":                   Or("sword", "boomerang L-2", "ember seeds", "gale slingshot", Hard("bombs")),
	"remove flower sustainable":       Or("sword", "boomerang L-2"),
	"kill shrouded stalfos (throw)":   Or("kill stalfos", "bracelet"),
	"kill like-like (pit, throw)":     Or("kill normal", "bracelet", "rod", Hard("shovel")),
	"kill water tektite (throw)":      Or("kill normal", "bracelet"),
	"damage agunima":                  Or("sword", "scent seeds", "bombs", "fool's ore", "punch"),
	"kill agunima":                    And("ember seeds", "damage agunima"),
	"hit very far lever":              Or("boomerang L-2", "slingshot"),
	"hit lever gap":                   Or("sword", "boomerang", "rod", "slingshot", "fool's ore"),
	"jump hit lever":                  And("jump", "hit lever gap"),
	"long jump hit lever":             And("long jump", "hit lever"),
	"hit far lever":                   Or("jump hit lever", "long jump hit lever", "boomerang", "slingshot"),
	"kill wizzrobe (pit, throw)":      Or("pit kill normal", "bracelet"),
	"kill gohma":                      Or("scent slingshot", "ember slingshot"),
	"remove mushroom":                 Or("boomerang L-2", "bracelet"),
	"kill moldorm":                    Or("sword", "bombs", "punch", "scent seeds"),
	"kill iron mask":                  Or("kill normal"),
	"kill armos":                      Or("sword", "bombs", "beams", "boomerang L-2", "scent seeds", "fool's ore"),
	"kill darknut":                    Or("sword", "bombs", "beams", "scent seeds", "fool's ore", "punch"),
	"kill darknut (pit)":              Or("sword", "bombs", "beams", "scent seeds", "fool's ore", "punch", "shield", "rod", Hard("shovel")),
	"kill syger":                      Or("sword", "bombs", "scent seeds", "fool's ore", "punch"),
	"kill digdogger":                  Or("magnet gloves"),
	"break crystal":                   Or("sword", "bombs", "punch", "bracelet"),
	"kill hardhat (magnet)":           Or("magnet gloves", "gale seed weapon"),
	"kill vire":                       Or("sword", "bombs", "fool's ore", "punch"),
	"finish manhandla":                Or("sword", "bombs", "slingshot", "fool's ore"),
	"kill manhandla":                  And("boomerang L-2", "finish manhandla"),
	"kill wizzrobe":                   Or("kill normal"),
	"kill keese":                      Or("kill normal", "boomerang"),
	"kill magunesu":                   Or("sword", "fool's ore", "punch"), // even bombs don't work!
	"kill poe sister":                 Or("sword", "beams", "ember seeds", "scent seeds", "bombs", "fool's ore", "punch"),
	"kill darknut (across pit)": Or(
		Or("beams", "toss bombs", "scent slingshot", "magnet gloves"),
		And("feather L-2", "kill darknut (pit)")),
	"kill wizzrobe (pit)":   Or("pit kill normal"),
	"kill stalfos (pit)":    Or("kill stalfos", "pit kill normal"),
	"kill gleeok":           Or("sword", "beams", "bombs", "fool's ore", "punch"),
	"hit switch":            Or("sword", "beams", "boomerang", "rod", "satchel", "slingshot", "bombs", "fool's ore", "punch", "shovel"),
	"kill frypolar":         And("mystery seeds", Or("bracelet", "ember seeds")),
	"kill pols voice (pit)": Or("sword", "beams", "boomerang", "rod", "scent seeds", "gale seed weapon", "bombs", "shield", "shovel", "fool's ore", "punch", "flute"),
	"kill medusa head":      Or("sword", "fool's ore"),
	"kill floormaster":      Or("kill normal"),
	"kill onox":             And("sword", "jump"), // probably, idc
}
