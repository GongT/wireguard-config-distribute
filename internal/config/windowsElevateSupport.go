// +build windows

package config

import (
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
)

type internalOptions struct {
	StandardOutputPath string `long:"pipe-output"`
	IsElevated         bool   `long:"is-elevate"`
}

var InternalOption = &internalOptions{}
var internalConfigGroup *flags.Group

func windowsAddProgramArguments() {
	if g, err := parser.AddGroup("Elevate Options", "", InternalOption); err != nil {
		tools.Die("internal error: %v", err)
	} else {
		g.Hidden = true
		internalConfigGroup = g
	}
}

func windowsCommitConfig() {
	if InternalOption.IsElevated && len(InternalOption.StandardOutputPath) > 0 {
		tools.Error("log will dup to pipes (%s).", InternalOption.StandardOutputPath)
		SetLogPipe(InternalOption.StandardOutputPath)
		tools.Error("[child] log start.")
	}
}

func IsSuRunning() bool {
	return InternalOption.IsElevated
}
