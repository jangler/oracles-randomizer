// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"crypto/sha1"
	"fmt"
	"log"
	"strings"
)

const bankSize = 0x4000

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

// Mutate changes the contents of loaded ROM bytes in place. It returns a
// checksum of the result or an error.
func Mutate(b []byte) ([]byte, error) {
	setSceneGfx("rod gift", "rod graphics")
	setSceneGfx("noble sword spot", "noble sword graphics")
	setSceneGfx("noble sword spot", "master sword graphics")
	setSceneGfx("d0 sword chest", "wooden sword graphics")
	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, Seasons["north horon season"].New[0]}

	setSeedData()

	log.Printf("old bytes: sha-1 %x", sha1.Sum(b))
	var err error
	for _, m := range getAllMutables() {
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
	for _, m := range constMutables {
		err = m.Mutate(b)
		if err != nil {
			return nil, err
		}
	}

	varMutables["initial season"].(*MutableRange).New =
		[]byte{0x2d, b[Seasons["north horon season"].Addr.FullOffset()]}

	// change seed mechanics based on the ROM's existing tree information
	for _, name := range []string{"ember tree", "scent tree", "mystery tree",
		"pegasus tree", "sunken gale tree", "tarm gale tree"} {
		ItemSlots[name].Treasure.id = b[ItemSlots[name].IDAddrs[0].FullOffset()]
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

	for k, m := range getAllMutables() {
		switch k {
		// special cases that will error normally
		case "maku key fall", "fool's ore", "noble sword spot",
			"ember tree seeds", "mystery tree seeds", "scent tree seeds",
			"pegasus tree seeds", "gale tree seeds 1", "gale tree seeds 2":
			break
		default:
			if strings.HasSuffix(k, " ring") {
				break
			}
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
	if gfx := sceneItemGfx[itemName]; gfx == 0 {
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
