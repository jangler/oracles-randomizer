package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

// fatals if the command got the wrong number of arguments
func checkNumArgs(op string, expected int) {
	if flag.NArg() != expected {
		log.Printf("%s takes %d argument(s); got %d",
			op, expected, flag.NArg())
		os.Exit(2)
	}
}

func main() {
	// init flags
	flagGoal := flag.String("goal", "done",
		"comma-separated list of nodes that must be reachable")
	flagForbid := flag.String("forbid", "",
		"comma-separated list of nodes that must not be reachable")
	flagMaxlen := flag.Int("maxlen", -1,
		"if >= 0, maximum number of slotted items in the route")
	flagDryrun := flag.Bool(
		"dryrun", false, "don't write an output ROM file")
	flagDevcmd := flag.String("devcmd", "", "if given, run developer command")
	flag.Parse()

	// perform given command (or default, randomize)
	switch *flagDevcmd {
	case "pregen":
		// auto-generate some graph nodes
		checkNumArgs(*flagDevcmd, 1)

		f, err := os.Create(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		generatePrenodes(f)
	case "verify":
		checkNumArgs(*flagDevcmd, 1)

		// load rom
		romData, err := readFileBytes(flag.Arg(0))
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

		rand.Seed(time.Now().UnixNano())

		// load rom
		romData, err := readFileBytes(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// split node params
		goal := strings.Split(*flagGoal, ",")
		forbid := []string{}
		if *flagForbid != "" {
			forbid = strings.Split(*flagForbid, ",")
		}

		summary := getSummaryChannel()

		// randomize according to params
		if errs := randomize(romData, flag.Arg(1), []string{"horon village"},
			goal, forbid, *flagMaxlen, summary); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}

		close(summary)

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

// return the contents of the names file as a slice of bytes
func readFileBytes(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ioutil.ReadAll(f)
}

// messes up rom data and writes it to a file. this also calls rom.Verify().
func randomize(romData []byte, outFilename string,
	start, goal, forbid []string, maxlen int, summary chan string) []error {
	// make sure rom data is a match first
	if errs := rom.Verify(romData); errs != nil {
		return errs
	}

	// find a viable random route
	r := NewRoute(start)
	usedItems, unusedItems, usedSlots :=
		findRoute(r, start, goal, forbid, maxlen, summary)

	// place selected treasures in slots
	usedLines := make([]string, 0, usedSlots.Len())
	for usedSlots.Len() > 0 {
		slotName := usedSlots.Remove(usedSlots.Front()).(*graph.Node).Name
		treasureName := usedItems.Remove(usedItems.Front()).(*graph.Node).Name
		rom.ItemSlots[slotName].Treasure = rom.Treasures[treasureName]

		usedLines =
			append(usedLines, fmt.Sprintf("%s <- %s", slotName, treasureName))
	}

	// do it! (but don't write anything)
	checksum, err := rom.Mutate(romData)
	if err != nil {
		return []error{err}
	}

	// write info to summary file
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	summary <- ""
	summary <- "used items, in order:"
	summary <- ""
	for _, usedLine := range usedLines {
		summary <- usedLine
	}
	summary <- ""
	summary <- "unused items:"
	summary <- ""
	for e := unusedItems.Front(); e != nil; e = e.Next() {
		summary <- e.Value.(*graph.Node).Name
	}

	return nil
}
