package logic

// "ricky", "dimitri", and "moosh" refer to accessing those animal companions
// in their designated regions (e.g. dimitri in sunken city). "x's flute" means
// being able to call the animal in general.

var holodrumNodes = map[string]*Node{
	"start": And(), // parent for nodes reachable by default

	// horon village
	"horon village": And("start"),
	"maku tree":     AndSlot("horon village", "sword"),
	"horon village seed tree": AndSlot("horon village", "seed item",
		Or("harvest tree", "dimitri's flute", Hard("remove bush"))),
	"horon village SE chest": AndSlot("horon village", "bombs"),
	"horon village SW chest": AndSlot("horon village",
		Or("remove mushroom", "dimitri's flute")),
	"shop, 20 rupees":  AndSlot("start"),
	"shop, 30 rupees":  AndSlot("start"),
	"shop, 150 rupees": AndSlot("start"),
	"member's shop 1":  AndSlot("member's card"),
	"member's shop 2":  AndSlot("member's card"),
	"member's shop 3":  AndSlot("member's card"),

	// western coast
	//
	// possible but not in logic: reaching the stump without feather
	"black beast's chest": AndSlot("horon village",
		Or("ember slingshot", Hard("mystery slingshot")),
		"mystery seeds", "kill moldorm"),
	"enter d0":    And("horon village"),
	"pirate ship": And("pirate's bell", "pirate house"),
	"coast stump": And("pirate ship", "bombs", "jump 2"),
	"enter d7": And("pirate ship",
		Or("jump 3", "western coast default summer",
			And("coast stump", "summer")),
		Or("shovel", "western coast default spring",
			"western coast default summer",
			"western coast default autumn",
			And("coast stump", Or("spring", "summer", "autumn")))),
	"western coast, beach chest": AndSlot("pirate ship"),
	"western coast, in house":    AndSlot("pirate ship"),

	// eastern suburbs
	"suburbs": Or( // this is the area south of the pool by sokra's stump
		And("horon village", "ember seeds"),
		And("rosa portal", Or("remove bush", "flute")),
		And("fairy fountain", Or("eastern suburbs default winter", "winter",
			"flippers", "jump 2", "ricky's flute", "dimitri's flute"))),
	"fairy fountain": Or(
		And("sunken city",
			Or("eastern suburbs default spring", "spring", "gale satchel")),
		And("suburbs", Or("eastern suburbs default winter", "winter",
			"flippers", "jump 2", "ricky's flute", "dimitri's flute"))),
	"moblin road": Or(
		And("fairy fountain", Or("eastern suburbs default winter", "winter")),
		And("sunken city", "flippers", Or(
			"sunken city default spring", "spring",
			"sunken city default summer", "summer",
			"sunken city default autumn", "autumn"),
			Or("gale satchel", And(
				Or("eastern suburbs default winter", "winter"),
				Or("eastern suburbs default spring", "spring"))))),
	"holly's house": AndSlot("moblin road",
		Or("woods of winter default winter", "winter")),
	"central woods of winter": And("fairy fountain", Or(
		"shovel", "jump 2", "flute", "spring", "summer", "autumn",
		And("flippers", Or(
			"eastern suburbs default spring",
			"eastern suburbs default summer",
			"eastern suburbs default autumn")))),
	"woods of winter seed tree": AndSlot("central woods of winter",
		"seed item", Or("harvest tree", "dimitri's flute")),
	"enter d2 A": And("central woods of winter", Or("remove bush", "flute")),
	"enter d2 B": Or(
		And("central woods of winter", "bracelet",
			Or("woods of winter default summer", "ricky's flute")),
		And("d2 blade chest", "bracelet")),
	"chest on top of D2": AndSlot("enter d2 B"),
	"cave outside D2": AndSlot("central woods of winter",
		Or("remove mushroom", "dimitri's flute"),
		Or("jump 4", "magnet gloves"),
		Or("woods of winter default autumn", And("autumn", Or("ricky's flute",
			"woods of winter default summer", "enter d2 B")))),
	"woods of winter, 1st cave": AndSlot("moblin road",
		Or("bombs", "ricky's flute"), "remove bush safe",
		Or("woods of winter default spring", "spring",
			"woods of winter default summer", "summer",
			"woods of winter default autumn", "autumn")),
	"eastern suburbs, on cliff": AndSlot("suburbs", "bracelet",
		Or("bomb jump 2", "magnet gloves"),
		Or("eastern suburbs default spring", "spring")),
	"woods of winter, 2nd cave": AndSlot("moblin road",
		Or("flippers", "bomb jump 3")),

	// holodrum plain
	"ghastly stump": Or("north swamp",
		And("blaino's gym", Or("jump 2", "ricky", "flute",
			And("flippers", "remove bush"), "holodrum plain default winter")),
		And("south swamp", Or("flippers", "dimitri's flute"), "remove bush")),
	"blaino's gym": Or(
		And("ghastly stump", Or("jump 2", "ricky", "flute", "winter",
			"holodrum plain default winter")),
		And("south swamp", Or("flippers", "dimitri's flute")),
		And("sunken city", Or(
			And("natzu prairie", "flute"),
			And("natzu river", "jump 2", Or("flippers", "flute")),
			And("natzu wasteland",
				Or("flute", And("remove bush", "bomb jump 3"))))),
		And("north horon stump", Or("bracelet",
			And(Or("remove bush", "flute"),
				Or("flippers", "dimitri's flute")))),
		And("temple remains", "jump 3"),
		And("goron mountain", "flippers")),
	"north horon seed tree": AndSlot("blaino's gym", "seed item",
		Or("harvest tree", "dimitri's flute")),
	"blaino prize": AndSlot("blaino's gym"),
	"ricky":        Or("ricky's flute"),
	"old man in treehouse": AndSlot("blaino's gym",
		Or("flippers", "dimitri's flute")),
	"cave south of mrs. ruul": AndSlot("blaino's gym", "flippers"),
	"cave north of D1": AndSlot("blaino's gym", "flippers",
		Or("remove mushroom", "dimitri's flute"),
		Or("holodrum plain default autumn", And("ghastly stump", "autumn"))),

	// spool swamp
	"north swamp": And("ghastly stump", Or("holodrum plain default summer",
		"summer", "jump 4", "ricky", "moosh's flute")),
	"spool swamp seed tree": AndSlot("north swamp", "seed item",
		Or("harvest tree", "dimitri's flute")),
	"floodgate keeper's house": AndSlot("north swamp", "hit lever"),
	"spool stump": And("north swamp", "hit lever", "bracelet", "floodgate key",
		Or("pegasus satchel", "flippers", "bomb jump 3")),
	"dry swamp": Or("spool swamp default summer",
		"spool swamp default autumn", "spool swamp default winter",
		And("spool stump", Or("summer", "autumn", "winter"))),
	"south swamp": Or(
		And("spool stump", Or("flippers", "dimitri's flute")),
		And("spool stump", "dry swamp", Or("jump 2", "flute")),
		And("ghastly stump", "remove bush", Or("flippers", "dimitri's flute")),
		And("blaino's gym", Or("flippers", "dimitri's flute")),
		And("swamp portal", "bracelet")),
	"spool swamp cave": AndSlot("south swamp",
		Or("spool swamp default winter", And("spool stump", "winter")),
		Or("shovel", "flute"), Or("bombs", "ricky's flute")),
	"enter d3": And("spool stump", Or("spool swamp default summer", "summer")),

	// north horon / eyeglass lake
	"not north horon default summer": Or("north horon default spring",
		"north horon default autumn", "north horon default winter"),
	"north horon stump": Or(
		And("horon village", Or("remove bush", "flute")),
		And("blaino's gym", "bracelet"),
		And("south swamp", Or("flippers", "dimitri's flute"),
			Or("remove bush", "flute")),
		And("lake portal", "not north horon default summer",
			"flippers", "jump 2"),
		And("lake portal", "jump 6", "north horon default winter")),
	"enter d1": And("gnarled key", Or(
		And("south swamp", Or("flippers", "dimitri's flute")),
		And("north horon stump", Or("remove bush", "flute")))),
	"wet eyeglass lake": Or("not north horon default summer",
		"spring", "autumn", "winter"),
	"d5 stump": Or(
		And("lake portal", "not north horon default summer",
			Or("flippers", And("north horon default winter", "jump 6"))),
		And("north horon stump", Or("jump 2", "ricky's flute", "moosh's flute"),
			Or("north horon default winter", "winter", "flippers",
				And("bracelet", "dimitri's flute")))),
	"enter d5": And("d5 stump", Or("remove mushroom", "dimitri's flute"),
		Or("autumn", And("north horon default autumn",
			Or("lake portal", "jump 2", "ricky's flute", "moosh's flute"),
			Or("flippers", And("dimitri's flute", "bracelet"))))),
	"eyeglass lake, across bridge": AndSlot("horon village", Or("jump 4",
		And("jump 2", Or("north horon default autumn",
			And("autumn", "north horon stump"))))),
	"dry eyeglass lake, east cave": AndSlot("d5 stump", "bracelet",
		Or("summer", And("enter d5", "north horon default summer"))),
	"dry eyeglass lake, west cave": AndSlot(
		Or("bombs", "ricky's flute"), "flippers",
		Or(And("north horon stump", Or("north horon default summer", "summer"),
			Or("jump 2", "ricky's flute", "moosh's flute")),
			And("d5 stump", "summer"),
			And("enter d5", "north horon default summer"))),

	// natzu
	"natzu prairie":   Root("start"),
	"natzu river":     Root(),
	"natzu wasteland": Root(),
	"moblin keep": AndSlot(Or("flippers", "bomb jump 4"),
		"bracelet", Or(
			And("natzu prairie", "sunken city"),
			And("natzu river", "blaino's gym", "dimitri's flute"),
			And("natzu wasteland", "blaino's gym",
				Or("flute", And("remove bush", "jump 3"))))),
	"natzu region, across water": OrSlot(
		And("blaino's gym", Or("flippers", "dimitri's flute")),
		And("sunken city", "natzu river", "jump 6")),

	// sunken city
	"sunken city": Or(
		And("mount cucco", "flippers",
			Or("summer", "sunken city default summer", "gale satchel")),
		And("fairy fountain", Or("eastern suburbs default spring", "spring")),
		And("blaino's gym", Or(
			And("natzu prairie", "flute"),
			And("natzu river", Or(And(Or("flippers", "flute"), "jump 2"),
				And("flute", "flippers", "gale satchel"))),
			And("natzu wasteland", Or("flute", And("remove bush",
				Or("bomb jump 3", And("jump 3", "flippers")))))))),
	"sunken city seed tree": AndSlot("sunken city", "seed item",
		Or("harvest tree", "dimitri"),
		Or("jump 2", "flippers", "dimitri's flute",
			"sunken city default winter")),
	"dimitri": And("sunken city", Or("dimitri's flute",
		And("bombs", Or("jump 2", "flippers", "sunken city default winter")))),
	"master diver's challenge": AndSlot("dimitri", "sword",
		Or("jump 2", "flippers")),
	"master diver's reward": AndSlot("dimitri", "master's plaque"),
	"sunken city, summer cave": AndSlot("sunken city", "flippers",
		"remove bush safe", Or("sunken city default summer", "summer")),
	"chest in master diver's cave": AndSlot("dimitri"),

	// mount cucco
	"mount cucco": Or("mountain portal",
		And("sunken city", "flippers",
			Or("sunken city default summer", "summer")),
		And("goron mountain", "bracelet", "shovel")),
	"spring banana tree": AndSlot("mount cucco", "remove flower", "bracelet",
		"jump 2", Or("sunken city default spring", "spring"),
		Or("sword", "fool's ore")),
	"moosh": And("mount cucco", "spring banana"),
	"goron mountain, across pits": AndSlot("mount cucco",
		Or("moosh", "jump 6", Hard("jump 4"))),
	"mt. cucco, talon's cave": AndSlot("mount cucco",
		Or("sunken city default spring", "spring")),
	"dragon keyhole": And("mt. cucco, talon's cave",
		"winter", "jump 2", "bracelet"),
	"enter d4":               And("dragon key", "dragon keyhole", "summer"),
	"diving spot outside D4": AndSlot("mt. cucco, talon's cave", "flippers"),

	// goron mountain
	"goron mountain": Or(
		And("mount cucco", Or("shovel", "spring banana"), "bracelet"),
		And("temple remains", "jump 3", Or("flippers", "bomb jump 4")),
		And("blaino's gym", "flippers")),
	"chest in goron mountain": AndSlot("goron mountain", "bombs", "bomb jump 3"),

	// tarm ruins
	"tarm ruins": And("north swamp",
		"square jewel", "pyramid jewel", "round jewel", "x-shaped jewel"),
	"lost woods": AndSlot("tarm ruins", "remove mushroom", "winter", "autumn",
		"spring", "summer"),
	"tarm ruins seed tree": AndSlot("lost woods", "seed item", "harvest tree"),
	"enter d6": And("lost woods", "remove bush",
		Or("tarm ruins default winter", "winter"),
		Or("shovel", "ember seeds")),
	"tarm ruins, under tree": AndSlot("lost woods", "remove mushroom",
		"ember seeds", Or("tarm ruins default autumn", "autumn")),

	// samasa desert
	"desert":              And("suburbs", "pirate house"),
	"samasa desert pit":   AndSlot("desert", "bracelet"),
	"samasa desert chest": AndSlot("desert", "flippers"),

	// temple remains (the important logic is in the portal nodes)
	"temple remains": Or(
		And("goron mountain", Or("flippers", "bomb jump 4"), "jump 3"),
		And("blaino's gym", "jump 3")),

	// northern peak
	"maku seed": And("sword", "d1 essence", "d2 essence", "d3 essence",
		"d4 essence", "d5 essence", "d6 essence", "d7 essence", "d8 essence"),
	"enter d9": And("blaino's gym", "maku seed"),

	// old men
	"goron mountain old man": And("goron mountain", "ember seeds"),
	"western coast old man":  And("pirate ship", "ember seeds"),
	"holodrum plain east old man": And("blaino's gym", "ember seeds",
		Or("ricky's flute", "holodrum plain default summer",
			And("ghastly stump", "summer"))),
	"horon village old man":       And("horon village", "ember seeds"),
	"north horon old man":         And("north horon stump", "ember seeds"),
	"tarm ruins old man":          And("enter d6", "ember seeds"),
	"woods of winter old man":     And("holly's house", "ember seeds"),
	"holodrum plain west old man": And("ghastly stump", "ember seeds"),
}

var seasonNodes = map[string]*Node{
	"north horon default spring": Root(),
	"north horon default summer": Root(),
	"north horon default autumn": Root(),
	"north horon default winter": Root("start"),

	"eastern suburbs default spring": Root(),
	"eastern suburbs default summer": Root(),
	"eastern suburbs default autumn": Root("start"),
	"eastern suburbs default winter": Root(),

	"woods of winter default spring": Root(),
	"woods of winter default summer": Root("start"),
	"woods of winter default autumn": Root(),
	"woods of winter default winter": Root(),

	"spool swamp default spring": Root(),
	"spool swamp default summer": Root(),
	"spool swamp default autumn": Root("start"),
	"spool swamp default winter": Root(),

	"holodrum plain default spring": Root("start"),
	"holodrum plain default summer": Root(),
	"holodrum plain default autumn": Root(),
	"holodrum plain default winter": Root(),

	"sunken city default spring": Root(),
	"sunken city default summer": Root("start"),
	"sunken city default autumn": Root(),
	"sunken city default winter": Root(),

	"lost woods default spring": Root(),
	"lost woods default summer": Root(),
	"lost woods default autumn": Root("start"),
	"lost woods default winter": Root(),

	"tarm ruins default spring": Root("start"),
	"tarm ruins default summer": Root(),
	"tarm ruins default autumn": Root(),
	"tarm ruins default winter": Root(),

	"western coast default spring": Root(),
	"western coast default summer": Root(),
	"western coast default autumn": Root(),
	"western coast default winter": Root("start"),

	"temple remains default spring": Root(),
	"temple remains default summer": Root(),
	"temple remains default autumn": Root(),
	"temple remains default winter": Root("start"),
}
