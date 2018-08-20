package prenode

// keep chests and chest items separate, so they can be altered later
// if possible
//
// the dungeon should rely on overworld information as little as possible.
// ideally "enter <dungeon>" is the only overworld item the dungeon nodes
// reference (and that node should not be defined here)
//
// make sure there's only *one* reference to each small key in a dungeon's
// requirements. it might make key counting easier for the routing algorithm.
//
// not that keys can NOT be numbered 1..n because of the code generation
// syntax; label them A..N instead.

var d0Prenodes = map[string]*Prenode{
	"d0 key chest":   And("enter d0"),
	"d0 sword chest": AndSlot("enter d0", "d0 small key"),
	"d0 rupee chest": AndSlot("remove bush"),

	"d0 small key": And("d0 key chest"),
}

var d1Prenodes = map[string]*Prenode{
	"d1 key fall":       And("enter d1", "kill stalfos (throw)"),
	"d1 map chest":      AndSlot("d1 key A", "kill stalfos"),
	"d1 compass chest":  AndSlot("d1 map chest"),
	"d1 gasha chest":    AndSlot("d1 map chest", "kill goriya"),
	"d1 bomb chest":     AndSlot("d1 map chest", "hit lever"),
	"d1 key chest":      And("d1 map chest", "hit lever"),
	"enter goriya bros": And("d1 bomb chest", "bombs", "d1 key B"),
	"d1 satchel spot":   AndSlot("enter goriya bros", "kill goriya bros"),
	"d1 boss key chest": And("d1 map chest", "ember seeds", "kill goriya (pit)"),
	"d1 ring chest":     AndSlot("enter d1", "ember seeds"),
	"enter aquamentus":  And("enter d1", "ember seeds", "d1 boss key"),
	"d1 essence":        AndStep("enter aquamentus", "kill aquamentus"),

	"d1 key A":    And("d1 key fall"),
	"d1 key B":    And("d1 key chest"),
	"d1 boss key": And("d1 boss key chest"),
}

// this is tricky because of the multiple entrances. the nexus is what
// i'll call the "arrow room" because of the arrow-shaped block arrangement in
// it. you can either get to this room by entering the main way and lighting
// the torches, or by entering the third way (into the roller room), pushing
// the rollers, and killing ropes and goriya.
//
// another weird thing about this dungeon is that if you enter via the
// secondary entrances, the save location is set to just outside the main
// entrance. this doesn't really matter because you need to remove bushes in
// order to use any entrance, though.
//
// you can actually complete this entire dungeon without ember seeds, since
// they're only required to open one door, which you can circumvent via the
// various entrances.
var d2Prenodes = map[string]*Prenode{
	"d2 5-rupee chest": AndSlot("d2 torch room"),
	"d2 key fall":      And("d2 torch room", "kill rope"),
	"d2 arrow room": Or(
		And("d2 torch room", "ember seeds"),
		And("enter d2 B", "bracelet")),
	"d2 rupee room":   And("d2 arrow room", "bombs"),
	"d2 hardhat room": And("d2 arrow room", "d2 key A"),
	"d2 map chest":    AndSlot("d2 hardhat room", "remove pot"),
	"d2 compass chest": OrSlot(
		And("d2 torch room", "ember seeds", "kill rope"),
		And("d2 arrow room", "kill goriya", "kill rope")),
	"d2 bracelet chest": AndSlot("d2 hardhat room", "kill hardhat (pit, throw)", "kill moblin (gap, throw)"),
	"d2 bomb key chest": And("enter d2 B", "remove bush", "bombs"),
	"d2 blade key chest": Or(
		And("enter d2 B", "bracelet"),
		And("d2 arrow room", "kill rope", "kill goriya")),

	"d2 bomb wall": And("d2 blade key chest"), // alias for external reference

	// from here on it's entirely linear
	"d2 10-rupee chest": AndSlot("d2 bomb wall", "bombs", "bracelet"),
	"enter facade":      And("d2 10-rupee chest", "remove pot", "d2 key B"),
	"d2 spinner":        And("enter facade", "kill facade", Or("d2 key C", "bombs")),
	"d2 boss key chest": And("enter facade", "kill facade", "d2 key C", "bombs"),
	"enter dodongo":     And("d2 boss key chest", "d2 boss key"),
	"d2 essence":        AndStep("enter dodongo", "kill dodongo"),

	"d2 key A":    And("d2 key fall"),
	"d2 key B":    And("d2 bomb key chest"),
	"d2 key C":    And("d2 blade key chest"),
	"d2 boss key": And("d2 boss key chest"),

	"d2 torch room": Or("enter d2 A", "d2 compass chest"),
}

var d3Prenodes = map[string]*Prenode{
	// first floor
	"d3 mimic stairs": Or("d3 feather stairs",
		And("enter d3", "kill spiked beetle (throw)", "bracelet")),
	"d3 roller key chest": And("d3 mimic stairs", "bracelet"),
	"d3 feather stairs": Or("d3 mimic stairs", "d3 basement B in",
		And("enter d3", "jump")),
	"d3 basement B in": And("jump", Or("d3 feather stairs", "d3 basement B out")),
	"d3 basement B out": Or(
		And("d3 basement B in", "jump"),
		And("d3 trampoline stairs", "bracelet")),
	"d3 rupee chest":    AndSlot("d3 feather stairs"),
	"enter omuai":       And("d3 mimic stairs", "jump", "d3 key B"),
	"d3 gasha chest":    AndSlot("d3 mimic stairs", "jump"),
	"d3 omuai stairs":   And("enter omuai", "kill omuai"),
	"d3 boss key chest": And("d3 omuai stairs", "jump"),
	"d3 basement A in":  And("jump", Or("d3 feather stairs", "d3 basement A out")),
	"d3 basement A out": Or("d3 trampoline stairs",
		And("d3 basement A in", "jump")),
	"d3 trampoline stairs": Or("d3 basement A out",
		And("d3 compass chest", "bracelet")),
	"d3 map chest": AndSlot("d3 basement B out", "jump"),

	// second floor
	"d3 bomb chest":           AndSlot("d3 mimic stairs"),
	"d3 compass chest":        AndSlot("d3 bomb chest", "bombs"),
	"d3 feather room":         And("d3 rupee chest", "d3 key A"),
	"d3 feather chest":        AndSlot("d3 feather room", "kill mimic"),
	"d3 trampoline key chest": And("d3 trampoline stairs", "jump"),
	"enter mothula":           And("d3 omuai stairs", "d3 boss key"),
	"d3 essence":              AndStep("enter mothula", "kill mothula"),

	// fixed items
	"d3 key A":    And("d3 roller key chest"),
	"d3 key B":    And("d3 trampoline key chest"),
	"d3 boss key": And("d3 boss key chest"),
}

// this whole dungeon is basically a tree so all the links are one-way
var d4Prenodes = map[string]*Prenode{
	// left branch from entrance
	"d4 bomb chest":     AndSlot("enter d4", "cross large pool"),
	"d4 pot key fall":   And("d4 bomb chest", "bombs", "bracelet"),
	"d4 statue stairs":  And("d4 bomb chest", "hit lever"),
	"d4 map chest":      AndSlot("d4 statue stairs"),
	"d4 dark key chest": And("d4 statue stairs", "jump"),

	// 2F (ground floor), right branch
	"d4 compass chest":   AndSlot("enter d4", "cross large pool", "d4 key A", "bombs"),
	"d4 roller minecart": And("enter d4", "flippers", "d4 key A", "jump"),
	"d4 water key fall":  And("d4 roller minecart", "hit lever", "kill water tektite (throw)", "kill like-like (pit, throw)", "flippers"),
	"d4 stalfos stairs":  And("d4 roller minecart", "kill shrouded stalfos (throw)", "jump", "d4 key B"),

	// 1F
	"d4 pre-mid key":     And("d4 stalfos stairs"),
	"enter agunima":      And("d4 pre-mid key", "jump"), // being nice
	"d4 final minecart":  And("enter agunima", "kill agunima"),
	"d4 torch key chest": And("enter agunima", "ember slingshot", "jump"),
	"d4 slingshot chest": AndSlot("d4 final minecart", "d4 key C"),
	"d4 boss key chest":  And("d4 final minecart", "hit very far lever", "jump", "d4 key D", "flippers"),
	"d4 basement stairs": And("d4 final minecart", "hit far lever", "kill wizzrobe (pit, throw)", "d4 key E"),

	// B1F
	"d4 cross bridge": Or("ember slingshot", "long jump"),
	"enter gohma":     And("d4 basement stairs", "d4 cross bridge", "d4 boss key"),
	"d4 essence":      AndStep("enter gohma", "kill gohma"),

	// fixed items
	"d4 key A":    And("d4 pot key fall"),
	"d4 key B":    And("d4 dark key chest"),
	"d4 key C":    And("d4 water key fall"),
	"d4 key D":    And("d4 pre-mid key"),
	"d4 key E":    And("d4 torch key chest"),
	"d4 boss key": And("d4 boss key chest"),
}

var d5Prenodes = map[string]*Prenode{
	// general
	"cross magnet gap":   Or("pegasus jump L-2", "magnet gloves"),
	"magnet jump":        And("jump", "magnet gloves"),
	"sidescroll magnets": Or("magnet jump", "pegasus jump L-2"),

	// 1F (it's the only F)
	"d5 cart bay":       And("enter d5", Or("flippers", "long jump")),
	"d5 cart key chest": And("d5 cart bay", "hit lever"),
	"d5 underground A":  Or("d5 stairs A in", "d5 stairs C in"),
	"d5 stairs A in":    And("d5 cart bay"),
	"d5 stairs B out": Or(
		And("d5 stairs A in", "jump"),
		And("d5 stairs C in", "bombs", "jump")),
	// stairs B out is one-way
	"d5 map chest":           AndSlot("d5 stairs B out"),
	"d5 magnet gloves chest": AndSlot("d5 stairs B out", "cross large pool", "d5 key A"),
	"d5 left key chest":      And("enter d5", "cross magnet gap"),
	"d5 stairs C out":        And("d5 underground A", "bombs", "jump"),
	"d5 stairs C in":         And("enter d5", "magnet gloves"),
	"d5 large rupee chest": OrSlot("d5 stairs C out",
		And("enter d5", "magnet gloves")),
	"d5 compass chest":     AndSlot("enter d5", "kill moldorm", "kill iron mask"),
	"d5 armos key chest":   And("d5 stairs C out", "kill moldorm", "kill iron mask", "kill armos"),
	"d5 float key chest":   And("d5 cart bay", "cross magnet gap"),
	"d5 drop ball":         And("d5 cart bay", "hit lever", "kill darknut (pit)"),
	"d5 pre-mid key chest": And("d5 cart bay", "cross magnet gap"),
	"enter syger":          And("d5 cart bay", "cross magnet gap", "d5 key B"),
	"d5 post-syger":        And("enter syger", "kill syger"),
	"d5 push ball":         And("d5 drop ball", "d5 post-syger", "d5 key C", "magnet gloves"),
	"d5 boss key spot":     And("d5 push ball", "d5 key D", "long jump", "sidescroll magnets"), // being nice
	"enter digdogger":      And("d5 post-syger", "d5 key E", "jump", "magnet gloves", "d5 boss key"),
	"d5 essence":           AndStep("enter digdogger", "kill digdogger"),

	// fixed items
	"d5 key A":    And("d5 cart key chest"),
	"d5 key B":    And("d5 left key chest"),
	"d5 key C":    And("d5 armos key chest"),
	"d5 key D":    And("d5 float key chest"),
	"d5 key E":    And("d5 pre-mid key chest"),
	"d5 boss key": And("d5 boss key spot"),
}

// all the rupee chests in this dungeon are trivial, so i'm ignoring which is
// which and just labeling them by letter.
var d6Prenodes = map[string]*Prenode{
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
	"d6 rupee chest C":   AndSlot("map chest"),
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
	"d6 boss key chest":    And("d6 torch stairs", "long jump"),
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
	"d6 key A":    And("d6 magnet key fall"),
	"d6 key B":    And("d6 vire key chest"),
	"d6 key C":    And("d6 skipped key chest"),
	"d6 boss key": And("d6 boss key chest"),
}

var d7Prenodes = map[string]*Prenode{
	// 1F
	"d7 wizzrobe key chest": And("enter d7", "kill wizzrobe"),
	"d7 ring chest":         AndSlot("enter d7", "d7 key A"),
	"enter poe A":           And("d7 ring chest", "ember slingshot"),
	"d7 compass chest":      AndSlot("enter d7", "bombs"),
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
	"d7 water stairs": And("enter poe B", "pegasus satchel", "ember seeds", "kill poe sister", "flippers"),
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
	"d7 boss key chest": And("d7 stairs room", "jump", "hit switch", "kill stalfos",
		Or("d7 key E", HardAnd("poe skip", "d7 key D", "d7 enter skipped"))),
	"enter gleeok": And("d7 stairs room", "d7 boss key"),
	"d7 essence":   AndStep("enter gleeok", "kill gleeok"),

	// fixed items
	"d7 key A":    And("d7 wizzrobe key chest"),
	"d7 key B":    And("d7 zol key fall"),
	"d7 key C":    And("d7 armos key fall"),
	"d7 key D":    And("d7 magunesu key chest"),
	"d7 key E":    And("d7 skipped key poof"),
	"d7 boss key": And("d7 boss key chest"),
}

// keys get pretty wonky; hopefully they're correct between HSS skip and
// non-HSS skip.
var d8Prenodes = map[string]*Prenode{
	// 1F
	"d8 eye key fall":     And("enter d8", "slingshot", "remove pot"),
	"d8 ring chest":       AndSlot("enter d8", "slingshot L-2", "jump"),
	"d8 hardhat room":     And("enter d8", "kill magunesu"),
	"d8 hardhat key fall": And("d8 hardhat room", "kill hardhat (magnet)"),
	"HSS skip":            HardAnd("d8 hardhat room", "kill zol", "bombs", "pegasus jump L-2"),
	"d8 cross bridge A":   Or(Hard("HSS skip"), And("d8 ice puzzle room", "long jump")),
	"d8 compass chest":    AndSlot("d8 blade room"),
	"d8 map chest":        AndSlot("d8 spinner"),
	"d8 bomb room":        And("d8 HSS stairs", "slingshot L-2"),
	"d8 bomb key chest":   And("d8 bomb room", "bombs", "kill darknut"),
	"d8 ice puzzle room": Or("d8 frypolar stairs",
		And("d8 cross bridge A", "long jump")),
	"d8 cross bridge B": And("d8 ice puzzle room", Or(
		"boomerang L-2", "pegasus jump L-2",
		And("jump", Or("boomerang", "slingshot L-2")),
		And("long jump", Or("bombs", "slingshot L-1")),
		HardAnd("long jump", Or("sword", "fool's ore")))),
	// technically there are pots to throw in this room but i don't care
	"d8 boss key chest": And("d8 cross bridge B", "kill keese", "kill pols voice (pit)"),
	"d8 crystal room": And("d8 ice puzzle room",
		Or("d8 key D", HardAnd("d8 key A", "HSS skip"))),
	"d8 armos key fall": And("d8 crystal room", "bombs"),
	"d8 NW crystal": And("d8 crystal room",
		Or("d8 key E", HardAnd("d8 key B", "HSS skip"))),
	"d8 NE crystal": And("d8 crystal room", "hit lever"),
	"d8 SE crystal": And("d8 crystal room"),
	"d8 SW crystal": And("d8 crystal room",
		Or("d8 key F", HardAnd("d8 key C", "HSS skip"))),
	"d8 pot key chest": And("d8 SE crystal", "d8 NE crystal", "remove pot"),

	// B1F
	"d8 cross pot path": Or(Hard("remove pot"), And("remove pot", "kill rope"), "jump"),
	"d8 double rollers": Or(
		And("d8 hardhat room", "d8 key A", "d8 cross pot path"),
		And("d8 blade room", "long jump")),
	"d8 blade room": Or("d8 spinner", And("d8 double rollers", "long jump")),
	"d8 spinner": Or(
		And("d8 blade room", "d8 key B"),
		And("enter frypolar", "kill frypolar", "magnet gloves"),
		And("d8 HSS stairs", "jump")),
	"d8 place ball":        And("d8 spinner", "magnet gloves"),
	"d8 HSS chest":         AndSlot("d8 spinner", "magnet gloves"),
	"d8 HSS stairs":        And("d8 HSS chest", Or("pegasus jump L-2", "d8 place ball")),
	"d8 spinner key chest": And("d8 HSS stairs"),
	"enter frypolar":       Or(And("d8 HSS stairs", "d8 key C"), "d8 frypolar stairs"),
	"d8 frypolar stairs": Or("d8 ice puzzle room",
		And("enter frypolar", "kill frypolar", "ember seeds", "slingshot L-2")),
	"d8 lava key chest": And("d8 SE crystal"),
	"d8 bomb chest":     AndSlot("d8 SW crystal"),
	"enter medusa head": And("d8 SW crystal", "d8 SE crystal", "d8 NW crystal", "d8 boss key",
		Or("d8 key G", HardAnd("d8 key D", "HSS skip"))),
	"d8 essence": AndStep("enter medusa head", "kill medusa head"),

	// fixed items
	"d8 key A":    And("d8 eye key fall"),
	"d8 key B":    And("d8 hardhat key fall"),
	"d8 key C":    And("d8 spinner key chest"),
	"d8 key D":    And("d8 bomb key chest"),
	"d8 key E":    And("d8 armos key fall"),
	"d8 key F":    And("d8 lava key chest"),
	"d8 key G":    And("d8 pot key chest"),
	"d8 boss key": And("d8 boss key chest"),
}

// onox's castle
var d9Prenodes = map[string]*Prenode{
	"enter onox": And("enter d9", "kill wizzrobe", "kill floormaster", "kill darknut", "kill facade"),
	"done":       AndStep("enter onox", "kill onox"),
}
