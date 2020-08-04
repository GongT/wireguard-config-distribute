package tools

import (
	"log"
	"os"
	"strconv"
)

var suddenDeath bool = true
var isDying bool = false
var quitCode int = 0
var handlers []func(int) = make([]func(int), 0, 5)

func HasNoError() {
	quitCode = 0
}

func HasError(code int) {
	if quitCode == 0 {
		quitCode = code
	}
}

func Exit() {
	log.Println("program terminated with code " + strconv.FormatInt(int64(quitCode), 10) + ".")
	os.Exit(quitCode)
}

func ExitMain() {
	callCleanup()
	log.Println("program exited with code " + strconv.FormatInt(int64(quitCode), 10) + ".")
	os.Exit(quitCode)
}

func Die(format string, a ...interface{}) {
	Error(format, a...)
	if isDying {
		panic("die in die")
	}
	isDying = true
	if suddenDeath {
		if quitCode == 0 {
			quitCode = 1
		}
		log.Printf("program will terminate with code %d\n", quitCode)
		os.Exit(1)
	} else {
		callCleanup()
		Exit()
	}
}

func WaitExit(handler func(int)) {
	if suddenDeath {
		suddenDeath = false
	}

	handlers = append(handlers, handler)
}

func callCleanup() {
	log.Println("program will exit after cleanup")
	for _, handler := range handlers {
		handler(quitCode)
	}
}
