package main

import (
	"github.com/jangler/oos-randomizer/graph"
)

// check for known softlock conditions
func canSoftlock(g *graph.Graph) bool {
	return canShovelSoftlock(g) || canRosaPortalSoftlock(g)
}

// make sure you can't reach the shovel gift without either having a shovel
// already or getting a shovel there
func canShovelSoftlock(g *graph.Graph) bool {
	gift, shovel := g.Map["shovel gift"], g.Map["shovel"]
	parents := shovel.Parents()

	// check whether gift *is* shovel
	if !graph.IsNodeInSlice(gift, shovel.Parents()) {
		// check whether gift is reachable if shovel is unreachable
		shovel.ClearParents()
		defer shovel.AddParents(parents...)
		return canReachTargets(g, gift.Name())
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
	return canReachTargets(g, exit.Name())
}
