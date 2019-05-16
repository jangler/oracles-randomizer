package main

import (
	"container/list"
	"math/rand"
	"sort"
)

// returns true iff the node is in the list.
func nodeInList(n *node, l *list.List) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*node) == n {
			return true
		}
	}
	return false
}

func trySlotRandomItem(r *Route, src *rand.Rand,
	itemPool, slotPool *list.List) (usedItem, usedSlot *list.Element) {
	// we're dead
	if slotPool.Len() == 0 || itemPool.Len() == 0 {
		return nil, nil
	}

	// try placing the first item in a slot until it fits
	for _, progressionItemsOnly := range []bool{true, false} {
		for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)
			if progressionItemsOnly && itemIsJunk(item.name) {
				continue
			}
			item.removeParent(r.Graph["start"])

			for es := slotPool.Front(); es != nil; es = es.Next() {
				slot := es.Value.(*node)

				if !itemFitsInSlot(item, slot, src) {
					continue
				}

				r.Graph.clearMarks()
				item.addParent(slot)
				if r.Graph["done"].getMark() != markTrue {
					item.removeParent(slot)
					continue
				}

				return ei, es
			}

			item.addParent(r.Graph["start"])
		}
	}

	return nil, nil
}

// maps should be looped through based on a sorted set of keys (which can be
// reordered before iteration, as long as it's ordered first); otherwise the
// same random seed can yield different results.
func getSortedKeys(g graph, src *rand.Rand) []string {
	keys := make([]string, 0, len(g))
	for k := range g {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

// checks whether the item fits in the slot due to things like seeds only going
// in trees, certain item slots not accomodating sub IDs. this doesn't check
// for softlocks or the availability of the slot and item.
func itemFitsInSlot(itemNode, slotNode *node, src *rand.Rand) bool {
	// dummy shop slots 1 and 2 can only hold their vanilla items.
	if slotNode.name == "shop, 20 rupees" && itemNode.name != "bombs, 10" {
		return false
	}
	if slotNode.name == "shop, 30 rupees" && itemNode.name != "wooden shield" {
		return false
	}
	if itemNode.name == "wooden shield" && slotNode.name != "shop, 30 rupees" {
		return false
	}

	// bomb flower has special graphics something
	// TODO: maybe this can be worked around like with the temple of seasons
	// item in seasons. not sure if it's super worth it but it'd be good to be
	// consistent.
	if itemNode.name == "bomb flower" {
		switch slotNode.name {
		case "cheval's test", "cheval's invention", "wild tokay game",
			"hidden tokay cave", "library present", "library past":
			return false
		}
	}

	// and only seeds can be slotted in seed trees, of course
	switch itemNode.name {
	case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
		"pegasus tree seeds", "gale tree seeds":
		return slotIsSeedTree(slotNode.name)
	default:
		return !slotIsSeedTree(slotNode.name)
	}
}

func slotIsSeedTree(name string) bool {
	switch name {
	case "horon village seed tree", "woods of winter seed tree",
		"north horon seed tree", "spool swamp seed tree",
		"sunken city seed tree", "tarm ruins seed tree", "south lynna tree",
		"deku forest tree", "crescent island tree", "symmetry city tree",
		"rolling ridge west tree", "rolling ridge east tree",
		"ambi's palace tree", "zora village tree":
		return true
	}
	return false
}
