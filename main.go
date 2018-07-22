package main

import (
	"container/list"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jangler/oos-randomizer/graph"
	"github.com/jangler/oos-randomizer/rom"
)

func usage() {
	fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	fmt.Fprintf(flag.CommandLine.Output(), `
Valid operations are checkGraph, verifyData, findPath, and randomize.
`[1:])
}

func main() {
	flag.Usage = usage
	flagOp := flag.String("op", "chedkGraph", "operation")
	flag.Parse()

	switch *flagOp {
	case "chedkGraph":
		if flag.NArg() != 0 {
			log.Fatalf("findPath takes 0 arguments; got %d", flag.NArg())
		}

		// validate
		if errs := checkGraph(); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}
	case "verifyData":
		if flag.NArg() != 1 {
			log.Fatalf("verify takes 1 argument; got %d", flag.NArg())
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// verify program data vs rom data
		if errs := rom.Verify(romData); err != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		} else {
			log.Print("everything OK")
		}
	case "findPath":
		if flag.NArg() != 1 {
			log.Fatalf("findPath takes 2 arguments; got %d", flag.NArg())
		}

		// try to find a valid path to the target node
		g, _ := initRoute() // ignore errors; they're diagnostic only
		target, ok := g.Map[flag.Arg(0)]
		if !ok {
			log.Fatal("target node not found")
		}

		path := list.New()
		mark := target.GetMark(path)
		if mark == graph.MarkTrue {
			for path.Len() > 0 {
				step := path.Remove(path.Front()).(string)
				log.Print(step)
			}
		} else {
			log.Print("path not found")
		}
	case "randomize":
		if flag.NArg() != 2 {
			log.Fatalf("randomize takes 2 arguments; got %d", flag.NArg())
		}

		// load rom
		romData, err := loadROM(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}

		// randomize
		if errs := randomize(romData, flag.Arg(1)); errs != nil {
			for _, err := range errs {
				log.Print(err)
			}
			os.Exit(1)
		}
	default:
		log.Fatalf("no such operation: %s", *flagOp)
	}
}

// make sure the base route graph is ok (before randomizing anything)
func checkGraph() []error {
	// TODO initRoute() does this automatically and i'm not sure it should
	_, errs := initRoute()
	return errs
}

// can be used for loading pretty much anything, really, as long as you want it
// as a slice of bytes.
func loadROM(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return rom.Load(f)
}

// messes up rom data and writes it to a file.
//
// this also calls verify.
func randomize(romData []byte, outFilename string) []error {
	if errs := rom.Verify(romData); errs != nil {
		return errs
	}

	// XXX old code, but could be used as reference for new code
	/*
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
	*/

	return []error{fmt.Errorf("NYI")}
}
