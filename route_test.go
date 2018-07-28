package main

import (
	"testing"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
)

// these tests actually (mostly) test the graph package, but this package
// combines the graph code with the actual node data. so it's better for more
// realistic benchmarking this way.

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
