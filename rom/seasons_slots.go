package rom

// seasonsChest constructs a MutableSlot from a treasure name and an address in
// bank $15, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively to chests.
func seasonsChest(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return basicSlot(treasure, 0x15, addr, addr+1, group, room, mode, coords)
}

// seasonsScriptItem constructs a MutableSlot from a treasure name and an
// address in bank $0b, where the ID and sub-ID are two consecutive bytes at
// that address. This applies to most items given by NPCs.
func seasonsScriptItem(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return basicSlot(treasure, 0x0b, addr, addr+1, group, room, mode, coords)
}

// seasonsFoundItem constructs a MutableSlot from a treasure name and an address in
// bank $09, where the sub-ID and ID (in that order) are two consecutive bytes
// at that address. This applies to most items that are found lying around.
func seasonsFoundItem(treasure string, addr uint16,
	group, room, mode, coords byte) *MutableSlot {
	return basicSlot(treasure, 0x09, addr+1, addr, group, room, mode, coords)
}

var seasonsSlots = map[string]*MutableSlot{
	// holodrum
	"eyeglass lake, across bridge": seasonsChest(
		"gasha seed", 0x4f92, 0x00, 0xb8, collectChest, 0xb8),
	"maku tree": &MutableSlot{
		treasureName: "gnarled key",
		idAddrs:      []Addr{{0x15, 0x613a}, {0x09, 0x7e16}},
		subIDAddrs:   []Addr{{0x15, 0x613d}, {0x09, 0x7e19}},
		group:        0x02,
		room:         0x0b,
		collectMode:  collectFall,
		mapCoords:    0xc9,
	},
	"horon village SW chest": seasonsChest(
		"rupees, 20", 0x4f7e, 0x00, 0xf5, collectChest, 0xf5),
	"horon village SE chest": seasonsChest(
		"rupees, 20", 0x4f82, 0x00, 0xf9, collectChest, 0xf9),
	"holly's house": seasonsScriptItem(
		"shovel", 0x6a6c, 0x03, 0xa3, collectFind2, 0x7f),
	"chest on top of D2": seasonsChest(
		"gasha seed", 0x4f86, 0x00, 0x8e, collectChest, 0x8e),
	"blaino prize": seasonsScriptItem(
		"gasha seed", 0x64cc, 0x03, 0xb4, collectFind1, 0x78),
	"floodgate keeper's house": seasonsFoundItem(
		"floodgate key", 0x6281, 0x03, 0xb5, collectFind1, 0x62),
	"spool swamp cave": &MutableSlot{
		treasureName: "square jewel",
		idAddrs:      []Addr{{0x0b, 0x7395}},
		subIDAddrs:   []Addr{{0x0b, 0x7399}},
		group:        0x04,
		room:         0xfa,
		collectMode:  collectChest,
		mapCoords:    0xc2,
	},
	"moblin keep": seasonsChest(
		"piece of heart", 0x4f8e, 0x00, 0x5b, collectChest, 0x5b),
	"master diver's challenge": seasonsChest(
		"master's plaque", 0x510a, 0x05, 0xbc, collectChest, 0x2e),
	"master diver's reward": seasonsScriptItem( // addr set at EOB
		"flippers", 0x0000, 0x05, 0xbd, collectNil, 0x2e), // special case
	"spring banana tree": seasonsFoundItem(
		"spring banana", 0x66c6, 0x00, 0x0f, collectFind2, 0x0f),
	"goron mountain, across pits": seasonsFoundItem(
		"dragon key", 0x62a3, 0x00, 0x1a, collectFind1, 0x1a),
	"mt. cucco, platform cave": seasonsFoundItem( // addr set at EOB
		"green joy ring", 0x0000, 0x05, 0xbb, collectFall, 0x1f),
	"diving spot outside D4": seasonsScriptItem(
		"pyramid jewel", 0x734e, 0x07, 0xe5, collectUnderwater, 0x1d),
	"black beast's chest": seasonsChest(
		"x-shaped jewel", 0x4f8a, 0x00, 0xf4, collectChest, 0xf4),
	"old man in treehouse": seasonsScriptItem(
		"round jewel", 0x7332, 0x03, 0x94, collectFind2, 0xb5),
	"lost woods": &MutableSlot{
		treasureName: "sword 2",
		idAddrs:      []Addr{{0x0b, 0x6418}, {0x0b, 0x641f}},
		subIDAddrs:   []Addr{{0x0b, 0x6419}, {0x0b, 0x6420}},
		group:        0x00,
		room:         0xc9,
		collectMode:  collectFind1,
		mapCoords:    0x40,
	},
	"samasa desert pit": &MutableSlot{
		treasureName: "rusty bell",
		idAddrs:      []Addr{{0x09, 0x648d}, {0x0b, 0x60b1}},
		subIDAddrs:   []Addr{{0x09, 0x648c}},
		group:        0x05,
		room:         0xd2,
		collectMode:  collectFind2,
		mapCoords:    0xbf,
	},
	"samasa desert chest": seasonsChest(
		"rang ring L-1", 0x4f9a, 0x00, 0xff, collectChest, 0xff),
	"western coast, beach chest": seasonsChest(
		"blast ring", 0x4f96, 0x00, 0xe3, collectChest, 0xe3),
	"western coast, in house": seasonsChest(
		"bombs, 10", 0x4fac, 0x03, 0x88, collectChest, 0xd2),
	"cave south of mrs. ruul": seasonsChest(
		"octo ring", 0x5081, 0x04, 0xe0, collectChest, 0xb3),
	"cave north of D1": seasonsChest(
		"quicksand ring", 0x5085, 0x04, 0xe1, collectChest, 0x87),
	"cave outside D2": seasonsChest(
		"moblin ring", 0x50fe, 0x05, 0xb3, collectChest, 0x8e),
	"woods of winter, 1st cave": seasonsChest(
		"rupees, 30", 0x5102, 0x05, 0xb4, collectChest, 0x7d),
	"sunken city, summer cave": seasonsChest(
		"gasha seed", 0x5106, 0x05, 0xb5, collectChest, 0x4f),
	"chest in master diver's cave": seasonsChest(
		"rupees, 50", 0x510e, 0x05, 0xbd, collectChest, 0x2e),
	"dry eyeglass lake, east cave": seasonsChest(
		"piece of heart", 0x5112, 0x05, 0xc0, collectChest, 0xaa),
	"chest in goron mountain": seasonsChest(
		"armor ring L-2", 0x511a, 0x05, 0xc8, collectChest, 0x18),
	"natzu region, across water": seasonsChest(
		"rupees, 50", 0x5122, 0x05, 0x0e, collectChest, 0x49),
	"mt. cucco, talon's cave": seasonsChest(
		"subrosian ring", 0x511e, 0x05, 0xb6, collectChest, 0x1b),
	"tarm ruins, under tree": seasonsChest(
		"gasha seed", 0x4fa8, 0x03, 0x9b, collectChest, 0x10),
	"eastern suburbs, on cliff": seasonsChest(
		"gasha seed", 0x5089, 0x04, 0xf7, collectChest, 0xcc),
	"dry eyeglass lake, west cave": &MutableSlot{
		treasureName: "rupees, 100",
		idAddrs:      []Addr{{0x0b, 0x73a1}},
		subIDAddrs:   []Addr{{0x0b, 0x73a5}},
		group:        0x04,
		room:         0xfb,
		collectMode:  collectChest,
		mapCoords:    0xa7,
	},
	"woods of winter, 2nd cave": &MutableSlot{
		treasureName: "gasha seed",
		idAddrs:      []Addr{{0x0a, 0x5003}},
		subIDAddrs:   []Addr{{0x0a, 0x5008}},
		group:        0x05,
		room:         0x12,
		collectMode:  collectChest,
		mapCoords:    0x7e,
	},

	// dummy slots for bombs and shield
	"shop, 20 rupees": &MutableSlot{
		treasureName: "bombs, 10",
		group:        0x03,
		room:         0xa6,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},
	"shop, 30 rupees": &MutableSlot{
		treasureName: "wooden shield",
		group:        0x03,
		room:         0xa6,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},

	"shop, 150 rupees": &MutableSlot{
		treasureName: "strange flute",
		idAddrs:      []Addr{{0x08, 0x4ce8}},
		subIDAddrs:   []Addr{{0x08, 0x4ce9}},
		group:        0x03,
		room:         0xa6,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},
	"member's shop 1": &MutableSlot{
		treasureName: "satchel 2",
		idAddrs:      []Addr{{0x08, 0x4cce}},
		subIDAddrs:   []Addr{{0x08, 0x4ccf}},
		group:        0x03,
		room:         0xb0,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},
	"member's shop 2": &MutableSlot{
		treasureName: "gasha seed",
		idAddrs:      []Addr{{0x08, 0x4cd2}},
		subIDAddrs:   []Addr{{0x08, 0x4cd3}},
		group:        0x03,
		room:         0xb0,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},
	"member's shop 3": &MutableSlot{
		treasureName: "treasure map",
		idAddrs:      []Addr{{0x08, 0x4cd8}},
		subIDAddrs:   []Addr{{0x08, 0x4cd9}},
		group:        0x03,
		room:         0xb0,
		collectMode:  collectNil,
		mapCoords:    0xe6,
	},

	// subrosia
	"tower of winter": seasonsScriptItem(
		"winter", 0x4fc5, 0x05, 0xf2, collectFind1, 0x9a),
	"tower of summer": seasonsScriptItem(
		"summer", 0x4fb9, 0x05, 0xf8, collectFind1, 0xb0),
	"tower of spring": seasonsScriptItem(
		"spring", 0x4fb5, 0x05, 0xf5, collectFind1, 0x1e),
	"tower of autumn": seasonsScriptItem(
		"autumn", 0x4fc1, 0x05, 0xfb, collectFind1, 0xb9),
	"subrosian dance hall": seasonsScriptItem(
		"boomerang 1", 0x6646, 0x03, 0x95, collectFind2, 0x9a),
	"temple of seasons": &MutableSlot{
		treasureName: "rod",
		idAddrs:      []Addr{{0x15, 0x70ce}},
		subIDAddrs:   []Addr{{0x15, 0x70cc}},
		group:        0x03,
		room:         0xac,
		collectMode:  collectNil,
		mapCoords:    0x9a,
	},
	"subrosia seaside": &MutableSlot{ // addrs set dynamically at EOB
		treasureName: "star ore",
		idAddrs:      []Addr{{0x08, 0x0000}},
		subIDAddrs:   []Addr{{0x08, 0x0000}},
		group:        0x01,
		room:         0x66,
		collectMode:  collectDig,
		mapCoords:    0xb0,
	},
	"subrosian wilds chest": seasonsChest(
		"blue ore", 0x4f9f, 0x01, 0x41, collectChest, 0x1e),
	"subrosia village chest": seasonsChest(
		"red ore", 0x4fa3, 0x01, 0x58, collectChest, 0xb9),
	"subrosia, open cave": seasonsChest(
		"gasha seed", 0x5095, 0x04, 0xf1, collectChest, 0x25),
	"subrosia, locked cave": seasonsChest(
		"gasha seed", 0x5116, 0x05, 0xc6, collectChest, 0xb0),
	"subrosia market, 1st item": &MutableSlot{
		treasureName: "ribbon",
		idAddrs:      []Addr{{0x09, 0x77da}},
		subIDAddrs:   []Addr{{0x09, 0x77db}},
		group:        0x03,
		room:         0xa0,
		collectMode:  collectNil,
		mapCoords:    0xb0,
	},
	"subrosia market, 2nd item": &MutableSlot{
		treasureName: "rare peach stone",
		idAddrs:      []Addr{{0x09, 0x77e2}},
		subIDAddrs:   []Addr{{0x09, 0x77e3}},
		group:        0x03,
		room:         0xa0,
		collectMode:  collectNil,
		mapCoords:    0xb0,
	},
	"subrosia market, 5th item": &MutableSlot{
		treasureName: "member's card",
		idAddrs:      []Addr{{0x09, 0x77f4}},
		subIDAddrs:   []Addr{{0x09, 0x77f5}},
		group:        0x03,
		room:         0xa0,
		collectMode:  collectNil,
		mapCoords:    0xb0,
	},
	"great furnace": &MutableSlot{ // addrs set dynamically at EOB
		treasureName: "hard ore",
		idAddrs:      []Addr{{0x15, 0x0000}, {0x09, 0x66eb}},
		subIDAddrs:   []Addr{{0x15, 0x0000}, {0x09, 0x66ea}},
		group:        0x03,
		room:         0x8e,
		collectMode:  collectFind2,
		mapCoords:    0xb9,
	},
	"subrosian smithy": &MutableSlot{
		treasureName: "shield L-2",
		idAddrs:      []Addr{{0x15, 0x62be}},
		subIDAddrs:   []Addr{{0x15, 0x62b4}},
		group:        0x03,
		room:         0x97,
		collectMode:  collectFind2,
		mapCoords:    0x25,
	},

	// hero's cave
	"d0 sword chest": &MutableSlot{
		treasureName: "sword 1",
		idAddrs:      []Addr{{0x0a, 0x7b90}},
		paramAddrs:   []Addr{{0x0a, 0x7b92}},
		textAddrs:    []Addr{{0x0a, 0x7b9c}},
		gfxAddrs:     []Addr{{0x3f, 0x6676}},
		group:        0x04,
		room:         0x06,
		collectMode:  collectNil,
		mapCoords:    0xd4,
	},
	"d0 rupee chest": seasonsChest(
		"rupees, 30", 0x4fb5, 0x04, 0x05, collectChest, 0xd4),

	// d1
	"d1 basement": seasonsFoundItem(
		"satchel 1", 0x66b1, 0x06, 0x09, collectFind2, 0x96),
	"d1 block-pushing room": seasonsChest(
		"gasha seed", 0x4fbd, 0x04, 0x0d, collectChest, 0x96),
	"d1 railway chest": seasonsChest(
		"bombs, 10", 0x4fc5, 0x04, 0x10, collectChest, 0x96),
	"d1 floormaster room": seasonsChest(
		"discovery ring", 0x4fd1, 0x04, 0x17, collectChest, 0x96),
	"d1 lever room": seasonsChest(
		"compass", 0x4fc1, 0x04, 0x0f, collectChest2, 0x96),
	"d1 stalfos chest": seasonsChest(
		"dungeon map", 0x4fd5, 0x04, 0x19, collectChest2, 0x96),
	"d1 goriya chest": seasonsChest(
		"d1 boss key", 0x4fcd, 0x04, 0x14, collectChest, 0x96),
	"d1 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0x12, collectAppear2, 0x96),

	// d2
	"d2 moblin chest": seasonsChest(
		"bracelet", 0x4fe1, 0x04, 0x2a, collectChest, 0x8d),
	"d2 roller chest": seasonsChest(
		"rupees, 10", 0x4fd9, 0x04, 0x1f, collectChest, 0x8d),
	"d2 left from entrance": seasonsChest(
		"rupees, 5", 0x4ff5, 0x04, 0x38, collectChest, 0x8d),
	"d2 pot chest": seasonsChest(
		"dungeon map", 0x4fe5, 0x04, 0x2b, collectChest2, 0x8d),
	"d2 rope chest": seasonsChest(
		"compass", 0x4ff1, 0x04, 0x36, collectChest2, 0x8d),
	"d2 terrace chest": seasonsChest(
		"d2 boss key", 0x4fdd, 0x04, 0x24, collectChest, 0x8d),
	"d2 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0x29, collectAppear2, 0x8d),

	// d3
	"d3 mimic chest": seasonsChest(
		"feather 1", 0x5015, 0x04, 0x50, collectChest, 0x60),
	"d3 water room": seasonsChest(
		"rupees, 30", 0x4ff9, 0x04, 0x41, collectChest, 0x60),
	"d3 quicksand terrace": seasonsChest(
		"gasha seed", 0x5001, 0x04, 0x44, collectChest, 0x60),
	"d3 moldorm chest": seasonsChest(
		"bombs, 10", 0x5019, 0x04, 0x54, collectChest, 0x60),
	"d3 trampoline chest": seasonsChest(
		"compass", 0x5009, 0x04, 0x4d, collectChest2, 0x60),
	"d3 bombed wall chest": seasonsChest(
		"dungeon map", 0x5011, 0x04, 0x51, collectChest2, 0x60),
	"d3 giant blade room": seasonsChest(
		"d3 boss key", 0x4ffd, 0x04, 0x46, collectChest, 0x60),
	"d3 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0x53, collectAppear2, 0x60),

	// d4
	"d4 cracked floor room": seasonsChest(
		"slingshot 1", 0x502d, 0x04, 0x73, collectChest, 0x1d),
	"d4 north of entrance": seasonsChest(
		"bombs, 10", 0x5031, 0x04, 0x7f, collectChest, 0x1d),
	"d4 maze chest": seasonsChest(
		"dungeon map", 0x5025, 0x04, 0x69, collectChest2, 0x1d),
	"d4 water ring room": seasonsChest(
		"compass", 0x5035, 0x04, 0x83, collectChest2, 0x1d),
	"d4 dive spot": seasonsScriptItem(
		"d4 boss key", 0x4c0b, 0x04, 0x6c, collectDive, 0x1d),
	"d4 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0x5f, collectAppear2, 0x1d),

	// d5
	"d5 magnet ball chest": seasonsChest(
		"magnet gloves", 0x503d, 0x04, 0x89, collectChest, 0x89),
	"d5 terrace chest": seasonsChest(
		"rupees, 100", 0x5041, 0x04, 0x97, collectChest, 0x8a),
	"d5 gibdo/zol chest": seasonsChest(
		"dungeon map", 0x5039, 0x04, 0x8f, collectChest2, 0x8f),
	"d5 spiral chest": seasonsChest(
		"compass", 0x5049, 0x04, 0x9d, collectChest2, 0x8a),
	"d5 basement": seasonsScriptItem(
		"d5 boss key", 0x4c22, 0x06, 0x8b, collectFind2, 0x8a),
	"d5 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0x8c, collectAppear2, 0x8a),

	// d6
	"d6 armos hall": seasonsChest(
		"boomerang 2", 0x507d, 0x04, 0xd0, collectChest, 0x00),
	"d6 crystal trap room": seasonsChest(
		"rupees, 10", 0x505d, 0x04, 0xaf, collectChest, 0x00),
	"d6 1F east": seasonsChest(
		"rupees, 5", 0x5065, 0x04, 0xb3, collectChest, 0x00),
	"d6 2F gibdo chest": seasonsChest(
		"bombs, 10", 0x5069, 0x04, 0xbf, collectChest, 0x00),
	"d6 2F armos chest": seasonsChest(
		"rupees, 5", 0x5075, 0x04, 0xc3, collectChest, 0x00),
	"d6 beamos room": seasonsChest(
		"compass", 0x5059, 0x04, 0xad, collectChest2, 0x00),
	"d6 1F terrace": seasonsChest(
		"dungeon map", 0x5061, 0x04, 0xb0, collectChest2, 0x00),
	"d6 escape room": seasonsChest(
		"d6 boss key", 0x5079, 0x04, 0xc4, collectChest, 0x00),
	"d6 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x04, 0xd5, collectAppear2, 0x00),

	// d7
	"d7 spike chest": seasonsChest(
		"feather 2", 0x509e, 0x05, 0x44, collectChest, 0xd0),
	"d7 maze chest": seasonsChest(
		"rupees, 1", 0x509a, 0x05, 0x43, collectChest, 0xd0),
	"d7 right of entrance": seasonsChest(
		"power ring L-1", 0x50b6, 0x05, 0x5a, collectChest, 0xd0),
	"d7 bombed wall chest": seasonsChest(
		"compass", 0x50aa, 0x05, 0x52, collectChest2, 0xd0),
	"d7 quicksand chest": seasonsChest(
		"dungeon map", 0x50b2, 0x05, 0x58, collectChest2, 0xd0),
	"d7 stalfos chest": seasonsChest(
		"d7 boss key", 0x50a6, 0x05, 0x48, collectChest, 0xd0),
	"d7 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x05, 0x50, collectAppear2, 0xd0),

	// d8
	"d8 armos chest": seasonsChest(
		"slingshot 2", 0x50da, 0x05, 0x8d, collectChest, 0x04),
	"d8 SW lava chest": seasonsChest(
		"bombs, 10", 0x50ba, 0x05, 0x6a, collectChest, 0x04),
	"d8 three eyes chest": seasonsChest(
		"steadfast ring", 0x50c6, 0x05, 0x7d, collectChest, 0x04),
	"d8 spike room": seasonsChest(
		"compass", 0x50d2, 0x05, 0x8b, collectChest2, 0x04),
	"d8 magnet ball room": seasonsChest(
		"dungeon map", 0x50de, 0x05, 0x8e, collectChest2, 0x04),
	"d8 pols voice chest": seasonsChest(
		"d8 boss key", 0x50ca, 0x05, 0x80, collectChest, 0x04),
	"d8 boss": seasonsChest( // EOB addr
		"heart container", 0x0000, 0x05, 0x64, collectAppear2, 0x04),

	// don't use this slot; no one knows about it and it's not required for
	// anything in a normal playthrough
	// "ring box L-2 gift": seasonsScriptItem("ring box L-2", 0x5c18),

	// these are "fake" item slots in that they don't slot real treasures
	"horon village seed tree": &MutableSlot{
		treasureName: "ember tree seeds",
		idAddrs:      []Addr{{0x0d, 0x68fb}},
	},
	"woods of winter seed tree": &MutableSlot{
		treasureName: "mystery tree seeds",
		idAddrs:      []Addr{{0x0d, 0x68fe}},
	},
	"north horon seed tree": &MutableSlot{
		treasureName: "scent tree seeds",
		idAddrs:      []Addr{{0x0d, 0x6901}},
	},
	"spool swamp seed tree": &MutableSlot{
		treasureName: "pegasus tree seeds",
		idAddrs:      []Addr{{0x0d, 0x6904}},
	},
	"sunken city seed tree": &MutableSlot{
		treasureName: "gale tree seeds",
		idAddrs:      []Addr{{0x0d, 0x6907}},
	},
	"tarm ruins seed tree": &MutableSlot{
		treasureName: "gale tree seeds",
		idAddrs:      []Addr{{0x0d, 0x690a}},
	},
}
