package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
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

func gameName(game int) string {
	switch game {
	case rom.GameAges:
		return "ooa"
	case rom.GameSeasons:
		return "oos"
	default:
		return "UNKNOWN"
	}
}

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [<original file> [<new file>]]\n", os.Args[0])
	flag.PrintDefaults()
}

// fatal prints the given error to stderr, waits for user input if `wait` is
// true, then exits with status 1.
func fatal(err error, wait bool) {
	fmt.Fprintf(os.Stderr, "fatal: %v.\n", err)
	if wait {
		fmt.Fprint(os.Stderr, "press enter to continue.")
		os.Stdin.Read(make([]byte, 1))
	}
	os.Exit(1)
}

func main() {
	// init flags
	flag.Usage = usage
	flagHard := flag.Bool(
		"hard", false, "require some plays outside normal logic")
	flagNoMusic := flag.Bool(
		"nomusic", false, "don't play any music in the modified ROM")
	flagProfile := flag.String(
		"profile", "", "write CPU profile to given filename")
	flagSeed := flag.String("seed", "",
		"specific random seed to use (32-bit hex number)")
	flagVerbose := flag.Bool(
		"verbose", false, "print more detailed output to terminal")
	flag.Parse()

	// turn profiling on if specified
	if *flagProfile != "" {
		profFile, err := os.Create(*flagProfile)
		if err != nil {
			fatal(err, false)
		}
		if err := pprof.StartCPUProfile(profFile); err != nil {
			fatal(err, false)
		}
		defer profFile.Close()
		defer pprof.StopCPUProfile()
	}

	switch flag.NArg() {
	case 0: // no specified files, invalid
		fatal(fmt.Errorf("no input ROM specified"), true)
	case 1: // specified input file only, assume not using command line
		b, game, err := readGivenROM(flag.Arg(0))
		if err != nil {
			fatal(err, true)
		}

		if err := handleFile(b, game, flag.Arg(0), *flagSeed, *flagNoMusic,
			*flagHard, *flagVerbose); err != nil {
			fatal(err, true)
		}

		fmt.Fprint(os.Stderr, "press enter to continue.")
		os.Stdin.Read(make([]byte, 1))
	case 2: // specified input and output file, so using command line
		b, game, err := readGivenROM(flag.Arg(0))
		if err != nil {
			fatal(err, false)
		}

		// operate on file
		var sum []byte
		var seed uint32
		var logFilename string
		fmt.Printf("randomizing %s\n", flag.Arg(0))
		seed, sum, logFilename, err = randomize(b, game, *flagSeed,
			*flagNoMusic, *flagHard, *flagVerbose)
		if err != nil {
			fatal(err, false)
		}

		// write file
		if err := writeROM(b, flag.Arg(1), logFilename, seed, sum); err != nil {
			fatal(err, false)
		}
	default:
		flag.Usage()
	}
}

// attempt to write rom data to a file and print summary info.
func writeROM(b []byte, filename, logFilename string, seed uint32,
	sum []byte) error {
	// write file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}

	// print summary
	fmt.Printf("seed: %08x\n", seed)
	fmt.Printf("sha-1 sum: %x\n", string(sum))
	fmt.Printf("wrote new rom to %s\n", filename)
	fmt.Printf("wrote log file to %s\n", logFilename)

	return nil
}

// search for a vanilla US seasons rom in the current directory, and return it
// as a byte slice if possible.
func findVanillaROM() ([]byte, error) {
	// read slice of file info from working dir
	dirName, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dir, err := os.Open(dirName)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, info := range files {
		// check file metadata
		if info.Size() != 1048576 {
			continue
		}

		// read file
		f, err := os.Open(info.Name())
		if err != nil {
			return nil, err
		}
		defer f.Close()
		b, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, err
		}

		// check file data
		if rom.IsSeasons(b) && rom.IsUS(b) && rom.IsVanilla(b) {
			fmt.Printf("found vanilla ROM: %s\n", info.Name())
			return b, nil
		}
	}

	return nil, fmt.Errorf("no vanilla ROM found in working directory")
}

// read the specified file into a slice of bytes, returning an error if the
// read fails or if the file is an invalid rom. also returns the game as an
// int.
func readGivenROM(filename string) ([]byte, int, error) {
	// read file
	f, err := os.Open(filename)
	if err != nil {
		return nil, rom.GameNil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, rom.GameNil, err
	}

	// check file data
	if !rom.IsAges(b) && !rom.IsSeasons(b) {
		return nil, rom.GameNil,
			fmt.Errorf("%s is not an Oracles ROM", filename)
	}
	if !rom.IsUS(b) {
		return nil, rom.GameNil,
			fmt.Errorf("%s is a JP ROM; only US is supported", filename)
	}
	if !rom.IsVanilla(b) {
		return nil, rom.GameNil,
			fmt.Errorf("%s is an unrecognized Oracles ROM", filename)
	}

	game := rom.GameAges
	if rom.IsSeasons(b) {
		game = rom.GameSeasons
	}
	return b, game, nil
}

// decide whether to randomize or update the file
func handleFile(romData []byte, game int, filename, seedFlag string,
	noMusic, hard, verbose bool) error {
	var seed uint32
	var sum []byte
	var err error
	var outName, logFilename string

	// operate on rom data
	fmt.Printf("randomizing %s\n", flag.Arg(0))
	seed, sum, logFilename, err =
		randomize(romData, game, seedFlag, noMusic, hard, verbose)
	if err != nil {
		return err
	}
	hardString := ""
	if hard {
		hardString = "_hard"
	}
	outName = fmt.Sprintf("%srando_%s_%08x%s.gbc",
		gameName(game), version, seed, hardString)

	// write to file
	return writeROM(romData, outName, logFilename, seed, sum)
}

// sets a 32-bit unsigned random seed based on a hexstring, if non-empty, or
// else the current time, and returns that seed.
func setRandomSeed(hexString string) uint32 {
	seed := uint32(time.Now().UnixNano())
	if hexString != "" {
		v, err := strconv.ParseUint(
			strings.Replace(hexString, "0x", "", 1), 16, 32)
		if err != nil {
			fatal(fmt.Errorf(`fatal: invalid seed "%s"`, hexString), false)
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
func randomize(romData []byte, game int, seedFlag string,
	noMusic, hard, verbose bool) (uint32, []byte, string, error) {
	// make sure rom data is a match first
	rom.Init(game)
	if errs := rom.Verify(romData, game); errs != nil {
		return 0, nil, "", errs[0]
	}

	if noMusic {
		rom.SetNoMusic()
	}

	seed := setRandomSeed(seedFlag)
	if seedFlag == "" {
		seed = 0 // none specified, not an actual zero seed
	}

	// give each routine its own random source, so that they can return the
	// seeds that they used. if a specific seed was specified, only use one
	// thread.
	numThreads := 1
	if !verbose && seed == 0 {
		numThreads = runtime.NumCPU()
	}
	fmt.Printf("using %d thread(s)\n", numThreads)
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
	routeChan := make(chan *RouteInfo)
	logChan := make(chan string)
	stopLogChan := make(chan int)
	doneChan := make(chan int)
	for i := 0; i < numThreads; i++ {
		go searchAsync(rand.New(sources[i]), game, seeds[i], hard, verbose,
			logChan, routeChan, doneChan)
	}

	// log messages from all threads
	go func() {
		for {
			select {
			case msg := <-logChan:
				fmt.Println(msg)
			case <-stopLogChan:
				return
			}
		}
	}()

	// get return values
	var ri *RouteInfo
	for i := 0; i < numThreads; i++ {
		ri = <-routeChan
		if ri != nil {
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
	if ri == nil {
		return 0, nil, "", fmt.Errorf("no route found")
	}

	// place selected treasures in slots
	for ri.UsedSlots.Len() > 0 {
		slotName :=
			ri.UsedSlots.Remove(ri.UsedSlots.Front()).(*graph.Node).Name
		treasureName :=
			ri.UsedItems.Remove(ri.UsedItems.Front()).(*graph.Node).Name
		if verbose {
			fmt.Printf("%s <- %s\n", slotName, treasureName)
		}
		rom.ItemSlots[slotName].Treasure = rom.Treasures[treasureName]
	}

	// set season data
	if game == rom.GameSeasons {
		for area, id := range ri.Seasons {
			rom.Seasons[fmt.Sprintf("%s season", area)].New = []byte{id}
		}
	}

	rom.SetAnimal(ri.Companion)

	// do it! (but don't write anything)
	checksum, err := rom.Mutate(romData, game)
	if err != nil {
		return 0, nil, "", err
	}

	hardString := ""
	if hard {
		hardString = "hard_"
	}
	logFilename := fmt.Sprintf("%srando_%s_%08x_%slog.txt",
		gameName(game), version, ri.Seed, hardString)
	summary, summaryDone := getSummaryChannel(logFilename)

	// write info to summary file
	summary <- fmt.Sprintf("seed: %08x", ri.Seed)
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	if hard {
		summary <- fmt.Sprintf("difficulty: hard")
	} else {
		summary <- fmt.Sprintf("difficulty: normal")
	}
	logItems(summary, "required items", ri.ProgressItems, ri.ProgressSlots)
	logItems(summary, "optional items", ri.ExtraItems, ri.ExtraSlots)
	if game == rom.GameSeasons {
		summary <- ""
		summary <- "default seasons:"
		summary <- ""
		for name, area := range rom.Seasons {
			summary <- fmt.Sprintf("%-15s <- %s",
				name[:len(name)-7], seasonsByID[int(area.New[0])])
		}
		summary <- ""
		summary <- fmt.Sprintf("natzu region <- %s", []string{
			"", "natzu prairie", "natzu river", "natzu wasteland",
		}[ri.Companion])
	} else {
		summary <- ""
		summary <- fmt.Sprintf("animal companion <- %s", []string{
			"", "ricky", "dimitri", "moosh",
		}[ri.Companion])
	}

	close(summary)
	<-summaryDone

	return ri.Seed, checksum, logFilename, nil
}

// searches for a route and logs and returns a route on the given channels.
func searchAsync(src *rand.Rand, game int, seed uint32, hard, verbose bool,
	logChan chan string, retChan chan *RouteInfo, doneChan chan int) {
	// find a viable random route
	retChan <- findRoute(src, game, seed, hard, verbose, logChan, doneChan)
}

// send lines of item/slot info to a summary channel. this is a destructive
// operation on the lists.
func logItems(summary chan string, title string, items, slots *list.List) {
	summary <- ""
	summary <- title + ":"
	summary <- ""

	for slots.Len() > 0 {
		slotName := slots.Remove(slots.Front()).(*graph.Node).Name
		itemName := items.Remove(items.Front()).(*graph.Node).Name
		summary <- fmt.Sprintf("%-28s <- %s",
			getNiceName(slotName), getNiceName(itemName))
	}
}
