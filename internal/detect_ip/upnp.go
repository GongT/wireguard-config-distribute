package detect_ip

import (
	"errors"
	"fmt"

	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func upnpGetPublicIp() (ret string, err error) {
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

	if !IsValidIPv4(ret) {
		err = errors.New("Invalid UPnP response.")
	}

	return
}
