package randomizer

import (
	"container/list"
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

// give up completely if routing fails too many times
const maxTries = 200

// names of portals from the subrosia side.
var subrosianPortalNames = map[string]string{
	"eastern suburbs":      "volcanoes east",
	"spool swamp":          "subrosia market",
	"mt. cucco":            "strange brothers",
	"eyeglass lake":        "great furnace",
	"horon village":        "house of pirates",
	"temple remains lower": "volcanoes west",
	"temple remains upper": "d8 entrance",
}

var dungeonNames = map[int][]string{
	gameSeasons: []string{
		"d0", "d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8"},
	gameAges: []string{
		"d1", "d2", "d3", "d4", "d5", "d6 present", "d6 past", "d7", "d8"},
}

// adds nodes to the map based on default contents of item slots.
func addDefaultItemNodes(rom *romState, nodes map[string]*prenode) {
	for _, slot := range rom.itemSlots {
		tName, _ := reverseLookup(rom.treasures, slot.treasure)
		nodes[tName.(string)] = rootPrenode()
	}
}

func addNodes(prenodes map[string]*prenode, g graph) {
	for key, pn := range prenodes {
		switch pn.nType {
		case andNode, orNode, rupeesNode:
			g[key] = newNode(key, pn.nType)
		case countNode:
			g[key] = newNode(key, countNode)
			g[key].minCount = pn.minCount
		default:
			panic("unknown logic type for " + key)
		}
	}
}

func addNodeParents(prenodes map[string]*prenode, g graph) {
	for k, pn := range prenodes {
		if g[k] == nil {
			continue
		}
		for _, parent := range pn.parents {
			if g[parent.(string)] == nil {
				continue
			}
			g.addParents(map[string][]string{k: []string{parent.(string)}})
		}
	}
}

type routeInfo struct {
	graph        graph
	slots        map[string]*node
	seed         uint32
	seasons      map[string]byte
	entrances    map[string]string
	portals      map[string]string
	companion    int // 1 to 3
	usedItems    *list.List
	usedSlots    *list.List
	ringMap      map[string]string
	attemptCount int
	src          *rand.Rand
}

const (
	ricky   = 1
	dimitri = 2
	moosh   = 3
)

func newRouteGraph(rom *romState) graph {
	g := newGraph()
	totalPrenodes := getPrenodes(rom.game)
	addDefaultItemNodes(rom, totalPrenodes)
	addNodes(totalPrenodes, g)
	addNodeParents(totalPrenodes, g)
	return g
}

// attempts to create a path to the given targets by placing different items in
// slots.
func findRoute(rom *romState, seed uint32, src *rand.Rand,
	ropts randomizerOptions, verbose bool, logf logFunc) (*routeInfo, error) {
	// make stacks out of the item names and slot names for backtracking
	var itemList, slotList *list.List

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	ri := &routeInfo{
		seed:      seed,
		usedItems: list.New(),
		usedSlots: list.New(),
		src:       src,
	}

	// try to find the route, retrying if needed
	tries := 0
	for tries = 0; tries < maxTries; tries++ {
		ri.graph = newRouteGraph(rom)
		ri.slots = make(map[string]*node, 0)
		for name := range rom.itemSlots {
			ri.slots[name] = ri.graph[name]
		}
		if ropts.hard {
			ri.graph["hard"].addParent(ri.graph["start"])
		}

		ri.companion = rollAnimalCompanion(ri.src, ri.graph, rom.game)
		ri.ringMap, _ = rom.randomizeRingPool(ri.src, nil)
		itemList, slotList = initRouteInfo(ri, rom)

		// attach free items to the "start" node until placed.
		for ei := itemList.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)
			ri.graph[item.name].addParent(ri.graph["start"])
		}

		// slot "world" nodes before items
		if rom.game == gameSeasons {
			ri.seasons = rollSeasons(ri.src, ri.graph)
			ri.portals = setPortals(ri.src, ri.graph, ropts.portals)
		}
		ri.entrances = setDungeonEntrances(
			ri.src, ri.graph, rom.game, ropts.dungeons)

		if tryPlaceItems(
			ri, itemList, slotList, rom.treasures, rom.game, verbose, logf) {
			ri.graph.reset()
			ri.graph["start"].explore()
			if ri.graph["done"].reached {
				// and we're done
				ri.attemptCount = tries + 1
				break
			} else if verbose {
				logf("all items placed but seed not completable")
			}
		}

		// clear placements and try again
		ri.usedItems, ri.usedSlots = list.New(), list.New()
	}

	if tries >= maxTries {
		return nil, fmt.Errorf("could not find route after %d tries", maxTries)
	}

	return ri, nil
}

var (
	seasonsById = []string{"spring", "summer", "autumn", "winter"}
	seasonAreas = []string{
		"north horon", "eastern suburbs", "woods of winter", "spool swamp",
		"holodrum plain", "sunken city", "lost woods", "tarm ruins",
		"western coast", "temple remains",
	}
)

// set the default seasons for all the applicable areas in the game, and return
// a mapping of area name to season value.
func rollSeasons(src *rand.Rand, g graph) map[string]byte {
	seasonMap := make(map[string]byte, len(seasonAreas))
	for _, area := range seasonAreas {
		id := src.Intn(len(seasonsById))
		season := seasonsById[id]
		g[fmt.Sprintf("%s default %s", area, season)].addParent(g["start"])
		seasonMap[area] = byte(id)
	}
	return seasonMap
}

// connect dungeon entrances, randomly or vanilla-ly.
func setDungeonEntrances(
	src *rand.Rand, g graph, game int, shuffle bool) map[string]string {
	dungeonEntranceMap := make(map[string]string)
	dungeons := make([]string, len(dungeonNames[game]))
	copy(dungeons, dungeonNames[game])
	if game == gameSeasons {
		dungeons = dungeons[1:]
	}

	if game == gameSeasons && !shuffle {
		g["d2 alt entrances enabled"].addParent(g["start"])
	}

	entrances := make([]string, len(dungeons))
	copy(entrances, dungeons)

	if shuffle {
		src.Shuffle(len(entrances), func(i, j int) {
			entrances[i], entrances[j] = entrances[j], entrances[i]
		})
	}

	for i := 0; i < len(dungeons); i++ {
		entranceName := fmt.Sprintf("%s entrance", entrances[i])
		dungeonEntranceMap[entrances[i]] = dungeons[i]
		g[fmt.Sprintf("enter %s", dungeons[i])].addParent(g[entranceName])
	}

	return dungeonEntranceMap
}

// connect subrosia portals, randomly or vanilla-ly.
func setPortals(src *rand.Rand, g graph, shuffle bool) map[string]string {
	portalMap := make(map[string]string)
	var portals = []string{
		"eastern suburbs", "spool swamp", "mt. cucco", "eyeglass lake",
		"horon village", "temple remains lower", "temple remains upper",
	}
	var connects = make([]string, len(portals))
	for i, portal := range portals {
		connects[i] = subrosianPortalNames[portal]
	}

	if shuffle {
		src.Shuffle(len(connects), func(i, j int) {
			connects[i], connects[j] = connects[j], connects[i]
		})
	}

	for i := 0; i < len(portals); i++ {
		portalMap[portals[i]] = connects[i]
		g[fmt.Sprintf("exit %s portal", connects[i])].
			addParent(g[fmt.Sprintf("enter %s portal", portals[i])])
		g[fmt.Sprintf("exit %s portal", portals[i])].
			addParent(g[fmt.Sprintf("enter %s portal", connects[i])])
	}

	return portalMap
}

// randomly determines animal companion and returns its ID (1 to 3)
func rollAnimalCompanion(src *rand.Rand, g graph, game int) int {
	companion := src.Intn(3) + 1

	if game == gameSeasons {
		switch companion {
		case ricky:
			g["natzu prairie"].addParent(g["start"])
		case dimitri:
			g["natzu river"].addParent(g["start"])
		case moosh:
			g["natzu wasteland"].addParent(g["start"])
		}
	} else {
		switch companion {
		case ricky:
			g["ricky nuun"].addParent(g["start"])
		case dimitri:
			g["dimitri nuun"].addParent(g["start"])
		case moosh:
			g["moosh nuun"].addParent(g["start"])
		}
	}

	return companion
}

var seedNames = []string{"ember tree seeds", "scent tree seeds",
	"pegasus tree seeds", "gale tree seeds", "mystery tree seeds"}

var seedTreeNames = map[string]bool{
	"horon village tree":      true,
	"woods of winter tree":    true,
	"north horon tree":        true,
	"spool swamp tree":        true,
	"sunken city tree":        true,
	"tarm ruins tree":         true,
	"south lynna tree":        true,
	"deku forest tree":        true,
	"crescent island tree":    true,
	"symmetry city tree":      true,
	"rolling ridge west tree": true,
	"rolling ridge east tree": true,
	"ambi's palace tree":      true,
	"zora village tree":       true,
}

// return shuffled lists of item and slot nodes
func initRouteInfo(
	ri *routeInfo, rom *romState) (itemList, slotList *list.List) {
	// get slices of names
	var itemNames []string
	slotNames := make([]string, 0, len(ri.slots))

	// get count of each seed tree from RNG
	nTrees := sora(rom.game, 6, 8).(int)
	thisSeeds := make([]int, 0, nTrees)
	seedCounts := make(map[int]int)
	for len(thisSeeds) < cap(thisSeeds) {
		id := ri.src.Intn(len(seedNames))
		for seedCounts[id] > len(seedCounts)/len(seedNames) {
			id = ri.src.Intn(len(seedNames))
		}
		thisSeeds = append(thisSeeds, id)
		seedCounts[id]++
	}

	for key, slot := range rom.itemSlots {
		switch {
		case seedTreeNames[key]:
			id := thisSeeds[0]
			thisSeeds = thisSeeds[1:]
			itemNames = append(itemNames, seedNames[id])
		default:
			// substitute identified flute for strange flute
			tName, _ := reverseLookup(rom.treasures, slot.treasure)
			treasureName := tName.(string)
			if strings.HasSuffix(treasureName, " flute") {
				switch ri.companion {
				case ricky:
					treasureName = "ricky's flute"
				case dimitri:
					treasureName = "dimitri's flute"
				case moosh:
					treasureName = "moosh's flute"
				}
			}

			// substitute ring pool
			if ringSub, ok := ri.ringMap[treasureName]; ok {
				treasureName = ringSub
			}

			itemNames = append(itemNames, treasureName)
		}
	}
	for key := range ri.slots {
		slotNames = append(slotNames, key)
	}

	// sort the slices so that order isn't dependent on map implementation,
	// then shuffle the sorted slices
	sort.Strings(itemNames)
	sort.Strings(slotNames)
	ri.src.Shuffle(len(itemNames), func(i, j int) {
		itemNames[i], itemNames[j] = itemNames[j], itemNames[i]
	})
	ri.src.Shuffle(len(slotNames), func(i, j int) {
		slotNames[i], slotNames[j] = slotNames[j], slotNames[i]
	})

	// push the graph nodes by name onto stacks
	itemList = list.New()
	slotList = list.New()
	for _, key := range itemNames {
		itemList.PushBack(ri.graph[key])
	}
	for _, key := range slotNames {
		slotList.PushBack(ri.graph[key])
	}

	return itemList, slotList
}

// returns true iff successful
func tryPlaceItems(ri *routeInfo, itemList, slotList *list.List,
	treasures map[string]*treasure, game int, verbose bool, logf logFunc) bool {
	for itemList.Len() > 0 && slotList.Len() > 0 {
		if verbose {
			logf("searching; filling %d more slots", slotList.Len())
			logf("(%d more items)", itemList.Len())
		}

		eItem, eSlot := trySlotRandomItem(
			ri.graph, ri.src, itemList, slotList, treasures, game)

		if eItem != nil {
			item := itemList.Remove(eItem).(*node)
			ri.usedItems.PushBack(item)
			slot := slotList.Remove(eSlot).(*node)
			ri.usedSlots.PushBack(slot)
			if verbose {
				logf("placing: %s <- %s", slot.name, item.name)
			}
		} else {
			if verbose {
				logf("search failed. unplaced items:")
				for ei := itemList.Front(); ei != nil; ei = ei.Next() {
					logf(ei.Value.(*node).name)
				}
				logf("unfilled slots:")
				for es := slotList.Front(); es != nil; es = es.Next() {
					logf(es.Value.(*node).name)
				}
			}
			return false
		}
	}
	return true
}

func trySlotRandomItem(g graph, src *rand.Rand, itemPool, slotPool *list.List,
	treasures map[string]*treasure, game int) (usedItem, usedSlot *list.Element) {
	// try placing the first item in a slot until it fits
	triedProgression := false
	for _, progressionItemsOnly := range []bool{true, false} {
		if !progressionItemsOnly && triedProgression {
			return nil, nil
		}

		for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)

			if progressionItemsOnly && itemIsInert(treasures, item.name) {
				continue
			}
			item.removeParent(g["start"])
			triedProgression = true

			for es := slotPool.Front(); es != nil; es = es.Next() {
				slot := es.Value.(*node)

				if !itemFitsInSlot(item, slot) {
					continue
				}

				// make sure enough space is left for remaining dungeon items
				if dungeonsOverfilled(game, ei, es, itemPool, slotPool) {
					continue
				}

				// test whether seed is still beatable w/ item placement
				g.reset()
				item.addParent(slot)
				g["start"].explore()
				if !g["done"].reached {
					item.removeParent(slot)
					continue
				}

				// make sure item didn't cause a forward-wise dead end
				if isDeadEnd(g, ei, es, itemPool, slotPool) {
					item.removeParent(slot)
					continue
				}

				return ei, es
			}

			item.addParent(g["start"])
		}
	}

	return nil, nil
}

// checks whether the item fits in the slot due to things like seeds only going
// in trees, certain item slots not accomodating sub IDs. this doesn't check
// for softlocks or the availability of the slot and item.
func itemFitsInSlot(itemNode, slotNode *node) bool {
	// dummy shop slots 1 and 2 can only hold their vanilla items.
	switch {
	case slotNode.name == "shop, 20 rupees" && itemNode.name != "bombs, 10":
		fallthrough
	case slotNode.name == "shop, 30 rupees" && itemNode.name != "wooden shield":
		fallthrough
	case itemNode.name == "wooden shield" && slotNode.name != "shop, 30 rupees":
		return false
	}

	// bomb flower has special graphics something. this could probably be
	// worked around like with the temple of seasons, but i'm not super
	// interested in doing that.
	if itemNode.name == "bomb flower" {
		switch slotNode.name {
		case "cheval's test", "cheval's invention", "wild tokay game",
			"hidden tokay cave", "library present", "library past":
			return false
		}
	}

	// dungeons can only hold their respective dungeon-specific items. the
	// HasPrefix is specifically for ages d6 boss key.
	dungeonName := getDungeonName(itemNode.name)
	if dungeonName != "" &&
		!strings.HasPrefix(getDungeonName(slotNode.name), dungeonName) {
		return false
	}

	// and only seeds can be slotted in seed trees, of course
	switch itemNode.name {
	case "ember tree seeds", "mystery tree seeds", "scent tree seeds",
		"pegasus tree seeds", "gale tree seeds":
		return seedTreeNames[slotNode.name]
	default:
		return !seedTreeNames[slotNode.name]
	}
}

// return the name of a dungeon associated with a given item or slot name. ages
// d6 boss key returns "d6". non-dungeon names return "".
func getDungeonName(name string) string {
	if strings.HasPrefix(name, "d6 present") {
		return "d6 present"
	} else if strings.HasPrefix(name, "d6 past") {
		return "d6 past"
	} else if strings.HasPrefix(name, "maku path") {
		return "d0"
	} else if name == "slate" {
		return "d8"
	}

	switch name[:2] {
	case "d0", "d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8":
		return name[:2]
	default:
		return ""
	}
}

// returns true iff no open slots beyond curSlot are reachable if all the items
// left in the pool, except for curItem, are assumed to be unreachable. returns
// false if only one slot remains in the pool, since that slot is assumed to be
// curSlot.
func isDeadEnd(g graph, curItem, curSlot *list.Element,
	itemPool, slotPool *list.List) bool {
	if slotPool.Len() == 1 {
		return false
	}

	for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
		if ei != curItem {
			ei.Value.(*node).removeParent(g["start"])
		}
	}
	g.reset()
	g["start"].explore()

	dead := true
	for es := slotPool.Front(); es != nil; es = es.Next() {
		if es != curSlot && es.Value.(*node).reached {
			dead = false
			break
		}
	}

	for ei := itemPool.Front(); ei != nil; ei = ei.Next() {
		if ei != curItem {
			ei.Value.(*node).addParent(g["start"])
		}
	}

	return dead
}

// returns true iff there are more items specific to any dungeon than there are
// slots remaining in that dungeon. elements item and slot are not counted.
func dungeonsOverfilled(game int, item, slot *list.Element,
	itemPool, slotPool *list.List) bool {
	for _, name := range dungeonNames[game] {
		// ages d6 boss key isn't correctly accounted for here. oh well.
		nItems := countList(itemPool, func(e *list.Element) bool {
			return e != item && getDungeonName(e.Value.(*node).name) == name
		})
		nSlots := countList(slotPool, func(e *list.Element) bool {
			return e != slot && getDungeonName(e.Value.(*node).name) == name
		})
		if nItems > nSlots {
			return true
		}
	}
	return false
}

// returns the number of elements in the list for which the given function
// returns true.
func countList(l *list.List, f func(*list.Element) bool) int {
	n := 0
	for e := l.Front(); e != nil; e = e.Next() {
		if f(e) {
			n++
		}
	}
	return n
}

// itemIsInert returns true iff the item with the given name can never be
// progression, regardless of context.
func itemIsInert(treasures map[string]*treasure, name string) bool {
	switch name {
	case "fist ring", "expert's ring", "energy ring", "toss ring",
		"swimmer's ring":
		return false
	}

	// non-default junk rings
	if treasures[name] == nil {
		return true
	}

	// not part of next switch since the ID is only junk in seasons
	if name == "treasure map" {
		return true
	}

	switch treasures[name].id {
	// heart refill, PoH, HC, ring, compass, dungeon map, gasha seed
	case 0x29, 0x2a, 0x2b, 0x2d, 0x32, 0x33, 0x34:
		return true
	}
	return false
}

// moves the first matching string in the slice to the end of the slice.
func moveStringToBack(a []string, s string) {
	for i, s2 := range a {
		if s2 == s {
			a = append(a[:i], append(a[i+1:], s)...)
		}
	}
}

// returns true iff a is a slice and v is a value in that slice. panics if a is
// not a slice.
func sliceContains(a interface{}, v interface{}) bool {
	aValue := reflect.ValueOf(a)
	for i := 0; i < aValue.Len(); i++ {
		v2 := aValue.Index(i).Interface()
		if reflect.DeepEqual(v, v2) {
			return true
		}
	}
	return false
}

// return alphabetically sorted string values from a map.
func orderedValues(m map[string]string) []string {
	a, i := make([]string, len(m)), 0
	for _, v := range m {
		a[i] = v
		i++
	}
	sort.Strings(a)
	return a
}
