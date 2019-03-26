package hints

import (
	"testing"

	"github.com/jangler/oracles-randomizer/rom"
)

// make sure that every item and check has a corresponding hint name.
func TestHintCoverage(t *testing.T) {
	for name, _ := range rom.SeasonsTreasures {
		if _, ok := itemMap[name]; !ok {
			t.Errorf("missing name for seasons treasure \"%s\"", name)
		}
	}
	for name, _ := range rom.AgesTreasures {
		if _, ok := itemMap[name]; !ok {
			t.Errorf("missing name for ages treasure \"%s\"", name)
		}
	}
	for name, _ := range rom.SeasonsSlots {
		if _, ok := seasonsAreaMap[name]; !ok {
			t.Errorf("missing name for seasons slot \"%s\"", name)
		}
	}
	for name, _ := range rom.AgesSlots {
		if _, ok := agesAreaMap[name]; !ok {
			t.Errorf("missing name for ages slot \"%s\"", name)
		}
	}
}

// make sure that every name fits on a line, with room for punctuation.
func TestHintLength(t *testing.T) {
	for _, name := range itemMap {
		if len(name.name) > 15 {
			t.Errorf("hint name \"%s\" too long", name.name)
		}
	}
	for _, name := range seasonsAreaMap {
		if len(name) > 15 {
			t.Errorf("hint name \"%s\" too long", name)
		}
	}
	for _, name := range agesAreaMap {
		if len(name) > 15 {
			t.Errorf("hint name \"%s\" too long", name)
		}
	}
}
