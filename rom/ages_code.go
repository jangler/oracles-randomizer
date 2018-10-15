package rom

func newAgesRomBanks() *romBanks {
	r := romBanks{
		endOfBank: make([]uint16, 0x40),
	}

	r.endOfBank[0x00] = 0x3ef8
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
