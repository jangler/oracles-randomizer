package randomizer

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

// treasure interaction spawn + collect modes. bits 0-2 are for collect
// animation; bit 3 sets room flag on collection; bits 4-6 determine how the
// treasure appears; bit 7 is used as the randomizer as a jump table index for
// special cases (it can't appear in the vanilla table).
var collectModes = map[string]byte{
	"touch":               0x0a, // standing and given items
	"poof":                0x1a, // boss HCs
	"drop":                0x29, // SK drops, maku tree, graveyard key
	"chest":               0x38, // most chests, rising animation
	"dive":                0x49, // SK and BK in seasons D4
	"dig":                 0x5a, // star ore, ricky's gloves (ages)
	"delay":               0x68, // map and compass chests
	"diver room":          0x80,
	"poe skip room":       0x81,
	"maku tree (seasons)": 0x82,
	"d4 pool":             0x83,
	"d5 armos":            0x84,
	"maku tree (ages)":    0x80,
	"target carts":        0x81,
	"big bang game":       0x82,
	"lava juice room":     0x83,
}

// data associated with a particular item ID and sub ID.
type treasure struct {
	displayName string // this can change based on ring replacement etc
	id, subid   byte
	addr        address

	// in order, starting at addr
	mode   byte // collection mode
	param  byte // parameter value to use for giveTreasure
	text   byte
	sprite byte
}

// returns a slice of consecutive bytes of treasure data, as they would appear
// in the ROM.
func (t treasure) bytes() []byte {
	return []byte{t.mode, t.param, t.text, t.sprite}
}

// implements `mutate()` from the `mutable` interface.
func (t treasure) mutate(b []byte) {
	// fake treasure
	if t.addr.offset == 0 {
		return
	}

	addr, data := t.addr.fullOffset(), t.bytes()
	for i := 0; i < 4; i++ {
		b[addr+i] = data[i]
	}
}

// implements `check()` from the `mutable` interface.
func (t treasure) check(b []byte) error {
	// fake treasure
	if t.addr.offset == 0 {
		return nil
	}

	addr, data := t.addr.fullOffset(), t.bytes()
	if bytes.Compare(b[addr:addr+4], data) != 0 {
		return fmt.Errorf("expected %x at %x; found %x",
			data, addr, b[addr:addr+4])
	}
	return nil
}

// returns the full offset of the treasure's four-byte entry in the rom.
func getTreasureAddr(b []byte, game int, id, subid byte) address {
	ptr := sora(game, address{0x15, 0x5129}, address{0x16, 0x5332}).(address)

	ptr.offset += uint16(id) * 4
	if b[ptr.fullOffset()]&0x80 != 0 {
		ptr.offset = uint16(b[ptr.fullOffset()+1]) +
			uint16(b[ptr.fullOffset()+2])*0x100
	}
	ptr.offset += uint16(subid) * 4

	return ptr
}

// return a map of treasure names to treasure data. if b is nil, only "static"
// data is loaded.
func loadTreasures(b []byte, game int) map[string]*treasure {
	allRawIds := make(map[string]map[string]uint16)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/treasures.yaml"), allRawIds); err != nil {
		panic(err)
	}

	rawIds := make(map[string]uint16)
	for k, v := range allRawIds["common"] {
		rawIds[k] = v
	}
	for k, v := range allRawIds[gameNames[game]] {
		rawIds[k] = v
	}

	m := make(map[string]*treasure)
	for name, rawId := range rawIds {
		t := &treasure{
			displayName: name,
			id:          byte(rawId >> 8),
			subid:       byte(rawId),
		}

		if b != nil {
			t.addr = getTreasureAddr(b, game, t.id, t.subid)
			t.mode = b[t.addr.fullOffset()]
			t.param = b[t.addr.fullOffset()+1]
			t.text = b[t.addr.fullOffset()+2]
			t.sprite = b[t.addr.fullOffset()+3]
		}

		m[name] = t
	}

	if game == gameSeasons {
		// these treasures don't exist as treasure interactions in the vanilla
		// game, so they're missing some data.
		m["fool's ore"].text = 0x36
		m["fool's ore"].sprite = 0x4a
		m["rare peach stone"].sprite = 0x4e
		m["ribbon"].text = 0x41
		m["ribbon"].sprite = 0x4f
		m["treasure map"].text = 0x6c
		m["treasure map"].sprite = 0x49
		m["member's card"].text = 0x45
		m["member's card"].sprite = 0x48

		// and seasons flutes aren't initially real treasures like ages ones are
		t := m["ricky's flute"]
		t.addr = address{}
		t.param = 0x0b
		t = m["dimitri's flute"]
		t.addr = address{}
		t.param = 0x0c
		t = m["moosh's flute"]
		t.addr = address{}
		t.param = 0x0d

		// give bracelet a level for ages multiworld compatibility
		m["bracelet"].param = 0x01
	} else {
		// give strange flutes identified flute text and palettes
		m["ricky's flute"].text = 0x38
		m["ricky's flute"].sprite = 0x6c
		m["dimitri's flute"].text = 0x39
		m["dimitri's flute"].sprite = 0x6d
		m["moosh's flute"].text = 0x3a
		m["moosh's flute"].sprite = 0x6e
	}

	// add dummy treasures for seed trees
	m["ember tree seeds"] = &treasure{id: 0x00}
	m["scent tree seeds"] = &treasure{id: 0x01}
	m["pegasus tree seeds"] = &treasure{id: 0x02}
	m["gale tree seeds"] = &treasure{id: 0x03}
	m["mystery tree seeds"] = &treasure{id: 0x04}

	return m
}
