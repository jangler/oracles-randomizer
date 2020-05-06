package randomizer

import (
	"testing"
)

func TestMutableOverlap(t *testing.T) {
	for _, game := range []int{gameSeasons, gameAges} {
		hitBytes := make(map[int]*string)
		for k, v := range newRomState(nil, game, 0, nil).getAllMutables() {
			k := k
			switch v := v.(type) {
			case *mutableRange:
				offset := v.addr.fullOffset()
				for i := offset; i < offset+len(v.new); i++ {
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

func TestProcessText(t *testing.T) {
	testExpect(t, processText("A\\xff # hello\nB"), []byte{'A', 0xff, 'B'})
}
