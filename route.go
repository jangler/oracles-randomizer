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

// A Point is a mapping of point strings that will become And or Or nodes in
// the graph.
type Point interface {
	Parents() []string
}

// the different types of points are all just string slices; the reason for
// having different ones is purely for type assertions

type And []string

func (p And) Parents() []string { return p }

type Or []string

func (p Or) Parents() []string { return p }

type AndSlot []string

func (p AndSlot) Parents() []string { return p }

type OrSlot []string

func (p OrSlot) Parents() []string { return p }

func initRoute() (*graph.Graph, map[string]Point, []error) {
	g := graph.NewGraph()

	totalPoints := make(map[string]Point, 0)
	appendNodes(totalPoints,
		baseItemNodes, itemNodesAnd, itemNodesOr,
		killNodesAnd, killNodesOr,
		holodrumNodesAnd, holodrumNodesOr,
		subrosiaNodesAnd, subrosiaNodesOr,
		portalNodesAnd, portalNodesOr,
		d0NodesAnd, d0NodesOr,
		d1NodesAnd, d1NodesOr,
		d2NodesAnd, d2NodesOr,
	)

	addPointNodes(g, totalPoints)
	addPointParents(g, totalPoints)

	openSlots := make(map[string]Point, 0)
	for name, point := range totalPoints {
		switch point.(type) {
		case AndSlot, OrSlot:
			openSlots[name] = point
		}
	}

	// validate
	var errs []error
	for name, node := range g.Map {
		if !node.HasParents() {
			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("orphan node: %s", name))
		}
	}

	return g, openSlots, errs
}

func appendNodes(total map[string]Point, pointMaps ...map[string]Point) {
	for _, pointMap := range pointMaps {
		for k, v := range pointMap {
			total[k] = v
		}
	}
}

func addPointNodes(g *graph.Graph, points map[string]Point) {
	for key, pt := range points {
		switch pt.(type) {
		case And, AndSlot:
			g.AddAndNodes(key)
		case Or, OrSlot:
			g.AddOrNodes(key)
		default:
			panic("unknown point type for " + key)
		}
	}
}

func addPointParents(g *graph.Graph, points map[string]Point) {
	// TODO optimize?
	for k, p := range points {
		g.AddParents(map[string][]string{k: p.Parents()})
	}
}
