package logic

var agesD1Nodes = map[string]*Node{
	// 0 keys
	"d1 east terrace": AndSlot("enter d1", "kill switch hook"),
	"d1 ghini drop":   AndSlot("d1 east terrace"),
	"d1 crossroads":   AndSlot("d1 east terrace"),
	"d1 crystal room": AndSlot("d1 east terrace", "ember seeds",
		"break crystal"),
	"d1 west terrace": AndSlot("enter d1", "break pot"),
	"d1 pot chest":    AndSlot("enter d1", "break pot"),

	// 2 keys
	"d1 wide room":        AndSlot("d1 ghini drop", Count(2, "d1 small key")),
	"d1 two-button chest": AndSlot("d1 wide room"),
	"d1 one-button chest": AndSlot("d1 wide room"),
	"d1 boss": AndSlot("d1 wide room", "break bush safe", "d1 boss key",
		"kill pumpkin head"),

	// potentially 3 keys w/ vanilla route
	"d1 U-room": Or("d1 west terrace",
		And("d1 wide room", "break bush safe", "kill giant ghini",
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
		Or("feather", HardAnd("cane", Or("pegasus satchel", "bombs")))),
	"d2 basement drop": AndSlot("d2 basement", "feather"),
	"d2 basement chest": AndSlot("d2 basement", "feather", "bombs",
		"hit lever from minecart", "kill normal"),

	// 3 keys
	"d2 moblin platform": AndSlot("d2 basement", "feather",
		Count(3, "d2 small key")),
	// push moblin into doorway, stand on button, use switch hook
	"d2 statue puzzle": AndSlot("d2 moblin platform",
		Or("bracelet", "cane", HardAnd("switch hook", "push enemy"))),

	// 4 keys
	"d2 rope room": AndSlot("enter d2", "kill switch hook",
		Count(4, "d2 small key")),
	"d2 ladder chest": AndSlot("enter d2", Count(4, "d2 small key"), "bombs"),

	// 5 keys
	"d2 color room":   AndSlot("d2 statue puzzle", Count(5, "d2 small key")),
	"d2 boss":         AndSlot("d2 color room", "d2 boss key"),
	"head thwomp owl": And("mystery seeds", "d2 boss"),
}

var agesD3Nodes = map[string]*Node{
	// 0 keys
	"d3 pols voice chest": AndSlot("enter d3", "bombs"),
	"d3 1F spinner":       And("enter d3", Or("kill moldorm", "bracelet")),
	"d3 S crystal":        And("d3 1F spinner"),
	"four crystals owl":   And("mystery seeds", "d3 1F spinner"),
	"d3 E crystal":        And("d3 1F spinner", "bombs"),
	"d3 statue drop":      AndSlot("d3 E crystal"),

	// 1 key
	// you can clip into the blocks enough to hit this crystal with switch hook
	"d3 W crystal": And("enter d3", "d3 small key"),
	"d3 N crystal": And("d3 W crystal",
		Or("any seed shooter", "boomerang", Hard("switch hook"))),
	"stone soldiers owl": And("mystery seeds", "d3 small key"),
	"d3 armos drop":      AndSlot("d3 W crystal"),
	"d3 six-block drop":  AndSlot("d3 W crystal"),
	"break crystal switch": Or("sword", "switch hook", "boomerang",
		"ember satchel", "scent satchel", "mystery satchel",
		"any seed shooter", "punch object"),
	"d3 B1F spinner": And("d3 S crystal", "d3 E crystal", "d3 N crystal",
		"d3 W crystal", "break crystal switch"),
	"d3 crossroads":         AndSlot("d3 B1F spinner"),
	"d3 conveyor belt room": AndSlot("d3 W crystal"),
	"d3 torch chest": AndSlot("d3 B1F spinner",
		Or("ember shooter", Hard("mystery shooter"))),
	"d3 bridge chest": AndSlot("d3 W crystal",
		Or("any seed shooter", "jump 3",
			HardAnd("d3 post-subterror", Count(4, "d3 small key"), "feather"),
			HardAnd(Or("boomerang", And("bracelet", "toss ring")),
				Or("feather", "pegasus satchel")))),
	"d3 B1F east": AndSlot("d3 B1F spinner", "kill subterror",
		Or("any seed shooter", Hard("sword"))), // spin slash through corner
	// post-subterror and boss door do not reference each other
	"d3 post-subterror": Or(
		"d3 boss door",
		And("d3 B1F spinner", "kill subterror"),
		And("d3 bridge chest", Count(4, "d3 small key"),
			Or("jump 3", Hard("feather")))),
	"d3 boss door": Or(
		And("d3 post-subterror", Or("jump 3", Hard("feather")),
			Or("any seed shooter", "boomerang",
				HardAnd(
					Or("sword", And("bomb jump 2",
						Or("ember seeds", "scent seeds", "mystery seeds"))),
					Or("jump 3", "switch hook",
						And("bracelet", Count(4, "d3 small key")))))),
		And("d3 bridge chest", Count(4, "d3 small key"),
			Or("any seed shooter", "boomerang"))),
	"d3 moldorm drop": AndSlot("kill moldorm", "d3 post-subterror"),
	"d3 boss": AndSlot("d3 boss door", "d3 boss key",
		Or("ember shooter", "scent shooter", "ember satchel",
			Hard("scent satchel"))),

	// 3 keys
	"d3 bush beetle room": AndSlot("enter d3", "kill switch hook",
		Count(3, "d3 small key")),

	// 4 keys
	"d3 mimic room": AndSlot("d3 bush beetle room", "kill normal",
		Count(4, "d3 small key")),
}

var agesD4Nodes = map[string]*Node{
	// 0 keys
	"d4 first chest": AndSlot("enter d4", Or("kill normal", "push enemy"),
		Or("feather", "switch hook")),
	"d4 cube chest": AndSlot("d4 first chest", "feather"),

	// 1 key
	"d4 minecart A": And("enter d4", "feather", "d4 small key"),
	"d4 first crystal switch": AndSlot("d4 minecart A",
		Or("any seed shooter", Hard("boomerang"))),
	"d4 minecart chest": AndSlot("d4 minecart A", "hit lever"),

	// 2 keys
	"d4 minecart B": And("d4 minecart A", "hit lever", "bracelet",
		"kill normal", Count(2, "d4 small key")),
	"d4 second crystal switch": AndSlot("d4 minecart B",
		Or("any seed shooter", Hard("boomerang"))),

	// 3 keys
	"d4 minecart C": And("d4 minecart B", Count(3, "d4 small key")),
	"d4 color tile drop": AndSlot("d4 minecart C",
		Or("sword", "ember seeds", "scent shooter", "gale shooter",
			HardOr("scent satchel", And("peace ring", "bombs")))),

	// 4 keys
	"d4 minecart D": And("d4 minecart C", Count(4, "d4 small key")),
	// these weapons are for the miniboss, not the moldorms
	"d4 small floor puzzle": AndSlot("d4 minecart D", "bombs",
		Or("sword", "switch hook", "scent shooter", "punch enemy", Hard())),
	"d4 large floor puzzle": AndSlot("d4 minecart D", "switch hook"),
	"d4 boss": AndSlot("d4 large floor puzzle", "d4 boss key",
		Or("sword", "boomerang", "punch enemy")),

	// 5 keys
	"d4 lava pot chest": AndSlot("d4 large floor puzzle",
		Count(5, "d4 small key")),
}

// every chest not behind a key door in d5 requires you to be able to hit a
// switch, so that's a requirement for the first node.
var agesD5Nodes = map[string]*Node{
	// 0 keys
	"d5 switch A": And("enter d5", "kill normal",
		Or("hit switch", Hard("bracelet"))),
	"d5 blue peg chest": AndSlot("d5 switch A"),
	"d5 dark room": AndSlot("d5 switch A", "hit switch", // can't use pots here
		Or("cane", "switch hook", HardOr("kill normal", "push enemy"))),
	"d5 like-like chest": AndSlot("d5 switch A",
		Or("hit switch ranged", Hard("bracelet"), HardAnd("feather", "cane",
			Or("ember seeds", "scent seeds", "mystery seeds")))),
	"d5 eyes chest": AndSlot("d5 switch A", Or("any seed shooter",
		HardAnd("pegasus satchel", "feather", "mystery seeds",
			Or("hit switch ranged", And("bracelet", "toss ring"), "cane")))),
	"d5 two-statue puzzle": AndSlot("d5 switch A", "break pot", "cane",
		"feather", Or("any seed shooter", "boomerang", Hard("sword"),
			HardAnd("bomb jump 2",
				Or("ember seeds", "scent seeds", "mystery seeds")))),
	"d5 boss": AndSlot("d5 switch A", "d5 boss key", "cane", "sword"),

	// 2 keys
	"d5 crossroads": And("d5 switch A", "feather", "bracelet",
		Count(2, "d5 small key"),
		Or("cane", Hard("jump 3"), HardAnd("sword", "switch hook"))),
	"d5 diamond chest": AndSlot("d5 crossroads", "switch hook"),

	// 5 keys
	"d5 three-statue puzzle": AndSlot("d5 switch A", "cane",
		Count(5, "d5 small key")),
	"d5 six-statue puzzle": AndSlot("d5 switch A", "ember shooter", "feather",
		Count(5, "d5 small key")),
	"d5 red peg chest": AndSlot("d5 crossroads", "hit switch ranged",
		Count(5, "d5 small key")),
	"crown dungeon owl": And("mystery seeds", "d5 red peg chest"),
	"d5 owl puzzle":     AndSlot("d5 red peg chest"),
}

var agesD6Nodes = map[string]*Node{
	// past, 0 keys
	"d6 past color room": AndSlot(
		"enter d6 past", "feather", "kill switch hook"),
	"d6 past wizzrobe chest": AndSlot("enter d6 past", "bombs",
		"kill normal"),
	"d6 past pool chest": AndSlot("enter d6 past", "bombs", "ember seeds",
		"flippers"),
	"d6 open wall":    And("enter d6 past", "bombs", "ember shooter"),
	"deep waters owl": And("mystery seeds", "d6 open wall"),
	"d6 past stalfos chest": AndSlot("enter d6 past", "ember seeds",
		Or("kill normal ranged", "scent satchel", "feather", Hard())),
	"d6 past rope chest": AndSlot("d6 open wall", "mermaid suit"),

	// past, 1 key
	"d6 past spinner": And("enter d6 past", "cane", "bracelet", "feather",
		"d6 past small key", "bombs"),
	"d6 past spear chest": AndSlot("d6 past spinner", "mermaid suit"),
	"d6 past diamond chest": AndSlot("d6 past spinner", "mermaid suit",
		"switch hook"),

	// past, 3 keys
	"d6 boss": AndSlot("d6 past spinner", "d6 boss key", "any seed shooter",
		Count(3, "d6 past small key")),

	// present, 0 keys
	"d6 present diamond chest": AndSlot("enter d6 present", "switch hook"),
	"d6 present rope room": And("enter d6 present",
		Or("flippers", "feather", "switch hook"),
		Or("any seed shooter", "boomerang", "jump 3")),
	"scent seduction owl":   And("mystery seeds", "d6 present rope room"),
	"d6 present rope chest": AndSlot("d6 present rope room", "scent satchel"),
	"d6 present hand room": And("enter d6 present",
		Or("flippers", "feather", "switch hook"),
		Or("any seed shooter", "boomerang",
			And("jump 3", Or("sword", "switch hook", "ember seeds",
				"scent seeds", "mystery seeds", Hard("bombs"))))),
	"d6 present cube chest": AndSlot("d6 present hand room", "bombs",
		"switch hook", Or("feather", Hard())),
	"d6 present spinner chest": AndSlot("d6 past spinner",
		"d6 present hand room", Or("feather", "switch hook")),
	"d6 present beamos chest": AndSlot("enter d6 present", "d6 open wall",
		"feather", Or("flippers",
			And("switch hook", Count(2, "d6 present small key")))),

	// present, 3 keys
	"luck test owl": And("mystery seeds", "d6 present beamos chest",
		Count(3, "d6 present small key")),
	// only sustainable weapons count for killing the ropes
	"d6 present RNG chest": AndSlot("d6 present beamos chest", "bracelet",
		Or("sword", "cane", "switch hook", "punch enemy"),
		Count(3, "d6 present small key")),
	"d6 present channel chest": AndSlot("enter d6 present", "d6 open wall",
		"switch hook", Count(3, "d6 present small key")),
	"d6 present vire chest": AndSlot("d6 present spinner chest",
		Count(3, "d6 present small key"),
		Or("sword", "expert's ring", Hard()), "switch hook"),
}

// assume mermaid suit.
// leaving/entering the dungeon (but not loading a file) resets the water
// level. this is necessary to make keys work out, since otherwise you can
// drain the water level without getting enough keys to refill it! there just
// aren't enough chests otherwise.
var agesD7Nodes = map[string]*Node{
	// 0 keys
	"d7 spike chest": AndSlot("enter d7"),
	"d7 crab chest": AndSlot("enter d7",
		Or("kill underwater", And("drain d7", "kill normal"))),
	"d7 diamond puzzle": AndSlot("enter d7", "switch hook"),
	"d7 flower room":    AndSlot("enter d7", "long hook"),
	"golden isle owl":   And("mystery seeds", "enter d7"),
	"d7 stairway chest": AndSlot("enter d7",
		Or("long hook", And("refill d7", "cane"))),
	"d7 right wing": AndSlot("d7 stairway chest", "kill moldorm"),

	// 3 keys - enough to drain dungeon (and also refill 1F)
	"drain d7":               And("enter d7", Count(3, "d7 small key")),
	"refill d7":              And("drain d7", "switch hook"),
	"jabu switch room owl":   And("mystery seeds", "drain d7"),
	"d7 boxed chest":         AndSlot("drain d7"),
	"d7 pot island chest":    AndSlot("drain d7", "switch hook"),
	"d7 cane/diamond puzzle": AndSlot("drain d7", "long hook", "cane"),

	// 4 keys - enough to choose any water level
	"flood d7":       And("refill d7", "long hook", Count(4, "d7 small key")),
	"d7 3F terrace":  AndSlot("flood d7"),
	"d7 left wing":   AndSlot("flood d7"),
	"plasmarine owl": And("mystery shooter", "flood d7"),
	"d7 boss":        AndSlot("d7 boss key", "flood d7"),

	// 5 keys
	"d7 hallway chest": AndSlot("drain d7", "long hook",
		Count(5, "d7 small key")),

	// 7 keys
	"d7 miniboss chest": AndSlot("d7 stairway chest", "feather",
		Or("sword", "boomerang", "scent shooter"), Count(7, "d7 small key")),
	"d7 post-hallway chest": AndSlot("flood d7", Count(7, "d7 small key")),
}

var agesD8Nodes = map[string]*Node{
	// 0 keys
	"open ears owl": And("mystery seeds", "enter d8"),
	"d8 1F chest":   AndSlot("enter d8", "bombs"),

	// 1 key - access B1F
	"d8 ghini chest": AndSlot("enter d8", "d8 small key", "switch hook",
		"cane", "seed shooter", Or("ember seeds", Hard("mystery seeds"))),
	"d8 B1F NW chest": AndSlot("d8 ghini chest"),

	// 2 keys - access SE spinner
	"d8 blue peg chest":    AndSlot("d8 ghini chest", Count(2, "d8 small key")),
	"d8 blade trap chest":  AndSlot("d8 blue peg chest"),
	"d8 sarcophagus chest": AndSlot("d8 blue peg chest", "power glove"),
	"d8 stalfos":           AndSlot("d8 blue peg chest"),

	// 4 keys - reach miniboss
	"d8 maze chest": AndSlot("d8 blue peg chest", "sword",
		Count(4, "d8 small key")),
	"d8 NW slate chest": AndSlot("d8 maze chest"),
	"d8 NE slate chest": AndSlot("d8 maze chest", "feather", "flippers",
		"ember seeds"),
	"ancient words owl": And("mystery shooter", "d8 NE slate chest"),
	"d8 B3F chest":      AndSlot("d8 maze chest", "power glove"),
	"d8 tile room":      AndSlot("d8 maze chest", "feather"),
	"d8 SE slate chest": AndSlot("d8 tile room"),
	"d8 boss": AndSlot("d8 boss key", "d8 tile room",
		Count(4, "slate")),

	// 5 keys
	"d8 floor puzzle": AndSlot("d8 blue peg chest", Count(5, "d8 small key")),
	"d8 SW slate chest": AndSlot("d8 maze chest", "bracelet",
		Count(5, "d8 small key")),
}

var agesD9Nodes = map[string]*Node{
	"black tower owl": And("mystery seeds", "maku seed"),
	"done": AndStep("maku seed", "mystery seeds", "switch hook",
		Or("sword", "punch enemy"), "bombs"), // bombs in case of spider form
}
