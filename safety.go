package main

import (
	"errors"

	"github.com/jangler/oos-randomizer/graph"
)

// TODO: write similar (?) functions to make sure dungeon navigation is
//       possible w/ the locations of small keys and given item state

// these are ordered, roughly, from least to most costly to check
var softlockChecks = [](func(graph.Graph) error){
	canFlowerSoftlock,
	canFeatherSoftlock,
	canEmberSeedSoftlock,
	canPiratesBellSoftlock,
}

// check for known softlock conditions
func canSoftlock(g graph.Graph) error {
	for _, check := range softlockChecks {
		if err := check(g); err != nil {
			return err
		}
	}
	return nil
}

// make sure you can't reach the spring banana cucco before having an item that
// removes flowers. you can still softlock if you forget to change season to
// spring, of course
func canFlowerSoftlock(g graph.Graph) error {
	// first check if cucco has been reached
	cucco := g["spring banana cucco"]
	if cucco.Mark != graph.MarkTrue {
		return nil
	}

	// temporarily make ammoless flower-removal items unavailable
	disabledNodes := append(g["remove flower sustainable"].Parents)
	disabledParents := make([][]*graph.Node, len(disabledNodes))
	for i, node := range disabledNodes {
		disabledParents[i] = node.Parents
		node.ClearParents()
	}
	defer g.ExploreFromStart()
	defer func() {
		for i, node := range disabledNodes {
			node.AddParents(disabledParents[i]...)
		}
	}()

	// see if you can still reach the exit
	g.ClearMarks()
	if cucco.GetMark(cucco, nil) == graph.MarkTrue {
		return errors.New("cucco softlock")
	}
	return nil
}

// make sure you can't reach the hide & seek area in subrosia without getting a
// shovel first. if your feather is stolen and you can't dig it back up, you
// can't exit that area.
func canFeatherSoftlock(g graph.Graph) error {
	// first check whether hide and seek has been reached
	hideAndSeek := g["hide and seek"]
	if hideAndSeek.Mark != graph.MarkTrue {
		return nil
	}
	// also test that you can jump, since you can't H&S without jumping (and it
	// would be beneficial even if you could)
	if g["jump"].Mark != graph.MarkTrue {
		return nil
	}

	shovel := g["shovel"]
	parents := shovel.Parents

	// check whether hide and seek is reachable if shovel is unreachable
	shovel.ClearParents()
	defer g.ExploreFromStart()
	defer shovel.AddParents(parents...)
	g.ClearMarks()
	if hideAndSeek.GetMark(hideAndSeek, nil) == graph.MarkTrue {
		return errors.New("feather softlock")
	}

	return nil
}

// since ember seeds can burn down bushes, make sure that the player doesn't
// have access to ember seeds without having a sustainable means of removing
// bushes first. flowers are covered by canFlowerSoftlock.
func canEmberSeedSoftlock(g graph.Graph) error {
	// first check if a seed item has been obtained (either gives ember seeds
	// currently, so just having one would be enough to softlock)
	seedItem := g["seed item"]
	if seedItem.Mark != graph.MarkTrue {
		return nil
	}

	// temporarily make ammoless bush-removal items unavailable
	disabledNodes := append(g["remove bush sustainable"].Parents)
	disabledParents := make([][]*graph.Node, len(disabledNodes))
	for i, node := range disabledNodes {
		disabledParents[i] = node.Parents
		node.ClearParents()
	}
	defer g.ExploreFromStart()
	defer func() {
		for i, node := range disabledNodes {
			node.AddParents(disabledParents[i]...)
		}
	}()

	// see if you can still reach the exit
	g.ClearMarks()
	if seedItem.GetMark(seedItem, nil) == graph.MarkTrue {
		return errors.New("ember seed softlock")
	}
	return nil
}

// if the player gives the captain his bell before talking to the ghost pirate,
// the ship will leave samasa desert and the player won't be able to get the
// item there. make sure the player can't get the bell before talking to the
// captain.
//
// "pirate house" as a prereq causes fewer problems for the algorithm than
// "samasa desert".
func canPiratesBellSoftlock(g graph.Graph) error {
	if canReachWithoutPrereq(g, g["rusty bell"], g["pirate house"]) {
		return errors.New("pirate's bell softlock")
	}
	return nil
}

// returns true iff the player can reach the goal node without the prerequisite
// one, given the current state of the graph.
func canReachWithoutPrereq(g graph.Graph, goal, prereq *graph.Node) bool {
	// check whether the goal node has been reached
	if goal.Mark != graph.MarkTrue {
		return false
	}

	// check whether goal node is reachable if prereq is not
	parents := prereq.Parents
	prereq.ClearParents()
	defer g.ExploreFromStart()
	defer prereq.AddParents(parents...)
	g.ClearMarks()
	if goal.GetMark(goal, nil) == graph.MarkTrue {
		return true
	}
	return false
}
