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
	setRodGraphics()

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
		if k == "maku key fall" || k == "rod gift" ||
			strings.HasSuffix(k, "ring") {
			continue // special cases that will error but we don't care about
		}
		if err := m.Check(b); err != nil {
			errors = append(errors, fmt.Errorf("%s: %v", k, err))
		}
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// sets the rod's cutscene graphics based on the treasure assigned to its slot
func setRodGraphics() {
	slot := ItemSlots["rod gift"]
	treasure := slot.Treasure
	itemName := treasureNameFromIDs(treasure.id, treasure.subID)
	if gfx := rodGraphics[itemName]; gfx == 0 {
		log.Fatalf("fatal: no rod graphics for %s (%02x%02x)",
			itemName, treasure.id, treasure.subID)
	} else {
		rodGfxMutable.New = []byte{byte(gfx >> 16), byte(gfx >> 8), byte(gfx)}
	}
}
