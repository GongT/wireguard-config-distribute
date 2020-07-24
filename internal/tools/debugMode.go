package tools

import "os"

var debugMode bool

func SetDebugMode(v bool) {
	debugMode = v
}

func IsDevelopmennt() bool {
	return debugMode
}

func init() {
	if os.Getenv("VSCODE_IPC_HOOK_CLI") != "" {
		debugMode = true
	} else {
		debugMode = false
	}
}
