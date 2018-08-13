package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sort"
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
	Graph, HardGraph graph.Graph
	Slots            map[string]*graph.Node
}

// NewRoute returns an initialized route with all prenodes, and those prenodes
// with the names in start functioning as givens (always satisfied).
func NewRoute(start []string) *Route {
	g, hg := graph.New(), graph.New()

	totalPrenodes := prenode.GetAll()

	// make start nodes given
	for _, key := range start {
		totalPrenodes[key] = prenode.And()
	}

	addNodes(g, totalPrenodes, false)
	addNodeParents(g, totalPrenodes)
	addNodes(hg, totalPrenodes, true)
	addNodeParents(hg, totalPrenodes)

	openSlots := make(map[string]*graph.Node, 0)
	for name, pn := range totalPrenodes {
		switch pn.Type {
		case prenode.AndSlotType, prenode.OrSlotType:
			openSlots[name] = g[name]
		}
	}

	return &Route{Graph: g, HardGraph: hg, Slots: openSlots}
}

func (r *Route) AddParent(child, parent string) {
	r.Graph[child].AddParents(r.Graph[parent])
	r.HardGraph[child].AddParents(r.HardGraph[parent])
}

func (r *Route) ClearParents(node string) {
	r.Graph[node].ClearParents()
	r.HardGraph[node].ClearParents()
}

// if hard is false, "hard" nodes are omitted
func addNodes(g graph.Graph, prenodes map[string]*prenode.Prenode, hard bool) {
	for key, pn := range prenodes {
		switch pn.Type {
		case prenode.AndType, prenode.AndSlotType, prenode.AndStepType,
			prenode.HardAndType:
			isStep := pn.Type == prenode.AndSlotType ||
				pn.Type == prenode.AndStepType
			if hard || pn.Type != prenode.HardAndType {
				g.AddNodes(graph.NewNode(key, graph.AndType, isStep))
			}
		case prenode.OrType, prenode.OrSlotType, prenode.OrStepType,
			prenode.RootType, prenode.HardOrType:
			isStep := pn.Type == prenode.OrSlotType ||
				pn.Type == prenode.OrStepType
			if hard || pn.Type != prenode.HardOrType {
				g.AddNodes(graph.NewNode(key, graph.OrType, isStep))
			}
		default:
			panic("unknown prenode type for " + key)
		}
	}
}

// nodes not in the graph are omitted (for example, "hard" nodes in a non-hard
// graph)
func addNodeParents(g graph.Graph, prenodes map[string]*prenode.Prenode) {
	for k, pn := range prenodes {
		if g[k] == nil {
			continue
		}
		for _, parent := range pn.Parents {
			if g[parent.(string)] == nil {
				continue
			}
			g.AddParents(map[string][]string{k: []string{parent.(string)}})
		}
	}
}

type RouteLists struct {
	Seed                              uint32
	Seasons                           map[string]byte
	UsedItems, UnusedItems, UsedSlots *list.List
}

// attempts to create a path to the given targets by placing different items in
// slots. returns nils if no route is found.
func findRoute(src *rand.Rand, seed uint32, r *Route, verbose bool,
	logChan chan string, doneChan chan int) *RouteLists {
	// make stacks out of the item names and slot names for backtracking
	itemList, slotList := initRouteLists(src, r)

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems := list.New()
	usedSlots := list.New()

	start := []*graph.Node{r.Graph["horon village"]}

	// try to find the route, retrying if needed
	var seasons map[string]byte
	iteration, tries := 0, 0
	for tries = 0; tries < maxTries; tries++ {
		// abort if route was already found on another thread
		select {
		case <-doneChan:
			return nil
		default:
		}

		seasons = rollSeasons(src, r)
		logChan <- fmt.Sprintf("searching for route (%d)", tries+1)

		if tryExploreTargets(r, nil, start, &iteration, itemList, usedItems,
			slotList, usedSlots, verbose, logChan) {
			if verbose {
				announceSuccessDetails(r, usedItems, usedSlots, logChan)
			}
			break
		} else if iteration > maxIterations {
			if verbose {
				logChan <- "routing took too long; retrying"
			}
			itemList, slotList = initRouteLists(src, r)
			usedItems, usedSlots = list.New(), list.New()
			iteration = 0
		} else {
			logChan <- "could not find route"
		}
	}
	if tries >= maxTries {
		logChan <- fmt.Sprintf("abort; could not find route after %d tries",
			maxTries)
		return nil
	}

	return &RouteLists{seed, seasons, usedItems, itemList, usedSlots}
}

var (
	seasonsByID = []string{"spring", "summer", "autumn", "winter"}
	seasonAreas = []string{
		"north horon", "eastern suburbs", "woods of winter", "spool swamp",
		"holodrum plain", "sunken city", "lost woods", "tarm ruins",
		"western coast", "temple remains",
	}
)

// set the default seasons for all the applicable areas in the game, and return
// a mapping of area name to season value.
func rollSeasons(src *rand.Rand, r *Route) map[string]byte {
	seasonMap := make(map[string]byte, len(seasonAreas))

	for _, area := range seasonAreas {
		// reset default seasons
		for _, season := range seasonsByID {
			r.ClearParents(fmt.Sprintf("%s default %s", area, season))
		}

		// roll new default season
		id := src.Intn(len(seasonsByID))
		season := seasonsByID[id]
		r.AddParent(fmt.Sprintf("%s default %s", area, season), "start")
		seasonMap[area] = byte(id)
	}

	return seasonMap
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
//
// the lists are lists of nodes.
func tryExploreTargets(r *Route, start map[*graph.Node]bool, add []*graph.Node,
	iteration *int, itemList, usedItems, slotList, usedSlots *list.List,
	verbose bool, logChan chan string) bool {
	*iteration++
	if verbose {
		logChan <- fmt.Sprintf("iteration %d", *iteration)
	}
	if *iteration > maxIterations {
		if verbose {
			logChan <- "false; maximum iterations reached"
		}
		return false
	}

	// explore given the old state and changes
	reached := r.Graph.Explore(start, add)
	if verbose {
		logChan <- fmt.Sprintf("%d steps reached", countSteps(reached))
	}

	// check whether to return right now
	fillUnused := false
	switch checkRouteState(
		r, start, reached, add, slotList, verbose, logChan) {
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

		// try a piece of heart if we've already slotted everything useful
		filledByPoH := false
		if fillUnused {
			if rom.ItemSlots[slotNode.Name].CollectMode == rom.CollectChest {
				usedItems.PushBack(graph.NewNode(
					"piece of heart (chest)", graph.RootType, false))

				if verbose {
					printItemSequence(usedItems, logChan)
					logChan <- fmt.Sprintf("trying slot %s", slotNode.Name)
				}
				if tryExploreTargets(r, reached, nil, iteration, itemList,
					usedItems, slotList, usedSlots, verbose, logChan) {
					return true
				}

				usedItems.Remove(usedItems.Back())
			}
		}

		// try placing each unused item into the slot
		jewelChecked := false
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemNode := itemList.Remove(itemList.Back()).(*graph.Node)
			usedItems.PushBack(itemNode)
			r.AddParent(itemNode.Name, slotNode.Name)

			if verbose {
				printItemSequence(usedItems, logChan)
			}

			// recurse unless the item should be skipped
			var skip bool
			skip, jewelChecked = shouldSkipItem(
				r.Graph, itemNode, slotNode, jewelChecked, fillUnused)
			if !skip {
				if verbose {
					logChan <- fmt.Sprintf("trying slot %s", slotNode.Name)
				}
				if tryExploreTargets(r, reached, []*graph.Node{itemNode},
					iteration, itemList, usedItems, slotList, usedSlots,
					verbose, logChan) {
					return true
				}
			}

			// item didn't work; unslot it and pop it onto the front of the
			// unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemNode)
			r.ClearParents(itemNode.Name)
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotNode)
	}

	// nothing worked
	if verbose {
		logChan <- "false; no slot/item combination worked"
	}
	return false
}

// return shuffled lists of item and slot nodes
func initRouteLists(src *rand.Rand, r *Route) (itemList, slotList *list.List) {
	// get slices of names
	itemNames := make([]string, 0, len(prenode.BaseItems()))
	slotNames := make([]string, 0, len(r.Slots))
	for key := range prenode.BaseItems() {
		itemNames = append(itemNames, key)
	}
	for key := range r.Slots {
		slotNames = append(slotNames, key)
	}

	// sort the slices so that order isn't dependent on map implementation,
	// then shuffle the sorted slices
	sort.Strings(itemNames)
	sort.Strings(slotNames)
	src.Shuffle(len(itemNames), func(i, j int) {
		itemNames[i], itemNames[j] = itemNames[j], itemNames[i]
	})
	src.Shuffle(len(slotNames), func(i, j int) {
		slotNames[i], slotNames[j] = slotNames[j], slotNames[i]
	})

	// push the graph nodes by name onto stacks
	itemList = list.New()
	slotList = list.New()
	for _, key := range itemNames {
		itemList.PushBack(r.Graph[key])
	}
	for _, key := range slotNames {
		slotList.PushBack(r.Graph[key])
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
func checkRouteState(r *Route, start, reached map[*graph.Node]bool,
	add []*graph.Node, slots *list.List, verbose bool,
	logChan chan string) RouteState {
	// check for softlocks
	r.HardGraph.ExploreFromStart()
	if err := canSoftlock(r.HardGraph); err != nil {
		if verbose {
			logChan <- fmt.Sprintf("false; %v", err)
		}
		return RouteInvalid
	}

	// success if all goal nodes are reached *and* all slots are filled
	if reached[r.Graph["done"]] {
		if verbose {
			logChan <- "goal reached"
		}
		if slots.Len() == 0 {
			if verbose {
				logChan <- "true; goal reached and slots filled"
			}
			return RouteSuccess
		}
		if verbose {
			logChan <- "filling extra slots"
		}
		return RouteFillUnused
	} else {
		if verbose {
			logChan <- "have not reached goal"
		}
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
		if reached[r.Graph["harvest tree"]] {
			switch add[0].Name {
			case "satchel", "slingshot L-1", "slingshot L-2":
				if !(reached[r.Graph["slingshot L-2"]] &&
					add[0].Name == "slingshot L-1") {
					needCount = false
				}
			case "gale tree seeds 1":
				if reached[r.Graph["seed item"]] {
					needCount = false
				}
			}
		}

		if needCount && countSteps(reached) <= countSteps(start) {
			if verbose {
				logChan <- fmt.Sprintf(
					"false; reached steps %d <= start steps %d",
					countSteps(reached), countSteps(start))
			}
			return RouteInvalid
		}
	}

	return RouteIndeterminate
}

// print the currently evaluating sequence of slotted items
func printItemSequence(usedItems *list.List, logChan chan string) {
	items := make([]string, 0, usedItems.Len())
	for e := usedItems.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(*graph.Node).Name)
	}
	logChan <- fmt.Sprintf("trying %s", strings.Join(items, " -> "))
}

// return skip = true iff conditions mean this item shouldn't be checked, and
// checked = true iff a jewel (round, square, pyramid, x-shaped) has been
// checked by now.
func shouldSkipItem(g graph.Graph, itemNode, slotNode *graph.Node,
	jewelChecked, fillUnused bool) (skip, checked bool) {
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
func announceSuccessDetails(r *Route, usedItems, usedSlots *list.List,
	logChan chan string) {
	logChan <- "slotted items:"

	// iterate by rotating again for some reason
	for i := 0; i < usedItems.Len(); i++ {
		logChan <- fmt.Sprintf("%v <- %v",
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
