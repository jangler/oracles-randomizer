package main

import (
	"errors"

	"github.com/jangler/oos-randomizer/graph"
)

var softlockChecks = [](func(graph.Graph) error){
	canEmberSeedSoftlock,
	canD7ExitSoftlock,
	canD2ExitSoftlock,
	canSpringSwampSoftlock,
	canD5ExitFlipperSoftlock,
	canD5ExitBraceletSoftlock,
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

// since ember seeds can burn down bushes, make sure that the player doesn't
// have access to ember seeds without having a sustainable means of removing
// bushes first. flowers are covered by canFlowerSoftlock.
func canEmberSeedSoftlock(g graph.Graph) error {
	emberSeeds := g["ember seeds"]
	if emberSeeds.GetMark(emberSeeds, nil) != graph.MarkTrue {
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
	if emberSeeds.GetMark(emberSeeds, nil) == graph.MarkTrue {
		return errors.New("ember seed softlock")
	}
	return nil
}

// the player needs shovel before they can enter D7, or else then can be stuck
// if the default season is winter when they exit.
func canD7ExitSoftlock(g graph.Graph) error {
	// no snow piles == no softlock
	winter := g["western coast default winter"]
	if winter.GetMark(winter, nil) == graph.MarkFalse {
		return nil
	}

	if canReachWithoutPrereq(g, g["enter d7"], g["shovel"]) {
		return errors.New("d7 exit softlock")
	}
	return nil
}

// same deal with d2, except that feather also works. technically it's not just
// d2; the player can even enter the d2 entrance screen without bracelet and
// they'll still get the default season when they walk back out.
func canD2ExitSoftlock(g graph.Graph) error {
	// no snow piles == no softlock
	winter := g["eastern suburbs default winter"]
	if winter.GetMark(winter, nil) == graph.MarkFalse {
		return nil
	}

	if canReachWithoutPrereq(g, g["central woods of winter"], g["shovel"]) &&
		(canReachWithoutPrereq(g, g["central woods of winter"], g["jump"]) ||
			canReachWithoutPrereq(g, g["central woods of winter"], g["bracelet"])) {
		return errors.New("d2 exit softlock")
	}
	return nil
}

// it's impossible to swim out of the currents in south spool swamp in the
// spring, so the area in spring needs to be blocked by either the floodgate
// key or the swamp portal.
func canSpringSwampSoftlock(g graph.Graph) error {
	spring := g["spool swamp default spring"]
	if spring.GetMark(spring, nil) == graph.MarkFalse {
		return nil
	}

	if canReachWithoutPrereq(g, g["south swamp"], g["floodgate key"]) &&
		canReachWithoutPrereq(g, g["south swamp"], g["swamp portal"]) {
		return errors.New("spring swamp softlock")
	}
	return nil
}

// exiting d5 without flippers if it's not default winter means you're stuck.
func canD5ExitFlipperSoftlock(g graph.Graph) error {
	winter := g["north horon default winter"]
	if winter.GetMark(winter, nil) == graph.MarkTrue {
		return nil
	}

	if canReachWithoutPrereq(g, g["enter d5"], g["flippers"]) {
		return errors.New("d5 exit flipper softlock")
	}
	return nil
}

// exiting d5 without bracelet if it's not default autumn means you're stuck.
func canD5ExitBraceletSoftlock(g graph.Graph) error {
	autumn := g["north horon default autumn"]
	if autumn.GetMark(autumn, nil) == graph.MarkTrue {
		return nil
	}

	if canReachWithoutPrereq(g, g["enter d5"], g["bracelet"]) {
		return errors.New("d5 exit bracelet softlock")
	}
	return nil
}

// returns true iff the player can reach the goal node without the prerequisite
// one, given the current state of the graph.
func canReachWithoutPrereq(g graph.Graph, goal, prereq *graph.Node) bool {
	// check whether the goal node has been reached
	if goal.GetMark(goal, nil) != graph.MarkTrue {
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
