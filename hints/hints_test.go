package hints

import (
	"strings"
	"testing"

	"github.com/jangler/oracles-randomizer/rom"
)

// TODO rings aren't covered by these tests, since they're separate (and
// private) in the rom package.

// make sure that every item and check has a corresponding hint name.
func TestHintCoverage(t *testing.T) {
	for name, _ := range rom.SeasonsTreasures {
		if _, ok := itemMap[name]; !ok {
			t.Errorf("missing name for seasons treasure %q", name)
		}
	}
	for name, _ := range rom.AgesTreasures {
		if _, ok := itemMap[name]; !ok {
			t.Errorf("missing name for ages treasure %q", name)
		}
	}
	for name, _ := range rom.SeasonsSlots {
		if _, ok := seasonsAreaMap[name]; !ok {
			t.Errorf("missing name for seasons slot %q", name)
		}
	}
	for name, _ := range rom.AgesSlots {
		if _, ok := agesAreaMap[name]; !ok {
			t.Errorf("missing name for ages slot %q", name)
		}
	}
}

// make sure that no hints refer to nothing.
func TestDanglingHints(t *testing.T) {
	for name := range itemMap {
		if rom.SeasonsTreasures[name] == nil &&
			rom.AgesTreasures[name] == nil &&
			!strings.Contains(name, " ring") {
			t.Errorf("dangling item name: %q", name)
		}
	}
	for name := range seasonsAreaMap {
		if rom.SeasonsSlots[name] == nil {
			t.Errorf("dangling area name: %q", name)
		}
	}
	for name := range agesAreaMap {
		if rom.AgesSlots[name] == nil {
			t.Errorf("dangling area name: %q", name)
		}
	}
}
