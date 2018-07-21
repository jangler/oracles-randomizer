// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"bytes"
	"io"
)

// A Mutable is a byte that can be changed by the randomizer. Addr is the
// offset within the bank, Old is the original value (for validation purposes),
// and New is the replacement value.
type Mutable struct {
	Bank, Addr int // each bank is a 0x4000-byte offset, starting at 2
	Old, New   byte
}

// A MutableWord is two consecutive mutable bytes (not necessarily aligned).
type MutableWord struct {
	Bank, Addr int
	Old, New   uint16
}

// XXX: so far, this file only handles items and obstacles enocuntered in
//      normal gameplay up through D2.

// holodrum
var (
	// want to have maku gate open from start
	MakuGateCheck = Mutable{0x04, 0x6a13, 0x7e, 0x66}

	// want to have the horon village shop stock *and* sell items from the
	// start; replace each with $02
	HoronShopStockCheck = Mutable{0x08, 0x4adb, 0x05, 0x02}
	HoronShopSellCheck  = Mutable{0x08, 0x48d0, 0x05, 0x02}

	// also stock the strange flute without needing essences
	HoronShopFluteCheck = Mutable{0x08, 0x4b02, 0xcb, 0xf6}

	// can replace the gnarled key with a different item
	MakuDropID      = Mutable{0x15, 0x657d, 0x42, 0x42}
	MakuDropSubID   = Mutable{0x15, 0x6580, 0x00, 0x00}
	MakuRedropID    = Mutable{0x09, 0x7dff, 0x42, 0x42}
	MakuRedropSubID = Mutable{0x09, 0x7e02, 0x01, 0x01}
	MakuRedropCheck = Mutable{0x09, 0x7de6, 0x42, 0x42}

	// spawn rosa without having an essence
	RosaEssenceCheck = Mutable{0x09, 0x678c, 0x40, 0x02}

	// swappable items
	Shovel = MutableWord{0x0b, 0x6a6f, 0x1500, 0x1500}

	// chests that could possibly matter in the overworld
	// TODO
)

// subrosia
// rod doesn't seem practical to swap, but maybe it could be placed somewhere
// in the overworld as a prerequisite to access subrosia.
var (
	BoomerangL1 = MutableWord{0x0b, 0x6648, 0x0600, 0x0600}
)

// hero's cave
var (
	D0KeyChest   = MutableWord{0x15, 0x53f4, 0x3303, 0x3303}
	D0RupeeChest = MutableWord{0x15, 0x53f8, 0x2804, 0x2804}
	D0SwordChest = MutableWord{0x15, 0x53fc, 0x5000, 0x5000}

	// disable the "get sword" event that messes up the chest.
	// unfortunately this also disables the fade to white.
	D0SwordEvent = Mutable{0x11, 0x70ec, 0xf2, 0xff}
)

// dungeon 1
var (
	D1KeyFall      = MutableWord{0x0b, 0x466f, 0x3301, 0x3301}
	D1MapChest     = MutableWord{0x15, 0x5418, 0x3302, 0x3302}
	D1CompassChest = MutableWord{0x15, 0x5404, 0x3202, 0x3202}
	D1GashaChest   = MutableWord{0x15, 0x5400, 0x3101, 0x3401}
	D1BombsChest   = MutableWord{0x15, 0x5408, 0x0300, 0x0300}
	D1KeyChest     = MutableWord{0x15, 0x540c, 0x3003, 0x3003}
	D1Satchel      = MutableWord{0x09, 0x669a, 0x1900, 0x1900}
	D1BossKeyChest = MutableWord{0x15, 0x5410, 0x3103, 0x3103}
	D1RingChest    = MutableWord{0x15, 0x5414, 0x2d04, 0x2d04}
)

// dungeon 2
var (
	D2Rupee5Chest   = MutableWord{0x15, 0x5438, 0x2801, 0x2801}
	D2KeyFall       = MutableWord{0x0b, 0x466f, 0x3001, 0x3001}
	D2CompassChest  = MutableWord{0x15, 0x5434, 0x3202, 0x3202}
	D2MapChest      = MutableWord{0x15, 0x5428, 0x3302, 0x3302}
	D2BraceletChest = MutableWord{0x15, 0x5424, 0x1600, 0x1600}
	D2BombKeyChest  = MutableWord{0x15, 0x542c, 0x3003, 0x3003}
	D2BladeKeyChest = MutableWord{0x15, 0x5430, 0x3003, 0x3003}
	D2Rupee10Chest  = MutableWord{0x15, 0x541c, 0x2802, 0x2802}
	D2BossKeyChest  = MutableWord{0x15, 0x5420, 0x3103, 0x3103}
)

// LoadROM reads ROM data from a reader into memory.
func LoadROM(f io.Reader) error {
	// read file into buffer
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, f)
	return err
}
