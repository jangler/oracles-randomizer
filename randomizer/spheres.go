package randomizer

import (
	"container/list"
	"fmt"
	"sort"
)

// getChecks converts a route info into a map of checks.
func getChecks(usedItems, usedSlots *list.List) map[*node]*node {
	checks := make(map[*node]*node)

	ei, es := usedItems.Front(), usedSlots.Front()
	for ei != nil {
		checks[es.Value.(*node)] = ei.Value.(*node)
		ei, es = ei.Next(), es.Next()
	}

	return checks
}

// getSpheres returns successive slices of checks that can be reached at a step
// in item collection. sphere 0 is the checks that can be reached from the
// start with no items; sphere 1 is the checks that can be reached using the
// items from sphere 0, and so on. each check only belongs to one sphere. it
// also returns a separate slice of checks that aren't reachable at all.
// returned slices are ordered alphabetically.
func getSpheres(g graph, checks map[*node]*node, resetFunc func()) ([][]*node, []*node) {
	reached := make(map[*node]bool)
	spheres := make([][]*node, 0)

	// need to track unreached items so that unreached dungeon items etc can
	// have their parents restored even if they're not reachable yet.
	unreachedChecks := make(map[*node]*node)
	for slot, item := range checks {
		// don't delimit spheres by intra-dungeon keys -- it obscures "actual"
		// progression in the log file.
		if !keyRegexp.MatchString(item.name) {
			unreachedChecks[slot] = item
			item.removeParent(slot)
		}
	}

	for {
		sphere := make([]*node, 0)
		g.reset()
		resetFunc()
		g["start"].explore()

		// get the set of newly reachable nodes
		for n, _ := range checks {
			if !reached[n] && n.reached {
				sphere = append(sphere, n)
			}
		}

		// mark nodes as reached and add item checks into the next iteration
		for _, n := range sphere {
			reached[n] = true
			if item := checks[n]; item != nil {
				if unreachedChecks[n] != nil {
					delete(unreachedChecks, n)
					item.addParent(n)
				}
				sphere = append(sphere, item)
				reached[item] = true
			}
		}

		if len(sphere) == 0 {
			break
		}
		spheres = append(spheres, sphere)
	}

	for slot, item := range unreachedChecks {
		item.addParent(slot)
	}

	extra := make([]*node, 0)
	for slot, item := range checks {
		if !reached[slot] {
			extra = append(extra, slot, item)
		}
	}

	for _, sphere := range append(spheres, extra) {
		sort.Slice(sphere, func(i, j int) bool {
			return sphere[i].name < sphere[j].name
		})
	}

	return spheres, extra
}

// logSpheres prints item placement by sphere to the summary channel.
func logSpheres(summary chan string, checks map[*node]*node,
	spheres [][]*node, extra []*node, game int, filter func(string) bool) {
	// don't print an extra newline before the first sphere in the section.
	firstSphere := true

	for i, sphere := range append(spheres, extra) {
		// get lines first, to make sure there are actual relevant items in
		// this sphere.
		lines := make([]string, 0)
		for slot, item := range checks {
			if filter != nil && !filter(item.name) {
				continue
			}
			for _, n := range sphere {
				if n == slot {
					if slot.player == 0 {
						lines = append(lines, fmt.Sprintf("%-28s <- %s",
							getNiceName(slot.name, game),
							getNiceName(item.name, game)))
					} else {
						lines = append(lines, fmt.Sprintf("P%d %-28s <- P%d %s",
							slot.player,
							getNiceName(slot.name, game),
							checks[slot].player,
							getNiceName(item.name, game)))
					}
					break
				}
			}
		}

		// then log the sphere if it's non-empty.
		if len(lines) > 0 {
			if firstSphere {
				firstSphere = false
			} else {
				summary <- ""
			}

			if i < len(spheres) {
				summary <- fmt.Sprintf("sphere %d:", i)
			} else {
				summary <- "inaccessible:"
			}

			sort.Strings(lines)
			for _, line := range lines {
				summary <- line
			}
		}
	}
}

// collates all the checks from multiple routes and returns check/sphere data.
// also returns a "master graph" which contains all the route graphs.
func getAllSpheres(routes []*routeInfo) (graph, map[*node]*node, [][]*node, []*node) {
	checks, spheres := make(map[*node]*node), make([][]*node, 0)
	for _, ri := range routes {
		for k, v := range getChecks(ri.usedItems, ri.usedSlots) {
			checks[k] = v
		}
	}
	g := newGraph()
	g["start"] = newNode("start", andNode)
	g["done"] = newNode("done", andNode)
	for _, ri := range routes {
		ri.graph["start"].addParent(g["start"])
		g["done"].addParent(ri.graph["done"])
	}
	spheres, extra := getSpheres(g, checks, func() {
		for _, ri := range routes {
			ri.graph.reset()
		}
	})
	return g, checks, spheres, extra
}
