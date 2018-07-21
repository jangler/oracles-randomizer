package main

import (
	"bytes"
	"container/list"
	"log"
	"os"

	"github.com/jangler/oos-randomizer/rom"
)

func main() {
	if false {
		r := NewRoute()
		path := list.New()
		target := r.Map[os.Args[1]]
		mark := target.GetMark(path)
		if mark == MarkTrue {
			for path.Len() > 0 {
				step := path.Remove(path.Front()).(string)
				log.Print(step)
			}
		} else {
			log.Print("path not found")
		}
	}

	var loadedRom *bytes.Buffer
	{
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		loadedRom, err = rom.Load(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	errs := rom.Verify(loadedRom)
	if len(errs) > 0 {
		log.Printf("%d verification errors:\n", len(errs))
		for _, err := range rom.Verify(loadedRom) {
			log.Print(err)
		}
		return
	} else {
		log.Print("everything OK")
	}

	if len(os.Args) > 2 {
		// randomize rom
		b := loadedRom.Bytes()
		rom.Mutate(b)

		// write to file
		f, err := os.Create(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		if _, err := f.Write(b); err != nil {
			log.Fatal(err)
		}
		log.Printf("wrote new ROM to %s", os.Args[2])
	}
}
