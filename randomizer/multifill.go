package randomizer

import (
	"container/list"
	"math/rand"
	"sort"
	"strings"
)

// a route struct with additional information for multiworld shuffle
type multiRoute struct {
	ri     *routeInfo
	local  map[*node]bool
	checks map[*node]*node
}

// returns true iff a given slot/item can be multiworld shuffled
func isMultiEligible(mr *multiRoute, slot, item *node) bool {
	return !mr.local[slot] && !strings.Contains(item.name, " ring") &&
		!strings.HasSuffix(item.name, "small key") &&
		!strings.HasSuffix(item.name, "boss key") &&
		!strings.HasSuffix(item.name, "dungeon map") &&
		!strings.HasSuffix(item.name, "compass") &&
		item.name != "slate"
}

// picks a random multiworld-eligible check from a route
func randomMultiCheck(src *rand.Rand, mr *multiRoute) (*node, *node) {
	i, slots := 0, make([]*node, 0, len(mr.checks))
	for slot, item := range mr.checks {
		if isMultiEligible(mr, slot, item) {
			slots = append(slots, slot)
			i++
		}
	}

	sort.Slice(slots, func(i, j int) bool {
		return slots[i].name < slots[j].name
	})

	slot := slots[src.Intn(len(slots))]
	return slot, mr.checks[slot]
}

func shuffleMultiworld(
	ris []*routeInfo, roms []*romState, verbose bool, logf logFunc) {
	mrs := make([]*multiRoute, len(ris))
	src := ris[0].src
	swapCounts := make(map[*node]int)
	swaps := 0

	for i, ri := range ris {
		// mark graph nodes as belonging to each player
		for _, v := range ri.graph {
			v.player = i + 1
		}

		// accumulate checks in a sane way
		mrs[i] = &multiRoute{
			ri:     ri,
			checks: getChecks(ri.usedItems, ri.usedSlots),
			local:  make(map[*node]bool, ri.usedItems.Len()),
		}
		for slot, item := range mrs[i].checks {
			mrs[i].local[slot] = roms[i].itemSlots[slot.name].localOnly
			if isMultiEligible(mrs[i], slot, item) {
				swapCounts[slot] = 0
			}
		}
	}

	// swap some random items ???
	consecutiveMisses := 0
	for consecutiveMisses < 1000 {
		slot1, item1 := randomMultiCheck(src, mrs[src.Intn(len(mrs))])
		slot2, item2 := randomMultiCheck(src, mrs[src.Intn(len(mrs))])

		// skip if slots are from the same player, or either of the treasures
		// aren't present in the other game, or the swapped treasures don't fit
		// in their new slots
		if slot1.player == slot2.player ||
			roms[slot1.player-1].treasures[item2.name] == nil ||
			roms[slot2.player-1].treasures[item1.name] == nil ||
			!itemFitsInSlot(item2, slot1) ||
			!itemFitsInSlot(item1, slot2) {
			continue
		}

		// only bother making swaps where at least one slot has yet to swap.
		// this is also used as a metric to determine when we've probably
		// reached about the maximum number of possible swaps
		if swapCounts[slot1] != 0 && swapCounts[slot2] != 0 {
			consecutiveMisses++
			continue
		}
		consecutiveMisses = 0

		if verbose {
			logf("swapping {p%d %s <- p%d %s} with {p%d %s <- p%d %s}",
				slot1.player, slot1.name, item1.player, item1.name,
				slot2.player, slot2.name, item2.player, item2.name)
		}

		// swap parents
		item1.removeParent(slot1)
		item2.removeParent(slot2)
		item1.addParent(slot2)
		item2.addParent(slot1)

		// test whether seeds are still beatable w/ item placement
		success := true
		for _, ri := range ris {
			ri.graph.reset()
		}
		for _, ri := range ris {
			ri.graph["start"].explore()
		}
		for _, ri := range ris {
			if !ri.graph["done"].reached {
				success = false
				break
			}
		}

		// make sure no player has to wait too long on progression from another
		mrs[slot1.player-1].checks[slot1] = item2
		mrs[slot2.player-1].checks[slot2] = item1
		for i, ri := range ris {
			ri.usedItems, ri.usedSlots = list.New(), list.New()
			for slot, item := range mrs[i].checks {
				ri.usedItems.PushBack(item)
				ri.usedSlots.PushBack(slot)
			}
		}
		mrs[slot1.player-1].checks[slot1] = item1
		mrs[slot2.player-1].checks[slot2] = item2
		if playerHasConsecutiveEmptySpheres(ris, 2) {
			success = false
		}

		// update check maps
		if success {
			mrs[slot1.player-1].checks[slot1] = item2
			mrs[slot2.player-1].checks[slot2] = item1
			swaps++
			swapCounts[slot1]++
			swapCounts[slot2]++
		} else {
			item1.removeParent(slot2)
			item2.removeParent(slot1)
			item1.addParent(slot1)
			item2.addParent(slot2)
			if verbose {
				logf("route no longer viable")
			}
		}
	}

	if verbose {
		logf("made %d successful swaps", swaps)
	}

	// reconstruct used item and slot lists
	for i, ri := range ris {
		ri.usedItems, ri.usedSlots = list.New(), list.New()
		for slot, item := range mrs[i].checks {
			ri.usedItems.PushBack(item)
			ri.usedSlots.PushBack(slot)
			roms[i].itemSlots[slot.name].player = byte(item.player)
		}
	}
}

// returns true iff any of the players have >= limit empty spheres *before*
// they're finished.
func playerHasConsecutiveEmptySpheres(routes []*routeInfo, limit int) bool {
	g, _, spheres, _ := getAllSpheres(routes)

	// clean up getAllSpheres master graph
	for _, ri := range routes {
		ri.graph["start"].removeParent(g["start"])
		g["done"].removeParent(ri.graph["done"])
	}

	// figure out what the final sphere for each player is
	finalSpheres := make([]int, len(routes))
	for i, sphere := range spheres {
		for _, check := range sphere {
			if i > finalSpheres[check.player-1] {
				finalSpheres[check.player-1] = i
			}
		}
	}

	// number of consecutive empty spheres a player has had
	drought := make([]int, len(routes))

	// check whether any player has an empty sphere before their final one
	for i, sphere := range spheres {
	playerLoop:
		for j, route := range routes {
			if i > finalSpheres[j] {
				continue
			}
			for _, node := range sphere {
				if route.slots[node.name] == node {
					drought[j] = 0
					continue playerLoop
				}
			}
			drought[j]++
			if drought[j] >= limit {
				return true
			}
		}
	}

	return false
}
