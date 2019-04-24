package main

import (
	"container/list"
	"math/rand"
	"sort"

	"github.com/jangler/oracles-randomizer/graph"
	"github.com/jangler/oracles-randomizer/logic"
	"github.com/jangler/oracles-randomizer/rom"
)

// returns true iff the node is in the list.
func nodeInList(n *graph.Node, l *list.List) bool {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*graph.Node) == n {
			return true
		}
	}
	return false
}

func trySlotRandomItem(r *Route, src *rand.Rand, itemPool,
	slotPool *list.List, countFunc func(*Route, bool) int, numUsedSlots int,
	hard, fillUnused bool) (usedItem, usedSlot *list.Element) {
	// we're dead
	if slotPool.Len() == 0 || itemPool.Len() == 0 {
		return nil, nil
	}

	// this is the last slot, so it has to open up progression
	var initialCount int
	if slotPool.Len() == numUsedSlots+1 && !fillUnused {
		initialCount = countFunc(r, hard)
	}

	// try placing an item in the first slot until one fits
	for es := slotPool.Front(); es != nil; es = es.Next() {
		slot := es.Value.(*graph.Node)

		r.Graph.ClearMarks()
		if slot.GetMark(slot, hard) != graph.MarkTrue ||
			!canAffordSlot(r, slot, hard) {
			continue
		}

		for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*graph.Node)

			if !itemFitsInSlot(item, slot, src) {
				continue
			}

			item.AddParents(slot)

			if slotPool.Len() == numUsedSlots+1 && !fillUnused {
				newCount := countFunc(r, hard)
				if newCount <= initialCount {
					item.RemoveParent(slot)
					continue
				}
			}

			return ei, es
		}
	}

	return nil, nil
}

// maps should be looped through based on a sorted set of keys (which can be
// reordered before iteration, as long as it's ordered first); otherwise the
// same random seed can yield different results.
func getSortedKeys(g graph.Graph, src *rand.Rand) []string {
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
func itemFitsInSlot(itemNode, slotNode *graph.Node, src *rand.Rand) bool {
	// dummy shop slots 1 and 2 can only hold their vanilla items.
	if slotNode.Name == "shop, 20 rupees" && itemNode.Name != "bombs, 10" {
		return false
	}
	if slotNode.Name == "shop, 30 rupees" && itemNode.Name != "wooden shield" {
		return false
	}
	if itemNode.Name == "wooden shield" && slotNode.Name != "shop, 30 rupees" {
		return false
	}

	// rod of seasons has special graphics something
	if slotNode.Name == "temple of seasons" &&
		!rom.CanSlotAsRod(itemNode.Name) {
		return false
	}

	// bomb flower has special graphics something
	if itemNode.Name == "bomb flower" {
		switch slotNode.Name {
		case "cheval's test", "cheval's invention", "wild tokay game",
			"hidden tokay cave", "library present", "library past":
			return false
		}
	}

	// and only seeds can be slotted in seed trees, of course
	switch itemNode.Name {
	case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
		"pegasus tree seeds", "gale tree seeds":
		return slotIsSeedTree(slotNode.Name)
	default:
		return !slotIsSeedTree(slotNode.Name)
	}

	return true
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

func canAffordSlot(r *Route, slot *graph.Node, hard bool) bool {
	// if it doesn't cost anything, of course it's affordable
	balance := logic.NodeValues[slot.Name]
	if balance >= 0 {
		return true
	}

	// in hard mode, 100 rupee manips with shovel are in logic
	if hard {
		shovel := r.Graph["shovel"]
		if shovel.GetMark(shovel, hard) == graph.MarkTrue {
			return true
		}
	}

	// otherwise, count the net rupees available to the player
	balance += r.Rupees
	for _, node := range r.Graph {
		value := logic.NodeValues[node.Name]
		if value != 0 && node != slot &&
			node.GetMark(node, hard) == graph.MarkTrue {
			balance += value
		}
	}

	return balance >= 0
}
