package graph

import (
	"fmt"
	"testing"
)

// helper functions

var andCounter, orCounter int

func makeAndNode() Node {
	andCounter++
	return NewAndNode(fmt.Sprintf("and%d", andCounter))
}

func makeOrNode() Node {
	orCounter++
	return NewOrNode(fmt.Sprintf("or%d", orCounter))
}

// tests

func TestNodeGetName(t *testing.T) {
	names := []string{"foo", "bar"}
	nodes := []Node{NewAndNode(names[0]), NewOrNode(names[1])}

	for i, node := range nodes {
		if node.GetName() != names[i] {
			t.Errorf("want %s, got %s", node.GetName(), names[i])
		}
	}
}

func TestNodeSetMark(t *testing.T) {
	for _, maker := range []func() Node{makeAndNode, makeOrNode} {
		node := maker()
		if node.PeekMark() != MarkNone {
			t.Errorf("want %d, got %d", MarkNone, node.PeekMark())
			continue
		}
		node.SetMark(MarkTrue)
		if node.PeekMark() != MarkTrue {
			t.Errorf("want %d, got %d", MarkTrue, node.PeekMark())
			continue
		}
	}
}

func TestNodeRelationships(t *testing.T) {
	permutations := [][]func() Node{
		[]func() Node{makeAndNode, makeOrNode},
		[]func() Node{makeOrNode, makeAndNode},
	}

	for _, perm := range permutations {
		n1, n2 := perm[0](), perm[1]()

		// new nodes shouldn't have relationships
		if n1.HasParents() {
			t.Errorf("node has parents: %+v", n1)
		}
		if n1.HasChildren() {
			t.Errorf("node has children: %+v", n1)
		}
		if t.Failed() {
			continue
		}

		// test adding a parent
		n1.AddParents(n2)
		if !n1.HasParents() {
			t.Errorf("node has no parents: %+v", n1)
		}
		if n1.HasChildren() {
			t.Errorf("node has children: %+v", n1)
		}
		if n2.HasParents() {
			t.Errorf("node has parents: %+v", n2)
		}
		if !n2.HasChildren() {
			t.Errorf("node has no children: %+v", n2)
		}
		if t.Failed() {
			continue
		}

		// test clearing parents
		n1.ClearParents()
		if n1.HasParents() {
			t.Errorf("node has parents: %+v", n1)
		}
		if n2.HasChildren() {
			t.Errorf("node has children: %+v", n2)
		}
	}
}

// make sure nodes convert to string correctly
func TestNodeString(t *testing.T) {
	andName, orName := "and1", "or1"
	and1, or1 := NewAndNode(andName), NewOrNode(orName)

	if s := and1.String(); s != andName {
		t.Errorf("want %s, got %s", andName, s)
	}
	if s := or1.String(); s != orName {
		t.Errorf("want %s, got %s", orName, s)
	}
}
