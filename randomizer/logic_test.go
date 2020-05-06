package randomizer

import (
	"regexp"
	"strings"
	"testing"
)

var dungeonEntranceRegexp = regexp.MustCompile(`d[1-8].* entrance`)
var portalEntranceRegexp = regexp.MustCompile(`enter .+ portal`)

func TestLinks(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		testLinksForGame(t, game)
	}
}

func testLinksForGame(t *testing.T, game int) {
	nodes := getPrenodes(game)
	rom := newRomState(nil, game, 0, nil)

	for key, slot := range rom.itemSlots {
		treasureName, _ := reverseLookup(rom.treasures, slot.treasure)
		if node, ok := nodes[treasureName.(string)]; ok {
			node.parents = append(node.parents, key)
		} else {
			n := &prenode{nType: andNode, parents: make([]interface{}, 1)}
			n.parents[0] = key
			nodes[treasureName.(string)] = n
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
		case "done", "gasha seed", "piece of heart", "rare peach stone",
			"treasure map", "heart container":
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
