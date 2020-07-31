package service

import (
	"os"
	"path/filepath"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func SetLogOutput(path string) *os.File {
	if len(path) == 0 {
		path = filepath.Join(filepath.Dir(os.Args[0]), "output.log")
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0644)

	if err != nil {
		tools.Die("Failed open log file: %s", err.Error())
	}

	os.Stdout = f
	os.Stderr = f
	/*
		multiWriter := io.MultiWriter(os.Stdout, f)
		rd, wr, err := os.Pipe()
		if err != nil {
			tools.Die("Failed open pipe: %s", err.Error())
		}
		os.Stdout = wr
		os.Stderr = wr

		go func() {
			scanner := bufio.NewScanner(rd)
			for scanner.Scan() {
				stdoutLine := scanner.Text()
				multiWriter.Write([]byte(stdoutLine + "\n"))
			}
		}()
	*/

	return f
}
