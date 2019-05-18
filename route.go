package main

import (
	"container/list"
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/jangler/oracles-randomizer/rom"
)

// give up completely if routing fails too many times
const maxTries = 50

// names of portals from the subrosia side. log and logic care about these, but
// rom code doesn't.
var subrosianPortalNames = map[string]string{
	"eastern suburbs":      "volcanoes east",
	"spool swamp":          "subrosia market",
	"mt. cucco":            "strange brothers",
	"eyeglass lake":        "great furnace",
	"horon village":        "house of pirates",
	"temple remains lower": "volcanoes west",
	"temple remains upper": "D8 entrance",
}

// adds nodes to the map based on default contents of item slots.
func addDefaultItemNodes(nodes map[string]*prenode) {
	for _, slot := range rom.ItemSlots {
		nodes[rom.FindTreasureName(slot.Treasure)] = rootPrenode()
	}
}

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph graph
	Slots map[string]*node
}

// NewRoute returns an initialized route with all nodes.
func NewRoute(game int) *Route {
	g := newGraph()

	totalPrenodes := getPrenodes(game)
	addDefaultItemNodes(totalPrenodes)

	addNodes(totalPrenodes, g)
	addNodeParents(totalPrenodes, g)

	openSlots := make(map[string]*node, 0)
	for name := range rom.ItemSlots {
		openSlots[name] = g[name]
	}

	return &Route{
		Graph: g,
		Slots: openSlots,
	}
}

func (r *Route) AddParent(child, parent string) {
	r.Graph[child].addParent(r.Graph[parent])
}

func (r *Route) ClearParents(node string) {
	r.Graph[node].clearParents()
}

func addNodes(prenodes map[string]*prenode, g graph) {
	for key, pn := range prenodes {
		switch pn.nType {
		case andNode:
			g[key] = newNode(key, andNode)
		case orNode:
			g[key] = newNode(key, orNode)
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

type RouteInfo struct {
	Route                *Route
	Seed                 uint32
	Seasons              map[string]byte
	Entrances            map[string]string
	Portals              map[string]string
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
func findRoute(game int, seed uint32, ropts randomizerOptions, verbose bool,
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
		if ropts.hard {
			r.AddParent("hard", "start")
		}
		ri.Companion = rollAnimalCompanion(ri.Src, r, game)
		ri.RingMap = rom.RandomizeRingPool(ri.Src, game)
		itemList, slotList = initRouteInfo(ri.Src, r, ri.RingMap, game,
			ri.Companion)

		// attach free items to start node - just assume we have them, until
		// they're placed
		for ei := itemList.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)
			r.AddParent(item.name, "start")
		}

		// slot "world" nodes before items
		if game == rom.GameSeasons {
			ri.Seasons = rollSeasons(ri.Src, r)
			ri.Portals = setPortals(ri.Src, r, ropts.portals)
		}
		ri.Entrances = setDungeonEntrances(ri.Src, r, game, ropts.dungeons)

		// place dungeon-specific items, then "regular" items
		dungeonItems, nonDungeonItems := list.New(), list.New()
		for ei := itemList.Front(); ei != nil; ei = ei.Next() {
			item := ei.Value.(*node)
			if getDungeonName(item.name) != "" {
				dungeonItems.PushBack(item)
			} else {
				nonDungeonItems.PushBack(item)
			}
		}
		if tryPlaceItems(ri, r, dungeonItems, slotList, verbose, logf) &&
			tryPlaceItems(ri, r, nonDungeonItems, slotList, verbose, logf) {
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

// connect dungeon entrances, randomly or vanilla-ly.
func setDungeonEntrances(src *rand.Rand,
	r *Route, game int, shuffle bool) map[string]string {
	dungeonEntranceMap := make(map[string]string)
	var dungeons []string

	if game == rom.GameSeasons {
		dungeons = []string{"d1", "d2", "d3", "d4", "d5", "d6", "d7", "d8"}
		if !shuffle {
			r.ClearParents("d2 alt entrances enabled")
			r.AddParent("d2 alt entrances enabled", "start")
		}
	} else {
		dungeons = []string{"d1", "d2", "d3", "d4", "d5",
			"d6 present", "d6 past", "d7", "d8"}
	}

	// reset entrances
	for _, dungeon := range dungeons {
		r.ClearParents(fmt.Sprintf("enter %s", dungeon))
	}

	var entrances = make([]string, len(dungeons))
	copy(entrances, dungeons)

	if shuffle {
		src.Shuffle(len(entrances), func(i, j int) {
			entrances[i], entrances[j] = entrances[j], entrances[i]
		})
	}

	for i := 0; i < len(dungeons); i++ {
		entranceName := fmt.Sprintf("%s entrance", entrances[i])
		dungeonEntranceMap[entrances[i]] = dungeons[i]
		r.AddParent(fmt.Sprintf("enter %s", dungeons[i]), entranceName)
	}

	return dungeonEntranceMap
}

// connect subrosia portals, randomly or vanilla-ly.
func setPortals(src *rand.Rand, r *Route, shuffle bool) map[string]string {
	portalMap := make(map[string]string)
	var portals = []string{
		"eastern suburbs", "spool swamp", "mt. cucco", "eyeglass lake",
		"horon village", "temple remains lower", "temple remains upper",
	}

	// reset exits
	for _, portal := range portals {
		r.ClearParents(fmt.Sprintf("exit %s portal", portal))
		r.ClearParents(fmt.Sprintf("exit %s portal",
			subrosianPortalNames[portal]))
	}

	var connects = make([]string, len(portals))
	copy(connects, portals)

	if shuffle {
		src.Shuffle(len(connects), func(i, j int) {
			connects[i], connects[j] = connects[j], connects[i]
		})
	}

	for i := 0; i < len(portals); i++ {
		portalMap[portals[i]] = connects[i]
		r.AddParent(fmt.Sprintf("exit %s portal", subrosianPortalNames[connects[i]]),
			fmt.Sprintf("enter %s portal", portals[i]))
		r.AddParent(fmt.Sprintf("exit %s portal", portals[i]),
			fmt.Sprintf("enter %s portal", subrosianPortalNames[connects[i]]))
	}

	return portalMap
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

var seedNames = []string{"ember tree seeds", "scent tree seeds",
	"pegasus tree seeds", "gale tree seeds", "mystery tree seeds"}

// return shuffled lists of item and slot nodes
func initRouteInfo(src *rand.Rand, r *Route, ringMap map[string]string,
	game, companion int) (itemList, slotList *list.List) {
	// get slices of names
	var itemNames []string
	if game == rom.GameSeasons {
		// TODO: do this differently. like put it in a regular slot. also does
		// this actually work like it's supposed to?
		itemNames = make([]string, 0, len(rom.ItemSlots)+1) // +1 for fool's ore
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
			if strings.HasSuffix(treasureName, " flute") {
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
		itemNames = append(itemNames, "fool's ore")
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

// returns true iff successful
func tryPlaceItems(ri *RouteInfo, r *Route, itemList, slotList *list.List,
	verbose bool, logf logFunc) bool {
	for itemList.Len() > 0 {
		if verbose {
			logf("searching; filling %d more slots", slotList.Len())
			logf("(%d more items)", itemList.Len())
		}

		eItem, eSlot := trySlotRandomItem(r, ri.Src, itemList, slotList)

		if eItem != nil {
			item := itemList.Remove(eItem).(*node)
			ri.UsedItems.PushBack(item)
			slot := slotList.Remove(eSlot).(*node)
			ri.UsedSlots.PushBack(slot)
			if verbose {
				logf("placing: %s <- %s", slot.name, item.name)
			}
		} else {
			if verbose {
				logf("search failed. unplaced non-junk items:")
				for ei := itemList.Front(); ei != nil; ei = ei.Next() {
					item := ei.Value.(*node)
					if !itemIsJunk(item.name) {
						logf(item.name)
					}
				}
			}
			return false
		}
	}

	return true
}
