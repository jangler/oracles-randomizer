package main

import (
	"fmt"

	"github.com/jangler/oos-randomizer/graph"
)

// this file contains the actual connection of nodes in the game graph, and
// tracks them as they update.

// XXX need to be careful about rings. i can't imagine a situation where you'd
//     need both energy ring and fist ring, but if you did, then you'd need to
//     have the L-2 ring box to do so without danger of soft locking.

// these types are just for readability
type And []string
type Or []string

func initRoute() (*graph.Graph, []error) {
	g := graph.NewGraph()

	g.AddOrNodes(baseItemNodes...)
	addAndOrNodes(g, itemNodesAnd, itemNodesOr)
	addAndOrNodes(g, killNodesAnd, killNodesOr)
	addAndOrNodes(g, d0NodesAnd, d0NodesOr)
	addAndOrNodes(g, d1NodesAnd, d1NodesOr)
	addAndOrNodes(g, d2NodesAnd, d2NodesOr)
	addAndOrNodes(g, subrosiaNodesAnd, subrosiaNodesOr)
	addAndOrNodes(g, portalNodesAnd, portalNodesOr)
	addAndOrNodes(g, holodrumNodesAnd, holodrumNodesOr)

	g.AddParents(itemNodesAnd)
	g.AddParents(itemNodesOr)
	g.AddParents(killNodesAnd)
	g.AddParents(killNodesOr)
	g.AddParents(d0NodesAnd)
	g.AddParents(d0NodesOr)
	g.AddParents(d1NodesAnd)
	g.AddParents(d1NodesOr)
	g.AddParents(d2NodesAnd)
	g.AddParents(d2NodesOr)
	g.AddParents(subrosiaNodesAnd)
	g.AddParents(subrosiaNodesOr)
	g.AddParents(portalNodesAnd)
	g.AddParents(portalNodesOr)
	g.AddParents(holodrumNodesAnd)
	g.AddParents(holodrumNodesOr)

	// validate
	var errs []error
	for name, node := range g.Map {
		switch nt := node.(type) {
		case graph.ChildNode:
			if !nt.HasParents() {
				if errs == nil {
					errs = make([]error, 0)
				}
				errs = append(errs, fmt.Errorf("orphan node: %s", name))
			}
		}
	}

	return g, errs
}

func addAndOrNodes(g *graph.Graph, andNodes, orNodes map[string][]string) {
	for key, _ := range andNodes {
		g.AddAndNodes(key)
	}
	for key, _ := range orNodes {
		g.AddOrNodes(key)
	}
}
