package rom

import (
	"fmt"
)

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure              *Treasure
	IDAddrs, SubIDAddrs   []Addr
	ParamAddrs, TextAddrs []Addr
	GfxAddrs              []Addr
	CollectMode           byte
}

// IsChest returns true iff the slot has a chest collection mode.
func IsChest(ms *MutableSlot) bool {
	return ms.CollectMode == CollectChest1 || ms.CollectMode == CollectChest2
}

// IsFound returns true iff the slot has a "normal" non-chest collection mode
// (they seem to be compatible).
func IsFound(ms *MutableSlot) bool {
	switch ms.CollectMode {
	case CollectRingBox, CollectUnderwater, CollectFind1, CollectFind2,
		CollectAppear:
		return true
	}
	return false
}

// Mutate replaces the given IDs and subIDs in the given ROM data, and changes
// the associated treasure's collection mode as appropriate.
func (ms *MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset()] = ms.Treasure.id
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

	// use a sub-ID based on slot (chest vs non-chest) for gasha seeds and
	// pieces of heart. for other treasures, use the set sub-ID and set the
	// treasure's collect mode accordingly.
	subID := ms.Treasure.subID
	switch ms.Treasure {
	case Treasures["gasha seed"], Treasures["piece of heart"]:
		if IsChest(ms) {
			subID = 1
		} else {
			subID = 0
		}
	default:
		ms.Treasure.mode = ms.CollectMode
		subID = ms.Treasure.subID
	}

	// for boss keys, override the map/compass chest collection mode with the
	// normal chest collection mode.
	if ms.Treasure.id == Treasures["boss key"].id {
		ms.Treasure.mode = CollectChest1
	}

	for _, addr := range ms.SubIDAddrs {
		b[addr.FullOffset()] = subID
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
	if ms.CollectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.CollectMode, ms.Treasure.mode)
	}

	return nil
}

// BasicSlot constucts a MutableSlot from a treasure name, bank number, and an
// address for each its ID and sub-ID. Most slots fit this pattern.
func BasicSlot(treasure string, bank byte,
	idOffset, subIDOffset uint16) *MutableSlot {
	return &MutableSlot{
		Treasure:   Treasures[treasure],
		IDAddrs:    []Addr{{bank, idOffset}},
		SubIDAddrs: []Addr{{bank, subIDOffset}},
	}
}

// MutableChest constructs a MutableSlot from a treasure name and an address in
// bank $15, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively (?) to chests.
func MutableChest(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x15, addr, addr+1)
}

// MutableGift constructs a MutableSlot from a treasure name and an address in
// bank $0b, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to most items given by NPCs.
func MutableGift(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x0b, addr, addr+1)
}

// MutableFind constructs a MutableSlot from a treasure name and an address in
// bank $09, where the sub-ID and ID (in that order) are two consecutive bytes
// at that address. This applies to most items that are found lying around.
func MutableFind(treasure string, addr uint16) *MutableSlot {
	return BasicSlot(treasure, 0x09, addr+1, addr)
}

func init() {
	// set item slot collect modes based on default treasures
	for _, slot := range ItemSlots {
		slot.CollectMode = slot.Treasure.mode
	}
}

var ItemSlots = map[string]*MutableSlot{
	// holodrum
	"lake chest": MutableChest("gasha seed", 0x4f92),
	"maku tree gift": &MutableSlot{
		Treasure:   Treasures["gnarled key"],
		IDAddrs:    []Addr{{0x15, 0x613a}, {0x09, 0x7e16}, {0x09, 0x7dfd}},
		SubIDAddrs: []Addr{{0x15, 0x613d}, {0x09, 0x7e19}},
	},
	"village SW chest":   MutableChest("rupees, 20", 0x4f7e),
	"village SE chest":   MutableChest("rupees, 20", 0x4f82),
	"shovel gift":        MutableGift("shovel", 0x6a6c),
	"outdoor d2 chest":   MutableChest("gasha seed", 0x4f86),
	"blaino gift":        MutableGift("gasha seed", 0x64cc),
	"floodgate key spot": MutableFind("floodgate key", 0x6281),
	"square jewel chest": &MutableSlot{
		Treasure:   Treasures["square jewel"],
		IDAddrs:    []Addr{{0x0b, 0x7395}},
		SubIDAddrs: []Addr{{0x0b, 0x7399}},
	},
	"great moblin chest":    MutableChest("piece of heart", 0x4f8e),
	"master's plaque chest": MutableChest("master's plaque", 0x510a),
	"diver gift": &MutableSlot{
		Treasure:   Treasures["flippers"],
		IDAddrs:    []Addr{{0x0b, 0x730e}, {0x0b, 0x72f1}},
		SubIDAddrs: []Addr{{0x0b, 0x730f}},
	},
	"spring banana tree":   MutableFind("spring banana", 0x66c6),
	"dragon key spot":      MutableFind("dragon key", 0x62a3),
	"pyramid jewel spot":   MutableGift("pyramid jewel", 0x734e),
	"x-shaped jewel chest": MutableChest("x-shaped jewel", 0x4f8a),
	"round jewel gift":     MutableGift("round jewel", 0x7332),
	"noble sword spot": &MutableSlot{
		Treasure:   Treasures["sword 2"],
		IDAddrs:    []Addr{{0x0b, 0x6418}, {0x0b, 0x641f}},
		SubIDAddrs: []Addr{{0x0b, 0x6419}, {0x0b, 0x6420}},
		GfxAddrs:   []Addr{{0x3f, 0x69f7}, {0x3f, 0x69fa}},
	},
	"desert pit": &MutableSlot{
		Treasure:   Treasures["rusty bell"],
		IDAddrs:    []Addr{{0x09, 0x648d}, {0x0b, 0x60b1}},
		SubIDAddrs: []Addr{{0x09, 0x648c}},
	},
	"desert chest":        MutableChest("blast ring", 0x4f9a),
	"western coast chest": MutableChest("rang ring L-1", 0x4f96),
	"coast house chest":   MutableChest("bombs, 10", 0x4fac),
	"water cave chest":    MutableChest("octo ring", 0x5081),
	"mushroom cave chest": MutableChest("quicksand ring", 0x5085),
	"mystery cave chest":  MutableChest("moblin ring", 0x50fe),
	"moblin road chest":   MutableChest("rupees, 30", 0x5102),
	"sunken cave chest":   MutableChest("gasha seed", 0x5106),
	"diver chest":         MutableChest("rupees, 50", 0x510e),
	"dry lake east chest": MutableChest("piece of heart", 0x5112),
	"goron chest":         MutableChest("armor ring L-2", 0x511a),
	"platform chest":      MutableChest("rupees, 50", 0x5122),
	"talon cave chest":    MutableChest("subrosian ring", 0x511e),
	"tarm gasha chest":    MutableChest("gasha seed", 0x4fa8),
	"moblin cliff chest":  MutableChest("gasha seed", 0x5089),
	"dry lake west chest": &MutableSlot{
		Treasure:   Treasures["rupees, 100"],
		IDAddrs:    []Addr{{0x0b, 0x73a1}},
		SubIDAddrs: []Addr{{0x0b, 0x73a5}},
	},
	"linked dive chest": &MutableSlot{
		Treasure:   Treasures["gasha seed"],
		IDAddrs:    []Addr{{0x0a, 0x5003}},
		SubIDAddrs: []Addr{{0x0a, 0x5008}},
	},

	// dummy slots for bombs and shield
	"village shop 1": &MutableSlot{
		Treasure: Treasures["bombs, 10"],
	},
	"village shop 2": &MutableSlot{
		Treasure: Treasures["shop shield L-1"],
	},

	"village shop 3": &MutableSlot{
		Treasure:   Treasures["strange flute"],
		IDAddrs:    []Addr{{0x08, 0x4ce8}, {0x08, 0x4af2}, {0x08, 0x4a8a}},
		SubIDAddrs: []Addr{{0x08, 0x4ce9}},
	},
	"member's shop 1": &MutableSlot{
		Treasure:   Treasures["satchel 2"],
		IDAddrs:    []Addr{{0x08, 0x4cce}},
		SubIDAddrs: []Addr{{0x08, 0x4ccf}},
	},
	"member's shop 2": &MutableSlot{
		Treasure:   Treasures["gasha seed"],
		IDAddrs:    []Addr{{0x08, 0x4cd2}},
		SubIDAddrs: []Addr{{0x08, 0x4cd3}},
	},
	"member's shop 3": &MutableSlot{
		Treasure:   Treasures["treasure map"],
		IDAddrs:    []Addr{{0x08, 0x4cd8}},
		SubIDAddrs: []Addr{{0x08, 0x4cd9}},
	},

	// subrosia
	"winter tower":     BasicSlot("winter", 0x0b, 0x4fc5, 0x4fc6),
	"summer tower":     BasicSlot("summer", 0x0b, 0x4fb9, 0x4fba),
	"spring tower":     BasicSlot("spring", 0x0b, 0x4fb5, 0x4fb6),
	"autumn tower":     BasicSlot("autumn", 0x0b, 0x4fc1, 0x4fc2),
	"dance hall prize": MutableGift("boomerang L-1", 0x6646),
	"rod gift": &MutableSlot{
		Treasure:   Treasures["rod"],
		IDAddrs:    []Addr{{0x15, 0x70ce}},
		SubIDAddrs: []Addr{{0x15, 0x70cc}},
		GfxAddrs:   []Addr{{0x3f, 0x6c25}},
	},
	"star ore spot": &MutableSlot{
		Treasure:   Treasures["star ore"],
		IDAddrs:    []Addr{{0x08, 0x62f4}, {0x08, 0x62fe}},
		SubIDAddrs: []Addr{}, // special case, not set at all
	},
	"blue ore chest":       MutableChest("blue ore", 0x4f9f),
	"red ore chest":        MutableChest("red ore", 0x4fa3),
	"non-rosa gasha chest": MutableChest("gasha seed", 0x5095),
	"rosa gasha chest":     MutableChest("gasha seed", 0x5116),
	"subrosian market 1": &MutableSlot{
		Treasure:   Treasures["ribbon"],
		IDAddrs:    []Addr{{0x09, 0x77da}},
		SubIDAddrs: []Addr{{0x09, 0x77db}},
	},
	"subrosian market 2": &MutableSlot{
		Treasure:   Treasures["rare peach stone"],
		IDAddrs:    []Addr{{0x09, 0x77e2}},
		SubIDAddrs: []Addr{{0x09, 0x77e3}},
	},
	"subrosian market 5": &MutableSlot{
		Treasure:   Treasures["member's card"],
		IDAddrs:    []Addr{{0x09, 0x77f4}, {0x09, 0x7755}},
		SubIDAddrs: []Addr{{0x09, 0x77f5}},
	},
	"hard ore slot": &MutableSlot{
		Treasure:   Treasures["hard ore"],
		IDAddrs:    []Addr{{0x15, 0x5b85}, {0x09, 0x66eb}},
		SubIDAddrs: []Addr{},
	},
	"iron shield gift": &MutableSlot{
		Treasure:   Treasures["shield L-2"],
		IDAddrs:    []Addr{{0x15, 0x62be}},
		ParamAddrs: []Addr{{0x15, 0x62b4}},
	},

	// hero's cave
	"d0 sword chest": &MutableSlot{
		Treasure:   Treasures["sword 1"],
		IDAddrs:    []Addr{{0x0a, 0x7b90}},
		ParamAddrs: []Addr{{0x0a, 0x7b92}},
		TextAddrs:  []Addr{{0x0a, 0x7b9c}},
		GfxAddrs:   []Addr{{0x3f, 0x6676}},
	},
	"d0 rupee chest": MutableChest("rupees, 30", 0x4fb5),

	// d1
	"d1 satchel spot":   MutableFind("satchel 1", 0x66b1),
	"d1 gasha chest":    MutableChest("gasha seed", 0x4fbd),
	"d1 bomb chest":     MutableChest("bombs, 10", 0x4fc5),
	"d1 ring chest":     MutableChest("discovery ring", 0x4fd1),
	"d1 compass chest":  MutableChest("compass", 0x4fc1),
	"d1 map chest":      MutableChest("dungeon map", 0x4fd5),
	"d1 boss key chest": MutableChest("d1 boss key", 0x4fcd),

	// d2
	"d2 bracelet chest": MutableChest("bracelet", 0x4fe1),
	"d2 10-rupee chest": MutableChest("rupees, 10", 0x4fd9),
	"d2 5-rupee chest":  MutableChest("rupees, 5", 0x4ff5),
	"d2 map chest":      MutableChest("dungeon map", 0x4fe5),
	"d2 compass chest":  MutableChest("compass", 0x4ff1),
	"d2 boss key chest": MutableChest("d2 boss key", 0x4fdd),

	// d3
	"d3 feather chest":  MutableChest("feather 1", 0x5015),
	"d3 rupee chest":    MutableChest("rupees, 30", 0x4ff9),
	"d3 gasha chest":    MutableChest("gasha seed", 0x5001),
	"d3 bomb chest":     MutableChest("bombs, 10", 0x5019),
	"d3 compass chest":  MutableChest("compass", 0x5009),
	"d3 map chest":      MutableChest("dungeon map", 0x5011),
	"d3 boss key chest": MutableChest("d3 boss key", 0x4ffd),

	// d4
	"d4 slingshot chest": MutableChest("slingshot 1", 0x502d),
	"d4 bomb chest":      MutableChest("bombs, 10", 0x5031),
	"d4 map chest":       MutableChest("dungeon map", 0x5025),
	"d4 compass chest":   MutableChest("compass", 0x5035),

	// d5
	"d5 magnet gloves chest": MutableChest("magnet gloves", 0x503d),
	"d5 large rupee chest":   MutableChest("rupees, 100", 0x5041),
	"d5 map chest":           MutableChest("dungeon map", 0x5039),
	"d5 compass chest":       MutableChest("compass", 0x5049),

	// d6
	"d6 boomerang chest": MutableChest("boomerang L-2", 0x507d),
	"d6 rupee chest A":   MutableChest("rupees, 10", 0x505d),
	"d6 rupee chest B":   MutableChest("rupees, 5", 0x5065),
	"d6 bomb chest":      MutableChest("bombs, 10", 0x5069),
	"d6 rupee chest C":   MutableChest("rupees, 5", 0x5075),
	"d6 compass chest":   MutableChest("compass", 0x5059),
	"d6 map chest":       MutableChest("dungeon map", 0x5061),
	"d6 boss key chest":  MutableChest("d6 boss key", 0x5079),

	// d7
	"d7 cape chest":     MutableChest("feather 2", 0x509e),
	"d7 rupee chest":    MutableChest("rupees, 1", 0x509a),
	"d7 ring chest":     MutableChest("power ring L-1", 0x50b6),
	"d7 compass chest":  MutableChest("compass", 0x50aa),
	"d7 map chest":      MutableChest("dungeon map", 0x50b2),
	"d7 boss key chest": MutableChest("d7 boss key", 0x50a6),

	// d8
	"d8 HSS chest":      MutableChest("slingshot 2", 0x50da),
	"d8 bomb chest":     MutableChest("bombs, 10", 0x50ba),
	"d8 ring chest":     MutableChest("steadfast ring", 0x50c6),
	"d8 compass chest":  MutableChest("compass", 0x50d2),
	"d8 map chest":      MutableChest("dungeon map", 0x50de),
	"d8 boss key chest": MutableChest("d8 boss key", 0x50ca),

	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	// "ring box L-2 gift": MutableGift("ring box L-2", 0x5c18),

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
