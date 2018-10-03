package rom

// this file is for fixed mutables that go at the end of banks. each should be
// a self-contained unit (i.e. don't jr to anywhere outside the byte string) so
// that they can be appended automatically with respect to their size.

// return e.g. "\x2d\x79" for 0x792d
func addrString(addr uint16) string {
	return string([]byte{byte(addr), byte(addr >> 8)})
}

// adds code at the given address, returning the length of the byte string.
func addCode(name string, bank byte, offset uint16, code string) uint16 {
	constMutables[name] = MutableString(Addr{bank, offset},
		string([]byte{bank}), code)
	return uint16(len(code))
}

func initCode() {
	endOfBank00 := uint16(0x3ec8)
	endOfBank02 := uint16(0x75bb)
	endOfBank03 := uint16(0x7dd7)
	endOfBank04 := uint16(0x7e02)
	endOfBank05 := uint16(0x7e2d)
	endOfBank06 := uint16(0x77d4)
	endOfBank07 := uint16(0x78f0)
	endOfBank09 := uint16(0x7f4e)
	endOfBank11 := uint16(0x7eb0)
	endOfBank15 := uint16(0x792d)
	endOfBank3f := uint16(0x714d)

	// try to order these first by bank, then by call location.

	// bank 00

	// don't play any music if the -nomusic flag is given. because of this,
	// this *must* be the first function at the end of bank zero.
	endOfBank00 += addCode("no music func", 0x00, endOfBank00,
		"\x67\xfe\x40\x30\x03\x3e\x08\xc9\xf0\xb5\xc9")

	// force the item in the temple of seasons cutscene to use normal item
	// animations.
	constMutables["rod cutscene gfx call"] = MutableString(Addr{0x00, 0x2600},
		"\x1e\x41\x1a", "\xcd"+addrString(endOfBank00))
	endOfBank00 += addCode("rod cutscene gfx func", 0x00, endOfBank00,
		"\x1e\x41\x1a\xfe\xe6\xc0\x1c\x1a\xfe\x02\x28\x03\x1d\x1a\xc9"+
			"\x3e\x60\xc9")

	// set hl = address of treasure data + 1 for item with ID a, sub ID c.
	endOfBank00 += addCode("treasure data func", 0x00, endOfBank00,
		"\xf5\xc5\xd5\x47\x1e\x15\x21"+addrString(endOfBank15)+
			"\xcd\x8a\x00\xd1\xc1\xf1\xc9")
	endOfBank15 += addCode("treasure data body", 0x15, endOfBank15,
		"\x78\xc5\x21\x29\x51\xcd\xc3\x01\x09"+ // add ID offset
			"\xcb\x7e\x28\x09\x23\x2a\x66\x6f"+ // load as address if bit 7 set
			"\xc1\x79\xc5\x18\xef"+ // use sub ID as second offset
			"\x23\x06\x03\xd5\x11\xfd\xcd\xcd\x62\x04"+ // copy data
			"\x21\xfd\xcd\xd1\xc1\xc9") // set hl and ret

	// change hl to point to different treasure data if the item is progressive
	// and needs to be upgraded. param a = treasure ID.
	endOfBank00 += addCode("progressive item func", 0x00, endOfBank00,
		"\xd5\x5f\xcd\x6a\x3f\x7b\xd1\xd0"+ // ret if you don't have L-1
			"\xfe\x05\x20\x04\x21\x12\x3f\xc9"+ // check sword
			"\xfe\x06\x20\x04\x21\x15\x3f\xc9"+ // check boomerang
			"\xfe\x13\x20\x04\x21\x18\x3f\xc9"+ // check slingshot
			"\xfe\x17\x20\x04\x21\x1b\x3f\xc9"+ // check feather
			"\xfe\x19\xc0\x21\x1e\x3f\xc9"+ // check satchel
			// treasure data
			"\x02\x1d\x11\x02\x23\x1d\x02\x2f\x22\x02\x28\x17\x00\x46\x20")
	// use cape graphics for stolen feather if applicable.
	endOfBank00 += addCode("upgrade stolen feather func", 0x00, endOfBank00,
		"\xcd\x17\x17\xd8\xf5\x7b"+ // ret if you have the item
			"\xfe\x17\x20\x13\xd5\x1e\x43\x1a\xfe\x02\xd1\x20\x0a"+ // check IDs
			"\xfa\xb4\xc6\xfe\x02\x20\x03"+ // check feather level
			"\x21\x89\x3f\xf1\xc9"+ // set hl if match
			"\x02\x37\x17") // treasure data

	// this is a replacement for giveTreasure that gives treasure, plays sound,
	// and sets text based on item ID a and sub ID c, and accounting for item
	// progression.
	endOfBank00 += addCode("give item func", 0x00, endOfBank00,
		"\xcd\xd3\x3e\xcd\xe3\x3e"+ // get treasure data
			"\x4e\xcd\xeb\x16\x28\x05\xe5\xcd\x74\x0c\xe1"+ // give, play sound
			"\x06\x00\x23\x4e\xcd\x4b\x18\xaf\xc9") // show text

	// utility function, call a function hl in bank 02, preserving af. e can't
	// be used as a parameter to that function, but it can be returned.
	endOfBank00 += addCode("call bank 02", 0x00, endOfBank00,
		"\xf5\x1e\x02\xcd\x8a\x00\xf1\xc9")

	// utility function, read a byte from hl in bank e into a and e.
	endOfBank00 += addCode("read byte from bank", 0x00, endOfBank00,
		"\xfa\x97\xff\xf5\x7b\xea\x97\xff\xea\x22\x22"+ // switch bank
			"\x5e\xf1\xea\x97\xff\xea\x22\x22\x7b\xc9") // read and switch back

	// bank 02

	// warp to ember tree if holding start when closing the map screen, using
	// the playtime counter as a cooldown. this also sets the player's respawn
	// point.
	constMutables["outdoor map jump"] = MutableString(Addr{0x02, 0x6089},
		"\xc2\x7b\x4f", "\xc4"+addrString(endOfBank02))
	constMutables["dungeon map jump"] = MutableString(Addr{0x02, 0x602c},
		"\xc2\x7b\x4f", "\xc4"+addrString(endOfBank02))
	endOfBank02 += addCode("tree warp", 0x02, endOfBank02,
		"\xfa\x81\xc4\xe6\x08\x28\x33"+ // close as normal if start not held
			"\xfa\x49\xcc\xfe\x02\x30\x07"+ // check if indoors
			"\x21\x25\xc6\xcb\x7e\x28\x06"+ // check if cooldown is up
			"\x3e\x5a\xcd\x74\x0c\xc9"+ // play error sound and ret
			"\x21\x22\xc6\x11\xf8\x75\x06\x04\xcd\x5b\x04"+ // copy playtime
			"\x21\x2b\xc6\x11\xfc\x75\x06\x06\xcd\x5b\x04"+ // copy save point
			"\x21\xb7\xcb\x36\x05\xaf\xcd\x7b\x5e\xc3\x7b\x4f"+ // close + warp
			"\x40\xb4\xfc\xff\x00\xf8\x02\x02\x34\x38") // data for copies

	// warp to room under cursor if wearing developer ring.
	constMutables["dev ring warp call"] = MutableString(Addr{0x02, 0x5e9b},
		"\x89\x0c", addrString(endOfBank02))
	endOfBank02 += addCode("dev ring warp func", 0x02, endOfBank02,
		"\xfa\xc5\xc6\xfe\x40\x20\x12\xfa\x49\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x63\xcc\xfa\xb6\xcb\xea\x64\xcc\x3e\x03\xcd\x89\x0c\xc9")

	// bank 03

	// allow skipping the capcom screen after one second by pressing start
	constMutables["skip capcom call"] = MutableString(Addr{0x03, 0x4d6c},
		"\x37\x02", addrString(endOfBank03))
	endOfBank03 += addCode("skip capcom func", 0x03, endOfBank03,
		"\xe5\xfa\xb3\xcb\xfe\x94\x30\x03\xcd\x62\x08\xe1\xcd\x37\x02\xc9")

	// bank 04

	// if entering certain warps blocked by snow piles, mushrooms, or bushes,
	// set the animal companion to appear right outside instead of where you
	// left them. table entries are {entered group, entered room, animal room,
	// saved y, saved x}.
	animalSavePointTable := addrString(endOfBank04)
	endOfBank04 += addCode("animal save point table", 0x04, endOfBank04,
		"\x04\xfa\xc2\x18\x68\x00"+ // square jewel cave
			"\x05\xcc\x2a\x38\x18\x00"+ // goron mountain cave
			"\x05\xb3\x8e\x58\x88\x00"+ // cave outside d2
			"\x04\xe1\x86\x48\x68\x00"+ // quicksand ring cave
			"\x05\xc9\x2a\x38\x18\x00"+ // goron mountain main
			"\x05\xba\x2f\x18\x68\x00"+ // spring banana cave
			"\x05\xbb\x2f\x18\x68\x00"+ // joy ring cave
			"\x01\x05\x9a\x38\x48\x00"+ // rosa portal
			"\x04\x39\x8d\x38\x38\x00"+ // d2 entrance
			"\xff") // end
	constMutables["animal save point call"] = MutableString(Addr{0x04, 0x461e},
		"\xfa\x64\xcc", "\xcd"+addrString(endOfBank04))
	endOfBank04 += addCode("animal save point func", 0x04, endOfBank04,
		// b = group, c = room, d = animal room, hl = table
		"\xc5\xd5\x47\xfa\x64\xcc\x4f\xfa\x42\xcc\x57\x21"+animalSavePointTable+
			"\x2a\xb8\x20\x12\x2a\xb9\x20\x0e\x7e\xba\x20\x0a"+ // check criteria
			"\x11\x42\xcc\x06\x03\xcd\x62\x04\x18\x0a"+ // set save pt, done
			"\x2a\xb7\x20\xfc\x7e\x3c\x28\x02\x18\xe0"+ // go to next table entry
			"\x79\xd1\xc1\xc9") // done

	// set room flags so that rosa never appears in the overworld, and her
	// portal is activated by default.
	constMutables["set portal flag call"] = MutableString(Addr{0x04, 0x45f5},
		"\xfa\x64\xcc", "\xcd"+addrString(endOfBank04))
	endOfBank04 += addCode("set portal flag func", 0x04, endOfBank04,
		"\xe5\x21\x9a\xc7\x7e\xf6\x60\x77\x2e\xcb\x7e\xf6\xc0\x77"+ // set flags
			"\xe1\xfa\x64\xcc\xc9") // do what the address normally does

	// bank 05

	// do this so that animals don't immediately stop walking on screen when
	// called on a bridge.
	constMutables["ricky enter call"] = MutableString(Addr{0x05, 0x71ea},
		"\xcd\xaa\x44\xb7", "\xcd"+addrString(endOfBank05)+"\x00")
	constMutables["non-ricky enter call"] = MutableString(Addr{0x05, 0x493b},
		"\xcd\xaa\x44\xb7", "\xcd"+addrString(endOfBank05)+"\x00")
	endOfBank05 += addCode("flute enter func", 0x05, endOfBank05,
		"\xcd\xaa\x44\xb7\xc8\xfe\x1a\xc8\xfe\x1b\xc9")

	// bank 06

	// create a warning interaction when breaking bushes and flowers under
	// certain circumstances.
	constMutables["break bush warning call"] = MutableString(Addr{0x06, 0x477b},
		"\x21\x26\xc6", "\xcd"+addrString(endOfBank06))
	endOfBank06 += addCode("break bush warning func", 0x06, endOfBank06,
		"\xf5\xc5\xd5\xcd\xe1\x77\x21\x26\xc6\xd1\xc1\xf1\xc9"+ // wrapper
			"\xfe\xc3\x28\x09\xfe\xc4\x28\x05\xfe\xe5\x28\x01\xc9"+ // tile
			"\xfa\x4c\xcc\xfe\xa7\x28\x0d\xfe\x97\x28\x09"+ // jump by room
			"\xfe\x8d\x28\x05\xfe\x9a\x28\x01\xc9"+ // (cont.)
			"\x3e\x09\xcd\x17\x17\xd8"+ // "already warned" flag
			"\x3e\x16\xcd\x17\x17\xd8"+ // bracelet
			"\x3e\x0e\xcd\x17\x17\xd8"+ // flute
			"\x3e\x05\xcd\x17\x17\xd8"+ // sword
			"\x3e\x06\xcd\x17\x17\xfe\x02\xc8"+ // boomerang, L-2
			"\x21\x92\xc6\x3e\x09\xcd\x0e\x02"+ // set "already warned" flag
			"\x21\xe0\xcf\x36\x01"+ // set warning text index
			"\xcd\xc6\x3a\xc0\x36\x22\x2c\x36\x0a"+ // create warning object
			"\x2e\x4a\x11\x0a\xd0\x06\x04\xcd\x5b\x04"+ // place it on link
			"\xc9") // ret

	// bank 07

	// don't warp link using gale seeds if no trees have been reached (the menu
	// gets stuck in an infinite loop)
	constMutables["call gale seed check"] = MutableString(Addr{0x07, 0x4f45},
		"\xfa\x50\xcc\x3d", "\xcd"+addrString(endOfBank07)+"\x00")
	endOfBank07 += addCode("gale seed check", 0x07, endOfBank07,
		"\xfa\x50\xcc\x3d\xc0\xaf\x21\xf8\xc7\xb6\x21\x9e\xc7\xb6\x21\x72\xc7"+
			"\xb6\x21\x67\xc7\xb6\x21\x5f\xc7\xb6\x21\x10\xc7\xb6\xcb\x67"+
			"\x20\x02\x3c\xc9\xaf\xc9")

	// if wearing dev ring, change season regardless of where link is standing.
	constMutables["dev ring season call"] = MutableString(Addr{0x07, 0x5b75},
		"\xfa\xb6\xcc\xfe\x08", "\xcd"+addrString(endOfBank07)+"\x00\x00")
	endOfBank07 += addCode("dev ring season func", 0x07, endOfBank07,
		"\xfa\xc5\xc6\xfe\x40\xc8\xfa\xb6\xcc\xfe\x08\xc9")

	// bank 08

	// ORs the default season in the given area (low byte b in bank 1) with the
	// seasons the rod has (c), then ANDs and compares the results with d.
	warningHelperAddr := addrString(endOfBank15)
	endOfBank15 += addCode("warning helper", 0x15, endOfBank15,
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
	endOfBank15 += addCode("warning func", 0x15, endOfBank15,
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

	// bank 09

	// animals called by flute normally veto any nonzero collision value for
	// the purposes of entering a screen, but this allows double-wide bridges
	// (1a and 1b) as well. this specifically fixes the problem of not being
	// able to call an animal on the d1 screen, or on the bridge to the screen
	// to the right. the vertical collision check isn't modified, since bridges
	// only run horizontally.
	constMutables["flute collision calls"] = MutableStrings(
		[]Addr{{0x09, 0x4d9a}, {0x09, 0x4dad}},
		"\xcd\xd9\x4e", "\xcd"+addrString(endOfBank09))
	endOfBank09 += addCode("flute collision func", 0x09, endOfBank09,
		"\x06\x01\x7e\xfe\x1a\x28\x06\xfe\x1b\x28\x02\xb7\xc0"+ // first tile
			"\x7d\x80\x6f\x7e\xfe\x1a\x28\x05\xfe\x1b\x28\x01\xb7"+ // second
			"\x7d\xc0\xcd\x89\x20\xaf\xc9") // vanilla stuff

	// if wearing dev ring, warp to animal companion if it's already in the
	// same room when playing the flute.
	constMutables["dev ring flute call"] = MutableString(Addr{0x09, 0x4e2c},
		"\xd9\x3a", addrString(endOfBank09))
	endOfBank09 += addCode("dev ring flute func", 0x09, endOfBank09,
		"\xd5\xfa\xc5\xc6\xfe\x40\x20\x07"+ // check dev ring
			"\xfa\x04\xd1\xfe\x01\x28\x04"+ // check animal companion
			"\xd1\xc3\xd9\x3a"+ // done
			"\xcd\xc6\x3a\x20\x0c\x36\x05"+ // create poof
			"\x11\x0a\xd0\x2e\x4a\x06\x04\xcd\x5b\x04"+ // move poof
			"\x11\x0a\xd1\x21\x0a\xd0\x06\x04\xcd\x5b\x04"+ // move animal
			"\x18\xde") // jump to done

	// remove star ore from inventory when buying the first subrosian market
	// item. this can't go in the gain/lose items table, since the given item
	// doesn't necessarily have a unique ID.
	constMutables["trade star ore call"] = MutableString(Addr{0x09, 0x7887},
		"\xdf\x2a\x4e", "\xcd"+addrString(endOfBank09))
	endOfBank09 += addCode("trade star ore func", 0x09, endOfBank09,
		"\xb7\x20\x07\xe5\x21\x9a\xc6\xcb\xae\xe1\xdf\x2a\x4e\xc9")

	// bank 11

	// the interaction on the mount cucco waterfall/vine screen
	constMutables["waterfall cliff interaction redirect"] =
		MutableString(Addr{0x11, 0x6c10},
			"\xf2\x1f\x08\x68", "\xf3"+addrString(endOfBank11)+"\xff")
	endOfBank11 += addCode("waterfall cliff interactions", 0x11, endOfBank11,
		"\xf2\x1f\x08\x68\x68\x22\x0a\x20\x18\xfe")
	// natzu / woods of winter cliff
	constMutables["flower cliff interaction redirect"] =
		MutableString(Addr{0x11, 0x6568},
			"\xf2\x9c\x00\x58", "\xf3"+addrString(endOfBank11)+"\xff")
	endOfBank11 += addCode("flower cliff interactions", 0x11, endOfBank11,
		"\xf2\x9c\x00\x58\x58\x22\x0a\x30\x58\xfe")
	// sunken city diving spot
	constMutables["diving spot interaction redirect"] =
		MutableString(Addr{0x11, 0x69cc},
			"\xf2\x1f\x0d\x68", "\xf3"+addrString(endOfBank11)+"\xff")
	endOfBank11 += addCode("diving spot interactions", 0x11, endOfBank11,
		"\xf2\x1f\x0d\x68\x68\x3e\x31\x18\x68\x22\x0a\x64\x68\xfe")
	// moblin keep -> sunken city
	constMutables["moblin keep interaction redirect"] =
		MutableString(Addr{0x11, 0x650b},
			"\xf2\xab\x00\x40", "\xf3"+addrString(endOfBank11)+"\xff")
	endOfBank11 += addCode("moblin keep interactions", 0x11, endOfBank11,
		"\xf2\xab\x00\x40\x70\x22\x0a\x58\x44\xf8\x2d\x00\x33\xfe")
	// hss skip room
	constMutables["hss skip room interaction redirect"] =
		MutableString(Addr{0x11, 0x7ada},
			"\xf3\x93\x55", "\xf3"+addrString(endOfBank11))
	endOfBank11 += addCode("hss skip room interactions", 0x11, endOfBank11,
		"\xf2\x22\x0a\x88\x98\xf3\x93\x55\xfe")

	// bank 15

	// upgrade normal items (interactions with ID 60) as necessary when they're
	// created.
	constMutables["set normal progressive call"] = MutableString(
		Addr{0x15, 0x465a}, "\x47\xcb\x37", "\xcd"+addrString(endOfBank15))
	endOfBank15 += addCode("normal progressive func", 0x15, endOfBank15,
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
	endOfBank15 += addCode("pirate flag func", 0x15, endOfBank15,
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7"+
			"\xcb\xf6\xfa"+addrString(endOfBank15-1)+"\xea\x4e\xcc\xc9")

	// set sub ID for hard ore
	constMutables["hard ore id call"] = MutableString(Addr{0x15, 0x5b83},
		"\x2c\x36\x52", "\xcd"+addrString(endOfBank15))
	endOfBank15 += addCode("hard ore id func", 0x15, endOfBank15,
		"\x2c\x36\x52\x2c\x36\x00\xc9")

	// bank 3f

	// have seed satchel inherently refill all seeds.
	constMutables["satchel refill call"] = MutableString(Addr{0x00, 0x16f6},
		"\xcd\xc8\x44", "\xcd"+addrString(endOfBank3f))
	endOfBank3f += addCode("satchel seed refill func", 0x3f, endOfBank3f,
		"\xc5\xcd\xc8\x44\x78\xc1\xf5\x78\xfe\x19\x20\x07"+
			"\xc5\xd5\xcd\xe5\x17\xd1\xc1\xf1\x47\xc9")

	// load gfx data for randomized shop and market items.
	constMutables["item gfx call"] = MutableString(Addr{0x3f, 0x443c},
		"\x4f\x06\x00", "\xcd"+addrString(endOfBank3f))
	endOfBank3f += addCode("item gfx func", 0x3f, endOfBank3f,
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
	endOfBank3f += addCode("check randomized shop item", 0x3f, endOfBank3f,
		"\x79\xfe\x47\xc0\x7b\xb7\xc8\xfe\x02\xc8\xfe\x05\xc8\xfe\x0d\xc9")
	// same as above but for subrosia market.
	endOfBank3f += addCode("check randomized market item", 0x3f, endOfBank3f,
		"\x79\xfe\x81\xc0\x7b\xb7\xc8\xfe\x04\xc8\xfe\x0d\xc9")
	// and rod of seasons.
	endOfBank3f += addCode("check rod", 0x3f, endOfBank3f,
		"\x79\xfe\xe6\xc0\x7b\xfe\x02\xc9")
	// returns c,e = treasure ID,subID
	endOfBank15 += addCode("rod lookup", 0x15, endOfBank15,
		"\x21\xcc\x70\x5e\x23\x23\x4e\xc9")

	// "activate" a flute by setting its icon and song when obtained. also
	// activates the corresponding animal companion.
	constMutables["flute set icon call"] = MutableString(Addr{0x3f, 0x452c},
		"\x4e\x45", addrString(endOfBank3f))
	endOfBank3f += addCode("flute set icon func", 0x3f, endOfBank3f,
		"\xf5\xd5\xe5\x78\xfe\x0e\x20\x0d\x1e\xaf\x79\xd6\x0a\x12\xc6\x42"+
			"\x26\xc6\x6f\xcb\xfe\xe1\xd1\xf1\xcd\x4e\x45\xc9")
}
