package main

import (
	"log"
	"os"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	log.Println("test run!")

	os.Setenv("GO_TEST_ENV", "true")
	os.Args = []string{os.Args[0]}

	go main()
	log.Println("test wait")
	time.Sleep(30 * time.Second)
	log.Println("test done")
}
