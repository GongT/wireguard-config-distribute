package tools

import (
	"log"
	"os"
	"runtime/debug"
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
	if !suddenDeath {
		callCleanup()
	}

	debug.PrintStack()

	if quitCode == 0 {
		quitCode = 1
	}
	Exit()
}

func WaitExit(handler func(int)) func() {
	if suddenDeath {
		suddenDeath = false
	}
	handlersIndex++

	currIndex := handlersIndex

	handlers[currIndex] = handler

	return func() {
		delete(handlers, currIndex)
		if len(handlers) == 0 {
			suddenDeath = true
		}
	}
}

func callCleanup() {
	if len(handlers) == 0 {
		log.Println("no cleanup handler")
		return
	}
	log.Println("program will exit after cleanup -", len(handlers))
	for index, handler := range handlers {
		println("--------------- quit handler " + strconv.FormatUint(uint64(index), 10) + ":")
		handler(quitCode)
	}
	handlers = make(map[uint]func(int))
	log.Println("--------------- all quit handler has been called")
}
