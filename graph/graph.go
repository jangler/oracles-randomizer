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
func (r *Graph) AddRootNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &RootNode{Name: name}
	}
}

// panics if node with name already exists
func (r *Graph) AddAndNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &AndNode{Name: name, Parents: make([]Node, 0)}
	}
}

// panics if node with name already exists
func (r *Graph) AddOrNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &OrNode{Name: name, Parents: make([]Node, 0)}
	}
}

func (r *Graph) CheckDuplicateName(name string) {
	if r.Map[name] != nil {
		panic("node named " + name + " already in route map")
	}
}

// panics if any of the given child nodes aren't actually child nodes
func (r *Graph) AddParents(links map[string][]string) {
	for childName, parentNames := range links {
		if child, ok := r.Map[childName]; ok {
			for _, parentName := range parentNames {
				if parent, ok := r.Map[parentName]; ok {
					child.(ChildNode).AddParents(parent)
				} else {
					panic("no parent named " + parentName)
				}
			}
		} else {
			panic("no child named " + childName)
		}
	}
}
