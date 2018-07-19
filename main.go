package main

import (
	"container/list"
	"log"
	"os"
)

func main() {
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
