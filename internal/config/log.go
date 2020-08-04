package config

import (
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

var switchedLog = false

func Cleanup() {
	if switchedLog {
		close(os.Stdout)
		close(os.Stderr)
	}
}

func close(logger *os.File) {
	if err := logger.Sync(); err != nil {
		tools.Error("file.Sync() fail: %s", err.Error())
	} else if err := logger.Close(); err != nil {
		tools.Error("file.Close() fail: %s", err.Error())
	}
}

func SetLogOutput(path string) {
	if len(path) == 0 {
		path = filepath.Join(filepath.Dir(os.Args[0]), "output.log")
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0644)

	if err != nil {
		tools.Die("Failed open log file: %s", err.Error())
	}

	os.Stdout = f
	os.Stderr = f
	switchedLog = true
}
