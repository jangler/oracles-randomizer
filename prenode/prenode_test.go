package prenode

import (
	"testing"
)

// returns true iff p1 is a parent of p2.
func isParent(p1Name string, p2 *Prenode) bool {
	for _, parent := range p2.Parents {
		if parent == p1Name {
			return true
		}
	}
	return false
}

func TestLinks(t *testing.T) {
	prenodes := GetAll()

	for name, p := range prenodes {
		// check if any non-root nodes are missing parents
		if p.Type != RootType && len(p.Parents) == 0 {
			t.Errorf("non-root prenode %s has no parents", name)
		}

		// check if any non-slot nodes are missing children
		switch p.Type {
		case AndType, OrType, AndSlotType, OrSlotType:
			break
		default:
			// ignore nodes which are present purely for -goal purposes
			if name == "done" || name == "enter d2" {
				break
			}

			hasChildren := false
			for _, p2 := range prenodes {
				if isParent(name, p2) {
					hasChildren = true
					break
				}
			}
			if !hasChildren {
				t.Errorf("non-slot prenode %s has no children", name)
			}
		}
	}
}
