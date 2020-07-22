package tools

import (
	"fmt"
	"os"
	"os/signal"
)

func WaitForCtrlC() {
	signal_channel := make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	<-signal_channel
	fmt.Println("")
}
