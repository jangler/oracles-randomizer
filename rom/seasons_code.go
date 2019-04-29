package rom

import (
	"strings"
)

func newSeasonsRomBanks() *romBanks {
	asm, err := newAssembler()
	if err != nil {
		panic(err)
	}

	r := romBanks{
		endOfBank: make([]uint16, 0x40),
		assembler: asm,
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

	// do this before loading asm files, since the size of this table varies
	// with the number of checks.
	r.appendToBank(0x06, "collectModeTable", makeSeasonsCollectModeTable())

	r.applyAsmFiles([]string{"/asm/common.yaml", "/asm/seasons.yaml"})

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
	globalRomBanks = r

	// bank 00

	r.replaceAsm(0x00, 0x0c76,
		"ld h,a; ld a,(ff00+b5)", "call filterMusic")
	r.replaceAsm(0x00, 0x39df,
		"push de; ld a,(ff00+8c)", "jp loadWinterLayout")

	giveItem := addrString(r.assembler.getDef("giveTreasureCustom"))

	callBank2 := addrString(r.assembler.getDef("callBank2"))
	searchValue := addrString(r.assembler.getDef("searchValue"))
	searchDoubleKey := addrString(r.assembler.getDef("searchDoubleKey"))
	readSeason := addrString(r.assembler.getDef("readDefaultSeason"))

	// bank 02

	r.replaceAsm(0x02, 0x602c,
		"jp nz,_closeMenu", "call nz,treeWarp")
	r.replaceAsm(0x02, 0x6089,
		"jp nz,_closeMenu", "call nz,treeWarp")
	r.replaceAsm(0x02, 0x5e9a,
		"call setMusicVolume", "call devWarp")

	r.replaceAsm(0x02, 0x5ec8,
		"call _mapMenu_checkRoomVisited", "call checkTreeVisited")
	r.replaceAsm(0x02, 0x65e1,
		"call _mapMenu_checkRoomVisited", "call checkTreeVisited")
	r.replaceAsm(0x02, 0x609b,
		"call _mapMenu_checkCursorRoomVisited", "call checkCursorVisited")

	r.replaceAsm(0x02, 0x56a1,
		"ld a,(wInventorySubmenu1CursorPos)", "call openRingList")
	r.replaceAsm(0x02, 0x6f4a,
		"call _ringMenu_updateSelectedRingFromList", "call autoEquipRing")
	r.replaceAsm(0x02, 0x5035,
		"call setMusicVolume", "call ringListGfxFix")

	// bank 03

	r.replaceAsm(0x03, 0x4d6b,
		"call decHlRef16WithCap", "call skipCapcom")

	// bank 04

	r.replaceAsm(0x00, 0x3854,
		"call applyAllTileSubstitutions", "call applyExtraTileSubstitutions")
	r.replaceAsm(0x04, 0x461e,
		"ld a,(wWarpDestIndex)", "call checkSetAnimalSavePoint")

	// bank 05

	r.replaceAsm(0x05, 0x5fe8,
		"call lookupCollisionTable", "call cliffLookup")

	r.replaceAsm(0x05, 0x493b,
		"call _specialObjectGetRelativeTileWithDirectionTable; or a",
		"call animalEntryIgnoreBridges; nop")
	r.replaceAsm(0x05, 0x71ea,
		"call _specialObjectGetRelativeTileWithDirectionTable; or a",
		"call animalEntryIgnoreBridges; nop")

	r.replaceAsm(0x05, 0x776b,
		"ld a,(wAnimalRegion)", "call checkFlute")
	r.replaceAsm(0x05, 0x7a65,
		"ld a,(wAnimalRegion)", "call checkFlute")

	// bank 06

	r.replaceAsm(0x06, 0x4774,
		"call setTile", "call checkBreakD6Flower")
	r.replaceAsm(0x06, 0x47f5,
		"call getFreePartSlot", "call dropExtraGalesOnEmpty")

	// use expert's or fist ring with only one button unequipped.
	r.replaceAsm(0x06, 0x490e, "ret nz", "nop")

	// bank 07

	r.replaceAsm(0x07, 0x5b75,
		"ld a,(wActiveTileType); cp a,08", "call devChangeSeason; nop; nop")

	// bank 08

	// use the custom "give item" function in the shop instead of the normal
	// one. this obviates some hard-coded shop data (sprite, text) and allows
	// the item to progressively upgrade.
	// param = b (item index/subID), returns c,e = treasure ID,subID
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
	// TODO: i'm not sure whether maku tree needs this anymore, since the
	//       collect mode func goes off script position now. is it still
	//       required to determine whether you picked up the drop or something?
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

	r.replaceAsm(0x0a, 0x66ed,
		"db 1e,78,1a,cb,7f,20", // dunno what this is
		"call setInitialFlags; jp objectDelete_useActiveObjectType")

	r.replaceAsm(0x0a, 0x7b93,
		"call giveTreasure", "call giveTreasureCustom")
	r.replaceAsm(0x0a, 0x7b9e,
		"jp showText", "ret; nop; nop")

	// remove generic "you got a ring" text for gasha nuts
	gashaNutRingText := r.appendToBank(0x0a, "remove ring text from gasha nut",
		"\x79\xfe\x04\xc2\x4b\x18\xe1\xc9")
	r.replace(0x0a, 0x4863, "remove ring text from gasha nut caller",
		"\xc3\x4b\x18", "\xc3"+gashaNutRingText)

	// bank 0b

	// command and corresponding address in jump table
	r.replace(0x0b, 0x4dea, "d1 entrance cmd byte", "\xa0", "\xb2")
	r.replace(0x0b, 0x406d, "jump d1EntranceScriptCmd",
		"\x03\x41", addrString(r.assembler.getDef("d1EntranceScriptCmd")))

	diverIDScript := r.appendToBank(0x0b, "diver fake id script",
		"\xde\x2e\x00\x92\x94\xc6\x02\xc1")
	r.replace(0x0b, 0x730d, "diver fake id call",
		"\xde\x2e\x00", "\xc0"+diverIDScript)

	// skip forced ring appraisal and ring list with vasu (prevents softlock)
	r.replace(0x0b, 0x4a2b, "skip vasu ring appraisal",
		"\x98\x33", "\x4a\x39")

	// this will be overwritten after randomization
	smallKeyDrops := r.appendToBank(0x3f, "small key drops",
		makeKeyDropTable())
	lookUpKeyDropBank3F := r.appendToBank(0x3f, "look up key drop bank 3f",
		"\xc5\xfa\x49\xcc\x47\xfa\x4c\xcc\x4f\x21"+smallKeyDrops+ // load group/room
			"\x1e\x02\xcd"+searchDoubleKey+"\xc1\xd0\x46\x23\x4e\xc9")
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

	r.replaceAsm(0x15, 0x465a,
		"ld b,a; swap a", "call modifyTreasure")

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
	r.replaceAsm(0x15, 0x70cf,
		"call giveTreasure", "call giveTreasureCustom")

	// some dungeons share the same script for spawning the HC.
	r.replaceAsm(0x0b, 0x4b8f,
		"db dd,2a,00", "db e0; dw spawnBossItem")
	r.replaceAsm(0x0b, 0x4bb1,
		"db dd,2a,00", "db e0; dw spawnBossItem")

	// bank 1f

	// replace ring appraisal text with "you got the {ring}"
	r.replace(0x1f, 0x5d99, "obtain ring text replacement",
		"\x03\x13\x20\x49\x04\x06", "\x02\x03\x0f\xfd\x21\x00")

	// bank 3f

	r.replaceAsm(0x00, 0x16f6,
		"call giveTreasure_body", "call satchelRefillSeeds")
	r.replaceAsm(0x3f, 0x4356,
		"call _interactionGetData", "call checkLoadCustomSprite")

	// "activate" a flute by setting its icon and song when obtained. also
	// activates the corresponding animal companion.
	setFluteIcon := r.appendToBank(0x3f, "flute set icon func",
		"\xf5\xd5\xe5\x78\xfe\x0e\x20\x15\x1e\xaf\x79\xd6\x0a\x12\xc6\x42"+
			"\x26\xc6\x6f\xfe\x45\x20\x04\xcb\xee\x18\x02\xcb\xfe"+
			"\xe1\xd1\xf1\xcd\x4e\x45\xc9")
	r.replace(0x3f, 0x452c, "flute set icon call", "\x4e\x45", setFluteIcon)

	r.replace(0x3f, 0x460e, "seed capacity pointer",
		"\x16\x46", addrString(r.assembler.getDef("seedCapacityTable")))

	r.replaceAsm(0x3f, 0x461a,
		"set 6,c; call realignUnappraisedRings", "nop; jp autoAppraiseRing")

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
