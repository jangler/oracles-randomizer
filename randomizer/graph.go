package randomizer

import (
	"fmt"
)

// map of names to set of nodes in a directed graph
type graph map[string]*node

// returns a new, initialized graph.
func newGraph() graph {
	return graph(make(map[string]*node))
}

// adds relationships in bulk between existing nodes in the graph, by name.
// attempting to link a name not in the graph results in a panic.
func (g graph) addParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := g[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := g[parentName]; ok {
					child.addParent(parent)
				} else {
					panic("no node named " + parentName)
				}
			}
		} else {
			panic("no child named " + childName)
		}
	}
}

// resets all the nodes in a graph to an unreached state (unless it's true
// without any reached parents). this is required any time relationships in the
// graph change.
func (g graph) reset() {
	for _, n := range g {
		n.indegree = 0
		n.reached = n.indegree >= n.mindegree()
	}
}

// determines the number of parents required for a node to be considered
// reached; see node.bounds().
type nodeType uint8

const (
	andNode nodeType = iota
	orNode
	countNode
	rupeesNode
)

// a single vertex in the graph.
type node struct {
	name     string
	ntype    nodeType
	reached  bool
	indegree int // number of reached parents
	minCount int // for countNodes, minimum parents to reach
	parents  []*node
	children []*node
	player   int
}

// returns a new unconnected graph node, not yet part of any graph.
func newNode(name string, nt nodeType) *node {
	// create node
	n := &node{
		name:     name,
		ntype:    nt,
		parents:  make([]*node, 0),
		children: make([]*node, 0),
	}
	n.reached = n.indegree >= n.mindegree()
	return n
}

// returns the minimum and reached parents at which a node is considered true.
func (n *node) mindegree() int {
	switch n.ntype {
	case andNode:
		return len(n.parents)
	case orNode, rupeesNode:
		return 1
	case countNode:
		return n.minCount
	default:
		panic("unknown type for node: " + n.name)
	}
}

// makes a node a parent of another node. a node can be a parent of another
// node multiple times; i.e. it can appear twice or more in the child's list of
// parents.
func (n *node) addParent(parent *node) {
	n.parents = append(n.parents, parent)
	parent.children = append(parent.children, n)
}

// removes the given node from this node's parents, once. it panics if the
// given node isn't actually a parent of this node.
func (n *node) removeParent(parent *node) {
	for i, p := range n.parents {
		if p == parent {
			for i, c := range p.children {
				if c == n {
					p.children = append(p.children[:i], p.children[i+1:]...)
					break
				}
			}
			n.parents = append(n.parents[:i], n.parents[i+1:]...)
			return
		}
	}
	panic(fmt.Sprintf("removeParent: %v is not a parent of %v", parent, n))
}

// removes all parent connections from the node.
func (n *node) clearParents() {
	for _, p := range n.parents {
		for i, c := range p.children {
			if c == n {
				p.children = append(p.children[:i], p.children[i+1:]...)
				break
			}
		}
	}
	n.parents = n.parents[:0]
}

// satisfies the fmt.Stringer interface.
func (n *node) String() string { return n.name }

// explores the graph starting from the given node, assuming the given node is
// reachable, marking successors as appropriate.
func (n *node) explore() {
	n.reached = true

	// rupees node sets indegree of children to total # of rupees reached
	if n.ntype == rupeesNode {
		for _, c := range n.children {
			c.indegree = n.indegree
			if !c.reached {
				c.exploreIfReachable()
			}
		}
		return
	}

	for _, c := range n.children {
		// add rupee value of parent to rupees node
		if c.ntype == rupeesNode {
			c.indegree += rupeeValues[n.name]
			c.explore()
			return
		} else {
			c.indegree++
		}

		if !c.reached {
			c.exploreIfReachable()
		} else {
			// handle count nodes
			for _, cc := range c.children {
				if !cc.reached && cc.ntype == countNode {
					cc.indegree++
					cc.exploreIfReachable()
				} else if cc.ntype == rupeesNode {
					cc.indegree += rupeeValues[c.name]
					cc.explore()
				}
			}
		}
	}
}

// like explore, but only works if the node's indegree is in bounds.
func (n *node) exploreIfReachable() {
	if n.indegree >= n.mindegree() {
		n.explore()
	}
}
