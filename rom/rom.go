// Package rom deals with the structure of the oracles ROM files themselves.
// The given addresses are for the English versions of the games, and if two
// are specified, Ages comes first.
package rom

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

const bankSize = 0x4000

const (
	GameNil = iota
	GameAges
	GameSeasons
)

var itemGfx map[string]int

func Init(game int) {
	if game == GameAges {
		ItemSlots = agesSlots
		Treasures = agesTreasures
		fixedMutables = agesFixedMutables
		varMutables = agesVarMutables
		itemGfx = agesItemGfx
		initAgesEOB()
	} else {
		ItemSlots = seasonsSlots
		Treasures = seasonsTreasures
		fixedMutables = seasonsFixedMutables
		varMutables = seasonsVarMutables
		itemGfx = seasonsItemGfx
		initSeasonsEOB()

		for k, v := range Seasons {
			varMutables[k] = v
		}
	}

	for _, slot := range ItemSlots {
		slot.Treasure = Treasures[slot.treasureName]
	}

	// rings and boss keys all have the same sprite
	for name, treasure := range Treasures {
		if treasure.id == 0x2d {
			itemGfx[name] = itemGfx["ring"]
		}
		if treasure.id == 0x31 {
			itemGfx[name] = itemGfx["boss key"]
		}
	}

	// use these graphics as default for progressive items
	itemGfx["sword 1"] = itemGfx["sword L-1"]
	itemGfx["sword 2"] = itemGfx["sword L-1"]
	itemGfx["boomerang 1"] = itemGfx["boomerang L-1"]
	itemGfx["boomerang 2"] = itemGfx["boomerang L-1"]
	itemGfx["slingshot 1"] = itemGfx["slingshot L-1"]
	itemGfx["slingshot 2"] = itemGfx["slingshot L-1"]
	itemGfx["feather 1"] = itemGfx["feather L-1"]
	itemGfx["feather 2"] = itemGfx["feather L-1"]

	// get set of unique items (to determine which can be slotted freely)
	treasureCounts := make(map[string]int)
	for _, slot := range ItemSlots {
		name := FindTreasureName(slot.Treasure)
		if treasureCounts[name] == 0 {
			treasureCounts[name] = 1
		} else {
			treasureCounts[name]++
		}
	}
	for name, count := range treasureCounts {
		if count == 1 {
			TreasureIsUnique[name] = true
		}
	}
	for _, name := range []string{"ricky's flute", "dimitri's flute",
		"moosh's flute"} {
		TreasureIsUnique[name] = true
	}
	for _, name := range []string{"d1 boss key", "d2 boss key", "d3 boss key",
		"d6 boss key", "d7 boss key", "d8 boss key"} {
		delete(TreasureIsUnique, name)
	}
}

// Addr is a fully-specified memory address.
type Addr struct {
	Bank   uint8
	Offset uint16
}

// FullOffset returns the actual offset of the address in the ROM, based on
// bank number and relative address.
func (a *Addr) FullOffset() int {
	var bankOffset int
	if a.Bank >= 2 {
		bankOffset = bankSize * (int(a.Bank) - 1)
	}
	return bankOffset + int(a.Offset)
}

func IsAges(b []byte) bool {
	return string(b[0x134:0x13f]) == "ZELDA NAYRU"
}

func IsSeasons(b []byte) bool {
	return string(b[0x134:0x13d]) == "ZELDA DIN"
}

func IsUS(b []byte) bool {
	return b[0x014a] != 0
}

func IsVanilla(b []byte) bool {
	knownSum := "\x88\x03\x74\xfb\x97\x8b\x18\xaf\x4a\xa5\x29\xe2\xe3\x2f\x7f" +
		"\xfb\x4d\x7d\xd2\xf4"
	if IsSeasons(b) {
		knownSum = "\xba\x12\x68\x29\x0f\xb2\xb1\xb7\x05\x05\xd2\xd7\xb5\x82" +
			"\x5f\xc8\xa4\x81\x6a\x4b"
	}
	sum := sha1.Sum(b)

	return string(sum[:]) == knownSum
}

// get mutables in order, so that sums are consistent with the same seed
func orderedKeys(m map[string]Mutable) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Mutate changes the contents of loaded ROM bytes in place. It returns a
// checksum of the result or an error.
func Mutate(b []byte, game int) ([]byte, error) {
	if game == GameSeasons {
		varMutables["initial season"].(*MutableRange).New =
			[]byte{0x2d, Seasons["north horon season"].New[0]}
		codeMutables["season after pirate cutscene"].(*MutableRange).New =
			[]byte{Seasons["western coast season"].New[0]}

		setTreasureMapData()

		// explicitly set these addresses and IDs after their functions
		codeAddr := codeMutables["star ore id func"].(*MutableRange).Addrs[0]
		ItemSlots["star ore spot"].IDAddrs[0].Offset = codeAddr.Offset + 2
		ItemSlots["star ore spot"].SubIDAddrs[0].Offset = codeAddr.Offset + 5
		codeAddr = codeMutables["hard ore id func"].(*MutableRange).Addrs[0]
		ItemSlots["hard ore slot"].IDAddrs[0].Offset = codeAddr.Offset + 2
		ItemSlots["hard ore slot"].SubIDAddrs[0].Offset = codeAddr.Offset + 5
		codeAddr = codeMutables["diver fake id script"].(*MutableRange).Addrs[0]
		ItemSlots["diver gift"].IDAddrs[0].Offset = codeAddr.Offset + 1
		ItemSlots["diver gift"].SubIDAddrs[0].Offset = codeAddr.Offset + 2
	} else {
		setAgesGfx("cheval's test", 0x6b)
		setAgesGfx("cheval's invention", 0x6b)
		setAgesGfx("tokay hut", 0x6b)
		setAgesGfx("wild tokay game", 0x63)
		setAgesGfx("shop, 150 rupees", 0x47)
		setAgesGfx("library present", 0x80)
		setAgesGfx("library past", 0x80)

		// explicitly set these addresses and IDs after their functions
		codeAddr := codeMutables["target carts flag"].(*MutableRange).Addrs[0]
		ItemSlots["target carts 2"].IDAddrs[1].Offset = codeAddr.Offset + 1
		ItemSlots["target carts 2"].SubIDAddrs[1].Offset = codeAddr.Offset + 2
	}

	setSeedData(game)

	var err error
	mutables := getAllMutables()
	for _, k := range orderedKeys(mutables) {
		err = mutables[k].Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	// explicitly set these IDs after their functions are written
	if game == GameSeasons {
		ItemSlots["star ore spot"].Mutate(b)
		ItemSlots["hard ore slot"].Mutate(b)
		ItemSlots["diver gift"].Mutate(b)

		setCompassData(b)
	} else {
		ItemSlots["nayru's house"].Mutate(b)
		ItemSlots["target carts 2"].Mutate(b)
	}

	outSum := sha1.Sum(b)
	return outSum[:], nil
}

// Verify checks all the package's data against the ROM to see if it matches.
// It returns a slice of errors describing each mismatch.
func Verify(b []byte, game int) []error {
	errors := make([]error, 0)
	for k, m := range getAllMutables() {
		// ignore special cases that would error even when correct
		switch k {
		// flutes
		case "ricky's flute", "moosh's flute", "dimitri's flute",
			"strange flute":
		// mystical seeds
		case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
			"pegasus tree seeds", "gale tree seeds":
		// progressive items
		case "noble sword spot", "d6 boomerang chest", "d8 HSS chest",
			"d7 cape chest", "member's shop 1", "sword 2", "boomerang 2",
			"slingshot 2", "feather 2", "satchel 2":
		// shop items (use sub ID instead of param, no text)
		case "village shop 1", "village shop 2", "village shop 3",
			"member's shop 2", "member's shop 3", "subrosian market 1",
			"subrosian market 2", "subrosian market 5", "zero shop text":
		// seasons misc.
		case "maku tree gift", "fool's ore", "member's card", "treasure map",
			"rod gift", "rare peach stone", "ribbon", "blaino gift",
			"star ore spot", "hard ore slot", "iron shield gift", "diver gift",
			"d5 boss key spot":
		// ages misc.
		case "sword 1", "nayru's house", "maku tree", "south shore dirt",
			"target carts 1", "target carts 2", "big bang game", "tokay hut":
		// ages, script item using collect mode other than 0a
		case "trade lava juice", "goron dancing past", "goron elder",
			"tingle's upgrade", "king zora":
			break
		default:
			if err := m.Check(b); err != nil {
				errors = append(errors, fmt.Errorf("%s: %v", k, err))
			}
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// set the initial satchel and slingshot seeds (and selections) based on what
// grows on the horon village tree, and set the map icon for each tree to match
// the seed type.
func setSeedData(game int) {
	var seedType byte
	if game == GameSeasons {
		seedType = ItemSlots["ember tree"].Treasure.id
	} else {
		seedType = ItemSlots["south lynna tree"].Treasure.id
	}

	if game == GameSeasons {
		for _, name := range []string{"satchel initial seeds",
			"carry seeds in slingshot"} {
			mut := varMutables[name].(*MutableRange)
			mut.New[0] = 0x20 + seedType
		}

		// slingshot starting seeds
		varMutables["edit gain/lose items tables"].(*MutableRange).New[1] =
			0x20 + seedType

		for _, name := range []string{
			"satchel initial selection", "slingshot initial selection"} {
			mut := varMutables[name].(*MutableRange)
			mut.New[1] = seedType
		}

		for _, name := range []string{"ember tree map icon",
			"scent tree map icon", "mystery tree map icon",
			"pegasus tree map icon", "sunken gale tree map icon",
			"tarm gale tree map icon"} {
			mut := varMutables[name].(*MutableRange)
			slotName := strings.Replace(name, " map icon", "", 1)
			id := ItemSlots[slotName].Treasure.id
			mut.New[0] = 0x15 + id
		}
	} else {
		// set high nybbles (seed types) of seed tree interactions
		setTreeNybble(varMutables["symmetry city tree sub ID"],
			ItemSlots["symmetry city tree"])
		setTreeNybble(varMutables["south lynna present tree sub ID"],
			ItemSlots["south lynna tree"])
		setTreeNybble(varMutables["crescent island tree sub ID"],
			ItemSlots["crescent island tree"])
		setTreeNybble(varMutables["zora village present tree sub ID"],
			ItemSlots["zora village tree"])
		setTreeNybble(varMutables["rolling ridge west tree sub ID"],
			ItemSlots["rolling ridge west tree"])
		setTreeNybble(varMutables["ambi's palace tree sub ID"],
			ItemSlots["ambi's palace tree"])
		setTreeNybble(varMutables["rolling ridge east tree sub ID"],
			ItemSlots["rolling ridge east tree"])
		setTreeNybble(varMutables["south lynna past tree sub ID"],
			ItemSlots["south lynna tree"])
		setTreeNybble(varMutables["deku forest tree sub ID"],
			ItemSlots["deku forest tree"])
		setTreeNybble(varMutables["zora village past tree sub ID"],
			ItemSlots["zora village tree"])

		// set map icons
		for _, name := range []string{"crescent island tree",
			"symmetry city tree", "south lynna tree", "zora village tree",
			"rolling ridge west tree", "ambi's palace tree",
			"rolling ridge east tree", "deku forest tree"} {
			mut := varMutables[name+" map icon"].(*MutableRange)
			mut.New[0] = 0x15 + ItemSlots[name].Treasure.id
		}
	}
}

// sets the high nybble (seed type) of a seed tree interaction in ages.
func setTreeNybble(subID Mutable, slot *MutableSlot) {
	mut := subID.(*MutableRange)
	mut.New[0] = (mut.Old[0] & 0x0f) | (slot.Treasure.id << 4)
}

// set the locations of the sparkles for the jewels on the treasure map.
func setTreasureMapData() {
	for _, name := range []string{"round", "pyramid", "square", "x-shaped"} {
		mut := varMutables[name+" jewel coords"].(*MutableRange)
		slot := lookupItemSlot(name + " jewel")
		mut.New[0] = slot.mapCoords
	}
}

// match the compass's beep beep beep boops to the actual boss key locations.
func setCompassData(b []byte) {
	// clear original boss key flags
	for _, name := range []string{"d1 boss key chest", "d2 boss key chest",
		"d3 boss key chest", "d4 boss key spot", "d5 boss key spot",
		"d6 boss key chest", "d7 boss key chest", "d8 boss key chest"} {
		slot := ItemSlots[name]
		offset := getDungeonPropertiesAddr(slot.group, slot.room).FullOffset()
		b[offset] = b[offset] & 0xef // reset bit 4
	}

	// add new boss key flags
	for i := 1; i <= 8; i++ {
		name := fmt.Sprintf("d%d boss key", i)
		slot := lookupItemSlot(name)
		offset := getDungeonPropertiesAddr(slot.group, slot.room).FullOffset()
		b[offset] = (b[offset] & 0xbf) | 0x10 // set bit 4, reset bit 6
	}
}

// returns the slot where the named item was placed. this only works for unique
// items, of course.
func lookupItemSlot(itemName string) *MutableSlot {
	t := Treasures[itemName]
	for _, slot := range ItemSlots {
		if slot.Treasure == t {
			return slot
		}
	}
	return nil
}

// get the location of the dungeon properties byte for a specific room.
func getDungeonPropertiesAddr(group, room byte) *Addr {
	offset := 0x4d41 + uint16(room)
	if group%2 != 0 {
		offset += 0x100
	}
	return &Addr{0x01, offset}
}

// some item-related interactions need explicit graphics changes, including
// interaction 6b (not needed in seasons).
func setAgesGfx(name string, interactionID byte) {
	mut := varMutables[name+" gfx"].(*MutableRange)
	treasureName := FindTreasureName(ItemSlots[name].Treasure)
	gfx := itemGfx[treasureName]
	if gfx == 0 {
		panic("no item graphics for " + treasureName)
	}
	mut.New = []byte{byte(gfx >> 16), byte(gfx >> 8), byte(gfx)}

	switch interactionID {
	case 0x6b:
		switch gfx & 0x0f {
		case 0x00, 0x02:
			mut.New[2]++
		case 0x03:
			mut.New[2]--
		default:
			panic(treasureName + " item gfx are incompatible with " + name)
		}
	case 0x80:
		switch gfx & 0x0f {
		case 0x00, 0x03:
			mut.New[2] = byte(gfx&0xf0) | 0x04
		case 0x02:
			break
		default:
			panic(treasureName + " item gfx are incompatible with " + name)
		}
	}
}
