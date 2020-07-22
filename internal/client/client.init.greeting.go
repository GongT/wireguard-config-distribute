package client

import (
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *clientStateHolder) startProtocol() {
	for i := 0; i < 5; i++ {
		fmt.Printf("Handshake (try %d)\n", i)

		result, err := s.rpc.Greeting(s.context, &protocol.ClientInfoRequest{})

		if err == nil {
			fmt.Printf("  * handshake complete. server offer ip address: %s\n", s.address)

			s.address = result.OfferIp

			return
		} else {
			tools.Error("Failed handshake: %s", err.Error())
		}
	}
	tools.Die("Failed to greet server (after 5 retry).")
}
