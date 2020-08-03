package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (s *ClientStateHolder) uploadInformation() bool {
	tools.Error(" ~  uploadInformation()")
	s.statusData.lock()
	defer s.statusData.unlock()

	// if s.isRunning {
	// 	return s.isRunning
	// }

	data := s.configData

	result, err := s.server.Greeting(&protocol.ClientInfoRequest{
		MachineId:    s.machineId,
		GroupName:    data.GroupName,
		Title:        data.Title,
		Hostname:     data.Hostname,
		Services:     s.statusData.services,
		RequestVpnIp: s.vpn.GetRequestedAddress(),
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
		tools.Error("  * complete.\n        server offer ip address: %s/%d\n        interface private key: %s", result.OfferIp, result.Subnet, result.PrivateKey)

		s.vpn.UpdateInterface(result.OfferIp, result.PrivateKey, uint16(result.Subnet))
		s.isRunning = true
		if s.machineId != result.MachineId {
			if len(s.machineId) > 0 {
				tools.Error("Machine ID is different between server and local (using %s)", result.MachineId)
			}
			s.machineId = result.MachineId
		}
		s.sessionId = types.DeSerialize(result.SessionId)

		// TODO create interface

		return true
	} else {
		tools.Error("  * failed: %s", err.Error())

		s.isRunning = false

		return false
	}
}
