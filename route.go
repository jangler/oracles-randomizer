package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

// this file contains the actual connection of nodes in the game graph, and
// tracks them as they update.

// XXX need to be careful about rings. i can't imagine a situation where you'd
//     need both energy ring and fist ring, but if you did, then you'd need to
//     have the L-2 ring box to do so without danger of soft locking.

type PointType int

const (
	RootType PointType = iota
	AndType
	OrType
	AndSlotType
	OrSlotType
	AndStepType
	OrStepType
)

// A Point is a mapping of point strings that will become And or Or nodes in
// the graph.
type Point struct {
	Parents []string
	Type    PointType
}

// the different types of points are all just string slices; the reason for
// having different ones is purely for type assertions.
//
// And, Or, and Root are pretty self-explanatory; one with a Slot suffix is
// an item slot; one with a Step suffix is treated as a milestone for routing
// purposes. Slot types are also treated as steps; see the Point.IsStep()
// function.

func Root(a ...string) Point { return Point{a, RootType} }

func And(a ...string) Point { return Point{a, AndType} }

func Or(a ...string) Point { return Point{a, OrType} }

func AndSlot(a ...string) Point { return Point{a, AndSlotType} }

func OrSlot(a ...string) Point { return Point{a, OrSlotType} }

func AndStep(a ...string) Point { return Point{a, AndStepType} }

func OrStep(a ...string) Point { return Point{a, OrStepType} }

func (p *Point) IsStep() bool {
	switch p.Type {
	case AndSlotType, OrSlotType, AndStepType, OrStepType:
		return true
	}
	return false
}

type Route struct {
	Graph  *graph.Graph
	Points map[string]Point
	Slots  map[string]Point
}

// total
var nonGeneratedPoints map[string]Point

func init() {
	nonGeneratedPoints = make(map[string]Point)
	appendPoints(nonGeneratedPoints,
		baseItemPoints, ignoredBaseItemPoints,
		itemPoints, killPoints,
		holodrumPoints, subrosiaPoints, portalPoints,
		d0Points, d1Points, d2Points, d3Points, d4Points,
		d5Points, d6Points, d7Points, d8Points, d9Points,
	)
}

func initRoute(start []string) *Route {
	g := graph.NewGraph()

	totalPoints := make(map[string]Point, 0)
	appendPoints(totalPoints, nonGeneratedPoints, generatedPoints)

	// ignore semicolon-delimited points; they're only used for generation
	for key := range totalPoints {
		if strings.ContainsRune(key, ';') {
			delete(totalPoints, key)
		}
	}

	// make start nodes given
	for _, key := range start {
		totalPoints[key] = And()
	}

	addPointNodes(g, totalPoints)
	addPointParents(g, totalPoints)

	openSlots := make(map[string]Point, 0)
	for name, point := range totalPoints {
		switch point.Type {
		case AndSlotType, OrSlotType:
			openSlots[name] = point
		}
	}

	return &Route{Graph: g, Points: totalPoints, Slots: openSlots}
}

// CheckGraph returns an error for each orphan and childless node in the graph,
// ignoring nodes which are *supposed* to be orphans or childless. If there are
// no errors, it returns nil.
func (r *Route) CheckGraph() []error {
	var errs []error

	for name, node := range r.Graph.Map {
		// check for parents and children
		if len(node.Parents()) == 0 {
			// root nodes are supposed to be parentless
			if r.Points[name].Type == RootType {
				// it's supposed to be orphan/childless; skip it
				continue
			}

			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("orphan node: %s", name))
		}
		if len(node.Children()) == 0 {
			// item slots are supposed to be childless
			switch r.Points[name].Type {
			case AndSlotType, OrSlotType:
				continue
			}

			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("childless node: %s", name))
		}
	}

	return errs
}

func appendPoints(total map[string]Point, pointMaps ...map[string]Point) {
	for _, pointMap := range pointMaps {
		for k, v := range pointMap {
			total[k] = v
		}
	}
}

func addPointNodes(g *graph.Graph, points map[string]Point) {
	for key, pt := range points {
		switch pt.Type {
		case AndType, AndSlotType, AndStepType:
			g.AddNodes(graph.NewAndNode(key, pt.IsStep()))
		case OrType, OrSlotType, OrStepType, RootType:
			g.AddNodes(graph.NewOrNode(key, pt.IsStep()))
		default:
			panic("unknown point type for " + key)
		}
	}
}

func addPointParents(g *graph.Graph, points map[string]Point) {
	// ugly but w/e
	for k, p := range points {
		g.AddParents(map[string][]string{k: p.Parents})
	}
}

// attempts to find a path from the start to the given node in the graph.
// returns nil if no path was found.
func findPath(g *graph.Graph, target graph.Node) *list.List {
	path := list.New()
	mark := target.GetMark(path)
	if mark == graph.MarkTrue {
		return path
	}
	return nil
}

// attempts to create a path to the given targets by placing different items in
// slots.
func makeRoute(r *Route, start, goal, forbid []string,
	maxlen int) (usedItems, usedSlots, itemList, slotList *list.List) {
	// make stacks out of the item names and slot names for backtracking
	itemList = list.New()
	slotList = list.New()
	{
		// shuffle names in slices
		items := make([]graph.Node, 0, len(baseItemPoints))
		slots := make([]graph.Node, 0, len(r.Slots))
		for itemName, _ := range baseItemPoints {
			items = append(items, r.Graph.Map[itemName])
		}
		for slotName, _ := range r.Slots {
			slots = append(slots, r.Graph.Map[slotName])
		}
		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		rand.Shuffle(len(slots), func(i, j int) {
			slots[i], slots[j] = slots[j], slots[i]
		})

		// push the shuffled items onto the stacks
		for _, itemNode := range items {
			itemList.PushBack(itemNode)
		}
		for _, slotNode := range slots {
			slotList.PushBack(slotNode)
		}
	}

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems = list.New()
	usedSlots = list.New()

	// convert name lists into node lists
	startNodes := make([]graph.Node, len(start))
	for i, name := range start {
		startNodes[i] = r.Graph.Map[name]
	}
	goalNodes := make([]graph.Node, len(goal))
	for i, name := range goal {
		goalNodes[i] = r.Graph.Map[name]
	}
	forbidNodes := make([]graph.Node, len(forbid))
	for i, name := range forbid {
		forbidNodes[i] = r.Graph.Map[name]
	}

	if tryExploreTargets(r.Graph, nil, startNodes, goalNodes,
		forbidNodes, maxlen, itemList, usedItems, slotList, usedSlots) {
		log.Print("-- success")
		for _, target := range goal {
			r.Graph.ClearMarks()
			if !canReachTargets(r.Graph, target) {
				log.Fatalf("fatal: no path to %s!", target)
			}
		}
		log.Print("-- slotted items")
		if usedItems.Len() != usedSlots.Len() {
			log.Fatalf("FATAL: usedItems.Len() == %d; usedSlots.Len() == %d", usedItems.Len(), usedSlots.Len())
		}
		for i := 0; i < usedItems.Len(); i++ {
			log.Printf("%v <- %v",
				usedItems.Front().Value.(graph.Node),
				usedSlots.Front().Value.(graph.Node))
			usedItems.MoveToBack(usedItems.Front())
			usedSlots.MoveToBack(usedSlots.Front())
		}
	} else {
		log.Fatal("-- could not find route")
	}

	return
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
//
// the lists are lists of nodes.
func tryExploreTargets(g *graph.Graph, start map[graph.Node]bool,
	add, goal, forbid []graph.Node, maxlen int,
	itemList, usedItems, slotList, usedSlots *list.List) bool {
	// explore given the old state and changes
	reached := g.Explore(start, add, nil)
	log.Print(len(reached), " steps reached")

	// abort if any forbidden node is reached
	for _, node := range forbid {
		if reached[node] {
			log.Printf("-- false; reached forbidden node %s", node)
			return false
		}
	}

	// success if all goal nodes are reached
	allReached := true
	for _, node := range goal {
		if !reached[node] {
			log.Printf("-- have not reached goal node %s", node)
			allReached = false
			break
		}
	}
	if allReached {
		return true
	}

	// if the new state doesn't reach any more steps, abandon this branch
	if countSteps(reached) <= countSteps(start) {
		log.Printf("-- false; reached steps %d <= start steps %d",
			countSteps(reached), countSteps(start))
		return false
	}

	// can't slot any more items
	if maxlen == 0 {
		log.Print("-- false; slotted maxlen items")
		return false
	}

	// TODO: check softlocks

	// TODO: move some of this logic to its own function(s)
	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slotElem := slotList.Back()
		slotList.MoveToFront(slotElem)

		// see if slot node has been reached
		slotNode := slotElem.Value.(graph.Node)
		if !reached[slotNode] {
			continue
		}

		// move slot from unused to used
		usedSlots.PushBack(slotNode)
		slotList.Remove(slotElem)

		// try placing each unused item into the slot
		jewelChecked := false
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemNode := itemList.Remove(itemList.Back()).(graph.Node)
			usedItems.PushBack(itemNode)
			g.Map[itemNode.Name()].AddParents(g.Map[slotNode.Name()])

			// print currently evaluating sequence of items
			{
				items := make([]string, 0, usedItems.Len())
				for e := usedItems.Front(); e != nil; e = e.Next() {
					items = append(items, e.Value.(graph.Node).Name())
				}
				log.Print("trying " + strings.Join(items, " -> "))
			}

			// check whether this item should be skipped
			skip := false
			// only check one jewel per loop, since they're functionally
			// identical
			if strings.HasSuffix(itemNode.Name(), " jewel") {
				if !jewelChecked {
					jewelChecked = true
				} else {
					skip = true
				}
			}
			// the star ore code is unique in that it doesn't set the sub ID at
			// all, leaving it zeroed. so if we're looking at the star ore
			// slot, then skip any items that have a nonzero sub ID.
			if slotNode.Name() == "star ore spot" &&
				rom.Treasures[itemNode.Name()].SubID() != 0 {
				skip = true
			}

			// recurse
			if !skip && tryExploreTargets(
				g, reached, []graph.Node{itemNode}, goal, forbid, maxlen-1,
				itemList, usedItems, slotList, usedSlots) {
				return true
			}

			// item didn't work; unslot it and pop it onto the front of the
			// unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemNode)
			g.Map[itemNode.Name()].ClearParents()
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotNode)

		// reachable slots usually equivalent in terms of routing, so don't
		// bother checking more at this point
		break
	}

	// nothing worked
	log.Print("-- false; no slot/item combination worked")
	return false
}

// return the number of "step" nodes in the given set
func countSteps(nodes map[graph.Node]bool) (count int) {
	for node := range nodes {
		if node.IsStep() {
			count++
		}
	}
	return
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
func tryReachTargets(g *graph.Graph, goal, forbid []string, maxlen int,
	itemList, slotList, usedItems, usedSlots *list.List) bool {
	// prevent any known softlocks
	g.ClearMarks() // not strictly necessary
	if canSoftlock(g) {
		return false
	}
	g.ClearMarks()
	// make sure no forbidden nodes are reachable
	for _, node := range forbid {
		if canReachTargets(g, node) {
			return false
		}
	}
	// try to reach all targets
	if canReachTargets(g, goal...) {
		return true
	}
	// can't slot any more items
	if maxlen == 0 {
		return false
	}

	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slot := slotList.Back()
		slotList.MoveToFront(slot)

		slotName := slot.Value.(string)
		if !canReachTargets(g, slotName) {
			continue
		}

		// move slot from unused to used
		usedSlots.PushBack(slotName)
		slotList.Remove(slot)

		// try placing each unused item into the slot
		jewelChecked := false
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemName := itemList.Remove(itemList.Back()).(string)
			usedItems.PushBack(itemName)
			g.Map[itemName].AddParents(g.Map[slotName])

			{
				items := make([]string, 0, usedItems.Len())
				for e := usedItems.Front(); e != nil; e = e.Next() {
					items = append(items, e.Value.(string))
				}
				log.Print("trying " + strings.Join(items, " -> "))
			}

			// check whether this item should be skipped
			skip := false
			// only check one jewel per loop, since they're functionally identical
			if strings.HasSuffix(itemName, " jewel") {
				if !jewelChecked {
					jewelChecked = true
				} else {
					skip = true
				}
			}
			// the star ore code is unique in that it doesn't set the sub ID at
			// all, leaving it zeroed. so if we're looking at the star ore
			// slot, then skip any items that have a nonzero sub ID.
			if slotName == "star ore spot" && rom.Treasures[itemName].SubID() != 0 {
				skip = true
			}
			// don't place a L-1 item if you can already get the L-2 one
			if strings.HasSuffix(itemName, " L-1") {
				if canReachTargets(g, strings.Replace(itemName, "L-1", "L-2", 1)) {
					skip = true
				}
			} else if itemName == "find fist ring" {
				if canReachTargets(g, "find expert's ring") {
					skip = true
				}
			}

			if !skip && tryReachTargets(g, goal, forbid, maxlen-1,
				itemList, slotList, usedItems, usedSlots) {
				return true
			}

			// item didn't work; unslot it and pop it onto the front of the unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemName)
			g.Map[itemName].ClearParents()
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotName)

		// reachable slots usually equivalent in terms of routing, so don't
		// bother checking more at this point
		break
	}

	// nothing worked
	return false
}

// check if the targets are reachable using the current graph state
func canReachTargets(g *graph.Graph, targets ...string) bool {
	for _, target := range targets {
		if g.Map[target].GetMark(nil) != graph.MarkTrue {
			return false
		}
	}
	return true
}
