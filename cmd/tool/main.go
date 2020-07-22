package main

import (
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/upnp"
)

func main() {
	err := upnp.Discover()
	if err != nil {
		tools.Error("no NAT-PMP: %s", err.Error())
	}

	// client := client.NewClient()

	// client.ConnectGrpc()

	// client.Quit()
}
