package rom

const (
	collectMakuTree    = 0x80
	collectTargetCarts = 0x81
	collectBigBang     = 0x82
)

// agesChest constructs a MutableSlot from a treasure name and an address in
// bank $16, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively to chests.
func agesChest(treasure string, addr uint16, group, room byte) *MutableSlot {
	if _, ok := agesTreasures[treasure]; !ok {
		panic("treasure " + treasure + " does not exist")
	}
	mode := agesTreasures[treasure].mode
	return BasicSlot(treasure, 0x16, addr, addr+1, group, room, mode, 0)
}

// for items given by script command de.
func agesScriptItem(treasure string, addr uint16,
	group, room byte) *MutableSlot {
	return BasicSlot(treasure, 0x0c, addr, addr+1, group, room, collectFind2, 0)
}

// same as agesScriptItem, but for scripts in bank 15 that are copied to the
// c3xx buffer. some non-scripted items in bank 15 also use this format.
func agesBufferItem(treasure string, addr uint16,
	group, room byte) *MutableSlot {
	return BasicSlot(treasure, 0x15, addr, addr+1, group, room, collectFind2, 0)
}

var agesSlots = map[string]*MutableSlot{
	// overworld present
	"starting chest": BasicSlot("sword 1", 0x00, 0x10f8, 0x10f7,
		0x00, 0x39, collectChest, 0x39),
	"nayru's house": BasicSlot("harp 1", 0x0b, 0x6828, 0x6827,
		0x03, 0xae, collectFind2, 0x3a),
	"maku tree": &MutableSlot{
		treasureName: "satchel 1",
		IDAddrs:      []Addr{{0x15, 0x70e0}, {0x15, 0x7115}},
		SubIDAddrs:   []Addr{{0x15, 0x70e3}, {0x15, 0x7118}},
		group:        0x00,
		room:         0x38,
		collectMode:  collectMakuTree,
	},
	"grave under tree": BasicSlot("graveyard key", 0x10, 0x750d, 0x750c,
		0x05, 0xed, collectFall, 0x8d),
	"cheval's test":      agesScriptItem("flippers 1", 0x723b, 0x05, 0xbf),
	"cheval's invention": agesScriptItem("cheval rope", 0x7232, 0x05, 0xb6),
	"south shore dirt": BasicSlot("ricky's gloves", 0x0a, 0x5e3d, 0x5e3c,
		0x00, 0x98, collectDigPile, 0x98),
	"tingle's gift":    agesScriptItem("island chart", 0x7e20, 0x00, 0x79),
	"tingle's upgrade": agesScriptItem("satchel 2", 0x7e7a, 0x00, 0x79),
	"shop, 150 rupees": BasicSlot("dimitri's flute", 0x09, 0x4511, 0x4512,
		0x02, 0x5e, collectFind2, 0x68),
	"defeat great moblin": agesScriptItem("bomb flower", 0x757d, 0x00, 0x09),
	"goron elder":         agesBufferItem("crown key", 0x7386, 0x05, 0xc3),
	"target carts 1": &MutableSlot{
		treasureName: "rock brisket",
		IDAddrs:      []Addr{{0x15, 0x66e8}, {0x0c, 0x6e71}},
		SubIDAddrs:   []Addr{{0x15, 0x66e9}, {0x0c, 0x6e72}},
		group:        0x05,
		room:         0xd8,
		collectMode:  collectTargetCarts,
	},
	"target carts 2": &MutableSlot{ // second addrs set dynamically at EOB
		treasureName: "boomerang",
		IDAddrs:      []Addr{{0x15, 0x66f0}, {0x0c, 0x0000}},
		SubIDAddrs:   []Addr{{0x15, 0x66f1}, {0x0c, 0x0000}},
		group:        0x05,
		room:         0xd8,
		collectMode:  collectTargetCarts,
	},
	"goron dancing":      agesScriptItem("brother emblem", 0x698c, 0x02, 0xed),
	"trade rock brisket": agesBufferItem("goron vase", 0x6b2c, 0x02, 0xfd),
	"trade goron vase":   agesBufferItem("goronade", 0x6b23, 0x02, 0xff),
	"big bang game": &MutableSlot{
		treasureName: "old mermaid key",
		IDAddrs:      []Addr{{0x15, 0x6742}, {0x0c, 0x707a}},
		SubIDAddrs:   []Addr{{0x15, 0x6743}, {0x0c, 0x707b}},
		group:        0x03,
		room:         0x3e,
		collectMode:  collectBigBang,
	},
	"shooting gallery":       agesBufferItem("lava juice", 0x5285, 0x03, 0xe7),
	"trade lava juice":       agesScriptItem("goron letter", 0x6ee9, 0x03, 0x1f),
	"rescue nayru":           agesBufferItem("harp 3", 0x54f1, 0x00, 0x38),
	"king zora":              agesScriptItem("library key", 0x7ae4, 0x05, 0xab),
	"library present":        agesBufferItem("book of seals", 0x5db9, 0x05, 0xc8),
	"zora's reward":          agesScriptItem("zora scale", 0x7c48, 0x02, 0xa0),
	"piratian captain":       agesBufferItem("tokay eyeball", 0x7969, 0x05, 0xf8),
	"old zora":               agesBufferItem("broken sword", 0x61ad, 0x02, 0xf5),
	"lynna city chest":       agesChest("rupees, 30", 0x511e, 0x00, 0x49),
	"fairies' woods chest":   agesChest("rupees, 50", 0x5122, 0x00, 0x84),
	"fairies' coast chest":   agesChest("green holy ring", 0x5126, 0x00, 0x91),
	"zora seas chest":        agesChest("whimsical ring", 0x512e, 0x00, 0xd5),
	"talus peaks chest":      agesChest("gasha seed", 0x5132, 0x00, 0x63),
	"ruined keep":            agesChest("armor ring L-1", 0x5144, 0x02, 0xbe),
	"nuun highlands cave":    agesChest("light ring L-1", 0x5154, 0x02, 0xf4),
	"zora village present":   agesChest("gasha seed", 0x515c, 0x02, 0xc0),
	"D6 entrance pool":       agesChest("toss ring", 0x5161, 0x03, 0x0e),
	"sea of storms present":  agesChest("gasha seed", 0x5169, 0x03, 0xe8),
	"mayor plen's house":     agesChest("green luck ring", 0x5171, 0x03, 0xf9),
	"crescent seafloor cave": agesChest("piece of heart", 0x5175, 0x03, 0xfd),
	"goron's hiding place":   agesChest("gold joy ring", 0x52f7, 0x05, 0xbd),
	"ridge base present":     agesChest("rupees, 50", 0x52fb, 0x05, 0xb9),
	"ridge NE cave present":  agesChest("gasha seed", 0x52ff, 0x05, 0xee),
	"goron diamond cave":     agesChest("bombs, 10", 0x5303, 0x05, 0xdd),
	"ridge west cave":        agesChest("rupees, 30", 0x5307, 0x05, 0xc0),
	"zora NW cave":           agesChest("blue luck ring", 0x531f, 0x05, 0xc7),
	"zora palace chest":      agesChest("rupees, 200", 0x532b, 0x05, 0xac),

	// overworld past
	"black tower worker": agesScriptItem("shovel", 0x65e3, 0x04, 0xe1),
	"deku forest soldier": agesScriptItem(
		"bombs, 10", 0x0000, 0x01, 0x72), // addr set dynamically at EOB
	"wild tokay game": agesBufferItem(
		"scent seedling", 0x5bbb, 0x02, 0xde), // not actually a script
	"hidden tokay cave":     agesBufferItem("iron shield", 0x5b36, 0x05, 0xe9),
	"symmetry city brother": agesBufferItem("tuni nut", 0x7929, 0x03, 0x6f),
	"tokkey's composition":  agesBufferItem("harp 2", 0x76cf, 0x03, 0x8f),
	"goron dancing past":    agesScriptItem("mermaid key", 0x699f, 0x02, 0xef),
	"library past":          agesBufferItem("fairy powder", 0x5dd8, 0x05, 0xe4),
	"sea of no return":      agesChest("blue ring", 0x5137, 0x01, 0x6d),
	"bomb goron head":       agesChest("rupees, 100", 0x5148, 0x02, 0xfc),
	"tokay bomb cave":       agesChest("gasha seed", 0x514c, 0x02, 0xce),
	"zora cave past":        agesChest("red holy ring", 0x5158, 0x02, 0x4f),
	"ridge bush cave":       agesChest("rupees, 100", 0x5165, 0x03, 0x1f),
	"sea of storms past":    agesChest("pegasus ring", 0x516d, 0x03, 0xff),
	"deku forest cave west": agesChest("rupees, 30", 0x52f3, 0x05, 0xb5),
	"ridge diamonds past":   agesChest("rupees, 50", 0x530f, 0x05, 0xe1),
	"ridge base past":       agesChest("gasha seed", 0x5313, 0x05, 0xe0),
	"deku forest cave east": agesChest("gasha seed", 0x5317, 0x05, 0xb3),
	"ambi's palace chest":   agesChest("gold luck ring", 0x531b, 0x05, 0xcb),
	"tokay crystal cave":    agesChest("gasha seed", 0x5323, 0x05, 0xca),
	"tokay pot cave":        agesChest("power ring L-2", 0x532f, 0x05, 0xf7),

	// dungeons
	"d1 button chest":          agesChest("gasha seed", 0x517e, 0x04, 0x15),
	"d1 crystal room":          agesChest("power ring L-1", 0x518a, 0x04, 0x1c),
	"d1 crossroads":            agesChest("compass", 0x518e, 0x04, 0x1d),
	"d1 west terrace":          agesChest("discovery ring", 0x5192, 0x04, 0x1f),
	"d1 pot chest":             agesChest("d1 boss key", 0x5196, 0x04, 0x23),
	"d1 east terrace":          agesChest("dungeon map", 0x519a, 0x04, 0x25),
	"d1 basement":              agesScriptItem("bracelet 1", 0x4bbb, 0x06, 0x10),
	"d2 color room":            agesChest("d2 boss key", 0x51a2, 0x04, 0x3e),
	"d2 bombed terrace":        agesChest("dungeon map", 0x51a6, 0x04, 0x40),
	"d2 moblin platform":       agesChest("gasha seed", 0x51aa, 0x04, 0x41),
	"d2 rope room":             agesChest("compass", 0x51ae, 0x04, 0x45),
	"d2 thwomp shelf":          agesScriptItem("rupees, 30", 0x4c0f, 0x06, 0x27),
	"d2 thwomp tunnel":         agesScriptItem("feather", 0x4c0a, 0x06, 0x28),
	"d3 bridge chest":          agesChest("rupees, 20", 0x51b6, 0x04, 0x4e),
	"d3 boss key chest":        agesChest("d3 boss key", 0x51ba, 0x04, 0x50),
	"d3 torch chest":           agesChest("gasha seed", 0x51be, 0x04, 0x55),
	"d3 compass chest":         agesChest("compass", 0x51c2, 0x04, 0x56),
	"d3 shooter chest":         agesChest("seed shooter", 0x51c6, 0x04, 0x58),
	"d3 hall of bushes":        agesChest("rupees, 30", 0x51ca, 0x04, 0x5c),
	"d3 crossroads":            agesChest("gasha seed", 0x51ce, 0x04, 0x60),
	"d3 pols voice chest":      agesChest("dungeon map", 0x51d2, 0x04, 0x65),
	"d4 lava pot chest":        agesChest("d4 boss key", 0x51de, 0x04, 0x7a),
	"d4 small floor puzzle":    agesChest("switch hook 1", 0x51e2, 0x04, 0x87),
	"d4 first chest":           agesChest("compass", 0x51e6, 0x04, 0x8b),
	"d4 minecart chest":        agesChest("dungeon map", 0x51ea, 0x04, 0x8f),
	"d5 red peg chest":         agesChest("rupees, 50", 0x51f6, 0x04, 0x99),
	"d5 owl puzzle":            agesChest("d5 boss key", 0x51fa, 0x04, 0x9b),
	"d5 six-statue puzzle":     agesChest("cane", 0x520a, 0x04, 0xa5),
	"d5 diamond chest":         agesChest("compass", 0x520e, 0x04, 0xad),
	"d5 blue peg chest":        agesChest("dungeon map", 0x521a, 0x04, 0xbe),
	"d6 present vire chest":    agesChest("flippers 2", 0x524f, 0x05, 0x13),
	"d6 present RNG chest":     agesChest("d6 boss key", 0x525b, 0x05, 0x1c),
	"d6 present diamond chest": agesChest("dungeon map", 0x525f, 0x05, 0x1d),
	"d6 present beamos chest":  agesChest("rupees, 10", 0x5263, 0x05, 0x1f),
	"d6 present channel chest": agesChest("compass", 0x526b, 0x05, 0x25),
	"d6 past spear chest":      agesChest("rupees, 30", 0x5273, 0x05, 0x2e),
	"d6 past color room":       agesChest("compass", 0x527f, 0x05, 0x3f),
	"d6 past pool chest":       agesChest("dungeon map", 0x5283, 0x05, 0x41),
	"d6 past wizzrobe chest":   agesChest("gasha seed", 0x5287, 0x05, 0x45),
	"d7 pot island chest":      agesChest("like-like ring", 0x528b, 0x05, 0x4c),
	"d7 stairway chest":        agesChest("gasha seed", 0x528f, 0x05, 0x4d),
	"d7 miniboss chest":        agesChest("switch hook 2", 0x5293, 0x05, 0x4e),
	"d7 crab chest":            agesChest("compass", 0x529b, 0x05, 0x54),
	"d7 spike chest":           agesChest("dungeon map", 0x52a7, 0x05, 0x65),
	"d7 hallway chest":         agesChest("gasha seed", 0x52ab, 0x05, 0x6a),
	"d7 post-hallway chest":    agesChest("d7 boss key", 0x52af, 0x05, 0x6c),
	"d8 B3F chest":             agesChest("d8 boss key", 0x52bb, 0x05, 0x79),
	"d8 isolated chest":        agesChest("dungeon map", 0x52cb, 0x05, 0x85),
	"d8 sarcophagus chest":     agesChest("gasha seed", 0x52db, 0x05, 0x9f),
	"d8 blue peg chest":        agesChest("compass", 0x52e3, 0x05, 0xa4),
	"d8 floor puzzle":          agesChest("bracelet 2", 0x52eb, 0x05, 0xa6),
	"d8 tile room":             agesChest("gasha seed", 0x52ef, 0x05, 0x91),

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
