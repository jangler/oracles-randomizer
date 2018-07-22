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

var itemNodesAnd = map[string][]string{
	"harvest ember seeds":   []string{"ember tree", "harvest seeds"},
	"harvest mystery seeds": []string{"mystery tree", "harvest seeds"},
	"harvest scent seeds":   []string{"scent tree", "harvest seeds"},
	"harvest pegasus seeds": []string{"pegasus tree", "harvest seeds"},
	"harvest gale seeds":    []string{"gale tree", "harvest seeds"},

	"find d1 ember seeds":   []string{"enter d1", "remove bush"},
	"find d2 ember seeds":   []string{"mystery tree", "remove bush"},
	"find d2 mystery seeds": []string{"d2 bomb wall", "remove bush"},
	"find d2 bombs":         []string{"d2 bomb wall", "remove bush"},

	"ember satchel":   []string{"ember seeds", "satchel"},
	"mystery satchel": []string{"mystery seeds", "satchel"},
	"scent satchel":   []string{"scent seeds", "satchel"},
	"pegasus satchel": []string{"pegasus seeds", "satchel"},
	"gale satchel":    []string{"gale seeds", "satchel"},

	"ember slingshot":   []string{"ember seeds", "slingshot"},
	"mystery slingshot": []string{"mystery seeds", "slingshot"},
	"scent slingshot":   []string{"scent seeds", "slingshot"},
	"pegasus slingshot": []string{"pegasus seeds", "slingshot"},
	"gale slingshot":    []string{"gale seeds", "slingshot"},

	"punch":           []string{"find punch ring", "rupees"},
	"energy ring":     []string{"find energy ring", "rupees"},
	"sword beams L-1": []string{"sword L-1", "energy ring"},

	"pegasus jump L-1": []string{"pegasus satchel", "feather L-1"},
	"pegasus jump L-2": []string{"pegasus satchel", "feather L-2"},

	"sword L-2": []string{"lost woods", "winter", "autumn", "spring", "summer"},

	"winter": []string{"rod", "winter tower"},
	"summer": []string{"rod", "summer tower"},
	"spring": []string{"rod", "spring tower"},
	"autumn": []string{"rod", "autumn tower"},
}

var itemNodesOr = map[string][]string{
	"rod":        Or{"temple"}, // keep in place for now
	"animal":     Or{"ricky"},  // TODO there may be other ways to get one
	"fool's ore": Or{},         // disregard for now
	"shield L-1": Or{"rupees"},
	"shield L-2": Or{}, // TODO as if it matters

	"sword":      []string{"sword L-1", "sword L-2"},
	"shield":     []string{"shield L-1", "shield L-2"},
	"beams":      []string{"sword L-1", "sword beams L-1"},
	"boomerang":  []string{"boomerang L-1", "boomerang L-2"},
	"slingshot":  []string{"slingshot L-1", "slingshot L-2"},
	"seed item":  []string{"satchel", "slingshot"},
	"find bombs": []string{"find d2 bombs"},
	"bombs":      []string{"rupees", "find bombs"},
	"jump":       []string{"feather L-1", "feather L-2"},

	"harvest seeds":      []string{"sword", "rod", "fool's ore", "punch"},
	"find ember seeds":   []string{"find d1 ember seeds", "find d2 ember seeds"}, // TODO
	"ember seeds":        []string{"harvest ember seeds", "find ember seeds"},
	"find mystery seeds": []string{"find d2 mystery seeds"}, // TODO
	"mystery seeds":      []string{"harvest mystery seeds", "find mystery seeds"},
	"find scent seeds":   []string{}, // TODO
	"scent seeds":        []string{"harvest scent seeds", "find scent seeds"},
	"find pegasus seeds": []string{}, // TODO
	"pegasus seeds":      []string{"harvest pegasus seeds", "find pegasus seeds"},
	"find gale seeds":    []string{}, // TODO
	"gale seeds":         []string{"harvest gale seeds", "find gale seeds"},

	"long jump": []string{"feather L-2", "pegasus jump L-1"},

	"find punch ring": []string{"find fist ring", "find expert's ring"},

	// technically the rod can kill certain enemies for rupees, but you can't
	// access those enemies without another item that already collects rupees.
	// i'm also not including expendable items in this list just because it
	// could be super tedious to farm rupees using them.
	"rupees": Or{"sword", "boomerang L-2", "shovel", "bracelet", "animal", "fool's ore", "punch"},
}
