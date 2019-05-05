package rom

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

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

// checks for multiple instances of the same asm change, so that duplicates can
// be merged to free code that can be `/include`d.
func TestIdenticalAsm(t *testing.T) {
	filenames := []string{
		"/asm/common.yaml",
		"/asm/seasons.yaml",
		"/asm/ages.yaml",
	}
	ads := make([]*asmData, len(filenames))

	for i, filename := range filenames {
		ads[i] = new(asmData)
		ads[i].filename = filename
		if err := yaml.Unmarshal(
			FSMustByte(false, filename), ads[i]); err != nil {
			panic(err)
		}
	}

	changes := make(map[string]*changeId)
	for _, ad := range ads {
		for _, bankItems := range ad.Appends {
			for _, item := range bankItems {
				for name, content := range item {
					checkAddChangeId(t, changes, content, &changeId{
						filename: ad.filename,
						name:     name,
					})
				}
			}
		}
		for name, content := range ad.FreeCode {
			checkAddChangeId(t, changes, content, &changeId{
				filename: ad.filename,
				name:     name,
			})
		}
	}
}

type changeId struct {
	filename, name string
}

func checkAddChangeId(t *testing.T, m map[string]*changeId, k string,
	v *changeId) {
	if k == "" || strings.HasPrefix(k, "/include") {
		return
	}

	dup := m[k]
	if dup != nil {
		t.Errorf("change %s:%s has same content as %s:%s",
			v.filename, v.name, dup.filename, dup.name)
	} else {
		m[k] = v
	}
}
