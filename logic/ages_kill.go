package logic

// remember to test:
// - sword
// - from satchel, then seed shooter if satchel doesn't work:
//   - ember seeds
//   - scent seeds
//   - gale seeds
// - cane
// - switch hook
// - thrown objects, if applicable
// - pit items, if applicable (shield, shovel, boomerang)
// - boomerang
// - flute
// - shovel

var agesKillNodes = map[string]*Node{
	"break crystal": Or("sword", "bombs", "bracelet"),
	"break pot":     Or("bracelet", "switch hook", "noble sword"),

	// in seasons, shovel hits levers. not in ages, apparently.
	"hit lever": Or("sword", "ember seeds", "scent seeds", "mystery seeds",
		"any seed shooter", "switch hook", "boomerang"),
	"hit switch": Or("sword", "bombs", "ember seeds", "scent seeds",
		"mystery seeds", "any seed shooter", "switch hook", "boomerang"),
	"hit switch ranged": Or("bombs", "any seed shooter", "switch hook", "boomerang"),

	// flute isn't included here since it's only available in some places.
	"break bush safe": Or("sword", "switch hook", "bracelet",
		"bombs", "ember seeds", "gale shooter"),
	"break bush": Or("sword", "switch hook", "bracelet"),

	"satchel weapon": And("satchel",
		Or("ember seeds", HardOr("scent seeds", "gale seeds"))),
	"shooter weapon": And("seed shooter",
		Or("ember seeds", "scent seeds", "gale seeds")),

	// most enemies are vulnerable to these items
	"kill normal": Or("sword", "satchel weapon", "shooter weapon", "cane",
		Hard("bombs")),
	"kill normal ranged": Or("shooter weapon", And("cane", "bracelet"),
		Hard("bombs")),
	"kill underwater": Or("sword", "shooter weapon"),
	"pit normal":      Or("shield", And("boomerang", "shovel")),

	"kill gel":     Or("kill normal", "switch hook", "boomerang", "shovel"),
	"kill stalfos": Or("kill normal"),
	"kill zol":     Or("kill normal", "switch hook"),
	"kill ghini":   Or("kill normal", "switch hook"),
	"kill giant ghini": Or("sword", "scent shooter", "switch hook",
		HardOr("bombs", "scent satchel")),
	"kill pumpkin head": And("bracelet",
		Or("sword", "ember seeds", "scent shooter",
			HardOr("bombs", "scent satchel"))),

	"kill spiked beetle": Or("gale shooter", Hard("gale satchel"),
		And(Or("shield", "shovel"), Or("kill normal", "switch hook"))),
	"kill swoop": Or("sword", "scent shooter", "switch hook",
		HardOr("bombs", "scent satchel")),

	"kill moldorm": Or("sword", "scent shooter", "cane", "switch hook",
		HardOr("bombs", "scent satchel")),
	"kill armos": Or("scent shooter", "cane", HardOr("bombs", "scent satchel")),

	"kill wizzrobe": Or("kill normal"),
}
