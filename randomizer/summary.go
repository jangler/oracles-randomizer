package randomizer

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"
)

// returns a channel that will write strings to a text file with CRLF line
// endings. the function will send on the int channel when finished printing.
func getSummaryChannel(filename string) (chan string, chan int) {
	c, done := make(chan string), make(chan int)

	go func() {
		logFile, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		defer logFile.Close()

		for line := range c {
			fmt.Fprintf(logFile, "%s\r\n", line)
		}
		done <- 1
	}()

	// header
	c <- fmt.Sprintf("oracles randomizer %s", version)
	c <- fmt.Sprintf("generated %s", time.Now().Format(time.RFC3339))

	return c, done
}

// separates a map of checks into progression checks and junk checks.
func filterJunk(g graph, checks map[*node]*node,
	treasures map[string]*treasure) (prog, junk map[*node]*node) {
	prog, junk = make(map[*node]*node), make(map[*node]*node)

	// start by assuming every item is progression
	for k, v := range checks {
		prog[k] = v
	}

	done := false
	for !done {
		spheres, _ := getSpheres(g, prog)
		done = true

		// create a copy to pass to functions so that the map we're iterating
		// over isn't modified
		progCopy := make(map[*node]*node)
		for k, v := range prog {
			progCopy[k] = v
		}

		// if item isn't required, move it to junk and reset iteration
		for slot, item := range prog {
			if itemIsInert(treasures, item.name) ||
				(!itemIsRequired(g, slot, item) &&
					!itemChangesProgression(g, progCopy, spheres, slot, item)) {
				delete(prog, slot)
				junk[slot] = item
				done = false
			}
		}
	}

	return
}

// returns true iff removing the slot/item combination from the graph would
// make the seed unbeatable.
func itemIsRequired(g graph, slot, item *node) bool {
	g.reset()
	item.removeParent(slot)
	g["start"].explore()
	item.addParent(slot)
	return !g["done"].reached
}

// returns true iff removing the slot/item combination from the graph would
// change the spheres in which other items appear.
func itemChangesProgression(g graph, checks map[*node]*node, spheres [][]*node,
	slot, item *node) bool {
	oldText := spheresToText(spheres, checks, slot)
	item.removeParent(slot)
	delete(checks, slot)
	newSpheres, _ := getSpheres(g, checks)
	item.addParent(slot)
	checks[slot] = item
	newText := spheresToText(newSpheres, checks, slot)
	return newText != oldText
}

// returns a sorted textual representation of the slots in each sphere (except
// for the slot `except`), for easier comparison.
func spheresToText(spheres [][]*node, checks map[*node]*node, except *node) string {
	b := new(strings.Builder)
	for _, sphere := range spheres {
		sort.Slice(sphere, func(i, j int) bool {
			return sphere[i].name < sphere[j].name
		})
		for _, n := range sphere {
			if checks[n] != nil && n != except {
				b.WriteString(n.name + "\n")
			}
		}
	}
	return b.String()
}

type summary struct {
	items    map[string]string
	dungeons map[string]string
	portals  map[string]string
	seasons  map[string]string
	hints    map[string]string
}

func newSummary() *summary {
	return &summary{
		items:    make(map[string]string),
		dungeons: make(map[string]string),
		portals:  make(map[string]string),
		seasons:  make(map[string]string),
		hints:    make(map[string]string),
	}
}

func orderedValues(m map[string]string) []string {
	a, i := make([]string, len(m)), 0
	for _, v := range m {
		a[i] = v
		i++
	}
	sort.Strings(a)
	return a
}

var conditionRegexp = regexp.MustCompile(`(.+?) +<- (.+)`)

// loads conditions from a log file.
func parseSummary(path string, game int) (*summary, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	sum := newSummary()
	section := sum.items
	for _, line := range strings.Split(string(b), "\n") {
		line = strings.Replace(line, "\r", "", 1)
		if strings.HasPrefix(line, "--") {
			switch line {
			case "-- items --", "-- progression items --",
				"-- small keys and boss keys --", "-- other items --":
				section = sum.items
			case "-- dungeon entrances --":
				section = sum.dungeons
			case "-- subrosia portals --":
				section = sum.portals
			case "-- default seasons --":
				section = sum.seasons
			case "-- hints --":
				section = sum.hints
			default:
				return nil, fmt.Errorf("unknown section: %q", line)
			}
		} else {
			submatches := conditionRegexp.FindStringSubmatch(line)
			if submatches != nil {
				if submatches[1] == "null" {
					var nullKey string
					for i := 0; true; i++ {
						nullKey = fmt.Sprintf("null %d", i)
						if section[nullKey] == "" {
							break
						}
					}
					section[nullKey] = ungetNiceName(submatches[2], game)
				} else {
					section[ungetNiceName(submatches[1], game)] =
						ungetNiceName(submatches[2], game)
				}
			}
		}
	}

	return sum, nil
}

// write a "spoiler log" to a file.
func writeSummary(path string, checksum []byte, ropts randomizerOptions,
	rom *romState, ri *routeInfo, checks map[*node]*node, spheres [][]*node,
	extra []*node, owlHints map[string]string, fast bool) {
	summary, summaryDone := getSummaryChannel(path)

	// header
	summary <- fmt.Sprintf("seed: %08x", ri.seed)
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	summary <- fmt.Sprintf("difficulty: %s",
		ternary(ropts.hard, "hard", "normal"))
	summary <- ""

	// items
	if fast {
		sendSectionHeader(summary, "items")
		logSpheres(summary, checks, spheres, extra, rom.game, nil)
	} else {
		sendSectionHeader(summary, "progression items")
		nonKeyChecks := make(map[*node]*node)
		for slot, item := range checks {
			if !keyRegexp.MatchString(item.name) {
				nonKeyChecks[slot] = item
			}
		}
		prog, junk := filterJunk(ri.graph, nonKeyChecks, rom.treasures)
		logSpheres(summary, prog, spheres, extra, rom.game, nil)
		sendSectionHeader(summary, "small keys and boss keys")
		logSpheres(summary, checks, spheres, extra, rom.game, keyRegexp.MatchString)
		sendSectionHeader(summary, "other items")
		logSpheres(summary, junk, spheres, extra, rom.game, nil)
	}

	// warps
	if ropts.dungeons {
		sendSectionHeader(summary, "dungeon entrances")
		sendSorted(summary, func(c chan string) {
			for entrance, dungeon := range ri.entrances {
				c <- fmt.Sprintf("%s entrance <- %s",
					"D"+entrance[1:], "D"+dungeon[1:])
			}
			close(c)
		})
		summary <- ""
	}
	if ropts.portals {
		sendSectionHeader(summary, "subrosia portals")
		sendSorted(summary, func(c chan string) {
			for in, out := range ri.portals {
				c <- fmt.Sprintf("%-20s <- %s",
					getNiceName(in, rom.game), getNiceName(out, rom.game))
			}
			close(c)
		})
		summary <- ""
	}

	// default seasons (oos only)
	if rom.game == gameSeasons {
		sendSectionHeader(summary, "default seasons")
		sendSorted(summary, func(c chan string) {
			for area, id := range ri.seasons {
				c <- fmt.Sprintf("%-15s <- %s", area, seasonsById[id])
			}
			close(c)
		})
		summary <- ""
	}

	// owl hints
	sendSectionHeader(summary, "hints")
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

// sends a section delimiter to the channel.
func sendSectionHeader(c chan string, name string) {
	c <- ""
	c <- fmt.Sprintf("-- %s --", name)
	c <- ""
}
