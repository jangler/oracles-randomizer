package main

import (
	"testing"

	"github.com/jangler/oracles-randomizer/graph"
	"github.com/jangler/oracles-randomizer/rom"
)

func TestGraph(t *testing.T) {
	testSeasonsGraph(t)
	// testAgesGraph(t)
}

// check that graph logic is working as expected
func testSeasonsGraph(t *testing.T) {
	rom.Init(nil, rom.GameSeasons)
	r := NewRoute(rom.GameSeasons)
	g := r.Graph

	// test basic start item
	checkReach(t, g, map[string]string{
		"d0 key chest": "feather",
	}, "maku tree", false, false)
	checkReach(t, g, map[string]string{
		"d0 key chest": "sword",
	}, "maku tree", false, true)

	// test hard logic via bombs as weapon
	testMap := map[string]string{
		"d0 key chest":           "moosh's flute",
		"d0 rupee chest":         "bombs, 10",
		"horon village SE chest": "gnarled key",
		"d1 entrance":            "enter d1",
	}
	checkReach(t, g, testMap, "d1 stalfos drop", false, false)
	checkReach(t, g, testMap, "d1 stalfos drop", true, true)

	// test key counting
	testMap = map[string]string{
		"d0 key chest":     "sword",
		"maku tree":        "gnarled key",
		"d1 entrance":      "enter d1",
		"d1 stalfos drop":  "d1 small key",
		"d1 stalfos chest": "bombs, 10",
	}
	checkReach(t, g, testMap, "d1 basement", false, false)
	testMap["d1 railway chest"] = "d1 small key"
	checkReach(t, g, testMap, "d1 basement", false, true)

	// check a subrosia portal
	testMap = map[string]string{
		"d0 key chest":   "sword",
		"d0 rupee chest": "boomerang",
		"maku tree":      "boomerang",
	}
	checkReach(t, g, testMap, "suburbs", false, false)
	testMap["enter horon village portal"] = "exit eastern suburbs portal"
	checkReach(t, g, testMap, "suburbs", false, true)
}

// check that graph logic is working as expected
func testAgesGraph(t *testing.T) {
	rom.Init(nil, rom.GameAges)
	r := NewRoute(rom.GameAges)
	g := r.Graph

	// test basic start item
	checkReach(t, g, map[string]string{
		"starting chest": "feather",
	}, "black tower worker", false, false)
	checkReach(t, g, map[string]string{
		"starting chest": "sword",
	}, "black tower worker", false, true)

	// test hard logic via bombs as weapon
	testMap := map[string]string{
		"starting chest":     "bombs, 10",
		"nayru's house":      "bracelet",
		"black tower worker": "shovel",
		"d2 entrance":        "enter d2",
	}
	checkReach(t, g, testMap, "d2 bombed terrace", false, false)
	checkReach(t, g, testMap, "d2 bombed terrace", true, true)

	// test key counting
	testMap = map[string]string{
		"starting chest":      "sword",
		"nayru's house":       "bombs, 10",
		"black tower worker":  "dimitri's flute",
		"d3 entrance":         "enter d3",
		"d3 pols voice chest": "d3 small key",
		"d3 statue drop":      "d3 small key",
	}
	checkReach(t, g, testMap, "d3 bush beetle room", false, false)
	testMap["d3 armos drop"] = "d3 small key"
	checkReach(t, g, testMap, "d3 bush beetle room", false, true)

	// test bombs from head thwomp in hard logic
	headThwompBombMap := map[string]string{
		"starting chest":        "bracelet",
		"nayru's house":         "harp",
		"black tower worker":    "harp",
		"lynna city chest":      "switch hook",
		"fairies' woods chest":  "iron shield",
		"symmetry city brother": "sword",
		"d2 entrance":           "enter d2",
		"d2 moblin drop":        "feather",
		"d2 basement drop":      "d2 small key",
		"d2 thwomp tunnel":      "d2 small key",
		"d2 thwomp shelf":       "d2 small key",
		"d2 moblin platform":    "d2 small key",
		"d2 rope room":          "d2 small key",
		"d2 statue puzzle":      "d2 boss key",
	}
	checkReach(t, g, headThwompBombMap, "d2 bombed terrace", false, false)
	checkReach(t, g, headThwompBombMap, "d2 bombed terrace", true, true)
}

// helper function for testing whether a node is reachable given a certain
// slotting
func checkReach(t *testing.T, g graph.Graph, links map[string]string,
	target string, hard, expect bool) {
	t.Helper()

	// add parents at the start of the function, and remove them at the end. if
	// a parent is blank, remove it at the start and add it at the end (only
	// useful for default seasons).
	for parent, child := range links {
		if parent == "" {
			g[child].ClearParents()
		} else {
			g[child].AddParents(g[parent])
		}
	}
	defer func() {
		for parent, child := range links {
			if parent == "" {
				g[child].AddParents(g["start"])
			} else {
				g[child].ClearParents()
			}
		}
	}()
	g.ClearMarks()

	if (g[target].GetMark(g[target], hard) == graph.MarkTrue) != expect {
		if expect {
			t.Errorf("expected to reach %s, but could not", target)
		} else {
			t.Errorf("expected not to reach %s, but could", target)
		}
	}
}
