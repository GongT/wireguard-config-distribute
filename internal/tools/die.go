package tools

import (
	"fmt"
	"os"
)

func Die(format string, a ...interface{}) {
	fmt.Fprint(os.Stderr, fmt.Sprintf(format, a...)+"\n")
	os.Exit(1)
}
