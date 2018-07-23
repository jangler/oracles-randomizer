package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), `
Valid operations are checkGraph, verifyData, findPath, and randomize.
`[1:])
}

func main() {
	flag.Usage = usage
	flagOp := flag.String("op", "checkGraph", "operation")
	flag.Parse()

	switch *flagOp {
	case "checkGraph":
		if flag.NArg() != 0 {
			log.Fatalf("checkGraph takes 0 arguments; got %d", flag.NArg())
		}

		// validate
		if errs := checkGraph(); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}
	case "verifyData":
		if flag.NArg() != 1 {
			log.Fatalf("verify takes 1 argument; got %d", flag.NArg())
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// verify program data vs rom data
		if errs := rom.Verify(romData); err != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		} else {
			log.Print("everything OK")
		}
	case "findPath":
		if flag.NArg() < 2 {
			log.Fatalf("findPath takes 2+ arguments; got %d", flag.NArg())
		}
		if flag.NArg()%2 != 0 {
			log.Fatalf("findPath requires an even number of arguments; got %d",
				flag.NArg())
		}

		r, _ := initRoute() // ignore errors; they're diagnostic only

		// get start and end nodes
		if start, ok := r.Graph.Map[flag.Arg(0)]; ok {
			start.ClearParents()
		} else {
			log.Fatal("node %s not found", flag.Arg(0))
		}
		dest, ok := r.Graph.Map[flag.Arg(1)]
		if !ok {
			log.Fatalf("node %s not found", flag.Arg(1))
		}

		// place items in slots
		for i := 2; i < flag.NArg(); i += 2 {
			if _, ok := baseItemNodes[flag.Arg(i)]; !ok {
				log.Fatalf("%s is not an item", flag.Arg(i))
			}
			if _, ok := r.Slots[flag.Arg(i+1)]; !ok {
				log.Fatalf("%s is not a slot", flag.Arg(i+1))
			}
			r.Graph.AddParents(
				map[string][]string{flag.Arg(i): []string{flag.Arg(i + 1)}})
		}

		// try to find a valid path
		if path := findPath(r.Graph, dest); path != nil {
			for path.Len() > 0 {
				step := path.Remove(path.Front()).(string)
				log.Print(step)
			}
		} else {
			log.Print("path not found")
		}
	case "makeRoute":
		if flag.NArg() != 0 {
			log.Fatalf("makeRoute takes 0 arguments; got %d", flag.NArg())
		}

		r, _ := initRoute()
		_ = makeRoute(r, []string{"d1 essence", "d2 essence"})
	case "randomize":
		if flag.NArg() != 2 {
			log.Fatalf("randomize takes 2 arguments; got %d", flag.NArg())
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// randomize (TODO)
		if errs := randomize(romData, flag.Arg(1)); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}
	default:
		log.Fatalf("no such operation: %s", *flagOp)
	}
}

// make sure the base route graph is ok (before randomizing anything)
func checkGraph() []error {
	// TODO initRoute() does this automatically and i'm not sure it should
	_, errs := initRoute()
	return errs
}

// can be used for loading pretty much anything, really, as long as you want it
// as a slice of bytes.
func loadROM(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return rom.Load(f)
}

// attempts to find a path from the start to the given node in the graph.
// returns nil if no path was found.
func findPath(g *graph.Graph, target graph.Node) *list.List {
	path := list.New()
	mark := target.GetMark(path)
	if mark == graph.MarkTrue {
		return path
	}
	return nil
}

// attempts to create a path to the given targets by placing different items in
// slots.
func makeRoute(r *Route, targets []string) *list.List {
	// make stacks out of the item names and slot names for backtracking
	itemList := list.New()
	slotList := list.New()
	{
		// shuffle names in slices
		items := make([]string, 0, len(baseItemNodes))
		slots := make([]string, 0, len(r.Slots))
		for itemName, _ := range baseItemNodes {
			items = append(items, itemName)
		}
		for slotName, _ := range r.Slots {
			slots = append(slots, slotName)
		}
		rand.Shuffle(len(items), func(i, j int) {
			items[i], items[j] = items[j], items[i]
		})
		rand.Shuffle(len(slots), func(i, j int) {
			slots[i], slots[j] = slots[j], slots[i]
		})

		// push the shuffled items onto the stacks
		for _, itemName := range items {
			itemList.PushBack(itemName)
		}
		for _, slotName := range slots {
			slotList.PushBack(slotName)
		}
	}

	// also keep track of which items we've popped off the stacks.
	// these lists are parallel; i.e. the first item is in the first slot
	usedItems := list.New()
	usedSlots := list.New()

	if tryReachTargets(r.Graph, targets, itemList, slotList, usedItems, usedSlots) {
		log.Print("-- success")
		for _, target := range targets {
			log.Print("-- path to " + target)
			r.Graph.ClearMarks()
			path := findPath(r.Graph, r.Graph.Map[target])
			for path.Len() > 0 {
				step := path.Remove(path.Front()).(string)
				log.Print(step)
			}
		}
		log.Print("-- slotted items")
		if usedItems.Len() != usedSlots.Len() {
			log.Fatalf("FATAL: usedItems.Len() == %d; usedSlots.Len() == %d", usedItems.Len(), usedSlots.Len())
		}
		for usedItems.Len() > 0 {
			log.Printf("%s <- %s", usedItems.Remove(usedItems.Front()), usedSlots.Remove(usedSlots.Front()))
		}
	} else {
		log.Print("-- failure; something is wrong")
	}

	return nil // TODO
}

// try to reach all the given targets using the current graph status. if
// targets are unreachable, try placing an unused item in a reachable unused
// slot, and call recursively. if no combination of slots and items works,
// return false.
func tryReachTargets(g *graph.Graph, targets []string, itemList, slotList, usedItems, usedSlots *list.List) bool {
	// try to reach all targets
	if canReachTargets(g, targets) {
		return true
	}

	// try to reach each unused slot
	for i := 0; i < slotList.Len(); i++ {
		// iterate by rotating the list
		slot := slotList.Back()
		slotList.MoveToFront(slot)

		slotName := slot.Value.(string)
		g.ClearMarks() // probably redundant but safe
		if !canReachTargets(g, []string{slotName}) {
			continue
		}

		// move slot from unused to used
		usedSlots.PushBack(slotName)
		slotList.Remove(slot)

		// try placing each unused item into the slot
		for j := 0; j < itemList.Len(); j++ {
			// slot the item and move it to the used list
			itemName := itemList.Remove(itemList.Back()).(string)
			usedItems.PushBack(itemName)
			g.Map[itemName].AddParents(g.Map[slotName])

			{
				items := make([]string, 0, usedItems.Len())
				for e := usedItems.Front(); e != nil; e = e.Next() {
					items = append(items, e.Value.(string))
				}
				log.Print("trying " + strings.Join(items, " -> "))
			}

			// recurse with new state
			if tryReachTargets(g, targets, itemList, slotList, usedItems, usedSlots) {
				return true
			}

			// item didn't work; unslot it and pop it onto the front of the unused list
			usedItems.Remove(usedItems.Back())
			itemList.PushFront(itemName)
			g.Map[itemName].ClearParents()
		}

		// slot didn't work; pop it onto the front of the unused list
		usedSlots.Remove(usedSlots.Back())
		slotList.PushFront(slotName)

		// reachable slots usually equivalent in terms of routing, so don't
		// bother checking more at this point
		break
	}

	// nothing worked
	return false
}

// check if the targets are reachable using the current graph state
func canReachTargets(g *graph.Graph, targets []string) bool {
	for _, target := range targets {
		g.ClearMarks()
		if g.Map[target].GetMark(nil) != graph.MarkTrue {
			return false
		}
	}
	return true
}

// messes up rom data and writes it to a file.
//
// this also calls verify.
func randomize(romData []byte, outFilename string) []error {
	if errs := rom.Verify(romData); errs != nil {
		return errs
	}

	// XXX old code, but could be used as reference for new code
	/*
		if len(os.Args) > 2 {
			// randomize rom
			b := loadedRom.Bytes()
			rom.Mutate(b)

			// write to file
			f, err := os.Create(os.Args[2])
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			if _, err := f.Write(b); err != nil {
				log.Fatal(err)
			}
			log.Printf("wrote new ROM to %s", os.Args[2])
		}
	*/

	return []error{fmt.Errorf("NYI")}
}
