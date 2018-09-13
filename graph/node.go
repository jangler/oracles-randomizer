package graph

import (
	"fmt"
)

// this file defines operations on a dependency graph for the game, but does
// not itself define the graph.

// Mark is the current state of a Node in its evaluation. When a non-root Node
// is evaluated, it is set to MarkFalse until proven otherwise. This is to
// prevent evaluating (infinite) loops in the graph.
type Mark int

// NodeType determines how a node approaches GetMark(). And nodes return
// MarkTrue only if all of their parents do, Or nodes return MarkTrue if any of
// their parents do, and Root act as Or nodes, but conventionally start without
// parents (Or nodes without parents return MarkFalse).
//
// An And node with no parents always returns MarkTrue.
type NodeType int

// See Mark and NodeType comments for information.
const (
	MarkNone    Mark = iota // satisfied depending on parents
	MarkTrue                // succeed an OrNode, continue an AndNode
	MarkFalse               // continue an OrNode, fail an AndNode
	MarkPending             // prevents circular dependencies

	// nodes will not ever set themselves to MarkFalse, but they will return it
	// if they are set to MarkNone and are not satisfied

	RootType NodeType = iota
	AndType
	OrType
)

// A Node is a single point in the directed graph.
type Node struct {
	Name       string
	Type       NodeType
	GetMark    func(*Node, bool) Mark
	IsStep     bool
	IsSlot     bool
	IsHard     bool
	IsOptional bool
	Mark       Mark
	parents    []*Node
	children   []*Node
}

// NewNode returns a new unconnected graph node, not yet part of any graph.
func NewNode(name string, nodeType NodeType,
	isStep, isSlot, isHard bool) *Node {
	// create node
	n := Node{
		Name:     name,
		Type:     nodeType,
		IsStep:   isStep,
		IsSlot:   isSlot,
		IsHard:   isHard,
		Mark:     MarkNone,
		parents:  make([]*Node, 0),
		children: make([]*Node, 0),
	}

	// set node's GetMark function based on type
	switch n.Type {
	case RootType:
		n.GetMark = getOrMark
	case AndType:
		n.GetMark = getAndMark
	case OrType:
		n.GetMark = getOrMark
	default:
		panic("unknown node type for node " + name)
	}

	return &n
}

func getAndMark(n *Node, hard bool) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		for _, parent := range n.parents {
			if !hard && parent.IsHard {
				continue
			}

			switch parent.GetMark(parent, hard) {
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

func getOrMark(n *Node, hard bool) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		allPending := true

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.parents {
			if !hard && parent.IsHard {
				continue
			}

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
				if !hard && parent.IsHard {
					continue
				}

				switch parent.GetMark(parent, hard) {
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

// Parents returns a copy of the node's slice of parents.
func (n *Node) Parents() []*Node {
	parents := make([]*Node, 0, len(n.parents))
	parents = append(parents, n.parents...)
	return parents
}

// NumParents returns the number of parents the node has. This is mode time-
// and memory-efficient than taking the length of node.Parents().
func (n *Node) NumParents() int {
	return len(n.parents)
}

// AddParents makes the given nodes parents of the node, and likewise adds this
// node to each parent's list of children. If a given parent is already a
// parent of the node, nothing is done.
func (n *Node) AddParents(parents ...*Node) {
	for _, parent := range parents {
		if !IsNodeInSlice(parent, n.parents) {
			n.parents = append(n.parents, parent)
			addChild(n, parent)
		}
	}
}

// RemoveParent removes the given node from this node's parents. It panics if
// the given node isn't actually a parent of this node.
func (n *Node) RemoveParent(parent *Node) {
	for i, p := range n.parents {
		if p == parent {
			n.parents = append(n.parents[:i], n.parents[i+1:]...)
			removeChild(parent, n)
			return
		}
	}

	panic(fmt.Sprintf("RemoveParent: %v is not a parent of %v", parent, n))
}

// PopParent removes and returns the last parent of the node.
func (n *Node) PopParent() *Node {
	p := n.parents[len(n.parents)-1]
	n.RemoveParent(p)
	return p
}

// ClearParents makes the node into an effective root node (though not a Root
// node).
func (n *Node) ClearParents() {
	removeChild(n, n.parents...)
	n.parents = n.parents[:0]
}

// String satisfies the fmt.Stringer interface.
func (n *Node) String() string { return n.Name }

// helper functions

// IsNodeInSlice returns true iff the node is in the slice of nodes.
func IsNodeInSlice(node *Node, slice []*Node) bool {
	for _, match := range slice {
		if node == match {
			return true
		}
	}
	return false
}

func addChild(child *Node, parents ...*Node) {
	for _, parent := range parents {
		if !IsNodeInSlice(child, parent.children) {
			parent.children = append(parent.children, child)
		}
	}
}

func removeChild(child *Node, parents ...*Node) {
	for _, parent := range parents {
		removeNodeFromSlice(child, &parent.children)
	}
}

func removeNodeFromSlice(node *Node, slice *[]*Node) {
	// O(n)
	for i, match := range *slice {
		if match == node {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}
