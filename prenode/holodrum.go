package prenode

// new set of holodrum prenodes, accounting for randomized seasons
//
// these use new nested prenode constructors, which are nice, but hard to read
// if you nest them in a way that requires multiple levels of indentation. so
// try not to do that too much.

var holodrumPrenodes = map[string]*Prenode{
	"start": And(), // parent for nodes reachable by default

	// horon village & western coast
	"horon village": And("start"),
	"x-shaped jewel chest": AndSlot("horon village", Or("ember slingshot", "mystery slingshot"),
		"mystery seeds", "kill moldorm"),
	"enter d0":      And("horon village"),
	"maku key fall": AndSlot("horon village", "pop maku bubble"),
	"ember tree":    AndSlot("horon village"),
	"pirate ship":   And("pirate's bell", "pirate house"),
	"graveyard": And("pirate ship", Or("long jump",
		"western coast default summer",
		And("bombs", "jump", "summer"))),
	"enter d7": And("graveyard", Or("shovel",
		Or("western coast default summer", "summer"))),

	// eastern suburbs
	"suburbs": Or( // this is the area south of the pool by sokra's stump
		And("horon village", "ember seeds"),
		And("rosa portal", "remove bush"),
		And("fairy fountain", Or("eastern suburbs default winter", "winter", "cross water gap"))),
	"fairy fountain": Or("sunken city",
		And("suburbs", Or("eastern suburbs default winter", "winter", "cross water gap"))),
	"shovel gift": AndSlot("fairy fountain", Or("eastern suburbs default winter", "winter"),
		Or("woods of winter default winter", "winter")),
	"mystery tree": OrSlot(
		And("fairy fountain", Or("eastern suburbs default winter", "winter", "shovel")),
		And("fairy fountain", Or(
			"eastern suburbs default spring", "spring",
			"eastern suburbs default summer", "summer",
			"eastern suburbs default autumn", "autumn"))),
	"enter d2 A": And("mystery tree", "remove bush"),
	"enter d2 B": And("mystery tree", "bracelet"),
	"enter d2 C": And("mystery tree", "bracelet"),
	"enter d2":   Or("enter d2 A", "enter d2 B", "enter d2 C"),

	// holodrum plain
	"ghastly stump": Or(
		"pegasus tree",
		And("scent tree", Or("jump", "ricky", "holodrum plain default winter")),
		And("south swamp", "flippers", "remove bush")),
	"scent tree": OrSlot(
		And("ghastly stump", Or("jump", "ricky", "winter", "holodrum plain default winter")),
		And("south swamp", "flippers"),
		And("sunken city", "animal flute"),
		And("north horon stump", "bracelet"),
		And("temple remains", "long jump"),
		And("goron mountain", "flippers")),
	"blaino gift": AndSlot("scent tree", "rupees"),
	"ricky":       And("scent tree", "ricky's gloves"),

	// spool swamp
	"pegasus tree": AndSlot("ghastly stump",
		Or("holodrum plain default summer", "summer", "feather L-2")),
	"floodgate key gift": AndSlot("pegasus tree", "hit lever"),
	"spool stump": And("pegasus tree", "hit lever",
		Or("pegasus satchel", "flippers", "feather L-2"), "bracelet", "floodgate key"),
	"dry swamp": Or(
		Or("spool swamp default summer", "spool swamp default autumn", "spool swamp default winter"),
		And("spool stump", Or("summer", "autumn", "winter"))),
	"south swamp": Or(
		And("spool stump", "flippers"),
		And("spool stump", "dry swamp", Or("long jump", "animal flute")),
		And("ghastly stump", "remove bush", "flippers"),
		And("scent tree", "flippers"),
		And("swamp portal", "bracelet")),
	"square jewel chest": AndSlot(Or("spool swamp default winter",
		And("spool stump", "winter")), "shovel", Or("animal flute", "bombs")),
	"enter d3": And("spool stump", Or("spool swamp default summer", "summer")),

	// north horon / eyeglass lake
	"not north horon default summer": Or(
		"north horon default spring", "north horon default autumn", "north horon default winter"),
	"north horon stump": Or(
		And("horon village", "remove bush"),
		And("scent tree", "bracelet"),
		And("south swamp", "flippers", "remove bush"),
		And("lake portal", "not north horon default summer", "flippers", "jump"),
		And("lake portal", "pegasus jump L-2", "north horon default winter")),
	"enter d1": And("gnarled key", Or(
		And("south swamp", "flippers"),
		And("north horon stump", "remove bush"))),
	"wet eyeglass lake": Or("not north horon default summer", "spring", "autumn", "winter"),
	"enter d5": And(Or("north horon default autumn", "autumn"), "remove mushroom", Or(
		And("lake portal", "not north horon default summer", "flippers"),
		And("north horon stump", "north horon default winter", Or("jump", "animal flute")))),

	// sunken city
	"sunken city": Or("fairy fountain",
		And("scent tree", Or("jump", "flippers"), "animal flute"),
		And("mount cucco", "flippers")),
	"sunken gale tree":      AndSlot("sunken city", "cross water gap"),
	"dimitri":               And("sunken gale tree", "bombs"),
	"master's plaque chest": AndSlot("dimitri", "sword", "cross water gap"),
	"flippers gift":         AndSlot("dimitri", "master's plaque"),

	// mount cucco
	"mount cucco": Or("mountain portal",
		And("sunken city", "flippers", Or("sunken city default summer", "summer")),
		And("goron mountain", "bracelet", "shovel")),
	"spring banana cucco": And("mount cucco", "bracelet"),
	"spring banana tree": AndSlot("spring banana cucco", "jump",
		Or("sunken city default spring", "spring"), Or("sword", "fool's ore")),
	"moosh":              And("mount cucco", "spring banana"),
	"dragon key spot":    AndSlot("mount cucco", Or("moosh", "feather L-2")),
	"mario cave":         And("mount cucco", Or("sunken city default spring", "spring")),
	"dragon keyhole":     And("mario cave", "winter", "jump", "bracelet"),
	"enter d4":           And("dragon key", "dragon keyhole", "summer", "cross water gap"),
	"pyramid jewel spot": AndSlot("mario cave", "flippers"),

	// goron mountain
	"goron mountain": Or(
		And("mount cucco", "shovel", "bracelet"),
		And("temple remains", "long jump", Or("flippers", "pegasus jump L-2")),
		And("scent tree", "flippers")),

	// tarm ruins
	"tarm ruins": And("pegasus tree",
		"square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	"lost woods": And("tarm ruins", Or("lost woods default summer", "summer"),
		Or("lost woods default winter", "winter"),
		Or("lost woods default autumn", "autumn"), "remove mushroom"),
	"noble sword spot": AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"tarm gale tree":   AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"enter d6": And("tarm gale tree", Or("tarm ruins default winter", "winter"),
		Or("tarm ruins default spring", "spring"), "shovel", "remove flower"),

	// samasa desert
	"desert":          And("suburbs", "pirate house"),
	"rusty bell spot": AndSlot("desert", "bracelet"),

	// temple remains (the important logic is in the portal prenodes)
	"temple remains": Or(
		And("goron mountain", Or("flippers", "pegasus jump L-2"), "long jump"),
		And("scent tree", "long jump")),

	// northern peak
	"maku seed": And("d1 essence", "d2 essence", "d3 essence", "d4 essence",
		"d5 essence", "d6 essence", "d7 essence", "d8 essence"),
	"enter d9": And("scent tree", "maku seed"),
}

var seasonPrenodes = map[string]*Prenode{
	"north horon default spring": Root(),
	"north horon default summer": Root(),
	"north horon default autumn": Root(),
	"north horon default winter": Root("start"),

	"eastern suburbs default spring": Root(),
	"eastern suburbs default summer": Root(),
	"eastern suburbs default autumn": Root("start"),
	"eastern suburbs default winter": Root(),

	"woods of winter default spring": Root(),
	"woods of winter default summer": Root("start"),
	"woods of winter default autumn": Root(),
	"woods of winter default winter": Root(),

	"spool swamp default spring": Root(),
	"spool swamp default summer": Root(),
	"spool swamp default autumn": Root("start"),
	"spool swamp default winter": Root(),

	"holodrum plain default spring": Root("start"),
	"holodrum plain default summer": Root(),
	"holodrum plain default autumn": Root(),
	"holodrum plain default winter": Root(),

	"sunken city default spring": Root(),
	"sunken city default summer": Root("start"),
	"sunken city default autumn": Root(),
	"sunken city default winter": Root(),

	"lost woods default spring": Root(),
	"lost woods default summer": Root(),
	"lost woods default autumn": Root("start"),
	"lost woods default winter": Root(),

	"tarm ruins default spring": Root("start"),
	"tarm ruins default summer": Root(),
	"tarm ruins default autumn": Root(),
	"tarm ruins default winter": Root(),

	"western coast default spring": Root(),
	"western coast default summer": Root(),
	"western coast default autumn": Root(),
	"western coast default winter": Root("start"),

	"temple remains default spring": Root(),
	"temple remains default summer": Root(),
	"temple remains default autumn": Root(),
	"temple remains default winter": Root("start"),
}
