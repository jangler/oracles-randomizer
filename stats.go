package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
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

// generate a bunch of seeds and print item configurations in YAML format.
func logStats(game, trials int, hard bool, logf logFunc) {
	// get `trials` routes
	routes := generateSeeds(trials, game, hard)

	// make a YAML-serializable slice of check maps
	stringChecks := make([]map[string]string, len(routes))
	for i, ri := range routes {
		stringChecks[i] = make(map[string]string)
		for k, v := range getChecks(ri) {
			stringChecks[i][k.Name] = v.Name
		}
	}

	// encode to stdout
	if err := yaml.NewEncoder(os.Stdout).Encode(stringChecks); err != nil {
		panic(err)
	}
}
