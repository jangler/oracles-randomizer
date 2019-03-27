package hints

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/jangler/oracles-randomizer/graph"
)

const (
	GameNil = iota
	GameAges
	GameSeasons
)

// Generate returns a randomly generated map of owl names to owl messages.
func Generate(src *rand.Rand, g graph.Graph,
	checks map[*graph.Node]*graph.Node, owlNames []string,
	game int) map[string]string {
	// function body starts here lol
	hints := make(map[string]string)
	slots := getOrderedSlots(src, checks)
	i := 0

	// keep track of which slots have been hinted at in order to avoid
	// duplicates. in practice the implementation of the hint loop makes this
	// very unlikely in the first place.
	hintedSlots := make(map[*graph.Node]bool)

	for _, owlName := range owlNames {
		for {
			slot, item := slots[i], checks[slots[i]]
			i = (i + 1) % len(slots)

			if hintedSlots[slot] {
				continue
			}

			// don't give hints about checks that are required to reach the owl
			// in the first place, *as dictated by hard logic*.
			item.RemoveParent(slot)
			g.ClearMarks()
			required := g[owlName].GetMark(g[owlName], true) == graph.MarkFalse
			item.AddParents(slot)

			if !required {
				hints[owlName] = formatMessage(slot, item, game)
				hintedSlots[slot] = true
				break
			}
		}
	}

	return hints
}

// getOrderedSlots returns a randomly ordered slice of slot nodes.
func getOrderedSlots(src *rand.Rand,
	checks map[*graph.Node]*graph.Node) []*graph.Node {
	// make slice of check names
	slots := make([]*graph.Node, len(checks)-8*3)
	i := 0
	for slot, item := range checks {
		// don't include dungeon items, since dungeon item hints would be
		// useless ("Level 7 holds a Boss Key")
		if item.Name == "dungeon map" ||
			item.Name == "compass" ||
			strings.HasSuffix(item.Name, "boss key") {
			continue
		}

		slots[i] = slot
		i++
	}

	// sort the slots before shuffling to get "deterministic" results
	sort.Sort(nodeSlice(slots))
	src.Shuffle(len(slots), func(i, j int) {
		slots[i], slots[j] = slots[j], slots[i]
	})

	return slots
}

// implement sort.Interface for []*graph.Node
type nodeSlice []*graph.Node

func (ns nodeSlice) Len() int {
	return len(ns)
}

func (ns nodeSlice) Less(i, j int) bool {
	return ns[i].Name < ns[j].Name
}

func (ns nodeSlice) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

// returns a message stating that an item is in an area, formatted for an owl
// text box. this doesn't include control characters.
func formatMessage(slot, item *graph.Node, game int) string {
	var areaMap map[string]string
	if game == GameSeasons {
		areaMap = seasonsAreaMap
	} else {
		areaMap = agesAreaMap
	}

	return fmt.Sprintf("%s\nholds %s\n%s.", areaMap[slot.Name],
		itemMap[item.Name].article, itemMap[item.Name].name)
}
