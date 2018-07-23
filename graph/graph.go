package graph

// this file contains facilities for linking nodes into a graph

type Graph struct {
	Map map[string]Node
}

func NewGraph() *Graph {
	return &Graph{
		Map: make(map[string]Node),
	}
}

// panics if node with name already exists
func (g *Graph) AddAndNodes(names ...string) {
	for _, name := range names {
		g.CheckDuplicateName(name)
		g.Map[name] = &AndNode{Name: name, Parents: make([]Node, 0), Children: make([]Node, 0)}
	}
}

// panics if node with name already exists
func (g *Graph) AddOrNodes(names ...string) {
	for _, name := range names {
		g.CheckDuplicateName(name)
		g.Map[name] = &OrNode{Name: name, Parents: make([]Node, 0), Children: make([]Node, 0)}
	}
}

func (g *Graph) CheckDuplicateName(name string) {
	if g.Map[name] != nil {
		panic("node named " + name + " already in route map")
	}
}

func (g *Graph) AddParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := g.Map[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := g.Map[parentName]; ok {
					child.AddParents(parent)
				} else {
					panic("no node named " + parentName)
				}
			}
		} else {
			panic("no child named " + childName)
		}
	}
}

func (g *Graph) ClearMarks() {
	for _, node := range g.Map {
		node.SetMark(MarkNone)
	}
}
