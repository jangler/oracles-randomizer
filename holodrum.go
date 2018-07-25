package main

// overworld route logic

// portal parents are defined here since they're mostly overworld nodes
// see subrosia.go for the note about "remove stuck bush"

var portalPoints = map[string]Point{
	"rosa portal in":         And{"sokra stump", "remove bush"},
	"rosa portal out":        And{"temple"},
	"rosa portal in wrapper": Or{"rosa portal in"}, // hack for safety.go
	"rosa portal":            Or{"rosa portal in wrapper", "rosa portal out"},

	"open floodgate 1": And{"pegasus tree", "hit lever", "floodgate key", "pegasus satchel", "bracelet"},
	"open floodgate 2": And{"pegasus tree", "hit lever", "floodgate key", "feather L-2", "bracelet"},
	"open floodgate 3": And{"floodgate key", "hit lever", "flippers", "bracelet"},
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

	"desert portal": And{"samasa desert"}, // one-way

	// effectively one-way
	"remains portal 1": And{"temple remains", "shovel", "remove bush", "pegasus jump L-2"},
	"remains portal 2": And{"temple remains", "spring", "remove flower", "remove bush", "pegasus jump L-2", "winter"},
	"remains portal 3": And{"temple remains", "summer", "remove bush", "pegasus jump L-2", "winter"},
	"remains portal 4": And{"temple remains", "autumn", "remove bush", "jump", "winter"},

	// dead end
	"d8 portal 1": And{"remains portal", "summer", "long jump", "magnet gloves"},
	"d8 portal 2": And{"remains portal", "summer", "pegasus jump L-2"},

	// exiting subrosia via the rosa portal without having activated it from
	// holodrum gets you stuck in a bush unless you have a way to cut it down.
	// usable items are: sword (spin slash), bombs, gale seeds, slingshot w/
	// ember seeds but NOT satchel, bracelet, and L-2 boomerang. bombs and
	// seeds are not renewable and it's possible to reach this portal via, say,
	// the village portal with only satchel, pegasus seeds, and cape. this node
	// is used for checking for softlocks, but should not be a parent of any
	// other node.
	"remove stuck bush": Or{"sword", "boomerang L-2", "bracelet"},
}

var holodrumPoints = map[string]Point{
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
	"mystery tree 1":  And{"post-d2 stump", "winter", "shovel"},
	"mystery tree 2":  And{"post-d2 stump", "jump"},
	"mystery tree 3":  And{"sokra stump", "cross water gap"},
	"mystery tree 4":  And{"sunken city"},
	"enter d2 A":      And{"mystery tree", "remove bush"},
	"enter d2 B":      And{"mystery tree", "bracelet", "remove bush"},
	"enter d2 C":      And{"mystery tree", "bracelet", "remove bush"},

	// d2->d3
	"north horon stump":  And{"horon village", "remove bush"},
	"scent tree 1":       And{"north horon stump", "bracelet"},
	"scent tree 2":       And{"natzu", "animal"},
	"scent tree 3":       And{"north horon stump", "flippers"},
	"blaino":             And{"scent tree"},
	"blaino gift":        AndSlot{"blaino", "rupees"},
	"ricky pen 1":        And{"scent tree"},
	"ricky pen 2":        And{"ghastly stump", "jump"},
	"ricky pen 3":        And{"pegasus tree", "jump"},
	"ghastly stump 1":    And{"horon village", "remove bush", "flippers"},
	"ghastly stump 2":    And{"ricky pen", "animal"},
	"ghastly stump 3":    And{"ricky pen", "jump"},
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

	// d3->d4
	"natzu 1":               And{"scent tree", "jump", "animal"},
	"natzu 2":               And{"goron mountain", "flippers"},
	"natzu 3":               And{"sunken city", "animal"},
	"sunken city 1":         And{"natzu", "animal"},
	"sunken city 2":         And{"mount cucco", "flippers"},
	"gale tree":             And{"sunken city", "cross water gap"},
	"dimitri":               And{"gale tree", "bombs"},
	"master's plaque chest": AndSlot{"gale tree", "dimitri", "sword", "cross water gap"},
	"flippers gift":         AndSlot{"gale tree", "dimitri", "master's plaque"},
	"mount cucco 1":         And{"sunken city", "flippers"},
	"mount cucco 2":         And{"goron mountain", "shovel", "bracelet"},
	"mount cucco 3":         And{"mountain portal"},
	"banana harvest item":   Or{"sword", "fool's ore"},
	"spring banana tree":    AndSlot{"mount cucco", "spring", "bracelet", "jump", "banana harvest item"},
	"moosh":                 And{"mount cucco", "spring banana"},
	"dragon key cross 1":    And{"mount cucco", "moosh"},
	"dragon key cross 2":    And{"mount cucco", "pegasus jump L-2"},
	"dragon key spot":       AndSlot{"dragon key cross"}, // wraps generated node
	"mario cave":            And{"mount cucco", "spring"},
	"dragon keyhole":        And{"mario cave", "winter", "jump", "bracelet"},
	"enter d4":              And{"dragon key", "dragon keyhole", "summer", "cross water gap"},
	"pyramid jewel spot":    AndSlot{"mario cave", "flippers"},

	// goron mountain
	"goron mountain 1": And{"mount cucco", "bracelet", "shovel"},
	"goron mountain 2": And{"temple remains", "flippers"},
	"goron mountain 3": And{"temple remains", "pegasus jump L-2"},
	"goron mountain 4": And{"natzu", "animal", "flippers"},

	// d4->d5 TODO
	"eyeglass lake": And{"north horon stump", "jump"},
	"enter d5":      And{"eyeglass lake", "autumn", "remove mushroom"},

	// d5->d6 TODO
	"x-shaped jewel chest": AndSlot{"horon village", "mystery slingshot", "kill moldorm"},

	// d6->d7 TODO
	"eastern coast": And{"horon village", "ember seeds"},
	"samasa desert": And{"pirate house", "eastern coast"},

	// d7->d8 TODO
	"temple remains 1": And{"goron mountain", "pegasus jump L-2"},
	"temple remains 2": And{"goron mountain", "flippers"},
	"temple remains 3": And{"ricky pen", "long jump"},

	// referenced things that i don't want to deal with yet
	"lost woods": Or{},
}
