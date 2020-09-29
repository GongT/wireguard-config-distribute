package config

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
)

type Sanitizable interface {
	Sanitize() error
}

var lastError error
var parser *flags.Parser

var ApplicationOption interface{}
var InternalOption = &internalOptions{}
var CommonOption = &commonOptions{}
var internalConfigGroup *flags.Group

func InitProgramArguments(opts interface{}) error {
	if parser != nil {
		panic(fmt.Errorf("duplicate call to InitProgramArguments()"))
	}

	ApplicationOption = opts

	parser = flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash)

	if _, err := parser.AddGroup("Common Options", "", CommonOption); err != nil {
		panic(err)
	}

	if g, err := parser.AddGroup("Elevate Options", "", InternalOption); err != nil {
		panic(err)
	} else {
		g.Hidden = true
		internalConfigGroup = g
	}

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
		ini := flags.NewIniParser(parser)
		bs, _ := base64.StdEncoding.DecodeString(os.Args[1][5:])
		iniData := bytes.NewReader(bs)
		err := ini.Parse(iniData)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed parse ini arguments: %s", err.Error()))
		}

		commitConfig()
	} else {
		_, err := parser.Parse()

		if CommonOption.ShowVersion {
			tools.ShowVersion()
			os.Exit(0)
		}

		if err != nil {
			parser.WriteHelp(os.Stderr)
			serr, ok := err.(*flags.Error)
			if ok {
				switch serr.Type {
				case flags.ErrHelp:
					os.Exit(0)
				}
			}
			return errors.New(fmt.Sprintf("Failed parse arguments.\n\t%s", err.Error()))
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

	if InternalOption.IsElevated && len(InternalOption.StandardOutputPath) > 0 {
		tools.Error("log will dup to pipes (%s).", InternalOption.StandardOutputPath)
		SetLogPipe(InternalOption.StandardOutputPath)
		tools.Error("[child] log start.")
	} else if len(CommonOption.LogFilePath) > 0 {
		tools.Error("log will dup to file (%s).", CommonOption.LogFilePath)
		SetLogOutput(CommonOption.LogFilePath)
		tools.Error("log start.")
	}
}
