package client

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
	"github.com/gongt/wireguard-config-distribute/internal/types"
)

func (s *ClientStateHolder) handshake() bool {
	s.sharedStatus.lock()
	defer s.sharedStatus.unlock()

	s.isRunning = false
	data := &s.privateStatus

	tools.Error("  1: register...")
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
	tools.Error("  * register: ok. Session Id: %v\n    server offer ip address: %s/%d\n    interface private key: %s", result1.SessionId, result1.OfferIp, result1.Subnet, result1.PrivateKey)

	s.vpn.UpdateInterfaceInfo(result1.SessionId, result1.OfferIp, result1.PrivateKey, uint8(result1.Subnet))
	s.isRunning = true
	if s.machineId != result1.MachineId {
		if len(s.machineId) > 0 {
			tools.Error("Machine ID is different between server and local (using %s)", result1.MachineId)
		}
		s.machineId = result1.MachineId
	}
	s.sessionId = types.DeSerializeSidType(result1.SessionId)

	if result1.GetEnableObfuse() {
		shadow, err := s.nat.Start(uint16(data.InternalPortDefault))
		if err != nil {
			tools.Die("failed create nat: %v", err)
		}
		data.InternalPort = uint32(shadow)
	} else {
		s.nat.Stop()
		data.InternalPort = data.InternalPortDefault
	}
	s.vpn.SetWireguardListenPort(data.InternalPort)

	tools.Error("  2: update info...")
	_, err = s.server.UpdateClientInfo(&protocol.ClientInfoRequest{
		SessionId: s.sessionId.Serialize(),
		Services:  s.sharedStatus.services,
		Network: &protocol.PhysicalNetwork{
			ExternalIp:   s.ipDetect.GetLast(),
			ExternalPort: data.ExternalPort,
			InternalIp:   data.InternalIp,
			InternalPort: data.InternalPort,
			MTU:          uint32(data.SelfMtu),
		},
	})
	if err != nil {
		tools.Error("  * update info: failed: %s", err.Error())
		return false
	}

	tools.Error("  * update info: ok.")

	return true

}
