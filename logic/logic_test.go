package logic

import (
	"regexp"
	"strings"
	"testing"

	"github.com/jangler/oracles-randomizer/rom"
)

var dungeonEntranceRegexp = regexp.MustCompile(`d[1-8] entrance`)

// returns true iff p1 is a parent of p2.
func isParent(p1Name string, p2 *Node) bool {
	for _, parent := range p2.Parents {
		if parent == p1Name {
			return true
		}
	}
	return false
}

func TestLinks(t *testing.T) {
	// need to be changed manually for now
	nodes := GetSeasons()
	rom.Init(rom.GameSeasons)

	for key, slot := range rom.ItemSlots {
		treasureName := rom.FindTreasureName(slot.Treasure)
		if node, ok := nodes[treasureName]; ok {
			node.Parents = append(node.Parents, key)
		} else {
			nodes[treasureName] = And(key)
		}
	}

	// check if any referenced nodes don't exist, and check if any nodes aren't
	// referenced.
	referenced := make(map[string]bool)
	for name, node := range nodes {
		for _, parentName := range node.Parents {
			if _, ok := nodes[parentName.(string)]; !ok {
				t.Errorf("node %s references nonexistent node %s",
					name, parentName)
			}
			referenced[parentName.(string)] = true
		}
	}
	for name := range nodes {
		switch name {
		case "done", "gasha seed", "piece of heart", "rare peach stone",
			"treasure map", "dungeon map", "compass", "strange flute",
			"heart container":
			continue
		case "pegasus seeds", "any satchel":
			// defined for consistency but unused
			continue
		case "ricky nuun", "dimitri nuun", "moosh nuun":
			continue
		}
		if strings.Contains(name, "rupee") ||
			strings.HasSuffix(name, "old man") ||
			strings.HasSuffix(name, " ring") ||
			strings.Contains(name, " ring L-") ||
			strings.Contains(name, " default ") ||
			strings.HasSuffix(name, " owl") {
			continue
		}
		if dungeonEntranceRegexp.MatchString(name) {
			continue
		}

		if _, ok := referenced[name]; !ok {
			t.Errorf("node %s is unreferenced", name)
		}
	}
}
