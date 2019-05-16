package main

import (
	"fmt"
	"sort"
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

	for {
		sphere := make([]*node, 0)
		g.clearMarks()

		// get the set of newly reachable nodes
		for _, n := range g {
			if !reached[n] && n.getMark() == markTrue {
				sphere = append(sphere, n)
			}
		}

		// mark nodes as reached and add item checks into the next iteration
		for _, n := range sphere {
			reached[n] = true
			delete(unreached, n)
			if item := checks[n]; item != nil {
				item.addParent(n)
				sphere = append(sphere, item)
				reached[item] = true
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
