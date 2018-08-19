package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
	"github.com/jangler/oos-randomizer/rom"
)

const (
	// a routing attempt fails if it fails to fill a slot this many times
	maxStrikes = 3

	maxTries = 50 // give up completely if routing fails too many times
)

// adds prenodes to the map based on default contents of item slots.
func addDefaultItemNodes(nodes map[string]*prenode.Prenode) {
	for key, slot := range rom.ItemSlots {
		if key != "rod gift" { // real rod is an Or, not a Root
			nodes[rom.FindTreasureName(slot.Treasure)] = prenode.Root()
		}
	}
}

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph, HardGraph graph.Graph
	Slots            map[string]*graph.Node
	Dungeons         []Dungeon
	KeyItemsTotal    int
	KeyItemsPlaced   int
}

type Dungeon struct {
	ItemsPlaced int
	HasMap      bool
	HasCompass  bool
}

// NewRoute returns an initialized route with all prenodes, and those prenodes
// with the names in start functioning as givens (always satisfied).
func NewRoute(start []string) *Route {
	g, hg := graph.New(), graph.New()

	totalPrenodes := prenode.GetAll()
	addDefaultItemNodes(totalPrenodes)

	// make start nodes given
	for _, key := range start {
		totalPrenodes[key] = prenode.And()
	}

	addNodes(g, totalPrenodes, false)
	addNodeParents(g, totalPrenodes)
	addNodes(hg, totalPrenodes, true)
	addNodeParents(hg, totalPrenodes)

	keyItemCount := 0
	openSlots := make(map[string]*graph.Node, 0)
	for name, pn := range totalPrenodes {
		switch pn.Type {
		case prenode.RootType:
			if keyItems[name] {
				keyItemCount++
			}
		case prenode.AndSlotType, prenode.OrSlotType:
			openSlots[name] = g[name]
		}
	}

	return &Route{
		Graph:          g,
		HardGraph:      hg,
		Slots:          openSlots,
		Dungeons:       make([]Dungeon, 9),
		KeyItemsTotal:  keyItemCount,
		KeyItemsPlaced: 0,
	}
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
func findRoute(src *rand.Rand, seed uint32, r *Route, keyonly, verbose bool,
	logChan chan string, doneChan chan int) *RouteLists {
	// make stacks out of the item names and slot names for backtracking
	itemList, slotList := initRouteLists(src, r, keyonly)

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems := list.New()
	usedSlots := list.New()

	start := []*graph.Node{r.Graph["horon village"]}

	// try to find the route, retrying if needed
	var seasons map[string]byte
	strikes, tries := 0, 0
	for tries = 0; tries < maxTries; tries++ {
		// abort if route was already found on another thread
		select {
		case <-doneChan:
			return nil
		default:
		}

		seasons = rollSeasons(src, r)
		logChan <- fmt.Sprintf("searching for route (%d)", tries+1)

		if !keyonly {
			placeDungeonItems(src, itemList, usedItems, slotList, usedSlots)
		}

		if tryExploreTargets(src, r, nil, start, &strikes, itemList,
			usedItems, slotList, usedSlots, verbose, logChan) {
			if verbose {
				announceSuccessDetails(r, usedItems, usedSlots, logChan)
			}
			break
		} else if strikes >= maxStrikes {
			if verbose {
				logChan <- "routing struck out; retrying"
			}
			itemList, slotList = initRouteLists(src, r, keyonly)
			usedItems, usedSlots = list.New(), list.New()
			strikes = 0
		} else {
			logChan <- "could not find route"
		}
	}
	if tries >= maxTries {
		logChan <- fmt.Sprintf("abort; could not find route after %d tries",
			maxTries)
		return nil
	}

	if verbose {
		logChan <- fmt.Sprintf("%d slots, %d strike(s)", usedSlots.Len(), strikes)
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

// dungeonIndex returns the index of a slot's dungeon if it's in a dungeon, or
// -1 if it's not.
func dungeonIndex(node *graph.Node) int {
	isInDungeon, _ := regexp.MatchString(`^d\d `, node.Name)
	if isInDungeon {
		index, _ := strconv.Atoi(string(node.Name[1]))
		return index
	}
	return -1
}

// place maps and compasses in chests in dungeons (before attempting to slot
// the other ones)
func placeDungeonItems(src *rand.Rand,
	itemList, usedItems, slotList, usedSlots *list.List) {
	for i := 1; i < 9; i++ {
		for _, itemName := range []string{"dungeon map", "compass"} {
			slotElem, itemElem, slotNode, itemNode :=
				getDungeonItem(i, itemName, slotList, itemList)

			usedSlots.PushBack(slotNode)
			slotList.Remove(slotElem)
			usedItems.PushBack(itemNode)
			itemList.Remove(itemElem)
		}
	}
}

func getDungeonItem(index int, itemName string, slotList,
	itemList *list.List) (slotElem, itemElem *list.Element, slotNode, itemNode *graph.Node) {
	for es := slotList.Front(); es != nil; es = es.Next() {
		slot := es.Value.(*graph.Node)
		if dungeonIndex(slot) != index {
			continue
		}

		for ei := itemList.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*graph.Node)
			if item.Name != itemName {
				continue
			}

			return es, ei, slot, item
		}
	}

	panic("could not place dungeon-specific items")
}

// because linked lists are really bad for this type of operation, sort them by
// emptying the list into a slice and then refilling it after sorting. this is
// only for lists of graph nodes!
func emptyList(l *list.List) []*graph.Node {
	// empty list into slice
	a := make([]*graph.Node, 0, l.Len())
	for l.Len() > 0 {
		value := l.Remove(l.Front()).(*graph.Node)
		a = append(a, value)
	}
	return a
}

// see emptyList comment
func refillList(l *list.List, a []*graph.Node) {
	for _, node := range a {
		l.PushBack(node)
	}
}

// check dungeon slots first if the respective dungeon has no items in it, so
// that key items are more likely to end up in dungeons. if the slot has a
// "special" collect mode, then check that next so that a unique item (i.e. one
// that can actually fit in it) is likely to end up there.
//
// the back of the list is checked first, so non-dungeon slots should be
// counted as "less".
func sortSlots(r *Route, l *list.List) {
	a := emptyList(l)

	sort.Slice(a, func(i, j int) bool {
		// dungeon chests go first
		di := dungeonIndex(a[j])
		if di >= 0 && r.Dungeons[di].ItemsPlaced == 0 {
			return true
		}

		// special item slots go second
		switch rom.ItemSlots[a[i].Name].CollectMode {
		case rom.CollectChest, rom.CollectFind1, rom.CollectFind2:
			return true
		default:
			return false
		}
	})

	refillList(l, a)
}

// check key items and rupees first, since other types of items aren't useful
// for progression.
//
// the back of the list is checked first, so non-key items should be counted as
// "less".
func sortItems(l *list.List) {
	a := emptyList(l)

	sort.Slice(a, func(i, j int) bool {
		return keyItems[a[j].Name] || strings.HasPrefix(a[j].Name, "rupees")
	})

	refillList(l, a)
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
//
// the lists are lists of nodes.
func tryExploreTargets(src *rand.Rand, r *Route, start map[*graph.Node]bool,
	add []*graph.Node, strikes *int, itemList, usedItems, slotList,
	usedSlots *list.List, verbose bool, logChan chan string) bool {
	// explore given the old state and changes
	reached := r.Graph.Explore(start, add)

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

	// get slot priotity (see function comment for details)
	sortSlots(r, slotList)

	// no point in trying to put items in multiple slots of the same type,
	// since they'll be equivalent if they're reachable
	triedCollectModes := map[byte]bool{}

	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slotElem := slotList.Back()
		slotList.MoveToFront(slotElem)

		// if we haven't reached the node yet, don't bother checking it, unless
		// we're just filling unused slots
		slotNode := slotElem.Value.(*graph.Node)
		if !reached[slotNode] && !fillUnused {
			if verbose {
				logChan <- fmt.Sprintf("can't reach slot %s", slotNode.Name)
			}
			continue
		}

		// continue if we're already tried a slot of this type
		if triedCollectModes[rom.ItemSlots[slotNode.Name].CollectMode] {
			continue
		}
		triedCollectModes[rom.ItemSlots[slotNode.Name].CollectMode] = true

		if verbose {
			logChan <- fmt.Sprintf("trying slot %s", slotNode.Name)
		}

		// move slot from unused to used
		usedSlots.PushBack(slotNode)
		slotList.Remove(slotElem)

		sortItems(itemList)

		// try placing each unused item into the slot
		jewelChecked := false
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemNode := itemList.Remove(itemList.Back()).(*graph.Node)
			usedItems.PushBack(itemNode)
			r.AddParent(itemNode.Name, slotNode.Name)

			// recurse unless the item should be skipped
			var skip bool
			skip, jewelChecked = shouldSkipItem(src, r.Graph, reached,
				itemNode, slotNode, jewelChecked, fillUnused)

			if keyItems[itemNode.Name] {
				r.KeyItemsPlaced++
			}

			// count that we're placing an item in a dungeon
			di := dungeonIndex(slotNode)
			if di >= 0 {
				r.Dungeons[di].ItemsPlaced++
			}

			if !skip {
				if verbose {
					logChan <- fmt.Sprintf("trying item %s", itemNode.Name)
				}
				if tryExploreTargets(src, r, reached, []*graph.Node{itemNode},
					strikes, itemList, usedItems, slotList, usedSlots,
					verbose, logChan) {
					return true
				}
			}

			// didn't place item in dungeon after all
			if di >= 0 {
				r.Dungeons[di].ItemsPlaced--
			}

			if keyItems[itemNode.Name] {
				r.KeyItemsPlaced--
			}

			// item didn't work; unslot it and pop it onto the front of the
			// unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemNode)
			r.ClearParents(itemNode.Name)

			if *strikes >= maxStrikes {
				if verbose {
					logChan <- "false; maximum strikes reached"
				}
				return false
			}
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotNode)
	}

	// nothing worked
	*strikes++
	if verbose {
		logChan <- "false; no slot/item combination worked"
	}
	return false
}

// return shuffled lists of item and slot nodes
func initRouteLists(src *rand.Rand, r *Route,
	keyonly bool) (itemList, slotList *list.List) {
	// get slices of names
	itemNames := make([]string, 0,
		len(rom.ItemSlots)+len(prenode.ExtraItems()))
	slotNames := make([]string, 0, len(r.Slots))
	for key, slot := range rom.ItemSlots {
		if key != "rod gift" { // don't slot vanilla, seasonless rod
			treasureName := rom.FindTreasureName(slot.Treasure)
			if !keyonly || keyItems[treasureName] {
				itemNames = append(itemNames, treasureName)
			}
		}
	}
	for key := range prenode.ExtraItems() {
		if !keyonly || keyItems[key] {
			itemNames = append(itemNames, key)
		}
	}
	for key := range r.Slots {
		if !keyonly || keySlots[key] {
			slotNames = append(slotNames, key)
		}
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

	// if the new state hasn't reached enough essences at this stage in the
	// process, it's invalid. this is to help prevent seeds from becoming
	// mostly overworld treks with d3, d4, and d6 always at the end.
	essencesReached := 0
	for node := range reached {
		if strings.HasSuffix(node.Name, "essence") {
			essencesReached++
		}
	}
	if r.KeyItemsPlaced > 0 &&
		essencesReached < 4*r.KeyItemsPlaced/r.KeyItemsTotal {
		if verbose {
			logChan <- "false; have not reached enough essences"
		}
		return RouteInvalid
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
//
// TODO look into why every "skip = true" isn't a "return true, checked"
func shouldSkipItem(src *rand.Rand, g graph.Graph,
	reached map[*graph.Node]bool, itemNode, slotNode *graph.Node, jewelChecked,
	fillUnused bool) (skip, checked bool) {
	// only check one jewel per loop, since they're functionally
	// identical.
	if strings.HasSuffix(itemNode.Name, " jewel") {
		if !jewelChecked {
			checked = true
		} else {
			skip = true
		}
	}

	// don't try non-progression items when trying to progress.
	if !fillUnused && (!keyItems[itemNode.Name] ||
		strings.HasPrefix(itemNode.Name, "rupees")) {
		skip = true
	}

	// gasha seeds and pieces of heart can be placed in either chests or
	// found/gift slots. beyond that, only unique items can be placed in
	// non-chest slots.
	if itemNode.Name == "gasha seed" || itemNode.Name == "piece of heart" {
		switch rom.ItemSlots[slotNode.Name].CollectMode {
		case rom.CollectFind1, rom.CollectFind2, rom.CollectChest:
			if slotNode.Name == "d0 sword chest" ||
				slotNode.Name == "rod gift" {
				skip = true
			}
		default:
			skip = true
		}
	} else if (rom.ItemSlots[slotNode.Name].CollectMode != rom.CollectChest ||
		slotNode.Name == "d0 sword chest" || slotNode.Name == "rod gift") &&
		!rom.TreasureIsUnique[itemNode.Name] {
		skip = true
	}

	// don't put gale seeds in the ember tree, since then gale seeds will come
	// with the satchel and the player can freeze the game by trying to warp
	// without having explored any trees.
	if slotNode.Name == "ember tree" &&
		strings.HasPrefix(itemNode.Name, "gale tree seeds") {
		skip = true
	}

	// give only a 1 in 2 change per sword of slotting in the hero's cave chest
	// to compensate for the fact that there are two of them. each season gets
	// a 1 in 4 chance for the same reason.
	if slotNode.Name == "d0 sword chest" {
		switch itemNode.Name {
		case "sword L-1", "sword L-2":
			if src.Intn(2) != 0 {
				skip = true
			}
		case "winter", "spring", "summer", "autumn":
			if src.Intn(4) != 0 {
				skip = true
			}
		}
	}

	// the star ore code is unique in that it doesn't set the sub ID at all,
	// leaving it zeroed. so if we're looking at the star ore slot, then skip
	// any items that have a nonzero sub ID.
	//
	// the master diver is similar in that he decides whether to give his item
	// based on whether you have on with that ID, meaning that he won't upgrade
	// L-1 items.
	switch slotNode.Name {
	case "star ore spot", "diver gift":
		if rom.Treasures[itemNode.Name].SubID() != 0 {
			skip = true
		}
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
