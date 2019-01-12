package logic

var portalNodes = map[string]*Node{
	"rosa portal": Or("temple", And("suburbs", Or("remove bush", "flute"))),

	"swamp portal": Or("beach",
		And("south swamp", "bracelet", Or("flute",
			"spool swamp default summer", "spool swamp default autumn",
			And("spool stump", Or("summer", "autumn")),
			And(Or("spool swamp default winter", And("spool stump", "winter")),
				"shovel"),
			And(Or("spool swamp default spring", And("spool stump", "spring")),
				"remove flower")))),

	// jump added since it's effectively useless otherwise
	"mountain portal": And("jump 2", Or("mount cucco", "hide and seek")),

	"lake portal": Or("furnace", And("north horon stump", Or(
		And("wet eyeglass lake", Or("jump 2", "ricky's flute", "moosh's flute"),
			Or("flippers", And("dimitri's flute", "bracelet"))),
		And(Or("north horon default winter", "winter"),
			Or("jump 6", And("jump 2", "dimitri's flute")))))),

	"village portal": Or(
		And("horon village", Or("boomerang L-2", Hard("jump 6"))),
		And("pirate house", "hit lever")),

	// effectively one-way
	"remains portal": And("temple remains",
		Or("temple remains default winter", "winter"), Or(
			HardAnd("shovel", "remove bush", "jump 6"),
			HardAnd(Or("temple remains default spring", "spring"),
				"remove flower", "remove bush", "jump 6", "winter"),
			HardAnd(Or("temple remains default summer", "summer"),
				"remove bush", "jump 6", "winter"),
			And(Or("temple remains default autumn", "autumn"),
				"remove bush", "jump 2", "winter"))),

	// dead end
	"d8 portal": And("remains portal", "bombs",
		Or("temple remains default summer", "summer"),
		Or("jump 6", And("bomb jump 2", "magnet gloves"))),
}
