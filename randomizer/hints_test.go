package randomizer

import (
	"strings"
	"testing"
)

// make sure that every item and check has a corresponding hint name.
func TestHintCoverage(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		initRom(nil, game)
		hinter := newHinter(game)

		for name := range Treasures {
			// dungeon items aren't hinted and don't need names
			if getDungeonName(name) != "" {
				continue
			}

			if _, ok := hinter.items[name]; !ok {
				t.Errorf("%s missing name for item %q",
					gameName(game), name)
			}
		}

		for name := range ItemSlots {
			if _, ok := hinter.areas[name]; !ok {
				t.Errorf("%s missing name for area %q",
					gameName(game), name)
			}
		}
	}
}

// make sure that no hints refer to nothing.
func TestDanglingHints(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		initRom(nil, game)
		hinter := newHinter(game)

		for name := range hinter.items {
			if Treasures[name] == nil &&
				!strings.Contains(name, " ring") {
				t.Errorf("dangling %s item name: %q",
					gameName(game), name)
			}
		}

		for name := range hinter.areas {
			if ItemSlots[name] == nil {
				t.Errorf("dangling %s area name: %q",
					gameName(game), name)
			}
		}
	}
}
