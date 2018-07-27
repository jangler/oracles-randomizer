package main

//go:generate go build
//go:generate ./oos-randomizer -devcmd pointgen generated.go
//go:generate go fmt
//go:generate go build

// this file contains logic for automaticaly generating graph "points" based on
// special syntax in the keys:
//
// - if a key ends with a number, e.g. "scent tree 1", that key is added as a
//   new Or point named "scent tree" in the graph, with the original key as a
//   parent.
//
// the generated points are written to a global map `generatedPoints` in a
// separate go source file in the working directory.

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"
)

const genPointTemplate = `package main

var generatedPoints = map[string]Point{
%s}
`

func generatePoints(w io.Writer) error {
	// get list of non-generated points
	nonGeneratedPoints := getNonGeneratedPoints()

	// get list of generated points
	resultPoints := makeNumberPoints(nonGeneratedPoints)

	// consistently order map keys to minimize diffs
	orderedKeys := make(sort.StringSlice, len(resultPoints))
	i := 0
	for key := range resultPoints {
		orderedKeys[i] = key
		i++
	}
	orderedKeys.Sort()

	// build contents of generated map string
	builder := new(strings.Builder)
	for _, key := range orderedKeys {
		builder.WriteString(strings.Replace(fmt.Sprintf("\t\"%s\": %#v,\n",
			key, resultPoints[key]), "main.", "", 1))
	}

	// write out result
	_, err := fmt.Fprintf(w, genPointTemplate, builder.String())
	return err
}

func makeNumberPoints(pointMaps ...map[string]Point) map[string]Point {
	numberPoints := make(map[string]Point)
	numberRegexp := regexp.MustCompile(`(^.+) \d+$`)

	for _, points := range pointMaps {
		for key := range points {
			if strings.ContainsRune(key, ';') {
				continue
			}

			matches := numberRegexp.FindAllStringSubmatch(key, 1)
			if matches != nil {
				realKey := matches[0][1]
				if pt, ok := numberPoints[realKey]; ok {
					// sort for consistent order and minimal difs
					parents := sort.StringSlice(append(pt.Parents, key))
					parents.Sort()
					numberPoints[realKey] = Or(parents...)
				} else {
					numberPoints[realKey] = Or(key)
				}
			}
		}
	}

	return numberPoints
}
