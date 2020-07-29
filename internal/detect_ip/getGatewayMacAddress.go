package detect_ip

import (
	"net"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jackpal/gateway"
)

var gatwayIp net.IP

func GetGatewayIP() (net.IP, error) {
	if gatwayIp == nil {
		gwIp, err := gateway.DiscoverGateway()
		if err != nil {
			return nil, err
		}

		tools.Debug("Gateway IP address is: %s", gwIp.String())
		gatwayIp = gwIp
	}

	return gatwayIp, nil
}

func GetGatewayMac() (string, error) {
	gwIp, err := GetGatewayIP()
	if err != nil {
		return "", err
	}

	return findGatewayAddr(gwIp)
}
