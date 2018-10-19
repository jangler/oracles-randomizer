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
}
