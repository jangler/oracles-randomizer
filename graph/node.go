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

const (
	MarkNone    Mark = iota // satisfied depending on parents
	MarkTrue                // succeed an OrNode, continue an AndNode
	MarkFalse               // continue an OrNode, fail an AndNode
	MarkPending             // prevents circular dependencies

	// nodes will not ever set themselves to MarkFalse, but they will return it
	// if they are set to MarkNone and are not satisfied
)

// Node is the general interface that encompasses everything in the graph.
type Node interface {
	fmt.Stringer

	Name() string
	GetMark(*list.List) Mark // list to append path to if non-nil
	PeekMark() Mark          // like GetMark but doesn't check parents
	SetMark(Mark)
	AddParents(...Node)
	ClearParents()
	Parents() []Node
	Children() []Node
}

// AndNode is satisfied if all of its parents are satisfied, or if it has no
// parents.
type AndNode struct {
	name     string
	mark     Mark
	parents  []Node
	children []Node
}

func NewAndNode(name string) *AndNode {
	return &AndNode{name: name,
		parents: make([]Node, 0), children: make([]Node, 0)}
}

func (n *AndNode) Name() string { return n.name }

func (n *AndNode) GetMark(path *list.List) Mark {
	if n.mark == MarkNone {
		var parentNames []string
		if path != nil {
			parentNames = make([]string, len(n.parents))
		}

		n.mark = MarkPending
		for i, parent := range n.parents {
			switch parent.GetMark(path) {
			case MarkPending, MarkFalse:
				n.mark = MarkNone
				return MarkFalse
			}
			if parentNames != nil {
				parentNames[i] = parent.Name()
			}
		}
		if n.mark == MarkPending {
			n.mark = MarkTrue
		}

		if path != nil && n.mark == MarkTrue {
			if len(parentNames) > 0 {
				path.PushBack(n.name + " <- " + strings.Join(parentNames, ", "))
			} else {
				path.PushBack(n.name)
			}
		}
	}

	return n.mark
}

func (n *AndNode) PeekMark() Mark { return n.mark }

func (n *AndNode) SetMark(m Mark) { n.mark = m }

func (n *AndNode) AddParents(parents ...Node) {
	n.parents = append(n.parents, parents...)
	addChild(n, parents...)
}

func (n *AndNode) ClearParents() {
	removeChild(n, n.parents...)
	n.parents = n.parents[:0]
}

func (n *AndNode) Parents() []Node { return n.parents }

func (n *AndNode) Children() []Node { return n.children }

func (n *AndNode) String() string { return n.name }

// OrNode is satisfied if any of its parents is satisfied, unless it has no
// parents.
type OrNode struct {
	name     string
	mark     Mark
	parents  []Node
	children []Node
}

func NewOrNode(name string) *OrNode {
	return &OrNode{name: name,
		parents: make([]Node, 0), children: make([]Node, 0)}
}

func (n *OrNode) Name() string { return n.name }

func (n *OrNode) GetMark(path *list.List) Mark {
	if n.mark == MarkNone {
		n.mark = MarkPending
		allPending := true
		var parentName string

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.parents {
			switch parent.PeekMark() {
			case MarkTrue:
				n.mark = MarkTrue
				allPending = false
				parentName = parent.Name()
				break OrPeekLoop
			case MarkFalse:
				allPending = false
			}
		}

		// then actually check them otherwise
		if n.mark == MarkPending {
		OrGetLoop:
			for _, parent := range n.parents {
				switch parent.GetMark(path) {
				case MarkTrue:
					n.mark = MarkTrue
					allPending = false
					parentName = parent.Name()
					break OrGetLoop
				case MarkFalse:
					allPending = false
				}
			}
		}

		if (allPending && len(n.parents) > 0) || n.mark == MarkPending {
			n.mark = MarkNone
			return MarkFalse
		}

		if path != nil && n.mark == MarkTrue {
			path.PushBack(fmt.Sprintf("%s <- %s", n.name, parentName))
		}
	}

	return n.mark
}

func (n *OrNode) PeekMark() Mark {
	return n.mark
}

func (n *OrNode) SetMark(m Mark) {
	n.mark = m
}

func (n *OrNode) AddParents(parents ...Node) {
	n.parents = append(n.parents, parents...)
	addChild(n, parents...)
}

func (n *OrNode) ClearParents() {
	removeChild(n, n.parents...)
	n.parents = n.parents[:0]
}

func (n *OrNode) Parents() []Node { return n.parents }

func (n *OrNode) Children() []Node { return n.children }

func (n *OrNode) String() string { return n.name }

// helper functions

func addChild(child Node, parents ...Node) {
	// both types don't work as a single case for whatever reason
	for _, parent := range parents {
		switch nt := parent.(type) {
		case *AndNode:
			nt.children = append(nt.children, child)
		case *OrNode:
			nt.children = append(nt.children, child)
		}
	}
}

func removeChild(child Node, parents ...Node) {
	// same deal as above
	for _, parent := range parents {
		switch nt := parent.(type) {
		case *AndNode:
			removeNodeFromSlice(child, &nt.children)
		case *OrNode:
			removeNodeFromSlice(child, &nt.children)
		}
	}
}

func removeNodeFromSlice(node Node, slice *[]Node) {
	// O(n)
	for i, match := range *slice {
		if match == node {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			break
		}
	}
}

// IsNodeInSlice returns true iff the node is in the slice.
func IsNodeInSlice(node Node, slice []Node) bool {
	for _, match := range slice {
		if node == match {
			return true
		}
	}
	return false
}
