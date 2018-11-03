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
	"shop shield L-1": "wooden shield",
	"shield L-2":      "iron shield",
	"cane":            "cane of somaria",
	"sword 1":         "wooden/noble sword",
	"sword 2":         "wooden/noble sword",
	"boomerang 1":     "(magic) boomerang",
	"boomerang 2":     "(magic) boomerang",
	"spring":          "rod of spring",
	"summer":          "rod of summer",
	"autumn":          "rod of autumn",
	"winter":          "rod of winter",
	"magnet gloves":   "magnetic gloves",
	"harp 1":          "tune of echoes/currents/ages",
	"harp 2":          "tune of echoes/currents/ages",
	"harp 3":          "tune of echoes/currents/ages",
	"switch hook 1":   "switch/long hook",
	"switch hook 2":   "switch/long hook",
	"slingshot 1":     "(hyper) slingshot",
	"slingshot 2":     "(hyper) slingshot",
	"bracelet":        "power bracelet",
	"bracelet 1":      "power bracelet/glove",
	"bracelet 2":      "power bracelet/glove",
	"feather 1":       "roc's feather/cape",
	"feather 2":       "roc's feather/cape",
	"satchel 1":       "seed satchel",
	"satchel 2":       "seed satchel",

	// collection items
	"flippers":         "zora's flippers",
	"flippers 1":       "zora's flippers / mermaid suit",
	"flippers 2":       "zora's flippers / mermaid suit",
	"star ore":         "star-shaped ore",
	"rare peach stone": "piece of heart",
	"goron letter":     "letter of introduction",

	// north horon / holodrum plain / eyeglass lake slots
	"scent tree":          "north horon seed tree",
	"blaino gift":         "blaino's gym",
	"round jewel gift":    "old man in treehouse",
	"lake chest":          "eyeglass lake, across bridge",
	"water cave chest":    "cave south of mrs. ruul",
	"mushroom cave chest": "cave north of D1",
	"dry lake east chest": "dry eyeglass lake, east cave",
	"dry lake west chest": "dry eyeglass lake, west cave",

	// horon village slots
	"maku tree gift":   "maku tree",
	"ember tree":       "horon village seed tree",
	"village shop 1":   "shop, 20 rupees",
	"village shop 2":   "shop, 30 rupees",
	"village shop 3":   "shop, 150 rupees",
	"member's shop 1":  "member's shop, 300 rupees",
	"member's shop 2":  "member's shop, 300 rupees",
	"member's shop 3":  "member's shop, 200 rupees",
	"village SW chest": "horon village SW chest",
	"village SE chest": "horon village SE chest",

	// western coast / graveyard slots
	"x-shaped jewel chest": "black beast's chest",
	"western coast chest":  "western coast, beach chest",
	"coast house chest":    "western coast, in house",

	// eastern suburbs / woods of winter slots
	"moblin road chest":  "woods of winter, 1st cave",
	"linked dive chest":  "woods of winter, 2nd cave",
	"shovel gift":        "holly's house",
	"mystery tree":       "woods of winter seed tree",
	"outdoor d2 chest":   "chest on top of D2",
	"mystery cave chest": "cave outside D2",
	"moblin cliff chest": "eastern suburbs, on cliff",

	// spool swamp slots
	"pegasus tree":       "spool swamp seed tree",
	"floodgate key spot": "floodgate keeper's house",
	"square jewel chest": "spool swamp cave",

	// natzu slots
	"platform chest":     "natzu region, across water",
	"great moblin chest": "moblin keep",

	// sunken city / mount cucco / goron mountain slots
	"sunken gale tree":      "sunken city seed tree",
	"master's plaque chest": "master diver's challenge",
	"diver gift":            "master diver's reward",
	"diver chest":           "chest in master diver's cave",
	"talon cave chest":      "mt. cucco, talon's cave",
	"dragon key spot":       "goron mountain, across pits",
	"pyramid jewel spot":    "diving spot outside D4",
	"sunken cave chest":     "sunken city, summer cave",
	"goron chest":           "chest in goron mountain",

	// tarm ruins / lost woods slots
	"noble sword spot": "lost woods",
	"tarm gale tree":   "tarm ruins seed tree",
	"tarm gasha chest": "tarm ruins, under tree",

	// samasa desert slots
	"desert pit":   "samasa desert pit",
	"desert chest": "samasa desert chest",

	// subrosia slotss
	"dance hall prize":     "subrosian dance hall",
	"rod gift":             "temple of seasons",
	"spring tower":         "tower of spring",
	"summer tower":         "tower of summer",
	"autumn tower":         "tower of autumn",
	"winter tower":         "tower of winter",
	"star ore spot":        "subrosia seaside",
	"subrosian market 1":   "subrosia market, 1st item",
	"subrosian market 2":   "subrosia market, 2nd item",
	"subrosian market 5":   "subrosia market, 5th item",
	"non-rosa gasha chest": "subrosia, open cave",
	"rosa gasha chest":     "subrosia, locked cave",
	"red ore chest":        "subrosia village chest",
	"blue ore chest":       "subrosian wilds chest",
	"hard ore slot":        "great furnace",
	"iron shield gift":     "subrosian smithy",

	// seasons dungeon slots
	"d0 sword chest":         "hero's cave sword chest",
	"d0 rupee chest":         "hero's cave rupee chest",
	"d1 satchel spot":        "D1 seed satchel spot",
	"d2 bracelet chest":      "D2 power bracelet chest",
	"d3 feather chest":       "D3 roc's feather chest",
	"d5 magnet gloves chest": "D5 magnetic gloves chest",
	"d6 boomerang chest":     "D6 magic boomerang chest",
	"d7 cape chest":          "D7 roc's cape chest",
	"d8 HSS chest":           "D8 hyper slingshot chest",

	// labrynna slots
	"crescent seafloor cave": "under crescent island",
	"ridge base chest":       "ridge base east present",
	"goron diamond chest":    "ridge hook cave present",
	"ridge west cave":        "ridge base west present",
	"ridge bush cave":        "ridge past bush cave",
	"ridge base past":        "ridge base west past",
	"zora cave past":         "fisher's island cave",
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
