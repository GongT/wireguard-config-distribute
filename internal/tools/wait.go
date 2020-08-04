package tools

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForCtrlC() chan bool {
	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)

	closed := false

	done := make(chan bool, 1)
	go func() {
		for range signal_channel {
			fmt.Println("")
			fmt.Println("Bye, bye.")

			if closed {
				Die("Terminate by double SIGINT.")
			} else {
				closed = true
				done <- true
				close(done)
			}
		}
	}()
	return done
}
