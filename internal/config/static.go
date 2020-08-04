package config

import (
	"bytes"
	"encoding/base64"
	"os"

	"github.com/jessevdk/go-flags"
)

func ResetInternalOption() {
	InternalOption.IsElevated = false
	InternalOption.StandardOutputPath = ""
}

func StringifyOptions() string {
	internalConfigGroup.Hidden = false

	ini := flags.NewIniParser(parser)
	buff := bytes.Buffer{}
	ini.Write(&buff, flags.IniIncludeDefaults)

	internalConfigGroup.Hidden = true

	encode := base64.StdEncoding.EncodeToString(buff.Bytes())
	return "data:" + encode
}

func dumpCurrentIni() {
	internalConfigGroup.Hidden = false
	ini := flags.NewIniParser(parser)
	ini.Write(os.Stderr, flags.IniIncludeDefaults)
	internalConfigGroup.Hidden = true
}

func IsSuRunning() bool {
	return InternalOption.IsElevated
}
