package main

import (
	"fmt"
	"strings"

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

type Route struct {
	Graph *graph.Graph
	Slots map[string]Point
}

// total
var nonGeneratedPoints map[string]Point

func init() {
	nonGeneratedPoints = make(map[string]Point)
	appendPoints(nonGeneratedPoints,
		baseItemNodes, ignoredBaseItemNodes,
		itemNodesAnd, itemNodesOr,
		killNodesAnd, killNodesOr,
		holodrumNodesAnd, holodrumNodesOr,
		subrosiaNodesAnd, subrosiaNodesOr,
		portalNodesAnd, portalNodesOr,
		d0NodesAnd, d0NodesOr,
		d1NodesAnd, d1NodesOr,
		d2NodesAnd, d2NodesOr,
	)
}

func initRoute() (*Route, []error) {
	g := graph.NewGraph()

	totalPoints := make(map[string]Point, 0)
	appendPoints(totalPoints, nonGeneratedPoints, generatedPoints)

	// ignore semicolon-delimited points; they're only used for generation
	for key := range totalPoints {
		if strings.ContainsRune(key, ';') {
			delete(totalPoints, key)
		}
	}

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
		if baseItemNodes[name] != nil || ignoredBaseItemNodes[name] != nil ||
			name == "horon village" {
			// it's supposed to be orphan/childless; skip it
			continue
		}

		// check for parents and children
		if len(node.Parents()) == 0 {
			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("orphan node: %s", name))
		}
		if len(node.Children()) == 0 {
			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("childless node: %s", name))
		}
	}

	return &Route{Graph: g, Slots: openSlots}, errs
}

func appendPoints(total map[string]Point, pointMaps ...map[string]Point) {
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
	// ugly but w/e
	for k, p := range points {
		g.AddParents(map[string][]string{k: p.Parents()})
	}
}
