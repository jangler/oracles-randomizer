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
)

// Node is the general interface that encompasses everything in the graph.
type Node interface {
	fmt.Stringer

	GetName() string
	GetMark(*list.List) Mark // list to append path to if non-nil
	PeekMark() Mark          // like GetMark but doesn't check parents
	SetMark(Mark)
	AddParents(...Node)
	ClearParents()
	HasParents() bool
	HasChildren() bool
}

// AndNode is satisfied if all of its parents are satisfied, or if it has no
// parents.
type AndNode struct {
	Name     string
	Mark     Mark
	Parents  []Node
	Children []Node
}

func NewAndNode(name string) *AndNode {
	return &AndNode{Name: name,
		Parents: make([]Node, 0), Children: make([]Node, 0)}
}

func (n *AndNode) GetName() string { return n.Name }

func (n *AndNode) GetMark(path *list.List) Mark {
	if n.Mark == MarkNone {
		var parentNames []string
		if path != nil {
			parentNames = make([]string, len(n.Parents))
		}

		n.Mark = MarkPending
	AndLoop:
		for i, parent := range n.Parents {
			switch parent.GetMark(path) {
			case MarkFalse:
				n.Mark = MarkFalse
				break AndLoop
			// if we encounter a pending node, this node isn't satisfied now,
			// but it could be in the future of the same graph.
			case MarkPending:
				n.Mark = MarkNone
				return MarkFalse
			}
			if parentNames != nil {
				parentNames[i] = parent.GetName()
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

func (n *AndNode) PeekMark() Mark {
	return n.Mark
}

func (n *AndNode) SetMark(m Mark) {
	n.Mark = m
}

func (n *AndNode) AddParents(parents ...Node) {
	n.Parents = append(n.Parents, parents...)
	addChild(n, parents...)
}

func (n *AndNode) ClearParents() {
	removeChild(n, n.Parents...)
	n.Parents = n.Parents[:0]
}

func (n *AndNode) HasParents() bool {
	return len(n.Parents) > 0
}

func (n *AndNode) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *AndNode) String() string {
	return n.Name
}

// OrNode is satisfied if any of its parents is satisfied, unless it has no
// parents.
type OrNode struct {
	Name     string
	Mark     Mark
	Parents  []Node
	Children []Node
}

func NewOrNode(name string) *OrNode {
	return &OrNode{Name: name,
		Parents: make([]Node, 0), Children: make([]Node, 0)}
}

func (n *OrNode) GetName() string { return n.Name }

func (n *OrNode) GetMark(path *list.List) Mark {
	if n.Mark == MarkNone {
		n.Mark = MarkPending
		allPending := true
		var parentName string

		// prioritize already satisfied nodes
	OrPeekLoop:
		for _, parent := range n.Parents {
			switch parent.PeekMark() {
			case MarkTrue:
				n.Mark = MarkTrue
				allPending = false
				parentName = parent.GetName()
				break OrPeekLoop
			case MarkFalse:
				allPending = false
			}
		}

		// then actually check them otherwise
		if n.Mark == MarkPending {
		OrGetLoop:
			for _, parent := range n.Parents {
				switch parent.GetMark(path) {
				case MarkTrue:
					n.Mark = MarkTrue
					allPending = false
					parentName = parent.GetName()
					break OrGetLoop
				case MarkFalse:
					allPending = false
				}
			}
		}

		if allPending && len(n.Parents) > 0 {
			// if everything else is pending; don't give up Forever
			n.Mark = MarkNone
			return MarkFalse
		} else if n.Mark == MarkPending {
			n.Mark = MarkFalse
		}

		if path != nil && n.Mark == MarkTrue {
			path.PushBack(fmt.Sprintf("%s <- %s", n.Name, parentName))
		}
	}

	return n.Mark
}

func (n *OrNode) PeekMark() Mark {
	return n.Mark
}

func (n *OrNode) SetMark(m Mark) {
	n.Mark = m
}

func (n *OrNode) AddParents(parents ...Node) {
	n.Parents = append(n.Parents, parents...)
	addChild(n, parents...)
}

func (n *OrNode) ClearParents() {
	removeChild(n, n.Parents...)
	n.Parents = n.Parents[:0]
}

func (n *OrNode) HasParents() bool {
	return len(n.Parents) > 0
}

func (n *OrNode) HasChildren() bool {
	return len(n.Children) > 0
}

func (n *OrNode) String() string {
	return n.Name
}

// helper functions

func addChild(child Node, parents ...Node) {
	// both types don't work as a single case for whatever reason
	for _, parent := range parents {
		switch nt := parent.(type) {
		case *AndNode:
			nt.Children = append(nt.Children, child)
		case *OrNode:
			nt.Children = append(nt.Children, child)
		}
	}
}

func removeChild(child Node, parents ...Node) {
	// same deal as above
	for _, parent := range parents {
		switch nt := parent.(type) {
		case *AndNode:
			removeNodeFromSlice(child, &nt.Children)
		case *OrNode:
			removeNodeFromSlice(child, &nt.Children)
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
