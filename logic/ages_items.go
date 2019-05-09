package logic

var agesItemNodes = map[string]*Node{
	"shield": Or("wooden shield", "iron shield"),

	"bombs": Or(And("bombs, 10", Or("break bush", "flute", "shovel")),
		And("hard", Or("d2 boss", "goron shooting gallery"))),

	"ricky's flute":   Root(),
	"dimitri's flute": Root(),
	"moosh's flute":   Root(),
	"flute":           Or("ricky's flute", "dimitri's flute", "moosh's flute"),

	// expert's ring can do some things that fist ring can't, so this is for
	// the lowest common denominator.
	"punch object": Or("fist ring", "expert's ring"),
	"punch enemy":  Or(And("hard", "fist ring"), "expert's ring"),

	// progressives
	"noble sword":  Count(2, "sword"),
	"long hook":    Count(2, "switch hook"),
	"echoes":       Count(1, "harp"),
	"currents":     Count(2, "harp"),
	"ages":         Count(3, "harp"),
	"power glove":  Count(2, "bracelet"),
	"mermaid suit": Count(2, "flippers"),

	"bomb jump 2": And("feather", Or("pegasus satchel", And("hard", "bombs"))),
	"jump 3":      And("feather", "pegasus satchel"),
	"bomb jump 3": And("hard", "feather", "pegasus satchel", "bombs"),

	"seed item": Or("satchel", "seed shooter"),

	"ember seeds": And("ember tree seeds"),
	// you can also get scent seeds from ramrock, but the requirements for
	// those are a superset of the requirements for the D3 ones.
	"scent seeds": Or("scent tree seeds",
		And("hard", "d3 E crystal", "seed item")),
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
