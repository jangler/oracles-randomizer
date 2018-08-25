package prenode

// new set of holodrum prenodes, accounting for randomized seasons

var holodrumPrenodes = map[string]*Prenode{
	"start": And(), // parent for nodes reachable by default

	// horon village
	"horon village":    And("start"),
	"maku tree gift":   AndSlot("horon village", "pop maku bubble"),
	"ember tree":       AndSlot("horon village"),
	"village SE chest": AndSlot("horon village", "bombs"),
	"village SW chest": AndSlot("horon village", "remove bush", "remove mushroom"),
	"member's shop 1":  AndSlot("member's card", "big rupees"),
	"member's shop 2":  AndSlot("member's card", "big rupees"),
	"member's shop 3":  AndSlot("member's card", "big rupees"),

	// western coast
	"x-shaped jewel chest": AndSlot("horon village", Or("ember slingshot", "mystery slingshot"),
		"mystery seeds", "kill moldorm"),
	"enter d0":    And("horon village"),
	"pirate ship": And("pirate's bell", "pirate house"),
	"graveyard": And("pirate ship", Or("long jump",
		"western coast default summer",
		And("bombs", Or("jump", Hard("start")), "summer"))),
	"enter d7": And("graveyard", Or("shovel",
		Or("western coast default summer", "summer"),
		Or("western coast default spring", "spring"),
		Or("western coast default autumn", "autumn"))),
	"western coast chest": AndSlot("pirate ship"),
	"coast house chest":   AndSlot("pirate ship"),

	// eastern suburbs
	"suburbs": Or( // this is the area south of the pool by sokra's stump
		And("horon village", "ember seeds"),
		And("rosa portal", "remove bush"),
		And("fairy fountain", Or("eastern suburbs default winter", "winter", "cross water gap"))),
	"fairy fountain": Or("sunken city",
		And("suburbs", Or("eastern suburbs default winter", "winter", "cross water gap"))),
	"shovel gift": AndSlot("fairy fountain", Or("eastern suburbs default winter", "winter"),
		Or("woods of winter default winter", "winter")),
	"central woods of winter": Or(
		And("fairy fountain", Or("eastern suburbs default winter", "winter"),
			Or("shovel", "jump")),
		And("fairy fountain", Or(
			"eastern suburbs default spring", "spring",
			"eastern suburbs default summer", "summer",
			"eastern suburbs default autumn", "autumn"))),
	"mystery tree": OrSlot(
		And("fairy fountain", Or("eastern suburbs default winter", "winter"),
			Or("shovel", And("jump", "bracelet"))),
		And("fairy fountain", Or(
			"eastern suburbs default spring", "spring",
			"eastern suburbs default summer", "summer",
			"eastern suburbs default autumn", "autumn"))),
	"enter d2 A": And("mystery tree", "remove bush"),
	"enter d2 B": Or(
		And("mystery tree", "woods of winter default summer", "bracelet"),
		And("d2 blade key chest", "bracelet")),
	"enter d2":         Or("enter d2 A", "enter d2 B"),
	"outdoor d2 chest": AndSlot("enter d2 B"),
	"mystery cave chest": AndSlot("mystery tree", "remove mushroom",
		Or("feather L-2", "magnet gloves"),
		Or("woods of winter default autumn", And("autumn",
			Or("woods of winter default summer", And("enter d2 B", "bracelet"))))),
	"moblin road chest": AndSlot("fairy fountain", "bombs", "remove bush",
		Or("eastern suburbs default winter", "winter"),
		Or("woods of winter default spring", "spring",
			"woods of winter default summer", "summer",
			"woods of winter default autumn", "autumn")),
	"moblin cliff chest": AndSlot("suburbs", "bracelet",
		Or("eastern suburbs default spring", "spring")),
	"linked dive chest": AndSlot("fairy fountain",
		Or("eastern suburbs default winter", "winter"),
		Or("flippers", "feather L-2")),

	// holodrum plain
	"ghastly stump": Or(
		"pegasus tree",
		And("scent tree", Or("jump", "ricky",
			And("flippers", "remove bush"), "holodrum plain default winter")),
		And("south swamp", Or("flippers", Hard("dimitri flute")), "remove bush")),
	"scent tree": OrSlot(
		And("ghastly stump", Or("jump", "ricky", Hard("moosh flute"), "winter",
			"holodrum plain default winter")),
		And("south swamp", Or("flippers", Hard("dimitri flute"))),
		And("sunken city", "animal flute"),
		And("north horon stump",
			Or("bracelet", And("remove bush", Or("flippers", Hard("dimitri flute"))))),
		And("temple remains", "long jump"),
		And("goron mountain", "flippers")),
	"blaino gift":      AndSlot("scent tree", "rupees"),
	"ricky":            And("scent tree", "ricky's gloves"),
	"round jewel gift": AndSlot("scent tree", Or("flippers", Hard("dimitri flute"))),
	"water cave chest": AndSlot("scent tree", "flippers"),
	"mushroom cave chest": AndSlot("scent tree", "remove mushroom", "flippers",
		Or("holodrum plain default autumn", And("ghastly stump", "autumn"))),

	// spool swamp
	"pegasus tree": AndSlot("ghastly stump",
		Or("holodrum plain default summer", "summer", "feather L-2", "ricky", Hard("moosh flute"))),
	"floodgate key spot": AndSlot("pegasus tree", Or("bracelet", "hit lever")),
	"spool stump": And("pegasus tree", Or("bracelet", "hit lever"),
		Or("pegasus satchel", "flippers", "feather L-2"), "bracelet", "floodgate key"),
	"dry swamp": Or(
		Or("spool swamp default summer", "spool swamp default autumn", "spool swamp default winter"),
		And("spool stump", Or("summer", "autumn", "winter"))),
	"south swamp": Or(
		And("spool stump", Or("flippers", Hard("dimitri flute"))),
		And("spool stump", "dry swamp", Or("long jump", "animal flute")),
		And("ghastly stump", "remove bush", Or("flippers", Hard("dimitri flute"))),
		And("scent tree", Or("flippers", Hard("dimitri flute"))),
		HardAnd("scent tree", "dimitri flute"),
		And("swamp portal", "bracelet")),
	"square jewel chest": AndSlot("south swamp", Or("spool swamp default winter",
		And("spool stump", "winter")), Or("shovel", "animal flute"), "bombs"),
	"enter d3": And("spool stump", Or("spool swamp default summer", "summer")),

	// north horon / eyeglass lake
	"not north horon default summer": Or(
		"north horon default spring", "north horon default autumn", "north horon default winter"),
	"north horon stump": Or(
		And("horon village", "remove bush"),
		And("scent tree", "bracelet"),
		And("south swamp", "flippers", "remove bush"),
		HardAnd("south swamp", "dimitri flute"),
		And("lake portal", "not north horon default summer", "flippers", "jump"),
		And("lake portal", "pegasus jump L-2", "north horon default winter")),
	"enter d1": And("gnarled key", Or(
		And("south swamp", Or("flippers", Hard("dimitri flute"))),
		And("north horon stump", "remove bush"))),
	"wet eyeglass lake": Or("not north horon default summer", "spring", "autumn", "winter"),
	"d5 stump": And(Or("north horon default autumn", "autumn"), Or(
		And("lake portal", "not north horon default summer", "flippers"),
		And("north horon stump", Or("north horon default winter", And("winter", "autumn")),
			Or("jump", Hard("ricky", "moosh flute"))))),
	"enter d5": And("d5 stump", "remove mushroom",
		Or("north horon default autumn", "autumn")),
	"lake chest": AndSlot("horon village", Or("feather L-2", And("jump",
		Or("north horon default autumn",
			And("autumn", "north horon stump"))))),
	"dry lake west chest": AndSlot("d5 stump", "bracelet",
		Or("summer", And("enter d5", "north horon default summer"))),
	"dry lake east chest": AndSlot(And("bombs", "flippers"), Or(
		And("north horon stump", "jump", Or("north horon default summer", "summer")),
		And("d5 stump", "summer"),
		And("enter d5", "north horon default summer", "bracelet", "jump"))),

	// natzu
	"great moblin chest": AndSlot(Or("flippers", "jump"), "bracelet",
		Or("animal flute", Hard("start")), Or("flippers", "pegasus jump L-2")),
	"platform chest": AndSlot("scent tree", "flippers"),

	// sunken city
	"sunken city": Or("fairy fountain",
		And("scent tree", Or("jump", "flippers"), "animal flute"),
		And("mount cucco", "flippers")),
	"sunken gale tree":      AndSlot("sunken city", "cross water gap"),
	"dimitri":               And("sunken gale tree", "bombs"),
	"master's plaque chest": AndSlot("dimitri", "sword", "cross water gap"),
	"diver gift":            AndSlot("dimitri", "master's plaque"),
	"sunken cave chest": AndSlot("sunken city", "flippers", "remove bush",
		Or("sunken city default summer", "summer")),
	"diver chest": AndSlot("dimitri"),

	// mount cucco
	"mount cucco": Or("mountain portal",
		And("sunken city", "flippers", Or("sunken city default summer", "summer")),
		And("goron mountain", "bracelet", "shovel")),
	"spring banana cucco": And("mount cucco", "bracelet"),
	"spring banana tree": AndSlot("spring banana cucco", "jump",
		Or("sunken city default spring", "spring"), Or("sword", "fool's ore")),
	"moosh": And("mount cucco", "spring banana"),
	"dragon key spot": AndSlot("mount cucco",
		Or("moosh", "pegasus jump L-2", Hard("feather L-2"))),
	"talon cave":         And("mount cucco", Or("sunken city default spring", "spring")),
	"dragon keyhole":     And("talon cave", "winter", "jump", "bracelet"),
	"enter d4":           And("dragon key", "dragon keyhole", "summer", "cross water gap"),
	"pyramid jewel spot": AndSlot("talon cave", "flippers"),
	"talon cave chest":   AndSlot("talon cave"),

	// goron mountain
	"goron mountain": Or(
		And("mount cucco", "shovel", "bracelet"),
		And("temple remains", "long jump", Or("flippers", "pegasus jump L-2")),
		And("scent tree", "flippers")),
	"goron chest": AndSlot("goron mountain", "feather L-2", "bombs"),

	// tarm ruins
	"tarm ruins": And("pegasus tree",
		"square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	"lost woods": And("tarm ruins", Or("lost woods default summer", "summer"),
		Or("lost woods default winter", "winter"), "autumn", "remove mushroom"),
	"noble sword spot": AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"tarm gale tree":   AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"enter d6": And("tarm gale tree", Or("tarm ruins default winter", "winter"),
		"shovel", "remove bush"),
	"tarm gasha chest": AndSlot("tarm gale tree", "remove mushroom", "ember seeds",
		Or("tarm ruins default autumn", "autumn")),

	// samasa desert
	"desert":       And("suburbs", "pirate house"),
	"desert pit":   AndSlot("desert", "bracelet"),
	"desert chest": AndSlot("desert", "flippers"),

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
