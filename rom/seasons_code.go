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
	globalRomBanks = newSeasonsRomBanks()
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
