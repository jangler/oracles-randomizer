package prenode

// these are *extra* items that can be shuffled around in the route as root
// nodes, in addition to the ones automatically added from checking default
// item slot contents.
var baseItemPrenodes = map[string]*Prenode{
	"fool's ore": Root(),

	// could be uncommented and function as a filler item
	// "bombchus": Root(),

	// could fill the four unused ring slots
	/*
		"find fist ring":     Root(),
		"find expert's ring": Root(),
		"find energy ring":   Root(),
		"find toss ring":     Root(),
	*/
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

	"punch":           And("find punch ring", "medium rupees"),
	"energy ring":     And("find energy ring", "medium rupees"),
	"toss ring":       And("find toss ring", "medium rupees"),
	"sword beams L-1": And("sword L-1", "energy ring"),

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
	"bombs": Or("medium rupees",
		HardAnd("harvest bush", Or("d2 bracelet chest", "d2 spinner"))),
	"jump": Or("feather L-1", "feather L-2"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	"find punch ring": Or("find fist ring", "find expert's ring"),

	// small rupees is ~1-10, and any item that can possibly yield rupees is
	// included.
	"small rupees": Or("sword", "boomerang", "shovel", "bracelet",
		"ember seeds", "scent seeds", "ricky", "moosh", Hard("dimitri flute"),
		"fool's ore", "punch"),
	// medium rupees is ~11-99, and only items that can reach rupee chests are
	// included. TODO update this for other chests as they're added
	"medium rupees": And("small rupees", Or(
		Hard("small rupees"), "big rupees", "rupees, 20", "rupees, 30")),
	// big rupees is ~100+, and only ember seeds are included since they can
	// burn down trees leading to generous old men.
	"big rupees": Or(Hard("medium rupees"), "ember seeds"),

	"ore chunks": Or("shovel", "temple"),
}
