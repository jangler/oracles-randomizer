package logic

// these are *extra* items that can be shuffled around in the route as root
// nodes, in addition to the ones automatically added from checking default
// item slot contents.
var baseItemNodes = map[string]*Node{
	"fool's ore": Root(),

	// could be uncommented and function as a filler item
	// "bombchus": Root(),
}

var itemNodes = map[string]*Node{
	"rod": Or("winter", "summer", "spring", "autumn"),

	"ricky's flute":   Root(),
	"dimitri's flute": Root(),
	"moosh's flute":   Root(),

	"gale tree seeds": Or("gale tree seeds 1", "gale tree seeds 2"),
	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "harvest tree"), "d5 armos key chest",
		And("harvest bush", Or("enter agunima", "enter d7")))),
	"harvest mystery seeds": And("seed item", Or(
		And("mystery tree seeds", "harvest tree"),
		And("enter frypolar", "harvest bush"))),
	"harvest scent seeds":   And("scent tree seeds", "seed item", "harvest tree"),
	"harvest pegasus seeds": And("pegasus tree seeds", "seed item", "harvest tree"),
	"harvest gale seeds":    And("gale tree seeds", "seed item", "harvest tree"),

	"satchel": Or("satchel 1", "satchel 2"),

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
	"pegasus seeds": Or(
		And("harvest pegasus seeds", "seed item"),
		HardAnd("beach", "shield", "ore chunks", "seed item")), // subrosian market
	"gale seeds": And("harvest gale seeds", "seed item"),

	"punch": Or("fist ring", "expert's ring"),

	"pegasus jump L-1": And("pegasus satchel", "feather L-1"),
	"pegasus jump L-2": And("pegasus satchel", "feather L-2"),
	"long jump":        Or("feather L-2", "pegasus jump L-1"),
	"cross water gap":  Or("flippers", "jump"),
	"cross large pool": Or("flippers", "pegasus jump L-2"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump", "bracelet"),

	"flute": Or("ricky's flute", "moosh's flute", "dimitri's flute"),

	"shield L-1": Or("shop shield L-1", Hard("spool stump"),
		And("beach", "ore chunks")),
	"shield L-2": And("shield L-1", "red ore", "blue ore"),

	"sword":  Or("sword L-1", "sword L-2"),
	"shield": Or("shield L-1", "shield L-2"),
	"beams": Or("energy ring", And("sword L-2", Or(Hard("start"),
		"light ring L-1", "light ring L-2", "heart ring L-2"))),
	"boomerang": Or("boomerang L-1", "boomerang L-2"),
	"slingshot": Or("slingshot L-1", "slingshot L-2"),
	"seed item": Or("satchel", "slingshot"),
	"kill for bombs": Or("sword", "ember seeds", "scent seeds",
		Hard("mystery seeds"), "fool's ore", "punch"),
	"bombs": Or(Hard("enter d2 B"),
		And("harvest bush", Or("d2 bracelet room", "d2 spinner room")),
		And("bombs, 10", Or("remove pot", "shovel", "remove flover", "flute",
			And("kill for bombs", Or("suburbs", "fairy fountain",
				And("mount cucco", Or("spring",
					"sunken city default spring"))))))),
	"jump": Or("feather L-1", "feather L-2"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	"ore chunks": Or("shovel", "temple"),
}
