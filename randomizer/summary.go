package randomizer

import (
	"fmt"
	"os"
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
func filterJunk(g graph, checks map[*node]*node, treasures map[string]*treasure,
	resetFunc func()) (prog, junk map[*node]*node) {
	prog, junk = make(map[*node]*node), make(map[*node]*node)

	// get all required items. if multiple instances of the same class exist
	// and any is skippable but some are required, the first instances are
	// considered required and the rest are considered unrequired.
	spheres, _ := getSpheres(g, checks, resetFunc)
	for _, class := range getAllItemClasses(checks) {
		// skip known inert items
		if class != "rupees" && itemIsInert(treasures, class) {
			continue
		}

		// start by removing all instances
		removed := make(map[*node]*node)
		for slot, item := range checks {
			if item.name == class ||
				(class == "rupees" && strings.HasPrefix(item.name, "rupees")) {
				removed[slot] = item
				item.removeParent(slot)
			}
		}

		// add instances back one at a time in sphere order
		resetFunc()
		g.reset()
		g["start"].explore()
		for !g["done"].reached {
		outerLoop:
			for _, sphere := range spheres {
				for _, node := range sphere {
					if item := removed[node]; item != nil {
						delete(removed, node)
						prog[node] = item
						item.addParent(node)
						break outerLoop
					}
				}
			}
			resetFunc()
			g.reset()
			g["start"].explore()
		}

		// add all other instances back
		for slot, item := range removed {
			item.addParent(slot)
		}
	}

	// remove denominations of rupees that were added but are actually too
	// small to matter.
	junkRupees := make(map[*node]*node)
	for slot, item := range checks {
		if strings.HasPrefix(item.name, "rupees") && prog[slot] == nil {
			item.removeParent(slot)
			junkRupees[slot] = item
		}
	}
	trivialRupees := make([]*node, 0, 10)
	for slot, item := range prog {
		if strings.HasPrefix(item.name, "rupees") {
			item.removeParent(slot)
			resetFunc()
			g.reset()
			g["start"].explore()
			if g["done"].reached {
				trivialRupees = append(trivialRupees, slot)
			}
			item.addParent(slot)
		}
	}
	for _, slot := range trivialRupees {
		delete(prog, slot)
	}
	for slot, item := range junkRupees {
		item.addParent(slot)
	}

	// the remainder is junk.
	for slot, item := range checks {
		if prog[slot] == nil {
			junk[slot] = item
		}
	}

	return
}

// return an ordered slice of names of different item classes. all rupees are
// considered a single class.
func getAllItemClasses(checks map[*node]*node) []string {
	allClasses := make(map[string]bool)
	for _, item := range checks {
		if strings.HasPrefix(item.name, "rupees") {
			allClasses["rupees"] = true
		} else {
			allClasses[item.name] = true
		}
	}
	return orderedKeys(allClasses)
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

// write a "spoiler log" to a file.
func writeSummary(path string, checksum []byte, ropts randomizerOptions,
	rom *romState, ri *routeInfo, checks map[*node]*node, spheres [][]*node,
	extra []*node, g graph, resetFunc func(), treasures map[string]*treasure,
	owlHints map[string]string) {
	summary, summaryDone := getSummaryChannel(path)

	// header
	summary <- fmt.Sprintf("seed: %08x", ri.seed)
	summary <- fmt.Sprintf("sha-1 sum: %x", checksum)
	summary <- fmt.Sprintf("difficulty: %s",
		ternary(ropts.hard, "hard", "normal"))

	// items
	nonKeyChecks := make(map[*node]*node)
	for slot, item := range checks {
		if !keyRegexp.MatchString(item.name) {
			nonKeyChecks[slot] = item
		}
	}
	prog, junk := filterJunk(g, nonKeyChecks, treasures, resetFunc)
	sendSectionHeader(summary, "progression items")
	logSpheres(summary, prog, spheres, extra, rom.game, nil)
	sendSectionHeader(summary, "small keys and boss keys")
	logSpheres(summary, checks, spheres, extra, rom.game, keyRegexp.MatchString)
	sendSectionHeader(summary, "other items")
	logSpheres(summary, junk, spheres, extra, rom.game, nil)

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
	}

	// owl hints
	if owlHints != nil {
		sendSectionHeader(summary, "hints")
		sendSorted(summary, func(c chan string) {
			for owlName, hint := range owlHints {
				oneLineHint := strings.ReplaceAll(hint, "\n", " ")
				oneLineHint = strings.ReplaceAll(oneLineHint, "  ", " ")
				c <- fmt.Sprintf("%-20s <- \"%s\"", owlName, oneLineHint)
			}
			close(c)
		})
	}

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
	c <- ""
	c <- fmt.Sprintf("-- %s --", name)
	c <- ""
}
