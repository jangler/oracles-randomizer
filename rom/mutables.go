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

// SetNoMusic sets music off in the modified rom
func SetNoMusic() {
	constMutables["no music call"].(*MutableRange).New = []byte("\xcd\xc8\x3e")
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

	// don't play any music if the -nomusic flag is given
	"no music call": MutableString(Addr{0x00, 0x0c76},
		"\x67\xf0\xb5", "\x67\xf0\xb5"), // modified only by SetNoMusic()
	"no music func": MutableString(Addr{0x00, 0x3ec8}, "\x00",
		"\x67\xfe\x40\x30\x03\x3e\x08\xc9\xf0\xb5\xc9"),

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

	// animals called by flute normally veto any nonzero collision value for
	// the purposes of entering a screen, but this allows double-wide bridges
	// (1a and 1b) as well. this specifically fixes the problem of not being
	// able to call an animal on the d1 screen, or on the bridge to the screen
	// to the right. the vertical collision check isn't modified, since bridges
	// only run horizontally.
	"flute collision call horizontal": MutableStrings([]Addr{{0x09, 0x4d9a},
		{0x09, 0x4dad}}, "\xcd\xd9\x4e", "\xcd\x7f\x7f"),
	"flute collision func": MutableString(Addr{0x09, 0x7f7f}, "\x09",
		"\x06\x01\x7e\xfe\x1a\x28\x06\xfe\x1b\x28\x02\xb7\xc0"+ // first tile
			"\x7d\x80\x6f\x7e\xfe\x1a\x28\x05\xfe\x1b\x28\x01\xb7"+ // second
			"\x7d\xc0\xcd\x89\x20\xaf\xc9"), // vanilla stuff
	// also need to do this so that dimitri and moosh don't immediately stop
	// walking at the edge of the screen. and do ricky for consistency.
	"ricky flute enter call": MutableString(Addr{0x05, 0x71ea},
		"\xcd\xaa\x44\xb7", "\xcd\x2d\x7e\x00"),
	"dimitri/moosh flute enter call": MutableString(Addr{0x05, 0x493b},
		"\xcd\xaa\x44\xb7", "\xcd\x2d\x7e\x00"),
	"flute enter func": MutableString(Addr{0x05, 0x7e2d}, "\x05",
		"\xcd\xaa\x44\xb7\xc8\xfe\x1a\xc8\xfe\x1b\xc9"),

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
	// spawn dimitri and kids in sunken city based on flute, not animal region.
	"spawn dimitri in sunken city": MutableStrings(
		[]Addr{{0x09, 0x4e4c}, {0x09, 0x6f08}, {0x09, 0x737e}},
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
	// a vanilla bug lets star ore be dug up on the first screen even if you
	// already have the item. so… make first try a second instance of second
	// try.
	"star ore bugfix": MutableWord(Addr{0x08, 0x62d5}, 0x6656, 0x7624),

	// remove star ore from inventory when buying the first subrosian market
	// item. this can't go in the gain/lose items table, since the given item
	// doesn't necessarily have a unique ID.
	"remove traded star ore call": MutableString(Addr{0x09, 0x7887},
		"\xdf\x2a\x4e", "\xcd\xa0\x7f"),
	"remove traded star ore func": MutableString(Addr{0x09, 0x7fa0}, "\x09",
		"\xb7\x20\x07\xe5\x21\x9a\xc6\xcb\xae\xe1\xdf\x2a\x4e\xc9"),

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

	// set room flags so that rosa never appears in the overworld, and her
	// portal is activated by default.
	"set portal room flag call": MutableString(Addr{0x04, 0x45f5},
		"\xfa\x64\xcc", "\xcd\x6c\x7e"),
	"set portal room flag func": MutableString(Addr{0x04, 0x7e6c}, "\x04",
		"\xe5\x21\x9a\xc7\x7e\xf6\x60\x77\x2e\xcb\x7e\xf6\xc0\x77"+ // set flags
			"\xe1\xfa\x64\xcc\xc9"), // do what the address normally does
	// a hack so that a different flag can be used to set the portal tile
	// replacement, allowing the bush-breaking warning interaction to be used
	// on this screen.
	"portal tile replacement": MutableString(Addr{0x04, 0x6016},
		"\x40\x33\xc5", "\x20\x33\xe6"),

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

	// getting the L-2 (or L-3) sword in the lost woods normally gives a second
	// "spin slash" item. remove this from the script.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0xde, 0xc1),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0xde, 0xc1),

	// remove the snow piles in front of the shovel house so that shovel isn't
	// required not to softlock there (it's still required not to softlock in
	// hide and seek 2)
	"remove snow piles": MutableByte(Addr{0x24, 0x6474}, 0xd9, 0x04),

	// restrict the area triggering sokra to talk to link in horon village to
	// the left side of the burnable trees (prevents softlock)
	"resize sokra trigger": MutableString(Addr{0x08, 0x5ba5},
		"\xfa\x0b\xd0\xfe\x3c\xd8\xfe\x60\xd0",
		"\xfe\x88\xd0\xfa\x0b\xd0\xfe\x3c\xd8"),

	// you can softlock in d6 misusing keys without magnet gloves, so just move
	// the magnet ball onto the button it needs to press to get the key the
	// speedrun skips.
	"move d6 magnet ball": MutableByte(Addr{0x15, 0x4f36}, 0x98, 0x58),

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move d8 magnet ball": MutableByte(Addr{0x15, 0x4f62}, 0x48, 0x38),

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

	// block the waterfalls from mt cucco to sunken city, so that there only
	// needs to be one warning interaction at the vines.
	"block waterfalls": MutableStrings([]Addr{{0x21, 0x5bd1}, {0x21, 0x5c17},
		{0x22, 0x58a4}, {0x22, 0x58ea}, {0x23, 0x5645}, {0x23, 0x568b},
		{0x24, 0x54fa}, {0x24, 0x5540}}, "\x36\xff\x35", "\x40\x40\x40"),

	// extend the railing on moblin keep to require only one warning
	// interaction for the potential one-way jump in dimitri's region. one
	// address per natzu region, then one for the ruined version.
	"moblin keep rail 1": MutableStrings([]Addr{{0x21, 0x63f8}, {0x22, 0x6050},
		{0x23, 0x5e56}, {0x24, 0x5bb9}}, "\x26", "\x48"),
	"moblin keep rail 2": MutableStrings([]Addr{{0x21, 0x63ff}, {0x22, 0x6057},
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

	// replace the stairs outside the portal in eyeglass lake in summer with a
	// railing, because if the player jumps off those stairs in summer they
	// fall into the noble sword room.
	"replace lake stairs": MutableString(Addr{0x22, 0x791b},
		"\x36\xd0\x35", "\x40\x40\x40"),

	// replace some currents in spool swamp in spring so that the player isn't
	// trapped by them.
	"replace currents 1": MutableWord(Addr{0x21, 0x7ab1}, 0xd2d2, 0xd3d3),
	"replace currents 2": MutableString(Addr{0x21, 0x7ab6},
		"\xd3\xd2\xd2", "\xd4\xd4\xd4"),
	"replace currents 3": MutableByte(Addr{0x21, 0x7abe}, 0xd3, 0xd1),

	// skip pirate cutscene. adds flag-setting code at the end of the bank.
	// includes setting flag $1b, which makes the pirate skull appear in the
	// desert, in case the player hasn't talked to the ghost.
	"pirate flag call": MutableWord(Addr{0x15, 0x5a0f}, 0xcd30, 0x2d79),
	"pirate flag func": MutableString(Addr{0x15, 0x792d}, "\x15",
		"\xcd\xcd\x30\x3e\x17\xcd\xcd\x30\x3e\x1b\xcd\xcd\x30\x21\xe2\xc7\xcb\xf6"+
			"\xfa\x46\x79\xea\x4e\xcc\xc9"),
	"pirate warp": MutableString(Addr{0x15, 0x5a1c},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),

	// if entering certain warps blocked by snow piles, mushrooms, or bushes,
	// set the animal companion to appear right outside instead of where you
	// left them. this requires adding some code at the end of the bank.
	"animal save point call": MutableString(Addr{0x04, 0x461e},
		"\xfa\x64\xcc", "\xcd\x02\x7e"),
	"animal save point func": MutableString(Addr{0x04, 0x7e02}, "\x04",
		// b = group, c = room, d = animal room, hl = table
		"\xc5\xd5\x47\xfa\x64\xcc\x4f\xfa\x42\xcc\x57\x21\x34\x7e"+
			"\x2a\xb8\x20\x12\x2a\xb9\x20\x0e\x7e\xba\x20\x0a"+ // check criteria
			"\x11\x42\xcc\x06\x03\xcd\x62\x04\x18\x0a"+ // set save pt, done
			"\x2a\xb7\x20\xfc\x7e\x3c\x28\x02\x18\xe0"+ // go to next table entry
			"\x79\xd1\xc1\xc9"), // done
	// table entries are {entered group, entered room, animal room, saved y,
	// saved x}.
	"animal save point table": MutableString(Addr{0x04, 0x7e34}, "\x04",
		"\x04\xfa\xc2\x18\x68\x00"+ // square jewel cave
			"\x05\xcc\x2a\x38\x18\x00"+ // goron mountain cave
			"\x05\xb3\x8e\x58\x88\x00"+ // cave outside d2
			"\x04\xe1\x86\x48\x68\x00"+ // quicksand ring cave
			"\x05\xc9\x2a\x38\x18\x00"+ // goron mountain main
			"\x05\xba\x2f\x18\x68\x00"+ // spring banana cave
			"\x05\xbb\x2f\x18\x68\x00"+ // joy ring cave
			"\x01\x05\x9a\x38\x48\x00"+ // rosa portal
			"\x04\x39\x8d\x38\x38\x00"+ // d2 entrance
			"\xff"), // end of table

	// moosh won't spawn in the mountains if you have the wrong number of
	// essences. bit 6 seems related to this, and needs to be zero too?
	"skip moosh essence check 1": MutableByte(Addr{0x0f, 0x7429}, 0x03, 0x00),
	"skip moosh essence check 2": MutableByte(Addr{0x09, 0x4e36}, 0xca, 0xc3),
	"skip moosh flag check":      MutableByte(Addr{0x09, 0x4ead}, 0x40, 0x00),

	// remove the moosh and dimitri events in spool swamp.
	"prevent moosh cutscene":   MutableByte(Addr{0x11, 0x6572}, 0xf1, 0xff),
	"prevent dimitri cutscene": MutableByte(Addr{0x11, 0x68d4}, 0xf1, 0xff),

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

	// overwrite unused maku gate interaction with warning interaction
	"warning script pointer": MutableWord(Addr{0x08, 0x5663}, 0x874e, 0x6d7f),
	"warning script": MutableString(Addr{0x0b, 0x7f6d}, "\x0b",
		"\xd0\xe0\x47\x79\xa0\xbd\xd7\x3c"+ // wait for collision and animation
			"\x87\xe0\xcf\x7e\x7f\x83\x7f\x88\x7f"+ // jump based on cfe0 bits
			"\x98\x26\x00\xbe\x00"+ // show cliff warning text
			"\x98\x26\x01\xbe\x00"+ // show bush warning text
			"\x98\x26\x02\xbe\x00"), // show hss skip warning text

	// helper function, takes b = high byte of season addr, returns season in b
	"read default season": MutableString(Addr{0x01, 0x7e89}, "\x01",
		"\x26\x7e\x68\x7e\x47\xc9"),
	// this communicates with the script by setting bit zero of $cfc0 if the
	// warning needs to be displayed (based on room, season, etc), and also
	// displays the exclamation mark if so.
	"warning func": MutableString(Addr{0x15, 0x7947}, "\x15",
		"\xc5\xd5\xcd\x4f\x79\xd1\xc1\xc9"+ // wrap function in push/pops
			"\xfa\x4e\xcc\x47\xfa\xb0\xc6\x4f\xfa\x4c\xcc"+ // load room, season, rod
			"\xfe\x7c\x28\x12\xfe\x6e\x28\x18\xfe\x3d\x28\x22"+ // jump by room
			"\xfe\x5c\x28\x28\xfe\x78\x28\x32\x18\x43"+ // (cont.)
			"\x06\x61\x16\x01\xcd\xd1\x79\xc8\x18\x35"+ // flower cliff
			"\x78\xfe\x03\xc8\x06\x61\x16\x09\xcd\xd1\x79\xc8\x18\x27"+ // diving spot
			"\x06\x65\x16\x02\xcd\xd1\x79\xc8\x18\x1d"+ // waterfall cliff
			"\xfa\x10\xc6\xfe\x0c\xc0\x3e\x17\xcd\x17\x17\xd8\x18\x0f"+ // keep
			"\xcd\x56\x19\xcb\x76\xc0\xcb\xf6\x3e\x02\xea\xe0\xcf\x18\x04"+ // hss skip room
			"\xaf\xea\xe0\xcf"+ // set cliff warning text
			"\xcd\xc6\x3a\xc0\x36\x9f\x2e\x46\x36\x3c"+ // init object
			"\x01\x00\xf1\x11\x0b\xd0\xcd\x1a\x22"+ // set position
			"\x3e\x50\xcd\x74\x0c"+ // play sound
			"\x21\xc0\xcf\xcb\xc6\xc9"), // set $cfc0 bit and ret
	// ORs the default season in the given area (low byte b in bank 1) with the
	// seasons the rod has (c), then ANDs and compares the results with d.
	"warning helper func": MutableString(Addr{0x15, 0x79d1}, "\x15",
		"\x1e\x01\x21\x89\x7e\xcd\x8a\x00"+ // get default season
			"\x78\xb7\x3e\x01\x28\x05\xcb\x27\x05\x20\xfb"+ // match rod format
			"\xb1\xa2\xba\xc9"), // OR with c, AND with d, compare with d, ret

	// all this text overwrites the text from the initial rosa encounter, which
	// runs from 1f:4533 to 1f:45c1 inclusive. the last entry is displayed at
	// the end of any warning message.
	"cliff warning text": MutableString(Addr{0x1f, 0x4533}, "\x0c\x21",
		"\x0c\x00\x02\x3b\x67\x6f\x20\x05\x73\x01"+ // If you go down
			"\x74\x68\x65\x72\x65\x2c\x04\x2d\x20\x77\x6f\x6e\x27\x74\x01"+ // there, you won't
			"\x62\x65\x20\x02\xa4\x05\x0f\x01"+ // be able to get
			"\x04\x9f\x20\x75\x70\x03\xa4"+ // back up.
			"\x07\x03"), // jump to end text
	"bush warning addr": MutableWord(Addr{0x1c, 0x6b50}, 0xfb91, 0xdb91),
	"bush warning text": MutableString(Addr{0x1f, 0x455d}, "\x55\x6d",
		"\x0c\x00\x42\x72\x65\x61\x6b\x03\xa6\x62\x75\x73\x68\x65\x73\x01"+ // Breaking bushes
			"\x03\x69\x04\xb5\x20\x04\xc4\x01"+ // with only those
			"\x04\xcc\x20\x05\xe5\x75\x6e\x73\x61\x66\x65\x03\xa4"+ // items is unsafe.
			"\x07\x03"), // jump to end text
	"hss skip warning addr": MutableWord(Addr{0x1c, 0x6b52}, 0x1192, 0x0292),
	"hss skip warning text": MutableString(Addr{0x1f, 0x4584}, "\x20\x05",
		"\x0c\x00\x02\x3b\x73\x6b\x69\x70\x01"+ // If you skip
			"\x6b\x65\x79\x73\x2c\x04\xaa\x03\x2c\x01"+ // keys, use them
			"\x03\x70\x6c\x79\x03\xa4"+ // carefully.
			"\x07\x03"), // jump to end text
	"end warning addr": MutableWord(Addr{0x1c, 0x6b54}, 0x2592, 0x1d92),
	"end warning text": MutableString(Addr{0x1f, 0x459f}, "\x01\x05",
		"\x0c\x00\x43\x6f\x6e\x74\x69\x6e\x75\x65\x20\x61\x74\x01"+ // Continue at
			"\x03\x0b\x6f\x77\x6e\x20\x72\x69\x73\x6b\x21\x00"), // your own risk!

	// remove the original maku gate interaction
	"maku gate interactions": MutableString(Addr{0x11, 0x634c},
		"\xf2\x22\x0a\x30\x50\x2d\x84\x48\x68\x37\x84\x48\x68",
		"\xf2\x2d\x84\x48\x68\x37\x84\x48\x68\xff"),
	// the interaction on the mount cucco waterfall/vine screen
	"waterfall cliff interaction redirect": MutableString(Addr{0x11, 0x6c10},
		"\xf2\x1f\x08\x68", "\xf3\xb0\x7e\xff"),
	"waterfall cliff interactions": MutableString(Addr{0x11, 0x7eb0}, "\x11",
		"\xf2\x1f\x08\x68\x68\x22\x0a\x20\x18\xfe"),
	// natzu / woods of winter cliff
	"flower cliff interaction redirect": MutableString(Addr{0x11, 0x6568},
		"\xf2\x9c\x00\x58", "\xf3\xba\x7e\xff"),
	"flower cliff interactions": MutableString(Addr{0x11, 0x7eba}, "\x11",
		"\xf2\x9c\x00\x58\x58\x22\x0a\x30\x58\xfe"),
	// sunken city diving spot
	"diving spot interaction redirect": MutableString(Addr{0x11, 0x69cc},
		"\xf2\x1f\x0d\x68", "\xf3\xc4\x7e\xff"),
	"diving spot interactions": MutableString(Addr{0x11, 0x7ec4}, "\x11",
		"\xf2\x1f\x0d\x68\x68\x3e\x31\x18\x68\x22\x0a\x64\x68\xfe"),
	// moblin keep -> sunken city
	"moblin keep interaction redirect": MutableString(Addr{0x11, 0x650b},
		"\xf2\xab\x00\x40", "\xf3\xd2\x7e\xff"),
	"moblin keep interactions": MutableString(Addr{0x11, 0x7ed2}, "\x11",
		"\xf2\xab\x00\x40\x70\x22\x0a\x58\x44\xf8\x2d\x00\x33\xfe"),
	// hss skip room
	"hss skip room interaction redirect": MutableString(Addr{0x11, 0x7ada},
		"\xf3\x93\x55", "\xf3\xe0\x7e"),
	"hss skip room interactions": MutableString(Addr{0x11, 0x7ee0}, "\x11",
		"\xf2\x22\x0a\x88\x98\xf3\x93\x55\xfe"),

	// create a warning interaction when breaking bushes and flowers under
	// certain circumstances.
	"break bush warning call": MutableString(Addr{0x06, 0x477b},
		"\x21\x26\xc6", "\xcd\xd4\x77"),
	"break bush warning func": MutableString(Addr{0x06, 0x77d4}, "\x06",
		"\xf5\xc5\xd5\xcd\xe1\x77\x21\x26\xc6\xd1\xc1\xf1\xc9"+ // wrapper
			"\xfe\xc3\x28\x09\xfe\xc4\x28\x05\xfe\xe5\x28\x01\xc9"+ // tile
			"\xfa\x4c\xcc\xfe\xa7\x28\x0d\xfe\x97\x28\x09"+ // jump by room
			"\xfe\x8d\x28\x05\xfe\x9a\x28\x01\xc9"+ // (cont.)
			"\x3e\x09\xcd\x17\x17\xd8"+ // "already warned" flag
			"\x3e\x16\xcd\x17\x17\xd8"+ // bracelet
			"\x3e\x0e\xcd\x17\x17\xd8"+ // flute
			"\x3e\x05\xcd\x17\x17\xd8"+ // sword
			"\x3e\x06\xcd\x17\x17\xfe\x02\xc8"+ // boomerang, L-2
			"\x21\x92\xc6\x3e\x09\xcd\x0e\x02"+ // set "already warned" flag
			"\x21\xe0\xcf\x36\x01"+ // set warning text index
			"\xcd\xc6\x3a\xc0\x36\x22\x2c\x36\x0a"+ // create warning object
			"\x2e\x4a\x11\x0a\xd0\x06\x04\xcd\x5b\x04"+ // place it on link
			"\xc9"), // ret

	// set hl = address of treasure data + 1 for item with ID a, sub ID c.
	"get treasure data func": MutableString(Addr{0x00, 0x3ed3}, "\x00",
		"\xf5\xc5\xd5\x47\x1e\x15\x21\xf4\x79\xcd\x8a\x00\xd1\xc1\xf1\xc9"),
	"get treasure data body": MutableString(Addr{0x15, 0x79f4}, "\x15",
		"\x78\xc5\x21\x29\x51\xcd\xc3\x01\x09"+ // add ID offset
			"\xcb\x7e\x28\x09\x23\x2a\x66\x6f"+ // load as address if bit 7 set
			"\xc1\x79\xc5\x18\xef"+ // use sub ID as second offset
			"\x23\x06\x03\xd5\x11\xfd\xcd\xcd\x62\x04"+ // copy data
			"\x21\xfd\xcd\xd1\xc1\xc9"), // set hl and ret

	// change hl to point to different treasure data if the item is progressive
	// and needs to be upgraded. param a = treasure ID.
	"progressive item func": MutableString(Addr{0x00, 0x3ee3}, "\x00",
		"\xd5\x5f\xcd\x6a\x3f\x7b\xd1\xd0"+ // ret if you don't have L-1
			"\xfe\x05\x20\x04\x21\x12\x3f\xc9"+ // check sword
			"\xfe\x06\x20\x04\x21\x15\x3f\xc9"+ // check boomerang
			"\xfe\x13\x20\x04\x21\x18\x3f\xc9"+ // check slingshot
			"\xfe\x17\x20\x04\x21\x1b\x3f\xc9"+ // check feather
			"\xfe\x19\xc0\x21\x1e\x3f\xc9"+ // check satchel
			// treasure data
			"\x02\x1d\x11\x02\x23\x1d\x02\x2f\x22\x02\x28\x17\x00\x46\x20"),
	// use cape graphics for stolen feather if applicable.
	"upgrade stolen feather func": MutableString(Addr{0x00, 0x3f6a}, "\x00",
		"\xcd\x17\x17\xd8\xf5\x7b"+ // ret if you have the item
			"\xfe\x17\x20\x13\xd5\x1e\x43\x1a\xfe\x02\xd1\x20\x0a"+ // check IDs
			"\xfa\xb4\xc6\xfe\x02\x20\x03"+ // check feather level
			"\x21\x89\x3f\xf1\xc9"+ // set hl if match
			"\x02\x37\x17"), // treasure data

	// this is a replacement for giveTreasure that gives treasure, plays sound,
	// and sets text based on item ID a and sub ID c, and accounting for item
	// progression.
	"give item func": MutableString(Addr{0x00, 0x3f21}, "\x00",
		"\xcd\xd3\x3e\xcd\xe3\x3e"+ // get treasure data
			"\x4e\xcd\xeb\x16\x28\x05\xe5\xcd\x74\x0c\xe1"+ // give, play sound
			"\x06\x00\x23\x4e\xcd\x4b\x18\xaf\xc9"), // show text

	// upgrade normal items (interactions with ID 60) as necessary when they're
	// created.
	"set normal progressive call": MutableString(Addr{0x15, 0x465a},
		"\x47\xcb\x37", "\xcd\xe8\x79"),
	"set normal progressive func": MutableString(Addr{0x15, 0x79e8}, "\x15",
		"\x47\xcb\x37\xf5\x1e\x42\x1a\xcd\xe3\x3e\xf1\xc9"),

	// utility function, call a function hl in bank 02, preserving af. e can't
	// be used as a parameter to that function, but it can be returned.
	"call bank 02": MutableString(Addr{0x00, 0x3f4d}, "\x00",
		"\xf5\x1e\x02\xcd\x8a\x00\xf1\xc9"),

	// utility function, read a byte from hl in bank e into a and e.
	"read byte from bank": MutableString(Addr{0x00, 0x3f55}, "\x00",
		"\xfa\x97\xff\xf5\x7b\xea\x97\xff\xea\x22\x22"+ // switch bank
			"\x5e\xf1\xea\x97\xff\xea\x22\x22\x7b\xc9"), // read and switch back

	// check fake treasure ID 0a instead of ID of maku tree item. the flag is
	// set in "bank 9 fake id call" below. this only matters if you leave the
	// room without picking up the item.
	"maku tree check fake id": MutableByte(Addr{0x09, 0x7dfd}, 0x42, 0x0a),

	// check fake treasure ID 0f instead of ID of shop item 3.
	"shop check fake id": MutableStrings([]Addr{{0x08, 0x4a8a},
		{0x08, 0x4af2}}, "\x0e", "\x0f"),
	"shop give fake id call": MutableString(Addr{0x08, 0x4bfe},
		"\x1e\x42\x1a", "\xcd\xef\x7f"),
	"shop give fake id func": MutableString(Addr{0x08, 0x7fef}, "\x08",
		"\x1e\x42\x1a\xfe\x0d\xc0\x21\x93\xc6\xcb\xfe\xc9"),

	// check fake treasure ID 10 instead of ID of market item 5. the function
	// is called as part of "market give item func" below.
	"market check fake id": MutableByte(Addr{0x09, 0x7755}, 0x53, 0x10),
	"market give fake id func": MutableString(Addr{0x09, 0x7fdb}, "\x09",
		"\xe5\x21\x94\xc6\xcb\xc6\xe1\x18\xe6"),

	// use fake treasure ID 11 instead of 2e for master diver.
	"diver check fake id": MutableByte(Addr{0x0b, 0x72f1}, 0x2e, 0x11),
	"diver give fake id call": MutableString(Addr{0x0b, 0x730d},
		"\xde\x2e\x00", "\xc0\x94\x7f"),
	"diver give fake id script": MutableString(Addr{0x0b, 0x7f94}, "\x0b",
		"\xde\x2e\x00\x92\x94\xc6\x02\xc1"),

	// not much room left in bank 9, so this calls a bank 2 function that sets
	// treasure ID 12 if applicable.
	"star ore fake id check": MutableByte(Addr{0x08, 0x62fe}, 0x45, 0x12),

	// shared by maku tree and star-shaped ore.
	"bank 9 fake id call": MutableWord(Addr{0x09, 0x42e1}, 0xeb16, 0xe47f),
	"bank 9 fake id func": MutableString(Addr{0x09, 0x7fe4}, "\x09",
		"\xf5\xe5\x21\x21\x76\xcd\x4d\x3f\xe1\xf1\xcd\xeb\x16\xc9"),
	"bank 2 fake id func": MutableString(Addr{0x02, 0x7621}, "\x02",
		"\xfa\x49\xcc\xfe\x01\x28\x05\xfe\x02\x28\x1b\xc9"+ // compare group
			"\xfa\x4c\xcc\xfe\x65\x28\x0d\xfe\x66\x28\x09"+ // compare room
			"\xfe\x75\x28\x05\xfe\x76\x28\x01\xc9"+ // cont.
			"\x21\x94\xc6\xcb\xd6\xc9"+ // set treasure id 12
			"\xfa\x4c\xcc\xfe\x0b\xc0\x21\x93\xc6\xcb\xd6\xc9"), // id 0a

	// use the custom "give item" function in the shop instead of the normal
	// one. this obviates some hard-coded shop data (sprite, text) and allows
	// the item to progressively upgrade.
	"shop give item call": MutableWord(Addr{0x08, 0x4bfc}, 0xeb16, 0xc07f),
	"shop give item func": MutableString(Addr{0x08, 0x7fc0}, "\x08",
		"\xc5\x47\x7d\xcd\xd2\x7f\x78\xc1\x28\x04\xcd\xeb\x16\xc9"+
			"\xcd\x21\x3f\xc9"+ // give item and ret
			"\xfe\xe9\xc8\xfe\xcf\xc8\xfe\xd3\xc8\xfe\xd9\xc9"), // check addr
	// and zero the original text IDs
	"zero shop text": MutableStrings([]Addr{{0x08, 0x4d53}, {0x08, 0x4d46},
		{0x08, 0x4d48}, {0x08, 0x4d4b}}, "\x00", "\x00"),
	// param = b (item index/subID), returns c,e = treasure ID,subID
	"shop item lookup": MutableString(Addr{0x08, 0x7fde}, "\x08",
		"\x21\xce\x4c\x78\x87\xd7\x4e\x23\x5e\xc9"),

	// do the same for the subrosian market.
	"market give item call": MutableString(Addr{0x09, 0x788a},
		"\xfe\x2d\x20\x03\xcd\xb9\x17\xcd\xeb\x16\x1e\x42",
		"\x00\x00\x00\x00\x00\x00\x00\xcd\xae\x7f\x38\x0b"), // jump on carry flag
	"market give item func": MutableString(Addr{0x09, 0x7fae}, "\x09",
		"\xf5\x7d\xfe\xdb\x28\x16\xfe\xe3\x28\x12\xfe\xf5\x28\x1f"+
			"\xf1\xfe\x2d\x20\x03\xcd\xb9\x17\xcd\xeb\x16\x1e\x42\xc9"+
			"\xf1\xcd\x21\x3f\xd1\x37\xc9"), // give item, scf, ret
	// param = b (item index/subID), returns c,e = treasure ID,subID
	"market item lookup": MutableString(Addr{0x09, 0x7fd1}, "\x09",
		"\x21\xda\x77\x78\x87\xd7\x4e\x23\x5e\xc9"),

	// use custom "give item" func in rod cutscene.
	"rod give item call": MutableString(Addr{0x15, 0x70cf},
		"\xcd\xeb\x16", "\xcd\x21\x3f"),
	"no rod text": MutableString(Addr{0x15, 0x70be},
		"\xcd\x4b\x18", "\x00\x00\x00"),
	// returns c,e = treasure ID,subID
	"rod lookup": MutableString(Addr{0x15, 0x7a1a}, "\x15",
		"\x21\xcc\x70\x5e\x23\x23\x4e\xc9"),

	// returns c,e = treasure ID,subID
	"noble sword lookup": MutableString(Addr{0x0b, 0x7f8d}, "\x0b",
		"\x21\x18\x64\x4e\x23\x5e\xc9"),

	// load gfx data for randomized shop and market items.
	"item gfx call": MutableString(Addr{0x3f, 0x443c},
		"\x4f\x06\x00", "\xcd\x69\x71"),
	"item gfx func": MutableString(Addr{0x3f, 0x7169}, "\x3f",
		// check for matching object
		"\x43\x4f\xcd\xdc\x71\x28\x17\x79\xfe\x59\x28\x19"+ // rod, woods
			"\xcd\xbf\x71\x28\x1b\xcd\xcf\x71\x28\x1d"+ // shops
			"\x79\xfe\x6e\x28\x1f\x06\x00\xc9"+ // feather
			// look up item ID, subID
			"\x1e\x15\x21\x1a\x7a\x18\x1d\x1e\x0b\x21\x8d\x7f\x18\x16"+
			"\x1e\x08\x21\xde\x7f\x18\x0f\x1e\x09\x21\xd1\x7f\x18\x08"+
			"\xfa\xb4\xc6\xc6\x15\x5f\x18\x0e"+ // feather
			"\xcd\x8a\x00"+ // get treasure
			"\x79\x4b\xcd\xd3\x3e\xcd\xe3\x3e\x23\x23\x5e"+ // get sprite
			"\x3e\x60\x4f\x06\x00\xc9"), // replace object gfx w/ treasure gfx
	// return z if object is randomized shop item.
	"check randomized shop item": MutableString(Addr{0x3f, 0x71bf}, "\x3f",
		"\x79\xfe\x47\xc0\x7b\xb7\xc8\xfe\x02\xc8\xfe\x05\xc8\xfe\x0d\xc9"),
	// same as above but for subrosia market.
	"check randomized market item": MutableString(Addr{0x3f, 0x71cf}, "\x3f",
		"\x79\xfe\x81\xc0\x7b\xb7\xc8\xfe\x04\xc8\xfe\x0d\xc9"),
	// and rod of seasons.
	"check rod": MutableString(Addr{0x3f, 0x71dc}, "\x3f",
		"\x79\xfe\xe6\xc0\x7b\xfe\x02\xc9"),

	// force the item in the temple of seasons cutscene to use normal item
	// animations.
	"rod cutscene gfx call": MutableString(Addr{0x00, 0x2600},
		"\x1e\x41\x1a", "\xcd\x3b\x3f"),
	"rod cutscene gfx func": MutableString(Addr{0x00, 0x3f3b}, "\x00",
		"\x1e\x41\x1a\xfe\xe6\xc0\x1c\x1a\xfe\x02\x28\x03\x1d\x1a\xc9"+
			"\x3e\x60\xc9"),

	// move the bushes on the rosa portal screen by one tile so that it's
	// possible to leave and re-enter without breaking bushes.
	"move rosa portal bushes": MutableStrings([]Addr{
		{0x21, 0x7454}, {0x22, 0x709d}, {0x23, 0x6ea9}, {0x24, 0x6b9f}},
		"\x0e\xc4\xf7\x4d\x5f\x11\x6e\x38\xc4\x11\x5e\xf7\x5d\x11",
		"\x38\xc4\xf7\x4d\x04\x5d\x6e\x38\xc4\x11\x5e\xf7\x4d\x5f"),

	// prevent the first member's shop item from always refilling all seeds.
	"no shop seed refill": MutableString(Addr{0x08, 0x4c02},
		"\xcc\xe5\x17", "\x00\x00\x00"),
	// instead, have any satchel refill all seeds.
	"satchel seed refill call": MutableString(Addr{0x00, 0x16f6},
		"\xcd\xc8\x44", "\xcd\xe4\x71"),
	"satchel seed refill func": MutableString(Addr{0x3f, 0x71e4}, "\x3f",
		"\xc5\xcd\xc8\x44\x78\xc1\xf5\x78\xfe\x19\x20\x07"+
			"\xc5\xd5\xcd\xe5\x17\xd1\xc1\xf1\x47\xc9"),
}

var (
	mapIconByTreeID  = []byte{0x15, 0x19, 0x16, 0x17, 0x18, 0x18}
	roomByTreeID     = []byte{0xf8, 0x9e, 0x67, 0x72, 0x5f, 0x10}
	roomNameByTreeID = []string{
		"ember tree room", "mystery tree room", "scent tree room",
		"pegasus tree room", "sunken gale tree room", "tarm gale tree room",
	}
)

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

	// seed tree rooms (need to match seed types for regrowth)
	"ember tree room":       MutableByte(Addr{0x01, 0x5eed}, 0xf8, 0xf8),
	"mystery tree room":     MutableByte(Addr{0x01, 0x5eef}, 0x9e, 0x9e),
	"scent tree room":       MutableByte(Addr{0x01, 0x5ef1}, 0x67, 0x67),
	"pegasus tree room":     MutableByte(Addr{0x01, 0x5ef3}, 0x72, 0x72),
	"sunken gale tree room": MutableByte(Addr{0x01, 0x5ef5}, 0x5f, 0x5f),
	"tarm gale tree room":   MutableByte(Addr{0x01, 0x5ef7}, 0x10, 0x10),

	// the satchel should contain the type of seeds that grow on the horon
	// village tree.
	"satchel initial seeds": MutableByte(Addr{0x3f, 0x453b}, 0x20, 0x20),

	// give the player seeds when they get the slingshot, and don't take the
	// player's: fool's ore when they get feather, star ore when they get
	// ribbon, or red and blue ore when they get hard ore (just zero the whole
	// "lose items" table). one byte of this is changed in setSeedData() to
	// change what type of seeds the slingshot gives.
	"edit gain/lose items tables": MutableString(Addr{0x3f, 0x4543},
		"\x00\x46\x45\x00\x52\x50\x51",
		"\x13\x20\x20\x00\x00\x00\x00"),
	"edit lose items table pointer": MutableByte(Addr{0x3f, 0x44cf},
		0x44, 0x47),

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

	// set sub ID for star ore
	"star ore id call": MutableString(Addr{0x08, 0x62f2},
		"\x2c\x36\x45", "\xcd\xe8\x7f"),
	"star ore id func": MutableString(Addr{0x08, 0x7fe8}, "\x08",
		"\x2c\x36\x45\x2c\x36\x00\xc9"),

	// set sub ID for hard ore
	"hard ore id call": MutableString(Addr{0x15, 0x5b83},
		"\x2c\x36\x52", "\xcd\x22\x7a"),
	"hard ore id func": MutableString(Addr{0x15, 0x7a22}, "\x15",
		"\x2c\x36\x52\x2c\x36\x00\xc9"),
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
