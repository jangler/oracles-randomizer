package main

import (
	"errors"

	"github.com/jangler/oos-randomizer/graph"
)

// TODO: write similar (?) functions to make sure dungeon navigation is
//       possible w/ the locations of small keys and given item state

var softlockChecks = [](func(graph.Graph) error){
	canShovelSoftlock,
	canFlowerSoftlock,
	canFeatherSoftlock,
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

// make sure you can't reach the shovel gift without either having a shovel
// already or getting a shovel there, *if* the shovel gift has been assigned
// yet.
func canShovelSoftlock(g graph.Graph) error {
	gift, shovel := g["shovel gift"], g["shovel"]
	parents := shovel.Parents

	// if the slot hasn't been assigned yet or it *is* the shovel, it's fine
	if len(gift.Children) > 0 &&
		!graph.IsNodeInSlice(gift, shovel.Parents) {
		// check whether gift is reachable if shovel is unreachable
		shovel.ClearParents()
		defer shovel.AddParents(parents...)
		g.ClearMarks()
		if gift.GetMark(gift, nil) == graph.MarkTrue {
			return errors.New("shovel softlock")
		}
	}

	return nil
}

// make sure you can't reach the spring banana cucco before having an item that
// removes flowers. you can still softlock if you forget to change season to
// spring, of course
func canFlowerSoftlock(g graph.Graph) error {
	// temporarily make entrance and bush items unavailable
	disabledNodes := append(g["remove flower sustainable"].Parents)
	disabledParents := make([][]*graph.Node, len(disabledNodes))
	for i, node := range disabledNodes {
		disabledParents[i] = node.Parents
		node.ClearParents()
	}
	defer func() {
		for i, node := range disabledNodes {
			node.AddParents(disabledParents[i]...)
		}
	}()

	// see if you can still reach the exit
	g.ClearMarks()
	cucco := g["spring banana cucco"]
	if cucco.GetMark(cucco, nil) == graph.MarkTrue {
		return errors.New("cucco softlock")
	}
	return nil
}

// make sure you can't reach the hide & seek area in subrosia without getting a
// shovel first. if your feather is stolen and you can't dig it back up, you
// can't exit that area.
func canFeatherSoftlock(g graph.Graph) error {
	hideAndSeek, shovel := g["hide and seek"], g["shovel"]
	parents := shovel.Parents

	// check whether hide and seek is reachable if shovel is unreachable
	shovel.ClearParents()
	defer shovel.AddParents(parents...)
	g.ClearMarks()
	if hideAndSeek.GetMark(hideAndSeek, nil) == graph.MarkTrue {
		return errors.New("feather softlock")
	}

	return nil
}
