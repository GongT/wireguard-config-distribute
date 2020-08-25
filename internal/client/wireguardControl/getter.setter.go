package wireguardControl

import (
	"strconv"
	"strings"

	"github.com/gongt/wireguard-config-distribute/internal/client/clientType"
	"github.com/gongt/wireguard-config-distribute/internal/client/wireguardControl/wgexe"
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
	mtu          uint16
}

func (wc *WireguardControl) UpdatePeers(list clientType.PeerDataList) {
	defer wc.mu.Lock("update peers")()

	tools.Error("Updating peers:")
	wc.peers = wc.peers[0:0]
	for _, client := range list {
		peer := client.GetPeer()
		selectedIp := client.GetSelectedAddress()
		remotePort := client.GetSelectedPort()
		tools.Error("  * <%d> %s", client.GetSessionId(), client.GetTitle())
		if len(selectedIp) == 0 {
			tools.Error("      endpoint: <no external ip>:%d, mtu: %v", remotePort, peer.GetMTU())
		} else {
			tools.Error("      endpoint: %s:%d, mtu: %v", selectedIp, remotePort, peer.GetMTU())
		}

		kl := uint(peer.GetKeepAlive())

		wc.peers = append(wc.peers, peerData{
			comment:      client.GetTitle(),
			publicKey:    peer.GetPublicKey(),
			presharedKey: "",
			ip:           selectedIp,
			port:         remotePort,
			keepAlive:    kl,
			privateIp:    peer.GetVpnIp(),
			mtu:          uint16(peer.GetMTU()),
		})
	}

	wc.lowestMtu = uint16(1420)
	for _, p := range wc.peers {
		if p.mtu > 0 && wc.lowestMtu > p.mtu {
			wc.lowestMtu = p.mtu
		}
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
	return wc.networkAddr
}

func (wc *WireguardControl) GetAddress() string {
	return wc.givenAddress + "/32"
}

func (wc *WireguardControl) GetMtu() uint16 {
	return wc.lowestMtu
}

func (wc *WireguardControl) GetRequestedAddress() string {
	return wc.requestedAddress
}

func (wc *WireguardControl) UpdateInterfaceInfo(address string, privateKey string, subnet uint8) {
	addrs := strings.Split(address, ".")
	var networkAddr string
	switch uint64(subnet) {
	case 24:
		networkAddr = strings.Join(addrs[:3], ".")
		networkAddr += ".0"
	case 16:
		networkAddr = strings.Join(addrs[:2], ".")
		networkAddr += ".0.0"
	case 8:
		networkAddr = addrs[1]
		networkAddr += ".0.0.0"
	default:
		tools.Die("server did not send subnet infomation")
	}
	networkAddr += "/" + strconv.FormatUint(uint64(subnet), 10)

	wc.givenAddress = address
	wc.networkAddr = networkAddr
	wc.privateKey = privateKey
	wc.subnet = subnet
}
