package graph

import (
	"testing"
)

func TestNodeGetters(t *testing.T) {
	names := []string{"foo", "bar"}
	marks := []Mark{MarkTrue, MarkFalse}
	nodes := []Node{
		&AndNode{Name: names[0], Mark: marks[0]},
		&OrNode{Name: names[1], Mark: marks[1]},
	}

	for i, node := range nodes {
		if node.GetName() != names[i] {
			t.Errorf("want %s, got %s", node.GetName(), names[i])
		}
		if node.PeekMark() != marks[i] {
			t.Errorf("want %d, got %d", node.PeekMark(), marks[i])
		}
		// GetMark isn't actually a getter, so it's not tested here
	}
}

func TestNodeRelationships(t *testing.T) {
	and1 := NewAndNode("and1")
	or1 := NewOrNode("or1")
	if and1.HasParents() {
		t.Errorf("new node has parents")
		t.Logf("%+v", and1.Parents)
	}
	if and1.HasChildren() {
		t.Errorf("new node has children")
		t.Logf("%+v", and1.Children)
	}
	// TODO add tests for when there actually are relationships
}
