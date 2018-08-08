package prenode

// overworld route logic

// portal parents are defined here since they're mostly overworld nodes
// see subrosia.go for the note about "remove stuck bush"

var portalPrenodes = map[string]*Prenode{
	"rosa portal in":         And("sokra stump", "remove bush"),
	"rosa portal out":        And("temple"),
	"rosa portal in wrapper": Or("rosa portal in"), // hack for safety.go
	"rosa portal":            Or("rosa portal in wrapper", "rosa portal out"),

	"open floodgate 1": And("pegasus tree", "hit lever", "floodgate key", "pegasus satchel", "bracelet"),
	"open floodgate 2": And("pegasus tree", "hit lever", "floodgate key", "feather L-2", "bracelet"),
	"open floodgate 3": And("floodgate key", "hit lever", "flippers", "bracelet"),
	"swamp portal 1":   And("horon village", "remove bush", "flippers", "bracelet"),
	"swamp portal 2":   And("open floodgate", "long jump", "bracelet"),
	"swamp portal 3":   And("open floodgate", "animal flute", "bracelet"),
	"swamp portal 4":   And("beach"),

	// jump added since it's effectively useless otherwise
	"mountain portal 1": And("mount cucco", "jump"),
	"mountain portal 2": And("hide and seek", "jump"),

	"lake portal 1": And("eyeglass lake", "flippers"),
	"lake portal 2": And("eyeglass lake", "pegasus jump L-2"),
	"lake portal 3": And("furnace"),

	"village portal 1": And("horon village", "boomerang L-2"),
	"village portal 2": And("horon village", "pegasus jump L-2"),
	"village portal 3": And("pirate house", "hit lever"),

	"desert portal": And("samasa desert"), // one-way

	// effectively one-way
	"remains portal 1": And("temple remains", "shovel", "remove bush", "pegasus jump L-2"),
	"remains portal 2": And("temple remains", "spring", "remove flower", "remove bush", "pegasus jump L-2", "winter"),
	"remains portal 3": And("temple remains", "summer", "remove bush", "pegasus jump L-2", "winter"),
	"remains portal 4": And("temple remains", "autumn", "remove bush", "jump", "winter"),

	// dead end
	"d8 portal 1": And("remains portal", "bombs", "summer", "long jump", "magnet gloves"),
	"d8 portal 2": And("remains portal", "bombs", "summer", "pegasus jump L-2"),

	// this is strictly for softlock checking; see safety.go
	"remove stuck bush": Or("sword", "boomerang L-2", "bracelet"),
}

var holodrumPrenodes = map[string]*Prenode{
	// start->d1
	"horon village 1": And("north horon stump", "remove bush"),
	"horon village 2": And("ghastly stump", "remove bush", "flippers"),
	"horon village 3": And("eastern coast", "ember seeds"),
	"horon village 4": And("sokra stump", "ember seeds"),
	"horon village 5": And("village portal", "hit lever"),
	"horon village 6": And("swamp portal", "bracelet", "flippers", "remove bush"),
	"enter d0":        AndStep("horon village"),
	"maku key fall":   AndSlot("horon village", "pop maku bubble"),
	"enter d1":        AndStep("horon village", "remove bush", "gnarled key"),

	// d1->d2
	"ember tree":      AndSlot("horon village"),
	"sokra stump 1":   And("horon village", "ember seeds"),
	"sokra stump 2":   And("rosa portal", "remove bush"),
	"sokra stump 3":   And("post-d2 stump", "winter"),
	"sokra stump 4":   And("post-d2 stump", "cross water gap"),
	"post-d2 stump 1": And("sokra stump", "winter"),
	"post-d2 stump 2": And("sokra stump", "cross water gap"),
	"post-d2 stump 3": And("sunken city"),
	"post-d2 stump 4": And("mystery tree"),
	"shovel gift":     AndSlot("post-d2 stump", "winter"),
	"mystery tree A":  And("post-d2 stump", "winter"),
	"mystery tree B":  And("post-d2 stump", "jump"),
	"mystery tree C":  And("sokra stump", "cross water gap"),
	"mystery tree D":  And("sunken city"),
	"mystery tree":    AndSlot("mystery tree A", "mystery tree B", "mystery tree C", "mystery tree D"),
	"enter d2 A":      And("mystery tree", "remove bush"),
	"enter d2 B":      And("mystery tree", "bracelet", "remove bush"),
	"enter d2 C":      And("mystery tree", "bracelet", "remove bush"),
	"enter d2":        OrStep("enter d2 A", "enter d2 B", "enter d2 C"),

	// d2->d3
	"north horon stump":    And("horon village", "remove bush"),
	"scent tree":           OrSlot("scent tree A", "scent tree B", "scent tree C"),
	"scent tree A":         And("north horon stump", "bracelet"),
	"scent tree B":         And("natzu", "animal flute"),
	"scent tree C":         And("north horon stump", "flippers"),
	"blaino":               And("scent tree"),
	"blaino gift":          AndSlot("blaino", "rupees"),
	"ricky pen 1":          And("scent tree"),
	"ricky pen 2":          And("ghastly stump", "jump"),
	"ricky pen 3":          And("pegasus tree", "jump"),
	"ghastly stump 1":      And("horon village", "remove bush", "flippers"),
	"ghastly stump 2":      And("ricky pen", "ricky"),
	"ghastly stump 3":      And("ricky pen", "jump"),
	"ghastly stump 4":      And("pegasus tree"),
	"ghastly stump 5":      And("swamp portal", "bracelet", "remove bush"),
	"pegasus tree":         OrSlot("pegasus tree A", "pegasus tree B", "pegasus tree C"),
	"pegasus tree A":       And("ghastly stump", "ricky"),
	"pegasus tree B":       And("ghastly stump", "feather L-2"),
	"pegasus tree C":       And("ghastly stump", "summer"),
	"floodgate key gift":   AndSlot("pegasus tree", "hit lever"),
	"spool swamp 1":        And("open floodgate"),
	"spool swamp 2":        And("ghastly stump", "remove bush", "flippers"),
	"spool swamp 3":        And("scent tree", "flippers"),
	"square jewel chest":   AndSlot("square jewel chest A", "square jewel chest B", "square jewel chest C"),
	"square jewel chest A": And("open floodgate", "winter", "animal flute"),
	"square jewel chest B": And("open floodgate", "winter", "long jump", "bombs"),
	"square jewel chest C": And("open floodgate", "winter", "flippers", "bombs"),
	"enter d3":             AndStep("open floodgate", "summer"),

	// d3->d4
	"natzu 1":               And("scent tree", "jump", "animal flute"),
	"natzu 2":               And("goron mountain", "flippers"),
	"natzu 3":               And("sunken city", "animal flute"),
	"sunken city 1":         And("natzu", "animal flute"),
	"sunken city 2":         And("mount cucco", "flippers"),
	"sunken city 3":         And("post-d2 stump"),
	"sunken gale tree":      AndSlot("sunken city", "cross water gap"),
	"dimitri":               And("sunken gale tree", "bombs"),
	"master's plaque chest": AndSlot("sunken gale tree", "dimitri", "sword", "cross water gap"),
	"flippers gift":         AndSlot("sunken gale tree", "dimitri", "master's plaque"),
	"mount cucco 1":         And("sunken city", "flippers"),
	"mount cucco 2":         And("goron mountain", "shovel", "bracelet"),
	"mount cucco 3":         And("mountain portal"),
	"banana harvest item":   Or("sword", "fool's ore"),
	"spring banana cucco":   And("mount cucco", "bracelet"),
	"spring banana tree":    AndSlot("spring banana cucco", "spring", "jump", "banana harvest item"),
	"moosh":                 And("mount cucco", "spring banana"),
	"dragon key cross 1":    And("mount cucco", "moosh"),
	"dragon key cross 2":    And("mount cucco", "pegasus jump L-2"),
	"dragon key spot":       AndSlot("dragon key cross"), // wraps generated node
	"mario cave":            And("mount cucco", "spring"),
	"dragon keyhole":        And("mario cave", "winter", "jump", "bracelet"),
	"enter d4":              AndStep("dragon key", "dragon keyhole", "summer", "cross water gap"),
	"pyramid jewel spot":    AndSlot("mario cave", "flippers"),

	// goron mountain
	"goron mountain 1": And("mount cucco", "bracelet", "shovel"),
	"goron mountain 2": And("temple remains", "flippers"),
	"goron mountain 3": And("temple remains", "pegasus jump L-2"),
	"goron mountain 4": And("natzu", "animal flute", "flippers"),

	// d4->d5
	"eyeglass lake": And("north horon stump", "jump"),
	"enter d5":      AndStep("eyeglass lake", "autumn", "remove mushroom"),

	// d5->d6; i'm treating tarm ruins like it's one way (like it normally is)
	"x-shaped jewel chest": AndSlot("horon village", "mystery slingshot", "kill moldorm"),
	"round jewel gift":     AndSlot("spool swamp", "flippers"),
	"tarm ruins":           And("pegasus tree", "square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	"lost woods":           And("tarm ruins", "summer", "winter", "autumn", "bracelet"),
	"noble sword spot":     AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"tarm gale tree":       AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"enter d6":             AndStep("tarm gale tree", "winter", "shovel", "spring", "remove flower"),

	// d6->d7
	"eastern coast":   And("horon village", "ember seeds"),
	"samasa desert":   And("pirate house", "eastern coast"),
	"rusty bell spot": AndSlot("samasa desert", "bracelet"),
	"pirate ship":     And("pirate's bell"),
	"graveyard 1":     And("pirate ship", "long jump"),
	"graveyard 2":     And("pirate ship", "bombs", "jump", "summer"),
	"enter d7":        OrStep("enter d7 A", "enter d7 B", "enter d7 C", "enter d7 D"),
	"enter d7 A":      And("graveyard", "shovel"),
	"enter d7 B":      And("graveyard", "spring"),
	"enter d7 C":      And("graveyard", "summer"),
	"enter d7 D":      And("graveyard", "autumn"),

	// d7->d8
	"temple remains 1": And("goron mountain", "pegasus jump L-2"),
	"temple remains 2": And("goron mountain", "flippers"),
	"temple remains 3": And("ricky pen", "long jump"),

	// d8->d9
	"maku seed": And("d1 essence", "d2 essence", "d3 essence", "d4 essence", "d5 essence", "d6 essence", "d7 essence", "d8 essence"),
	"enter d9":  AndStep("scent tree", "maku seed"),
}
