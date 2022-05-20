//////////////////////////////////////////////////////////////////////
//
// Given is a mock process which runs indefinitely and blocks the
// program. Right now the only way to stop the program is to send a
// SIGINT (Ctrl-C). Killing a process like that is not graceful, so we
// want to try to gracefully stop the process first.
//
// Change the program to do the following:
//   1. On SIGINT try to gracefully stop the process using
//          `proc.Stop()`
//   2. If SIGINT is called again, just kill the program (last resort)
//

package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	// Create a process
	proc := MockProcess{}
	// create a channel for relaying the SIGINT signal
	c := make(chan os.Signal, 1)
	// make this channel notified when a SIGINT has received
	signal.Notify(c, os.Interrupt)

	// waiting to receive from the channel (<-c)is a blocking operation
	// so create a go routine for it
	go func(c <-chan os.Signal) {
		// blocks untill signal is received
		<-c
		fmt.Println("SIGINT captured, C^c again to exit forcefully")
		go proc.Stop()

		// if process has not stopped kill the process
		<-c
		fmt.Println("Exiting")
		os.Exit(0)
		return
	}(c)

	// Run the process (blocking)
	proc.Run()
}
