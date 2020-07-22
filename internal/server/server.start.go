package server

import "github.com/gongt/wireguard-config-distribute/internal/protocol"

func (s *serverImplement) Start(*protocol.ClientInfoRequest, protocol.WireguardApi_StartServer) error {
	return nil
}
