package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

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
			log.Fatalf("findPath takes 0 arguments; got %d", flag.NArg())
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
		if flag.NArg() != 1 {
			log.Fatalf("findPath takes 2 arguments; got %d", flag.NArg())
		}

		// try to find a valid path to the target node
		g, _, _ := initRoute() // ignore errors; they're diagnostic only
		target, ok := g.Map[flag.Arg(0)]
		if !ok {
			log.Fatal("target node not found")
		}
		if path := findPath(g, target); path != nil {
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

		g, openSlots, _ := initRoute()
		_ = makeRoute(g, openSlots, []string{"d1 essence", "d2 essence"})
	case "randomize":
		if flag.NArg() != 2 {
			log.Fatalf("randomize takes 2 arguments; got %d", flag.NArg())
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// randomize
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
	_, _, errs := initRoute()
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
func makeRoute(g *graph.Graph, openSlots map[string]Point, targets []string) *list.List {
	// make stacks out of the item names and slot names for backtracking
	itemList := list.New()
	slotList := list.New()
	{
		// shuffle names in slices
		items := make([]string, 0, len(baseItemNodes))
		slots := make([]string, 0, len(openSlots))
		for itemName, _ := range baseItemNodes {
			items = append(items, itemName)
		}
		for slotName, _ := range openSlots {
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

	// loop until all targets are reachable
	for {
		// try to reach an open slot
		reachedSlot := false
		for i := 0; i < slotList.Len(); i++ {
			// iterate the unused slot list by rotating it
			slot := slotList.Back()
			slotName := slot.Value.(string)
			slotList.MoveToFront(slot)

			// see if a path to the slot exists
			g.ClearMarks()
			if path := findPath(g, g.Map[slotName]); path != nil {
				// move slot from unused to used list
				usedSlots.PushBack(slotName)
				slotList.Remove(slot)

				// try placing the item at the back of the item list
				itemName := itemList.Remove(itemList.Back()).(string)
				usedItems.PushBack(itemName)
				g.Map[itemName].AddParents(g.Map[slotName])
				log.Printf("-- placing %s in %s", itemName, slotName)

				// we're good
				reachedSlot = true
				break
			}
		}
		if !reachedSlot {
			// no slot could be reached w/ new item in new slot; pop them both
			// off and put them at the front of their respective lists. iterate
			// by rotating again
			//
			// and don't forget to remove the item's parents
			log.Print("-- failure to find open slot")
			item := usedItems.Remove(usedItems.Back()).(string)
			g.Map[item].ClearParents()
			itemList.PushFront(item)
			slotList.PushFront(usedSlots.Remove(usedSlots.Back()))
			continue
		}
		// TODO the above isn't proper about actually bracktracking if *no*
		//      item gets you anywhere from the current slot.

		// if you can reach all targets from here, yr done
		reachAll := true
		for _, target := range targets {
			g.ClearMarks()
			if findPath(g, g.Map[target]) == nil {
				reachAll = false
				break
			}
		}
		if reachAll {
			log.Print("-- success")
			for _, target := range targets {
				log.Print("-- path to " + target)
				g.ClearMarks()
				path := findPath(g, g.Map[target])
				for path.Len() > 0 {
					step := path.Remove(path.Front()).(string)
					log.Print(step)
				}
			}
			log.Print("-- used items")
			if usedItems.Len() != usedSlots.Len() {
				log.Fatalf("FATAL: usedItems.Len() == %d; usedSlots.Len() == %d", usedItems.Len(), usedSlots.Len())
			}
			for usedItems.Len() > 0 {
				log.Printf("%s <- %s", usedItems.Remove(usedItems.Front()), usedSlots.Remove(usedSlots.Front()))
			}
			break
		}
	}

	return nil // TODO
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
