package rom

// agesChest constructs a MutableSlot from a treasure name and an address in
// bank $16, where the ID and sub-ID are two consecutive bytes at that address.
// This applies to almost all chests, and exclusively to chests.
func agesChest(treasure string, addr uint16,
	group, room, coords byte) *MutableSlot {
	return BasicSlot(treasure, 0x16, addr, addr+1,
		group, room, collectChest, coords)
}

var agesSlots = map[string]*MutableSlot{
	// overworld present
	/*
		"impa's gift": CustomSlot("sword 1",
			0x00, 0x39, collectFind2, 0x39),
		"nayru's house": CustomSlot("tune of echoes",
			0x03, 0xae, collectFind2, 0x3a),
	*/
	"lynna city chest": agesChest("rupees, 30", 0x511e,
		0x00, 0x49, 0x49),
	"fairies' woods chest": agesChest("rupees, 50", 0x5122,
		0x00, 0x84, 0x84),
	"fairies' coast chest": agesChest("green holy ring", 0x5126,
		0x00, 0x91, 0x91),
	"zora seas chest": agesChest("whimsical ring", 0x512e,
		0x00, 0xd5, 0xd5),
	"talus peaks chest": agesChest("gasha seed", 0x5132,
		0x00, 0x63, 0x63),
	"ruined keep": agesChest("armor ring L-1", 0x5144,
		0x02, 0xbe, 0x09),
	"nuun cave": agesChest("light ring L-1", 0x5154,
		0x02, 0xf4, 0x37),
	"zora village present": agesChest("gasha seed", 0x515c,
		0x02, 0xc0, 0xc0),
	"D6 entrance pool": agesChest("toss ring", 0x5161,
		0x03, 0x0e, 0x3c),
	// "sea of storms cave present": nil, // special linked chest?
	"mayor plen's house": agesChest("green luck ring", 0x5171,
		0x03, 0xf9, 0x57),
	"cresent island underwater cave": agesChest("piece of heart", 0x5175,
		0x03, 0xfd, 0xda),

	// overworld past
	"sea of no return": agesChest("blue ring", 0x5137,
		0x01, 0x6d, 0x6d),
	"bomb goron head": agesChest("rupees, 100", 0x5148,
		0x02, 0xfc, 0x0d),
	"tokay cave past": agesChest("gasha seed", 0x514c,
		0x02, 0xce, 0xcd),
	"zora cave past": agesChest("red holy ring", 0x5158,
		0x02, 0x4f, 0xc5),
	"ridge bush cave": agesChest("rupees, 100", 0x5165,
		0x03, 0x1f, 0x1c),
	// "sea of storms cave past": nil, // special linked chest?

	// dungeons TODO
}
