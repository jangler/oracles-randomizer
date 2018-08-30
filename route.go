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

// give up completely if routing fails too many times
const maxTries = 50

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
	OldSlots         map[*graph.Node]bool
	DungeonItems     []int
	KeyItemsTotal    int
	KeyItemsPlaced   int
	Costs            int
}

// NewRoute returns an initialized route with all prenodes, and those prenodes
// with the names in start functioning as givens (always satisfied). If no
// names are given, only the normal start node functions as a given.
func NewRoute(start ...string) *Route {
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
		OldSlots:       make(map[*graph.Node]bool),
		DungeonItems:   make([]int, 9),
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
			isSlot := pn.Type == prenode.AndSlotType
			if hard || pn.Type != prenode.HardAndType {
				g.AddNodes(graph.NewNode(key, graph.AndType, isStep, isSlot))
			}
		case prenode.OrType, prenode.OrSlotType, prenode.OrStepType,
			prenode.RootType, prenode.HardOrType:
			isStep := pn.Type == prenode.OrSlotType ||
				pn.Type == prenode.OrStepType
			isSlot := pn.Type == prenode.OrSlotType
			nodeType := graph.OrType
			if pn.Type == prenode.RootType {
				nodeType = graph.RootType
			}
			if hard || pn.Type != prenode.HardOrType {
				g.AddNodes(graph.NewNode(key, nodeType, isStep, isSlot))
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
	Companion                         int // 1 to 3
	UsedItems, UnusedItems, UsedSlots *list.List
}

const (
	ricky   = 1
	dimitri = 2
	moosh   = 3
)

// attempts to create a path to the given targets by placing different items in
// slots. returns nils if no route is found.
func findRoute(src *rand.Rand, seed uint32, r *Route, keyonly, verbose bool,
	logChan chan string, doneChan chan int) *RouteLists {
	// make stacks out of the item names and slot names for backtracking
	var itemList, slotList *list.List

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems := list.New()
	usedSlots := list.New()

	// try to find the route, retrying if needed
	var seasons map[string]byte
	var companion int
	tries := 0
	for tries = 0; tries < maxTries; tries++ {
		// abort if route was already found on another thread
		select {
		case <-doneChan:
			return nil
		default:
		}

		companion = rollAnimalCompanion(src, r)
		itemList, slotList = initRouteLists(src, r, companion, keyonly)
		logChan <- fmt.Sprintf("trying seed %08x", seed)

		// slot initial nodes before algorithm slots progression items
		seasons = rollSeasons(src, r)
		if !keyonly {
			placeDungeonItems(src, r, itemList, usedItems, slotList, usedSlots)
		}

		// clear "old" slots and item counts, since we're starting fresh
		for k := range r.OldSlots {
			delete(r.OldSlots, k)
		}
		for i := range r.DungeonItems {
			r.DungeonItems[i] = 0
		}
		r.Costs = 0

		// slot progression items
		done := r.Graph["done"]
		for done.GetMark(done, nil) != graph.MarkTrue {
			if verbose {
				logChan <- fmt.Sprintf("searching; have %d more slots",
					slotList.Len())
			}

			// try to find a new combination of items that opens progression
			items, slots := trySlotItemSet(r, src, itemList, slotList,
				countSteps, false)

			if items != nil {
				for items.Len() > 0 {
					usedItems.PushBack(items.Remove(items.Front()))
					slot := slots.Remove(slots.Front()).(*graph.Node)
					usedSlots.PushBack(slot)

					match := dungeonRegexp.FindStringSubmatch(slot.Name)
					if match != nil {
						di, _ := strconv.Atoi(match[1])
						r.DungeonItems[di]++
					}
				}
			} else {
				break
			}
		}

		// if goal was reached, fill unused slots
		if done.GetMark(done, nil) == graph.MarkTrue {
			for slotList.Len() > 0 {
				if verbose {
					logChan <- fmt.Sprintf("done; filling %d more slots",
						slotList.Len())
				}

				items, slots := trySlotItemSet(r, src, itemList, slotList,
					countSteps, true)
				if items != nil {
					for items.Len() > 0 {
						usedItems.PushBack(items.Remove(items.Front()))
						usedSlots.PushBack(slots.Remove(slots.Front()))
					}
				} else {
					break
				}
			}
		}

		if slotList.Len() == 0 {
			break
		} else if verbose {
			logChan <- "unfilled slots:"
			for e := slotList.Front(); e != nil; e = e.Next() {
				logChan <- e.Value.(*graph.Node).Name
			}
			logChan <- "unused items:"
			for e := itemList.Front(); e != nil; e = e.Next() {
				logChan <- e.Value.(*graph.Node).Name
			}
		}

		itemList, slotList = initRouteLists(src, r, companion, keyonly)
		for e := itemList.Front(); e != nil; e = e.Next() {
			e.Value.(*graph.Node).ClearParents()
		}
		r.Graph.ClearMarks()
		r.HardGraph.ClearMarks()
		usedItems, usedSlots = list.New(), list.New()

		// get a new seed for the next iteration
		seed = uint32(src.Int31())
		src = rand.New(rand.NewSource(int64(seed)))
	}

	if tries >= maxTries {
		logChan <- fmt.Sprintf("abort; could not find route after %d tries",
			maxTries)
		return nil
	}

	return &RouteLists{seed, seasons, companion, usedItems, itemList, usedSlots}
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

// randomly determines animal companion and returns its ID (1 to 3)
func rollAnimalCompanion(src *rand.Rand, r *Route) int {
	companion := src.Intn(3) + 1

	r.ClearParents("natzu prairie")
	r.ClearParents("natzu river")
	r.ClearParents("natzu wasteland")

	switch companion {
	case ricky:
		r.AddParent("natzu prairie", "start")
	case dimitri:
		r.AddParent("natzu river", "start")
	case moosh:
		r.AddParent("natzu wasteland", "start")
	}

	return companion
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

// place maps, compasses, and boss keys in chests in dungeons (before
// attempting to slot the other ones)
func placeDungeonItems(src *rand.Rand, r *Route,
	itemList, usedItems, slotList, usedSlots *list.List) {
	// place boss keys first
	for i := 1; i < 9; i++ {
		if i == 4 || i == 5 {
			continue
		}

		slotted := false
		for ei := itemList.Front(); ei != nil && !slotted; ei = ei.Next() {
			item := ei.Value.(*graph.Node)
			if item.Name == fmt.Sprintf("d%d boss key", i) {
				for es := slotList.Front(); es != nil; es = es.Next() {
					slot := es.Value.(*graph.Node)
					if dungeonIndex(slot) == i &&
						rom.IsChest(rom.ItemSlots[slot.Name]) {
						r.AddParent(item.Name, slot.Name)

						usedSlots.PushBack(slot)
						slotList.Remove(es)
						usedItems.PushBack(item)
						itemList.Remove(ei)

						slotted = true
						break
					}
				}
			}
		}
	}

	// then place maps and compasses
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
		if dungeonIndex(slot) != index || !rom.IsChest(rom.ItemSlots[slot.Name]) {
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

func emptyList(l *list.List) []*graph.Node {
	a := make([]*graph.Node, l.Len())
	i := 0
	for l.Len() > 0 {
		a[i] = l.Remove(l.Front()).(*graph.Node)
		i++
	}
	return a
}

func fillList(l *list.List, a []*graph.Node) {
	for _, node := range a {
		l.PushBack(node)
	}
}

// return shuffled lists of item and slot nodes
func initRouteLists(src *rand.Rand, r *Route, companion int,
	keyonly bool) (itemList, slotList *list.List) {
	// get slices of names
	itemNames := make([]string, 0,
		len(rom.ItemSlots)+len(prenode.ExtraItems()))
	slotNames := make([]string, 0, len(r.Slots))
	for key, slot := range rom.ItemSlots {
		if key != "rod gift" { // don't slot vanilla, seasonless rod
			treasureName := rom.FindTreasureName(slot.Treasure)
			if !keyonly || keyItems[treasureName] {
				// substitute identified flute for strange flute
				if treasureName == "strange flute" {
					switch companion {
					case ricky:
						treasureName = "ricky's flute"
					case dimitri:
						treasureName = "dimitri's flute"
					case moosh:
						treasureName = "moosh's flute"
					}
				}

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

// print the currently evaluating sequence of slotted items
func printItemSequence(usedItems *list.List, logChan chan string) {
	items := make([]string, 0, usedItems.Len())
	for e := usedItems.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value.(*graph.Node).Name)
	}
	logChan <- fmt.Sprintf("trying %s", strings.Join(items, " -> "))
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

// return the number of "step" nodes in the given set which are not also slots
func countOnlySteps(nodes map[*graph.Node]bool) int {
	count := 0
	for node := range nodes {
		if node.IsStep && !node.IsSlot {
			count++
		}
	}
	return count
}
