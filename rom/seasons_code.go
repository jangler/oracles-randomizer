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
	r.endOfBank[0x14] = 0x6fd0 // amazing
	r.endOfBank[0x15] = 0x792d
	r.endOfBank[0x3f] = 0x714d

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	r.appendToBank(0x06, "collectModeTable", makeSeasonsCollectModeTable())
	r.appendToBank(0x3f, "smallKeyDrops", makeKeyDropTable())
	r.appendToBank(0x3f, "owlTextOffsets", string(make([]byte, 0x1e*2)))

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

	r.replaceAsm(0x08, 0x4bfb,
		"call giveTreasure", "call shopGiveTreasure")
	r.replaceAsm(0x08, 0x62a7,
		"call checkGlobalFlag", "call checkBeachItemObtained")
	r.replaceAsm(0x08, 0x62f2,
		"inc l; ld (hl),45", "call setStarOreIds")

	// allow desert pits to work even if player has the actual bell already.
	r.replaceAsm(0x08, 0x73a2, "jr c,09", "nop; nop")

	r.replaceAsm(0x08, 0x5663, "dw 4e87", "dw script_checkDisplayWarning")

	// (volcano cutscene skip)
	r.replaceAsm(0x08, 0x7d07,
		"ld a,(cd18); or a; ret nz; call getThisRoomFlags;"+
			"set 6,(hl); ld a,0b; ld (cc04),a",
		"ld a,(d244); cp a,01; ret nz; call interactionDelete;"+
			"ld hl,skipVolcanoCutscene; jp callBank2")

	// enable exit from volcano room after skipping cutscene.
	r.replaceAsm(0x08, 0x7cf5, "ld (ccab),a", "nop; nop; nop")

	// remove generic "you got a ring" text for rings from shops
	r.replace(0x08, 0x4d55, "obtain ring text replacement (shop) 1", "\x54", "\x00")
	r.replace(0x08, 0x4d56, "obtain ring text replacement (shop) 2", "\x54", "\x00")

	// bank 09

	r.replaceAsm(0x09, 0x42e0,
		"call giveTreasure", "call setFakeIdsForStarOreAndMakuTree")
	r.replaceAsm(0x09, 0x4b4f,
		"ld (wWarpTransition2),a", "call essenceWarp")
	r.replaceAsm(0x09, 0x4d9a,
		"call 4ed9", "call checkFluteCollisions")
	r.replaceAsm(0x09, 0x4dad,
		"call 4ed9", "call checkFluteCollisions")
	r.replaceAsm(0x09, 0x641a,
		"ld bc,2701", "jp createMtCuccoItem")
	r.replaceAsm(0x09, 0x7887,
		"rst 18; ldi a,(hl); ld c,(hl)", "call tradeStarOre")
	r.replaceAsm(0x09, 0x7d95,
		"call checkTreasureObtained", "call makuTreeCheckItem")

	// use custom "give item" func in subrosia market.
	r.replaceAsm(0x09, 0x788a,
		"cp a,2d; jr nz,03; call getRandomRingOfGivenTier; call giveTreasure; ld e,42",
		"db 00,00,00,00,00,00,00; call marketGiveTreasure; jr c,0b")

	// bank 0a

	r.replaceAsm(0x0a, 0x4863,
		"jp showText", "jp removeGashaNutRingText")
	r.replaceAsm(0x0a, 0x7b93,
		"call giveTreasure", "call giveTreasureCustom")
	r.replaceAsm(0x0a, 0x7b9e,
		"jp showText", "ret; nop; nop")

	r.replaceAsm(0x0a, 0x66ed,
		"db 1e,78,1a,cb,7f,20", // dunno what this is
		"call setInitialFlags; jp objectDelete_useActiveObjectType")

	// bank 0b

	// command and corresponding address in jump table
	r.replace(0x0b, 0x4dea, "d1 entrance cmd byte", "\xa0", "\xb2")
	r.replaceAsm(0x0b, 0x406d,
		"dw scriptEnd", "dw d1EntranceScriptCmd")

	r.replaceAsm(0x0b, 0x730d,
		"db giveitem; db TREASURE_FLIPPERS; db 00",
		"db callscript; dw script_diverGiveItem")

	// skip forced ring appraisal and ring list with vasu (prevents softlock)
	r.replaceAsm(0x0b, 0x4a2b,
		"db showtext; db 33", "dw 394a")

	r.replaceAsm(0x0b, 0x4416,
		"ld (hl),60; inc l", "call lookupKeyDropBank0b")

	// some dungeons share the same script for spawning the HC.
	r.replaceAsm(0x0b, 0x4b8f,
		"db dd,2a,00", "db e0; dw spawnBossItem")
	r.replaceAsm(0x0b, 0x4bb1,
		"db dd,2a,00", "db e0; dw spawnBossItem")

	// bank 11

	// these are all for adding warning interactions
	r.replaceAsm(0x11, 0x6c10,
		"db f2,1f,08,68", "db f3; dw waterfallStaticObjects; db ff")
	r.replaceAsm(0x11, 0x6568,
		"db f2,9c,00,58", "db f3; dw flowerCliffStaticObjects; db ff")
	r.replaceAsm(0x11, 0x69cc,
		"db f2,1f,0d,68", "db f3; dw divingSpotStaticObjects; db ff")
	r.replaceAsm(0x11, 0x650b,
		"db f2,ab,00,40", "db f3; dw moblinKeepStaticObjects; db ff")
	r.replaceAsm(0x11, 0x7ada,
		"db f3,93,55", "db f3; dw hssSkipStaticObjects; db ff")

	// bank 14

	r.replaceAsm(0x00, 0x25d9,
		"ld e,41; ld a,(de)", "call overrideAnimationId")
	r.replaceAsm(0x00, 0x2600,
		"ld e,41; ld a,(de)", "call overrideAnimationId")

	// bank 15

	r.replaceAsm(0x15, 0x465a,
		"ld b,a; swap a", "call modifyTreasure")
	r.replaceAsm(0x15, 0x5a0e,
		"call setGlobalFlag", "call setPirateCutsceneFlags")
	r.replaceAsm(0x15, 0x5b83,
		"inc l; ld (hl),52", "call setHardOreIds")

	// use custom "give item" func in rod cutscene.
	r.replaceAsm(0x15, 0x70cf,
		"call giveTreasure", "call giveTreasureCustom")

	// skip "you got all four seasons" text from season spirits.
	r.replaceAsm(0x15, 0x57c2, "cp a,04", "cp a,05")

	// bank 1f

	// replace ring appraisal text with "you got the {ring}"
	r.replace(0x1f, 0x5d99, "obtain ring text replacement",
		"\x03\x13\x20\x49\x04\x06", "\x02\x03\x0f\xfd\x21\x00")

	// bank 3f

	r.replaceAsm(0x00, 0x16f6,
		"call giveTreasure_body", "call satchelRefillSeeds")
	r.replaceAsm(0x3f, 0x452b,
		"call applyParameter", "call activateFlute")
	r.replaceAsm(0x3f, 0x4356,
		"call _interactionGetData", "call checkLoadCustomSprite")
	r.replaceAsm(0x3f, 0x4535,
		"call playSound", "call playSoundExceptForLinkedStartItem")
	r.replaceAsm(0x3f, 0x460d,
		"ld hl,4616", "ld hl,seedCapacityTable")
	r.replaceAsm(0x3f, 0x461a,
		"set 6,c; call realignUnappraisedRings", "nop; jp autoAppraiseRing")
	r.replaceAsm(0x3f, 0x4fd9,
		"ld (w7ActiveBank),a", "call useOwlText")

	// this *MUST* be the last thing in the bank, since it's going to grow
	// dynamically later.
	r.appendToBank(0x3f, "owlText", "")
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
		b.Write([]byte{0x02, room, collectSeasonsMakuTree})
	}

	// add linked hero's cave chest
	b.Write([]byte{0x05, 0x2c, collectChest})

	b.Write([]byte{0xff})
	return b.String()
}
