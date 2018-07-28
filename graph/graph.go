package graph

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
func (g Graph) Explore(start map[*Node]bool, add []*Node) map[*Node]bool {
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

		for _, child := range node.Children {
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
			if node.GetMark(node, nil) == MarkTrue {
				reached[node] = true
				node.Mark = MarkTrue
				for _, child := range node.Children {
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
