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

	// also stock the strange flute from the start
	HoronShopFluteCheck = Mutable{} // TODO

	// can replace the gnarled key with a different item
	MakuDropID      = Mutable{0x15, 0x657d, 0x42, 0x42}
	MakuDropSubID   = Mutable{0x15, 0x6580, 0x00, 0x00}
	MakuRedropID    = Mutable{0x09, 0x7dff, 0x42, 0x42}
	MakuRedropSubID = Mutable{0x09, 0x7e02, 0x01, 0x01}
	MakuRedropCheck = Mutable{0x09, 0x7de6, 0x42, 0x42}

	// spawn rosa without having an essence
	RosaEssenceCheck = Mutable{0x09, 0x678c, 0x40, 0x02}

	// swappable items
	Shovel = MutableWord{} // TODO
)

// subrosia
// rod doesn't seem practical to swap
var (
	BoomerangL1 = MutableWord{0x0b, 0x6648, 0x0600, 0x0600}
)

// hero's cave
var (
	D0KeyChest   = MutableWord{} // TODO
	D0SwordChest = MutableWord{} // TODO
	D0RupeeChest = MutableWord{} // TODO
)

// dungeon 1
var (
	D1KeyFall      = MutableWord{} // TODO
	D1MapChest     = MutableWord{0x15, 0x5418, 0x3302, 0x3302}
	D1CompassChest = MutableWord{} // TODO
	D1GashaChest   = MutableWord{} // TODO
	D1BombsChest   = MutableWord{} // TODO
	D1KeyChest     = MutableWord{} // TODO
	D1Item         = MutableWord{0x09, 0x669a, 0x1900, 0x1900}
	D1BossKeyChest = MutableWord{0x15, 0x5410, 0x3103, 0x3103}
	D1RingChest    = MutableWord{0x15, 0x5414, 0x2d04, 0x2d04}
)

type treasure struct {
	id, param byte
	addr      uint16 // bank 15, value of hl at $15:466b

	// in order, starting at addr - 1
	mode, value, text, sprite byte
}

var (
	// treasure item info
	shieldL1     = treasure{0x01, 0x00, 0x5701, 0x0a, 0x01, 0x1f, 0x13}
	bombs        = treasure{0x03, 0x00, 0x570d, 0x38, 0x10, 0x4d, 0x05}
	swordL1      = treasure{0x05, 0x00, 0x577d, 0x38, 0x01, 0x1c, 0x10}
	swordL2      = treasure{0x05, 0x01, 0x5721, 0x09, 0x02, 0x1d, 0x11}
	boomerangL1  = treasure{0x06, 0x00, 0x5735, 0x0a, 0x01, 0x22, 0x1c}
	boomerangL2  = treasure{0x06, 0x01, 0x5739, 0x38, 0x02, 0x23, 0x1d}
	rod          = treasure{0x07, 0x00, 0x573d, 0x38, 0x07, 0x0a, 0x1e}
	magnetGloves = treasure{0x08, 0x00, 0x558d, 0x38, 0x00, 0x30, 0x18}
	strangeFlute = treasure{0x0e, 0x00, 0x55a5, 0x0a, 0x0c, 0x3b, 0x23}
	slingshotL1  = treasure{0x13, 0x00, 0x5769, 0x38, 0x01, 0x2e, 0x21}
	slingshotL2  = treasure{0x13, 0x01, 0x576d, 0x38, 0x02, 0x2f, 0x22}
	shovel       = treasure{0x15, 0x00, 0x55c1, 0x0a, 0x00, 0x25, 0x1b}
	bracelet     = treasure{0x16, 0x00, 0x55c5, 0x38, 0x00, 0x26, 0x19}
	featherL1    = treasure{0x17, 0x00, 0x5771, 0x38, 0x01, 0x27, 0x16}
	featherL2    = treasure{0x16, 0x01, 0x5775, 0x38, 0x02, 0x28, 0x17}
	satchel      = treasure{0x19, 0x00, 0x56f9, 0x0a, 0x01, 0x2d, 0x20}
	foolsOre     = treasure{0x1e, 0x00, 0x55e5, 0x00, 0x00, 0xff, 0x00}
	flippers     = treasure{0x2e, 0x00, 0x5625, 0x0a, 0x00, 0x31, 0x31}

	// rings actually have various entries based on param
	ring = treasure{0x2d, 0x00, 0x57fd, 0x09, 0xff, 0x54, 0x0e}

	// TODO
	smallKey   = treasure{}
	bossKey    = treasure{}
	compass    = treasure{}
	dungeonMap = treasure{}
	gnarledKey = treasure{}

	// not until after D2
	floodgateKey  = treasure{}
	dragonKey     = treasure{}
	starOre       = treasure{}
	ribbon        = treasure{}
	springBanana  = treasure{}
	rickysGloves  = treasure{}
	bombFlower    = treasure{}
	piratesBell   = treasure{}
	roundJewel    = treasure{}
	pyramidJewel  = treasure{}
	squareJewel   = treasure{}
	xShapedJewel  = treasure{}
	membersCard   = treasure{}
	mastersPlaque = treasure{}
)
