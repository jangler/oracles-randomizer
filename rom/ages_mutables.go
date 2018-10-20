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

var agesVarMutables = map[string]Mutable{}
