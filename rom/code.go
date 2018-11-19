package rom

import (
	"fmt"
	"strings"
)

// this file is for mutables that go at the end of banks. each should be a
// self-contained unit (i.e. don't jr to anywhere outside the byte string) so
// that they can be appended automatically with respect to their size.

// return e.g. "\x2d\x79" for 0x792d
func addrString(addr uint16) string {
	return string([]byte{byte(addr), byte(addr >> 8)})
}

// adds code at the given address, returning the length of the byte string.
func addCode(name string, bank byte, offset uint16, code string) uint16 {
	codeMutables[name] = MutableString(Addr{bank, offset},
		string([]byte{bank}), code)
	return uint16(len(code))
}

type romBanks struct {
	endOfBank []uint16
}

var codeMutables = map[string]Mutable{}

// appendToBank appends the given data to the end of the given bank, associates
// it with the given name, and returns the address of the data as a string such
// as "\xc8\x3e" for 0x3ec8. it panics if the end of the bank is zero or if the
// data would overflow the bank.
func (r *romBanks) appendToBank(bank byte, name, data string) string {
	eob := r.endOfBank[bank]

	if eob == 0 {
		panic(fmt.Sprintf("end of bank %02x undefined for %s", bank, name))
	}

	if eob+uint16(len(data)) > 0x8000 {
		panic(fmt.Sprintf("not enough space for %s in bank %02x", name, bank))
	}

	codeMutables[name] = MutableString(Addr{bank, eob}, "", data)
	r.endOfBank[bank] += uint16(len(data))

	return addrString(eob)
}

// replace replaces the old data at the given address with the new data, and
// associates the change with the given name. actual replacement will fail at
// runtime if the old data does not match the original data in the ROM.
func (r *romBanks) replace(bank byte, offset uint16, name, old, new string) {
	codeMutables[name] = MutableString(Addr{bank, offset}, old, new)
}

// replaceMultiple acts as replace, but operates on multiple addresses.
func (r *romBanks) replaceMultiple(addrs []Addr, name, old, new string) {
	codeMutables[name] = MutableStrings(addrs, old, new)
}

// returns a byte table of (group, room, collect mode) entries for randomized
// items. in ages, a mode >7f means to use &7f as an index to a jump table for
// special cases.
func makeCollectModeTable() string {
	b := new(strings.Builder)

	for _, slot := range ItemSlots {
		// trees and slots where it doesn't matter (shops, rod)
		if slot.collectMode == 0 {
			continue
		}

		_, err := b.Write([]byte{slot.group, slot.room, slot.collectMode})
		if err != nil {
			panic(err)
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}
