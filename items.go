package main

// these are items that can be shuffled around in the route as OR nodes
//
// OR nodes without parents are false. if these are given parents (i.e. a
// chest/gift/whatever node), then they can become a valid part of the
// route
//
// TODO there are some optimizations that could be made here: the four jewels
//      are interchangeable, fist ring and expert's ring are interchangeable.
//      it would also be nice if the randomizer checked whether a L-1 item is
//      sufficient before trying a L-2 one.
var baseItemPoints = map[string]Point{
	// ring box L-1 is free, but these nodes are "find" because it costs
	// rupees to appraise (and therefore use) rings
	"find energy ring":   Or{},
	"find fist ring":     Or{},
	"find expert's ring": Or{},

	// shield, bombs, and flute can be bought
	"sword L-1":     Or{},
	"gnarled key":   Or{},
	"satchel":       Or{},
	"boomerang L-1": Or{},
	// rod?
	"shovel":         Or{},
	"bracelet":       Or{},
	"ricky's gloves": Or{},
	"floodgate key":  Or{},
	// member's card?
	"star ore":    Or{},
	"feather L-1": Or{},
	"flippers":    Or{},
	// no fool's ore, see comment in subrosia.go
	"slingshot L-1": Or{},
	"magnet gloves": Or{},
	// sword L-2 is fixed
	"boomerang L-2":   Or{},
	"feather L-2":     Or{},
	"master's plaque": Or{},
	"spring banana":   Or{},
	"dragon key":      Or{},
	"slingshot L-2":   Or{},
}

// don't slot these for now; they don't satisfy anything
var ignoredBaseItemPoints = map[string]Point{
	"ring box L-2":   Or{},
	"find toss ring": Or{},
	"square jewel":   Or{},
	"pyramid jewel":  Or{},
	"x-shaped jewel": Or{},
	"round jewel":    Or{},
	"rusty bell":     Or{},
}

var itemPoints = map[string]Point{
	"harvest ember seeds":   And{"ember tree", "satchel", "harvest item"},
	"harvest mystery seeds": And{"mystery tree", "satchel", "harvest item"},
	"harvest scent seeds":   And{"scent tree", "satchel", "harvest item"},
	"harvest pegasus seeds": And{"pegasus tree", "satchel", "harvest item"},
	"harvest gale seeds":    And{"gale tree", "satchel", "harvest item"},

	"find d1 ember seeds": And{"enter d1", "satchel", "remove bush"},
	"find d2 ember seeds": And{"mystery tree", "satchel", "remove bush"},
	"find d2 bombs":       And{"d2 bomb wall", "satchel", "remove bush"},

	// you can usually only get seed drops if you already have that type of
	// seed, meaning the first way to obtain it has to be harvesting from a
	// tree. the exception is ember seeds, since the satchel comes loaded with
	// them from the start.

	"ember satchel":   And{"get ember seeds", "satchel"},
	"mystery satchel": And{"harvest mystery seeds", "satchel"},
	"scent satchel":   And{"harvest scent seeds", "satchel"},
	"pegasus satchel": And{"harvest pegasus seeds", "satchel"},
	"gale satchel":    And{"harvest gale seeds", "satchel"},

	"ember slingshot":   And{"get ember seeds", "slingshot"},
	"mystery slingshot": And{"harvest mystery seeds", "slingshot"},
	"scent slingshot":   And{"harvest scent seeds", "slingshot"},
	"pegasus slingshot": And{"harvest pegasus seeds", "slingshot"},
	"gale slingshot":    And{"harvest gale seeds", "slingshot"},

	"ember seeds":   And{"get ember seeds", "seed item"},
	"mystery seeds": And{"harvest mystery seeds", "seed item"},
	"scent seeds":   And{"harvest scent seeds", "seed item"},
	"pegasus seeds": And{"harvest pegasus seeds", "seed item"},
	"gale seeds":    And{"harvest gale seeds", "seed item"},

	"punch":           And{"find punch ring", "rupees"},
	"energy ring":     And{"find energy ring", "rupees"},
	"sword beams L-1": And{"sword L-1", "energy ring"},

	"pegasus jump L-1": And{"pegasus satchel", "feather L-1"},
	"pegasus jump L-2": And{"pegasus satchel", "feather L-2"},
	"long jump":        Or{"feather L-2", "pegasus jump L-1"},
	"cross water gap":  Or{"flippers", "jump"},
	"cross large pool": Or{"flippers", "pegasus jump L-2"},

	"sword L-2": And{"lost woods", "winter", "autumn", "spring", "summer"},

	"ribbon":      And{"star ore", "beach"},
	"bomb flower": And{"furnace", "jump"},

	"winter": And{"rod", "winter tower"},
	"summer": And{"rod", "summer tower"},
	"spring": And{"rod", "spring tower"},
	"autumn": And{"rod", "autumn tower"},

	"animal": And{"ricky pen", "ricky's gloves"}, // TODO flute stuff

	"rod":        Or{"temple"}, // keep in place for now
	"fool's ore": Or{},         // disregard for now
	"shield L-1": Or{"rupees"},
	"shield L-2": Or{}, // TODO as if it matters

	"sword":          Or{"sword L-1", "sword L-2"},
	"shield":         Or{"shield L-1", "shield L-2"},
	"beams":          Or{"sword L-1", "sword beams L-1"},
	"boomerang":      Or{"boomerang L-1", "boomerang L-2"},
	"find slingshot": Or{"slingshot L-1", "slingshot L-2"},
	"slingshot":      And{"find slingshot", "satchel"}, // need satchel to use
	"seed item":      Or{"satchel", "slingshot"},
	"find bombs":     Or{"find d2 bombs"},
	"bombs":          Or{"rupees", "find bombs"},
	"jump":           Or{"feather L-1", "feather L-2"},

	"harvest item":     Or{"sword", "rod", "fool's ore", "punch"},
	"find ember seeds": Or{"find d1 ember seeds", "find d2 ember seeds"}, // TODO
	"get ember seeds":  Or{"harvest ember seeds", "find ember seeds"},

	"find punch ring": Or{"find fist ring", "find expert's ring"},

	// technically the rod can kill certain enemies for rupees, but you can't
	// access those enemies without another item that already collects rupees.
	// i'm also not including expendable items in this list just because it
	// could be super tedious to farm rupees using them.
	"rupees": Or{"sword", "boomerang", "shovel", "bracelet", "animal", "fool's ore", "punch"},
}
