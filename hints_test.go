package main

import (
	"strings"
	"testing"

	"github.com/jangler/oracles-randomizer/rom"
)

// make sure that every item and check has a corresponding hint name.
func TestHintCoverage(t *testing.T) {
	for _, game := range []int{rom.GameSeasons, rom.GameAges} {
		rom.Init(nil, game)
		hinter := newHinter(game)

		for name := range rom.Treasures {
			if _, ok := hinter.items[name]; !ok {
				t.Errorf("%s missing name for item %q",
					gameName(game), name)
			}
		}

		for name := range rom.ItemSlots {
			if _, ok := hinter.areas[name]; !ok {
				t.Errorf("%s missing name for area %q",
					gameName(game), name)
			}
		}
	}
}

// make sure that no hints refer to nothing.
func TestDanglingHints(t *testing.T) {
	for _, game := range []int{rom.GameSeasons, rom.GameAges} {
		rom.Init(nil, game)
		hinter := newHinter(game)

		for name := range hinter.items {
			if rom.Treasures[name] == nil &&
				!strings.Contains(name, " ring") {
				t.Errorf("dangling %s item name: %q",
					gameName(game), name)
			}
		}

		for name := range hinter.areas {
			if rom.ItemSlots[name] == nil {
				t.Errorf("dangling %s area name: %q",
					gameName(game), name)
			}
		}
	}
}
