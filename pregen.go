package main

//go:generate go build
//go:generate ./oos-randomizer -devcmd pregen prenode/generated.go
//go:generate go fmt github.com/jangler/oos-randomizer/prenode
//go:generate go build

// this file contains logic for automaticaly generating graph prenodes based on
// special syntax in the keys:
//
// - if a key ends with a number, e.g. "scent tree 1", that key is added as a
//   new Or prenode named "scent tree" in the graph, with the original key as
//   a parent.
//
// the generated prenodes are written to a global map `generatedPrenodes` in a
// separate go source file in the working directory.

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/jangler/oos-randomizer/prenode"
)

const genPrenodeTemplate = `package prenode

var generatedPrenodes = map[string]*Prenode{
%s}
`

func generatePrenodes(w io.Writer) error {
	// get list of non-generated prenodes
	nonGeneratedPrenodes := prenode.GetNonGenerated()

	// get list of generated prenodes
	resultPrenodes := makeNumberPrenodes(nonGeneratedPrenodes)

	// consistently order map keys to minimize diffs
	orderedKeys := make(sort.StringSlice, len(resultPrenodes))
	i := 0
	for key := range resultPrenodes {
		orderedKeys[i] = key
		i++
	}
	orderedKeys.Sort()

	// build contents of generated map string
	builder := new(strings.Builder)
	for _, key := range orderedKeys {
		builder.WriteString(strings.Replace(fmt.Sprintf("\t\"%s\": %#v,\n",
			key, resultPrenodes[key]), "prenode.", "", 1))
	}

	// write out result
	_, err := fmt.Fprintf(w, genPrenodeTemplate, builder.String())
	return err
}

func makeNumberPrenodes(
	maps ...map[string]*prenode.Prenode) map[string]*prenode.Prenode {
	numberPrenodes := make(map[string]*prenode.Prenode)
	numberRegexp := regexp.MustCompile(`(^.+) \d+$`)

	for _, prenodes := range maps {
		for key := range prenodes {
			if strings.ContainsRune(key, ';') {
				continue
			}

			matches := numberRegexp.FindAllStringSubmatch(key, 1)
			if matches != nil {
				realKey := matches[0][1]
				if pt, ok := numberPrenodes[realKey]; ok {
					// sort for consistent order and minimal difs
					parents := sort.StringSlice(append(pt.Parents, key))
					parents.Sort()
					numberPrenodes[realKey] = prenode.Or(parents...)
				} else {
					numberPrenodes[realKey] = prenode.Or(key)
				}
			}
		}
	}

	return numberPrenodes
}
