package main

import (
	"container/list"
	"log"
	"math/rand"
	"sort"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

// attempts to reach new steps from the given graph state by slotting available
// items in available slots. it returns a list of slotted items if it succeeds,
// or nil if it fails.
func trySlotItemSet(r *Route, src *rand.Rand) *list.List {
	slots := getAvailableSlots(r, src)
	items := getAvailableItems(r, src)
	slottedItems := list.New()
	initialCount := countSteps(r.Graph.ExploreFromStart())
	newCount := 0

	if slots.Len() == 0 || items.Len() == 0 {
		return nil
	}

	// try placing each item in each slot, until no more slots are available
	n := items.Len()
	for i := 0; i < n; i++ {
		usedSlots := list.New()

		for slots.Len() > 0 && newCount <= initialCount {
			slot := slots.Remove(slots.Front()).(*graph.Node)
			usedSlots.PushBack(slot)
			for items.Len() > 0 && newCount <= initialCount {
				item := items.Remove(items.Front()).(*graph.Node)
				item.AddParents(slot)
				slottedItems.PushBack(item)

				if canSoftlock(r.HardGraph) != nil {
					item.ClearParents()
					slottedItems.Remove(slottedItems.Back())
					items.PushBack(item)
					continue
				}

				newCount = countSteps(r.Graph.ExploreFromStart())
				break
			}
		}

		if newCount > initialCount {
			break
		} else {
			for usedSlots.Len() > 0 {
				slots.PushBack(usedSlots.Remove(usedSlots.Front()))
			}
			for slottedItems.Len() > 0 {
				item := slottedItems.Remove(slottedItems.Front()).(*graph.Node)
				item.ClearParents()
				items.PushBack(item)
			}
		}
	}

	log.Printf("initial count %d, new count %d", initialCount, newCount)

	if newCount <= initialCount {
		return nil
	}

	// try removing each item from each slot to see if the path can still be
	// reached without them
	retry := true
	for retry {
		retry = false

		for i := 0; i < slottedItems.Len(); i++ {
			item := slottedItems.Remove(slottedItems.Front()).(*graph.Node)
			parents := item.Parents
			item.ClearParents()

			testCount := countSteps(r.Graph.ExploreFromStart())

			if testCount == newCount && canSoftlock(r.HardGraph) == nil {
				retry = true
				break
			} else {
				item.AddParents(parents...)
				slottedItems.PushBack(item)
			}
		}
	}

	log.Printf("slotted %d item(s):", slottedItems.Len())
	for e := slottedItems.Front(); e != nil; e = e.Next() {
		log.Printf("- %s", e.Value.(*graph.Node).Name)
	}

	if newCount > initialCount {
		return slottedItems
	}
	return nil
}

// get item slots that are reachable and unused, sorted by slot priority.
func getAvailableSlots(r *Route, src *rand.Rand) *list.List {
	names := make([]string, 0, len(r.Graph))
	for name := range r.Graph {
		names = append(names, name)
	}
	sort.Strings(names)

	slots := make([]*graph.Node, 0)
	for _, name := range names {
		node := r.Graph[name]
		if node.IsSlot && len(node.Children) == 0 &&
			node.GetMark(node, nil) == graph.MarkTrue {
			slots = append(slots, node)
		}
	}

	src.Shuffle(len(slots), func(i, j int) {
		slots[i], slots[j] = slots[j], slots[i]
	})

	sort.Slice(slots, func(i, j int) bool {
		// dungeon chests go first
		di := dungeonIndex(slots[i])
		if di >= 0 && r.Dungeons[di].ItemsPlaced == 0 {
			return true
		}

		// special item slots go second
		slot := rom.ItemSlots[slots[j].Name]
		return rom.IsChest(slot) || rom.IsFound(slot)
	})

	l := list.New()
	refillList(l, slots)
	return l
}

// get unused item nodes, sorted by placement priority.
func getAvailableItems(r *Route, src *rand.Rand) *list.List {
	names := make([]string, 0, len(r.Graph))
	for name := range r.Graph {
		names = append(names, name)
	}
	sort.Strings(names)

	items := make([]*graph.Node, 0)
	for _, name := range names {
		node := r.Graph[name]
		if node.Type == graph.RootType && len(node.Parents) == 0 &&
			!strings.Contains(node.Name, "default") &&
			node.Name != "compass" && node.Name != "dungeon map" {
			items = append(items, node)
		}
	}

	src.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})

	sort.Slice(items, func(i, j int) bool {
		return keyItems[items[i].Name] ||
			strings.HasPrefix(items[i].Name, "rupees") ||
			items[i].Name == "member's card" || items[i].Name == "red ore" ||
			items[i].Name == "blue ore"
	})

	l := list.New()
	refillList(l, items)
	return l
}
