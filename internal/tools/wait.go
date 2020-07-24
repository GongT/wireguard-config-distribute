package tools

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForCtrlC() chan bool {
	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)

	done := make(chan bool, 1)
	go func() {
		<-signal_channel
		fmt.Println("")
		fmt.Println("Bye, bye.")
		done <- true
	}()
	return done
}
