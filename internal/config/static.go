package config

import (
	"bytes"
	"encoding/base64"
	"os"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jessevdk/go-flags"
)

func StringifyOptions(withInternal bool) string {
	internalConfigGroup.Hidden = !withInternal

	ini := flags.NewIniParser(parser)
	buff := bytes.Buffer{}
	ini.Write(&buff, flags.IniIncludeDefaults)

	internalConfigGroup.Hidden = true

	str := strings.TrimSpace(buff.String())

	tools.Error("================================================================")
	tools.Error(str)
	tools.Error("================================================================")

	encode := base64.StdEncoding.EncodeToString([]byte(str))
	return "data:" + encode
}

func dumpCurrentIni() {
	internalConfigGroup.Hidden = false
	ini := flags.NewIniParser(parser)
	ini.Write(os.Stderr, flags.IniIncludeDefaults)
	internalConfigGroup.Hidden = true
}
