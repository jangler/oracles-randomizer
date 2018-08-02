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
	var softlock error

	// make sure that getting there with a shovel does not error
	softlock = checkSoftlockWithSlots(t, canShovelSoftlock, g,
		map[string]string{
			"sword L-1":        "d0 sword chest",
			"satchel":          "maku key fall",
			"ember tree seeds": "ember tree",
			"feather L-1":      "boomerang gift",
			"shovel":           "rod gift",
			"rod":              "star ore spot",
			"fool's ore":       "shovel gift",
		})
	if g["shovel gift"].GetMark(g["shovel gift"], nil) != graph.MarkTrue {
		t.Error("test invalid: cannot reach shovel gift")
	}
	if softlock != nil {
		t.Error("false positive shovel softlock w/ shovel prereq")
	}

	// make sure that getting a shovel there does not error
	softlock = checkSoftlockWithSlots(t, canShovelSoftlock, g,
		map[string]string{
			"sword L-1":        "d0 sword chest",
			"satchel":          "maku key fall",
			"ember tree seeds": "ember tree",
			"rod":              "rod gift",
			"shovel":           "shovel gift",
		})
	if g["shovel gift"].GetMark(g["shovel gift"], nil) != graph.MarkTrue {
		t.Error("test invalid: cannot reach shovel gift")
	}
	if softlock != nil {
		t.Error("false positive shovel softlock w/ shovel as gift")
	}

	// and make sure that getting there with an optional shovel errors
	softlock = checkSoftlockWithSlots(t, canShovelSoftlock, g,
		map[string]string{
			"sword L-1":        "d0 sword chest",
			"satchel":          "maku key fall",
			"ember tree seeds": "ember tree",
			"rod":              "rod gift",
			"fool's ore":       "shovel gift",
			"shovel":           "boomerang gift",
		})
	if g["shovel gift"].GetMark(g["shovel gift"], nil) != graph.MarkTrue {
		t.Error("test invalid: cannot reach shovel gift")
	}
	if softlock == nil {
		t.Error("false negative shovel softlock w/ optional shovel")
	}
}

func TestFeatherLockCheck(t *testing.T) {
	r := NewRoute([]string{"horon village"})
	g := r.Graph
	var softlock error

	// make sure reaching H&S with mandatory shovel does not error
	softlock = checkSoftlockWithSlots(t, canFeatherSoftlock, g,
		map[string]string{
			"bracelet":           "d0 sword chest",
			"flippers":           "maku key fall",
			"shovel":             "blaino gift",
			"feather L-2":        "star ore spot",
			"rod":                "rod gift",
			"satchel":            "shovel gift",
			"pegasus tree seeds": "ember tree",
		})
	if softlock != nil {
		t.Error("false positive feather softlock w/ shovel prereq")
	}

	// make sure reaching H&S with optional shovel errors
	softlock = checkSoftlockWithSlots(t, canFeatherSoftlock, g,
		map[string]string{
			"bracelet":           "d0 sword chest",
			"flippers":           "maku key fall",
			"feather L-2":        "blaino gift",
			"rod":                "rod gift",
			"pegasus tree seeds": "ember tree",
			"shovel":             "boomerang gift",
		})
	if softlock == nil {
		t.Error("false negative feather softlock w/ optional shovel")
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

func BenchmarkCanEmberSeedSoftlock(b *testing.B) {
	benchGraphCheck(b, canEmberSeedSoftlock)
}

func BenchmarkCanPiratesBellSoftlock(b *testing.B) {
	benchGraphCheck(b, canPiratesBellSoftlock)
}

// the keys in the "parents" map MUST be root nodes. this function errors if
func checkSoftlockWithSlots(t *testing.T, check func(g graph.Graph) error,
	g graph.Graph, parents map[string]string) error {
	t.Helper()

	// add parents at the start of the function, and remove them at the end
	for child, parent := range parents {
		g[child].AddParents(g[parent])
	}
	defer func() {
		for child := range parents {
			g[child].ClearParents()
		}
	}()
	g.ExploreFromStart()

	return check(g)
}
