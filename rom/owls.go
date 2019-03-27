package rom

import (
	"strings"
)

// this file is so small that i don't really want to make a separate one for
// each game. it maps owl statue names (matching those in the logic package) to
// the low bytes of their respective text IDs.

var seasonsOwls = map[string]byte{
	"dodongo owl":          0x0d,
	"gohma owl":            0x0e,
	"armos knights owl":    0x0f,
	"silent watch owl":     0x10,
	"magical ice owl":      0x11,
	"woods of winter owl":  0x14,
	"omuai owl":            0x15,
	"poe curse owl":        0x16,
	"spiked beetles owl":   0x17,
	"trampoline owl":       0x18,
	"greater distance owl": 0x19,
	"frypolar owl":         0x1a,
	"shining blue owl":     0x1c,
	"floodgate keeper owl": 0x1d,
	// "roller owl":        0x12, // unused
	// "guide owl":         0x13, // unused
}

var agesOwls = map[string]byte{
	"crown dungeon owl":    0x00,
	"spiked beetles owl":   0x01,
	"ancient words owl":    0x02,
	"blue wing owl":        0x03,
	"talus peaks owl":      0x05,
	"deku forest owl":      0x06,
	"head thwomp owl":      0x07,
	"scent seduction owl":  0x08,
	"deep waters owl":      0x09,
	"open ears owl":        0x0a,
	"luck test owl":        0x0b,
	"stone soldiers owl":   0x0c,
	"four crystals owl":    0x0d,
	"black tower owl":      0x0e,
	"rolling ridge owl":    0x0f,
	"jabu switch room owl": 0x10,
	"plasmarine owl":       0x11,
	"golden isle owl":      0x12,
	"mermaid legend owl":   0x13,
}

// updates the owl statue text data based on the given hints. does not mutate
// anything.
func SetOwlData(owlHints map[string]string, game int) {
	table := codeMutables["owl text offsets"].(*MutableRange)
	text := codeMutables["owl text"].(*MutableRange)
	builder := new(strings.Builder)
	addr := text.Addrs[0].offset

	var owlTextIDs map[string]byte
	if game == GameSeasons {
		owlTextIDs = seasonsOwls
	} else {
		owlTextIDs = agesOwls
	}

	for owlName, hint := range owlHints {
		textID := owlTextIDs[owlName]
		str := "\x0c\x00" + strings.ReplaceAll(hint, "\n", "\x01") + "\x00"
		table.New[textID*2] = addrString(addr)[0]
		table.New[textID*2+1] = addrString(addr)[1]
		addr += uint16(len(str))
		builder.WriteString(str)
	}

	text.New = []byte(builder.String())

	codeMutables["owl text offsets"] = table
	codeMutables["owl text"] = text
}

// returns an array of owl statue names for the given game (matching those in
// the logic package).
func GetOwlNames(game int) []string {
	var src map[string]byte
	if game == GameSeasons {
		src = seasonsOwls
	} else {
		src = agesOwls
	}

	a := make([]string, len(src))
	i := 0
	for k := range src {
		a[i] = k
		i++
	}

	return a
}
