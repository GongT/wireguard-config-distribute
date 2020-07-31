package tools

import (
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func GetSystemHostName(name *string) bool {
	if len(*name) > 0 {
		Error("hostname=%s", *name)
		return true
	}
	*name = os.Getenv("HOSTNAME")
	if len(*name) > 0 {
		Error("using HOSTNAME")
		return true
	}
	*name = os.Getenv("COMPUTERNAME")
	if len(*name) > 0 {
		Error("using COMPUTERNAME")
		return true
	}
	if runtime.GOOS == "linux" {
		if data, err := ioutil.ReadFile("/etc/hostname"); err == nil {
			Error("using /etc/hostname")
			*name = strings.TrimSpace(string(data))
		} else {
			Error("failed reading /etc/hostname (%s)", err.Error())
		}
	}

	return len(*name) > 0
}
