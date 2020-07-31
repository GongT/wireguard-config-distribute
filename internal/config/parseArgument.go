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

type DebugOption interface {
	GetDebugMode() bool
}

type Sanitizable interface {
	Sanitize() error
}

type wrappedParser struct {
	Parser      *flags.Parser
	isSuRunning bool
	options     interface{}
	parseError  error
}

var parserCached *wrappedParser

func InitProgramArguments(opts interface{}) (*wrappedParser, error) {
	if parserCached != nil {
		return parserCached, nil
	}

	r := wrappedParser{
		Parser:  flags.NewParser(opts, flags.HelpFlag|flags.PassDoubleDash),
		options: opts,
	}

	r.parseError = r.ParseCommandline()

	parserCached = &r
	return parserCached, r.parseError
}

func (wp *wrappedParser) Err() error {
	return wp.parseError
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

func (wp *wrappedParser) ParseCommandline() error {
	if len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "data:") {
		ini := flags.NewIniParser(wp.Parser)
		bs, _ := base64.StdEncoding.DecodeString(os.Args[1][5:])
		iniData := bytes.NewReader(bs)
		err := ini.Parse(iniData)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed parse ini arguments: %s", err.Error()))
		}
		wp.isSuRunning = true

		wp.updateDebug()
	} else {
		_, err := wp.Parser.Parse()

		if err != nil {
			wp.Parser.WriteHelp(os.Stderr)
			serr, ok := err.(*flags.Error)
			if ok {
				switch serr.Type {
				case flags.ErrHelp:
					os.Exit(0)
				}
			}
			return errors.New(fmt.Sprintf("Failed parse arguments.\n\t%s", err.Error()))
		}

		wp.updateDebug()

		if opts, ok := wp.options.(Sanitizable); ok {
			err := opts.Sanitize()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (wp *wrappedParser) updateDebug() {
	if UpdateDebug != nil {
		if opts, ok := wp.options.(DebugOption); ok {
			UpdateDebug(opts.GetDebugMode())
		}
	}
}

var UpdateDebug func(debug bool)

func StringifyOptions() string {
	ini := flags.NewIniParser(parserCached.Parser)
	buff := bytes.Buffer{}
	ini.Write(&buff, flags.IniIncludeDefaults)
	encode := base64.StdEncoding.EncodeToString(buff.Bytes())
	return "data:" + encode
}

func IsSuRunning() bool {
	return parserCached.isSuRunning
}
