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

// A Point is a mapping of point strings that will become And or Or nodes in
// the graph.
type Point interface {
	Parents() []string
}

// the different types of points are all just string slices; the reason for
// having different ones is purely for type assertions

type And []string

func (p And) Parents() []string { return p }

type Or []string

func (p Or) Parents() []string { return p }

type AndSlot []string

func (p AndSlot) Parents() []string { return p }

type OrSlot []string

func (p OrSlot) Parents() []string { return p }

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
		totalPoints[key] = And{}
	}

	addPointNodes(g, totalPoints)
	addPointParents(g, totalPoints)

	openSlots := make(map[string]Point, 0)
	for name, point := range totalPoints {
		switch point.(type) {
		case AndSlot, OrSlot:
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
			// base items are supposed to be parentless
			if baseItemPoints[name] != nil || ignoredBaseItemPoints[name] != nil {
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
			switch r.Points[name].(type) {
			case AndSlot, OrSlot:
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
		switch pt.(type) {
		case And, AndSlot:
			g.AddAndNodes(key)
		case Or, OrSlot:
			g.AddOrNodes(key)
		default:
			panic("unknown point type for " + key)
		}
	}
}

func addPointParents(g *graph.Graph, points map[string]Point) {
	// ugly but w/e
	for k, p := range points {
		g.AddParents(map[string][]string{k: p.Parents()})
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
func makeRoute(r *Route, goal, forbid []string,
	maxlen int) (usedItems, usedSlots, itemList, slotList *list.List) {
	// make stacks out of the item names and slot names for backtracking
	itemList = list.New()
	slotList = list.New()
	{
		// shuffle names in slices
		items := make([]string, 0, len(baseItemPoints))
		slots := make([]string, 0, len(r.Slots))
		for itemName, _ := range baseItemPoints {
			items = append(items, itemName)
		}
		for slotName, _ := range r.Slots {
			slots = append(slots, slotName)
		}
		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		rand.Shuffle(len(slots), func(i, j int) {
			slots[i], slots[j] = slots[j], slots[i]
		})

		// push the shuffled items onto the stacks
		for _, itemName := range items {
			itemList.PushBack(itemName)
		}
		for _, slotName := range slots {
			slotList.PushBack(slotName)
		}
	}

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems = list.New()
	usedSlots = list.New()

	if tryReachTargets(r.Graph, goal, forbid, maxlen,
		itemList, slotList, usedItems, usedSlots) {
		log.Print("-- success")
		for _, target := range goal {
			log.Print("-- path to " + target)
			r.Graph.ClearMarks()
			path := findPath(r.Graph, r.Graph.Map[target])
			for path.Len() > 0 {
				step := path.Remove(path.Front()).(string)
				log.Print(step)
			}
		}
		log.Print("-- slotted items")
		if usedItems.Len() != usedSlots.Len() {
			log.Fatalf("FATAL: usedItems.Len() == %d; usedSlots.Len() == %d", usedItems.Len(), usedSlots.Len())
		}
		for i := 0; i < usedItems.Len(); i++ {
			log.Printf("%s <- %s", usedItems.Front().Value.(string), usedSlots.Front().Value.(string))
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
