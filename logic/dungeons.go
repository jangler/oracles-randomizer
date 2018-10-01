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
	"d0 rupee chest": OrSlot("remove bush safe", "flute"),

	"d0 small key": And("d0 key chest"),
}

var d1Nodes = map[string]*Node{
	"d1 key fall":       And("enter d1", "kill stalfos (throw)"),
	"d1 map chest":      AndSlot("d1 key A", "kill stalfos"),
	"d1 compass chest":  AndSlot("d1 map chest"),
	"d1 gasha chest":    AndSlot("d1 map chest", "kill goriya"),
	"d1 bomb chest":     AndSlot("d1 map chest", "hit lever"),
	"d1 key chest":      And("d1 bomb chest"),
	"enter goriya bros": And("d1 bomb chest", "bombs", "d1 key B"),
	"d1 satchel spot":   AndSlot("enter goriya bros", "kill goriya bros"),
	"d1 boss key chest": AndSlot("d1 map chest",
		Or("ember seeds", "mystery seeds"), "kill goriya (pit)"),
	"d1 ring chest": AndSlot("enter d1", Or("ember seeds", "mystery seeds")),
	"d1 essence": AndStep("d1 ring chest", "d1 boss key",
		"kill aquamentus"),

	"d1 key A": And("d1 key fall"),
	"d1 key B": And("d1 key chest"),
}

var d2Nodes = map[string]*Node{
	"d2 5-rupee chest": AndSlot("d2 torch room"),
	"d2 rope room":     And("d2 torch room", "kill rope"),
	"d2 arrow room": Or("enter d2 B",
		And("d2 torch room", Or("ember seeds", "mystery seeds"))),
	"d2 rupee room":     And("d2 arrow room", "bombs"),
	"d2 hardhat room":   And("d2 arrow room", Or("d2 2 keys", Hard("d2 1 key"))),
	"d2 map chest":      AndSlot("d2 hardhat room", "remove pot"),
	"d2 compass chest":  AndSlot("d2 arrow room", "kill normal"),
	"d2 bracelet room":  And("d2 hardhat room", "kill hardhat (pit, throw)"),
	"d2 bracelet chest": AndSlot("d2 bracelet room", "kill moblin (gap, throw)"),
	"d2 spiral chest":   And("enter d2 B", "bombs"),
	"d2 blade chest":    Or("enter d2 B", And("d2 arrow room", "kill normal")),

	// from here on it's entirely linear.
	"d2 10-rupee chest": AndSlot("d2 bomb wall", "bombs", "bracelet"),
	"d2 spinner": And("d2 10-rupee chest",
		Or("d2 2 keys", Hard("d2 1 key"))),
	"d2 boss key chest": AndSlot("d2 spinner", "d2 3 keys"),
	"d2 essence":        And("d2 spinner", "d2 boss key"),

	"d2 key A": And("d2 rope room"),
	"d2 key B": And("d2 spiral chest"),
	"d2 key C": And("d2 blade chest"),
	"d2 1 key": Or("d2 key A", "d2 key B", "d2 key C"),
	"d2 2 keys": Or(And("d2 key A", "d2 key B"), And("d2 key A", "d2 key C"),
		And("d2 key B", "d2 key C")),
	"d2 3 keys": And("d2 key A", "d2 key B", "d2 key C"),

	"d2 torch room": Or("enter d2 A", "d2 compass chest"),
	"d2 bomb wall":  And("d2 blade chest"), // alias for external reference
}

var d3Nodes = map[string]*Node{
	// first floor
	"d3 center":       And("enter d3", "kill spiked beetle (throw)"),
	"d3 mimic stairs": Or("d3 rupee chest", And("d3 center", "bracelet")),
	"d3 roller chest": And("d3 mimic stairs", "bracelet"),
	"d3 rupee chest":  OrSlot("d3 mimic stairs", And("d3 center", "jump")),
	"d3 gasha chest":  AndSlot("d3 mimic stairs", "jump"),
	"d3 omuai stairs": And("d3 mimic stairs", "jump",
		Or("d3 2 keys", Hard("d3 1 key")), "kill omuai"),
	"d3 boss key chest": AndSlot("d3 omuai stairs"),

	// second floor
	"d3 bomb chest": AndSlot("d3 mimic stairs"),
	"d3 map chest":  AndSlot("d3 bomb chest", "bombs"),
	"d3 feather chest": AndSlot("d3 rupee chest",
		Or("d3 2 keys", Hard("d3 1 key")), "kill mimic"),
	"d3 trampoline chest": And("d3 center", "jump"),
	"d3 compass chest":    AndSlot("d3 center", "jump"),
	"enter mothula":       And("d3 omuai stairs", "d3 boss key"),
	"d3 essence":          AndStep("enter mothula", "kill mothula"),

	// fixed items
	"d3 key A":  And("d3 roller chest"),
	"d3 key B":  And("d3 trampoline chest"),
	"d3 1 key":  Or("d3 key A", "d3 key B"),
	"d3 2 keys": And("d3 key A", "d3 key B"),
}

var d4Nodes = map[string]*Node{
	// left branch from entrance
	"d4 bomb chest": AndSlot("enter d4", "cross large pool"),
	"d4 pot room":   And("d4 bomb chest", "bombs", "bracelet"),
	"d4 map chest":  AndSlot("d4 bomb chest", "hit lever"),
	"d4 dark chest": And("d4 map chest", "jump"),

	// 2F (ground floor), right branch
	"d4 compass chest": AndSlot("enter d4", "cross large pool", "bombs",
		"d4 1 key"),
	"d4 roller minecart": And("enter d4", "flippers", "jump", "d4 1 key"),
	"d4 water key key room": And("d4 roller minecart", "hit lever", "flippers",
		Or("kill normal", "bracelet")),
	"d4 stalfos stairs": And("d4 roller minecart", "d4 2 keys",
		Or("kill stalfos", "bracelet")),

	// 1F
	"d4 pre-mid chest":  And("d4 stalfos stairs"),
	"d4 final minecart": And("d4 stalfos stairs", "kill agunima"),
	"d4 torch chest":    And("d4 stalfos stairs", "ember slingshot"),
	"d4 slingshot chest": AndSlot("d4 final minecart",
		Or("d4 5 keys", Hard("d4 2 keys"))),
	"d4 boss key spot": AndSlot("d4 final minecart", "hit very far lever",
		Or("d4 5 keys", Hard("d4 2 keys"))),
	"d4 basement stairs": And("d4 final minecart", "hit far lever",
		Or("d4 5 keys", Hard("d4 2 keys"))),

	// B1F
	"enter gohma": And("d4 basement stairs", "d4 boss key",
		Or("ember slingshot", "mystery slingshot", "long jump")),
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
	// general
	"cross magnet gap":   Or("pegasus jump L-2", "magnet gloves"),
	"magnet jump":        And("jump", "magnet gloves"),
	"sidescroll magnets": Or("magnet jump", "pegasus jump L-2"),

	// 1F (it's the only F)
	"d5 cart bay":   And("enter d5", Or("flippers", "long jump")),
	"d5 cart chest": And("d5 cart bay", "hit lever"),
	"d5 pot room": And("enter d5", "jump",
		Or("d5 cart bay", And("magnet gloves", "bombs"))),
	"d5 map chest": AndSlot("d5 pot room", "kill gibdo", "kill zol"),
	"d5 magnet gloves chest": AndSlot("d5 pot room", "cross large pool",
		"d5 5 keys"),
	"d5 left chest": And("enter d5", "cross magnet gap"),
	"d5 rupee chest": AndSlot("enter d5",
		Or("magnet gloves", And("d5 cart bay", "jump", "bombs"))),
	"d5 compass chest": AndSlot("enter d5", "kill moldorm", "kill iron mask"),
	"d5 armos chest": And("d5 rupee chest", "kill moldorm", "kill iron mask",
		"kill armos"),
	"d5 spinner chest": And("d5 cart bay", "cross magnet gap"),
	"d5 drop ball":     And("d5 cart bay", "hit lever", "kill darknut (pit)"),
	"d5 pre-mid chest": And("d5 cart bay", Or("magnet gloves", "feather L-2")),
	"d5 post-syger":    And("d5 pre-mid chest", "kill syger"), // keys after
	"d5 boss key spot": AndSlot("d5 drop ball", "d5 post-syger",
		"magnet gloves", Or("long jump", "kill magunesu"), "d5 5 keys",
		"sidescroll magnets"),
	"d5 essence": AndStep("d5 post-syger", "magnet gloves", "d5 boss key",
		"d5 5 keys"),

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
	"d6 spinner":         And("enter d6"),
	"d6 rupee chest A":   AndSlot("d6 spinner"),
	"d6 rupee room":      And("d6 spinner", "bombs"),
	"d6 magkey ball":     And("d6 spinner", "magnet gloves", "jump"),
	"d6 magkey jump":     And("pegasus jump L-2"),
	"d6 magnet key fall": Or("d6 magkey ball", "d6 magkey jump"),
	"d6 compass chest":   AndSlot("d6 spinner", "d6 key A"),
	"d6 crumble stairs":  And("d6 spinner", "d6 key A", "long jump"),
	"d6 key skip":        And("d6 armos room", "jump", "break crystal"),
	"d6 map chest":       OrSlot("d6 key skip", "d6 spinner"),
	"d6 rupee chest C":   AndSlot("d6 map chest"),
	"avoid traps":        Or("pegasus satchel", "jump"),
	"d6 switch stairs":   And("d6 map chest", "break crystal", "avoid traps", "boomerang L-2"),
	"d6 U-room":          And("d6 cracked room", "boomerang L-2"),
	"d6 torch stairs":    And("d6 U-room", "ember seeds"),

	// 2F
	"d6 skipped key chest": And("d6 spinner", "magnet gloves", "break crystal", "jump"), // being nice
	"d6 bomb chest":        AndSlot("d6 crumble stairs"),
	"d6 rupee chest B":     AndSlot("d6 armos room"),
	"d6 armos room":        And("d6 crumble stairs", "bombs"),
	"d6 boomerang chest":   AndSlot("d6 armos room", "jump"),
	"d6 cracked room":      And("d6 switch stairs"),
	"d6 boss key chest":    AndSlot("d6 torch stairs", "long jump"),
	"d6 gauntlet stairs":   And("d6 boss key chest"),

	// 3F
	"d6 vire key chest": And("d6 gauntlet stairs", "kill stalfos", "jump"),
	"enter vire":        And("d6 gauntlet stairs", "kill stalfos", "d6 key B"),
	"d6 rng stairs":     And("enter vire", "kill vire"),

	// 4F
	"d6 3-switch room": And("d6 rng stairs", "kill hardhat (magnet)"),

	// 5F
	"d6 pre-boss room": And("d6 3-switch room", "hit very far switch"),
	"enter manhandla":  And("d6 pre-boss room", "jump", "hit far switch", "d6 boss key"),
	"d6 essence":       AndStep("enter manhandla", "kill manhandla"),

	// fixed items
	"d6 key A": And("d6 magnet key fall"),
	"d6 key B": And("d6 vire key chest"),
	"d6 key C": And("d6 skipped key chest"),
}

var d7Nodes = map[string]*Node{
	// 1F
	"d7 wizzrobe key chest": And("enter d7", "kill wizzrobe"),
	"d7 ring chest":         AndSlot("enter d7", "d7 key A"),
	"enter poe A": And("d7 ring chest",
		Or("ember slingshot", "mystery slingshot")),
	"d7 compass chest": AndSlot("enter d7", "bombs"),
	"d7 map chest": AndSlot("d7 pot room", "jump",
		Or("d7 key B", HardAnd("d7 key A", "poe skip"))),
	"poe skip": HardAnd("enter d7", "bombs", "bracelet", "feather", "pegasus satchel"),

	// B1F
	"d7 armos room": And("enter d7", "bracelet",
		Or(And("enter poe A", "kill poe sister"), Hard("poe skip"))),
	"d7 zol key fall":       And("d7 armos room", "jump"),
	"d7 pot room":           And("d7 armos room"),
	"d7 magunesu key chest": And("d7 magunesu room", "kill magunesu", "jump", "magnet gloves"),
	"enter poe B": And("d7 pot room", "d7 key B",
		Or("d7 key C", HardAnd("d7 key A", "poe skip"))),
	"d7 water stairs": And("enter poe B", "ember seeds", "kill poe sister", "flippers"),
	"d7 cape chest":   AndSlot("d7 trampoline pair", "jump", "kill stalfos (pit)"),

	// B2F
	"d7 fool's gap":     Or("long jump", "magnet gloves"),
	"d7 armos puzzle":   And("d7 pot room", "kill keese", "d7 fool's gap"), // being nice
	"d7 armos key fall": And("d7 armos puzzle"),
	"d7 magunesu room":  And("d7 armos puzzle", "long jump"),
	"d7 cross bridge": Or("feather L-2", "kill darknut (across pit)",
		And("jump", "magnet gloves")),
	"d7 trampoline pair": And("d7 water stairs", "d7 cross bridge"),
	"d7 moldorm room": And("d7 water stairs", "feather L-2",
		Or("d7 key D", HardAnd("d7 key C", "poe skip"))),
	"enter poe sisters": And("d7 moldorm room", "kill moldorm", "feather L-2"),
	"d7 stairs room":    And("enter poe sisters", "kill poe sister"),
	"d7 rupee chest":    AndSlot("d7 stairs room"),
	"d7 enter skipped": And("d7 stairs room", Or(
		And("magnet gloves", "jump"), HardAnd("pegasus jump L-2"))),
	"d7 skipped key poof": And("d7 enter skipped", "kill wizzrobe (pit)", "kill stalfos (pit)"),
	"d7 boss key chest": AndSlot("d7 stairs room", "jump", "hit switch",
		"kill stalfos", Or("d7 key E",
			HardAnd("poe skip", "d7 key D", "d7 enter skipped"))),
	"enter gleeok": And("d7 stairs room", "d7 boss key"),
	"d7 essence":   AndStep("enter gleeok", "kill gleeok"),

	// fixed items
	"d7 key A": And("d7 wizzrobe key chest"),
	"d7 key B": And("d7 zol key fall"),
	"d7 key C": And("d7 armos key fall"),
	"d7 key D": And("d7 magunesu key chest"),
	"d7 key E": And("d7 skipped key poof"),
}

// this does *not* account for HSS skip.
var d8Nodes = map[string]*Node{
	// 1F
	"d8 eye room":      And("enter d8", "any slingshot", "remove pot"),
	"d8 ring chest":    AndSlot("enter d8", "any slingshot L-2", "jump"),
	"d8 hardhat room":  And("enter d8", "kill magunesu"),
	"d8 hardhat key":   And("d8 hardhat room", "kill hardhat (magnet)"),
	"d8 compass chest": AndSlot("d8 hardhat room", "d8 1 key", "long jump"),
	"d8 map chest":     AndSlot("d8 spinner"),
	"d8 bomb chest": And("d8 HSS stairs", "any slingshot L-2", "bombs",
		"kill darknut"),
	"d8 ice puzzle room": And("d8 HSS stairs", "kill frypolar", "ember seeds",
		"slingshot L-2"),
	"d8 boss key chest": AndSlot("d8 ice puzzle room",
		Or("feather L-2", Hard("long jump")),
		Or("pegasus jump L-2", "boomerang", "bombs", "any slingshot",
			HardAnd("feather L-2", Or("sword", "fool's ore")))),
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
