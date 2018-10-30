package logic

// remember to test:
// - sword
// - bombs (*usually* don't require except with the bomb bag)
// - from satchel, then seed shooter if satchel doesn't work:
//   - ember seeds
//   - scent seeds
//   - gale seeds
//   - mystery seeds (don't require except with the satchel upgrade)
// - cane
// - switch hook
// - thrown objects, if applicable
// - pit items, if applicable (shield, shovel, boomerang)
// - boomerang
// - flute
// - shovel

// TODO some of this should be in hard logic, like scent seeds from satchel

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
	"break bush": Or("sword", "switch hook", "bracelet",
		HardOr("bombs", "ember seeds", "gale shooter")),

	"bomb weapon": Hard("bombs"),
	"satchel weapon": And("satchel",
		Or("ember seeds", HardOr("scent seeds", "gale seeds"))),
	"shooter weapon": And("seed shooter",
		Or("ember seeds", "scent seeds", "gale seeds")),

	// most enemies are vulnerable to these items
	"kill normal": Or("sword", "bomb weapon", "satchel weapon",
		"shooter weapon", "cane"),
	"kill normal ranged": Or("bomb weapon", "shooter weapon",
		And("cane", "bracelet")),
	"kill underwater": Or("sword", "shooter weapon"),
	"pit normal":      Or("shield", And("boomerang", "shovel")),

	"kill gel":     Or("kill normal", "switch hook", "boomerang", "shovel"),
	"kill stalfos": Or("kill normal"),
	"kill zol":     Or("kill normal", "switch hook"),
	"kill ghini":   Or("kill normal", "switch hook"),
	"kill giant ghini": Or("sword", "bomb weapon", "scent shooter",
		HardOr("scent satchel", "mystery seeds"), "switch hook"),
	"kill pumpkin head": And("bracelet", Or("sword", "bomb weapon",
		"ember seeds", "scent seeds")),

	"kill spiked beetle": Or("gale shooter", Hard("gale satchel"),
		And(Or("shield", "shovel"), Or("kill normal", "switch hook"))),
	"kill swoop": Or("sword", "bomb weapon", "scent seeds", "switch hook",
		Hard("mystery seeds")),

	"kill moldorm": Or("sword", "bomb weapon", "scent seeds", "cane",
		"switch hook", Hard("mystery seeds")),
	"kill armos": Or("bomb weapon", "scent seeds", "cane",
		Hard("mystery seeds")),

	"kill wizzrobe": Or("kill normal"),
}
