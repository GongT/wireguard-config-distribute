package config

import (
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
)

type DebugOption struct {
	DebugMode bool
}

func ParseProgramArguments(opts interface{}) {
	_, err := flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash).Parse()

	if err != nil {
		serr, ok := err.(*flags.Error)
		if ok {
			switch serr.Type {
			case flags.ErrHelp:
				os.Exit(0)
			}
		}
		tools.Die("Failed parse arguments.\n\t%s", err.Error())
	}
}
