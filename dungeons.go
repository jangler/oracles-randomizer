package main

// keep chests and chest items separate, so they can be altered later
// if possible
//
// the dungeon should rely on overworld information as little as possible.
// ideally "enter <dungeon>" is the only overworld item the dungeon nodes
// reference.

// TODO: d0

var d1ChestItems = []string{
	"", // 1-indexed
	"d1 map",
	"d1 compass",
	"gasha seed",
	"", // bombs, can replace with w/e
	"d1 key 2",
	"d1 boss key",
	"ring",
}

var d1NodesAnd = map[string][]string{
	"d1 key 1":          []string{"enter d1", "kill stalfos (throw)"},
	"d1 chest 1":        []string{"d1 key 1", "kill stalfos"},
	d1ChestItems[1]:     []string{"d1 chest 1"}, // map
	"d1 chest 2":        []string{"d1 chest 1"},
	d1ChestItems[2]:     []string{"d1 chest 2"}, // compass
	"d1 chest 3":        []string{"d1 chest 1", "kill goriya"},
	d1ChestItems[3]:     []string{"d1 chest 3"}, // gasha seed
	"d1 chest 4":        []string{"d1 chest 1", "hit lever"},
	d1ChestItems[4]:     []string{"d1 chest 4"}, // bombs
	"d1 chest 5":        []string{"d1 chest 1", "hit lever"},
	d1ChestItems[5]:     []string{"d1 chest 5"}, // small key
	"enter goriya bros": []string{"d1 key 2", "bombs"},
	"d1 warp":           []string{"enter goriya bros", "kill goriya bros"},
	"satchel":           []string{"d1 warp"},
	"d1 chest 6":        []string{"d1 chest 1", "ember seeds", "kill goriya (pit)"},
	d1ChestItems[6]:     []string{"d1 chest 6"}, // boss key
	"d1 chest 7":        []string{"enter d1", "ember seeds"},
	d1ChestItems[7]:     []string{"d1 chest 7"}, // ring
	"enter aquamentus":  []string{"ember seeds", "d1 boss key"},
	"d1 essence":        []string{"enter aquamentus", "kill aquamentus"},
}

var d1NodesOr = map[string][]string{}

// TODO: d2
