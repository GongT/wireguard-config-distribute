package client

import (
	"fmt"

	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (s *ClientStateHolder) handshake() bool {
	tools.Error("handshake:")
	s.statusData.lock()
	defer s.statusData.unlock()

	s.isRunning = false
	data := &s.configData

	result1, err := s.server.RegisterClient(&protocol.RegisterClientRequest{
		MachineId:    s.machineId,
		VpnGroup:     data.VpnGroupName,
		Title:        data.Title,
		Hostname:     data.Hostname,
		RequestVpnIp: s.vpn.GetRequestedAddress(),
		LocalGroup:   data.LocalNetworkName,
	})
	if err != nil {
		tools.Error("  * register: failed: %s", err.Error())
		return false
	}

	s.vpn.UpdateInterfaceInfo(result1.SessionId, result1.OfferIp, result1.PrivateKey, uint8(result1.Subnet))
	s.isRunning = true
	if s.machineId != result1.MachineId {
		if len(s.machineId) > 0 {
			tools.Error("Machine ID is different between server and local (using %s)", result1.MachineId)
		}
		s.machineId = result1.MachineId
	}
	s.sessionId = types.DeSerializeSidType(result1.SessionId)
	tools.Error("  * register: ok. Session Id: %v\n    server offer ip address: %s/%d\n    interface private key: %s", result1.SessionId, result1.OfferIp, result1.Subnet, result1.PrivateKey)

	if result1.GetEnableObfuse() {
		shadow, err := s.nat.Start(uint16(data.InternalPortDefault))
		if err != nil {
			panic(fmt.Errorf("failed create nat: %v", err))
		}
		data.InternalPort = uint32(shadow)
	} else {
		s.nat.Stop()
		data.InternalPort = data.InternalPortDefault
	}
	s.vpn.SetWireguardListenPort(data.InternalPort)

	_, err = s.server.UpdateClientInfo(&protocol.ClientInfoRequest{
		SessionId: s.sessionId.Serialize(),
		Services:  s.statusData.services,
		Network: &protocol.PhysicalNetwork{
			ExternalEnabled: data.ExternalEnabled,
			ExternalIp:      data.ExternalIp,
			ExternalPort:    data.ExternalPort,
			InternalIp:      data.InternalIp,
			InternalPort:    data.InternalPort,
			MTU:             uint32(data.SelfMtu),
		},
	})
	if err != nil {
		tools.Error("  * update info: failed: %s", err.Error())
		return false
	}

	tools.Error("  * update info: ok.")

	return true

}
