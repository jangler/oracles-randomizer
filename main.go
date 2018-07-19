package main

import (
	"container/list"
	"log"
)

func main() {
	r := NewRoute()
	path := list.New()
	target := r.Map["kill facade"]
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
