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
	r.appendToBank(0x38, "roomTreasures", makeRoomTreasureTable())
	r.appendToBank(0x3f, "owlTextOffsets", string(make([]byte, 0x14*2)))

	r.applyAsmFiles(GameAges,
		[]string{
			"/asm/common.yaml",
			"/asm/ages.yaml",
		},
		[]string{
			"/asm/item_lookup.yaml",
			"/asm/layouts.yaml",
			"/asm/rings.yaml",

			"/asm/text.yaml", // must go last
		})

	return &r
}

func initAgesEOB() {
	globalRomBanks = newAgesRomBanks()
}

// makes ages-specific additions to the collection mode table.
func makeAgesCollectModeTable() string {
	b := new(strings.Builder)
	table := makeCollectModeTable()
	b.WriteString(table[:len(table)-1]) // strip final ff

	// add eatern symmetry city brother
	b.Write([]byte{0x03, 0x6f, collectModes["touch"]})

	// add ricky and dimitri nuun caves
	b.Write([]byte{0x02, 0xec, collectModes["chest"],
		0x05, 0xb8, collectModes["chest"]})

	b.Write([]byte{0xff})
	return b.String()
}
