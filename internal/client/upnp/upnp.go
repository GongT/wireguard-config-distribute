package upnp

import (
	"fmt"
	"time"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func TryAddPortMapping(external int) (err error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return
	}

	client := natpmp.NewClient(gatewayIP)
	response, err := client.GetExternalAddress()
	if err != nil {
		return
	}
	fmt.Printf("External IP address: %v\n", response.ExternalIPAddress)

	ret, err := client.AddPortMapping("udp", external, external, int(60*24*time.Hour))

	if err != nil {
		return
	}
	fmt.Println(ret)

	return
}
