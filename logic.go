package main

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// A prenode is the precursor to a graph node; its parents can be either
// strings (the names of other prenodes) or other prenodes. The main difference
// between a prenode and a graph node is that prenodes are trees, not graphs.
// String references to other prenodes become pointers when converting from
// prenodes to nodes, thus forming the graph.
type prenode struct {
	parents  []interface{}
	nType    nodeType
	minCount int
}

// rootPrenode returns a new prenode which does not have parents, and which
// will remain false until it does.
func rootPrenode(parents ...interface{}) *prenode {
	return &prenode{parents: parents, nType: orNode}
}

var seasonsPrenodes, agesPrenodes map[string]*prenode

func init() {
	seasonsPrenodes = make(map[string]*prenode)
	appendPrenodes(seasonsPrenodes, loadLogic("rings.yaml"),
		loadLogic("seasons_items.yaml"), loadLogic("seasons_kill.yaml"),
		loadLogic("holodrum.yaml"), loadLogic("subrosia.yaml"),
		loadLogic("portals.yaml"), loadLogic("seasons_dungeons.yaml"))
	flattenNestedPrenodes(seasonsPrenodes)

	agesPrenodes = make(map[string]*prenode)
	appendPrenodes(agesPrenodes, loadLogic("rings.yaml"),
		loadLogic("ages_items.yaml"), loadLogic("ages_kill.yaml"),
		loadLogic("labrynna.yaml"), loadLogic("ages_dungeons.yaml"))
	flattenNestedPrenodes(agesPrenodes)
}

// add nested nodes to the map and turn their references into strings.
func flattenNestedPrenodes(nodes map[string]*prenode) {
	done := true

	for name, pn := range nodes {
		subID := 0
		for i, parent := range pn.parents {
			switch parent := parent.(type) {
			case *prenode:
				subID++
				subName := fmt.Sprintf("%s %d", name, subID)
				pn.parents[i] = subName
				nodes[subName] = parent
				done = false
			}
		}
	}

	// recurse if nodes were added
	if !done {
		flattenNestedPrenodes(nodes)
	}
}

// returns a copy of all prenodes for the given game.
func getPrenodes(game int) map[string]*prenode {
	src := seasonsPrenodes
	if gameNames[game] == "ages" {
		src = agesPrenodes
	}

	dst := make(map[string]*prenode, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// merges the given prenode maps into the first argument.
func appendPrenodes(total map[string]*prenode, maps ...map[string]*prenode) {
	for _, nodeMap := range maps {
		for k, v := range nodeMap {
			if _, ok := total[k]; ok {
				panic("fatal: duplicate logic key: " + k)
			}
			total[k] = v
		}
	}
}

// loads a logic map from yaml.
func loadLogic(filename string) map[string]*prenode {
	raw := make(map[string]interface{})
	if err := yaml.Unmarshal(
		FSMustByte(false, "/logic/"+filename), raw); err != nil {
		panic(err)
	}

	m := make(map[string]*prenode)
	for k, v := range raw {
		m[k] = loadNode(v)
	}
	return m
}

// loads a node (and any of its explicit parents, recursively) from yaml.
func loadNode(v interface{}) *prenode {
	n := new(prenode)

	switch v := v.(type) {
	case []interface{}: // and node
		n.nType = andNode
		n.parents = make([]interface{}, len(v))
		for i, parent := range v {
			switch parent.(type) {
			case string:
				n.parents[i] = parent
			default:
				n.parents[i] = loadNode(parent)
			}
		}
	case map[interface{}]interface{}: // other node
		switch {
		case v["or"] != nil:
			n.nType = orNode
			n.parents = loadParents(v["or"])
		case v["not"] != nil:
			n.nType = nandNode
			n.parents = loadParents(v["not"])
		case v["nor"] != nil:
			n.nType = norNode
			n.parents = loadParents(v["nor"])
		case v["count"] != nil:
			n.nType = countNode
			n.minCount = v["count"].([]interface{})[0].(int)
			n.parents = make([]interface{}, 1)
			n.parents[0] = v["count"].([]interface{})[1].(string)
		default:
			println("unknown map type")
		}
	}

	return n
}

// loads a node's parents from yaml.
func loadParents(v interface{}) []interface{} {
	var parents []interface{}

	switch v := v.(type) {
	case []interface{}: // and node
		parents = make([]interface{}, len(v))
		for i, parent := range v {
			switch parent.(type) {
			case string:
				parents[i] = parent
			default:
				parents[i] = loadNode(parent)
			}
		}
	default: // single parent, other node
		parents = make([]interface{}, 1)
		parents[0] = loadNode(v)
	}

	return parents
}

var locationValues = map[string]int{
	"shop, 20 rupees":  -20,
	"shop, 30 rupees":  -30,
	"shop, 150 rupees": -150,
	"member's shop 1":  -300,
	"member's shop 2":  -300,
	"member's shop 3":  -200,

	"blaino prize": -10,

	"syrup":                    -300,
	"wild tokay game":          -10,
	"goron dance, with letter": -20,
	"goron dance present":      -10,
	"target carts 1":           -10,
	"target carts 2":           -10,
	"goron shooting gallery":   -20,

	"goron mountain old man":      300,
	"western coast old man":       300,
	"holodrum plain east old man": 200,
	"horon village old man":       100,
	"north horon old man":         100,

	"tarm ruins old man":          -200,
	"woods of winter old man":     -50,
	"holodrum plain west old man": -100,

	// rng is involved; each rupee is either worth 1, 5, or 10
	"d2 rupee room": 200,
	"d6 rupee room": 150,
}

var itemValues = map[string]int{
	"rupees, 1":   1,
	"rupees, 5":   5,
	"rupees, 10":  10,
	"rupees, 20":  20,
	"rupees, 30":  30,
	"rupees, 50":  50,
	"rupees, 100": 100,
	"rupees, 200": 200,
}
