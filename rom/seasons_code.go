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
		endOfBank: loadBankEnds("seasons"),
		assembler: asm,
	}

	// do this before loading asm files, since the sizes of the tables vary
	// with the number of checks.
	r.appendToBank(0x06, "collectModeTable", makeSeasonsCollectModeTable())
	r.appendToBank(0x3f, "roomTreasures", makeRoomTreasureTable())
	r.appendToBank(0x3f, "owlTextOffsets", string(make([]byte, 0x1e*2)))

	r.applyAsmFiles(GameSeasons,
		[]string{
			"/asm/common.yaml",
			"/asm/seasons.yaml",
		},
		[]string{
			"/asm/item_lookup.yaml",
			"/asm/layouts.yaml",
			"/asm/rings.yaml",

			"/asm/text.yaml", // must go last
		})

	return &r
}

// for some reason the maku tree has a different room for every number of
// essences you have.
var (
	makuTreeRooms = []byte{0x0b, 0x0c, 0x7b, 0x2b, 0x2c, 0x2d, 0x5b, 0x5c, 0x5d}
	starOreRooms  = []byte{0x66, 0x76, 0x75, 0x65}
)

func initSeasonsEOB() {
	globalRomBanks = newSeasonsRomBanks()
}

// makes seasons-specific additions to the collection mode table.
func makeSeasonsCollectModeTable() string {
	b := new(strings.Builder)
	table := makeCollectModeTable()
	b.WriteString(table[:len(table)-1]) // strip final ff

	// add other three star ore screens
	for _, room := range starOreRooms[1:] {
		b.Write([]byte{0x01, room, collectModes["dig"]})
	}

	// add other eight maku tree screens
	for _, room := range makuTreeRooms[1:] {
		b.Write([]byte{0x02, room, collectModes["maku tree (seasons)"]})
	}

	// add linked hero's cave chest
	b.Write([]byte{0x05, 0x2c, collectModes["chest"]})

	b.Write([]byte{0xff})
	return b.String()
}
