package randomizer

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"

	"gopkg.in/yaml.v2"
)

// generate a bunch of seeds.
func generateSeeds(n, game int, gopts *globalOptions) []*routeInfo {
	threads := runtime.NumCPU()
	dummyLogf := func(string, ...interface{}) {}

	// search for routes
	routeChan := make(chan *routeInfo)
	attempts := 0
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n/threads; i++ {
				for {
					// i don't know if a new romState actually *needs* to be
					// created for each iteration.
					seed := uint32(rand.Int())
					roms := []*romState{
						newRomState(nil, game, 1, gopts.include),
					}
					routes, _ := findRoutes(roms, seed, gopts, false, dummyLogf)
					if routes[0] != nil {
						attempts += routes[0].attemptCount
						routeChan <- routes[0]
						break
					}
				}
			}
		}()
	}

	// receive found routes
	routes := make([]*routeInfo, n/threads*threads)
	for i := 0; i < len(routes); i++ {
		routes[i] = <-routeChan
		fmt.Fprintf(os.Stderr, "%d routes found\n", i+1)
	}
	fmt.Fprintf(os.Stderr, "%.01f%% of seeds succeeded\n",
		100*float64(n)/float64(attempts))

	return routes
}

// generate a bunch of seeds and print item configurations in YAML format.
func logStats(game, trials int, gopts *globalOptions, logf logFunc) {
	// get `trials` routes
	routes := generateSeeds(trials, game, gopts)

	// make a YAML-serializable slice of check maps
	stringChecks := make([]map[string]string, len(routes))
	for i, ri := range routes {
		stringChecks[i] = make(map[string]string)
		for k, v := range getChecks(ri.usedItems, ri.usedSlots) {
			stringChecks[i][k.name] = v.name
		}
		if game == gameSeasons {
			for area, seasonId := range ri.seasons {
				// make sure not to overwrite info about lost woods item
				if area == "lost woods" {
					area = "lost woods (season)"
				}
				stringChecks[i][area] = seasonsById[int(seasonId)]
			}
		}
		stringChecks[i]["_seed"] = fmt.Sprintf("%08x", ri.seed)
	}

	// encode to stdout
	if err := yaml.NewEncoder(os.Stdout).Encode(stringChecks); err != nil {
		panic(err)
	}
}
