package prenode

// new set of holodrum prenodes, accounting for randomized seasons
//
// these use new nested prenode constructors, which are nice, but hard to read
// if you nest them in a way that requires multiple levels of indentation. so
// try not to do that too much.

var holodrum2Prenodes = map[string]*Prenode{
	// horon village & western coast
	"village": Root(),
	"x-shaped jewel chest": AndSlot("village", Or("ember slingshot", "mystery slingshot"),
		"mystery seeds", "kill moldorm"),
	"enter d0": Root(),
	"pirate portal": Or(
		And("village", Or("boomerang L-2", "pegasus jump L-2")),
		And("hide and seek", "jump")),
	"maku key gift": And("village", "pop maku bubble"),

	// eastern commons
	"commons": Or(
		And("village", "ember seeds"),
		And("rosa portal", "remove bush"),
		And("fairy fountain", Or("commons winter", "cross water gap"))),
	"fairy fountain": Or("sunken city"), // TODO

	// holodrum plain
	"ghastly stump": Or(
		"pegasus tree",
		And("scent tree", Or("jump", "ricky", "holodrum plain default winter")),
		And("south swamp", "flippers", "remove bush")),
	"scent tree": OrSlot(
		And("ghastly stump", Or("jump", "ricky", "winter", "holodrum plain default winter")),
		And("south swamp", "flippers"),
		And("natzu"), // TODO natzu is complicated
		And("north horon stump", "bracelet"),
		And("temple remains", "long jump")),

	// spool swamp
	"pegasus tree": And("ghastly stump",
		Or("holodrum plain default summer", "summer", "feather L-2")),
	"spool stump": And("pegasus tree", "hit lever",
		Or("pegasus satchel", "flippers", "feather L-2"), "bracelet", "floodgate key"),
	"dry swamp": Or(
		Not("spool swamp default spring"),
		And("spool stump", Or("summer", "autumn", "winter"))),
	"south swamp": Or(
		And("spool stump", "flippers"),
		And("spool stump", "dry swamp", Or("long jump", "animal flute")),
		And("ghastly stump", "remove bush", "flippers"),
		And("scent tree", "flippers"),
		And("swamp portal", "bracelet")),
	"square jewel chest": And(Or("spool swamp default winter", And("spool stump", "winter")),
		"shovel", Or("animal flute", "bombs")),
	// TODO swamp portal from other end
	"swamp portal": And("south swamp", "bracelet"),
	"enter d3":     And("spool stump", Or("spool swamp default summer", "summer")),

	// north horon / eyeglass lake
	"north horon default spring": Root(),
	"north horon default summer": Root(),
	"north horon default autumn": Root(),
	"north horon default winter": Root(),
	"north horon stump": Or(
		And("village", "remove bush"),
		And("scent tree", "bracelet"),
		And("south swamp", "flippers", "remove bush"),
		And("lake portal", Not("north horon default summer"), "flippers", "jump"),
		And("lake portal", "pegasus jump L-2", "north horon default winter")),
	"enter d1": And("gnarled key", Or(
		And("south swamp", "flippers"),
		And("north horon stump", "remove bush"))),
	"wet eyeglass lake": Or(Not("north horon default summer"), "spring", "autumn", "winter"),
	// TODO lake portal from the other end
	"lake portal": And("north horon stump", Or(
		And("wet eyeglass lake", Or("jump", "animal flute"), "flippers"),
		And(Or("north horon default winter", "winter"), "pegasus jump L-2"))),
	"enter d5": And(Or("north horon default autumn", "autumn"), "remove mushroom", Or(
		And("lake portal", Not("north horon default summer"), "flippers"),
		And("north horon stump", "north horon default winter", Or("jump", "animal flute")))),

	// natzu
	// TODO it's complicated

	// tarm ruins
	"tarm ruins": And("pegasus tree",
		"square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	// TODO how exactly does the statue/season "puzzle" work?
	"lost woods": And("tarm ruins", Or("tarm ruins default summer", "summer"),
		"winter", "autumn", "remove mushroom"),

	// samasa desert
	"desert":          And("commons", "pirate house"),
	"rusty bell slot": And("desert", "bracelet"),
	"desert portal":   And("desert"),
}
