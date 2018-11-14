package main

import (
	"fmt"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/logic"
)

// getChecks converts a route info into a slice of checks.
func getChecks(ri *RouteInfo) map[*graph.Node]*graph.Node {
	checks := make(map[*graph.Node]*graph.Node)

	ei, es := ri.UsedItems.Front(), ri.UsedSlots.Front()
	for ei != nil {
		checks[es.Value.(*graph.Node)] = ei.Value.(*graph.Node)
		ei, es = ei.Next(), es.Next()
	}

	return checks
}

// getSpheres returns successive slices of nodes that can be reached at a step
// in item collection. sphere 0 is the nodes that can be reached from the
// start with no items; sphere 1 is the nodes that can be reached using the
// items from sphere 0, and so on. each node only belongs to one sphere.
func getSpheres(g graph.Graph, checks map[*graph.Node]*graph.Node,
	hard bool) [][]*graph.Node {
	reached := make(map[*graph.Node]bool)
	spheres := make([][]*graph.Node, 0)

	for slot, item := range checks {
		item.RemoveParent(slot)
	}

	rupees := 0
	for {
		sphere := make([]*graph.Node, 0)
		g.ClearMarks()

		// get the set of newly reachable nodes
		for _, node := range g {
			if !reached[node] && node.GetMark(node, hard) == graph.MarkTrue {
				cost := logic.NodeCosts[node.Name]
				if checks[node] != nil && cost+rupees < 0 {
					continue
				}
				rupees += cost

				sphere = append(sphere, node)
				reached[node] = true
			}
		}

		// add reached item checks into the next iteration
		for _, node := range sphere {
			if item := checks[node]; item != nil {
				item.AddParents(node)
				sphere = append(sphere, item)
				reached[item] = true
				rupees += logic.RupeeValues[item.Name]
			}
		}

		if len(sphere) == 0 {
			break
		}
		spheres = append(spheres, sphere)
	}

	return spheres
}

// logSpheres prints item placement by sphere to the summary channel.
func logSpheres(summary chan string, checks map[*graph.Node]*graph.Node,
	spheres [][]*graph.Node, filter func(string) bool) {
	for i, sphere := range spheres {
		// get lines first, to make sure there are actual relevant items in
		// this sphere.
		lines := make([]string, 0)
		for slot, item := range checks {
			if !filter(item.Name) {
				continue
			}
			for _, node := range sphere {
				if node == slot {
					lines = append(lines, fmt.Sprintf("%-28s <- %s",
						getNiceName(slot.Name),
						getNiceName(item.Name)))
					break
				}
			}
		}

		// then log the sphere if it's non-empty.
		if len(lines) > 0 {
			summary <- fmt.Sprintf("sphere %d:", i)
			for _, line := range lines {
				summary <- line
			}
			summary <- ""
		}
	}
}
