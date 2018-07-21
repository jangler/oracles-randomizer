// Package rom deals with the structure of the OOS ROM file itself. The given
// addresses are for the Japanese version of the game.
package rom

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
)

const bankSize = 0x4000

// bankOffset returns the offset of the given bank in the ROM.
func bankOffset(bank int) int64 {
	if bank < 2 {
		return 0
	}
	return int64(bankSize * (bank - 1))
}

// Load reads ROM data from a reader into memory.
func Load(f io.Reader) (*bytes.Buffer, error) {
	// read file into buffer
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, f)
	return buf, err
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
func Verify(buf *bytes.Buffer) []error {
	errors := make([]error, 0)
	reader := bytes.NewReader(buf.Bytes())

	// check mutables TODO
	for k, m := range Mutables {
		if err := verifyMutable(reader, m, k); err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func verifyMutable(r io.ReaderAt, m Mutable, name string) error {
	mData := m.Bytes()
	romData := make([]byte, len(mData))
	if _, err := r.ReadAt(romData, m.RealAddr()); err != nil {
		return err
	}
	if bytes.Compare(romData, mData) != 0 {
		return fmt.Errorf("%s: at %x, expected %x, got %x",
			name, m.RealAddr(), mData, romData)
	}
	return nil
}
