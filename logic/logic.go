package logic

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

// This package contains definitions of nodes and node relationships before
// they are inserted into the graph. This is necessary because nodes
// relationships can't be made until the nodes are added first (and it's nice
// not to clutter the other packages with all these definitions).

// A Type identifies whether a node is an And, Or, or Count node, whether it is
// an item slot, and whether it is a non-item slot milestone.
type Type int

// The following functions are half syntactic sugar for declaring large lists
// of node relationships.
const (
	AndType Type = iota
	OrType
	CountType
)

// A Node is a mapping of strings that will become And or Or nodes in the
// graph. A node can have nested nodes as parents instead of strings.
type Node struct {
	Parents  []interface{}
	Type     Type
	MinCount int
}

// Root returns a new root node - one which does not have parents, and remains
// false until it does.
func Root(parents ...interface{}) *Node {
	return &Node{Parents: parents, Type: OrType}
}

var seasonsNodes, agesNodes map[string]*Node

func init() {
	seasonsNodes = make(map[string]*Node)
	appendNodes(seasonsNodes, loadLogic("rings.yaml"),
		loadLogic("seasons_items.yaml"), loadLogic("seasons_kill.yaml"),
		loadLogic("holodrum.yaml"), loadLogic("subrosia.yaml"),
		loadLogic("portals.yaml"), loadLogic("seasons_dungeons.yaml"))
	flattenNestedNodes(seasonsNodes)

	agesNodes = make(map[string]*Node)
	appendNodes(agesNodes, loadLogic("rings.yaml"),
		loadLogic("ages_items.yaml"), loadLogic("ages_kill.yaml"),
		loadLogic("labrynna.yaml"), loadLogic("ages_dungeons.yaml"))
	flattenNestedNodes(agesNodes)
}

// add nested nodes to the map and turn their references into strings
func flattenNestedNodes(nodes map[string]*Node) {
	done := true

	for name, pn := range nodes {
		subID := 0
		for i, parent := range pn.Parents {
			switch parent := parent.(type) {
			case *Node:
				subID++
				subName := fmt.Sprintf("%s %d", name, subID)
				pn.Parents[i] = subName
				nodes[subName] = parent
				done = false
			}
		}
	}

	// recurse if nodes were added
	if !done {
		flattenNestedNodes(nodes)
	}
}

// GetSeasons returns a copy of all seasons nodes.
func GetSeasons() map[string]*Node {
	return copyMap(seasonsNodes)
}

// GetAges returns a copy of all ages nodes.
func GetAges() map[string]*Node {
	return copyMap(agesNodes)
}

// merge the given maps into the first argument
func appendNodes(total map[string]*Node, maps ...map[string]*Node) {
	for _, nodeMap := range maps {
		for k, v := range nodeMap {
			if _, ok := total[k]; ok {
				panic("fatal: duplicate logic key: " + k)
			}
			total[k] = v
		}
	}
}

// returns a shallow copy of a string/node map
func copyMap(src map[string]*Node) map[string]*Node {
	dst := make(map[string]*Node, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// loads a logic map from yaml.
func loadLogic(filename string) map[string]*Node {
	raw := make(map[string]interface{})
	if err := yaml.Unmarshal(
		FSMustByte(false, "/logic/"+filename), raw); err != nil {
		panic(err)
	}

	m := make(map[string]*Node)
	for k, v := range raw {
		m[k] = loadNode(v)
	}
	return m
}

// loads a node (and any of its explicit parents, recursively) from yaml.
func loadNode(v interface{}) *Node {
	n := new(Node)

	switch v := v.(type) {
	case string:
		n.Type = AndType
		n.Parents = make([]interface{}, 1)
		n.Parents = append(n.Parents, v)
	case []interface{}:
		n.Type = AndType
		n.Parents = make([]interface{}, len(v))
		for i, parent := range v {
			switch parent.(type) {
			case string:
				n.Parents[i] = parent
			default:
				n.Parents[i] = loadNode(parent)
			}
		}
	case map[interface{}]interface{}:
		if v["or"] != nil {
			n.Type = OrType
			n.Parents = loadParents(v["or"])
		} else if v["count"] != nil {
			n.Type = CountType
			n.MinCount = v["count"].([]interface{})[0].(int)
			n.Parents = make([]interface{}, 1)
			n.Parents[0] = v["count"].([]interface{})[1].(string)
		} else {
			println("unknown map type")
		}
	}

	return n
}

// loads a node's parents from yaml.
func loadParents(v interface{}) []interface{} {
	var parents []interface{}

	switch v := v.(type) {
	case string:
		parents = make([]interface{}, 1)
		parents[0] = v
	case []interface{}:
		parents = make([]interface{}, len(v))
		for i, parent := range v {
			switch parent.(type) {
			case string:
				parents[i] = parent
			default:
				parents[i] = loadNode(parent)
			}
		}
	default:
		parents = make([]interface{}, 1)
		parents[0] = loadNode(v)
	}

	return parents
}
