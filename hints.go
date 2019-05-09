package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/jangler/oracles-randomizer/graph"
	"gopkg.in/yaml.v2"
)

// names corresponding to the constant indices from the rom package.
var gameNames = []string{"", "ages", "seasons"}

type hinter struct {
	areas map[string]string
	items map[string]string
}

// returns a new hinter initialized for the given game.
func newHinter(game int) *hinter {
	h := &hinter{
		areas: make(map[string]string),
		items: make(map[string]string),
	}

	// load item names
	itemFiles := []string{
		"/hints/common_items.yaml",
		fmt.Sprintf("/hints/%s_items.yaml", gameNames[game]),
	}
	for _, filename := range itemFiles {
		if err := yaml.Unmarshal(
			FSMustByte(false, filename), h.items); err != nil {
			panic(err)
		}
	}

	// load area names
	rawAreas := make(map[string][]string)
	areasFilename := fmt.Sprintf("/hints/%s_areas.yaml", gameNames[game])
	if err := yaml.Unmarshal(
		FSMustByte(false, areasFilename), rawAreas); err != nil {
		panic(err)
	}

	// transform the areas map from: {final: [internal 1, internal 2]}
	// to: {internal 1: final, internal 2: final}
	for k, a := range rawAreas {
		for _, v := range a {
			h.areas[v] = k
		}
	}

	return h
}

// returns a randomly generated map of owl names to owl messages.
func (h *hinter) generate(src *rand.Rand, g graph.Graph,
	checks map[*graph.Node]*graph.Node, owlNames []string) map[string]string {
	// function body starts here lol
	hints := make(map[string]string)
	slots := getShuffledSlots(src, checks)
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
			required := g[owlName].GetMark() == graph.MarkFalse
			item.AddParent(slot)

			if !required {
				hints[owlName] = h.format(slot, item)
				hintedSlots[slot] = true
				break
			}
		}
	}

	return hints
}

// returns a message stating that an item is in an area, formatted for an owl
// text box. this doesn't include control characters (except newlines).
func (h *hinter) format(slot, item *graph.Node) string {
	// split message into words to be wrapped
	words := strings.Split(h.areas[slot.Name], " ")
	words = append(words, "holds")
	words = append(words, strings.Split(h.items[item.Name], " ")...)
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

// getShuffledSlots returns a randomly ordered slice of slot nodes.
func getShuffledSlots(src *rand.Rand,
	checks map[*graph.Node]*graph.Node) []*graph.Node {
	// make slice of check names
	slots := make([]*graph.Node, len(checks))
	i := 0
	for slot, item := range checks {
		// don't include dungeon items, since dungeon item hints would be
		// useless ("Level 7 holds a Boss Key")
		if item.Name == "dungeon map" ||
			item.Name == "compass" ||
			strings.HasPrefix(item.Name, "slate") ||
			strings.HasSuffix(item.Name, "small key") ||
			strings.HasSuffix(item.Name, "boss key") {
			continue
		}

		// and don't include these checks, since they're dummy slots that
		// aren't actually randomized, or seed trees that the player is
		// guaranteed to know about if they're using seeds.
		switch slot.Name {
		case "shop, 20 rupees", "shop, 30 rupees",
			"horon village seed tree", "south lynna tree":
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
