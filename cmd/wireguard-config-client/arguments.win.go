//go:generate go-generate-struct-interface
// +build windows

package main

type clientProgramOptions struct {
	clientProgramOptionsBase
	notMoveArguments

	InstallService   bool `long:"install" description:"install windows service"`
	UnInstallService bool `long:"uninstall" description:"uninstall windows service"`
}
