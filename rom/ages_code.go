package rom

import (
	"strings"
)

func newAgesRomBanks() *romBanks {
	r := romBanks{
		endOfBank: make([]uint16, 0x40),
	}

	r.endOfBank[0x00] = 0x3ef8
	r.endOfBank[0x01] = 0x7fc3
	r.endOfBank[0x02] = 0x7e93
	r.endOfBank[0x03] = 0x7ebd
	r.endOfBank[0x04] = 0x7edb
	r.endOfBank[0x05] = 0x7d9d
	r.endOfBank[0x06] = 0x7a31
	r.endOfBank[0x09] = 0x7dee
	r.endOfBank[0x0a] = 0x7e09
	r.endOfBank[0x0b] = 0x7fa8
	r.endOfBank[0x0c] = 0x7f94
	r.endOfBank[0x10] = 0x7ef4
	r.endOfBank[0x11] = 0x7f73
	r.endOfBank[0x12] = 0x7e8f
	r.endOfBank[0x15] = 0x7bfb
	r.endOfBank[0x16] = 0x7e03
	r.endOfBank[0x3f] = 0x7d0a

	return &r
}

func initAgesEOB() {
	r := newAgesRomBanks()

	// bank 00

	// don't play any music if the -nomusic flag is given.
	noMusicFunc := r.appendToBank(0x00, "no music func",
		"\x67\xfe\x47\x30\x03\x3e\x08\xc9\xf0\xb7\xc9")
	r.replace(0x00, 0x0c9a, "no music call",
		"\x67\xf0\xb7", "\xcd"+noMusicFunc)

	// only increment the maku tree's state if on the maku tree screen, or if
	// all essences are obtained, set it to the value it would normally have at
	// that point in the game. this allows getting the maku tree's item as long
	// as you haven't collected all essences.
	makuStateCheck := r.appendToBank(0x00, "maku state check",
		"\xfa\x2d\xcc\xfe\x02\x30\x0e\xfa\x30\xcc\xfe\x38\x20\x07"+
			"\xfa\xe8\xc6\x3c\xfe\x11\xc9\xfa\xbf\xc6\x3c\x37\x20\x03"+
			"\x3e\x0e\xc9\xfa\xe8\xc6\xc9")
	r.replace(0x00, 0x3e56, "call maku state check",
		"\x3c\xfe\x11", "\xcd"+makuStateCheck)

	// return z iff the current group and room match c and b.
	compareRoom := r.appendToBank(0x00, "compare room",
		"\xfa\x2d\xcc\xb9\xc0\xfa\x30\xcc\xb8\xc9")

	// read 2 bytes from bank e at hl into bc.
	readWord := r.appendToBank(0x00, "read word",
		"\xfa\x97\xff\xf5\x7b\xea\x97\xff\xea\x22\x22"+ // switch bank
			"\x2a\x47\x7e\x4f"+ // read
			"\xf1\xea\x97\xff\xea\x22\x22\xc9") // switch back

	// searches for a value in a table starting at hl, with an entry matching
	// keys b and subkey c, and values e bytes long. sets c if found. a key of
	// ff ends the table.
	searchDoubleKey := r.appendToBank(0x00, "search by double key",
		"\x2a\xfe\xff\xc8\xb8\x20\x06\x2a\xb9\x20\x03\x37\xc9"+
			"\x23\x7b\xd7\x18\xee")

	// bank 01

	// use a different invalid tile table for time warping if link doesn't have
	// flippers.
	noFlippersTable := r.appendToBank(0x01, "no flippers table",
		"\xf3\x00\xfe\x00\xff\x00\xe4\x00\xe5\x00\xe6\x00\xe7\x00\xe8\x00"+
			"\xe9\x00\xfc\x01\xfa\x00\xe0\x00\xe1\x00\xe2\x00\xe3\x00\x00")
	dontDrownLink := r.appendToBank(0x01, "don't drown link",
		"\x21\x17\x63\xfa\x9f\xc6\xe6\x40\xc0\x21"+noFlippersTable+"\xc9")
	r.replace(0x01, 0x6301, "call don't drown link",
		"\x21\x17\x63", "\xcd"+dontDrownLink)

	// bank 02

	// warp to ember tree if holding start when closing the map screen.
	treeWarp := r.appendToBank(0x02, "tree warp",
		"\xfa\x81\xc4\xe6\x08\x28\x1b"+ // close as normal if start not held
			"\xfa\x2d\xcc\xfe\x02\x38\x06"+ // check if indoors
			"\x3e\x5a\xcd\x98\x0c\xc9"+ // play error sound and ret
			"\x21\x47\xcc\x36\x80\x23\x36\x78\x2e\x4a\x36\x55"+ // set dest
			"\xcd\xbf\x5f\xc3\xba\x4f") // close + warp
	r.replaceMultiple([]Addr{{0x02, 0x6133}, {0x02, 0x618b}}, "tree warp jump",
		"\xc2\xba\x4f", "\xc4"+treeWarp)

	// warp to room under cursor if wearing developer ring.
	devWarp := r.appendToBank(0x02, "dev ring warp func",
		"\xfa\xcb\xc6\xfe\x40\x20\x12\xfa\x2d\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x47\xcc\xfa\xb6\xcb\xea\x48\xcc\x3e\x03\xcd\xad\x0c\xc9")
	r.replace(0x02, 0x5fcc, "dev ring warp call", "\xad\x0c", devWarp)

	// allow warping to south lynna tree even if it hasn't been visited (warp
	// menu locks otherwise).
	checkTreeVisited := r.appendToBank(0x02, "check tree visited",
		"\xfe\x78\xc2\x39\x66\xb7\xc9")
	r.replace(0x02, 0x5ff9, "call check tree visited 1",
		"\xcd\x39\x66", "\xcd"+checkTreeVisited)
	r.replace(0x02, 0x66a9, "call check tree visited 2",
		"\xcd\x39\x66", "\xcd"+checkTreeVisited)
	checkCursorVisited := r.appendToBank(0x02, "check cursor visited",
		"\xfa\xb6\xcb\xfe\x78\xc2\x39\x66\xb7\xc9")
	r.replace(0x02, 0x619d, "call check cursor visited",
		"\xcd\x36\x66", "\xcd"+checkCursorVisited)

	// bank 03

	// allow skipping the capcom screen after one second by pressing start
	skipCapcom := r.appendToBank(0x03, "skip capcom func",
		"\xe5\xfa\xb3\xcb\xfe\x94\x30\x03\xcd\x86\x08\xe1\xcd\x37\x02\xc9")
	r.replace(0x03, 0x4d6c, "skip capcom call", "\x37\x02", skipCapcom)

	// set flags to skip opening, abbreviate poe clock sequence, set
	// fairies' woods in order, and open the other half of the key door in d7.
	skipOpening := r.appendToBank(0x03, "skip opening",
		"\xcd\xf9\x31\x3e\x0a\xcd\xf9\x31\x3e\x23\xcd\xf9\x31"+ // global flags
			"\x3e\x2b\xcd\xf9\x31"+ // more
			"\xfa\xff\x7f\xea\x10\xc6"+ // set animal region
			"\xe5\x21\x7a\xc7\xcb\xf6\x2e\x6a\xcb\xf6"+ // set room flags
			"\x2e\x59\xcb\xf6\x2e\x39\x36\xc8\x2e\x7c\xcb\xf6"+ // more
			"\x2e\x2e\xcb\xf6\x21\x6d\xca\x36\x02\xe1\xc9") // more
	r.replace(0x03, 0x6e97, "call skip opening",
		"\xc3\xf9\x31", "\xc3"+skipOpening)

	// bank 04

	// look up tiles in custom replacement table after loading a room. the
	// format is (group, room, bitmask, YX, tile ID), with ff ending the table.
	// if the bitmask AND the current room flags is nonzero, the replacement is
	// not made.
	tileReplaceTable := r.appendToBank(0x04, "tile replace table",
		"\x01\x48\x00\x45\xd7"+ // portal south of past maku tree
			"\x00\x39\x00\x63\xf0"+ // open chest on intro screen
			"\x00\x39\x20\x63\xf1"+ // closed chest on intro screen
			"\x00\x6b\x00\x42\x3a"+ // removed tree in yoll graveyard
			"\x00\x6b\x02\x42\xce"+ // not removed tree in yoll graveyard
			"\x00\x83\x00\x43\xa4"+ // rock outside D2
			"\x03\x0f\x00\x66\xf9"+ // water in d6 past entrance
			"\x04\x1b\x01\x03\x78"+ // key door in D1
			"\x01\x13\x00\x61\xd7"+ // portal in symmetry city past
			"\x04\xbd\x00\x50\xa0"+ // remove key door in D5
			"\x04\xa6\x80\x54\x1e"+ // add key block in room below
			"\x05\x6e\x00\x50\xa0"+ // remove blue button key door in D7
			"\x05\x5a\x80\x66\x1e"+ // add key block on floor below
			"\x00\x25\x00\x37\xd7"+ // portal in nuun highlands
			"\xff")
	tileReplaceFunc := r.appendToBank(0x04, "tile replace body",
		"\xc5\xd5\xcd\x7d\x19\x5f\x21"+tileReplaceTable+"\xfa\x2d\xcc\x47"+
			"\xfa\x30\xcc\x4f"+ // load room flags, table addr, group, room
			"\x2a\xfe\xff\x28\x1b\xb8\x20\x12\x2a\xb9\x20\x0f"+
			"\x2a\xa3\x20\x0c"+ // compare group, room, flags
			"\xd5\x16\xcf\x2a\x5f\x2a\x12\xd1\x18\xe6"+ // replace
			"\x23\x23\x23\x23\x18\xe0\xd1\xc1\xcd\xef\x5f\xc9")
	r.replace(0x00, 0x38c0, "tile replace call",
		"\xcd\xef\x5f", "\xcd"+tileReplaceFunc)

	// treat the d2 present entrance like the d2 past entrance.
	replaceWarpEnter := r.appendToBank(0x04, "replace warp enter",
		"\xc5\x01\x00\x83\xcd"+compareRoom+"\xc1\xfa\x2d\xcc\xc0\x3c\xc9")
	r.replace(0x04, 0x4630, "call replace warp enter",
		"\xfa\x2d\xcc", "\xcd"+replaceWarpEnter)
	// and exit into the present if there's a portal in the present entrance.
	replaceWarpExit := r.appendToBank(0x00, "replace warp exit",
		"\xea\x48\xcc\xfe\x83\xc0\xfa\x2d\xcc\xfe\x04\xc0"+
			"\xfa\x3e\xc6\xb7\xc0\xfa\x3f\xc6\xfe\x83\xc0"+
			"\xfa\x47\xcc\xe6\xf0\xea\x47\xcc\xc9")
	r.replace(0x04, 0x45e8, "call replace warp exit normal",
		"\xea\x48\xcc", "\xcd"+replaceWarpExit)
	r.replace(0x0a, 0x4738, "call replace warp exit essence",
		"\xea\x48\xcc", "\xcd"+replaceWarpExit)

	// bank 05

	// if wearing dev ring, jump over any tile like a ledge by pressing B with
	// no B item equipped.
	devJump := r.appendToBank(0x05, "dev jump",
		"\xf5\xfa\xcb\xc6\xfe\x40\x20\x13"+ // check ring
			"\xfa\x88\xc6\xb7\x20\x0d"+ // check B item
			"\xfa\x81\xc4\xe6\x02\x28\x06"+ // check input
			"\xf1\xfa\x09\xd0\x37\xc9\xf1\xc9") // jump over ledge
	graveJump := r.appendToBank(0x05, "grave jump",
		"\xf5\xfe\x84\x20\x0f"+ // tile
			"\xc5\x01\x00\x5b\xcd"+compareRoom+"\xc1\x20\x05"+ // room
			"\xf1\x3e\x10\x37\xc9\xf1\xc9") // jump
	cliffLookup := r.appendToBank(0x05, "cliff lookup",
		"\xcd"+devJump+"\xd8\xcd"+graveJump+"\xd8\xc3\x1f\x1e")
	r.replace(0x05, 0x6083, "call cliff lookup",
		"\xcd\x1f\x1e", "\xcd"+cliffLookup)

	// prevent link from surfacing from underwater without mermaid suit. this
	// is probably only relevant for the sea of no return.
	preventSurface := r.appendToBank(0x05, "prevent surface",
		"\xfa\x91\xcc\xb7\xc0\xfa\xa3\xc6\xe6\x04\xfe\x04\xc9")
	r.replace(0x05, 0x516c, "call prevent surface",
		"\xfa\x91\xcc\xb7", "\xcd"+preventSurface+"\x00")

	// bank 06

	// burning the first tree in yoll graveyard should set room flag 1 so that
	// it can be gone for good.
	removeYollTree := r.appendToBank(0x06, "remove yoll tree",
		"\xf5\xf0\x8f\xfe\x0c\x20\x0f"+
			"\xc5\x01\x00\x6b\xcd"+compareRoom+"\x20\x05"+
			"\x21\x6b\xc7\xcb\xce\xc1\xf1\x21\x26\xc6\xc9")
	r.replace(0x06, 0x47aa, "call remove yoll tree",
		"\x21\x26\xc6", "\xcd"+removeYollTree)

	// bank 16 (pt. 1)

	// upgraded item data (one byte for old ID, one for new ID two for address):
	progItemAddrs := r.appendToBank(0x16, "progressive item addrs",
		"\x05\x05\xea\x54"+ // noble sword
			"\x0a\x0a\x12\x55"+ // long switch
			"\x16\x16\x52\x55"+ // power glove
			"\x19\x19\x76\x55"+ // satchel upgrade
			"\x25\x26\xca\x53"+ // tune of currents
			"\x26\x27\xce\x53"+ // tune of ages
			"\x2e\x4a\x5a\x54"+ // mermaid suit
			"\xff")
	// given a treasure ID in b, make hl = the start of the upgraded treasure
	// data + 1, if the treasure needs to be upgraded, and returns the new
	// treasure ID in b.
	getUpgradedTreasure := r.appendToBank(0x16, "get upgraded treasure",
		"\x78\xcd\x48\x17\x78\xd0"+ // check obtained
			"\xfe\x25\x20\x09\x3e\x26\x5f\xcd\x48\x17\x30\x01\x43"+ // harp
			"\xe5\x21"+progItemAddrs+"\x2a\xfe\xff\x28\x13"+ // search
			"\xb8\x20\x06\x2a\x47\x2a\x5e\x18\x05\x23\x23\x23\x18\xed"+
			"\xe1\x63\x6f\x23\xc9\xe1\xc9") // done
	// load the address of a treasure's 4-byte data entry + 1 into hl, using b
	// as the ID and c as sub ID, accounting for progressive upgrades.
	getTreasureDataBody := r.appendToBank(0x16, "get treasure data body",
		"\x21\x32\x53\x78\x87\xd7\x78\x87\xd7\xcb\x7e\x28\x04"+
			"\x23\x2a\x66\x6f\x79\x87\x87\xd7\x23\xc3"+getUpgradedTreasure)
	// do the above and put the ID, param, and text in b, c, and e.
	getTreasureDataBCE := r.appendToBank(0x16, "get treasure data bc",
		"\xcd"+getTreasureDataBody+"\x4e\x23\x5e\xc9")
	getTreasureData := r.appendToBank(0x00, "get treasure data",
		"\x1e\x16\x21"+getTreasureDataBCE+"\xc3\x8a\x00")

	// bank 09

	// set treasure ID 07 (rod of seasons) when buying the 150 rupee shop item,
	// so that the shop can check this specific ID.
	shopSetFakeID := r.appendToBank(0x09, "shop set fake ID",
		"\xfe\x0d\x20\x05\x21\x9a\xc6\xcb\xfe\x21\xf7\x44\xc9")
	r.replace(0x09, 0x4418, "call shop set fake ID",
		"\x21\xf7\x44", "\xcd"+shopSetFakeID)

	// set treasure ID 08 (magnet gloves) when getting item from south shore
	// dirt pile.
	digSetFakeID := r.appendToBank(0x09, "dirt set fake ID",
		"\xc5\x01\x00\x98\xcd"+compareRoom+"\xc1\xc0\xe5\x21\x9b\xc6\xcb\xc6"+
			"\xe1\xc9")
	// set treasure ID 13 (slingshot) when getting first item from tingle.
	tingleSetFakeID := r.appendToBank(0x09, "tingle set fake ID",
		"\xc5\x01\x00\x79\xcd"+compareRoom+"\xc1\xc0\xe5\x21\x9c\xc6\xcb\xde"+
			"\xe1\xc9")
	// set treasure ID 1e (fool's ore) for symmetry city brother.
	brotherSetFakeID := r.appendToBank(0x09, "brother set fake ID",
		"\xc5\x01\x03\x6e\xcd"+compareRoom+"\x28\x04\x04\xcd"+compareRoom+
			"\xc1\xc0\xe5\x21\x9d\xc6\xcb\xf6\xe1\xc9")
	// set treasure ID 10 (nothing) for king zora.
	kingZoraSetFakeID := r.appendToBank(0x09, "king zora set fake ID",
		"\xc5\x01\x05\xab\xcd"+compareRoom+"\xc1\xc0\xe5\x21\x9c\xc6\xcb\xc6"+
			"\xe1\xc9")
	// set treasure ID 12 (nothing) for first goron dance, and 14 (nothing) for
	// the second. if you're in the present, it's always 12. if you're in the
	// past, it's 12 iff you don't have letter of introduction.
	goronDanceSetFakeID := r.appendToBank(0x09, "dance 1 set fake ID",
		"\xc5\x01\x02\xed\xcd"+compareRoom+"\xc1\x28\x12"+ // present
			"\xc5\x01\x02\xef\xcd"+compareRoom+"\xc1\xc0"+ // past
			"\x3e\x59\xcd\x48\x17\x3e\x10\x38\x02\x3e\x04"+
			"\xe5\x21\x9c\xc6\xb6\x77\xe1\xc9")
	// set flag for d6 past and present boss keys whether you get the key in
	// past or present.
	setD6BossKey := r.appendToBank(0x09, "set d6 boss key",
		"\x7b\xfe\x31\xc0\xfa\x39\xcc\xfe\x06\x28\x03\xfe\x0c\xc0"+
			"\xe5\x21\x82\xc6\xcb\xf6\x23\xcb\xe6\xe1\xc9")
	// refill all seeds when picking up a seed satchel.
	refillSeedSatchel := r.appendToBank(0x09, "refill seed satchel",
		"\x7b\xfe\x19\xc0"+
			"\xc5\xd5\xe5\x21\xb4\xc6\x34\xcd\x0c\x18\x35\xe1\xd1\xc1\xc9")
	// give 20 seeds when picking up the seed shooter.
	fillSeedShooter := r.appendToBank(0x09, "fill seed shooter",
		"\x7b\xfe\x0f\xc0\xc5\x3e\x20\x0e\x20\xcd\x1c\x17\xc1\xc9")
	// give flute the correct icon and make it functional from the start.
	activateFlute := r.appendToBank(0x09, "activate flute",
		"\x7b\xfe\x0e\xc0"+
			"\x79\xd6\x0a\xea\xb5\xc6\xe5\x26\xc6\xc6\x45\x6f\x36\xc3\xe1\xc9")
	// this function checks all the above conditions when collecting an item.
	handleGetItem := r.appendToBank(0x09, "handle get item",
		"\x5f\xcd"+digSetFakeID+"\xcd"+setD6BossKey+"\xcd"+refillSeedSatchel+
			"\xcd"+fillSeedShooter+"\xcd"+activateFlute+"\xcd"+tingleSetFakeID+
			"\xcd"+brotherSetFakeID+"\xcd"+kingZoraSetFakeID+
			"\xcd"+goronDanceSetFakeID+"\x7b\xc3\x1c\x17")
	r.replace(0x09, 0x4c4e, "call handle get item",
		"\xcd\x1c\x17", "\xcd"+handleGetItem)

	// don't set room's item flag if it's nayru's item on the maku tree screen,
	// since link still might not have taken the maku tree's item.
	makuTreeItemFlag := r.appendToBank(0x09, "maku tree item flag",
		"\xcd\x7d\x19\xc5\x01\x38\xc7\xcd\xd6\x01\xc1\x20\x06\xfa\x0d\xd0"+
			"\xfe\x50\xc8\xcb\xee\xc9")
	r.replace(0x09, 0x4c82, "call maku tree item flag",
		"\xcd\x7d\x19", "\xc3"+makuTreeItemFlag)

	// give correct ID and param for shop item.
	shopGiveTreasure := r.appendToBank(0x09, "shop give treasure",
		"\x47\x1a\xfe\x0d\x78\x20\x04\xcd"+getTreasureData+
			"\x78\xc3"+handleGetItem)
	r.replace(0x09, 0x4425, "call shop give treasure",
		"\xcd\x1c\x17", "\xcd"+shopGiveTreasure)
	// and display correct text.
	shopShowText := r.appendToBank(0x09, "shop show text",
		"\x1a\xfe\x0d\x79\x20\x0c\x21\x11\x45\x46\x23\x4e\xcd"+getTreasureData+
			"\x06\x00\x4b\xc3\x72\x18")
	r.replace(0x09, 0x4443, "call shop show text",
		"\xc2\x72\x18", "\xc2"+shopShowText)

	// bank 0a

	// make ricky appear if you have his gloves, without giving rafton rope.
	checkRickyAppear := r.appendToBank(0x0a, "check ricky appear",
		"\xcd\xf3\x31\xc0\xfa\xa3\xc6\xcb\x47\xc0\xfa\x46\xc6\xb7\xc9")
	r.replace(0x0a, 0x4bb8, "call check ricky appear",
		"\xcd\xf3\x31", "\xcd"+checkRickyAppear)

	// require giving rafton rope, even if you have the island chart.
	checkRaftonRope := r.appendToBank(0x0a, "check rafton rope",
		"\xcd\x48\x17\xd0\x3e\x15\xcd\xf3\x31\xc8\x37\xc9")
	r.replace(0x0a, 0x4d5f, "call check rafton rope",
		"\xcd\x48\x17", "\xcd"+checkRaftonRope)

	// set sub ID for south shore dig item.
	dirtSpawnItem := r.appendToBank(0x0a, "dirt spawn item",
		"\xcd\xd4\x27\xc0\xcd\x42\x22\xaf\xc9")
	r.replace(0x0a, 0x5e3e, "call dirt spawn item",
		"\xcd\xc5\x24", "\xcd"+dirtSpawnItem)

	// automatically save maku tree when saving nayru.
	saveMakuTreeWithNayru := r.appendToBank(0x0a, "save maku tree with nayru",
		"\xcd\xf9\x31\xfa\xe8\xc6\xfe\x0e\x28\x02\x3e\x02\x3d\xea\xe8\xc6"+
			"\x3e\x0c\xcd\xf9\x31\x3e\x12\xcd\xf9\x31\x3e\x3f\xcd\xf9\x31"+
			"\xe5\x21\x38\xc7\xcb\x86\x24\xcb\xfe\x2e\x48\xcb\xc6\xe1\xc9")
	r.replace(0x0a, 0x5541, "call save maku tree with nayru",
		"\xcd\xf9\x31", "\xcd"+saveMakuTreeWithNayru)

	// set link's simulated input differently if entering the first maku treen
	// cutscene from the right. this prevents him from being forced to move
	// offscreen, where the script never restores player control.
	makuTreeInput := r.appendToBank(0x0a, "maku tree input",
		"\x2d\x00\x00\x10\x00\x20\x04\x00\x00\x20\x00\x80\x04\x00\x00"+
			"\x30\x00\x20\x04\x00\x00\x04\x00\x40\x2d\x00\x00\xff\xff")
	setMakuTreeInput := r.appendToBank(0x0a, "set maku tree input",
		"\x5f\xfa\x02\xcd\xfe\x03\x7b\x20\x03\x21"+makuTreeInput+
			"\xcd\x1d\x2a\xc9")
	r.replace(0x0a, 0x66e0, "call set maku tree input",
		"\xcd\x1d\x2a", "\xcd"+setMakuTreeInput)

	// bank 0b

	// always get item from king zora before permission to enter jabu-jabu.
	kingZoraCheck := r.appendToBank(0x0b, "king zora check",
		"\xcd\xf3\x31\xc8\x3e\x10\xcd\x48\x17\x3e\x00\xd0\x3c\xc9")
	r.replace(0x0b, 0x5464, "call king zora check",
		"\xcd\xf3\x31", "\xcd"+kingZoraCheck)

	// bank 0c

	// use custom script for soldier in deku forest with sub ID 0; they should
	// give an item in exchange for mystery seeds.
	soldierScriptAfter := r.appendToBank(0x0c, "soldier script after item",
		"\x97\x59\x08\x00")
	soldierScriptGive := r.appendToBank(0x0c, "soldier script give item",
		"\xeb\x9e\x98\x59\x0b\xb4\xbd\x00\x92\xe9\xcb\x02\xde\x00\x00\xb1\x20"+
			"\xc4"+soldierScriptAfter)
	soldierScriptCheck := r.appendToBank(0x0c, "soldier script check count",
		"\xb3\xbd\xff"+soldierScriptGive+"\x5d\xee")
	soldierScript := r.appendToBank(0x0c, "soldier script",
		"\xb0\x20"+soldierScriptAfter+"\xdf\x24"+soldierScriptCheck+"\x5d\xee")
	r.replace(0x09, 0x5207, "soldier script pointer", "\xee\x5d", soldierScript)

	// bank 10

	// keep black tower in initial state until the player got the item from the
	// worker.
	blackTowerCheck := r.appendToBank(0x10, "black tower check",
		"\x21\x27\x79\xc8\xfa\xe1\xc9\xe6\x20\xc9")
	r.replace(0x10, 0x7914, "call black tower check",
		"\x21\x27\x79", "\xcd"+blackTowerCheck)

	// bank 11

	// allow collection of seeds with only shooter and no satchel
	checkSeedHarvest := r.appendToBank(0x11, "check seed harvest",
		"\xcd\x48\x17\xd8\x3e\x0f\xc3\x48\x17")
	r.replace(0x11, 0x4aba, "call check seed harvest",
		"\xcd\x48\x17", "\xcd"+checkSeedHarvest)

	// bank 12

	// add time portal interaction in symmetry city past, to avoid softlock if
	// player only has echoes.
	symmetryPastPortal := r.appendToBank(0x12, "symmetry past portal",
		"\xf1\xdc\x05\xf2\xe1\x00\x68\x18\xfe")
	r.replace(0x12, 0x5e91, "symmetry past portal pointer",
		"\xf1\xdc\x05", "\xf3"+symmetryPastPortal)
	// add one to nuun highlands too.
	nuunPortalOtherObjects := r.appendToBank(0x12, "nuun portal other objects",
		"\xf2\x9a\x00\x68\x48\x9a\x01\x58\x58\x9a\x02\x58\x68\x9a\x03\x48\x58"+
			"\x9a\x04\x38\x58\xfe")
	r.replace(0x12, 0x5a7b, "nuun portal", "\xf2\x9a\x00",
		"\xf2\xe1\x00\x38\x78\xf3"+nuunPortalOtherObjects+"\xff")

	// bank 15

	// don't equip sword for shooting galleries if player don't have it
	// (doesn't work anyway).
	shootingGalleryEquip := r.appendToBank(0x15, "shooting gallery equip",
		"\x3e\x05\xcd\x48\x17\x3e\x00\x22\xd0\x2b\x3e\x05\x22\xc9")
	r.replace(0x15, 0x50ae, "call shooting gallery equip",
		"\x3e\x05\x22", "\xcd"+shootingGalleryEquip)

	// always make "boomerang" second prize for target carts, checking room
	// flag 6 to track it.
	targetCartsItem := r.appendToBank(0x15, "target carts item",
		"\xcd\x7d\x19\xcb\x77\x3e\x04\xca\xbb\x66\xcd\x3e\x04\xc3\xa5\x66")
	r.replace(0x15, 0x66a2, "call target carts item",
		"\xcd\x3e\x04", "\xc3"+targetCartsItem)
	// set room flag 6 when "boomerang" is given in script.
	targetCartsFlag := r.appendToBank(0x0c, "target carts flag",
		"\xde\x06\x02\xb1\x40\xc1")
	r.replace(0x0c, 0x6e6e, "jump target carts flag",
		"\x88\x6e", targetCartsFlag)

	// bank 16

	// given a treasure ID in dx42, return hl = the start of the treasure data
	// + 1, accounting for progressive upgrades. also writes the new treasure
	// ID to d070, which is used to set the treasure obtained flag.
	upgradeTreasure := r.appendToBank(0x16, "upgrade treasure",
		"\x1e\x42\x1a\x47\xcd"+getUpgradedTreasure+"\x1e\x70\x78\x12\xc9")

	// just get item bc's sprite index in e.
	getItemSpriteIndexBody := r.appendToBank(0x16, "get item sprite index body",
		"\xcd"+getTreasureDataBody+"\x23\x23\x7e\x5f\xc9")
	getItemSpriteIndex := r.appendToBank(0x00, "get item sprite index",
		"\x1e\x16\x21"+getItemSpriteIndexBody+"\xc3\x8a\x00")

	// return collection mode in a and e, based on current room. call is in
	// bank 16, func is in bank 00, body is in bank 06.
	collectModeTable := r.appendToBank(0x06, "collect mode table",
		makeAgesCollectModeTable())
	// maku tree item falls or exists on floor depending on script position.
	collectMakuTreeFunc := r.appendToBank(0x06, "collect maku tree",
		"\xfa\x58\xd2\xfe\x84\x1e\x29\xc8\x1e\x0a\xc9")
	// target carts items appear with a poof if they're in the enclosure.
	collectTargetCartsFunc := r.appendToBank(0x06, "collect target carts",
		"\x1e\x4d\x1a\xfe\x78\x1e\x19\xc8\x1e\x0a\xc9")
	// big bang game items appear with a poof if they're above the goron.
	collectBigBangFunc := r.appendToBank(0x06, "collect big bang game",
		"\x1e\x4b\x1a\xfe\x38\x1e\x19\xc8\x1e\x0a\xc9")
	// lava juice trading goron also has a chest in the room.
	collectLavaJuiceFunc := r.appendToBank(0x06, "collect lava juice room",
		"\x1e\x4d\x1a\xfe\x68\x1e\x0a\xd8\x1e\x38\xc9")
	collectModeJumpTable := r.appendToBank(0x06, "collect mode jump table",
		collectMakuTreeFunc+collectTargetCartsFunc+collectBigBangFunc+
			collectLavaJuiceFunc)
	collectModeLookupBody := r.appendToBank(0x06, "collect mode lookup body",
		"\xfa\x2d\xcc\x47\xfa\x30\xcc\x4f\x1e\x01\x21"+collectModeTable+
			"\xcd"+searchDoubleKey+"\x5f\xd0\x7e\x5f\xfe\x80\xd8"+
			"\x21"+collectModeJumpTable+"\xe6\x7f\x87\xd7\x2a\x66\x6f\xe9")
	collectModeLookup := r.appendToBank(0x00, "collect mode lookup",
		"\xc5\xd5\xe5\x1e\x06\x21"+collectModeLookupBody+"\xcd\x8a\x00\x7b"+
			"\xe1\xfe\xff\x20\x02\x2b\x2a\xd1\xc1\xc9")
	// return treasure data address and collect mode modified as necessary,
	// given a treasure ID in dx42.
	modifyTreasure := r.appendToBank(0x16, "modify treasure",
		"\xcd"+upgradeTreasure+"\xcd"+collectModeLookup+"\x47\xcb\x37\xc9")
	r.replace(0x16, 0x4539, "call modify treasure",
		"\x47\xcb\x37", "\xcd"+modifyTreasure)

	// bank 3f

	// set hl to the address of the item sprite with ID a.
	calcItemSpriteAddr := r.appendToBank(0x3f, "get item sprite addr",
		"\x21\xdb\x66\x5f\x87\xd7\x7b\xd7\xc9")
	// set hl to the address of the item sprite for the item at hl in bank e.
	lookupItemSpriteAddr := r.appendToBank(0x3f, "look up item sprite addr",
		"\xc5\xcd"+readWord+"\xcd"+getItemSpriteIndex+
			"\x7b\xcd"+calcItemSpriteAddr+"\xc1\xc9")

	// copy three bytes at hl to a temporary ram address, and set hl to the
	// address of the last byte, with a as the value.
	copySpriteData := r.appendToBank(0x3f, "copy sprite data",
		"\xd5\x11\xf0\xcf\x2a\x12\x13\x2a\x12\x13\x7e\x12"+
			"\x62\x6b\xd1\xc9")

	// make the deku forest soldier that gives the item red instead of blue.
	soldierSprite := r.appendToBank(0x3f, "soldier sprite", "\x4d\x00\x22")
	setSoldierSprite := r.appendToBank(0x3f, "set soldier sprite",
		"\x21"+soldierSprite+"\xf1\xc9")
	// these interactions use the same flags as regular items
	setShopItemSprite := r.appendToBank(0x3f, "set shop item sprite",
		"\x1e\x09\x21\x11\x45\xcd"+lookupItemSpriteAddr+"\xf1\xc9")
	setHiddenTokaySprite := r.appendToBank(0x3f, "set hidden tokay sprite",
		"\x1e\x15\x21\x36\x5b\xcd"+lookupItemSpriteAddr+"\xf1\xc9")
	setWildTokaySprite := r.appendToBank(0x3f, "set wild tokay sprite",
		"\x1e\x15\x21\xbb\x5b\xcd"+lookupItemSpriteAddr+"\xf1\xc9")
	// interaction 6b, can't handle bomb flower and needs different flags
	set6BSprite := r.appendToBank(0x3f, "set interaction 6b sprite",
		"\xcd"+lookupItemSpriteAddr+"\xcd"+copySpriteData+
			"\xcb\x46\x20\x03\x34\x18\x01\x35\x2b\x2b\xf1\xc9")
	setInventionSprite := r.appendToBank(0x3f, "set invention sprite",
		"\x1e\x0c\x21\x32\x72\xc3"+set6BSprite)
	setChevalTestSprite := r.appendToBank(0x3f, "set cheval's test sprite",
		"\x1e\x0c\x21\x3b\x72\xc3"+set6BSprite)
	// interaction 80, can't handle bomb flower and needs different flags
	set80Sprite := r.appendToBank(0x3f, "set interaction 80 sprite",
		"\xcd"+lookupItemSpriteAddr+"\xcd"+copySpriteData+
			"\x5f\xe6\x0f\x20\x05\x7b\xc6\x03\x18\x06"+
			"\xfe\x02\x7b\x28\x01\x3c\x77\x2b\x2b\xf1\xc9")
	setLibraryPastSprite := r.appendToBank(0x3f, "set library past sprite",
		"\x1e\x15\x21\xd8\x5d\xc3"+set80Sprite)
	setLibrarySprite := r.appendToBank(0x3f, "set library sprite",
		"\x1e\x15\x21\xb9\x5d\xc3"+set80Sprite)
	// table of ID, sub ID, jump address
	customSpriteTable := r.appendToBank(0x3f, "custom sprite table",
		"\x40\x00"+setSoldierSprite+
			"\x47\x0d"+setShopItemSprite+ // 150 rupees only
			"\x63\x14"+setHiddenTokaySprite+ // iron shield
			"\x63\x15"+setHiddenTokaySprite+ // mirror shield
			"\x63\x3e"+setWildTokaySprite+ // wild tokay game prize
			"\x6b\x0b"+setInventionSprite+ // cheval's invention
			"\x6b\x0c"+setChevalTestSprite+
			"\x80\x07"+setLibraryPastSprite+
			"\x80\x08"+setLibrarySprite+
			"\xff")
	// override the sprites loaded for certain ID / sub ID pairs.
	loadCustomSprite := r.appendToBank(0x3f, "load custom sprite",
		"\xcd\x37\x44\xf5\xc5\xe5\x1e\x41\x1a\x47\x1c\x1a\x4f"+
			"\x1e\x02\x21"+customSpriteTable+"\xcd"+searchDoubleKey+
			"\x30\x08\x2a\x47\x7e\xe1\x67\x68\xc1\xe9\xe1\xc1\xf1\xc9")
	r.replace(0x3f, 0x4356, "call load custom sprite",
		"\xcd\x37\x44", "\xcd"+loadCustomSprite)
}

// makes ages-specific additions to the collection mode table.
func makeAgesCollectModeTable() string {
	b := new(strings.Builder)
	table := makeCollectModeTable()
	b.WriteString(table[:len(table)-1]) // strip final ff

	// add eatern symmetry city brother
	b.Write([]byte{0x03, 0x6f, collectFind2})

	// add ricky and dimitri nuun caves
	b.Write([]byte{0x02, 0xec, collectChest, 0x05, 0xb8, collectChest})

	b.Write([]byte{0xff})
	return b.String()
}
