package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const version = "1.1.0"

// returns a channel that will write strings to a text file with CRLF line
// endings.
func getSummaryChannel() chan string {
	c := make(chan string)

	go func() {
		logFile, err := os.Create("oos-randomizer_log_" +
			time.Now().Format("2006-01-02_150405") + ".txt")
		if err != nil {
			log.Fatal(err)
		}
		defer logFile.Close()

		for line := range c {
			fmt.Fprintf(logFile, "%s\r\n", line)
		}
	}()

	// header
	c <- fmt.Sprintf("oos-randomizer %s", version)
	c <- fmt.Sprintf("generated %s", time.Now().Format(time.RFC3339))

	return c
}
