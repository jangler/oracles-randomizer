package rom

import (
	"strings"
)

func newSeasonsRomBanks() *romBanks {
	r := romBanks{
		endOfBank: make([]uint16, 0x40),
	}

	r.endOfBank[0x00] = 0x3ec8
	r.endOfBank[0x01] = 0x7e89
	r.endOfBank[0x02] = 0x75bb
	r.endOfBank[0x03] = 0x7dd7
	r.endOfBank[0x04] = 0x7e02
	r.endOfBank[0x05] = 0x7e2d
	r.endOfBank[0x06] = 0x77d4
	r.endOfBank[0x07] = 0x78f0
	r.endOfBank[0x08] = 0x7fc0
	r.endOfBank[0x09] = 0x7f4e
	r.endOfBank[0x0a] = 0x7bea
	r.endOfBank[0x0b] = 0x7f6d
	r.endOfBank[0x11] = 0x7eb0
	r.endOfBank[0x15] = 0x792d
	r.endOfBank[0x3f] = 0x714d

	return &r
}

// for some reason the maku tree has a different room for every number of
// essences you have.
var (
	makuTreeRooms = []byte{0x0b, 0x0c, 0x7b, 0x2b, 0x2c, 0x2d, 0x5b, 0x5c, 0x5d}
	starOreRooms  = []byte{0x66, 0x76, 0x75, 0x65}
)

func initSeasonsEOB() {
	r := newSeasonsRomBanks()

	// try to order these first by bank, then by call location. maybe group
	// them into subfunctions when applicable?

	// bank 00

	// don't play any music if the -nomusic flag is given.
	noMusicFunc := r.appendToBank(0x00, "no music func",
		"\x67\xfe\x49\x30\x03\x3e\x08\xc9\xf0\xb5\xc9")
	r.replace(0x00, 0x0c76, "no music call",
		"\x67\xf0\xb5", "\xcd"+noMusicFunc)

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
	// treasure data
	progData := r.endOfBank[0x00]
	r.appendToBank(0x00, "progressive item data",
		"\x03\x21\x15"+ // mirror shield
			"\x02\x1d\x11"+ // noble sword
			"\x03\x1e\x12"+ // master sword
			"\x02\x23\x1d"+ // boomerang
			"\x02\x2f\x22"+ // slingshot
			"\x02\x28\x17"+ // feather
			"\x00\x46\x20") // satchel
	// change hl to point to different treasure data if the item is progressive
	// and needs to be upgraded. param a = treasure ID.
	progressiveItemFunc := r.appendToBank(0x00, "progressive item func",
		"\xd5\x5f\xcd"+upgradeFeather+"\x7b\xd1\xd0"+ // ret if missing L-1
			"\xfe\x01\x20\x0b\xfa\xa9\xc6\xfe\x02\x20\x04"+ // shield
			"\x21"+addrString(progData)+"\xc9"+ // to mirror shield
			"\xfe\x05\x20\x0f\xfa\xac\xc6\xfe\x02\x28\x04"+ // sword
			"\x21"+addrString(progData+3)+"\xc9"+ // to noble sword
			"\x21"+addrString(progData+6)+"\xc9"+ // to master sword
			"\xfe\x06\x20\x04\x21"+addrString(progData+9)+"\xc9"+ // boomerang
			"\xfe\x13\x20\x04\x21"+addrString(progData+12)+"\xc9"+ // slingshot
			"\xfe\x17\x20\x04\x21"+addrString(progData+15)+"\xc9"+ // feather
			"\xfe\x19\xc0\x21"+addrString(progData+18)+"\xc9") // satchel

	// this is a replacement for giveTreasure that gives treasure, plays sound,
	// and sets text based on item ID a and sub ID c, and accounting for item
	// progression.
	giveItem := r.appendToBank(0x00, "give item func",
		"\xcd"+getTreasureData+"\x47\xcd"+progressiveItemFunc+"\x78"+ // load data
			"\x4e\xcd\xeb\x16\x28\x05\xe5\xcd\x74\x0c\xe1"+ // give, play sound
			"\x06\x00\x23\x4e\x79\xfe\xff\xc8\xcd\x4b\x18\xaf\xc9") // show text

	// utility function, call a function hl in bank 02, preserving af. e can't
	// be used as a parameter to that function, but it can be returned.
	callBank2 := r.appendToBank(0x00, "call bank 02",
		"\xf5\x1e\x02\xcd\x8a\x00\xf1\xc9")

	// increment hl until (hl) equals either register a or ff. returns z if a
	// match was found.
	searchValue := r.appendToBank(0x00, "search value",
		"\xc5\x47\x2a\xb8\x28\x06\x3c\x28\x02\x18\xf7\x3c\x78\xc1\xc9")

	// search for a key bc in a dictionary starting at hl. the dictionary is a
	// series of two-byte keys and two-byte values. if a match is found, hl is
	// the address of the value, and z is set. the dictionary ends when $ff is
	// encountered at the beginning of an entry.
	lookupWord := r.appendToBank(0x00, "lookup word",
		"\x2a\xfe\xff\x20\x02\xb7\xc9"+ // ret if key is $ff
			"\xb8\x2a\x20\x02\xb9\xc8"+ // compare key bytes
			"\x23\x23\x18\xef") // loop

	// bank 01

	// helper function, takes b = high byte of season addr, returns season in b
	readSeason := r.appendToBank(0x01, "read default season",
		"\x26\x7e\x68\x7e\x47\xc9")

	// bank 02

	// warp to ember tree if holding start when closing the map screen.
	treeWarp := r.appendToBank(0x02, "tree warp",
		"\xfa\x81\xc4\xe6\x08\x28\x16"+ // close as normal if start not held
			"\xfa\x50\xcc\xe6\x01\x20\x06"+ // check if indoors
			"\x3e\x5a\xcd\x74\x0c\xc9"+ // play error sound and ret
			"\x21\xb7\xcb\x36\x05\xaf\xcd\x7b\x5e\xc3\x7b\x4f") // close + warp
	r.replaceMultiple([]Addr{{0x02, 0x6089}, {0x02, 0x602c}}, "tree warp jump",
		"\xc2\x7b\x4f", "\xc4"+treeWarp)

	// warp to room under cursor if wearing developer ring.
	devWarp := r.appendToBank(0x02, "dev ring warp func",
		"\xfa\xc5\xc6\xfe\x40\x20\x12\xfa\x49\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x63\xcc\xfa\xb6\xcb\xea\x64\xcc\x3e\x03\xcd\x89\x0c\xc9")
	r.replace(0x02, 0x5e9b, "dev ring warp call", "\x89\x0c", devWarp)

	// load a custom room layout for the problematic woods of winter screen in
	// winter. the code here is one 8-tile compression block per line.
	winterLayout := r.appendToBank(0x02, "winter layout",
		"\x55\x80\x81\x81\x81\x81"+
			"\x7c\x16\x80\x82\x17"+
			"\xf0\x1b\xc4\xc4\x70\x72"+
			"\x00\x01\x0d\x17\xc4\x80\x81\x70\x71"+
			"\x60\x04\x70\x71\x1a\x1b\x1c\xf7"+
			"\x05\x80\x81\x81\x70\x71\x9e\x9e"+
			"\x1c\x16\x04\x15\x17\x80\x81"+
			"\x30\x1b\x99\x9b\xd9\x1a\x01\x19"+
			"\x00\x70\x71\x15\x16\x17\xf7\x7a\x8c"+
			"\x11\x18\x19\x80\x81\x01\x19\x70")
	loadWinterLayout := r.appendToBank(0x00, "load winter layout",
		"\xd5\xfa\x4c\xcc\xfe\x9d\x20\x14\xfa\x4e\xcc\xfe\x03\x20\x0d"+
			"\xfa\x49\xcc\xb7\x20\x07\x3e\x02\xe0\x8c\x21"+winterLayout+
			"\xf0\x8c\xc3\xe2\x39")
	r.replace(0x00, 0x39df, "jump to winter layout",
		"\xd5\xf0\x8c", "\xc3"+loadWinterLayout)

	// allow ring list to be accessed through the ring box icon
	ringListOpener := r.appendToBank(0x02, "ring list opener",
		"\xfa\xd1\xcb\xfe\x0f\xc0\x3e\x81\xea\xd3\xcb\x3e\x04\xcd\x76\x1a\xe1\xc9")
	r.replace(0x02, 0x56a1, "call ring list opener",
		"\xfa\xd1\xcb", "\xcd"+ringListOpener)

	// auto-equip selected ring from ring list
	autoEquipRing := r.appendToBank(0x02, "auto-equip ring",
		"\xcd\x6c\x71\xea\xc5\xc6\xc9")
	r.replace(0x02, 0x6f4a, "call auto-equip ring",
		"\xcd\x6c\x71", "\xcd"+autoEquipRing)

	// don't save gfx when opening ring list from subscreen (they were already saved when
	// opening the item menu), and clear screen scroll variables (which are saved anyway)
	ringListGfxFix := r.appendToBank(0x02, "ring list gfx fix",
		"\xcd\x89\x0c\xfa\xd3\xcb\xcb\x7f\xc8\xe6\x7f\xea\xd3\xcb"+
			"\xaf\xe0\xa8\xe0\xaa\x21\x08\xcd\x22\x22\xc3\x72\x50")
	r.replace(0x02, 0x5035, "call ring list gfx fix",
		"\xcd\x89\x0c", "\xcd"+ringListGfxFix)

	// allow warping to horon village tree even if it hasn't been visited (warp
	// menu locks otherwise).
	checkTreeVisited := r.appendToBank(0x02, "check tree visited",
		"\xfe\xf8\xc2\x60\x65\xb7\xc9")
	r.replace(0x02, 0x5ec8, "call check tree visited 1",
		"\xcd\x60\x65", "\xcd"+checkTreeVisited)
	r.replace(0x02, 0x65e1, "call check tree visited 2",
		"\xcd\x60\x65", "\xcd"+checkTreeVisited)
	checkCursorVisited := r.appendToBank(0x02, "check cursor visited",
		"\xfa\xb6\xcb\xc3"+checkTreeVisited)
	r.replace(0x02, 0x609b, "call check cursor visited",
		"\xcd\x5d\x65", "\xcd"+checkCursorVisited)

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

	// bank 05

	// do this so that animals don't immediately stop walking on screen when
	// called on a bridge.
	fluteEnterFunc := r.appendToBank(0x05, "flute enter func",
		"\xcd\xaa\x44\xb7\xc8\xfe\x1a\xc8\xfe\x1b\xc9")
	r.replaceMultiple([]Addr{{0x05, 0x71ea}, {0x05, 0x493b}},
		"animal enter call", "\xcd\xaa\x44\xb7", "\xcd"+fluteEnterFunc+"\x00")

	// let link jump down the cliff outside d7, in case of winter sans shovel.
	// also let link jump down the snow cliff added in woods of winter. also
	// lets link jump over any tile if wearing dev ring while shielding.
	cliffLookupFunc := r.appendToBank(0x05, "cliff lookup func",
		"\xf5\xfa\xc5\xc6\xfe\x40\x20\x0c\xfa\x89\xcc\xb7\x28\x06"+ // dev
			"\xf1\xfa\x09\xd0\x37\xc9"+ // always jump if dev ring + shield
			"\xfa\x49\xcc\xb7\x20\x21"+ // cp group
			"\xfa\x4c\xcc\xfe\xd0\x20\x09\xf1"+ // d7 entrance
			"\xfe\xa8\x20\x16\x3e\x08\x37\xc9"+ // cp tile
			"\xfe\x9d\x20\x0d\xf1"+ // woods of winter
			"\xfe\x99\x28\x04\xfe\x9b\x20\x05\x3e\x10\x37\xc9"+ // cp tile
			"\xf1\xc3\xdd\x1d") // jp to normal lookup
	r.replace(0x05, 0x5fe8, "cliff lookup call",
		"\xcd\xdd\x1d", "\xcd"+cliffLookupFunc)

	// make moosh unrideable on mt cucco in the case of not having flute in a
	// moosh seed.
	checkMooshFlute := r.appendToBank(0x05, "check moosh flute",
		"\x3e\x0e\xc3\x17\x17")
	r.replace(0x05, 0x776b, "call check moosh flute 1",
		"\xfa\x10\xc6", "\xcd"+checkMooshFlute)
	r.replace(0x05, 0x7a65, "call check moosh flute 2",
		"\xfa\x10\xc6", "\xcd"+checkMooshFlute)

	// bank 06

	// replace a random item drop with gale seeds 1/4 of the time if the player
	// is out of gale seeds. this is important so that the one-way cliffs can
	// be in logic with gale seeds.
	galeDrop := r.appendToBank(0x06, "gale drop func",
		"\x3e\x23\xcd\x17\x17\xd0\x2e\xb8\xb6\xc0"+
			"\xcd\x1a\x04\xfe\x40\xd0\x0e\x08\xc9")
	galeDropWrapper := r.appendToBank(0x06, "gale drop wrapper",
		"\xcd"+galeDrop+"\xcd\xa7\x3e\xc9")
	r.replace(0x06, 0x47f5, "gale drop call",
		"\xcd\xa7\x3e", "\xcd"+galeDropWrapper)

	// Use expert's or fist ring with only one button unequipped
	r.replace(0x06, 0x490e, "punch with 1 button", "\xc0", "\x00")

	// bank 07

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
			"\xcd"+giveItem+"\xc9") // give item and ret
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
	// returns c if the player has gale seeds and the seed satchel. used for
	// warnings for cliffs and diving.
	checkGaleSatchel := r.appendToBank(0x15, "check gale satchel",
		"\xc5\x47\x3e\x19\xcd\x17\x17\x30\x05\x3e\x23\xcd\x17\x17\x78\xc1\xc9")
	warnGeneric := r.appendToBank(0x15, "warn generic",
		"\xcd\xc6\x3a\xc0\x36\x9f\x2e\x46\x36\x3c"+ // init object
			"\x01\x00\xf1\x11\x0b\xd0\xcd\x1a\x22"+ // set position
			"\x3e\x50\xcd\x74\x0c"+ // play sound
			"\x21\xc0\xcf\xcb\xc6\xc9") // set $cfc0 bit and ret
	warnCliff := r.appendToBank(0x15, "warn cliff",
		"\xaf\xea\xe0\xcf\xc3"+warnGeneric)
	warnFlowerCliff := r.appendToBank(0x15, "warn flower cliff",
		"\xcd"+checkGaleSatchel+"\xd8"+
			"\x06\x61\x16\x01\xcd"+warningHelper+"\xc8\xc3"+warnCliff)
	warnDivingSpot := r.appendToBank(0x15, "warn diving spot",
		"\x78\xfe\x03\xc8\xcd"+checkGaleSatchel+"\xd8"+
			"\x06\x61\x16\x09\xcd"+warningHelper+"\xc8\xc3"+warnCliff)
	warnWaterfallCliff := r.appendToBank(0x15, "warn waterfall cliff",
		"\xcd"+checkGaleSatchel+"\xd8"+
			"\x06\x65\x16\x02\xcd"+warningHelper+"\xc8\xc3"+warnCliff)
	warnMoblinKeep := r.appendToBank(0x15, "warn moblin keep",
		"\xcd"+checkGaleSatchel+"\xd8"+
			"\xfa\x10\xc6\xfe\x0c\xc0\x3e\x17\xcd\x17\x17\xd8\xc3"+warnCliff)
	warnHSSSkip := r.appendToBank(0x15, "warn hss skip",
		"\xfa\x86\xca\xb7\xc0\xcd\x56\x19\xcb\x76\xc0\xcb\xf6"+
			"\x3e\x02\xea\xe0\xcf\xc3"+warnGeneric)
	// this communicates with the warning script by setting bit zero of $cfc0
	// if the warning needs to be displayed (based on room, season, etc), and
	// also displays the exclamation mark if so.
	warningFunc := r.appendToBank(0x15, "warning func",
		"\xc5\xd5\xcd"+addrString(r.endOfBank[0x15]+8)+"\xd1\xc1\xc9"+ // wrap
			"\xfa\x4e\xcc\x47\xfa\xb0\xc6\x4f\xfa\x4c\xcc"+ // load env data
			"\xfe\x7c\xca"+warnFlowerCliff+
			"\xfe\x6e\xca"+warnDivingSpot+"\xfe\x3d\xca"+warnWaterfallCliff+
			"\xfe\x5c\xca"+warnMoblinKeep+"\xfe\x78\xca"+warnHSSSkip+
			"\xc3"+warnGeneric)
	warnCliffText := r.appendToBank(0x0b, "cliff warning script",
		"\x98\x26\x00\xbe\x00") // show cliff warning text
	warnBushText := r.appendToBank(0x0b, "bush warning script",
		"\x00") // impossible since 2.2.0
	warnSkipText := r.appendToBank(0x0b, "skip warning script",
		"\x98\x26\x02\xbe\x00") // show key skip warning text
	// point to this script instead of the normal maku gate script
	warningScript := r.appendToBank(0x0b, "warning script",
		"\xcb\x4c\xcc\xd9\x87\x4e"+ // use maku gate script if on that screen
			"\xd0\xe0"+warningFunc+"\xa0\xbd\xd7\x3c"+ // wait for collision
			"\x87\xe0\xcf"+warnCliffText+warnBushText+warnSkipText) // jp table
	r.replace(0x08, 0x5663, "warning script pointer", "\x87\x4e", warningScript)

	// set sub ID for star ore
	starOreIDFunc := r.appendToBank(0x08, "star ore id func",
		"\x2c\x36\x45\x2c\x36\x00\xc9")
	r.replace(0x08, 0x62f2, "star ore id call",
		"\x2c\x36\x45", "\xcd"+starOreIDFunc)

	// remove volcano cutscene.
	rmVolcano := r.appendToBank(0x02, "remove volcano scene",
		"\xcd\x56\x19\xcb\xf6\x11\x44\xd2\x3e\x02\x12\x21\x14\x63\xcd\xfe\x24"+
			"\x3e\x15\xc3\xcd\x30")
	r.replace(0x08, 0x7d07, "call remove volcano scene",
		"\xfa\x18\xcd\xb7\xc0\xcd\x56\x19\xcb\xf6\x3e\x0b\xea\x04\xcc\xcd",
		"\xfa\x44\xd2\xfe\x01\xc0\xcd\xd9\x3a\x21"+rmVolcano+"\xc3"+callBank2)
	r.replace(0x08, 0x7cf5, "enable volcano exit",
		"\xea\xab\xcc", "\x00\x00\x00")

	// remove generic "you got a ring" text for rings from shops
	r.replace(0x08, 0x4d55, "obtain ring text replacement (shop) 1", "\x54", "\x00")
	r.replace(0x08, 0x4d56, "obtain ring text replacement (shop) 2", "\x54", "\x00")

	// bank 09

	// shared by maku tree and star-shaped ore.
	starOreRoomTable := r.appendToBank(0x02, "star ore room table",
		string(starOreRooms)+"\xff")
	makuTreeRoomTable := r.appendToBank(0x02, "maku tree room table",
		string(makuTreeRooms)+"\xff")
	bank2IDFunc := r.appendToBank(0x02, "bank 2 fake id func",
		"\xfa\x49\xcc\xfe\x01\x28\x05\xfe\x02\x28\x11\xc9"+ // compare group
			"\xfa\x4c\xcc\x21"+starOreRoomTable+"\xcd"+searchValue+
			"\xc0\x21\x94\xc6\xcb\xd6\xc9"+
			"\xfa\x4c\xcc\x21"+makuTreeRoomTable+"\xcd"+searchValue+
			"\xc0\x21\x93\xc6\xcb\xd6\xc9")
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

	// check treasure id 0a to determine whether the maku tree gives its intro
	// speech and item, but return the number of essences in a.
	makuTreeCheckItem := r.appendToBank(0x09, "maku tree check item",
		"\xcd\x17\x17\xfa\xbb\xc6\xc9")
	r.replace(0x09, 0x7d93, "maku tree check item call",
		"\x3e\x40\xcd\x17\x17", "\x3e\x0a\xcd"+makuTreeCheckItem)

	// use a non-cutscene screen transition for exiting a dungeon via essence,
	// so that overworld music plays, and set maku tree state.
	essenceWarp := r.appendToBank(0x09, "essence warp",
		"\x3e\x81\xea\x67\xcc\xfa\xbb\xc6\xcd\x76\x01\xea\xdf\xc6\xc9")
	r.replace(0x09, 0x4b4f, "call essence warp",
		"\xea\x67\xcc", "\xcd"+essenceWarp)

	// use createTreasure for mt. cucco platform cave item, not
	// createRingTreasure.
	createMtCuccoItem := r.appendToBank(0x09, "create mt. cucco item",
		"\x01\x00\x00\xcd\x1b\x27\xc3\x21\x64")
	r.replace(0x09, 0x641a, "call create mt. cucco item",
		"\x01\x01\x27", "\xc3"+createMtCuccoItem)

	// bank 0a

	// set global flags and room flags that would be set during the intro, as
	// well as some other flags to skip cutscenes, etc.
	initialGlobalFlags := r.appendToBank(0x0a, "initial global flags",
		"\x0a\x1c\xff")
	giveLinkedStartItem := r.appendToBank(0x0a, "give linked starting item",
		"\xfa\xfe\x7f\x4f\xfa\xfd\x7f\xea\x82\xc6\xcd\xeb\x16"+
			"\xaf\xea\x80\xc6\xc9")
	setStartingFlags := r.appendToBank(0x0a, "set starting flags",
		"\xe5\x21"+initialGlobalFlags+"\x2a\xfe\xff\x28\x07"+
			"\xe5\xcd\xcd\x30\xe1\x18\xf4\xe1"+ // init global flags
			"\xfa\xff\x7f\xea\x10\xc6\x3e\xff\xea\x46\xc6"+ // animal stuff
			"\xfa\x01\xcc\xb7\xc4"+giveLinkedStartItem+ // in place of sword
			"\x3e\x50\xea\xa7\xc7"+ // bits 4 + 6
			"\x3e\x60\xea\x9a\xc7"+ // bits 5 + 6
			"\x3e\xc0\xea\x98\xc7\xea\xcb\xc7"+ // bits 6 + 7
			"\x3e\x40\xea\xb6\xc7\xea\x2a\xc8\xea\x00\xc8"+ // bit 6
			"\xea\x00\xc7\xea\x96\xc7\xea\x8d\xc7\xea\x60\xc7\xea\xd0\xc7"+
			"\xea\x1d\xc7\xea\x8a\xc7\xea\xe9\xc7\xea\x9b\xc7\xea\x29\xc8"+
			"\x3e\x10\xea\x97\xc6\x3e\x03\xea\xc6\xc6"+ // give L-3 ring box
			"\xc9")
	r.replace(0x0a, 0x66ed, "call set starting flags",
		"\x1e\x78\x1a\xcb\x7f\x20", "\xcd"+setStartingFlags+"\xc3\x9e\x21")

	// remove generic "you got a ring" text for gasha nuts
	gashaNutRingText := r.appendToBank(0x0a, "remove ring text from gasha nut",
		"\x79\xfe\x04\xc2\x4b\x18\xe1\xc9")
	r.replace(0x0a, 0x4863, "remove ring text from gasha nut caller",
		"\xc3\x4b\x18", "\xc3"+gashaNutRingText)

	// bank 0b

	// custom script command to use on d1 entrance screen: wait until bit of
	// cfc0 is set, and set ccaa to 01 meanwhile. fixes a vanilla bug where
	// dismounting an animal on that screen allowed you to enter without key.
	r.replace(0x0b, 0x4dea, "d1 entrance script cmd", "\xa0", "\xb2")
	d1EntranceFunc := r.appendToBank(0x0b, "d1 entrance cmd func",
		"\xe1\xfa\x49\xcc\xfe\x00\xc0\xfa\x4c\xcc\xfe\x96\xc0"+ // check room
			"\x3e\x01\xea\xaa\xcc\xaf\xc3\x2d\x43")
	r.replace(0x0b, 0x406d, "d1 entrance cmd jump", "\x03\x41", d1EntranceFunc)

	diverIDScript := r.appendToBank(0x0b, "diver fake id script",
		"\xde\x2e\x00\x92\x94\xc6\x02\xc1")
	r.replace(0x0b, 0x730d, "diver fake id call",
		"\xde\x2e\x00", "\xc0"+diverIDScript)

	// returns c,e = treasure ID,subID
	nobleSwordLookup := r.appendToBank(0x0b, "noble sword lookup",
		"\x21\x18\x64\x4e\x23\x5e\xc9")

	// skip forced ring appraisal and ring list with vasu (prevents softlock)
	r.replace(0x0b, 0x4a2b, "skip vasu ring appraisal",
		"\x98\x33", "\x4a\x39")

	// this will be overwritten after randomization
	smallKeyDrops := r.appendToBank(0x3f, "small key drops",
		makeKeyDropTable())
	lookUpKeyDropBank3F := r.appendToBank(0x3f, "look up key drop bank 3f",
		"\xc5\xfa\x49\xcc\x47\xfa\x4c\xcc\x4f\x21"+smallKeyDrops+ // load group/room
			"\xcd"+lookupWord+"\xc1\xc0\x46\x23\x4e\xc9")
	lookUpKeyDrop := r.appendToBank(0x0b, "look up key drop",
		"\x36\x60\x2c\xd5\xe5\x1e\x3f\x21"+lookUpKeyDropBank3F+
			"\xcd\x8a\x00\xe1\xd1\xc9")
	r.replace(0x0b, 0x4416, "call look up key drop",
		"\x36\x60\x2c", "\xcd"+lookUpKeyDrop)

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

	// look up item collection mode in a table based on room. if no entry is
	// found, the original mode (a) is preserved. the table is three bytes per
	// entry, (group, room, collect mode). ff ends the table. rooms that
	// contain more than one item are special cases.
	collectModeTable := r.appendToBank(0x15, "collection mode table",
		makeSeasonsCollectModeTable())
	// cp link's position if in diver room, set mode to 02 if on right side,
	// ret z if set
	collectModeDiver := r.appendToBank(0x15, "diver collect mode",
		"\x3e\x05\xb8\xc0\x3e\xbd\xb9\xc0\xfa\x0d\xd0\xfe\x80\xd8"+
			"\xaf\x3e\x02\xc9")
	// if in d7 compass room, set mode based on link's x position: left =
	// key drop; right = chest.
	collectModeD7Key := r.appendToBank(0x15, "d7 key collect mode",
		"\x3e\x05\xb8\xc0\x3e\x52\xb9\xc0"+ // cp room
			"\xfa\x0d\xd0\xfe\x80\x30\x04"+ // cp x
			"\xaf\x3e\x28\xc9\xaf\x3e\x38\xc9")
	// if link already has the maku tree's item, use default mode.
	collectModeMakuSeed := r.appendToBank(0x15, "maku seed collect mode",
		"\x3e\x02\xb8\xc0\x3e\x5d\xb9\xc0\x3e\x0a\xcd\x17\x17\x38\x02"+
			"\x3c\xc9\xaf\x7b\xc9")
	collectModeLookup := r.appendToBank(0x15, "collection mode lookup func",
		"\x5f\xc5\xe5\xfa\x49\xcc\x47\xfa\x4c\xcc\x4f\x21"+collectModeTable+
			"\x2a\xfe\xff\x28\x1d\xb8\x20\x16\x2a\xb9\x20\x13"+
			"\xcd"+collectModeDiver+"\x28\x12\xcd"+collectModeD7Key+"\x28\x0d"+
			"\xcd"+collectModeMakuSeed+"\x28\x08"+
			"\x2a\x18\x05\x23\x23\x18\xde\x7b\xe1\xc1\xc9")

	// upgrade normal items (interactions with ID 60) as necessary when they're
	// created, and set collection mode.
	normalProgressiveFunc := r.appendToBank(0x15, "normal progressive func",
		"\xcd"+collectModeLookup+"\x47\xcb\x37\xf5"+
			"\x1e\x43\x1a\xfe\x02\x30\x05"+ // don't upgrade spin slash
			"\x1b\x1a\xcd"+progressiveItemFunc+"\xf1\xc9")
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

	// call function to spawn item based on room instead of spawning a heart
	// container.
	bossItemTable := r.appendToBank(0x15, "boss item table", "\x00\x00"+
		"\x2a\x00\x2a\x00\x2a\x00\x2a\x00\x2a\x00\x2a\x00\x2a\x00\x2a\x00")
	spawnBossItem := r.appendToBank(0x15, "spawn boss item",
		"\xe5\x21"+bossItemTable+"\xfa\x55\xcc\xdf"+
			"\x46\x23\x4e\xcd\x1b\x27\xcd\xfd\x21\xe1\xc9")
	// some dungeons share the same script for spawning the HC.
	r.replace(0x0b, 0x4b8f, "HC call 1", "\xdd\x2a\x00", "\xe0"+spawnBossItem)
	r.replace(0x0b, 0x4bb1, "HC call 2", "\xdd\x2a\x00", "\xe0"+spawnBossItem)

	// bank 1f

	// replace ring appraisal text with "you got the {ring}"
	r.replace(0x1f, 0x5d99, "obtain ring text replacement",
		"\x03\x13\x20\x49\x04\x06", "\x02\x03\x0f\xfd\x21\x00")

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
		"\xf5\xd5\xe5\x78\xfe\x0e\x20\x15\x1e\xaf\x79\xd6\x0a\x12\xc6\x42"+
			"\x26\xc6\x6f\xfe\x45\x20\x04\xcb\xee\x18\x02\xcb\xfe"+
			"\xe1\xd1\xf1\xcd\x4e\x45\xc9")
	r.replace(0x3f, 0x452c, "flute set icon call", "\x4e\x45", setFluteIcon)

	// use different seed capacity table, so that level zero satchel can still
	// hold 20 seeds.
	seedCapTable := r.appendToBank(0x3f, "seed capacity table",
		"\x20\x20\x50\x99")
	r.replace(0x3f, 0x460e, "seed capacity pointer", "\x16\x46", seedCapTable)

	// put obtained rings directly into ring list (no need for appraisal), and tell the
	// player what type of ring it is
	r.replace(0x3f, 0x461a, "auto ring appraisal",
		"\xcb\xf1\xcd\x75\x46\xfe\x64\x38",
		"\x21\x16\xc6\x79\xe6\x3f\xcd\x0e\x02\x79\xc6\x40\xea\xb1\xcb\x01\x1c\x30\xcd\x4b\x18\xc9")

	// don't play a sound for obtaining an item if it's on the starting screen,
	// so that the linked starting item can be given silently.
	giveItemSilently := r.appendToBank(0x3f, "give item silently",
		"\x47\xfa\x4c\xcc\xfe\xa7\x78\xc2\x74\x0c"+
			"\xfa\x49\xcc\xb7\x78\xc2\x74\x0c\xc9")
	r.replace(0x3f, 0x4535, "call give item silently",
		"\xcd\x74\x0c", "\xcd"+giveItemSilently)

	// use different addresses for owl statue text.
	owlTextOffsets := r.appendToBank(0x3f, "owl text offsets",
		string(make([]byte, 0x1e*2))) // to be set later
	useOwlText := r.appendToBank(0x3f, "use owl text",
		"\xea\xd4\xd0\xfa\xa3\xcb\xfe\x3d\xc0"+ // ret if normal text
			"\x21"+owlTextOffsets+"\xfa\xa2\xcb\xdf\x2a\x66\x6f"+ // set addr
			"\x3e\x3f\xea\xd4\xd0\xc9") // set bank
	r.replace(0x3f, 0x4fd9, "call use owl text",
		"\xea\xd4\xd0", "\xcd"+useOwlText)

	// this *MUST* be the last thing in the bank, since it's going to grow
	// dynamically later.
	r.appendToBank(0x3f, "owl text", "")
}

// makes seasons-specific additions to the collection mode table.
func makeSeasonsCollectModeTable() string {
	b := new(strings.Builder)
	table := makeCollectModeTable()
	b.WriteString(table[:len(table)-1]) // strip final ff

	// add other three star ore screens
	for _, room := range starOreRooms[1:] {
		b.Write([]byte{0x01, room, collectDig})
	}

	// add other eight maku tree screens
	for _, room := range makuTreeRooms[1:] {
		b.Write([]byte{0x02, room, collectFall})
	}

	// add linked hero's cave chest
	b.Write([]byte{0x05, 0x2c, collectChest})

	b.Write([]byte{0xff})
	return b.String()
}
