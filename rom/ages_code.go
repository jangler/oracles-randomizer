package rom

import (
	"strings"
)

func newAgesRomBanks() *romBanks {
	asm, err := newAssembler()
	if err != nil {
		panic(err)
	}

	r := romBanks{
		endOfBank: make([]uint16, 0x40),
		assembler: asm,
	}

	r.endOfBank[0x00] = 0x3ef8
	r.endOfBank[0x01] = 0x7fc3
	r.endOfBank[0x02] = 0x7e93
	r.endOfBank[0x03] = 0x7ebd
	r.endOfBank[0x04] = 0x7edb
	r.endOfBank[0x05] = 0x7d9d
	r.endOfBank[0x06] = 0x7a31
	r.endOfBank[0x08] = 0x7f60
	r.endOfBank[0x09] = 0x7dee
	r.endOfBank[0x0a] = 0x7e09
	r.endOfBank[0x0b] = 0x7fa8
	r.endOfBank[0x0c] = 0x7f94
	r.endOfBank[0x0f] = 0x7f90
	r.endOfBank[0x10] = 0x7ef4
	r.endOfBank[0x11] = 0x7f73
	r.endOfBank[0x12] = 0x7e8f
	r.endOfBank[0x15] = 0x7bfb
	r.endOfBank[0x16] = 0x7e03
	r.endOfBank[0x38] = 0x6b00 // to be safe
	r.endOfBank[0x3f] = 0x7d0a

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	r.appendToBank(0x06, "collectModeTable", makeAgesCollectModeTable())
	r.appendToBank(0x38, "smallKeyDrops", makeKeyDropTable())
	r.appendToBank(0x3f, "owlTextOffsets", string(make([]byte, 0x14*2)))

	r.applyAsmFiles([]string{"/asm/common.yaml", "/asm/ages.yaml"})

	return &r
}

func initAgesEOB() {
	r := newAgesRomBanks()
	globalRomBanks = r

	// bank 00

	r.replaceAsm(0x00, 0x3e56,
		"inc a; cp a,11", "call checkMakuState")

	compareRoom := addrString(r.assembler.getDef("compareRoom"))
	findObjectWithId := addrString(r.assembler.getDef("findObjectWithId"))

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

	r.replaceMultiple([]Addr{{0x02, 0x6133}, {0x02, 0x618b}}, "tree warp jump",
		"\xc2\xba\x4f", "\xc4"+addrString(r.assembler.getDef("treeWarp")))
	r.replaceAsm(0x02, 0x5fcb, "call setMusicVolume", "call devWarp")

	r.replaceAsm(0x02, 0x5ff9,
		"call _mapMenu_checkRoomVisited", "call checkTreeVisited")
	r.replaceAsm(0x02, 0x66a9,
		"call _mapMenu_checkRoomVisited", "call checkTreeVisited")
	r.replaceAsm(0x02, 0x619d,
		"call _mapMenu_checkCursorRoomVisited", "call checkCursorVisited")
	r.replaceAsm(0x02, 0x6245,
		"ld a,(wMapMenu_cursorIndex)", "jp displayPortalPopups")

	r.replaceAsm(0x02, 0x56dd,
		"ld a,(wInventorySubmenu1CursorPos)", "call openRingList")
	r.replaceAsm(0x02, 0x7019,
		"call _ringMenu_updateSelectedRingFromList", "call autoEquipRing")
	r.replaceAsm(0x02, 0x5074,
		"call setMusicVolume", "call ringListGfxFix")

	// bank 03

	r.replaceAsm(0x03, 0x4d6b,
		"call decHlRef16WithCap", "call skipCapcom")
	r.replaceAsm(0x03, 0x6e97,
		"jp setGlobalFlag", "jp setInitialFlags")

	// bank 04

	r.replaceAsm(0x00, 0x38c0,
		"call applyAllTileSubstitutions", "call applyExtraTileSubstitutions")

	// treat the d2 present entrance like the d2 past entrance, and reset the
	// water level when entering jabu (see logic comments).
	replaceWarpEnter := r.appendToBank(0x04, "replace warp enter",
		"\xc5\x01\x00\x83\xcd"+compareRoom+"\x20\x04\xc1\x3e\x01\xc9"+
			"\x01\x02\x90\xcd"+compareRoom+"\xc1\x20\x05\x3e\x21\xea\xe9\xc6"+
			"\xfa\x2d\xcc\xc9")
	r.replace(0x04, 0x4630, "call replace warp enter",
		"\xfa\x2d\xcc", "\xcd"+replaceWarpEnter)
	// d2: exit into the present if the past entrance is closed.
	replaceWarpExit := r.appendToBank(0x00, "replace warp exit",
		"\xea\x48\xcc\xfe\x83\xc0\xfa\x83\xc8\xe6\x80\xc0"+
			"\xfa\x47\xcc\xe6\x0f\xfe\x01\xc0"+
			"\xfa\x47\xcc\xe6\xf0\xea\x47\xcc\xc9")
	r.replace(0x04, 0x45e8, "call replace warp exit normal",
		"\xea\x48\xcc", "\xcd"+replaceWarpExit)
	r.replace(0x0a, 0x4738, "call replace warp exit essence",
		"\xea\x48\xcc", "\xcd"+replaceWarpExit)

	// bank 05

	r.replaceAsm(0x05, 0x516c,
		"ld a,(cc91); or a", "call checkPreventSurface; nop")
	r.replaceAsm(0x05, 0x6083,
		"call lookupCollisionTable", "call cliffLookup")

	// bank 06

	// burning the first tree in yoll graveyard should set room flag 1 so that
	// it can be gone for good.
	removeYollTree := r.appendToBank(0x06, "remove yoll tree",
		"\xf5\xf0\x8f\xfe\x0c\x20\x0f"+
			"\xc5\x01\x00\x6b\xcd"+compareRoom+"\x20\x05"+
			"\x21\x6b\xc7\xcb\xce\xc1\xf1\x21\x26\xc6\xc9")
	r.replace(0x06, 0x47aa, "call remove yoll tree",
		"\x21\x26\xc6", "\xcd"+removeYollTree)

	// reenter a warp tile that link is standing on when playing the tune of
	// currents (useful if you warp into a patch of bushes). also activate the
	// west present crescent island portal.
	reenterCurrentsWarp := r.appendToBank(0x06, "special currents actions",
		"\xc5\x01\x00\xa9\xcd"+compareRoom+"\xc1\x20\x11"+ // island portal
			"\xd5\x3e\xe1\xcd"+findObjectWithId+"\x20\x05"+ // cont.
			"\x1e\x44\x3e\x02\x12\xd1\xc3\x08\x4e"+ // cont.
			"\xfa\x34\xcc\xf5\xd5\x3e\xde\xcd"+findObjectWithId+ // reenter
			"\x20\x05\x1e\x44\x3e\x02\x12\xd1\xf1\xc3\x37\x4e") // cont.
	r.replace(0x06, 0x4e34, "call special currents actions",
		"\xfa\x34\xcc", "\xc3"+reenterCurrentsWarp)

	// set text index for portal sign on crescent island.
	setPortalSignText := r.appendToBank(0x06, "set portal sign text",
		"\x01\x00\xa9\xcd"+compareRoom+"\x01\x01\x09\xc0\x01\x01\x56\xc9")
	r.replace(0x06, 0x40e7, "call set portal sign text",
		"\x01\x01\x09", "\xcd"+setPortalSignText)

	// use expert's or fist ring with only one button unequipped.
	r.replaceAsm(0x06, 0x4969, "ret nz", "nop")

	// bank 16 (pt. 1)

	getTreasureDataBCE := addrString(r.assembler.getDef("getTreasureDataBCE"))

	// bank 09

	// set treasure ID 07 (rod of seasons) when buying the 150 rupee shop item,
	// so that the shop can check this specific ID.
	shopSetFakeID := r.appendToBank(0x09, "shop set fake ID",
		"\xfe\x0d\x20\x05\x21\x9a\xc6\xcb\xfe\x21\xf7\x44\xc9")
	r.replace(0x09, 0x4418, "call shop set fake ID",
		"\x21\xf7\x44", "\xcd"+shopSetFakeID)

	r.replaceAsm(0x09, 0x4c4e,
		"call giveTreasure", "call handleGetItem")
	handleGetItem := addrString(r.assembler.getDef("handleGetItem"))

	// remove generic "you got a ring" text for rings from shops
	r.replace(0x09, 0x4580, "obtain ring text replacement (shop) 1", "\x54", "\x00")
	r.replace(0x09, 0x458a, "obtain ring text replacement (shop) 2", "\x54", "\x00")
	r.replace(0x09, 0x458b, "obtain ring text replacement (shop) 3", "\x54", "\x00")

	// remove generic "you got a ring" text for gasha nuts
	gashaNutRingText := r.appendToBank(0x0b, "remove ring text from gasha nut",
		"\x79\xfe\x04\xc2\x72\x18\xe1\xc9")
	r.replace(0x0b, 0x45bb, "remove ring text from gasha nut caller",
		"\xc3\x72\x18", "\xc3"+gashaNutRingText)

	// don't set room's item flag if it's nayru's item on the maku tree screen,
	// since link still might not have taken the maku tree's item.
	makuTreeItemFlag := r.appendToBank(0x09, "maku tree item flag",
		"\xcd\x7d\x19\xc5\x01\x38\xc7\xcd\xd6\x01\xc1\x20\x06\xfa\x0d\xd0"+
			"\xfe\x50\xc8\xcb\xee\xc9")
	r.replace(0x09, 0x4c82, "call maku tree item flag",
		"\xcd\x7d\x19", "\xc3"+makuTreeItemFlag)

	// give correct ID and param for shop item, play sound, and load correct
	// text index into temp wram address.
	shopGiveTreasure := r.appendToBank(0x09, "shop give treasure",
		"\x47\x1a\xfe\x0d\x78\x20\x08\xcd"+getTreasureDataBCE+"\x7b\xea\x0d\xcf"+
			"\x78\xcd"+handleGetItem+"\xc2\x98\x0c\x3e\x4c\xc3\x98\x0c")
	r.replace(0x09, 0x4425, "call shop give treasure",
		"\xcd\x1c\x17", "\xcd"+shopGiveTreasure)
	// display text based on above temp wram address.
	shopShowText := r.appendToBank(0x09, "shop show text",
		"\x1a\xfe\x0d\xc2\x72\x18\xfa\x0d\xcf\x06\x00\x4f"+
			"\x79\xfe\xff\xc8\xc3\x72\x18") // text $ff is ring
	r.replace(0x09, 0x4443, "call shop show text",
		"\xc2\x72\x18", "\xc2"+shopShowText)

	// bank 0a

	r.replaceAsm(0x0a, 0x4d5f,
		"call checkTreasureObtained", "call checkRaftonHasRope")
	r.replaceAsm(0x0a, 0x4bb8,
		"call checkGlobalFlag", "call checkShouldRickyAppear")
	r.replaceAsm(0x0a, 0x5541,
		"call setGlobalFlag", "call saveMakuTreeWithNayru")

	// set sub ID for south shore dig item.
	dirtSpawnItem := r.appendToBank(0x0a, "dirt spawn item",
		"\xcd\xd4\x27\xc0\xcd\x42\x22\xaf\xc9")
	r.replace(0x0a, 0x5e3e, "call dirt spawn item",
		"\xcd\xc5\x24", "\xcd"+dirtSpawnItem)

	// use a non-cutscene screen transition for exiting a dungeon via essence,
	// so that overworld music plays, and set maku tree state.
	essenceWarp := r.appendToBank(0x0a, "essence warp",
		"\x3e\x81\xea\x4b\xcc\xc3\x53\x3e")
	r.replace(0x0a, 0x4745, "call essence warp",
		"\xea\x4b\xcc", "\xcd"+essenceWarp)

	// bank 0b

	r.replaceAsm(0x0b, 0x5464,
		"call checkGlobalFlag", "call checkKingZoraSequence")
	r.replaceAsm(0x0b, 0x61d7,
		"call getThisRoomFlags", "call checkZoraGuards")
	r.replaceAsm(0x0b, 0x7954,
		"ld (wCutsceneTrigger),a", "call skipFairyQueenCutscene")

	// bank 0c

	r.replaceAsm(0x0c, 0x442e,
		"ld (hl),60; inc l", "call lookupKeyDropBank0c")
	r.replaceAsm(0x08, 0x5087,
		"ld bc,3001", "call lookupKeyDropBank08")
	r.replaceAsm(0x0a, 0x7075,
		"ld bc,3001", "call lookupKeyDropBank0a")

	r.replaceAsm(0x0c, 0x4bd8,
		"db dd,2a,00", "db e0; dw spawnBossItem")

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

	// set room flags for other side of symmetry city bridge at end of building
	// cutscene.
	setBridgeFlag := r.appendToBank(0x15, "set bridge flag",
		"\xe5\xaf\xea\x8a\xcc\x3e\x25\xcd\xf9\x31"+
			"\x21\x24\xc7\xcb\xce\xe1\xc9")
	r.replace(0x0c, 0x7a6f, "call set bridge flag",
		"\xb9\xb6\x25", "\xe0"+setBridgeFlag)

	// skip forced ring appraisal and ring list with vasu (prevents softlock)
	r.replace(0x0c, 0x4a27, "skip vasu ring appraisal",
		"\x98\x33", "\x4a\x35")

	// bank 0f

	// set room flag for tunnel behind keep when defeating great moblin.
	setTunnelFlag := r.appendToBank(0x0f, "set tunnel flag",
		"\x21\x09\xc7\xcb\xc6\x21\xda\xca\xc9")
	r.replace(0x0f, 0x7f3e, "call set tunnel flag",
		"\x21\x09\xc7", "\xcd"+setTunnelFlag)

	// bank 10

	r.replaceAsm(0x10, 0x7914,
		"ld hl,7927", "call checkBlackTowerState")
	r.replaceAsm(0x10, 0x7d88,
		"ld a,(cc8d)", "call checkActivateEchoesPortal")

	// bank 11

	r.replaceAsm(0x11, 0x4aba,
		"call checkTreasureObtained", "call checkCanHarvestSeeds")

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
	// and outside D2 present.
	d2PresentPortal := r.appendToBank(0x12, "d2 present portal",
		"\xf2\xdc\x02\x48\x38\xe1\x00\x48\x48\xfe")
	r.replace(0x12, 0x5d42, "d2 present portal pointer",
		"\xdc\x02\x48\x38", "\xf3"+d2PresentPortal+"\xff")

	// bank 15

	r.replaceAsm(0x15, 0x50ae,
		"ld a,TREASURE_SWORD; ldi (hl),a", "call setShootingGalleryEquips")

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

	r.replaceAsm(0x16, 0x4539,
		"ld b,a; swap a", "call modifyTreasure")

	// bank 21

	// replace ring appraisal text with "you got the {ring}"
	r.replace(0x21, 0x76a0, "obtain ring text replacement",
		"\x04\x2c\x20\x04\x96\x21", "\x02\x06\x0f\xfd\x21\x00")

	// bank 3f

	r.replaceAsm(0x3f, 0x4356,
		"call _interactionGetData", "call checkLoadCustomSprite")
	r.replaceAsm(0x3f, 0x4607,
		"ld hl,4610", "ld hl,seedCapacityTable")
	r.replaceAsm(0x3f, 0x4614,
		"set 6,c; call realignUnappraisedRings", "nop; jp autoAppraiseRing")
	r.replaceAsm(0x3f, 0x4faa,
		"ld (w7ActiveBank),a", "call useOwlText")

	// this *MUST* be the last thing in the bank, since it's going to grow
	// dynamically later.
	r.appendToBank(0x38, "owlText", "")
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
