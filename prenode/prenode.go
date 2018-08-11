package prenode

import (
	"fmt"
)

// This package contains definitions of nodes and node relationships before
// they are inserted into the graph. This is necessary because nodes
// relationships can't be made until the nodes are added first (and it's nice
// not to clutter the other packages with all these definitions).

// XXX need to be careful about rings. i can't imagine a situation where you'd
//     need both energy ring and fist ring, but if you did, then you'd need to
//     have the L-2 ring box to do so without danger of soft locking.

// A Type identifies whether a prenode is and And, Or, or Root node, whether it
// is an item slot, and whether it is a non-item slot milestone.
type Type int

// And, Or, and Root are pretty self-explanatory. One with a Slot suffix is an
// item slot, and one with a Step suffix is treated as a milestone for routing
// purposes. Slot types are also treated as steps; see the Point.IsStep()
// function.
//
// The following function are half syntactic sugar for declaring large lists of
// node relationships.
const (
	RootType Type = iota
	AndType
	OrType
	AndSlotType
	OrSlotType
	AndStepType
	OrStepType
)

// A Prenode is a mapping of strings that will become And or Or nodes in the
// graph. A prenode can have nested prenodes as parents instead of strings.
type Prenode struct {
	Parents []interface{}
	Type    Type
}

// CreateFunc returns a function that creates graph nodes from a list of key
// strings or sub-prenodes, based on the given prenode type.
func CreateFunc(prenodeType Type) func(parents ...interface{}) *Prenode {
	return func(parents ...interface{}) *Prenode {
		return &Prenode{Parents: parents, Type: prenodeType}
	}
}

// Convenience functions for creating prenodes succinctly. See the Type const
// comment for information on the various types.
var (
	Root    = CreateFunc(RootType)
	And     = CreateFunc(AndType)
	AndSlot = CreateFunc(AndSlotType)
	AndStep = CreateFunc(AndStepType)
	Or      = CreateFunc(OrType)
	OrSlot  = CreateFunc(OrSlotType)
	OrStep  = CreateFunc(OrStepType)
)

var allPrenodes map[string]*Prenode

func init() {
	allPrenodes = make(map[string]*Prenode)
	appendPrenodes(allPrenodes,
		itemPrenodes, baseItemPrenodes, ignoredBaseItemPrenodes, killPrenodes,
		holodrumPrenodes, subrosiaPrenodes, portalPrenodes, seasonPrenodes,
		d0Prenodes, d1Prenodes, d2Prenodes, d3Prenodes, d4Prenodes,
		d5Prenodes, d6Prenodes, d7Prenodes, d8Prenodes, d9Prenodes)
	flattenNestedPrenodes(allPrenodes)
}

// add nested prenodes to the map and turn their references into strings
func flattenNestedPrenodes(prenodes map[string]*Prenode) {
	done := true

	for name, pn := range prenodes {
		subID := 0
		for i, parent := range pn.Parents {
			switch parent := parent.(type) {
			case *Prenode:
				subID++
				subName := fmt.Sprintf("%s %d", name, subID)
				pn.Parents[i] = subName
				prenodes[subName] = parent
				done = false
			}
		}
	}

	// recurse if prenodes were added
	if !done {
		flattenNestedPrenodes(prenodes)
	}
}

// BaseItems returns a map of item prenodes that may be assigned to slots.
func BaseItems() map[string]*Prenode {
	return copyMap(baseItemPrenodes)
}

// GetAll returns a copy of all prenodes.
func GetAll() map[string]*Prenode {
	return copyMap(allPrenodes)
}

// merge the given maps into the first argument
func appendPrenodes(total map[string]*Prenode, maps ...map[string]*Prenode) {
	for _, prenodeMap := range maps {
		for k, v := range prenodeMap {
			if _, ok := total[k]; ok {
				panic("fatal: duplicate prenode key: " + k)
			}
			total[k] = v
		}
	}
}

// returns a shallow copy of a string/prenode map
func copyMap(src map[string]*Prenode) map[string]*Prenode {
	dst := make(map[string]*Prenode, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
