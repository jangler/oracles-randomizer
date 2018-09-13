package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"time"

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
	Graph        graph.Graph
	Slots        map[string]*graph.Node
	TurnsReached map[*graph.Node]int
	DungeonItems []int
	Costs        int
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
		Graph:        g,
		Slots:        openSlots,
		TurnsReached: make(map[*graph.Node]int),
		DungeonItems: make([]int, 9),
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

type RouteLists struct {
	Seed                         uint32
	Seasons                      map[string]byte
	Companion                    int // 1 to 3
	UsedItems, UsedSlots         *list.List
	RequiredItems, RequiredSlots *list.List
	OptionalItems, OptionalSlots *list.List
}

const (
	ricky   = 1
	dimitri = 2
	moosh   = 3
)

// attempts to create a path to the given targets by placing different items in
// slots. returns nils if no route is found.
func findRoute(src *rand.Rand, seed uint32, verbose bool, logChan chan string,
	doneChan chan int) *RouteLists {
	// make stacks out of the item names and slot names for backtracking
	var itemList, slotList *list.List

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	rl := &RouteLists{
		Seed:          seed,
		UsedItems:     list.New(),
		UsedSlots:     list.New(),
		RequiredItems: list.New(),
		RequiredSlots: list.New(),
		OptionalItems: list.New(),
		OptionalSlots: list.New(),
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
		rl.Companion = rollAnimalCompanion(src, r)
		itemList, slotList = initRouteLists(src, r, rl.Companion)
		logChan <- fmt.Sprintf("trying seed %08x", seed)

		// slot initial nodes before algorithm slots progression items
		rl.Seasons = rollSeasons(src, r)
		placeDungeonItems(src, r,
			itemList, rl.UsedItems, slotList, rl.UsedSlots)

		startTime := time.Now()

		// slot progression items
		done := r.Graph["done"]
		success := true
		for done.GetMark(done, false) != graph.MarkTrue {
			if verbose {
				logChan <- fmt.Sprintf("searching; have %d more slots",
					slotList.Len())
			}

			// check to make sure this step isn't taking too long
			if time.Now().Sub(startTime) > time.Second*10 {
				success = false
				break
			}

			// try to find a new combination of items that opens progression
			items, slots := trySlotItemSet(r, src, itemList, slotList,
				countSteps, false)

			if items != nil {
				for items.Len() > 0 {
					rl.UsedItems.PushBack(items.Remove(items.Front()))
					slot := slots.Remove(slots.Front()).(*graph.Node)
					rl.UsedSlots.PushBack(slot)
					r.Costs += logic.Rupees[slot.Name]

					match := dungeonRegexp.FindStringSubmatch(slot.Name)
					if match != nil {
						di, _ := strconv.Atoi(match[1])
						r.DungeonItems[di]++
					}
				}
			} else {
				success = false
				break
			}

			r.Graph.ClearMarks()
		}

		if success {
			arrangeListsForLog(r, rl, verbose)

			// fill unused slots
			for slotList.Len() > 0 {
				if verbose {
					logChan <- fmt.Sprintf("done; filling %d more slots",
						slotList.Len())
				}

				// check to make sure this step isn't taking too long
				if time.Now().Sub(startTime) > time.Second*10 {
					break
				}

				items, slots := trySlotItemSet(r, src, itemList, slotList,
					countSteps, true)
				if items != nil {
					for items.Len() > 0 {
						item := items.Remove(items.Front())
						slot := slots.Remove(slots.Front())
						rl.UsedItems.PushBack(item)
						rl.UsedSlots.PushBack(slot)
						rl.OptionalItems.PushBack(item)
						rl.OptionalSlots.PushBack(slot)
					}
				} else {
					break
				}
			}
		}

		if slotList.Len() == 0 {
			// rotate dungeon items to the back of the lists
			items, slots := rl.RequiredItems, rl.RequiredSlots
			for i := 0; i < 6; i++ {
				items.PushBack(items.Remove(items.Front()))
				slots.PushBack(slots.Remove(slots.Front()))
			}
			items, slots = rl.OptionalItems, rl.OptionalSlots
			for i := 0; i < 16; i++ {
				items.PushBack(items.Remove(items.Front()))
				slots.PushBack(slots.Remove(slots.Front()))
			}

			// and we're done
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

		rl.UsedItems, rl.UsedSlots = list.New(), list.New()
		rl.RequiredItems, rl.RequiredSlots = list.New(), list.New()
		rl.OptionalItems, rl.OptionalSlots = list.New(), list.New()

		// get a new seed for the next iteration
		seed = uint32(src.Int31())
		src = rand.New(rand.NewSource(int64(seed)))
	}

	if tries >= maxTries {
		logChan <- fmt.Sprintf("abort; could not find route after %d tries",
			maxTries)
		return nil
	}

	return rl
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
func initRouteLists(src *rand.Rand, r *Route,
	companion int) (itemList, slotList *list.List) {
	// get slices of names
	itemNames := make([]string, 0,
		len(rom.ItemSlots)+len(logic.ExtraItems()))
	slotNames := make([]string, 0, len(r.Slots))
	for key, slot := range rom.ItemSlots {
		if key != "rod gift" { // don't slot vanilla, seasonless rod
			treasureName := rom.FindTreasureName(slot.Treasure)
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
			node.Name != "village shop 2" &&
			(r.TurnsReached[node] > 0 || canAffordSlot(r, node)) {
			count++
		}
	}
	return count
}

// break down the used items into required and optional items, so that the log
// makes sense.
func arrangeListsForLog(r *Route, rl *RouteLists, verbose bool) {
	done := r.Graph["done"]

	// figure out which items aren't necessary
	ei, es := rl.UsedItems.Front(), rl.UsedSlots.Front()
	for i := 0; i < rl.UsedItems.Len(); i++ {
		item, slot := ei.Value.(*graph.Node), es.Value.(*graph.Node)

		// remove parent provisionally
		item.RemoveParent(slot)

		// ask if anyone misses it
		r.Graph.ClearMarks()
		if logic.Rupees[item.Name] == 0 &&
			done.GetMark(done, false) == graph.MarkTrue {
			if verbose {
				fmt.Printf("%s (in %s) is extra\n", item.Name, slot.Name)
			}
			rl.OptionalItems.PushBack(item)
			rl.OptionalSlots.PushBack(slot)
		} else {
			item.AddParents(slot)
			rl.RequiredItems.PushBack(item)
			rl.RequiredSlots.PushBack(slot)
		}

		ei, es = ei.Next(), es.Next()
	}

	// attach removed parents back to optional items
	ei, es = rl.OptionalItems.Front(), rl.OptionalSlots.Front()
	for i := 0; i < rl.OptionalItems.Len(); i++ {
		item, slot := ei.Value.(*graph.Node), es.Value.(*graph.Node)
		item.AddParents(slot)
		ei, es = ei.Next(), es.Next()
	}
}
