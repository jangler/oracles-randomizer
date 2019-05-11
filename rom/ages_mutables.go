package rom

var agesFixedMutables = map[string]Mutable{
	// first and second time portals (near maku tree) are always active
	"first portal active": MutableString(Addr{0x10, 0x7d4e},
		"\x20\x0e", "\x20\x00"),
	"second portal active": MutableString(Addr{0x10, 0x7d57},
		"\x38\x05", "\x38\x00"),

	// prevent stairs disappearing in event where maku tree is attacked by
	// moblins, preventing softlock if player gets there with seed satchel and
	// no sword or something stupid
	"maku tree stairs": MutableString(Addr{0x15, 0x6bf3},
		"\x84\x05\x00", "\xc4\x15\xc3"),

	// never spawn hide and seek event in fairies' woods. apparently you're
	// frozen if you enter on an animal?
	"don't spawn fairies": MutableByte(Addr{0x0a, 0x52bf}, 0xc2, 0xc3),

	// make guy in front of d2 go away if you have bombs
	"d2 guy flag check": MutableString(Addr{0x09, 0x5242},
		"\x3e\x0b\xcd\xf3\x31\xc2", "\x3e\x03\xcd\x48\x17\xda"),
	// and center him on a tile so you can't get stuck in a currents loop
	"d2 guy position": MutableByte(Addr{0x12, 0x611c}, 0x4e, 0x48),

	// cut off the end of deku forest soldier's text so it makes sense when
	// giving item
	"soldier text end": MutableByte(Addr{0x23, 0x6656}, 0x01, 0x00),
	// and position the "you may go now" text correctly on screen
	"soldier text position": MutableByte(Addr{0x23, 0x65d8}, 0x22, 0x00),
	// and remove the usual soldier event (taken to palace etc)
	"remove soldier event": MutableByte(Addr{0x12, 0x58f5}, 0xcd, 0xc9),

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

	// prevent bridge-building foreman from setting flag 22 so that
	// animal/flute event doesn't happen in fairies' woods
	"bridge foreman script": MutableString(Addr{0x15, 0x75bf},
		"\xb6\x22", "\xb6\xa2"),

	// stop d6 boss key chest from setting past boss key flag
	"stop d6 boss key chest": MutableString(Addr{0x10, 0x793c},
		"\xc3\x0e\x02", "\xc9\x00\x00"),

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

	// buy tokay trader's shield if you have scent seeds but not satchel
	"tokay trader satchel check": MutableString(Addr{0x0a, 0x629c},
		"\x30\x16", "\x30\x00"),

	// add railing to ricky nuun screen and move worker off the "roof"
	"ricky nuun railing": MutableString(Addr{0x23, 0x718e},
		"\x69\x07\x07\x6a", "\x72\x50\x50\x73"),
	"move nuun worker": MutableString(Addr{0x12, 0x5a9e},
		"\x28\x50", "\x68\x48"),

	// text for special crescent island present portal
	"portal sign text": MutableString(Addr{0x23, 0x583f}, "\x0c\x20\x02\x18",
		"\x0c\x00C\x04\x23s only.\x01"+ // Currents only.
			" -\x04\x56Management\x00"), // -The Management

	// vanilla bug: compass doesn't show D6 boss key chest.
	"fix d6 compass": MutableByte(Addr{0x01, 0x4eea}, 0x14, 0x34),
}
