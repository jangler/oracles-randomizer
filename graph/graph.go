package graph

import (
	"fmt"
)

// A Graph maps names to a set of (hopeully) connected nodes. The graph is
// directed.
type Graph map[string]*Node

// New returns an initialized, empty graph.
func New() Graph {
	return Graph(make(map[string]*Node))
}

// AddNodes adds the given nodes to the graph. Name collision is a fatal error.
func (g Graph) AddNodes(nodes ...*Node) {
	for _, node := range nodes {
		if g[node.Name] != nil {
			panic("node name already in graph: " + node.Name)
		}
		g[node.Name] = node
	}
}

// AddParents adds relationships in bulk between existing nodes in the graph,
// by name. Attempting to link a name not in the graph results in a panic.
func (g Graph) AddParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := g[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := g[parentName]; ok {
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

// ClearMarks resets all the nodes in a graph to an "unknown" state. This is
// required any time relationships in the graph change.
func (g Graph) ClearMarks() {
	for _, node := range g {
		node.Mark = MarkNone
	}
}

// Explore returns a new set of all nodes reachable from the set of nodes in
// start, adding the nodes in add. This is a destructive operation; at the end,
// the graph will have all nodes in the return set set to MarkTrue and the rest
// set to MarkNone.
func (g Graph) Explore(start map[*Node]bool, hard bool,
	add ...*Node) map[*Node]bool {
	// copy set, and mark nodes accordingly
	g.ClearMarks()
	reached := make(map[*Node]bool, len(start))
	for node := range start {
		reached[node] = true
		node.Mark = MarkTrue
	}

	// make set of unchecked children
	frontier := make(map[*Node]bool)

	// add nodes to the reached set and add their unreached children to the
	// frontier.
	for _, node := range add {
		reached[node] = true
		node.Mark = MarkTrue

		for _, child := range node.children {
			if !reached[child] {
				frontier[child] = true
			}
		}
	}

	// make set of nodes that were already processed as parents
	tried := make(map[*Node]bool)

	// explore. done when no new nodes are reached in an iteration
	for len(frontier) > 0 {
		for node := range frontier {
			// if we can reach the node, add it to the reached set and add its
			// (previously unchecked) children to the frontier
			if node.GetMark(node, hard) == MarkTrue {
				reached[node] = true
				node.Mark = MarkTrue
				for _, child := range node.children {
					if !reached[child] && !tried[child] {
						frontier[child] = true
					}
				}
			}

			// get this node out of my sight
			delete(frontier, node)
			tried[node] = true
		}
	}

	return reached
}

// ExploreFromStart calls explore without specific start and add nodes, instead
// exploring the entirety of the existing graph.
func (g Graph) ExploreFromStart(hard bool) map[*Node]bool {
	return g.Explore(nil, hard, g["start"])
}

// Reduce returns a version of the graph that is 1. only relevant to the given
// target and 2. reduced to as few nodes as possible.
func (g Graph) Reduce(target string) (Graph, error) {
	if g[target] == nil {
		return nil, fmt.Errorf("target node %s not in graph", target)
	}

	// copy graph but remove start node
	reduced := copyGraph(g)
	if start := g["start"]; start != nil {
		for _, child := range start.children {
			removeParent(child, start)
		}
		delete(g, "start")
	}

	// iteratively cut out parents with only one child, or zero children, as
	// long as the parent type matches the type of its single child. direct
	// parents of the target node also don't need to be parents of any other
	// node. (TODO this principle can be applied recursively)
	done := false
	for !done {
		done = true

		// collapse single-parent lines
		for name, node := range reduced {
			if name == target || node.Type == RootType {
				continue
			}

			switch len(node.children) {
			case 0:
				node.ClearParents()
				delete(reduced, name)
			case 1:
				if len(node.parents) == 1 ||
					node.Type == node.children[0].Type {
					done = false
					node.children[0].AddParents(node.parents...)
					removeParent(node.children[0], node)
					removeChild(node, node.parents...)
					delete(reduced, name)
				}
			}
		}

		// make direct parents of the target node parents only of that node
		for _, node := range reduced {
			if IsNodeInSlice(reduced[target], node.children) {
				for i := 0; i < len(node.children); i++ {
					if node.children[i] != reduced[target] {
						done = false
						removeParent(node.children[i], node)
						node.children =
							append(node.children[:i], node.children[i+1:]...)
						i--
					}
				}
			}
		}
	}

	return reduced, nil
}

// returns a new copy of the graph with new but identical nodes and
// relationships.
func copyGraph(old Graph) Graph {
	new := New()

	// add nodes
	for name, node := range old {
		new[name] = NewNode(node.Name, node.Type, node.IsStep, node.IsSlot,
			node.IsHard)
	}

	// add relationships
	for name, node := range old {
		for _, parent := range node.parents {
			newNode := new[name]
			newNode.parents = append(newNode.parents, new[parent.Name])
			children := new[parent.Name].children
			new[parent.Name].children = append(children, newNode)
		}
	}

	return new
}

// doesn't do anything if the child already doesn't have the parent
func removeParent(child, removal *Node) {
	for i, parent := range child.parents {
		if parent == removal {
			child.parents = append(child.parents[:i], child.parents[i+1:]...)
			break
		}
	}
}
