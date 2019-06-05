package randomizer

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadYaml(t *testing.T) {
	// make sure all yaml files are well-formed.
	dirnames := []string{"asm", "hints", "logic", "romdata"}
	for _, dirname := range dirnames {
		// get list of files in directory
		dir, err := FS(false).Open("/" + dirname + "/")
		if err != nil {
			t.Fatal(err)
		}
		files, err := dir.Readdir(-1)
		if err != nil {
			t.Fatal(err)
		}

		for _, file := range files {
			// ignore non-yaml files like readmes
			if !strings.HasSuffix(file.Name(), ".yaml") {
				continue
			}

			path := "/" + dirname + "/" + file.Name()

			// either a slice or a map should work
			m := make(map[interface{}]interface{})
			mapErr := yaml.Unmarshal(FSMustByte(false, path), &m)
			a := make([]interface{}, 0)
			sliceErr := yaml.Unmarshal(FSMustByte(false, path), &a)

			if mapErr != nil && sliceErr != nil {
				t.Errorf("failed to unmarshal %s into map or slice", path)
				t.Error(mapErr)
				t.Error(sliceErr)
			}
		}
	}
}
