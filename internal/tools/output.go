package tools

import (
	"fmt"
	"os"
)

func Error(format string, a ...interface{}) {
	fmt.Fprint(os.Stderr, fmt.Sprintf(format, a...)+"\n")
}

func Debug(format string, a ...interface{}) {
	if debugMode {
		fmt.Fprint(os.Stderr, fmt.Sprintf(format, a...)+"\n")
	}
}
