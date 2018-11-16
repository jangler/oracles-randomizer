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
	paramAddrs, textAddrs    []Addr
	gfxAddrs                 []Addr
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
	for _, addr := range ms.paramAddrs {
		b[addr.fullOffset()] = ms.Treasure.param
	}
	for _, addr := range ms.textAddrs {
		b[addr.fullOffset()] = ms.Treasure.text
	}
	for _, addr := range ms.gfxAddrs {
		gfx := itemGfx[FindTreasureName(ms.Treasure)]
		for i := 0; i < 3; i++ {
			b[addr.fullOffset()+i] = byte(gfx >> (8 * uint(2-i)))
		}
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
	for _, addr := range ms.idAddrs {
		if err := check(b, addr, ms.Treasure.id); err != nil {
			return err
		}
	}
	for _, addr := range ms.subIDAddrs {
		if err := check(b, addr, ms.Treasure.subID); err != nil {
			return err
		}
	}
	for _, addr := range ms.paramAddrs {
		if err := check(b, addr, ms.Treasure.param); err != nil {
			return err
		}
	}
	for _, addr := range ms.textAddrs {
		if err := check(b, addr, ms.Treasure.text); err != nil {
			return err
		}
	}
	for _, addr := range ms.gfxAddrs {
		gfx := itemGfx[FindTreasureName(ms.Treasure)]
		for i := uint16(0); i < 3; i++ {
			addr := Addr{addr.bank, addr.offset + i}
			if err := check(b, addr, byte(gfx>>(8*(2-i)))); err != nil {
				return err
			}
		}
	}
	if ms.collectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.collectMode, ms.Treasure.mode)
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

var ItemSlots map[string]*MutableSlot
