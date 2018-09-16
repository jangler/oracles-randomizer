package main

import (
	"container/list"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/logic"
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
	countFunc func(*Route) int,
	fillUnused bool) (usedItems, usedSlots *list.List) {

	startTime := time.Now()

	// at least one item must fit in a slot that has already been reached for
	// no more than this many steps
	maxStaleness := int(math.Abs(src.NormFloat64()))

	// get a list of slots that are actually reachable; see what can be reached
	// before slotting anything more
	freeSlots := getAvailableSlots(r, src, slotPool, maxStaleness, fillUnused)
	initialCount := countFunc(r)
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
					continue
				}
				if !itemFitsInSlot(item, slot, src) {
					continue
				}

				item.AddParents(slot)

				if canSoftlock(r.Graph) != nil {
					item.RemoveParent(slot)
				} else {
					newCount = countFunc(r)
					usedItems.PushBack(item)
					usedSlots.PushBack(slot)
					break
				}
			}
		}

		if newCount == initialCount && !fillUnused {
			for usedItems.Len() > 0 {
				item := usedItems.Remove(usedItems.Front()).(*graph.Node)
				slot := usedSlots.Remove(usedSlots.Front()).(*graph.Node)
				item.RemoveParent(slot)
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
	usedUnstaleTurnSlot := false
	for ei := usedItems.Front(); ei != nil; ei = ei.Next() {
		item := ei.Value.(*graph.Node)
		item.PopParent()

		for es := freeSlots.Front(); es != nil; es = es.Next() {
			slot := es.Value.(*graph.Node)
			if nodeInList(slot, usedSlots) {
				continue
			}

			if itemFitsInSlot(item, slot, nil) {
				item.AddParents(slot)

				if canSoftlock(r.Graph) != nil {
					item.RemoveParent(slot)
				} else {
					usedSlots.PushBack(slot)
					if r.TurnsReached[slot] <= maxStaleness {
						usedUnstaleTurnSlot = true
					}
					break
				}
			}
		}
	}

	// retry if none of the slotted items were placed in new slots (because
	// they couldn't fit). shops and trees are exempt from this check.
	if !fillUnused && !usedUnstaleTurnSlot {
		for usedItems.Len() > 0 {
			item := usedItems.Remove(usedItems.Front()).(*graph.Node)
			item.PopParent()
		}
		return list.New(), list.New()
	}

	// increment staleness, capping turn-neutral items at staleness 2.
	for e := freeSlots.Front(); e != nil; e = e.Next() {
		slot := e.Value.(*graph.Node)
		if r.TurnsReached[slot] < 2 || !isTurnNeutralItem(slot) {
			r.TurnsReached[slot]++
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

// some items shouldn't be prioritized based on number of turns reached,
// because other requirements have to be met in order for them to be used.
func isTurnNeutralItem(node *graph.Node) bool {
	// like shop items
	if logic.Rupees[node.Name] != 0 {
		return true
	}

	// and seed trees
	switch node.Name {
	case "ember tree", "mystery tree", "scent tree", "pegasus tree",
		"sunken gale tree", "tarm gale tree":
		return true
	}

	return false
}

// filter a list of item slots by those that can be reached, shuffle them, and
// sort them by priority, returning a new list.
func getAvailableSlots(r *Route, src *rand.Rand, pool *list.List,
	maxStaleness int, fillUnused bool) *list.List {
	a := make([]*graph.Node, 0)
	r.Graph.ClearMarks()
	for e := pool.Front(); e != nil; e = e.Next() {
		node := e.Value.(*graph.Node)
		if node.GetMark(node, false) == graph.MarkTrue &&
			canAffordSlot(r, node) {
			a = append(a, node)
		}
	}

	src.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	// prioritize newer slots
	sort.Slice(a, func(i, j int) bool {
		return r.TurnsReached[a[i]] <= maxStaleness &&
			r.TurnsReached[a[j]] > maxStaleness
	})

	// if filling unused, return only one slot at a time
	if fillUnused {
		l := list.New()

		// first prioritize especially restrictive slots
		for _, node := range a {
			switch node.Name {
			case "d0 sword chest", "rod gift", "star ore spot", "hard ore slot",
				"diver gift", "subrosian market 5", "village shop 1",
				"village shop 2", "village shop 3", "iron shield gift":
				l.PushBack(node)
				return l
			}
		}

		// then prioritize non-chests
		for _, node := range a {
			if !rom.IsChest(rom.ItemSlots[node.Name]) {
				l.PushBack(node)
				return l
			}
		}

		if len(a) > 0 {
			l.PushBack(a[0])
		}
		return l
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
			slotNode.Name == "iron shield gift" ||
			slotNode.Name == "hard ore slot" ||
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
			case "feather L-1", "feather L-2":
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

	switch slotNode.Name {
	// hard ore is a special case because it doesn't set sub ID.
	case "hard ore slot":
		if item.SubID() != 0 && !(itemNode.Name == "piece of heart" ||
			itemNode.Name == "gasha seed") {
			return false
		}
	// these slots won't give you the item if you already have one with that
	// ID, so only use items that have unique IDs and can't be lost.
	case "diver gift", "subrosian market 5", "village shop 3":
		if !rom.TreasureHasUniqueID(itemNode.Name) ||
			rom.TreasureCanBeLost(itemNode.Name) {
			return false
		}
	// star ore is the above two cases combined.
	case "star ore spot":
		if item.SubID() != 0 || !rom.TreasureHasUniqueID(itemNode.Name) ||
			rom.TreasureCanBeLost(itemNode.Name) {
			return false
		}
	// this slot apparently checks ID too (and will give you a fake rusty bell
	// if you already have the ID it checks), but loseable items are ok since
	// you trade hard ore for it.
	case "iron shield gift":
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
	countFunc func(*Route) int, fillUnused bool) {
	// try removing items
	retry := true
	for retry && !fillUnused {
		retry = false

		for e := usedItems.Front(); e != nil; e = e.Next() {
			item := e.Value.(*graph.Node)
			parent := item.PopParent()

			testCount := countFunc(r)
			if testCount > initialCount && canSoftlock(r.Graph) == nil {
				// remove the item and cycle again if it can be omitted
				retry = true
				usedItems.Remove(e)
				break
			}

			item.AddParents(parent)
		}
	}

	// try downgrading items, as long as that doesn't affect the number of new
	// slots reached
	targetCount := countFunc(r)
	retry = true
	triedDowngradeToSatchel := false
	for retry && !fillUnused {
		retry = false

		for e := usedItems.Front(); e != nil; e = e.Next() {
			item := e.Value.(*graph.Node)
			var downgrade *graph.Node

			// don't use a slingshot where a satchel will do
			if !triedDowngradeToSatchel &&
				strings.HasPrefix(item.Name, "slingshot") &&
				r.Graph["satchel 1"].NumParents() == 0 &&
				r.Graph["satchel 2"].NumParents() == 0 {
				downgrade = r.Graph["satchel 1"]
			}

			// and don't use a L-2 item where a L-1 item will do
			if downgrade == nil {
				if !strings.HasSuffix(item.Name, "L-2") {
					continue
				}
				if item.Name == "shield L-2" || item.Name == "armor ring L-2" {
					continue // no L-1 equivalents in pool
				}
				downgrade =
					r.Graph[strings.Replace(item.Name, "L-2", "L-1", 1)]
			}

			if downgrade == nil || downgrade.NumParents() > 0 {
				continue
			}

			parent := item.PopParent()
			downgrade.AddParents(parent)

			testCount := countFunc(r)
			if testCount >= targetCount && canSoftlock(r.Graph) == nil {
				// downgrade item and cycle again
				retry = true
				usedItems.InsertAfter(downgrade, e)
				usedItems.Remove(e)
				break
			}

			downgrade.RemoveParent(parent)
			item.AddParents(parent)

			// if L-2 slingshot downgrade to satchel didn't work, still try to
			// downgrade it to L-1 slingshot.
			if !triedDowngradeToSatchel && item.Name == "slingshot L-2" &&
				strings.HasPrefix(downgrade.Name, "satchel") {
				triedDowngradeToSatchel, retry = true, true
				break
			}
		}
	}
}

func canAffordSlot(r *Route, slot *graph.Node) bool {
	// if it doesn't cost anything, of course it's affordable
	balance := logic.Rupees[slot.Name]
	if balance >= 0 {
		return true
	}

	// otherwise, count the net rupees available to the player
	balance += r.Costs
	for _, node := range r.Graph {
		value := logic.Rupees[node.Name]
		if value > 0 && node.GetMark(node, false) == graph.MarkTrue {
			balance += value
		}
	}

	return balance >= 0
}
