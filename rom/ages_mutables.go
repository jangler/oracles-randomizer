package rom

var agesFixedMutables = map[string]Mutable{
	// first and second time portals (near maku tree) are always active
	"first portal active": MutableString(Addr{0x10, 0x7d4e},
		"\x20\x0e", "\x20\x00"),
	"second portal active": MutableString(Addr{0x10, 0x7d57},
		"\x38\x05", "\x38\x00"),

	// allow access to nayru's house check from start
	"move impa": MutableByte(Addr{0x09, 0x6567}, 0xd0, 0x00),

	// prevent stairs disappearing in event where maku tree is attacked by
	// moblins, preventing softlock if player gets there with seed satchel and
	// no sword or something stupid
	"maku tree stairs": MutableString(Addr{0x15, 0x6bf3},
		"\x84\x05\x00", "\xc4\x15\xc3"),

	// set seed capacity by level to 20/20/50/cb instead of c9/20/50/99 so that
	// level zero (shooter only) can still carry 20 seeds.
	"seed capacity pointer": MutableByte(Addr{0x3f, 0x4608}, 0x10, 0x11),
	"seed capacity table": MutableString(Addr{0x3f, 0x4611},
		"\x20\x50\x99", "\x20\x20\x50"),

	// change harp interaction to allow sub ID
	"create harp with sub ID": MutableString(Addr{0x0b, 0x6825},
		"\xcd\xef\x3a\xc0\x36\x60\x2c\x36\x11",
		"\xc5\x01\x00\x11\xcd\xd4\x27\xc1\xc0"),

	// delete cutscene interaction in nayru's basement after it's done
	// initializing
	"delete harp cutscene": MutableString(Addr{0x0b, 0x684a},
		"\xc3\xe0\x23", "\xc3\xe0\x21"),

	// remove essence check for fairies' hide and seek game
	"fairies' essence check": MutableString(Addr{0x0a, 0x52b4},
		"\xcd\x48\x17", "\x37\x37\x37"),

	// make guy in front of d2 go away if you have bombs
	"d2 guy flag check": MutableString(Addr{0x09, 0x5242},
		"\x3e\x0b\xcd\xf3\x31\xc2", "\x3e\x03\xcd\x48\x17\xda"),

	// clear rubble from rolling ridge base present without d4 essence
	"clear rubble": MutableByte(Addr{0x04, 0x6a44}, 0xc8, 0x00),

	// cut off the end of deku forest soldier's text so it makes sense when
	// giving item
	"soldier text end": MutableByte(Addr{0x23, 0x6656}, 0x01, 0x00),
	// and position the "you may go now" text correctly on screen
	"soldier text position": MutableByte(Addr{0x23, 0x65d8}, 0x22, 0x00),
	// and remove the usual soldier event (taken to palace etc)
	"remove soldier event": MutableByte(Addr{0x12, 0x58f5}, 0xcd, 0xc9),

	// skip essence checks for the following events:
	"rafton essence check":    MutableByte(Addr{0x0a, 0x4d7a}, 0x20, 0x18),
	"dimitri essence check 1": MutableByte(Addr{0x09, 0x5816}, 0x13, 0x00),
	"dimitri essence check 2": MutableByte(Addr{0x0a, 0x4bb3}, 0xc8, 0x00),
	"open palace": MutableString(Addr{0x09, 0x51f8},
		"\x3e\x40", "\xaf\xc9"),

	// moosh should always appear in the graveyard
	"moosh essence checks": MutableStrings([]Addr{{0x0a, 0x5dd5},
		{0x0a, 0x5943}, {0x0a, 0x4b85}}, "\xcb\x4f", "\xf6\x01"),
	"moosh rope checks": MutableStrings([]Addr{{0x05, 0x78b8}, {0x0a, 0x4ba3}},
		"\xcd\x48\x17", "\x3f\x3f\x3f"),
	"moosh cheval checks": MutableStrings([]Addr{{0x0a, 0x5ddc},
		{0x0a, 0x594b}, {0x0a, 0x4b8c}}, "\xcb\x77", "\xf6\x01"),
	// allow exiting moosh/ghost cutscene screen without killing ghosts
	"transition from moosh cutscene": MutableString(Addr{0x0a, 0x595a},
		"\xea\x91\xcc", "\x00\x00\x00"),
	// vanilla bug: moosh is leaves forever if you enter the screen south of
	// cheval's grave if you dismount, transition, and come back
	"vanilla moosh disappear bug": MutableByte(Addr{0x05, 0x78bb}, 0x30, 0x38),
	// another: moosh appears on the screen south of cheval's grave after
	// visiting the cheval's grave screen, whether you've obtained him or not
	"vanilla moosh appear bug": MutableByte(Addr{0x12, 0x5c5d}, 0xf1, 0xff),

	// ricky shouldn't leave after talking to tingle
	"end tingle script": MutableString(Addr{0x0c, 0x7e2a},
		"\x91\x03\xd1\x02", "\xbe\x7d\xfe\xba"),
	// and check fake treasure ID 13 (slingshot) instead of island chart
	"tingle fake ID": MutableByte(Addr{0x0c, 0x7e00}, 0x54, 0x13),
	// dig up item on south shore regardless of ricky state
	"south shore ricky check 1": MutableByte(Addr{0x04, 0x6b77}, 0x0a, 0x00),
	"south shore ricky check 2": MutableByte(Addr{0x04, 0x6b7b}, 0x06, 0x00),
	"south shore ricky check 3": MutableByte(Addr{0x0a, 0x5e2f}, 0x12, 0x00),
	"south shore ricky check 4": MutableByte(Addr{0x0a, 0x5e33}, 0x0e, 0x00),
	// and check fake treasure ID 08 (magnet gloves) instead of ricky's gloves
	"south shore fake ID": MutableStrings([]Addr{{0x04, 0x6b7d},
		{0x0a, 0x5e35}}, "\x48", "\x08"),

	// don't refill seeds when getting item from tingle
	"tingle seed refill": MutableString(Addr{0x0c, 0x7e7d},
		"\xe0\x0c\x18", "\xc4\x80\x7e"),

	// remove storm event that washes link up on crescent island without raft,
	// and the event where tokays steal link's items
	"remove storm event": MutableByte(Addr{0x0b, 0x52e3}, 0xc2, 0xc3),
	"remove tokay event": MutableStrings([]Addr{{0x09, 0x5756}, {0x09, 0x5731},
		{0x0a, 0x4fe1}}, "\xc2", "\xc3"),
	"remove tokay items": MutableString(Addr{0x09, 0x57a5},
		"\xcb\x77", "\x3c\x3c"),
	"tokay trading hut": MutableStrings([]Addr{{0x0a, 0x623a}, {0x0a, 0x62d7}},
		"\xcd\xf3\x31", "\xb7\xb7\xb7"),
	// don't have an item in the chicken hut
	"tokay bomb hut": MutableString(Addr{0x12, 0x638f},
		"\xf2\x6b\x0a\x28", "\xf3\x57\x41\xff"),

	// sell 150 rupee item from lynna city shop from the start
	"shop flute flag check": MutableString(Addr{0x09, 0x4333},
		"\x28\x04", "\x00\x00"),
	// check for fake treasure ID 07 (rod) so that non-unique items can be sold
	"shop fake ID": MutableByte(Addr{0x09, 0x4328}, 0x0e, 0x07),

	// remove flute item from shooting gallery prizes
	"shooting gallery script": MutableString(Addr{0x15, 0x51d8},
		"\xdf\x0e", "\xdf\x02"),
	// prevent bridge-building foreman from setting flag 22 so that
	// animal/flute event doesn't happen in fairies' woods
	"bridge foreman script": MutableString(Addr{0x15, 0x75bf},
		"\xb6\x22", "\xb6\xa2"),

	// skip normal boomerang check in target carts, since EOB code handles it
	"skip target carts boomerang check": MutableString(Addr{0x15, 0x66ae},
		"\x20\x0b", "\x18\x0b"),
	// and remove "boomerang" from random prizes
	"target carts prize table": MutableString(Addr{0x15, 0x66e5},
		"\x04\x04\x04", "\x03\x03\x03"),

	// trade lava juice without mermaid key
	"trade lava juice without key": MutableString(Addr{0x15, 0x6879},
		"\x30\x07", "\x30\x00"),

	// stop d6 boss key chest from setting past boss key flag
	"stop d6 boss key chest": MutableString(Addr{0x10, 0x793c},
		"\xc3\x0e\x02", "\xc9\x00\x00"),

	// skip ralph cutscene entering palace
	"skip ralph at palace": MutableString(Addr{0x08, 0x6e61},
		"\xcb\x6f", "\xe6\x00"),
	// and get rid of the intangible guard standing outside
	"remove intangible guard": MutableByte(Addr{0x09, 0x5152}, 0xc2, 0xc3),

	// don't require talking to queen fairy before getting book of seals
	"skip library flag check": MutableString(Addr{0x15, 0x5da6},
		"\xb5\x20\xac\x5d", "\xc4\xac\x5d\x00"),

	// remove special interaction from caves in sea of storms so that the
	// chests can be normal chests
	"normalize sea of storms chests": MutableStrings(
		[]Addr{{0x12, 0x6417}, {0x12, 0x6421}}, "\xf1", "\xff"),

	// fix pickup text for harp tunes
	"tune of echoes text": MutableString(Addr{0x1e, 0x4c3e}, "\x49",
		"\x02\x06"+ // You got the
			"\x09\x01Tune\x04\xceE\x05\x0d\x04\x91"+ // Tune of Echoes!
			"Play\x04\x0f\x01"+ // Play it to
			"awaken \x04\xa8\x04\x5a"+ // awaken sleeping
			"\x09\x03Time Portals\x09\x00!\x00"), // Time Portals!
	"tune of currents text": MutableString(Addr{0x1d, 0x7e48}, "\x59",
		"\x02\x06"+ // You got the
			"\x09\x01Tune\x04\xce\x01"+ // Tune of
			"Currents\x05\x95Play\x01"+ // Currents! Play
			"it\x04\x57\x05\x5b\x03\x50"+ // it to move from
			"\x02\x81 \x02\x64\x01"+ // the past to the
			"\x03\x2e!\x00"), // the present!
	"tune of ages text": MutableString(Addr{0x1d, 0x7e8e}, "\x59",
		"\x02\x06"+ // You got the
			"\x09\x01Tune \x03\x31\x04\x91"+ // Tune of Ages!
			"Play\x04\x0f\x04\xdf"+ // Play it to move
			"freely \x02\x77\x01"+ // freely through
			"\x04\xdd!\x00"), // time!

	// make tokay who gives iron shield always give the same item, and in a
	// format compatible with lookupItemSpriteAddr.
	"give hidden tokay item": MutableString(Addr{0x15, 0x5b35},
		"\x06\x01\x0e\x01\xfa\xaf\xc6\xfe\x02",
		"\x01\x01\x01\x78\x41\x4f\x37\x00\x00"),

	// game has zora scale palette in item gfx wrong for some reason
	"fix zora scale palette": MutableByte(Addr{0x3f, 0x67d0}, 0x13, 0x43),

	// put a bush on the other side of the syrup's shop screen so that long
	// hook isn't a softlock
	"syrup screen fix 1": MutableString(Addr{0x23, 0x7ea0},
		"\x01\x27", "\x27\xc8"),
	"syrup screen fix 2": MutableByte(Addr{0x23, 0x7ead}, 0x27, 0x22),
}

var agesVarMutables = map[string]Mutable{
	// seed tree types
	"symmetry city tree sub ID": MutableByte(Addr{0x12, 0x59a1}, 0x35, 0x35),
	"south lynna present tree sub ID": MutableByte(Addr{0x12, 0x5ca4},
		0x06, 0x06),
	"crescent island tree sub ID": MutableByte(Addr{0x12, 0x59b8}, 0x17, 0x17),
	"zora village present tree sub ID": MutableByte(Addr{0x12, 0x59bf},
		0x38, 0x38),
	"rolling ridge west tree sub ID": MutableByte(Addr{0x12, 0x5e4d},
		0x29, 0x29),
	"ambi's palace tree sub ID": MutableByte(Addr{0x12, 0x5e5b}, 0x1a, 0x1a),
	"rolling ridge east tree sub ID": MutableByte(Addr{0x12, 0x5f46},
		0x4b, 0x4b),
	"south lynna past tree sub ID": MutableByte(Addr{0x12, 0x5e62}, 0x0c, 0x0c),
	"deku forest tree sub ID":      MutableByte(Addr{0x12, 0x6101}, 0x4d, 0x4d),
	"zora village past tree sub ID": MutableByte(Addr{0x12, 0x5e6f},
		0x3e, 0x3e),

	// first satchel should give the seeds on the south lynna tree.
	"satchel initial seeds": MutableByte(Addr{0x3f, 0x453b}, 0x20, 0x20),

	// set default satchel and shooter selection based on south lynna tree.
	// overwrites unimportant bytes in file initialization.
	"satchel initial selection": MutableWord(Addr{0x07, 0x418e}, 0x0700, 0xc400),
	"shooter initial selection": MutableWord(Addr{0x07, 0x4190}, 0x0e00, 0xc500),

	// map pop-up icons for seed trees
	"crescent island tree map icon": MutableByte(
		Addr{0x02, 0x6d05}, 0x16, 0x16),
	"symmetry city tree map icon": MutableByte(
		Addr{0x02, 0x6d08}, 0x18, 0x18),
	"south lynna tree map icon": MutableStrings(
		[]Addr{{0x02, 0x6d0b}, {0x02, 0x6d29}}, "\x15", "\x15"),
	"zora village tree map icon": MutableStrings(
		[]Addr{{0x02, 0x6d0e}, {0x02, 0x6d2f}}, "\x18", "\x18"),
	"rolling ridge west tree map icon": MutableByte(
		Addr{0x02, 0x6d20}, 0x17, 0x17),
	"ambi's palace tree map icon": MutableByte(
		Addr{0x02, 0x6d23}, 0x16, 0x16),
	"rolling ridge east tree map icon": MutableByte(
		Addr{0x02, 0x6d26}, 0x19, 0x19),
	"deku forest tree map icon": MutableByte(
		Addr{0x02, 0x6d2c}, 0x19, 0x19),

	// 33 for ricky, 23 for dimitri, 13 for moosh
	"flute palette": MutableByte(Addr{0x3f, 0x6746}, 0x03, 0x03),
	// 0b for ricky, 0c for dimitri, 0d for moosh
	"animal region": MutableByte(Addr{0xa, 0x5ac6}, 0x0d, 0x0d),
}
