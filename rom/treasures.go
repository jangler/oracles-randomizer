package rom

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v2"
)

// collection modes
// i don't know what the difference between the two find modes is
const (
	collectNil        = 0x00 // custom, for shop items
	collectBuySatchel = 0x01
	collectFind0      = 0x02 // sets room flag 6 instead of 5
	collectUnderwater = 0x08 // pyramid jewel
	collectFind1      = 0x09
	collectFind2      = 0x0a
	collectEnemyDrop  = 0x18 // stalfos in ages D8
	collectAppear1    = 0x19 // seasons D5 boss key
	collectAppear2    = 0x1a // heart containers
	collectKeyFall    = 0x28 // no fanfare
	collectFall       = 0x29
	collectChest      = 0x38 // most chests
	collectDive       = 0x49
	collectChest2     = 0x68 // map and compass
	collectDigPile    = 0x51
	collectDig        = 0x5a
)

// A Treasure is data associated with a particular item ID and sub ID.
type Treasure struct {
	id, subID byte
	addr      Addr

	// in order, starting at addr
	mode   byte // collection mode
	param  byte // parameter value to use for giveTreasure
	text   byte
	sprite byte
}

// ID returns the item ID of the treasure.
func (t Treasure) ID() byte {
	return t.id
}

// SubID returns the item sub ID of the treasure.
func (t Treasure) SubID() byte {
	return t.subID
}

// Bytes returns a slice of consecutive bytes of treasure data, as they would
// appear in the ROM.
func (t Treasure) Bytes() []byte {
	return []byte{t.mode, t.param, t.text, t.sprite}
}

// Mutate replaces the associated treasure in the given ROM data with this one.
func (t Treasure) Mutate(b []byte) error {
	// fake treasure
	if t.addr.offset == 0 {
		return nil
	}

	addr, data := t.addr.fullOffset(), t.Bytes()
	for i := 0; i < 4; i++ {
		b[addr+i] = data[i]
	}
	return nil
}

// Check verifies that the treasure's data matches the given ROM data.
func (t Treasure) Check(b []byte) error {
	addr, data := t.addr.fullOffset(), t.Bytes()
	if bytes.Compare(b[addr:addr+4], data) != 0 {
		return fmt.Errorf("expected %x at %x; found %x",
			data, addr, b[addr:addr+4])
	}
	return nil
}

// Treasures maps item names to associated treasure data.
var Treasures map[string]*Treasure

// FindTreasureName does a reverse lookup of the treasure in the map to return
// its name. It returns an empty string if not found.
func FindTreasureName(t *Treasure) string {
	for k, v := range Treasures {
		if v == t {
			return k
		}
	}
	return ""
}

// returns the full offset of the treasure's four-byte entry in the rom.
func getTreasureAddr(b []byte, game int, id, subid byte) Addr {
	var ptr Addr
	if game == GameSeasons {
		ptr = Addr{0x15, 0x5129}
	} else {
		ptr = Addr{0x16, 0x5332}
	}

	ptr.offset += uint16(id) * 4
	if b[ptr.fullOffset()]&0x80 != 0 {
		ptr.offset = uint16(b[ptr.fullOffset()+1]) +
			uint16(b[ptr.fullOffset()+2])*0x100
	}
	ptr.offset += uint16(subid) * 4

	return ptr
}

// loads treasure data from the rom into the given map.
func loadTreasures(b []byte, game int, m map[string]*Treasure) {
	rawIds := make(map[string]uint16)

	if err := yaml.Unmarshal(
		FSMustByte(false, "/rom/treasures.yaml"), rawIds); err != nil {
		panic(err)
	}

	for name, rawId := range rawIds {
		if m[name] != nil {
			panic("duplicate treasure name: " + name)
		}

		t := &Treasure{
			id:    byte(rawId >> 8),
			subID: byte(rawId),
		}
		t.addr = getTreasureAddr(b, game, t.id, t.subID)
		t.mode = b[t.addr.fullOffset()]
		t.param = b[t.addr.fullOffset()+1]
		t.text = b[t.addr.fullOffset()+2]
		t.sprite = b[t.addr.fullOffset()+3]

		m[name] = t
	}

	if game == GameSeasons {
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

		// and seasons flutes aren't real treasures like ages ones are
		t := m["ricky's flute"]
		t.param = 0x0b
		t.text = 0x38
		t.sprite = 0x23
		t = m["dimitri's flute"]
		t.subID = 0x00
		t.addr = m["ricky's flute"].addr
		t.param = 0x0c
		t.text = 0x39
		t.sprite = 0x23
		t = m["moosh's flute"]
		t.subID = 0x00
		t.addr = m["ricky's flute"].addr
		t.param = 0x0d
		t.text = 0x3a
		t.sprite = 0x23

		// need to set this before Verify() is called
		varMutables["above d7 zol button"].(*MutableSlot).Treasure =
			m["d7 small key"]
	} else {
		// give strange flute ricky's flute text
		m["ricky's flute"].text = 0x38
	}

	// add dummy treasures for seed trees
	m["ember tree seeds"] = &Treasure{id: 0x00}
	m["scent tree seeds"] = &Treasure{id: 0x01}
	m["pegasus tree seeds"] = &Treasure{id: 0x02}
	m["gale tree seeds"] = &Treasure{id: 0x03}
	m["mystery tree seeds"] = &Treasure{id: 0x04}
}
