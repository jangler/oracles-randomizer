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

// Mutate changes the contents of loaded ROM bytes in place.
func Mutate(b []byte) error {
	setSceneGfx("rod gift", "rod graphics")
	setSceneGfx("noble sword spot", "noble sword graphics")
	setSceneGfx("noble sword spot", "master sword graphics")
	setSceneGfx("d0 sword chest", "wooden sword graphics")

	log.Printf("old bytes: sha-1 %x", sha1.Sum(b))
	var err error
	for _, m := range getAllMutables() {
		err = m.Mutate(b)
		if err != nil {
			return err
		}
	}
	log.Printf("new bytes: sha-1 %x", sha1.Sum(b))
	return nil
}

// Verify checks all the package's data against the ROM to see if it matches.
// It returns a slice of errors describing each mismatch.
func Verify(b []byte) []error {
	errors := make([]error, 0)

	for k, m := range getAllMutables() {
		switch k {
		// special cases that will error normally
		case "d0 sword chest", "maku key fall", "rod gift", "fool's ore",
			"noble sword spot", "ember tree seeds", "mystery tree seeds",
			"scent tree seeds", "pegasus tree seeds", "gale tree seeds 1",
			"gale tree seeds 2":
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
		mut := dataMutables[gfxName].(MutableRange)
		mut.New = []byte{byte(gfx >> 16), byte(gfx >> 8), byte(gfx)}
		dataMutables[gfxName] = mut
	}
}
