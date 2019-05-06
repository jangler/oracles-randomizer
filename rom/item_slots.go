package rom

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// A MutableSlot is an item slot (chest, gift, etc). It references room data
// and treasure data.
type MutableSlot struct {
	Treasure *Treasure

	treasureName             string
	idAddrs, subIDAddrs      []Addr
	group, room, collectMode byte
	mapCoords                byte // overworld map coords, yx
}

// Mutate replaces the given IDs, subIDs, and other applicable data in the ROM.
func (ms *MutableSlot) Mutate(b []byte) error {
	for _, addr := range ms.idAddrs {
		b[addr.fullOffset()] = ms.Treasure.id
	}
	for _, addr := range ms.subIDAddrs {
		b[addr.fullOffset()] = ms.Treasure.subID
	}

	return ms.Treasure.Mutate(b)
}

// helper function for MutableSlot.Check
func check(b []byte, addr Addr, value byte) error {
	if b[addr.fullOffset()] != value {
		return fmt.Errorf("expected %x at %x; found %x",
			value, addr.fullOffset(), b[addr.fullOffset()])
	}
	return nil
}

// Check verifies that the slot's data matches the given ROM data.
func (ms *MutableSlot) Check(b []byte) error {
	// skip zero addresses
	if len(ms.idAddrs) == 0 || ms.idAddrs[0].offset == 0 {
		return nil
	}

	// only check ID addresses, since situational variants and progressive
	// items mess with everything else.
	for _, addr := range ms.idAddrs {
		if err := check(b, addr, ms.Treasure.id); err != nil {
			return err
		}
	}

	return nil
}

// basicSlot constucts a MutableSlot from a treasure name, bank number, and an
// address for each its ID and sub-ID. Most slots fit this pattern.
func basicSlot(treasure string, bank byte, idOffset, subIDOffset uint16,
	group, room, mode, coords byte) *MutableSlot {
	return &MutableSlot{
		treasureName: treasure,
		idAddrs:      []Addr{{bank, idOffset}},
		subIDAddrs:   []Addr{{bank, subIDOffset}},
		group:        group,
		room:         room,
		collectMode:  mode,
		mapCoords:    coords,
	}
}

// keyDropSlot constructs a MutableSlot for a small key drop. the mutable
// itself is a dummy and does not have an address; the data is used to
// construct a table of small key drops.
func keyDropSlot(treasure string, group, room, coords byte) *MutableSlot {
	return &MutableSlot{
		treasureName: treasure,
		group:        group,
		room:         room,
		collectMode:  collectFall,
		mapCoords:    coords,
	}
}

var ItemSlots map[string]*MutableSlot

type rawSlot struct {
	// required
	Treasure string
	Room     uint16

	// required if not == low byte of room or in dungeon.
	MapTile *byte

	// pick one, or default to chest
	Addr        rawAddr // for id, then subid
	ReverseAddr rawAddr // for subid, then id
	Ids, SubIds []rawAddr
	KeyDrop     bool
	Dummy       bool

	// optional override
	Collect string
}

// has exported fields for loading from yaml.
type rawAddr struct {
	Bank   byte
	Offset uint16 `yaml:"addr"`
}

// data that can be inferred from a room's music.
type musicData struct {
	MapTile byte
}

// loads slot data from the rom into the given map.
func loadSlots(b []byte, game int, m map[string]*MutableSlot) {
	raws := make(map[string]*rawSlot)

	filename := fmt.Sprintf("/rom/%s_slots.yaml", gameNames[game])
	if err := yaml.Unmarshal(
		FSMustByte(false, filename), raws); err != nil {
		panic(err)
	}

	allMusic := make(map[string](map[byte]musicData))
	if err := yaml.Unmarshal(
		FSMustByte(false, "/rom/music.yaml"), allMusic); err != nil {
		panic(err)
	}
	musicMap := allMusic[gameNames[game]]

	for name, raw := range raws {
		if raw.Room == 0 && !raw.Dummy {
			panic(name + " room is zero")
		}

		slot := &MutableSlot{
			treasureName: raw.Treasure,
			group:        byte(raw.Room >> 8),
			room:         byte(raw.Room),
		}

		// unspecified map tile = assume overworld
		if raw.MapTile == nil {
			musicIndex := getMusicIndex(b, game, slot.group, slot.room)
			if music, ok := musicMap[musicIndex]; ok {
				slot.mapCoords = music.MapTile
			} else {
				if slot.group > 2 || (slot.group == 2 &&
					(slot.room&0x0f > 0x0d || slot.room&0xf0 > 0xd0)) {
					panic(fmt.Sprintf("invalid room for %s: %04x",
						name, raw.Room))
				}
				slot.mapCoords = slot.room
			}
		} else {
			slot.mapCoords = *raw.MapTile
		}

		if raw.KeyDrop {
			slot.collectMode = collectFall
		} else if raw.Addr != (rawAddr{}) {
			slot.idAddrs = []Addr{{raw.Addr.Bank, raw.Addr.Offset}}
			slot.subIDAddrs = []Addr{{raw.Addr.Bank, raw.Addr.Offset + 1}}
		} else if raw.ReverseAddr != (rawAddr{}) {
			slot.idAddrs = []Addr{{
				raw.ReverseAddr.Bank, raw.ReverseAddr.Offset + 1}}
			slot.subIDAddrs = []Addr{{
				raw.ReverseAddr.Bank, raw.ReverseAddr.Offset}}
		} else if raw.Ids != nil {
			slot.idAddrs = make([]Addr, len(raw.Ids))
			for i, id := range raw.Ids {
				slot.idAddrs[i] = Addr{id.Bank, id.Offset}
			}

			// allow no subIds, only because of seed trees
			if raw.SubIds != nil {
				slot.subIDAddrs = make([]Addr, len(raw.SubIds))
				for i, subid := range raw.SubIds {
					slot.subIDAddrs[i] = Addr{subid.Bank, subid.Offset}
				}
			}
		} else if !raw.Dummy && raw.Collect != "d4 pool" {
			// try to get chest data for room
			addr := getChestAddr(b, game, slot.group, slot.room)
			if addr != (Addr{}) {
				slot.idAddrs = []Addr{{addr.bank, addr.offset}}
				slot.subIDAddrs = []Addr{{addr.bank, addr.offset + 1}}
			} else {
				panic(fmt.Sprintf("invalid raw slot: %s: %#v", name, raw))
			}
		}

		// TODO: have a Dx small key (drop) treasure or something??
		// TODO: even better, just have a small key (drop) treasure (etc) and
		// get dungeon automatically based on music or something
		if slot.collectMode == 0 { // key drop slots have theirs set already
			if raw.Collect != "" {
				if mode, ok := collectModesByName[raw.Collect]; ok {
					slot.collectMode = mode
				} else {
					panic("collect mode not found: " + raw.Collect)
				}
			} else {
				if t, ok := Treasures[raw.Treasure]; ok {
					slot.collectMode = t.mode
				} else {
					panic("treasure not found: " + raw.Treasure)
				}
			}
		}

		m[name] = slot
	}
}

// returns the full offset of a room's chest's two-byte entry in the rom.
// returns a zero addr if no chest data is found.
func getChestAddr(b []byte, game int, group, room byte) Addr {
	var ptr Addr
	if game == GameSeasons {
		ptr = Addr{0x15, 0x4f6c}
	} else {
		ptr = Addr{0x16, 0x5108}
	}

	ptr.offset += uint16(group) * 2
	ptr.offset = uint16(b[ptr.fullOffset()]) +
		uint16(b[ptr.fullOffset()+1])*0x100

	for {
		info := b[ptr.fullOffset()]
		if info == 0xff {
			break
		}

		chest_room := b[ptr.fullOffset()+1]
		if chest_room == room {
			ptr.offset += 2
			return ptr
		}

		ptr.offset += 4
	}

	return Addr{}
}

// returns the music index for the given room.
func getMusicIndex(b []byte, game int, group, room byte) byte {
	var ptr Addr
	if game == GameSeasons {
		ptr = Addr{0x04, 0x483c}
	} else {
		ptr = Addr{0x04, 0x495c}
	}

	ptr.offset += uint16(group) * 2
	ptr.offset = uint16(b[ptr.fullOffset()]) +
		uint16(b[ptr.fullOffset()+1])*0x100
	ptr.offset += uint16(room)

	return b[ptr.fullOffset()]
}
