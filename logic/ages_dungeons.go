package logic

var agesD1Nodes = map[string]*Node{
	// 0 keys
	"d1 east terrace": AndSlot("enter d1", "kill zol"),
	"d1 ghini drop":   AndSlot("d1 east terrace", "kill ghini"),
	"d1 crossroads":   AndSlot("d1 east terrace"),
	"d1 crystal room": AndSlot("d1 east terrace", "ember seeds",
		"break crystal"),
	"d1 west terrace": AndSlot("enter d1", "break pot"),
	"d1 pot chest":    AndSlot("enter d1", "break pot"),

	// 2 keys
	"d1 wide room chest":  AndSlot("d1 ghini drop", Count(2, "d1 small key")),
	"d1 two-button chest": AndSlot("d1 wide room chest"),
	"d1 one-button chest": AndSlot("d1 wide room chest"),
	"d1 boss": AndSlot("d1 wide room chest", "break bush safe", "d1 boss key",
		"kill pumpkin head"),

	// potentially 3 keys w/ vanilla route
	"d1 U-room": Or("d1 west terrace",
		And("d1 wide room chest", "break bush safe", "kill giant ghini",
			Count(3, "d1 small key"))),
	"d1 basement": AndSlot("d1 U-room", "ember seeds"),
}

var agesD2Nodes = map[string]*Node{
	// 0 keys
	"spiked beetles owl": And("mystery seeds", "enter d2"),
	"d2 bombed terrace":  AndSlot("enter d2", "kill spiked beetle", "bombs"),
	"d2 moblin drop": AndSlot("enter d2", "kill spiked beetle",
		"kill normal"),

	// potentially 2 keys w/ vanilla route
	"enter swoop": Or(And("enter d2", "kill spiked beetle", "feather"),
		Count(2, "d2 small key")),
	"d2 basement":      And("enter swoop", "kill swoop"),
	"blue wing owl":    And("mystery seeds", "d2 basement"),
	"d2 thwomp tunnel": AndSlot("d2 basement"),
	"d2 thwomp shelf": AndSlot("d2 basement",
		Or("feather", HardAnd("cane", "pegasus satchel"))),
	"d2 color key": And("d2 basement", "feather"), // TODO
	"d2 basement key": And("d2 basement", "feather", "bombs",
		"hit lever from minecart", "kill normal"), // TODO

	// 3 keys
	"d2 moblin platform": AndSlot("d2 basement", Count(3, "d2 small key")),
	// push moblin into doorway, stand on button, use switch hook
	"d2 statue room": AndSlot("d2 moblin platform",
		Or("bracelet", "cane", HardAnd("switch hook", "push enemy"))),

	// 4 keys
	"d2 rope room": AndSlot("enter d2", "kill rope",
		Count(4, "d2 small key")),
	"d2 ladder chest": AndSlot("enter d2", Count(4, "d2 small key"), "bombs"),

	// 5 keys
	"d2 color room":   AndSlot("d2 statue room", Count(5, "d2 small key")),
	"d2 boss":         AndSlot("d2 color room", "d2 boss key"),
	"head thwomp owl": And("mystery seeds", "d2 boss"),
}

var agesD3Nodes = map[string]*Node{
	"d3 pols voice chest": AndSlot("enter d3", "bombs"),
	"d3 1F spinner":       And("enter d3", Or("kill moldorm", "bracelet")),
	"d3 S crystal":        And("d3 1F spinner"),
	"four crystals owl":   And("mystery seeds", "d3 1F spinner"),
	"d3 E crystal":        And("d3 1F spinner", "bombs"),
	"d3 statue key":       And("d3 E crystal"),
	// you can clip into the blocks enough to hit this crystal with switch hook
	"d3 N crystal": And("d3 statue key",
		Or("any seed shooter", "boomerang", Hard("switch hook"))),
	"stone soldiers owl":  And("mystery seeds", "d3 statue key"),
	"d3 armos key":        And("d3 statue key"),
	"d3 bush beetle room": AndSlot("d3 armos key"),
	"d3 W crystal":        And("d3 statue key"),
	"d3 compass key":      And("d3 statue key"),
	"d3 mimic room":       AndSlot("d3 all keys"),

	"break crystal switch": Or("sword", "switch hook", "boomerang",
		"ember satchel", "scent satchel", "mystery satchel",
		"any seed shooter", "punch object"),
	"d3 B1F spinner": And("d3 S crystal", "d3 E crystal", "d3 N crystal",
		"d3 W crystal", "break crystal switch"),
	"d3 crossroads":         AndSlot("d3 B1F spinner"),
	"d3 conveyor belt room": AndSlot("d3 statue key"),
	"d3 bridge chest": AndSlot("d3 statue key",
		Or("any seed shooter", "jump 3", HardAnd("d3 all keys", "feather"),
			HardAnd(Or("boomerang", And("bracelet", "toss ring")),
				Or("feather", "pegasus satchel")))),
	"d3 torch chest": AndSlot("d3 B1F spinner",
		Or("ember shooter", Hard("mystery shooter"))),
	"kill subterror": And("shovel", Or("sword", "switch hook", "scent seeds",
		"punch enemy", Hard("bombs"))),
	"d3 B1F east": AndSlot("d3 B1F spinner", "kill subterror",
		Or("any seed shooter", Hard("sword"))), // spin slash through corner
	"d3 moldorm key": And("d3 B1F spinner", "kill subterror"),
	"d3 all keys":    And("d3 armos key", "d3 compass key", "d3 moldorm key"),
	"d3 boss": AndSlot("d3 boss key", "d3 all keys",
		Or("ember seeds", "scent seeds"),
		Or("seed shooter", And(
			Or("ember seeds", Hard()),
			Or("d3 bridge chest", "jump 3", Hard("feather")),
			Or("boomerang", Hard("jump 3"), HardAnd("feather", "sword",
				Or("switch hook", And("bracelet", "toss ring"))))))),
}

var agesD4Nodes = map[string]*Node{
	"d4 first chest": AndSlot("enter d4", Or("kill stalfos", "push enemy"),
		Or("feather", "switch hook")),
	"d4 key chest A": And("d4 first chest", "feather"),
	"d4 minecart A":  And("enter d4", "feather", "d4 key A"),
	"d4 key chest B": And("d4 minecart A",
		Or("any seed shooter", Hard("boomerang"))),
	"d4 minecart chest": AndSlot("d4 minecart A", "hit lever"),
	"d4 minecart B": And("d4 minecart A", "hit lever",
		"d4 key B", "bracelet", "kill stalfos"),
	"d4 key chest C": And("d4 minecart B",
		Or("any seed shooter", Hard("boomerang"))),
	"d4 minecart C": And("d4 minecart B", "d4 key C"),
	"d4 minecart D": And("d4 minecart C", "d4 key D"),
	// these weapons are for the miniboss, not the moldorms
	"d4 small floor puzzle": AndSlot("d4 minecart D", "bombs",
		Or("sword", "switch hook", "scent shooter", "punch enemy", Hard())),
	"d4 key chest E":    And("d4 minecart D", "switch hook"),
	"d4 lava pot chest": AndSlot("d4 key chest E", "d4 key E"),
	"d4 boss": AndSlot("d4 key chest E", "d4 boss key",
		Or("sword", "boomerang", "punch enemy")),

	"d4 key A": And("d4 key chest A"),
	"d4 key B": And("d4 key chest B"),
	"d4 key C": And("d4 key chest C"),
	"d4 key D": And("d4 minecart C", Or("sword", "ember seeds", "scent shooter",
		"gale shooter", HardOr("scent satchel", "peace ring"))),
	"d4 key E": And("d4 key chest E"),
}

// every chest not behind a key door in d5 requires you to be able to hit a
// switch, so that's a requirement for the first node.
var agesD5Nodes = map[string]*Node{
	"d5 switch A": And("enter d5", "kill normal",
		Or("hit switch", Hard("bracelet"))),
	"d5 blue peg chest": AndSlot("d5 switch A"),
	"d5 dark chest": And("d5 switch A", "hit switch", // can't use pots here
		Or("cane", "switch hook", HardOr("kill normal", "push enemy"))),
	"d5 boxed chest": And("d5 switch A",
		Or("hit switch ranged", Hard("bracelet"))),
	"d5 eyes chest": And("d5 switch A", Or("any seed shooter",
		HardAnd("pegasus satchel", "feather",
			Or("hit switch ranged", And("bracelet", "toss ring")),
			Or("ember seeds", "scent seeds", "mystery seeds")))),
	"d5 2-statue chest": And("d5 switch A", "break pot", "cane", "feather"),
	"d5 boss":           AndSlot("d5 switch A", "d5 boss key", "cane", "sword"),

	// require 1 small key minimum, 2 maximum.
	// keys A (dark chest) and E (3-statue chest) are always available by now.
	// sword is to manip RNG to switch with the darknut.
	"d5 crossroads": And("d5 switch A", "feather", "bracelet",
		Or("cane", Hard("jump 3"), HardAnd("sword", "switch hook"))),
	"d5 diamond chest": AndSlot("d5 crossroads", "switch hook"),

	// require 1 small key minimum, 5 maximum.
	"d5 3-statue chest":    And("d5 switch A", "cane"),
	"d5 six-statue puzzle": AndSlot("d5 all keys", "ember shooter", "feather"),

	// require 4 small keys minimum, 5 maximum.
	"d5 red peg chest": AndSlot("d5 crossroads", "d5 all keys",
		"hit switch ranged"),
	"crown dungeon owl": And("mystery seeds", "d5 red peg chest"),
	"d5 owl puzzle":     AndSlot("d5 red peg chest"),

	"d5 key A": And("d5 dark chest"),
	"d5 key B": And("d5 boxed chest"),
	"d5 key C": And("d5 eyes chest"),
	"d5 key D": And("d5 2-statue chest"),
	"d5 key E": And("d5 3-statue chest"),
	"d5 all keys": And("d5 key A", "d5 key B", "d5 key C", "d5 key D",
		"d5 key E"),
}

var agesD6Nodes = map[string]*Node{
	// past, 0 keys
	"d6 past color room": AndSlot("enter d6 past", "feather", "kill color gel"),
	"d6 past wizzrobe chest": AndSlot("enter d6 past", "bombs",
		"kill wizzrobe"),
	"d6 past pool chest": AndSlot("enter d6 past", "bombs", "ember seeds",
		"flippers"),
	"d6 open wall":    And("enter d6 past", "bombs", "ember shooter"),
	"deep waters owl": And("mystery seeds", "d6 open wall"),
	"d6 past stalfos chest": And("enter d6 past", "ember seeds",
		Or("kill normal ranged", "scent satchel", "feather", Hard())),
	"d6 past rope chest": And("d6 open wall", "mermaid suit"),

	// past, 1 key
	"d6 past spinner": And("enter d6 past", "cane", "bracelet", "feather",
		Or("d6 past key A", "d6 past key B"), "bombs"),
	"d6 past spear chest": AndSlot("d6 past spinner", "mermaid suit"),
	"d6 past diamond chest": And("d6 past spinner", "mermaid suit",
		"switch hook"),

	// past, 3 keys
	"d6 boss": AndSlot("d6 past spinner", "d6 past key A", "d6 past key B",
		"d6 past key C", "d6 boss key", "any seed shooter"),

	"d6 past key A": And("d6 past stalfos chest"),
	"d6 past key B": And("d6 past rope chest"),
	"d6 past key C": And("d6 past diamond chest"),

	// present, 0 keys
	"d6 present diamond chest": AndSlot("enter d6 present", "switch hook"),
	"d6 present rope room": And("enter d6 present",
		Or("flippers", "feather", "switch hook"),
		Or("any seed shooter", "boomerang", "jump 3")),
	"scent seduction owl":   And("mystery seeds", "d6 present rope room"),
	"d6 present rope chest": And("d6 present rope room", "scent satchel"),
	"d6 present hand room": And("enter d6 present",
		Or("flippers", "feather", "switch hook"),
		Or("any seed shooter", "boomerang",
			And("jump 3", Or("sword", "switch hook", "ember seeds",
				"scent seeds", "mystery seeds", Hard("bombs"))))),
	"d6 present color room": And("d6 present hand room", "bombs",
		"switch hook", Or("feather", Hard())),
	"d6 present spinner chest": And("d6 past spinner", "d6 present hand room",
		Or("feather", "switch hook")),
	"d6 present beamos chest": AndSlot("enter d6 present", "d6 open wall",
		Or("flippers", And("d6 present 2 keys", "switch hook")), "feather"),

	// present, 1+ keys (keys can be used in any order if player has flippers)
	"luck test owl": And("mystery seeds", "d6 present beamos chest",
		"d6 present all keys"),
	"d6 present RNG chest": AndSlot("d6 present beamos chest",
		"d6 present all keys", "bracelet", Or("sword", "cane", "switch hook")),
	"d6 present channel chest": AndSlot("enter d6 present", "d6 open wall",
		"d6 present all keys", "switch hook"),
	"d6 present vire chest": AndSlot("d6 present spinner chest",
		"d6 present all keys",
		Or("sword", "expert's ring", Hard()), "switch hook"),

	"d6 present key A": And("d6 present rope chest"),
	"d6 present key B": And("d6 present color room"),
	"d6 present key C": And("d6 present spinner chest"),
	"d6 present 2 keys": Or(
		And("d6 present key A", "d6 present key B"),
		And("d6 present key A", "d6 present key C"),
		And("d6 present key B", "d6 present key C")),
	"d6 present all keys": And("d6 present key A", "d6 present key B",
		"d6 present key C"),
}

// assume mermaid suit
// stating this logic in terms of small keys is not really viable since it's
// possible to cut off access to some of them by changing the water level
var agesD7Nodes = map[string]*Node{
	// compass chest is potentially free
	"d7 crab chest": AndSlot("enter d7",
		Or("kill underwater", And("drain d7", "kill normal"))),

	// but everything except compass chest needs to be locked behind this
	"refill d7": And("enter d7", Or("long hook", And("switch hook", "cane"))),
	// and those requirements are also enough to drain the dungeon
	"drain d7":             And("refill d7"),
	"jabu switch room owl": And("mystery seeds", "drain d7"),
	// and get these chests
	"d7 spike chest":    AndSlot("refill d7"),
	"d7 stairway chest": AndSlot("refill d7"),
	"golden isle owl":   And("mystery seeds", "refill d7"),
	// but cane is needed here, since long hook can skip it initially.
	"d7 pot island chest": AndSlot("refill d7", "cane"),

	// this chest requires feather to reach its key block, and one of the small
	// keys in the dungeon also requires feather. so it's always possible to
	// reach the other chests without feather?
	"d7 miniboss chest": AndSlot("refill d7", "feather", "cane",
		Or("sword", "boomerang", "scent shooter")),

	// long hook is required to flood the dungeon (1F and 2F submerged)
	"flood d7": And("refill d7", "long hook"),
	// which is enough to get everything but the boss key chest
	"d7 hallway chest": AndSlot("flood d7"),
	// which also requires cane, since it's needed to get a small key
	"d7 post-hallway chest": AndSlot("flood d7", "cane"),

	"plasmarine owl": And("mystery shooter", "flood d7"),
	"d7 boss":        AndSlot("d7 boss key", "flood d7"),
}

// small keys aren't randomized, so the items alone here are enough to get the
// required keys.
var agesD8Nodes = map[string]*Node{
	"open ears owl": And("mystery seeds", "enter d8"),
	"d8 group A":    And("enter d8", "bombs"), // only has small key

	"d8 group B": And("d8 group A", "switch hook", "cane",
		"seed shooter", Or("ember seeds", Hard("mystery seeds"))), // +1 key
	"d8 ghini chest": AndSlot("d8 group B"),

	"d8 group C":           And("d8 group B"), // +1 small key
	"d8 blue peg chest":    AndSlot("d8 group C"),
	"d8 floor puzzle":      AndSlot("d8 group C"),
	"d8 sarcophagus chest": AndSlot("d8 group C", "power glove"),

	"d8 group D":        And("d8 group C", "sword"), // post-miniboss
	"d8 NW slate chest": AndSlot("d8 group D"),
	"d8 SW slate chest": AndSlot("d8 group D", "bracelet"), // +1 small key
	"d8 NE slate chest": AndSlot("d8 group D", "feather", "flippers",
		"ember seeds"),
	"ancient words owl": And("mystery shooter", "d8 NE slate chest"),

	"d8 group G":        And("d8 group D", "power glove"),
	"d8 B3F chest":      AndSlot("d8 group G"),
	"d8 tile room":      AndSlot("d8 group G", "feather"),
	"d8 SE slate chest": AndSlot("d8 group G", "feather"),
	"d8 boss": AndSlot("d8 boss key", "d8 group G", "slate 1", "slate 2",
		"slate 3", "slate 4"),
}

var agesD9Nodes = map[string]*Node{
	"black tower owl": And("mystery seeds", "maku seed"),
	"done": AndStep("maku seed", "mystery seeds", "switch hook",
		Or("sword", "punch enemy"), "bombs"), // bombs in case of spider form
}
