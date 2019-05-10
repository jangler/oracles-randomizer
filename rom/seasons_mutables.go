package rom

var seasonsFixedMutables = map[string]Mutable{
	// make link actionable as soon as he drops into the world.
	"link immediately actionable": MutableString(Addr{0x05, 0x4d98},
		"\x3e\x08\xcd\x16", "\xcd\x16\x2a\xc9"),

	// this all has to do with animals and flutes:
	// this edits ricky's script so that he never gives his flute.
	"ricky skip flute script":  MutableByte(Addr{0x0b, 0x6b7a}, 0x0b, 0x7f),
	"don't give ricky's flute": MutableByte(Addr{0x09, 0x6e6c}, 0xc0, 0xc9),
	// this prevents subrosian dancing from giving dimitri's flute.
	"don't give dimitri's flute": MutableByte(Addr{0x09, 0x5e37}, 0xe6, 0xf6),
	// this prevents holodrum plain from changing the animal region.
	"don't change animal region": MutableWord(Addr{0x09, 0x6f79},
		0x3804, 0x1808),
	// this keeps ricky in his pen based on flute, not animal region.
	"keep ricky in pen": MutableString(Addr{0x09, 0x4e77},
		"\x10\xc6\xfe\x0b", "\xaf\xc6\xfe\x01"),
	// and this does the same for saying goodbye once reaching spool swamp.
	"ricky say goodbye": MutableString(Addr{0x09, 0x6ccd},
		"\x10\xc6\xfe\x0b", "\xaf\xc6\xfe\x01"),
	// spawn dimitri and kids in sunken city based on flute, not animal region.
	"spawn dimitri in sunken city": MutableStrings(
		[]Addr{{0x09, 0x4e4c}, {0x09, 0x6f08}, {0x09, 0x737e}},
		"\x10\xc6\xfe\x0c", "\xaf\xc6\xfe\x02"),

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

	// skip shield check for forging hard ore
	"skip iron shield check": MutableByte(Addr{0x0b, 0x75c7}, 0x01, 0x02),
	// and skip the check for what level shield you currently have
	"skip iron shield level check": MutableString(Addr{0x15, 0x62ac},
		"\x38\x01", "\x18\x05"),

	// check fake treasure ID 0a for maku tree item. this only matters if you
	// leave the screen without picking up the item.
	"maku tree check fake id": MutableByte(Addr{0x09, 0x7dfd}, 0x42, 0x0a),
	// check fake treasure ID 10 for market item 5.
	"market check fake id": MutableByte(Addr{0x09, 0x7755}, 0x53, 0x10),
	// check fake treasure ID 11 for master diver.
	"diver check fake id": MutableByte(Addr{0x0b, 0x72f1}, 0x2e, 0x11),

	// banks 08-0a (most interaction-specific non-script behavior?)

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e72}, 0xcb, 0xf6),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e40}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e56}, 0x2e, 0x04),
	"master essence check 1":    MutableByte(Addr{0x0a, 0x4bf5}, 0x02, 0x00),
	"master essence check 2":    MutableByte(Addr{0x0a, 0x4bea}, 0x40, 0x02),
	"round jewel essence check": MutableByte(Addr{0x0a, 0x4f8b}, 0x05, 0x00),

	// prevent leaving sunken city with dimitri unless you have his flute, in
	// order to prevent a variety of softlocks.
	"block dimitri exit": MutableString(Addr{0x09, 0x6f34},
		"\xfa\x10\xc6\xfe\x0c", "\xfa\xaf\xc6\xfe\x02"),

	// moosh won't spawn in the mountains if you have the wrong number of
	// essences. bit 6 seems related to this, and needs to be zero too?
	"skip moosh essence check 1": MutableByte(Addr{0x0f, 0x7429}, 0x03, 0x00),
	"skip moosh essence check 2": MutableByte(Addr{0x09, 0x4e36}, 0xca, 0xc3),

	// sell member's card in subrosian market before completing d3
	"member's card essence check": MutableWord(Addr{0x09, 0x7750},
		0xcb57, 0xf601),

	// count number of essences, not highest numbered essence.
	"maku seed check 1": MutableByte(Addr{0x09, 0x7da4}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7da6}, 0x30, 0x18),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7bb9}, 0x05, 0x3f),

	// bank 0b (scripts)

	// don't set a ricky flag when buying the "flute".
	"shop no set ricky flag": MutableByte(Addr{0x0b, 0x4826}, 0x20, 0x00),

	// don't require rod to get items from season spirits.
	"season spirit rod check": MutableByte(Addr{0x0b, 0x4eb2}, 0x07, 0x02),

	// getting the L-2 (or L-3) sword in the lost woods normally gives a second
	// "spin slash" item. remove this from the script.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0xde, 0xc1),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0xde, 0xc1),

	// end maku seed script as soon as link gets the seed.
	"abbreviate maku seed cutscene": MutableString(Addr{0x0b, 0x71ec},
		"\xe1\x23\x61\x01", "\xb6\x19\xbe\x00"),
	// end northen peak barrier cutscene as soon as the barrier is broken.
	"abbreviate barrier cutscene": MutableString(Addr{0x0b, 0x79f1},
		"\x88\x18\x50\xf8", "\xb6\x1d\xbe\x00"),

	// bank 0d

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b5}, 0xb8, 0xbf),

	// bank 11 (interactions)

	// remove the moosh and dimitri events in spool swamp.
	"prevent moosh cutscene":   MutableByte(Addr{0x11, 0x6572}, 0xf1, 0xff),
	"prevent dimitri cutscene": MutableByte(Addr{0x11, 0x68d4}, 0xf1, 0xff),

	// bank 14 (scripts loaded into c3xx block, also some OAM pointers?)

	// skip the great furnace dance. for some reason command c4 (jumpalways)
	// doesn't work here, so a jump based on c6xx is used instead.
	"skip furnace dance": MutableString(Addr{0x14, 0x4b15},
		"\xe4\x31\xd7\x7d\x80", "\xb3\x92\xff\x3f\xc3"),

	// allow screen transitions away from linked great moblin event, since
	// seeds could be uncompleteable without a way to kill the moblins.
	"allow transitions from linked moblins": MutableString(Addr{0x14, 0x5247},
		"\x91\xab\xcc\x01", "\x91\xab\xcc\x00"),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x53d7}, 0x5959, 0x1957),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x55a7}, 0xf36b, 0x4f68),

	// bank 15 (script functions)

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move d8 magnet ball": MutableByte(Addr{0x15, 0x4f62}, 0x48, 0x38),

	// change destination of initial transition in pirate cutscene.
	"pirate warp": MutableString(Addr{0x15, 0x5a1c},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),

	// zero normal rod text.
	"no rod text": MutableString(Addr{0x15, 0x70be},
		"\xcd\x4b\x18", "\x00\x00\x00"),

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

	// bank 3f

	// give items that don't normally appear as treasure interactions entries
	// in the treasure graphics table.
	"member's card gfx": MutableString(Addr{0x3f, 0x67b4},
		"\x00\x00\x00", "\x5d\x0c\x13"),
	"treasure map gfx": MutableString(Addr{0x3f, 0x67b7},
		"\x00\x00\x00", "\x65\x14\x33"),
	"fool's ore gfx": MutableString(Addr{0x3f, 0x67ba},
		"\x00\x00\x00", "\x60\x14\x00"),
	"rare peach stone gfx": MutableString(Addr{0x3f, 0x67c6},
		"\x00\x00\x00", "\x5d\x10\x26"),
	"ribbon gfx": MutableString(Addr{0x3f, 0x67c9},
		"\x00\x00\x00", "\x65\x0c\x23"),
}

var seasonsVarMutables = map[string]Mutable{
	// set initial season correctly in the init variables. this replaces
	// null-terminating whoever's son's name, which *should* be zeroed anyway.
	"initial season": MutableWord(Addr{0x07, 0x4188}, 0x0e00, 0x2d00),

	// map pop-up icons for seed trees
	"tarm ruins seed tree map icon":      MutableByte(Addr{0x02, 0x6c51}, 0x18, 0x18),
	"sunken city seed tree map icon":     MutableByte(Addr{0x02, 0x6c54}, 0x18, 0x18),
	"north horon seed tree map icon":     MutableByte(Addr{0x02, 0x6c57}, 0x16, 0x16),
	"spool swamp seed tree map icon":     MutableByte(Addr{0x02, 0x6c5a}, 0x17, 0x17),
	"woods of winter seed tree map icon": MutableByte(Addr{0x02, 0x6c5d}, 0x19, 0x19),
	"horon village seed tree map icon":   MutableByte(Addr{0x02, 0x6c60}, 0x15, 0x15),

	// locations of sparkles on treasure map
	"round jewel coords":    MutableByte(Addr{0x02, 0x6663}, 0xb5, 0xb5),
	"pyramid jewel coords":  MutableByte(Addr{0x02, 0x6664}, 0x1d, 0x1d),
	"square jewel coords":   MutableByte(Addr{0x02, 0x6665}, 0xc2, 0xc2),
	"x-shaped jewel coords": MutableByte(Addr{0x02, 0x6666}, 0xf4, 0xf4),

	// the satchel should contain the type of seeds that grow on the horon
	// village tree.
	"satchel initial seeds": MutableByte(Addr{0x3f, 0x453b}, 0x20, 0x20),

	// give the player seeds when they get the slingshot, and don't take the
	// player's: fool's ore when they get feather, star ore when they get
	// ribbon, or red and blue ore when they get hard ore (just zero the whole
	// "lose items" table). one byte of this is changed in setSeedData() to
	// change what type of seeds the slingshot gives.
	"edit gain/lose items tables": MutableString(Addr{0x3f, 0x4543},
		"\x00\x46\x45\x00\x52\x50\x51",
		"\x13\x20\x20\x00\x00\x00\x00"),
	"edit lose items table pointer": MutableByte(Addr{0x3f, 0x44cf},
		0x44, 0x47),

	// the correct type of seed needs to be selected by default, otherwise the
	// player may be unable to use seeds when they only have one type. there
	// could also be serious problems with the submenu when they *do* obtain a
	// second type if the selection isn't either of them.
	//
	// this works by overwriting a couple of unimportant bytes in file
	// initialization.
	"satchel initial selection":   MutableWord(Addr{0x07, 0x418e}, 0xa210, 0xbe00),
	"slingshot initial selection": MutableWord(Addr{0x07, 0x419a}, 0x2e02, 0xbf00),

	// allow seed collection if you have a slingshot, by checking for the given
	// initial seed type
	"carry seeds in slingshot": MutableByte(Addr{0x10, 0x4b19}, 0x19, 0x20),

	// 33 for ricky, 23 for dimitri, 13 for moosh
	"flute palette": MutableByte(Addr{0x3f, 0x6747}, 0x03, 0x03),
	// 0b for ricky, 0c for dimitri, 0d for moosh
	"animal region": MutableByte(Addr{0x0a, 0x7fff}, 0x0a, 0x0b),

	// for the item dropped in the room *above* the trampoline
	"above d7 zol button": &MutableSlot{
		idAddrs:    []Addr{{0x15, 0x55d8}},
		subIDAddrs: []Addr{{0x15, 0x55db}},
	},
}

var Seasons = map[string]*MutableRange{
	// randomize default seasons (before routing). sunken city also applies to
	// mt. cucco; eastern suburbs applies to the vertical part of moblin road
	// but not the horizontal part. note that "tarm ruins" here refers only to
	// the part beyond the lost woods.
	//
	// horon village is random, natzu and desert can only be summer, and goron
	// mountain can only be winter. not sure about northern peak but it doesn't
	// matter.
	"north horon season":     MutableByte(Addr{0x01, 0x7e60}, 0x03, 0x03),
	"eastern suburbs season": MutableByte(Addr{0x01, 0x7e61}, 0x02, 0x02),
	"woods of winter season": MutableByte(Addr{0x01, 0x7e62}, 0x01, 0x01),
	"spool swamp season":     MutableByte(Addr{0x01, 0x7e63}, 0x02, 0x02),
	"holodrum plain season":  MutableByte(Addr{0x01, 0x7e64}, 0x00, 0x00),
	"sunken city season":     MutableByte(Addr{0x01, 0x7e65}, 0x01, 0x01),
	"lost woods season":      MutableByte(Addr{0x01, 0x7e67}, 0x02, 0x02),
	"tarm ruins season":      MutableByte(Addr{0x01, 0x7e68}, 0x00, 0x00),
	"western coast season":   MutableByte(Addr{0x01, 0x7e6b}, 0x03, 0x03),
	"temple remains season":  MutableByte(Addr{0x01, 0x7e6c}, 0x03, 0x03),
}
