package logic

// skipping keys is often possible in ages dungeons, but it doesn't matter in
// logic because you could misuse the keys anyway. the logic always assumes the
// worst possible key usage.

// third key door is changed to a shutter from the "wrong" side, since if you
// use your first key on it, you can't reach the other two keys.
var agesD1Nodes = map[string]*Node{
	"d1 east terrace": AndSlot("enter d1", "kill zol"),
	"d1 ghini key":    And("d1 east terrace", "kill ghini"),
	"d1 crossroads":   AndSlot("d1 east terrace"),
	"d1 crystal room": AndSlot("d1 east terrace", "ember seeds",
		"break crystal"),
	"d1 free key chest":     And("d1 ghini key"),
	"d1 platform key chest": And("d1 ghini key"),
	"d1 button chest":       AndSlot("d1 ghini key"),
	"d1 U-room": Or("d1 west terrace", And("d1 free key chest",
		"d1 platform key chest", "break bush safe", "kill giant ghini")),
	"d1 basement":     AndSlot("d1 U-room", "ember seeds"),
	"d1 west terrace": AndSlot("enter d1", "break pot"),
	"d1 pot chest":    AndSlot("enter d1", "break pot"),
	"d1 essence": AndStep("d1 free key chest", "break bush safe", "d1 boss key",
		"kill pumpkin head"),
}

var agesD2Nodes = map[string]*Node{
	"d2 bombed terrace": AndSlot("enter d2", "kill spiked beetle", "bombs"),
	"d2 rope room": AndSlot("enter d2", "d2 key 1", "d2 color key",
		"d2 basement key", "d2 statue key"),
	"enter swoop": Or(And("enter d2", "kill spiked beetle", "feather"),
		And("d2 key 1", "d2 key 2")),
	"d2 basement":      And("enter swoop", "kill swoop"),
	"d2 thwomp tunnel": AndSlot("d2 basement"),
	"d2 thwomp shelf": AndSlot("d2 basement",
		Or("feather", And("cane", "pegasus satchel"))),
	"d2 moblin platform": AndSlot("d2 3 keys"),
	// push moblin into doorway, stand on button, use switch hook
	"d2 statue room": And("d2 moblin platform", Or("bracelet", "cane",
		HardAnd("switch hook", "push enemy"))),
	"d2 color room": AndSlot("d2 all keys"),
	"d2 essence":    AndStep("d2 all keys", "d2 boss key"),

	"d2 key 1":     And("enter d2", "kill spiked beetle", "kill normal"),
	"d2 key 2":     And("enter d2", "d2 key 1", "bombs"),
	"d2 color key": And("d2 basement", "feather"),
	"d2 basement key": And("d2 basement", "feather", "bombs",
		"hit lever from minecart", "kill normal"),
	"d2 3 keys": Or(
		And("d2 key 1",
			Or(And("d2 key 2", Or("d2 color key", "d2 basement key")),
				And("d2 color key", "d2 basement key"))),
		And("d2 key 2", "d2 color key", "d2 basement key")),
	"d2 statue key": And("d2 statue room", "feather"),
	"d2 all keys": And("d2 key 1", "d2 key 2", "d2 color key",
		"d2 basement key", "d2 statue key"),
}

// killing armos is an exception to the "bombs are hard logic" rule, and since
// you need bombs to do anything in d3, they're not even relevant to logic.
var agesD3Nodes = map[string]*Node{
	"d3 pols voice chest": AndSlot("enter d3", "bombs"),
	"d3 1F spinner":       And("enter d3", Or("kill moldorm", "bracelet")),
	"d3 S crystal":        And("d3 1F spinner"),
	"d3 E crystal":        And("d3 1F spinner", "bombs"),
	"d3 statue key":       And("d3 E crystal"),
	// you can clip into the blocks enough to hit this crystal with switch hook
	"d3 N crystal": And("d3 statue key",
		Or("any seed shooter", "boomerang", Hard("switch hook"))),
	"d3 armos key":        And("d3 statue key"),
	"d3 bush beetle room": AndSlot("d3 armos key"),
	"d3 W crystal":        And("d3 statue key"),
	"d3 compass key":      And("d3 statue key"),
	"d3 mimic room": AndSlot("d3 armos key", "d3 compass key",
		"kill moldorm"),

	"break crystal switch": Or("sword", "switch hook", "boomerang",
		"ember satchel", "scent satchel", "mystery satchel",
		"any seed shooter"),
	"d3 B1F spinner": And("d3 S crystal", "d3 E crystal", "d3 N crystal",
		"d3 W crystal", "break crystal switch"),
	"d3 crossroads":         AndSlot("d3 B1F spinner"),
	"d3 conveyor belt room": AndSlot("d3 statue key"),
	"d3 bridge chest": AndSlot("d3 statue key",
		Or("any seed shooter", "jump 3", HardAnd("d3 all keys", "feather"),
			And("boomerang", Or("feather", "pegasus satchel")))),
	"d3 torch chest": AndSlot("d3 B1F spinner",
		Or("ember shooter", Hard("mystery shooter"))),
	"kill subterror": And("shovel",
		Or("sword", "switch hook", "scent seeds", Hard("bombs"))),
	"d3 B1F east": AndSlot("d3 B1F spinner", "kill subterror",
		"any seed shooter"),
	"d3 block key": And("d3 B1F spinner", "kill subterror"),
	"d3 all keys":  And("d3 armos key", "d3 compass key", "d3 block key"),
	"d3 essence": AndStep("d3 boss key", "d3 all keys",
		Or("ember seeds", "scent seeds"),
		Or("seed shooter", And(
			Or("ember seeds", Hard()),
			Or("boomerang", Hard("jump 3"),
				HardAnd("feather", "sword", "switch hook"))))),
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
		Or("any seed shooter", HardAnd("jump 3", "boomerang"))),
	"d4 minecart C": And("d4 minecart B", "d4 key C"),
	"d4 minecart D": And("d4 minecart C", "d4 key D"),
	// these weapons are for the miniboss, not the moldorms
	"d4 small floor puzzle": AndSlot("d4 minecart D", "bombs",
		Or("sword", "switch hook", "scent shooter", Hard())),
	"d4 key chest E":    And("d4 minecart D", "switch hook"),
	"d4 lava pot chest": AndSlot("d4 key chest E", "d4 key E"),
	"d4 essence": AndStep("d4 key chest E", "d4 boss key",
		Or("sword", "boomerang")),

	"d4 key A": And("d4 key chest A"),
	"d4 key B": And("d4 key chest B"),
	"d4 key C": And("d4 key chest C"),
	"d4 key D": And("d4 minecart C", Or("sword", "ember seeds",
		"scent shooter", "gale shooter", Hard("scent satchel"))),
	"d4 key E": And("d4 key chest E"),
}

// every chest not behind a key door in d5 requires you to be able to hit a
// switch, so that's a requirement for the first node.
var agesD5Nodes = map[string]*Node{
	"d5 switch A":       And("enter d5", "kill normal", "hit switch"),
	"d5 blue peg chest": AndSlot("d5 switch A"),
	"d5 dark chest": And("d5 switch A",
		Or("cane", "switch hook", HardOr("kill normal", "push enemy"))),
	"d5 boxed chest": And("d5 switch A"),
	"d5 eyes chest":  And("d5 switch A", "any seed shooter"),
	"d5 2-statue chest": And("d5 switch A", "break pot", "cane", "feather",
		Or("any seed shooter", "boomerang", HardAnd("feather", "sword"))),
	"d5 essence": AndStep("d5 switch A", "d5 boss key", "cane", "sword"),

	// require 1 small key minimum, 2 maximum.
	// keys A (dark chest) and E (3-statue chest) are always available by now.
	"d5 crossroads": And("d5 switch A", "feather", "bracelet",
		Or("cane", Hard("jump 3"))),
	"d5 diamond chest": AndSlot("d5 crossroads", "switch hook"),

	// require 1 small key minimum, 5 maximum.
	"d5 3-statue chest":    And("d5 switch A", "cane"),
	"d5 six-statue puzzle": AndSlot("d5 all keys", "ember shooter", "feather"),

	// require 4 small keys minimum, 5 maximum.
	"d5 red peg chest": AndSlot("d5 crossroads", "d5 all keys",
		"hit switch ranged"),
	"d5 owl puzzle": AndSlot("d5 red peg chest"),

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
	"d6 past color room": AndSlot("enter d6 past", "feather", "kill gel"),
	"d6 past wizzrobe chest": AndSlot("enter d6 past", "bombs",
		"kill wizzrobe"),
	"d6 past pool chest": AndSlot("enter d6 past", "bombs", "ember seeds",
		"flippers"),
	"d6 open wall": And("enter d6 past", "bombs", "ember shooter"),
	"d6 past stalfos chest": And("enter d6 past", "ember seeds",
		Or("kill normal ranged", "scent satchel", "feather")),
	"d6 past rope chest": And("d6 open wall", "mermaid suit"),

	// past, 1 key
	"d6 past spinner": And("enter d6 past", "cane", "bracelet", "feather",
		Or("d6 past key A", "d6 past key B"), "bombs"),
	"d6 past spear chest": AndSlot("d6 past spinner", "mermaid suit"),
	"d6 past diamond chest": And("d6 past spinner", "mermaid suit",
		"switch hook"),

	// past, 3 keys
	"d6 essence": AndStep("d6 past spinner", "d6 past key A", "d6 past key B",
		"d6 past key C", "d6 boss key", "any seed shooter"),

	"d6 past key A": And("d6 past stalfos chest"),
	"d6 past key B": And("d6 past rope chest"),
	"d6 past key C": And("d6 past diamond chest"),

	// present, 0 keys
	"d6 present diamond chest": AndSlot("enter d6 present", "switch hook"),
	"d6 present rope chest": And("enter d6 present", "scent satchel",
		Or("flippers", "feather", "switch hook"),
		Or("any seed shooter", "boomerang", "jump 3")),
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
	"d6 present RNG chest": AndSlot("d6 present beamos chest",
		"d6 present all keys", "bracelet", Or("sword", "cane", "switch hook")),
	"d6 present channel chest": AndSlot("enter d6 present", "d6 open wall",
		"d6 present all keys", "switch hook"),
	"d6 present vire chest": AndSlot("d6 present spinner chest",
		"d6 present all keys", Or("sword", Hard()), "switch hook"),

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
	"drain d7": And("refill d7"),
	// and get these chests
	"d7 spike chest":    AndSlot("refill d7"),
	"d7 stairway chest": AndSlot("refill d7"),
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

	"d7 essence": AndStep("d7 boss key", "flood d7"),
}

// small keys aren't randomized, so the items alone here are enough to get the
// required keys.
var agesD8Nodes = map[string]*Node{
	"d8 group A": And("enter d8", "bombs"), // only has small key

	"d8 group B": And("d8 group A", "switch hook", "cane",
		"seed shooter", Or("ember seeds", Hard("mystery seeds"))), // +1 key
	"d8 isolated chest": AndSlot("d8 group B"),

	"d8 group C":           And("d8 group B"), // +1 small key
	"d8 blue peg chest":    AndSlot("d8 group C"),
	"d8 floor puzzle":      AndSlot("d8 group C"),
	"d8 sarcophagus chest": AndSlot("d8 group C", "power glove"),

	"d8 group D":  And("d8 group C", "sword"), // post-miniboss
	"d8 NW slate": And("d8 group D"),
	"d8 SW slate": And("d8 group D", "bracelet"), // +1 small key
	"d8 NE slate": And("d8 group D", "feather", "flippers", "ember seeds"),

	"d8 group G":   And("d8 group D", "power glove"),
	"d8 B3F chest": AndSlot("d8 group G"),
	"d8 tile room": AndSlot("d8 group G", "feather"),
	"d8 SE slate":  And("d8 group G", "feather"),
	"d8 essence": AndStep("d8 boss key", "d8 group G", "d8 NW slate",
		"d8 NE slate", "d8 SW slate", "d8 SE slate"),
}

var agesD9Nodes = map[string]*Node{
	"done": AndStep("maku seed", "mystery seeds", "switch hook", "sword"),
}
