package main

import (
	"github.com/jangler/oos-randomizer/graph"
	"log"
)

// TODO: write similar (?) functions to make sure dungeon navigation is
//       possible w/ the locations of small keys and given item state

// check for known softlock conditions
func canSoftlock(g *graph.Graph) bool {
	return canShovelSoftlock(g) || canRosaPortalSoftlock(g) ||
		canFlowerSoftlock(g)
}

// make sure you can't reach the shovel gift without either having a shovel
// already or getting a shovel there, *if* the shovel gift has been assigned
// yet.
func canShovelSoftlock(g *graph.Graph) bool {
	gift, shovel := g.Map["shovel gift"], g.Map["shovel"]
	parents := shovel.Parents()

	// if the slot hasn't been assigned yet or it *is* the shovel, it's fine
	if len(gift.Children()) > 0 &&
		!graph.IsNodeInSlice(gift, shovel.Parents()) {
		// check whether gift is reachable if shovel is unreachable
		shovel.ClearParents()
		defer shovel.AddParents(parents...)
		g.ClearMarks()
		if canReachTargets(g, gift.Name()) {
			log.Print("shovel softlock")
			return true
		}
	}

	return false
}

// make sure you can't exit subrosia via the rosa portal without either having
// activated it or getting an item that removes a bush you're stuck in
func canRosaPortalSoftlock(g *graph.Graph) bool {
	// temporarily make entrance and bush items unavailable
	entrance, exit := g.Map["rosa portal in wrapper"], g.Map["rosa portal out"]
	disabledNodes := append(g.Map["remove stuck bush"].Parents(), entrance)
	disabledParents := make([][]graph.Node, len(disabledNodes))
	for i, node := range disabledNodes {
		disabledParents[i] = node.Parents()
		node.ClearParents()
	}
	defer func() {
		for i, node := range disabledNodes {
			node.AddParents(disabledParents[i]...)
		}
	}()

	// see if you can still reach the exit
	g.ClearMarks()
	if canReachTargets(g, exit.Name()) {
		log.Print("portal softlock")
		return true
	}
	return false
}

// make sure you can't reach the spring banana cucco before having an item that
// removes flowers. you can still softlock if you forget to change season to
// spring, of course
func canFlowerSoftlock(g *graph.Graph) bool {
	// temporarily make entrance and bush items unavailable
	disabledNodes := append(g.Map["remove flower sustainable"].Parents())
	disabledParents := make([][]graph.Node, len(disabledNodes))
	for i, node := range disabledNodes {
		disabledParents[i] = node.Parents()
		node.ClearParents()
	}
	defer func() {
		for i, node := range disabledNodes {
			node.AddParents(disabledParents[i]...)
		}
	}()

	// see if you can still reach the exit
	g.ClearMarks()
	if canReachTargets(g, "spring banana cucco") {
		log.Print("cucco softlock")
		return true
	}
	return false
}
