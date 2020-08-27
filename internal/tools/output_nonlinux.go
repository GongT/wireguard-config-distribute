// +build !linux

package tools

import (
	"log"
)

func Error(format string, a ...interface{}) {
	log.Printf(format+"\n", a...)
}

func Debug(format string, a ...interface{}) {
	if debugMode {
		log.Printf(format+"\n", a...)
	}
}
