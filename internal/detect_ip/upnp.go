// +build !android

package detect_ip

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/jackpal/gateway"
	natpmp "github.com/jackpal/go-nat-pmp"
)

func upnpGetPublicIp() (net.IP, error) {
	gatewayIP, err := gateway.DiscoverGateway()
	if err != nil {
		return nil, err
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
		ip := net.ParseIP(ret)
		if !tools.IsIPv4(ip) {
			return nil, errors.New("Invalid UPnP response.")
		}
		return ip, nil
	case e := <-ech:
		return nil, e
	case <-time.After(2 * time.Second):
		return nil, errors.New("UPnP timed out")
	}
}
