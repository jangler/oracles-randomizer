package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/jangler/oos-randomizer/graph"
)

// generate a bunch of seeds.
func generateSeeds(n, game int, hard bool) []*RouteInfo {
	threads := runtime.NumCPU()

	// ignore messages from route searches
	logChan := make(chan string)
	doneChan := make(chan int)
	go func() {
		for {
			select {
			case <-logChan:
			case <-doneChan:
				break
			}
		}
	}()

	// search for routes
	routeChan := make(chan *RouteInfo)
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n/threads; i++ {
				seed := uint32(rand.Int())
				routeChan <- findRoute(game, seed, hard, false,
					logChan, doneChan)
			}
		}()
	}

	// receive found routes
	routes := make([]*RouteInfo, n/threads*threads)
	for i := 0; i < len(routes); i++ {
		routes[i] = <-routeChan
		fmt.Fprintf(os.Stderr, "%d routes found\n", i+1)
	}

	return routes
}

// generate a bunch of seeds and print information about how often items are
// required, and what spheres they're normally in.
func logStats(game int, hard bool, logf func(string, ...interface{})) {
	routes := generateSeeds(1000, game, hard)

	// aggregate data on required items
	requiredCounts := make(map[string]int)
	for _, route := range routes {
		for e := route.ProgressItems.Front(); e != nil; e = e.Next() {
			itemName := e.Value.(*graph.Node).Name
			if itemName == "satchel 1" || itemName == "satchel 2" {
				itemName = "seed satchel"
			}
			if requiredCounts[itemName] == 0 {
				requiredCounts[itemName] = 0
			}
			requiredCounts[itemName]++
		}
	}

	for item, count := range requiredCounts {
		logf("%5.1f%% - %s", 100*float64(count)/float64(len(routes)),
			getNiceName(item))
	}
}
