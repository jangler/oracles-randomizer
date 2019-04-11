package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/jangler/oracles-randomizer/graph"
	"github.com/jangler/oracles-randomizer/logic"
	"github.com/jangler/oracles-randomizer/rom"
)

// give up completely if routing fails too many times
const maxTries = 50

// adds nodes to the map based on default contents of item slots.
func addDefaultItemNodes(nodes map[string]*logic.Node) {
	for key, slot := range rom.ItemSlots {
		if key != "temple of seasons" { // real rod is an Or, not a Root
			nodes[rom.FindTreasureName(slot.Treasure)] = logic.Root()
		}
	}
}

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph  graph.Graph
	Slots  map[string]*graph.Node
	Rupees int
}

// NewRoute returns an initialized route with all nodes, and those nodes with
// the names in start functioning as givens (always satisfied). If no names are
// given, only the normal start node functions as a given.
func NewRoute(game int, start ...string) *Route {
	g := graph.New()

	var totalPrenodes map[string]*logic.Node
	if game == rom.GameSeasons {
		totalPrenodes = logic.GetSeasons()
	} else {
		totalPrenodes = logic.GetAges()
	}
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
		case logic.CountType:
			node := graph.NewNode(key, graph.CountType, false, false, false)
			node.MinCount = pn.MinCount
			g.AddNodes(node)
		default:
			panic("unknown logic type for " + key)
		}
	}
}

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
	Route                *Route
	Seed                 uint32
	Seasons              map[string]byte
	Companion            int // 1 to 3
	UsedItems, UsedSlots *list.List
	RingMap              map[string]string
	AttemptCount         int
	Src                  *rand.Rand
}

const (
	ricky   = 1
	dimitri = 2
	moosh   = 3
)

// attempts to create a path to the given targets by placing different items in
// slots. returns nils if no route is found.
func findRoute(game int, seed uint32, hard, verbose bool,
	logf logFunc) *RouteInfo {
	// make stacks out of the item names and slot names for backtracking
	var itemList, slotList *list.List

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	ri := &RouteInfo{
		Seed:      seed,
		UsedItems: list.New(),
		UsedSlots: list.New(),
	}

	// try to find the route, retrying if needed
	tries := 0
	for tries = 0; tries < maxTries; tries++ {
		ri.Src = rand.New(rand.NewSource(int64(ri.Seed)))
		logf("trying seed %08x", ri.Seed)

		r := NewRoute(game)
		ri.Companion = rollAnimalCompanion(ri.Src, r, game)
		ri.RingMap = rom.RandomizeRingPool(ri.Src, game)
		itemList, slotList = initRouteInfo(ri.Src, r, ri.RingMap, game,
			ri.Companion)

		// slot initial nodes before algorithm slots progression items
		if game == rom.GameSeasons {
			ri.Seasons = rollSeasons(ri.Src, r)
		}
		placeDungeonItems(ri.Src, r, game, hard,
			itemList, ri.UsedItems, slotList, ri.UsedSlots)

		slotRecord := 0
		i, maxIterations := 0, 1+itemList.Len()

		// slot progression items
		done := r.Graph["done"]
		success := true
		for done.GetMark(done, hard) != graph.MarkTrue {
			if verbose {
				logf("searching; have %d more slots", slotList.Len())
				logf("%d/%d iterations", i, maxIterations)
			}

			eItem, eSlot := trySlotRandomItem(r, ri.Src, itemList, slotList,
				countSteps, ri.UsedSlots.Len(), hard, false)

			if eItem != nil {
				item := itemList.Remove(eItem).(*graph.Node)
				ri.UsedItems.PushBack(item)
				slot := slotList.Remove(eSlot).(*graph.Node)
				ri.UsedSlots.PushBack(slot)
				r.Rupees += logic.RupeeValues[item.Name]

				if ri.UsedSlots.Len() > slotRecord {
					slotRecord = ri.UsedSlots.Len()
					i, maxIterations = 0, 1+itemList.Len()
				}
			} else {
				item := ri.UsedItems.Remove(ri.UsedItems.Back()).(*graph.Node)
				slot := ri.UsedSlots.Remove(ri.UsedSlots.Back()).(*graph.Node)
				r.Rupees -= logic.RupeeValues[item.Name]
				itemList.PushBack(item)
				slotList.PushBack(slot)
				item.RemoveParent(slot)
			}

			r.Graph.ClearMarks()

			i++
			if i > maxIterations {
				success = false
				if verbose {
					logf("maximum iterations reached")
				}
				break
			}
		}

		if success {
			// fill unused slots
			for slotList.Len() > 0 {
				if verbose {
					logf("done; filling %d more slots", slotList.Len())
					logf("%d/%d iterations", i, maxIterations)
				}

				eItem, eSlot := trySlotRandomItem(r, ri.Src, itemList, slotList,
					countSteps, ri.UsedSlots.Len(), hard, true)

				if eItem != nil {
					item := itemList.Remove(eItem).(*graph.Node)
					ri.UsedItems.PushBack(item)
					slot := slotList.Remove(eSlot).(*graph.Node)
					ri.UsedSlots.PushBack(slot)
					r.Rupees += logic.RupeeValues[item.Name]

					if ri.UsedSlots.Len() > slotRecord {
						slotRecord = ri.UsedSlots.Len()
						i, maxIterations = 0, 1+itemList.Len()
					}
				} else {
					item := ri.UsedItems.Remove(ri.UsedItems.Back()).(*graph.Node)
					slot := ri.UsedSlots.Remove(ri.UsedSlots.Back()).(*graph.Node)
					r.Rupees -= logic.RupeeValues[item.Name]
					itemList.PushBack(item)
					slotList.PushBack(slot)
					item.RemoveParent(slot)
				}

				i++
				if i > maxIterations {
					if verbose {
						logf("maximum iterations reached")
					}
					break
				}
			}
		}

		if success && slotList.Len() == 0 {
			// and we're done
			ri.Route = r
			ri.AttemptCount = tries + 1
			break
		}

		ri.UsedItems, ri.UsedSlots = list.New(), list.New()

		// get a new seed for the next iteration
		ri.Seed = uint32(ri.Src.Int31())
	}

	if tries >= maxTries {
		logf("abort; could not find route after %d tries", maxTries)
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
func rollAnimalCompanion(src *rand.Rand, r *Route, game int) int {
	companion := src.Intn(3) + 1

	if game == rom.GameSeasons {
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
	} else {
		r.ClearParents("ricky nuun")
		r.ClearParents("dimitri nuun")
		r.ClearParents("moosh nuun")

		switch companion {
		case ricky:
			r.AddParent("ricky nuun", "start")
		case dimitri:
			r.AddParent("dimitri nuun", "start")
		case moosh:
			r.AddParent("moosh nuun", "start")
		}
	}

	return companion
}

// place maps, compasses, small keys, boss keys, and slates in chests in
// dungeons (before attempting to slot the other items).
func placeDungeonItems(src *rand.Rand, r *Route, game int, hard bool,
	itemList, usedItems, slotList, usedSlots *list.List) {
	g := r.Graph

	// place boss keys first
	for i := 1; i < 9; i++ {
		prefix := fmt.Sprintf("d%d", i)
		itemName := prefix + " boss key"

		slotElem, itemElem, slotNode, itemNode :=
			getDungeonItem(prefix, itemName, slotList, itemList, g, hard)
		placeItem(slotNode, itemNode, slotElem, itemElem,
			usedSlots, slotList, usedItems, itemList)
	}

	// then place small keys
	for i := 0; i < 9; i++ {
		prefix := fmt.Sprintf("d%d", i)
		itemName := prefix + " small key"

		for {
			slotElem, itemElem, slotNode, itemNode :=
				getDungeonItem(prefix, itemName, slotList, itemList, g, hard)
			if itemNode == nil {
				// no more small keys to place for this dungeon
				break
			}

			placeItem(slotNode, itemNode, slotElem, itemElem,
				usedSlots, slotList, usedItems, itemList)
		}
	}

	// place slates in ages
	if game == rom.GameAges {
		for i := 1; i <= 4; i++ {
			itemName := fmt.Sprintf("slate %d", i)
			slotElem, itemElem, slotNode, itemNode :=
				getDungeonItem("d8", itemName, slotList, itemList, g, hard)
			placeItem(slotNode, itemNode, slotElem, itemElem,
				usedSlots, slotList, usedItems, itemList)
		}
	}

	prefixes := []string{"d1", "d2", "d3", "d4", "d5"}
	if game == rom.GameSeasons {
		prefixes = append(prefixes, "d6")
	} else {
		prefixes = append(prefixes, "d6 present", "d6 past")
	}
	prefixes = append(prefixes, "d7", "d8")

	// then place maps and compasses
	for _, prefix := range prefixes {
		for _, itemName := range []string{"dungeon map", "compass"} {
			slotElem, itemElem, slotNode, itemNode :=
				getDungeonItem(prefix, itemName, slotList, itemList, g, hard)
			placeItem(slotNode, itemNode, slotElem, itemElem,
				usedSlots, slotList, usedItems, itemList)
		}
	}
}

// find a valid position for a dungeon item
func getDungeonItem(prefix, itemName string, slotList, itemList *list.List,
	g graph.Graph, hard bool) (slotElem, itemElem *list.Element, slotNode, itemNode *graph.Node) {
	for es := slotList.Front(); es != nil; es = es.Next() {
		slot := es.Value.(*graph.Node)
		if !strings.HasPrefix(slot.Name, prefix) {
			continue
		}
		if (strings.HasSuffix(itemName, "boss key") ||
			strings.HasPrefix(itemName, "slate")) &&
			strings.HasSuffix(slot.Name, "boss") {
			continue
		}
		if strings.HasSuffix(itemName, "small key") &&
			!canReachViaKeys(g, slot, hard) {
			continue
		}

		for ei := itemList.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*graph.Node)
			if item.Name != itemName {
				continue
			}

			return es, ei, slot, item
		}

		// return nil when there are no more small keys to place, since this is
		// how the caller determines whether it needs to place more keys.
		if slot != nil && strings.HasSuffix(itemName, "small key") {
			return nil, nil, nil, nil
		}
	}

	panic("could not place dungeon-specific items")
}

// place item in the given slot, and remove it from the pool
func placeItem(slotNode, itemNode *graph.Node,
	slotElem, itemElem *list.Element,
	usedSlots, slotList, usedItems, itemList *list.List) {
	usedSlots.PushBack(slotNode)
	slotList.Remove(slotElem)
	usedItems.PushBack(itemNode)
	itemList.Remove(itemElem)

	itemNode.AddParents(slotNode)
}

// returns true iff the target node can be reached if the player has automatic
// access to every item that isn't a small key or boss key.
func canReachViaKeys(g graph.Graph, target *graph.Node, hard bool) bool {
	g.ClearMarks()

	for _, itemSlot := range rom.ItemSlots {
		treasureName := rom.FindTreasureName(itemSlot.Treasure)
		if !strings.HasSuffix(treasureName, "small key") &&
			!strings.HasSuffix(treasureName, "boss key") {
			g[treasureName].Mark = graph.MarkTrue
		}
	}

	return target.GetMark(target, hard) == graph.MarkTrue
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
func initRouteInfo(src *rand.Rand, r *Route, ringMap map[string]string,
	game, companion int) (itemList, slotList *list.List) {
	// get slices of names
	var itemNames []string
	if game == rom.GameSeasons {
		itemNames = make([]string, 0,
			len(rom.ItemSlots)+len(logic.SeasonsExtraItems()))
	} else {
		itemNames = make([]string, 0, len(rom.ItemSlots))
	}
	slotNames := make([]string, 0, len(r.Slots))
	thisSeedNames := make([]string, len(seedNames))
	copy(thisSeedNames, seedNames)
	for key, slot := range rom.ItemSlots {
		switch key {
		case "temple of seasons": // don't slot vanilla, seasonless rod
			break
		case "tarm ruins seed tree", "ambi's palace tree",
			"rolling ridge east tree", "zora village tree":
			// use random duplicate seed types, but only duplicate a seed type
			// once
			index := src.Intn(len(thisSeedNames))
			treasureName := thisSeedNames[index]
			itemNames = append(itemNames, treasureName)
			thisSeedNames = append(thisSeedNames[:index],
				thisSeedNames[index+1:]...)
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

			// substitute ring pool
			if ringSub, ok := ringMap[treasureName]; ok {
				treasureName = ringSub
			}

			itemNames = append(itemNames, treasureName)
		}
	}
	if game == rom.GameSeasons {
		for key := range logic.SeasonsExtraItems() {
			itemNames = append(itemNames, key)
		}
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
func countSteps(r *Route, hard bool) int {
	r.Graph.ClearMarks()
	reached := r.Graph.ExploreFromStart(hard)
	count := 0
	for node := range reached {
		if node.IsStep && canAffordSlot(r, node, hard) {
			count++
		}
	}
	return count
}
