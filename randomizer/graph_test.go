package randomizer

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
	checkMark(t, and1, markTrue)
	// orphan OrNodes are false
	checkMark(t, or1, markFalse)

	and2 := makeAndNode()
	and1.addParent(or1)
	and1.addParent(and2)
	clearMarks(and1, or1)

	// AndNodes need all parents to succeed
	checkMark(t, and1, markFalse)

	or2 := makeOrNode()
	or1.addParent(and1)
	or1.addParent(or2)
	clearMarks(and1, or1, and2)

	// OrNodes need one
	checkMark(t, or1, markFalse)
	// make sure the OrNode gets the same results by peeking
	or1.mark = markNone
	checkMark(t, or1, markFalse)

	// (clear and re-add w/ true child in front to make sure breaks in switch
	// statements are breaking to loop labels)
	or1.clearParents()
	or1.addParent(and2)
	or1.addParent(and1)
	or1.addParent(or2)
	clearMarks(and1, or1, and2, or2)

	// and only one
	checkMark(t, or1, markTrue)
	// make sure the OrNode gets the same results by peeking
	or1.mark = markNone
	checkMark(t, or1, markTrue)
	// and now the AndNode should be satisfied
	checkMark(t, and1, markTrue)

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
	checkMark(t, and1, markFalse)
	checkMark(t, or1, markFalse)
}

func TestCountNodes(t *testing.T) {
	count := newNode("count", countNode)
	count.minCount = 2
	child := newNode("child", andNode)
	parent := newNode("parent", andNode)
	count.addParent(child)

	// if child has only one parent, count should be 1 (< 2)
	child.addParent(parent)
	checkMark(t, count, markFalse)

	// two parents should suffice
	child.addParent(parent)
	checkMark(t, count, markTrue)
}

func TestNegatedNodes(t *testing.T) {
	and := newNode("tn", andNode)
	or := newNode("fn", orNode)
	not := newNode("not", nandNode)
	nor := newNode("nor", norNode)

	not.addParent(or)
	checkMark(t, not, markTrue)
	nor.addParent(or)
	checkMark(t, nor, markTrue)
	not.addParent(and)
	checkMark(t, not, markTrue)
	nor.addParent(and)
	checkMark(t, nor, markFalse)
}

func TestEitherNodes(t *testing.T) {
	either := newNode("either", eitherNode)
	start := newNode("start", andNode)
	root := newNode("root", orNode)

	and := newNode("and", andNode)
	and.addParent(either)
	checkMark(t, and, markEither)
	and.addParent(start)
	checkMark(t, and, markEither)
	and.addParent(root)
	checkMark(t, and, markFalse)

	or := newNode("or", orNode)
	or.addParent(either)
	checkMark(t, or, markEither)
	or.addParent(root)
	checkMark(t, or, markEither)
	or.addParent(start)
	checkMark(t, or, markTrue)

	count := newNode("count", countNode)
	count.minCount = 2
	counted := newNode("counted", orNode)
	count.addParent(counted)
	counted.addParent(either)
	checkMark(t, count, markFalse)
	counted.addParent(start)
	checkMark(t, count, markEither)
	counted.addParent(or)
	checkMark(t, count, markTrue)

	nand := newNode("nand", nandNode)
	nand.addParent(either)
	checkMark(t, nand, markEither)
	nand.addParent(start)
	checkMark(t, nand, markEither)
	nand.addParent(root)
	checkMark(t, nand, markTrue)
}

func checkMark(t *testing.T, n *node, expected nodeMark) {
	t.Helper()
	n.mark = markNone
	if actual := n.getMark(false); expected != actual {
		t.Fatalf("want %s, got %s", expected, actual)
	}
}
