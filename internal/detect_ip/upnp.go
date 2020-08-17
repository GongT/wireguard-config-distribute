// +build !android

package detect_ip

import (
	"errors"
	"fmt"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func upnpGetPublicIp() (string, error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return "", err
	}

	client := natpmp.NewClient(gatewayIP)

	rch := make(chan *natpmp.GetExternalAddressResult, 1)
	ech := make(chan error, 1)
	defer close(rch)
	defer close(ech)

	go func() {
		response, err := client.GetExternalAddress()
		if err == nil {
			rch <- response
		} else {
			ech <- err
		}
	}()

	select {
	case response := <-rch:
		ret := fmt.Sprintf("%x.%x.%x.%x", response.ExternalIPAddress[0], response.ExternalIPAddress[1], response.ExternalIPAddress[2], response.ExternalIPAddress[3])
		if !tools.IsValidIPv4(ret) {
			return "", errors.New("Invalid UPnP response.")
		}
		return ret, nil
	case e := <-ech:
		return "", e
	case <-time.After(2 * time.Second):
		return "", errors.New("UPnP timed out")
	}
}
