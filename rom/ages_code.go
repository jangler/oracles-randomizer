package rom

func newAgesRomBanks() *romBanks {
	r := romBanks{
		endOfBank: make([]uint16, 0x40),
	}

	r.endOfBank[0x00] = 0x3ef8
	r.endOfBank[0x02] = 0x7e95
	r.endOfBank[0x04] = 0x7ee2

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
}
