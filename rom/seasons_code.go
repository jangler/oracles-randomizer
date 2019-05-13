package rom

// TODO: combine this and newAgesRomBanks into newRomBanks.
func newSeasonsRomBanks() *romBanks {
	asm, err := newAssembler()
	if err != nil {
		panic(err)
	}

	r := romBanks{
		endOfBank: loadBankEnds("seasons"),
		assembler: asm,
	}

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	r.replaceRaw(Addr{0x06, 0}, "collectModeTable", makeCollectModeTable())
	r.replaceRaw(Addr{0x3f, 0}, "roomTreasures", makeRoomTreasureTable(GameSeasons))
	r.replaceRaw(Addr{0x3f, 0}, "owlTextOffsets", string(make([]byte, 0x1e*2)))

	r.applyAsmFiles(GameSeasons,
		[]string{
			"/asm/common.yaml",
			"/asm/seasons.yaml",
		},
		[]string{
			"/asm/animals.yaml",
			"/asm/cutscenes.yaml",
			"/asm/gfx.yaml",
			"/asm/item_events.yaml",
			"/asm/item_lookup.yaml",
			"/asm/layouts.yaml",
			"/asm/linked.yaml",
			"/asm/misc.yaml",
			"/asm/rings.yaml",
			"/asm/triggers.yaml",
			"/asm/vars.yaml",

			"/asm/text.yaml", // must go last
		})

	return &r
}
