package main

//go:generate bash scripts/generate.sh
//go:generate esc -o embed.go hints/ logic/
//go:generate esc -o rom/embed.go -pkg rom asm/ hints/ romdata/ lgbtasm/lgbtasm.lua

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jangler/oracles-randomizer/rom"
)

type logFunc func(string, ...interface{})

var keyRegexp = regexp.MustCompile("(slate|(small|boss) key)$")

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
	flagCpuProf  string
	flagDevCmd   string
	flagDungeons bool
	flagHard     bool
	flagNoMusic  bool
	flagNoUI     bool
	flagPlan     string
	flagPortals  bool
	flagSeed     string
	flagTreewarp bool
	flagVerbose  bool
)

type randomizerOptions struct {
	dungeons bool
	hard     bool
	plan     *summary
	portals  bool
	seed     string // given seed, not necessarily final seed
}

// initFlags initializes the CLI/TUI option values and variables.
func initFlags() {
	flag.Usage = usage
	flag.StringVar(&flagCpuProf, "cpuprofile", "",
		"write CPU profile to file")
	flag.StringVar(&flagDevCmd, "devcmd", "",
		"subcommands are 'findaddr', 'showasm', and 'stats'")
	flag.BoolVar(&flagDungeons, "dungeons", false,
		"shuffle dungeon entrances")
	flag.BoolVar(&flagHard, "hard", false,
		"enable more difficult logic")
	flag.BoolVar(&flagNoMusic, "nomusic", false,
		"don't play any music in the modified ROM")
	flag.BoolVar(&flagNoUI, "noui", false,
		"use command line output without option prompts")
	flag.StringVar(&flagPlan, "plan", "",
		"use fixed 'randomization' from a file")
	flag.BoolVar(&flagPortals, "portals", false,
		"shuffle subrosia portal connections (seasons)")
	flag.StringVar(&flagSeed, "seed", "",
		"specific random seed to use (32-bit hex number)")
	flag.BoolVar(&flagTreewarp, "treewarp", false,
		"warp to ember tree by pressing start+B on map screen")
	flag.BoolVar(&flagVerbose, "verbose", false,
		"print more detailed output to terminal")
	flag.Parse()
}

// main is the program's entry point.
func main() {
	initFlags()

	if flagCpuProf != "" {
		f, err := os.Create(flagCpuProf)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	ropts := randomizerOptions{
		dungeons: flagDungeons,
		hard:     flagHard,
		portals:  flagPortals,
		seed:     flagSeed,
	}

	switch flagDevCmd {
	case "findaddr":
		// print the name of the mutable/etc that modifies an address
		var game int

		tokens := strings.Split(flag.Arg(0), ":")
		if len(tokens) != 3 {
			panic("findaddr: invalid argument: " + flag.Arg(0))
		}

		switch tokens[0] {
		case "seasons":
			game = rom.GameSeasons
		case "ages":
			game = rom.GameAges
		default:
			panic("invalid game name: " + tokens[0])
		}

		bank, err := strconv.ParseUint(tokens[1], 16, 8)
		if err != nil {
			panic(err)
		}
		addr, err := strconv.ParseUint(tokens[2], 16, 16)
		if err != nil {
			panic(err)
		}

		// optionall specify path of rom to load
		if flag.Arg(1) == "" {
			rom.Init(nil, game)
		} else {
			f, err := os.Open(flag.Arg(1))
			if err != nil {
				panic(err)
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}
			rom.Init(b, game)
		}

		fmt.Println(rom.FindAddr(byte(bank), uint16(addr)))
	case "stats":
		// do stats instead of randomizing
		var game int

		switch flag.Arg(0) {
		case "seasons":
			game = rom.GameSeasons
		case "ages":
			game = rom.GameAges
		default:
			panic("invalid game: " + flag.Arg(0))
		}

		numTrials, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			panic(err)
		}

		rom.Init(nil, game)
		rand.Seed(time.Now().UnixNano())
		logStats(game, numTrials, ropts, func(s string, a ...interface{}) {
			fmt.Printf(s, a...)
			fmt.Println()
		})
	case "showasm":
		// print the asm for the named function/etc
		var game int

		tokens := strings.Split(flag.Arg(0), ":")
		if len(tokens) != 2 {
			panic("showasm: invalid argument: " + flag.Arg(0))
		}

		switch tokens[0] {
		case "seasons":
			game = rom.GameSeasons
		case "ages":
			game = rom.GameAges
		default:
			panic("invalid game name: " + tokens[0])
		}

		rom.Init(nil, game)
		if err := rom.ShowAsm(tokens[1], os.Stdout); err != nil {
			panic(err)
		}
	case "":
		if flag.NArg()+flag.NFlag() > 1 { // CLI used
			// run randomizer on main goroutine
			runRandomizer(nil, ropts, func(s string, a ...interface{}) {
				fmt.Printf(s, a...)
				fmt.Println()
			})
		} else { // CLI maybe not used
			// run TUI on main goroutine and randomizer on alternate goroutine
			ui := newUI("oracles randomizer " + version)
			go runRandomizer(ui, ropts, func(s string, a ...interface{}) {
				ui.printf(s, a...)
			})
			ui.run()
		}
	default:
		fmt.Printf("invalid dev command: %s\n", flagDevCmd)
	}
}

// run the main randomizer routine, printing messages via logf, which should
// act analogously to fmt.Printf with added newline.
func runRandomizer(ui *uiInstance, ropts randomizerOptions, logf logFunc) {
	// close TUI after randomizer is done
	defer func() {
		if ui != nil {
			ui.done()
		}
	}()

	// if rom is to be randomized, infile must be non-empty after switch
	var dirName, infile, outfile string
	switch flag.NArg() {
	case 0: // no specified files, search in executable's directory
		var seasons, ages string
		var err error
		dirName, seasons, ages, err = findVanillaROMs(ui)
		if err != nil {
			fatal(err, logf)
			break
		}

		// print which files, if any, are found.
		if seasons != "" {
			ui.printPath("found vanilla US seasons ROM: ", seasons, "")
		} else {
			ui.printf("no vanilla US seasons ROM found.")
		}
		if ages != "" {
			ui.printPath("found vanilla US ages ROM: ", ages, "")
		} else {
			ui.printf("no vanilla US ages ROM found.")
		}
		ui.printf("")

		// determine which filename to use based on what roms are found, and on
		// user input.
		if seasons == "" && ages == "" {
			ui.printf("no ROMs found in program's directory, " +
				"and no ROMs specified.")
		} else if seasons != "" && ages != "" {
			which := ui.doPrompt("randomize (s)easons or (a)ges?")
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
			rom.Init(b, game)
		}
		logf("randomizing %s.", infile)

		getAndLogOptions(ui, logf)

		if ui != nil {
			logf("")
		}

		rom.SetMusic(!flagNoMusic)
		rom.SetTreewarp(flagTreewarp)

		if flagPlan != "" {
			var err error
			ropts.plan, err = parseSummary(flagPlan, game)
			if err != nil {
				fatal(err, logf)
				return
			}
			ropts.dungeons = ropts.dungeons || len(ropts.plan.dungeons) > 0
			ropts.portals = ropts.portals || len(ropts.plan.portals) > 0
		} else {
			ropts.plan = newSummary()
		}

		if err := randomizeFile(
			b, game, dirName, outfile, ropts, flagVerbose, logf); err != nil {
			fatal(err, logf)
			return
		}
	}
}

// getAndLogOptions logs values of selected options, prompting for them first
// if the TUI is used.
func getAndLogOptions(ui *uiInstance, logf logFunc) {
	if ui != nil {
		if ui.doPrompt("use specific seed? (y/n)") == 'y' {
			flagSeed = ui.promptSeed("enter seed: (8-digit hex number)")
			logf("using seed %s.", flagSeed)
		}
	}

	if ui != nil {
		flagHard = ui.doPrompt("enable hard difficulty? (y/n)") == 'y'
	}
	if flagHard {
		logf("using hard difficulty.")
	} else {
		logf("using normal difficulty.")
	}

	if ui != nil {
		flagNoMusic = ui.doPrompt("disable music? (y/n)") == 'y'
	}
	if flagNoMusic {
		logf("music off.")
	} else {
		logf("music on.")
	}

	if ui != nil {
		flagTreewarp = ui.doPrompt("enable tree warp? (y/n)") == 'y'
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
func findVanillaROMs(
	ui *uiInstance) (dirName, seasons, ages string, err error) {
	// read slice of file info from executable's dir
	exe, err := os.Executable()
	if err != nil {
		return
	}

	dirName = filepath.Dir(exe)
	ui.printPath("searching ", dirName, " for ROMs.")
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
		if rom.IsNonJP(b) && rom.IsVanilla(b) {
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
	if !rom.IsNonJP(b) {
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

func randomizeFile(romData []byte, game int, dirName, outfile string,
	ropts randomizerOptions, verbose bool, logf logFunc) error {
	var seed uint32
	var sum []byte
	var err error
	var logFilename string

	if ropts.portals && game == rom.GameAges {
		return fmt.Errorf("portal randomization does not apply to ages")
	}

	// operate on rom data
	if outfile != "" {
		logFilename = outfile[:len(outfile)-4] + "_log.txt"
	}
	seed, sum, logFilename, err = randomize(
		romData, game, dirName, logFilename, ropts, verbose, logf)
	if err != nil {
		return err
	}
	hardString := ""
	if ropts.hard {
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
func randomize(romData []byte, game int, dirName, logFilename string,
	ropts randomizerOptions, verbose bool,
	logf logFunc) (uint32, []byte, string, error) {
	// sanity check beforehand
	if errs := rom.Verify(romData, game); errs != nil {
		if verbose {
			for _, err := range errs {
				logf(err.Error())
			}
		}
		return 0, nil, "", errs[0]
	}

	seed, err := setRandomSeed(ropts.seed)
	if err != nil {
		return 0, nil, "", err
	}

	// search for route
	ri, err := findRoute(game, seed, ropts, verbose, logf)
	if err != nil {
		return 0, nil, "", err
	}

	checks := getChecks(ri)
	spheres, extra := getSpheres(ri.Route.Graph, checks)
	owlHints, err := newHinter(game).generate(ri.Src, ri.Route.Graph, checks,
		rom.GetOwlNames(game), ropts.plan.hints)
	if err != nil {
		return 0, nil, "", err
	}

	checksum, err := setROMData(
		romData, game, ri, owlHints, ropts, logf, verbose)
	if err != nil {
		return 0, nil, "", err
	}

	hardString := ""
	if ropts.hard {
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
	if ropts.hard {
		summary <- fmt.Sprintf("difficulty: hard")
	} else {
		summary <- fmt.Sprintf("difficulty: normal")
	}
	summary <- ""
	summary <- ""
	summary <- "-- progression items --"
	summary <- ""
	nonKeyChecks := make(map[*node]*node)
	for slot, item := range checks {
		if !keyRegexp.MatchString(item.name) {
			nonKeyChecks[slot] = item
		}
	}
	prog, junk := filterJunk(ri.Route.Graph, nonKeyChecks)
	logSpheres(summary, prog, spheres, extra, game, nil)
	summary <- ""
	summary <- "-- small keys and boss keys --"
	summary <- ""
	logSpheres(summary, checks, spheres, extra, game, keyRegexp.MatchString)
	summary <- ""
	summary <- "-- other items --"
	summary <- ""
	logSpheres(summary, junk, spheres, extra, game, nil)
	if ropts.dungeons {
		summary <- ""
		summary <- "-- dungeon entrances --"
		summary <- ""
		sendSorted(summary, func(c chan string) {
			for entrance, dungeon := range ri.Entrances {
				c <- fmt.Sprintf("%s entrance <- %s",
					"D"+entrance[1:], "D"+dungeon[1:])
			}
			close(c)
		})
		summary <- ""
	}
	if ropts.portals {
		summary <- ""
		summary <- "-- subrosia portals --"
		summary <- ""
		sendSorted(summary, func(c chan string) {
			for in, out := range ri.Portals {
				c <- fmt.Sprintf("%-20s <- %s",
					getNiceName(in, game), getNiceName(out, game))
			}
			close(c)
		})
		summary <- ""
	}
	if game == rom.GameSeasons {
		summary <- ""
		summary <- "-- default seasons --"
		summary <- ""
		sendSorted(summary, func(c chan string) {
			for area, id := range ri.Seasons {
				c <- fmt.Sprintf("%-15s <- %s", area, seasonsById[id])
			}
			close(c)
		})
		summary <- ""
	}
	summary <- ""
	summary <- "-- hints --"
	summary <- ""
	sendSorted(summary, func(c chan string) {
		for owlName, hint := range owlHints {
			oneLineHint := strings.ReplaceAll(hint, "\n", " ")
			oneLineHint = strings.ReplaceAll(oneLineHint, "  ", " ")
			c <- fmt.Sprintf("%-20s <- \"%s\"", owlName, oneLineHint)
		}
		close(c)
	})

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

	// not part of next switch since the ID is only junk in seasons
	if name == "treasure map" {
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
	owlHints map[string]string, ropts randomizerOptions, logf logFunc,
	verbose bool) ([]byte, error) {
	// place selected treasures in slots
	checks := getChecks(ri)
	for slot, item := range checks {
		if verbose {
			logf("%s <- %s", slot.name, item.name)
		}

		romItemName := item.name
		if ringName, ok := reverseLookup(ri.RingMap, item.name); ok {
			romItemName = ringName
		}
		rom.ItemSlots[slot.name].Treasure = rom.Treasures[romItemName]
	}

	// set season data
	if game == rom.GameSeasons {
		for area, id := range ri.Seasons {
			// dumb camel case transformation
			key := fmt.Sprintf("%c%sSeason", area[0],
				strings.ReplaceAll(strings.Title(area)[1:], " ", ""))
			rom.SetSeason(key, id)
		}
	}

	rom.SetAnimal(ri.Companion)
	rom.SetOwlData(owlHints, game)

	warps := make(map[string]string)
	if ropts.dungeons {
		for k, v := range ri.Entrances {
			warps[k] = v
		}
	}
	if ropts.portals {
		for k, v := range ri.Portals {
			holodrumV, _ := reverseLookup(subrosianPortalNames, v)
			warps[fmt.Sprintf("%s portal", k)] =
				fmt.Sprintf("%s portal", holodrumV)
		}
	}

	// do it! (but don't write anything)
	return rom.Mutate(romData, game, warps, ropts.dungeons)
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

// get the output of a function that sends strings to a given channel, sort
// those strings, and send them to the `out` channel.
func sendSorted(out chan string, generate func(chan string)) {
	in := make(chan string)
	lines := make([]string, 0, 20) // should be enough capacity for most cases

	go generate(in)
	for s := range in {
		lines = append(lines, s)
	}

	sort.Strings(lines)
	for _, line := range lines {
		out <- line
	}
}
