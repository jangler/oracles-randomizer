package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
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

	startTime := time.Now()

	// get a list of slots that are actually reachable; see what can be reached
	// before slotting anything more
	freeSlots := getAvailableSlots(r, src, slotPool, fillUnused)
	initialCount := countFunc(r.Graph.ExploreFromStart())
	newCount := initialCount

	sortItemPool(itemPool, src)

	if freeSlots.Len() == 0 || itemPool.Len() == 0 {
		return nil, nil
	}

	// try placing each item in each slot, until no more slots are available.
	usedItems = list.New()
	usedSlots = list.New()
	for i := 0; i < itemPool.Len() && newCount == initialCount; i++ {
		// check to make sure this step isn't taking too long
		if time.Now().Sub(startTime) > time.Second*10 {
			return nil, nil
		}

		for e := freeSlots.Front(); e != nil &&
			newCount == initialCount; e = e.Next() {
			slot := e.Value.(*graph.Node)
			if nodeInList(slot, usedSlots) {
				continue
			}

			for e := itemPool.Front(); e != nil; e = e.Next() {
				item := e.Value.(*graph.Node)
				if nodeInList(item, usedItems) {
					// break if filling unused because only one gasha seed can
					// be slotted per iteration
					if fillUnused {
						break
					}
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

			a := emptyList(usedItems)
			newCount = countFunc(r.Graph.ExploreFromStart())
			fillList(usedItems, a)

			// hack to make sure gasha seeds and such don't pile up at the end
			if fillUnused && len(a) > 0 &&
				rom.TreasureIsUnique[a[len(a)-1].Name] {
				break
			}
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
				item.Parents = append(item.Parents, slot)

				if canSoftlock(r.HardGraph) != nil {
					item.Parents = item.Parents[:len(item.Parents)-1]
				} else {
					usedSlots.PushBack(slot)
					break
				}
			}
		}
	}

	// abort if it's impossible to pay for the slotted items
	if !tryMeetCosts(r, usedItems, itemPool, usedSlots, freeSlots, src) {
		return nil, nil
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

// shuffle the item pool and arrange the items by priority.
func sortItemPool(pool *list.List, src *rand.Rand) {
	a := emptyList(pool)

	src.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	fillList(pool, a)

	// place all jewels together
	var jewelMark *list.Element
	for e := pool.Front(); e != nil; e = e.Next() {
		switch e.Value.(*graph.Node).Name {
		case "square jewel", "pyramid jewel", "round jewel",
			"x-shaped jewel":
			if jewelMark == nil {
				jewelMark = e
			} else {
				next := e.Next()
				pool.MoveAfter(e, jewelMark)
				if next != nil {
					e = next.Prev()
				}
			}
		}
	}
}

// filter a list of item slots by those that can be reached, shuffle them, and
// sort them by priority, returning a new list.
func getAvailableSlots(r *Route, src *rand.Rand, pool *list.List,
	fillUnused bool) *list.List {
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

	// prioritize, in order:
	// 1. if filling extra items, slots with specific restrictions
	// 2. dungeons with no items over anything else
	// 3. anything over slots that were already reached in a previous iteration
	sort.Slice(a, func(i, j int) bool {
		if fillUnused {
			switch a[i].Name {
			case "diver gift", "subrosian market 2", "subrosian market 5",
				"village shop 3", "d0 sword chest", "rod gift",
				"star ore spot", "hard ore slot":
				return true
			}
		}

		match := dungeonRegexp.FindStringSubmatch(a[i].Name)
		if match != nil {
			di, _ := strconv.Atoi(match[1])
			return r.DungeonItems[di] == 0
		}

		if !r.OldSlots[a[i]] && r.OldSlots[a[j]] &&
			prenode.Rupees[a[j].Name] == 0 {
			return true
		}

		return false
	})

	for _, slot := range a {
		r.OldSlots[slot] = true
	}

	l := list.New()
	fillList(l, a)

	// place the ember tree, if present, before any other tree
	var treeMark *list.Element
	for e := l.Front(); e != nil; e = e.Next() {
		node := e.Value.(*graph.Node)
		if treeMark == nil {
			switch node.Name {
			case "mystery tree", "scent tree", "pegasus tree",
				"sunken gale tree", "tarm gale tree":
				treeMark = e
			}
		} else if node.Name == "ember tree" {
			next := e.Next()
			l.MoveBefore(e, treeMark)
			if next != nil {
				e = next.Prev()
			}
		}
	}

	return l
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

	// dummy shop slots 1 and 2 can only hold their vanilla items.
	if slotNode.Name == "village shop 1" && itemNode.Name != "bombs, 10" {
		return false
	}
	if slotNode.Name == "village shop 2" && itemNode.Name != "shop shield L-1" {
		return false
	}
	if itemNode.Name == "shop shield L-1" && slotNode.Name != "village shop 2" {
		return false
	}

	// give proportionally reduced chances of roughly equivalent items
	// appearing in the d0 sword chest.
	if src != nil {
		if slotNode.Name == "d0 sword chest" {
			switch itemNode.Name {
			case "sword L-1", "sword L-2":
				if src.Intn(2) != 0 {
					return false
				}
			case "ricky's flute", "dimitri's flute", "moosh's flute":
				if src.Intn(3) != 0 {
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
	case "diver gift", "subrosian market 5", "village shop 3":
		if !rom.TreasureHasUniqueID(itemNode.Name) {
			return false
		}
	}

	// rod of seasons has special graphics something
	if slotNode.Name == "rod gift" && !rom.CanSlotAsRod(itemNode.Name) {
		return false
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

// try removing and downgrading slotted items until the minimal set neceessary
// for progression is reached.
func cutExtraItems(r *Route, usedItems *list.List, initialCount int,
	countFunc func(map[*graph.Node]bool) int, fillUnused bool) {
	// try removing items
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

	// try downgrading L-2 items
	retry = true
	for retry && !fillUnused {
		retry = false

		for e := usedItems.Front(); e != nil; e = e.Next() {
			item := e.Value.(*graph.Node)
			if !strings.HasSuffix(item.Name, "L-2") {
				continue
			}
			if item.Name == "shield L-2" {
				continue // no L-1 shield in pool
			}
			downgrade := r.Graph[strings.Replace(item.Name, "L-2", "L-1", 1)]
			if len(downgrade.Parents) > 0 {
				continue
			}

			parent := item.Parents[len(item.Parents)-1]
			item.Parents = item.Parents[:len(item.Parents)-1]
			downgrade.Parents = append(downgrade.Parents, parent)

			testCount := countFunc(r.Graph.ExploreFromStart())
			if testCount > initialCount && canSoftlock(r.HardGraph) == nil {
				// downgrade item and cycle again
				retry = true
				usedItems.InsertAfter(downgrade, e)
				usedItems.Remove(e)
				break
			}

			downgrade.Parents = downgrade.Parents[:len(downgrade.Parents)-1]
			item.Parents = append(item.Parents, parent)
		}
	}
}

func tryMeetCosts(r *Route, usedItems, itemPool, usedSlots,
	slotPool *list.List, src *rand.Rand) bool {
	// count the new costs of the used items
	costs := 0
	for e := usedSlots.Front(); e != nil; e = e.Next() {
		node := e.Value.(*graph.Node)
		costs -= prenode.Rupees[node.Name]
	}
	if costs <= 0 {
		return true
	}
	r.Costs += costs
	balance := -r.Costs

	// count the net rupees available to the player
	for node := range r.Graph.ExploreFromStart() {
		// don't subtract shops, since the player can see what they're buying
		if !strings.Contains(node.Name, "shop") {
			value := prenode.Rupees[node.Name]
			if !nodeInList(node, usedSlots) {
				balance += value
			}
		}
	}

	// if possible, add rupees until the player can afford the items
	for ei := itemPool.Front(); balance < 0 && ei != nil; ei = ei.Next() {
		item := ei.Value.(*graph.Node)
		value := prenode.Rupees[item.Name]
		if value < 10 {
			continue
		}
		if nodeInList(item, usedItems) {
			continue
		}

		for es := slotPool.Front(); es != nil; es = es.Next() {
			slot := es.Value.(*graph.Node)
			if !nodeInList(slot, usedSlots) && itemFitsInSlot(item, slot, src) {
				item.AddParents(slot)
				usedItems.PushFront(item)
				usedSlots.PushFront(slot)
				balance += value
				break
			}
		}
	}

	return balance >= 0
}
