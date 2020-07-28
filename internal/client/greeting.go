package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

func (s *clientStateHolder) uploadInformation() bool {
	tools.Error(" ~  uploadInformation()")
	s.statusData.lock()
	defer s.statusData.unlock()

	// if s.isRunning {
	// 	return s.isRunning
	// }

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
		tools.Error("  * complete. server offer ip address: %s", result.OfferIp)

		s.vpn.givenAddress = result.OfferIp
		s.vpn.interfacePrivateKey = result.PrivateKey
		s.SessionId = result.SessionId
		s.isRunning = true

		// TODO create interface

		return true
	} else {
		tools.Error("  * failed: %s", err.Error())

		s.SessionId = 0
		s.isRunning = false

		return false
	}
}
