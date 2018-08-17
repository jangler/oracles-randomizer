package rom

import (
	"bytes"
	"fmt"
)

// collection modes
// i don't know what the difference between the two find modes are
const (
	CollectGoronGift  = 0x02 // for the L-2 ring box only ??
	CollectUnderwater = 0x08
	CollectFind1      = 0x09
	CollectFind2      = 0x0a
	CollectFall       = 0x29
	CollectChest      = 0x38
	CollectDig        = 0x5a
)

// A Treasure is data associated with a particular item ID and sub ID.
type Treasure struct {
	id, subID byte
	addr      uint16 // bank 15, value of hl at $15:466b, minus one

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

func (t Treasure) CollectMode() byte {
	return t.mode
}

// RealAddr returns the total offset of the treasure data in the ROM.
func (t Treasure) RealAddr() int {
	return (&Addr{0x15, t.addr}).FullOffset()
}

// Bytes returns a slice of consecutive bytes of treasure data, as they would
// appear in the ROM.
func (t Treasure) Bytes() []byte {
	return []byte{t.mode, t.param, t.text, t.sprite}
}

// Mutate replaces the associated treasure in the given ROM data with this one.
func (t Treasure) Mutate(b []byte) error {
	// fake treasure
	if t.addr == 0 {
		return nil
	}

	addr, data := t.RealAddr(), t.Bytes()
	for i := 0; i < 4; i++ {
		b[addr+i] = data[i]
	}
	return nil
}

// Check verifies that the treasure's data matches the given ROM data.
func (t Treasure) Check(b []byte) error {
	addr, data := t.RealAddr(), t.Bytes()
	if bytes.Compare(b[addr:addr+4], data) != 0 {
		return fmt.Errorf("expected %x at %x; found %x",
			data, addr, b[addr:addr+4])
	}
	return nil
}

// Treasures maps item names to associated treasure data.
var Treasures = map[string]*Treasure{
	"shield L-1":    &Treasure{0x01, 0x00, 0x5700, 0x0a, 0x01, 0x1f, 0x13},
	"bombs":         &Treasure{0x03, 0x00, 0x570c, 0x38, 0x10, 0x4d, 0x05},
	"sword L-1":     &Treasure{0x05, 0x00, 0x571c, 0x38, 0x01, 0x1c, 0x10},
	"sword L-2":     &Treasure{0x05, 0x01, 0x5720, 0x09, 0x02, 0x1d, 0x11},
	"boomerang L-1": &Treasure{0x06, 0x00, 0x5734, 0x0a, 0x01, 0x22, 0x1c},
	"boomerang L-2": &Treasure{0x06, 0x01, 0x5738, 0x38, 0x02, 0x23, 0x1d},
	"rod":           &Treasure{0x07, 0x00, 0x573c, 0x38, 0x07, 0x0a, 0x1e},
	"magnet gloves": &Treasure{0x08, 0x00, 0x558c, 0x38, 0x00, 0x30, 0x18},
	"bombchus":      &Treasure{0x0d, 0x00, 0x5760, 0x0a, 0x10, 0x32, 0x24},
	"strange flute": &Treasure{0x0e, 0x00, 0x55a4, 0x0a, 0x0c, 0x3b, 0x23},
	"slingshot L-1": &Treasure{0x13, 0x00, 0x5768, 0x38, 0x01, 0x2e, 0x21},
	"slingshot L-2": &Treasure{0x13, 0x01, 0x576c, 0x38, 0x02, 0x2f, 0x22},
	"shovel":        &Treasure{0x15, 0x00, 0x55c0, 0x0a, 0x00, 0x25, 0x1b},
	"bracelet":      &Treasure{0x16, 0x00, 0x55c4, 0x38, 0x00, 0x26, 0x19},
	"feather L-1":   &Treasure{0x17, 0x00, 0x5770, 0x38, 0x01, 0x27, 0x16},
	"feather L-2":   &Treasure{0x17, 0x01, 0x5774, 0x38, 0x02, 0x28, 0x17},
	"satchel":       &Treasure{0x19, 0x00, 0x56f8, 0x0a, 0x01, 0x2d, 0x20},
	"fool's ore":    &Treasure{0x1e, 0x00, 0x55e4, 0x00, 0x00, 0xff, 0x1a},
	"flippers":      &Treasure{0x2e, 0x00, 0x5624, 0x0a, 0x00, 0x31, 0x31},

	// seasons are obtained by giving the rod of seasons with differet sub-IDs
	"winter": &Treasure{0x07, 0x05, 0x5750, 0x09, 0x03, 0x0a, 0x1e},
	"summer": &Treasure{0x07, 0x03, 0x5748, 0x09, 0x01, 0x0b, 0x1e},
	"spring": &Treasure{0x07, 0x02, 0x5744, 0x09, 0x00, 0x0d, 0x1e},
	"autumn": &Treasure{0x07, 0x04, 0x574c, 0x09, 0x02, 0x0c, 0x1e},

	"small key": &Treasure{0x30, 0x03, 0x584c, 0x38, 0x01, 0x1a, 0x42},
	"boss key":  &Treasure{0x31, 0x03, 0x585c, 0x38, 0x00, 0x1b, 0x43},
	"compass":   &Treasure{0x32, 0x02, 0x5868, 0x68, 0x00, 0x19, 0x41},
	"map":       &Treasure{0x33, 0x02, 0x5874, 0x68, 0x00, 0x18, 0x40},

	"gnarled key":     &Treasure{0x42, 0x00, 0x58a8, 0x29, 0x00, 0x42, 0x44},
	"ricky's gloves":  &Treasure{0x48, 0x00, 0x568c, 0x09, 0x01, 0x67, 0x55},
	"floodgate key":   &Treasure{0x43, 0x00, 0x5678, 0x09, 0x00, 0x43, 0x45},
	"star ore":        &Treasure{0x45, 0x00, 0x5680, 0x5a, 0x00, 0x40, 0x57},
	"square jewel":    &Treasure{0x4e, 0x00, 0x56a4, 0x38, 0x00, 0x48, 0x38},
	"master's plaque": &Treasure{0x54, 0x00, 0x56bc, 0x38, 0x00, 0x70, 0x26},
	"spring banana":   &Treasure{0x47, 0x00, 0x5688, 0x0a, 0x00, 0x66, 0x54},
	"dragon key":      &Treasure{0x44, 0x00, 0x567c, 0x09, 0x00, 0x44, 0x46},
	"pyramid jewel":   &Treasure{0x4d, 0x00, 0x58bc, 0x08, 0x00, 0x4a, 0x37},
	"x-shaped jewel":  &Treasure{0x4f, 0x00, 0x56a8, 0x38, 0x00, 0x49, 0x39},
	"round jewel":     &Treasure{0x4c, 0x00, 0x569c, 0x0a, 0x00, 0x47, 0x36},
	"rusty bell":      &Treasure{0x4a, 0x00, 0x58b0, 0x0a, 0x00, 0x55, 0x5b},
	"red ore":         &Treasure{0x50, 0x00, 0x56ac, 0x38, 0x00, 0x3f, 0x59},
	"blue ore":        &Treasure{0x51, 0x00, 0x56b0, 0x38, 0x00, 0x3e, 0x58},
	"ring box L-2":    &Treasure{0x2c, 0x02, 0x57f0, 0x02, 0x03, 0x34, 0x35},

	"bombs, 10":      &Treasure{0x03, 0x00, 0x570c, 0x38, 0x10, 0x4d, 0x05},
	"gasha seed":     &Treasure{0x34, 0x01, 0x5784, 0x38, 0x01, 0x4b, 0x0d},
	"rupees, 1":      &Treasure{0x28, 0x00, 0x5798, 0x38, 0x01, 0x01, 0x28},
	"rupees, 5":      &Treasure{0x28, 0x01, 0x579c, 0x38, 0x03, 0x02, 0x29},
	"rupees, 10":     &Treasure{0x28, 0x02, 0x57a0, 0x38, 0x04, 0x03, 0x2a},
	"rupees, 20":     &Treasure{0x28, 0x03, 0x57a4, 0x38, 0x05, 0x04, 0x2b},
	"rupees, 30":     &Treasure{0x28, 0x04, 0x57a8, 0x38, 0x07, 0x05, 0x2b},
	"rupees, 50":     &Treasure{0x28, 0x05, 0x57ac, 0x38, 0x0b, 0x06, 0x2c},
	"rupees, 100":    &Treasure{0x28, 0x06, 0x57b0, 0x38, 0x0c, 0x07, 0x2d},
	"piece of heart": &Treasure{0x2b, 0x01, 0x57d4, 0x38, 0x01, 0x17, 0x3a},
	"discovery ring": &Treasure{0x2d, 0x04, 0x580c, 0x38, 0x28, 0x54, 0x0e},
	"moblin ring":    &Treasure{0x2d, 0x05, 0x5810, 0x38, 0x2b, 0x54, 0x0e},
	"steadfast ring": &Treasure{0x2d, 0x06, 0x5814, 0x38, 0x10, 0x54, 0x0e},
	"rang ring L-1":  &Treasure{0x2d, 0x07, 0x5818, 0x38, 0x0c, 0x54, 0x0e},
	"blast ring":     &Treasure{0x2d, 0x08, 0x581c, 0x38, 0x0d, 0x54, 0x0e},
	"octo ring":      &Treasure{0x2d, 0x09, 0x5820, 0x38, 0x2a, 0x54, 0x0e},
	"quicksand ring": &Treasure{0x2d, 0x0a, 0x5824, 0x38, 0x23, 0x54, 0x0e},
	"armor ring L-2": &Treasure{0x2d, 0x0b, 0x5828, 0x38, 0x05, 0x54, 0x0e},
	"power ring L-1": &Treasure{0x2d, 0x0e, 0x5834, 0x38, 0x01, 0x54, 0x0e},
	"subrosian ring": &Treasure{0x2d, 0x10, 0x583c, 0x38, 0x2d, 0x54, 0x0e},

	// these seeds are "fake" treasures. real treasures corresponding to each
	// type of seed exist, but those can't be used for changing which tree
	// yields which seeds.
	"ember tree seeds":   &Treasure{id: 0x00},
	"mystery tree seeds": &Treasure{id: 0x01},
	"scent tree seeds":   &Treasure{id: 0x02},
	"pegasus tree seeds": &Treasure{id: 0x03},
	"gale tree seeds 1":  &Treasure{id: 0x04},
	"gale tree seeds 2":  &Treasure{id: 0x05},
}

var seedIndexByTreeID = []byte{0, 4, 1, 2, 3, 3}

// reverse lookup the treasure name; returns empty string if not found. this
// ignores fake seed treasures.
func treasureNameFromIDs(id, subID byte) string {
	for k, v := range Treasures {
		if v.addr != 0 && v.id == id && v.subID == subID {
			return k
		}
	}
	return ""
}

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

// CanSlotOutsideChest is a map indicating whether an item can be given by an
// NPC, found on the ground, etc as opposed to being found in a chest. This
// essentially breaks down to whether an item is unique, as in there's only one
// of them and therefore its collect mode can be safely altered.
//
// All the entries in this map are going to be true since items not in the map
// with evaluate to false anyway.
var CanSlotOutsideChest = map[string]bool{
	// equip items
	"sword L-1":     true,
	"sword L-2":     true,
	"boomerang L-1": true,
	"boomerang L-2": true,
	"winter":        true,
	"summer":        true,
	"spring":        true,
	"autumn":        true,
	"magnet gloves": true,
	"bombchus":      true,
	"slingshot L-1": true,
	"slingshot L-2": true,
	"shovel":        true,
	"bracelet":      true,
	"feather L-1":   true,
	"feather L-2":   true,
	"satchel":       true,
	"fool's ore":    true,

	// collection items
	"ring box L-2":    true,
	"flippers":        true,
	"gnarled key":     true,
	"floodgate key":   true,
	"dragon key":      true,
	"star ore":        true,
	"spring banana":   true,
	"ricky's gloves":  true,
	"rusty bell":      true,
	"round jewel":     true,
	"pyramid jewel":   true,
	"square jewel":    true,
	"x-shaped jewel":  true,
	"master's plaque": true,

	// seeds
	"ember tree seeds":   true,
	"mystery tree seeds": true,
	"scent tree seeds":   true,
	"pegasus tree seeds": true,
	"gale tree seeds 1":  true,
	"gale tree seeds 2":  true,

	// rings
	"discovery ring": true,
	"moblin ring":    true,
	"steadfast ring": true,
	"rang ring L-1":  true,
	"blast ring":     true,
	"octo ring":      true,
	"quicksand ring": true,
	"armor ring L-2": true,
	"power ring L-1": true,
	"subrosian ring": true,
}
