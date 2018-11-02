package main

import (
	"container/list"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
	"github.com/jangler/oos-randomizer/ui"
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

// fatal prints an error to the UI.
func fatal(err error) {
	ui.Printf("fatal: %v.", err)
}

var (
	flagHard, flagNoMusic, flagTreewarp, flagVerbose bool

	flagSeed string
)

func main() {
	// init flags
	flag.Usage = usage
	flag.BoolVar(&flagHard, "hard", false,
		"require some plays outside normal logic")
	flag.BoolVar(&flagNoMusic, "nomusic", false,
		"don't play any music in the modified ROM")
	flag.StringVar(&flagSeed, "seed", "",
		"specific random seed to use (32-bit hex number)")
	flag.BoolVar(&flagTreewarp, "treewarp", false,
		"warp to ember tree by pressing start+B on map screen")
	flag.BoolVar(&flagVerbose, "verbose", false,
		"print more detailed output to terminal")
	flag.Parse()

	ui.Init("oracles randomizer " + version)
	go runRandomizer()
	ui.Run()
}

func runRandomizer() {
	defer ui.Done()

	// if rom is to be randomized, infile must be non-empty after switch
	var infile, outfile string

	switch flag.NArg() {
	case 0: // no specified files, search in executable's directory
		dir, seasons, ages, err := findVanillaROMs()
		if err != nil {
			fatal(err)
			break
		}

		// print which files, if any, are found.
		if seasons != "" {
			ui.PrintPath("found vanilla US seasons ROM: ", seasons, "")
		} else {
			ui.Printf("no vanilla US seasons ROM found.")
		}
		if ages != "" {
			ui.PrintPath("found vanilla US ages ROM: ", ages, "")
		} else {
			ui.Printf("no vanilla US ages ROM found.")
		}
		ui.Printf("")

		// determine which filename to use based on what roms are found, and on
		// user input.
		if seasons == "" && ages == "" {
			ui.Printf("no ROMs found in program's directory, " +
				"and no ROMs specified.")
		} else if seasons != "" && ages != "" {
			which := ui.Prompt("randomize (s)easons or (a)ges?")
			if which == 's' {
				infile = seasons
			} else {
				infile = ages
			}
		} else if seasons != "" {
			infile = seasons
		} else {
			infile = ages
		}

		seasons = filepath.Join(dir, seasons)
		ages = filepath.Join(dir, ages)
	case 1: // specified input file only
		infile = flag.Arg(0)
	case 2: // specified input and output file
		infile, outfile = flag.Arg(0), flag.Arg(1)
	default:
		flag.Usage()
	}

	if infile != "" {
		if _, game, err := readGivenROM(infile); err != nil {
			fatal(err)
			return
		} else {
			rom.Init(game)
		}
		ui.Printf("randomizing %s.", infile)

		// prompt for options if it wasn't necessarily a CLI invocation

		if flag.NArg() != 2 {
			difficulty := ui.Prompt("difficulty: (n)ormal or (h)ard?")
			flagHard = difficulty == 'h'
		}
		if flagHard {
			ui.Printf("using hard difficulty.")
		} else {
			ui.Printf("using normal difficulty.")
		}

		if flag.NArg() != 2 {
			music := ui.Prompt("(m)usic or (n)o music?")
			flagNoMusic = music == 'n'
		}
		if flagNoMusic {
			ui.Printf("music off.")
		} else {
			ui.Printf("music on.")
		}

		if flag.NArg() != 2 {
			treewarp := ui.Prompt("enable tree warp? (y/n)")
			flagTreewarp = treewarp == 'y'
		}
		if flagTreewarp {
			ui.Printf("tree warp on.")
		} else {
			ui.Printf("tree warp off.")
		}

		ui.Printf("")

		rom.SetMusic(!flagNoMusic)
		rom.SetTreewarp(flagTreewarp)

		if err := randomizeFile(infile, outfile); err != nil {
			fatal(err)
		}
	}
}

func randomizeFile(infile, outfile string) error {
	b, game, err := readGivenROM(infile)
	if err != nil {
		return err
	}

	if err := handleFile(b, game, infile, flagSeed,
		flagHard, flagVerbose); err != nil {
		return err
	}

	return nil
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
	ui.Printf("seed: %08x\n", seed)
	ui.Printf("SHA-1 sum: %x\n", string(sum))
	ui.Printf("wrote new ROM to %s\n", filename)
	ui.Printf("wrote log file to %s\n", logFilename)

	return nil
}

// search for a vanilla US seasons and ages ROMs in the executable's directory,
// and return their filenames.
func findVanillaROMs() (dirName, seasons, ages string, err error) {
	// read slice of file info from executable's dir
	exe, err := os.Executable()
	if err != nil {
		return
	}

	dirName = filepath.Dir(exe)
	ui.PrintPath("searching ", dirName, " for ROMs.")
	dir, err := os.Open(dirName)
	if err != nil {
		return
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return
	}

	for _, info := range files {
		// check file metadata
		if info.Size() != 1048576 {
			continue
		}

		// read file
		var f *os.File
		f, err = os.Open(filepath.Join(dirName, info.Name()))
		if err != nil {
			return
		}
		defer f.Close()
		var b []byte
		b, err = ioutil.ReadAll(f)
		if err != nil {
			return
		}

		// check file data
		if rom.IsUS(b) && rom.IsVanilla(b) {
			if rom.IsAges(b) {
				ages = info.Name()
			} else {
				seasons = info.Name()
			}
		}

		if ages != "" && seasons != "" {
			break
		}
	}

	return
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
			fmt.Errorf("%s is not an oracles ROM", filename)
	}
	if !rom.IsUS(b) {
		return nil, rom.GameNil,
			fmt.Errorf("%s is a JP ROM; only US is supported", filename)
	}
	if !rom.IsVanilla(b) {
		return nil, rom.GameNil,
			fmt.Errorf("%s is an unrecognized oracles ROM", filename)
	}

	game := rom.GameAges
	if rom.IsSeasons(b) {
		game = rom.GameSeasons
	}
	return b, game, nil
}

// decide whether to randomize or update the file
func handleFile(romData []byte, game int, filename, seedFlag string,
	hard, verbose bool) error {
	var seed uint32
	var sum []byte
	var err error
	var outName, logFilename string

	// operate on rom data
	seed, sum, logFilename, err =
		randomize(romData, game, seedFlag, hard, verbose)
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
func setRandomSeed(hexString string) (uint32, error) {
	seed := uint32(time.Now().UnixNano())
	if hexString != "" {
		v, err := strconv.ParseUint(
			strings.Replace(hexString, "0x", "", 1), 16, 32)
		if err != nil {
			return 0, fmt.Errorf(`invalid seed "%s"`, hexString)
		}
		seed = uint32(v)
	}
	rand.Seed(int64(seed))

	return seed, nil
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

// messes up rom data and writes it to a file.
func randomize(romData []byte, game int, seedFlag string,
	hard, verbose bool) (uint32, []byte, string, error) {
	seed, err := setRandomSeed(seedFlag)
	if err != nil {
		return 0, nil, "", err
	}
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
	ui.Printf("using %d thread(s).", numThreads)
	seeds := make([]uint32, numThreads)
	for i := 0; i < numThreads; i++ {
		if seed == 0 {
			randSeed := uint32(rand.Int63())
			seeds[i] = randSeed
		} else {
			seeds[i] = seed
		}
	}

	// search for route, parallelized
	routeChan := make(chan *RouteInfo)
	logChan := make(chan string)
	stopLogChan := make(chan int)
	doneChan := make(chan int)
	for i := 0; i < numThreads; i++ {
		go searchAsync(game, seeds[i], hard, verbose,
			logChan, routeChan, doneChan)
	}

	// log messages from all threads
	go func() {
		for {
			select {
			case msg := <-logChan:
				ui.Printf(msg)
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
			ui.Printf("%s <- %s\n", slotName, treasureName)
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
func searchAsync(game int, seed uint32, hard, verbose bool,
	logChan chan string, retChan chan *RouteInfo, doneChan chan int) {
	// find a viable random route
	retChan <- findRoute(game, seed, hard, verbose, logChan, doneChan)
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
