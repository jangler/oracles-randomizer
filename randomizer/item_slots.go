package randomizer

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// an item slot (chest, gift, etc). it references room data and treasure data.
type itemSlot struct {
	treasure                         *treasure
	idAddrs, subidAddrs              []address
	group, room, collectMode, player byte
	moreRooms                        []uint16 // high = group, low = room
	mapTile                          byte     // overworld map coords, yx
	localOnly                        bool     // multiworld
}

// implementes `mutate` from the `mutable` interface.
func (mut *itemSlot) mutate(b []byte) {
	for _, addr := range mut.idAddrs {
		b[addr.fullOffset()] = mut.treasure.id
	}
	for _, addr := range mut.subidAddrs {
		b[addr.fullOffset()] = mut.treasure.subid
	}
	mut.treasure.mutate(b)
}

// helper function for itemSlot.check()
func checkByte(b []byte, addr address, value byte) error {
	if b[addr.fullOffset()] != value {
		return fmt.Errorf("expected %x at %x; found %x",
			value, addr.fullOffset(), b[addr.fullOffset()])
	}
	return nil
}

// implements `check()` from the `mutable` interface.
func (mut *itemSlot) check(b []byte) error {
	// skip zero addresses
	if len(mut.idAddrs) == 0 || mut.idAddrs[0].offset == 0 {
		return nil
	}

	// only check ID addresses, since situational variants and progressive
	// items mess with everything else.
	for _, addr := range mut.idAddrs {
		if err := checkByte(b, addr, mut.treasure.id); err != nil {
			return err
		}
	}

	return nil
}

// raw slot data loaded from yaml.
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

	// optional additional rooms
	MoreRooms []uint16

	Local bool // dummy implies true
}

// like address, but has exported fields for loading from yaml.
type rawAddr struct {
	Bank   byte
	Offset uint16 `yaml:"addr"`
}

// data that can be inferred from a room's music.
type musicData struct {
	MapTile byte
}

// return a map of slot names to slot data. if romState.data is nil, only
// "static" data is loaded.
func (rom *romState) loadSlots() map[string]*itemSlot {
	raws := make(map[string]*rawSlot)

	filename := fmt.Sprintf("/romdata/%s_slots.yaml", gameNames[rom.game])
	if err := yaml.Unmarshal(
		FSMustByte(false, filename), raws); err != nil {
		panic(err)
	}

	allMusic := make(map[string](map[byte]musicData))
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/music.yaml"), allMusic); err != nil {
		panic(err)
	}
	musicMap := allMusic[gameNames[rom.game]]

	m := make(map[string]*itemSlot)
	for name, raw := range raws {
		if raw.Room == 0 && !raw.Dummy {
			panic(name + " room is zero")
		}

		slot := &itemSlot{
			treasure:  rom.treasures[raw.Treasure],
			group:     byte(raw.Room >> 8),
			room:      byte(raw.Room),
			moreRooms: raw.MoreRooms,
			localOnly: raw.Local || raw.Dummy,
		}

		// unspecified map tile = assume overworld
		if raw.MapTile == nil && rom.data != nil {
			musicIndex := getMusicIndex(rom.data, rom.game, slot.group, slot.room)
			if music, ok := musicMap[musicIndex]; ok {
				slot.mapTile = music.MapTile
			} else {
				// nope, definitely not overworld.
				if slot.group > 2 || (slot.group == 2 &&
					(slot.room&0x0f > 0x0d || slot.room&0xf0 > 0xd0)) {
					panic(fmt.Sprintf("invalid room for %s: %04x",
						name, raw.Room))
				}
				slot.mapTile = slot.room
			}
		} else if raw.MapTile != nil {
			slot.mapTile = *raw.MapTile
		}

		if raw.KeyDrop {
			if raw.Collect == "" {
				slot.collectMode = collectModes["drop"]
			}
		} else if raw.Addr != (rawAddr{}) {
			slot.idAddrs = []address{{raw.Addr.Bank, raw.Addr.Offset}}
			slot.subidAddrs = []address{{raw.Addr.Bank, raw.Addr.Offset + 1}}
		} else if raw.ReverseAddr != (rawAddr{}) {
			slot.idAddrs = []address{{
				raw.ReverseAddr.Bank, raw.ReverseAddr.Offset + 1}}
			slot.subidAddrs = []address{{
				raw.ReverseAddr.Bank, raw.ReverseAddr.Offset}}
		} else if raw.Ids != nil {
			slot.idAddrs = make([]address, len(raw.Ids))
			for i, id := range raw.Ids {
				slot.idAddrs[i] = address{id.Bank, id.Offset}
			}

			// allow absence of subids, only because of seed trees
			if raw.SubIds != nil {
				slot.subidAddrs = make([]address, len(raw.SubIds))
				for i, subid := range raw.SubIds {
					slot.subidAddrs[i] = address{subid.Bank, subid.Offset}
				}
			}
		} else if !raw.Dummy && raw.Collect != "d4 pool" && rom.data != nil {
			// try to get chest data for room
			addr := getChestAddr(rom.data, rom.game, slot.group, slot.room)
			if addr != (address{}) {
				slot.idAddrs = []address{{addr.bank, addr.offset}}
				slot.subidAddrs = []address{{addr.bank, addr.offset + 1}}
			} else {
				panic(fmt.Sprintf("invalid raw slot: %s: %#v", name, raw))
			}
		}

		if slot.collectMode != collectModes["drop"] { // drops already set
			if raw.Collect != "" {
				if mode, ok := collectModes[raw.Collect]; ok {
					slot.collectMode = mode
				} else {
					panic("collect mode not found: " + raw.Collect)
				}
			} else {
				if t, ok := rom.treasures[raw.Treasure]; ok {
					slot.collectMode = t.mode
				} else {
					panic("treasure not found: " + raw.Treasure)
				}
			}
		}

		m[name] = slot
	}

	return m
}

// returns the full offset of a room's chest's two-byte entry in the rom.
// returns a zero addr if no chest data is found.
func getChestAddr(b []byte, game int, group, room byte) address {
	ptr := sora(game, address{0x15, 0x4f6c}, address{0x16, 0x5108}).(address)
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

	return address{}
}

// returns the music index for the given room.
func getMusicIndex(b []byte, game int, group, room byte) byte {
	ptr := sora(game, address{0x04, 0x483c}, address{0x04, 0x495c}).(address)

	ptr.offset += uint16(group) * 2
	ptr.offset = uint16(b[ptr.fullOffset()]) +
		uint16(b[ptr.fullOffset()+1])*0x100
	ptr.offset += uint16(room)

	return b[ptr.fullOffset()]
}
