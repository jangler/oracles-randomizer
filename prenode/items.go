package prenode

// these are items that can be shuffled around in the route as root nodes
var baseItemPrenodes = map[string]*Prenode{
	// shield, bombs, and flute can be bought
	"sword L-1":       Root(),
	"gnarled key":     Root(),
	"satchel":         Root(),
	"boomerang L-1":   Root(),
	"rod":             Root(),
	"shovel":          Root(),
	"bracelet":        Root(),
	"ricky's gloves":  Root(),
	"floodgate key":   Root(),
	"star ore":        Root(),
	"feather L-1":     Root(),
	"flippers":        Root(),
	"fool's ore":      Root(),
	"slingshot L-1":   Root(),
	"magnet gloves":   Root(),
	"sword L-2":       Root(),
	"boomerang L-2":   Root(),
	"feather L-2":     Root(),
	"master's plaque": Root(),
	"spring banana":   Root(),
	"dragon key":      Root(),
	"slingshot L-2":   Root(),
	"square jewel":    Root(),
	"pyramid jewel":   Root(),
	"x-shaped jewel":  Root(),
	"round jewel":     Root(),
	"rusty bell":      Root(),

	// these can only be placed in seed tree slots
	"ember tree seeds":   Root(),
	"mystery tree seeds": Root(),
	"scent tree seeds":   Root(),
	"pegasus tree seeds": Root(),
	"gale tree seeds 1":  Root(),
	"gale tree seeds 2":  Root(),

	// could be uncommented and function as a filler item
	// "bombchus": Root(),
}

// don't slot these for now; they don't satisfy anything or cause problems
var ignoredBaseItemPrenodes = map[string]*Prenode{
	// "ring box L-2": Root(),

	"find fist ring":     Root(),
	"find expert's ring": Root(),
	"find energy ring":   Root(),
	"find toss ring":     Root(),
}

var itemPrenodes = map[string]*Prenode{
	"gale tree seeds": Or("gale tree seeds 1", "gale tree seeds 2"),
	"harvest ember seeds": And("seed item", Or(
		And("ember tree seeds", "harvest tree"),
		HardAnd("harvest bush", Or("enter agunima", "enter d7")))),
	"harvest mystery seeds": And("seed item", Or(
		And("mystery tree seeds", "harvest tree"),
		HardAnd("enter frypolar", "harvest bush"))),
	"harvest scent seeds":   And("scent tree seeds", "seed item", "harvest tree"),
	"harvest pegasus seeds": And("pegasus tree seeds", "seed item", "harvest tree"),
	"harvest gale seeds":    And("gale tree seeds", "seed item", "harvest tree"),

	// has to be a different node from the slottable one
	"buy satchel": HardAnd("beach", "ore chunks", "rupees"),

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

	"punch":           And("find punch ring", "rupees"),
	"energy ring":     And("find energy ring", "rupees"),
	"toss ring":       And("find toss ring", "rupees"),
	"sword beams L-1": And("sword L-1", "energy ring"),

	"pegasus jump L-1": And("pegasus satchel", "feather L-1"),
	"pegasus jump L-2": And("pegasus satchel", "feather L-2"),
	"long jump":        Or("feather L-2", "pegasus jump L-1"),
	"cross water gap":  Or("flippers", "jump"),
	"cross large pool": Or("flippers", "pegasus jump L-2"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump", "bracelet"),

	"winter": AndStep("rod", "winter tower"),
	"summer": AndStep("rod", "summer tower"),
	"spring": AndStep("rod", "spring tower"),
	"autumn": AndStep("rod", "autumn tower"),

	"strange flute": Or("rupees", "temple"),
	"moosh flute":   And("rupees", "south swamp", "kill moblin"),
	"dimitri flute": And("temple", "south swamp", "rupees"),
	"animal flute":  OrStep("ricky", "moosh flute", "dimitri flute"),
	"flute":         OrStep("strange flute", "animal flute"),

	"shield L-1": Or("rupees"),
	"shield L-2": Root(), // TODO as if it matters

	"sword":     Or("sword L-1", "sword L-2"),
	"shield":    Or("shield L-1", "shield L-2"),
	"beams":     Or("sword L-2", "sword beams L-1"),
	"boomerang": Or("boomerang L-1", "boomerang L-2"),
	"slingshot": Or("slingshot L-1", "slingshot L-2"),
	"seed item": Or("satchel", "slingshot", Hard("buy satchel")),
	"bombs": Or("rupees",
		HardAnd("harvest bush", Or("d2 bracelet chest", "d2 spinner"))),
	"jump": Or("feather L-1", "feather L-2"),

	"harvest tree": Or("sword", "rod", "fool's ore", "punch"),
	"harvest bush": Or("sword", "bombs", "fool's ore"),

	"find punch ring": Or("find fist ring", "find expert's ring"),

	// technically the rod can kill certain enemies for rupees, but you can't
	// access those enemies without another item that already collects rupees.
	// most expendable items aren't included, but ember seeds are because
	// they're free and can burn down trees leading to old men who give rupees.
	"rupees": OrStep("sword", Hard("boomerang L-1"), "boomerang L-2", "shovel",
		"bracelet", "ember seeds", "ricky", Hard("moosh"), "animal flute",
		"fool's ore", "punch"),

	"ore chunks": OrStep("shovel", "temple"),
}
