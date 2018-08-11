package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
	"github.com/jangler/oos-randomizer/rom"
)

const (
	maxIterations = 200 // restart if routing runs for too long
	maxTries      = 50  // give up if routing fails too many times
)

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph graph.Graph
	Slots map[string]*graph.Node
}

// NewRoute returns an initialized route with all prenodes, and those prenodes
// with the names in start functioning as givens (always satisfied).
func NewRoute(start []string) *Route {
	g := graph.New()

	totalPrenodes := prenode.GetAll()
	expandNestedPrenodes(totalPrenodes)

	// make start nodes given
	for _, key := range start {
		totalPrenodes[key] = prenode.And()
	}

	addNodes(g, totalPrenodes)
	addNodeParents(g, totalPrenodes)

	openSlots := make(map[string]*graph.Node, 0)
	for name, pn := range totalPrenodes {
		switch pn.Type {
		case prenode.AndSlotType, prenode.OrSlotType:
			openSlots[name] = g[name]
		}
	}

	return &Route{Graph: g, Slots: openSlots}
}

// flatten nested prenodes
func expandNestedPrenodes(prenodes map[string]*prenode.Prenode) {
	done := true

	for name, pn := range prenodes {
		subID := 0
		for i, parent := range pn.Parents {
			switch parent := parent.(type) {
			case *prenode.Prenode:
				subID++
				subName := fmt.Sprintf("%s %d", name, subID)
				pn.Parents[i] = subName
				prenodes[subName] = parent
				done = false
			}
		}
	}

	// recurse if prenodes were added
	if !done {
		expandNestedPrenodes(prenodes)
	}
}

func addNodes(g graph.Graph, prenodes map[string]*prenode.Prenode) {
	for key, pn := range prenodes {
		switch pn.Type {
		case prenode.AndType, prenode.AndSlotType, prenode.AndStepType:
			isStep := pn.Type == prenode.AndSlotType ||
				pn.Type == prenode.AndStepType
			g.AddNodes(graph.NewNode(key, graph.AndType, isStep))
		case prenode.OrType, prenode.OrSlotType, prenode.OrStepType,
			prenode.RootType:
			isStep := pn.Type == prenode.OrSlotType ||
				pn.Type == prenode.OrStepType
			g.AddNodes(graph.NewNode(key, graph.OrType, isStep))
		default:
			panic("unknown prenode type for " + key)
		}
	}
}

func addNodeParents(g graph.Graph, prenodes map[string]*prenode.Prenode) {
	for k, pn := range prenodes {
		for _, parent := range pn.Parents {
			g.AddParents(map[string][]string{k: []string{parent.(string)}})
		}
	}
}

// attempts to create a path to the given targets by placing different items in
// slots.
func findRoute(r *Route, start, goal, forbid []string,
	maxlen int) (usedItems, itemList, usedSlots *list.List) {
	// make stacks out of the item names and slot names for backtracking
	var slotList *list.List
	itemList, slotList = initRouteLists(r)

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems = list.New()
	usedSlots = list.New()

	// convert name lists into node lists
	startNodes := make([]*graph.Node, len(start))
	for i, name := range start {
		startNodes[i] = r.Graph[name]
	}
	goalNodes := make([]*graph.Node, len(goal))
	for i, name := range goal {
		goalNodes[i] = r.Graph[name]
	}
	forbidNodes := make([]*graph.Node, len(forbid))
	for i, name := range forbid {
		forbidNodes[i] = r.Graph[name]
	}

	// try to find the route, retrying if needed
	iteration, tries := 0, 0
	for tries = 0; tries < maxTries; tries++ {
		rollSeasons(r.Graph)

		if tryExploreTargets(r.Graph, nil, startNodes, goalNodes, forbidNodes,
			maxlen, &iteration, itemList, usedItems, slotList, usedSlots) {
			log.Print("-- success")
			announceSuccessDetails(r, goal, usedItems, usedSlots)
			break
		} else if iteration > maxIterations {
			log.Print("-- routing took too long; retrying")
			itemList, slotList = initRouteLists(r)
			usedItems, usedSlots = list.New(), list.New()
			iteration = 0
		} else {
			log.Fatal("-- fatal: could not find route")
		}
	}
	if tries >= maxTries {
		log.Fatalf("-- fatal: could not find route after %d tries", maxTries)
	}

	return
}

var (
	seasonsByID = []string{"spring", "summer", "autumn", "winter"}
	seasonAreas = []string{
		"north horon", "eastern suburbs", "woods of winter", "spool swamp",
		"holodrum plain", "sunken city", "lost woods", "tarm ruins",
		"western coast", "temple remains",
	}
)

// set the default seasons for all the applicable areas in the game. also set
// the values in the ROM.
func rollSeasons(g graph.Graph) {
	for _, area := range seasonAreas {
		// reset default seasons
		for _, season := range seasonsByID {
			g[fmt.Sprintf("%s default %s", area, season)].ClearParents()
		}

		// roll new default season
		id := rand.Intn(len(seasonsByID))
		season := seasonsByID[id]
		g.AddParents(map[string][]string{fmt.Sprintf(
			"%s default %s", area, season): []string{"start"}})
		rom.Seasons[fmt.Sprintf("%s season", area)].New = []byte{byte(id)}
	}
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
//
// the lists are lists of nodes.
func tryExploreTargets(g graph.Graph, start map[*graph.Node]bool,
	add, goal, forbid []*graph.Node, maxlen int, iteration *int,
	itemList, usedItems, slotList, usedSlots *list.List) bool {
	*iteration++
	log.Print("iteration ", *iteration)
	if *iteration > maxIterations {
		log.Print("-- false; maximum iterations reached")
		return false
	}

	// explore given the old state and changes
	reached := g.Explore(start, add)
	log.Print(countSteps(reached), " steps reached")

	// check whether to return right now
	fillUnused := false
	switch checkRouteState(
		g, start, reached, add, goal, forbid, slotList, maxlen) {
	case RouteFillUnused:
		fillUnused = true
	case RouteSuccess:
		return true
	case RouteInvalid:
		return false
	}

	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slotElem := slotList.Back()
		slotList.MoveToFront(slotElem)

		// see if slot node has been reached OR we don't care anymore
		slotNode := slotElem.Value.(*graph.Node)
		if !reached[slotNode] && !fillUnused {
			continue
		}

		// move slot from unused to used
		usedSlots.PushBack(slotNode)
		slotList.Remove(slotElem)

		// try placing each unused item into the slot
		jewelChecked := false
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemNode := itemList.Remove(itemList.Back()).(*graph.Node)
			usedItems.PushBack(itemNode)
			g[itemNode.Name].AddParents(g[slotNode.Name])

			printItemSequence(usedItems)

			// recurse unless the item should be skipped
			var skip bool
			skip, jewelChecked = shouldSkipItem(
				g, reached, itemNode, slotNode, jewelChecked, fillUnused)
			if !skip {
				log.Print("trying slot " + slotNode.Name)
				if tryExploreTargets(g, reached, []*graph.Node{itemNode}, goal,
					forbid, maxlen-1, iteration, itemList, usedItems, slotList,
					usedSlots) {
					return true
				}
			}

			// item didn't work; unslot it and pop it onto the front of the
			// unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemNode)
			g[itemNode.Name].ClearParents()
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotNode)
	}

	// nothing worked
	log.Print("-- false; no slot/item combination worked")
	return false
}

// return shuffled lists of item and slot nodes
func initRouteLists(r *Route) (itemList, slotList *list.List) {
	// shuffle names in slices
	items := make([]*graph.Node, 0, len(prenode.BaseItems()))
	slots := make([]*graph.Node, 0, len(r.Slots))
	for itemName := range prenode.BaseItems() {
		items = append(items, r.Graph[itemName])
	}
	for slotName := range r.Slots {
		slots = append(slots, r.Graph[slotName])
	}
	rand.Shuffle(len(items), func(i, j int) {
		items[i], items[j] = items[j], items[i]
	})
	rand.Shuffle(len(slots), func(i, j int) {
		slots[i], slots[j] = slots[j], slots[i]
	})

	// push the shuffled items onto stacks
	itemList = list.New()
	slotList = list.New()
	for _, itemNode := range items {
		itemList.PushBack(itemNode)
	}
	for _, slotNode := range slots {
		slotList.PushBack(slotNode)
	}

	return itemList, slotList
}

// possible return values of checkRouteState
type RouteState int

// possible return values of checkRouteState
const (
	RouteIndeterminate = iota
	RouteFillUnused    // goals reached, some slots still open
	RouteSuccess
	RouteInvalid
)

// returns a RouteState based on whether the route is complete, invalid, or
// needs more work
func checkRouteState(g graph.Graph, start, reached map[*graph.Node]bool,
	add, goal, forbid []*graph.Node, slots *list.List, maxlen int) RouteState {
	// abort if any forbidden node is reached
	for _, node := range forbid {
		if reached[node] {
			log.Printf("-- false; reached forbidden node %s", node)
			return RouteInvalid
		}
	}

	// check for softlocks
	if err := canSoftlock(g); err != nil {
		log.Print("-- false; ", err)
		return RouteInvalid
	}

	// success if all goal nodes are reached *and* all slots are filled
	allReached := true
	for _, node := range goal {
		if !reached[node] {
			log.Printf("-- have not reached goal node %s", node)
			allReached = false
			break
		}
	}
	if allReached {
		log.Print("-- all goals reached")
		if slots.Len() == 0 {
			log.Print("-- true; all goals reached and slots filled")
			return RouteSuccess
		}
		log.Print("-- filling extra slots")
		return RouteFillUnused
	}

	// if the new state doesn't reach any more steps, abandon this branch,
	// *unless* the new item is a jewel, seed item, gale seed, or we've already
	// reached the goals. jewels need this logic because they won't reach any
	// more steps until all four have been slotted, and seed items need this
	// logic because they're useless until seeds have been slotted too.
	//
	// gale seeds don't *need* this logic, strictly speaking, but they're very
	// convenient for the player to have. but still don't slot them until the
	// player already has a seed item, or else they'll probably end up in horon
	// village a lot. also only slot the first one this way! the second one can
	// be filler.
	if !strings.HasSuffix(add[0].Name, " jewel") {
		needCount := true

		// still, don't slot seed stuff until the player can at least harvest
		if reached[g["harvest item"]] {
			switch add[0].Name {
			case "satchel", "slingshot L-1", "slingshot L-2":
				if !(reached[g["slingshot L-2"]] &&
					add[0].Name == "slingshot L-1") {
					needCount = false
				}
			case "gale tree seeds 1":
				if reached[g["seed item"]] {
					needCount = false
				}
			}
		}

		if needCount && countSteps(reached) <= countSteps(start) {
			log.Printf("-- false; reached steps %d <= start steps %d",
				countSteps(reached), countSteps(start))
			return RouteInvalid
		}
	}

	// can't slot any more items
	if maxlen == 0 {
		log.Print("-- false; slotted maxlen items")
		return RouteInvalid
	}

	return RouteIndeterminate
}

// print the currently evaluating sequence of slotted items
func printItemSequence(usedItems *list.List) {
	items := make([]string, 0, usedItems.Len())
	for e := usedItems.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(*graph.Node).Name)
	}
	log.Print("trying " + strings.Join(items, " -> "))
}

// return skip = true iff conditions mean this item shouldn't be checked, and
// checked = true iff a jewel (round, square, pyramid, x-shaped) has been
// checked by now.
func shouldSkipItem(g graph.Graph, reached map[*graph.Node]bool, itemNode,
	slotNode *graph.Node, jewelChecked, fillUnused bool) (skip, checked bool) {
	// only check one jewel per loop, since they're functionally
	// identical.
	if strings.HasSuffix(itemNode.Name, " jewel") {
		if !jewelChecked {
			checked = true
		} else {
			skip = true
		}
	}

	// the star ore code is unique in that it doesn't set the sub ID at
	// all, leaving it zeroed. so if we're looking at the star ore
	// slot, then skip any items that have a nonzero sub ID.
	if slotNode.Name == "star ore spot" &&
		rom.Treasures[itemNode.Name].SubID() != 0 {
		skip = true
	}
	// some items can't be drawn correctly in "scene" item slots.
	switch slotNode.Name {
	case "d0 sword chest", "rod gift", "noble sword spot":
		if !rom.CanSlotInScene(itemNode.Name) {
			skip = true
		}
	}
	// and only seeds can be slotted in seed trees, of course
	switch itemNode.Name {
	case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
		"pegasus tree seeds", "gale tree seeds 1", "gale tree seeds 2":
		switch slotNode.Name {
		case "ember tree", "mystery tree", "scent tree",
			"pegasus tree", "sunken gale tree", "tarm gale tree":
			break
		default:
			skip = true
		}
	default:
		switch slotNode.Name {
		case "ember tree", "mystery tree", "scent tree",
			"pegasus tree", "sunken gale tree", "tarm gale tree":
			skip = true
		}
	}

	return
}

// print item/slot info on a succeeded route
func announceSuccessDetails(
	r *Route, goal []string, usedItems, usedSlots *list.List) {
	log.Print("-- slotted items")

	// iterate by rotating again for some reason
	for i := 0; i < usedItems.Len(); i++ {
		log.Printf("%v <- %v",
			usedItems.Front().Value.(*graph.Node),
			usedSlots.Front().Value.(*graph.Node))
		usedItems.MoveToBack(usedItems.Front())
		usedSlots.MoveToBack(usedSlots.Front())
	}
}

// return the number of "step" nodes in the given set
func countSteps(nodes map[*graph.Node]bool) int {
	count := 0
	for node := range nodes {
		if node.IsStep {
			count++
		}
	}
	return count
}
