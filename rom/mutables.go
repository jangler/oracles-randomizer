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
	addr := mr.Addr.FullOffset()
	for i, value := range mr.New {
		b[addr+i] = value
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr *MutableRange) Check(b []byte) error {
	addr := mr.Addr.FullOffset()
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
		constMutables["tree warp"].(*MutableRange).New[12] = 0x18
	} else {
		constMutables["tree warp"].(*MutableRange).New[12] = 0x28
	}
}

// consider these mutables constants; they aren't changed in the randomization
// process.
var constMutables = map[string]Mutable{
	// start game with link below bushes, not above
	"initial link placement": MutableByte(Addr{0x07, 0x4197}, 0x38, 0x58),
	// make link actionable as soon as he drops into the world.
	"link immediately actionable": MutableString(Addr{0x05, 0x4d98},
		"\x3e\x08\xcd\x15", "\xcd\x15\x2a\xc9"),
	// set global flags and room flags that would be set during the intro,
	// overwriting the initial din interaction.
	"set intro flags": MutableString(Addr{0x0a, 0x66ed},
		"\x1e\x78\x1a\xcb\x7f\x20\x08\xe6\x7f\xc4\xb7\x25\xcd\xb7\x25\xcd\x0b\x25\xd0",
		"\x3e\x0a\xcd\xb9\x30\x21\x98\xc7\x36\xc0\x2e\xa7\x36\x50\x2e\xb6\x36\x40\xc9"),

	// warp to ember tree if holding start when closing the map screen, using
	// the playtime counter as a cooldown. this requires adding some code at
	// the end of the bank.
	"outdoor map jump redirect": MutableString(Addr{0x02, 0x60eb},
		"\xc2\xdd\x4f", "\xc4\x1d\x76"),
	"dungeon map jump redirect": MutableString(Addr{0x02, 0x608e},
		"\xc2\xdd\x4f", "\xc4\x1d\x76"),
	"tree warp": MutableString(Addr{0x02, 0x761d},
		"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02"+
			"\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02\x02",
		"\xfa\x81\xc4\xe6\x08\x28\x21\x21\x25\xc6\xcb\x7e\x28\x06\x3e\x5a"+
			"\xcd\x74\x0c\xc9\x36\xff\x2b\x36\xfc\x2b\x36\xb4\x2b\x36\x40"+
			"\x21\xb7\xcb\x36\x05\xaf\xcd\xdd\x5e\xc3\xdd\x4f"),

	// have maku gate open from start
	"maku gate check": MutableByte(Addr{0x04, 0x61a3}, 0x7e, 0x66),

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't disable the flute appearing until actually getting
	// ricky's flute; normally it disappears as soon as you enter the screen
	// northeast of d1 (or ricky's spot, whichever comes first).
	"horon shop stock check":   MutableByte(Addr{0x08, 0x4adb}, 0x05, 0x02),
	"horon shop sell check":    MutableByte(Addr{0x08, 0x48d0}, 0x05, 0x02),
	"horon shop flute check 1": MutableByte(Addr{0x08, 0x4b02}, 0xcb, 0xf6),
	"horon shop flute check 2": MutableByte(Addr{0x08, 0x4afc}, 0x6f, 0x7f),

	// subrosian dancing's flute prize is normally disabled by visiting the
	// same areas as the horon shop's flute.
	"dance hall flute check": MutableByte(Addr{0x09, 0x5e21}, 0x20, 0x80),

	// don't require rod to get items from season spirits
	"season spirit rod check": MutableByte(Addr{0x0b, 0x4eb2}, 0x07, 0x02),

	// the member's card isn't in the normal logic currently, but remove the
	// essence check anyway
	"member's card essence check": MutableWord(Addr{0x09, 0x7739},
		0xcb57, 0xf601),

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e68}, 0xcb, 0xf6),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e36}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e4c}, 0x2e, 0x04),
	"master essence check 1":    MutableByte(Addr{0x0a, 0x4bf5}, 0x02, 0x00),
	"master essence check 2":    MutableByte(Addr{0x0a, 0x4bea}, 0x40, 0x02),
	"master essence check 3":    MutableByte(Addr{0x08, 0x5887}, 0x40, 0x02),
	"round jewel essence check": MutableByte(Addr{0x0a, 0x4f8b}, 0x05, 0x00),
	"pirate essence check":      MutableByte(Addr{0x08, 0x6c32}, 0x20, 0x00),
	"eruption check 1":          MutableByte(Addr{0x08, 0x7c41}, 0x07, 0x00),
	"eruption check 2":          MutableByte(Addr{0x08, 0x7cd3}, 0x07, 0x00),

	// stop rosa from spawning and activate her portal by default. the first is
	// an essence check and the second is an edit to tile replacement data. the
	// *third* sets the room to explored before loading its tile replacement
	// data, which ordinarily happens during normal screen transitions but not
	// portal ones.
	"rosa spawn check": MutableByte(Addr{0x09, 0x678c}, 0x40, 0x04),
	"activate rosa portal": &MutableRange{Addr{0x04, 0x6016},
		[]byte{0x40, 0x33, 0xc5}, []byte{0x10, 0x33, 0xe6}},
	"set explored before load": &MutableRange{Addr{0x04, 0x5fdf},
		[]byte{0x55, 0x19, 0x4f}, []byte{0x23, 0x2d, 0x4e}},

	// count number of essences, not highest number essence
	"maku seed check 1": MutableByte(Addr{0x09, 0x7d8d}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7d8f}, 0x30, 0x18),

	// move sleeping talon and his mushroom so they don't block the chest
	"move talon":    MutableWord(Addr{0x11, 0x6d28}, 0x6858, 0x88a8),
	"move mushroom": MutableWord(Addr{0x0b, 0x607f}, 0x6848, 0x78a8),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore 1": MutableByte(Addr{0x14, 0x4111}, 0xe0, 0xf0),
	"get fools ore 2": MutableByte(Addr{0x14, 0x4112}, 0x2e, 0xf0),
	"get fools ore 3": MutableByte(Addr{0x14, 0x4113}, 0x5d, 0xf0),
	// There are tables indicating extra items to "get" and "lose" upon getting
	// an item. We remove the "lose fools ore" entry and insert a "get seeds
	// from slingshot" entry.
	"lose fools, get seeds from slingshot 1": MutableByte(Addr{0x3f, 0x4543}, 0x00, 0x13),
	"lose fools, get seeds from slingshot 2": &MutableRange{Addr{0x3f, 0x4545},
		[]byte{0x45, 0x00, 0x52, 0x50, 0x51, 0x17, 0x1e, 0x00},
		[]byte{0x20, 0x00, 0x46, 0x45, 0x00, 0x52, 0x50, 0x51}},
	"lose fools, get seeds from slingshot 3": MutableByte(Addr{0x3f, 0x44cf}, 0x44, 0x47),
	// since slingshot doesn't increment seed capacity, set the level-zero
	// capacity of seeds to 20, and move the pointer up by one byte.
	"satchel capacity": &MutableRange{Addr{0x3f, 0x4617},
		[]byte{0x20, 0x50, 0x99}, []byte{0x20, 0x20, 0x50}},
	"satchel capacity pointer": MutableByte(Addr{0x3f, 0x460e}, 0x16, 0x17),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7baf}, 0x05, 0x3f),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x4c67}, 0xe951, 0xa94f),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x4e37}, 0x8364, 0xdf60),

	// getting the L-2 (or L-3) sword in the lost woods gives you two items;
	// one for the item itself and another that gives you the item and also
	// makes you do a spin slash animation. change the second ID bytes to a
	// fake item so that one slot doesn't give two items / the same item twice.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0x05, 0x3f),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0x05, 0x3f),

	// remove the snow piles in front of the shovel house so that shovel isn't
	// required not to softlock there (it's still required not to softlock in
	// hide and seek 2)
	"remove snow piles": MutableByte(Addr{0x24, 0x5dfe}, 0xd9, 0x04),

	// restrict the area triggering sokra to talk to link in horon village to
	// the left side of the burnable trees (prevents softlock)
	"resize sokra trigger": &MutableRange{Addr{0x08, 0x5ba5},
		[]byte{0xfa, 0x0b, 0xd0, 0xfe, 0x3c, 0xd8, 0xfe, 0x60, 0xd0},
		[]byte{0xfe, 0x88, 0xd0, 0xfa, 0x0b, 0xd0, 0xfe, 0x3c, 0xd8}},

	// remove one-way diving spot on the south end of sunken city to prevent
	// softlock on moblin road without winter. this requires moving
	// interactions around.
	"remove diving spot": &MutableRange{Addr{0x11, 0x69ca},
		[]byte{0x1f, 0x0d, 0x68, 0x68, 0x3e, 0x31, 0x18, 0x68},
		[]byte{0x3e, 0x31, 0x18, 0x68, 0xff, 0xff, 0xff, 0xff}},

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move magnet ball": MutableByte(Addr{0x15, 0x53a5}, 0x48, 0x38),

	// move the trigger for the bridge from holodrum plain to natzu to the
	// top-left corner of the screen, where it can't be hit, and replace the
	// lever tile as well. this prevents the bridge from blocking the waterway.
	"remove bridge trigger": MutableWord(Addr{0x11, 0x6734}, 0x6868, 0x0000),
	"remove bridge lever":   MutableByte(Addr{0x21, 0x5bf1}, 0xb1, 0x04),

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b3}, 0xb8, 0xbf),

	// the one-way sunken city -> eastern suburbs cliff makes routing
	// complicated. this replaces the flower and wall with stairs, so that the
	// wall can be climbed in all seasons.
	"remove cliff flower":    MutableByte(Addr{0x11, 0x6566}, 0x9c, 0xff),
	"replace cliff spring 1": MutableByte(Addr{0x21, 0x65d5}, 0xce, 0xd0),
	"replace cliff spring 2": MutableByte(Addr{0x21, 0x65df}, 0x54, 0xd0),
	"replace cliff spring 3": MutableByte(Addr{0x21, 0x65e9}, 0x2c, 0x04),
	"replace cliff summer 1": MutableByte(Addr{0x22, 0x621c}, 0xce, 0xd0),
	"replace cliff summer 2": MutableByte(Addr{0x22, 0x6226}, 0x54, 0xd0),
	"replace cliff summer 3": MutableByte(Addr{0x22, 0x6230}, 0x93, 0x04),
	"replace cliff autumn 1": MutableByte(Addr{0x23, 0x6035}, 0xce, 0xd0),
	"replace cliff autumn 2": MutableByte(Addr{0x23, 0x603f}, 0x54, 0xd0),
	"replace cliff autumn 3": MutableByte(Addr{0x23, 0x6049}, 0x93, 0x04),
	"replace cliff winter 1": MutableByte(Addr{0x24, 0x5d15}, 0xce, 0xd0),
	"replace cliff winter 2": MutableByte(Addr{0x24, 0x5d1f}, 0x54, 0xd0),
	"replace cliff winter 3": MutableByte(Addr{0x24, 0x5d29}, 0x93, 0x04),

	// normally none of the desert pits will work if the player already has the
	// rusty bell
	"desert item check": MutableByte(Addr{0x08, 0x739e}, 0x4a, 0x04),

	// replace the rock/flower outside of d6 with a normal bush so that the
	// player doesn't get softlocked if they exit d6 without gale satchel or
	// default spring.
	"replace d6 flower spring": MutableByte(Addr{0x21, 0x47fd}, 0xd8, 0xc4),
	"replace d6 flower summer": MutableByte(Addr{0x22, 0x450d}, 0x92, 0xc4),
	"replace d6 flower autumn": MutableByte(Addr{0x23, 0x42fd}, 0x92, 0xc4),
	"replace d6 flower winter": MutableByte(Addr{0x23, 0x7f5a}, 0x92, 0xc4),

	// replace the stairs outside the portal in eyeglass lake in summer with a
	// railing, because if the player jumps off those stairs in summer they
	// fall into the noble sword room.
	"replace lake stairs": MutableString(Addr{0x22, 0x72a5},
		"\x36\xd0\x35", "\x40\x40\x40"),

	// skip pirate cutscene. adds flag-setting code at the end of the bank.
	// includes setting flag $1b, which makes the pirate skull appear in the
	// desert, in case the player hasn't talked to the ghost.
	"pirate flag call": MutableWord(Addr{0x15, 0x5e52}, 0xb930, 0x707d),
	"pirate flag func": MutableString(Addr{0x15, 0x7d70},
		"\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15\x15",
		"\xcd\xb9\x30\x3e\x17\xcd\xb9\x30\x3e\x1b\xcd\xb9\x30\x21\xe2\xc7\xcb\xf6\xc9"),
	"pirate warp": MutableString(Addr{0x15, 0x5e5f},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),
}

var mapIconByTreeID = []byte{0x15, 0x19, 0x16, 0x17, 0x18, 0x18}

// like the item slots, these are (usually) no-ops until the randomizer touches
// them.
var varMutables = map[string]Mutable{
	// set initial season correctly in the init variables. this replaces
	// null-terminating whoever's son's name, which *should* be zeroed anyway.
	"initial season": MutableWord(Addr{0x07, 0x4188}, 0x0e00, 0x2d00),

	// map pop-up icons for seed trees
	"tarm gale tree map icon":   MutableByte(Addr{0x02, 0x6cb3}, 0x18, 0x18),
	"sunken gale tree map icon": MutableByte(Addr{0x02, 0x6cb6}, 0x18, 0x18),
	"scent tree map icon":       MutableByte(Addr{0x02, 0x6cb9}, 0x16, 0x16),
	"pegasus tree map icon":     MutableByte(Addr{0x02, 0x6cbc}, 0x17, 0x17),
	"mystery tree map icon":     MutableByte(Addr{0x02, 0x6cbf}, 0x19, 0x19),
	"ember tree map icon":       MutableByte(Addr{0x02, 0x6cc2}, 0x15, 0x15),

	// these scenes use specific item sprites not tied to treasure data
	"wooden sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x65f4},
		Old:  []byte{0x60, 0x00, 0x00},
		New:  []byte{0x60, 0x00, 0x00},
	},
	"rod graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6ba3},
		Old:  []byte{0x60, 0x10, 0x21},
		New:  []byte{0x60, 0x10, 0x21},
	},
	"noble sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6975},
		Old:  []byte{0x4e, 0x1a, 0x50},
		New:  []byte{0x4e, 0x1a, 0x50},
	},
	"master sword graphics": &MutableRange{
		Addr: Addr{0x3f, 0x6978},
		Old:  []byte{0x4e, 0x1a, 0x40},
		New:  []byte{0x4e, 0x1a, 0x40},
	},

	// the satchel and slingshot should contain the type of seeds that grow on
	// the horon village tree.
	"satchel initial seeds":   MutableByte(Addr{0x3f, 0x453b}, 0x20, 0x20),
	"slingshot initial seeds": MutableByte(Addr{0x3f, 0x4544}, 0x46, 0x20),

	// the correct type of seed needs to be selected by default, otherwise the
	// player may be unable to use seeds when they only have one type. there
	// could also be serious problems with the submenu when they *do* obtain a
	// second type if the selection isn't either of them.
	//
	// this works by overwriting a couple of unimportant bytes in file
	// initialization.
	"satchel initial selection":   MutableWord(Addr{0x07, 0x418e}, 0xa210, 0xbe00),
	"slingshot initial selection": MutableWord(Addr{0x07, 0x419a}, 0x2e02, 0xbf00),

	// allow seed collection if you have a slingshot, by checking for the given
	// initial seed type
	"carry seeds in slingshot": MutableByte(Addr{0x10, 0x4b19}, 0x19, 0x20),
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
	"north horon season":     MutableByte(Addr{0x01, 0x7e42}, 0x03, 0x03),
	"eastern suburbs season": MutableByte(Addr{0x01, 0x7e43}, 0x02, 0x02),
	"woods of winter season": MutableByte(Addr{0x01, 0x7e44}, 0x01, 0x01),
	"spool swamp season":     MutableByte(Addr{0x01, 0x7e45}, 0x02, 0x02),
	"holodrum plain season":  MutableByte(Addr{0x01, 0x7e46}, 0x00, 0x00),
	"sunken city season":     MutableByte(Addr{0x01, 0x7e47}, 0x01, 0x01),
	"lost woods season":      MutableByte(Addr{0x01, 0x7e49}, 0x02, 0x02),
	"tarm ruins season":      MutableByte(Addr{0x01, 0x7e4a}, 0x00, 0x00),
	"western coast season":   MutableByte(Addr{0x01, 0x7e4d}, 0x03, 0x03),
	"temple remains season":  MutableByte(Addr{0x01, 0x7e4e}, 0x03, 0x03),
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
