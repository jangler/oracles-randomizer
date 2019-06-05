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
	flagNoMusic  bool
	flagNoUI     bool
	flagPlan     string
	flagPortals  bool
	flagSeed     string
	flagTreewarp bool
	flagVerbose  bool
)

type randomizerOptions struct {
	treewarp bool
	hard     bool
	dungeons bool
	portals  bool
	plan     *plan
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

	ropts := randomizerOptions{
		treewarp: flagTreewarp,
		hard:     flagHard,
		dungeons: flagDungeons,
		portals:  flagPortals,
		seed:     flagSeed,
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
			rom = newRomState(nil, game)
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
			rom = newRomState(b, game)
		}

		fmt.Println(rom.findAddr(byte(bank), uint16(addr)))
	case "stats":
		// do stats instead of randomizing
		game := reverseLookupOrPanic(gameNames, flag.Arg(0)).(int)
		numTrials, err := strconv.Atoi(flag.Arg(1))
		if err != nil {
			fatal(err, printErrf)
			return
		}

		rand.Seed(time.Now().UnixNano())
		logStats(game, numTrials, ropts, func(s string, a ...interface{}) {
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

		rom := newRomState(nil, game)
		if err := rom.showAsm(tokens[1], os.Stdout); err != nil {
			fatal(err, printErrf)
			return
		}
	case "":
		// no devcmd, run randomizer normally
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
	dirName, infile, outfile := getRomPaths(ui, logf)
	if infile != "" {
		var rom *romState
		b, game, err := readGivenRom(filepath.Join(dirName, infile))
		if err != nil {
			fatal(err, logf)
			return
		} else {
			rom = newRomState(b, game)
		}

		logf("randomizing %s.", infile)
		getAndLogOptions(ui, logf)
		if ui != nil {
			logf("")
		}

		rom.setMusic(!flagNoMusic)
		rom.setTreewarp(flagTreewarp)

		if flagPlan != "" {
			var err error
			ropts.plan, err = parseSummary(flagPlan, game)
			if err != nil {
				fatal(err, logf)
				return
			}
		}

		if err := randomizeFile(
			rom, dirName, outfile, ropts, flagVerbose, logf); err != nil {
			fatal(err, logf)
			return
		}
	}
}

// returns the target directory and filenames of input and output files. the
// output filename may be empty, in which case it will be automatically
// determined.
func getRomPaths(ui *uiInstance, logf logFunc) (dir, in, out string) {
	switch flag.NArg() {
	case 0: // no specified files, search in executable's directory
		var seasons, ages string
		var err error
		dir, seasons, ages, err = findVanillaRoms(ui)
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
			in = ternary(which == 's', seasons, ages).(string)
		} else if seasons != "" {
			in = seasons
		} else {
			in = ages
		}
	case 1: // specified input file only
		in = flag.Arg(0)
	case 2: // specified input and output file
		in, out = flag.Arg(0), flag.Arg(1)
	default:
		flag.Usage()
	}

	return dir, in, out
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
	logf("using %s difficulty.", ternary(flagHard, "hard", "normal"))

	if ui != nil {
		flagNoMusic = ui.doPrompt("disable music? (y/n)") == 'y'
	}
	logf("music %s.", ternary(flagNoMusic, "off", "on"))

	if ui != nil {
		flagTreewarp = ui.doPrompt("enable tree warp? (y/n)") == 'y'
	}
	logf("tree warp %s.", ternary(flagTreewarp, "on", "off"))
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
	if flagPlan == "" {
		logf("seed: %08x", seed)
	}
	logf("SHA-1 sum: %x", string(sum))
	logf("wrote new ROM to %s", filename)
	if flagPlan == "" {
		logf("wrote log file to %s", logFilename)
	}

	return nil
}

// search for a vanilla US seasons and ages roms in the executable's directory,
// and return their filenames.
func findVanillaRoms(
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

// finds a valid seed/configuration and writes it to the output file.
func randomizeFile(rom *romState, dirName, outfile string,
	ropts randomizerOptions, verbose bool, logf logFunc) error {
	var seed uint32
	var sum []byte
	var err error
	var logFilename string

	if ropts.portals && rom.game == gameAges {
		return fmt.Errorf("portal randomization does not apply to ages")
	}

	// operate on rom data
	if outfile != "" {
		logFilename = outfile[:len(outfile)-4] + "_log.txt"
	}
	seed, sum, logFilename, err = randomize(
		rom, dirName, logFilename, ropts, verbose, logf)
	if err != nil {
		return err
	}
	if outfile == "" {
		gamePrefix := sora(rom.game, "oos", "ooa")
		outfile = fmt.Sprintf("%srando_%s_%s.gbc",
			gamePrefix, version, optString(seed, ropts))
	}

	// write to file
	return writeRom(rom.data, dirName, outfile, logFilename, seed, sum, logf)
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
func randomize(rom *romState, dirName, logFilename string,
	ropts randomizerOptions, verbose bool,
	logf logFunc) (uint32, []byte, string, error) {
	// sanity check beforehand
	if errs := rom.verify(); errs != nil {
		if verbose {
			for _, err := range errs {
				logf(err.Error())
			}
		}
		return 0, nil, "", errs[0]
	}

	// search for valid configuration
	var ri *routeInfo
	if ropts.plan == nil {
		seed, err := setRandomSeed(ropts.seed)
		if err != nil {
			return 0, nil, "", err
		}
		ri, err = findRoute(rom, seed, ropts, verbose, logf)
		if err != nil {
			return 0, nil, "", err
		}
	} else {
		var err error
		ri, err = makePlannedRoute(rom, ropts.plan)
		if err != nil {
			return 0, nil, "", err
		}
		if ri.entrances != nil && len(ri.entrances) > 0 {
			ropts.dungeons = true
		}
		if ri.portals != nil && len(ri.portals) > 0 {
			ropts.portals = true
		}
	}

	// configuration found; come up with auxiliary data
	checks := getChecks(ri.usedItems, ri.usedSlots)
	spheres, extra := getSpheres(ri.graph, checks)
	owlNames := orderedKeys(getOwlIds(rom.game))
	owlHinter := newHinter(rom.game)
	owlHints := owlHinter.generate(ri.src, ri.graph, checks, owlNames)
	if ropts.plan != nil {
		if err := planOwlHints(ropts.plan, owlHinter, owlHints); err != nil {
			return 0, nil, "", err
		}
	}

	checksum, err := setRomData(rom, ri, owlHints, ropts, logf, verbose)
	if err != nil {
		return 0, nil, "", err
	}

	// write spoiler log
	if ropts.plan == nil {
		if logFilename == "" {
			gamePrefix := sora(rom.game, "oos", "ooa")
			logFilename = fmt.Sprintf("%srando_%s_%s_log.txt",
				gamePrefix, version, optString(ri.seed, ropts))
		}
		writeSummary(filepath.Join(dirName, logFilename), checksum,
			ropts, rom, ri, checks, spheres, extra, owlHints)
	}

	return ri.seed, checksum, logFilename, nil
}

// mutates the rom data in-place based on the given route. this doesn't write
// the file.
func setRomData(rom *romState, ri *routeInfo, owlHints map[string]string,
	ropts randomizerOptions, logf logFunc, verbose bool) ([]byte, error) {
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
	rom.setOwlData(owlHints)

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
func optString(seed uint32, ropts randomizerOptions) string {
	s := ""

	if ropts.plan != nil {
		// -plan gets a hash based on source file rather than a seed
		sum := sha1.Sum([]byte(ropts.plan.source))
		s += fmt.Sprintf("plan-%04x", sum[:2])

		// treewarp is the only option that makes a difference in plando
		if ropts.treewarp {
			s += "+t"
		}

		return s
	}

	s += fmt.Sprintf("%08x", seed)

	if ropts.treewarp || ropts.hard || ropts.dungeons || ropts.portals {
		// these are in chronological order of introduction, for no particular
		// reason.
		s += "+"
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
