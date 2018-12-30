// Package rom deals with the structure of the oracles ROM files themselves.
// The given addresses are for the English versions of the games, and if two
// are specified, Ages comes first.
package rom

import (
	"crypto/sha1"
	"fmt"
	"math/rand"
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

	// use these graphics as default for progressive items (seasons)
	itemGfx["sword 1"] = itemGfx["sword L-1"]
	itemGfx["sword 2"] = itemGfx["sword L-1"]
	itemGfx["boomerang 1"] = itemGfx["boomerang L-1"]
	itemGfx["boomerang 2"] = itemGfx["boomerang L-1"]
	itemGfx["slingshot 1"] = itemGfx["slingshot L-1"]
	itemGfx["slingshot 2"] = itemGfx["slingshot L-1"]
	itemGfx["feather 1"] = itemGfx["feather L-1"]
	itemGfx["feather 2"] = itemGfx["feather L-1"]

	// (ages)
	itemGfx["sword 1"] = itemGfx["sword L-1"]
	itemGfx["switch hook 1"] = itemGfx["switch hook"]
	itemGfx["switch hook 2"] = itemGfx["long hook"]
	itemGfx["bracelet 1"] = itemGfx["bracelet"]
	itemGfx["bracelet 2"] = itemGfx["power glove"]
	itemGfx["harp 1"] = itemGfx["tune of echoes"]
	itemGfx["harp 2"] = itemGfx["tune of currents"]
	itemGfx["harp 3"] = itemGfx["tune of ages"]
	itemGfx["flippers 1"] = itemGfx["flippers"]
	itemGfx["flippers 2"] = itemGfx["mermaid suit"]
}

// Addr is a fully-specified memory address.
type Addr struct {
	bank   uint8
	offset uint16
}

// fullOffset returns the actual offset of the address in the ROM, based on
// bank number and relative address.
func (a *Addr) fullOffset() int {
	var bankOffset int
	if a.bank >= 2 {
		bankOffset = bankSize * (int(a.bank) - 1)
	}
	return bankOffset + int(a.offset)
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
		ItemSlots["subrosia seaside"].idAddrs[0].offset = codeAddr.offset + 2
		ItemSlots["subrosia seaside"].subIDAddrs[0].offset = codeAddr.offset + 5
		codeAddr = codeMutables["hard ore id func"].(*MutableRange).Addrs[0]
		ItemSlots["great furnace"].idAddrs[0].offset = codeAddr.offset + 2
		ItemSlots["great furnace"].subIDAddrs[0].offset = codeAddr.offset + 5
		codeAddr = codeMutables["diver fake id script"].(*MutableRange).Addrs[0]
		ItemSlots["master diver's reward"].idAddrs[0].offset = codeAddr.offset + 1
		ItemSlots["master diver's reward"].subIDAddrs[0].offset = codeAddr.offset + 2
	} else {
		// explicitly set these addresses and IDs after their functions
		mut := codeMutables["soldier script give item"].(*MutableRange)
		slot := ItemSlots["deku forest soldier"]
		slot.idAddrs[0].offset = mut.Addrs[0].offset + 13
		slot.subIDAddrs[0].offset = mut.Addrs[0].offset + 14
		codeAddr := codeMutables["target carts flag"].(*MutableRange).Addrs[0]
		ItemSlots["target carts 2"].idAddrs[1].offset = codeAddr.offset + 1
		ItemSlots["target carts 2"].subIDAddrs[1].offset = codeAddr.offset + 2
	}

	setSeedData(game)

	// set the text IDs for all rings to $ff (blank), since custom code deals
	// with text
	for _, treasure := range Treasures {
		if treasure.id == 0x2d {
			treasure.text = 0xff
		}
	}

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
		ItemSlots["subrosia seaside"].Mutate(b)
		ItemSlots["great furnace"].Mutate(b)
		ItemSlots["master diver's reward"].Mutate(b)
	} else {
		ItemSlots["nayru's house"].Mutate(b)
		ItemSlots["deku forest soldier"].Mutate(b)
		ItemSlots["target carts 2"].Mutate(b)
		ItemSlots["hidden tokay cave"].Mutate(b)
	}

	setCompassData(b, game)

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
		case "lost woods", "d6 boomerang chest", "d8 HSS chest",
			"d7 cape chest", "member's shop 1", "sword 2", "boomerang 2",
			"slingshot 2", "feather 2", "satchel 2":
		// shop items (use sub ID instead of param, no text)
		case "shop, 20 rupees", "shop, 30 rupees", "shop, 150 rupees",
			"member's shop 2", "member's shop 3", "subrosia market, 1st item",
			"subrosia market, 2nd item", "subrosia market, 5th item",
			"zero shop text":
		// seasons misc.
		case "maku tree", "fool's ore", "member's card", "treasure map",
			"temple of seasons", "rare peach stone", "ribbon", "blaino prize",
			"subrosia seaside", "great furnace", "subrosian smithy",
			"master diver's reward", "d5 basement":
		// ages misc.
		case "sword 1", "nayru's house", "south shore dirt", "target carts 1",
			"target carts 2", "big bang game", "harp 1", "harp 2", "harp 3",
			"sea of storms present", "sea of storms past", "starting chest",
			"deku forest soldier", "hidden tokay cave", "ridge bush cave",
			"graveyard poe":
		// ages, script item using collect mode other than 0a
		case "trade lava juice", "goron dance, with letter", "goron elder",
			"balloon guy's upgrade", "king zora", "d2 thwomp shelf":
		// ages, progressive items/slots not covered elsewhere
		case "d6 present vire chest", "d7 miniboss chest", "d8 floor puzzle",
			"tokkey's composition", "rescue nayru", "bracelet 2", "flippers 2":
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
		seedType = ItemSlots["horon village seed tree"].Treasure.id
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

		for _, name := range []string{
			"horon village seed tree map icon",
			"north horon seed tree map icon",
			"woods of winter seed tree map icon",
			"spool swamp seed tree map icon",
			"sunken city seed tree map icon",
			"tarm ruins seed tree map icon",
		} {
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

		// satchel and shooter come with south lynna tree seeds
		mut := varMutables["satchel initial seeds"].(*MutableRange)
		mut.New[0] = 0x20 + seedType
		mut = codeMutables["fill seed shooter"].(*MutableRange)
		mut.New[6] = 0x20 + seedType
		for _, name := range []string{"satchel initial selection",
			"shooter initial selection"} {
			mut := varMutables[name].(*MutableRange)
			mut.New[1] = seedType
		}

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
func setCompassData(b []byte, game int) {
	var names []string
	if game == GameSeasons {
		names = []string{"d1 goriya chest", "d2 terrace chest",
			"d3 giant blade room", "d4 dive spot", "d5 basement",
			"d6 escape room", "d7 stalfos chest", "d8 pols voice chest"}
	} else {
		names = []string{"d1 pot chest", "d2 color room", "d3 B1F east",
			"d4 lava pot chest", "d5 owl puzzle", "d6 present RNG chest",
			"d7 post-hallway chest", "d8 B3F chest"}
	}

	// clear original boss key flags
	for _, name := range names {
		slot := ItemSlots[name]
		offset :=
			getDungeonPropertiesAddr(game, slot.group, slot.room).fullOffset()
		b[offset] = b[offset] & 0xef // reset bit 4
	}

	// add new boss key flags
	for i := 1; i <= 8; i++ {
		name := fmt.Sprintf("d%d boss key", i)
		slot := lookupItemSlot(name)
		offset :=
			getDungeonPropertiesAddr(game, slot.group, slot.room).fullOffset()
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
func getDungeonPropertiesAddr(game int, group, room byte) *Addr {
	offset := uint16(room)
	if game == GameSeasons {
		offset += 0x4d41
	} else {
		offset += 0x4dce
	}
	if group%2 != 0 {
		offset += 0x100
	}
	return &Addr{0x01, offset}
}

// RandomizeRingPool randomizes the types of rings in the item pool, returning
// a map of vanilla ring names to the randomized ones.
func RandomizeRingPool(src *rand.Rand) map[string]string {
	nameMap := make(map[string]string)
	usedRings := make([]bool, 0x40)

	for _, slot := range ItemSlots {
		if slot.Treasure.id == 0x2d {
			oldName := rings[slot.Treasure.param]

			// loop until we get a ring that's not literally useless, and which
			// we haven't used before.
			done := false
			for !done {
				param := byte(src.Intn(0x40))
				switch rings[param] {
				case "friendship ring", "GBA time ring", "GBA nature ring",
					"slayer's ring", "rupee ring", "victory ring", "sign ring",
					"100th ring":
					break
				default:
					if !usedRings[param] {
						slot.Treasure.param = param
						usedRings[param] = true
						done = true
					}
				}
			}

			nameMap[oldName] = rings[slot.Treasure.param]
		}
	}

	return nameMap
}
