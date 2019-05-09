package logic

// to "enter" a portal means to walk into it; to "exit" means to walk out.
// exits are root nodes since portal connections can be randomized.
var portalNodes = map[string]*Node{
	// pair 1

	"enter eastern suburbs portal": And("suburbs", Or("remove bush", "flute")),
	"exit eastern suburbs portal":  Root(),

	"enter volcanoes east portal": And("temple"),
	"exit volcanoes east portal":  Root(),

	// pair 2

	"enter spool swamp portal": And("south swamp", "bracelet", Or(
		"flute", "spool swamp default summer", "spool swamp default autumn",
		And("spool stump", Or("summer", "autumn")),
		And(Or("spool swamp default winter", And("spool stump", "winter")),
			"shovel"),
		And(Or("spool swamp default spring", And("spool stump", "spring")),
			"remove flower"))),
	"exit spool swamp portal": Root(),

	"enter subrosia market portal": And("beach"),
	"exit subrosia market portal":  Root(),

	// pair 3

	"enter mt. cucco portal": And("mount cucco"),
	"exit mt. cucco portal":  Root(),

	"enter strange brothers portal": And("feather", "hide and seek"),
	"exit strange brothers portal":  Root(),

	// pair 4

	"enter eyeglass lake portal": And("north horon stump", Or(
		And("wet eyeglass lake", Or("jump 2", "ricky's flute", "moosh's flute"),
			Or("flippers", And("dimitri's flute", "bracelet"))),
		And(Or("north horon default winter", "winter"),
			Or("jump 6", And("jump 2", "dimitri's flute"))))),
	"exit eyeglass lake portal": Root(),

	"enter great furnace portal": And("furnace"),
	"exit great furnace portal":  Root(),

	// pair 5

	"enter horon village portal": And("horon village",
		Or("magic boomerang", And("hard", "jump 6"))),
	"exit horon village portal": Root(),

	"enter house of pirates portal": And("pirate house", "hit lever"),
	"exit house of pirates portal":  Root(),

	// pair 6

	"enter temple remains lower portal": And("temple remains lower stump", Or(
		And("bomb temple remains", "feather"),
		And(Or("temple remains default winter", "winter"), Or(
			And("exit temple remains upper portal",
				Or("temple remains default winter", "feather")),
			And("hard", "shovel", "remove bush", "jump 6"),
			And("hard", Or("temple remains default spring", "spring"),
				"remove flower", "remove bush", "jump 6", "winter"),
			And("hard", Or("temple remains default summer", "summer"),
				"remove bush", "jump 6", "winter"),
			And(Or("temple remains default autumn", "autumn"),
				"remove bush", "jump 2", "winter"))))),
	"exit temple remains lower portal": Root(),

	"enter volcanoes west portal": Or(), // not possible without exiting first
	"exit volcanoes west portal":  Root(),

	// pair 7

	"enter temple remains upper portal": And(
		"temple remains lower stump", "bomb temple remains",
		Or("temple remains default summer", "summer"),
		Or("jump 6", And("bomb jump 2", "magnet gloves"))),
	"exit temple remains upper portal": Root(),

	"enter D8 entrance portal": And("d8 entrance"),
	"exit D8 entrance portal":  Root(),

	// idk where to put this
	"bomb temple remains": And("exit volcanoes west portal", "bombs"),
}
