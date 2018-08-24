package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
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

// attempts to reach new steps from the given graph state by slotting available
// items in available slots. it returns a list of slotted items if it succeeds,
// or nil if it fails.
func trySlotItemSet(r *Route, src *rand.Rand, itemPool, slotPool *list.List,
	countFunc func(map[*graph.Node]bool) int,
	fillUnused bool) (usedItems, usedSlots *list.List) {

	// get a list of slots that are actually reachable; see what can be reached
	// before slotting anything more
	freeSlots := getAvailableSlots(r, src, slotPool)
	initialCount := countFunc(r.Graph.ExploreFromStart())
	newCount := initialCount

	if freeSlots.Len() == 0 || itemPool.Len() == 0 {
		return nil, nil
	}

	// try placing each item in each slot, until no more slots are available.
	usedItems = list.New()
	usedSlots = list.New()
	for i := 0; i < itemPool.Len() && newCount == initialCount; i++ {
		for e := freeSlots.Front(); e != nil &&
			newCount == initialCount; e = e.Next() {
			slot := e.Value.(*graph.Node)
			if nodeInList(slot, usedSlots) {
				continue
			}

			for e := itemPool.Front(); e != nil; e = e.Next() {
				item := e.Value.(*graph.Node)
				if nodeInList(item, usedItems) {
					// XXX this is not really accurate since a gasha seed could
					//     be slotted twice in one iteration
					continue
				}
				if !itemFitsInSlot(item, slot, src) {
					continue
				}

				item.Parents = append(item.Parents, slot)

				if canSoftlock(r.HardGraph) != nil {
					item.Parents = item.Parents[:len(item.Parents)-1]
				} else {
					usedItems.PushBack(item)
					usedSlots.PushBack(slot)
					break
				}
			}

			newCount = countFunc(r.Graph.ExploreFromStart())
		}

		if newCount == initialCount && !fillUnused {
			for usedItems.Len() > 0 {
				item := usedItems.Remove(usedItems.Front()).(*graph.Node)
				slot := usedSlots.Remove(usedSlots.Front()).(*graph.Node)
				removeNodeFromSlice(slot, &item.Parents)
			}
		}
		itemPool.PushBack(itemPool.Remove(itemPool.Front()))
	}

	// couldn't find any progression; fail
	if newCount == initialCount && !fillUnused {
		return nil, nil
	}

	// omit items not necessary for progression, then slot again from the start
	cutExtraItems(r, usedItems, initialCount, countFunc, fillUnused)
	usedSlots.Init()
	for ei := usedItems.Front(); ei != nil; ei = ei.Next() {
		item := ei.Value.(*graph.Node)
		item.Parents = item.Parents[:len(item.Parents)-1]

		for es := freeSlots.Front(); es != nil; es = es.Next() {
			slot := es.Value.(*graph.Node)
			if nodeInList(slot, usedSlots) {
				continue
			}

			if itemFitsInSlot(item, slot, nil) {
				usedSlots.PushBack(slot)
				item.Parents = append(item.Parents, slot)
				break
			}
		}
	}

	// remove the used nodes from the persistent pools
	if newCount > initialCount || (fillUnused && usedItems.Len() > 0) {
		for e := usedItems.Front(); e != nil; e = e.Next() {
			removeNodeFromList(e.Value.(*graph.Node), itemPool)
		}
		for e := usedSlots.Front(); e != nil; e = e.Next() {
			removeNodeFromList(e.Value.(*graph.Node), slotPool)
		}
		return usedItems, usedSlots
	}
	return nil, nil
}

func removeNodeFromList(n *graph.Node, l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value.(*graph.Node) == n {
			l.Remove(e)
			return
		}
	}
	panic(fmt.Sprintf("node %v not in list", n))
}

func removeNodeFromSlice(n *graph.Node, a *[]*graph.Node) {
	for i, v := range *a {
		if v == n {
			*a = append((*a)[:i], (*a)[i+1:]...)
			return
		}
	}
	panic(fmt.Sprintf("node %v not in slice", n))
}

var dungeonRegexp = regexp.MustCompile(`^d(\d) `)

// filter a list of item slots by those that can be reached, shuffle them, and
// sort them by priority, returning a new list.
func getAvailableSlots(r *Route, src *rand.Rand, pool *list.List) *list.List {
	a := make([]*graph.Node, 0)
	for e := pool.Front(); e != nil; e = e.Next() {
		node := e.Value.(*graph.Node)
		if node.GetMark(node, nil) == graph.MarkTrue {
			a = append(a, node)
		}
	}

	src.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	// prioritize slots from higher dungeon numbers, and prioritize anything
	// over slots that were already reached in a previous iteration
	sort.Slice(a, func(i, j int) bool {
		if !r.OldSlots[a[i]] && r.OldSlots[a[j]] {
			return true
		}

		iMatch := dungeonRegexp.FindStringSubmatch(a[i].Name)
		jMatch := dungeonRegexp.FindStringSubmatch(a[j].Name)

		if iMatch != nil && jMatch != nil {
			di, _ := strconv.Atoi(iMatch[1])
			dj, _ := strconv.Atoi(jMatch[1])
			return di > dj
		}

		return false
	})

	for _, slot := range a {
		r.OldSlots[slot] = true
	}

	return listFromSlice(a)
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
	slot := rom.ItemSlots[slotNode.Name]
	item := rom.Treasures[itemNode.Name]

	// gasha seeds and pieces of heart can be placed in either chests or
	// found/gift slots. beyond that, only unique items can be placed in
	// non-chest slots.
	if itemNode.Name == "gasha seed" || itemNode.Name == "piece of heart" {
		if slotNode.Name == "d0 sword chest" || slotNode.Name == "rod gift" ||
			!(rom.IsChest(slot) || rom.IsFound(slot)) {
			return false
		}
	} else if (!rom.IsChest(slot) ||
		slotNode.Name == "d0 sword chest" || slotNode.Name == "rod gift") &&
		!rom.TreasureIsUnique[itemNode.Name] {
		return false
	}

	// don't put gale seeds in the ember tree, since then gale seeds will come
	// with the satchel and the player can freeze the game by trying to warp
	// without having explored any trees.
	if slotNode.Name == "ember tree" &&
		strings.HasPrefix(itemNode.Name, "gale tree seeds") {
		return false
	}

	// give only a 1 in 2 change per sword of slotting in the hero's cave chest
	// to compensate for the fact that there are two of them. each season gets
	// a 1 in 4 chance for the same reason.
	if src != nil {
		if slotNode.Name == "d0 sword chest" {
			switch itemNode.Name {
			case "sword L-1", "sword L-2":
				if src.Intn(2) != 0 {
					return false
				}
			case "winter", "spring", "summer", "autumn":
				if src.Intn(4) != 0 {
					return false
				}
			}
		}
	}

	// star ore and hard ore are special cases because they doesn't set sub ID
	// at all, so only slot zero-ID treasures there.
	//
	// the other slots won't give you the item if you already have one with
	// that ID, so only use items with unique IDs there.
	switch slotNode.Name {
	case "star ore spot", "hard ore slot":
		if item.SubID() != 0 && !(itemNode.Name == "piece of heart" ||
			itemNode.Name == "gasha seed") {
			return false
		}
	case "diver gift", "subrosian market 5":
		if !rom.TreasureHasUniqueID(itemNode.Name) {
			return false
		}
	}

	// some items can't be drawn correctly in certain item slots.
	switch slotNode.Name {
	case "d0 sword chest", "rod gift", "noble sword spot":
		if !rom.CanSlotInScene(itemNode.Name) {
			return false
		}
	case "member's shop 1", "member's shop 2", "member's shop 3":
		if !rom.CanSlotInShop(itemNode.Name) {
			return false
		}
	case "subrosian market 2", "subrosian market 5":
		if !rom.CanSlotInMarket(itemNode.Name) {
			return false
		}
	}

	// and only seeds can be slotted in seed trees, of course
	switch itemNode.Name {
	case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
		"pegasus tree seeds", "gale tree seeds 1", "gale tree seeds 2":
		switch slotNode.Name {
		case "ember tree", "mystery tree", "scent tree",
			"pegasus tree", "sunken gale tree", "tarm gale tree":
			break
		default:
			return false
		}
	default:
		switch slotNode.Name {
		case "ember tree", "mystery tree", "scent tree",
			"pegasus tree", "sunken gale tree", "tarm gale tree":
			return false
		}
	}

	return true
}

// try removing each item from each slot to see if progression can still be
// reached without it
func cutExtraItems(r *Route, usedItems *list.List, initialCount int,
	countFunc func(map[*graph.Node]bool) int, fillUnused bool) {
	retry := true
	for retry && !fillUnused {
		retry = false

		for e := usedItems.Front(); e != nil; e = e.Next() {
			item := e.Value.(*graph.Node)
			parent := item.Parents[len(item.Parents)-1]
			item.Parents = item.Parents[:len(item.Parents)-1]

			testCount := countFunc(r.Graph.ExploreFromStart())
			if testCount > initialCount && canSoftlock(r.HardGraph) == nil {
				// remove the item and cycle again if it can be omitted
				retry = true
				usedItems.Remove(e)
				break
			}

			item.Parents = append(item.Parents, parent)
		}
	}
}
