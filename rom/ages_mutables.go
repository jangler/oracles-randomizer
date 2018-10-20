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
}
