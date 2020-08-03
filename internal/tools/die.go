package tools

import (
	"fmt"
	"log"
	"os"
)

var suddenDeath bool = true
var dieChan chan string
var isDying bool = false

func Die(format string, a ...interface{}) {
	Error(format, a...)
	if isDying {
		panic("die in die")
	}
	isDying = true
	if suddenDeath {
		log.Println("program will exit with code 1")
		os.Exit(1)
	} else {
		log.Println("program will exit after cleanup")
		dieChan <- fmt.Sprintf(format, a...)
		close(dieChan)
		log.Println("  * waitting")
		for {
		}
	}
}

func WaitDie() <-chan string {
	if !suddenDeath {
		panic("fatal error handler listen twice")
	}
	suddenDeath = false
	dieChan = make(chan string)
	return dieChan
}
