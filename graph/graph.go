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
					child.AddParent(parent)
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

// Mark is the current state of a Node in its evaluation. MarkNone means the
// Node has not been evaludated, MarkTrue and MarkFalse mean that the value of
// the node has been determined, and MarkPending is used temporarily to prevent
// evaluating infinite loops in the graph.
type Mark uint8

const (
	MarkNone Mark = iota
	MarkTrue
	MarkFalse
	MarkPending
)

// NodeType determines how a node approaches GetMark(). And nodes return
// MarkTrue iff all of their parents do, Or nodes return MarkTrue iff any of
// their parents do, and Count nodes return true iff at least a certain number
// of their parents do.
type NodeType uint8

const (
	AndType NodeType = iota
	OrType
	CountType
)

// A Node is a single point in the directed graph.
type Node struct {
	Name     string
	Type     NodeType
	Mark     Mark
	MinCount int
	parents  []*Node
}

// NewNode returns a new unconnected graph node, not yet part of any graph.
func NewNode(name string, nodeType NodeType) *Node {
	// create node
	return &Node{
		Name:    name,
		Type:    nodeType,
		Mark:    MarkNone,
		parents: make([]*Node, 0),
	}
}

func (n *Node) GetMark() Mark {
	switch n.Type {
	case AndType:
		return getAndMark(n)
	case OrType:
		return getOrMark(n)
	case CountType:
		return getCountMark(n)
	default:
		panic("unknown node type for node " + n.Name)
	}
}

// returns true iff all parents are true (no parents == true).
func getAndMark(n *Node) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		for _, parent := range n.parents {
			switch parent.GetMark() {
			case MarkPending, MarkFalse:
				n.Mark = MarkNone
				return MarkFalse
			}
		}
		if n.Mark == MarkPending {
			n.Mark = MarkTrue
		}
	}
	return n.Mark
}

// returns true iff any parent is true (no parents == false).
func getOrMark(n *Node) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		allPending := true

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.parents {
			switch parent.Mark {
			case MarkTrue:
				n.Mark = MarkTrue
				allPending = false
				break OrPeekLoop
			case MarkFalse:
				allPending = false
			}
		}

		// then actually check them otherwise
		if n.Mark == MarkPending {
		OrGetLoop:
			for _, parent := range n.parents {
				switch parent.GetMark() {
				case MarkTrue:
					n.Mark = MarkTrue
					allPending = false
					break OrGetLoop
				case MarkFalse:
					allPending = false
				}
			}
		}

		if (allPending && len(n.parents) > 0) || n.Mark == MarkPending {
			n.Mark = MarkNone
			return MarkFalse
		}
	}

	return n.Mark
}

// returns true iff at least x parents are true.
func getCountMark(n *Node) Mark {
	count := 0

	if n.Mark == MarkNone {
		n.Mark = MarkPending

		for _, parent := range n.parents[0].parents {
			switch parent.GetMark() {
			case MarkPending, MarkFalse:
				continue
			default:
				count++
			}

			if count >= n.MinCount {
				n.Mark = MarkTrue
				return n.Mark
			}
		}

		if n.Mark == MarkPending {
			n.Mark = MarkNone
			return MarkFalse
		}
	}

	return n.Mark
}

// AddParent makes the given parent a parent of the node. If it is already a
// parent of the node, nothing is done.
func (n *Node) AddParent(parent *Node) {
	for _, p := range n.parents {
		if p == parent {
			return
		}
	}
	n.parents = append(n.parents, parent)
}

// RemoveParent removes the given node from this node's parents. It panics if
// the given node isn't actually a parent of this node.
func (n *Node) RemoveParent(parent *Node) {
	for i, p := range n.parents {
		if p == parent {
			n.parents = append(n.parents[:i], n.parents[i+1:]...)
			return
		}
	}
	panic(fmt.Sprintf("RemoveParent: %v is not a parent of %v", parent, n))
}

// ClearParents makes the node into an effective root node (though not a Root
// node).
func (n *Node) ClearParents() {
	n.parents = n.parents[:0]
}

// String satisfies the fmt.Stringer interface.
func (n *Node) String() string { return n.Name }
