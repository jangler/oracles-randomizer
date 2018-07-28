package main

import (
	"math/rand"
	"testing"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
)

func TestShovelLockCheck(t *testing.T) {
	r := NewRoute([]string{"horon village"})
	g := r.Graph

	// make sure that needing a shovel in advance passes
	// this route is via the swamp portal
	g["shovel"].AddParents(g["d0 sword chest"])
	g["bracelet"].AddParents(g["maku key fall"])
	g["flippers"].AddParents(g["blaino gift"])
	g["feather L-1"].AddParents(g["star ore spot"])
	if canShovelSoftlock(g) != nil {
		t.Error("false positive shovel softlock w/ shovel prereq")
	}
	// and make sure the shovel's parents are unchanged
	if len(g["shovel"].Parents) != 1 {
		t.Fatal("shovel parents altered by safety check")
	}

	// make sure that getting there with no shovel fails
	g["shovel"].ClearParents()
	g["bracelet"].ClearParents()
	g["bracelet"].AddParents(g["d0 sword chest"])
	g["feather L-1"].ClearParents()
	g["feather L-1"].AddParents(g["maku key fall"])
	g["sword L-1"].ClearParents()
	g["sword L-1"].AddParents(g["shovel gift"])
	if canShovelSoftlock(g) == nil {
		t.Error("false negative shovel softlock w/ no shovel")
	}

	// make sure that getting a shovel as the gift passes
	g["shovel"].ClearParents()
	g["shovel"].AddParents(g["shovel gift"])
	if canShovelSoftlock(g) != nil {
		t.Error("false positive shovel softlock w/ shovel as gift")
	}

	// and make sure that getting there with an optional shovel fails
	g["shovel"].ClearParents()
	g["shovel"].AddParents(g["boomerang gift"])
	if canShovelSoftlock(g) == nil {
		t.Error("false negative shovel softlock w/ optional shovel")
	}
}

func TestFeatherLockCheck(t *testing.T) {
	r := NewRoute([]string{"horon village"})
	g := r.Graph

	// make sure that it doesn't detect softlock if you can't reach H&S
	g["sword L-1"].AddParents(g["d0 sword chest"])
	g["gnarled key"].AddParents(g["maku key fall"])
	g["satchel"].AddParents(g["d1 satchel"])
	if canFeatherSoftlock(g) != nil {
		t.Error("false positive feather softlock w/o reaching H&S")
	}

	// make sure that it detects softlock if you don't have shovel before H&S
	g["bracelet"].AddParents(g["boomerang gift"])
	g["feather L-2"].AddParents(g["blaino gift"])
	if canFeatherSoftlock(g) == nil {
		t.Error("false negative feather softlock")
	}

	// make sure that it doesn't detect softlock if you must have shovel first
	g["feather L-2"].ClearParents()
	g["shovel"].AddParents(g["blaino gift"])
	g["feather L-2"].AddParents(g["d2 bracelet chest"])
	if canFeatherSoftlock(g) != nil {
		t.Error("false positive feather softlock reaching H&S after shovel")
	}
}

// helper function used for the other benchmarks
func benchGraphCheck(b *testing.B, check func(graph.Graph) error) {
	// make a list of base item nodes to use for testing
	r := NewRoute([]string{"horon village"})
	g := r.Graph
	baseItems := make([]*graph.Node, 0, len(prenode.BaseItems()))
	for name := range prenode.BaseItems() {
		baseItems = append(baseItems, g[name])
	}

	for i := 0; i < b.N; i++ {
		// create a fresh graph and shuffle the item list
		b.StopTimer()
		r = NewRoute([]string{"horon village"})
		g = r.Graph
		reached := map[*graph.Node]bool{g["horon village"]: true}

		rand.Shuffle(len(baseItems), func(i, j int) {
			baseItems[i], baseItems[j] = baseItems[j], baseItems[i]
		})
		b.StartTimer()

		// gradually add items to the graph to get a picture of performance at
		// various stages in the exploration
		for _, itemNode := range baseItems {
			itemNode.AddParents(g["d0 sword chest"])
			reached = g.Explore(reached, []*graph.Node{itemNode})

			// run 10 times to get a better proportion of check runtime vs
			// explore runtime. just ignoring the explore runtime results in
			// really long tests
			for j := 0; j < 10; j++ {
				check(g)
			}
		}
	}
}

func BenchmarkCanSoftlock(b *testing.B) {
	benchGraphCheck(b, canSoftlock)
}

func BenchmarkCanShovelSoftlock(b *testing.B) {
	benchGraphCheck(b, canShovelSoftlock)
}

func BenchmarkCanFlowerSoftlock(b *testing.B) {
	benchGraphCheck(b, canFlowerSoftlock)
}

func BenchmarkCanFeatherSoftlock(b *testing.B) {
	benchGraphCheck(b, canFeatherSoftlock)
}
