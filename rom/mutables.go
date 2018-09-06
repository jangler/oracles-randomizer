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
	Addrs    []Addr
	Old, New []byte
}

// MutableByte returns a special case of MutableRange with a range of a single
// byte.
func MutableByte(addr Addr, old, new byte) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   []byte{old},
		New:   []byte{new},
	}
}

// MutableWord returns a special case of MutableRange with a range of a two
// bytes.
func MutableWord(addr Addr, old, new uint16) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   []byte{byte(old >> 8), byte(old)},
		New:   []byte{byte(new >> 8), byte(new)},
	}
}

// MutableString returns a MutableRange constructed from the bytes in two
// strings.
func MutableString(addr Addr, old, new string) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   bytes.NewBufferString(old).Bytes(),
		New:   bytes.NewBufferString(new).Bytes(),
	}
}

// MutableStrings returns a MutableRange constructed from the bytes in two
// strings, at multiple addresses.
func MutableStrings(addrs []Addr, old, new string) *MutableRange {
	return &MutableRange{
		Addrs: addrs,
		Old:   bytes.NewBufferString(old).Bytes(),
		New:   bytes.NewBufferString(new).Bytes(),
	}
}

// Mutate replaces bytes in its range.
func (mr *MutableRange) Mutate(b []byte) error {
	for _, addr := range mr.Addrs {
		offset := addr.FullOffset()
		for i, value := range mr.New {
			b[offset+i] = value
		}
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr *MutableRange) Check(b []byte) error {
	for _, addr := range mr.Addrs {
		offset := addr.FullOffset()
		for i, value := range mr.Old {
			if b[offset+i] != value {
				return fmt.Errorf("expected %x at %x; found %x",
					mr.Old[i], offset+i, b[offset+i])
			}
		}
	}
	return nil
}

// SetFreewarp sets whether tree warp in the generated ROM will have a
// cooldown (true = no cooldown).
func SetFreewarp(freewarp bool) {
	if freewarp {
		constMutables["tree warp"].(*MutableRange).New[19] = 0x18
	} else {
		constMutables["tree warp"].(*MutableRange).New[19] = 0x28
	}
}

// SetAnimal sets the flute type and Natzu region type based on a companion
// number 1 to 3.
func SetAnimal(companion int) {
	varMutables["animal region"].(*MutableRange).New =
		[]byte{byte(companion + 0x0a)}
}

// consider these mutables constants; they aren't changed in the randomization
// process.
var constMutables = map[string]Mutable{
	// allow skipping the capcom screen after one second by pressing start
	"skip capcom call": MutableWord(Addr{0x03, 0x4d6c}, 0x3702, 0xd77d),
	"skip capcom func": MutableString(Addr{0x03, 0x7dd7}, "\x03",
		"\xe5\xfa\xb3\xcb\xfe\x94\x30\x03\xcd\x62\x08\xe1\xcd\x37\x02\xc9"),

	// start game with link below bushes, not above
	"initial link placement": MutableByte(Addr{0x07, 0x4197}, 0x38, 0x58),
	// make link actionable as soon as he drops into the world.
	"link immediately actionable": MutableString(Addr{0x05, 0x4d98},
		"\x3e\x08\xcd\x16", "\xcd\x16\x2a\xc9"),
	// set global flags and room flags that would be set during the intro,
	// overwriting the initial din interaction.
	"set intro flags": MutableString(Addr{0x0a, 0x66ed},
		"\x1e\x78\x1a\xcb\x7f\x20\x08\xe6\x7f\xc4\xb8\x25\xcd\xb8\x25\xcd"+
			"\x0c\x25\xd0",
		"\x3e\x0a\xcd\xcd\x30\x21\x98\xc7\x36\xc0\x2e\xa7\x36\x50\x2e\xb6"+
			"\x36\x40\xc9"),

	// warp to ember tree if holding start when closing the map screen, using
	// the playtime counter as a cooldown. this also sets the player's respawn
	// point.
	"outdoor map jump redirect": MutableString(Addr{0x02, 0x6089},
		"\xc2\x7b\x4f", "\xc4\xbb\x75"),
	"dungeon map jump redirect": MutableString(Addr{0x02, 0x602c},
		"\xc2\x7b\x4f", "\xc4\xbb\x75"),
	"tree warp": MutableString(Addr{0x02, 0x75bb}, "\x02",
		"\xfa\x81\xc4\xe6\x08\x28\x33"+ // close as normal if start not held
			"\xfa\x49\xcc\xfe\x02\x30\x07"+ // check if indoors
			"\x21\x25\xc6\xcb\x7e\x28\x06"+ // check if cooldown is up
			"\x3e\x5a\xcd\x74\x0c\xc9"+ // play error sound and ret
			"\x21\x22\xc6\x11\xf8\x75\x06\x04\xcd\x5b\x04"+ // copy playtime
			"\x21\x2b\xc6\x11\xfc\x75\x06\x06\xcd\x5b\x04"+ // copy save point
			"\x21\xb7\xcb\x36\x05\xaf\xcd\x7b\x5e\xc3\x7b\x4f"+ // close + warp
			"\x40\xb4\xfc\xff\x00\xf8\x02\x02\x34\x38"), // data for copies

	// warp to room under cursor if wearing developer ring. this goes right
	// after the normal tree warp code (but doesn't fall through from it).
	"dev ring tree warp call": MutableWord(Addr{0x02, 0x5e9b},
		0x890c, 0x0276),
	"dev ring tree warp func": MutableString(Addr{0x02, 0x7602}, "\x02",
		"\xfa\xc5\xc6\xfe\x40\x20\x12\xfa\x49\xcc\xfe\x02\x30\x0b\xf6\x80"+
			"\xea\x63\xcc\xfa\xb6\xcb\xea\x64\xcc\x3e\x03\xcd\x89\x0c\xc9"),

	// if wearing dev ring, warp to animal companion if it's already in the
	// same room when playing the flute.
	"dev ring flute call": MutableWord(Addr{0x09, 0x4e2c}, 0xd93a, 0x4e7f),
	"dev ring flute func": MutableString(Addr{0x09, 0x7f4e}, "\x09",
		"\xd5\xfa\xc5\xc6\xfe\x40\x20\x07"+ // check dev ring
			"\xfa\x04\xd1\xfe\x01\x28\x04"+ // check animal companion
			"\xd1\xc3\xd9\x3a"+ // done
			"\xcd\xc6\x3a\x20\x0c\x36\x05"+ // create poof
			"\x11\x0a\xd0\x2e\x4a\x06\x04\xcd\x5b\x04"+ // move poof
			"\x11\x0a\xd1\x21\x0a\xd0\x06\x04\xcd\x5b\x04"+ // move animal
			"\x18\xde"), // jump to done

	// if wearing dev ring, change season regardless of where link is standing.
	"dev ring season call": MutableString(Addr{0x07, 0x5b75},
		"\xfa\xb6\xcc\xfe\x08", "\xcd\x16\x79\x00\x00"),
	"dev ring season func": MutableString(Addr{0x07, 0x7916}, "\x07",
		"\xfa\xc5\xc6\xfe\x40\xc8\xfa\xb6\xcc\xfe\x08\xc9"),

	// have maku gate open from start
	"maku gate check": MutableByte(Addr{0x04, 0x61a3}, 0x7e, 0x66),

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't stop the flute from appearing because of animal
	// flags, since it probably won't be a flute at all.
	"horon shop stock check":   MutableByte(Addr{0x08, 0x4adb}, 0x05, 0x02),
	"horon shop sell check":    MutableByte(Addr{0x08, 0x48d0}, 0x05, 0x02),
	"horon shop flute check 1": MutableByte(Addr{0x08, 0x4b02}, 0xcb, 0xf6),
	"horon shop flute check 2": MutableWord(Addr{0x08, 0x4afb},
		0xcb6f, 0xafaf),
	// and don't set a ricky flag when buying the "flute"
	"shop no set ricky flag": MutableByte(Addr{0x0b, 0x4826}, 0x20, 0x00),

	// this all has to do with animals and flutes:
	// this edits ricky's script so that he never gives his flute.
	"ricky skip flute script":  MutableByte(Addr{0x0b, 0x6b7a}, 0x0b, 0x7f),
	"don't give ricky's flute": MutableByte(Addr{0x09, 0x6e6c}, 0xc0, 0xc9),
	// this prevents subrosian dancing from giving dimitri's flute.
	"don't give dimitri's flute": MutableByte(Addr{0x09, 0x5e37}, 0xe6, 0xf6),
	// this prevents holodrum plain from changing the animal region.
	"don't change animal region": MutableWord(Addr{0x09, 0x6f79},
		0x3804, 0x1808),
	// this keeps ricky in his pen based on flute, not animal region.
	"keep ricky in pen": MutableString(Addr{0x09, 0x4e77},
		"\x10\xc6\xfe\x0b", "\xaf\xc6\xfe\x01"),
	// and this does the same for saying goodbye once reaching spool swamp.
	"ricky say goodbye": MutableString(Addr{0x09, 0x6ccd},
		"\x10\xc6\xfe\x0b", "\xaf\xc6\xfe\x01"),
	// spawn dimitri in sunken city based on flute, not animal region.
	"spawn dimitri in sunken city": MutableString(Addr{0x09, 0x4e4c},
		"\x10\xc6\xfe\x0c", "\xaf\xc6\xfe\x02"),

	// "activate" a flute by setting its icon and song when obtained. also
	// activates the corresponding animal companion.
	"flute set icon call": MutableWord(Addr{0x3f, 0x452c}, 0x4e45, 0x4d71),
	"flute set icon func": MutableString(Addr{0x3f, 0x714d}, "\x3f",
		"\xf5\xd5\xe5\x78\xfe\x0e\x20\x0d\x1e\xaf\x79\xd6\x0a\x12\xc6\x42"+
			"\x26\xc6\x6f\xcb\xfe\xe1\xd1\xf1\xcd\x4e\x45\xc9"),

	// don't require rod to get items from season spirits
	"season spirit rod check": MutableByte(Addr{0x0b, 0x4eb2}, 0x07, 0x02),

	// i don't know what global flag 0e is. it's only checked in for star ore
	// digging, and disabling the check seems to be sometimes necessary (?)
	"star ore flag check": MutableString(Addr{0x08, 0x62aa},
		"\xc2\xd9\x3a", "\x00\x00\x00"),

	// sell member's card in subrosian market before completing d3
	"member's card essence check": MutableWord(Addr{0x09, 0x7750},
		0xcb57, 0xf601),

	// give member's card, treasure map, fool's ore, and identified flutes
	// graphics in treasure sprite table
	"member's card gfx": MutableString(Addr{0x3f, 0x67b4},
		"\x00\x00\x00", "\x5d\x0c\x13"),
	"treasure map gfx": MutableString(Addr{0x3f, 0x67b7},
		"\x00\x00\x00", "\x65\x14\x33"),
	"fool's ore gfx": MutableString(Addr{0x3f, 0x67ba},
		"\x00\x00\x00", "\x60\x14\x00"),
	"ricky's flute gfx": MutableString(Addr{0x3f, 0x67bd},
		"\x00\x00\x00", "\x5f\x16\x33"),
	"dimitri's flute gfx": MutableString(Addr{0x3f, 0x67c0},
		"\x00\x00\x00", "\x5f\x16\x23"),
	"moosh's flute gfx": MutableString(Addr{0x3f, 0x67c3},
		"\x00\x00\x00", "\x5f\x16\x13"),
	"rare peach stone gfx": MutableString(Addr{0x3f, 0x67c6},
		"\x00\x00\x00", "\x5d\x10\x26"),
	"ribbon gfx": MutableString(Addr{0x3f, 0x67c9},
		"\x00\x00\x00", "\x65\x0c\x23"),

	// initiate all these events without requiring essences
	"ricky spawn check":         MutableByte(Addr{0x09, 0x4e72}, 0xcb, 0xf6),
	"dimitri essence check":     MutableByte(Addr{0x09, 0x4e40}, 0xcb, 0xf6),
	"dimitri flipper check":     MutableByte(Addr{0x09, 0x4e56}, 0x2e, 0x04),
	"master essence check 1":    MutableByte(Addr{0x0a, 0x4bf5}, 0x02, 0x00),
	"master essence check 2":    MutableByte(Addr{0x0a, 0x4bea}, 0x40, 0x02),
	"master essence check 3":    MutableByte(Addr{0x08, 0x5887}, 0x40, 0x02),
	"round jewel essence check": MutableByte(Addr{0x0a, 0x4f8b}, 0x05, 0x00),
	"pirate essence check":      MutableByte(Addr{0x08, 0x6c32}, 0x20, 0x00),
	"eruption check 1":          MutableByte(Addr{0x08, 0x7c41}, 0x07, 0x00),
	"eruption check 2":          MutableByte(Addr{0x08, 0x7cd3}, 0x07, 0x00),

	// stop rosa from spawning and activate her portal by default. the first is
	// an essence check, and the second is a function that sets the portal's
	// room flags to do the tile replacement.
	"rosa spawn check": MutableByte(Addr{0x09, 0x67a3}, 0x40, 0x04),
	"call set portal room flag": MutableString(Addr{0x04, 0x45f5},
		"\xfa\x64\xcc", "\xcd\x4b\x7e"),
	"set portal room flag func": MutableString(Addr{0x04, 0x7e4b}, "\x04",
		"\xe5\x21\x9a\xc7\x7e\xf6\xc0\x77\xe1\xfa\x64\xcc\xc9"),

	// count number of essences, not highest number essence
	"maku seed check 1": MutableByte(Addr{0x09, 0x7da4}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7da6}, 0x30, 0x18),

	// move sleeping talon and his mushroom so they don't block the chest
	"move talon":    MutableWord(Addr{0x11, 0x6d2b}, 0x6858, 0x88a8),
	"move mushroom": MutableWord(Addr{0x0b, 0x6080}, 0x6848, 0x78a8),

	// feather game: don't give fools ore, and don't return fools ore
	"get fools ore": MutableString(Addr{0x14, 0x4881},
		"\xe0\xeb\x58", "\xf0\xf0\xf0"),
	// but always give up feather if the player doesn't have it
	"give stolen feather": MutableString(Addr{0x15, 0x5dcf},
		"\xcd\x56\x19\xcb\x6e\x20", "\x3e\x17\xcd\x17\x17\x38"),
	// and make the feather appear without needing to be dug up
	"stolen feather appears": MutableByte(Addr{0x15, 0x5335}, 0x5a, 0x1a),
	// AND allow transition away from the screen once the feather is retrieved
	// (not once the hole is dug)
	"leave H&S screen": MutableString(Addr{0x09, 0x65a0},
		"\xcd\x32\x14\x1e\x49\x1a\xbe", "\xcd\x56\x19\xcb\x6e\x00\x00"),

	// There are tables indicating extra items to "get" and "lose" upon getting
	// an item. We remove the "lose fools ore" entry and insert a "get seeds
	// from slingshot" entry.
	"lose fools, get seeds from slingshot 1": MutableByte(Addr{0x3f, 0x4543},
		0x00, 0x13),
	"lose fools, get seeds from slingshot 2": MutableString(Addr{0x3f, 0x4545},
		"\x45\x00\x52\x50\x51\x17\x1e\x00", "\x20\x00\x46\x45\x00\x52\x50\x51"),
	"lose fools, get seeds from slingshot 3": MutableByte(Addr{0x3f, 0x44cf},
		0x44, 0x47),

	// since slingshot doesn't increment seed capacity, set the level-zero
	// capacity of seeds to 20, and move the pointer up by one byte.
	"satchel capacity": MutableString(Addr{0x3f, 0x4617},
		"\x20\x50\x99", "\x20\x20\x50"),
	"satchel capacity pointer": MutableByte(Addr{0x3f, 0x460e}, 0x16, 0x17),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7bb9}, 0x05, 0x3f),

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x53d7}, 0x5959, 0x1957),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x55a7}, 0xf36b, 0x4f68),

	// getting the L-2 (or L-3) sword in the lost woods gives you two items;
	// one for the item itself and another that gives you the item and also
	// makes you do a spin slash animation. change the second ID bytes to a
	// fake item so that one slot doesn't give two items / the same item twice.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641b}, 0x05, 0x3f),
	"master sword second item": MutableByte(Addr{0x0b, 0x6422}, 0x05, 0x3f),

	// remove the snow piles in front of the shovel house so that shovel isn't
	// required not to softlock there (it's still required not to softlock in
	// hide and seek 2)
	"remove snow piles": MutableByte(Addr{0x24, 0x6474}, 0xd9, 0x04),

	// restrict the area triggering sokra to talk to link in horon village to
	// the left side of the burnable trees (prevents softlock)
	"resize sokra trigger": MutableString(Addr{0x08, 0x5ba5},
		"\xfa\x0b\xd0\xfe\x3c\xd8\xfe\x60\xd0",
		"\xfe\x88\xd0\xfa\x0b\xd0\xfe\x3c\xd8"),

	// remove one-way diving spot on the south end of sunken city to prevent
	// softlock on moblin road without winter. this requires moving
	// interactions around.
	"remove diving spot": MutableString(Addr{0x11, 0x69cd},
		"\x1f\x0d\x68\x68\x3e\x31\x18\x68", "\x3e\x31\x18\x68\xff\xff\xff\xff"),

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move magnet ball": MutableByte(Addr{0x15, 0x4f62}, 0x48, 0x38),

	// move the trigger for the bridge from holodrum plain to natzu to the
	// top-left corner of the screen, where it can't be hit, and replace the
	// lever tile as well. this prevents the bridge from blocking the waterway.
	"remove bridge trigger": MutableWord(Addr{0x11, 0x6737},
		0x6868, 0x0000),
	"remove prairie bridge lever": MutableByte(Addr{0x21, 0x6267},
		0xb1, 0x04),
	"remove wasteland bridge lever": MutableByte(Addr{0x23, 0x5cb7},
		0xb1, 0x04),

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b5}, 0xb8, 0xbf),

	// block the sunken city / eastern suburbs cliff with a spring flower, and
	// place a stump at the top so that you can still travel down the cliff if
	// you have spring.
	"block cliff 1": MutableStrings([]Addr{{0x21, 0x6c2b}, {0x22, 0x6872},
		{0x23, 0x668b}, {0x24, 0x636b}}, "\x5d\x5e", "\x6d\x6e"),
	"block cliff 2": MutableStrings([]Addr{{0x21, 0x6c33}, {0x22, 0x687a},
		{0x23, 0x6693}, {0x24, 0x6373}},
		"\x47\x12\x6d\x11\x5f", "\x1f\x20\x21\x04\x04"),
	"block cliff spring 3": MutableString(Addr{0x21, 0x6c3d},
		"\x52\x12\x12\x5d\x11", "\x22\x23\x24\x04\xd8"),
	"block cliff non-spring 3": MutableStrings([]Addr{{0x22, 0x6884},
		{0x23, 0x669d}, {0x24, 0x637d}},
		"\x52\x12\x12\x5d\x11", "\x22\x23\x24\x04\x92"),
	"block cliff 4": MutableStrings([]Addr{{0x21, 0x6c47}, {0x22, 0x688e},
		{0x23, 0x66a7}, {0x24, 0x6387}},
		"\x62\x40\x40\x58\xce", "\x40\x40\x40\x58\xcf"),
	"block cliff 5": MutableStrings([]Addr{{0x21, 0x6c51}, {0x22, 0x6898},
		{0x23, 0x66b1}, {0x24, 0x6391}}, "\x50", "\x54"),

	// extend the railing on moblin keep to prevent a one-way jump down to
	// sunken city (player needs feather to get back to natzu river). one
	// address per natzu region, then one for the ruined version.
	"moblin keep rail 1": MutableStrings([]Addr{{0x21, 0x63e6}, {0x22, 0x603e},
		{0x23, 0x5e44}, {0x24, 0x5ba6}}, "\xb3", "\x48"),
	"moblin keep rail 2": MutableStrings([]Addr{{0x21, 0x63ee}, {0x22, 0x6046},
		{0x23, 0x5e4c}, {0x24, 0x5baf}}, "\x26", "\x53"),
	"moblin keep rail 3": MutableStrings([]Addr{{0x21, 0x63f8}, {0x22, 0x6050},
		{0x23, 0x5e56}, {0x24, 0x5bb9}}, "\x26", "\x53"),
	"moblin keep rail 4": MutableStrings([]Addr{{0x21, 0x63ff}, {0x22, 0x6057},
		{0x23, 0x5e5d}, {0x24, 0x5bc3}}, "\x48", "\x53"),
	// and remove the cannon near the stairs so that players without flippers
	// can exit (if they arrived by jumping and ran out of pegasus seeds).
	"remove keep cannon object": MutableByte(Addr{0x11, 0x6563}, 0xf8, 0xff),
	"replace moblin keep cannon tiles": MutableStrings([]Addr{{0x21, 0x6bee},
		{0x22, 0x6835}, {0x23, 0x664e}},
		"\xa4\x06\x18\xb9\xa5\xb2\x0d\x1c\xf2\x1a\x25\xb5\xb6",
		"\xb9\x06\x18\xb9\xb9\xb2\x0d\x1c\xf2\x1a\x25\xb9\xb9"),
	"replace ruined keep cannon tiles": MutableString(Addr{0x24, 0x632c},
		"\xa6\x04\x08\x83\xa7\xb9\xb2\x0d\x1c\xf2\x1a\x25\xa9\xb6",
		"\xb9\x04\x08\x83\xb9\xb9\xb2\x0d\x1c\xf2\x1a\x25\xb9\xb9"),

	// normally none of the desert pits will work if the player already has the
	// rusty bell
	"desert item check": MutableByte(Addr{0x08, 0x739e}, 0x4a, 0x04),

	// replace the rock/flower outside of d6 with a normal bush so that the
	// player doesn't get softlocked if they exit d6 without gale satchel or
	// default spring.
	"replace d6 flower spring": MutableByte(Addr{0x21, 0x4e73}, 0xd8, 0xc4),
	"replace d6 flower non-spring": MutableStrings(
		[]Addr{{0x22, 0x4b83}, {0x23, 0x4973}, {0x24, 0x45d0}},
		"\x92", "\xc4"),

	// remove a flower on the way to the spring banana tree, since the player
	// could remove it with moosh and then be stuck behind it. it doesn't lock
	// any items anyway, since only sword can cut the item from the tree.
	"remove mt. cucco flower": MutableByte(Addr{0x21, 0x58fd}, 0xd8, 0x04),

	// replace the stairs outside the portal in eyeglass lake in summer with a
	// railing, because if the player jumps off those stairs in summer they
	// fall into the noble sword room.
	"replace lake stairs": MutableString(Addr{0x22, 0x791b},
		"\x36\xd0\x35", "\x40\x40\x40"),

	// skip pirate cutscene. adds flag-setting code at the end of the bank.
	// includes setting flag $1b, which makes the pirate skull appear in the
	// desert, in case the player hasn't talked to the ghost.
	"pirate flag call": MutableWord(Addr{0x15, 0x5a0f}, 0xcd30, 0x2d79),
	"pirate flag func": MutableString(Addr{0x15, 0x792d}, "\x15",
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7\xcb\xf6"+
			"\xfa\x46\x79\xea\x4e\xcc\xc9"),
	"pirate warp": MutableString(Addr{0x15, 0x5a1c},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),

	// if entering certain warps blocked by snow piles or mushrooms, set the
	// animal companion to appear right outside instead of where you left them.
	// this requires adding some code at the end of the bank.
	"animal save point call": MutableString(Addr{0x04, 0x461e},
		"\xfa\x64\xcc", "\xcd\x02\x7e"),
	"set animal save point func": MutableString(Addr{0x04, 0x7e02}, "\x04",
		// b = group, c = room, hl = table
		"\xc5\xd5\x47\xfa\x64\xcc\x4f\x21\x2c\x7e"+
			"\x2a\xb8\x20\x0e\x2a\xb9\x20\x0a"+ // check group and room
			"\x11\x42\xcc\x06\x03\xcd\x62\x04\x18\x0a"+ // set save pt, done
			"\x2a\xb7\x20\xfc\x7e\x3c\x28\x02\x18\xe4"+ // go to next table entry
			"\x79\xd1\xc1\xc9"), // done
	// table entries are {entered group, entered room, saved room, saved y,
	// saved x}.
	"animal save point table": MutableString(Addr{0x04, 0x7e2c}, "\x04",
		"\x04\xfa\xc2\x18\x68\x00"+ // square jewel cave
			"\x01\x57\xb0\x48\x48\x00"+ // spool swamp subrosia portal
			"\x05\xcc\x29\x28\x68\x00"+ // goron mountain cave
			"\x05\xb3\x8e\x58\x88\x00"+ // cave outside d2
			"\x04\xe1\x86\x48\x68\x00"+ // quicksand ring cave
			"\xff"), // end of table

	// moosh won't spawn in the mountains if you have the wrong number of
	// essences. bit 6 seems related to this, and needs to be zero too?
	"skip moosh essence check 1": MutableByte(Addr{0x0f, 0x7429}, 0x03, 0x00),
	"skip moosh essence check 2": MutableByte(Addr{0x09, 0x4e36}, 0xca, 0xc3),
	"skip moosh flag check":      MutableByte(Addr{0x09, 0x4ead}, 0x40, 0x00),

	// don't warp link using gale seeds if no trees have been reached (the menu
	// gets stuck in an infinite loop)
	"call gale seed check": MutableString(Addr{0x07, 0x4f45},
		"\xfa\x50\xcc\x3d", "\xcd\xf0\x78\x00"),
	"gale seed check": MutableString(Addr{0x07, 0x78f0}, "\x07",
		"\xfa\x50\xcc\x3d\xc0\xaf\x21\xf8\xc7\xb6\x21\x9e\xc7\xb6\x21\x72\xc7"+
			"\xb6\x21\x67\xc7\xb6\x21\x5f\xc7\xb6\x21\x10\xc7\xb6\xcb\x67"+
			"\x20\x02\x3c\xc9\xaf\xc9"),

	// end maku seed script as soon as link gets the seed
	"abbreviate maku seed cutscene": MutableString(Addr{0x0b, 0x71ec},
		"\xe1\x23\x61\x01", "\xb6\x19\xbe\x00"),
	// end northen peak barrier cutscene as soon as the barrier is broken
	"abbreviate barrier cutscene": MutableString(Addr{0x0b, 0x79f1},
		"\x88\x18\x50\xf8", "\xb6\x1d\xbe\x00"),

	// skip shield check for forging hard ore
	"skip iron shield check": MutableByte(Addr{0x0b, 0x75c7}, 0x01, 0x02),
	// and skip the check for what level shield you currently have
	"skip iron shield level check": MutableString(Addr{0x15, 0x62ac},
		"\x38\x01", "\x18\x05"),
	// for some reason the param of iron shield is incremented before giving it
	"keep iron shield param": MutableByte(Addr{0x09, 0x42de}, 0x0c, 0x00),
}

var mapIconByTreeID = []byte{0x15, 0x19, 0x16, 0x17, 0x18, 0x18}

// like the item slots, these are (usually) no-ops until the randomizer touches
// them.
var varMutables = map[string]Mutable{
	// set initial season correctly in the init variables. this replaces
	// null-terminating whoever's son's name, which *should* be zeroed anyway.
	"initial season": MutableWord(Addr{0x07, 0x4188}, 0x0e00, 0x2d00),

	// map pop-up icons for seed trees
	"tarm gale tree map icon":   MutableByte(Addr{0x02, 0x6c51}, 0x18, 0x18),
	"sunken gale tree map icon": MutableByte(Addr{0x02, 0x6c54}, 0x18, 0x18),
	"scent tree map icon":       MutableByte(Addr{0x02, 0x6c57}, 0x16, 0x16),
	"pegasus tree map icon":     MutableByte(Addr{0x02, 0x6c5a}, 0x17, 0x17),
	"mystery tree map icon":     MutableByte(Addr{0x02, 0x6c5d}, 0x19, 0x19),
	"ember tree map icon":       MutableByte(Addr{0x02, 0x6c60}, 0x15, 0x15),

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

	// determines what natzu looks like and what animal the flute calls
	"animal region": MutableByte(Addr{0x07, 0x41a6}, 0x0b, 0x0b),

	// should be set to match the western coast season
	"season after pirate cutscene": MutableByte(Addr{0x15, 0x7946}, 0x15, 0x15),
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
	"north horon season":     MutableByte(Addr{0x01, 0x7e60}, 0x03, 0x03),
	"eastern suburbs season": MutableByte(Addr{0x01, 0x7e61}, 0x02, 0x02),
	"woods of winter season": MutableByte(Addr{0x01, 0x7e62}, 0x01, 0x01),
	"spool swamp season":     MutableByte(Addr{0x01, 0x7e63}, 0x02, 0x02),
	"holodrum plain season":  MutableByte(Addr{0x01, 0x7e64}, 0x00, 0x00),
	"sunken city season":     MutableByte(Addr{0x01, 0x7e65}, 0x01, 0x01),
	"lost woods season":      MutableByte(Addr{0x01, 0x7e67}, 0x02, 0x02),
	"tarm ruins season":      MutableByte(Addr{0x01, 0x7e68}, 0x00, 0x00),
	"western coast season":   MutableByte(Addr{0x01, 0x7e6b}, 0x03, 0x03),
	"temple remains season":  MutableByte(Addr{0x01, 0x7e6c}, 0x03, 0x03),
}

// get a collated map of all mutables
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	treasureMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		if v.Treasure.addr != 0 {
			treasureMutables[FindTreasureName(v.Treasure)] = v.Treasure
		}
		slotMutables[k] = v
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
