package randomizer

// map internal names to descriptive names for log file

var commonNiceNames = map[string]string{
	// seeds
	"ember tree seeds":   "ember seeds",
	"mystery tree seeds": "mystery seeds",
	"scent tree seeds":   "scent seeds",
	"pegasus tree seeds": "pegasus seeds",
	"gale tree seeds":    "gale seeds",

	// items
	"sword":   "wooden/noble sword",
	"satchel": "seed satchel",
}

var seasonsNiceNames = map[string]string{
	// items
	"boomerang":     "(magic) boomerang",
	"spring":        "rod of spring",
	"summer":        "rod of summer",
	"autumn":        "rod of autumn",
	"winter":        "rod of winter",
	"magnet gloves": "magnetic gloves",
	"slingshot":     "(hyper) slingshot",
	"bracelet":      "power bracelet",
	"feather":       "roc's feather/cape",
	"flippers":      "zora's flippers",
	"star ore":      "star-shaped ore",

	// checks
	"d0 key chest":   "hero's cave key chest",
	"d0 sword chest": "hero's cave sword chest",
	"d0 rupee chest": "hero's cave rupee chest",
	"blaino prize":   "blaino's gym",
}

var agesNiceNames = map[string]string{
	// items
	"cane":         "cane of somaria",
	"harp":         "tune of echoes/currents/ages",
	"switch hook":  "switch/long hook",
	"bracelet":     "power bracelet/glove",
	"feather":      "roc's feather",
	"flippers":     "zora's flippers / mermaid suit",
	"goron letter": "letter of introduction",

	// checks
	"ridge base chest":    "ridge west top present",
	"goron diamond chest": "ridge hook cave present",
	"ridge west cave":     "ridge base west present",
	"ridge bush cave":     "ridge past bush cave",
	"ridge base past":     "ridge base west past",
}

// get a user-friendly equivalent of the given internal item or slot name.
func getNiceName(name string, game int) string {
	if name := commonNiceNames[name]; name != "" {
		return name
	}
	niceNames := sora(game, seasonsNiceNames, agesNiceNames).(map[string]string)
	if name := niceNames[name]; name != "" {
		return name
	}

	if name[0] == 'd' && (len(name) == 2 || name[2] == ' ') {
		name = "D" + name[1:]
	}

	return name
}

// turn a spoiler log name into an internal name.
func ungetNiceName(name string, game int) string {
	if name[0] == 'D' && (len(name) == 2 || name[2] == ' ') {
		name = "d" + name[1:]
	}

	reverseNiceNames := make(map[string]string)
	for k, v := range commonNiceNames {
		reverseNiceNames[v] = k
	}
	gameNiceNames := sora(
		game, seasonsNiceNames, agesNiceNames).(map[string]string)
	for k, v := range gameNiceNames {
		reverseNiceNames[v] = k
	}
	if v, ok := reverseNiceNames[name]; ok {
		return v
	}

	return name
}
