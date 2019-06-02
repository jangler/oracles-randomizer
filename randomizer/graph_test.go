package randomizer

import (
	"reflect"
	"testing"
)

func resetNodes(nodes ...*node) {
	for _, n := range nodes {
		n.reached, n.indegree = false, 0
	}
}

func TestNodeRelationships(t *testing.T) {
	n1 := newNode("n1", andNode)
	testExpect(t, len(n1.parents), 0)
	testExpect(t, len(n1.children), 0)

	n2 := newNode("n2", orNode)
	n1.addParent(n2)
	testExpect(t, len(n1.parents), 1)
	testExpect(t, len(n1.children), 0)
	testExpect(t, len(n2.parents), 0)
	testExpect(t, len(n2.children), 1)
}

func TestExplore(t *testing.T) {
	and := newNode("and", andNode)
	testExpect(t, and.indegree, 0)
	testExpect(t, and.reached, true)

	or := newNode("or", orNode)
	testExpect(t, or.indegree, 0)
	testExpect(t, or.reached, false)

	count := newNode("count", countNode)
	count.minCount = 2
	resetNodes(and, or, count)
	testExpect(t, count.indegree, 0)
	testExpect(t, count.reached, false)

	// there used to be more tests here
	// but then the graph code changed
	// oh well
}

// helper for more concise testing
func testExpect(t *testing.T, x, y interface{}) {
	t.Helper()
	if !reflect.DeepEqual(x, y) {
		t.Errorf("expected %v to equal %v", x, y)
	}
}
