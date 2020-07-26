package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *clientStateHolder) UploadInformation() bool {
	s.statusData.lock()
	defer s.statusData.unlock()

	if s.isRunning {
		return s.isRunning
	}

	data := s.configData

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
		tools.Error("  * complete. server offer ip address: %s\n", result.OfferIp)

		s.vpn.givenAddress = result.OfferIp
		s.SessionId = result.SessionId
		s.isRunning = true

		return true
	} else {
		tools.Error("  * failed: %s", err.Error())

		s.SessionId = 0
		s.isRunning = false

		return false
	}
}
