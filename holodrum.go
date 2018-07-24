package main

// overworld route logic

// portal parents are defined here since they're mostly overworld nodes
// see subrosia.go for the note about "remove stuck bush"

var portalNodesAnd = map[string]Point{
	"rosa portal in":  And{"sokra stump", "remove bush"},
	"rosa portal out": And{"temple"},

	"open floodgate 1": And{"pegasus tree", "floodgate key", "pegasus satchel", "bracelet"},
	"open floodgate 2": And{"pegasus tree", "floodgate key", "feather L-2", "bracelet"},
	"open floodgate 3": And{"floodgate key", "flippers", "bracelet"},
	"swamp portal 1":   And{"horon village", "remove bush", "flippers", "bracelet"},
	"swamp portal 2":   And{"open floodgate", "long jump", "bracelet"},
	"swamp portal 3":   And{"open floodgate", "animal", "bracelet"},
	"swamp portal 4":   And{"beach"},

	// jump added since it's effectively useless otherwise
	"mountain portal 1": And{"mount cucco", "jump"},
	"mountain portal 2": And{"hide and seek", "jump"},

	"lake portal 1": And{"eyeglass lake", "flippers"},
	"lake portal 2": And{"eyeglass lake", "pegasus jump L-2"},
	"lake portal 3": And{"furnace"},

	"village portal 1": And{"horon village", "boomerang L-2"},
	"village portal 2": And{"horon village", "pegasus jump L-2"},
	"village portal 3": And{"pirate house", "hit lever"},

	"desert portal": And{"samasa desert", "remove stuck bush"}, // one-way

	// effectively one-way
	"remains portal 1": And{"temple remains", "shovel", "remove bush", "pegasus jump L-2"},
	"remains portal 2": And{"temple remains", "spring", "remove flower", "remove bush", "pegasus jump L-2", "winter"},
	"remains portal 3": And{"temple remains", "summer", "remove bush", "pegasus jump L-2", "winter"},
	"remains portal 4": And{"temple remains", "autumn", "remove bush", "jump", "winter"},

	// dead end
	"d8 portal 1": And{"remains portal", "summer", "long jump", "magnet gloves"},
	"d8 portal 2": And{"remains portal", "summer", "pegasus jump L-2"},
}

var portalNodesOr = map[string]Point{
	"open floodgate": Or{"open floodgate 1", "open floodgate 2", "open floodgate 3"},

	// "unsafe" refers to the "remove stuck bush" issue
	"rosa portal in wrapper": Or{"rosa portal in"}, // dumb hack for safety checking; see safety.go
	"rosa portal":            Or{"rosa portal in wrapper", "rosa portal out"},
	"swamp portal":           Or{"swamp portal 1", "swamp portal 2", "swamp portal 3", "swamp portal 4"},
	"mountain portal":        Or{"mountain portal 1", "mountain portal 2"},
	"lake portal":            Or{"lake portal 1", "lake portal 2", "lake portal 3"},
	"village portal":         Or{"village portal 1", "village portal 2", "village portal 3"},
	"remains portal":         Or{"remains portal 1", "remains portal 2", "remains portal 3", "remains portal 4"},
	"d8 portal":              Or{"d8 portal 1", "d8 portal 2"},
}

var holodrumNodesAnd = map[string]Point{
	// start->d1
	"horon village": And{}, // start
	"enter d0":      And{"horon village"},
	"maku key fall": AndSlot{"horon village", "pop maku bubble"},
	"enter d1":      And{"horon village", "remove bush", "gnarled key"},

	// d1->d2
	"ember tree":      And{"horon village"},
	"sokra stump 1":   And{"horon village", "ember seeds"},
	"sokra stump 2":   And{"rosa portal", "remove bush"},
	"sokra stump 3":   And{"post-d2 stump", "winter"},
	"sokra stump 4":   And{"post-d2 stump", "cross water gap"},
	"post-d2 stump 1": And{"sokra stump", "winter"},
	"post-d2 stump 2": And{"sokra stump", "cross water gap"},
	"post-d2 stump 3": And{"sunken city"},
	"post-d2 stump 4": And{"mystery tree"},
	"shovel gift":     AndSlot{"post-d2 stump", "winter"},
	"mystery tree 1":  And{"shovel gift", "shovel", "winter"},
	"mystery tree 2":  And{"post-d2 stump", "winter", "shovel"},
	"mystery tree 3":  And{"post-d2 stump", "jump"},
	"mystery tree 4":  And{"sokra stump", "cross water gap"},
	"mystery tree 5":  And{"sunken city"},
	"enter d2 1":      And{"mystery tree", "remove bush"},
	"enter d2 2":      And{"mystery tree", "bracelet", "remove bush"},
	"enter d2 3":      And{"mystery tree", "bracelet", "remove bush"},

	// d2->d3
	"north horon stump":  And{"horon village", "remove bush"},
	"scent tree 1":       And{"north horon stump", "bracelet"},
	"scent tree 2":       And{"natzu", "animal"},
	"scent tree 3":       And{"natzu", "remove bush"}, // defaults to prairie if no animal
	"scent tree 4":       And{"north horon stump", "flippers"},
	"blaino":             And{"scent tree"},
	"blaino gift":        AndSlot{"blaino", "rupees"},
	"ricky 1":            And{"scent tree"},
	"ricky 2":            And{"ghastly stump", "jump"},
	"ricky 3":            And{"pegasus tree", "jump"},
	"ghastly stump 1":    And{"horon village", "remove bush", "flippers"},
	"ghastly stump 2":    And{"ricky", "animal"},
	"ghastly stump 3":    And{"ricky", "jump"},
	"ghastly stump 4":    And{"pegasus tree"},
	"ghastly stump 5":    And{"swamp portal", "bracelet", "remove bush"},
	"pegasus tree 1":     And{"ghastly stump", "animal"},
	"pegasus tree 2":     And{"ghastly stump", "feather L-2"},
	"pegasus tree 3":     And{"ghastly stump", "summer"},
	"floodgate key gift": AndSlot{"pegasus tree", "hit lever"},
	"square jewel 1":     And{"open floodgate", "winter", "animal"},
	"square jewel 2":     And{"open floodgate", "winter", "long jump", "bombs"},
	"square jewel 3":     And{"open floodgate", "winter", "flippers", "bombs"},
	"enter d3":           And{"open floodgate", "summer"},

	// d3->d4 TODO
	"natzu":       And{"scent tree", "jump"},
	"gale tree 1": And{"sunken city", "cross water gap"},
	"gale tree 2": And{"mount cucco", "flippers"},

	// d4->d5 TODO
	"eyeglass lake": And{"north horon stump", "jump"},
	"enter d5":      And{"eyeglass lake", "autumn", "remove mushroom"},

	// d5->d6 TODO
	"x-shaped jewel chest": AndSlot{"horon village", "mystery slingshot", "kill moldorm"},

	// d6->d7 TODO
	"eastern coast": And{"horon village", "ember seeds"},
	"samasa desert": And{"pirate house", "eastern coast"},

	// d7->d8 TODO
}

var holodrumNodesOr = map[string]Point{
	"cross water gap":    Or{"flippers", "jump"},
	"sokra stump":        Or{"sokra stump 1", "sokra stump 2", "sokra stump 3", "sokra stump 4"},
	"post-d2 stump":      Or{"post-d2 stump 1", "post-d2 stump 2", "post-d2 stump 3", "post-d2 stump 4"},
	"mystery tree":       Or{"mystery tree 1", "mystery tree 2", "mystery tree 3", "mystery tree 4", "mystery tree 5"},
	"scent tree":         Or{"scent tree 1", "scent tree 2", "scent tree 3", "scent tree 4"},
	"ricky":              Or{"ricky 1", "ricky 2", "ricky 3"},
	"ghastly stump":      Or{"ghastly stump 1", "ghastly stump 2", "ghastly stump 3", "ghastly stump 4", "ghastly stump 5"},
	"pegasus tree":       Or{"pegasus tree 1", "pegasus tree 2", "pegasus tree 3"},
	"square jewel chest": OrSlot{"square jewel 1", "square jewel 2", "square jewel 3"},
	"gale tree":          Or{"gale tree 1", "gale tree 2"},

	// referenced things that i don't want to deal with yet
	"sunken city":    Or{},
	"mount cucco":    Or{},
	"lost woods":     Or{},
	"temple remains": Or{},
}
