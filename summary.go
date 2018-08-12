package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

const version = "1.3.3"

// returns a channel that will write strings to a text file with CRLF line
// endings. the function will send on the int channel when finished printing.
func getSummaryChannel() (chan string, chan int) {
	c, done := make(chan string), make(chan int)

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
		done <- 1
	}()

	// header
	c <- fmt.Sprintf("oos-randomizer %s", version)
	c <- fmt.Sprintf("generated %s", time.Now().Format(time.RFC3339))

	return c, done
}
