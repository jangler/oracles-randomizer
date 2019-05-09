package graph

import (
	"fmt"
	"testing"
)

// helper functions

var andCounter, orCounter int

func newNormalNode(name string, nodeType NodeType) *Node {
	return NewNode(name, nodeType)
}

func makeAndNode() *Node {
	andCounter++
	return newNormalNode(fmt.Sprintf("and%d", andCounter), AndType)
}

func makeOrNode() *Node {
	orCounter++
	return newNormalNode(fmt.Sprintf("or%d", orCounter), OrType)
}

func clearMarks(nodes ...*Node) {
	for _, n := range nodes {
		n.Mark = MarkNone
	}
}

// tests

func TestNodeRelationships(t *testing.T) {
	permutations := [][]func() *Node{
		[]func() *Node{makeAndNode, makeOrNode},
		[]func() *Node{makeOrNode, makeAndNode},
	}

	for _, perm := range permutations {
		n1, n2 := perm[0](), perm[1]()

		// new nodes shouldn't have relationships
		if len(n1.parents) > 0 {
			t.Errorf("node has parents: %+v", n1)
		}
		if len(n1.children) > 0 {
			t.Errorf("node has children: %+v", n1)
		}
		if t.Failed() {
			continue
		}

		// test adding a parent
		n1.AddParents(n2)
		if len(n1.parents) == 0 {
			t.Errorf("node has no parents: %+v", n1)
		}
		if len(n1.children) > 0 {
			t.Errorf("node has children: %+v", n1)
		}
		if len(n2.parents) > 0 {
			t.Errorf("node has parents: %+v", n2)
		}
		if len(n2.children) == 0 {
			t.Errorf("node has no children: %+v", n2)
		}
		if t.Failed() {
			continue
		}

		// test clearing parents
		n1.ClearParents()
		if len(n1.parents) > 0 {
			t.Errorf("node has parents: %+v", n1)
		}
		if len(n2.children) > 0 {
			t.Errorf("node has children: %+v", n2)
		}
	}
}

// this is the big one…
func TestNodeGetMark(t *testing.T) {
	and1, or1 := makeAndNode(), makeOrNode()

	// orphan AndNodes are true
	if mark := and1.GetMark(and1); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// orphan OrNodes are false
	if mark := or1.GetMark(or1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	and2 := makeAndNode()
	and1.AddParents(or1, and2)
	clearMarks(and1, or1)

	// AndNodes need all parents to succeed
	if mark := and1.GetMark(and1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	or2 := makeOrNode()
	or1.AddParents(and1, or2)
	clearMarks(and1, or1, and2)

	// OrNodes need one
	if mark := or1.GetMark(or1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.Mark = MarkNone
	if mark := or1.GetMark(or1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}

	// (clear and re-add w/ true child in front to make sure breaks in switch
	// statements are breaking to loop labels)
	or1.ClearParents()
	or1.AddParents(and2, and1, or2)
	clearMarks(and1, or1, and2, or2)

	// and only one
	if mark := or1.GetMark(or1); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.Mark = MarkNone
	if mark := or1.GetMark(or1); mark != MarkTrue {
		t.Fatalf("want %d, got %d", MarkTrue, mark)
	}
	// and now the AndNode should be satisfied
	if mark := and1.GetMark(and1); mark != MarkTrue {
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
	if mark := and1.GetMark(and1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
	if mark := or1.GetMark(or1); mark != MarkFalse {
		t.Fatalf("want %d, got %d", MarkFalse, mark)
	}
}
