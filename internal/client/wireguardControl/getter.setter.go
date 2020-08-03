package wireguardControl

import (
	"github.com/gongt/wireguard-config-distribute/internal/protocol"
	"github.com/gongt/wireguard-config-distribute/internal/tools"
)

type peerData struct {
	comment      string
	publicKey    string
	presharedKey string
	ip           string
	port         uint16
	keepAlive    uint
	privateIp    string
}

func (wc *WireguardControl) UpdatePeers(list []*protocol.Peers_Peer) {
	defer wc.mu.Lock("update peers")()

	tools.Error("Updating peers:")
	wc.peers = wc.peers[0:0]
	for _, peer := range list {
		tools.Error("  * <%d> %s -> %v", peer.GetSessionId(), peer.GetHostname(), peer.GetPeer().GetAddress())
		selectedIp := selectIp(peer.GetPeer().GetAddress())
		tools.Error("      -> %s:%d", selectedIp, peer.GetPeer().GetPort())

		kl := uint(peer.GetPeer().GetKeepAlive())

		wc.peers = append(wc.peers, peerData{
			comment:      peer.GetTitle(),
			publicKey:    peer.GetPeer().GetPublicKey(),
			presharedKey: "",
			ip:           selectedIp,
			port:         uint16(peer.GetPeer().GetPort()),
			keepAlive:    kl,
			privateIp:    peer.GetPeer().GetVpnIp(),
		})
	}

	err := wc.creatConfigFile()
	if err != nil {
		tools.Error("Failed creating config file: %s", err.Error())
		return
	}

	wc.updateInterface()
}

func (wc *WireguardControl) GetRequestedAddress() string {
	return wc.requestedAddress
}

func (wc *WireguardControl) UpdateInterface(address string, privateKey string, subnet uint16) {
	wc.givenAddress = address
	wc.privateKey = privateKey
	wc.subnet = subnet
}
