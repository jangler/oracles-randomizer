package main

import (
	"fmt"
)

// map of names to set of nodes in a directed graph
type graph map[string]*node

// returns a new, initialized graph.
func newGraph() graph {
	return graph(make(map[string]*node))
}

// adds relationships in bulk between existing nodes in the graph, by name.
// attempting to link a name not in the graph results in a panic.
func (g graph) addParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := g[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := g[parentName]; ok {
					child.addParent(parent)
				} else {
					panic("no node named " + parentName)
				}
			}
		} else {
			panic("no child named " + childName)
		}
	}
}

// resets all the nodes in a graph to an "unknown" state. this is required any
// time relationships in the graph change.
func (g graph) clearMarks() {
	for _, n := range g {
		n.mark = markNone
	}
}

// the current state of a node in its evaluation. markNone means the node has
// not been evaludated, markTrue and markFalse mean that the value of the node
// has been determined, and markPending is used temporarily to prevent
// evaluating infinite loops in the graph.
type nodeMark uint8

const (
	markNone nodeMark = iota
	markTrue
	markFalse
	markPending
)

// determines how a node approaches getMark(). an andNode returns markTrue iff
// all of its parents do, an orNode returns markTrue iff any of its parents do,
// and a countNode returns true iff at least a certain number of its parents
// do.
type nodeType uint8

const (
	andNode nodeType = iota
	orNode
	countNode
)

// a single vertex in the graph.
type node struct {
	name     string
	nType    nodeType
	mark     nodeMark
	minCount int
	parents  []*node
}

// returns a new unconnected graph node, not yet part of any graph.
func newNode(name string, nType nodeType) *node {
	// create node
	return &node{
		name:    name,
		nType:   nType,
		mark:    markNone,
		parents: make([]*node, 0),
	}
}

func (n *node) getMark() nodeMark {
	switch n.nType {
	case andNode:
		return getAndMark(n)
	case orNode:
		return getOrMark(n)
	case countNode:
		return getCountMark(n)
	default:
		panic("unknown node type for node " + n.name)
	}
}

// returns true iff all parents are true (no parents == true).
func getAndMark(n *node) nodeMark {
	if n.mark == markNone {
		n.mark = markPending
		for _, parent := range n.parents {
			switch parent.getMark() {
			case markPending, markFalse:
				n.mark = markNone
				return markFalse
			}
		}
		if n.mark == markPending {
			n.mark = markTrue
		}
	}
	return n.mark
}

// returns true iff any parent is true (no parents == false).
func getOrMark(n *node) nodeMark {
	if n.mark == markNone {
		n.mark = markPending
		allPending := true

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.parents {
			switch parent.mark {
			case markTrue:
				n.mark = markTrue
				allPending = false
				break OrPeekLoop
			case markFalse:
				allPending = false
			}
		}

		// then actually check them otherwise
		if n.mark == markPending {
		OrGetLoop:
			for _, parent := range n.parents {
				switch parent.getMark() {
				case markTrue:
					n.mark = markTrue
					allPending = false
					break OrGetLoop
				case markFalse:
					allPending = false
				}
			}
		}

		if (allPending && len(n.parents) > 0) || n.mark == markPending {
			n.mark = markNone
			return markFalse
		}
	}

	return n.mark
}

// returns true iff at least x parents are true.
func getCountMark(n *node) nodeMark {
	count := 0

	if n.mark == markNone {
		n.mark = markPending

		for _, parent := range n.parents[0].parents {
			switch parent.getMark() {
			case markPending, markFalse:
				continue
			default:
				count++
			}

			if count >= n.minCount {
				n.mark = markTrue
				return n.mark
			}
		}

		if n.mark == markPending {
			n.mark = markNone
			return markFalse
		}
	}

	return n.mark
}

// makes a node a parent of another node. a node can be a parent of another
// node multiple times; i.e. it can appear twice or more in the child's list of
// parents.
func (n *node) addParent(parent *node) {
	n.parents = append(n.parents, parent)
}

// removes the given node from this node's parents, once. it panics if the
// given node isn't actually a parent of this node.
func (n *node) removeParent(parent *node) {
	for i, p := range n.parents {
		if p == parent {
			n.parents = append(n.parents[:i], n.parents[i+1:]...)
			return
		}
	}
	panic(fmt.Sprintf("removeParent: %v is not a parent of %v", parent, n))
}

// removes all parent connections from the node.
func (n *node) clearParents() {
	n.parents = n.parents[:0]
}

// satisfies the fmt.Stringer interface.
func (n *node) String() string { return n.name }
