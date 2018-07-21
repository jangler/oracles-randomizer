package main

import (
	"container/list"
	"fmt"
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

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	loadedRom, err := rom.Load(f)
	if err != nil {
		log.Fatal(err)
	}

	errs := rom.Verify(loadedRom)
	if len(errs) > 0 {
		fmt.Printf("%d errors:\n", len(errs))
		for _, err := range rom.Verify(loadedRom) {
			fmt.Println(err)
		}
	} else {
		fmt.Println("everything OK")
	}
}
