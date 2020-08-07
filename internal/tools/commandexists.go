package tools

import "os/exec"

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func EnsureCommandExists(cmd string, extrainfo string) {
	if !CommandExists(cmd) {
		if len(extrainfo) == 0 {
			Die("Can not find `%s' executable from PATH", cmd)
		} else {
			Die("Can not find `%s' executable from PATH\n    %s", cmd, extrainfo)
		}
	}
}
