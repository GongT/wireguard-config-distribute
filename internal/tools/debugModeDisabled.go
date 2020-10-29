// +build nodebug

package tools

const debugMode = false

func SetDebugMode(debug bool) {
	if debug {
		Error("this is none debug build, --debug has no effect!")
	}
}

func IsDevelopmennt() bool {
	return debugMode
}

func init() {
}
