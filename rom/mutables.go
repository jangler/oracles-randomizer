package rom

import (
	"bytes"
	"fmt"
	"log"
)

// A Mutable is a memory data that can be changed by the randomizer.
type Mutable interface {
	Mutate([]byte) error // change ROM bytes
	Check([]byte) error  // verify that the mutable matches the ROM
}

// A MutableRange is a length of mutable bytes starting at a given address.
type MutableRange struct {
	Addr     Addr
	Old, New []byte
}

// MutableByte returns a special case of MutableRange with a range of a single
// byte.
func MutableByte(addr Addr, old, new byte) *MutableRange {
	return &MutableRange{Addr: addr, Old: []byte{old}, New: []byte{new}}
}

// MutableWord returns a special case of MutableRange with a range of a two
// bytes.
func MutableWord(addr Addr, old, new uint16) *MutableRange {
	return &MutableRange{
		Addr: addr,
		Old:  []byte{byte(old >> 8), byte(old)},
		New:  []byte{byte(new >> 8), byte(new)},
	}
}

// MutableString returns a MutableRange constructed from the bytes in two
// strings.
func MutableString(addr Addr, old, new string) *MutableRange {
	return &MutableRange{
		Addr: addr,
		Old:  bytes.NewBufferString(old).Bytes(),
		New:  bytes.NewBufferString(new).Bytes(),
	}
}

// Mutate replaces bytes in its range.
func (mr *MutableRange) Mutate(b []byte) error {
	addr := mr.Addr.FullOffset(isEn(b))
	for i, value := range mr.New {
		b[addr+i] = value
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr *MutableRange) Check(b []byte) error {
	addr := mr.Addr.FullOffset(isEn(b))
	for i, value := range mr.Old {
		if b[addr+i] != value {
			return fmt.Errorf("expected %x at %x; found %x",
				mr.Old[i], addr+i, b[addr+i])
		}
	}
	return nil
}

// SetFreewarp sets whether tree warp in the generated ROM will have a
// cooldown (true = no cooldown).
func SetFreewarp(freewarp bool) {
	if freewarp {
		constMutables["tree warp (jp)"].(*MutableRange).New[12] = 0x18
		constMutables["tree warp (en)"].(*MutableRange).New[12] = 0x18
	} else {
		constMutables["tree warp (jp)"].(*MutableRange).New[12] = 0x28
		constMutables["tree warp (en)"].(*MutableRange).New[12] = 0x28
	}
}

// most of the tree warp code between jp and en is the same; only the last two
// instructions (six bytes) differ
const treeWarpCommon = "\xfa\x81\xc4\xe6\x08\x28\x21\x21\x25\xc6\xcb\x7e" +
	"\x28\x06\x3e\x5a\xcd\x74\x0c\xc9\x36\xff\x2b\x36\xfc\x2b\x36\xb4\x2b" +
	"\x36\x40\x21\xb7\xcb\x36\x05\xaf"

// consider these mutables constants; they aren't changed in the randomization
// process.
var constMutables = map[string]Mutable{
	// allow skipping the capcom screen after half a second by pressing start
	"skip capcom call (en)": MutableWord(Addr{0x03, 0, 0x4d6c}, 0x3702, 0xd77d),
	"skip capcom func (en)": MutableString(Addr{0x03, 0, 0x7dd7},
		"\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03\x03",
		"\xe5\xfa\xb3\xcb\xfe\xb2\x30\x03\xcd\x62\x08\xe1\xcd\x37\x02\xc9"),

	// start game with link below bushes, not above
	"initial link placement": MutableByte(sameAddr(0x07, 0x4197), 0x38, 0x58),
	// make link actionable as soon as he drops into the world.
	"link immediately actionable (jp)": MutableString(sameAddr(0x05, 0x4d98),
		"\x3e\x08\xcd\x15", "\xcd\x15\x2a\xc9"),
	"link immediately actionable (en)": MutableString(sameAddr(0x05, 0x4d98),
		"\x3e\x08\xcd\x16", "\xcd\x16\x2a\xc9"),
	// set global flags and room flags that would be set during the intro,
	// overwriting the initial din interaction.
	"set intro flags (jp)": MutableString(sameAddr(0x0a, 0x66ed),
		"\x1e\x78\x1a\xcb\x7f\x20\x08\xe6\x7f\xc4\xb7\x25\xcd\xb7\x25\xcd\x0b\x25\xd0",
		"\x3e\x0a\xcd\xb9\x30\x21\x98\xc7\x36\xc0\x2e\xa7\x36\x50\x2e\xb6\x36\x40\xc9"),
	"set intro flags (en)": MutableString(sameAddr(0x0a, 0x66ed),
		"\x1e\x78\x1a\xcb\x7f\x20\x08\xe6\x7f\xc4\xb8\x25\xcd\xb8\x25\xcd\x0c\x25\xd0",
		"\x3e\x0a\xcd\xcd\x30\x21\x98\xc7\x36\xc0\x2e\xa7\x36\x50\x2e\xb6\x36\x40\xc9"),

	// warp to ember tree if holding start when closing the map screen, using
	// the playtime counter as a cooldown. this requires adding some code at
	// the end of the bank.
	"outdoor map jump redirect (jp)": MutableString(Addr{0x02, 0x60eb, 0x6089},
		"\xc2\xdd\x4f", "\xc4\x1d\x76"),
	"dungeon map jump redirect (jp)": MutableString(Addr{0x02, 0x608e, 0x602c},
		"\xc2\xdd\x4f", "\xc4\x1d\x76"),
	"outdoor map jump redirect (en)": MutableString(Addr{0x02, 0x60eb, 0x6089},
		"\xc2\x7b\x4f", "\xc4\xbb\x75"),
	"dungeon map jump redirect (en)": MutableString(Addr{0x02, 0x608e, 0x602c},
		"\xc2\x7b\x4f", "\xc4\xbb\x75"),
	"tree warp (jp)": MutableString(Addr{0x02, 0x761d, 0},
		"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02",
		treeWarpCommon+"\xcd\xdd\x5e\xc3\xdd\x4f"),
	"tree warp (en)": MutableString(Addr{0x02, 0, 0x75bb},
		"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02",
		treeWarpCommon+"\xcd\x7b\x5e\xc3\x7b\x4f"),

	// have maku gate open from start
	"maku gate check": MutableByte(sameAddr(0x04, 0x61a3), 0x7e, 0x66),

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't disable the flute appearing until actually getting
	// ricky's flute; normally it disappears as soon as you enter the screen
	// northeast of d1 (or ricky's spot, whichever comes first).
	"horon shop stock check":   MutableByte(sameAddr(0x08, 0x4adb), 0x05, 0x02),
	"horon shop sell check":    MutableByte(sameAddr(0x08, 0x48d0), 0x05, 0x02),
	"horon shop flute check 1": MutableByte(sameAddr(0x08, 0x4b02), 0xcb, 0xf6),
	"horon shop flute check 2": MutableByte(sameAddr(0x08, 0x4afc), 0x6f, 0x7f),

	// subrosian dancing's flute prize is normally disabled by visiting the
	// same areas as the horon shop's flute.
	"dance hall flute check": MutableByte(Addr{0x09, 0x5e21, 0x5e38}, 0x20, 0x80),

	// don't require rod to get items from season spirits
	"season spirit rod check": MutableByte(sameAddr(0x0b, 0x4eb2), 0x07, 0x02),

	// i don't know what global flag 0e is. it's only checked in for star ore
	// digging, and disabling the check seems to be sometimes necessary (?)
	"star ore flag check (jp)": MutableString(Addr{0x08, 0x62aa, 0},
		"\xc2\xc5\x3a", "\x00\x00\x00"),
	"star ore flag check (en)": MutableString(Addr{0x08, 0, 0x62aa},
		"\xc2\xd9\x3a", "\x00\x00\x00"),

	// the member's card isn't in the normal logic currently, but remove the
	// essence check anyway
	"member's card essence check": MutableWord(Addr{0x09, 0x7739, 0x7750},
		0xcb57, 0xf601),

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e68, 0x4e72}, 0xcb, 0xf6),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e36, 0x4e40}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e4c, 0x4e56}, 0x2e, 0x04),
	"master essence check 1":    MutableByte(sameAddr(0x0a, 0x4bf5), 0x02, 0x00),
	"master essence check 2":    MutableByte(sameAddr(0x0a, 0x4bea), 0x40, 0x02),
	"master essence check 3":    MutableByte(sameAddr(0x08, 0x5887), 0x40, 0x02),
	"round jewel essence check": MutableByte(sameAddr(0x0a, 0x4f8b), 0x05, 0x00),
	"pirate essence check":      MutableByte(sameAddr(0x08, 0x6c32), 0x20, 0x00),
	"eruption check 1":          MutableByte(sameAddr(0x08, 0x7c41), 0x07, 0x00),
	"eruption check 2":          MutableByte(sameAddr(0x08, 0x7cd3), 0x07, 0x00),

	// stop rosa from spawning and activate her portal by default. the first is
	// an essence check and the second is an edit to tile replacement data. the
	// *third* sets the room to explored before loading its tile replacement
	// data, which ordinarily happens during normal screen transitions but not
	// portal ones. the third one isn't needed in the en/us version and causes
	// problems like getting stuck in doors.
	"rosa spawn check": MutableByte(Addr{0x09, 0x678c, 0x67a3}, 0x40, 0x04),
	"activate rosa portal": &MutableRange{sameAddr(0x04, 0x6016),
		[]byte{0x40, 0x33, 0xc5}, []byte{0x10, 0x33, 0xe6}},
	"set explored before load (jp)": &MutableRange{sameAddr(0x04, 0x5fdf),
		[]byte{0x55, 0x19, 0x4f}, []byte{0x23, 0x2d, 0x4e}},

	// count number of essences, not highest number essence
	"maku seed check 1": MutableByte(Addr{0x09, 0x7d8d, 0x7da4}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7d8f, 0x7da6}, 0x30, 0x18),

	// move sleeping talon and his mushroom so they don't block the chest
	"move talon":    MutableWord(Addr{0x11, 0x6d28, 0x6d2b}, 0x6858, 0x88a8),
	"move mushroom": MutableWord(Addr{0x0b, 0x607f, 0x6080}, 0x6848, 0x78a8),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore (jp)": MutableString(Addr{0x14, 0x4111, 0},
		"\xe0\x2e\x5d", "\xf0\xf0\xf0"),
	"get fools ore (en)": MutableString(Addr{0x14, 0, 0x4881},
		"\xe0\xeb\x58", "\xf0\xf0\xf0"),
	// but always give up feather if the player doesn't have it
	"give stolen feather (jp)": MutableString(Addr{0x15, 0x6212, 0},
		"\xcd\x55\x19\xcb\x6e\x20", "\x3e\x17\xcd\x17\x17\x38"),
	"give stolen feather (en)": MutableString(Addr{0x15, 0, 0x5dcf},
		"\xcd\x56\x19\xcb\x6e\x20", "\x3e\x17\xcd\x17\x17\x38"),
	// and make the feather appear without needing to be dug up
	"stolen feather appears": MutableByte(Addr{0x15, 0x5778, 0x5335}, 0x5a, 0x1a),
	// There are tables indicating extra items to "get" and "lose" upon getting
	// an item. We remove the "lose fools ore" entry and insert a "get seeds
	// from slingshot" entry.
	"lose fools, get seeds from slingshot 1": MutableByte(sameAddr(0x3f, 0x4543), 0x00, 0x13),
	"lose fools, get seeds from slingshot 2": &MutableRange{sameAddr(0x3f, 0x4545),
		[]byte{0x45, 0x00, 0x52, 0x50, 0x51, 0x17, 0x1e, 0x00},
		[]byte{0x20, 0x00, 0x46, 0x45, 0x00, 0x52, 0x50, 0x51}},
	"lose fools, get seeds from slingshot 3": MutableByte(sameAddr(0x3f, 0x44cf), 0x44, 0x47),
	// since slingshot doesn't increment seed capacity, set the level-zero
	// capacity of seeds to 20, and move the pointer up by one byte.
	"satchel capacity": &MutableRange{sameAddr(0x3f, 0x4617),
		[]byte{0x20, 0x50, 0x99}, []byte{0x20, 0x20, 0x50}},
	"satchel capacity pointer": MutableByte(sameAddr(0x3f, 0x460e), 0x16, 0x17),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7baf, 0x7bb9}, 0x05, 0x3f),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1 (jp)": MutableWord(Addr{0x14, 0x4c67, 0}, 0xe951, 0xa94f),
	"noble sword anim 2 (jp)": MutableWord(Addr{0x14, 0x4e37, 0}, 0x8364, 0xdf60),
	"noble sword anim 1 (en)": MutableWord(Addr{0x14, 0, 0x53d7}, 0x5959, 0x1957),
	"noble sword anim 2 (en)": MutableWord(Addr{0x14, 0, 0x55a7}, 0xf36b, 0x4f68),

	// getting the L-2 (or L-3) sword in the lost woods gives you two items;
	// one for the item itself and another that gives you the item and also
	// makes you do a spin slash animation. change the second ID bytes to a
	// fake item so that one slot doesn't give two items / the same item twice.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a, 0x641b}, 0x05, 0x3f),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421, 0x6422}, 0x05, 0x3f),

	// remove the snow piles in front of the shovel house so that shovel isn't
	// required not to softlock there (it's still required not to softlock in
	// hide and seek 2)
	"remove snow piles": MutableByte(Addr{0x24, 0x5dfe, 0x6474}, 0xd9, 0x04),

	// restrict the area triggering sokra to talk to link in horon village to
	// the left side of the burnable trees (prevents softlock)
	"resize sokra trigger": &MutableRange{sameAddr(0x08, 0x5ba5),
		[]byte{0xfa, 0x0b, 0xd0, 0xfe, 0x3c, 0xd8, 0xfe, 0x60, 0xd0},
		[]byte{0xfe, 0x88, 0xd0, 0xfa, 0x0b, 0xd0, 0xfe, 0x3c, 0xd8}},

	// remove one-way diving spot on the south end of sunken city to prevent
	// softlock on moblin road without winter. this requires moving
	// interactions around.
	"remove diving spot": &MutableRange{Addr{0x11, 0x69ca, 0x69cd},
		[]byte{0x1f, 0x0d, 0x68, 0x68, 0x3e, 0x31, 0x18, 0x68},
		[]byte{0x3e, 0x31, 0x18, 0x68, 0xff, 0xff, 0xff, 0xff}},

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move magnet ball": MutableByte(Addr{0x15, 0x53a5, 0x4f62}, 0x48, 0x38),

	// move the trigger for the bridge from holodrum plain to natzu to the
	// top-left corner of the screen, where it can't be hit, and replace the
	// lever tile as well. this prevents the bridge from blocking the waterway.
	"remove bridge trigger": MutableWord(Addr{0x11, 0x6734, 0x6737}, 0x6868, 0x0000),
	"remove bridge lever":   MutableByte(Addr{0x21, 0x5bf1, 0x6267}, 0xb1, 0x04),

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b3, 0x68b5}, 0xb8, 0xbf),

	// the one-way sunken city -> eastern suburbs cliff makes routing
	// complicated. this replaces the flower and wall with stairs, so that the
	// wall can be climbed in all seasons.
	"remove cliff flower":    MutableByte(Addr{0x11, 0x6566, 0x6569}, 0x9c, 0xff),
	"replace cliff spring 1": MutableByte(Addr{0x21, 0x65d5, 0x6c4b}, 0xce, 0xd0),
	"replace cliff spring 2": MutableByte(Addr{0x21, 0x65df, 0x6c55}, 0x54, 0xd0),
	"replace cliff spring 3": MutableByte(Addr{0x21, 0x65e9, 0x6c5f}, 0x2c, 0x04),
	"replace cliff summer 1": MutableByte(Addr{0x22, 0x621c, 0x6892}, 0xce, 0xd0),
	"replace cliff summer 2": MutableByte(Addr{0x22, 0x6226, 0x689c}, 0x54, 0xd0),
	"replace cliff summer 3": MutableByte(Addr{0x22, 0x6230, 0x68a6}, 0x93, 0x04),
	"replace cliff autumn 1": MutableByte(Addr{0x23, 0x6035, 0x66ab}, 0xce, 0xd0),
	"replace cliff autumn 2": MutableByte(Addr{0x23, 0x603f, 0x66b5}, 0x54, 0xd0),
	"replace cliff autumn 3": MutableByte(Addr{0x23, 0x6049, 0x66bf}, 0x93, 0x04),
	"replace cliff winter 1": MutableByte(Addr{0x24, 0x5d15, 0x638b}, 0xce, 0xd0),
	"replace cliff winter 2": MutableByte(Addr{0x24, 0x5d1f, 0x6395}, 0x54, 0xd0),
	"replace cliff winter 3": MutableByte(Addr{0x24, 0x5d29, 0x639f}, 0x93, 0x04),

	// normally none of the desert pits will work if the player already has the
	// rusty bell
	"desert item check": MutableByte(sameAddr(0x08, 0x739e), 0x4a, 0x04),

	// replace the rock/flower outside of d6 with a normal bush so that the
	// player doesn't get softlocked if they exit d6 without gale satchel or
	// default spring.
	"replace d6 flower spring":      MutableByte(Addr{0x21, 0x47fd, 0x4e73}, 0xd8, 0xc4),
	"replace d6 flower summer":      MutableByte(Addr{0x22, 0x450d, 0x4b83}, 0x92, 0xc4),
	"replace d6 flower autumn":      MutableByte(Addr{0x23, 0x42fd, 0x4973}, 0x92, 0xc4),
	"replace d6 flower winter (jp)": MutableByte(Addr{0x23, 0x7f5a, 0}, 0x92, 0xc4),
	"replace d6 flower winter (en)": MutableByte(Addr{0x24, 0, 0x45d0}, 0x92, 0xc4),

	// remove a flower on the way to the spring banana tree, since the player
	// could remove it with moosh and then be stuck behind it. it doesn't lock
	// any items anyway, since only sword can cut the item from the tree.
	"remove mt. cucco flower": MutableByte(Addr{0x21, 0x5287, 0x58fd}, 0xd8, 0x04),

	// replace the stairs outside the portal in eyeglass lake in summer with a
	// railing, because if the player jumps off those stairs in summer they
	// fall into the noble sword room.
	"replace lake stairs": MutableString(Addr{0x22, 0x72a5, 0x791b},
		"\x36\xd0\x35", "\x40\x40\x40"),

	// skip pirate cutscene. adds flag-setting code at the end of the bank.
	// includes setting flag $1b, which makes the pirate skull appear in the
	// desert, in case the player hasn't talked to the ghost.
	"pirate flag call (jp)": MutableWord(Addr{0x15, 0x5e52, 0x5a0f}, 0xb930, 0x707d),
	"pirate flag call (en)": MutableWord(Addr{0x15, 0x5e52, 0x5a0f}, 0xcd30, 0x2d79),
	"pirate flag func (jp)": MutableString(Addr{0x15, 0x7d70, 0x792d},
		"\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15",
		"\xcd\xb9\x30\x3e\x17\xcd\xb9\x30\x3e\x1b\xcd\xb9\x30\x21\xe2\xc7\xcb\xf6\xc9"),
	"pirate flag func (en)": MutableString(Addr{0x15, 0x7d70, 0x792d},
		"\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15",
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7\xcb\xf6\xc9"),
	"pirate warp": MutableString(Addr{0x15, 0x5e5f, 0x5a1c},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),

	// if entering certain warps blocked by snow piles in winter, set the
	// animal companion to appear right outside instead of where you left them.
	// this requires adding some code at the end of the bank.
	"animal save point call": MutableString(sameAddr(0x04, 0x461e),
		"\xfa\x64\xcc", "\xcd\x02\x7e"),
	"set animal save point": MutableString(sameAddr(0x04, 0x7e02),
		"\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04"+
			"\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04"+
			"\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04\x04"+
			"\x04\x04\x04",
		"\xc5\x47\xfa\x64\xcc\x4f\x78\xfe\x04\x20\x05\x79\xfe\xfa\x28\x14\x78"+
			"\xfe\x05\x20\x05\x79\xfe\xcc\x28\x0a\x78\xfe\x01\x20\x11\x79"+
			"\xfe\x57\x20\x0c\xfa\x4c\xcc\x21\x42\xcc\x22\x36\x28\x23\x36\x68"+
			"\x79\xc1\xc9"),

	// moosh won't spawn in the mountains if you have the wrong number of
	// essences. bit 6 seems related to this, and needs to be zero too?
	"skip moosh essence check 1": MutableByte(sameAddr(0x0f, 0x7429), 0x03, 0x00),
	"skip moosh essence check 2": MutableByte(Addr{0x09, 0x4e2c, 0x4e36}, 0xca, 0xc3),
	"skip moosh flag check":      MutableByte(Addr{0x09, 0x4ea3, 0x4ead}, 0x40, 0x00),
}

var mapIconByTreeID = []byte{0x15, 0x19, 0x16, 0x17, 0x18, 0x18}

// like the item slots, these are (usually) no-ops until the randomizer touches
// them.
var varMutables = map[string]Mutable{
	// set initial season correctly in the init variables. this replaces
	// null-terminating whoever's son's name, which *should* be zeroed anyway.
	"initial season": MutableWord(sameAddr(0x07, 0x4188), 0x0e00, 0x2d00),

	// map pop-up icons for seed trees
	"tarm gale tree map icon":   MutableByte(Addr{0x02, 0x6cb3, 0x6c51}, 0x18, 0x18),
	"sunken gale tree map icon": MutableByte(Addr{0x02, 0x6cb6, 0x6c54}, 0x18, 0x18),
	"scent tree map icon":       MutableByte(Addr{0x02, 0x6cb9, 0x6c57}, 0x16, 0x16),
	"pegasus tree map icon":     MutableByte(Addr{0x02, 0x6cbc, 0x6c5a}, 0x17, 0x17),
	"mystery tree map icon":     MutableByte(Addr{0x02, 0x6cbf, 0x6c5d}, 0x19, 0x19),
	"ember tree map icon":       MutableByte(Addr{0x02, 0x6cc2, 0x6c60}, 0x15, 0x15),

	// these scenes use specific item sprites not tied to treasure data
	"wooden sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x65f4, 0x6676},
		Old:  []byte{0x60, 0x00, 0x00},
		New:  []byte{0x60, 0x00, 0x00},
	},
	"rod graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6ba3, 0x6c25},
		Old:  []byte{0x60, 0x10, 0x21},
		New:  []byte{0x60, 0x10, 0x21},
	},
	"noble sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6975, 0x69f7},
		Old:  []byte{0x4e, 0x1a, 0x50},
		New:  []byte{0x4e, 0x1a, 0x50},
	},
	"master sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6978, 0x69fa},
		Old:  []byte{0x4e, 0x1a, 0x40},
		New:  []byte{0x4e, 0x1a, 0x40},
	},

	// the satchel and slingshot should contain the type of seeds that grow on
	// the horon village tree.
	"satchel initial seeds":   MutableByte(sameAddr(0x3f, 0x453b), 0x20, 0x20),
	"slingshot initial seeds": MutableByte(sameAddr(0x3f, 0x4544), 0x46, 0x20),

	// the correct type of seed needs to be selected by default, otherwise the
	// player may be unable to use seeds when they only have one type. there
	// could also be serious problems with the submenu when they *do* obtain a
	// second type if the selection isn't either of them.
	//
	// this works by overwriting a couple of unimportant bytes in file
	// initialization.
	"satchel initial selection":   MutableWord(sameAddr(0x07, 0x418e), 0xa210, 0xbe00),
	"slingshot initial selection": MutableWord(sameAddr(0x07, 0x419a), 0x2e02, 0xbf00),

	// allow seed collection if you have a slingshot, by checking for the given
	// initial seed type
	"carry seeds in slingshot": MutableByte(sameAddr(0x10, 0x4b19), 0x19, 0x20),
}

var Seasons = map[string]*MutableRange{
	// randomize default seasons (before routing). sunken city also applies to
	// mt. cucco; eastern suburbs applies to the vertical part of moblin road
	// but not the horizontal part. note that "tarm ruins" here refers only to
	// the part beyond the lost woods.
	//
	// horon village is random, natzu and desert can only be summer, and goron
	// mountain can only be winter. not sure about northern peak but it doesn't
	// matter.
	"north horon season":     MutableByte(Addr{0x01, 0x7e42, 0x7e60}, 0x03, 0x03),
	"eastern suburbs season": MutableByte(Addr{0x01, 0x7e43, 0x7e61}, 0x02, 0x02),
	"woods of winter season": MutableByte(Addr{0x01, 0x7e44, 0x7e62}, 0x01, 0x01),
	"spool swamp season":     MutableByte(Addr{0x01, 0x7e45, 0x7e63}, 0x02, 0x02),
	"holodrum plain season":  MutableByte(Addr{0x01, 0x7e46, 0x7e64}, 0x00, 0x00),
	"sunken city season":     MutableByte(Addr{0x01, 0x7e47, 0x7e65}, 0x01, 0x01),
	"lost woods season":      MutableByte(Addr{0x01, 0x7e49, 0x7e67}, 0x02, 0x02),
	"tarm ruins season":      MutableByte(Addr{0x01, 0x7e4a, 0x7e68}, 0x00, 0x00),
	"western coast season":   MutableByte(Addr{0x01, 0x7e4d, 0x7e6b}, 0x03, 0x03),
	"temple remains season":  MutableByte(Addr{0x01, 0x7e4e, 0x7e6c}, 0x03, 0x03),
}

// get a collated map of all mutables
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		slotMutables[k] = v
	}
	treasureMutables := make(map[string]Mutable)
	for k, v := range Treasures {
		treasureMutables[k] = v
	}
	seasonMutables := make(map[string]Mutable)
	for k, v := range Seasons {
		seasonMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		constMutables,
		treasureMutables,
		slotMutables,
		varMutables,
		seasonMutables,
	}

	// initialize master map w/ adequate capacity
	count := 0
	for _, set := range mutableSets {
		count += len(set)
	}
	allMutables := make(map[string]Mutable, count)

	// add mutables to master map
	for _, set := range mutableSets {
		for k, v := range set {
			if _, ok := allMutables[k]; ok {
				log.Fatalf("duplicate mutable key: %s", k)
			}
			allMutables[k] = v
		}
	}

	return allMutables
}
