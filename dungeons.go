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

var d0NodesAnd = map[string][]string{
	"d0 key chest":   []string{"enter d0"},
	"d0 sword chest": []string{"enter d0", "d0 small key"},
	"d0 rupee chest": []string{"remove bush"},
}

var d0NodesOr = map[string][]string{}

var d1NodesAnd = map[string][]string{
	"d1 key fall":       []string{"enter d1", "kill stalfos (throw)"},
	"d1 map chest":      []string{"d1 key 1", "kill stalfos"},
	"d1 compass chest":  []string{"d1 map chest"},
	"d1 gasha chest":    []string{"d1 map chest", "kill goriya"},
	"d1 bomb chest":     []string{"d1 map chest", "hit lever"},
	"d1 key chest":      []string{"d1 map chest", "hit lever"},
	"enter goriya bros": []string{"d1 bomb chest", "bombs", "d1 key 2"},
	"d1 satchel":        []string{"enter goriya bros", "kill goriya bros"},
	"d1 boss key chest": []string{"d1 map chest", "ember seeds", "kill goriya (pit)"},
	"d1 ring chest":     []string{"enter d1", "ember seeds"},
	"enter aquamentus":  []string{"enter d1", "ember seeds", "d1 boss key"},
	"d1 essence":        []string{"enter aquamentus", "kill aquamentus"},
}

var d1NodesOr = map[string][]string{}

// this is tricky because of the multiple entrances. the nexus is what
// i'll call the "arrow room" because of the arrow-shaped block arrangement in
// it. you can either get to this room by entering the main way and lighting
// the torches, or by entering the third way (into the roller room), pushing
// the rollers, and killing ropes and goyira.
//
// another weird thing about this dungeon is that if you enter via the
// secondary entrances, the save location is set to just outside the main
// entrance. this doesn't really matter because you need to remove bushes in
// order to use any entrance, though.
//
// you can actually complete this entire dungeon without ember seeds, since
// they're only required to open one door, which you can circumvent via the
// various entrances.
var d2NodesAnd = map[string][]string{
	"d2 5-rupee chest":     []string{"d2 torch room"},
	"d2 key fall":          []string{"d2 torch room", "kill rope"},
	"d2 arrow room 1":      []string{"d2 torch room", "ember seeds"},
	"d2 arrow room 2":      []string{"enter d2 3", "bracelet"},
	"d2 hardhat room":      []string{"d2 arrow room", "d2 key 1"},
	"d2 map chest":         []string{"d2 hardhat room"},
	"d2 compass chest 1":   []string{"d2 torch room", "ember seeds", "kill rope"},
	"d2 compass chest 2":   []string{"d2 arrow room", "kill goriya", "kill rope"},
	"d2 bracelet chest":    []string{"d2 hardhat room", "kill hardhat (pit, throw)", "kill moblin (pit, throw)"},
	"d2 bomb key chest":    []string{"enter d2 2", "remove bush", "bombs"},
	"d2 blade key chest 1": []string{"enter d2 3", "bracelet"},
	"d2 blade key chest 2": []string{"d2 arrow room", "kill rope", "kill goyira"},

	// TODO AND nodes can never require each other. write a routine to check
	//      for mutual dependencies in the raw graph.
	"d2 bomb wall": []string{"d2 blade key chest"}, // alias for external reference

	// from here on it's entirely linear
	"d2 10-rupee chest": []string{"d2 bomb wall", "bombs", "bracelet"},
	"enter facade":      []string{"d2 10-rupee chest", "bracelet", "d2 key 2"},
	"d2 boss key chest": []string{"enter facade", "d2 key 3", "bombs"},
	"enter dodongo":     []string{"d2 boss key chest", "d2 boss key"},
	"d2 essence":        []string{"enter dodongo", "kill dodongo"},
}

var d2NodesOr = map[string][]string{
	"d2 torch room":      []string{"enter d2 1", "d2 compass chest"},
	"d2 compass chest":   []string{"d2 compass chest 1", "d2 compass chest 2"},
	"d2 arrow room":      []string{"d2 arrow room 1", "d2 arrow room 2"},
	"d2 blade key chest": []string{"d2 blade key chest 1", "d2 blade key chest 2"},
}
