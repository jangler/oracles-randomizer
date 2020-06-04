package randomizer

import (
	"crypto/sha1"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"runtime/pprof"
	"sort"
	"strconv"
	"strings"
	"time"
)

type logFunc func(string, ...interface{})

var keyRegexp = regexp.MustCompile("(slate|(small|boss) key)$")

const (
	gameNil = iota
	gameAges
	gameSeasons
)

var gameNames = map[int]string{
	gameNil:     "nil",
	gameAges:    "ages",
	gameSeasons: "seasons",
}

// usage is called when an invalid CLI invocation is used, or if the -h flag is
// passed.
func usage() {
	fmt.Fprintf(flag.CommandLine.Output(),
		"Usage: %s [<original file> [<new file>]]\n", os.Args[0])
	flag.PrintDefaults()
}

// fatal prints an error to whichever UI is used. this doesn't exit the
// program, since that would destroy the TUI.
func fatal(err error, logf logFunc) {
	logf("fatal: %v.", err)
}

// a quick and dirty type of logFunc.
func printErrf(s string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", a...)
}

// options specified on the command line or via the TUI
var (
	flagCpuProf  string
	flagDevCmd   string
	flagDungeons bool
	flagHard     bool
	flagIncludes string
	flagNoUI     bool
	flagPlan     string
	flagMulti    string
	flagPortals  bool
	flagSeed     string
	flagRace     bool
	flagTreewarp bool
	flagVerbose  bool
)

type randomizerOptions struct {
	treewarp bool
	hard     bool
	dungeons bool
	portals  bool
	plan     *plan
	race     bool
	seed     string
	include  []string
	game     int
	players  int
}

// initFlags initializes the CLI/TUI option values and variables.
func initFlags() {
	flag.Usage = usage
	flag.StringVar(&flagCpuProf, "cpuprofile", "",
		"write CPU profile to file")
	flag.StringVar(&flagDevCmd, "devcmd", "",
		"subcommands are 'findaddr', 'showasm', 'stats', and 'hardstats'")
	flag.BoolVar(&flagDungeons, "dungeons", false,
		"shuffle dungeon entrances")
	flag.BoolVar(&flagHard, "hard", false,
		"enable more difficult logic")
	flag.StringVar(&flagIncludes, "include", "",
		"comma-separated list of additional asm files to include")
	flag.BoolVar(&flagNoUI, "noui", false,
		"use command line without prompts if input file is given")
	flag.StringVar(&flagPlan, "plan", "",
		"use fixed 'randomization' from a file")
	flag.StringVar(&flagMulti, "multi", "",
		"comma-separated list of strings such as s+hdp or a+ht")
	flag.BoolVar(&flagPortals, "portals", false,
		"shuffle subrosia portal connections (seasons)")
	flag.BoolVar(&flagRace, "race", false,
		"don't print full seed in file select screen or filename")
	flag.StringVar(&flagSeed, "seed", "",
		"specific random seed to use (32-bit hex number)")
	flag.BoolVar(&flagTreewarp, "treewarp", false,
		"warp to ember tree by pressing start+B on map screen")
	flag.BoolVar(&flagVerbose, "verbose", false,
		"print more detailed output to terminal")
	flag.Parse()
}

// parses options from a string like "s+dp" or "ages+hk" in a ropts.
func roptsFromString(s string, ropts *randomizerOptions) error {
	a := strings.Split(s, "+")
	if len(a) == 0 || len(a) > 2 {
		return fmt.Errorf("bad option string: %s", s)
	}

	// game name
	switch a[0] {
	case "s", "seasons":
		ropts.game = gameSeasons
	case "a", "ages":
		ropts.game = gameAges
	default:
		return fmt.Errorf("unknown game: %s", a[0])
	}

	// flags
	if len(a) == 2 {
		for _, c := range a[1] {
			switch c {
			case 'd':
				ropts.dungeons = true
			case 'h':
				ropts.hard = true
			case 'p':
				ropts.portals = true
			case 't':
				ropts.treewarp = true
			default:
				return fmt.Errorf("unknown flag: %v", c)
			}
		}
	}

	return nil
}

// the program's entry point.
func Main() {
	initFlags()

	if flagCpuProf != "" {
		f, err := os.Create(flagCpuProf)
		if err != nil {
			fatal(err, printErrf)
			return
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// get options
	optsList := make([]*randomizerOptions, 0, 1)
	include := strings.Split(flagIncludes, ",")
	if flagMulti != "" {
		for i, s := range strings.Split(flagMulti, ",") {
			optsList = append(optsList, &randomizerOptions{
				race:    flagRace,
				seed:    flagSeed,
				include: include,
			})
			if err := roptsFromString(s, optsList[i]); err != nil {
				fatal(err, printErrf)
				return
			}
		}
	} else {
		optsList = append(optsList, &randomizerOptions{
			race:     flagRace,
			seed:     flagSeed,
			treewarp: flagTreewarp,
			hard:     flagHard,
			dungeons: flagDungeons,
			portals:  flagPortals,
			include:  include,
		})
	}
	for _, ropts := range optsList {
		ropts.players = len(optsList)
	}

	switch flagDevCmd {
	case "findaddr":
		// print the name of the mutable/etc that modifies an address
		tokens := strings.Split(flag.Arg(0), "/")
		if len(tokens) != 3 {
			fatal(fmt.Errorf("findaddr: invalid argument: %s", flag.Arg(0)),
				printErrf)
			return
		}
		game := reverseLookupOrPanic(gameNames, tokens[0]).(int)
		bank, err := strconv.ParseUint(tokens[1], 16, 8)
		if err != nil {
			fatal(err, printErrf)
			return
		}
		addr, err := strconv.ParseUint(tokens[2], 16, 16)
		if err != nil {
			fatal(err, printErrf)
			return
		}

		// optionally specify path of rom to load.
		// i forget why or whether this is useful.
		var rom *romState
		if flag.Arg(1) == "" {
			rom = newRomState(nil, game, 1, optsList[0].include)
		} else {
			f, err := os.Open(flag.Arg(1))
			if err != nil {
				fatal(err, printErrf)
				return
			}
			defer f.Close()
			b, err := ioutil.ReadAll(f)
			if err != nil {
				fatal(err, printErrf)
				return
			}
			rom = newRomState(b, game, 1, optsList[0].include)
		}

		fmt.Println(rom.findAddr(byte(bank), uint16(addr)))
	case "stats", "hardstats":
		// do stats instead of randomizing
		game := reverseLookupOrPanic(gameNames, flag.Arg(0)).(int)
		numTrials, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fatal(err, printErrf)
			return
		}

		rand.Seed(time.Now().UnixNano())

		statFunc := logStats
		if flagDevCmd == "hardstats" {
			statFunc = logHardStats
		}
		statFunc(game, numTrials, *optsList[0],
			func(s string, a ...interface{}) {
				fmt.Printf(s, a...)
				fmt.Println()
			})
	case "showasm":
		// print the asm for the named function/etc
		tokens := strings.Split(flag.Arg(0), "/")
		if len(tokens) != 2 {
			fatal(fmt.Errorf("showasm: invalid argument: %s", flag.Arg(0)),
				printErrf)
			return
		}
		game := reverseLookupOrPanic(gameNames, tokens[0]).(int)

		rom := newRomState(nil, game, 1, optsList[0].include)
		if err := rom.showAsm(tokens[1], os.Stdout); err != nil {
			fatal(err, printErrf)
			return
		}
	case "":
		// no devcmd, run randomizer normally
		if flagMulti != "" ||
			(flag.NArg() > 0 && flag.NArg()+flag.NFlag() > 1) { // CLI used
			// run randomizer on main goroutine
			runRandomizer(nil, optsList, func(s string, a ...interface{}) {
				fmt.Printf(s, a...)
				fmt.Println()
			})
		} else { // CLI maybe not used
			// run TUI on main goroutine and randomizer on alternate goroutine
			ui := newUI("oracles randomizer " + version)
			go runRandomizer(ui, optsList, func(s string, a ...interface{}) {
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
func runRandomizer(ui *uiInstance, optsList []*randomizerOptions, logf logFunc) {
	// close TUI after randomizer is done
	defer func() {
		if ui != nil {
			ui.done()
		}
	}()

	// if rom is to be randomized, infile must be non-empty after switch
	dirName, infiles, outfiles := getRomPaths(ui, optsList, logf)
	if infiles != nil {
		roms := make([]*romState, len(infiles))
		routes := make([]*routeInfo, len(infiles))

		if ui != nil {
			if ui.doPrompt("use specific seed? (y/n)") == 'y' {
				optsList[0].seed =
					ui.promptSeed("enter seed: (8-digit hex number)")
				logf("using seed %s.", optsList[0].seed)
			}
		}
		seed, err := setRandomSeed(optsList[0].seed)
		if err != nil {
			fatal(err, logf)
			return
		}
		src := rand.New(rand.NewSource(int64(seed)))

		// get input for instance
		for i, infile := range infiles {
			ropts := optsList[i]

			b, game, err := readGivenRom(filepath.Join(dirName, infile))
			if err != nil {
				fatal(err, logf)
				return
			} else {
				roms[i] = newRomState(b, game, i+1, ropts.include)
			}

			// sanity check beforehand
			if errs := roms[i].verify(); errs != nil {
				if flagVerbose {
					for _, err := range errs {
						logf(err.Error())
					}
				}
				fatal(errs[0], logf)
				return
			}

			logf("randomizing %s.", infile)
			getAndLogOptions(game, ui, ropts, logf)
			if ui != nil {
				logf("")
			}

			roms[i].setTreewarp(ropts.treewarp)

			if flagPlan != "" {
				var err error
				ropts.plan, err = parseSummary(flagPlan, game)
				if err != nil {
					fatal(err, logf)
					return
				}
			}

			// find routes
			if ropts.plan == nil {
				route, err := findRoute(
					roms[i], seed, src, *ropts, flagVerbose, logf)
				if err != nil {
					fatal(err, logf)
					return
				}
				routes[i] = route
			} else {
				route, err := makePlannedRoute(roms[i], ropts.plan)
				if err != nil {
					fatal(err, logf)
					return
				}
				routes[i] = route
				ropts.dungeons = route.entrances != nil && len(route.entrances) > 0
				ropts.portals = route.portals != nil && len(route.portals) > 0
			}
		}

		if len(routes) > 1 {
			shuffleMultiworld(routes, roms, flagVerbose, logf)
		}

		// come up with log data
		g, checks, spheres, extra := getAllSpheres(routes)
		resetFunc := func() {
			for _, ri := range routes {
				ri.graph.reset()
			}
		}
		if flagVerbose {
			logf("%d checks", len(checks))
			logf("%d spheres", len(spheres))
		}

		// accumulate all treasures for reference by log functions
		treasures := make(map[string]*treasure)
		for _, rom := range roms {
			for k, v := range rom.treasures {
				treasures[k] = v
			}
		}

		// write roms
		for i, rom := range roms {
			ropts := optsList[i]

			gamePrefix := sora(rom.game, "oos", "ooa")
			var outfile string
			if outfiles != nil && len(outfiles) > i {
				outfile = outfiles[i]
			} else if len(roms) == 1 {
				outfile = fmt.Sprintf("%srando_%s_%s.gbc", gamePrefix, version,
					optString(seed, ropts, "-"))
			} else {
				outfile = fmt.Sprintf("%srando_%s_%s_p%d.gbc", gamePrefix, version,
					optString(seed, ropts, "-"), i+1)
			}
			logFilename := strings.Replace(outfile, ".gbc", "", 1) + "_log.txt"

			sum, err := applyRoute(rom, routes[i], dirName, logFilename, ropts,
				checks, spheres, extra, g, resetFunc, treasures, flagVerbose, logf)
			if err != nil {
				fatal(err, logf)
				return
			}

			if writeRom(rom.data, dirName, outfile, logFilename, seed, sum, logf); err != nil {
				fatal(err, logf)
				return
			}
		}

		for _, ri := range routes {
			ri.graph["start"].removeParent(g["start"])
			g["done"].removeParent(ri.graph["done"])
		}
	}
}

// returns the target directory and filenames of input and output files. the
// output filename may be empty, in which case it will be automatically
// determined.
func getRomPaths(ui *uiInstance, optsList []*randomizerOptions,
	logf logFunc) (dir string, in, out []string) {
	switch flag.NArg() {
	case 0: // no specified files, search in executable's directory
		var seasons, ages string
		var err error
		dir, seasons, ages, err = findVanillaRoms(ui, logf)
		if err != nil {
			fatal(err, logf)
			break
		}

		// print which files, if any, are found.
		if seasons != "" {
			if ui != nil {
				ui.printPath("found vanilla US seasons ROM: ", seasons, "")
			} else {
				logf("found vanilla US seasons ROM: %s", seasons)
			}
		} else {
			logf("no vanilla US seasons ROM found.")
		}
		if ages != "" {
			if ui != nil {
				ui.printPath("found vanilla US ages ROM: ", ages, "")
			} else {
				logf("found vanilla US ages ROM: %s", ages)
			}
		} else {
			logf("no vanilla US ages ROM found.")
		}
		if ui != nil {
			logf("")
		}

		// determine which filename to use based on what roms are found, and on
		// user input.
		in = make([]string, len(optsList))
		for i, ropts := range optsList {
			if seasons == "" && ages == "" {
				logf("no ROMs found in program's directory, " +
					"and no ROMs specified.")
				in = nil
				break
			} else if ropts.game != 0 {
				in[i] = ternary(ropts.game == gameSeasons, seasons, ages).(string)
				if in[i] == "" {
					logf("ROM for game not found")
					in = nil
					break
				}
			} else if seasons != "" && ages != "" {
				which := ui.doPrompt("randomize (s)easons or (a)ges?")
				in[i] = ternary(which == 's', seasons, ages).(string)
			} else if seasons != "" {
				in[i] = seasons
			} else {
				in[i] = ages
			}
		}
	case 1: // specified input file only
		in = strings.Split(flag.Arg(0), ",")
	case 2: // specified input and output file
		in = strings.Split(flag.Arg(0), ",")
		out = strings.Split(flag.Arg(1), ",")
	default:
		flag.Usage()
	}

	return dir, in, out
}

// getAndLogOptions logs values of selected options, prompting for them first
// if the TUI is used.
func getAndLogOptions(game int, ui *uiInstance, ropts *randomizerOptions,
	logf logFunc) {
	if ui != nil {
		ropts.hard = ui.doPrompt("enable hard difficulty? (y/n)") == 'y'
	}
	logf("using %s difficulty.", ternary(ropts.hard, "hard", "normal"))

	if ui != nil {
		ropts.treewarp = ui.doPrompt("enable tree warp? (y/n)") == 'y'
	}
	logf("tree warp %s.", ternary(ropts.treewarp, "on", "off"))

	if ui != nil {
		ropts.dungeons = ui.doPrompt("shuffle dungeons? (y/n)") == 'y'
	}
	logf("dungeon shuffle %s.", ternary(ropts.dungeons, "on", "off"))

	if game == gameSeasons {
		if ui != nil {
			ropts.portals = ui.doPrompt("shuffle portals? (y/n)") == 'y'
		}
		logf("portal shuffle %s.", ternary(ropts.portals, "on", "off"))
	}
}

// attempt to write rom data to a file and print summary info.
func writeRom(b []byte, dirName, filename, logFilename string, seed uint32,
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
	if flagPlan == "" && !flagRace {
		logf("seed: %08x", seed)
	}
	logf("SHA-1 sum: %x", string(sum))
	logf("wrote new ROM to %s", filename)
	if flagPlan == "" && !flagRace {
		logf("wrote log file to %s", logFilename)
	}

	return nil
}

// search for a vanilla US seasons and ages roms in the executable's directory,
// and return their filenames.
func findVanillaRoms(
	ui *uiInstance, logf logFunc) (dirName, seasons, ages string, err error) {
	// read slice of file info from executable's dir
	exe, err := os.Executable()
	if err != nil {
		return
	}
	dirName = filepath.Dir(exe)
	if ui != nil {
		ui.printPath("searching ", dirName, " for ROMs.")
	} else {
		logf("searching %s for ROMs.", dirName)
	}
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
		if !romIsJp(b) && romIsVanilla(b) {
			if romIsAges(b) {
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
func readGivenRom(filename string) ([]byte, int, error) {
	// read file
	f, err := os.Open(filename)
	if err != nil {
		return nil, gameNil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, gameNil, err
	}

	// check file data
	if !romIsAges(b) && !romIsSeasons(b) {
		return nil, gameNil,
			fmt.Errorf("%s is not an oracles ROM", filename)
	}
	if romIsJp(b) {
		return nil, gameNil,
			fmt.Errorf("%s is a JP ROM; only US is supported", filename)
	}
	if !romIsVanilla(b) {
		return nil, gameNil,
			fmt.Errorf("%s is an unrecognized oracles ROM", filename)
	}

	game := ternary(romIsSeasons(b), gameSeasons, gameAges).(int)
	return b, game, nil
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
func applyRoute(rom *romState, ri *routeInfo, dirName, logFilename string,
	ropts *randomizerOptions, checks map[*node]*node, spheres [][]*node,
	extra []*node, g graph, resetFunc func(), treasures map[string]*treasure,
	verbose bool, logf logFunc) ([]byte, error) {
	checksum, err := setRomData(rom, ri, ropts, logf, verbose)
	if err != nil {
		return nil, err
	}

	// write spoiler log
	if ropts.plan == nil && !ropts.race {
		writeSummary(filepath.Join(dirName, logFilename), checksum, *ropts,
			rom, ri, checks, spheres, extra, g, resetFunc, treasures, nil)
	}

	return checksum, nil
}

// mutates the rom data in-place based on the given route. this doesn't write
// the file.
func setRomData(rom *romState, ri *routeInfo, ropts *randomizerOptions,
	logf logFunc, verbose bool) ([]byte, error) {
	// place selected treasures in slots
	checks := getChecks(ri.usedItems, ri.usedSlots)
	for slot, item := range checks {
		if verbose {
			logf("%s <- %s", slot.name, item.name)
		}

		romItemName := item.name
		if ringName, ok := reverseLookup(ri.ringMap, item.name); ok {
			romItemName = ringName.(string)
		}
		rom.itemSlots[slot.name].treasure = rom.treasures[romItemName]
	}

	// set season data
	if rom.game == gameSeasons {
		for area, id := range ri.seasons {
			rom.setSeason(inflictCamelCase(area+"Season"), id)
		}
	}

	rom.setAnimal(ri.companion)

	warps := make(map[string]string)
	if ropts.dungeons {
		for k, v := range ri.entrances {
			warps[k] = v
		}
	}
	if ropts.portals {
		for k, v := range ri.portals {
			holodrumV, _ := reverseLookup(subrosianPortalNames, v)
			warps[fmt.Sprintf("%s portal", k)] =
				fmt.Sprintf("%s portal", holodrumV)
		}
	}

	// do it! (but don't write anything)
	return rom.mutate(warps, ri.seed, ropts)
}

// returns a string representing a seed/has plus the randomizer options that
// affect the generated seed or how it's played - so not including things like
// music on/off.
func optString(seed uint32, ropts *randomizerOptions, flagSep string) string {
	s := ""

	if ropts.plan != nil {
		// -plan gets a hash based on source file rather than a seed
		sum := sha1.Sum([]byte(ropts.plan.source))
		s += fmt.Sprintf("plan-%03x", ((int(sum[0])<<8)+int(sum[1]))>>4)

		// treewarp is the only option that makes a difference in plando
		if ropts.treewarp {
			s += flagSep + "t"
		}

		return s
	}

	if ropts.race {
		s += fmt.Sprintf("race-%03x", seed>>20)
	} else {
		s += fmt.Sprintf("%08x", seed)
	}

	if ropts.treewarp || ropts.hard || ropts.dungeons || ropts.portals {
		// these are in chronological order of introduction, for no particular
		// reason.
		s += flagSep
		if ropts.treewarp {
			s += "t"
		}
		if ropts.hard {
			s += "h"
		}
		if ropts.dungeons {
			s += "d"
		}
		if ropts.portals {
			s += "p"
		}
	}

	return s
}

// reverseLookup looks up the key for a given map value. If multiple keys are
// associated with the same value, it will return one of those keys at random.
func reverseLookup(m, match interface{}) (interface{}, bool) {
	iter := reflect.ValueOf(m).MapRange()
	for iter.Next() {
		k, v := iter.Key(), iter.Value()
		if reflect.DeepEqual(v.Interface(), match) {
			return k.Interface(), true
		}
	}
	return nil, false
}

// guess what this does.
func reverseLookupOrPanic(m, match interface{}) interface{} {
	i, ok := reverseLookup(m, match)
	if !ok {
		panic(fmt.Sprintf("reverse lookup failed for value %v", match))
	}
	return i
}

// returns a sorted slice of string keys from a map.
func orderedKeys(m interface{}) []string {
	v := reflect.ValueOf(m)
	a := make([]string, v.Len())
	for i, key := range v.MapKeys() {
		a[i] = key.String()
	}
	sort.Strings(a)
	return a
}

// sora = Seasons OR Ages: returns the first value if the game is seasons, and
// the second if the game is ages. panics if the game is neither.
func sora(game int, sOption, aOption interface{}) interface{} {
	switch game {
	case gameSeasons:
		return sOption
	case gameAges:
		return aOption
	}
	panic("invalid game provided to sora()")
}

// equivalent to the ternary operation (a ? b : c) in C, etc.
func ternary(expr bool, trueOpt, falseOpt interface{}) interface{} {
	if expr {
		return trueOpt
	}
	return falseOpt
}
