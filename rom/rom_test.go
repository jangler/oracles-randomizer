package rom

import "testing"

func init() {
	Init(GameAges) // XXX have to change this manually to test each game
}

func TestGraphicsPresent(t *testing.T) {
	for name, _ := range Treasures {
		if itemGfx[name] == 0 {
			t.Errorf("no graphics for %s", name)
		}
	}
}

func TestMutableOverlap(t *testing.T) {
	hitBytes := make(map[int]*string)

	for k, v := range getAllMutables() {
		k := k
		switch v := v.(type) {
		case *MutableRange:
			for _, addr := range v.Addrs {
				offset := addr.fullOffset()
				for i := offset; i < offset+len(v.New); i++ {
					if hitBytes[i] != nil {
						t.Errorf("%s collides with %s at %d",
							k, *hitBytes[i], i)
					}
					hitBytes[i] = &k
				}
			}
		}
	}
}
