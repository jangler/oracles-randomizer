package rom

import "testing"

func TestGraphicsPresent(t *testing.T) {
	for name, _ := range Treasures {
		if itemGfx[name] == 0 {
			t.Errorf("no graphics for %s", name)
		}
	}
}
