package rom

import (
	"sort"
	"strings"

	"gopkg.in/yaml.v2"
)

// returns a map of owl names to text indexes for the given game.
func getOwlIds(game int) map[string]byte {
	owls := make(map[string]map[string]byte)
	if err := yaml.Unmarshal(
		FSMustByte(false, "/romdata/owls.yaml"), owls); err != nil {
		panic(err)
	}
	return owls[gameNames[game]]
}

// updates the owl statue text data based on the given hints. does not mutate
// anything.
func SetOwlData(owlHints map[string]string, game int) {
	table := codeMutables["owlTextOffsets"]
	text := codeMutables["owlText"]
	builder := new(strings.Builder)
	addr := text.Addrs[0].offset
	owlTextIds := getOwlIds(game)

	for _, owlName := range GetOwlNames(game) {
		hint := owlHints[owlName]
		textID := owlTextIds[owlName]
		str := "\x0c\x00" + strings.ReplaceAll(hint, "\n", "\x01") + "\x00"
		table.New[textID*2] = addrString(addr)[0]
		table.New[textID*2+1] = addrString(addr)[1]
		addr += uint16(len(str))
		builder.WriteString(str)
	}

	text.New = []byte(builder.String())

	codeMutables["owlTextOffsets"] = table
	codeMutables["owlText"] = text
}

// returns a sorted array of owl statue names for the given game (matching
// those in the logic package).
func GetOwlNames(game int) []string {
	src := getOwlIds(game)

	a := make([]string, len(src))
	i := 0
	for k := range getOwlIds(game) {
		a[i] = k
		i++
	}

	sort.Strings(a)

	return a
}
