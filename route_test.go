package main

import (
	"testing"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
)

// make sure the route's "normal" and "hard" graphs are behaving appropriately
func TestNormalVsHard(t *testing.T) {
	r := NewRoute([]string{"horon village"})

	// references var in safety_test.go
	for child, parent := range testData2 {
		if parent == "" {
			r.ClearParents(child)
		} else {
			r.AddParent(child, parent)
		}
	}

	// make sure at least root nodes are identical
	for name := range testData2 {
		node := r.Graph[name]
		for i, parent := range node.Parents {
			if r.HardGraph[name].Parents[i].Name != parent.Name {
				t.Errorf("parent mismatch: %s (normal) vs %s (hard)",
					parent.Name, r.HardGraph[name].Parents[i].Name)
			}
		}
	}
}

// check that graph logic is working as expected
func TestGraph(t *testing.T) {
	r := NewRoute([]string{"horon village"})
	g := r.Graph

	checkReach(t, g,
		map[string]string{
			"feather L-2": "d0 sword chest",
		}, "maku tree gift", false)

	checkReach(t, g,
		map[string]string{
			"feather L-2": "d0 sword chest",
		}, "lake chest", true)

	checkReach(t, g,
		map[string]string{
			"winter": "d0 sword chest",
		}, "maku tree gift", true)

	checkReach(t, g,
		map[string]string{
			"feather L-2": "d0 sword chest",
			"winter":      "lake chest",
		}, "maku tree gift", true)

	checkReach(t, g,
		map[string]string{
			"boomerang L-2": "d0 sword chest",
		}, "maku tree gift", false)

	checkReach(t, g,
		map[string]string{
			"boomerang L-2": "d0 sword chest",
			"rupees, 20":    "d0 rupee chest",
		}, "maku tree gift", true)

	checkReach(t, g,
		map[string]string{
			"sword L-1":        "d0 sword chest",
			"ember tree seeds": "ember tree",
			"satchel 1":        "maku tree gift",
			"member's card":    "village SE chest",
		}, "member's shop 1", true)
}

func BenchmarkGraphExplore(b *testing.B) {
	// init graph
	r := NewRoute([]string{"horon village"})
	b.ResetTimer()

	// explore all items from the d0 sword chest
	for name := range prenode.ExtraItems() {
		r.Graph.Explore(
			make(map[*graph.Node]bool), []*graph.Node{r.Graph[name]})
	}
}

// helper function for testing whether a node is reachable given a certain
// slotting
//
// TODO refactor this and checkSoftlockWithSlots, since they share most of
//      their code
func checkReach(t *testing.T, g graph.Graph, parents map[string]string,
	target string, expect bool) {
	t.Helper()

	// add parents at the start of the function, and remove them at the end. if
	// a parent is blank, remove it at the start and add it at the end (only
	// useful for default seasons).
	for child, parent := range parents {
		if parent == "" {
			g[child].ClearParents()
		} else {
			g[child].AddParents(g[parent])
		}
	}
	defer func() {
		for child, parent := range parents {
			if parent == "" {
				g[child].AddParents(g["start"])
			} else {
				g[child].ClearParents()
			}
		}
	}()
	g.ExploreFromStart()

	if (g[target].GetMark(g[target], nil) == graph.MarkTrue) != expect {
		if expect {
			t.Errorf("expected to reach %s, but could not", target)
		} else {
			t.Errorf("expected not to reach %s, but could", target)
		}
	}
}
