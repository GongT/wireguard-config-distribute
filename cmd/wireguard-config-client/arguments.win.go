//go:generate go-generate-struct-interface
// +build windows

package main

type clientProgramOptions struct {
	clientProgramOptionsBase
	notMoveArguments
}
