package rom

import (
	"fmt"
)

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure *Treasure

	treasureName             string
	idAddrs, subIDAddrs      []Addr
	group, room, collectMode byte
	mapCoords                byte // overworld map coords, yx
}

// Mutate replaces the given IDs, subIDs, and other applicable data in the ROM.
func (ms *MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.idAddrs {
		b[addr.fullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.subIDAddrs {
		b[addr.fullOffset()] = ms.Treasure.subID
	}

	return ms.Treasure.Mutate(b)
}

// helper function for MutableSlot.Check
func check(b []byte, addr Addr, value byte) error {
	if b[addr.fullOffset()] != value {
		return fmt.Errorf("expected %x at %x; found %x",
			value, addr.fullOffset(), b[addr.fullOffset()])
	}
	return nil
}

// Check verifies that the slot's data matches the given ROM data.
func (ms *MutableSlot) Check(b []byte) error {
	// skip zero addresses
	if len(ms.idAddrs) == 0 || ms.idAddrs[0].offset == 0 {
		return nil
	}

	// only check ID addresses, since situational variants and progressive
	// items mess with everything else.
	for _, addr := range ms.idAddrs {
		if err := check(b, addr, ms.Treasure.id); err != nil {
			return err
		}
	}

	return nil
}

// basicSlot constucts a MutableSlot from a treasure name, bank number, and an
// address for each its ID and sub-ID. Most slots fit this pattern.
func basicSlot(treasure string, bank byte, idOffset, subIDOffset uint16,
	group, room, mode, coords byte) *MutableSlot {
	return &MutableSlot{
		treasureName: treasure,
		idAddrs:      []Addr{{bank, idOffset}},
		subIDAddrs:   []Addr{{bank, subIDOffset}},
		group:        group,
		room:         room,
		collectMode:  mode,
		mapCoords:    coords,
	}
}

// keyDropSlot constructs a MutableSlot for a small key drop. the mutable
// itself is a dummy and does not have an address; the data is used to
// construct a table of small key drops.
func keyDropSlot(treasure string, group, room, coords byte) *MutableSlot {
	return &MutableSlot{
		treasureName: treasure,
		group:        group,
		room:         room,
		collectMode:  collectFall,
		mapCoords:    coords,
	}
}

var ItemSlots map[string]*MutableSlot
