package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/logic"
	"github.com/jangler/oos-randomizer/rom"
)

// give up completely if routing fails too many times
const maxTries = 50

// adds nodes to the map based on default contents of item slots.
func addDefaultItemNodes(nodes map[string]*logic.Node) {
	for key, slot := range rom.ItemSlots {
		if key != "rod gift" { // real rod is an Or, not a Root
			nodes[rom.FindTreasureName(slot.Treasure)] = logic.Root()
		}
	}
}

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph graph.Graph
	Slots map[string]*graph.Node
	Costs int
}

// NewRoute returns an initialized route with all nodes, and those nodes with
// the names in start functioning as givens (always satisfied). If no names are
// given, only the normal start node functions as a given.
func NewRoute(start ...string) *Route {
	g := graph.New()

	totalPrenodes := logic.GetAll()
	addDefaultItemNodes(totalPrenodes)

	// make start nodes given
	for _, key := range start {
		totalPrenodes[key] = logic.And()
	}

	addNodes(totalPrenodes, g)
	addNodeParents(totalPrenodes, g)

	openSlots := make(map[string]*graph.Node, 0)
	for name, pn := range totalPrenodes {
		switch pn.Type {
		case logic.AndSlotType, logic.OrSlotType:
			openSlots[name] = g[name]
		}
	}

	return &Route{
		Graph: g,
		Slots: openSlots,
	}
}

func (r *Route) AddParent(child, parent string) {
	r.Graph[child].AddParents(r.Graph[parent])
}

func (r *Route) ClearParents(node string) {
	r.Graph[node].ClearParents()
}

// if hard is false, "hard" nodes are omitted
func addNodes(prenodes map[string]*logic.Node, g graph.Graph) {
	for key, pn := range prenodes {
		switch pn.Type {
		case logic.AndType, logic.AndSlotType, logic.AndStepType,
			logic.HardAndType:
			isStep := pn.Type == logic.AndSlotType ||
				pn.Type == logic.AndStepType
			isSlot := pn.Type == logic.AndSlotType
			isHard := pn.Type == logic.HardAndType

			node := graph.NewNode(key, graph.AndType, isStep, isSlot, isHard)
			g.AddNodes(node)
		case logic.OrType, logic.OrSlotType, logic.OrStepType, logic.RootType,
			logic.HardOrType:
			isStep := pn.Type == logic.OrSlotType ||
				pn.Type == logic.OrStepType
			isSlot := pn.Type == logic.OrSlotType
			nodeType := graph.OrType
			if pn.Type == logic.RootType {
				nodeType = graph.RootType
			}
			isHard := pn.Type == logic.HardOrType

			node := graph.NewNode(key, nodeType, isStep, isSlot, isHard)
			g.AddNodes(node)
		default:
			panic("unknown logic type for " + key)
		}
	}
}

// nodes not in the graph are omitted (for example, "hard" nodes in a non-hard
// graph)
func addNodeParents(prenodes map[string]*logic.Node, gs ...graph.Graph) {
	for k, pn := range prenodes {
		for _, g := range gs {
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
}

type RouteInfo struct {
	Route                        *Route
	Seed                         uint32
	Seasons                      map[string]byte
	Companion                    int // 1 to 3
	UsedItems, UsedSlots         *list.List
	ProgressItems, ProgressSlots *list.List
	ExtraItems, ExtraSlots       *list.List
	AttemptCount                 int
}

const (
	ricky   = 1
	dimitri = 2
	moosh   = 3
)

// attempts to create a path to the given targets by placing different items in
// slots. returns nils if no route is found.
func findRoute(src *rand.Rand, seed uint32, verbose bool, logChan chan string,
	doneChan chan int) *RouteInfo {
	// make stacks out of the item names and slot names for backtracking
	var itemList, slotList *list.List

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	ri := &RouteInfo{
		Seed:          seed,
		UsedItems:     list.New(),
		UsedSlots:     list.New(),
		ProgressItems: list.New(),
		ProgressSlots: list.New(),
		ExtraItems:    list.New(),
		ExtraSlots:    list.New(),
	}

	// try to find the route, retrying if needed
	tries := 0
	for tries = 0; tries < maxTries; tries++ {
		// abort if route was already found on another thread
		select {
		case <-doneChan:
			return nil
		default:
		}

		r := NewRoute()
		ri.Companion = rollAnimalCompanion(src, r)
		itemList, slotList = initRouteInfo(src, r, ri.Companion)
		logChan <- fmt.Sprintf("trying seed %08x", ri.Seed)

		// slot initial nodes before algorithm slots progression items
		ri.Seasons = rollSeasons(src, r)
		placeDungeonItems(src, r,
			itemList, ri.UsedItems, slotList, ri.UsedSlots)

		slotRecord := 0
		i, maxIterations := 0, 1+itemList.Len()

		// slot progression items
		done := r.Graph["done"]
		success := true
		for done.GetMark(done, false) != graph.MarkTrue {
			if verbose {
				logChan <- fmt.Sprintf("searching; have %d more slots",
					slotList.Len())
				logChan <- fmt.Sprintf("%d/%d iterations", i, maxIterations)
			}

			eItem, eSlot := trySlotRandomItem(r, src, itemList, slotList,
				countSteps, ri.UsedSlots.Len(), false)

			if eItem != nil {
				ri.UsedItems.PushBack(itemList.Remove(eItem))
				slot := slotList.Remove(eSlot).(*graph.Node)
				ri.UsedSlots.PushBack(slot)
				r.Costs += logic.Rupees[slot.Name]

				if ri.UsedSlots.Len() > slotRecord {
					slotRecord = ri.UsedSlots.Len()
					i, maxIterations = 0, 1+itemList.Len()
				}
			} else {
				item := ri.UsedItems.Remove(ri.UsedItems.Back()).(*graph.Node)
				slot := ri.UsedSlots.Remove(ri.UsedSlots.Back()).(*graph.Node)
				itemList.PushBack(item)
				slotList.PushBack(slot)
				item.RemoveParent(slot)
			}

			r.Graph.ClearMarks()

			i++
			if i > maxIterations {
				success = false
				if verbose {
					logChan <- "maximum iterations reached"
				}
				break
			}
		}

		if success {
			// fill unused slots
			for slotList.Len() > 0 {
				if verbose {
					logChan <- fmt.Sprintf("done; filling %d more slots",
						slotList.Len())
					logChan <- fmt.Sprintf("%d/%d iterations", i, maxIterations)
				}

				eItem, eSlot := trySlotRandomItem(r, src, itemList, slotList,
					countSteps, ri.UsedSlots.Len(), true)

				if eItem != nil {
					ri.UsedItems.PushBack(itemList.Remove(eItem))
					slot := slotList.Remove(eSlot).(*graph.Node)
					ri.UsedSlots.PushBack(slot)
					r.Costs += logic.Rupees[slot.Name]

					if ri.UsedSlots.Len() > slotRecord {
						slotRecord = ri.UsedSlots.Len()
						i, maxIterations = 0, 1+itemList.Len()
					}
				} else {
					item := ri.UsedItems.Remove(ri.UsedItems.Back()).(*graph.Node)
					slot := ri.UsedSlots.Remove(ri.UsedSlots.Back()).(*graph.Node)
					itemList.PushBack(item)
					slotList.PushBack(slot)
					item.RemoveParent(slot)
				}

				i++
				if i > maxIterations {
					if verbose {
						logChan <- "maximum iterations reached"
					}
					break
				}
			}
		}

		if slotList.Len() == 0 {
			arrangeListsForLog(r, ri, verbose)

			// rotate dungeon items to the back of the lists
			items, slots := ri.ProgressItems, ri.ProgressSlots
			for i := 0; i < 8; i++ {
				items.PushBack(items.Remove(items.Front()))
				slots.PushBack(slots.Remove(slots.Front()))
			}
			items, slots = ri.ExtraItems, ri.ExtraSlots
			for i := 0; i < 16; i++ {
				items.PushBack(items.Remove(items.Front()))
				slots.PushBack(slots.Remove(slots.Front()))
			}

			// and we're done
			ri.Route = r
			ri.AttemptCount = tries + 1
			break
		}

		ri.UsedItems, ri.UsedSlots = list.New(), list.New()
		ri.ProgressItems, ri.ProgressSlots = list.New(), list.New()
		ri.ExtraItems, ri.ExtraSlots = list.New(), list.New()

		// get a new seed for the next iteration
		ri.Seed = uint32(src.Int31())
		src = rand.New(rand.NewSource(int64(ri.Seed)))
	}

	if tries >= maxTries {
		logChan <- fmt.Sprintf("abort; could not find route after %d tries",
			maxTries)
		return nil
	}

	return ri
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
// attempting to slot the other ones).
func placeDungeonItems(src *rand.Rand, r *Route,
	itemList, usedItems, slotList, usedSlots *list.List) {

	// place boss keys first
	for i := 1; i < 9; i++ {
		slotted := false
		for ei := itemList.Front(); ei != nil && !slotted; ei = ei.Next() {
			item := ei.Value.(*graph.Node)
			if item.Name == fmt.Sprintf("d%d boss key", i) {
				for es := slotList.Front(); es != nil; es = es.Next() {
					slot := es.Value.(*graph.Node)
					if dungeonIndex(slot) == i {
						item.AddParents(slot)

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

			itemNode.AddParents(slotNode)
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

var seedNames = []string{"ember tree seeds", "scent tree seeds",
	"pegasus tree seeds", "gale tree seeds", "mystery tree seeds"}

// return shuffled lists of item and slot nodes
func initRouteInfo(src *rand.Rand, r *Route,
	companion int) (itemList, slotList *list.List) {
	// get slices of names
	itemNames := make([]string, 0,
		len(rom.ItemSlots)+len(logic.ExtraItems()))
	slotNames := make([]string, 0, len(r.Slots))
	for key, slot := range rom.ItemSlots {
		switch key {
		case "rod gift": // don't slot vanilla, seasonless rod
			break
		case "tarm gale tree": // use random duplicate seed type
			treasureName := seedNames[rand.Intn(len(seedNames))]
			itemNames = append(itemNames, treasureName)
		default:
			// substitute identified flute for strange flute
			treasureName := rom.FindTreasureName(slot.Treasure)
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
	for key := range logic.ExtraItems() {
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

// return the number of "step" nodes in the given set
func countSteps(r *Route) int {
	r.Graph.ClearMarks()
	reached := r.Graph.ExploreFromStart(false)
	count := 0
	for node := range reached {
		if node.IsStep && node.Name != "village shop 1" &&
			node.Name != "village shop 2" && canAffordSlot(r, node) {
			count++
		}
	}
	return count
}

// break down the used items into required and optional items, so that the log
// makes sense.
func arrangeListsForLog(r *Route, ri *RouteInfo, verbose bool) {
	done := r.Graph["done"]

	// figure out which items aren't necessary
	ei, es := ri.UsedItems.Front(), ri.UsedSlots.Front()
	for i := 0; i < ri.UsedItems.Len(); i++ {
		item, slot := ei.Value.(*graph.Node), es.Value.(*graph.Node)

		// remove parent provisionally
		item.RemoveParent(slot)

		// ask if anyone misses it
		r.Graph.ClearMarks()
		if done.GetMark(done, false) == graph.MarkTrue {
			if verbose {
				fmt.Printf("%s (in %s) is extra\n", item.Name, slot.Name)
			}
			ri.ExtraItems.PushBack(item)
			ri.ExtraSlots.PushBack(slot)
		} else {
			item.AddParents(slot)
			ri.ProgressItems.PushBack(item)
			ri.ProgressSlots.PushBack(slot)
		}

		ei, es = ei.Next(), es.Next()
	}

	// attach removed parents back to optional items
	ei, es = ri.ExtraItems.Front(), ri.ExtraSlots.Front()
	for i := 0; i < ri.ExtraItems.Len(); i++ {
		item, slot := ei.Value.(*graph.Node), es.Value.(*graph.Node)
		item.AddParents(slot)
		ei, es = ei.Next(), es.Next()
	}
}
