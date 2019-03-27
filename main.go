package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jangler/oracles-randomizer/hints"
	"github.com/jangler/oracles-randomizer/rom"
	"github.com/jangler/oracles-randomizer/ui"
)

type logFunc func(string, ...interface{})

// gameName returns the short name associated with a game number.
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

// usage is called when an invalid CLI invocation is used, or if the -h flag is
// passed.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [<original file> [<new file>]]\n", os.Args[0])
	flag.PrintDefaults()
}

// fatal prints an error to whichever UI is used.
func fatal(err error, logf logFunc) {
	logf("fatal: %v.", err)
}

// options specified on the command line or via the TUI
var (
	flagHard     bool
	flagN        int
	flagNoMusic  bool
	flagNoUI     bool
	flagSeed     string
	flagStats    string
	flagTreewarp bool
	flagVerbose  bool
)

// initFlags initializes the CLI/TUI option values and variables.
func initFlags() {
	flag.Usage = usage
	flag.BoolVar(&flagHard, "hard", false,
		"require some plays outside normal logic")
	flag.IntVar(&flagN, "n", 100,
		"number of trials for stats")
	flag.BoolVar(&flagNoMusic, "nomusic", false,
		"don't play any music in the modified ROM")
	flag.BoolVar(&flagNoUI, "noui", false,
		"use command line output without option prompts")
	flag.StringVar(&flagSeed, "seed", "",
		"specific random seed to use (32-bit hex number)")
	flag.StringVar(&flagStats, "stats", "",
		"test routes and print stats for 'seasons' or 'ages'")
	flag.BoolVar(&flagTreewarp, "treewarp", false,
		"warp to ember tree by pressing start+B on map screen")
	flag.BoolVar(&flagVerbose, "verbose", false,
		"print more detailed output to terminal")
	flag.Parse()
}

// main is the program's entry point.
func main() {
	initFlags()

	if flagStats != "" {
		// do stats instead of randomizing
		var game int

		if flagStats == "seasons" {
			game = rom.GameSeasons
		} else if flagStats == "ages" {
			game = rom.GameAges
		} else {
			fmt.Printf("'%s' is invalid. try 'seasons' or 'ages'.\n", flagStats)
			return
		}

		rom.Init(game)
		rand.Seed(time.Now().UnixNano())
		logStats(game, flagN, flagHard, func(s string, a ...interface{}) {
			fmt.Printf(s, a...)
			fmt.Println()
		})
	} else if flag.NArg()+flag.NFlag() > 1 { // CLI used
		// run randomizer on main goroutine
		runRandomizer(false, func(s string, a ...interface{}) {
			fmt.Printf(s, a...)
			fmt.Println()
		})
	} else { // CLI maybe not used
		// run TUI on main goroutine and randomizer on alternate goroutine
		ui.Init("oracles randomizer " + version)
		go runRandomizer(true, func(s string, a ...interface{}) {
			ui.Printf(s, a...)
		})
		ui.Run()
	}
}

// run the main randomizer routine, printing messages via logf, which should
// act analogously to fmt.Printf with added newline.
func runRandomizer(useTUI bool, logf logFunc) {
	// close TUI after randomizer is done
	defer func() {
		if useTUI {
			ui.Done()
		}
	}()

	// if rom is to be randomized, infile must be non-empty after switch
	var dirName, infile, outfile string
	switch flag.NArg() {
	case 0: // no specified files, search in executable's directory
		var seasons, ages string
		var err error
		dirName, seasons, ages, err = findVanillaROMs()
		if err != nil {
			fatal(err, logf)
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
	case 1: // specified input file only
		infile = flag.Arg(0)
	case 2: // specified input and output file
		infile, outfile = flag.Arg(0), flag.Arg(1)
	default:
		flag.Usage()
	}

	if infile != "" {
		b, game, err := readGivenROM(filepath.Join(dirName, infile))
		if err != nil {
			fatal(err, logf)
			return
		} else {
			rom.Init(game)
		}
		logf("randomizing %s.", infile)

		getAndLogOptions(useTUI, logf)

		if useTUI {
			logf("")
		}

		rom.SetMusic(!flagNoMusic)
		rom.SetTreewarp(flagTreewarp)

		if err := randomizeFile(b, game, dirName, outfile, flagSeed,
			flagHard, flagVerbose, logf); err != nil {
			fatal(err, logf)
			return
		}
	}
}

// getAndLogOptions logs values of selected options, prompting for them first
// if the TUI is used.
func getAndLogOptions(useTUI bool, logf logFunc) {
	if useTUI {
		if ui.Prompt("use specific seed? (y/n)") == 'y' {
			flagSeed = ui.PromptSeed("enter seed: (8-digit hex number)")
			logf("using seed %s.", flagSeed)
		}
	}

	if useTUI {
		flagHard = ui.Prompt("enable hard difficulty? (y/n)") == 'y'
	}
	if flagHard {
		logf("using hard difficulty.")
	} else {
		logf("using normal difficulty.")
	}

	if useTUI {
		flagNoMusic = ui.Prompt("disable music? (y/n)") == 'y'
	}
	if flagNoMusic {
		logf("music off.")
	} else {
		logf("music on.")
	}

	if useTUI {
		flagTreewarp = ui.Prompt("enable tree warp? (y/n)") == 'y'
	}
	if flagTreewarp {
		logf("tree warp on.")
	} else {
		logf("tree warp off.")
	}
}

// attempt to write rom data to a file and print summary info.
func writeROM(b []byte, dirName, filename, logFilename string, seed uint32,
	sum []byte, logf logFunc) error {
	// write file
	f, err := os.Create(filepath.Join(dirName, filename))
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(b); err != nil {
		return err
	}

	// print summary
	logf("seed: %08x", seed)
	logf("SHA-1 sum: %x", string(sum))
	logf("wrote new ROM to %s", filename)
	logf("wrote log file to %s", logFilename)

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

func randomizeFile(romData []byte, game int, dirName, outfile, seedFlag string,
	hard, verbose bool, logf logFunc) error {
	var seed uint32
	var sum []byte
	var err error
	var logFilename string

	// operate on rom data
	if outfile != "" {
		logFilename = outfile[:len(outfile)-4] + "_log.txt"
	}
	seed, sum, logFilename, err = randomize(
		romData, game, dirName, logFilename, seedFlag, hard, verbose, logf)
	if err != nil {
		return err
	}
	hardString := ""
	if hard {
		hardString = "_hard"
	}
	if outfile == "" {
		outfile = fmt.Sprintf("%srando_%s_%08x%s.gbc",
			gameName(game), version, seed, hardString)
	}

	// write to file
	return writeROM(romData, dirName, outfile, logFilename, seed, sum, logf)
}

// setRandomSeed sets a 32-bit unsigned random seed based on a hexstring, if
// non-empty, or else the current time, and returns that seed.
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

// messes up rom data and writes it to a file.
func randomize(romData []byte, game int, dirName, logFilename, seedFlag string,
	hard, verbose bool, logf logFunc) (uint32, []byte, string, error) {
	// sanity check beforehand
	if errs := rom.Verify(romData, game); errs != nil {
		if verbose {
			for _, err := range errs {
				logf(err.Error())
			}
		}
		return 0, nil, "", errs[0]
	}

	seed, err := setRandomSeed(seedFlag)
	if err != nil {
		return 0, nil, "", err
	}

	// search for route
	ri := findRoute(game, seed, hard, verbose, logf)
	if ri == nil {
		return 0, nil, "", fmt.Errorf("no route found")
	}

	checks := getChecks(ri)
	spheres := getSpheres(ri.Route.Graph, checks, hard)
	owlHints := hints.Generate(ri.Src, ri.Route.Graph, checks,
		rom.GetOwlNames(game), game, hard)

	checksum, err := setROMData(romData, game, ri, owlHints, logf, verbose)
	if err != nil {
		return 0, nil, "", err
	}

	hardString := ""
	if hard {
		hardString = "hard_"
	}
	if logFilename == "" {
		logFilename = fmt.Sprintf("%srando_%s_%08x_%slog.txt",
			gameName(game), version, ri.Seed, hardString)
	}
	summary, summaryDone := getSummaryChannel(
		filepath.Join(dirName, logFilename))

	// write info to summary file
	summary <- fmt.Sprintf("seed: %08x", ri.Seed)
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	if hard {
		summary <- fmt.Sprintf("difficulty: hard")
	} else {
		summary <- fmt.Sprintf("difficulty: normal")
	}
	summary <- ""
	summary <- ""
	summary <- "-- progression items --"
	summary <- ""
	logSpheres(summary, checks, spheres,
		func(name string) bool { return !itemIsJunk(name) })
	summary <- ""
	summary <- "-- other items --"
	summary <- ""
	logSpheres(summary, checks, spheres, itemIsJunk)
	if game == rom.GameSeasons {
		summary <- ""
		summary <- "-- default seasons --"
		summary <- ""
		for name, area := range rom.Seasons {
			summary <- fmt.Sprintf("%-15s <- %s",
				name[:len(name)-7], seasonsByID[int(area.New[0])])
		}
	}
	summary <- ""
	summary <- ""
	summary <- "-- hints --"
	summary <- ""
	for owlName, hint := range owlHints {
		summary <- fmt.Sprintf("%-20s <- \"%s\"", owlName,
			strings.ReplaceAll(strings.ReplaceAll(hint, "\n", " "), "  ", " "))
	}

	close(summary)
	<-summaryDone

	return ri.Seed, checksum, logFilename, nil
}

// itemIsJunk returns true iff the item with the given name can never be
// progression, regardless of context.
func itemIsJunk(name string) bool {
	switch name {
	case "fist ring", "expert's ring", "energy ring", "toss ring",
		"swimmer's ring":
		return false
	}

	// non-default junk rings
	if rom.Treasures[name] == nil {
		return true
	}

	switch rom.Treasures[name].ID() {
	// heart refill, PoH, HC, ring, compass, dungeon map, gasha seed
	case 0x29, 0x2a, 0x2b, 0x2d, 0x32, 0x33, 0x34:
		return true
	}
	return false
}

// setROMData mutates the ROM data in-place based on the given route.
func setROMData(romData []byte, game int, ri *RouteInfo,
	owlHints map[string]string, logf logFunc, verbose bool) ([]byte, error) {
	// place selected treasures in slots
	checks := getChecks(ri)
	for slot, item := range checks {
		if verbose {
			logf("%s <- %s", slot.Name, item.Name)
		}

		romItemName := item.Name
		if ringName, ok := reverseLookup(ri.RingMap, item.Name); ok {
			romItemName = ringName
		}
		rom.ItemSlots[slot.Name].Treasure = rom.Treasures[romItemName]
	}

	// set season data
	if game == rom.GameSeasons {
		for area, id := range ri.Seasons {
			rom.Seasons[fmt.Sprintf("%s season", area)].New = []byte{id}
		}
	}

	rom.SetAnimal(ri.Companion)
	rom.SetOwlData(owlHints, game)

	// do it! (but don't write anything)
	return rom.Mutate(romData, game)
}

// reverseLookup looks up the key for a given map value. Note that this is only
// "safe" if each value has only one key.
func reverseLookup(m map[string]string, v string) (string, bool) {
	for k, v2 := range m {
		if v2 == v {
			return k, true
		}
	}
	return "", false
}
