package main

import (
	"strings"
)

// map internal names to descriptive names for log file

var niceNames = map[string]string{
	// seeds
	"ember tree seeds":   "ember seeds",
	"mystery tree seeds": "mystery seeds",
	"scent tree seeds":   "scent seeds",
	"pegasus tree seeds": "pegasus seeds",
	"gale tree seeds":    "gale seeds",

	// equip items
	"cane":          "cane of somaria",
	"sword 1":       "wooden/noble sword",
	"sword 2":       "wooden/noble sword",
	"boomerang 1":   "(magic) boomerang",
	"boomerang 2":   "(magic) boomerang",
	"spring":        "rod of spring",
	"summer":        "rod of summer",
	"autumn":        "rod of autumn",
	"winter":        "rod of winter",
	"magnet gloves": "magnetic gloves",
	"harp 1":        "tune of echoes/currents/ages",
	"harp 2":        "tune of echoes/currents/ages",
	"harp 3":        "tune of echoes/currents/ages",
	"switch hook 1": "switch/long hook",
	"switch hook 2": "switch/long hook",
	"slingshot 1":   "(hyper) slingshot",
	"slingshot 2":   "(hyper) slingshot",
	"bracelet":      "power bracelet",
	"bracelet 1":    "power bracelet/glove",
	"bracelet 2":    "power bracelet/glove",
	"feather 1":     "roc's feather/cape",
	"feather 2":     "roc's feather/cape",
	"satchel 1":     "seed satchel",
	"satchel 2":     "seed satchel",

	// collection items
	"flippers":         "zora's flippers",
	"flippers 1":       "zora's flippers / mermaid suit",
	"flippers 2":       "zora's flippers / mermaid suit",
	"star ore":         "star-shaped ore",
	"rare peach stone": "piece of heart",
	"goron letter":     "letter of introduction",
	"slate 1":          "slate",
	"slate 2":          "slate",
	"slate 3":          "slate",
	"slate 4":          "slate",

	// seasons slots
	"d0 sword chest":  "hero's cave sword chest",
	"d0 rupee chest":  "hero's cave rupee chest",
	"blaino prize":    "blaino's gym",
	"member's shop 1": "member's shop, 300 rupees",
	"member's shop 2": "member's shop, 300 rupees",
	"member's shop 3": "member's shop, 200 rupees",

	// ages slots
	"ridge base chest":    "ridge west top present",
	"goron diamond chest": "ridge hook cave present",
	"ridge west cave":     "ridge base west present",
	"ridge bush cave":     "ridge past bush cave",
	"ridge base past":     "ridge base west past",
}

// get a user-friendly equivalent of the given internal item or slot name.
func getNiceName(name string) string {
	if name := niceNames[name]; name != "" {
		return name
	}

	if name[0] == 'd' && name[2] == ' ' {
		name = "D" + name[1:]
	}
	name = strings.Replace(name, "map chest", "dungeon map chest", 1)
	name = strings.Replace(name, "gasha chest", "gasha seed chest", 1)

	return name
}
