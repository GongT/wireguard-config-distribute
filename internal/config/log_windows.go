package config

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/natefinch/npipe"
)

func SetLogPipe(path string) {
	os.Stdout = pipe(path+".stdout", os.Stdout)
	os.Stderr = pipe(path+".stderr", os.Stderr)
	log.SetOutput(os.Stderr)

	controlSocket, err := npipe.Dial(path + ".control")
	if err != nil {
		tools.Die("Failed connect pipes (%s): %s", path+".control", err.Error())
	}

	tools.WaitExit(func(code int) {
		os.Stdout = originalOut
		os.Stderr = originalErr
		log.SetOutput(originalErr)
		controlSocket.Write([]byte("exit:" + strconv.FormatInt(int64(code), 10) + "\n"))
		controlSocket.Close()
	})

	go func() {
		fscanner := bufio.NewScanner(controlSocket)
		for fscanner.Scan() {
			text := fscanner.Text()
			if text == "sigint" {
				tools.Error("receive sigint from control pipe")
				tools.Exit()
			}
		}
	}()
}

func pipe(from string, out *os.File) *os.File {
	conn, err := npipe.Dial(from)
	if err != nil {
		tools.Die("Failed connect pipes (%s): %s", from, err.Error())
	}
	rd, wr, err := os.Pipe()
	if err != nil {
		tools.Die("Failed create pipe: %s", err.Error())
	}

	ch := make(chan bool)
	go func() {
		io.Copy(conn, rd)
		ch <- true
	}()
	go tools.WaitExit(func(code int) {
		fmt.Fprintf(wr, "[child] exit %d", code)
		my_close(wr)
		<-ch
	})

	return wr
}
