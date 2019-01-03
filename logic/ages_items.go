package logic

var agesItemNodes = map[string]*Node{
	"shield": Or("wooden shield", "iron shield"),

	"sword":       Or("sword 1", "sword 2"),
	"noble sword": And("sword 1", "sword 2"),

	"bombs": Or(And("bombs, 10", Or("break bush", "flute", "shovel")),
		Hard("goron shooting gallery")),

	"switch hook": Or("switch hook 1", "switch hook 2"),
	"long hook":   And("switch hook 1", "switch hook 2"),

	"ricky's flute":   Root(),
	"dimitri's flute": Root(),
	"moosh's flute":   Root(),
	"flute":           Or("ricky's flute", "dimitri's flute", "moosh's flute"),

	// ring things
	"fist ring":     Root(),
	"expert's ring": Root(),
	"energy ring":   Root(),

	// expert's ring can do some things that fist ring can't, so this is for
	// the lowest common denominator.
	"punch object": Or("fist ring", "expert's ring"),
	"punch enemy":  Or(Hard("fist ring"), "expert's ring"),

	"harp":   Or("harp 1", "harp 2", "harp 3"),
	"echoes": And("harp"),
	"currents": Or(And("harp 1", Or("harp 2", "harp 3")),
		And("harp 2", "harp 3")),
	"ages": And("harp 1", "harp 2", "harp 3"),

	"bracelet":    Or("bracelet 1", "bracelet 2"),
	"power glove": And("bracelet 1", "bracelet 2"),

	"satchel": Or("satchel 1", "satchel 2"),

	"flippers":     Or("flippers 1", "flippers 2"),
	"mermaid suit": And("flippers 1", "flippers 2"),

	"bomb jump 2": And("feather", Or("pegasus satchel", Hard("bombs"))),
	"jump 3":      And("feather", "pegasus satchel"),
	"bomb jump 3": HardAnd("feather", "pegasus satchel", "bombs"),

	"seed item": Or("satchel", "seed shooter"),

	"ember seeds": And("ember tree seeds"),
	// you can also get scent seeds from ramrock, but the requirements for
	// those are a superset of the requirements for the D3 ones.
	"scent seeds": Or("scent tree seeds",
		HardAnd("d3 E crystal", "seed item")),
	"pegasus seeds": And("pegasus tree seeds"),
	"gale seeds":    And("gale tree seeds"),
	"mystery seeds": And("mystery tree seeds"),

	"ember satchel":   And("ember seeds", "satchel"),
	"scent satchel":   And("scent seeds", "satchel"),
	"pegasus satchel": And("pegasus seeds", "satchel"),
	"gale satchel":    And("gale seeds", "satchel"),
	"mystery satchel": And("mystery seeds", "satchel"),

	"ember shooter":   And("ember seeds", "seed shooter"),
	"scent shooter":   And("scent seeds", "seed shooter"),
	"pegasus shooter": And("pegasus seeds", "seed shooter"),
	"gale shooter":    And("gale seeds", "seed shooter"),
	"mystery shooter": And("mystery seeds", "seed shooter"),
	"any seed shooter": And("seed shooter", Or("ember seeds", "scent seeds",
		"pegasus seeds", "gale seeds", "mystery seeds")),
}
