package main

// these nodes do not define items, only which items can kill which enemies
// under what circumstances, assuming that you've arrived in the room
// containing the enemy.
//
// anything that can be destroyed in more than one way is also included in
// here. bushes, flowers, mushrooms, etc.
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

var killNodesAnd = map[string]Point{
	"slingshot kill normal":   And{"slingshot", "seed kill normal"},
	"jump kill normal":        And{"jump", "kill normal"},
	"jump pit normal":         And{"jump", "pit kill normal"},
	"slingshot gale seeds":    And{"slingshot", "gale seeds"},
	"slingshot mystery seeds": And{"slingshot", "mystery seeds"},
	"kill dodongo":            And{"bombs", "bracelet"},
}

var killNodesOr = map[string]Point{
	// required enemies in normal route-ish order, but with prereqs first
	"seed kill normal":          Or{"ember seeds", "scent seeds", "gale seeds", "mystery seeds"},
	"pop maku bubble":           Or{"sword", "rod", "seed kill normal", "pegasus slingshot", "bombs", "fool's ore"},
	"remove bush":               Or{"sword", "boomerang L-2", "ember seeds", "gale slingshot", "bombs", "bracelet"},
	"kill normal":               Or{"sword", "bombs", "beams", "seed kill normal", "fool's ore", "punch"},
	"pit kill normal":           Or{"sword", "beams", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "punch"},
	"kill normal (pit)":         Or{"kill normal", "pit kill normal"},
	"kill stalfos":              Or{"kill normal", "rod"},
	"kill stalfos (throw)":      Or{"kill stalfos", "bracelet"},
	"hit lever":                 Or{"sword", "boomerang", "rod", "ember seeds", "scent seeds", "fool's ore", "punch"},
	"kill goriya bros":          Or{"sword", "bombs", "fool's ore", "punch"},
	"kill goriya":               Or{"kill normal"},
	"kill goriya (pit)":         Or{"kill goriya", "pit kill normal"},
	"kill aquamentus":           Or{"sword", "beams", "scent seeds", "bombs", "fool's ore", "punch"},
	"kill rope":                 Or{"kill normal"},
	"kill hardhat (pit, throw)": Or{"gale seeds", "sword", "beams", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "bracelet"},
	"kill moblin (gap, throw)":  Or{"sword", "beams", "scent seeds", "slingshot kill normal", "bombs", "fool's ore", "punch", "jump kill normal", "jump pit normal"},
	"kill gel":                  Or{"sword", "beams", "ember seeds", "slingshot gale seeds", "slingshot mystery seeds", "bombs", "fool's ore", "punch"},
	"kill facade":               Or{"bombs"},
	"kill beetle":               Or{"kill normal"},

	// enemies not required to kill until later
	"remove flower":   Or{"sword", "boomerang L-2", "ember seeds", "gale slingshot", "bombs"},
	"remove mushroom": Or{"boomerang L-2", "bracelet"},
	"kill moldorm":    Or{"sword", "bombs", "punch", "scent seeds"},
}
