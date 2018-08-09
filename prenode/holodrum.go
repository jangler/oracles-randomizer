package prenode

// overworld route logic

// portal parents are defined here since they're mostly overworld nodes

var portalPrenodes = map[string]*Prenode{
	"rosa portal in":         And("sokra stump", "remove bush"),
	"rosa portal out":        And("temple"),
	"rosa portal in wrapper": Or("rosa portal in"), // hack for safety.go
	"rosa portal":            Or("rosa portal in wrapper", "rosa portal out"),

	"open floodgate": Or(
		And("pegasus tree", "hit lever", "floodgate key", "pegasus satchel", "bracelet"),
		And("pegasus tree", "hit lever", "floodgate key", "feather L-2", "bracelet"),
		And("floodgate key", "hit lever", "flippers", "bracelet")),
	"swamp portal": Or("beach",
		And("horon village", "remove bush", "flippers", "bracelet"),
		And("open floodgate", "long jump", "bracelet"),
		And("open floodgate", "animal flute", "bracelet")),

	// jump added since it's effectively useless otherwise
	"mountain portal": And("jump", Or("mount cucco", "hide and seek")),

	"lake portal": Or("furnace",
		And("eyeglass lake", "flippers"),
		And("eyeglass lake", "pegasus jump L-2")),

	"village portal": Or(
		And("horon village", "boomerang L-2"),
		And("horon village", "pegasus jump L-2"),
		And("pirate house", "hit lever")),

	"desert portal": And("samasa desert"), // one-way

	// effectively one-way
	"remains portal": Or(
		And("temple remains", "shovel", "remove bush", "pegasus jump L-2"),
		And("temple remains", "spring", "remove flower", "remove bush", "pegasus jump L-2", "winter"),
		And("temple remains", "summer", "remove bush", "pegasus jump L-2", "winter"),
		And("temple remains", "autumn", "remove bush", "jump", "winter")),

	// dead end
	"d8 portal": Or(
		And("remains portal", "bombs", "summer", "long jump", "magnet gloves"),
		And("remains portal", "bombs", "summer", "pegasus jump L-2")),
}

var holodrumPrenodes = map[string]*Prenode{
	// start->d1
	"horon village": Or(
		And("north horon stump", "remove bush"),
		And("ghastly stump", "remove bush", "flippers"),
		And("eastern coast", "ember seeds"),
		And("sokra stump", "ember seeds"),
		And("village portal", "hit lever"),
		And("swamp portal", "bracelet", "flippers", "remove bush")),
	"enter d0":      AndStep("horon village"),
	"maku key fall": AndSlot("horon village", "pop maku bubble"),
	"enter d1":      AndStep("horon village", "remove bush", "gnarled key"),

	// d1->d2
	"ember tree": AndSlot("horon village"),
	"sokra stump": Or(
		And("horon village", "ember seeds"),
		And("rosa portal", "remove bush"),
		And("post-d2 stump", "winter"),
		And("post-d2 stump", "cross water gap")),
	"post-d2 stump": Or("sunken city", "mystery tree",
		And("sokra stump", "winter"),
		And("sokra stump", "cross water gap")),
	"shovel gift": AndSlot("post-d2 stump", "winter"),
	"mystery tree": OrSlot("sunken city",
		And("post-d2 stump", "winter", "shovel"),
		And("post-d2 stump", "jump"),
		And("sokra stump", "cross water gap")),
	"enter d2 A": And("mystery tree", "remove bush"),
	"enter d2 B": And("mystery tree", "bracelet", "remove bush"),
	"enter d2 C": And("mystery tree", "bracelet", "remove bush"),
	"enter d2":   OrStep("enter d2 A", "enter d2 B", "enter d2 C"),

	// d2->d3
	"north horon stump": And("horon village", "remove bush"),
	"scent tree": OrSlot(
		And("north horon stump", "bracelet"),
		And("natzu", "animal flute"),
		And("north horon stump", "flippers"),
		And("goron mountain", "flippers")),
	"blaino":      And("scent tree"),
	"blaino gift": AndSlot("blaino", "rupees"),
	"ricky pen": Or("scent tree",
		And("ghastly stump", "jump"),
		And("pegasus tree", "jump")),
	"ghastly stump": Or("pegasus tree",
		And("horon village", "remove bush", "flippers"),
		And("ricky pen", "ricky"),
		And("ricky pen", "jump"),
		And("swamp portal", "bracelet", "remove bush")),
	"pegasus tree": OrSlot(
		And("ghastly stump", "ricky"),
		And("ghastly stump", "feather L-2"),
		And("ghastly stump", "summer")),
	"floodgate key gift": AndSlot("pegasus tree", "hit lever"),
	"spool swamp": Or("open floodgate",
		And("ghastly stump", "remove bush", "flippers"),
		And("scent tree", "flippers")),
	"square jewel chest": OrSlot(
		And("open floodgate", "winter", "animal flute"),
		And("open floodgate", "winter", "long jump", "bombs"),
		And("open floodgate", "winter", "flippers", "bombs")),
	"enter d3": AndStep("open floodgate", "summer"),

	// d3->d4
	"natzu": Or(
		And("scent tree", "jump", "animal flute"),
		And("goron mountain", "flippers"),
		And("sunken city", "animal flute")),
	"sunken city": Or("post-d2 stump",
		And("natzu", "animal flute"),
		And("mount cucco", "flippers")),
	"sunken gale tree":      AndSlot("sunken city", "cross water gap"),
	"dimitri":               And("sunken gale tree", "bombs"),
	"master's plaque chest": AndSlot("sunken gale tree", "dimitri", "sword", "cross water gap"),
	"flippers gift":         AndSlot("sunken gale tree", "dimitri", "master's plaque"),
	"mount cucco": Or("mountain portal",
		And("sunken city", "flippers"),
		And("goron mountain", "shovel", "bracelet")),
	"banana harvest item": Or("sword", "fool's ore"),
	"spring banana cucco": And("mount cucco", "bracelet"),
	"spring banana tree":  AndSlot("spring banana cucco", "spring", "jump", "banana harvest item"),
	"moosh":               And("mount cucco", "spring banana"),
	"dragon key spot": Or(
		And("mount cucco", "moosh"),
		And("mount cucco", "feather L-2")),
	"mario cave":         And("mount cucco", "spring"),
	"dragon keyhole":     And("mario cave", "winter", "jump", "bracelet"),
	"enter d4":           AndStep("dragon key", "dragon keyhole", "summer", "cross water gap"),
	"pyramid jewel spot": AndSlot("mario cave", "flippers"),

	// goron mountain
	"goron mountain": Or(
		And("mount cucco", "bracelet", "shovel"),
		And("scent tree", "flippers"),
		And("temple remains", "pegasus jump L-2"),
		And("natzu", "animal flute", "flippers")),

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
	"graveyard": Or(
		And("pirate ship", "long jump"),
		And("pirate ship", "bombs", "jump", "summer")),
	"enter d7": AndStep("graveyard", "shovel"),

	// d7->d8
	"temple remains": Or(
		And("goron mountain", "pegasus jump L-2"),
		And("goron mountain", "flippers"),
		And("ricky pen", "long jump")),

	// d8->d9
	"maku seed": And("d1 essence", "d2 essence", "d3 essence", "d4 essence", "d5 essence", "d6 essence", "d7 essence", "d8 essence"),
	"enter d9":  AndStep("scent tree", "maku seed"),
}
