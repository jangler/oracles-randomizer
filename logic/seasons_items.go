package logic

// these are *extra* items that can be shuffled around in the route as root
// nodes, in addition to the ones automatically added from checking default
// item slot contents.
var seasonsBaseItemNodes = map[string]*Node{
	"fool's ore": Root(),

	// could be uncommented and function as a filler item
	// "bombchus": Root(),
}

var seasonsItemNodes = map[string]*Node{
	"rod": Or("winter", "summer", "spring", "autumn"),

	"ricky's flute":   Root(),
	"dimitri's flute": Root(),
	"moosh's flute":   Root(),

	// TODO: this is a dumb placeholder for until new fill is more developed
	"rupees": Or("rupees, 30", "rupees, 50", "rupees, 100", "ember seeds",
		And("hard", "shovel")),

	// expert's ring can do some things that fist ring can't, so this is for
	// the lowest common denominator.
	"punch object": Or("fist ring", "expert's ring"),
	"punch enemy":  Or(And("hard", "fist ring"), "expert's ring"),

	// progressives
	"noble sword":     Count(2, "sword"),
	"magic boomerang": Count(2, "boomerang"),
	"hyper slingshot": Count(2, "slingshot"),
	"cape":            Count(2, "feather"),

	// this of course doesn't apply to all trees, but trees won't have any
	// seeds attached to them unless they can be harvested. so it works out.
	"refill seeds": Or("harvest tree", "dimitri's flute", "dimitri",
		And("hard", "remove bush")),

	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "refill seeds"), And("hard", "d5 armos chest"),
		And("hard", "harvest bush", Or("enter agunima", "enter d7")))),
	"harvest mystery seeds": And("seed item", Or(
		And("mystery tree seeds", "refill seeds"),
		And("hard", "d8 armos chest", "harvest bush"))),
	"harvest scent seeds": And("scent tree seeds",
		"seed item", "refill seeds"),
	"harvest pegasus seeds": And("seed item", Or(
		And("pegasus tree seeds", "refill seeds"),
		And("hard", "beach", "shield", "ore chunks"))), // market
	"harvest gale seeds": And("gale tree seeds",
		"seed item", "refill seeds"),

	"ember satchel":   And("harvest ember seeds", "satchel"),
	"mystery satchel": And("harvest mystery seeds", "satchel"),
	"scent satchel":   And("harvest scent seeds", "satchel"),
	"pegasus satchel": And("harvest pegasus seeds", "satchel"),
	"gale satchel":    And("harvest gale seeds", "satchel"),
	"any satchel": Or("ember satchel", "mystery satchel", "scent satchel",
		"pegasus satchel", "gale satchel"),

	"ember slingshot":   And("harvest ember seeds", "slingshot"),
	"mystery slingshot": And("harvest mystery seeds", "slingshot"),
	"scent slingshot":   And("harvest scent seeds", "slingshot"),
	"pegasus slingshot": And("harvest pegasus seeds", "slingshot"),
	"gale slingshot":    And("harvest gale seeds", "slingshot"),
	"any slingshot": Or("ember slingshot", "mystery slingshot",
		"scent slingshot", "pegasus slingshot", "gale slingshot"),
	"any hyper slingshot": And("hyper slingshot", "any slingshot"),

	"ember seeds":   And("harvest ember seeds", "seed item"),
	"mystery seeds": And("harvest mystery seeds", "seed item"),
	"scent seeds":   And("harvest scent seeds", "seed item"),
	"pegasus seeds": And("harvest pegasus seeds", "seed item"),
	"gale seeds":    And("harvest gale seeds", "seed item"),

	"bomb flower": And("furnace", "jump 2", "bracelet"),

	"flute": Or("ricky's flute", "moosh's flute", "dimitri's flute"),

	"shield": Or("wooden shield", "iron shield",
		And("beach", "ember seeds")),
	"seed item": Or("satchel", "slingshot"),
	"bombs": Or(
		And("hard", "d2 blade chest", "bracelet"), // deku scrub, TODO: rupees
		And("hard", "harvest bush", "d2 bracelet room"),
		And("bombs, 10", Or("shovel", "remove flower", "flute"))),

	// jump x pit tiles
	"jump 2":      And("feather"),
	"jump 3":      Or(And("feather", "pegasus satchel"), "cape"),
	"bomb jump 2": Or("jump 3", And("hard", "jump 2", "bombs")),
	"bomb jump 3": Or("jump 4", And("hard", "jump 3", "bombs")),
	"jump 4":      And("cape"),
	"bomb jump 4": Or("jump 6", And("hard", "jump 4", "bombs")),
	"jump 6":      And("cape", "pegasus satchel"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch object"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	// technically the player can always get ore chunks if they can make it to
	// subrosia, but shovel is the only way that isn't annoying.
	"ore chunks": Or("shovel", "hard"),
}
