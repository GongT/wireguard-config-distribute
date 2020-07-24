package upnp

import (
	"fmt"
	"time"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func TryAddPortMapping(port int) (err error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return
	}

	client := natpmp.NewClient(gatewayIP)

	ret, err := client.AddPortMapping("udp", port, port, int(60*24*time.Hour))

	if err != nil {
		return
	}
	fmt.Println(ret)

	return
}

func GetPublicIp() (ret string, err error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		return
	}

	ret = fmt.Sprintf("%x.%x.%x.%x", response.ExternalIPAddress[0], response.ExternalIPAddress[1], response.ExternalIPAddress[2], response.ExternalIPAddress[3])

	return
}
