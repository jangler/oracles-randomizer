package randomizer

import (
	"regexp"
	"strings"
	"testing"

	"github.com/jangler/oracles-randomizer/rom"
)

var dungeonEntranceRegexp = regexp.MustCompile(`d[1-8].* entrance`)
var portalEntranceRegexp = regexp.MustCompile(`enter .+ portal`)

// returns true iff p1 is a parent of p2.
func isParent(p1Name string, p2 *prenode) bool {
	for _, parent := range p2.parents {
		if parent == p1Name {
			return true
		}
	}
	return false
}

func TestLinks(t *testing.T) {
	// TODO: needs to be changed manually for now
	game := rom.GameSeasons

	nodes := getPrenodes(game)
	rom.Init(nil, game)

	for key, slot := range rom.ItemSlots {
		treasureName := rom.FindTreasureName(slot.Treasure)
		if node, ok := nodes[treasureName]; ok {
			node.parents = append(node.parents, key)
		} else {
			n := &prenode{nType: andNode, parents: make([]interface{}, 1)}
			n.parents[0] = key
			nodes[treasureName] = n
		}
	}

	// check if any referenced nodes don't exist, and check if any nodes aren't
	// referenced.
	referenced := make(map[string]bool)
	for name, node := range nodes {
		for _, parentName := range node.parents {
			if _, ok := nodes[parentName.(string)]; !ok {
				t.Errorf("node %s references nonexistent node %s",
					name, parentName)
			}
			referenced[parentName.(string)] = true
		}
	}
	for name := range nodes {
		switch name {
		case "done", "unknown", "gasha seed", "piece of heart",
			"rare peach stone", "treasure map", "heart container":
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
			strings.HasSuffix(name, " compass") ||
			strings.HasSuffix(name, " dungeon map") ||
			strings.Contains(name, " ring L-") ||
			strings.Contains(name, " default ") ||
			strings.HasSuffix(name, " owl") {
			continue
		}
		if dungeonEntranceRegexp.MatchString(name) ||
			portalEntranceRegexp.MatchString(name) {
			continue
		}

		if _, ok := referenced[name]; !ok {
			t.Errorf("node %s is unreferenced", name)
		}
	}
}
