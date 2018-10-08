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

var ItemSlots = map[string]*MutableSlot{
	// holodrum
	"lake chest": MutableChest(
		"gasha seed", 0x4f92, 0x00, 0xb8, collectChest, 0xb8),
	"maku tree gift": &MutableSlot{
		Treasure:    Treasures["gnarled key"],
		IDAddrs:     []Addr{{0x15, 0x613a}, {0x09, 0x7e16}},
		SubIDAddrs:  []Addr{{0x15, 0x613d}, {0x09, 0x7e19}},
		group:       0x02,
		room:        0x0b,
		collectMode: collectFall,
		mapCoords:   0xc9,
	},
	"village SW chest": MutableChest(
		"rupees, 20", 0x4f7e, 0x00, 0xf5, collectChest, 0xf5),
	"village SE chest": MutableChest(
		"rupees, 20", 0x4f82, 0x00, 0xf9, collectChest, 0xf9),
	"shovel gift": MutableScriptItem(
		"shovel", 0x6a6c, 0x03, 0xa3, collectFind2, 0x7f),
	"outdoor d2 chest": MutableChest(
		"gasha seed", 0x4f86, 0x00, 0x8e, collectChest, 0x8e),
	"blaino gift": MutableScriptItem(
		"gasha seed", 0x64cc, 0x03, 0xb4, collectFind1, 0x78),
	"floodgate key spot": MutableFind(
		"floodgate key", 0x6281, 0x03, 0xb5, collectFind1, 0x62),
	"square jewel chest": &MutableSlot{
		Treasure:    Treasures["square jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7395}},
		SubIDAddrs:  []Addr{{0x0b, 0x7399}},
		group:       0x04,
		room:        0xfa,
		collectMode: collectChest,
		mapCoords:   0xc2,
	},
	"great moblin chest": MutableChest(
		"piece of heart", 0x4f8e, 0x00, 0x5b, collectChest, 0x5b),
	"master's plaque chest": MutableChest(
		"master's plaque", 0x510a, 0x05, 0xbc, collectChest, 0x2e),
	"diver gift": MutableScriptItem( // addr set at EOB
		"flippers", 0x0000, 0x05, 0xbd, collectNil, 0x2e), // special case
	"spring banana tree": MutableFind(
		"spring banana", 0x66c6, 0x00, 0x0f, collectFind2, 0x0f),
	"dragon key spot": MutableFind(
		"dragon key", 0x62a3, 0x00, 0x1a, collectFind1, 0x1a),
	"pyramid jewel spot": MutableScriptItem(
		"pyramid jewel", 0x734e, 0x07, 0xe5, collectUnderwater, 0x1d),
	"x-shaped jewel chest": MutableChest(
		"x-shaped jewel", 0x4f8a, 0x00, 0xf4, collectChest, 0xf4),
	"round jewel gift": MutableScriptItem(
		"round jewel", 0x7332, 0x03, 0x94, collectFind2, 0xb5),
	"noble sword spot": &MutableSlot{
		Treasure:    Treasures["sword 2"],
		IDAddrs:     []Addr{{0x0b, 0x6418}, {0x0b, 0x641f}},
		SubIDAddrs:  []Addr{{0x0b, 0x6419}, {0x0b, 0x6420}},
		group:       0x00,
		room:        0xc9,
		collectMode: collectFind1,
		mapCoords:   0x40,
	},
	"desert pit": &MutableSlot{
		Treasure:    Treasures["rusty bell"],
		IDAddrs:     []Addr{{0x09, 0x648d}, {0x0b, 0x60b1}},
		SubIDAddrs:  []Addr{{0x09, 0x648c}},
		group:       0x05,
		room:        0xd2,
		collectMode: collectFind2,
		mapCoords:   0xbf,
	},
	"desert chest": MutableChest(
		"blast ring", 0x4f9a, 0x00, 0xff, collectChest, 0xff),
	"western coast chest": MutableChest(
		"rang ring L-1", 0x4f96, 0x00, 0xe3, collectChest, 0xe3),
	"coast house chest": MutableChest(
		"bombs, 10", 0x4fac, 0x03, 0x88, collectChest, 0xd2),
	"water cave chest": MutableChest(
		"octo ring", 0x5081, 0x04, 0xe0, collectChest, 0xb3),
	"mushroom cave chest": MutableChest(
		"quicksand ring", 0x5085, 0x04, 0xe1, collectChest, 0x87),
	"mystery cave chest": MutableChest(
		"moblin ring", 0x50fe, 0x05, 0xb3, collectChest, 0x8e),
	"moblin road chest": MutableChest(
		"rupees, 30", 0x5102, 0x05, 0xb4, collectChest, 0x7d),
	"sunken cave chest": MutableChest(
		"gasha seed", 0x5106, 0x05, 0xb5, collectChest, 0x4f),
	"diver chest": MutableChest( // TODO this shares room w/ diver gift
		"rupees, 50", 0x510e, 0x05, 0xbd, collectChest, 0x2e),
	"dry lake east chest": MutableChest(
		"piece of heart", 0x5112, 0x05, 0xc0, collectChest, 0xaa),
	"goron chest": MutableChest(
		"armor ring L-2", 0x511a, 0x05, 0xc8, collectChest, 0x18),
	"platform chest": MutableChest(
		"rupees, 50", 0x5122, 0x05, 0x0e, collectChest, 0x49),
	"talon cave chest": MutableChest(
		"subrosian ring", 0x511e, 0x05, 0xb6, collectChest, 0x1b),
	"tarm gasha chest": MutableChest(
		"gasha seed", 0x4fa8, 0x03, 0x9b, collectChest, 0x10),
	"moblin cliff chest": MutableChest(
		"gasha seed", 0x5089, 0x04, 0xf7, collectChest, 0xcc),
	"dry lake west chest": &MutableSlot{
		Treasure:    Treasures["rupees, 100"],
		IDAddrs:     []Addr{{0x0b, 0x73a1}},
		SubIDAddrs:  []Addr{{0x0b, 0x73a5}},
		group:       0x04,
		room:        0xfb,
		collectMode: collectChest,
		mapCoords:   0xa7,
	},
	"linked dive chest": &MutableSlot{
		Treasure:    Treasures["gasha seed"],
		IDAddrs:     []Addr{{0x0a, 0x5003}},
		SubIDAddrs:  []Addr{{0x0a, 0x5008}},
		group:       0x05,
		room:        0x12,
		collectMode: collectChest,
		mapCoords:   0x7e,
	},

	// dummy slots for bombs and shield
	"village shop 1": &MutableSlot{
		Treasure:    Treasures["bombs, 10"],
		group:       0x03,
		room:        0xa6,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},
	"village shop 2": &MutableSlot{
		Treasure:    Treasures["shop shield L-1"],
		group:       0x03,
		room:        0xa6,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},

	"village shop 3": &MutableSlot{
		Treasure:    Treasures["strange flute"],
		IDAddrs:     []Addr{{0x08, 0x4ce8}},
		SubIDAddrs:  []Addr{{0x08, 0x4ce9}},
		group:       0x03,
		room:        0xa6,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},
	"member's shop 1": &MutableSlot{
		Treasure:    Treasures["satchel 2"],
		IDAddrs:     []Addr{{0x08, 0x4cce}},
		SubIDAddrs:  []Addr{{0x08, 0x4ccf}},
		group:       0x03,
		room:        0xb0,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},
	"member's shop 2": &MutableSlot{
		Treasure:    Treasures["gasha seed"],
		IDAddrs:     []Addr{{0x08, 0x4cd2}},
		SubIDAddrs:  []Addr{{0x08, 0x4cd3}},
		group:       0x03,
		room:        0xb0,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},
	"member's shop 3": &MutableSlot{
		Treasure:    Treasures["treasure map"],
		IDAddrs:     []Addr{{0x08, 0x4cd8}},
		SubIDAddrs:  []Addr{{0x08, 0x4cd9}},
		group:       0x03,
		room:        0xb0,
		collectMode: collectNil,
		mapCoords:   0xe6,
	},

	// subrosia
	"winter tower": MutableScriptItem(
		"winter", 0x4fc5, 0x05, 0xf2, collectFind1, 0x9a),
	"summer tower": MutableScriptItem(
		"summer", 0x4fb9, 0x05, 0xf8, collectFind1, 0xb0),
	"spring tower": MutableScriptItem(
		"spring", 0x4fb5, 0x05, 0xf5, collectFind1, 0x1e),
	"autumn tower": MutableScriptItem(
		"autumn", 0x4fc1, 0x05, 0xfb, collectFind1, 0xb9),
	"dance hall prize": MutableScriptItem(
		"boomerang 1", 0x6646, 0x03, 0x95, collectFind2, 0x9a),
	"rod gift": &MutableSlot{
		Treasure:    Treasures["rod"],
		IDAddrs:     []Addr{{0x15, 0x70ce}},
		SubIDAddrs:  []Addr{{0x15, 0x70cc}},
		group:       0x03,
		room:        0xac,
		collectMode: collectNil,
		mapCoords:   0x9a,
	},
	"star ore spot": &MutableSlot{ // addrs set dynamically at EOB
		Treasure:    Treasures["star ore"],
		IDAddrs:     []Addr{{0x08, 0x0000}},
		SubIDAddrs:  []Addr{{0x08, 0x0000}},
		group:       0x01,
		room:        0x66,
		collectMode: collectDig,
		mapCoords:   0xb0,
	},
	"blue ore chest": MutableChest(
		"blue ore", 0x4f9f, 0x01, 0x41, collectChest, 0x1e),
	"red ore chest": MutableChest(
		"red ore", 0x4fa3, 0x01, 0x58, collectChest, 0xb9),
	"non-rosa gasha chest": MutableChest(
		"gasha seed", 0x5095, 0x04, 0xf1, collectChest, 0x25),
	"rosa gasha chest": MutableChest(
		"gasha seed", 0x5116, 0x05, 0xc6, collectChest, 0xb0),
	"subrosian market 1": &MutableSlot{
		Treasure:    Treasures["ribbon"],
		IDAddrs:     []Addr{{0x09, 0x77da}},
		SubIDAddrs:  []Addr{{0x09, 0x77db}},
		group:       0x03,
		room:        0xa0,
		collectMode: collectNil,
		mapCoords:   0xb0,
	},
	"subrosian market 2": &MutableSlot{
		Treasure:    Treasures["rare peach stone"],
		IDAddrs:     []Addr{{0x09, 0x77e2}},
		SubIDAddrs:  []Addr{{0x09, 0x77e3}},
		group:       0x03,
		room:        0xa0,
		collectMode: collectNil,
		mapCoords:   0xb0,
	},
	"subrosian market 5": &MutableSlot{
		Treasure:    Treasures["member's card"],
		IDAddrs:     []Addr{{0x09, 0x77f4}},
		SubIDAddrs:  []Addr{{0x09, 0x77f5}},
		group:       0x03,
		room:        0xa0,
		collectMode: collectNil,
		mapCoords:   0xb0,
	},
	"hard ore slot": &MutableSlot{ // addrs set dynamically at EOB
		Treasure:    Treasures["hard ore"],
		IDAddrs:     []Addr{{0x15, 0x0000}, {0x09, 0x66eb}},
		SubIDAddrs:  []Addr{{0x15, 0x0000}, {0x09, 0x66ea}},
		group:       0x03,
		room:        0x8e,
		collectMode: collectFind2,
		mapCoords:   0xb9,
	},
	"iron shield gift": &MutableSlot{
		Treasure:    Treasures["shield L-2"],
		IDAddrs:     []Addr{{0x15, 0x62be}},
		SubIDAddrs:  []Addr{{0x15, 0x62b4}},
		group:       0x03,
		room:        0x97,
		collectMode: collectFind2,
		mapCoords:   0x25,
	},

	// hero's cave
	"d0 sword chest": &MutableSlot{
		Treasure:    Treasures["sword 1"],
		IDAddrs:     []Addr{{0x0a, 0x7b90}},
		ParamAddrs:  []Addr{{0x0a, 0x7b92}},
		TextAddrs:   []Addr{{0x0a, 0x7b9c}},
		GfxAddrs:    []Addr{{0x3f, 0x6676}},
		group:       0x04,
		room:        0x06,
		collectMode: collectChest,
		mapCoords:   0xd4,
	},
	"d0 rupee chest": MutableChest(
		"rupees, 30", 0x4fb5, 0x04, 0x05, collectChest, 0xd4),

	// d1
	"d1 satchel spot": MutableFind(
		"satchel 1", 0x66b1, 0x06, 0x09, collectFind2, 0x96),
	"d1 gasha chest": MutableChest(
		"gasha seed", 0x4fbd, 0x04, 0x0d, collectChest, 0x96),
	"d1 bomb chest": MutableChest(
		"bombs, 10", 0x4fc5, 0x04, 0x10, collectChest, 0x96),
	"d1 ring chest": MutableChest(
		"discovery ring", 0x4fd1, 0x04, 0x17, collectChest, 0x96),
	"d1 compass chest": MutableChest(
		"compass", 0x4fc1, 0x04, 0x0f, collectChest2, 0x96),
	"d1 map chest": MutableChest(
		"dungeon map", 0x4fd5, 0x04, 0x19, collectChest2, 0x96),
	"d1 boss key chest": MutableChest(
		"d1 boss key", 0x4fcd, 0x04, 0x14, collectChest, 0x96),

	// d2
	"d2 bracelet chest": MutableChest(
		"bracelet", 0x4fe1, 0x04, 0x2a, collectChest, 0x8d),
	"d2 10-rupee chest": MutableChest(
		"rupees, 10", 0x4fd9, 0x04, 0x1f, collectChest, 0x8d),
	"d2 5-rupee chest": MutableChest(
		"rupees, 5", 0x4ff5, 0x04, 0x38, collectChest, 0x8d),
	"d2 map chest": MutableChest(
		"dungeon map", 0x4fe5, 0x04, 0x2b, collectChest2, 0x8d),
	"d2 compass chest": MutableChest(
		"compass", 0x4ff1, 0x04, 0x36, collectChest2, 0x8d),
	"d2 boss key chest": MutableChest(
		"d2 boss key", 0x4fdd, 0x04, 0x24, collectChest, 0x8d),

	// d3
	"d3 feather chest": MutableChest(
		"feather 1", 0x5015, 0x04, 0x50, collectChest, 0x60),
	"d3 rupee chest": MutableChest(
		"rupees, 30", 0x4ff9, 0x04, 0x41, collectChest, 0x60),
	"d3 gasha chest": MutableChest(
		"gasha seed", 0x5001, 0x04, 0x44, collectChest, 0x60),
	"d3 bomb chest": MutableChest(
		"bombs, 10", 0x5019, 0x04, 0x54, collectChest, 0x60),
	"d3 compass chest": MutableChest(
		"compass", 0x5009, 0x04, 0x4d, collectChest2, 0x60),
	"d3 map chest": MutableChest(
		"dungeon map", 0x5011, 0x04, 0x51, collectChest2, 0x60),
	"d3 boss key chest": MutableChest(
		"d3 boss key", 0x4ffd, 0x04, 0x46, collectChest, 0x60),

	// d4
	"d4 slingshot chest": MutableChest(
		"slingshot 1", 0x502d, 0x04, 0x73, collectChest, 0x1d),
	"d4 bomb chest": MutableChest(
		"bombs, 10", 0x5031, 0x04, 0x7f, collectChest, 0x1d),
	"d4 map chest": MutableChest(
		"dungeon map", 0x5025, 0x04, 0x69, collectChest2, 0x1d),
	"d4 compass chest": MutableChest(
		"compass", 0x5035, 0x04, 0x83, collectChest2, 0x1d),
	"d4 boss key spot": MutableScriptItem(
		"d4 boss key", 0x4c0b, 0x04, 0x6c, collectDive, 0x1d),

	// d5
	"d5 magnet gloves chest": MutableChest(
		"magnet gloves", 0x503d, 0x04, 0x89, collectChest, 0x89),
	"d5 rupee chest": MutableChest(
		"rupees, 100", 0x5041, 0x04, 0x97, collectChest, 0x8a),
	"d5 map chest": MutableChest(
		"dungeon map", 0x5039, 0x04, 0x8f, collectChest2, 0x8f),
	"d5 compass chest": MutableChest(
		"compass", 0x5049, 0x04, 0x9d, collectChest2, 0x8a),
	"d5 boss key spot": MutableScriptItem(
		"d5 boss key", 0x4c22, 0x06, 0x8b, collectFind2, 0x8a),

	// d6
	"d6 boomerang chest": MutableChest(
		"boomerang 2", 0x507d, 0x04, 0xd0, collectChest, 0x00),
	"d6 rupee chest A": MutableChest(
		"rupees, 10", 0x505d, 0x04, 0xaf, collectChest, 0x00),
	"d6 rupee chest B": MutableChest(
		"rupees, 5", 0x5065, 0x04, 0xb3, collectChest, 0x00),
	"d6 bomb chest": MutableChest(
		"bombs, 10", 0x5069, 0x04, 0xbf, collectChest, 0x00),
	"d6 rupee chest C": MutableChest(
		"rupees, 5", 0x5075, 0x04, 0xc3, collectChest, 0x00),
	"d6 compass chest": MutableChest(
		"compass", 0x5059, 0x04, 0xad, collectChest2, 0x00),
	"d6 map chest": MutableChest(
		"dungeon map", 0x5061, 0x04, 0xb0, collectChest2, 0x00),
	"d6 boss key chest": MutableChest(
		"d6 boss key", 0x5079, 0x04, 0xc4, collectChest, 0x00),

	// d7
	"d7 cape chest": MutableChest(
		"feather 2", 0x509e, 0x05, 0x44, collectChest, 0xd0),
	"d7 rupee chest": MutableChest(
		"rupees, 1", 0x509a, 0x05, 0x43, collectChest, 0xd0),
	"d7 ring chest": MutableChest(
		"power ring L-1", 0x50b6, 0x05, 0x5a, collectChest, 0xd0),
	"d7 compass chest": MutableChest(
		"compass", 0x50aa, 0x05, 0x52, collectChest2, 0xd0),
	"d7 map chest": MutableChest(
		"dungeon map", 0x50b2, 0x05, 0x59, collectChest2, 0xd0),
	"d7 boss key chest": MutableChest(
		"d7 boss key", 0x50a6, 0x05, 0x48, collectChest, 0xd0),

	// d8
	"d8 HSS chest": MutableChest(
		"slingshot 2", 0x50da, 0x05, 0x8d, collectChest, 0x04),
	"d8 SW lava chest": MutableChest(
		"bombs, 10", 0x50ba, 0x05, 0x6a, collectChest, 0x04),
	"d8 ring chest": MutableChest(
		"steadfast ring", 0x50c6, 0x05, 0x7d, collectChest, 0x04),
	"d8 compass chest": MutableChest(
		"compass", 0x50d2, 0x05, 0x8b, collectChest2, 0x04),
	"d8 map chest": MutableChest(
		"dungeon map", 0x50de, 0x05, 0x8e, collectChest2, 0x04),
	"d8 boss key chest": MutableChest(
		"d8 boss key", 0x50ca, 0x05, 0x80, collectChest, 0x04),

	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	// "ring box L-2 gift": MutableScriptItem("ring box L-2", 0x5c18),

	// these are "fake" item slots in that they don't slot real treasures
	"ember tree": &MutableSlot{
		Treasure: Treasures["ember tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x64ce}},
	},
	"mystery tree": &MutableSlot{
		Treasure: Treasures["mystery tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x67e0}},
	},
	"scent tree": &MutableSlot{
		Treasure: Treasures["scent tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x685f}},
	},
	"pegasus tree": &MutableSlot{
		Treasure: Treasures["pegasus tree seeds"],
		IDAddrs:  []Addr{{0x11, 0x6873}},
	},
	"sunken gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 1"],
		IDAddrs:  []Addr{{0x11, 0x69b3}},
	},
	"tarm gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 2"],
		IDAddrs:  []Addr{{0x11, 0x6a49}},
	},
}
