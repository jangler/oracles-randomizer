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

	// not actually placed
	"fist ring":      Root(),
	"expert's ring":  Root(),
	"toss ring":      Root(),
	"energy ring":    Root(),
	"light ring L-1": Root(),
	"light ring L-2": Root(),

	"sword L-1":     Or("sword 1", "sword 2"),
	"sword L-2":     And("sword 1", "sword 2"),
	"boomerang L-1": Or("boomerang 1", "boomerang 2"),
	"boomerang L-2": And("boomerang 1", "boomerang 2"),
	"slingshot L-1": Or("slingshot 1", "slingshot 2"),
	"slingshot L-2": And("slingshot 1", "slingshot 2"),
	"feather L-1":   Or("feather 1", "feather 2"),
	"feather L-2":   And("feather 1", "feather 2"),
	"satchel":       Or("satchel 1", "satchel 2"),

	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "harvest tree"), Hard("d5 armos chest"),
		HardAnd("harvest bush", Or("enter agunima", "enter d7")))),
	"harvest mystery seeds": And("seed item", Or(
		And("mystery tree seeds", "harvest tree"),
		HardAnd("d8 armos chest", "harvest bush"))),
	"harvest scent seeds": And("scent tree seeds",
		"seed item", "harvest tree"),
	"harvest pegasus seeds": And("seed item", Or(
		And("pegasus tree seeds", "harvest tree"),
		HardAnd("beach", "shield", "ore chunks", "seed item"))), // market
	"harvest gale seeds": And("gale tree seeds",
		"seed item", "harvest tree"),

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
	"any slingshot L-2": And("slingshot L-2", "any slingshot"),

	"ember seeds":   And("harvest ember seeds", "seed item"),
	"mystery seeds": And("harvest mystery seeds", "seed item"),
	"scent seeds":   And("harvest scent seeds", "seed item"),
	"pegasus seeds": And("harvest pegasus seeds", "seed item"),
	"gale seeds":    And("harvest gale seeds", "seed item"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump 2", "bracelet"),

	"flute": Or("ricky's flute", "moosh's flute", "dimitri's flute"),

	"shield L-1": Or("shop shield L-1", And("beach", "ore chunks")),
	"shield L-2": And("shield L-1", "red ore", "blue ore"),

	"sword":     Or("sword L-1", "sword L-2"),
	"shield":    Or("shield L-1", "shield L-2"),
	"boomerang": Or("boomerang L-1", "boomerang L-2"),
	"slingshot": Or("slingshot L-1", "slingshot L-2"),
	"seed item": Or("satchel", "slingshot"),
	"kill for bombs": Or("sword", "ember seeds",
		Or("scent slingshot", Hard("scent seeds")), "fool's ore"),
	"bombs": Or(Hard("enter d2 B"),
		HardAnd("harvest bush", "d2 bracelet room"),
		And("bombs, 10", Or("remove pot", "shovel", "remove flower", "flute",
			And("kill for bombs", Or("suburbs", "fairy fountain",
				And("mount cucco", Or("spring",
					"sunken city default spring"))))))),

	// jump x pit tiles
	"jump 2":      Or("feather L-1", "feather L-2"),
	"jump 3":      Or(And("feather L-1", "pegasus satchel"), "feather L-2"),
	"bomb jump 2": Or("jump 3", HardAnd("jump 2", "bombs")),
	"bomb jump 3": Or("jump 4", HardAnd("jump 3", "bombs")),
	"jump 4":      And("feather L-2"),
	"bomb jump 4": Or("jump 6", HardAnd("jump 4", "bombs")),
	"jump 6":      And("feather L-2", "pegasus satchel"),

	"harvest tree": Or("sword", "rod", "fool's ore"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	// technically the player can always get ore chunks if they can make it to
	// subrosia, but shovel is the only way that isn't annoying.
	"ore chunks": Or("shovel", Hard("start")),
}
