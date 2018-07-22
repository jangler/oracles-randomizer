package main

// overworld route logic

// portal parents are defined here since they're mostly overworld nodes
// see subrosia.go for the note about "remove stuck bush"

var portalNodesAnd = map[string][]string{
	"rosa portal 1": And{"eastern suburbs", "remove bush"},
	"rosa portal 2": And{"temple"},

	"swamp portal 1": And{"spool swamp"}, // TODO
	"swamp portal 2": And{"beach"},

	// jump added since it's effectively useless otherwise
	"mountain portal 1": And{"mount cucco", "jump"},
	"mountain portal 2": And{"hide and seek", "jump"},

	"lake portal 1": And{"eyeglass lake", "flippers"},
	"lake portal 2": And{"furnace"},

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
	"d8 portal 1": And{"remains portal", "summer", "long jump", "magnet glove"},
	"d8 portal 2": And{"remains portal", "summer", "pegasus jump L-2"},

	"swamp portal":    And{"swamp portal unsafe", "remove stuck bush"},
	"mountain portal": And{"mountain portal unsafe", "remove stuck bush"},
	"lake portal":     And{"lake portal unsafe", "remove stuck bush"},
	"village portal":  And{"village portal unsafe", "remove stuck bush"},
}

var portalNodesOr = map[string][]string{
	// "unsafe" refers to the "remove stuck bush" issue
	"rosa portal":            Or{"rosa portal 1", "rosa portal 2"},
	"swamp portal unsafe":    Or{"swamp portal 1", "swamp portal 2"},
	"mountain portal unsafe": Or{"mountain portal 1", "mountain portal 2"},
	"lake portal unsafe":     Or{"lake portal 1", "lake portal 2"},
	"village portal unsafe":  Or{"village portal 1", "village portal 2", "village portal 3", "village portal 4"},
	"remains portal":         Or{"remains portal 1", "remains portal 2", "remains portal 3", "remains portal 4"},
	"d8 portal":              Or{"d8 portal 1", "d8 portal 2"},
}

var horonVillageNodesAnd = map[string][]string{}

var horonVillageNodesOr = map[string][]string{
	"pop bubble": {"sword", "bombs", "ember seeds", "scent seeds", "normal slingshot", "pegasus slingshot"},

	// TODO
}
