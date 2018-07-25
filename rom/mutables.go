package rom

import (
	"fmt"
	"log"
)

// A Mutable is a memory data that can be changed by the randomizer.
type Mutable interface {
	Mutate([]byte) error // change ROM bytes
	Check([]byte) error  // verify that the mutable matches the ROM
}

// A MutableByte is a single mutable byte.
type MutableByte struct {
	Addr     Addr
	Old, New byte
}

func (mb MutableByte) Mutate(b []byte) error {
	b[mb.Addr.FullOffset()] = mb.New
	return nil
}

func (mb MutableByte) Check(b []byte) error {
	addr := mb.Addr.FullOffset()
	if b[addr] == mb.Old {
		return nil
	}
	return fmt.Errorf("expected %x at %x; found %x", mb.Old, addr, b[addr])
}

// A MutableWord is two consecutive mutable bytes (not necessarily aligned).
type MutableWord struct {
	Addr     Addr
	Old, New uint16
}

func (mw MutableWord) Mutate(b []byte) error {
	addr := mw.Addr.FullOffset()
	b[addr] = byte(mw.New >> 8)
	b[addr+1] = byte(mw.New)
	return nil
}

func (mw MutableWord) Check(b []byte) error {
	addr := mw.Addr.FullOffset()
	if b[addr] == byte(mw.Old>>8) && b[addr+1] == byte(mw.Old) {
		return nil
	}
	return fmt.Errorf("expected %x at %x; found %x",
		mw.Old, addr, b[addr:addr+2])
}

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure            *Treasure
	IDAddrs, SubIDAddrs []Addr
	CollectMode         byte
}

func (ms MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.IDAddrs {
		b[addr.FullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.SubIDAddrs {
		b[addr.FullOffset()] = ms.Treasure.subID
	}

	// use normal chest mode instead of the x-shaped jewel one
	// since the x-shaped jewel one can lock the game
	if ms.CollectMode == CollectXChest {
		ms.Treasure.mode = CollectChest
	} else {
		ms.Treasure.mode = ms.CollectMode
	}

	return ms.Treasure.Mutate(b)
}

func (ms MutableSlot) Check(b []byte) error {
	for _, addr := range ms.IDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.id {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.id, addr.FullOffset(), b[addr.FullOffset()])
		}
	}
	for _, addr := range ms.SubIDAddrs {
		if b[addr.FullOffset()] != ms.Treasure.subID {
			return fmt.Errorf("expected %x at %x; found %x",
				ms.Treasure.subID, addr.FullOffset(), b[addr.FullOffset()])
		}
	}
	if ms.CollectMode != ms.Treasure.mode {
		return fmt.Errorf("slot/treasure collect mode mismatch: %x/%x",
			ms.CollectMode, ms.Treasure.mode)
	}

	return nil
}

var ItemSlots = map[string]*MutableSlot{
	"d0 sword chest": &MutableSlot{
		Treasure:    Treasures["sword L-1"],
		IDAddrs:     []Addr{{0x15, 0x53fc}},
		SubIDAddrs:  []Addr{{0x15, 0x53fd}},
		CollectMode: CollectChest,
	},
	"maku key fall": &MutableSlot{
		Treasure:    Treasures["gnarled key"],
		IDAddrs:     []Addr{{0x15, 0x657d}, {0x09, 0x7dff}, {0x09, 0x7de6}},
		SubIDAddrs:  []Addr{{0x15, 0x6580}, {0x09, 0x7e02}},
		CollectMode: CollectFall,
	},
	"boomerang gift": &MutableSlot{
		Treasure:    Treasures["boomerang L-1"],
		IDAddrs:     []Addr{{0x0b, 0x6648}},
		SubIDAddrs:  []Addr{{0x0b, 0x6649}},
		CollectMode: CollectFind2,
	},
	"shovel gift": &MutableSlot{
		Treasure:    Treasures["shovel"],
		IDAddrs:     []Addr{{0x0b, 0x6a6e}},
		SubIDAddrs:  []Addr{{0x0b, 0x6a6f}},
		CollectMode: CollectFind2,
	},
	"d1 satchel": &MutableSlot{
		// addresses are backwards from a normal slot
		Treasure:    Treasures["satchel"],
		IDAddrs:     []Addr{{0x09, 0x669b}},
		SubIDAddrs:  []Addr{{0x09, 0x669a}},
		CollectMode: CollectFind2,
	},
	"d2 bracelet chest": &MutableSlot{
		Treasure:    Treasures["bracelet"],
		IDAddrs:     []Addr{{0x15, 0x5424}},
		SubIDAddrs:  []Addr{{0x15, 0x5425}},
		CollectMode: CollectChest,
	},
	"blaino gift": &MutableSlot{
		Treasure:    Treasures["ricky's gloves"],
		IDAddrs:     []Addr{{0x0b, 0x64ce}},
		SubIDAddrs:  []Addr{{0x0b, 0x64cf}},
		CollectMode: CollectFind1,
	},
	"floodgate key gift": &MutableSlot{
		Treasure:    Treasures["floodgate key"],
		IDAddrs:     []Addr{{0x09, 0x626b}},
		SubIDAddrs:  []Addr{{0x09, 0x626a}},
		CollectMode: CollectFind1,
	},
	"square jewel chest": &MutableSlot{
		Treasure:    Treasures["square jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7397}},
		SubIDAddrs:  []Addr{{0x0b, 0x739b}},
		CollectMode: CollectChest,
	},
	"x-shaped jewel chest": &MutableSlot{
		Treasure:    Treasures["x-shaped jewel"],
		IDAddrs:     []Addr{{0x15, 0x53cd}},
		SubIDAddrs:  []Addr{{0x15, 0x53ce}},
		CollectMode: CollectXChest,
	},
	"star ore spot": &MutableSlot{
		Treasure:    Treasures["star ore"],
		IDAddrs:     []Addr{{0x08, 0x62f4}, {0x08, 0x62fe}},
		SubIDAddrs:  []Addr{}, // special case, not set at all
		CollectMode: CollectDig,
	},
	"d3 feather chest": &MutableSlot{
		Treasure:    Treasures["feather L-1"],
		IDAddrs:     []Addr{{0x15, 0x5458}},
		SubIDAddrs:  []Addr{{0x15, 0x5459}},
		CollectMode: CollectChest,
	},
	"master's plaque chest": &MutableSlot{
		Treasure:    Treasures["master's plaque"],
		IDAddrs:     []Addr{{0x15, 0x554d}},
		SubIDAddrs:  []Addr{{0x15, 0x554e}},
		CollectMode: CollectChest,
	},
	"flippers gift": &MutableSlot{
		Treasure:    Treasures["flippers"],
		IDAddrs:     []Addr{{0x0b, 0x7310}},
		SubIDAddrs:  []Addr{{0x0b, 0x7311}},
		CollectMode: CollectFind2,
	},
	"spring banana tree": &MutableSlot{
		Treasure:    Treasures["spring banana"],
		IDAddrs:     []Addr{{0x09, 0x66b0}},
		SubIDAddrs:  []Addr{{0x09, 0x66af}},
		CollectMode: CollectFind2,
	},
	"dragon key spot": &MutableSlot{
		Treasure:    Treasures["dragon key"],
		IDAddrs:     []Addr{{0x09, 0x628d}}, // i *think* this is the only one
		SubIDAddrs:  []Addr{{0x09, 0x628c}},
		CollectMode: CollectFind1,
	},
	"pyramid jewel spot": &MutableSlot{
		Treasure:    Treasures["pyramid jewel"],
		IDAddrs:     []Addr{{0x0b, 0x7350}},
		SubIDAddrs:  []Addr{{0x0b, 0x7351}},
		CollectMode: CollectUnderwater,
	},
	"ring box L-2 gift": &MutableSlot{
		Treasure:    Treasures["ring box L-2"],
		IDAddrs:     []Addr{{0x0b, 0x5c1a}},
		SubIDAddrs:  []Addr{{0x0b, 0x5c1b}},
		CollectMode: CollectGoronGift,
	},
	"d4 slingshot chest": &MutableSlot{
		Treasure:    Treasures["slingshot L-1"],
		IDAddrs:     []Addr{{0x15, 0x5470}},
		SubIDAddrs:  []Addr{{0x15, 0x5471}},
		CollectMode: CollectChest,
	},
}

var holodrumMutables = map[string]Mutable{
	// want to have maku gate open from start
	"maku gate check": MutableByte{Addr{0x04, 0x61a3}, 0x7e, 0x66},

	// want to have the horon village shop stock *and* sell items from the
	// start; replace each with $02
	"horon shop stock check": MutableByte{Addr{0x08, 0x4adb}, 0x05, 0x02},
	"horon shop sell check":  MutableByte{Addr{0x08, 0x48d0}, 0x05, 0x02},

	// remove checks against specific essence flags. changing 0xcb to 0xf6
	// turns the check from a bit check to an or against next byte in the code,
	// which is always nonzero.
	"horon shop flute check": MutableByte{Addr{0x08, 0x4b02}, 0xcb, 0xf6},
	"ricky spawn check":      MutableByte{Addr{0x09, 0x4e68}, 0xcb, 0xf6},

	// spawn rosa without having an essence (she checks for the treasure ID)
	"rosa spawn check": MutableByte{Addr{0x09, 0x678c}, 0x40, 0x02},

	// spawn dimitri without having essences, and even with flippers
	"dimitri essence check": MutableByte{Addr{0x09, 0x4e36}, 0xcb, 0xf6},
	"dimitri flipper check": MutableByte{Addr{0x09, 0x4e4c}, 0x2e, 0x00},

	// master who gives flippers checks essences two ways
	"master essence check 2": MutableByte{Addr{0x0a, 0x4bea}, 0x40, 0x02},
	"master essence check 1": MutableByte{Addr{0x0a, 0x4bf5}, 0x02, 0x00},
}

// hero's cave
var d0Mutables = map[string]Mutable{
	// unused
	"d0 key chest":   MutableWord{Addr{0x15, 0x53f4}, 0x3003, 0x3003},
	"d0 rupee chest": MutableWord{Addr{0x15, 0x53f8}, 0x2804, 0x2804},

	// disable the "get sword" event that messes up the chest.
	// unfortunately this also disables the fade to white.
	"d0 sword event": MutableByte{Addr{0x11, 0x70ec}, 0xf2, 0xff},
}

// dungeon 1 (unused)
var d1Mutables = map[string]Mutable{
	"d1 key fall":       MutableWord{Addr{0x0b, 0x466f}, 0x3001, 0x3001},
	"d1 map chest":      MutableWord{Addr{0x15, 0x5418}, 0x3302, 0x3302},
	"d1 compass chest":  MutableWord{Addr{0x15, 0x5404}, 0x3202, 0x3202},
	"d1 gasha chest":    MutableWord{Addr{0x15, 0x5400}, 0x3401, 0x3401},
	"d1 bomb chest":     MutableWord{Addr{0x15, 0x5408}, 0x0300, 0x0300},
	"d1 key chest":      MutableWord{Addr{0x15, 0x540c}, 0x3003, 0x3003},
	"d1 boss key chest": MutableWord{Addr{0x15, 0x5410}, 0x3103, 0x3103},
	"d1 ring chest":     MutableWord{Addr{0x15, 0x5414}, 0x2d04, 0x2d04},
}

// dungeon 2 (unused)
var d2Mutables = map[string]Mutable{
	"d2 5-rupee chest":   MutableWord{Addr{0x15, 0x5438}, 0x2801, 0x2801},
	"d2 key fall":        MutableWord{Addr{0x0b, 0x466f}, 0x3001, 0x3001},
	"d2 compass chest":   MutableWord{Addr{0x15, 0x5434}, 0x3202, 0x3202},
	"d2 map chest":       MutableWord{Addr{0x15, 0x5428}, 0x3302, 0x3302},
	"d2 bomb key chest":  MutableWord{Addr{0x15, 0x542c}, 0x3003, 0x3003},
	"d2 blade key chest": MutableWord{Addr{0x15, 0x5430}, 0x3003, 0x3003},
	"d2 10-rupee chest":  MutableWord{Addr{0x15, 0x541c}, 0x2802, 0x2802},
	"d2 boss key chest":  MutableWord{Addr{0x15, 0x5420}, 0x3103, 0x3103},
}

var Mutables map[string]Mutable

func init() {
	slotMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		slotMutables[k] = v
	}
	treasureMutables := make(map[string]Mutable)
	for k, v := range Treasures {
		treasureMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		holodrumMutables,
		d0Mutables,
		d1Mutables,
		d2Mutables,
		treasureMutables,
		slotMutables,
	}

	// initialize master map w/ adequate capacity
	count := 0
	for _, set := range mutableSets {
		count += len(set)
	}
	Mutables = make(map[string]Mutable, count)

	// add mutables to master map
	for _, set := range mutableSets {
		for k, v := range set {
			if _, ok := Mutables[k]; ok {
				log.Fatalf("duplicate mutable key: %s", k)
			}
			Mutables[k] = v
		}
	}
}
