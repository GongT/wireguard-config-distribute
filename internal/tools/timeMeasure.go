package tools

import (
	"fmt"
	"time"
)

func TimeMeasure(title string) func() {
	if debugMode {
		start := time.Now()
		return func() {
			duration := time.Since(start)
			fmt.Printf("[timing] %v: %v\n", title, duration)
		}
	} else {
		return func() {}
	}
}
