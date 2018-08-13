package prenode

var portalPrenodes = map[string]*Prenode{
	"rosa portal": Or("temple",
		And("suburbs", "remove bush")),

	"swamp portal": Or("beach",
		And("south swamp", "bracelet")),

	// jump added since it's effectively useless otherwise
	"mountain portal": And("jump", Or("mount cucco", "hide and seek")),

	// TODO maybe can use dimitri from d5 stump area?
	"lake portal": Or("furnace", And("north horon stump", Or(
		And("wet eyeglass lake", Or("jump", "animal flute"), "flippers"),
		And(Or("north horon default winter", "winter"), "pegasus jump L-2")))),

	"village portal": Or(
		And("horon village", "boomerang L-2"),
		And("horon village", "pegasus jump L-2"),
		And("pirate house", "hit lever")),

	// effectively one-way
	"remains portal": And("temple remains", Or("temple remains default winter", "winter"), Or(
		And("shovel", "remove bush", "pegasus jump L-2"),
		And(Or("temple remains default spring", "spring"),
			"remove flower", "remove bush", "pegasus jump L-2", "winter"),
		And(Or("temple remains default summer", "summer"),
			"remove bush", "pegasus jump L-2", "winter"),
		And(Or("temple remains default autumn", "autumn"),
			"remove bush", "jump", "winter"))),

	// dead end
	"d8 portal": And("remains portal", "bombs",
		Or("temple remains default summer", "summer"),
		Or("pegasus jump L-2", And("long jump", "magnet gloves"))),
}
