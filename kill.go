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

// when testing how to kill enemies, remember to try:
// - sword
// - beams
// - boomerang L-1
// - boomerang L-2
// - rod
// - seeds (satchel first, slingshot if satchel doesn't work)
// - bombs
// - thrown objects (if applicable)
// - magnet ball (if applicable)
// - fool's ore
// - punch
// - what pushes them into pits (if applicable)
//   - sword
//   - beams
//   - shield
//   - NOT boomerangs; they only push if they damage already
//   - seeds (satchel first, slingshot if satchel doesn't work)
//   - rod
//   - bombs
//   - shovel
//   - thrown objects (if applicable)
//   - NOT magnet ball; it kills anything pittable
//   - fool's ore
//   - punch

var killNodesAnd = map[string][]string{
	"slingshot kill normal":   []string{"slingshot", "seed kill normal"},
	"jump kill normal":        []string{"jump", "kill normal"},
	"slingshot gale seeds":    []string{"slingshot", "gale seeds"},
	"slingshot mystery seeds": []string{"slingshot", "mystery seeds"},
	"kill dodongo":            []string{"bombs", "bracelet"},
}

var killNodesOr = map[string][]string{
	// enemies in normal route-ish order, but with prereqs first
	"seed kill normal":          []string{"ember seeds", "scent seeds", "gale seeds", "mystery seeds"},
	"kill normal":               []string{"sword", "bombs", "beams", "seed kill normal", "fool's ore", "punch"},
	"pit kill normal":           []string{"sword", "beams", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "punch"},
	"kill normal (pit)":         []string{"kill normal", "pit kill normal"},
	"kill stalfos":              []string{"kill normal", "rod"},
	"kill stalfos (throw)":      []string{"kill stalfos", "bracelet"},
	"kill goriya bros":          []string{"sword", "bombs", "fool's ore", "punch"},
	"kill goriya":               []string{"kill normal"},
	"kill goriya (pit)":         []string{"kill goriya", "pit kill normal"},
	"kill aquamentus":           []string{"sword", "beams", "scent seeds", "bombs", "fool's ore", "punch"},
	"kill rope":                 []string{"kill normal"},
	"kill hardhat (pit, throw)": []string{"gale seeds", "sword", "beams", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "bracelet"},
	"kill moblin (gap, throw)":  []string{"sword", "beams", "scent seeds", "slingshot kill normal", "bombs", "fool's ore", "punch"},
	"kill gel":                  []string{"sword", "beams", "ember seeds", "slingshot gale seeds", "slingshot mystery seeds", "bombs", "fool's ore", "punch"},
	"kill facade":               []string{"bombs"},
	"kill beetle":               []string{"kill normal"},

	// enemies not required to kill until later
	"kill moldorm": []string{"sword", "bombs", "punch", "scent seeds"},
}
