package main

// this file contains the actual connection of graphs in the game node, and
// tracks them as they update.

// XXX XXX XXX BOOMERANG L-2 CAN BREAK MUSHROOMS XXX XXX XXX

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

		// later
		"slingshot",
		"scent seeds",
		"gale seeds",
		"feather",
		"cape",
		"fists",
	)

	r.Map["slingshot"].SetMark(MarkFalse)
	r.Map["scent seeds"].SetMark(MarkFalse)
	r.Map["gale seeds"].SetMark(MarkFalse)
	r.Map["feather"].SetMark(MarkFalse)
	r.Map["cape"].SetMark(MarkFalse)
	r.Map["fists"].SetMark(MarkFalse)

	r.AddAndNodes(
		"gnarled key",
		"enter d1",
		"ember seeds",
		"portal 1",
		"winter",
		"mystery tree",
		"enter d2 1",
		"enter d2 2",
		"d2 mystery seeds 1",
		"d2 mystery seeds 2",
		"mystery seeds",
		"d2 key 1",
		"bracelet",
		"d2 key 2", // TODO: what are the d2 keys for?
		"d2 key 3 1",
		"d2 key 3 2",
		"enter facade",
		"d2 warp",
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
		"hit lever",
		"seed item",
		"find ember seeds",
		"harvest seeds",
		"harvest ember seeds",
		"get ember seeds",
		"boomerang",
		"rod",
		"hit switch (far)",
		"shovel",
		"find mystery seeds",
		"harvest mystery seeds",
		"get mystery seeds",
		"d2 key 3",

		// later
		"jump",
	)

	for key, _ := range killNodesAnd {
		r.AddAndNodes(key)
	}
	for key, _ := range killNodesOr {
		r.AddOrNodes(key)
	}

	for key, _ := range d1NodesAnd {
		r.AddAndNodes(key)
	}
	for key, _ := range d1NodesOr {
		r.AddOrNodes(key)
	}

	// AND nodes only
	r.AddParents(map[string][]string{
		// TODO: non-any% stuff
		"gnarled key":         []string{"sword", "pop bubble"},
		"enter d1":            []string{"gnarled key", "remove bush"},
		"harvest ember seeds": []string{"satchel", "ember tree", "harvest seeds"},
		"ember seeds":         []string{"satchel", "get ember seeds"},
		"portal 1":            []string{"d1 essence", "ember seeds", "remove bush"},

		// TODO: account for sequence breaking
		"mystery tree": []string{"winter", "shovel"},
		"enter d2 1":   []string{"winter", "shovel", "remove bush"},
		"enter d2 2":   []string{"winter", "shovel", "bracelet"},

		"harvest mystery seeds": []string{"satchel", "mystery tree", "harvest seeds"},

		// same seeds, different route
		"d2 mystery seeds 1": []string{"enter d2 1", "kill rope", "ember seeds", "kill gel", "remove bush"},
		"d2 mystery seeds 2": []string{"enter d2 2", "bracelet", "remove bush"},

		"mystery seeds": []string{"satchel", "get mystery seeds"},

		"d2 key 1":     []string{"enter d2 1", "kill rope"},
		"bracelet":     []string{"d2 key 1", "ember seeds", "kill hardhat (pit, throw)", "kill moblin (gap, throw)"},
		"d2 key 2":     []string{"enter d2 2", "remove bush", "bombs"},
		"d2 key 3 1":   []string{"enter d2 1", "kill rope", "ember seeds", "kill gel"},
		"d2 key 3 2":   []string{"enter d2 2", "bracelet"},
		"enter facade": []string{"d2 key 3", "bombs", "bracelet"},
		"d2 warp":      []string{"enter facade", "kill facade"},
	})

	// OR nodes only
	r.AddParents(map[string][]string{
		"find bombs":         []string{}, // have to be careful not to soft lock
		"d0 key 1":           []string{"enter d0"},
		"sword":              []string{"d0 key 1"},
		"rupees":             []string{"sword", "ember seeds", "rod", "shovel", "bracelet", "bombs"}, // XXX could bombs soft lock?
		"bombs":              []string{"rupees", "find bombs"},
		"shield":             []string{"rupees"},
		"pop bubble":         []string{"sword", "bombs", "ember seeds"},
		"remove bush":        []string{"sword", "ember seeds", "bracelet"}, // bombs could soft lock
		"hit lever":          []string{"sword", "ember seeds"},
		"seed item":          []string{"satchel", "slingshot"},
		"harvest seeds":      []string{"sword", "rod"},
		"find ember seeds":   []string{"enter d1"}, // TODO: among other places
		"get ember seeds":    []string{"find ember seeds", "harvest ember seeds"},
		"boomerang":          []string{"portal 1"},
		"rod":                []string{"portal 1"},
		"hit switch (far)":   []string{"boomerang", "bombs"},
		"winter":             []string{"rod", "hit switch (far)"},
		"shovel":             []string{"winter"},
		"find mystery seeds": []string{"d2 mystery seeds 1", "d2 mystery seeds 2"},
		"get mystery seeds":  []string{"find mystery seeds", "harvest mystery seeds"},
		"d2 key 3":           []string{"d2 key 3 1", "d2 key 3 2"},

		// later
		"jump": []string{"feather", "cape"},
	})

	r.AddParents(killNodesAnd)
	r.AddParents(killNodesOr)
	r.AddParents(d1NodesAnd)
	r.AddParents(d1NodesOr)

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
