package tools

import (
	"net"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/constants"
)

func NormalizeServerString(pstr *string) {
	str := *pstr
	var rs net.Addr
	var err error

	if str[0:1] == "/" {
		rs, err = net.ResolveUnixAddr("unix", str)
		if err != nil {
			Die("Failed resolve unix socket '%s': %s", str, err.Error())
		}
		*pstr = rs.String()
		return
	}
	if ip := net.ParseIP(str); ip != nil {
		if ip.To4() == nil {
			*pstr = "[" + ip.String() + "]:" + constants.DEFAULT_PORT
		} else {
			*pstr = ip.String() + ":" + constants.DEFAULT_PORT
		}
	} else if !strings.Contains(str, ":") {
		*pstr += ":" + constants.DEFAULT_PORT
	}
}
