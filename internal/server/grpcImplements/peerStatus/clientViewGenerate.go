package peerStatus

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (peersList *vpnPeersMap) generateAllView(viewer *PeerData) *protocol.Peers {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(*peersList)-1)

	for cid, peer := range *peersList {
		hosts[peer.VpnIp] = peer.HostsLine

		if viewer.sessionId == cid {
			continue
		}

		list = append(list, peersList.generateOneView(viewer, peer))
	}

	return &protocol.Peers{
		List:  list,
		Hosts: hosts,
	}
}

func (peersList *vpnPeersMap) generateOneView(viewer *PeerData, peer *PeerData) *protocol.Peers_Peer {
	var keepAlive uint32 = 0
	port := peer.ExternalPort
	ip := peer.ExternalIp
	if viewer.NetworkId == peer.NetworkId && len(viewer.NetworkId) > 0 {
		// same local network
		port = peer.InternalPort
		ip = []string{peer.InternalIp}
	} else if len(viewer.ExternalIp) == 0 && len(peer.ExternalIp) != 0 {
		keepAlive = 25
	}

	p := protocol.Peers_Peer{
		SessionId: peer.sessionId.Serialize(),
		Title:     peer.Title,
		Hostname:  peer.Hostname,
		Peer: &protocol.Peers_ConnectionTarget{
			PublicKey: peer.PublicKey,
			Address:   ip,
			Port:      port,
			VpnIp:     peer.VpnIp,
			KeepAlive: keepAlive,
			MTU:       peer.MTU,
		},
	}

	return &p
}
