package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
)

// generate a bunch of seeds.
func generateSeeds(n, game int, hard bool) []*RouteInfo {
	threads := runtime.NumCPU()
	dummyLogf := func(string, ...interface{}) {}

	// search for routes
	routeChan := make(chan *RouteInfo)
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n/threads; i++ {
				seed := uint32(rand.Int())
				routeChan <- findRoute(game, seed, hard, false, dummyLogf)
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
func logStats(game, trials int, hard bool, logf logFunc) {
	routes := generateSeeds(trials, game, hard)

	// aggregate data on required items
	meanSpheres := make(map[string]float64)
	for _, ri := range routes {
		// total spheres
		checks := getChecks(ri)
		spheres := getSpheres(ri.Route.Graph, checks, hard)
		for i, sphere := range spheres {
			for _, node := range sphere {
				if !node.IsStep {
					continue
				}
				if meanSpheres[node.Name] == 0 {
					meanSpheres[node.Name] = 0
				}
				meanSpheres[node.Name] += float64(i)
			}
		}
	}

	for item, totalSpheres := range meanSpheres {
		logf("%s - %4.1f", getNiceName(item),
			float64(totalSpheres)/float64(len(routes)))
	}
}
