package rom

import (
	"bytes"
	"fmt"
	"log"
)

// A Mutable is a memory data that can be changed by the randomizer.
type Mutable interface {
	Mutate([]byte) error // change ROM bytes
	Check([]byte) error  // verify that the mutable matches the ROM
}

// A MutableRange is a length of mutable bytes starting at a given address.
type MutableRange struct {
	Addrs    []Addr
	Old, New []byte
}

// MutableByte returns a special case of MutableRange with a range of a single
// byte.
func MutableByte(addr Addr, old, new byte) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   []byte{old},
		New:   []byte{new},
	}
}

// MutableWord returns a special case of MutableRange with a range of a two
// bytes.
func MutableWord(addr Addr, old, new uint16) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   []byte{byte(old >> 8), byte(old)},
		New:   []byte{byte(new >> 8), byte(new)},
	}
}

// MutableString returns a MutableRange constructed from the bytes in two
// strings.
func MutableString(addr Addr, old, new string) *MutableRange {
	return &MutableRange{
		Addrs: []Addr{addr},
		Old:   bytes.NewBufferString(old).Bytes(),
		New:   bytes.NewBufferString(new).Bytes(),
	}
}

// MutableStrings returns a MutableRange constructed from the bytes in two
// strings, at multiple addresses.
func MutableStrings(addrs []Addr, old, new string) *MutableRange {
	return &MutableRange{
		Addrs: addrs,
		Old:   bytes.NewBufferString(old).Bytes(),
		New:   bytes.NewBufferString(new).Bytes(),
	}
}

// Mutate replaces bytes in its range.
func (mr *MutableRange) Mutate(b []byte) error {
	for _, addr := range mr.Addrs {
		offset := addr.fullOffset()
		for i, value := range mr.New {
			b[offset+i] = value
		}
	}
	return nil
}

// Check verifies that the range matches the given ROM data.
func (mr *MutableRange) Check(b []byte) error {
	for _, addr := range mr.Addrs {
		offset := addr.fullOffset()
		for i, value := range mr.Old {
			if b[offset+i] != value {
				return fmt.Errorf("expected %x at %x; found %x",
					mr.Old[i], offset+i, b[offset+i])
			}
		}
	}
	return nil
}

// SetMusic sets music on or off in the modified ROM. By default, it is off.
func SetMusic(music bool) {
	if music {
		mut := codeMutables["filterMusic"].(*MutableRange)
		mut.New[3] = 0x18
	}
}

// SetTreewarp sets treewarp on or off in the modified ROM. By default, it is
// on.
func SetTreewarp(treewarp bool) {
	if !treewarp {
		mut := codeMutables["treeWarp"].(*MutableRange)
		mut.New[5] = 0x18
	}
}

// SetAnimal sets the flute type and Natzu region type based on a companion
// number 1 to 3.
func SetAnimal(companion int) {
	codeMutables["animalRegion"].(*MutableRange).New =
		[]byte{byte(companion + 0x0a)}
	codeMutables["flutePalette"].(*MutableRange).New =
		[]byte{byte(0x10*(4-companion) + 3)}
}

// these mutables have fixed addresses and don't reference other mutables. try
// to generally order them by address, unless a grouping between mutables in
// different banks makes more sense.
var fixedMutables map[string]Mutable

// get a collated map of all mutables, *except* for treasures which do not
// appear in the seed. this allows things like the three seasons flutes having
// different data but the same address.
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	treasureMutables := make(map[string]Mutable)
	for k, v := range ItemSlots {
		if v.Treasure == nil {
			log.Fatalf("treasure for %s is nil", k)
		}
		if v.Treasure.addr.offset != 0 {
			treasureMutables[FindTreasureName(v.Treasure)] = v.Treasure
		}
		slotMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		fixedMutables,
		treasureMutables,
		slotMutables,
		codeMutables,
	}

	// initialize master map w/ adequate capacity
	count := 0
	for _, set := range mutableSets {
		count += len(set)
	}
	allMutables := make(map[string]Mutable, count)

	// add mutables to master map
	for _, set := range mutableSets {
		for k, v := range set {
			if _, ok := allMutables[k]; ok {
				log.Fatalf("duplicate mutable key: %s", k)
			}
			allMutables[k] = v
		}
	}

	return allMutables
}

// FindAddr returns the name of a mutable that covers the given address, or an
// empty string if none is found.
func FindAddr(bank byte, addr uint16) string {
	muts := getAllMutables()
	offset := (&Addr{bank, addr}).fullOffset()

	for name, mut := range muts {
		switch mut := mut.(type) {
		case *MutableRange:
			for _, addr := range mut.Addrs {
				if offset >= addr.fullOffset() &&
					offset < addr.fullOffset()+len(mut.New) {
					return name
				}
			}
		case *MutableSlot:
			for _, addr := range mut.idAddrs {
				if offset == addr.fullOffset() {
					return name
				}
			}
			for _, addr := range mut.subIDAddrs {
				if offset == addr.fullOffset() {
					return name
				}
			}
		case *Treasure:
			if offset >= mut.addr.fullOffset() &&
				offset < mut.addr.fullOffset()+4 {
				return name
			}
		default:
			panic("unknown type for mutable: " + name)
		}
	}

	return ""
}
