// +build !android

package detect_ip

import (
	"errors"
	"fmt"
	"net"

	"github.com/davecgh/go-spew/spew"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func ListAllLocalNetworkIp() (ret []net.IP) {
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
				if v.IP.IsLoopback() || v.IP.Equal(net.IPv4bcast) {
					continue
				}

				ret = append(ret, v.IP)
			case *net.IPAddr:
				fmt.Printf("  * %s::%s -> %s\n", iface.Name, addr.String(), spew.Sdump(addr))
			}
		}
	}

	return
}

func findRouteFromIp(target net.IP) (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.New("Failed get local network interface: " + err.Error())
	}

	for _, iface := range ifaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if ipnet.IP.IsLinkLocalUnicast() || ipnet.IP.IsLoopback() {
					continue
				}

				if ipnet.Contains(target) {
					return ipnet.IP, nil
				}
			}
		}
	}

	return nil, errors.New("No route to " + target.String())
}

func GetDefaultNetworkIp() (net.IP, error) {
	gatewayIP, err := GetGatewayIP()
	if err != nil {
		return nil, errors.New("Failed get default gateway: " + err.Error())
	}

	result, err := findRouteFromIp(gatewayIP)
	if err != nil {
		return nil, err
	}
	return result, nil
}
