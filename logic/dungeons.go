package logic

// keep small keys and their chests separate, so that they can be changed into
// slots if small keys are ever randomized.
//
// dungeons should rely on overworld information as little as possible.
// ideally "enter <dungeon>" is the only overworld item the dungeon nodes
// reference (and that node should not be defined here).

var d0Nodes = map[string]*Node{
	"d0 key chest":   And("enter d0"),
	"d0 sword chest": AndSlot("enter d0", "d0 small key"),
	"d0 rupee chest": OrSlot("remove bush safe", "flute", "dimitri to west"),

	"d0 small key": And("d0 key chest"),
}

var d1Nodes = map[string]*Node{
	"d1 key fall":       And("enter d1", "kill stalfos"),
	"d1 map chest":      AndSlot("d1 key A", "kill stalfos"),
	"d1 compass chest":  AndSlot("d1 map chest"),
	"d1 gasha chest":    AndSlot("d1 map chest", "kill goriya"),
	"d1 bomb chest":     AndSlot("d1 map chest", "hit lever"),
	"d1 key chest":      And("d1 bomb chest"),
	"enter goriya bros": And("d1 bomb chest", "bombs", "d1 key B"),
	"d1 satchel spot":   AndSlot("enter goriya bros", "kill goriya bros"),
	"d1 boss key chest": AndSlot("d1 map chest",
		Or("ember seeds", Hard("mystery seeds")), "kill goriya (pit)"),
	"d1 ring chest": AndSlot("enter d1",
		Or("ember seeds", Hard("mystery seeds"))),
	"d1 essence": AndStep("d1 ring chest", "d1 boss key",
		"kill aquamentus"),

	"d1 key A": And("d1 key fall"),
	"d1 key B": And("d1 key chest"),
}

var d2Nodes = map[string]*Node{
	"d2 5-rupee chest": AndSlot("d2 torch room"),
	"d2 rope room":     And("d2 torch room", "kill rope"),
	"d2 arrow room": Or("enter d2 B",
		And("d2 torch room", Or("ember seeds", Hard("mystery seeds")))),
	"d2 rupee room":     And("d2 arrow room", "bombs"),
	"d2 hardhat room":   And("d2 arrow room", "d2 2 keys"), // min. 1 key
	"d2 map chest":      AndSlot("d2 hardhat room", "remove pot"),
	"d2 compass chest":  AndSlot("d2 arrow room", "kill normal"),
	"d2 bracelet room":  And("d2 hardhat room", "kill hardhat (pit)"),
	"d2 bracelet chest": AndSlot("d2 bracelet room", "kill moblin (gap)"),
	"d2 spiral chest":   And("enter d2 B", "bombs"),
	"d2 blade chest":    Or("enter d2 B", And("d2 arrow room", "kill normal")),

	// from here on it's entirely linear.
	"d2 10-rupee chest": AndSlot("d2 bomb wall", "bombs", "bracelet"),
	"d2 spinner":        And("d2 10-rupee chest", "d2 2 keys"), // min. 1 key
	"d2 boss key chest": AndSlot("d2 spinner", "d2 3 keys"),
	"d2 essence":        And("d2 spinner", "d2 boss key"),

	"d2 key A": And("d2 rope room"),
	"d2 key B": And("d2 spiral chest"),
	"d2 key C": And("d2 blade chest"),
	"d2 2 keys": Or(And("d2 key A", "d2 key B"), And("d2 key A", "d2 key C"),
		And("d2 key B", "d2 key C")),
	"d2 3 keys": And("d2 key A", "d2 key B", "d2 key C"),

	"d2 torch room": Or("enter d2 A", "d2 compass chest"),
	"d2 bomb wall":  And("d2 blade chest"), // alias for external reference
}

var d3Nodes = map[string]*Node{
	// first floor
	"d3 center":       And("enter d3", "kill spiked beetle"),
	"d3 mimic stairs": Or("d3 rupee chest", And("d3 center", "bracelet")),
	"d3 roller chest": And("d3 mimic stairs", "bracelet"),
	"d3 rupee chest":  OrSlot("d3 mimic stairs", And("d3 center", "jump 2")),
	"d3 gasha chest":  AndSlot("d3 mimic stairs", "jump 2"),
	"d3 omuai stairs": And("d3 mimic stairs", "jump 2", "kill omuai",
		"d3 2 keys"), // min. 1 key
	"d3 boss key chest": AndSlot("d3 omuai stairs"),

	// second floor
	"d3 bomb chest": AndSlot("d3 mimic stairs"),
	"d3 map chest":  AndSlot("d3 bomb chest", "bombs"),
	"d3 feather chest": AndSlot("d3 rupee chest", "kill mimic",
		"d3 2 keys"), // min. 1 key
	"d3 trampoline chest": And("d3 center", "jump 2"),
	"d3 compass chest":    AndSlot("d3 center", "jump 2"),
	"enter mothula":       And("d3 omuai stairs", "d3 boss key"),
	"d3 essence":          AndStep("enter mothula", "kill mothula"),

	// fixed items
	"d3 key A":  And("d3 roller chest"),
	"d3 key B":  And("d3 trampoline chest"),
	"d3 2 keys": And("d3 key A", "d3 key B"),
}

var d4Nodes = map[string]*Node{
	// left branch from entrance
	"d4 bomb chest": AndSlot("enter d4", Or("flippers", "jump 6")),
	"d4 pot room":   And("d4 bomb chest", "bombs", "bracelet"),
	"d4 map chest":  AndSlot("d4 bomb chest", "hit lever"),
	"d4 dark chest": And("d4 map chest", "jump 2"),

	// 2F (ground floor), right branch
	"d4 compass chest": AndSlot("enter d4", Or("flippers", "jump 6"), "bombs",
		"d4 1 key"), // might be possible w/ cape + bomb boost?
	"d4 roller minecart": And("enter d4", "flippers", "jump 2", "d4 1 key"),
	"d4 water key room": And("d4 roller minecart", "hit lever", "flippers",
		"kill normal"),
	"d4 stalfos stairs": And("d4 roller minecart", "d4 2 keys", "kill stalfos"),

	// 1F
	"d4 pre-mid chest":   And("d4 stalfos stairs"),
	"d4 final minecart":  And("d4 stalfos stairs", "kill agunima"),
	"d4 torch chest":     And("d4 stalfos stairs", "ember slingshot"),
	"d4 slingshot chest": AndSlot("d4 final minecart", "d4 5 keys"), // min. 2 keys
	"d4 boss key spot": AndSlot("d4 final minecart", "hit very far lever",
		"d4 5 keys"), // min. 2 keys
	"d4 basement stairs": And("d4 final minecart", "hit far lever",
		"d4 5 keys"), // min. 2 keys

	// B1F
	"enter gohma": And("d4 basement stairs", "d4 boss key",
		Or("ember slingshot", Hard("mystery slingshot"), "jump 3")),
	"d4 essence": AndStep("enter gohma", "kill gohma"),

	// fixed items
	"d4 key A": And("d4 pot room"),
	"d4 key B": And("d4 dark chest"),
	"d4 key C": And("d4 water key room"),
	"d4 key D": And("d4 pre-mid chest"),
	"d4 key E": And("d4 torch chest"),
	"d4 1 key": Or("d4 key A", "d4 key B"),
	"d4 2 keys": Or(And("d4 key A", "d4 key B"), And("d4 key A", "d4 key C"),
		And("d4 key B", "d4 key C")),
	"d4 5 keys": And("d4 key A", "d4 key B", "d4 key C", "d4 key D",
		"d4 key E"),

	"enter agunima": And("d4 pre-mid chest"), // alias for external reference
}

// the keys in this dungeon suck, so i'm not even going to bother with "hard"
// logic for them.
var d5Nodes = map[string]*Node{
	// 1F (it's the only F)
	"d5 cart bay":   And("enter d5", Or("flippers", "jump 3")),
	"d5 cart chest": And("d5 cart bay", "hit lever"),
	"d5 pot room": And("enter d5", Or("magnet gloves", "bombs", "jump 2"),
		And("d5 cart bay", Or("jump 2", Hard("pegasus satchel")))),
	"d5 map chest": AndSlot("d5 pot room", "kill gibdo", "kill zol"),
	"d5 magnet gloves chest": AndSlot("d5 pot room", Or("flippers", "jump 6"),
		"d5 5 keys"),
	"d5 left chest": And("enter d5", Or("magnet gloves", "jump 4")),
	"d5 rupee chest": AndSlot("enter d5", Or("magnet gloves",
		And("d5 cart bay", Or("jump 2", Hard("pegasus satchel")), "bombs"))),
	"d5 compass chest": AndSlot("enter d5", "kill moldorm", "kill iron mask"),
	"d5 armos chest": And("d5 rupee chest", "kill moldorm", "kill iron mask",
		"kill armos"),
	"d5 spinner chest": And("d5 cart bay", Or("magnet gloves", "jump 6")),
	"d5 drop ball":     And("d5 cart bay", "hit lever", "kill darknut (pit)"),
	"d5 pre-mid chest": And("d5 cart bay", Or("magnet gloves", "jump 4")),
	"d5 post-syger":    And("d5 pre-mid chest", "kill syger"), // keys after
	"d5 boss key spot": AndSlot("d5 drop ball", "d5 post-syger",
		"magnet gloves", Or("kill magunesu", Hard("jump 2")), "d5 5 keys"),
	"d5 essence": AndStep("d5 post-syger", "magnet gloves",
		Or("jump 2", Hard("start")), "d5 boss key", "d5 5 keys"),

	// fixed items
	"d5 key A": And("d5 cart chest"),
	"d5 key B": And("d5 left chest"),
	"d5 key C": And("d5 armos chest"),
	"d5 key D": And("d5 spinner chest"),
	"d5 key E": And("d5 pre-mid chest"),
	"d5 5 keys": And("d5 key A", "d5 key B", "d5 key C", "d5 key D",
		"d5 key E"),
}

// all the rupee chests in this dungeon are trivial, so i'm ignoring which is
// which and just labeling them by letter.
var d6Nodes = map[string]*Node{
	// 1F
	"d6 rupee chest A": AndSlot("enter d6"),
	"d6 rupee room":    And("enter d6", "bombs"),
	"d6 magkey room": And("enter d6",
		Or(And("magnet gloves", "jump 2"), "jump 4")),
	"d6 compass chest": AndSlot("enter d6", "d6 3 keys"), // min. 1 key
	"d6 map chest":     AndSlot("enter d6"),
	"d6 rupee chest C": AndSlot("enter d6"),
	"d6 U-room":        And("enter d6", "break crystal", "boomerang L-2"),
	"d6 torch stairs":  And("d6 U-room", "ember seeds"),

	// 2F
	"d6 skipped chest":   And("enter d6", "magnet gloves", "break crystal"),
	"d6 bomb chest":      AndSlot("d6 compass chest"),
	"d6 rupee chest B":   AndSlot("d6 bomb chest", "bombs"),
	"d6 boomerang chest": AndSlot("d6 rupee chest B"),
	"d6 boss key chest": AndSlot("d6 torch stairs",
		"pegasus satchel", "jump 2"),

	// 3F
	"d6 vire chest": And("d6 boss key chest", "kill stalfos"),
	"enter vire":    And("d6 vire chest", "d6 3 keys"), // min. 1 key

	// 5F
	"d6 pre-boss room": And("enter vire", "kill vire", "kill hardhat (magnet)"),
	"d6 essence": AndStep("d6 pre-boss room", "d6 boss key",
		"kill manhandla"),

	// fixed items
	"d6 key A":  And("d6 magkey room"),
	"d6 key B":  And("d6 vire chest"),
	"d6 key C":  And("d6 skipped chest"),
	"d6 3 keys": And("d6 key A", "d6 key B", "d6 key C"),
}

// this does *not* account for poe skip.
var d7Nodes = map[string]*Node{
	// 1F
	"d7 wizzrobe chest": And("enter d7", "kill wizzrobe"),
	"d7 ring chest":     AndSlot("enter d7", "d7 key A"),
	"enter poe A": And("d7 ring chest",
		Or("ember slingshot", Hard("mystery slingshot"))),
	"d7 compass chest": AndSlot("enter d7", "bombs"),
	"d7 map chest":     AndSlot("d7 pot room", "jump 2", "d7 key B"),

	// B1F
	"d7 pot room": And("enter d7", "bracelet", "enter poe A",
		"kill poe sister"),
	"d7 zol button": And("d7 pot room", "jump 2"),
	"d7 magunesu chest": And("d7 armos puzzle", "jump 3", "kill magunesu",
		"magnet gloves"),
	"enter poe B": And("d7 pot room", "d7 3 keys", "ember seeds",
		Or("pegasus satchel", "slingshot L-2", Hard("start"))),
	"d7 water stairs": And("enter poe B", "flippers"),
	"d7 cape chest":   AndSlot("d7 water stairs", "d7 cross bridge"),

	// B2F
	"d7 armos puzzle": And("d7 pot room", Or("jump 3", "magnet gloves")),
	"d7 cross bridge": Or("jump 4", "kill darknut (across pit)",
		And("jump 2", "magnet gloves")),
	"d7 moldorm room":   And("d7 water stairs", "jump 3", "d7 4 keys"),
	"d7 rupee chest":    AndSlot("d7 moldorm room", "kill moldorm"),
	"d7 skipped room":   And("d7 rupee chest"),
	"d7 boss key chest": AndSlot("d7 rupee chest", "d7 key E"),
	"d7 essence": AndStep("d7 rupee chest", "d7 boss key",
		"kill gleeok"),

	// fixed items
	"d7 key A": And("d7 wizzrobe chest"),
	"d7 key B": And("d7 zol button"),
	"d7 key C": And("d7 armos puzzle"),
	"d7 key D": And("d7 magunesu chest"),
	"d7 key E": And("d7 skipped room"),
	"d7 3 keys": Or(
		And("d7 key B", Or("d7 key C", "d7 key D")),
		And("d7 key C", "d7 key D")),
	"d7 4 keys": And("d7 key A", "d7 key B", "d7 key C", "d7 key D"),
}

// this does *not* account for HSS skip.
var d8Nodes = map[string]*Node{
	// 1F
	"d8 eye room":     And("enter d8", "any slingshot", "remove pot"),
	"d8 ring chest":   AndSlot("enter d8", "any slingshot L-2", "jump 2"),
	"d8 hardhat room": And("enter d8", "kill magunesu"),
	"d8 hardhat key":  And("d8 hardhat room", "kill hardhat (magnet)"),
	"d8 compass chest": AndSlot("d8 hardhat room", "d8 1 key",
		Or("jump 4", Hard("jump 3"))),
	"d8 map chest": AndSlot("d8 spinner"),
	"d8 bomb chest": And("d8 HSS chest", "any slingshot L-2", "bombs",
		"kill darknut"),
	"d8 ice puzzle room": And("d8 HSS chest", "kill frypolar", "ember seeds",
		"slingshot L-2"),
	"d8 boss key chest": AndSlot("d8 ice puzzle room",
		Or("jump 6", "boomerang L-2", Hard("start"))),
	"d8 crystal room": And("d8 ice puzzle room", "d8 4 keys"),
	"d8 ghost armos":  And("d8 crystal room"),
	"d8 NW crystal":   And("d8 crystal room", "d8 7 keys"),
	"d8 NE crystal":   And("d8 crystal room", "hit lever"),
	"d8 SE crystal":   And("d8 crystal room"),
	"d8 SW crystal":   And("d8 crystal room", "d8 7 keys"),
	"d8 pot chest":    And("d8 SE crystal", "d8 NE crystal", "remove pot"),

	// B1F
	"d8 spinner":       And("d8 compass chest", "d8 2 keys"),
	"d8 HSS chest":     AndSlot("d8 spinner", "magnet gloves"),
	"d8 spinner chest": And("d8 HSS chest"),
	"d8 SE lava chest": And("d8 SE crystal"),
	"d8 SW lava chest": AndSlot("d8 SE crystal"),
	"d8 essence": AndStep("d8 SW crystal", "d8 SE crystal", "d8 NW crystal",
		"d8 7 keys", "d8 boss key", "kill medusa head"),

	// fixed items
	"d8 key A":  And("d8 eye room"),
	"d8 key B":  And("d8 hardhat key"),
	"d8 key C":  And("d8 spinner chest"),
	"d8 key D":  And("d8 bomb chest"),
	"d8 key E":  And("d8 ghost armos"),
	"d8 key F":  And("d8 SE lava chest"),
	"d8 key G":  And("d8 pot chest"),
	"d8 1 key":  Or("d8 key A", "d8 key B"),
	"d8 2 keys": And("d8 key A", "d8 key B"),
	"d8 4 keys": And("d8 key C", "d8 key D"),
	"d8 7 keys": And("d8 key E", "d8 key F", "d8 key G"),
}

// onox's castle
var d9Nodes = map[string]*Node{
	"enter onox": And("enter d9", "kill wizzrobe", "kill floormaster", "kill darknut", "kill facade"),
	"done":       AndStep("enter onox", "kill onox"),
}
