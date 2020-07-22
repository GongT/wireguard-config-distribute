package tools

import (
	"os"
)

func Die(format string, a ...interface{}) {
	Error(format, a...)
	os.Exit(1)
}
