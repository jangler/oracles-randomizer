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
// - punch

var agesKillNodes = map[string]*Node{
	"break crystal": Or("sword", "bombs", "bracelet", "ember seeds",
		"expert's ring"),
	"break pot": Or("bracelet", "switch hook", "noble sword"),

	// obviously this only works on standard enemies
	"push enemy": Or("shield",
		And("shovel", Or("boomerang", "pegasus shooter"))),

	// unlike in seasons, shovel doesn't hit levers.
	"hit lever": Or("sword", "ember seeds", "scent seeds", "mystery seeds",
		"any seed shooter", "switch hook", "boomerang", "punch object"),
	// and seeds from minecart don't hit levers. not sure if this is because
	// the ones in seasons are horizontal and the one in ages D2 is vertical.
	"hit lever from minecart": Or("sword", "any seed shooter", "boomerang"),
	"hit switch": Or("sword", "bombs", "punch object", "ember seeds",
		"scent seeds", "mystery seeds", "any seed shooter", "switch hook",
		"boomerang"),
	"hit switch ranged": Or("bombs", "any seed shooter", "switch hook",
		"boomerang", And("sword", "energy ring")),

	// flute isn't included here since it's only available in some places.
	"break bush safe": Or("sword", "switch hook", "bracelet",
		"bombs", "ember seeds", "gale shooter"),
	"break bush": Or("sword", "switch hook", "bracelet"),

	"satchel weapon": And("satchel",
		Or("ember seeds", And("hard", Or("scent seeds", "gale seeds")))),
	"shooter weapon": And("seed shooter",
		Or("ember seeds", "scent seeds", "gale seeds")),

	// most enemies are vulnerable to these items
	"kill normal": Or("sword", "satchel weapon", "shooter weapon", "cane",
		"punch enemy", And("hard", "bombs")),
	"kill normal ranged": Or("shooter weapon", And("cane", "bracelet"),
		And("hard", "bombs")),
	"kill underwater":  Or("sword", "shooter weapon", "punch enemy"),
	"kill switch hook": Or("kill normal", "switch hook"),

	"kill giant ghini": Or("sword", "scent shooter", "switch hook",
		"punch enemy", And("hard", Or("bombs", "scent satchel"))),
	"kill pumpkin head": And("bracelet",
		Or("sword", "punch enemy", "ember seeds", "scent shooter",
			And("hard", Or("bombs", "scent satchel")))),

	// spiked beetles can't be punched for some reason
	"kill spiked beetle": Or("gale shooter", And("hard", "gale satchel"),
		And(Or("shield", "shovel"), Or("sword", "satchel weapon",
			"shooter weapon", "cane", And("hard", "bombs"), "switch hook"))),
	"kill swoop": Or("sword", "scent shooter", "switch hook", "punch enemy",
		And("hard", Or("bombs", "scent satchel"))),

	"kill moldorm": Or("sword", "scent shooter", "cane", "switch hook",
		"punch enemy", And("hard", Or("bombs", "scent satchel"))),
	"kill subterror": And("shovel", Or("sword", "switch hook", "scent seeds",
		"punch enemy", And("hard", "bombs"))),
}
