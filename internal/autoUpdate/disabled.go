// +build noupdate

package autoUpdate

const updateEnabled = false

func checkUpdate() bool {
	return false
}

func applyUpdateAndRestart() {
}
