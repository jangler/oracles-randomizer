package rom

// special collection modes that jump to custom code, for when there are
// multiple modes required in the same room.
const (
	collectMakuTree    = 0x80
	collectTargetCarts = 0x81
	collectBigBang     = 0x82
	collectLavaJuice   = 0x83
)

// agesChest constructs a MutableSlot from a treasure name and an address in
// bank $16, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively to chests.
func agesChest(treasure string, addr uint16, group, room byte) *MutableSlot {
	if _, ok := AgesTreasures[treasure]; !ok {
		panic("treasure " + treasure + " does not exist")
	}
	mode := AgesTreasures[treasure].mode
	return basicSlot(treasure, 0x16, addr, addr+1, group, room, mode, 0)
}

// for boss items.
func agesHC(group, room byte) *MutableSlot {
	return basicSlot("heart container", 0x15, 0, 0,
		group, room, collectAppear2, 0)
}

// for items given by script command de.
func agesScriptItem(treasure string, addr uint16,
	group, room byte) *MutableSlot {
	return basicSlot(treasure, 0x0c, addr, addr+1, group, room, collectFind2, 0)
}

// same as agesScriptItem, but for scripts in bank 15 that are copied to the
// c3xx buffer. some non-scripted items in bank 15 also use this format.
func agesBufferItem(treasure string, addr uint16,
	group, room byte) *MutableSlot {
	return basicSlot(treasure, 0x15, addr, addr+1, group, room, collectFind2, 0)
}

var AgesSlots = map[string]*MutableSlot{
	// overworld present
	"starting chest": basicSlot("sword 1", 0x00, 0x10f8, 0x10f7,
		0x00, 0x39, collectChest, 0x39),
	"nayru's house": basicSlot("harp 1", 0x0b, 0x6828, 0x6827,
		0x03, 0xae, collectFind2, 0x3a),
	"maku tree": &MutableSlot{
		treasureName: "satchel 1",
		idAddrs:      []Addr{{0x15, 0x70e0}, {0x15, 0x7115}},
		subIDAddrs:   []Addr{{0x15, 0x70e3}, {0x15, 0x7118}},
		group:        0x00,
		room:         0x38,
		collectMode:  collectMakuTree,
	},
	"grave under tree": basicSlot("graveyard key", 0x10, 0x750d, 0x750c,
		0x05, 0xed, collectFall, 0x8d),
	"graveyard poe": &MutableSlot{
		treasureName: "sword 2",
		idAddrs:      []Addr{{0x15, 0x6188}},
		subIDAddrs:   []Addr{{0x15, 0x6189}},
		group:        0x00,
		room:         0x7c,
		collectMode:  collectFind2,
	},
	"cheval's test":      agesScriptItem("flippers 1", 0x723b, 0x05, 0xbf),
	"cheval's invention": agesScriptItem("cheval rope", 0x7232, 0x05, 0xb6),
	"south shore dirt": basicSlot("ricky's gloves", 0x0a, 0x5e3d, 0x5e3c,
		0x00, 0x98, collectDigPile, 0x98),
	"balloon guy's gift":    agesScriptItem("island chart", 0x7e20, 0x00, 0x79),
	"balloon guy's upgrade": agesScriptItem("satchel 2", 0x7e7a, 0x00, 0x79),
	"shop, 150 rupees": basicSlot("strange flute", 0x09, 0x4511, 0x4512,
		0x02, 0x5e, collectFind2, 0x68),
	"defeat great moblin": agesScriptItem("bomb flower", 0x757d, 0x00, 0x09),
	"goron elder":         agesBufferItem("crown key", 0x7386, 0x05, 0xc3),
	"target carts 1": &MutableSlot{
		treasureName: "rock brisket",
		idAddrs:      []Addr{{0x15, 0x66e8}, {0x0c, 0x6e71}},
		subIDAddrs:   []Addr{{0x15, 0x66e9}, {0x0c, 0x6e72}},
		group:        0x05,
		room:         0xd8,
		collectMode:  collectTargetCarts,
	},
	"target carts 2": &MutableSlot{ // second addrs set dynamically at EOB
		treasureName: "boomerang",
		idAddrs:      []Addr{{0x15, 0x66f0}, {0x0c, 0x0000}},
		subIDAddrs:   []Addr{{0x15, 0x66f1}, {0x0c, 0x0000}},
		group:        0x05,
		room:         0xd8,
		collectMode:  collectTargetCarts,
	},
	"goron dance present": agesScriptItem("brother emblem", 0x698c, 0x02, 0xed),
	"trade rock brisket":  agesBufferItem("goron vase", 0x6b2c, 0x02, 0xfd),
	"trade goron vase":    agesBufferItem("goronade", 0x6b23, 0x02, 0xff),
	"big bang game": &MutableSlot{
		treasureName: "old mermaid key",
		idAddrs:      []Addr{{0x15, 0x6742}, {0x0c, 0x707a}},
		subIDAddrs:   []Addr{{0x15, 0x6743}, {0x0c, 0x707b}},
		group:        0x03,
		room:         0x3e,
		collectMode:  collectBigBang,
	},
	"goron shooting gallery": agesBufferItem("lava juice", 0x5285, 0x03, 0xe7),
	"trade lava juice": basicSlot("goron letter", 0x0c, 0x6ee9, 0x6eea,
		0x03, 0x1f, collectLavaJuice, 0x1c),
	"rescue nayru": basicSlot("harp 3", 0x15, 0x54f1, 0x54f2,
		0x00, 0x38, collectMakuTree, 0x38),
	"king zora":            agesScriptItem("library key", 0x7ae4, 0x05, 0xab),
	"library present":      agesBufferItem("book of seals", 0x5db9, 0x05, 0xc8),
	"zora's reward":        agesScriptItem("zora scale", 0x7c48, 0x02, 0xa0),
	"piratian captain":     agesBufferItem("tokay eyeball", 0x7969, 0x05, 0xf8),
	"lynna city chest":     agesChest("rupees, 30", 0x511e, 0x00, 0x49),
	"fairies' woods chest": agesChest("rupees, 50", 0x5122, 0x00, 0x84),
	"fairies' coast chest": agesChest("green holy ring", 0x5126, 0x00, 0x91),
	"zora seas chest":      agesChest("whimsical ring", 0x512e, 0x00, 0xd5),
	"talus peaks chest":    agesChest("gasha seed", 0x5132, 0x00, 0x63),
	"under moblin keep":    agesChest("armor ring L-1", 0x5144, 0x02, 0xbe),
	"nuun highlands cave": &MutableSlot{
		// has three different rooms depending on animal
		treasureName: "light ring L-1",
		idAddrs:      []Addr{{0x16, 0x5150}, {0x16, 0x5154}, {0x16, 0x5327}},
		subIDAddrs:   []Addr{{0x16, 0x5151}, {0x16, 0x5155}, {0x16, 0x5328}},
		group:        0x02,
		room:         0xf4,
		collectMode:  collectChest,
	},
	"zora village present":  agesChest("gasha seed", 0x515c, 0x02, 0xc0),
	"pool in d6 entrance":   agesChest("toss ring", 0x5161, 0x03, 0x0e),
	"mayor plen's house":    agesChest("green luck ring", 0x5171, 0x03, 0xf9),
	"under crescent island": agesChest("piece of heart", 0x5175, 0x03, 0xfd),
	"goron's hiding place":  agesChest("gold joy ring", 0x52f7, 0x05, 0xbd),
	"ridge base chest":      agesChest("rupees, 50", 0x52fb, 0x05, 0xb9),
	"ridge NE cave present": agesChest("gasha seed", 0x52ff, 0x05, 0xee),
	"goron diamond cave":    agesChest("bombs, 10", 0x5303, 0x05, 0xdd),
	"ridge west cave":       agesChest("rupees, 30", 0x5307, 0x05, 0xc0),
	"zora NW cave":          agesChest("blue luck ring", 0x531f, 0x05, 0xc7),
	"zora palace chest":     agesChest("rupees, 200", 0x532b, 0x05, 0xac),

	// overworld past
	"black tower worker": agesScriptItem("shovel", 0x65e3, 0x04, 0xe1),
	"deku forest soldier": agesScriptItem(
		"bombs, 10", 0x0000, 0x01, 0x72), // addr set dynamically at EOB
	"wild tokay game": agesBufferItem(
		"scent seedling", 0x5bbb, 0x02, 0xde), // not actually a script
	"hidden tokay cave":        agesBufferItem("iron shield", 0x5b36, 0x05, 0xe9),
	"symmetry city brother":    agesBufferItem("tuni nut", 0x7929, 0x03, 0x6e),
	"tokkey's composition":     agesBufferItem("harp 2", 0x76cf, 0x03, 0x8f),
	"goron dance, with letter": agesScriptItem("mermaid key", 0x699f, 0x02, 0xef),
	"library past":             agesBufferItem("fairy powder", 0x5dd8, 0x05, 0xe4),
	"sea of no return":         agesChest("blue ring", 0x5137, 0x01, 0x6d),
	"bomb goron head":          agesChest("rupees, 100", 0x5148, 0x02, 0xfc),
	"tokay bomb cave":          agesChest("gasha seed", 0x514c, 0x02, 0xce),
	"fisher's island cave":     agesChest("red holy ring", 0x5158, 0x02, 0x4f),
	"ridge bush cave": basicSlot("rupees, 100", 0x16, 0x5165, 0x5166,
		0x03, 0x1f, collectLavaJuice, 0x1c),
	"sea of storms past":    agesChest("pegasus ring", 0x516d, 0x03, 0xff),
	"deku forest cave west": agesChest("rupees, 30", 0x52f3, 0x05, 0xb5),
	"ridge diamonds past":   agesChest("rupees, 50", 0x530f, 0x05, 0xe1),
	"ridge base past":       agesChest("gasha seed", 0x5313, 0x05, 0xe0),
	"deku forest cave east": agesChest("gasha seed", 0x5317, 0x05, 0xb3),
	"ambi's palace chest":   agesChest("gold luck ring", 0x531b, 0x05, 0xcb),
	"tokay crystal cave":    agesChest("gasha seed", 0x5323, 0x05, 0xca),
	"tokay pot cave":        agesChest("power ring L-2", 0x532f, 0x05, 0xf7),

	// dungeons
	"d1 one-button chest":      agesChest("gasha seed", 0x517e, 0x04, 0x15),
	"d1 two-button chest":      agesChest("d1 small key", 0x5182, 0x04, 0x16),
	"d1 wide room":             agesChest("d1 small key", 0x5186, 0x04, 0x1a),
	"d1 crystal room":          agesChest("power ring L-1", 0x518a, 0x04, 0x1c),
	"d1 crossroads":            agesChest("compass", 0x518e, 0x04, 0x1d),
	"d1 west terrace":          agesChest("discovery ring", 0x5192, 0x04, 0x1f),
	"d1 pot chest":             agesChest("d1 boss key", 0x5196, 0x04, 0x23),
	"d1 east terrace":          agesChest("dungeon map", 0x519a, 0x04, 0x25),
	"d1 ghini drop":            keyDropSlot("d1 small key", 0x04, 0x1e, 0),
	"d1 basement":              agesScriptItem("bracelet 1", 0x4bbb, 0x06, 0x10),
	"d1 boss":                  agesHC(0x04, 0x13),
	"d2 basement chest":        agesChest("d2 small key", 0x519e, 0x04, 0x30),
	"d2 color room":            agesChest("d2 boss key", 0x51a2, 0x04, 0x3e),
	"d2 bombed terrace":        agesChest("dungeon map", 0x51a6, 0x04, 0x40),
	"d2 moblin platform":       agesChest("gasha seed", 0x51aa, 0x04, 0x41),
	"d2 rope room":             agesChest("compass", 0x51ae, 0x04, 0x45),
	"d2 ladder chest":          agesChest("d2 small key", 0x51b2, 0x04, 0x48),
	"d2 basement drop":         keyDropSlot("d2 small key", 0x04, 0x2e, 0),
	"d2 moblin drop":           keyDropSlot("d2 small key", 0x04, 0x39, 0),
	"d2 statue puzzle":         keyDropSlot("d2 small key", 0x04, 0x42, 0),
	"d2 thwomp shelf":          agesScriptItem("rupees, 30", 0x4c0f, 0x06, 0x27),
	"d2 thwomp tunnel":         agesScriptItem("feather", 0x4c0a, 0x06, 0x28),
	"d2 boss":                  agesHC(0x06, 0x2b),
	"d3 bridge chest":          agesChest("rupees, 20", 0x51b6, 0x04, 0x4e),
	"d3 B1F east":              agesChest("d3 boss key", 0x51ba, 0x04, 0x50),
	"d3 torch chest":           agesChest("gasha seed", 0x51be, 0x04, 0x55),
	"d3 conveyor belt room":    agesChest("compass", 0x51c2, 0x04, 0x56),
	"d3 mimic room":            agesChest("seed shooter", 0x51c6, 0x04, 0x58),
	"d3 bush beetle room":      agesChest("rupees, 30", 0x51ca, 0x04, 0x5c),
	"d3 crossroads":            agesChest("gasha seed", 0x51ce, 0x04, 0x60),
	"d3 pols voice chest":      agesChest("dungeon map", 0x51d2, 0x04, 0x65),
	"d3 moldorm drop":          keyDropSlot("d3 small key", 0x04, 0x4b, 0),
	"d3 armos drop":            keyDropSlot("d3 small key", 0x04, 0x5e, 0),
	"d3 statue drop":           keyDropSlot("d3 small key", 0x04, 0x61, 0),
	"d3 six-block drop":        keyDropSlot("d3 small key", 0x04, 0x64, 0),
	"d3 boss":                  agesHC(0x04, 0x4a),
	"d4 large floor puzzle":    agesChest("d4 small key", 0x51d6, 0x04, 0x6f),
	"d4 second crystal switch": agesChest("d4 small key", 0x51da, 0x04, 0x74),
	"d4 lava pot chest":        agesChest("d4 boss key", 0x51de, 0x04, 0x7a),
	"d4 small floor puzzle":    agesChest("switch hook 1", 0x51e2, 0x04, 0x87),
	"d4 first chest":           agesChest("compass", 0x51e6, 0x04, 0x8b),
	"d4 minecart chest":        agesChest("dungeon map", 0x51ea, 0x04, 0x8f),
	"d4 cube chest":            agesChest("d4 small key", 0x51ee, 0x04, 0x90),
	"d4 first crystal switch":  agesChest("d4 small key", 0x51f2, 0x04, 0x92),
	"d4 color tile drop":       keyDropSlot("d4 small key", 0x04, 0x7b, 0),
	"d4 boss":                  agesHC(0x04, 0x6b),
	"d5 red peg chest":         agesChest("rupees, 50", 0x51f6, 0x04, 0x99),
	"d5 owl puzzle":            agesChest("d5 boss key", 0x51fa, 0x04, 0x9b),
	"d5 two-statue puzzle":     agesChest("d5 small key", 0x51fe, 0x04, 0x9e),
	"d5 like-like chest":       agesChest("d5 small key", 0x5202, 0x04, 0x9f),
	"d5 dark room":             agesChest("d5 small key", 0x5206, 0x04, 0xa3),
	"d5 six-statue puzzle":     agesChest("cane", 0x520a, 0x04, 0xa5),
	"d5 diamond chest":         agesChest("compass", 0x520e, 0x04, 0xad),
	"d5 eyes chest":            agesChest("d5 small key", 0x5212, 0x04, 0xba),
	"d5 three-statue puzzle":   agesChest("d5 small key", 0x5216, 0x04, 0xbc),
	"d5 blue peg chest":        agesChest("dungeon map", 0x521a, 0x04, 0xbe),
	"d5 boss":                  agesHC(0x04, 0xbf),
	"d6 present vire chest":    agesChest("flippers 2", 0x524f, 0x05, 0x13),
	"d6 present spinner chest": agesChest("d6 present small key", 0x5253, 0x05, 0x14),
	"d6 present rope chest":    agesChest("d6 present small key", 0x5257, 0x05, 0x1b),
	"d6 present RNG chest":     agesChest("d6 boss key", 0x525b, 0x05, 0x1c),
	"d6 present diamond chest": agesChest("dungeon map", 0x525f, 0x05, 0x1d),
	"d6 present beamos chest":  agesChest("rupees, 10", 0x5263, 0x05, 0x1f),
	"d6 present cube chest":    agesChest("d6 present small key", 0x5267, 0x05, 0x21),
	"d6 present channel chest": agesChest("compass", 0x526b, 0x05, 0x25),
	"d6 past diamond chest":    agesChest("d6 past small key", 0x526f, 0x05, 0x2c),
	"d6 past spear chest":      agesChest("rupees, 30", 0x5273, 0x05, 0x2e),
	"d6 past rope chest":       agesChest("d6 past small key", 0x5277, 0x05, 0x31),
	"d6 past stalfos chest":    agesChest("d6 past small key", 0x527b, 0x05, 0x3c),
	"d6 past color room":       agesChest("compass", 0x527f, 0x05, 0x3f),
	"d6 past pool chest":       agesChest("dungeon map", 0x5283, 0x05, 0x41),
	"d6 past wizzrobe chest":   agesChest("gasha seed", 0x5287, 0x05, 0x45),
	"d6 boss":                  agesHC(0x05, 0x36),
	"d7 pot island chest":      agesChest("like-like ring", 0x528b, 0x05, 0x4c),
	"d7 stairway chest":        agesChest("gasha seed", 0x528f, 0x05, 0x4d),
	"d7 miniboss chest":        agesChest("switch hook 2", 0x5293, 0x05, 0x4e),
	"d7 cane/diamond puzzle":   agesChest("d7 small key", 0x5297, 0x05, 0x53),
	"d7 crab chest":            agesChest("compass", 0x529b, 0x05, 0x54),
	"d7 left wing":             agesChest("d7 small key", 0x529f, 0x05, 0x5f),
	"d7 right wing":            agesChest("d7 small key", 0x52a3, 0x05, 0x64),
	"d7 spike chest":           agesChest("dungeon map", 0x52a7, 0x05, 0x65),
	"d7 hallway chest":         agesChest("gasha seed", 0x52ab, 0x05, 0x6a),
	"d7 post-hallway chest":    agesChest("d7 boss key", 0x52af, 0x05, 0x6c),
	"d7 3F terrace":            agesChest("d7 small key", 0x52b3, 0x05, 0x72),
	"d7 boxed chest":           agesChest("d7 small key", 0x52b7, 0x05, 0x50),
	"d7 flower room":           keyDropSlot("d7 small key", 0x05, 0x4b, 0),
	"d7 diamond puzzle":        keyDropSlot("d7 small key", 0x05, 0x55, 0),
	"d7 boss":                  agesHC(0x05, 0x62),
	"d8 B3F chest":             agesChest("d8 boss key", 0x52bb, 0x05, 0x79),
	"d8 maze chest":            agesChest("d8 small key", 0x52bf, 0x05, 0x7b),
	"d8 NW slate chest":        agesChest("slate", 0x52c3, 0x05, 0x7c),
	"d8 NE slate chest":        agesChest("slate", 0x52c7, 0x05, 0x7e),
	"d8 ghini chest":           agesChest("dungeon map", 0x52cb, 0x05, 0x85),
	"d8 SE slate chest":        agesChest("slate", 0x52cf, 0x05, 0x92),
	"d8 SW slate chest":        agesChest("slate", 0x52d3, 0x05, 0x94),
	"d8 B1F NW chest":          agesChest("d8 small key", 0x52d7, 0x05, 0x97),
	"d8 sarcophagus chest":     agesChest("gasha seed", 0x52db, 0x05, 0x9f),
	"d8 blade trap chest":      agesChest("d8 small key", 0x52df, 0x05, 0xa3),
	"d8 blue peg chest":        agesChest("compass", 0x52e3, 0x05, 0xa4),
	"d8 1F chest":              agesChest("d8 small key", 0x52e7, 0x05, 0xa7),
	"d8 floor puzzle":          agesChest("bracelet 2", 0x52eb, 0x05, 0xa6),
	"d8 tile room":             agesChest("gasha seed", 0x52ef, 0x05, 0x91),
	"d8 stalfos": basicSlot(
		"d8 small key", 0x0a, 0x6078, 0x6077, 0x05, 0x98, collectEnemyDrop, 0),
	"d8 boss": agesHC(0x05, 0x78),

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

	// this one's just a dummy; it'll always be shield
	"shop, 30 rupees": &MutableSlot{
		treasureName: "wooden shield",
		collectMode:  collectFind2,
	},
}
