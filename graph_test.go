package main

import (
	"fmt"
	"testing"
)

// helper functions

var andCounter, orCounter int

func makeAndNode() *node {
	andCounter++
	return newNode(fmt.Sprintf("and%d", andCounter), andNode)
}

func makeOrNode() *node {
	orCounter++
	return newNode(fmt.Sprintf("or%d", orCounter), orNode)
}

func clearMarks(nodes ...*node) {
	for _, n := range nodes {
		n.mark = markNone
	}
}

// tests

func TestNodeRelationships(t *testing.T) {
	permutations := [][]func() *node{
		[]func() *node{makeAndNode, makeOrNode},
		[]func() *node{makeOrNode, makeAndNode},
	}

	for _, perm := range permutations {
		n1, n2 := perm[0](), perm[1]()

		// new nodes shouldn't have relationships
		if len(n1.parents) > 0 {
			t.Errorf("node has parents: %+v", n1)
		}
		if t.Failed() {
			continue
		}

		// test adding a parent
		n1.addParent(n2)
		if len(n1.parents) == 0 {
			t.Errorf("node has no parents: %+v", n1)
		}
		if len(n2.parents) > 0 {
			t.Errorf("node has parents: %+v", n2)
		}
		if t.Failed() {
			continue
		}

		// test clearing parents
		n1.clearParents()
		if len(n1.parents) > 0 {
			t.Errorf("node has parents: %+v", n1)
		}
	}
}

// this is the big one…
func TestNodeGetMark(t *testing.T) {
	and1, or1 := makeAndNode(), makeOrNode()

	// orphan AndNodes are true
	if mark := and1.getMark(); mark != markTrue {
		t.Fatalf("want %d, got %d", markTrue, mark)
	}
	// orphan OrNodes are false
	if mark := or1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}

	and2 := makeAndNode()
	and1.addParent(or1)
	and1.addParent(and2)
	clearMarks(and1, or1)

	// AndNodes need all parents to succeed
	if mark := and1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}

	or2 := makeOrNode()
	or1.addParent(and1)
	or1.addParent(or2)
	clearMarks(and1, or1, and2)

	// OrNodes need one
	if mark := or1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.mark = markNone
	if mark := or1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}

	// (clear and re-add w/ true child in front to make sure breaks in switch
	// statements are breaking to loop labels)
	or1.clearParents()
	or1.addParent(and2)
	or1.addParent(and1)
	or1.addParent(or2)
	clearMarks(and1, or1, and2, or2)

	// and only one
	if mark := or1.getMark(); mark != markTrue {
		t.Fatalf("want %d, got %d", markTrue, mark)
	}
	// make sure the OrNode gets the same results by peeking
	or1.mark = markNone
	if mark := or1.getMark(); mark != markTrue {
		t.Fatalf("want %d, got %d", markTrue, mark)
	}
	// and now the AndNode should be satisfied
	if mark := and1.getMark(); mark != markTrue {
		t.Fatalf("want %d, got %d", markTrue, mark)
	}

	// but make sure loops don't satisfy nodes
	and1.clearParents()
	and2.clearParents()
	or1.clearParents()
	or2.clearParents()
	and1.addParent(and2)
	and2.addParent(and1)
	or1.addParent(or2)
	or2.addParent(or1)
	clearMarks(and1, and2, or1, or2)
	if mark := and1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}
	if mark := or1.getMark(); mark != markFalse {
		t.Fatalf("want %d, got %d", markFalse, mark)
	}
}
