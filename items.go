package main

// these are items that can be shuffled around in the route as OR nodes
//
// OR nodes without parents are false. if these are given parents (i.e. a
// chest/gift/whatever node), then they can become a valid part of the
// route
var baseItemNodes = []string{
	// ring box L-1 is free, but these nodes are "find" because it costs
	// rupees to appraise (and therefore use) rings
	"find energy ring",
	"find fist ring",
	"find expert's ring",
	"find toss ring",

	// shield, bombs, and flute can be bought
	"sword L-1",
	"gnarled key",
	"satchel",
	"boomerang L-1",
	// rod?
	"shovel",
	"bracelet",
	"ricky's gloves",
	"floodgate key",
	"square jewel",
	// member's card?
	"star-shaped ore",
	"ribbon",
	"feather L-1",
	"master's plaque",
	"flippers",
	// no fool's ore, see comment in subrosia.go
	"spring banana",
	"dragon key",
	"ring box L-2", // TODO where is this?
	"slingshot L-1",
	"pyramid jewel",
	"bomb flower",
	"magnet gloves",
	"x-shaped jewel",
	"round jewel",
	// sword L-2 is fixed
	"boomerang L-2",
	"rusty bell",
	"feather L-2",
	"slingshot L-2",
}

var itemNodesAnd = map[string]Point{
	"harvest ember seeds":   And{"ember tree", "harvest seeds"},
	"harvest mystery seeds": And{"mystery tree", "harvest seeds"},
	"harvest scent seeds":   And{"scent tree", "harvest seeds"},
	"harvest pegasus seeds": And{"pegasus tree", "harvest seeds"},
	"harvest gale seeds":    And{"gale tree", "harvest seeds"},

	"find d1 ember seeds":   And{"enter d1", "remove bush"},
	"find d2 ember seeds":   And{"mystery tree", "remove bush"},
	"find d2 mystery seeds": And{"d2 bomb wall", "remove bush"},
	"find d2 bombs":         And{"d2 bomb wall", "remove bush"},

	"ember satchel":   And{"ember seeds", "satchel"},
	"mystery satchel": And{"mystery seeds", "satchel"},
	"scent satchel":   And{"scent seeds", "satchel"},
	"pegasus satchel": And{"pegasus seeds", "satchel"},
	"gale satchel":    And{"gale seeds", "satchel"},

	"ember slingshot":   And{"ember seeds", "slingshot"},
	"mystery slingshot": And{"mystery seeds", "slingshot"},
	"scent slingshot":   And{"scent seeds", "slingshot"},
	"pegasus slingshot": And{"pegasus seeds", "slingshot"},
	"gale slingshot":    And{"gale seeds", "slingshot"},

	"punch":           And{"find punch ring", "rupees"},
	"energy ring":     And{"find energy ring", "rupees"},
	"sword beams L-1": And{"sword L-1", "energy ring"},

	"pegasus jump L-1": And{"pegasus satchel", "feather L-1"},
	"pegasus jump L-2": And{"pegasus satchel", "feather L-2"},

	"sword L-2": And{"lost woods", "winter", "autumn", "spring", "summer"},

	"winter": And{"rod", "winter tower"},
	"summer": And{"rod", "summer tower"},
	"spring": And{"rod", "spring tower"},
	"autumn": And{"rod", "autumn tower"},
}

var itemNodesOr = map[string]Point{
	"rod":        Or{"temple"}, // keep in place for now
	"animal":     Or{"ricky"},  // TODO there may be other ways to get one
	"fool's ore": Or{},         // disregard for now
	"shield L-1": Or{"rupees"},
	"shield L-2": Or{}, // TODO as if it matters

	"sword":      Or{"sword L-1", "sword L-2"},
	"shield":     Or{"shield L-1", "shield L-2"},
	"beams":      Or{"sword L-1", "sword beams L-1"},
	"boomerang":  Or{"boomerang L-1", "boomerang L-2"},
	"slingshot":  Or{"slingshot L-1", "slingshot L-2"},
	"seed item":  Or{"satchel", "slingshot"},
	"find bombs": Or{"find d2 bombs"},
	"bombs":      Or{"rupees", "find bombs"},
	"jump":       Or{"feather L-1", "feather L-2"},

	"harvest seeds":      Or{"sword", "rod", "fool's ore", "punch"},
	"find ember seeds":   Or{"find d1 ember seeds", "find d2 ember seeds"}, // TODO
	"ember seeds":        Or{"harvest ember seeds", "find ember seeds"},
	"find mystery seeds": Or{"find d2 mystery seeds"}, // TODO
	"mystery seeds":      Or{"harvest mystery seeds", "find mystery seeds"},
	"find scent seeds":   Or{}, // TODO
	"scent seeds":        Or{"harvest scent seeds", "find scent seeds"},
	"find pegasus seeds": Or{}, // TODO
	"pegasus seeds":      Or{"harvest pegasus seeds", "find pegasus seeds"},
	"find gale seeds":    Or{}, // TODO
	"gale seeds":         Or{"harvest gale seeds", "find gale seeds"},

	"long jump": Or{"feather L-2", "pegasus jump L-1"},

	"find punch ring": Or{"find fist ring", "find expert's ring"},

	// technically the rod can kill certain enemies for rupees, but you can't
	// access those enemies without another item that already collects rupees.
	// i'm also not including expendable items in this list just because it
	// could be super tedious to farm rupees using them.
	"rupees": Or{"sword", "boomerang L-2", "shovel", "bracelet", "animal", "fool's ore", "punch"},
}
