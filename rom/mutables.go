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
		codeMutables["tree warp"].(*MutableRange).New[19] = 0x18
	} else {
		codeMutables["tree warp"].(*MutableRange).New[19] = 0x28
	}
}

// SetNoMusic sets music off in the modified ROM.
func SetNoMusic() {
	mut := codeMutables["no music func"].(*MutableRange)
	funcAddr := addrString(mut.Addrs[0].Offset)
	codeMutables["no music call"].(*MutableRange).New =
		[]byte("\xcd" + funcAddr)
}

// SetAnimal sets the flute type and Natzu region type based on a companion
// number 1 to 3.
func SetAnimal(companion int) {
	varMutables["animal region"].(*MutableRange).New =
		[]byte{byte(companion + 0x0a)}
}

// these mutables have fixed addresses and don't reference other mutables. try
// to generally order them by address, unless a grouping between mutables in
// different banks makes more sense.
var fixedMutables = map[string]Mutable{
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
	// AND allow transition away from the screen if you have feather (not once
	// the hole is dug)
	"leave H&S screen": MutableString(Addr{0x09, 0x65a0},
		"\xcd\x32\x14\x1e\x49\x1a\xbe\xc8",
		"\x3e\x17\xcd\x17\x17\x00\x00\xd0"),

	// move the trigger for the bridge from holodrum plain to natzu to the
	// top-left corner of the screen, where it can't be hit, and replace the
	// lever tile as well. this prevents the bridge from blocking the waterway.
	"remove bridge trigger": MutableWord(Addr{0x11, 0x6737},
		0x6868, 0x0000),
	"remove prairie lever":   MutableByte(Addr{0x21, 0x6267}, 0xb1, 0x04),
	"remove wasteland lever": MutableByte(Addr{0x23, 0x5cb7}, 0xb1, 0x04),

	// skip shield check for forging hard ore
	"skip iron shield check": MutableByte(Addr{0x0b, 0x75c7}, 0x01, 0x02),
	// and skip the check for what level shield you currently have
	"skip iron shield level check": MutableString(Addr{0x15, 0x62ac},
		"\x38\x01", "\x18\x05"),

	// check fake treasure ID 0a for maku tree item. this only matters if you
	// leave the screen without picking up the item.
	"maku tree check fake id": MutableByte(Addr{0x09, 0x7dfd}, 0x42, 0x0a),
	// check fake treasure ID 0f for shop item 3.
	"shop check fake id": MutableStrings([]Addr{{0x08, 0x4a8a},
		{0x08, 0x4af2}}, "\x0e", "\x0f"),
	// check fake treasure ID 10 for market item 5.
	"market check fake id": MutableByte(Addr{0x09, 0x7755}, 0x53, 0x10),
	// check fake treasure ID 11 for master diver.
	"diver check fake id": MutableByte(Addr{0x0b, 0x72f1}, 0x2e, 0x11),
	// check fake treasure ID 12 for subrosia seaside,
	"star ore fake id check": MutableByte(Addr{0x08, 0x62fe}, 0x45, 0x12),

	// bank 00

	// blaino normally sets bit 6 of active ring to "unequip" it instead of
	// setting it to $ff. this only matters for the dev ring.
	"fix blaino ring unequip": MutableWord(Addr{0x00, 0x2376}, 0xcbf6, 0x36ff),

	// bank 01

	// the d5 boss key room is hard-coded to make a compass beep, even though
	// the room's can beep based on dungeon room properties.
	"fix d5 boss key beep": MutableByte(Addr{0x01, 0x4a0a}, 0x0c, 0x00),

	// bank 04

	// a hack so that a different flag can be used to set the rosa portal tile
	// replacement, allowing the bush-breaking warning interaction to be used
	// on this screen.
	"portal tile replacement": MutableString(Addr{0x04, 0x6016},
		"\x40\x33\xc5", "\x20\x33\xe6"),

	// banks 08-0a (most interaction-specific non-script behavior?)

	// have horon village shop stock *and* sell items from the start, including
	// the flute. also don't stop the flute from appearing because of animal
	// flags, since it probably won't be a flute at all.
	"horon shop stock check":   MutableByte(Addr{0x08, 0x4adb}, 0x05, 0x02),
	"horon shop sell check":    MutableByte(Addr{0x08, 0x48d0}, 0x05, 0x02),
	"horon shop flute check 1": MutableByte(Addr{0x08, 0x4b02}, 0xcb, 0xf6),
	"horon shop flute check 2": MutableWord(Addr{0x08, 0x4afb},
		0xcb6f, 0xafaf),

	// prevent the first member's shop item from always refilling all seeds.
	"no shop seed refill": MutableString(Addr{0x08, 0x4c02},
		"\xcc\xe5\x17", "\x00\x00\x00"),

	// zero the original shop item text (don't remember if this is actually
	// necessary).
	"zero shop text": MutableStrings([]Addr{{0x08, 0x4d53}, {0x08, 0x4d46},
		{0x08, 0x4d48}, {0x08, 0x4d4b}}, "\x00", "\x00"),

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

	// restrict the area triggering sokra to talk to link in horon village to
	// the left side of the burnable trees (prevents softlock).
	"resize sokra trigger": MutableString(Addr{0x08, 0x5ba5},
		"\xfa\x0b\xd0\xfe\x3c\xd8\xfe\x60\xd0",
		"\xfe\x88\xd0\xfa\x0b\xd0\xfe\x3c\xd8"),

	// i don't know what global flag 0e is. it's only checked in for star ore
	// digging, and disabling the check seems to be sometimes necessary (?)
	"star ore flag check": MutableString(Addr{0x08, 0x62aa},
		"\xc2\xd9\x3a", "\x00\x00\x00"),
	// a vanilla bug lets star ore be dug up on the first screen even if you
	// already have the item. so… make first try a second instance of second
	// try.
	"star ore bugfix": MutableWord(Addr{0x08, 0x62d5}, 0x6656, 0x7624),

	// prevent leaving sunken city with dimitri unless you have his flute, in
	// order to prevent a variety of softlocks.
	"block dimitri exit": MutableString(Addr{0x09, 0x6f34},
		"\xfa\x10\xc6\xfe\x0c", "\xfa\xaf\xc6\xfe\x02"),

	// normally none of the desert pits will work if the player already has the
	// rusty bell.
	"desert item check": MutableByte(Addr{0x08, 0x739e}, 0x4a, 0x04),

	// moosh won't spawn in the mountains if you have the wrong number of
	// essences. bit 6 seems related to this, and needs to be zero too?
	"skip moosh essence check 1": MutableByte(Addr{0x0f, 0x7429}, 0x03, 0x00),
	"skip moosh essence check 2": MutableByte(Addr{0x09, 0x4e36}, 0xca, 0xc3),
	"skip moosh flag check":      MutableByte(Addr{0x09, 0x4ead}, 0x40, 0x00),

	// sell member's card in subrosian market before completing d3
	"member's card essence check": MutableWord(Addr{0x09, 0x7750},
		0xcb57, 0xf601),

	// count number of essences, not highest numbered essence.
	"maku seed check 1": MutableByte(Addr{0x09, 0x7da4}, 0xea, 0x76),
	"maku seed check 2": MutableByte(Addr{0x09, 0x7da6}, 0x30, 0x18),

	// stop the hero's cave event from giving you a second wooden sword that
	// you use to spin slash
	"wooden sword second item": MutableByte(Addr{0x0a, 0x7bb9}, 0x05, 0x3f),

	// bank 0b (scripts)

	// don't set a ricky flag when buying the "flute".
	"shop no set ricky flag": MutableByte(Addr{0x0b, 0x4826}, 0x20, 0x00),

	// don't require rod to get items from season spirits.
	"season spirit rod check": MutableByte(Addr{0x0b, 0x4eb2}, 0x07, 0x02),

	// getting the L-2 (or L-3) sword in the lost woods normally gives a second
	// "spin slash" item. remove this from the script.
	"noble sword second item":  MutableByte(Addr{0x0b, 0x641a}, 0xde, 0xc1),
	"master sword second item": MutableByte(Addr{0x0b, 0x6421}, 0xde, 0xc1),

	// end maku seed script as soon as link gets the seed.
	"abbreviate maku seed cutscene": MutableString(Addr{0x0b, 0x71ec},
		"\xe1\x23\x61\x01", "\xb6\x19\xbe\x00"),
	// end northen peak barrier cutscene as soon as the barrier is broken.
	"abbreviate barrier cutscene": MutableString(Addr{0x0b, 0x79f1},
		"\x88\x18\x50\xf8", "\xb6\x1d\xbe\x00"),

	// bank 0d

	// grow seeds in all seasons
	"seeds grow always": MutableByte(Addr{0x0d, 0x68b5}, 0xb8, 0xbf),

	// bank 11 (interactions)

	// remove the moosh and dimitri events in spool swamp.
	"prevent moosh cutscene":   MutableByte(Addr{0x11, 0x6572}, 0xf1, 0xff),
	"prevent dimitri cutscene": MutableByte(Addr{0x11, 0x68d4}, 0xf1, 0xff),

	// bank 14

	// change the noble sword's animation pointers to match regular items
	"noble sword anim 1": MutableWord(Addr{0x14, 0x53d7}, 0x5959, 0x1957),
	"noble sword anim 2": MutableWord(Addr{0x14, 0x55a7}, 0xf36b, 0x4f68),

	// bank 15 (script functions)

	// you can softlock in d6 misusing keys without magnet gloves, so just move
	// the magnet ball onto the button it needs to press to get the key the
	// speedrun skips.
	"move d6 magnet ball": MutableByte(Addr{0x15, 0x4f36}, 0x98, 0x58),

	// if you go up the stairs into the room in d8 with the magnet ball and
	// can't move it, you don't have room to go back down the stairs. this
	// moves the magnet ball's starting position one more tile away.
	"move d8 magnet ball": MutableByte(Addr{0x15, 0x4f62}, 0x48, 0x38),

	// change destination of initial transition in pirate cutscene.
	"pirate warp": MutableString(Addr{0x15, 0x5a1c},
		"\x81\x74\x00\x42", "\x80\xe2\x00\x66"),

	// zero normal rod text.
	"no rod text": MutableString(Addr{0x15, 0x70be},
		"\xcd\x4b\x18", "\x00\x00\x00"),

	// banks 1c-1f (text)

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

	// banks 21-24 (room layouts)

	// replace the rock/flower outside of d6 with a normal bush so that the
	// player doesn't get softlocked if they exit d6 without gale satchel or
	// default spring.
	"replace d6 flower spring": MutableByte(Addr{0x21, 0x4e73}, 0xd8, 0xc4),
	"replace d6 flower non-spring": MutableStrings(
		[]Addr{{0x22, 0x4b83}, {0x23, 0x4973}, {0x24, 0x45d0}},
		"\x92", "\xc4"),

	// change water tiles outside d4 from deep to shallow (prevents softlock
	// from entering without flippers or default summer).
	"change d4 water tiles": MutableStrings(
		[]Addr{{0x21, 0x54a9}, {0x22, 0x5197}, {0x23, 0x4f6c}},
		"\xfd\x6b\x6b\x53\xfa\x3f\xfd", "\xfa\x6b\x6b\x53\xfa\x3f\xfa"),
	"change d4 water tiles winter": MutableString(Addr{0x24, 0x4cec},
		"\xfd\x00\xfc\x06\xfd\xfd\xfd\xfd",
		"\xdc\x00\xfc\x06\xdc\xdc\xdc\xdc"),

	// move the bushes on the rosa portal screen by one tile so that it's
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

	// remove the rock blocking exit from D5, since it makes no difference in
	// logic and is a softlock unless otherwise prevented.
	"remove rock outside d5": MutableStrings([]Addr{
		{0x21, 0x7448}, {0x22, 0x7091}, {0x23, 0x6e9d}, {0x24, 0x6b93}},
		"\xc0", "\x12"),
	// but add an extra (non-interactable) rock on the actual D5 screen so that
	// ricky can't jump up the cliff.
	"add rock outside d5": MutableStrings([]Addr{
		{0x21, 0x7031}, {0x22, 0x6c6e}, {0x23, 0x6a7c}, {0x24, 0x677d}},
		"\x12", "\x64"),

	// make it possible to leave and re-enter rosa's portal without breaking
	// bushes.
	"move rosa portal bushes": MutableStrings([]Addr{
		{0x21, 0x7454}, {0x22, 0x709d}, {0x23, 0x6ea9}, {0x24, 0x6b9f}},
		"\x0e\xc4\xf7\x4d\x5f\x11\x6e\x38\xc4\x11\x5e\xf7\x5d\x11",
		"\x38\xc4\xf7\x4d\x04\x5d\x6e\x38\xc4\x11\x5e\xf7\x4d\x5f"),

	// replace some currents in spool swamp in spring so that the player isn't
	// trapped by them.
	"replace currents 1": MutableWord(Addr{0x21, 0x7ab1}, 0xd2d2, 0xd3d3),
	"replace currents 2": MutableString(Addr{0x21, 0x7ab6},
		"\xd3\xd2\xd2", "\xd4\xd4\xd4"),
	"replace currents 3": MutableByte(Addr{0x21, 0x7abe}, 0xd3, 0xd1),

	// replace the stairs outside the portal in eyeglass lake in summer with a
	// railing, because if the player jumps off those stairs in summer they
	// fall into the noble sword room.
	"replace lake stairs": MutableString(Addr{0x22, 0x791b},
		"\x36\xd0\x35", "\x40\x40\x40"),

	// remove the snow piles in front of holly's house so that shovel isn't
	// required not to softlock there.
	"remove holly snow piles": MutableByte(Addr{0x24, 0x6474}, 0xd9, 0x04),
	// remove some snow piles outside D7 for the same reason.
	"remove d7 snow piles": MutableString(Addr{0x24, 0x7910},
		"\xd9\xa0\xb9\xd9", "\x2b\xa0\xb9\x2b"),

	// bank 3f

	// since slingshot doesn't increment seed capacity, set the level-zero
	// capacity of seeds to 20, and move the pointer up by one byte.
	"satchel capacity": MutableString(Addr{0x3f, 0x4617},
		"\x20\x50\x99", "\x20\x20\x50"),
	"satchel capacity pointer": MutableByte(Addr{0x3f, 0x460e}, 0x16, 0x17),

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
}

// like the item slots, these are (usually) no-ops until the randomizer touches
// them. these are also fixed, but generally need to have their values set
// elsewhere in order to do anything.
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

	// locations of sparkles on treasure map
	"round jewel coords":    MutableByte(Addr{0x02, 0x6663}, 0xb5, 0xb5),
	"pyramid jewel coords":  MutableByte(Addr{0x02, 0x6664}, 0x1d, 0x1d),
	"square jewel coords":   MutableByte(Addr{0x02, 0x6665}, 0xc2, 0xc2),
	"x-shaped jewel coords": MutableByte(Addr{0x02, 0x6666}, 0xf4, 0xf4),

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
		fixedMutables,
		treasureMutables,
		slotMutables,
		varMutables,
		seasonMutables,
		codeMutables,
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
