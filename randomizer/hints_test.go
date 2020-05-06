package randomizer

import (
	"strings"
	"testing"
)

// make sure that every item and check has a corresponding hint name.
func TestHintCoverage(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		rom := newRomState(nil, game, 0, nil)
		hinter := newHinter(game)

		for name := range rom.treasures {
			// dungeon items aren't hinted and don't need names
			if getDungeonName(name) != "" {
				continue
			}

			if _, ok := hinter.items[name]; !ok {
				t.Errorf("%s missing name for item %q",
					gameNames[game], name)
			}
		}

		for name := range rom.itemSlots {
			if _, ok := hinter.areas[name]; !ok {
				t.Errorf("%s missing name for area %q",
					gameNames[game], name)
			}
		}
	}
}

// make sure that no hints refer to nothing.
func TestDanglingHints(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		rom := newRomState(nil, game, 0, nil)
		hinter := newHinter(game)

		for name := range hinter.items {
			if rom.treasures[name] == nil &&
				!strings.Contains(name, " ring") {
				t.Errorf("dangling %s item name: %q",
					gameNames[game], name)
			}
		}

		for name := range hinter.areas {
			if rom.itemSlots[name] == nil {
				t.Errorf("dangling %s area name: %q",
					gameNames[game], name)
			}
		}
	}
}
