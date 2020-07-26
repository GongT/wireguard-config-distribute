package detect_ip

import (
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func DetectLocalNetwork() (ret []string) {
	ifaces, err := net.Interfaces()
	if err != nil {
		tools.Die("Failed get local network interface: %s", err.Error())
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if v.IP.IsLinkLocalUnicast() || v.IP.IsLoopback() {
					continue
				}

				ret = append(ret, v.IP.String())
			case *net.IPAddr:
				fmt.Printf("  * %s::%s -> %s\n", iface.Name, addr.String(), spew.Sdump(addr))
			}
		}
	}

	return
}
