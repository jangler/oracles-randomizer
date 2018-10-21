package rom

import (
	"bytes"
	"fmt"
)

// collection modes
// i don't know what the difference between the two find modes is
const (
	collectNil        = 0x00 // custom, for shop items
	collectBuySatchel = 0x01
	collectFind0      = 0x02 // flippers, ring box, maku seed, idk
	collectUnderwater = 0x08 // pyramid jewel
	collectFind1      = 0x09
	collectFind2      = 0x0a
	collectAppear1    = 0x19 // d5 boss key
	collectAppear2    = 0x1a // heart containers
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

// SubID returns item sub ID of the treasure.
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
	if t.addr.Offset == 0 {
		return nil
	}

	addr, data := t.addr.FullOffset(), t.Bytes()
	for i := 0; i < 4; i++ {
		b[addr+i] = data[i]
	}
	return nil
}

// Check verifies that the treasure's data matches the given ROM data.
func (t Treasure) Check(b []byte) error {
	addr, data := t.addr.FullOffset(), t.Bytes()
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

// initialized automatically in init() based on contents of item slots
var TreasureIsUnique = map[string]bool{}

// returns true iff a treasure can be lost permanently (i.e. outside of hide
// and seek).
func TreasureCanBeLost(name string) bool {
	switch name {
	case "shop shield L-1", "shield L-2", "star ore", "ribbon",
		"spring banana", "ricky's gloves", "round jewel", "pyramid jewel",
		"square jewel", "x-shapred jewel", "red ore", "blue ore", "hard ore":
		return true
	}
	return false
}
