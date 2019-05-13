package rom

func newAgesRomBanks() *romBanks {
	asm, err := newAssembler()
	if err != nil {
		panic(err)
	}

	r := romBanks{
		endOfBank: loadBankEnds("ages"),
		assembler: asm,
	}

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	r.replaceRaw(Addr{0x06, 0}, "collectModeTable", makeCollectModeTable())
	r.replaceRaw(Addr{0x38, 0}, "roomTreasures", makeRoomTreasureTable(GameAges))
	r.replaceRaw(Addr{0x3f, 0}, "owlTextOffsets", string(make([]byte, 0x14*2)))

	r.applyAsmFiles(GameAges,
		[]string{
			"/asm/common.yaml",
			"/asm/ages.yaml",
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
