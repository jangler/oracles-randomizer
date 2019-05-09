package main

import (
	"fmt"
	"sort"

	"github.com/jangler/oracles-randomizer/logic"
)

// getChecks converts a route info into a map of checks.
func getChecks(ri *RouteInfo) map[*node]*node {
	checks := make(map[*node]*node)

	ei, es := ri.UsedItems.Front(), ri.UsedSlots.Front()
	for ei != nil {
		checks[es.Value.(*node)] = ei.Value.(*node)
		ei, es = ei.Next(), es.Next()
	}

	return checks
}

// getSpheres returns successive slices of nodes that can be reached at a step
// in item collection. sphere 0 is the nodes that can be reached from the
// start with no items; sphere 1 is the nodes that can be reached using the
// items from sphere 0, and so on. each node only belongs to one sphere.
func getSpheres(g graph, checks map[*node]*node, hard bool) [][]*node {
	reached := make(map[*node]bool)
	spheres := make([][]*node, 0)

	// need to track unreached items so that unreached dungeon items etc can
	// have their parents restored even if they're not reachable yet.
	unreached := make(map[*node]*node)
	for slot, item := range checks {
		// don't delimit spheres by intra-dungeon keys -- it obscured "actual"
		// progression in the log file.
		if !keyRegexp.MatchString(item.name) {
			unreached[slot] = item
			item.removeParent(slot)
		}
	}

	rupees := 0
	for {
		sphere := make([]*node, 0)
		g.clearMarks()

		// get the set of newly reachable nodes
		for _, n := range g {
			if !reached[n] && n.getMark() == markTrue {
				if logic.NodeValues[n.name] > 0 {
					rupees += logic.NodeValues[n.name]
				}
				sphere = append(sphere, n)
			}
		}

		// remove the most expensive nodes that can't be afforded
		sphere, rupees = filterUnaffordableNodes(sphere, rupees)

		// mark nodes as reached and add item checks into the next iteration
		for _, n := range sphere {
			reached[n] = true
			delete(unreached, n)
			if item := checks[n]; item != nil {
				item.addParent(n)
				sphere = append(sphere, item)
				reached[item] = true
				rupees += logic.RupeeValues[item.name]

				// shovel is worth infinite rupees in hard difficulty
				if hard && item.name == "shovel" {
					rupees += 2000
				}
			}
		}

		if len(sphere) == 0 {
			break
		}
		spheres = append(spheres, sphere)
	}

	for slot, item := range unreached {
		item.addParent(slot)
	}

	return spheres
}

// logSpheres prints item placement by sphere to the summary channel.
func logSpheres(summary chan string, checks map[*node]*node,
	spheres [][]*node, game int, filter func(string) bool) {
	for i, sphere := range spheres {
		// get lines first, to make sure there are actual relevant items in
		// this sphere.
		lines := make([]string, 0)
		for slot, item := range checks {
			if !filter(item.name) {
				continue
			}
			for _, n := range sphere {
				if n == slot {
					lines = append(lines, fmt.Sprintf("%-28s <- %s",
						getNiceName(slot.name, game),
						getNiceName(item.name, game)))
					break
				}
			}
		}

		// then log the sphere if it's non-empty.
		if len(lines) > 0 {
			summary <- fmt.Sprintf("sphere %d:", i)
			sort.Strings(lines)
			for _, line := range lines {
				summary <- line
			}
			summary <- ""
		}
	}
}

// filterUnaffordableNodes removes nodes that the player can't currently afford
// from the slice, starting with the most expensive ones. it also sorts the
// slice from least to most expensive.
func filterUnaffordableNodes(sphere []*node, rupees int) ([]*node, int) {
	// sort first by name (to break ties), then by cost
	sort.Slice(sphere, func(i, j int) bool {
		return sphere[i].name < sphere[j].name
	})
	sort.Slice(sphere, func(i, j int) bool {
		return logic.NodeValues[sphere[i].name] >
			logic.NodeValues[sphere[j].name]
	})

	for i := 0; i < len(sphere); i++ {
		value := logic.NodeValues[sphere[i].name]
		if value < 0 {
			rupees += value
			if rupees < 0 {
				rupees -= value
				sphere = sphere[:i]
				break
			}
		}
	}

	return sphere, rupees
}
