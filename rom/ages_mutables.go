package rom

var agesFixedMutables = map[string]Mutable{
	// first and second time portals (near maku tree) are always active
	"first portal active": MutableString(Addr{0x10, 0x7d4e},
		"\x20\x0e", "\x20\x00"),
	"second portal active": MutableString(Addr{0x10, 0x7d57},
		"\x38\x05", "\x38\x00"),
}

var agesVarMutables = map[string]Mutable{}
