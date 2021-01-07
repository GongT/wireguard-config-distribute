package config

import (
	"bytes"
	"encoding/base64"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
)

type Sanitizable interface {
	Sanitize() error
}

var lastError error
var parser *flags.Parser

var ApplicationOption interface{}
var CommonOption = &commonOptions{}

func InitProgramArguments(opts interface{}) error {
	if parser != nil {
		tools.Die("duplicate call to InitProgramArguments()")
	}

	ApplicationOption = opts

	parser = flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash)

	if _, err := parser.AddGroup("Common Options", "", CommonOption); err != nil {
		tools.Die("internal invalid arguments: %v", err)
	}

	windowsAddProgramArguments()

	lastError := parseCommandline()

	return lastError
}

func Err() error {
	return lastError
}

func CommandActive(name string) bool {
	for curr := parser.Active; curr != nil; curr = curr.Active {
		if curr.Name == name || (curr.Aliases != nil && tools.ArrayContains(curr.Aliases, name)) {
			return true
		}
	}

	return false
}

func DieUsage() {
	parser.WriteHelp(os.Stderr)
	tools.Die("")
}

func parseCommandline() error {
	if len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "data:") {
		tools.Debug("use arguments data:")
		ini := flags.NewIniParser(parser)
		bs, _ := base64.StdEncoding.DecodeString(os.Args[1][5:])
		iniData := bytes.NewReader(bs)
		err := ini.Parse(iniData)
		if err != nil {
			return errors.Wrap(err, "failed parse ini arguments")
		}

		commitConfig()
	} else {
		tools.Debug("parse arguments: %s", strings.Join(os.Args, " "))
		extra, err := parser.Parse()

		if CommonOption.ShowVersion {
			tools.ShowVersion(os.Stdout)
			os.Exit(0)
		}

		if len(extra) > 0 && err == nil {
			err = errors.New("unknown positional argument: " + extra[0])
		}

		if err != nil {
			serr, ok := err.(*flags.Error)
			if ok {
				switch serr.Type {
				case flags.ErrHelp:
					parser.WriteHelp(os.Stdout)
					os.Exit(0)
				}
			}
			return err
		}

		commitConfig()

		if opts, ok := ApplicationOption.(Sanitizable); ok {
			err := opts.Sanitize()
			if err != nil {
				return err
			}
		}
	}

	if CommonOption.DebugMode {
		dumpCurrentIni()
	}

	return nil
}

func commitConfig() {
	tools.SetDebugMode(CommonOption.DebugMode)

	windowsCommitConfig()
}
