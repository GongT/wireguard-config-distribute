package tools

import "os/exec"

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func EnsureCommandExists(cmd string) {
	if !CommandExists(cmd) {
		Die("Can not find `%s' executable from PATH", cmd)
	}
}
