package main

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

var d0Points = map[string]Point{
	"d0 key chest":   And{"enter d0"},
	"d0 sword chest": AndSlot{"enter d0", "d0 small key"},
	"d0 rupee chest": And{"remove bush"},

	"d0 small key": And{"d0 key chest"},
}

var d1Points = map[string]Point{
	"d1 key fall":       And{"enter d1", "kill stalfos (throw)"},
	"d1 map chest":      And{"d1 key A", "kill stalfos"},
	"d1 compass chest":  And{"d1 map chest"},
	"d1 gasha chest":    And{"d1 map chest", "kill goriya"},
	"d1 bomb chest":     And{"d1 map chest", "hit lever"},
	"d1 key chest":      And{"d1 map chest", "hit lever"},
	"enter goriya bros": And{"d1 bomb chest", "bombs", "d1 key B"},
	"d1 satchel":        AndSlot{"enter goriya bros", "kill goriya bros"},
	"d1 boss key chest": And{"d1 map chest", "ember seeds", "kill goriya (pit)"},
	"d1 ring chest":     And{"enter d1", "ember seeds"},
	"enter aquamentus":  And{"enter d1", "ember seeds", "d1 boss key"},
	"d1 essence":        And{"enter aquamentus", "kill aquamentus"},

	"d1 key A":    And{"d1 key fall"},
	"d1 key B":    And{"d1 key chest"},
	"d1 boss key": And{"d1 boss key chest"},
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
var d2Points = map[string]Point{
	"d2 5-rupee chest":     And{"d2 torch room"},
	"d2 key fall":          And{"d2 torch room", "kill rope"},
	"d2 arrow room 1":      And{"d2 torch room", "ember seeds"},
	"d2 arrow room 2":      And{"enter d2 C", "bracelet"},
	"d2 hardhat room":      And{"d2 arrow room", "d2 key A"},
	"d2 map chest":         And{"d2 hardhat room", "remove pot"},
	"d2 compass chest 1":   And{"d2 torch room", "ember seeds", "kill rope"},
	"d2 compass chest 2":   And{"d2 arrow room", "kill goriya", "kill rope"},
	"d2 bracelet chest":    AndSlot{"d2 hardhat room", "kill hardhat (pit, throw)", "kill moblin (gap, throw)"},
	"d2 bomb key chest":    And{"enter d2 B", "remove bush", "bombs"},
	"d2 blade key chest 1": And{"enter d2 C", "bracelet"},
	"d2 blade key chest 2": And{"d2 arrow room", "kill rope", "kill goriya"},

	"d2 bomb wall": And{"d2 blade key chest"}, // alias for external reference

	// from here on it's entirely linear
	"d2 10-rupee chest": And{"d2 bomb wall", "bombs", "bracelet"},
	"enter facade":      And{"d2 10-rupee chest", "remove pot", "d2 key B"},
	"d2 boss key chest": And{"enter facade", "kill facade", "d2 key C", "bombs"},
	"enter dodongo":     And{"d2 boss key chest", "d2 boss key"},
	"d2 essence":        And{"enter dodongo", "kill dodongo"},

	"d2 key A":    And{"d2 key fall"},
	"d2 key B":    And{"d2 bomb key chest"},
	"d2 key C":    And{"d2 blade key chest"},
	"d2 boss key": And{"d2 boss key chest"},

	"d2 torch room": Or{"enter d2 A", "d2 compass chest"},
}

var d3Points = map[string]Point{
	// first floor
	"d3 mimic stairs 1":      And{"enter d3", "kill spiked beetle (throw)", "bracelet"},
	"d3 mimic stairs 2":      And{"d3 feather stairs"},
	"d3 roller key chest":    And{"d3 mimic stairs", "bracelet"},
	"d3 feather stairs 1":    And{"enter d3", "jump"},
	"d3 feather stairs 2":    And{"d3 mimic stairs"},
	"d3 feather stairs 3":    And{"d3 basement B in"},
	"d3 basement B in 1":     And{"d3 feather stairs", "jump"},
	"d3 basement B in 2":     And{"d3 basement B out", "jump"},
	"d3 basement B out 1":    And{"d3 basement B in", "jump"},
	"d3 basement B out 2":    And{"d3 trampoline stairs", "bracelet"},
	"d3 rupee chest":         And{"d3 feather stairs"},
	"enter omuai":            And{"d3 mimic stairs", "jump", "d3 key B"},
	"d3 gasha chest":         And{"d3 mimic stairs", "jump"},
	"d3 omuai stairs":        And{"enter omuai", "kill omuai"},
	"d3 boss key chest":      And{"d3 omuai stairs", "jump"},
	"d3 basement A in 1":     And{"d3 feather stairs", "jump"},
	"d3 basement A in 2":     And{"d3 basement A out", "jump"},
	"d3 basement A out 1":    And{"d3 basement A in", "jump"},
	"d3 basement A out 2":    And{"d3 trampoline stairs"},
	"d3 trampoline stairs 1": And{"d3 basement A out"},
	"d3 trampoline stairs 2": And{"d3 compass chest", "bracelet"},
	"d3 map chest":           And{"d3 basement B out", "jump"},

	// second floor
	"d3 bomb chest":           And{"d3 mimic stairs"},
	"d3 compass chest":        And{"d3 bomb chest", "bombs"},
	"d3 feather room":         And{"d3 rupee chest", "d3 key A"},
	"d3 feather chest":        AndSlot{"d3 feather room", "kill mimic"},
	"d3 trampoline key chest": And{"d3 trampoline stairs", "jump"},
	"enter mothula":           And{"d3 omuai stairs", "d3 boss key"},
	"d3 essence":              And{"enter mothula", "kill mothula"},

	// fixed items
	"d3 key A":    And{"d3 roller key chest"},
	"d3 key B":    And{"d3 trampoline key chest"},
	"d3 boss key": And{"d3 boss key chest"},
}

// this whole dungeon is basically a tree so all the links are one-way
var d4Points = map[string]Point{
	// left branch from entrance
	"d4 bomb chest":     And{"enter d4", "cross large pool"},
	"d4 pot key fall":   And{"d4 bomb chest", "bombs", "bracelet"},
	"d4 statue stairs":  And{"d4 bomb chest", "hit lever"},
	"d4 map chest":      And{"d4 statue stairs"},
	"d4 dark key chest": And{"d4 statue stairs", "jump"},

	// 2F (ground floor), right branch
	"d4 compass chest":   And{"enter d4", "cross large pool", "d4 key A", "bombs"},
	"d4 roller minecart": And{"enter d4", "flippers", "d4 key A", "jump"},
	"d4 water key fall":  And{"d4 roller minecart", "hit lever", "kill water tektite (throw)", "kill like-like (pit, throw)", "flippers"},
	"d4 stalfos stairs":  And{"d4 roller minecart", "kill shrouded stalfos (throw)", "jump", "d4 key B"},

	// 1F
	"d4 pre-mid key":     And{"d4 stalfos stairs"},
	"enter agunima":      And{"d4 pre-mid key", "jump"}, // being nice
	"d4 final minecart":  And{"enter agunima", "kill agunima"},
	"d4 torch key chest": And{"enter agunima", "ember slingshot", "jump"},
	"d4 slingshot chest": And{"d4 final minecart", "d4 key C"},
	"d4 boss key chest":  And{"d4 final minecart", "hit very far lever", "jump", "d4 key D", "flippers"},
	"d4 basement stairs": And{"d4 final minecart", "hit far lever", "kill wizzrobe (pit, throw)", "d4 key E"},

	// B1F
	"d4 cross bridge": Or{"ember slingshot", "long jump"},
	"enter gohma":     And{"d4 basement stairs", "d4 cross bridge", "d4 boss key"},
	"d4 essence":      And{"enter gohma", "kill gohma"},

	// fixed items
	"d4 key A":    And{"d4 pot key fall"},
	"d4 key B":    And{"d4 dark key chest"},
	"d4 key C":    And{"d4 water key fall"},
	"d4 key D":    And{"d4 pre-mid key"},
	"d4 key E":    And{"d4 torch key chest"},
	"d4 boss key": And{"d4 boss key chest"},
}

var d5Points = map[string]Point{
	// general
	"cross magnet gap":   Or{"pegasus jump L-2", "magnet gloves"},
	"magnet jump":        And{"jump", "magnet gloves"},
	"sidescroll magnets": Or{"magnet jump", "pegasus jump L-2"},

	// 1F (it's the only F)
	"d5 cart bay":       And{"enter d5", "cross large pool"},
	"d5 cart key chest": And{"d5 cart bay", "hit lever"},
	"d5 underground A":  Or{"d5 stairs A in", "d5 stairs C in"},
	"d5 stairs A in":    And{"d5 cart bay"},
	"d5 stairs A out 1": And{"d5 stairs A in", "jump"},
	"d5 stairs A out 2": And{"d5 stairs C in", "bombs", "jump"},
	"d5 stairs B out 1": And{"d5 stairs A in", "jump"},
	"d5 stairs B out 2": And{"d5 stairs C in", "bombs", "jump"},
	// stairs B out is one-way
	"d5 map chest":           And{"d5 stairs B out"},
	"d5 magnet gloves chest": AndSlot{"d5 stairs B out", "cross large pool", "d5 key A"},
	"d5 left key chest":      And{"enter d5", "cross magnet gap"},
	"d5 stairs C out":        And{"d5 underground A", "bombs", "jump"},
	"d5 stairs C in":         And{"enter d5", "magnet gloves"},
	"d5 large rupee chest 1": And{"d5 stairs C out"},
	"d5 large rupee chest 2": And{"enter d5", "magnet gloves"},
	"d5 compass chest":       And{"enter d5", "kill moldorm", "kill iron mask"},
	"d5 armos key chest":     And{"d5 stairs C out", "kill moldorm", "kill iron mask", "kill armos"},
	"d5 float key chest":     And{"d5 cart bay", "cross magnet gap"},
	"d5 drop ball":           And{"d5 cart bay", "hit lever", "kill darknut (pit)"},
	"d5 pre-mid key chest":   And{"d5 cart bay", "cross magnet gap"},
	"enter syger":            And{"d5 cart bay", "cross magnet gap", "d5 key B"},
	"d5 post-syger":          And{"enter syger", "kill syger"},
	"d5 push ball":           And{"d5 drop ball", "d5 post-syger", "d5 key C", "magnet gloves"},
	"d5 boss key spot":       And{"d5 push ball", "d5 key D", "long jump", "sidescroll magnets"}, // being nice
	"enter digdogger":        And{"d5 post-syger", "d5 key E", "jump", "magnet gloves", "d5 boss key"},
	"d5 essence":             And{"enter digdogger", "kill digdogger"},

	// fixed items
	"d5 key A":    And{"d5 cart key chest"},
	"d5 key B":    And{"d5 left key chest"},
	"d5 key C":    And{"d5 armos key chest"},
	"d5 key D":    And{"d5 float key chest"},
	"d5 key E":    And{"d5 pre-mid key chest"},
	"d5 boss key": And{"d5 boss key spot"},
}

// TODO d6
// TODO d7
// TODO d8
