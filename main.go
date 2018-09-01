package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"runtime/pprof"
	"strconv"
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
	flagFreewarp := flag.Bool(
		"freewarp", false, "allow unlimited tree warp (no cooldown)")
	flagKeyonly := flag.Bool(
		"keyonly", false, "only randomize key item locations")
	flagProfile := flag.String(
		"profile", "", "write CPU profile to given filename")
	flagSeed := flag.String("seed", "",
		"specific random seed to use (32-bit hex number)")
	flagUpdate := flag.Bool(
		"update", false, "update already randomized ROM to this version")
	flagVerbose := flag.Bool(
		"verbose", false, "print more detailed output to terminal")
	flag.Parse()

	checkNumArgs("randomizer", 2)

	if *flagProfile != "" {
		profFile, err := os.Create(*flagProfile)
		if err != nil {
			log.Fatal(err)
		}
		if err := pprof.StartCPUProfile(profFile); err != nil {
			log.Fatal(err)
		}
		defer profFile.Close()
		defer pprof.StopCPUProfile()
	}

	// load rom
	romData, err := readFileBytes(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	rom.SetFreewarp(*flagFreewarp)

	// randomize according to params, unless we're just updating
	if *flagUpdate {
		_, err := rom.Update(romData)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		seed := setRandomSeed(*flagSeed)
		if *flagSeed == "" {
			seed = 0 // none specified, not an actual zero seed
		}

		summary, summaryDone := getSummaryChannel()

		if errs := randomize(romData, flag.Arg(1), *flagKeyonly, *flagVerbose,
			seed, summary); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}

		close(summary)
		<-summaryDone
	}

	// write to file
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

// parses a delimited (e.g. with comma) command-line argument, stripping spaces
// around each entry.
func parseDelimitedArg(arg, delimiter string) []string {
	a := make([]string, 0)

	for _, s := range strings.Split(arg, delimiter) {
		a = append(a, strings.TrimSpace(s))
	}

	return a
}

// sets a 32-bit unsigned random seed based on a hexstring, if non-empty, or
// else the current time, and returns that seed.
func setRandomSeed(hexString string) uint32 {
	seed := uint32(time.Now().UnixNano())
	if hexString != "" {
		v, err := strconv.ParseUint(
			strings.Replace(hexString, "0x", "", 1), 16, 32)
		if err != nil {
			log.Fatalf(`fatal: invalid seed "%s"`, hexString)
		}
		seed = uint32(v)
	}
	rand.Seed(int64(seed))

	return seed
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
func randomize(romData []byte, outFilename string, keyonly, verbose bool,
	seed uint32, summary chan string) []error {
	// make sure rom data is a match first
	if errs := rom.Verify(romData); errs != nil {
		return errs
	}

	// give each routine its own random source, so that they can return the
	// seeds that they used. if a specific seed was specified, only use one
	// thread.
	numThreads := 1
	if !verbose && seed == 0 {
		numThreads = runtime.NumCPU()
	}
	log.Printf("using %d thread(s)", numThreads)
	sources := make([]rand.Source, numThreads)
	seeds := make([]uint32, numThreads)
	for i := 0; i < numThreads; i++ {
		if seed == 0 {
			randSeed := uint32(rand.Int63())
			sources[i] = rand.NewSource(int64(randSeed))
			seeds[i] = randSeed
		} else {
			sources[i] = rand.NewSource(int64(seed))
			seeds[i] = seed
		}
	}

	// search for route, parallelized
	routeChan := make(chan *RouteLists)
	logChan := make(chan string)
	stopLogChan := make(chan int)
	doneChan := make(chan int)
	for i := 0; i < numThreads; i++ {
		go searchAsync(rand.New(sources[i]), seeds[i], keyonly, verbose,
			logChan, routeChan, doneChan)
	}

	// log messages from all threads
	go func() {
		for {
			select {
			case msg := <-logChan:
				log.Print(msg)
			case <-stopLogChan:
				return
			}
		}
	}()

	// get return values
	var rl *RouteLists
	for i := 0; i < numThreads; i++ {
		rl = <-routeChan
		if rl != nil {
			break
		}
	}

	// tell all the other routines to stop
	stopLogChan <- 1
	go func() {
		for {
			doneChan <- 1
		}
	}()

	// didn't find any route
	if rl == nil {
		log.Fatal("fatal: no route found")
	}
	log.Printf("route found; seed %08x", rl.Seed)

	// place selected treasures in slots
	usedLines := make([]string, 0, rl.UsedSlots.Len())
	for rl.UsedSlots.Len() > 0 {
		slotName :=
			rl.UsedSlots.Remove(rl.UsedSlots.Front()).(*graph.Node).Name
		treasureName :=
			rl.UsedItems.Remove(rl.UsedItems.Front()).(*graph.Node).Name
		rom.ItemSlots[slotName].Treasure = rom.Treasures[treasureName]

		usedLines = append(usedLines, fmt.Sprintf("%-28s <- %s",
			getNiceName(slotName), getNiceName(treasureName)))
	}

	// set rom seasons and animal data
	for area, id := range rl.Seasons {
		rom.Seasons[fmt.Sprintf("%s season", area)].New = []byte{id}
	}
	rom.SetAnimal(rl.Companion)

	// do it! (but don't write anything)
	checksum, err := rom.Mutate(romData)
	if err != nil {
		return []error{err}
	}

	// write info to summary file
	summary <- fmt.Sprintf("seed: %08x", rl.Seed)
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	summary <- ""
	summary <- "used items, in (one possible) order:"
	summary <- ""
	if keyonly {
		for _, usedLine := range usedLines {
			summary <- usedLine
		}
	} else {
		// print boss keys, maps, and compasses last, even though they're
		// slotted first
		for _, usedLine := range usedLines[22:] {
			summary <- usedLine
		}
		for _, usedLine := range usedLines[:22] {
			summary <- usedLine
		}
	}
	if rl.UnusedItems.Len() > 0 {
		summary <- ""
		summary <- "unused items:"
		summary <- ""
		for e := rl.UnusedItems.Front(); e != nil; e = e.Next() {
			summary <- e.Value.(*graph.Node).Name
		}
	}

	summary <- ""
	summary <- "default seasons:"
	summary <- ""
	for name, area := range rom.Seasons {
		summary <- fmt.Sprintf("%s - %s", name, seasonsByID[int(area.New[0])])
	}

	return nil
}

// searches for a route and logs and returns a route on the given channels.
func searchAsync(src *rand.Rand, seed uint32, keyonly, verbose bool,
	logChan chan string, retChan chan *RouteLists, doneChan chan int) {
	// find a viable random route
	r := NewRoute()
	retChan <- findRoute(src, seed, r, keyonly, verbose, logChan, doneChan)
}
