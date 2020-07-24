package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *clientStateHolder) UploadInformation() bool {
	data := s.configData
	s.statusData.lock()

	result, err := s.server.Greeting(&protocol.ClientInfoRequest{
		GroupName:    data.GroupName,
		Title:        data.Title,
		Hostname:     data.Hostname,
		Services:     s.statusData.services,
		RequestVpnIp: s.vpn.requestedAddress,
		Network: &protocol.PhysicalNetwork{
			NetworkId:       data.NetworkId,
			ExternalEnabled: data.ExternalEnabled,
			ExternalIp:      data.ExternalIp,
			ExternalPort:    data.ExternalPort,
			InternalIp:      data.InternalIp,
			InternalPort:    data.InternalPort,
		},
	})

	if err == nil {
		tools.Error("  * handshake complete. server offer ip address: %s\n", result.OfferIp)

		s.vpn.givenAddress = result.OfferIp

		s.isRunning = true

		s.statusData.unlock()
		return true
	} else {
		tools.Error("  * failed handshake: %s", err.Error())

		s.isRunning = false

		s.statusData.unlock()
		return false
	}
}
