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
		ItemSlots = AgesSlots
		Treasures = AgesTreasures
		fixedMutables = agesFixedMutables
		varMutables = agesVarMutables
		itemGfx = agesItemGfx
		initAgesEOB()
	} else {
		ItemSlots = SeasonsSlots
		Treasures = SeasonsTreasures
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
		if treasure.id == 0x30 {
			itemGfx[name] = itemGfx["small key"]
		}
	}
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
		codeAddr = codeMutables["create mt. cucco item"].(*MutableRange).Addrs[0]
		ItemSlots["mt. cucco, platform cave"].idAddrs[0].offset = codeAddr.offset + 2
		ItemSlots["mt. cucco, platform cave"].subIDAddrs[0].offset = codeAddr.offset + 1
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

	setBossItemAddrs()
	setSeedData(game)
	setSmallKeyData(game)
	setCollectModeData(game)

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

	// explicitly set these items after their functions are written
	writeBossItems(b)
	if game == GameSeasons {
		ItemSlots["subrosia seaside"].Mutate(b)
		ItemSlots["great furnace"].Mutate(b)
		ItemSlots["master diver's reward"].Mutate(b)

		// annoying special case to prevent text on key drop
		mut := ItemSlots["d7 armos puzzle"]
		if mut.Treasure.id == SeasonsTreasures["d7 small key"].id {
			b[mut.subIDAddrs[0].fullOffset()] = 0x01
		}
	} else {
		ItemSlots["nayru's house"].Mutate(b)
		ItemSlots["deku forest soldier"].Mutate(b)
		ItemSlots["target carts 2"].Mutate(b)
		ItemSlots["hidden tokay cave"].Mutate(b)

		// other special case to prevent text on key drop
		mut := ItemSlots["d8 stalfos"]
		if mut.Treasure.id == AgesTreasures["d8 small key"].id {
			b[mut.subIDAddrs[0].fullOffset()] = 0x00
		}
	}

	setCompassData(b, game)
	setLinkedData(b, game)

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
		// mystical seeds
		case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
			"pegasus tree seeds", "gale tree seeds":
		// seasons shop items
		case "strange flute", "zero shop text", "member's card", "treasure map",
			"rare peach stone", "ribbon":
		// seasons misc.
		case "temple of seasons", "blaino prize", "green joy ring",
			"mt. cucco, platform cave", "diving spot outside D4":
		// ages progressive w/ different item IDs
		case "nayru's house", "tokkey's composition", "rescue nayru",
			"d6 present vire chest":
		// ages misc.
		case "south shore dirt", "target carts 1", "target carts 2",
			"sea of storms past", "starting chest", "graveyard poe":
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

// fill table—initial table is blank, since it's created before items are
// placed.
func setSmallKeyData(game int) {
	mut := codeMutables["small key drops"].(*MutableRange)
	mut.New = []byte(makeKeyDropTable())

	if game == GameSeasons {
		mut := varMutables["above d7 zol button"].(*MutableSlot)
		mut.Treasure = ItemSlots["d7 zol button"].Treasure
	}
}

// regenerate collect mode table to accommodate changes based on contents.
func setCollectModeData(game int) {
	mut := codeMutables["collectModeTable"].(*MutableRange)
	if game == GameSeasons {
		mut.New = []byte(makeSeasonsCollectModeTable())
	} else {
		mut.New = []byte(makeAgesCollectModeTable())
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

// set dungeon properties so that the compass beeps in the rooms actually
// containing small keys and boss keys.
func setCompassData(b []byte, game int) {
	var prefixes []string
	if game == GameSeasons {
		prefixes = []string{"d0", "d1", "d2", "d3", "d4", "d5", "d6", "d7",
			"d8"}
	} else {
		prefixes = []string{"d0", "d1", "d2", "d3", "d4", "d5", "d6 present",
			"d6 past", "d7", "d8"}
	}

	// clear key flags
	for _, prefix := range prefixes {
		for name, slot := range ItemSlots {
			if strings.HasPrefix(name, prefix+" ") {
				offset := getDungeonPropertiesAddr(
					game, slot.group, slot.room).fullOffset()
				b[offset] = b[offset] & 0xed // reset bit 4
			}
		}
	}

	// set key flags
	for _, prefix := range prefixes {
		slots := lookupAllItemSlots(fmt.Sprintf("%s small key", prefix))
		switch prefix {
		case "d0", "d6 present":
			break
		case "d6 past":
			slots = append(slots, lookupItemSlot("d6 boss key"))
		default:
			slots = append(slots,
				lookupItemSlot(fmt.Sprintf("%s boss key", prefix)))
		}

		for _, slot := range slots {
			offset := getDungeonPropertiesAddr(
				game, slot.group, slot.room).fullOffset()
			b[offset] = (b[offset] & 0xbf) | 0x10 // set bit 4, reset bit 6
		}
	}
}

// returns the slot where the named item was placed. this only works for unique
// items, of course.
func lookupItemSlot(itemName string) *MutableSlot {
	if slots := lookupAllItemSlots(itemName); len(slots) > 0 {
		return slots[0]
	} else {
		return nil
	}
}

// returns all slots where the named item was placed.
func lookupAllItemSlots(itemName string) []*MutableSlot {
	t := Treasures[itemName]
	slots := make([]*MutableSlot, 0)
	for _, slot := range ItemSlots {
		if slot.Treasure == t {
			slots = append(slots, slot)
		}
	}
	return slots
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
func RandomizeRingPool(src *rand.Rand, game int) map[string]string {
	nameMap := make(map[string]string)
	usedRings := make([]bool, 0x40)

	keys := make([]string, len(ItemSlots))
	i := 0
	for key, _ := range ItemSlots {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	for _, key := range keys {
		slot := ItemSlots[key]

		if slot.Treasure.id == 0x2d {
			oldName := FindTreasureName(slot.Treasure)

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
				case "rang ring L-1", "rang ring L-2":
					// these rings are literally useless in ages.
					if game == GameAges {
						break
					}
					fallthrough
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

func setBossItemAddrs() {
	table := codeMutables["bossItemTable"].(*MutableRange)

	for i := uint16(1); i <= 8; i++ {
		slot := ItemSlots[fmt.Sprintf("d%d boss", i)]
		slot.idAddrs[0].offset = table.Addrs[0].offset + i*2
		slot.subIDAddrs[0].offset = table.Addrs[0].offset + i*2 + 1
	}
}

func writeBossItems(b []byte) {
	for i := 1; i <= 8; i++ {
		ItemSlots[fmt.Sprintf("d%d boss", i)].Mutate(b)
	}
}

// set data to make linked playthroughs isomorphic to unlinked ones.
func setLinkedData(b []byte, game int) {
	if game == GameSeasons {
		// set linked starting / hero's cave terrace items based on which items
		// in unlinked hero's cave aren't keys. order matters.
		var tStart, tCave *Treasure
		if ItemSlots["d0 key chest"].Treasure.id == 0x30 {
			tStart = ItemSlots["d0 sword chest"].Treasure
			tCave = ItemSlots["d0 rupee chest"].Treasure
		} else {
			tStart = ItemSlots["d0 key chest"].Treasure
			tCave = ItemSlots["d0 sword chest"].Treasure
		}

		// give this item at start
		linkedStartItem := &MutableSlot{
			idAddrs:    []Addr{{0x0a, 0x7ffd}},
			subIDAddrs: []Addr{{0x0a, 0x7ffe}},
			Treasure:   tStart,
		}
		linkedStartItem.Mutate(b)

		// create slot for linked hero's cave terrace
		linkedChest := seasonsChest(
			"rupees, 20", 0x50e2, 0x05, 0x2c, collectChest, 0xd4)
		linkedChest.Treasure = tCave
		linkedChest.Mutate(b)
	}
}
