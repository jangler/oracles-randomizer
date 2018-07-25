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
//   - boomerangs (they work on hardhats!)
//   - seeds (satchel first, slingshot if satchel doesn't work)
//   - rod
//   - bombs
//   - shovel
//   - thrown objects (if applicable)
//   - NOT magnet ball; it kills anything pittable
//   - fool's ore
//   - punch

var killPoints = map[string]Point{
	"slingshot kill normal":   And{"slingshot", "seed kill normal"},
	"jump kill normal":        And{"jump", "kill normal"},
	"jump pit normal":         And{"jump", "pit kill normal"},
	"slingshot gale seeds":    And{"slingshot", "gale seeds"},
	"slingshot mystery seeds": And{"slingshot", "mystery seeds"},
	"kill dodongo":            And{"bombs", "bracelet"},

	// required enemies in normal route-ish order, but with prereqs first
	"seed kill normal":                Or{"ember seeds", "scent seeds", "gale seeds", "mystery seeds"},
	"pop maku bubble":                 Or{"sword", "rod", "seed kill normal", "pegasus slingshot", "bombs", "fool's ore"},
	"remove bush":                     Or{"sword", "boomerang L-2", "ember seeds", "gale slingshot", "bracelet"},
	"kill normal":                     Or{"sword", "bombs", "beams", "seed kill normal", "fool's ore", "punch"},
	"pit kill normal":                 Or{"sword", "beams", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "punch"},
	"kill stalfos":                    Or{"kill normal", "rod"},
	"kill stalfos (throw)":            Or{"kill stalfos", "bracelet"},
	"hit lever":                       Or{"sword", "boomerang", "rod", "ember seeds", "scent seeds", "slingshot", "fool's ore", "punch"},
	"kill goriya bros":                Or{"sword", "bombs", "fool's ore", "punch"},
	"kill goriya":                     Or{"kill normal"},
	"kill goriya (pit)":               Or{"kill goriya", "pit kill normal"},
	"kill aquamentus":                 Or{"sword", "beams", "scent seeds", "bombs", "fool's ore", "punch"},
	"hit far switch":                  Or{"beams", "boomerang", "bombs", "slingshot"},
	"kill rope":                       Or{"kill normal"},
	"kill hardhat (pit, throw)":       Or{"gale seeds", "sword", "beams", "boomerang", "shield", "scent seeds", "rod", "bombs", "shovel", "fool's ore", "bracelet"},
	"kill moblin (gap, throw)":        Or{"sword", "beams", "scent seeds", "slingshot kill normal", "bombs", "fool's ore", "punch", "jump kill normal", "jump pit normal"},
	"kill zol":                        Or{"sword", "beams", "ember seeds", "slingshot gale seeds", "slingshot mystery seeds", "bombs", "fool's ore", "punch"},
	"remove pot":                      Or{"sword L-2", "bracelet"},
	"kill facade":                     Or{"bombs"},
	"flip spiked beetle":              Or{"shield", "shovel"},
	"damage spiked beetle (throw)":    Or{"sword", "bombs", "beams", "seed kill normal", "bracelet", "fool's ore"},
	"flip kill spiked beetle (throw)": And{"flip spiked beetle", "damage spiked beetle (throw)"},
	"gale kill spiked beetle":         And{"gale seeds"},
	"kill spiked beetle (throw)":      Or{"flip kill spiked beetle (throw)", "gale kill spiked beetle"},
	"kill mimic":                      Or{"kill normal"},
	"damage omuai":                    Or{"sword", "bombs", "fool's ore", "punch"},
	"kill omuai":                      And{"damage omuai", "bracelet"},
	"damage mothula":                  Or{"sword", "bombs", "scent seeds", "fool's ore", "punch"},
	"kill mothula":                    And{"damage mothula", "jump"}, // you will basically die without feather
	"kill shrouded stalfos (throw)":   Or{"kill stalfos", "bracelet"},
	"kill like-like (pit, throw)":     Or{"kill normal", "bracelet", "rod", "shovel"},
	"kill water tektite (throw)":      Or{"kill normal", "bracelet"},
	"damage agunima":                  Or{"sword", "scent seeds", "bombs", "fool's ore", "punch"},
	"kill agunima":                    And{"ember seeds", "damage agunima"},
	"hit very far lever":              Or{"boomerang L-2", "slingshot"},
	"hit lever gap":                   Or{"sword", "boomerang", "rod", "slingshot", "fool's ore"},
	"jump hit lever":                  And{"jump", "hit lever gap"},
	"long jump hit lever":             And{"long jump", "hit lever"},
	"hit far lever":                   Or{"jump hit lever", "long jump hit lever", "boomerang", "slingshot"},
	"kill wizzrobe (pit, throw)":      Or{"pit kill normal", "bracelet"},
	"kill gohma":                      Or{"scent slingshot", "ember slingshot"},
	"remove mushroom":                 Or{"boomerang L-2", "bracelet"},
	"kill moldorm":                    Or{"sword", "bombs", "punch", "scent seeds"},
	"kill iron mask":                  Or{"sword", "bombs", "beams", "ember seeds", "scent seeds", "fool's ore", "punch"},
	"kill armos":                      Or{"sword", "bombs", "beams", "boomerang L-2", "scent seeds", "fool's ore"},
	"kill darknut (pit)":              Or{"sword", "bombs", "beams", "scent seeds", "fool's ore", "punch", "shield", "rod", "shovel"},
	"kill syger":                      Or{"sword", "bombs", "scent seeds", "fool's ore", "punch"},
	"kill digdogger":                  Or{"magnet gloves"},

	// enemies not required to kill until later
	"remove flower": Or{"sword", "boomerang L-2", "ember seeds", "gale slingshot"},
}
