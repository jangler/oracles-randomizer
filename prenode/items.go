package prenode

// these are *extra* items that can be shuffled around in the route as root
// nodes, in addition to the ones automatically added from checking default
// item slot contents.
var baseItemPrenodes = map[string]*Prenode{
	"fool's ore": Root(),

	// could be uncommented and function as a filler item
	// "bombchus": Root(),
}

var itemPrenodes = map[string]*Prenode{
	"rod": Or("winter", "summer", "spring", "autumn"),

	"ricky's flute":   Root(),
	"dimitri's flute": Root(),
	"moosh's flute":   Root(),

	"gale tree seeds": Or("gale tree seeds 1", "gale tree seeds 2"),
	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "harvest tree"),
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

	"ember slingshot":   And("harvest ember seeds", "slingshot"),
	"mystery slingshot": And("harvest mystery seeds", "slingshot"),
	"scent slingshot":   And("harvest scent seeds", "slingshot"),
	"pegasus slingshot": And("harvest pegasus seeds", "slingshot"),
	"gale slingshot":    And("harvest gale seeds", "slingshot"),

	"ember seeds":   And("harvest ember seeds", "seed item"),
	"mystery seeds": And("harvest mystery seeds", "seed item"),
	"scent seeds":   And("harvest scent seeds", "seed item"),
	"pegasus seeds": Or(
		And("harvest pegasus seeds", "seed item"),
		HardAnd("beach", "shield", "ore chunks", "seed item")), // subrosian market
	"gale seeds": And("harvest gale seeds", "seed item"),

	"punch":           And("punch ring"),
	"use energy ring": And("energy ring"),
	"use toss ring":   And("toss ring"),
	"sword beams L-1": And("sword L-1", "use energy ring"),

	"pegasus jump L-1": And("pegasus satchel", "feather L-1"),
	"pegasus jump L-2": And("pegasus satchel", "feather L-2"),
	"long jump":        Or("feather L-2", "pegasus jump L-1"),
	"cross water gap":  Or("flippers", "jump"),
	"cross large pool": Or("flippers", "pegasus jump L-2"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump", "bracelet"),

	"flute": Or("ricky's flute", "moosh's flute", "dimitri's flute"),

	"shield L-1": Root(), // TODO
	"shield L-2": And("shield L-1", "red ore", "blue ore"),

	"sword":     Or("sword L-1", "sword L-2"),
	"shield":    Or("shield L-1", "shield L-2"),
	"beams":     Or("sword L-2", "sword beams L-1"),
	"boomerang": Or("boomerang L-1", "boomerang L-2"),
	"slingshot": Or("slingshot L-1", "slingshot L-2"),
	"seed item": Or("satchel", "slingshot"),
	"buy bombs": Root(), // TODO
	"kill for bombs": Or("sword", "ember seeds", "scent seeds",
		Hard("mystery seeds"), "fool's ore", "punch"),
	"bombs": Or("buy bombs",
		And("harvest bush", Or("d2 bracelet room", "d2 spinner room")),
		And("bombs, 10", Or("remove pot", "shovel", And("kill for bombs",
			Or("suburbs", "fairy fountain", And("mount cucco",
				Or("spring", "sunken city default spring"))))))),
	"jump": Or("feather L-1", "feather L-2"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	"punch ring": Or("fist ring", "expert's ring"),

	"ore chunks": Or("shovel", "temple"),
}
