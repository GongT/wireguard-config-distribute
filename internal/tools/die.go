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

func HasFatalError() {
	quitCode = 233
}

func HasError(code int) {
	if quitCode == 0 {
		quitCode = code
	}
}

func Exit() {
	log.Println("program terminated with code " + strconv.FormatInt(int64(quitCode), 10) + ".")
	if quitCode == 0 {
		quitCode = 1
	}
	os.Exit(quitCode)
}

func ExitMain() {
	log.Println("exit() called")
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
		log.Println("program died.")
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
		Debug("no cleanup handler")
		return
	}
	Debug("program will exit after cleanup - %d handlers registered", len(handlers))
	for index, handler := range handlers {
		Debug("--------------- quit handler <%d>:", uint64(index))
		handler(quitCode)
	}
	handlers = make(map[uint]func(int))
	Debug("--------------- all quit handler has been called")
}
