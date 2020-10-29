// +build !windows

package config

type dummygroup struct {
	Hidden bool
}

var internalConfigGroup *dummygroup

func windowsAddProgramArguments() {
	internalConfigGroup = &dummygroup{}
}
func windowsCommitConfig() {
}
func IsSuRunning() bool {
	return false
}
