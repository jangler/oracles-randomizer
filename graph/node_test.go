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

func clearMarks(nodes ...Node) {
	for _, n := range nodes {
		n.SetMark(MarkNone)
	}
}

// tests

func TestNodeGetName(t *testing.T) {
	names := []string{"foo", "bar"}
	nodes := []Node{NewAndNode(names[0]), NewOrNode(names[1])}

	for i, node := range nodes {
		if node.GetName() != names[i] {
			t.Errorf("want %s, got %s", names[i], node.GetName())
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

// this is the big one…
func TestNodeGetMark(t *testing.T) {
	and1, or1 := makeAndNode(), makeOrNode()

	// orphan AndNodes are true
	if mark := and1.GetMark(nil); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// orphan OrNodes are false
	if mark := or1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	and2 := makeAndNode()
	and1.AddParents(or1, and2)
	clearMarks(and1, or1)

	// AndNodes need all parents to succeed
	if mark := and1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	or2 := makeOrNode()
	or1.AddParents(and1, or2)
	clearMarks(and1, or1, and2)

	// OrNodes need one
	if mark := or1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.SetMark(MarkNone)
	if mark := or1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	// (clear and re-add w/ true child in front to make sure breaks in switch
	// statements are breaking to loop labels)
	or1.ClearParents()
	or1.AddParents(and2, and1, or2)
	clearMarks(and1, or1, and2, or2)

	// and only one
	if mark := or1.GetMark(nil); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.SetMark(MarkNone)
	if mark := or1.GetMark(nil); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// and now the AndNode should be satisfied
	if mark := and1.GetMark(nil); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}

	// but make sure loops don't satisfy nodes
	and1.ClearParents()
	and2.ClearParents()
	or1.ClearParents()
	or2.ClearParents()
	and1.AddParents(and2)
	and2.AddParents(and1)
	or1.AddParents(or2)
	or2.AddParents(or1)
	clearMarks(and1, and2, or1, or2)
	if mark := and1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
	if mark := or1.GetMark(nil); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
}
