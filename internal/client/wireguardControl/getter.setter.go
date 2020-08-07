package wireguardControl

import (
	"strconv"

	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/wgexe"
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

	if err := wc.createConfigFile(); err != nil {
		tools.Error("failed creating config file: %s", err.Error())
	}

	if err := wc.nativeInterface.CreateOrUpdateInterface(wc); err != nil {
		tools.Error("failed update interface: %s", err.Error())
	}

	if wc.dryRun {
		return
	}

	if err := wgexe.GetWireguardCli().SmallChange(wc.interfaceName, wc.configFile); err != nil {
		tools.Error("failed update peers: %s", err.Error())
	}
}

func (wc *WireguardControl) GetNetwork() string {
	if wc.subnet > 0 {
		return wc.givenAddress + strconv.FormatUint(uint64(wc.subnet), 10)
	} else {
		return wc.givenAddress + "/32"
	}
}

func (wc *WireguardControl) GetMtu() int {
	return int(wc.interfaceMTU)
}

func (wc *WireguardControl) GetRequestedAddress() string {
	return wc.requestedAddress
}

func (wc *WireguardControl) UpdateInterfaceInfo(address string, privateKey string, subnet uint8) {
	wc.givenAddress = address
	wc.privateKey = privateKey
	wc.subnet = subnet
}