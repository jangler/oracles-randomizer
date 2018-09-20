package logic

// "ricky", "dimitri", and "moosh" refer to accessing those animal companions
// in their designated regions (e.g. dimitri in sunken city). "x's flute" means
// being able to call the animal in general.

var holodrumNodes = map[string]*Node{
	"start": And(), // parent for nodes reachable by default

	// horon village
	"horon village":    And("start"),
	"maku tree gift":   AndSlot("horon village", "pop maku bubble"),
	"ember tree":       AndSlot("horon village"),
	"village SE chest": AndSlot("horon village", "bombs"),
	"village SW chest": AndSlot("horon village", Or("remove bush safe", "flute"),
		Or("remove mushroom", "dimitri's flute")),
	"village shop 1":  AndSlot("start"),
	"village shop 2":  AndSlot("start"),
	"village shop 3":  AndSlot("start"),
	"member's shop 1": AndSlot("member's card"),
	"member's shop 2": AndSlot("member's card"),
	"member's shop 3": AndSlot("member's card"),

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
		And("rosa portal", Or("remove bush", "flute")),
		And("fairy fountain", Or("eastern suburbs default winter", "winter",
			Or("cross water gap", "dimitri's flute")))),
	"fairy fountain": Or(
		And("sunken city",
			Or(Hard("start"), "eastern suburbs default spring", "spring")),
		And("suburbs", Or("eastern suburbs default winter", "winter",
			"cross water gap", "ricky's flute", "dimitri's flute"))),
	"moblin road": Or(
		And("fairy fountain", Or("eastern suburbs default winter", "winter")),
		And("sunken city", "flippers", Or(
			"sunken city default spring", "spring",
			"sunken city default summer", "summer",
			"sunken city default autumn", "autumn"),
			Or(Hard("start"), And(
				Or("eastern suburbs default winter", "winter"),
				Or("eastern suburbs default spring", "spring"))))),
	"shovel gift": AndSlot("moblin road",
		Or("woods of winter default winter", "winter")),
	"central woods of winter": Or(
		And("fairy fountain", Or("eastern suburbs default winter", "winter"),
			Or("shovel", "jump", "flute")),
		And("fairy fountain", Or(
			"eastern suburbs default spring", "spring",
			"eastern suburbs default summer", "summer",
			"eastern suburbs default autumn", "autumn"))),
	"mystery tree": OrSlot(
		And("fairy fountain", Or("eastern suburbs default winter", "winter"),
			Or("shovel", And("jump", "bracelet"), "flute")),
		And("fairy fountain", Or(
			"eastern suburbs default spring", "spring",
			"eastern suburbs default summer", "summer",
			"eastern suburbs default autumn", "autumn"))),
	"enter d2 A": And("mystery tree", Or("remove bush", "flute")),
	"enter d2 B": Or(
		And("mystery tree", "bracelet",
			Or("woods of winter default summer", "ricky's flute")),
		And("d2 blade key chest", "bracelet")),
	"enter d2":         Or("enter d2 A", "enter d2 B"),
	"outdoor d2 chest": AndSlot("enter d2 B"),
	"mystery cave chest": AndSlot("mystery tree",
		Or("remove mushroom", "dimitri's flute"),
		Or("feather L-2", "magnet gloves"),
		Or("woods of winter default autumn", And("autumn",
			Or("ricky's flute", "woods of winter default summer",
				And("enter d2 B", "bracelet"))))),
	"moblin road chest": AndSlot("moblin road",
		Or("bombs", "ricky's flute"), "remove bush safe",
		Or("woods of winter default spring", "spring",
			"woods of winter default summer", "summer",
			"woods of winter default autumn", "autumn")),
	"moblin cliff chest": AndSlot("suburbs", "bracelet",
		Or("long jump", "magnet gloves"),
		Or("eastern suburbs default spring", "spring")),
	"linked dive chest": AndSlot("moblin road", Or("flippers", "feather L-2")),

	// holodrum plain
	"ghastly stump": Or(
		"pegasus tree",
		And("scent tree", Or("jump", "ricky", "flute",
			And("flippers", "remove bush"), "holodrum plain default winter")),
		And("south swamp", Or("flippers", "dimitri's flute"), "remove bush")),
	"scent tree": OrSlot(
		And("ghastly stump", Or("jump", "ricky", "flute", "winter",
			"holodrum plain default winter")),
		And("south swamp", Or("flippers", "dimitri's flute")),
		And("sunken city", Or(
			And("natzu prairie", "flute"),
			And("natzu river", "jump", Or("flippers", "flute")),
			And("natzu wasteland", Or("flute", And("remove bush",
				Or("feather L-2", Hard("long jump"))))))),
		And("north horon stump", Or("bracelet",
			And(Or("remove bush", "flute"),
				Or("flippers", "dimitri's flute")))),
		And("temple remains", "long jump"),
		And("goron mountain", "flippers")),
	"blaino gift":      AndSlot("scent tree"),
	"ricky":            Or("ricky's flute"),
	"round jewel gift": AndSlot("scent tree", Or("flippers", "dimitri's flute")),
	"water cave chest": AndSlot("scent tree", "flippers"),
	"mushroom cave chest": AndSlot("scent tree", "flippers",
		Or("remove mushroom", "dimitri's flute"),
		Or("holodrum plain default autumn", And("ghastly stump", "autumn"))),

	// spool swamp
	"pegasus tree": AndSlot("ghastly stump", Or("holodrum plain default summer",
		"summer", "feather L-2", "ricky", "moosh's flute")),
	"floodgate key spot": AndSlot("pegasus tree", "hit lever"),
	"spool stump": And("pegasus tree", Or("bracelet", "hit lever"),
		Or("pegasus satchel", "flippers", "feather L-2"), "bracelet", "floodgate key"),
	"dry swamp": Or(Or("spool swamp default summer",
		"spool swamp default autumn", "spool swamp default winter"),
		And("spool stump", Or("summer", "autumn", "winter"))),
	"south swamp": Or(
		And("spool stump", Or("flippers", "dimitri's flute")),
		And("spool stump", "dry swamp", Or("long jump", "flute")),
		And("ghastly stump", "remove bush", Or("flippers", "dimitri's flute")),
		And("scent tree", Or("flippers", "dimitri's flute")),
		And("swamp portal", "bracelet")),
	"square jewel chest": AndSlot("south swamp",
		Or("spool swamp default winter", And("spool stump", "winter")),
		Or("shovel", "flute"), Or("bombs", "ricky's flute")),
	"enter d3": And("spool stump", Or("spool swamp default summer", "summer")),

	// north horon / eyeglass lake
	"not north horon default summer": Or("north horon default spring",
		"north horon default autumn", "north horon default winter"),
	"north horon stump": Or(
		And("horon village", Or("remove bush", "flute")),
		And("scent tree", "bracelet"),
		And("south swamp", Or("flippers", "dimitri's flute"), Or("remove bush", "flute")),
		And("lake portal", "not north horon default summer", "flippers", "jump"),
		And("lake portal", "pegasus jump L-2", "north horon default winter")),
	"enter d1": And("gnarled key", Or(
		And("south swamp", Or("flippers", "dimitri's flute")),
		And("north horon stump", Or("remove bush", "flute")))),
	"wet eyeglass lake": Or("not north horon default summer", "spring", "autumn", "winter"),
	"d5 stump": Or(And("lake portal", "not north horon default summer", "flippers"),
		And("north horon stump", Or("jump", "ricky's flute", "moosh's flute"),
			Or("north horon default winter", "winter", "flippers",
				And("bracelet", "dimitri's flute")))),
	"enter d5": And("d5 stump", Or("remove mushroom", "dimitri's flute"),
		Or(And("north horon default autumn",
			Or("lake portal", "jump", "ricky's flute", "moosh's flute"),
			Or("flippers", And("dimitri's flute", "bracelet"))), "autumn")),
	"lake chest": AndSlot("horon village", Or("feather L-2", And("jump",
		Or("north horon default autumn",
			And("autumn", "north horon stump"))))),
	"dry lake east chest": AndSlot("d5 stump", "bracelet",
		Or("summer", And("enter d5", "north horon default summer"))),
	"dry lake west chest": AndSlot(And(Or("bombs", "ricky's flute")), "flippers", Or(
		And("north horon stump", Or("jump", "ricky's flute", "moosh's flute"),
			Or("north horon default summer", "summer")),
		And("d5 stump", "summer", "flippers"),
		And("enter d5", "north horon default summer", "bracelet", "jump"))),

	// natzu
	"natzu prairie":   Root("start"),
	"natzu river":     Root(),
	"natzu wasteland": Root(),
	"natzu":           Or("natzu prairie", "natzu river", "natzu wasteland"),
	"great moblin chest": AndSlot(Or("flippers", "pegasus jump L-2"), "bracelet", Or(
		And("natzu prairie", "sunken city"),
		And("natzu river", "scent tree", "dimitri's flute"),
		And("natzu wasteland", And("sunken city",
			Or("flute", And("remove bush", "long jump")))))),
	"platform chest": OrSlot(And("scent tree", Or("flippers", "dimitri's flute")),
		And("sunken city", "natzu river", "pegasus jump L-2")),

	// sunken city
	"sunken city": Or(
		And("mount cucco", "flippers",
			Or(Hard("start"), "summer", "sunken city default summer")),
		And("fairy fountain", Or("eastern suburbs default spring", "spring")),
		And("scent tree", Or(
			And("natzu prairie", "flute"),
			And("natzu river", Or(And("jump", Or("flippers", "flute")),
				HardAnd("flute", "flippers"))),
			And("natzu wasteland", Or("flute",
				And("remove bush", Or("feather L-2", Hard("long jump")))))))),
	"sunken gale tree": AndSlot("sunken city",
		Or("cross water gap", "sunken city default winter")),
	"dimitri":               And("sunken gale tree", Or("dimitri's flute", "bombs")),
	"master's plaque chest": AndSlot("dimitri", "sword", "cross water gap"),
	"diver gift":            AndSlot("dimitri", "master's plaque"),
	"sunken cave chest": AndSlot("sunken city", "flippers", "remove bush safe",
		Or("sunken city default summer", "summer")),
	"diver chest": AndSlot("dimitri"),

	// mount cucco
	"mount cucco": Or("mountain portal",
		And("sunken city", "flippers", Or("sunken city default summer", "summer")),
		And("goron mountain", "bracelet", "shovel")),
	"spring banana cucco": And("mount cucco", "bracelet"),
	"spring banana tree": AndSlot("spring banana cucco", "jump",
		Or("sunken city default spring", "spring"), Or("sword", "fool's ore")),
	"moosh": And("mount cucco", Or("moosh's flute", "spring banana")),
	"dragon key spot": AndSlot("mount cucco",
		Or("moosh", "pegasus jump L-2", Hard("feather L-2"))),
	"talon cave":         And("mount cucco", Or("sunken city default spring", "spring")),
	"dragon keyhole":     And("talon cave", "winter", "jump", "bracelet"),
	"enter d4":           And("dragon key", "dragon keyhole", "summer", "cross water gap"),
	"pyramid jewel spot": AndSlot("talon cave", "flippers"),
	"talon cave chest":   AndSlot("talon cave"),

	// goron mountain
	"goron mountain": Or(
		And("mount cucco", Or("shovel", "spring banana"), "bracelet"),
		And("temple remains", "long jump", Or("flippers", "pegasus jump L-2")),
		And("scent tree", "flippers")),
	"goron chest": AndSlot("goron mountain", "bombs",
		Or("feather L-2", Hard("long jump"))),

	// tarm ruins
	"tarm ruins": And("pegasus tree",
		"square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	"lost woods": And("tarm ruins", Or("lost woods default summer", "summer"),
		Or("lost woods default winter", "winter"), "autumn", "remove mushroom"),
	"noble sword spot": AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"tarm gale tree":   AndSlot("lost woods", "winter", "autumn", "spring", "summer"),
	"enter d6": And("tarm gale tree", Or("tarm ruins default winter", "winter"),
		Or("shovel", "ember seeds"), "remove bush"),
	"tarm gasha chest": AndSlot("tarm gale tree", "remove mushroom", "ember seeds",
		Or("tarm ruins default autumn", "autumn")),

	// samasa desert
	"desert":       And("suburbs", "pirate house"),
	"desert pit":   AndSlot("desert", "bracelet"),
	"desert chest": AndSlot("desert", "flippers"),

	// temple remains (the important logic is in the portal nodes)
	"temple remains": Or(
		And("goron mountain", Or("flippers", "pegasus jump L-2"), "long jump"),
		And("scent tree", "long jump")),

	// northern peak
	"maku seed": And("d1 essence", "d2 essence", "d3 essence", "d4 essence",
		"d5 essence", "d6 essence", "d7 essence", "d8 essence"),
	"enter d9": And("scent tree", "maku seed"),

	// old men
	"goron mountain old man": And("goron mountain", "ember seeds"),
	"western coast old man":  And("pirate ship", "ember seeds"),
	"holodrum plain east old man": And("scent tree", "ember seeds",
		Or("ricky's flute", "holodrum plain default summer",
			And("ghastly stump", "summer"))),
	"horon village old man":       And("horon village", "ember seeds"),
	"north horon old man":         And("north horon stump", "ember seeds"),
	"tarm ruins old man":          And("enter d6", "ember seeds"),
	"woods of winter old man":     And("shovel gift", "ember seeds"),
	"holodrum plain west old man": And("ghastly stump", "ember seeds"),
}

var seasonNodes = map[string]*Node{
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
