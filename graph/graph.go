package graph

import (
	"log"
)

// this file contains facilities for linking nodes into a graph

type Graph struct {
	Map map[string]Node
}

func NewGraph() *Graph {
	return &Graph{
		Map: make(map[string]Node),
	}
}

// AddNodes adds the given nodes to the graph. It panics if a given node has
// the same name as one already in the graph.
func (g *Graph) AddNodes(nodes ...Node) {
	for _, node := range nodes {
		g.CheckDuplicateName(node.Name())
		g.Map[node.Name()] = node
	}
}

func (g *Graph) CheckDuplicateName(name string) {
	if g.Map[name] != nil {
		panic("node named " + name + " already in route map")
	}
}

func (g *Graph) AddParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := g.Map[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := g.Map[parentName]; ok {
					child.AddParents(parent)
				} else {
					panic("no node named " + parentName)
				}
			}
		} else {
			panic("no child named " + childName)
		}
	}
}

func (g *Graph) ClearMarks() {
	for _, node := range g.Map {
		node.SetMark(MarkNone)
	}
}

// Explore returns a new set of all nodes reachable from the set of nodes in
// start, adding the nodes in add and subtracting the nodes in sub. This is a
// destructive operation; at the end, the graph will have all nodes in the
// return set set to MarkTrue and the rest set to MarkNone.
//
// add and sub can be nil.
func (g *Graph) Explore(start map[Node]bool, add, sub []Node) map[Node]bool {
	// copy set, and mark nodes accordingly
	g.ClearMarks()
	reached := make(map[Node]bool, len(start))
	for node := range start {
		reached[node] = true
		node.SetMark(MarkTrue)
	}

	// make set of unchecked children
	frontier := make(map[Node]bool)

	// add nodes to the reached set and add their unreached children to the
	// frontier.
	for _, node := range add {
		reached[node] = true
		node.SetMark(MarkTrue)

		for _, child := range node.Children() {
			if !reached[child] {
				frontier[child] = true
			}
		}
	}

	// subtract nodes from the reached set and add their reached children to
	// the frontier
	for _, node := range sub {
		delete(reached, node)
		for _, child := range node.Children() {
			if reached[child] {
				frontier[child] = true
				child.SetMark(MarkNone)
			}
		}
	}

	log.Print(len(frontier), " nodes in frontier")
	tried := make(map[Node]bool) // and set of nodes that were already tried

	// explore. done when no new nodes are reached in an iteration
	for len(frontier) > 0 {
		for node := range frontier {
			if node.GetMark(nil) == MarkTrue {
				reached[node] = true
				node.SetMark(MarkTrue)
				for _, child := range node.Children() {
					if !reached[child] && !tried[child] {
						frontier[child] = true
					}
				}
			}
			delete(frontier, node)
			tried[node] = true
		}
	}

	return reached
}
