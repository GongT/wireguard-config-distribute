package tools

import (
	"fmt"
	"os"
)

func Error(format string, a ...interface{}) {
	fmt.Fprint(os.Stderr, fmt.Sprintf(format, a...)+"\n")
}
