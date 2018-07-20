// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

// A Bank corresponds to an MBC5 memory bank offset, labeled by the BGB
// reckoning.
type Bank int

// only banks which are used in the randomizer are defined
const (
	ROM4 Bank = 0x0c000
	ROM8 Bank = 0x1c000
	RO15 Bank = 0x50000
)

// A Mutable is a byte that can be changed by the randomizer. Addr is the
// offset within the bank, Old is the original value (for validation purposes),
// and New is the replacement value.
type Mutable struct {
	Bank     Bank
	Addr     int
	Old, New byte
}

// code/data positions
var (
	// want to have maku gate open from start
	MakuGateCheck = Mutable{ROM4, 0x6a13, 0x7e, 0x66}

	// want to have the horon village shop stock *and* sell items from the
	// start; replace each with $02
	HoronShopStockCheck = Mutable{ROM8, 0x4adb, 0x05, 0x02}
	HoronShopSellCheck  = Mutable{ROM8, 0x48d0, 0x05, 0x02}

	// can replace the gnarled key with a different item
	MakuDropID    = Mutable{RO15, 0x657d, 0x42, 0x42}
	MakuDropSubID = Mutable{RO15, 0x6580, 0x00, 0x00}
)
