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

var d0NodesAnd = map[string]Point{
	"d0 key chest":   And{"enter d0"},
	"d0 sword chest": And{"enter d0", "d0 small key"},
	"d0 rupee chest": And{"remove bush"},

	// TODO randomize
	"d0 small key": And{"d0 key chest"},
}

var d0NodesOr = map[string]Point{}

var d1NodesAnd = map[string]Point{
	"d1 key fall":       And{"enter d1", "kill stalfos (throw)"},
	"d1 map chest":      And{"d1 key 1", "kill stalfos"},
	"d1 compass chest":  And{"d1 map chest"},
	"d1 gasha chest":    And{"d1 map chest", "kill goriya"},
	"d1 bomb chest":     And{"d1 map chest", "hit lever"},
	"d1 key chest":      And{"d1 map chest", "hit lever"},
	"enter goriya bros": And{"d1 bomb chest", "bombs", "d1 key 2"},
	"d1 satchel":        And{"enter goriya bros", "kill goriya bros"},
	"d1 boss key chest": And{"d1 map chest", "ember seeds", "kill goriya (pit)"},
	"d1 ring chest":     And{"enter d1", "ember seeds"},
	"enter aquamentus":  And{"enter d1", "ember seeds", "d1 boss key"},
	"d1 essence":        And{"enter aquamentus", "kill aquamentus"},

	// TODO randomize
	"d1 key 1":    And{"d1 key fall"},
	"d1 key 2":    And{"d1 key chest"},
	"d1 boss key": And{"d1 boss key chest"},
}

var d1NodesOr = map[string]Point{}

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
var d2NodesAnd = map[string]Point{
	"d2 5-rupee chest":     And{"d2 torch room"},
	"d2 key fall":          And{"d2 torch room", "kill rope"},
	"d2 arrow room 1":      And{"d2 torch room", "ember seeds"},
	"d2 arrow room 2":      And{"enter d2 3", "bracelet"},
	"d2 hardhat room":      And{"d2 arrow room", "d2 key 1"},
	"d2 map chest":         And{"d2 hardhat room"},
	"d2 compass chest 1":   And{"d2 torch room", "ember seeds", "kill rope"},
	"d2 compass chest 2":   And{"d2 arrow room", "kill goriya", "kill rope"},
	"d2 bracelet chest":    And{"d2 hardhat room", "kill hardhat (pit, throw)", "kill moblin (gap, throw)"},
	"d2 bomb key chest":    And{"enter d2 2", "remove bush", "bombs"},
	"d2 blade key chest 1": And{"enter d2 3", "bracelet"},
	"d2 blade key chest 2": And{"d2 arrow room", "kill rope", "kill goriya"},

	// TODO AND nodes can never require each other. write a routine to check
	//      for mutual dependencies in the raw graph.
	"d2 bomb wall": And{"d2 blade key chest"}, // alias for external reference

	// from here on it's entirely linear
	"d2 10-rupee chest": And{"d2 bomb wall", "bombs", "bracelet"},
	"enter facade":      And{"d2 10-rupee chest", "bracelet", "d2 key 2"},
	"d2 boss key chest": And{"enter facade", "d2 key 3", "bombs"},
	"enter dodongo":     And{"d2 boss key chest", "d2 boss key"},
	"d2 essence":        And{"enter dodongo", "kill dodongo"},

	// TODO randomize
	"d2 key 1":    And{"d2 key fall"},
	"d2 key 2":    And{"d2 bomb key chest"},
	"d2 key 3":    And{"d2 blade key chest"},
	"d2 boss key": And{"d2 boss key chest"},
}

var d2NodesOr = map[string]Point{
	"d2 torch room":      Or{"enter d2 1", "d2 compass chest"},
	"d2 compass chest":   Or{"d2 compass chest 1", "d2 compass chest 2"},
	"d2 arrow room":      Or{"d2 arrow room 1", "d2 arrow room 2"},
	"d2 blade key chest": Or{"d2 blade key chest 1", "d2 blade key chest 2"},
}
