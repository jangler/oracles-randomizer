package rom

import (
	"log"
)

// A Mutable is a byte string that can be changed by the randomizer.
type Mutable interface {
	RealAddr() int64 // return actual offset in ROM, accounting for bank
	Bytes() []byte   // "old" mutable data

	Mutate([]byte) error // change ROM bytes
}

// A MutableByte is a single mutable byte.
type MutableByte struct {
	Bank, Addr int // each bank is a 0x4000-byte offset, starting at 2
	Old, New   byte
}

func (mb MutableByte) RealAddr() int64 {
	return bankOffset(mb.Bank) + int64(mb.Addr)
}

func (mb MutableByte) Bytes() []byte {
	return []byte{mb.Old}
}

func (mb MutableByte) Mutate(b []byte) error {
	b[mb.RealAddr()] = mb.New
	return nil
}

// A MutableWord is two consecutive mutable bytes (not necessarily aligned).
type MutableWord struct {
	Bank, Addr int
	Old, New   uint16
}

func (mw MutableWord) RealAddr() int64 {
	return bankOffset(mw.Bank) + int64(mw.Addr)
}

func (mw MutableWord) Bytes() []byte {
	return []byte{byte(mw.Old >> 8), byte(mw.Old)}
}

func (mw MutableWord) Mutate(b []byte) error {
	addr, data := mw.RealAddr(), mw.Bytes()
	for i := 0; i < 2; i++ {
		b[addr+int64(i)] = data[i]
	}
	return nil
}

// XXX: so far, this file only handles items and obstacles enocuntered in
//      normal gameplay up through D2.

var holodrumMutables = map[string]Mutable{
	// want to have maku gate open from start
	"maku gate check": MutableByte{0x04, 0x61a3, 0x7e, 0x66},

	// want to have the horon village shop stock *and* sell items from the
	// start; replace each with $02
	"horon shop stock check": MutableByte{0x08, 0x4adb, 0x05, 0x02},
	"horon shop sell check":  MutableByte{0x08, 0x48d0, 0x05, 0x02},

	// also stock the strange flute without needing essences
	"horon shop flute check": MutableByte{0x08, 0x4b02, 0xcb, 0xf6},

	// can replace the gnarled key with a different item
	"maku drop ID":      MutableByte{0x15, 0x657d, 0x42, 0x42},
	"maku drop subID":   MutableByte{0x15, 0x6580, 0x00, 0x00},
	"maku redrop ID":    MutableByte{0x09, 0x7dff, 0x42, 0x42},
	"maku redrop subID": MutableByte{0x09, 0x7e02, 0x01, 0x01},
	"maku refall check": MutableByte{0x09, 0x7de6, 0x42, 0x42},

	// spawn rosa without having an essence
	"rosa spawn check": MutableByte{0x09, 0x678c, 0x40, 0x02},

	// swappable items
	"shovel gift": MutableWord{0x0b, 0x6a6e, 0x1500, 0x1500},

	// chests that could possibly matter in the overworld
	// TODO
}

// rod doesn't seem practical to swap, but maybe it could be placed somewhere
// in the overworld as a prerequisite to access subrosia.
var subrosiaMutables = map[string]Mutable{
	"boomerang L-1 gift": MutableWord{0x0b, 0x6648, 0x0600, 0x0600},
}

// hero's cave
var d0Mutables = map[string]Mutable{
	"d0 key chest":   MutableWord{0x15, 0x53f4, 0x3003, 0x3003},
	"d0 rupee chest": MutableWord{0x15, 0x53f8, 0x2804, 0x2804},
	"d0 sword chest": MutableWord{0x15, 0x53fc, 0x0500, 0x0500},

	// disable the "get sword" event that messes up the chest.
	// unfortunately this also disables the fade to white.
	"d0 sword event": MutableByte{0x11, 0x70ec, 0xf2, 0xff},
}

// dungeon 1
var d1Mutables = map[string]Mutable{
	"d1 key fall":      MutableWord{0x0b, 0x466f, 0x3001, 0x3001},
	"d1 map chest":     MutableWord{0x15, 0x5418, 0x3302, 0x3302},
	"d1 compass chest": MutableWord{0x15, 0x5404, 0x3202, 0x3202},
	"d1 gasha chest":   MutableWord{0x15, 0x5400, 0x3401, 0x3401},
	"d1 bomb chest":    MutableWord{0x15, 0x5408, 0x0300, 0x0300},
	"d1 key chest":     MutableWord{0x15, 0x540c, 0x3003, 0x3003},

	// this is backwards from a normal chest
	"d1 satchel": MutableWord{0x09, 0x669a, 0x0019, 0x0019},

	"d1 boss key chest": MutableWord{0x15, 0x5410, 0x3103, 0x3103},
	"d1 ring chest":     MutableWord{0x15, 0x5414, 0x2d04, 0x2d04},
}

// dungeon 2
var d2Mutables = map[string]Mutable{
	"d2 5-rupee chest":   MutableWord{0x15, 0x5438, 0x2801, 0x2801},
	"d2 key fall":        MutableWord{0x0b, 0x466f, 0x3001, 0x3001},
	"d2 compass chest":   MutableWord{0x15, 0x5434, 0x3202, 0x3202},
	"d2 map chest":       MutableWord{0x15, 0x5428, 0x3302, 0x3302},
	"d2 bracelet chest":  MutableWord{0x15, 0x5424, 0x1600, 0x1600},
	"d2 bomb key chest":  MutableWord{0x15, 0x542c, 0x3003, 0x3003},
	"d2 blade key chest": MutableWord{0x15, 0x5430, 0x3003, 0x3003},
	"d2 10-rupee chest":  MutableWord{0x15, 0x541c, 0x2802, 0x2802},
	"d2 boss key chest":  MutableWord{0x15, 0x5420, 0x3103, 0x3103},
}

var Mutables map[string]Mutable

func init() {
	mutableSets := []map[string]Mutable{
		holodrumMutables,
		subrosiaMutables,
		d0Mutables,
		d1Mutables,
		d2Mutables,
		Treasures,
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
