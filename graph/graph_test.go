package graph

import "testing"

// tests Graph.Reduce on a graph that is effectively a linked list
func TestListReduce(t *testing.T) {
	g := New()

	// nodes with only one parent can always be collapsed
	a := NewNode("A", OrType, false)
	b := NewNode("B", AndType, false)
	c := NewNode("C", OrType, false)
	d := NewNode("D", RootType, false)
	g.AddNodes(a, b, c, d)

	// so this graph is just |A <- &B <- |C <- .D
	g.AddParents(map[string][]string{
		"A": []string{"B"},
		"B": []string{"C"},
		"C": []string{"D"},
	})

	// and we want it to collapse to |A <- .D
	expectedGraph := New()
	expectedGraph.AddNodes(
		NewNode("A", OrType, false), NewNode("D", RootType, false))
	expectedGraph.AddParents(map[string][]string{"A": []string{"D"}})

	reduced, err := g.Reduce("A")
	if err != nil {
		t.Fatal(err)
	}
	compareGraphs(t, reduced, expectedGraph)
}

// tests Graph.Reduce on a graph without loops
func TestTreeReduce(t *testing.T) {
	given := New()

	// only nodes of the same type should be collapsed
	a := NewNode("A", OrType, false)
	b := NewNode("B", OrType, false)
	c := NewNode("C", AndType, false)
	d := NewNode("D", RootType, false)
	e := NewNode("E", RootType, false)
	f := NewNode("F", OrType, false)
	g := NewNode("G", AndType, false)
	h := NewNode("H", RootType, false)
	i := NewNode("I", RootType, false)
	j := NewNode("J", RootType, false)
	k := NewNode("K", RootType, false)
	given.AddNodes(a, b, c, d, e, f, g, h, i, j, k)

	// this graph is:
	// |A <- |B <- .D
	//          <- .E
	//    <- &C <- |F <- .H
	//                <- .I
	//          <- &G <- .J
	//                <- .K
	given.AddParents(map[string][]string{
		"A": []string{"B", "C"},
		"B": []string{"D", "E"},
		"C": []string{"F", "G"},
		"F": []string{"H", "I"},
		"G": []string{"J", "K"},
	})

	// and we want it to collapse to:
	// |A <- .D
	//    <- .E
	//    <- &C <- |F <- .H
	//                <- .I
	//          <- .J
	//          <- .K
	expected := New()
	expected.AddNodes(
		NewNode("A", OrType, false),
		NewNode("D", RootType, false),
		NewNode("E", RootType, false),
		NewNode("C", AndType, false),
		NewNode("F", OrType, false),
		NewNode("H", RootType, false),
		NewNode("I", RootType, false),
		NewNode("J", RootType, false),
		NewNode("K", RootType, false))
	expected.AddParents(map[string][]string{
		"A": []string{"D", "E", "C"},
		"C": []string{"F", "J", "K"},
		"F": []string{"H", "I"},
	})

	reduced, err := given.Reduce("A")
	if err != nil {
		t.Fatal(err)
	}
	compareGraphs(t, reduced, expected)
}

// tests Graph.Reduce on a graph with loops
func TestGraphReduce(t *testing.T) {
	given := New()

	a := NewNode("A", AndType, false)
	b := NewNode("B", OrType, false)
	c := NewNode("C", AndType, false)
	d := NewNode("D", RootType, false)
	e := NewNode("E", RootType, false)
	f := NewNode("F", RootType, false)
	given.AddNodes(a, b, c, d, e, f)

	// this graph is:
	// &A <- |B <- &C
	//          <- .D
	//          <- .E
	//    <- &C <- .E
	//          <- .F
	given.AddParents(map[string][]string{
		"A": []string{"B", "C"},
		"B": []string{"C", "D", "E"},
		"C": []string{"E", "F"},
	})

	// and we want it to collapse to:
	// &A <- .D
	//    <- .E
	//    <- .F
	expected := New()
	expected.AddNodes(
		NewNode("A", AndType, false),
		NewNode("D", RootType, false),
		NewNode("E", RootType, false),
		NewNode("F", RootType, false))
	expected.AddParents(map[string][]string{
		"A": []string{"D", "E", "F"},
	})

	reduced, err := given.Reduce("A")
	if err != nil {
		t.Fatal(err)
	}
	compareGraphs(t, reduced, expected)
}

// report errors if graphs don't match
func compareGraphs(t *testing.T, given, expected Graph) {
	t.Helper()

	// compare presence of nodes
	for name := range expected {
		if given[name] == nil {
			t.Errorf("node %s missing from graph", name)
		}
	}
	for name := range given {
		if expected[name] == nil {
			t.Errorf("node %s present in graph", name)
		}
	}
	if t.Failed() {
		t.FailNow()
	}

	// compare node relationships
	for name, node := range expected {
		expectedParents, givenParents := node.Parents, given[name].Parents
		if len(expectedParents) == len(givenParents) {
			for _, parent := range expectedParents {
				if !isEquivalentNodeInSlice(parent, givenParents) {
					t.Errorf("expected %s parents %v, given %v",
						name, expectedParents, givenParents)
				}
			}
			for _, parent := range givenParents {
				if !isEquivalentNodeInSlice(parent, expectedParents) {
					t.Errorf("expected %s parents %v, given %v",
						name, expectedParents, givenParents)
				}
			}
		} else {
			t.Errorf("expected %s parents %v, given %v",
				name, expectedParents, givenParents)
		}
	}
}

func isEquivalentNodeInSlice(node *Node, slice []*Node) bool {
	for _, match := range slice {
		if match.Name == node.Name && match.Type == node.Type {
			return true
		}
	}
	return false
}
