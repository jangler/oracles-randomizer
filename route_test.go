package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/logic"
)

// check that graph logic is working as expected
func TestGraph(t *testing.T) {
	r := NewRoute()
	g := r.Graph

	checkReach(t, g,
		map[string]string{
			"feather 1": "d0 sword chest",
		}, "maku tree gift", false)

	checkReach(t, g,
		map[string]string{
			"winter": "d0 sword chest",
		}, "maku tree gift", true)

	checkReach(t, g,
		map[string]string{
			"sword 1":          "d0 sword chest",
			"ember tree seeds": "ember tree",
			"satchel 1":        "maku tree gift",
			"member's card":    "d0 rupee chest",
		}, "member's shop 1", true)

	checkReach(t, g,
		map[string]string{
			"sword 1":  "d0 sword chest",
			"bracelet": "maku tree gift",
		}, "floodgate key spot", false)

	checkReach(t, g,
		map[string]string{
			"bracelet":         "d0 sword chest",
			"spring":           "village SW chest",
			"flippers":         "maku tree gift",
			"satchel 1":        "platform chest",
			"ember tree seeds": "ember tree",

			"woods of winter default summer": "",
			"woods of winter default winter": "start",
		}, "shovel gift", false)
}

func BenchmarkGraphExplore(b *testing.B) {
	// init graph
	r := NewRoute()
	b.ResetTimer()

	// explore all items from the d0 sword chest
	for name := range logic.ExtraItems() {
		r.Graph.Explore(make(map[*graph.Node]bool), false, r.Graph[name])
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
	g.ExploreFromStart(false)

	if (g[target].GetMark(g[target], false) == graph.MarkTrue) != expect {
		if expect {
			t.Errorf("expected to reach %s, but could not", target)
		} else {
			t.Errorf("expected not to reach %s, but could", target)
		}
	}
}

func TestFindRoute(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	totalRoutes, totalAttempts := 10, 0
	rand.Seed(time.Now().UnixNano())

	logChan, doneChan := make(chan string), make(chan int)
	go func() {
		for {
			select {
			case <-logChan:
			case <-doneChan:
				println("received")
				break
			}
		}
	}()

	for i := 0; i < totalRoutes; i++ {
		println(fmt.Sprintf("finding route %d/%d", i, totalRoutes))
		seed := uint32(rand.Int())
		src := rand.New(rand.NewSource(int64(seed)))
		totalAttempts += findRoute(src, seed, false,
			logChan, doneChan).AttemptCount
	}

	println(fmt.Sprintf("average %d attempts per route",
		totalAttempts/totalRoutes))
}
