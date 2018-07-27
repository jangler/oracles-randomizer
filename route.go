package main

import (
	"container/list"
	"fmt"
	"log"
	"math/rand"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/prenode"
	"github.com/jangler/oos-randomizer/rom"
)

// A Route is a set of information needed for finding an item placement route.
type Route struct {
	Graph graph.Graph
	Slots map[string]*graph.Node
}

// NewRoute returns an initialized route with all prenodes, and those prenodes
// with the names in start functioning as givens (always satisfied).
func NewRoute(start []string) *Route {
	g := graph.New()
	totalPrenodes := prenode.GetAll()

	// make start nodes given
	for _, key := range start {
		totalPrenodes[key] = prenode.And()
	}

	addNodes(g, totalPrenodes)
	addNodeParents(g, totalPrenodes)

	openSlots := make(map[string]*graph.Node, 0)
	for name, pn := range totalPrenodes {
		switch pn.Type {
		case prenode.AndSlotType, prenode.OrSlotType:
			openSlots[name] = g[name]
		}
	}

	return &Route{Graph: g, Slots: openSlots}
}

// CheckGraph returns an error for each orphan and childless node in the graph,
// ignoring nodes which are *supposed* to be orphans or childless. If there are
// no errors, it returns nil.
func (r *Route) CheckGraph() []error {
	var errs []error

	for name, node := range r.Graph {
		// check for parents and children
		if len(node.Parents) == 0 {
			// root nodes are supposed to be parentless
			if node.Type == graph.RootType {
				// it's supposed to be orphan/childless; skip it
				continue
			}

			if errs == nil {
				errs = make([]error, 0)
			}
			errs = append(errs, fmt.Errorf("orphan node: %s", name))
		}
		if len(node.Children) == 0 {
			// item slots are supposed to be childless
			if r.Slots[name] != nil {
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

func addNodes(g graph.Graph, prenodes map[string]*prenode.Prenode) {
	for key, pt := range prenodes {
		switch pt.Type {
		case prenode.AndType, prenode.AndSlotType, prenode.AndStepType:
			isStep := pt.Type == prenode.AndSlotType ||
				pt.Type == prenode.AndStepType
			g.AddNodes(graph.NewNode(key, graph.AndType, isStep))
		case prenode.OrType, prenode.OrSlotType, prenode.OrStepType,
			prenode.RootType:
			isStep := pt.Type == prenode.OrSlotType ||
				pt.Type == prenode.OrStepType
			g.AddNodes(graph.NewNode(key, graph.OrType, isStep))
		default:
			panic("unknown prenode type for " + key)
		}
	}
}

func addNodeParents(g graph.Graph, prenodes map[string]*prenode.Prenode) {
	// ugly but w/e
	for k, p := range prenodes {
		g.AddParents(map[string][]string{k: p.Parents})
	}
}

// attempts to create a path to the given targets by placing different items in
// slots.
func makeRoute(r *Route, start, goal, forbid []string,
	maxlen int) (usedItems, usedSlots *list.List) {
	// make stacks out of the item names and slot names for backtracking
	itemList := list.New()
	slotList := list.New()
	{
		// shuffle names in slices
		items := make([]*graph.Node, 0, len(prenode.BaseItems()))
		slots := make([]*graph.Node, 0, len(r.Slots))
		for itemName := range prenode.BaseItems() {
			items = append(items, r.Graph[itemName])
		}
		for slotName := range r.Slots {
			slots = append(slots, r.Graph[slotName])
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
	log.Print(len(r.Graph))

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems = list.New()
	usedSlots = list.New()

	// convert name lists into node lists
	startNodes := make([]*graph.Node, len(start))
	for i, name := range start {
		startNodes[i] = r.Graph[name]
	}
	goalNodes := make([]*graph.Node, len(goal))
	for i, name := range goal {
		goalNodes[i] = r.Graph[name]
	}
	forbidNodes := make([]*graph.Node, len(forbid))
	for i, name := range forbid {
		forbidNodes[i] = r.Graph[name]
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
			log.Fatalf("FATAL: usedItems.Len() == %d; usedSlots.Len() == %d",
				usedItems.Len(), usedSlots.Len())
		}
		for i := 0; i < usedItems.Len(); i++ {
			log.Printf("%v <- %v",
				usedItems.Front().Value.(*graph.Node),
				usedSlots.Front().Value.(*graph.Node))
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
func tryExploreTargets(g graph.Graph, start map[*graph.Node]bool,
	add, goal, forbid []*graph.Node, maxlen int,
	itemList, usedItems, slotList, usedSlots *list.List) bool {
	// explore given the old state and changes
	reached := g.Explore(start, add, nil)
	log.Print(len(reached), " nodes reached")
	log.Print(countSteps(reached), " steps reached")

	// abort if any forbidden node is reached
	for _, node := range forbid {
		if reached[node] {
			log.Printf("-- false; reached forbidden node %s", node)
			return false
		}
	}

	// success if all goal nodes are reached *and* all slots are filled
	allReached := true
	for _, node := range goal {
		if !reached[node] {
			log.Printf("-- have not reached goal node %s", node)
			allReached = false
			break
		}
	}
	if allReached {
		log.Printf("-- all goals reached")
		if slotList.Len() == 0 {
			log.Printf("-- true; all goals reached and slots filled")
			return true
		}
	}

	// if the new state doesn't reach any more steps, abandon this branch,
	// *unless* the new item is a jewel, or we've already reached the goals.
	if !allReached && !strings.HasSuffix(add[0].Name, " jewel") {
		if countSteps(reached) <= countSteps(start) {
			log.Printf("-- false; reached steps %d <= start steps %d",
				countSteps(reached), countSteps(start))
			return false
		}
	}

	// can't slot any more items
	if !allReached && maxlen == 0 {
		log.Print("-- false; slotted maxlen items")
		return false
	}

	// check for softlocks
	if canSoftlock(g) {
		log.Print("-- false; route blocked by softlock")
		return false
	}

	// TODO: move some of this logic to its own function(s)
	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slotElem := slotList.Back()
		slotList.MoveToFront(slotElem)

		// see if slot node has been reached
		slotNode := slotElem.Value.(*graph.Node)
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
			itemNode := itemList.Remove(itemList.Back()).(*graph.Node)
			usedItems.PushBack(itemNode)
			g[itemNode.Name].AddParents(g[slotNode.Name])

			// print currently evaluating sequence of items
			{
				items := make([]string, 0, usedItems.Len())
				for e := usedItems.Front(); e != nil; e = e.Next() {
					items = append(items, e.Value.(*graph.Node).Name)
				}
				log.Print("trying " + strings.Join(items, " -> "))
			}

			// check whether this item should be skipped
			skip := false
			// only check one jewel per loop, since they're functionally
			// identical
			if strings.HasSuffix(itemNode.Name, " jewel") {
				if !jewelChecked {
					jewelChecked = true
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

			// recurse
			if !skip && tryExploreTargets(
				g, reached, []*graph.Node{itemNode}, goal, forbid, maxlen-1,
				itemList, usedItems, slotList, usedSlots) {
				return true
			}

			// item didn't work; unslot it and pop it onto the front of the
			// unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemNode)
			g[itemNode.Name].ClearParents()
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
func countSteps(nodes map[*graph.Node]bool) int {
	count := 0
	for node := range nodes {
		if node.IsStep {
			count++
		}
	}
	return count
}

// check if the targets are reachable using the current graph state
func canReachTargets(g graph.Graph, targets ...string) bool {
	for _, target := range targets {
		if g[target].GetMark(nil) != graph.MarkTrue {
			return false
		}
	}
	return true
}
