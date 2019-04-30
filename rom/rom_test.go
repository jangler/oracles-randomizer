package rom

import "testing"

func init() {
	Init(GameSeasons) // XXX have to change this manually to test each game
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
