package main

// see generate.go in the directory above. this needs to be in a separate
// directory so that it's `go run`-able by `go generate`.

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

const (
	mainTemplate = `package main

// Code generated - DO NOT EDIT.

import "github.com/{{username}}/oracles-randomizer/randomizer"

func main() {
	randomizer.Main()
}
`
	versionTemplate = `package randomizer

// Code generated - DO NOT EDIT.

const version = {{version}}
`
)

var (
	usernameRegexp = regexp.MustCompile(`github.com/(.+)/oracles-randomizer`)
	versionRegexp  = regexp.MustCompile(`/(.+)-\d+-g(.+)`)
)

func main() {
	generateMain()
	generateVersion()
}

func generateMain() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	matches := usernameRegexp.FindStringSubmatch(wd)
	if matches == nil {
		panic("error getting import path from working directory")
	}

	s := strings.ReplaceAll(mainTemplate, "{{username}}", matches[1])
	if err := ioutil.WriteFile("main.go", []byte(s), 0644); err != nil {
		panic(err)
	}
}

func generateVersion() {
	// try matching an exact tag first
	describeCmd := exec.Command("git", "describe")
	output, err := describeCmd.Output()
	if err != nil {
		panic(err)
	}
	version := fmt.Sprintf(`"%s"`, strings.TrimSpace(string(output)))

	// not an exact tag; use long format
	if strings.Contains(string(output), "-") {
		describeCmd = exec.Command("git", "describe", "--all", "--long")
		if output, err = describeCmd.Output(); err != nil {
			panic(err)
		}
		matches := versionRegexp.FindStringSubmatch(string(output))
		if matches == nil {
			panic("error getting version string from git")
		}
		version = fmt.Sprintf(`"%s-%s"`, matches[1], matches[2])
	}

	s := strings.ReplaceAll(versionTemplate, "{{version}}", version)
	err = ioutil.WriteFile("randomizer/version.go", []byte(s), 0644)
	if err != nil {
		panic(err)
	}
}
