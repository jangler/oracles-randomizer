// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"crypto/sha1"
	"fmt"
	"sort"
	"strings"
)

const bankSize = 0x4000

func init() {
	// rings and boss keys all have the same sprite
	for name, treasure := range Treasures {
		if treasure.id == 0x2d {
			itemGfx[name] = itemGfx["ring"]
		}
		if treasure.id == 0x31 {
			itemGfx[name] = itemGfx["boss key"]
		}
	}

	// override blaino's default collection mode
	ItemSlots["blaino gift"].CollectMode = CollectFind2

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
		if (idCounts[t.id] == 1 && name != "gasha seed" &&
			name != "piece of heart") || strings.HasSuffix(name, "flute") {
			uniqueIDTreasures[name] = true
		}
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

func IsSeasons(b []byte) bool {
	return string(b[0x134:0x13d]) == "ZELDA DIN"
}

func IsUS(b []byte) bool {
	return b[0x014a] != 0
}

func IsVanilla(b []byte) bool {
	knownSum := "\xba\x12\x68\x29\x0f\xb2\xb1\xb7\x05\x05\xd2\xd7\xb5\x82\x5f" +
		"\xc8\xa4\x81\x6a\x4b"
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
func Mutate(b []byte) ([]byte, error) {
	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, Seasons["north horon season"].New[0]}
	varMutables["season after pirate cutscene"].(*MutableRange).New =
		[]byte{Seasons["western coast season"].New[0]}

	setSeedData()

	var err error
	mutables := getAllMutables()
	for _, k := range orderedKeys(mutables) {
		err = mutables[k].Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	// then fix rod graphics
	b[ItemSlots["rod gift"].GfxAddrs[0].FullOffset()+2] += 1

	outSum := sha1.Sum(b)
	return outSum[:], nil
}

// Update changes the content of loaded ROM bytes, but does not re-randomize
// any fields.
func Update(b []byte) ([]byte, error) {
	var err error

	// change fixed mutables
	for _, k := range orderedKeys(constMutables) {
		err = constMutables[k].Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, b[Seasons["north horon season"].Addrs[0].FullOffset()]}
	varMutables["season after pirate cutscene"].(*MutableRange).New =
		[]byte{b[Seasons["western coast season"].Addrs[0].FullOffset()]}

	// change seed mechanics based on the ROM's existing tree information
	for _, name := range []string{"ember tree", "scent tree", "mystery tree",
		"pegasus tree", "sunken gale tree", "tarm gale tree"} {
		ItemSlots[name].Treasure.id =
			b[ItemSlots[name].IDAddrs[0].FullOffset()]
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
	return outSum[:], nil
}

// Verify checks all the package's data against the ROM to see if it matches.
// It returns a slice of errors describing each mismatch.
func Verify(b []byte) []error {
	errors := make([]error, 0)

	for k, m := range getAllMutables() {
		switch k {
		// special cases that will error normally.
		// (flippers' collect mode is different between regions)
		case "maku tree gift", "fool's ore", "noble sword spot", "flippers",
			"ember tree seeds", "mystery tree seeds", "scent tree seeds",
			"pegasus tree seeds", "gale tree seeds 1", "gale tree seeds 2",
			"expert's ring", "energy ring", "toss ring", "fist ring",
			"member's card", "treasure map", "member's shop 3",
			"subrosian market 5", "member's shop 1", "ricky's flute",
			"moosh's flute", "dimitri's flute", "strange flute", "rod gift",
			"rare peach stone", "ribbon", "blaino gift":
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

	for i, name := range []string{"ember tree", "mystery tree", "scent tree",
		"pegasus tree", "sunken gale tree", "tarm gale tree"} {
		slot := ItemSlots[name]
		mut := varMutables[roomNameByTreeID[slot.Treasure.id]].(*MutableRange)
		mut.New[0] = roomByTreeID[i]
	}
}
