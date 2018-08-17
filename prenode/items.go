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

	"gale tree seeds": Or("gale tree seeds 1", "gale tree seeds 2"),
	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "harvest tree"),
		HardAnd("harvest bush", Or(
			"enter agunima", "d5 armos key chest", "enter d7")))),
	"harvest mystery seeds": And("seed item", Or(
		And("mystery tree seeds", "harvest tree"),
		HardAnd("enter frypolar", "harvest bush"))),
	"harvest scent seeds":   And("scent tree seeds", "seed item", "harvest tree"),
	"harvest pegasus seeds": And("pegasus tree seeds", "seed item", "harvest tree"),
	"harvest gale seeds":    And("gale tree seeds", "seed item", "harvest tree"),

	// has to be a different node from the slottable one
	"buy satchel": HardAnd("beach", "ore chunks", "big rupees"),

	"ember satchel":   And("harvest ember seeds", "satchel", Hard("buy satchel")),
	"mystery satchel": And("harvest mystery seeds", "satchel", Hard("buy satchel")),
	"scent satchel":   And("harvest scent seeds", "satchel", Hard("buy satchel")),
	"pegasus satchel": And("harvest pegasus seeds", "satchel", Hard("buy satchel")),
	"gale satchel":    And("harvest gale seeds", "satchel", Hard("buy satchel")),

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

	"punch":           And("punch ring", "medium rupees"),
	"use energy ring": And("energy ring", "medium rupees"),
	"use toss ring":   And("toss ring", "medium rupees"),
	"sword beams L-1": And("sword L-1", "use energy ring"),

	"pegasus jump L-1": And("pegasus satchel", "feather L-1"),
	"pegasus jump L-2": And("pegasus satchel", "feather L-2"),
	"long jump":        Or("feather L-2", "pegasus jump L-1"),
	"cross water gap":  Or("flippers", "jump"),
	"cross large pool": Or("flippers", "pegasus jump L-2"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump", "bracelet"),

	"strange flute": Or("big rupees", "temple"),
	"moosh flute":   And("big rupees", "south swamp", "kill moblin"),
	"dimitri flute": HardAnd("temple", "south swamp", "medium rupees"),
	"animal flute":  Or("ricky", "moosh flute", "dimitri flute"),
	"flute":         Or("strange flute", "animal flute"),

	"shield L-1": Or("medium rupees"),
	"shield L-2": And("shield L-1", "red ore", "blue ore"),

	"sword":     Or("sword L-1", "sword L-2"),
	"shield":    Or("shield L-1", "shield L-2"),
	"beams":     Or("sword L-2", "sword beams L-1"),
	"boomerang": Or("boomerang L-1", "boomerang L-2"),
	"slingshot": Or("slingshot L-1", "slingshot L-2"),
	"seed item": Or("satchel", "slingshot", Hard("buy satchel")),
	"bombs": Or("medium rupees", HardOr("bombs, 10",
		And("harvest bush", Or("d2 bracelet chest", "d2 spinner")))),
	"jump": Or("feather L-1", "feather L-2"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	"punch ring": Or("fist ring", "expert's ring"),

	// small rupees is ~10, and any item that can possibly yield rupees is
	// included.
	"small rupees": Or("sword", "boomerang", "shovel", "bracelet",
		"ember seeds", "scent seeds", "ricky", "moosh", Hard("dimitri flute"),
		"fool's ore", "punch"),
	// medium rupees is ~11-99, and only items that can reach rupee chests are
	// included.
	"medium rupees": And("small rupees", Or(Hard("small rupees"), "rupees, 20",
		"rupees, 30", "rupees, 50", "rupees, 100", "ember seeds",
		"d2 rupee room", "d6 rupee room")),
	// big rupees is ~100+, ember seeds are included since they can burn down
	// trees leading to generous old men.
	"big rupees": And("medium rupees", Or(Hard("medium rupees"), "ember seeds",
		"rupees, 100", "d2 rupee room", "d6 rupee room")),

	"ore chunks": Or("shovel", "temple"),
}
