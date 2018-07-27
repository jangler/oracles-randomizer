package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

func checkNumArgs(op string, expected int) {
	if flag.NArg() != expected {
		log.Printf("%s takes %d argument(s); got %d",
			op, expected, flag.NArg())
		os.Exit(2)
	}
}

func main() {
	// init flags
	flagStart := flag.String("start", "horon village",
		"comma-separated list of nodes treated as given")
	flagGoal := flag.String("goal", "done",
		"comma-separated list of nodes that must be reachable")
	flagForbid := flag.String("forbid", "",
		"comma-separated list of nodes that must not be reachable")
	flagMaxlen := flag.Int("maxlen", -1,
		"if >= 0, maximum number of slotted items in the route")
	flagDryrun := flag.Bool(
		"dryrun", false, "don't write an output file for any operation")
	flagDevcmd := flag.String("devcmd", "", "if given, run developer command")
	flag.Parse()

	// perform given command (or default, randomize)
	switch *flagDevcmd {
	case "checkgraph":
		checkNumArgs(*flagDevcmd, 0)

		// check for orphan/childless nodes
		r := initRoute(strings.Split(*flagStart, ","))
		if errs := r.CheckGraph(); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
		}
	case "explore":
		r := initRoute([]string{})

		// args are starting nodes
		reached := make(map[graph.Node]bool)
		for _, arg := range flag.Args() {
			log.Print("adding ", arg)
			reached = r.Graph.Explore(reached,
				[]graph.Node{r.Graph.Map[arg]}, nil)
			/*
				log.Print(len(reached), " nodes")
				log.Print(reached)
			*/
			steps := make(map[graph.Node]bool)
			for node := range reached {
				if node.IsStep() {
					steps[node] = true
				}
			}
			log.Print(len(steps), " step nodes")
			log.Print(steps)
		}
	case "pointgen":
		checkNumArgs(*flagDevcmd, 1)

		f, err := os.Create(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		generatePoints(f)
	case "verify":
		checkNumArgs(*flagDevcmd, 1)

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// verify program data vs rom data
		if errs := rom.Verify(romData); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		} else {
			log.Print("everything OK")
		}
	case "": // normal behavior (randomize)
		if *flagDryrun {
			checkNumArgs("dryrun", 1)
		} else {
			checkNumArgs("randomizer", 2)
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// split node params
		start := strings.Split(*flagStart, ",")
		goal := strings.Split(*flagGoal, ",")
		forbid := []string{}
		if *flagForbid != "" {
			forbid = strings.Split(*flagForbid, ",")
		}

		// randomize according to params
		if errs := randomize(romData, flag.Arg(1),
			start, goal, forbid, *flagMaxlen); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}

		// write to file unless it's a dry run
		if !*flagDryrun {
			f, err := os.Create(flag.Arg(1))
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()
			if _, err := f.Write(romData); err != nil {
				log.Fatal(err)
			}
			log.Printf("wrote new ROM to %s", flag.Arg(1))
		}
	default:
		log.Printf("no such devcmd: %s", *flagDevcmd)
		os.Exit(2)
	}
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

// messes up rom data and writes it to a file. this also calls rom.Verify().
func randomize(romData []byte, outFilename string,
	start, goal, forbid []string, maxlen int) []error {
	// make sure rom data is a match first
	if errs := rom.Verify(romData); errs != nil {
		return errs
	}

	// find a viable random route
	r := initRoute(start)
	usedItems, usedSlots := makeRoute(r, start, goal, forbid, maxlen)

	// place selected treasures in slots
	for usedItems.Len() > 0 {
		slotName := usedSlots.Remove(usedSlots.Front()).(graph.Node).Name()
		treasureName := usedItems.Remove(usedItems.Front()).(graph.Node).Name()
		if err := placeTreasureInSlot(treasureName, slotName); err != nil {
			return []error{err}
		}
	}

	rom.Mutate(romData)

	return nil
}

func placeTreasureInSlot(treasureName, slotName string) error {
	if slot, ok := rom.ItemSlots[slotName]; ok {
		if treasure, ok := rom.Treasures[treasureName]; ok {
			slot.Treasure = treasure
		} else {
			return fmt.Errorf("no treasure '%s' in ROM code", treasureName)
		}
	} else {
		return fmt.Errorf("no item slot '%s' in ROM code", slotName)
	}

	return nil
}
