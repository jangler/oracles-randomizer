package rom

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

// this file is for mutables that go at the end of banks. each should be a
// self-contained unit (i.e. don't jr to anywhere outside the byte string) so
// that they can be appended automatically with respect to their size.

// return e.g. "\x2d\x79" for 0x792d
func addrString(addr uint16) string {
	return string([]byte{byte(addr), byte(addr >> 8)})
}

// return e.g. 0x792d for "\x2d\x79"
func stringAddr(addr string) uint16 {
	return uint16([]byte(addr)[0]) + uint16([]byte(addr)[1])<<8
}

// adds code at the given address, returning the length of the byte string.
func addCode(name string, bank byte, offset uint16, code string) uint16 {
	codeMutables[name] = MutableString(Addr{bank, offset},
		string([]byte{bank}), code)
	return uint16(len(code))
}

type romBanks struct {
	endOfBank []uint16
	assembler *assembler
	addrs     map[string]uint16
}

// used for unmarshaling asm data from yaml.
type asmData struct {
	Addrs map[string]uint16
	Banks map[byte][]map[string]string
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

// perform substitutions on labels in asm.
func (r *romBanks) subLabels(s string) string {
	// perform substitutions from other entries
	// TODO cache these regexps?
	for k, v := range r.addrs {
		re := regexp.MustCompile(`\b` + k + `\b`)
		s = re.ReplaceAllString(s, fmt.Sprintf("%04x", v))
	}

	return s
}

// appendAsm acts as appendToBank, but by compiling a block of asm. additional
// arguments are formatted into `asm` by fmt.Sprintf. the returned address is
// also given as a uint16 rather than a big-endian word in string form.
func (r *romBanks) appendAsm(bank byte, name, asm string,
	a ...interface{}) uint16 {
	asm = r.subLabels(asm)

	var err error
	asm, err = r.assembler.compile(fmt.Sprintf(asm, a...), ";\n")
	if err != nil {
		panic(err)
	}

	as := r.appendToBank(bank, name, asm)
	r.addrs[name] = stringAddr(as)
	return r.addrs[name]
}

// replace replaces the old data at the given address with the new data, and
// associates the change with the given name. actual replacement will fail at
// runtime if the old data does not match the original data in the ROM.
func (r *romBanks) replace(bank byte, offset uint16, name, old, new string) {
	codeMutables[name] = MutableString(Addr{bank, offset}, old, new)
}

// replaceAsm acts as replace, but treating the old and new strings as asm
// instead of machine code. additional arguments are formatted into `new` by
// fmt.Sprintf.
func (r *romBanks) replaceAsm(bank byte, offset uint16, old, new string,
	a ...interface{}) {
	old = r.subLabels(old)
	new = r.subLabels(new)

	var err error
	old, err = r.assembler.compile(old, ";\n")
	if err != nil {
		panic(err)
	}
	new, err = r.assembler.compile(fmt.Sprintf(new, a...), ";\n")
	if err != nil {
		panic(err)
	}

	name := fmt.Sprintf("replacement at %02x:%04x", bank, offset)

	r.replace(bank, offset, name, old, new)
}

// replaceMultiple acts as replace, but operates on multiple addresses.
func (r *romBanks) replaceMultiple(addrs []Addr, name, old, new string) {
	codeMutables[name] = MutableStrings(addrs, old, new)
}

// returns an ordered slice of keys for slot names, so that dentical seeds
// produce identical checksums.
func getOrderedSlotKeys() []string {
	keys := make([]string, 0, len(ItemSlots))
	for k := range ItemSlots {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// returns a byte table of (group, room, collect mode) entries for randomized
// items. in ages, a mode >7f means to use &7f as an index to a jump table for
// special cases.
func makeCollectModeTable() string {
	b := new(strings.Builder)

	for _, key := range getOrderedSlotKeys() {
		slot := ItemSlots[key]

		// trees and slots where it doesn't matter (shops, rod)
		if slot.collectMode == 0 {
			continue
		}

		var err error
		if slot.collectMode == collectFall && slot.Treasure != nil &&
			slot.Treasure.id == 0x30 {
			// use falling key mode (no fanfare) if falling item is a key
			_, err = b.Write([]byte{slot.group, slot.room, collectKeyFall})
		} else {
			_, err = b.Write([]byte{slot.group, slot.room, slot.collectMode})
		}
		if err != nil {
			panic(err)
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// returns a byte table (group, room, ID, subID) entries for randomized small
// key drops (and other falling items, but those entries won't be used).
func makeKeyDropTable() string {
	b := new(strings.Builder)

	for _, key := range getOrderedSlotKeys() {
		slot := ItemSlots[key]

		if slot.collectMode != collectFall {
			continue
		}

		// accommodate nil treasures when creating the dummy table before
		// treasures have actually been assigned.
		var err error
		if slot.Treasure == nil {
			_, err = b.Write([]byte{slot.group, slot.room, 0x00, 0x00})
		} else if slot.Treasure.id == 0x30 {
			// make small keys the normal falling variety, with no text box.
			_, err = b.Write([]byte{slot.group, slot.room, 0x30, 0x01})
		} else {
			_, err = b.Write([]byte{slot.group, slot.room,
				slot.Treasure.id, slot.Treasure.subID})
		}
		if err != nil {
			panic(err)
		}
	}

	b.Write([]byte{0xff})
	return b.String()
}

// applies the labels and EOB declarations in the given asmData sets.
func (r *romBanks) applyAsmData(ads []*asmData) {
	// get preset addrs
	for _, ad := range ads {
		for k, v := range ad.Addrs {
			r.addrs[k] = v
		}
	}

	// make placeholders for EOB labels
	for _, ad := range ads {
		for _, items := range ad.Banks {
			for _, item := range items {
				for name := range item {
					r.addrs[name] = 0
				}
			}
		}
	}

	// save original EOB boundaries
	originalEOBs := make([]uint16, 0x40)
	copy(originalEOBs, r.endOfBank)

	// write asm using placeholders for labels, in order to get real addrs
	for _, ad := range ads {
		for bank, items := range ad.Banks {
			for _, item := range items {
				for name, body := range item {
					r.appendAsm(bank, name, body)
				}
			}
		}
	}

	// reset EOB boundaries
	copy(r.endOfBank, originalEOBs)

	// rewrite asm, using real addresses for labels
	for _, ad := range ads {
		for bank, items := range ad.Banks {
			for _, item := range items {
				for name, body := range item {
					r.appendAsm(bank, name, body)
				}
			}
		}
	}
}
