package detect_ip

import (
	"fmt"
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type options interface {
	GetIpHttpDisable() bool
	GetIpUpnpDisable() bool

	GetIpApi6() string
	GetIpApi4() string
}

func RunDetect(ipv4 *net.IP, options options) {
	var err error
	if len(*ipv4) == 0 && !options.GetIpUpnpDisable() {
		tools.Error("  * try to get ip from UPnP")
		*ipv4, err = upnpGetPublicIp()
		if err == nil {
			tools.Error("      -> %s", *ipv4)
		} else {
			tools.Error("      x> %s", err.Error())
		}
	}
	if len(*ipv4) == 0 && !options.GetIpHttpDisable() {
		fmt.Println("  * try to get ipv4 from http")
		*ipv4, err = httpGetPublicIp(options.GetIpApi4())
		if err == nil {
			tools.Error("      -> %s", *ipv4)
		} else {
			tools.Error("      x> %s", err.Error())
		}
	}
}
