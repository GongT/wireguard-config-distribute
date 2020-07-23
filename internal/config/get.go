package config

import (
	"os"
	"strconv"
)

func GetConfig(name string, defaultVal string) string {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultVal
	} else {
		return value
	}
}

func GetConfigNumber(name string, defaultVal int64) int64 {
	value := os.Getenv(name)
	if len(value) == 0 {
		return defaultVal
	} else {
		num, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic(err)
		}
		return num
	}
}

var debugMode bool = false

func setDebugMode(v bool) {
	debugMode = v
}
func IsDevelopmennt() bool {
	return debugMode
}

func getEnvDevelopmennt() bool {
	if os.Getenv("VSCODE_IPC_HOOK_CLI") != "" {
		return true
	}
	return false
}
