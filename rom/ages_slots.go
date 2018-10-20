package rom

// agesChest constructs a MutableSlot from a treasure name and an address in
// bank $16, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively to chests.
func agesChest(treasure string, addr uint16, group, room byte) *MutableSlot {
	return BasicSlot(treasure, 0x16, addr, addr+1,
		group, room, collectChest, 0)
}

// for items given by script command de.
func agesScriptItem(treasure string, addr uint16,
	group, room byte) *MutableSlot {
	return BasicSlot(treasure, 0x0c, addr, addr+1, group, room, collectFind2, 0)
}

var agesSlots = map[string]*MutableSlot{
	// overworld present
	// "impa's gift": CustomSlot("sword 1", 0x00, 0x39, collectFind2, 0x39),
	"nayru's house": BasicSlot("harp", 0x0b, 0x6828, 0x6827,
		0x03, 0xae, collectFind2, 0x3a),
	"maku tree": &MutableSlot{
		treasureName: "satchel 1",
		IDAddrs:      []Addr{{0x15, 0x70e0}, {0x15, 0x7115}},
		SubIDAddrs:   []Addr{{0x15, 0x70e3}, {0x15, 0x7118}},
		group:        0x00,
		room:         0x38,
		collectMode:  collectFall,
	},
	"grave under tree": BasicSlot("graveyard key", 0x10, 0x750d, 0x750c,
		0x05, 0xed, collectFall, 0x8d),
	"cheval's test":        agesScriptItem("flippers", 0x723b, 0x05, 0xbf),
	"cheval's invention":   agesScriptItem("cheval rope", 0x7232, 0x05, 0xb6),
	"tingle's gift":        agesScriptItem("island chart", 0x7e20, 0x00, 0x79),
	"tingle's upgrade":     agesScriptItem("satchel 2", 0x7e7a, 0x00, 0x79),
	"lynna city chest":     agesChest("rupees, 30", 0x511e, 0x00, 0x49),
	"fairies' woods chest": agesChest("rupees, 50", 0x5122, 0x00, 0x84),
	"fairies' coast chest": agesChest("green holy ring", 0x5126, 0x00, 0x91),
	"zora seas chest":      agesChest("whimsical ring", 0x512e, 0x00, 0xd5),
	"talus peaks chest":    agesChest("gasha seed", 0x5132, 0x00, 0x63),
	"ruined keep":          agesChest("armor ring L-1", 0x5144, 0x02, 0xbe),
	"nuun cave":            agesChest("light ring L-1", 0x5154, 0x02, 0xf4),
	"zora village present": agesChest("gasha seed", 0x515c, 0x02, 0xc0),
	"D6 entrance pool":     agesChest("toss ring", 0x5161, 0x03, 0x0e),
	// "sea of storms cave present": nil, // special linked chest?
	"mayor plen's house":      agesChest("green luck ring", 0x5171, 0x03, 0xf9),
	"crescent seafloor cave":  agesChest("piece of heart", 0x5175, 0x03, 0xfd),
	"goron's hiding place":    agesChest("gold joy ring", 0x52f7, 0x05, 0xbd),
	"ridge base cave present": agesChest("rupees, 50", 0x52fb, 0x05, 0xb9),
	"ridge NE cave present":   agesChest("gasha seed", 0x52ff, 0x05, 0xee),
	"goron diamond cave":      agesChest("bombs, 10", 0x5303, 0x05, 0xdd),
	"ridge west cave":         agesChest("rupees, 30", 0x5307, 0x05, 0xc0),
	"zora NW cave":            agesChest("blue luck ring", 0x531f, 0x05, 0xc7),
	"zora palace chest":       agesChest("rupees, 200", 0x532b, 0x05, 0xac),

	// overworld past
	"black tower worker": agesScriptItem("shovel", 0x65e3, 0x04, 0xe1),
	"sea of no return":   agesChest("blue ring", 0x5137, 0x01, 0x6d),
	"bomb goron head":    agesChest("rupees, 100", 0x5148, 0x02, 0xfc),
	"tokay cave past":    agesChest("gasha seed", 0x514c, 0x02, 0xce),
	"zora cave past":     agesChest("red holy ring", 0x5158, 0x02, 0x4f),
	"ridge bush cave":    agesChest("rupees, 100", 0x5165, 0x03, 0x1f),
	// "sea of storms cave past": nil, // special linked chest?
	"deku forest cave west": agesChest("rupees, 30", 0x52f3, 0x05, 0xb5),
	"ridge past diamonds":   agesChest("rupees, 50", 0x530f, 0x05, 0xe1),
	"ridge past base":       agesChest("gasha seed", 0x5313, 0x05, 0xe0),
	"deku forest cave east": agesChest("gasha seed", 0x5317, 0x05, 0xb3),
	"palace chest":          agesChest("gold luck ring", 0x531b, 0x05, 0xcb),
	"tokay crystal cave":    agesChest("gasha seed", 0x5323, 0x05, 0xca),
	"tokay pot cave":        agesChest("power ring L-2", 0x532f, 0x05, 0xf7),

	// dungeons TODO

	// seed trees work differently in ages; the seed type is determined by the
	// high nybble of the tree sub ID, and the low nybble is used to identify it
	// for regrowth purposes. so these can't be set directly like ordinary item
	// slots can.
	"symmetry city tree":      &MutableSlot{treasureName: "gale tree seeds"},
	"south lynna tree":        &MutableSlot{treasureName: "ember tree seeds"},
	"crescent island tree":    &MutableSlot{treasureName: "scent tree seeds"},
	"zora village tree":       &MutableSlot{treasureName: "gale tree seeds"},
	"rolling ridge west tree": &MutableSlot{treasureName: "pegasus tree seeds"},
	"ambi's palace tree":      &MutableSlot{treasureName: "scent tree seeds"},
	"rolling ridge east tree": &MutableSlot{treasureName: "mystery tree seeds"},
	"deku forest tree":        &MutableSlot{treasureName: "mystery tree seeds"},
}
