package rom

import (
	"fmt"
)

// this file is for mutables that go at the end of banks. each should be a
// self-contained unit (i.e. don't jr to anywhere outside the byte string) so
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

type ROM struct {
	endOfBank []uint16
	mutables  map[string]Mutable
}

// New returns a newly initialized ROM.
func New() *ROM {
	r := ROM{
		endOfBank: make([]uint16, 0x40),
		mutables:  make(map[string]Mutable),
	}

	r.endOfBank[0x00] = 0x3ec8
	r.endOfBank[0x01] = 0x7e89
	r.endOfBank[0x02] = 0x75bb
	r.endOfBank[0x03] = 0x7dd7
	r.endOfBank[0x04] = 0x7e02
	r.endOfBank[0x05] = 0x7e2d
	r.endOfBank[0x06] = 0x77d4
	r.endOfBank[0x07] = 0x78f0
	r.endOfBank[0x09] = 0x7f4e
	r.endOfBank[0x11] = 0x7eb0
	r.endOfBank[0x15] = 0x792d
	r.endOfBank[0x3f] = 0x714d

	for k, v := range constMutables {
		r.mutables[k] = v
	}
	r.initEndOfBank()

	return &r
}

// appendToBank appends the given data to the end of the given bank, associates
// it with the given name, and returns the address of the data as a string such
// as "\xc8\x3e" for 0x3ec8. it panics if the end of the bank is zero or if the
// data would overflow the bank.
func (r *ROM) appendToBank(bank byte, name, data string) string {
	eob := r.endOfBank[bank]

	if eob == 0 {
		panic(fmt.Sprintf("end of bank %02x undefined for %s", bank, name))
	}

	if eob+uint16(len(data)) > 0x8000 {
		panic(fmt.Sprintf("not enough space for %s in bank %02x", name, bank))
	}

	r.mutables[name] =
		MutableString(Addr{bank, eob}, string([]byte{bank}), data)
	r.endOfBank[bank] += uint16(len(data))

	return addrString(eob)
}

// replace replaces the old data at the given address with the new data, and
// associates the change with the given name. actual replacement will fail at
// runtime if the old data does not match the original data in the ROM.
func (r *ROM) replace(bank byte, offset uint16, name, old, new string) {
	r.mutables[name] = MutableString(Addr{bank, offset}, old, new)
}

// replaceMultiple acts as replace, but operates on multiple addresses.
func (r *ROM) replaceMultiple(addrs []Addr, name, old, new string) {
	r.mutables[name] = MutableStrings(addrs, old, new)
}

// initEndOfBank adds end-of-bank mutables and mutables that point to them.
func (r *ROM) initEndOfBank() {
	// try to order these first by bank, then by call location. maybe group
	// them into subfunctions when applicable?

	// bank 00

	// don't play any music if the -nomusic flag is given. because of this,
	// this *must* be the first function at the end of bank zero (for now).
	r.appendToBank(0x00, "no music func",
		"\x67\xfe\x40\x30\x03\x3e\x08\xc9\xf0\xb5\xc9")
	r.replace(0x00, 0x0c76, "no music call",
		"\x67\xf0\xb5", "\x67\xf0\xb5") // modified only by SetNoMusic()

	// force the item in the temple of seasons cutscene to use normal item
	// animations.
	rodCutsceneGfxFunc := r.appendToBank(0x00, "rod cutscene gfx func",
		"\x1e\x41\x1a\xfe\xe6\xc0\x1c\x1a\xfe\x02\x28\x03\x1d\x1a\xc9"+
			"\x3e\x60\xc9")
	r.replace(0x00, 0x2600, "rod cutscene gfx call",
		"\x1e\x41\x1a", "\xcd"+rodCutsceneGfxFunc)

	// set hl = address of treasure data + 1 for item with ID a, sub ID c.
	treasureDataBody := r.appendToBank(0x15, "treasure data body",
		"\x78\xc5\x21\x29\x51\xcd\xc3\x01\x09"+ // add ID offset
			"\xcb\x7e\x28\x09\x23\x2a\x66\x6f"+ // load as address if bit 7 set
			"\xc1\x79\xc5\x18\xef"+ // use sub ID as second offset
			"\x23\x06\x03\xd5\x11\xfd\xcd\xcd\x62\x04"+ // copy data
			"\x21\xfd\xcd\xd1\xc1\xc9") // set hl and ret
	getTreasureData := r.appendToBank(0x00, "treasure data func",
		"\xf5\xc5\xd5\x47\x1e\x15\x21"+treasureDataBody+
			"\xcd\x8a\x00\xd1\xc1\xf1\xc9")

	// use cape graphics for stolen feather if applicable.
	upgradeFeather := r.appendToBank(0x00, "upgrade stolen feather func",
		"\xcd\x17\x17\xd8\xf5\x7b"+ // ret if you have the item
			"\xfe\x17\x20\x13\xd5\x1e\x43\x1a\xfe\x02\xd1\x20\x0a"+ // check IDs
			"\xfa\xb4\xc6\xfe\x02\x20\x03"+ // check feather level
			"\x21\x89\x3f\xf1\xc9"+ // set hl if match
			"\x02\x37\x17") // treasure data
	// change hl to point to different treasure data if the item is progressive
	// and needs to be upgraded. param a = treasure ID.
	progressiveItemFunc := r.appendToBank(0x00, "progressive item func",
		"\xd5\x5f\xcd"+upgradeFeather+"\x7b\xd1\xd0"+ // ret if missing L-1
			"\xfe\x05\x20\x04\x21\x12\x3f\xc9"+ // check sword
			"\xfe\x06\x20\x04\x21\x15\x3f\xc9"+ // check boomerang
			"\xfe\x13\x20\x04\x21\x18\x3f\xc9"+ // check slingshot
			"\xfe\x17\x20\x04\x21\x1b\x3f\xc9"+ // check feather
			"\xfe\x19\xc0\x21\x1e\x3f\xc9"+ // check satchel
			// treasure data
			"\x02\x1d\x11\x02\x23\x1d\x02\x2f\x22\x02\x28\x17\x00\x46\x20")

	// this is a replacement for giveTreasure that gives treasure, plays sound,
	// and sets text based on item ID a and sub ID c, and accounting for item
	// progression.
	giveItem := r.appendToBank(0x00, "give item func",
		"\xcd"+getTreasureData+"\xcd"+progressiveItemFunc+ // get treasure data
			"\x4e\xcd\xeb\x16\x28\x05\xe5\xcd\x74\x0c\xe1"+ // give, play sound
			"\x06\x00\x23\x4e\xcd\x4b\x18\xaf\xc9") // show text

	// utility function, call a function hl in bank 02, preserving af. e can't
	// be used as a parameter to that function, but it can be returned.
	callBank2 := r.appendToBank(0x00, "call bank 02",
		"\xf5\x1e\x02\xcd\x8a\x00\xf1\xc9")

	// utility function, read a byte from hl in bank e into a and e.
	r.appendToBank(0x00, "read byte from bank",
		"\xfa\x97\xff\xf5\x7b\xea\x97\xff\xea\x22\x22"+ // switch bank
			"\x5e\xf1\xea\x97\xff\xea\x22\x22\x7b\xc9") // read and switch back

	// bank 01

	// helper function, takes b = high byte of season addr, returns season in b
	readSeason := r.appendToBank(0x01, "read default season",
		"\x26\x7e\x68\x7e\x47\xc9")

	// bank 02

	// warp to ember tree if holding start when closing the map screen, using
	// the playtime counter as a cooldown. this also sets the player's respawn
	// point.
	treeWarp := r.appendToBank(0x02, "tree warp",
		"\xfa\x81\xc4\xe6\x08\x28\x33"+ // close as normal if start not held
			"\xfa\x49\xcc\xfe\x02\x30\x07"+ // check if indoors
			"\x21\x25\xc6\xcb\x7e\x28\x06"+ // check if cooldown is up
			"\x3e\x5a\xcd\x74\x0c\xc9"+ // play error sound and ret
			"\x21\x22\xc6\x11\xf8\x75\x06\x04\xcd\x5b\x04"+ // copy playtime
			"\x21\x2b\xc6\x11\xfc\x75\x06\x06\xcd\x5b\x04"+ // copy save point
			"\x21\xb7\xcb\x36\x05\xaf\xcd\x7b\x5e\xc3\x7b\x4f"+ // close + warp
			"\x40\xb4\xfc\xff\x00\xf8\x02\x02\x34\x38") // data for copies
	r.replaceMultiple([]Addr{{0x02, 0x6089}, {0x02, 0x602c}}, "tree warp jump",
		"\xc2\x7b\x4f", "\xc4"+treeWarp)

	// warp to room under cursor if wearing developer ring.
	devWarp := r.appendToBank(0x02, "dev ring warp func",
		"\xfa\xc5\xc6\xfe\x40\x20\x12\xfa\x49\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x63\xcc\xfa\xb6\xcb\xea\x64\xcc\x3e\x03\xcd\x89\x0c\xc9")
	r.replace(0x02, 0x5e9b, "dev ring warp call", "\x89\x0c", devWarp)

	// bank 03

	// allow skipping the capcom screen after one second by pressing start
	skipCapcom := r.appendToBank(0x03, "skip capcom func",
		"\xe5\xfa\xb3\xcb\xfe\x94\x30\x03\xcd\x62\x08\xe1\xcd\x37\x02\xc9")
	r.replace(0x03, 0x4d6c, "skip capcom call", "\x37\x02", skipCapcom)

	// bank 04

	// if entering certain warps blocked by snow piles, mushrooms, or bushes,
	// set the animal companion to appear right outside instead of where you
	// left them. table entries are {entered group, entered room, animal room,
	// saved y, saved x}.
	animalSaveTable := r.appendToBank(0x04, "animal save point table",
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
	animalSaveFunc := r.appendToBank(0x04, "animal save point func",
		// b = group, c = room, d = animal room, hl = table
		"\xc5\xd5\x47\xfa\x64\xcc\x4f\xfa\x42\xcc\x57\x21"+animalSaveTable+
			"\x2a\xb8\x20\x12\x2a\xb9\x20\x0e\x7e\xba\x20\x0a"+ // check criteria
			"\x11\x42\xcc\x06\x03\xcd\x62\x04\x18\x0a"+ // set save pt, done
			"\x2a\xb7\x20\xfc\x7e\x3c\x28\x02\x18\xe0"+ // go to next table entry
			"\x79\xd1\xc1\xc9") // done
	r.replace(0x04, 0x461e, "animal save point call",
		"\xfa\x64\xcc", "\xcd"+animalSaveFunc)

	// set room flags so that rosa never appears in the overworld, and her
	// portal is activated by default.
	setPortalFlags := r.appendToBank(0x04, "set portal flag func",
		"\xe5\x21\x9a\xc7\x7e\xf6\x60\x77\x2e\xcb\x7e\xf6\xc0\x77"+ // set flags
			"\xe1\xfa\x64\xcc\xc9") // do what the address normally does
	r.replace(0x04, 0x45f5, "set portal flag call",
		"\xfa\x64\xcc", "\xcd"+setPortalFlags)

	// bank 05

	// do this so that animals don't immediately stop walking on screen when
	// called on a bridge.
	fluteEnterFunc := r.appendToBank(0x05, "flute enter func",
		"\xcd\xaa\x44\xb7\xc8\xfe\x1a\xc8\xfe\x1b\xc9")
	r.replaceMultiple([]Addr{{0x05, 0x71ea}, {0x05, 0x493b}},
		"animal enter call", "\xcd\xaa\x44\xb7", "\xcd"+fluteEnterFunc+"\x00")

	// bank 06

	// create a warning interaction when breaking bushes and flowers under
	// certain circumstances.
	bushWarningBody := r.appendToBank(0x06, "break bush warning body",
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
	bushWarning := r.appendToBank(0x06, "break bush warning func",
		"\xf5\xc5\xd5\xcd"+bushWarningBody+"\x21\x26\xc6\xd1\xc1\xf1\xc9")
	r.replace(0x06, 0x477b, "bush warning call",
		"\x21\x26\xc6", "\xcd"+bushWarning)

	// bank 07

	// don't warp link using gale seeds if no trees have been reached (the menu
	// gets stuck in an infinite loop)
	galeSeedCheck := r.appendToBank(0x07, "gale seed check",
		"\xfa\x50\xcc\x3d\xc0\xaf\x21\xf8\xc7\xb6\x21\x9e\xc7\xb6\x21\x72\xc7"+
			"\xb6\x21\x67\xc7\xb6\x21\x5f\xc7\xb6\x21\x10\xc7\xb6\xcb\x67"+
			"\x20\x02\x3c\xc9\xaf\xc9")
	r.replace(0x07, 0x4f45, "call gale seed check",
		"\xfa\x50\xcc\x3d", "\xcd"+galeSeedCheck+"\x00")

	// if wearing dev ring, change season regardless of where link is standing.
	devChangeSeason := r.appendToBank(0x07, "dev ring season func",
		"\xfa\xc5\xc6\xfe\x40\xc8\xfa\xb6\xcc\xfe\x08\xc9")
	r.replace(0x07, 0x5b75, "dev ring season call",
		"\xfa\xb6\xcc\xfe\x08", "\xcd"+devChangeSeason+"\x00\x00")

	// bank 08

	// use the custom "give item" function in the shop instead of the normal
	// one. this obviates some hard-coded shop data (sprite, text) and allows
	// the item to progressively upgrade.
	// param = b (item index/subID), returns c,e = treasure ID,subID
	shopLookup := r.appendToBank(0x08, "shop item lookup",
		"\x21\xce\x4c\x78\x87\xd7\x4e\x23\x5e\xc9")
	shopCheckAddr := r.appendToBank(0x08, "shop check addr",
		"\xfe\xe9\xc8\xfe\xcf\xc8\xfe\xd3\xc8\xfe\xd9\xc9")
	shopGiveItem := r.appendToBank(0x08, "shop give item func",
		"\xc5\x47\x7d\xcd"+shopCheckAddr+"\x78\xc1\x28\x04\xcd\xeb\x16\xc9"+
			"\xcd\x21\x3f\xc9") // give item and ret
	r.replace(0x08, 0x4bfc, "shop give item call",
		"\xeb\x16", shopGiveItem)

	// give fake treasure 0f for the strange flute item.
	shopIDFunc := r.appendToBank(0x08, "shop give fake id func",
		"\x1e\x42\x1a\xfe\x0d\xc0\x21\x93\xc6\xcb\xfe\xc9")
	r.replace(0x08, 0x4bfe, "shop give fake id call",
		"\x1e\x42\x1a", "\xcd"+shopIDFunc)

	// ORs the default season in the given area (low byte b in bank 1) with the
	// seasons the rod has (c), then ANDs and compares the results with d.
	warningHelper := r.appendToBank(0x15, "warning helper",
		"\x1e\x01\x21"+readSeason+"\xcd\x8a\x00"+ // get default season
			"\x78\xb7\x3e\x01\x28\x05\xcb\x27\x05\x20\xfb"+ // match rod format
			"\xb1\xa2\xba\xc9") // OR with c, AND with d, compare with d, ret
	// this communicates with the warning script by setting bit zero of $cfc0
	// if the warning needs to be displayed (based on room, season, etc), and
	// also displays the exclamation mark if so.
	warningFunc := r.appendToBank(0x15, "warning func",
		"\xc5\xd5\xcd"+addrString(r.endOfBank[0x15]+8)+"\xd1\xc1\xc9"+ // wrap
			"\xfa\x4e\xcc\x47\xfa\xb0\xc6\x4f\xfa\x4c\xcc"+ // load env data
			"\xfe\x7c\x28\x12\xfe\x6e\x28\x18\xfe\x3d\x28\x22"+ // jump by room
			"\xfe\x5c\x28\x28\xfe\x78\x28\x32\x18\x43"+ // (cont.)
			"\x06\x61\x16\x01\xcd"+warningHelper+"\xc8\x18\x35"+ // flower cliff
			"\x78\xfe\x03\xc8\x06\x61\x16\x09\xcd"+warningHelper+
			"\xc8\x18\x27"+ // diving spot
			"\x06\x65\x16\x02\xcd"+warningHelper+
			"\xc8\x18\x1d"+ // waterfall cliff
			"\xfa\x10\xc6\xfe\x0c\xc0\x3e\x17\xcd\x17\x17\xd8\x18\x0f"+ // keep
			"\xcd\x56\x19\xcb\x76\xc0\xcb\xf6\x3e\x02\xea\xe0\xcf"+
			"\x18\x04"+ // hss skip room
			"\xaf\xea\xe0\xcf"+ // set cliff warning text
			"\xcd\xc6\x3a\xc0\x36\x9f\x2e\x46\x36\x3c"+ // init object
			"\x01\x00\xf1\x11\x0b\xd0\xcd\x1a\x22"+ // set position
			"\x3e\x50\xcd\x74\x0c"+ // play sound
			"\x21\xc0\xcf\xcb\xc6\xc9") // set $cfc0 bit and ret
	// overwrite unused maku gate interaction with warning interaction
	warningScript := r.appendToBank(0x0b, "warningScript",
		"\xd0\xe0"+warningFunc+
			"\xa0\xbd\xd7\x3c"+ // wait for collision and animation
			"\x87\xe0\xcf\x7e\x7f\x83\x7f\x88\x7f"+ // jump based on cfe0 bits
			"\x98\x26\x00\xbe\x00"+ // show cliff warning text
			"\x98\x26\x01\xbe\x00"+ // show bush warning text
			"\x98\x26\x02\xbe\x00") // show hss skip warning text
	r.replace(0x08, 0x5663, "warning script pointer", "\x87\x4e", warningScript)

	// bank 09

	// shared by maku tree and star-shaped ore.
	bank2IDFunc := r.appendToBank(0x02, "bank 2 fake id func",
		"\xfa\x49\xcc\xfe\x01\x28\x05\xfe\x02\x28\x1b\xc9"+ // compare group
			"\xfa\x4c\xcc\xfe\x65\x28\x0d\xfe\x66\x28\x09"+ // compare room
			"\xfe\x75\x28\x05\xfe\x76\x28\x01\xc9"+ // cont.
			"\x21\x94\xc6\xcb\xd6\xc9"+ // set treasure id 12
			"\xfa\x4c\xcc\xfe\x0b\xc0\x21\x93\xc6\xcb\xd6\xc9") // id 0a
	bank9IDFunc := r.appendToBank(0x09, "bank 9 fake id func",
		"\xf5\xe5\x21"+bank2IDFunc+"\xcd"+callBank2+"\xe1\xf1\xcd\xeb\x16\xc9")
	r.replace(0x09, 0x42e1, "bank 9 fake id call", "\xeb\x16", bank9IDFunc)

	// animals called by flute normally veto any nonzero collision value for
	// the purposes of entering a screen, but this allows double-wide bridges
	// (1a and 1b) as well. this specifically fixes the problem of not being
	// able to call an animal on the d1 screen, or on the bridge to the screen
	// to the right. the vertical collision check isn't modified, since bridges
	// only run horizontally.
	fluteCollisionFunc := r.appendToBank(0x09, "flute collision func",
		"\x06\x01\x7e\xfe\x1a\x28\x06\xfe\x1b\x28\x02\xb7\xc0"+ // first tile
			"\x7d\x80\x6f\x7e\xfe\x1a\x28\x05\xfe\x1b\x28\x01\xb7"+ // second
			"\x7d\xc0\xcd\x89\x20\xaf\xc9") // vanilla stuff
	r.replaceMultiple([]Addr{{0x09, 0x4d9a}, {0x09, 0x4dad}},
		"flute collision calls", "\xcd\xd9\x4e", "\xcd"+fluteCollisionFunc)

	// if wearing dev ring, warp to animal companion if it's already in the
	// same room when playing the flute.
	devFluteWarp := r.appendToBank(0x09, "dev ring flute func",
		"\xd5\xfa\xc5\xc6\xfe\x40\x20\x07"+ // check dev ring
			"\xfa\x04\xd1\xfe\x01\x28\x04"+ // check animal companion
			"\xd1\xc3\xd9\x3a"+ // done
			"\xcd\xc6\x3a\x20\x0c\x36\x05"+ // create poof
			"\x11\x0a\xd0\x2e\x4a\x06\x04\xcd\x5b\x04"+ // move poof
			"\x11\x0a\xd1\x21\x0a\xd0\x06\x04\xcd\x5b\x04"+ // move animal
			"\x18\xde") // jump to done
	r.replace(0x09, 0x4e2c, "dev ring flute call", "\xd9\x3a", devFluteWarp)

	// remove star ore from inventory when buying the first subrosian market
	// item. this can't go in the gain/lose items table, since the given item
	// doesn't necessarily have a unique ID.
	tradeStarOre := r.appendToBank(0x09, "trade star ore func",
		"\xb7\x20\x07\xe5\x21\x9a\xc6\xcb\xae\xe1\xdf\x2a\x4e\xc9")
	r.replace(0x09, 0x7887, "trade star ore call",
		"\xdf\x2a\x4e", "\xcd"+tradeStarOre)

	// use custom "give item" func in the subrosian market.
	marketFinalGiveItem := r.appendToBank(0x09, "market final give item",
		"\xf1\xcd"+giveItem+"\xd1\x37\xc9") // give item, scf, ret
	marketIDFunc := r.appendToBank(0x09, "market give fake id func",
		"\xe5\x21\x94\xc6\xcb\xc6\xe1\xca"+marketFinalGiveItem)
	// param = b (item index/subID), returns c,e = treasure ID,subID
	marketLookup := r.appendToBank(0x09, "market item lookup",
		"\x21\xda\x77\x78\x87\xd7\x4e\x23\x5e\xc9")
	marketGiveItem := r.appendToBank(0x09, "market give item func",
		"\xf5\x7d\xfe\xdb\xca"+marketFinalGiveItem+
			"\xfe\xe3\xca"+marketFinalGiveItem+"\xfe\xf5\xca"+marketIDFunc+
			"\xf1\xfe\x2d\x20\x03\xcd\xb9\x17\xcd\xeb\x16\x1e\x42\xc9")
	r.replace(0x09, 0x788a, "market give item call",
		"\xfe\x2d\x20\x03\xcd\xb9\x17\xcd\xeb\x16\x1e\x42",
		"\x00\x00\x00\x00\x00\x00\x00\xcd"+marketGiveItem+"\x38\x0b")

	// bank 0b

	diverIDScript := r.appendToBank(0x0b, "diver give fake id script",
		"\xde\x2e\x00\x92\x94\xc6\x02\xc1")
	r.replace(0x0b, 0x730d, "diver give fake id call",
		"\xde\x2e\x00", "\xc0"+diverIDScript)

	// returns c,e = treasure ID,subID
	nobleSwordLookup := r.appendToBank(0x0b, "noble sword lookup",
		"\x21\x18\x64\x4e\x23\x5e\xc9")

	// bank 11

	// the interaction on the mount cucco waterfall/vine screen
	waterfallInteractions := r.appendToBank(0x11, "waterfall interactions",
		"\xf2\x1f\x08\x68\x68\x22\x0a\x20\x18\xfe")
	r.replace(0x11, 0x6c10, "waterfall cliff interaction jump",
		"\xf2\x1f\x08\x68", "\xf3"+waterfallInteractions+"\xff")
	// natzu / woods of winter cliff
	flowerCliffInteractions := r.appendToBank(0x11, "flower cliff interactions",
		"\xf2\x9c\x00\x58\x58\x22\x0a\x30\x58\xfe")
	r.replace(0x11, 0x6568, "flower cliff interaction jump",
		"\xf2\x9c\x00\x58", "\xf3"+flowerCliffInteractions+"\xff")
	// sunken city diving spot
	divingSpotInteractions := r.appendToBank(0x11, "diving spot interactions",
		"\xf2\x1f\x0d\x68\x68\x3e\x31\x18\x68\x22\x0a\x64\x68\xfe")
	r.replace(0x11, 0x69cc, "diving spot interaction jump",
		"\xf2\x1f\x0d\x68", "\xf3"+divingSpotInteractions+"\xff")
	// moblin keep -> sunken city
	moblinKeepInteractions := r.appendToBank(0x11, "moblin keep interactions",
		"\xf2\xab\x00\x40\x70\x22\x0a\x58\x44\xf8\x2d\x00\x33\xfe")
	r.replace(0x11, 0x650b, "moblin keep interaction jump",
		"\xf2\xab\x00\x40", "\xf3"+moblinKeepInteractions+"\xff")
	// hss skip room
	hssSkipInteractions := r.appendToBank(0x11, "hss skip interactions",
		"\xf2\x22\x0a\x88\x98\xf3\x93\x55\xfe")
	r.replace(0x11, 0x7ada, "hss skip interaction jump",
		"\xf3\x93\x55", "\xf3"+hssSkipInteractions)

	// bank 15

	// upgrade normal items (interactions with ID 60) as necessary when they're
	// created.
	normalProgressiveFunc := r.appendToBank(0x15, "normal progressive func",
		"\x47\xcb\x37\xf5\x1e\x42\x1a\xcd"+progressiveItemFunc+"\xf1\xc9")
	r.replace(0x15, 0x465a, "set normal progressive call",
		"\x47\xcb\x37", "\xcd"+normalProgressiveFunc)

	// should be set to match the western coast season
	pirateSeason := r.appendToBank(0x15, "season after pirate cutscene", "\x15")
	// skip pirate cutscene. includes setting flag $1b, which makes the pirate
	// skull appear in the desert in case the player hasn't talked to the
	// ghost yet.
	pirateFlagFunc := r.appendToBank(0x15, "pirate flag func",
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7"+
			"\xcb\xf6\xfa"+pirateSeason+"\xea\x4e\xcc\xc9")
	r.replace(0x15, 0x5a0f, "pirate flag call", "\xcd\x30", pirateFlagFunc)

	// set sub ID for hard ore
	hardOreFunc := r.appendToBank(0x15, "hard ore id func",
		"\x2c\x36\x52\x2c\x36\x00\xc9")
	r.replace(0x15, 0x5b83, "hard ore id call",
		"\x2c\x36\x52", "\xcd"+hardOreFunc)

	// use custom "give item" func in rod cutscene.
	r.replace(0x15, 0x70cf, "rod give item call",
		"\xcd\xeb\x16", "\xcd"+giveItem)

	// bank 3f

	// have seed satchel inherently refill all seeds.
	satchelRefill := r.appendToBank(0x3f, "satchel seed refill func",
		"\xc5\xcd\xc8\x44\x78\xc1\xf5\x78\xfe\x19\x20\x07"+
			"\xc5\xd5\xcd\xe5\x17\xd1\xc1\xf1\x47\xc9")
	r.replace(0x00, 0x16f6, "satchel refill call",
		"\xcd\xc8\x44", "\xcd"+satchelRefill)

	// returns c,e = treasure ID,subID
	rodLookup := r.appendToBank(0x15, "rod lookup",
		"\x21\xcc\x70\x5e\x23\x23\x4e\xc9")
	// return z if object is randomized shop item.
	checkShopItem := r.appendToBank(0x3f, "check randomized shop item",
		"\x79\xfe\x47\xc0\x7b\xb7\xc8\xfe\x02\xc8\xfe\x05\xc8\xfe\x0d\xc9")
	// same as above but for subrosia market.
	checkMarketItem := r.appendToBank(0x3f, "check randomized market item",
		"\x79\xfe\x81\xc0\x7b\xb7\xc8\xfe\x04\xc8\xfe\x0d\xc9")
	// and rod of seasons.
	checkRod := r.appendToBank(0x3f, "check rod",
		"\x79\xfe\xe6\xc0\x7b\xfe\x02\xc9")
	// load gfx data for randomized shop and market items.
	itemGfxFunc := r.appendToBank(0x3f, "item gfx func",
		// check for matching object
		"\x43\x4f\xcd"+checkRod+"\x28\x17\x79\xfe\x59\x28\x19"+ // rod, woods
			"\xcd"+checkShopItem+"\x28\x1b\xcd"+
			checkMarketItem+"\x28\x1d"+ // shops
			"\x79\xfe\x6e\x28\x1f\x06\x00\xc9"+ // feather
			// look up item ID, subID
			"\x1e\x15\x21"+rodLookup+"\x18\x1d"+
			"\x1e\x0b\x21"+nobleSwordLookup+"\x18\x16"+
			"\x1e\x08\x21"+shopLookup+"\x18\x0f"+
			"\x1e\x09\x21"+marketLookup+"\x18\x08"+
			"\xfa\xb4\xc6\xc6\x15\x5f\x18\x0e"+ // feather
			"\xcd\x8a\x00\x79\x4b\xcd"+getTreasureData+ // get treasure
			"\xcd"+progressiveItemFunc+"\x23\x23\x5e"+ // get sprite
			"\x3e\x60\x4f\x06\x00\xc9") // replace object gfx w/ treasure gfx
	r.replace(0x3f, 0x443c, "item gfx call", "\x4f\x06\x00", "\xcd"+itemGfxFunc)

	// "activate" a flute by setting its icon and song when obtained. also
	// activates the corresponding animal companion.
	setFluteIcon := r.appendToBank(0x3f, "flute set icon func",
		"\xf5\xd5\xe5\x78\xfe\x0e\x20\x0d\x1e\xaf\x79\xd6\x0a\x12\xc6\x42"+
			"\x26\xc6\x6f\xcb\xfe\xe1\xd1\xf1\xcd\x4e\x45\xc9")
	r.replace(0x3f, 0x452c, "flute set icon call", "\x4e\x45", setFluteIcon)
}
