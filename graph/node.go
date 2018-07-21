package graph

import (
	"container/list"
)

// this file defines operations on a dependency graph for the game, but does
// not itself define the graph.

// Mark is the current state of a Node in its evaluation. When a non-root Node
// is evaluated, it is set to MarkFalse until proven otherwise. This is to
// prevent evaluating (infinite) loops in the graph.
type Mark int

const (
	MarkNone  Mark = iota // satisfied depending on parents
	MarkTrue              // succeed an OrNode, continue an AndNode
	MarkFalse             // continue an OrNode, fail an AndNode
)

// Node is the general interface that encompasses everything in the graph.
type Node interface {
	GetName() string
	GetMark(*list.List) Mark  // list to append path to if non-nil
	PeekMark(*list.List) Mark // like GetMark but doesn't check parents
	SetMark(Mark)
}

// ChildNode is a node with parent(s).
type ChildNode interface {
	Node
	AddParents(...Node)
	HasParents() bool
}

// RootNode has no parents and is set to MarkTrue when queried if its mark is
// MarkNone.
type RootNode struct {
	Name string
	Mark Mark
}

func (n *RootNode) GetName() string { return n.Name }

func (n *RootNode) GetMark(path *list.List) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkTrue
	}

	if path != nil && n.Mark != MarkFalse {
		path.PushBack(n.Name)
	}

	return n.Mark
}

func (n *RootNode) PeekMark(path *list.List) Mark {
	return n.Mark
}

func (n *RootNode) SetMark(m Mark) {
	n.Mark = m
}

// AndNode is satisfied if all of its parents are satisfied, or if it has no
// parents.
type AndNode struct {
	Name    string
	Mark    Mark
	Parents []Node
}

func (n *AndNode) GetName() string { return n.Name }

func (n *AndNode) GetMark(path *list.List) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkTrue
		for _, parent := range n.Parents {
			if parent.GetMark(path) == MarkFalse {
				n.Mark = MarkFalse
				break
			}
		}
		if path != nil && n.Mark == MarkTrue {
			path.PushBack(n.Name)
		}
	}

	return n.Mark
}

func (n *AndNode) PeekMark(path *list.List) Mark {
	return n.Mark
}

func (n *AndNode) SetMark(m Mark) {
	n.Mark = m
}

func (n *AndNode) AddParents(parents ...Node) {
	n.Parents = append(n.Parents, parents...)
}

func (n *AndNode) HasParents() bool {
	return len(n.Parents) > 0
}

// OrNode is satisfied if any of its parents is satisfied, unless it has no
// parents.
type OrNode struct {
	Name    string
	Mark    Mark
	Parents []Node
}

func (n *OrNode) GetName() string { return n.Name }

func (n *OrNode) GetMark(path *list.List) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkFalse

		// prioritize already satisfied nodes
		for _, parent := range n.Parents {
			if parent.PeekMark(path) == MarkTrue {
				n.Mark = MarkTrue
				break
			}
		}

		// then actually check them otherwise
		if n.Mark == MarkFalse {
			for _, parent := range n.Parents {
				if parent.GetMark(path) == MarkTrue {
					n.Mark = MarkTrue
					break
				}
			}
		}

		if path != nil && n.Mark == MarkTrue {
			path.PushBack(n.Name)
		}
	}

	return n.Mark
}

func (n *OrNode) PeekMark(path *list.List) Mark {
	return n.Mark
}

func (n *OrNode) SetMark(m Mark) {
	n.Mark = m
}

func (n *OrNode) AddParents(parents ...Node) {
	n.Parents = append(n.Parents, parents...)
}

func (n *OrNode) HasParents() bool {
	return len(n.Parents) > 0
}
