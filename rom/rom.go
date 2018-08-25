// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"crypto/sha1"
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	bankSize   = 0x4000
	regionAddr = 0x014a // 0 = JP, 1 = US
)

func init() {
	// rings all have the same sprite
	for name, treasure := range Treasures {
		if treasure.id == 0x2d {
			narrowItemGfx[name] = narrowItemGfx["ring"]
		}
	}

	// accumulate all item sprites into map
	for name, sprite := range narrowItemGfx {
		itemGfx[name] = sprite
	}
	for name, sprite := range wideItemGfx {
		itemGfx[name] = sprite
	}

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

	// get set of items with unique IDs (more restrictive than the above)
	idCounts := make(map[byte]int)
	for _, t := range Treasures {
		if idCounts[t.id] == 0 {
			idCounts[t.id] = 1
		} else {
			idCounts[t.id]++
		}
	}
	for name, t := range Treasures {
		// TODO do other items need special expection?
		if idCounts[t.id] == 1 &&
			name != "gasha seed" && name != "piece of heart" {
			uniqueIDTreasures[name] = true
		}
	}
}

func isEn(b []byte) bool {
	return b[regionAddr] != 0
}

// Addr is a fully-specified memory address.
type Addr struct {
	Bank     uint8
	JpOffset uint16
	EnOffset uint16
}

// FullOffset returns the actual offset of the address in the ROM, based on
// bank number and relative address.
func (a *Addr) FullOffset(en bool) int {
	var bankOffset int
	if a.Bank >= 2 {
		bankOffset = bankSize * (int(a.Bank) - 1)
	}
	if en {
		return bankOffset + int(a.EnOffset)
	}
	return bankOffset + int(a.JpOffset)
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
func Mutate(b []byte) ([]byte, error) {
	setSceneGfx("rod gift", "rod graphics")
	setSceneGfx("noble sword spot", "noble sword graphics")
	setSceneGfx("noble sword spot", "master sword graphics")
	setSceneGfx("d0 sword chest", "wooden sword graphics")
	setSceneGfx("member's shop 1", "member's shop 1 graphics")
	setSceneGfx("member's shop 2", "member's shop 2 graphics")
	setSceneGfx("member's shop 3", "member's shop 3 graphics")
	setSceneGfx("subrosian market 2", "subrosian market 2 graphics")
	setSceneGfx("subrosian market 5", "subrosian market 5 graphics")
	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, Seasons["north horon season"].New[0]}

	setSeedData()

	en := isEn(b)
	log.Printf("old bytes: sha-1 %x", sha1.Sum(b))
	var err error
	mutables := getAllMutables()
	for _, k := range orderedKeys(mutables) {
		m := mutables[k]
		if (strings.HasSuffix(k, "(en)") && !en) ||
			(strings.HasSuffix(k, "(jp)") && en) {
			continue
		}

		err = m.Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	outSum := sha1.Sum(b)
	log.Printf("new bytes: sha-1 %x", outSum)
	return outSum[:], nil
}

// Update changes the content of loaded ROM bytes, but does not re-randomize
// any fields.
func Update(b []byte) ([]byte, error) {
	var err error
	log.Printf("old bytes: sha-1 %x", sha1.Sum(b))

	// change fixed mutables
	en := isEn(b)
	for _, k := range orderedKeys(constMutables) {
		m := constMutables[k]
		if (strings.HasSuffix(k, "(en)") && !en) ||
			(strings.HasSuffix(k, "(jp)") && en) {
			continue
		}

		err = m.Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, b[Seasons["north horon season"].Addr.FullOffset(en)]}

	// change seed mechanics based on the ROM's existing tree information
	for _, name := range []string{"ember tree", "scent tree", "mystery tree",
		"pegasus tree", "sunken gale tree", "tarm gale tree"} {
		ItemSlots[name].Treasure.id =
			b[ItemSlots[name].IDAddrs[0].FullOffset(en)]
	}
	setSeedData()
	for _, name := range []string{"satchel initial seeds",
		"slingshot initial seeds", "satchel initial selection",
		"slingshot initial selection", "carry seeds in slingshot",
		"ember tree map icon", "scent tree map icon", "mystery tree map icon",
		"pegasus tree map icon", "sunken gale tree map icon",
		"tarm gale tree map icon", "initial season"} {
		err = varMutables[name].Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	outSum := sha1.Sum(b)
	log.Printf("new bytes: sha-1 %x", outSum)
	return outSum[:], nil
}

// Verify checks all the package's data against the ROM to see if it matches.
// It returns a slice of errors describing each mismatch.
func Verify(b []byte) []error {
	errors := make([]error, 0)

	en := isEn(b)
	for k, m := range getAllMutables() {
		if (strings.HasSuffix(k, "(en)") && !en) ||
			(strings.HasSuffix(k, "(jp)") && en) {
			continue
		}

		switch k {
		// special cases that will error normally.
		// (flippers' collect mode is different between regions)
		case "maku tree gift", "fool's ore", "noble sword spot", "flippers",
			"ember tree seeds", "mystery tree seeds", "scent tree seeds",
			"pegasus tree seeds", "gale tree seeds 1", "gale tree seeds 2",
			"expert's ring", "energy ring", "toss ring", "fist ring",
			"member's card", "treasure map", "member's shop 3",
			"subrosian market 5", "member's shop 1":
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

// sets a mutable's cutscene graphics from the treasure assigned to its slot
func setSceneGfx(slotName, gfxName string) {
	slot := ItemSlots[slotName]
	treasure := slot.Treasure
	itemName := treasureNameFromIDs(treasure.id, treasure.subID)
	if gfx := itemGfx[itemName]; gfx == 0 {
		log.Fatalf("fatal: no %s for %s (%02x%02x)",
			gfxName, itemName, treasure.id, treasure.subID)
	} else {
		mut := varMutables[gfxName].(*MutableRange)
		mut.New = []byte{byte(gfx >> 16), byte(gfx >> 8), byte(gfx)}
	}
}

// set the initial satchel and slingshot seeds (and selections) based on what
// grows on the horon village tree, and set the map icon for each tree to match
// the seed type.
func setSeedData() {
	seedIndex := seedIndexByTreeID[int(ItemSlots["ember tree"].Treasure.id)]

	for _, name := range []string{"satchel initial seeds",
		"slingshot initial seeds", "carry seeds in slingshot"} {
		mut := varMutables[name].(*MutableRange)
		mut.New[0] = 0x20 + seedIndex
	}

	for _, name := range []string{
		"satchel initial selection", "slingshot initial selection"} {
		mut := varMutables[name].(*MutableRange)
		mut.New[1] = seedIndex
	}

	for _, name := range []string{"ember tree map icon", "scent tree map icon",
		"mystery tree map icon", "pegasus tree map icon",
		"sunken gale tree map icon", "tarm gale tree map icon"} {
		mut := varMutables[name].(*MutableRange)
		id := ItemSlots[strings.Replace(name, " map icon", "", 1)].Treasure.id
		mut.New[0] = mapIconByTreeID[int(id)]
	}
}
