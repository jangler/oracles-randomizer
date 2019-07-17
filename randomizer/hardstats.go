package randomizer

import (
	"fmt"
	"io"
	"os"
	"sort"
)

// these names are for ~4.0.2 and are NOT future-proof.
var seasonsTrickNames = map[string]string{
	"shop, 150 rupees 1 2":                      "shovel manip for 150-rupee item",
	"shop, 20 rupees 1 2":                       "shovel manip for bombs",
	"bomb jump 3 1":                             "bomb jump with seeds",
	"member's shop 1 1 2":                       "shovel manip for member's shop",
	"d8 spike room 1 1":                         "capeless d8 lava sidescroller",
	"d7 pot room 1 2":                           "poe skip",
	"d7 maze chest 1 1":                         "capeless d7",
	"black beast's chest 1 1":                   "black beast torch with mystery seeds",
	"harvest pegasus seeds 1 2":                 "pegasus seeds from market",
	"d1 floormaster room 1 1":                   "d1 east torches with mystery seeds",
	"bomb jump 4 1":                             "bomb jump with cape",
	"d2 arrow room 1 1 1":                       "d2 torches with mystery seeds",
	"blaino prize 1 2":                          "shovel manip for blaino",
	"enter temple remains lower portal 1 2 2 2": "temple remains jump, winter",
	"bomb jump 2 1":                             "bomb jump without seeds or cape",
	"harvest ember seeds 1 2":                   "ember seeds from d5",
	"mt. cucco, platform cave 1 1":              "cucco clip",
	"d8 darknut chest 1 1":                      "HSS-less d8 second triple eyes",
	"harvest ember seeds 1 3":                   "ember seeds from respawnable bushes",
	"d1 goriya chest 1 1":                       "d1 east torches with mystery seeds",
	"eastern suburbs, on cliff 1 1":             "eastern suburbs cliff feather-only",
	"d8 three eyes chest 1 1":                   "HSS-less d8 first triple eyes",
	"punch enemy 1":                             "fist ring as weapon",
	"enter horon village portal 1 1":            "village portal jump",
	"d5 pot room 1 2 1 1":                       "featherless d5 thwomp sidescroller",
	"shop, 30 rupees 1 2":                       "shovel manip for shield",
	"horon village tree 1 1":                    "use starting seeds without tree access",
	"d8 eye drop 1 1":                           "d8 first chest without slingshot",
	"refill seeds 1":                            "use starting seeds without tree access",
	"furnace 1 1 1":                             "subrosia jump to furnace",
	"harvest mystery seeds 1 2":                 "mystery seeds from frypolar",
	"enter temple remains lower portal 1 2 2 4": "temple remains jump, summer",
	"satchel kill normal 1 1":                   "scent or gale seeds from satchel, standard",
	"goron mountain, across pits 1 1":           "goron mountains pit jump",
	"moblin keep 2 3 1 1":                       "natzu wasteland moblin keep jump",
	"beach 3 1 1":                               "subrosia jump to beach",
	"kill armored 1 1 1":                        "scent seeds from satchel, 'armored'",
	"d7 B2F drop 1 1":                           "magnetless d7",
	"d5 basement 1 1":                           "jump through d5 fire trap",
	"d1 block-pushing room 1 1":                 "kill the two d1 goriya with bush",
	"d5 magnet ball chest 1 1":                  "cape-only jump to d5 magnet gloves chest",
	"temple remains lower stump 1 1 2 1 1":      "reverse temple remains jump, summer",
	"d2 blade chest 1 1 1":                      "use pots as weapons in d2",
}

// these names are for ~4.0.2 and are NOT future-proof.
var agesTrickNames = map[string]string{
	"south lynna tree 1 1":         "use starting seeds without tree access",
	"shop, 150 rupees 1 2":         "shovel manip for shop",
	"bombs 2":                      "bombs from head thwomp or shooting gallery",
	"satchel weapon 1 1":           "scent or gale seeds from satchel, standard",
	"ambi's palace chest 1 1":      "guard skip",
	"d4 first crystal switch 1 1":  "d4 first switch with boomerang",
	"d3 N crystal 1 1":             "d3 north crystal with switch hook",
	"d4 second crystal switch 1 1": "d4 second switch with boomerang",
	"kill spiked beetle 1":         "gale seeds from satchel for spiked beetles",
	"d5 eyes chest 1 1":            "shooterless d5 eyes chest",
	"scent seeds 1":                "scent seeds from d3",
	"d3 B1F east 1 1":              "shooterless d3 bk chest",
	"balloon guy 2 1 1 1":          "tingle bridge room with boomerang",
	"d8 ghini chest 1 1":           "d8 torch with mystery seeds",
	"d3 torch chest 1 1":           "d3 torch chest with mystery seeds",
	"ridge upper present 1 1 1":    "d2 skip, but with cane",
	"punch enemy 1":                "fist ring as weapon",
	"bomb jump 3":                  "bomb jump with seeds",
	"patch 1 1":                    "swordless patch",
	"d2 thwomp shelf 1 1":          "d2 thwomp shelf with cane",
	"d3 bridge chest 1 2":          "weird d3 bridge chest, non-backdoor",
	"d2 statue puzzle 1 1":         "d2 moblin door clip",
	"kill swoop 1":                 "scent seeds from satchel for swoop",
	"d3 boss door 1 1 1":           "seedless d3 boss door jump",
	"d6 past entrance 1 2":         "bomb jump into d6 past",
	"d3 boss door 1 2 1":           "d3 boss door without shooter or boomerang",
	"d3 boss 1 1":                  "scent seeds from satchel for shadow hag",
	"d5 crossroads 2 1":            "d5 bridge jump",
	"d5 crossroads 2 2":            "d5 darknut switch manip",
	"d4 color tile drop 1 1":       "scent seeds from satchel for d4 color room",
	"kill moldorm 1":               "scent seeds from satchel for moldorms",
	"bomb jump 2 1 1":              "bomb jump without seeds",
	"d3 bridge chest 1 1":          "weird d3 bridge chest, backdoor",
	"crescent island tree 2 1 1 1": "visit crescent island tree underwater",
	"d4 large floor puzzle 2":      "d4 bomb jump onto bridge",
}

func hasParent(n1, n2 *node) bool {
	for _, p := range n1.parents {
		if p == n2 {
			return true
		}
	}
	return false
}

// returns counts of which tricks are required how often, hopefully using nice
// names for them.
func getHardStats(routes []*routeInfo) map[string]int {
	hardReqs := make(map[string]int)

	for _, r := range routes {
		// create a "null" node that is never true
		hard := r.graph["hard"]
		null := newNode("null", orNode)
		r.graph["null"] = null
		r.graph["start"].explore()

		anything := false // track if any hard trick was required

		for k, v := range r.graph {
			if v.reached && hasParent(v, hard) {
				v.addParent(null)
				r.graph.reset()
				r.graph["start"].explore()

				if !r.graph["done"].reached {
					hardReqs[k]++
					anything = true
				}

				v.removeParent(null)
				r.graph.reset()
				r.graph["start"].explore()
			}
		}

		if anything {
			hardReqs["anything"]++
		}
	}

	return hardReqs
}

// print required hard tricks in descending order of frequency.
func printOrderedHardStats(w io.Writer, counts map[string]int,
	trials int, nameMap map[string]string) {
	sorted := make([]string, 0, len(counts))
	for k := range counts {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return counts[sorted[i]] > counts[sorted[j]]
	})

	for _, k := range sorted {
		fmt.Fprintf(w, "%.1f%%\t%s\n", 100*float64(counts[k])/float64(trials),
			ternary(nameMap[k] != "", nameMap[k], k))
	}
}

// generate a bunch of seeds and print info about frequency of required hard
// logic tricks.
func logHardStats(game, trials int, ropts randomizerOptions, logf logFunc) {
	// get `trials` routes
	routes := generateSeeds(trials, game, ropts)
	nameMap := ternary(game == gameSeasons,
		seasonsTrickNames, agesTrickNames).(map[string]string)
	printOrderedHardStats(os.Stdout, getHardStats(routes), trials, nameMap)
}
