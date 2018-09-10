package main

import (
	"errors"

	"github.com/jangler/oos-randomizer/graph"
)

// there's another softlock i can think of, but it's far-fetched enough that
// i'm not going to write a test for it. if spool swamp defaults to winter and
// you get bracelet, flippers, floodgate key, and spring/summer/autumn but NOT
// shovel or feather/cape, you can enter the spool swamp subrosia portal and be
// stuck there when you come out and the area defaults back to winter. you'd
// also need a way to get to the floodgate keyhole without ricky, which means
// either holodrum plain defaults to winter and you get summer, or it defaults
// to winter as well.

var softlockChecks = [](func(graph.Graph) error){
	canD7ExitSoftlock,
	canD2ExitSoftlock,
	canD5ExitBraceletSoftlock,
}

// check for known softlock conditions
func canSoftlock(g graph.Graph) error {
	g.ExploreFromStart()
	for _, check := range softlockChecks {
		if err := check(g); err != nil {
			return err
		}
	}
	return nil
}

// the player needs shovel before they can enter D7, or else they can be stuck
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
		canReachWithoutPrereq(g, g["central woods of winter"], g["flute"]) &&
		(canReachWithoutPrereq(g, g["central woods of winter"], g["jump"]) ||
			canReachWithoutPrereq(g, g["central woods of winter"], g["bracelet"])) {
		return errors.New("d2 exit softlock")
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
