package rom

func newAgesRomBanks() *romBanks {
	r := romBanks{
		endOfBank: make([]uint16, 0x40),
	}

	r.endOfBank[0x00] = 0x3ef8
	r.endOfBank[0x02] = 0x7e93
	r.endOfBank[0x03] = 0x7ebd
	r.endOfBank[0x04] = 0x7edb

	return &r
}

func initAgesEOB() {
	r := newAgesRomBanks()

	// bank 00

	// don't play any music if the -nomusic flag is given. because of this,
	// this *must* be the first function at the end of bank zero (for now).
	r.appendToBank(0x00, "no music func",
		"\x67\xfe\x40\x30\x03\x3e\x08\xc9\xf0\xb7\xc9")
	r.replace(0x00, 0x0c9a, "no music call",
		"\x67\xf0\xb7", "\x67\xf0\xb7") // modified only by SetNoMusic()

	// bank 02

	// warp to ember tree if holding start when closing the map screen.
	treeWarp := r.appendToBank(0x02, "tree warp",
		"\xfa\x81\xc4\xe6\x08\x28\x1b"+ // close as normal if start not held
			"\xfa\x2d\xcc\xfe\x02\x38\x06"+ // check if indoors
			"\x3e\x5a\xcd\x98\x0c\xc9"+ // play error sound and ret
			"\x21\xb7\xcb\x36\x02\xb7\x28\x02\x36\x03"+ // set tree based on age
			"\xaf\xcd\xac\x5f\xc3\xba\x4f") // close + warp
	r.replaceMultiple([]Addr{{0x02, 0x6133}, {0x02, 0x618b}}, "tree warp jump",
		"\xc2\xba\x4f", "\xc4"+treeWarp)

	// warp to room under cursor if wearing developer ring.
	devWarp := r.appendToBank(0x02, "dev ring warp func",
		"\xfa\xcb\xc6\xfe\x40\x20\x12\xfa\x2d\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x47\xcc\xfa\xb6\xcb\xea\x48\xcc\x3e\x03\xcd\xad\x0c\xc9")
	r.replace(0x02, 0x5fcc, "dev ring warp call", "\xad\x0c", devWarp)

	// bank 03

	// allow skipping the capcom screen after one second by pressing start
	skipCapcom := r.appendToBank(0x03, "skip capcom func",
		"\xe5\xfa\xb3\xcb\xfe\x94\x30\x03\xcd\x86\x08\xe1\xcd\x37\x02\xc9")
	r.replace(0x03, 0x4d6c, "skip capcom call", "\x37\x02", skipCapcom)

	// bank 04

	// look up tiles in custom replacement table after loading a room. the
	// format is (group, room, YX, tile ID), with ff ending the table. if bit 0
	// of the room is set, no replacements are made.
	tileReplaceTable := r.appendToBank(0x04, "tile replace table",
		"\x03\x0f\x66\xf9"+ // water in d6 past entrance
			"\x04\x1b\x03\x78"+ // key door in D1
			"\xff")
	tileReplaceFunc := r.appendToBank(0x04, "tile replace body",
		"\xcd\x7d\x19\xe6\x01\x20\x28"+
			"\xc5\x21"+tileReplaceTable+"\xfa\x2d\xcc\x47\xfa\x30\xcc\x4f"+
			"\x2a\xfe\xff\x28\x16\xb8\x20\x0e\x2a\xb9\x20\x0b"+
			"\xd5\x16\xcf\x2a\x5f\x2a\x12\xd1\x18\xea"+
			"\x23\x23\x23\x18\xe5\xc1\xcd\xef\x5f\xc9")
	r.replace(0x00, 0x38c0, "tile replace call",
		"\xcd\xef\x5f", "\xcd"+tileReplaceFunc)
}
