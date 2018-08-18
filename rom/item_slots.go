package rom

import (
	"fmt"
)

const (
	enChestOffset = -0x443
	enGiftOffset  = -2
	enFindOffset  = 0x17
)

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure                        *Treasure
	IDAddrs, SubIDAddrs, ParamAddrs []Addr
	CollectMode                     byte
}

// Mutate replaces the given IDs and subIDs in the given ROM data, and changes
// the associated treasure's collection mode as appropriate.
func (ms *MutableSlot) Mutate(b []byte) error {
	en := isEn(b)
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset(en)] = ms.Treasure.id
	}
	for _, addr := range ms.SubIDAddrs {
		b[addr.FullOffset(en)] = ms.Treasure.subID
	}
	for _, addr := range ms.ParamAddrs {
		b[addr.FullOffset(en)] = ms.Treasure.param
	}
	ms.Treasure.mode = ms.CollectMode
	return ms.Treasure.Mutate(b)
}

// helper function for MutableSlot.Check
func check(b []byte, addr Addr, value byte) error {
	en := isEn(b)

	if b[addr.FullOffset(en)] != value {
		return fmt.Errorf("expected %x at %x; found %x",
			value, addr.FullOffset(en), b[addr.FullOffset(en)])
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
	if ms.CollectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.CollectMode, ms.Treasure.mode)
	}

	return nil
}

// relativeAddrs constructs a slice of Addr from a relative english offset and
// a series of jp addresses.
func relativeAddrs(bank byte, enOffset int, jpAddrs ...uint16) []Addr {
	addrs := make([]Addr, len(jpAddrs))
	for i, jpAddr := range jpAddrs {
		addrs[i] = Addr{bank, jpAddr, uint16(int(jpAddr) + enOffset)}
	}
	return addrs
}

// sameAddr returns an address that's the same in both JP and US versions.
func sameAddr(bank byte, offset uint16) Addr {
	return Addr{bank, offset, offset}
}

// BasicSlot constucts a MutableSlot from a treasure name, bank number, and an
// address for each its ID and sub-ID. Most slots fit this pattern.
func BasicSlot(treasure string, bank byte,
	idOffset, subIDOffset uint16, enOffset int) *MutableSlot {
	return &MutableSlot{
		Treasure:   Treasures[treasure],
		IDAddrs:    relativeAddrs(bank, enOffset, idOffset),
		SubIDAddrs: relativeAddrs(bank, enOffset, subIDOffset),
	}
}

// MutableChest constructs a MutableSlot from a treasure name and an address in
// bank $15, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively (?) to chests.
func MutableChest(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x15, addr, addr+1, enChestOffset)
}

// MutableGift constructs a MutableSlot from a treasure name and an address in
// bank $0b, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to most items given by NPCs.
func MutableGift(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x0b, addr, addr+1, enGiftOffset)
}

// MutableFind constructs a MutableSlot from a treasure name and an address in
// bank $09, where the sub-ID and ID (in that order) are two consecutive bytes
// at that address. This applies to most items that are found lying around.
func MutableFind(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x09, addr+1, addr, enFindOffset)
}

func init() {
	// set item slot collect modes based on default treasures
	for _, slot := range ItemSlots {
		slot.CollectMode = slot.Treasure.mode
	}
}

var ItemSlots = map[string]*MutableSlot{
	// holodrum
	"lake chest": MutableChest("gasha seed", 0x53d5),
	"maku tree gift": &MutableSlot{
		Treasure: Treasures["gnarled key"],
		IDAddrs: []Addr{{0x15, 0x657d, 0x613a},
			{0x09, 0x7dff, 0x7e16}, {0x09, 0x7de6, 0x7dfd}},
		SubIDAddrs: []Addr{{0x15, 0x6580, 0x613d}, {0x09, 0x7e02, 0x7e19}},
	},
	"village SW chest":   MutableChest("rupees, 20", 0x53c1),
	"village SE chest":   MutableChest("rupees, 20", 0x53c5),
	"shovel gift":        MutableGift("shovel", 0x6a6e),
	"outdoor d2 chest":   MutableChest("gasha seed", 0x53c9),
	"blaino gift":        MutableGift("ricky's gloves", 0x64ce),
	"floodgate key spot": MutableFind("floodgate key", 0x626a),
	"square jewel chest": &MutableSlot{
		Treasure:   Treasures["square jewel"],
		IDAddrs:    relativeAddrs(0x0b, -2, 0x7397),
		SubIDAddrs: relativeAddrs(0x0b, -2, 0x739b),
	},
	"great moblin chest":    MutableChest("piece of heart", 0x53d1),
	"master's plaque chest": MutableChest("master's plaque", 0x554d),
	"diver gift": &MutableSlot{
		Treasure:   Treasures["flippers"],
		IDAddrs:    relativeAddrs(0x0b, -2, 0x7310, 0x72f3),
		SubIDAddrs: relativeAddrs(0x0b, -2, 0x7311),
	},
	"spring banana tree":   MutableFind("spring banana", 0x66af),
	"dragon key spot":      MutableFind("dragon key", 0x628c),
	"pyramid jewel spot":   MutableGift("pyramid jewel", 0x7350),
	"x-shaped jewel chest": MutableChest("x-shaped jewel", 0x53cd),
	"round jewel gift":     MutableGift("round jewel", 0x7334),
	"noble sword spot": &MutableSlot{
		// two cases depending on which sword you enter with
		Treasure:   Treasures["sword L-2"],
		IDAddrs:    relativeAddrs(0x0b, 1, 0x6417, 0x641e),
		SubIDAddrs: relativeAddrs(0x0b, 1, 0x6418, 0x641f),
	},
	"desert pit": &MutableSlot{
		Treasure:   Treasures["rusty bell"],
		IDAddrs:    []Addr{{0x09, 0x6476, 0x648d}, {0x0b, 0x60b0, 0x60b1}},
		SubIDAddrs: []Addr{{0x09, 0x6475, 0x648c}},
	},
	"desert chest":        MutableChest("blast ring", 0x53dd),
	"western coast chest": MutableChest("rang ring L-1", 0x53d9),
	"coast house chest":   MutableChest("bombs, 10", 0x53ef),
	"water cave chest":    MutableChest("octo ring", 0x54c4),
	"mushroom cave chest": MutableChest("quicksand ring", 0x54c8),
	"mystery cave chest":  MutableChest("moblin ring", 0x5541),
	"moblin road chest":   MutableChest("rupees, 30", 0x5545),
	"sunken cave chest":   MutableChest("gasha seed", 0x5549),
	"diver chest":         MutableChest("rupees, 50", 0x5551),
	"dry lake chest":      MutableChest("piece of heart", 0x5555),
	"goron chest":         MutableChest("armor ring L-2", 0x555d),
	"platform chest":      MutableChest("rupees, 50", 0x5565),
	"talon cave chest":    MutableChest("subrosian ring", 0x5561),
	"tarm gasha chest":    MutableChest("gasha seed", 0x53eb),
	"moblin cliff chest":  MutableChest("gasha seed", 0x54cc),

	// subrosia
	"winter tower":     BasicSlot("winter", 0x0b, 0x4fc5, 0x4fc6, 0),
	"summer tower":     BasicSlot("summer", 0x0b, 0x4fb9, 0x4fba, 0),
	"spring tower":     BasicSlot("spring", 0x0b, 0x4fb5, 0x4fb6, 0),
	"autumn tower":     BasicSlot("autumn", 0x0b, 0x4fc1, 0x4fc2, 0),
	"dance hall prize": MutableGift("boomerang L-1", 0x6648),
	"rod gift": &MutableSlot{
		Treasure:   Treasures["rod"],
		IDAddrs:    []Addr{{0x15, 0x7511, 0x70ce}},
		ParamAddrs: []Addr{{0x15, 0x750f, 0x70cc}},
	},
	"star ore spot": &MutableSlot{
		Treasure:   Treasures["star ore"],
		IDAddrs:    relativeAddrs(0x08, 0, 0x62f4, 0x62fe),
		SubIDAddrs: []Addr{}, // special case, not set at all
	},
	"blue ore chest":       MutableChest("blue ore", 0x53e2),
	"red ore chest":        MutableChest("red ore", 0x53e6),
	"non-rosa gasha chest": MutableChest("gasha seed", 0x54d8),
	"rosa gasha chest":     MutableChest("gasha seed", 0x5559),

	// hero's cave
	"d0 sword chest": &MutableSlot{
		Treasure:   Treasures["sword L-1"],
		IDAddrs:    relativeAddrs(0x0a, 0x0a, 0x7b86),
		ParamAddrs: relativeAddrs(0x0a, 0x0a, 0x7b88),
	},
	"d0 rupee chest": MutableChest("rupees, 30", 0x53f8),

	// d1
	"d1 satchel spot":  MutableFind("satchel", 0x669a),
	"d1 gasha chest":   MutableChest("gasha seed", 0x5400),
	"d1 bomb chest":    MutableChest("bombs, 10", 0x5408),
	"d1 ring chest":    MutableChest("discovery ring", 0x5414),
	"d1 compass chest": MutableChest("compass", 0x5404),
	"d1 map chest":     MutableChest("dungeon map", 0x5418),

	// d2
	"d2 bracelet chest": MutableChest("bracelet", 0x5424),
	"d2 10-rupee chest": MutableChest("rupees, 10", 0x541c),
	"d2 5-rupee chest":  MutableChest("rupees, 5", 0x5438),
	"d2 map chest":      MutableChest("dungeon map", 0x5428),
	"d2 compass chest":  MutableChest("compass", 0x5434),

	// d3
	"d3 feather chest": MutableChest("feather L-1", 0x5458),
	"d3 rupee chest":   MutableChest("rupees, 30", 0x543c),
	"d3 gasha chest":   MutableChest("gasha seed", 0x5444),
	"d3 bomb chest":    MutableChest("bombs, 10", 0x545c),
	"d3 compass chest": MutableChest("compass", 0x544c),
	"d3 map chest":     MutableChest("dungeon map", 0x5454),

	// d4
	"d4 slingshot chest": MutableChest("slingshot L-1", 0x5470),
	"d4 bomb chest":      MutableChest("bombs, 10", 0x5474),
	"d4 map chest":       MutableChest("dungeon map", 0x5468),
	"d4 compass chest":   MutableChest("compass", 0x5478),

	// d5
	"d5 magnet gloves chest": MutableChest("magnet gloves", 0x5480),
	"d5 large rupee chest":   MutableChest("rupees, 100", 0x5484),
	"d5 map chest":           MutableChest("dungeon map", 0x547c),
	"d5 compass chest":       MutableChest("compass", 0x548c),

	// d6
	"d6 boomerang chest": MutableChest("boomerang L-2", 0x54c0),
	"d6 rupee chest A":   MutableChest("rupees, 10", 0x54a0),
	"d6 rupee chest B":   MutableChest("rupees, 5", 0x54a8),
	"d6 bomb chest":      MutableChest("bombs, 10", 0x54ac),
	"d6 rupee chest C":   MutableChest("rupees, 5", 0x54b8),
	"d6 compass chest":   MutableChest("compass", 0x549c),
	"d6 map chest":       MutableChest("dungeon map", 0x54a4),

	// d7
	"d7 cape chest":    MutableChest("feather L-2", 0x54e1),
	"d7 rupee chest":   MutableChest("rupees, 1", 0x54dd),
	"d7 ring chest":    MutableChest("power ring L-1", 0x54f9),
	"d7 compass chest": MutableChest("compass", 0x54ed),
	"d7 map chest":     MutableChest("dungeon map", 0x54f5),

	//d8
	"d8 HSS chest":     MutableChest("slingshot L-2", 0x551d),
	"d8 bomb chest":    MutableChest("bombs, 10", 0x54fd),
	"d8 ring chest":    MutableChest("steadfast ring", 0x5509),
	"d8 compass chest": MutableChest("compass", 0x5515),
	"d8 map chest":     MutableChest("dungeon map", 0x5521),

	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	// "ring box L-2 gift": MutableGift("ring box L-2", 0x5c1a),

	// these are "fake" item slots in that they don't slot real treasures
	"ember tree": &MutableSlot{
		Treasure: Treasures["ember tree seeds"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x64cb),
	},
	"mystery tree": &MutableSlot{
		Treasure: Treasures["mystery tree seeds"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x67dd),
	},
	"scent tree": &MutableSlot{
		Treasure: Treasures["scent tree seeds"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x685c),
	},
	"pegasus tree": &MutableSlot{
		Treasure: Treasures["pegasus tree seeds"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x6870),
	},
	"sunken gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 1"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x69b0),
	},
	"tarm gale tree": &MutableSlot{
		Treasure: Treasures["gale tree seeds 2"],
		IDAddrs:  relativeAddrs(0x11, 3, 0x6a46),
	},
}
