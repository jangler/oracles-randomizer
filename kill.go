package main

// these nodes do not define items, only which items can kill which enemies
// under what circumstances, assuming that you've arrived in the room
// containing the enemy.
//
// technically mystery seeds can be used to kill many enemies that can be
// killed by ember, scent, or gale seeds. mystery seeds are only included as a
// kill option if at all three of these seed types work.
//
// if an enemy is in the same room as a throwable object and is vulnerable to
// thrown objects, than just adding "bracelet" as an OR is sufficient.
//
// animal companions are not (yet?) included in this logic.
//
// these conditions are added only as necessary: if there's no point in the
// route so far that requires killing ropes in a room with pits, there won't be
// any "kill rope (pit)" node.

var killNodesAnd = map[string][]string{
	"slingshot kill normal":   []string{"slingshot", "seed kill normal"},
	"jump kill normal":        []string{"jump", "kill normal"},
	"slingshot gale seeds":    []string{"slingshot", "gale seeds"},
	"slingshot mystery seeds": []string{"slingshot", "gale seeds"},
	"kill dodongo":            []string{"bombs", "bracelet"},
}

var killNodesOr = map[string][]string{
	// enemies in normal route-ish order, but with prereqs first
	"seed kill normal":          []string{"ember seeds", "scent seeds", "gale seeds", "mystery seeds"},
	"kill normal":               []string{"sword", "bombs", "fists", "seed kill normal"},
	"pit kill normal":           []string{"sword", "shield", "bombs", "rod", "shovel", "fists", "scent seeds"},
	"kill normal (pit)":         []string{"kill normal", "pit kill normal"},
	"kill stalfos":              []string{"kill normal", "rod"},
	"kill stalfos (throw)":      []string{"kill stalfos", "bracelet"},
	"kill goriya bros":          []string{"sword", "bombs", "fists"},
	"kill goriya":               []string{"kill normal"},
	"kill goriya (pit)":         []string{"kill goriya", "pit kill normal"},
	"kill aquamentus":           []string{"sword", "bombs", "fists", "scent seeds"},
	"kill rope":                 []string{"kill normal"},
	"kill hardhat (pit, throw)": []string{"gale seeds", "sword", "shield", "bombs", "rod", "shovel", "scent seeds", "bracelet"},
	"kill moblin (gap, throw)":  []string{"sword", "bombs", "fists", "bracelet", "scent seeds", "slingshot kill normal", "jump kill normal"},
	"kill gel":                  []string{"sword", "bombs", "fists", "ember seeds", "scent seeds", "slingshot gale seeds", "slingshot mystery seeds"},
	"kill facade":               []string{"bombs"},
	"kill beetle":               []string{"kill normal"},
	// TODO: required enemies after d2

	// enemies not required to kill until later
	"kill moldorm": []string{"sword", "bombs", "fists", "scent seeds"},
}
