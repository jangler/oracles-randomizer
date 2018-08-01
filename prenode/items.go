package prenode

// these are items that can be shuffled around in the route as OR nodes
//
// OR nodes without parents are false. if these are given parents (i.e. a
// chest/gift/whatever node), then they can become a valid part of the
// route
var baseItemPrenodes = map[string]*Prenode{
	// ring box L-1 is free, but these nodes are "find" because it costs
	// rupees to appraise (and therefore use) rings
	"find fist ring":     Root(),
	"find expert's ring": Root(),

	// shield, bombs, and flute can be bought
	"sword L-1":        Root(),
	"gnarled key":      Root(),
	"satchel":          Root(),
	"boomerang L-1":    Root(),
	"rod":              Root(),
	"shovel":           Root(),
	"bracelet":         Root(),
	"ricky's gloves":   Root(),
	"floodgate key":    Root(),
	"star ore":         Root(),
	"feather L-1":      Root(),
	"flippers":         Root(),
	"fool's ore":       Root(),
	"slingshot L-1":    Root(),
	"magnet gloves":    Root(),
	"sword L-2":        Root(),
	"boomerang L-2":    Root(),
	"feather L-2":      Root(),
	"master's plaque":  Root(),
	"spring banana":    Root(),
	"dragon key":       Root(),
	"slingshot L-2":    Root(),
	"square jewel":     Root(),
	"pyramid jewel":    Root(),
	"x-shaped jewel":   Root(),
	"round jewel":      Root(),
	"rusty bell":       Root(),
	"find energy ring": Root(),
	"find toss ring":   Root(),

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

// don't slot these for now; they don't satisfy anything
var ignoredBaseItemPrenodes = map[string]*Prenode{
	"ring box L-2": Root(),
}

var itemPrenodes = map[string]*Prenode{
	// XXX this will need to change once slingshot gives seeds
	"harvest ember seeds":   And("ember tree seeds", "satchel", "harvest item"),
	"harvest mystery seeds": And("mystery tree seeds", "satchel", "harvest item"),
	"harvest scent seeds":   And("scent tree seeds", "satchel", "harvest item"),
	"harvest pegasus seeds": And("pegasus tree seeds", "satchel", "harvest item"),
	"harvest gale seeds":    And("gale tree seeds", "satchel", "harvest item"),

	// you can usually only get seed drops if you already have that type of
	// seed, meaning the first way to obtain it has to be harvesting from a
	// tree. the exception is ember seeds, since the satchel comes loaded with
	// them from the start.

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
	"pegasus seeds": And("harvest pegasus seeds", "seed item"),
	"gale seeds":    And("harvest gale seeds", "seed item"),

	"punch":           And("find punch ring", "rupees"),
	"energy ring":     And("find energy ring", "rupees"),
	"toss ring":       And("find toss ring", "rupees"),
	"sword beams L-1": And("sword L-1", "energy ring"),

	"pegasus jump L-1": And("pegasus satchel", "feather L-1"),
	"pegasus jump L-2": And("pegasus satchel", "feather L-2"),
	"long jump":        Or("feather L-2", "pegasus jump L-1"),
	"cross water gap":  Or("flippers", "jump"),
	"cross large pool": Or("flippers", "pegasus jump L-2"),

	"sword L-2": And("lost woods", "winter", "autumn", "spring", "summer"),

	"ribbon":      And("star ore", "beach"),
	"bomb flower": And("furnace", "jump", "bracelet"),

	"winter": AndStep("rod", "winter tower"),
	"summer": AndStep("rod", "summer tower"),
	"spring": AndStep("rod", "spring tower"),
	"autumn": AndStep("rod", "autumn tower"),

	"ricky":         And("ricky pen", "ricky's gloves"),
	"strange flute": Or("rupees", "temple"),
	"moosh flute":   And("rupees", "spool swamp", "kill moblin"),
	"dimitri flute": And("temple", "spool swamp", "rupees"),
	"animal flute":  OrStep("ricky", "moosh flute", "dimitri flute"),
	"flute":         OrStep("strange flute", "animal flute"),

	"shield L-1": Or("rupees"),
	"shield L-2": Or(), // TODO as if it matters

	"sword":          Or("sword L-1", "sword L-2"),
	"shield":         Or("shield L-1", "shield L-2"),
	"beams":          Or("sword L-2", "sword beams L-1"),
	"boomerang":      Or("boomerang L-1", "boomerang L-2"),
	"find slingshot": Or("slingshot L-1", "slingshot L-2"),
	"slingshot":      And("find slingshot", "satchel"), // need satchel to use
	"seed item":      Or("satchel", "slingshot"),
	"bombs":          Or("rupees"),
	"jump":           Or("feather L-1", "feather L-2"),

	"harvest item": Or("sword", "rod", "fool's ore", "punch"),

	"find punch ring": Or("find fist ring", "find expert's ring"),

	// technically the rod can kill certain enemies for rupees, but you can't
	// access those enemies without another item that already collects rupees.
	// i'm also not including expendable items in this list just because it
	// could be super tedious to farm rupees using them.
	"rupees": OrStep("sword", "boomerang", "shovel", "bracelet", "ricky", "animal flute", "fool's ore", "punch"),
}
