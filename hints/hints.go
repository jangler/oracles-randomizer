package hints

import (
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
	game int, hard bool) map[string]string {
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
			// in the first place, as dictated by the logic of the seed.
			item.RemoveParent(slot)
			g.ClearMarks()
			required := g[owlName].GetMark(g[owlName], hard) == graph.MarkFalse
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
	slots := make([]*graph.Node, len(checks))
	i := 0
	for slot, item := range checks {
		// don't include dungeon items, since dungeon item hints would be
		// useless ("Level 7 holds a Boss Key")
		if item.Name == "dungeon map" ||
			item.Name == "compass" ||
			item.Name == "slate" ||
			strings.HasSuffix(item.Name, "boss key") {
			continue
		}

		// and don't include these checks, since they're dummy slots that
		// aren't actually randomized.
		switch slot.Name {
		case "shop, 20 rupees", "shop, 30 rupees":
			continue
		}

		slots[i] = slot
		i++
	}
	slots = slots[:i] // trim to the number of actually hintable checks

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
// text box. this doesn't include control characters (except newlines).
func formatMessage(slot, item *graph.Node, game int) string {
	var areaMap map[string]string
	if game == GameSeasons {
		areaMap = seasonsAreaMap
	} else {
		areaMap = agesAreaMap
	}

	// split message into words to be wrapped
	words := strings.Split(areaMap[slot.Name], " ")
	words = append(words, "holds")
	if itemMap[item.Name].article != "" {
		words = append(words, itemMap[item.Name].article)
	}
	words = append(words, strings.Split(itemMap[item.Name].name, " ")...)
	words[len(words)-1] = words[len(words)-1] + "."

	// build message line by line
	msg := new(strings.Builder)
	line := ""
	for _, word := range words {
		if len(line) == 0 {
			line += word
		} else if len(line)+len(word) <= 15 {
			line += " " + word
		} else {
			msg.WriteString(line + "\n")
			line = word
		}
	}
	msg.WriteString(line)

	return msg.String()
}
