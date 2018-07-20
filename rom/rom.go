// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

// A Mutable is a byte that can be changed by the randomizer. Addr is the
// offset within the bank, Old is the original value (for validation purposes),
// and New is the replacement value.
type Mutable struct {
	Bank, Addr int // each bank is a 0x4000-byte offset, starting at 2
	Old, New   byte
}

// code/data positions
var (
	// want to have maku gate open from start
	MakuGateCheck = Mutable{0x04, 0x6a13, 0x7e, 0x66}

	// want to have the horon village shop stock *and* sell items from the
	// start; replace each with $02
	HoronShopStockCheck = Mutable{0x08, 0x4adb, 0x05, 0x02}
	HoronShopSellCheck  = Mutable{0x08, 0x48d0, 0x05, 0x02}

	// can replace the gnarled key with a different item
	MakuDropID      = Mutable{0x15, 0x657d, 0x42, 0x42}
	MakuDropSubID   = Mutable{0x15, 0x6580, 0x00, 0x00}
	MakuRedropID    = Mutable{0x09, 0x7dff, 0x42, 0x42}
	MakuRedropSubID = Mutable{0x09, 0x7e02, 0x01, 0x01}
	MakuRedropCheck = Mutable{0x09, 0x7de6, 0x42, 0x42}
)
