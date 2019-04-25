package logic

// dungeons should rely on overworld information as little as possible.
// ideally "enter <dungeon>" is the only overworld item the dungeon nodes
// reference (and that node should not be defined here).
//
// bush- and pot-throwing is in hard logic, but with an arbitrary limit of
// three screen transitions per carry, and no more than two enemies can be
// required to be killed with one throw.

var seasonsD0Nodes = map[string]*Node{
	// 0 keys
	"d0 key chest":   AndSlot("enter d0"),
	"d0 rupee chest": OrSlot("remove bush safe", "flute"),

	// 1 key
	"d0 sword chest": AndSlot("enter d0", "d0 small key"),
}

// bush-throwing is in hard logic for a few rooms. goriya one only matters if
// you killed the stalfos with rod, and the lever one only matters if you
// killed the stalfos with bombs. bush-throwing is *not* in logic for the
// vanilla BK room, since you need to relight the torches every time you leave.
var seasonsD1Nodes = map[string]*Node{
	// 0 keys
	"d1 stalfos drop": AndSlot("enter d1", Or("kill stalfos", "bracelet")),
	"d1 floormaster room": AndSlot("enter d1",
		Or("ember seeds", Hard("mystery seeds"))),
	"d1 boss": AndSlot("d1 floormaster room", "d1 boss key", "kill armored"),

	// 1 key
	"d1 stalfos chest": AndSlot("enter d1", "d1 small key", "kill stalfos"),
	"d1 goriya chest": AndSlot("d1 stalfos chest",
		Or("ember seeds", Hard("mystery seeds")), "kill normal (pit)"),
	"d1 lever room": AndSlot("d1 stalfos chest"),
	"d1 block-pushing room": AndSlot("d1 stalfos chest",
		Or("kill normal", Hard("bracelet"))),
	"d1 railway chest": AndSlot("d1 stalfos chest",
		Or("hit lever", Hard("bracelet"))),
	"d1 button chest": AndSlot("d1 railway chest"),

	// 2 keys
	"d1 basement": AndSlot("d1 railway chest", "bombs",
		Count(2, "d1 small key"), "kill armored"),
}

var seasonsD2Nodes = map[string]*Node{
	// 0 keys
	"d2 torch room":         Or("enter d2 A", "d2 rope chest"),
	"d2 left from entrance": AndSlot("d2 torch room"),
	"d2 rope drop":          AndSlot("d2 torch room", "kill normal"),
	"d2 arrow room": Or("enter d2 B",
		And("d2 torch room", Or("ember seeds", Hard("mystery seeds")))),
	"d2 rope chest": AndSlot("d2 arrow room", "kill normal"),
	"d2 rupee room": And("d2 arrow room", "bombs"),
	"d2 blade chest": OrSlot("enter d2 B",
		And("d2 arrow room", Or("kill normal", Hard("bracelet")))),
	"d2 roller chest": AndSlot("d2 bomb wall", "bombs", "bracelet"),
	"d2 spiral chest": AndSlot("enter d2 B", "bombs"),

	// 2 keys
	"d2 spinner":  And("d2 roller chest", Count(2, "d2 small key")),
	"dodongo owl": And("mystery seeds", "d2 spinner"),
	"d2 boss":     AndSlot("d2 spinner", "d2 boss key"),

	// 3 keys
	"d2 hardhat room": And("d2 arrow room", Count(3, "d2 small key")),
	"d2 pot chest":    AndSlot("d2 hardhat room", "remove pot"),
	"d2 moblin chest": AndSlot("d2 hardhat room",
		Or(Hard("bracelet"), And("kill hardhat (pit)", "kill moblin (gap)"))),
	"d2 terrace chest": AndSlot("d2 spinner", Count(3, "d2 small key")),

	// alias for external reference
	"d2 bomb wall":     And("d2 blade chest"),
	"d2 bracelet room": And("d2 moblin chest"),
}

var seasonsD3Nodes = map[string]*Node{
	// 0 keys
	"spiked beetles owl": And("mystery seeds", "enter d3"),
	"d3 center": And("enter d3",
		Or("kill spiked beetle", And("flip spiked beetle", "bracelet"))),
	"d3 mimic stairs": Or("d3 water room", And("d3 center", "bracelet")),
	"d3 roller chest": AndSlot("d3 mimic stairs", "bracelet"),
	"d3 water room": OrSlot("d3 mimic stairs",
		And("d3 center", "jump 2")),
	"d3 quicksand terrace": AndSlot("d3 mimic stairs", "jump 2"),
	"omuai owl":            And("mystery seeds", "d3 mimic stairs", "jump 2"),
	"d3 moldorm chest":     AndSlot("d3 mimic stairs", "kill armored"),
	"d3 bombed wall chest": AndSlot("d3 moldorm chest", "bombs"),
	"trampoline owl":       And("mystery seeds", "d3 center", "jump 2"),
	"d3 trampoline chest":  AndSlot("d3 center", "jump 2"),
	"d3 zol chest":         AndSlot("d3 center", "jump 2"),

	// 2 keys
	"d3 mimic chest": AndSlot("d3 water room", "kill normal",
		Count(2, "d3 small key")),
	"d3 omuai stairs": And("d3 mimic stairs", "jump 2", "kill omuai",
		Count(2, "d3 small key")),
	"d3 giant blade room": AndSlot("d3 omuai stairs"),
	"d3 boss":             AndSlot("d3 omuai stairs", "d3 boss key"),
}

var seasonsD4Nodes = map[string]*Node{
	// 0 keys
	"d4 north of entrance": AndSlot("enter d4", Or("flippers", "jump 4")),
	"d4 pot puzzle":        AndSlot("d4 north of entrance", "bombs", "bracelet"),
	"d4 maze chest": AndSlot("d4 north of entrance",
		"hit lever from minecart"),
	"d4 dark room": AndSlot("d4 maze chest", "jump 2"),

	// 1 key
	"d4 water ring room": AndSlot("enter d4", Or("flippers", "jump 4"), "bombs",
		"d4 small key", Or("bracelet", "kill normal", And("rod", "boomerang"))),
	"d4 roller minecart": And("enter d4", "flippers", "jump 2", "d4 small key"),
	"d4 pool": AndSlot("d4 roller minecart", "hit lever from minecart",
		Or("kill normal", "bracelet")),

	// 2 keys
	"greater distance owl": And("mystery seeds", "d4 roller minecart",
		Count(2, "d4 small key")),
	"d4 stalfos stairs": And("d4 roller minecart", Count(2, "d4 small key"),
		Or("kill stalfos", "bracelet")),
	"d4 terrace":        AndSlot("d4 stalfos stairs"),
	"d4 final minecart": And("d4 stalfos stairs", "kill agunima"),
	"d4 torch chest":    AndSlot("d4 stalfos stairs", "ember slingshot"),

	// 5 keys
	"d4 cracked floor room": AndSlot("d4 final minecart",
		Count(5, "d4 small key")),
	"d4 dive spot": AndSlot("d4 final minecart", "hit very far lever",
		Count(5, "d4 small key")),
	"d4 basement stairs": And("d4 final minecart", "hit far lever",
		Count(5, "d4 small key")),
	"gohma owl": And("mystery seeds", "d4 basement stairs"),
	"enter gohma": And("d4 basement stairs", "d4 boss key",
		Or("ember slingshot", Hard("mystery slingshot"), "jump 3",
			HardAnd("jump 2", Or("ember seeds", "mystery seeds")))),
	"d4 boss": AndSlot("enter gohma", "kill gohma"),

	// alias for external reference
	"enter agunima": And("d4 terrace"),
}

var seasonsD5Nodes = map[string]*Node{
	// 0 keys
	"d5 cart bay":   And("enter d5", Or("flippers", "bomb jump 2")),
	"d5 cart chest": AndSlot("d5 cart bay", "hit lever from minecart"),
	"d5 pot room": And("enter d5", Or(And("magnet gloves", "bombs", "jump 2"),
		And("d5 cart bay", Or("jump 2", Hard("pegasus satchel"))))),
	"d5 gibdo/zol chest": AndSlot("d5 pot room", "kill normal"),
	"d5 left chest":      AndSlot("enter d5", Or("magnet gloves", "jump 4")),
	"d5 terrace chest": AndSlot("enter d5", Or("magnet gloves",
		And("d5 cart bay", "jump 2", "bombs"))),
	"armos knights owl": And("mystery seeds", "d5 terrace chest"),
	"d5 spiral chest":   AndSlot("enter d5", Or("shield", "kill armored")),
	// if you can kill the moldorms, you can kill the armos
	"d5 armos chest":   AndSlot("d5 terrace chest", "kill armored"),
	"d5 spinner chest": AndSlot("d5 cart bay", Or("magnet gloves", "jump 6")),
	"d5 drop ball": And("d5 cart bay", "hit lever from minecart",
		"kill armored (pit)"),
	// stalfos room means right side (where the chest is)
	"d5 stalfos room": AndSlot("d5 cart bay", Or("magnet gloves", "jump 4")),

	// 5 keys
	"d5 post-syger": And("d5 stalfos room", "kill armored"),
	"d5 magnet ball chest": AndSlot("d5 pot room",
		Or("flippers", "jump 6", Hard("jump 4")), Count(5, "d5 small key")),
	"d5 basement": AndSlot("d5 drop ball", "d5 post-syger", "magnet gloves",
		Count(5, "d5 small key"), Or("kill magunesu", Hard("jump 2"))),
	"d5 boss": AndSlot("d5 post-syger", "magnet gloves", "d5 boss key",
		Count(5, "d5 small key"), Or("jump 2", Hard())),
}

var seasonsD6Nodes = map[string]*Node{
	// 0 keys
	"d6 1F east":    AndSlot("enter d6"),
	"d6 rupee room": And("enter d6", "bombs"),
	"d6 magnet ball drop": AndSlot("enter d6",
		Or(And("magnet gloves", "jump 2"), "jump 4")),
	"d6 1F terrace":        AndSlot("enter d6"),
	"d6 crystal trap room": AndSlot("enter d6"),
	"d6 U-room":            And("enter d6", "break crystal", "boomerang L-2"),
	"d6 torch stairs":      And("d6 U-room", "ember seeds"),
	"d6 escape room":       AndSlot("d6 torch stairs", "jump 2"),
	"d6 vire chest":        AndSlot("d6 escape room", "kill stalfos"),

	// 3 keys
	"d6 beamos room":    AndSlot("enter d6", Count(3, "d6 small key")),
	"d6 2F gibdo chest": AndSlot("d6 beamos room"),
	"d6 2F armos chest": AndSlot("d6 2F gibdo chest", "bombs"),
	"d6 armos hall":     AndSlot("d6 2F armos chest", "jump 2"),
	"d6 spinner north": AndSlot("enter d6", "magnet gloves", "break crystal",
		Or("jump 2", Hard()), Or(
			And("kill normal", Count(3, "d6 small key")),
			And("bombs", "jump 2", Count(2, "d6 small key")))),
	"enter vire":       And("d6 vire chest", Count(3, "d6 small key")),
	"d6 pre-boss room": And("enter vire", "kill vire", "kill hardhat (magnet)"),
	"d6 boss": AndSlot("d6 pre-boss room", "d6 boss key",
		"kill manhandla"),
}

// the chest to the right of the entrance is locked behind all keys in logic
// because of poe skip, but the poe fight itself is only locked behind one.
var seasonsD7Nodes = map[string]*Node{
	// 0 keys
	"poe curse owl":        And("mystery seeds", "enter d7"),
	"d7 wizzrobe chest":    AndSlot("enter d7", "kill normal"),
	"d7 bombed wall chest": AndSlot("enter d7", "bombs"),

	// 1 key
	"enter poe A": And("enter d7", "d7 small key",
		Or("ember slingshot", Hard("mystery slingshot"))),
	"d7 pot room": And("enter d7", "bracelet", Or(
		And("enter poe A", "kill poe sister"),
		HardAnd("bombs", "jump 2", "pegasus satchel"))),
	"d7 zol button":   AndSlot("d7 pot room", "jump 2"),
	"d7 armos puzzle": AndSlot("d7 pot room", Or("jump 3", "magnet gloves")),
	"d7 magunesu chest": AndSlot("d7 armos puzzle", "jump 3", "kill magunesu",
		"magnet gloves"),

	// 2 keys
	"d7 quicksand chest": AndSlot("d7 pot room", "jump 2",
		Count(2, "d7 small key")),

	// 3 keys
	"enter poe B": And("d7 pot room", "ember seeds", Count(3, "d7 small key"),
		Or("pegasus satchel", "slingshot L-2", Hard())),
	"d7 water stairs": And("enter poe B", "flippers"),
	"d7 spike chest": AndSlot("d7 water stairs",
		Or("jump 4",
			And("jump 3", "kill armored (across pit)"),
			And("jump 2", "magnet gloves"))),

	// 4 keys
	"d7 maze chest": AndSlot("d7 water stairs", "kill poe sister", "jump 4",
		Count(4, "d7 small key")),
	"d7 B2F drop": AndSlot("d7 maze chest",
		Or("magnet gloves", Hard("jump 6"))),
	"shining blue owl": And("mystery seeds", "d7 stalfos chest"),
	"d7 boss":          AndSlot("d7 maze chest", "d7 boss key", "kill gleeok"),

	// 5 keys
	"d7 right of entrance": AndSlot("enter d7", Count(5, "d7 small key")),
	"d7 stalfos chest": AndSlot("d7 maze chest", Count(5, "d7 small key"),
		Or("pegasus satchel", Hard())),
}

// HSS skip is out of logic and can lead to key-locks.
// pots don't hurt magunesu, thank goodness.
var seasonsD8Nodes = map[string]*Node{
	// 0 keys
	"d8 eye drop": AndSlot("enter d8", "remove pot", Or("any slingshot",
		HardAnd("jump 2",
			Or("ember satchel", "scent satchel", "mystery satchel")))),
	"d8 three eyes chest": AndSlot("enter d8", "jump 2",
		Or("any slingshot L-2",
			HardOr("ember satchel", "scent satchel", "mystery satchel"))),
	"d8 hardhat room": And("enter d8", "kill magunesu"),
	"d8 hardhat drop": AndSlot("d8 hardhat room", "kill hardhat (magnet)"),

	// 1 key
	"d8 spike room": AndSlot("d8 hardhat room", "d8 small key",
		Or("jump 4", Hard("jump 3"))),

	// 2 keys
	"d8 spinner":          And("d8 spike room", Count(2, "d8 small key")),
	"silent watch owl":    And("mystery seeds", "d8 spinner"),
	"d8 magnet ball room": AndSlot("d8 spinner"),
	"d8 armos chest":      AndSlot("d8 spinner", "magnet gloves"),
	"d8 spinner chest":    AndSlot("d8 armos chest"),
	"frypolar owl":        And("mystery seeds", "d8 armos chest"),
	"d8 darknut chest": AndSlot("d8 armos chest", "bombs", "kill armored",
		Or("any slingshot L-2",
			HardOr("ember satchel", "scent satchel", "mystery satchel"))),

	// 3 keys
	"d8 ice puzzle room": And("d8 armos chest", Count(3, "d8 small key"),
		"kill frypolar", "ember seeds", "slingshot L-2"),
	"d8 pols voice chest": AndSlot("d8 ice puzzle room",
		Or("jump 6", "boomerang L-2", Hard())),

	// 4 keys
	"d8 crystal room":     And("d8 ice puzzle room", Count(4, "d8 small key")),
	"magical ice owl":     And("mystery seeds", "d8 crystal room"),
	"d8 ghost armos drop": AndSlot("d8 crystal room", "bombs"),
	"d8 NE crystal":       And("d8 crystal room", "bracelet", "hit lever"),
	"d8 SE crystal":       And("d8 crystal room", "bracelet"),
	"d8 SE lava chest":    AndSlot("d8 SE crystal"),
	"d8 spark chest": AndSlot("d8 SE crystal", "d8 NE crystal",
		"remove pot"),
	"d8 SW lava chest": AndSlot("d8 crystal room"),

	// 6 keys
	"d8 NW crystal": And("d8 crystal room", "bracelet",
		Count(6, "d8 small key")),
	"d8 SW crystal": And("d8 crystal room", "bracelet",
		Count(6, "d8 small key")),

	// 7 keys
	"d8 boss": AndSlot("d8 SW crystal", "d8 NW crystal",
		Count(7, "d8 small key"), "d8 boss key", "kill medusa head"),
}

// onox's castle
var seasonsD9Nodes = map[string]*Node{
	"done": AndStep("enter d9", "kill armored", "bombs", "kill onox"),
}
