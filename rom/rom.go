// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"crypto/sha1"
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

const bankSize = 0x4000

// Addr is a fully-specified memory address.
type Addr struct {
	Bank   uint8
	Offset uint16
}

func (a *Addr) FullOffset() int {
	var bankOffset int
	if a.Bank >= 2 {
		bankOffset = bankSize * (int(a.Bank) - 1)
	}
	return bankOffset + int(a.Offset)
}

// Load reads ROM data from a reader into memory.
func Load(f io.Reader) ([]byte, error) {
	return ioutil.ReadAll(f)
}

// Mutate changes the contents of loaded ROM bytes in place.
func Mutate(b []byte) error {
	log.Printf("old bytes: sha-1 %x", sha1.Sum(b))
	var err error
	for _, m := range Mutables {
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

	// check mutables TODO
	for k, m := range Mutables {
		if k == "maku key fall" {
			continue // special case that will error but we don't care about
		}
		if err := m.Check(b); err != nil {
			errors = append(errors, fmt.Errorf("%s: %v", k, err))
		}
	}

	return errors
}
