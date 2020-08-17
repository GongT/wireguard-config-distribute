// +build !android

package upnp

import (
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func TryAddPortMapping(port int, pubPort int) (uint16, error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return 0, err
	}

	client := natpmp.NewClient(gatewayIP)

	ret, err := client.AddPortMapping("udp", port, pubPort, int(60*24*time.Hour))

	if err != nil {
		return 0, err
	}

	tools.Debug("port forward timeout: %d", ret.PortMappingLifetimeInSeconds)

	return ret.MappedExternalPort, nil
}
