package logic

// keep small keys and their chests separate, so that they can be changed into
// slots if small keys are ever randomized.
//
// dungeons should rely on overworld information as little as possible.
// ideally "enter <dungeon>" is the only overworld item the dungeon nodes
// reference (and that node should not be defined here).
//
// bush- and pot-throwing is in hard logic, but with an arbitrary limit of
// three screen transitions per carry, and no more than two enemies can be
// required to be killed with one throw.

var seasonsD0Nodes = map[string]*Node{
	"d0 key chest":   And("enter d0"),
	"d0 sword chest": AndSlot("enter d0", "d0 small key"),
	"d0 rupee chest": OrSlot("remove bush safe", "flute"),

	"d0 small key": And("d0 key chest"),
}

// bush-throwing is in hard logic for a few rooms, though the first stalfos one
// doesn't matter, the goriya one only matters if you killed the stalfos with
// rod, and the lever one only matters if you killed the stalfos with bombs.
// bush-throwing is *not* in logic for the vanilla BK room, since you need to
// relight the torches every time you leave.
var seasonsD1Nodes = map[string]*Node{
	"d1 key fall": And("enter d1",
		Or("kill stalfos", Hard("bracelet"))),
	"d1 stalfos chest": AndSlot("d1 key A", "kill stalfos"),
	"d1 lever room":    AndSlot("d1 stalfos chest"),
	"d1 block-pushing room": AndSlot("d1 stalfos chest",
		Or("kill goriya", Hard("bracelet"))),
	"d1 railway chest": AndSlot("d1 stalfos chest",
		Or("hit lever", Hard("bracelet"))),
	"d1 key chest":      And("d1 railway chest"),
	"enter goriya bros": And("d1 railway chest", "bombs", "d1 key B"),
	"d1 basement":       AndSlot("enter goriya bros", "kill goriya bros"),
	"d1 goriya chest": AndSlot("d1 stalfos chest",
		Or("ember seeds", Hard("mystery seeds")), "kill goriya (pit)"),
	"d1 floormaster room": AndSlot("enter d1",
		Or("ember seeds", Hard("mystery seeds"))),
	"d1 boss": AndSlot("d1 floormaster room", "d1 boss key",
		"kill aquamentus"),

	"d1 key A": And("d1 key fall"),
	"d1 key B": And("d1 key chest"),
}

var seasonsD2Nodes = map[string]*Node{
	"d2 left from entrance": AndSlot("d2 torch room"),
	"d2 rope room":          And("d2 torch room", "kill normal"),
	"d2 arrow room": Or("enter d2 B",
		And("d2 torch room", Or("ember seeds", Hard("mystery seeds")))),
	"d2 rupee room":   And("d2 arrow room", "bombs"),
	"d2 hardhat room": And("d2 arrow room", "d2 2 keys"), // min. 1 key
	"d2 pot chest":    AndSlot("d2 hardhat room", "remove pot"),
	"d2 rope chest":   AndSlot("d2 arrow room", "kill normal"),
	"d2 bracelet room": And("d2 hardhat room",
		Or("kill hardhat (pit)", Hard("bracelet"))),
	"d2 moblin chest": AndSlot("d2 bracelet room",
		Or("kill moblin (gap)", Hard("bracelet"))),
	"d2 spiral chest": And("enter d2 B", "bombs"),
	"d2 blade chest": Or("enter d2 B",
		And("d2 arrow room", Or("kill normal", Hard("bracelet")))),

	// from here on it's entirely linear.
	"d2 roller chest":  AndSlot("d2 bomb wall", "bombs", "bracelet"),
	"d2 spinner":       And("d2 roller chest", "d2 2 keys"), // min. 1 key
	"d2 terrace chest": AndSlot("d2 spinner", "d2 3 keys"),
	"d2 boss":          AndSlot("d2 spinner", "d2 boss key"),

	"d2 key A": And("d2 rope room"),
	"d2 key B": And("d2 spiral chest"),
	"d2 key C": And("d2 blade chest"),
	"d2 2 keys": Or(And("d2 key A", "d2 key B"), And("d2 key A", "d2 key C"),
		And("d2 key B", "d2 key C")),
	"d2 3 keys": And("d2 key A", "d2 key B", "d2 key C"),

	"d2 torch room": Or("enter d2 A", "d2 rope chest"),
	"d2 bomb wall":  And("d2 blade chest"), // alias for external reference
}

var seasonsD3Nodes = map[string]*Node{
	// first floor
	"d3 center": And("enter d3",
		Or("kill spiked beetle", HardAnd("flip spiked beetle", "bracelet"))),
	"d3 mimic stairs":      Or("d3 water room", And("d3 center", "bracelet")),
	"d3 roller chest":      And("d3 mimic stairs", "bracelet"),
	"d3 water room":        OrSlot("d3 mimic stairs", And("d3 center", "jump 2")),
	"d3 quicksand terrace": AndSlot("d3 mimic stairs", "jump 2"),
	"d3 omuai stairs": And("d3 mimic stairs", "jump 2", "kill omuai",
		"d3 2 keys"), // min. 1 key
	"d3 giant blade room": AndSlot("d3 omuai stairs"),

	// second floor
	"d3 moldorm chest":     AndSlot("d3 mimic stairs", "kill moldorm"),
	"d3 bombed wall chest": AndSlot("d3 moldorm chest", "bombs"),
	"d3 mimic chest": AndSlot("d3 water room", "kill mimic",
		"d3 2 keys"), // min. 1 key
	"d3 trampoline chest": AndSlot("d3 center", "jump 2"),
	"enter mothula":       And("d3 omuai stairs", "d3 boss key"),
	"d3 boss":             AndSlot("enter mothula", "kill mothula"),

	// fixed items
	"d3 key A":  And("d3 roller chest"),
	"d3 key B":  And("d3 trampoline chest"),
	"d3 2 keys": And("d3 key A", "d3 key B"),
}

var seasonsD4Nodes = map[string]*Node{
	// left branch from entrance
	"d4 north of entrance": AndSlot("enter d4", Or("flippers", "jump 4")),
	"d4 pot room":          And("d4 north of entrance", "bombs", "bracelet"),
	"d4 maze chest": AndSlot("d4 north of entrance",
		"hit lever from minecart"),
	"d4 dark chest": And("d4 maze chest", "jump 2"),

	// 2F (ground floor), right branch
	"d4 water ring room": AndSlot("enter d4", Or("flippers", "jump 4"), "bombs",
		"d4 1 key"),
	"d4 roller minecart": And("enter d4", "flippers", "jump 2", "d4 1 key"),
	"d4 water key room": And("d4 roller minecart", "hit lever from minecart",
		Or("kill normal", Hard("bracelet"))),
	"d4 stalfos stairs": And("d4 roller minecart", "d4 2 keys",
		Or("kill stalfos", Hard("bracelet"))),

	// 1F
	"d4 pre-mid chest":      And("d4 stalfos stairs"),
	"d4 final minecart":     And("d4 stalfos stairs", "kill agunima"),
	"d4 torch chest":        And("d4 stalfos stairs", "ember slingshot"),
	"d4 cracked floor room": AndSlot("d4 final minecart", "d4 5 keys"), // min. 2 keys
	"d4 dive spot": AndSlot("d4 final minecart", "hit very far lever",
		"d4 5 keys"), // min. 2 keys
	"d4 basement stairs": And("d4 final minecart", "hit far lever",
		"d4 5 keys"), // min. 2 keys

	// B1F
	"enter gohma": And("d4 basement stairs", "d4 boss key",
		Or("ember slingshot", Hard("mystery slingshot"), "jump 3",
			HardAnd("jump 2", Or("ember seeds", "mystery seeds")))),
	"d4 boss": AndSlot("enter gohma", "kill gohma"),

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

// if you can reach the key door after the fire trap, you can necessarily reach
// the rest of the keys in this dungeon. this means that the other doors need
// no more than four out of five keys in logic.
var seasonsD5Nodes = map[string]*Node{
	// 1F (it's the only F)
	"d5 cart bay":   And("enter d5", Or("flippers", "bomb jump 2")),
	"d5 cart chest": And("d5 cart bay", "hit lever from minecart"),
	"d5 pot room": And("enter d5", Or(And("magnet gloves", "bombs", "jump 2"),
		And("d5 cart bay", Or("jump 2", Hard("pegasus satchel"))))),
	"d5 gibdo/zol chest": AndSlot("d5 pot room", "kill gibdo", "kill zol"),
	"d5 magnet ball chest": AndSlot("d5 pot room",
		Or("flippers", "jump 6", Hard("jump 4")), "d5 4 keys"),
	"d5 left chest": And("enter d5", Or("magnet gloves", "jump 4")),
	"d5 terrace chest": AndSlot("enter d5", Or("magnet gloves",
		And("d5 cart bay", "jump 2", "bombs"))),
	"d5 spiral chest": AndSlot("enter d5", Or("shield",
		And("kill moldorm", "kill iron mask"))),
	"d5 armos chest": And("d5 terrace chest", "kill moldorm", "kill iron mask",
		"kill armos"),
	"d5 spinner chest": And("d5 cart bay", Or("magnet gloves", "jump 6")),
	"d5 drop ball": And("d5 cart bay", "hit lever from minecart",
		"kill darknut (pit)"),
	"d5 pre-mid chest": And("d5 cart bay", Or("magnet gloves", "jump 4")),
	"d5 post-syger":    And("d5 pre-mid chest", "kill syger"),
	// you always have access to enough small keys for these nodes:
	"d5 basement": AndSlot("d5 drop ball", "d5 post-syger",
		"magnet gloves", Or("kill magunesu", Hard("jump 2"))),
	"d5 boss": AndSlot("d5 post-syger", "magnet gloves",
		Or("jump 2", Hard()), "d5 boss key"),

	// fixed items
	"d5 key A": And("d5 cart chest"),
	"d5 key B": And("d5 left chest"),
	"d5 key C": And("d5 armos chest"),
	"d5 key D": And("d5 spinner chest"),
	"d5 key E": And("d5 pre-mid chest"),
	"d5 4 keys": Or(
		And("d5 key A", Or(
			And("d5 key B", Or(
				And("d5 key C", Or("d5 key D", "d5 key E")),
				And("d5 key D", "d5 key E"))),
			And("d5 key C", "d5 key D", "d5 key E"))),
		And("d5 key B", "d5 key C", "d5 key D", "d5 key E")),
}

var seasonsD6Nodes = map[string]*Node{
	// 1F
	"d6 1F east":    AndSlot("enter d6"),
	"d6 rupee room": And("enter d6", "bombs"),
	"d6 beetle key room": And("enter d6",
		Or(And("magnet gloves", "jump 2"), "jump 4")),
	"d6 beamos room":       AndSlot("enter d6", "d6 key A", "d6 key C"),
	"d6 1F terrace":        AndSlot("enter d6"),
	"d6 crystal trap room": AndSlot("enter d6"),
	"d6 U-room":            And("enter d6", "break crystal", "boomerang L-2"),
	"d6 torch stairs":      And("d6 U-room", "ember seeds"),

	// 2F
	"d6 skipped chest": And("enter d6", "kill normal", "magnet gloves",
		"break crystal", Or("jump 2", Hard())),
	"d6 2F gibdo chest": AndSlot("d6 beamos room"),
	"d6 2F armos chest": AndSlot("d6 2F gibdo chest", "bombs"),
	"d6 escape room":    AndSlot("d6 torch stairs", "jump 2"),

	// 3F
	"d6 armos hall": AndSlot("d6 2F armos chest"),
	"d6 vire chest": And("d6 escape room", "kill stalfos"),
	"enter vire":    And("d6 vire chest", "d6 3 keys"), // min. 1 key

	// 5F
	"d6 pre-boss room": And("enter vire", "kill vire", "kill hardhat (magnet)"),
	"d6 boss": AndSlot("d6 pre-boss room", "d6 boss key",
		"kill manhandla"),

	// fixed items
	"d6 key A":  And("d6 beetle key room"),
	"d6 key B":  And("d6 vire chest"),
	"d6 key C":  And("d6 skipped chest"),
	"d6 3 keys": And("d6 key A", "d6 key B", "d6 key C"),
}

// poe skip with magnet gloves is possible in hard logic since you can't
// keylock that way and there's no warning, but you can still waste a key on
// the first key door, so the only difference it makes is that you don't have
// to kill the first poe.
var seasonsD7Nodes = map[string]*Node{
	// 1F
	"d7 wizzrobe chest":    And("enter d7", "kill wizzrobe"),
	"d7 right of entrance": AndSlot("enter d7", "d7 key A"),
	"enter poe A": And("d7 right of entrance",
		Or("ember slingshot", Hard("mystery slingshot"))),
	"d7 bombed wall chest": AndSlot("enter d7", "bombs"),
	"d7 quicksand chest":   AndSlot("d7 pot room", "jump 2", "d7 key B"),

	// B1F
	"d7 pot room": And("enter d7", "bracelet", Or(
		And("enter poe A", "kill poe sister"),
		HardAnd("magnet gloves", "jump 2", "pegasus satchel"))),
	"d7 zol button": And("d7 pot room", "jump 2"),
	"d7 magunesu chest": And("d7 armos puzzle", "jump 3", "kill magunesu",
		"magnet gloves"),
	"enter poe B": And("d7 pot room", "d7 3 keys", "ember seeds",
		Or("pegasus satchel", "slingshot L-2", Hard())),
	"d7 water stairs": And("enter poe B", "flippers"),
	"d7 spike chest":  AndSlot("d7 water stairs", "d7 cross bridge"),

	// B2F
	"d7 armos puzzle": And("d7 pot room", Or("jump 3", "magnet gloves")),
	"d7 cross bridge": Or("jump 4", "kill darknut (across pit)",
		And("jump 2", "magnet gloves")),
	"d7 maze chest": AndSlot("d7 water stairs", "kill moldorm", "jump 4",
		"d7 4 keys"),
	"d7 skipped room":  And("d7 maze chest"),
	"d7 stalfos chest": AndSlot("d7 maze chest", "d7 key E"),
	"d7 boss":          AndSlot("d7 maze chest", "d7 boss key", "kill gleeok"),

	// fixed items
	"d7 key A": And("d7 wizzrobe chest"),
	"d7 key B": And("d7 zol button"),
	"d7 key C": And("d7 armos puzzle"),
	"d7 key D": And("d7 magunesu chest"),
	"d7 key E": And("d7 skipped room"),
	"d7 3 keys": And("d7 key A", Or(
		And("d7 key B", Or("d7 key C", "d7 key D")),
		And("d7 key C", "d7 key D"))),
	"d7 4 keys": And("d7 key A", "d7 key B", "d7 key C", "d7 key D"),
}

// this does *not* account for HSS skip.
//
// possible but not in logic: hitting the sets of three eye statues quickly
// enough to make the chest/stairs appear, without HSS.
//
// pots don't hurt magunesu, thank goodness.
var seasonsD8Nodes = map[string]*Node{
	// 1F
	"d8 eye room": And("enter d8", "remove pot", Or("any slingshot",
		HardAnd("jump 2",
			Or("ember satchel", "scent satchel", "mystery satchel")))),
	"d8 three eyes chest": AndSlot("enter d8", "any slingshot L-2", "jump 2"),
	"d8 hardhat room":     And("enter d8", "kill magunesu"),
	"d8 hardhat key":      And("d8 hardhat room", "kill hardhat (magnet)"),
	"d8 spike room": AndSlot("d8 hardhat room", "d8 1 key",
		Or("jump 4", Hard("jump 3"))),
	"d8 magnet ball room": AndSlot("d8 spinner"),
	"d8 bomb chest": And("d8 armos chest", "any slingshot L-2", "bombs",
		"kill darknut"),
	"d8 ice puzzle room": And("d8 armos chest", "kill frypolar", "ember seeds",
		"slingshot L-2"),
	"d8 pols voice chest": AndSlot("d8 ice puzzle room",
		Or("jump 6", "boomerang L-2", Hard())),
	"d8 crystal room": And("d8 ice puzzle room", "d8 4 keys"),
	"d8 ghost armos":  And("d8 crystal room"),
	"d8 NW crystal":   And("d8 crystal room", "bracelet", "d8 7 keys"),
	"d8 NE crystal":   And("d8 crystal room", "bracelet", "hit lever"),
	"d8 SE crystal":   And("d8 crystal room", "bracelet"),
	"d8 SW crystal":   And("d8 crystal room", "bracelet", "d8 7 keys"),
	"d8 pot chest":    And("d8 SE crystal", "d8 NE crystal", "remove pot"),

	// B1F
	"d8 spinner":       And("d8 spike room", "d8 2 keys"),
	"d8 armos chest":   AndSlot("d8 spinner", "magnet gloves"),
	"d8 spinner chest": And("d8 armos chest"),
	"d8 SE lava chest": And("d8 SE crystal"),
	"d8 SW lava chest": AndSlot("d8 crystal room"),
	"d8 boss": AndSlot("d8 SW crystal", "d8 SE crystal", "d8 NW crystal",
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
var seasonsD9Nodes = map[string]*Node{
	"enter onox": And("enter d9", "kill wizzrobe", "kill floormaster",
		"kill darknut", "kill facade"),
	"done": AndStep("enter onox", "kill onox"),
}
