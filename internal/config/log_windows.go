package config

import (
	"io"
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/natefinch/npipe"
)

func SetLogPipe(path string) {
	os.Stdout = pipe(path+".stdout", os.Stdout)
	os.Stderr = pipe(path+".stderr", os.Stderr)
	switchedLog = true
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

	go io.Copy(conn, rd)

	return wr
}
