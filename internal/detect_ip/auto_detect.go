package detect_ip

import (
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func Detect(ipv4 *string, ipv6 *string, allowHttp bool, allowUPnP bool) {
	var err error
	if len(*ipv4) == 0 && allowUPnP {
		fmt.Println("  * try to get ip from UPnP")
		*ipv4, err = upnpGetPublicIp()
		if err == nil {
			tools.Error("      -> %s", *ipv4)
		} else {
			tools.Error("      -> %s", err.Error())
		}
	}
	if len(*ipv4) == 0 && allowHttp {
		fmt.Println("  * try to get ipv4 from http")
		*ipv4, err = httpGetPublicIp4()
		if err == nil {
			tools.Error("      -> %s", *ipv4)
		} else {
			tools.Error("      -> %s", err.Error())
		}
	}

	if len(*ipv6) == 0 && allowHttp {
		fmt.Println("  * try to get ipv6 from http")
		*ipv6, err = httpGetPublicIp6()
		if err == nil {
			tools.Error("      -> %s", *ipv6)
		} else {
			tools.Error("      -> %s", err.Error())
		}
	}
}
