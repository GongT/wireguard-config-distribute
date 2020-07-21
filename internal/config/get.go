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

func IsDevelopmennt() bool {
	if os.Getenv("VSCODE_IPC_HOOK_CLI") != "" {
		return true
	}
	if os.Getenv(config.CONFIG_FORCE_DEVELOPMENT) != "" {
		return true
	}
	return false
}
