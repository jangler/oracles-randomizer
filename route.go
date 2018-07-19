package main

// this file contains the actual connection of graphs in the game node, and
// tracks them as they update.

type Route struct {
	Map map[string]Node
}

func NewRoute() *Route {
	r := &Route{
		Map: make(map[string]Node),
	}

	// for now, just do the regular order

	r.AddRootNodes(
		"enter d0",
	)

	r.AddAndNodes(
		"gnarled key",
		"enter d1",
		"d1 key 1",
		"d1 key 2",
		"satchel",
		"ember seeds",
		"d1 boss key",
		"d1 essence",
	)

	// also include single-parent nodes
	r.AddOrNodes(
		"d0 key 1",
		"sword",
		"rupees",
		"bombs",
		"pop bubble",
		"remove bushes",
		"kill stalfos",
		"hit lever",
		"kill goriya bros",
		"kill goriya (pit)",
		"kill aquamentus",
	)

	// AND nodes only
	r.AddParents(map[string][]string{
		"gnarled key": []string{"sword", "pop bubble"},
		"enter d1":    []string{"gnarled key", "remove bushes"},
		"d1 key 1":    []string{"enter d1", "kill stalfos"},
		"d1 key 2":    []string{"d1 key 1", "kill stalfos", "hit lever"},
		"satchel":     []string{"d1 key 2", "bombs", "kill goriya bros"},
		"ember seeds": []string{"satchel", "sword"},
		"d1 boss key": []string{"ember seeds", "kill goriya (pit)"},
		"d1 essence":  []string{"d1 boss key", "sword"},
	})

	// OR nodes only
	r.AddParents(map[string][]string{
		"d0 key 1":          []string{"enter d0"},
		"sword":             []string{"d0 key 1"},
		"rupees":            []string{"sword", "ember seeds"},
		"bombs":             []string{"rupees"},
		"pop bubble":        []string{"sword", "bombs", "ember seeds"},
		"remove bushes":     []string{"sword", "bombs", "ember seeds"},
		"kill stalfos":      []string{"sword", "bombs", "ember seeds"},
		"hit lever":         []string{"sword", "ember seeds"},
		"kill goriya bros":  []string{"sword", "bombs"},
		"kill goriya (pit)": []string{"sword", "bombs", "ember seeds"},
		"kill aquamentus":   []string{"sword", "bombs"},
	})

	// validate
	for name, node := range r.Map {
		switch nt := node.(type) {
		case ChildNode:
			if !nt.HasParents() {
				panic("node with no parents: " + name)
			}
		}
	}

	return r
}

// panics if node with name already exists
func (r *Route) AddRootNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &RootNode{Name: name}
	}
}

// panics if node with name already exists
func (r *Route) AddAndNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &AndNode{Name: name, Parents: make([]Node, 0)}
	}
}

// panics if node with name already exists
func (r *Route) AddOrNodes(names ...string) {
	for _, name := range names {
		r.CheckDuplicateName(name)
		r.Map[name] = &OrNode{Name: name, Parents: make([]Node, 0)}
	}
}

func (r *Route) CheckDuplicateName(name string) {
	if r.Map[name] != nil {
		panic("node named " + name + " already in route map")
	}
}

// panics if any of the given child nodes aren't actually child nodes
func (r *Route) AddParents(links map[string][]string) {
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
