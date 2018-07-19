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
	// TODO: this is only up to facade for now

	r.AddRootNodes(
		"enter d0",
		"ember tree",
	)

	r.AddAndNodes(
		"gnarled key",
		"enter d1",
		"d1 key 1",
		"d1 key 2",
		"enter goriya bros",
		"kill goriya bros",
		"satchel",
		"d1 boss key",
		"d1 essence",
		"portal 1",
		"winter",
		"mystery tree",
		"enter d2 1",
		"enter d2 2",
		"d2 mystery seeds 1",
		"d2 mystery seeds 2",
		"d2 key 1",
		"bracelet",
		"d2 key 2", // TODO: what are the d2 keys for?
		"d2 key 3 1",
		"d2 key 3 2",
		"enter facade",
	)

	// also include single-parent nodes
	r.AddOrNodes(
		"d0 key 1",
		"sword",
		"rupees",
		"bombs",
		"shield",
		"pop bubble",
		"remove bush",
		"kill stalfos",
		"hit lever",
		"fight goriya bros",
		"find ember seeds",
		"harvest seeds",
		"harvest ember seeds",
		"ember seeds",
		"kill goriya (pit)",
		"kill aquamentus",
		"boomerang",
		"rod",
		"hit switch (far)",
		"shovel",
		"find mystery seeds",
		"harvest mystery seeds",
		"mystery seeds",
		"kill rope",
		"kill hardhat beetle (pit)",
		"kill moblin (gap)",
		"kill gel",
		"d2 key 3",
		"kill facade",
	)

	// AND nodes only
	r.AddParents(map[string][]string{
		"gnarled key":         []string{"sword", "pop bubble"},
		"enter d1":            []string{"gnarled key", "remove bush"},
		"d1 key 1":            []string{"enter d1", "kill stalfos"},
		"d1 key 2":            []string{"d1 key 1", "kill stalfos", "hit lever"},
		"enter goriya bros":   []string{"d1 key 2", "bombs"},
		"kill goriya bros":    []string{"enter goriya bros", "fight goriya bros"},
		"satchel":             []string{"kill goriya bros"},
		"harvest ember seeds": []string{"satchel", "ember tree", "harvest seeds"},
		"d1 boss key":         []string{"ember seeds", "kill goriya (pit)"},
		"d1 essence":          []string{"d1 boss key", "kill aquamentus"},
		"portal 1":            []string{"d1 essence", "ember seeds", "remove bush"},
		"hit switch (far)":    []string{"boomerang", "bombs"},

		// TODO: account for sequence breaking
		"mystery tree": []string{"winter", "shovel"},
		"enter d2 1":   []string{"winter", "shovel", "remove bush"},
		"enter d2 2":   []string{"winter", "shovel", "bracelet"},

		"harvest mystery seeds": []string{"satchel", "mystery tree", "harvest seeds"},

		// same seeds, different route
		"d2 mystery seeds 1": []string{"enter d2 1", "kill rope", "ember seeds", "kill gel", "remove bush"},
		"d2 mystery seeds 2": []string{"enter d2 2", "bracelet", "remove bush"},

		"d2 key 1":     []string{"enter d2 1", "kill rope"},
		"bracelet":     []string{"d2 key 1", "ember seeds", "kill hardhat beetle (pit)", "kill moblin (gap)"},
		"d2 key 2":     []string{"enter d2 2", "remove bush", "bombs"},
		"d2 key 3 1":   []string{"enter d2 1", "kill rope", "ember seeds", "kill gel"},
		"d2 key 3 2":   []string{"enter d2 2", "bracelet"},
		"enter facade": []string{"d2 key 3", "bombs", "bracelet"},
		"kill facade":  []string{"enter facade", "bombs"},
	})

	// OR nodes only
	r.AddParents(map[string][]string{
		"d0 key 1":                  []string{"enter d0"},
		"sword":                     []string{"d0 key 1"},
		"rupees":                    []string{"sword", "ember seeds"},
		"bombs":                     []string{"rupees"},
		"shield":                    []string{"rupees"},
		"pop bubble":                []string{"sword", "bombs", "ember seeds"},
		"remove bush":               []string{"sword", "bombs", "ember seeds"},
		"kill stalfos":              []string{"sword", "bombs", "ember seeds", "rod"},
		"hit lever":                 []string{"sword", "ember seeds"},
		"fight goriya bros":         []string{"sword", "bombs"},
		"harvest seeds":             []string{"sword", "rod"},
		"find ember seeds":          []string{"enter d1"}, // TODO: among other places
		"ember seeds":               []string{"satchel", "find ember seeds", "harvest ember seeds"},
		"kill goriya (pit)":         []string{"sword", "bombs", "ember seeds"},
		"kill aquamentus":           []string{"sword", "bombs"},
		"boomerang":                 []string{"portal 1"},
		"rod":                       []string{"portal 1"},
		"winter":                    []string{"rod", "hit switch (far)"},
		"shovel":                    []string{"winter"},
		"find mystery seeds":        []string{"d2 mystery seeds 1", "d2 mystery seeds 2"},
		"mystery seeds":             []string{"satchel", "find mystery seeds", "harvest mystery seeds"},
		"kill rope":                 []string{"sword", "bombs", "ember seeds"},
		"kill hardhat beetle (pit)": []string{"sword", "shield", "bombs", "rod", "shovel", "bracelet"},
		"kill moblin (gap)":         []string{"sword", "bombs", "bracelet"},
		"kill gel":                  []string{"sword", "bombs", "ember seeds"},
		"d2 key 3":                  []string{"d2 key 3 1", "d2 key 3 2"},
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
