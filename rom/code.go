package rom

// this file is for fixed mutables that go at the end of banks. each should be
// a self-contained unit (i.e. don't jr to anywhere outside the byte string) so
// that they can be appended automatically with respect to their size.

// return e.g. "\x2d\x79" for 0x792d
func addrString(addr uint16) string {
	return string([]byte{byte(addr), byte(addr >> 8)})
}

// adds code at the end of the bank, returning the new end of the bank.
func appendCode(name string, bank byte, endOfBank uint16, code string) uint16 {
	constMutables[name] = MutableString(Addr{bank, endOfBank},
		string([]byte{bank}), code)
	return endOfBank + uint16(len(code))
}

func initCode() {
	endOfBank15 := uint16(0x792d)

	// set hl = address of treasure data + 1 for item with ID a, sub ID c.
	constMutables["treasure data func"] = MutableString(Addr{0x00, 0x3ed3},
		"\x00", "\xf5\xc5\xd5\x47\x1e\x15\x21"+addrString(endOfBank15)+
			"\xcd\x8a\x00\xd1\xc1\xf1\xc9")
	endOfBank15 = appendCode("treasure data body", 0x15, endOfBank15,
		"\x78\xc5\x21\x29\x51\xcd\xc3\x01\x09"+ // add ID offset
			"\xcb\x7e\x28\x09\x23\x2a\x66\x6f"+ // load as address if bit 7 set
			"\xc1\x79\xc5\x18\xef"+ // use sub ID as second offset
			"\x23\x06\x03\xd5\x11\xfd\xcd\xcd\x62\x04"+ // copy data
			"\x21\xfd\xcd\xd1\xc1\xc9") // set hl and ret

	// ORs the default season in the given area (low byte b in bank 1) with the
	// seasons the rod has (c), then ANDs and compares the results with d.
	warningHelperAddr := addrString(endOfBank15)
	endOfBank15 = appendCode("warning helper", 0x15, endOfBank15,
		"\x1e\x01\x21\x89\x7e\xcd\x8a\x00"+ // get default season
			"\x78\xb7\x3e\x01\x28\x05\xcb\x27\x05\x20\xfb"+ // match rod format
			"\xb1\xa2\xba\xc9") // OR with c, AND with d, compare with d, ret
	// overwrite unused maku gate interaction with warning interaction
	constMutables["warning script pointer"] = MutableWord(Addr{0x08, 0x5663},
		0x874e, 0x6d7f)
	constMutables["warning script"] = MutableString(Addr{0x0b, 0x7f6d}, "\x0b",
		"\xd0\xe0"+addrString(endOfBank15)+
			"\xa0\xbd\xd7\x3c"+ // wait for collision and animation
			"\x87\xe0\xcf\x7e\x7f\x83\x7f\x88\x7f"+ // jump based on cfe0 bits
			"\x98\x26\x00\xbe\x00"+ // show cliff warning text
			"\x98\x26\x01\xbe\x00"+ // show bush warning text
			"\x98\x26\x02\xbe\x00") // show hss skip warning text
	// this communicates with the warning script by setting bit zero of $cfc0
	// if the warning needs to be displayed (based on room, season, etc), and
	// also displays the exclamation mark if so.
	endOfBank15 = appendCode("warning func", 0x15, endOfBank15,
		"\xc5\xd5\xcd"+addrString(endOfBank15+8)+"\xd1\xc1\xc9"+ // wrap in push/pops
			"\xfa\x4e\xcc\x47\xfa\xb0\xc6\x4f\xfa\x4c\xcc"+ // load room, season, rod
			"\xfe\x7c\x28\x12\xfe\x6e\x28\x18\xfe\x3d\x28\x22"+ // jump by room
			"\xfe\x5c\x28\x28\xfe\x78\x28\x32\x18\x43"+ // (cont.)
			"\x06\x61\x16\x01\xcd"+warningHelperAddr+"\xc8\x18\x35"+ // flower cliff
			"\x78\xfe\x03\xc8\x06\x61\x16\x09\xcd"+warningHelperAddr+"\xc8\x18\x27"+ // diving spot
			"\x06\x65\x16\x02\xcd"+warningHelperAddr+"\xc8\x18\x1d"+ // waterfall cliff
			"\xfa\x10\xc6\xfe\x0c\xc0\x3e\x17\xcd\x17\x17\xd8\x18\x0f"+ // keep
			"\xcd\x56\x19\xcb\x76\xc0\xcb\xf6\x3e\x02\xea\xe0\xcf\x18\x04"+ // hss skip room
			"\xaf\xea\xe0\xcf"+ // set cliff warning text
			"\xcd\xc6\x3a\xc0\x36\x9f\x2e\x46\x36\x3c"+ // init object
			"\x01\x00\xf1\x11\x0b\xd0\xcd\x1a\x22"+ // set position
			"\x3e\x50\xcd\x74\x0c"+ // play sound
			"\x21\xc0\xcf\xcb\xc6\xc9") // set $cfc0 bit and ret

	// upgrade normal items (interactions with ID 60) as necessary when they're
	// created.
	constMutables["set normal progressive call"] = MutableString(
		Addr{0x15, 0x465a}, "\x47\xcb\x37", "\xcd"+addrString(endOfBank15))
	endOfBank15 = appendCode("normal progressive func", 0x15, endOfBank15,
		"\x47\xcb\x37\xf5\x1e\x42\x1a\xcd\xe3\x3e\xf1\xc9")

	// should be set to match the western coast season
	varMutables["season after pirate cutscene"] = MutableByte(
		Addr{0x15, endOfBank15}, 0x15, 0x15)
	endOfBank15++
	// skip pirate cutscene. includes setting flag $1b, which makes the pirate
	// skull appear in the desert in case the player hasn't talked to the
	// ghost yet.
	constMutables["pirate flag call"] = MutableString(Addr{0x15, 0x5a0f},
		"\xcd\x30", addrString(endOfBank15))
	endOfBank15 = appendCode("pirate flag func", 0x15, endOfBank15,
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7"+
			"\xcb\xf6\xfa"+addrString(endOfBank15-1)+"\xea\x4e\xcc\xc9")

	// set sub ID for hard ore
	constMutables["hard ore id call"] = MutableString(Addr{0x15, 0x5b83},
		"\x2c\x36\x52", "\xcd"+addrString(endOfBank15))
	endOfBank15 = appendCode("hard ore id func", 0x15, endOfBank15,
		"\x2c\x36\x52\x2c\x36\x00\xc9")

	// load gfx data for randomized shop and market items.
	constMutables["item gfx call"] = MutableString(Addr{0x3f, 0x443c},
		"\x4f\x06\x00", "\xcd\x69\x71")
	constMutables["item gfx func"] = MutableString(Addr{0x3f, 0x7169}, "\x3f",
		// check for matching object
		"\x43\x4f\xcd\xdc\x71\x28\x17\x79\xfe\x59\x28\x19"+ // rod, woods
			"\xcd\xbf\x71\x28\x1b\xcd\xcf\x71\x28\x1d"+ // shops
			"\x79\xfe\x6e\x28\x1f\x06\x00\xc9"+ // feather
			// look up item ID, subID
			"\x1e\x15\x21"+addrString(endOfBank15)+"\x18\x1d"+
			"\x1e\x0b\x21\x8d\x7f\x18\x16"+
			"\x1e\x08\x21\xde\x7f\x18\x0f\x1e\x09\x21\xd1\x7f\x18\x08"+
			"\xfa\xb4\xc6\xc6\x15\x5f\x18\x0e"+ // feather
			"\xcd\x8a\x00"+ // get treasure
			"\x79\x4b\xcd\xd3\x3e\xcd\xe3\x3e\x23\x23\x5e"+ // get sprite
			"\x3e\x60\x4f\x06\x00\xc9") // replace object gfx w/ treasure gfx
	// return z if object is randomized shop item.
	constMutables["check randomized shop item"] = MutableString(
		Addr{0x3f, 0x71bf}, "\x3f",
		"\x79\xfe\x47\xc0\x7b\xb7\xc8\xfe\x02\xc8\xfe\x05\xc8\xfe\x0d\xc9")
	// same as above but for subrosia market.
	constMutables["check randomized market item"] = MutableString(
		Addr{0x3f, 0x71cf}, "\x3f",
		"\x79\xfe\x81\xc0\x7b\xb7\xc8\xfe\x04\xc8\xfe\x0d\xc9")
	// and rod of seasons.
	constMutables["check rod"] = MutableString(Addr{0x3f, 0x71dc}, "\x3f",
		"\x79\xfe\xe6\xc0\x7b\xfe\x02\xc9")
	// returns c,e = treasure ID,subID
	endOfBank15 = appendCode("rod lookup", 0x15, endOfBank15,
		"\x21\xcc\x70\x5e\x23\x23\x4e\xc9")
}
