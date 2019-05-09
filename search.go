package main

import (
	"container/list"
	"math/rand"
	"sort"

	"github.com/jangler/oracles-randomizer/logic"
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

func trySlotRandomItem(r *Route, src *rand.Rand, itemPool,
	slotPool *list.List, numUsedSlots int,
	hard, fillUnused bool) (usedItem, usedSlot *list.Element) {
	// we're dead
	if slotPool.Len() == 0 || itemPool.Len() == 0 {
		return nil, nil
	}

	// try placing an item in the first slot until one fits
	for es := slotPool.Front(); es != nil; es = es.Next() {
		slot := es.Value.(*node)

		r.Graph.clearMarks()
		if !fillUnused && (slot.getMark() != markTrue ||
			!canAffordSlot(r, slot, hard)) {
			continue
		}

		for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)

			if !itemFitsInSlot(item, slot, src) {
				continue
			}

			item.addParent(slot)

			return ei, es
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

func canAffordSlot(r *Route, slot *node, hard bool) bool {
	// if it doesn't cost anything, of course it's affordable
	balance := logic.NodeValues[slot.name]
	if balance >= 0 {
		return true
	}

	// in hard mode, 100 rupee manips with shovel are in logic
	if hard {
		if r.Graph["shovel"].getMark() == markTrue {
			return true
		}
	}

	// otherwise, count the net rupees available to the player
	balance += r.Rupees
	for _, n := range r.Graph {
		value := logic.NodeValues[n.name]
		if value != 0 && n != slot && n.getMark() == markTrue {
			balance += value
		}
	}

	return balance >= 0
}
