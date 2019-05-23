package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"github.com/jangler/oracles-randomizer/rom"
	"gopkg.in/yaml.v2"
)

// generate a bunch of seeds.
func generateSeeds(n, game int, ropts randomizerOptions) []*RouteInfo {
	threads := runtime.NumCPU()
	dummyLogf := func(string, ...interface{}) {}

	// search for routes
	routeChan := make(chan *RouteInfo)
	attempts := 0
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n/threads; i++ {
				for {
					seed := uint32(rand.Int())
					route, _ := findRoute(game, seed, ropts, false, dummyLogf)
					if route != nil {
						attempts += route.AttemptCount
						routeChan <- route
						break
					}
				}
			}
		}()
	}

	// receive found routes
	routes := make([]*RouteInfo, n/threads*threads)
	for i := 0; i < len(routes); i++ {
		routes[i] = <-routeChan
		fmt.Fprintf(os.Stderr, "%d routes found\n", i+1)
	}
	fmt.Fprintf(os.Stderr, "%.01f%% of seeds succeeded\n",
		100*float64(n)/float64(attempts))

	return routes
}

// generate a bunch of seeds and print item configurations in YAML format.
func logStats(game, trials int, ropts randomizerOptions, logf logFunc) {
	// get `trials` routes
	routes := generateSeeds(trials, game, ropts)

	// make a YAML-serializable slice of check maps
	stringChecks := make([]map[string]string, len(routes))
	for i, ri := range routes {
		stringChecks[i] = make(map[string]string)
		for k, v := range getChecks(ri) {
			stringChecks[i][k.name] = v.name
		}
		if game == rom.GameSeasons {
			for area, seasonId := range ri.Seasons {
				stringChecks[i][area] = seasonsById[int(seasonId)]
			}
		}
		stringChecks[i]["_seed"] = fmt.Sprintf("%08x", ri.Seed)
	}

	// encode to stdout
	if err := yaml.NewEncoder(os.Stdout).Encode(stringChecks); err != nil {
		panic(err)
	}
}
