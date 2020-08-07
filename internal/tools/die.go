package tools

import (
	"log"
	"os"
	"strconv"
)

var suddenDeath bool = true
var isDying bool = false
var quitCode int = 0
var handlersIndex uint = 0
var handlers map[uint]func(int) = make(map[uint]func(int))

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

func WaitExit(handler func(int)) func() {
	if suddenDeath {
		suddenDeath = false
	}
	handlersIndex++

	currIndex := handlersIndex

	handlers[currIndex] = handler

	return func() {
		if len(handlers) == 0 {
			suddenDeath = true
		}
		delete(handlers, currIndex)
	}
}

func callCleanup() {
	if len(handlers) == 0 {
		log.Println("no cleanup handler")
		return
	}
	log.Println("program will exit after cleanup -", len(handlers))
	for _, handler := range handlers {
		println("---------------")
		handler(quitCode)
	}
	println("---------------")
	handlers = make(map[uint]func(int))
	log.Println("all cleanup has been called")
}
