package peerStatus

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
)

func (peers *PeerStatus) createAllView(viewer lPeerData) *protocol.Peers {
	hosts := make(map[string]string)
	list := make([]*protocol.Peers_Peer, 0, len(peers.list)-1)

	for cid, peer := range peers.list {
		if viewer.sessionId == cid {
			continue
		}

		list = append(list, peers.createOneView(viewer, peer))

		for _, host := range peer.Hosts {
			hosts[host+"."+peer.Hostname] = peer.VpnIp
		}
	}

	return &protocol.Peers{
		List:  list,
		Hosts: hosts,
	}
}

func (peers *PeerStatus) createOneView(viewer lPeerData, peer lPeerData) *protocol.Peers_Peer {
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
