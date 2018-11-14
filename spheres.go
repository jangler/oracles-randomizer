package main

import (
	"fmt"
	"sort"

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

	// need to track unreached items so that unreached dungeon items etc can
	// have their parents restored even if they're not reachable yet.
	unreached := make(map[*graph.Node]*graph.Node)
	for slot, item := range checks {
		unreached[slot] = item
		item.RemoveParent(slot)
	}

	rupees := 0
	for {
		sphere := make([]*graph.Node, 0)
		g.ClearMarks()

		// get the set of newly reachable nodes
		for _, node := range g {
			if !reached[node] && node.GetMark(node, hard) == graph.MarkTrue {
				if logic.NodeValues[node.Name] > 0 {
					rupees += logic.NodeValues[node.Name]
				}
				sphere = append(sphere, node)
			}
		}

		// remove the most expensive nodes that can't be afforded
		sphere, rupees = filterUnaffordableNodes(sphere, rupees)

		// mark nodes as reached and add item checks into the next iteration
		for _, node := range sphere {
			reached[node] = true
			delete(unreached, node)
			if item := checks[node]; item != nil {
				item.AddParents(node)
				sphere = append(sphere, item)
				reached[item] = true
				rupees += logic.RupeeValues[item.Name]

				// shovel is worth infinite rupees in hard difficulty
				if hard && item.Name == "shovel" {
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
		item.AddParents(slot)
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

// filterUnaffordableNodes removes nodes that the player can't currently afford
// from the slice, starting with the most expensive ones. it also sorts the
// slice from least to most expensive.
func filterUnaffordableNodes(
	sphere []*graph.Node, rupees int) ([]*graph.Node, int) {
	// sort first by cost, then by name to break ties
	sort.Slice(sphere, func(i, j int) bool {
		return logic.NodeValues[sphere[i].Name] >
			logic.NodeValues[sphere[j].Name]
	})
	sort.Slice(sphere, func(i, j int) bool {
		return sphere[i].Name < sphere[j].Name
	})

	for i := 0; i < len(sphere); i++ {
		value := logic.NodeValues[sphere[i].Name]
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
