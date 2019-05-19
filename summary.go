package main

import (
	"fmt"
	"log"
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
			log.Fatal(err)
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
func filterJunk(g graph, checks map[*node]*node) (prog, junk map[*node]*node) {
	prog, junk = make(map[*node]*node), make(map[*node]*node)

	// start by assuming everything is progression
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
			if !itemIsRequired(g, slot, item) &&
				!itemChangesProgression(g, progCopy, spheres, slot, item) {
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
	g.clearMarks()
	item.removeParent(slot)
	mark := g["done"].getMark()
	item.addParent(slot)
	return mark != markTrue
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
