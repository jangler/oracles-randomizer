package rom

var seasonsFixedMutables = map[string]Mutable{
	// make link actionable as soon as he drops into the world.
	"link immediately actionable": MutableString(Addr{0x05, 0x4d98},
		"\x3e\x08\xcd\x16", "\xcd\x16\x2a\xc9"),

	// move sleeping talon and his mushroom so they don't block the chest
	"move talon":    MutableWord(Addr{0x11, 0x6d2b}, 0x6858, 0x88a8),
	"move mushroom": MutableWord(Addr{0x0b, 0x6080}, 0x6848, 0x78a8),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore": MutableString(Addr{0x14, 0x4881},
		"\xe0\xeb\x58", "\xf0\xf0\xf0"),
	// but always give up feather if the player doesn't have it
	"give stolen feather": MutableString(Addr{0x15, 0x5dcf},
		"\xcd\x56\x19\xcb\x6e\x20", "\x3e\x17\xcd\x17\x17\x38"),
	// and make the feather appear without needing to be dug up
	"stolen feather appears": MutableByte(Addr{0x15, 0x5335}, 0x5a, 0x1a),
	// AND allow transition away from the screen if you have feather (not once
	// the hole is dug)
	"leave H&S screen": MutableString(Addr{0x09, 0x65a0},
		"\xcd\x32\x14\x1e\x49\x1a\xbe\xc8",
		"\x3e\x17\xcd\x17\x17\x00\x00\xd0"),

	// move the trigger for the bridge from holodrum plain to natzu to the
	// top-left corner of the screen, where it can't be hit, and replace the
	// lever tile as well. this prevents the bridge from blocking the waterway.
	"remove bridge trigger": MutableWord(Addr{0x11, 0x6737},
		0x6868, 0x0000),
	"remove prairie lever":   MutableByte(Addr{0x21, 0x6267}, 0xb1, 0x04),
	"remove wasteland lever": MutableByte(Addr{0x23, 0x5cb7}, 0xb1, 0x04),

	// bank 0d

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b5}, 0xb8, 0xbf),

	// bank 15 (script functions)

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move d8 magnet ball": MutableByte(Addr{0x15, 0x4f62}, 0x48, 0x38),
}

// key = area name (as in asm/vars.yaml), id = season index (spring -> winter).
func SetSeason(key string, id byte) {
	codeMutables[key].(*MutableRange).New[0] = id
}
