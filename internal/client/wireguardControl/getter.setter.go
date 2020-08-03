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
	wc.mu.Lock()
	defer wc.mu.Unlock()

	tools.Error("Updating peers:")
	wc.peers = wc.peers[0:0]
	for _, peer := range list {
		selectedIp := selectIp(peer.GetPeer().GetAddress())
		if len(selectedIp) == 0 {
			tools.Error("  * DROP <%s>, failed ping any of %v", peer.GetTitle(), peer.GetPeer().GetAddress())
			continue
		}

		tools.Error("  * <%d> %s -> %s:%d", peer.GetSessionId(), peer.GetHostname(), selectedIp, peer.GetPeer().GetPort())
		wc.peers = append(wc.peers, peerData{
			comment:      peer.GetTitle(),
			publicKey:    peer.GetPeer().GetPublicKey(),
			presharedKey: "",
			ip:           selectedIp,
			port:         uint16(peer.GetPeer().GetPort()),
			keepAlive:    uint(peer.GetPeer().GetKeepAlive()),
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
