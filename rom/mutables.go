package rom

import (
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
		codeMutables["filterMusic"].New[3] = 0x18
	}
}

// SetTreewarp sets treewarp on or off in the modified ROM. By default, it is
// on.
func SetTreewarp(treewarp bool) {
	if !treewarp {
		codeMutables["treeWarp"].New[5] = 0x18
	}
}

// SetAnimal sets the flute type and Natzu region type based on a companion
// number 1 to 3.
func SetAnimal(companion int) {
	codeMutables["animalRegion"].New =
		[]byte{byte(companion + 0x0a)}
	codeMutables["flutePalette"].New =
		[]byte{byte(0x10*(4-companion) + 3)}
}

// key = area name (as in asm/vars.yaml), id = season index (spring -> winter).
func SetSeason(key string, id byte) {
	codeMutables[key].New[0] = id
}

// get a collated map of all mutables, *except* for treasures which do not
// appear in the seed. this allows things like the three seasons flutes having
// different data but the same address.
func getAllMutables() map[string]Mutable {
	slotMutables := make(map[string]Mutable)
	treasureMutables := make(map[string]Mutable)
	otherMutables := make(map[string]Mutable, len(codeMutables))
	for k, v := range ItemSlots {
		if v.Treasure == nil {
			log.Fatalf("treasure for %s is nil", k)
		}
		if v.Treasure.addr.offset != 0 {
			treasureMutables[FindTreasureName(v.Treasure)] = v.Treasure
		}
		slotMutables[k] = v
	}
	for k, v := range codeMutables {
		otherMutables[k] = v
	}

	mutableSets := []map[string]Mutable{
		treasureMutables,
		slotMutables,
		otherMutables,
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
