package config

import (
	"os"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
)

type DebugOption struct {
	DebugMode bool
}

type wrappedParser struct {
	Parser *flags.Parser

	options interface{}
}

func InitProgramArguments(opts interface{}) *wrappedParser {
	r := wrappedParser{
		Parser:  flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash),
		options: opts,
	}

	r.ParseCommandline()

	return &r
}
func (wp *wrappedParser) Exists(name string) bool {
	for curr := wp.Parser.Active; curr != nil; curr = curr.Active {
		if curr.Name == name || (curr.Aliases != nil && tools.ArrayContains(curr.Aliases, name)) {
			return true
		}
	}

	return false
}

func (wp *wrappedParser) DieUsage() {
	wp.Parser.WriteHelp(os.Stderr)
	os.Exit(1)
}

func (wp *wrappedParser) ParseCommandline() []string {
	extras, err := wp.Parser.Parse()

	if err != nil {
		wp.Parser.WriteHelp(os.Stderr)
		serr, ok := err.(*flags.Error)
		if ok {
			switch serr.Type {
			case flags.ErrHelp:
				os.Exit(0)
			}
		}
		tools.Die("Failed parse arguments.\n\t%s", err.Error())
	}
	return extras
}
