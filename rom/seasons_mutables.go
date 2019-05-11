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

	// banks 1c-1f (text)

	// all this text overwrites the text from the initial rosa encounter, which
	// runs from 1f:4533 to 1f:45c1 inclusive. the last entry is displayed at
	// the end of any warning message.
	"cliff warning text": MutableString(Addr{0x1f, 0x4533}, "\x0c\x21",
		"\x0c\x00\x02\x3b\x67\x6f\x20\x05\x73\x01"+ // If you go down
			"\x74\x68\x65\x72\x65\x2c\x04\x2d\x20\x77\x6f\x6e\x27\x74\x01"+ // there, you won't
			"\x62\x65\x20\x02\xa4\x05\x0f\x01"+ // be able to get
			"\x04\x9f\x20\x75\x70\x03\xa4"+ // back up.
			"\x07\x03"), // jump to end text
	"hss skip warning addr": MutableWord(Addr{0x1c, 0x6b52}, 0x1192, 0x0292),
	"hss skip warning text": MutableString(Addr{0x1f, 0x4584}, "\x20\x05",
		"\x0c\x00\x02\x3b\x73\x6b\x69\x70\x01"+ // If you skip
			"\x6b\x65\x79\x73\x2c\x04\xaa\x03\x2c\x01"+ // keys, use them
			"\x03\x70\x6c\x79\x03\xa4"+ // carefully.
			"\x07\x03"), // jump to end text
	"end warning addr": MutableWord(Addr{0x1c, 0x6b54}, 0x2592, 0x1d92),
	"end warning text": MutableString(Addr{0x1f, 0x459f}, "\x01\x05",
		"\x0c\x00\x43\x6f\x6e\x74\x69\x6e\x75\x65\x20\x61\x74\x01"+ // Continue at
			"\x03\x0b\x6f\x77\x6e\x20\x72\x69\x73\x6b\x21\x00"), // your own risk!
}

// key = area name (as in asm/vars.yaml), id = season index (spring -> winter).
func SetSeason(key string, id byte) {
	codeMutables[key].(*MutableRange).New[0] = id
}
