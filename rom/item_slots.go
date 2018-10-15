package rom

import (
	"fmt"
)

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure                 *Treasure
	IDAddrs, SubIDAddrs      []Addr
	ParamAddrs, TextAddrs    []Addr
	GfxAddrs                 []Addr
	group, room, collectMode byte
	mapCoords                byte // overworld map coords, yx
}

// Mutate replaces the given IDs and subIDs in the given ROM data, and changes
// the associated treasure's collection mode as appropriate.
func (ms *MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.SubIDAddrs {
		b[addr.FullOffset()] = ms.Treasure.subID
	}
	for _, addr := range ms.ParamAddrs {
		b[addr.FullOffset()] = ms.Treasure.param
	}
	for _, addr := range ms.TextAddrs {
		b[addr.FullOffset()] = ms.Treasure.text
	}
	for _, addr := range ms.GfxAddrs {
		gfx := itemGfx[FindTreasureName(ms.Treasure)]
		for i := 0; i < 3; i++ {
			b[addr.FullOffset()+i] = byte(gfx >> (8 * uint(2-i)))
		}
	}

	return ms.Treasure.Mutate(b)
}

// helper function for MutableSlot.Check
func check(b []byte, addr Addr, value byte) error {
	if b[addr.FullOffset()] != value {
		return fmt.Errorf("expected %x at %x; found %x",
			value, addr.FullOffset(), b[addr.FullOffset()])
	}
	return nil
}

// Check verifies that the slot's data matches the given ROM data.
func (ms *MutableSlot) Check(b []byte) error {
	for _, addr := range ms.IDAddrs {
		if err := check(b, addr, ms.Treasure.id); err != nil {
			return err
		}
	}
	for _, addr := range ms.SubIDAddrs {
		if err := check(b, addr, ms.Treasure.subID); err != nil {
			return err
		}
	}
	for _, addr := range ms.ParamAddrs {
		if err := check(b, addr, ms.Treasure.param); err != nil {
			return err
		}
	}
	for _, addr := range ms.TextAddrs {
		if err := check(b, addr, ms.Treasure.text); err != nil {
			return err
		}
	}
	for _, addr := range ms.GfxAddrs {
		gfx := itemGfx[FindTreasureName(ms.Treasure)]
		for i := uint16(0); i < 3; i++ {
			addr := Addr{addr.Bank, addr.Offset + i}
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

// BasicSlot constucts a MutableSlot from a treasure name, bank number, and an
// address for each its ID and sub-ID. Most slots fit this pattern.
func BasicSlot(treasure string, bank byte, idOffset, subIDOffset uint16,
	group, room, mode, coords byte) *MutableSlot {
	return &MutableSlot{
		Treasure:    Treasures[treasure],
		IDAddrs:     []Addr{{bank, idOffset}},
		SubIDAddrs:  []Addr{{bank, subIDOffset}},
		group:       group,
		room:        room,
		collectMode: mode,
		mapCoords:   coords,
	}
}

// MutableChest constructs a MutableSlot from a treasure name and an address in
// bank $15, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively (?) to chests.
func MutableChest(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return BasicSlot(treasure, 0x15, addr, addr+1, group, room, mode, coords)
}

// MutableScriptItem constructs a MutableSlot from a treasure name and an
// address in bank $0b, where the ID and sub-ID are two consecutive bytes at
// that address.  This applies to most items given by NPCs.
func MutableScriptItem(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return BasicSlot(treasure, 0x0b, addr, addr+1, group, room, mode, coords)
}

// MutableFind constructs a MutableSlot from a treasure name and an address in
// bank $09, where the sub-ID and ID (in that order) are two consecutive bytes
// at that address. This applies to most items that are found lying around.
func MutableFind(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return BasicSlot(treasure, 0x09, addr+1, addr, group, room, mode, coords)
}

var ItemSlots map[string]*MutableSlot
