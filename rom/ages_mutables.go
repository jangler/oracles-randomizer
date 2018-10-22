package rom

var agesFixedMutables = map[string]Mutable{
	// first and second time portals (near maku tree) are always active
	"first portal active": MutableString(Addr{0x10, 0x7d4e},
		"\x20\x0e", "\x20\x00"),
	"second portal active": MutableString(Addr{0x10, 0x7d57},
		"\x38\x05", "\x38\x00"),

	// allow access to nayru's house check from start
	"move impa": MutableByte(Addr{0x09, 0x6567}, 0xd0, 0x00),

	// change harp interaction to allow sub ID
	"create harp with sub ID": MutableString(Addr{0x0b, 0x6825},
		"\xcd\xef\x3a\xc0\x36\x60\x2c\x36\x11",
		"\xc5\x01\x00\x11\xcd\xd4\x27\xc1\xc0"),

	// delete cutscene interaction in nayru's basement after it's done
	// initializing.
	"delete harp cutscene": MutableString(Addr{0x0b, 0x684a},
		"\xc3\xe0\x23", "\xc3\xe0\x21"),
	// dig up item on south shore regardless of ricky state
	"south shore ricky check 1": MutableByte(Addr{0x04, 0x6b77}, 0x0a, 0x00),
	"south shore ricky check 2": MutableByte(Addr{0x04, 0x6b7b}, 0x06, 0x00),
	"south shore ricky check 3": MutableByte(Addr{0x0a, 0x5e2f}, 0x12, 0x00),
	"south shore ricky check 4": MutableByte(Addr{0x0a, 0x5e33}, 0x0e, 0x00),
	// and check fake treasure ID 08 (magnet gloves) instead of ricky's gloves
	"south shore fake ID": MutableStrings([]Addr{{0x04, 0x6b7d},
		{0x0a, 0x5e35}}, "\x48", "\x08"),

	// sell flute from lynna city shop from the start
	"shop flute flag check": MutableString(Addr{0x09, 0x4333},
		"\x28\x04", "\x00\x00"),
	// check for fake treasure ID 07 (rod) so that non-unique items can be sold
	"shop fake ID": MutableByte(Addr{0x09, 0x4328}, 0x0e, 0x07),

	// skip normal boomerang check in target carts, since EOB code handles it
	"skip target carts boomerang check": MutableString(Addr{0x15, 0x66ae},
		"\x20\x0b", "\x18\x0b"),
	// and remove "boomerang" from random prizes
	"target carts prize table": MutableString(Addr{0x15, 0x66e5},
		"\x04\x04\x04", "\x03\x03\x03"),

	// skip ralph cutscene entering palace
	"skip ralph at palace": MutableString(Addr{0x08, 0x6e61},
		"\xcb\x6f", "\xe6\x00"),

	// remove special interaction from caves in sea of storms so that the
	// chests can be normal chests
	"normalize sea of storms chests": MutableStrings(
		[]Addr{{0x12, 0x6417}, {0x12, 0x6421}}, "\xf1", "\xff"),
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

	// item graphics
	"cheval's test gfx": MutableString(Addr{0x3f, 0x6a56},
		"\x79\x04\x52", "\x79\x04\x52"),
	"cheval's invention gfx": MutableString(Addr{0x3f, 0x6a53},
		"\x81\x10\x32", "\x81\x10\x32"),
	"tokay hut gfx": MutableString(Addr{0x3f, 0x6a50},
		"\x78\x10\x41", "\x78\x10\x41"),
	"wild tokay game gfx": MutableString(Addr{0x3f, 0x6795},
		"\x83\x00\x03", "\x83\x00\x03"),
	"shop, 150 rupees gfx": MutableString(Addr{0x3f, 0x69c6},
		"\x7c\x16\x03", "\x7c\x16\x03"),
	"library present gfx": MutableString(Addr{0x3f, 0x6894},
		"\x7a\x16\x04", "\x7a\x16\x04"),
	"library past gfx": MutableString(Addr{0x3f, 0x6891},
		"\x82\x12\x32", "\x82\x12\x32"),
}
