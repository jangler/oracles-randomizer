package graph

import (
	"container/list"
	"fmt"
	"strings"
)

// this file defines operations on a dependency graph for the game, but does
// not itself define the graph.

// Mark is the current state of a Node in its evaluation. When a non-root Node
// is evaluated, it is set to MarkFalse until proven otherwise. This is to
// prevent evaluating (infinite) loops in the graph.
type Mark int

// NodeType determines how a node approaches GetMark(). And nodes return
// MarkTrue only if all of their parents do, Or nodes return MarkTrue if any of
// their parents do, and Root nodes always return MarkTrue.
//
// Technically an And node with no parents functions the same as a Root node,
// and an Or node with no parents always returns MarkFalse.
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
	Name     string
	Type     NodeType
	GetMark  func(*Node, *list.List) Mark
	IsStep   bool
	Mark     Mark
	Parents  []*Node
	Children []*Node
}

// NewNode returns a new unconnected graph node, not yet part of any graph.
func NewNode(name string, nodeType NodeType, isStep bool) *Node {
	// create node
	n := Node{
		Name:     name,
		Type:     nodeType,
		IsStep:   isStep,
		Mark:     MarkNone,
		Parents:  make([]*Node, 0),
		Children: make([]*Node, 0),
	}

	// set node's GetMark function based on type
	switch n.Type {
	case RootType:
		n.GetMark = getRootMark
	case AndType:
		n.GetMark = getAndMark
	case OrType:
		n.GetMark = getOrMark
	default:
		panic("unknown node type for node " + name)
	}

	return &n
}

func getRootMark(n *Node, path *list.List) Mark {
	return MarkFalse
}

func getAndMark(n *Node, path *list.List) Mark {
	if n.Mark == MarkNone {
		var parentNames []string
		if path != nil {
			parentNames = make([]string, len(n.Parents))
		}

		n.Mark = MarkPending
		for i, parent := range n.Parents {
			switch parent.GetMark(parent, path) {
			case MarkPending, MarkFalse:
				n.Mark = MarkNone
				return MarkFalse
			}
			if parentNames != nil {
				parentNames[i] = parent.Name
			}
		}
		if n.Mark == MarkPending {
			n.Mark = MarkTrue
		}

		if path != nil && n.Mark == MarkTrue {
			if len(parentNames) > 0 {
				path.PushBack(n.Name + " <- " + strings.Join(parentNames, ", "))
			} else {
				path.PushBack(n.Name)
			}
		}
	}

	return n.Mark
}

func getOrMark(n *Node, path *list.List) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		allPending := true
		var parentName string

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.Parents {
			switch parent.Mark {
			case MarkTrue:
				n.Mark = MarkTrue
				allPending = false
				parentName = parent.Name
				break OrPeekLoop
			case MarkFalse:
				allPending = false
			}
		}

		// then actually check them otherwise
		if n.Mark == MarkPending {
		OrGetLoop:
			for _, parent := range n.Parents {
				switch parent.GetMark(parent, path) {
				case MarkTrue:
					n.Mark = MarkTrue
					allPending = false
					parentName = parent.Name
					break OrGetLoop
				case MarkFalse:
					allPending = false
				}
			}
		}

		if (allPending && len(n.Parents) > 0) || n.Mark == MarkPending {
			n.Mark = MarkNone
			return MarkFalse
		}

		if path != nil && n.Mark == MarkTrue {
			path.PushBack(fmt.Sprintf("%s <- %s", n.Name, parentName))
		}
	}

	return n.Mark
}

// AddParents makes the given nodes parents of the node, and likewise adds this
// node to each parent's list of children.
func (n *Node) AddParents(parents ...*Node) {
	n.Parents = append(n.Parents, parents...)
	addChild(n, parents...)
}

// ClearParents makes the node into an effective root node (though not a Root
// node).
func (n *Node) ClearParents() {
	removeChild(n, n.Parents...)
	n.Parents = n.Parents[:0]
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
		parent.Children = append(parent.Children, child)
		parent.Children = append(parent.Children, child)
	}
}

func removeChild(child *Node, parents ...*Node) {
	for _, parent := range parents {
		removeNodeFromSlice(child, &parent.Children)
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
