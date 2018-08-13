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

func BenchmarkGraphExplore(b *testing.B) {
	// init graph
	r := NewRoute([]string{"horon village"})
	b.ResetTimer()

	// explore all items from the d0 sword chest
	for name := range prenode.BaseItems() {
		r.Graph.Explore(
			make(map[*graph.Node]bool), []*graph.Node{r.Graph[name]})
	}
}
